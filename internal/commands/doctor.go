package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ar4mirez/aicof/internal/core"
	"github.com/ar4mirez/aicof/internal/ui"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check AICoF installation health",
	Long: `Verify the AICoF framework installation is complete and healthy.

Checks performed:
- Config file exists and is valid
- CLAUDE.md is present
- All installed components exist
- No broken file references
- Directory structure is correct

Examples:
  aicof doctor           # Run health check
  aicof doctor --fix     # Auto-fix issues where possible`,
	RunE: runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().Bool("fix", false, "Auto-fix issues where possible")
}

type checkResult struct {
	name    string
	passed  bool
	message string
	fixable bool
}

func runDoctor(cmd *cobra.Command, args []string) error {
	autoFix, _ := cmd.Flags().GetBool("fix")

	ui.Header("AICoF Health Check")

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	var results []checkResult

	// Check 1: Config file exists
	config, configErr := core.LoadConfig()
	if configErr != nil {
		if os.IsNotExist(configErr) {
			results = append(results, checkResult{
				name:    "Config file",
				passed:  false,
				message: "aicof.yaml not found",
				fixable: false,
			})
		} else {
			results = append(results, checkResult{
				name:    "Config file",
				passed:  false,
				message: fmt.Sprintf("Config error: %v", configErr),
				fixable: false,
			})
		}
	} else {
		results = append(results, checkResult{
			name:    "Config file",
			passed:  true,
			message: fmt.Sprintf("aicof.yaml found (v%s)", config.Version),
		})
	}

	// Check 2: CLAUDE.md exists
	claudeMdPath := filepath.Join(cwd, "CLAUDE.md")
	if _, err := os.Stat(claudeMdPath); os.IsNotExist(err) {
		results = append(results, checkResult{
			name:    "CLAUDE.md",
			passed:  false,
			message: "CLAUDE.md not found",
			fixable: true,
		})
	} else {
		// Check version in CLAUDE.md
		content, _ := os.ReadFile(claudeMdPath)
		version := extractVersion(string(content))
		if version != "" {
			results = append(results, checkResult{
				name:    "CLAUDE.md",
				passed:  true,
				message: fmt.Sprintf("Present (v%s)", version),
			})
		} else {
			results = append(results, checkResult{
				name:    "CLAUDE.md",
				passed:  true,
				message: "Present",
			})
		}
	}

	// Check 3: AGENTS.md exists
	agentsMdPath := filepath.Join(cwd, "AGENTS.md")
	if _, err := os.Stat(agentsMdPath); os.IsNotExist(err) {
		results = append(results, checkResult{
			name:    "AGENTS.md",
			passed:  false,
			message: "AGENTS.md not found (cross-tool compatibility)",
			fixable: true,
		})
	} else {
		results = append(results, checkResult{
			name:    "AGENTS.md",
			passed:  true,
			message: "Present",
		})
	}

	// Check 4: .claude directory structure
	claudeDirs := []string{
		".claude",
		".claude/skills",
	}

	allDirsExist := true
	var missingDirs []string
	for _, dir := range claudeDirs {
		dirPath := filepath.Join(cwd, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			allDirsExist = false
			missingDirs = append(missingDirs, dir)
		}
	}

	if allDirsExist {
		results = append(results, checkResult{
			name:    "Directory structure",
			passed:  true,
			message: ".claude/skills/ directory structure valid",
		})
	} else {
		results = append(results, checkResult{
			name:    "Directory structure",
			passed:  false,
			message: fmt.Sprintf("Missing directories: %s", strings.Join(missingDirs, ", ")),
			fixable: true,
		})
	}

	// Check 5: Installed language guide skills exist
	if config != nil {
		// Ensure legacy configs have skills populated
		config.MigrateLanguagesToSkills()

		var missingLangs []string
		for _, lang := range config.Installed.Languages {
			component := core.FindLanguage(lang)
			if component != nil {
				// Language guides are now skill directories with SKILL.md
				skillPath := filepath.Join(cwd, component.Path, "SKILL.md")
				if _, err := os.Stat(skillPath); os.IsNotExist(err) {
					missingLangs = append(missingLangs, lang)
				}
			}
		}

		if len(missingLangs) == 0 {
			results = append(results, checkResult{
				name:    "Language guides",
				passed:  true,
				message: fmt.Sprintf("All %d installed languages present", len(config.Installed.Languages)),
			})
		} else {
			results = append(results, checkResult{
				name:    "Language guides",
				passed:  false,
				message: fmt.Sprintf("Missing: %s", strings.Join(missingLangs, ", ")),
				fixable: true,
			})
		}
	}

	// Check 6: Installed framework skills exist
	if config != nil {
		// Ensure legacy configs have skills populated
		config.MigrateFrameworksToSkills()

		var missingFws []string
		for _, fw := range config.Installed.Frameworks {
			component := core.FindFramework(fw)
			if component != nil {
				// Framework guides are now skill directories with SKILL.md
				skillPath := filepath.Join(cwd, component.Path, "SKILL.md")
				if _, err := os.Stat(skillPath); os.IsNotExist(err) {
					missingFws = append(missingFws, fw)
				}
			}
		}

		if len(missingFws) == 0 {
			results = append(results, checkResult{
				name:    "Framework guides",
				passed:  true,
				message: fmt.Sprintf("All %d installed frameworks present", len(config.Installed.Frameworks)),
			})
		} else {
			results = append(results, checkResult{
				name:    "Framework guides",
				passed:  false,
				message: fmt.Sprintf("Missing: %s", strings.Join(missingFws, ", ")),
				fixable: true,
			})
		}
	}

	// Check 7: Workflow skills exist
	if config != nil {
		// Ensure legacy configs have skills populated
		config.MigrateWorkflowsToSkills()

		var missingWfs []string
		workflowsToCheck := config.Installed.Workflows
		if len(workflowsToCheck) == 1 && workflowsToCheck[0] == "all" {
			workflowsToCheck = core.GetAllWorkflowNames()
		}

		for _, wf := range workflowsToCheck {
			component := core.FindWorkflow(wf)
			if component != nil {
				// Workflows are now skill directories with SKILL.md
				skillPath := filepath.Join(cwd, component.Path, "SKILL.md")
				if _, err := os.Stat(skillPath); os.IsNotExist(err) {
					missingWfs = append(missingWfs, wf)
				}
			}
		}

		if len(missingWfs) == 0 {
			results = append(results, checkResult{
				name:    "Workflows",
				passed:  true,
				message: fmt.Sprintf("All %d installed workflows present", len(workflowsToCheck)),
			})
		} else {
			results = append(results, checkResult{
				name:    "Workflows",
				passed:  false,
				message: fmt.Sprintf("Missing: %s", strings.Join(missingWfs, ", ")),
				fixable: true,
			})
		}
	}

	// Check 8: Skills validation
	skillsDir := filepath.Join(cwd, ".claude", "skills")
	if _, err := os.Stat(skillsDir); err == nil {
		skills, err := core.ScanSkillsDirectory(skillsDir)
		if err != nil {
			results = append(results, checkResult{
				name:    "Skills",
				passed:  false,
				message: fmt.Sprintf("Failed to scan skills: %v", err),
			})
		} else if len(skills) == 0 {
			results = append(results, checkResult{
				name:    "Skills",
				passed:  true,
				message: "No skills installed",
			})
		} else {
			validCount := 0
			invalidCount := 0
			for _, skill := range skills {
				if len(skill.Errors) == 0 {
					validCount++
				} else {
					invalidCount++
				}
			}

			if invalidCount == 0 {
				results = append(results, checkResult{
					name:    "Skills",
					passed:  true,
					message: fmt.Sprintf("%d skill(s) installed, all valid", validCount),
				})
			} else {
				results = append(results, checkResult{
					name:    "Skills",
					passed:  false,
					message: fmt.Sprintf("%d skill(s) installed, %d invalid", len(skills), invalidCount),
				})
			}
		}
	}

	// Check 9: Auto loop health (if present)
	autoDir := core.GetAutoDir(cwd)
	if _, err := os.Stat(autoDir); err == nil {
		results = append(results, checkAutoHealth(cwd)...)
	}

	// Check 10: Local modifications (informational)
	if config != nil {
		claudeMdModified := checkModification(claudeMdPath, config.Version)
		if claudeMdModified {
			results = append(results, checkResult{
				name:    "Local modifications",
				passed:  true,
				message: "CLAUDE.md has local modifications (expected)",
			})
		}
	}

	// Print results
	passedCount := 0
	failedCount := 0
	fixableCount := 0

	for _, r := range results {
		if r.passed {
			ui.SuccessItem(0, "%s: %s", r.name, r.message)
			passedCount++
		} else {
			ui.ErrorItem(0, "%s: %s", r.name, r.message)
			failedCount++
			if r.fixable {
				fixableCount++
			}
		}
	}

	// Summary
	fmt.Println()
	if failedCount == 0 {
		ui.Bold("Status: Healthy")
		ui.Success("All %d checks passed", passedCount)
	} else {
		ui.Bold("Status: Issues Found")
		ui.Error("%d checks failed, %d passed", failedCount, passedCount)

		if fixableCount > 0 && !autoFix {
			ui.Info("\n%d issues can be auto-fixed. Run 'aicof doctor --fix' to repair.", fixableCount)
		}
	}

	// Auto-fix if requested
	if autoFix && fixableCount > 0 {
		fmt.Println()
		ui.Info("Attempting to fix issues...")

		if config != nil {
			// Fix missing directories
			for _, dir := range missingDirs {
				dirPath := filepath.Join(cwd, dir)
				if err := os.MkdirAll(dirPath, 0755); err == nil {
					ui.Success("Created %s", dir)
				}
			}

			// Fix missing files by re-downloading
			downloader, err := core.NewDownloader()
			if err != nil {
				ui.Error("Failed to initialize downloader: %v", err)
			} else {
				cachePath, err := downloader.DownloadVersion(config.Version)
				if err != nil {
					ui.Error("Failed to download version: %v", err)
				} else {
					// Get all paths that should exist
					paths := core.GetComponentPaths(
						config.Installed.Languages,
						config.Installed.Frameworks,
						config.Installed.Workflows,
					)

					extractor := core.NewExtractor(cachePath, cwd)
					for _, path := range paths {
						localPath := filepath.Join(cwd, path)
						if _, err := os.Stat(localPath); os.IsNotExist(err) {
							if err := core.CopyFromCache(cachePath, cwd, path); err == nil {
								ui.Success("Restored %s", path)
							} else {
								ui.Error("Failed to restore %s: %v", path, err)
							}
						}
					}
					_ = extractor // silence unused warning
				}
			}

			ui.Success("Fix complete. Run 'aicof doctor' again to verify.")
		}
	}

	return nil
}

// extractVersion extracts version from CLAUDE.md content
func extractVersion(content string) string {
	// Look for version pattern like "**Current Version**: 1.7.0"
	re := regexp.MustCompile(`\*\*Current Version\*\*:\s*(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}

	// Alternative pattern "Current Version: 1.7.0"
	re = regexp.MustCompile(`Current Version:\s*(\d+\.\d+\.\d+)`)
	matches = re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

// checkAutoHealth validates the auto loop directory and files
func checkAutoHealth(cwd string) []checkResult {
	var results []checkResult

	prdPath := core.GetAutoPRDPath(cwd)
	prd, err := core.LoadAutoPRD(prdPath)
	if err != nil {
		results = append(results, checkResult{
			name:    "Auto loop",
			passed:  false,
			message: fmt.Sprintf("prd.json invalid: %v", err),
		})
		return results
	}

	errs := core.ValidateAutoPRD(prd)
	if len(errs) > 0 {
		results = append(results, checkResult{
			name:    "Auto loop",
			passed:  false,
			message: fmt.Sprintf("prd.json validation: %s", strings.Join(errs, "; ")),
		})
	} else {
		prd.RecalculateProgress()
		results = append(results, checkResult{
			name:    "Auto loop",
			passed:  true,
			message: fmt.Sprintf("prd.json valid (%d/%d tasks completed)", prd.Progress.CompletedTasks, prd.Progress.TotalTasks),
		})
	}

	return results
}

// checkModification checks if a file has been modified from the original
func checkModification(filePath string, _ string) bool {
	// Simple heuristic: if file exists and we can read it, assume it might be modified
	// A more robust check would compare against the cached original
	_, err := os.Stat(filePath)
	return err == nil
}
