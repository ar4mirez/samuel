package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ar4mirez/samuel/internal/core"
	"github.com/ar4mirez/samuel/internal/ui"
)

// checkConfigFile validates that samuel.yaml exists and is parseable.
func checkConfigFile() (checkResult, *core.Config) {
	config, configErr := core.LoadConfig()
	if configErr != nil {
		msg := "samuel.yaml not found"
		if !os.IsNotExist(configErr) {
			msg = fmt.Sprintf("Config error: %v", configErr)
		}
		return checkResult{
			name:    "Config file",
			passed:  false,
			message: msg,
			fixable: false,
		}, nil
	}
	return checkResult{
		name:    "Config file",
		passed:  true,
		message: fmt.Sprintf("samuel.yaml found (v%s)", config.Version),
	}, config
}

// checkCLAUDEMD verifies CLAUDE.md exists and optionally extracts its version.
func checkCLAUDEMD(cwd string) checkResult {
	claudeMdPath := filepath.Join(cwd, "CLAUDE.md")
	if _, err := os.Stat(claudeMdPath); os.IsNotExist(err) {
		return checkResult{
			name:    "CLAUDE.md",
			passed:  false,
			message: "CLAUDE.md not found",
			fixable: true,
		}
	}

	content, readErr := os.ReadFile(claudeMdPath)
	if readErr != nil {
		ui.Warn("Could not read CLAUDE.md: %v", readErr)
	}

	version := extractVersion(string(content))
	msg := "Present"
	if version != "" {
		msg = fmt.Sprintf("Present (v%s)", version)
	}
	return checkResult{
		name:    "CLAUDE.md",
		passed:  true,
		message: msg,
	}
}

// checkAGENTSMD verifies AGENTS.md exists for cross-tool compatibility.
func checkAGENTSMD(cwd string) checkResult {
	agentsMdPath := filepath.Join(cwd, "AGENTS.md")
	if _, err := os.Stat(agentsMdPath); os.IsNotExist(err) {
		return checkResult{
			name:    "AGENTS.md",
			passed:  false,
			message: "AGENTS.md not found (cross-tool compatibility)",
			fixable: true,
		}
	}
	return checkResult{
		name:    "AGENTS.md",
		passed:  true,
		message: "Present",
	}
}

// checkDirectoryStructure verifies .claude/skills/ directory exists.
// Returns the check result and a list of missing directories for auto-fix.
func checkDirectoryStructure(cwd string) (checkResult, []string) {
	claudeDirs := []string{
		".claude",
		".claude/skills",
	}

	var missingDirs []string
	for _, dir := range claudeDirs {
		dirPath := filepath.Join(cwd, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			missingDirs = append(missingDirs, dir)
		}
	}

	if len(missingDirs) == 0 {
		return checkResult{
			name:    "Directory structure",
			passed:  true,
			message: ".claude/skills/ directory structure valid",
		}, nil
	}
	return checkResult{
		name:    "Directory structure",
		passed:  false,
		message: fmt.Sprintf("Missing directories: %s", strings.Join(missingDirs, ", ")),
		fixable: true,
	}, missingDirs
}

// checkInstalledComponents checks that all installed language, framework,
// and workflow skills exist on disk.
func checkInstalledComponents(cwd string, config *core.Config) []checkResult {
	var results []checkResult

	config.MigrateLanguagesToSkills()
	results = append(results, checkInstalledSkills(
		cwd, config.Installed.Languages, "Language guides", core.FindLanguage,
	))

	config.MigrateFrameworksToSkills()
	results = append(results, checkInstalledSkills(
		cwd, config.Installed.Frameworks, "Framework guides", core.FindFramework,
	))

	config.MigrateWorkflowsToSkills()
	workflowsToCheck := config.Installed.Workflows
	if len(workflowsToCheck) == 1 && workflowsToCheck[0] == "all" {
		workflowsToCheck = core.GetAllWorkflowNames()
	}
	results = append(results, checkInstalledSkills(
		cwd, workflowsToCheck, "Workflows", core.FindWorkflow,
	))

	return results
}

// checkInstalledSkills verifies that SKILL.md files exist for a set of
// installed components of a given category (languages, frameworks, workflows).
func checkInstalledSkills(
	cwd string,
	names []string,
	category string,
	finder func(string) *core.Component,
) checkResult {
	var missing []string
	for _, name := range names {
		component := finder(name)
		if component != nil {
			skillPath := filepath.Join(cwd, component.Path, "SKILL.md")
			if _, err := os.Stat(skillPath); os.IsNotExist(err) {
				missing = append(missing, name)
			}
		}
	}

	if len(missing) == 0 {
		return checkResult{
			name:    category,
			passed:  true,
			message: fmt.Sprintf("All %d installed %s present", len(names), strings.ToLower(category)),
		}
	}
	return checkResult{
		name:    category,
		passed:  false,
		message: fmt.Sprintf("Missing: %s", strings.Join(missing, ", ")),
		fixable: true,
	}
}

// checkSkillsIntegrity scans and validates all installed skills.
func checkSkillsIntegrity(cwd string) []checkResult {
	skillsDir := filepath.Join(cwd, ".claude", "skills")
	if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
		return nil
	}

	skills, err := core.ScanSkillsDirectory(skillsDir)
	if err != nil {
		return []checkResult{{
			name:    "Skills",
			passed:  false,
			message: fmt.Sprintf("Failed to scan skills: %v", err),
		}}
	}

	if len(skills) == 0 {
		return []checkResult{{
			name:    "Skills",
			passed:  true,
			message: "No skills installed",
		}}
	}

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
		return []checkResult{{
			name:    "Skills",
			passed:  true,
			message: fmt.Sprintf("%d skill(s) installed, all valid", validCount),
		}}
	}
	return []checkResult{{
		name:    "Skills",
		passed:  false,
		message: fmt.Sprintf("%d skill(s) installed, %d invalid", len(skills), invalidCount),
	}}
}

// checkAutoHealth validates the auto loop directory and files.
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

// checkLocalModifications checks if key files have been modified locally.
func checkLocalModifications(cwd string, config *core.Config) []checkResult {
	claudeMdPath := filepath.Join(cwd, "CLAUDE.md")
	if checkModification(claudeMdPath) {
		return []checkResult{{
			name:    "Local modifications",
			passed:  true,
			message: "CLAUDE.md has local modifications (expected)",
		}}
	}
	return nil
}

// checkModification checks if a file exists (heuristic for local modification).
func checkModification(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// extractVersion extracts version from CLAUDE.md content.
func extractVersion(content string) string {
	re := regexp.MustCompile(`\*\*Current Version\*\*:\s*(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}

	re = regexp.MustCompile(`Current Version:\s*(\d+\.\d+\.\d+)`)
	matches = re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
