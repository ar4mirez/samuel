package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ar4mirez/aicof/internal/core"
	"github.com/ar4mirez/aicof/internal/ui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize AICoF framework in a project",
	Long: `Initialize the AICoF (Artificial Intelligence Coding Framework) in a new or existing project.

This command downloads and installs framework files including:
- CLAUDE.md (core guardrails and methodology)
- Language guides (TypeScript, Python, Go, etc.)
- Framework guides (React, FastAPI, etc.)
- Workflows (PRD creation, code review, etc.)

Examples:
  aicof init my-project              # Create new project
  aicof init .                       # Initialize in current directory
  aicof init --template minimal      # Use minimal template
  aicof init --languages ts,py,go    # Select specific languages`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("template", "t", "", "Template: full, starter, minimal")
	initCmd.Flags().StringSlice("languages", nil, "Languages to install (comma-separated)")
	initCmd.Flags().StringSlice("frameworks", nil, "Frameworks to install (comma-separated)")
	initCmd.Flags().BoolP("force", "f", false, "Overwrite existing files")
	initCmd.Flags().Bool("non-interactive", false, "Skip prompts, use defaults")
}

func runInit(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	templateName, _ := cmd.Flags().GetString("template")
	languageFlags, _ := cmd.Flags().GetStringSlice("languages")
	frameworkFlags, _ := cmd.Flags().GetStringSlice("frameworks")

	// Determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	// Resolve to absolute path
	absTargetDir, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Check if directory exists or needs to be created
	createDir := false
	if targetDir != "." {
		if _, err := os.Stat(absTargetDir); os.IsNotExist(err) {
			createDir = true
		}
	}

	// Check for existing config
	if core.ConfigExists(absTargetDir) && !force {
		return fmt.Errorf("AICoF already initialized in %s. Use --force to reinitialize", absTargetDir)
	}

	// Variables to collect user choices
	var selectedTemplate *core.Template
	var selectedLanguages []string
	var selectedFrameworks []string

	// Interactive mode
	if !nonInteractive && templateName == "" && len(languageFlags) == 0 {
		// Select template
		templateOptions := make([]ui.SelectOption, len(core.Templates))
		for i, t := range core.Templates {
			templateOptions[i] = ui.SelectOption{
				Name:        t.Name,
				Description: t.Description,
				Value:       t.Name,
			}
		}

		selected, err := ui.Select("Select template", templateOptions)
		if err != nil {
			return fmt.Errorf("template selection cancelled: %w", err)
		}
		templateName = selected.Value
	}

	// Get template or use defaults
	if templateName != "" {
		selectedTemplate = core.FindTemplate(templateName)
		if selectedTemplate == nil {
			return fmt.Errorf("unknown template: %s", templateName)
		}
		selectedLanguages = selectedTemplate.Languages
		selectedFrameworks = selectedTemplate.Frameworks
	}

	// Override with flags if provided
	if len(languageFlags) > 0 {
		selectedLanguages = expandLanguages(languageFlags)
	}
	if len(frameworkFlags) > 0 {
		selectedFrameworks = expandFrameworks(frameworkFlags)
	}

	// Interactive language selection if not full template and not specified
	if !nonInteractive && selectedTemplate != nil && selectedTemplate.Name != "full" && len(languageFlags) == 0 {
		langOptions := make([]ui.SelectOption, len(core.Languages))
		for i, l := range core.Languages {
			langOptions[i] = ui.SelectOption{
				Name:        l.Name,
				Description: l.Description,
				Value:       l.Name,
			}
		}

		// Pre-select languages from template
		defaults := selectedLanguages

		selected, err := ui.MultiSelect("Select languages", langOptions, defaults)
		if err != nil {
			return fmt.Errorf("language selection cancelled: %w", err)
		}

		selectedLanguages = make([]string, len(selected))
		for i, s := range selected {
			selectedLanguages[i] = s.Value
		}
	}

	// Interactive framework selection if not full template and not specified
	if !nonInteractive && selectedTemplate != nil && selectedTemplate.Name != "full" && len(frameworkFlags) == 0 && len(selectedLanguages) > 0 {
		// Filter frameworks based on selected languages
		relevantFrameworks := getRelevantFrameworks(selectedLanguages)

		if len(relevantFrameworks) > 0 {
			fwOptions := make([]ui.SelectOption, len(relevantFrameworks))
			for i, fw := range relevantFrameworks {
				fwOptions[i] = ui.SelectOption{
					Name:        fw.Name,
					Description: fw.Description,
					Value:       fw.Name,
				}
			}

			selected, err := ui.MultiSelect("Select frameworks (optional)", fwOptions, nil)
			if err != nil {
				// User cancelled - continue without frameworks
				selectedFrameworks = []string{}
			} else {
				selectedFrameworks = make([]string, len(selected))
				for i, s := range selected {
					selectedFrameworks[i] = s.Value
				}
			}
		}
	}

	// Default to starter template if nothing selected
	if selectedTemplate == nil && len(selectedLanguages) == 0 {
		selectedTemplate = core.FindTemplate("starter")
		selectedLanguages = selectedTemplate.Languages
		selectedFrameworks = selectedTemplate.Frameworks
	}

	// Show what will be installed
	ui.Header("AICoF Initialization")
	ui.TableRow("Target", absTargetDir)
	ui.TableRow("Languages", fmt.Sprintf("%d selected", len(selectedLanguages)))
	ui.TableRow("Frameworks", fmt.Sprintf("%d selected", len(selectedFrameworks)))
	ui.TableRow("Workflows", "all (13)")

	// Confirm in interactive mode
	if !nonInteractive {
		confirmed, err := ui.Confirm("\nProceed with installation?", true)
		if err != nil || !confirmed {
			ui.Info("Installation cancelled")
			return nil
		}
	}

	// Start installation
	spinner := ui.NewSpinner("Downloading framework...")
	spinner.Start()

	// Initialize downloader
	downloader, err := core.NewDownloader()
	if err != nil {
		spinner.Error("Failed to initialize")
		return fmt.Errorf("failed to initialize downloader: %w", err)
	}

	// Get latest version
	version, err := downloader.GetLatestVersion()
	if err != nil {
		spinner.Error("Failed to get latest version")
		return fmt.Errorf("failed to get latest version: %w", err)
	}

	// Download to cache
	cachePath, err := downloader.DownloadVersion(version)
	if err != nil {
		spinner.Error("Download failed")
		return fmt.Errorf("failed to download framework: %w", err)
	}
	spinner.Success(fmt.Sprintf("Downloaded AICoF v%s", version))

	// Create target directory if needed
	if createDir {
		if err := os.MkdirAll(absTargetDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		ui.Success("Created %s/", filepath.Base(absTargetDir))
	}

	// Get paths to extract
	workflows := []string{"all"}
	paths := core.GetComponentPaths(selectedLanguages, selectedFrameworks, workflows)

	// Extract files
	extractor := core.NewExtractor(cachePath, absTargetDir)
	result, err := extractor.Extract(paths, force)
	if err != nil {
		return fmt.Errorf("failed to extract files: %w", err)
	}

	// Report results
	ui.Success("Installed CLAUDE.md (v%s)", version)
	ui.Success("Installed %d language guides", len(selectedLanguages))
	ui.Success("Installed %d framework guides", len(selectedFrameworks))
	ui.Success("Installed %d workflows", len(core.Workflows))

	if len(result.FilesSkipped) > 0 {
		ui.Warn("Skipped %d existing files (use --force to overwrite)", len(result.FilesSkipped))
	}

	if len(result.Errors) > 0 {
		for _, e := range result.Errors {
			ui.Error("%v", e)
		}
	}

	// Create config file
	config := core.NewConfig(version)
	config.Installed.Languages = selectedLanguages
	config.Installed.Frameworks = selectedFrameworks
	config.Installed.Workflows = workflows

	if err := config.Save(absTargetDir); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	ui.Success("Created aicof.yaml")

	// Show next steps
	fmt.Println()
	ui.Bold("Next steps:")
	if createDir {
		ui.ListItem(1, "cd %s", filepath.Base(absTargetDir))
	}
	ui.ListItem(1, "Start coding with AI assistance!")
	ui.ListItem(1, "Run 'aicof doctor' to verify installation")

	return nil
}

// expandLanguages expands short language names
func expandLanguages(flags []string) []string {
	aliases := map[string]string{
		"ts":         "typescript",
		"js":         "typescript",
		"javascript": "typescript",
		"py":         "python",
		"cs":         "csharp",
		"c++":        "cpp",
		"c":          "cpp",
		"rb":         "ruby",
		"sh":         "shell",
		"bash":       "shell",
	}

	var result []string
	for _, f := range flags {
		// Handle comma-separated values
		for _, name := range strings.Split(f, ",") {
			name = strings.TrimSpace(strings.ToLower(name))
			if alias, ok := aliases[name]; ok {
				name = alias
			}
			// Verify it exists
			if core.FindLanguage(name) != nil {
				result = append(result, name)
			}
		}
	}
	return result
}

// expandFrameworks expands short framework names
func expandFrameworks(flags []string) []string {
	aliases := map[string]string{
		"next":   "nextjs",
		"spring": "spring-boot-java",
	}

	var result []string
	for _, f := range flags {
		for _, name := range strings.Split(f, ",") {
			name = strings.TrimSpace(strings.ToLower(name))
			if alias, ok := aliases[name]; ok {
				name = alias
			}
			if core.FindFramework(name) != nil {
				result = append(result, name)
			}
		}
	}
	return result
}

// getRelevantFrameworks returns frameworks related to selected languages
func getRelevantFrameworks(languages []string) []core.Component {
	// Map languages to their frameworks
	languageFrameworks := map[string][]string{
		"typescript": {"react", "nextjs", "express"},
		"python":     {"django", "fastapi", "flask"},
		"go":         {"gin", "echo", "fiber"},
		"rust":       {"axum", "actix-web", "rocket"},
		"kotlin":     {"spring-boot-kotlin", "ktor", "android-compose"},
		"java":       {"spring-boot-java", "quarkus", "micronaut"},
		"csharp":     {"aspnet-core", "blazor", "unity"},
		"php":        {"laravel", "symfony", "wordpress"},
		"swift":      {"swiftui", "uikit", "vapor"},
		"ruby":       {"rails", "sinatra", "hanami"},
		"dart":       {"flutter", "shelf", "dart-frog"},
	}

	var result []core.Component
	seen := make(map[string]bool)

	for _, lang := range languages {
		if fws, ok := languageFrameworks[lang]; ok {
			for _, fwName := range fws {
				if !seen[fwName] {
					if fw := core.FindFramework(fwName); fw != nil {
						result = append(result, *fw)
						seen[fwName] = true
					}
				}
			}
		}
	}

	return result
}
