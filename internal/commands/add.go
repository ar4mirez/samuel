package commands

import (
	"fmt"
	"os"

	"github.com/ar4mirez/aicof/internal/core"
	"github.com/ar4mirez/aicof/internal/ui"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <type> <name>",
	Short: "Add a component to your project",
	Long: `Add a language guide, framework guide, or workflow to your project.

Types:
  language   Add a language guide (e.g., rust, kotlin)
  framework  Add a framework guide (e.g., django, rails)
  workflow   Add a workflow (e.g., security-audit)

Examples:
  aicof add language rust
  aicof add framework django
  aicof add workflow security-audit`,
	Args: cobra.ExactArgs(2),
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	componentType := args[0]
	componentName := args[1]

	// Load config
	config, err := core.LoadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no AICoF installation found. Run 'aicof init' first")
		}
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate component type and find component
	var component *core.Component
	var alreadyInstalled bool

	switch componentType {
	case "language", "lang", "l":
		component = core.FindLanguage(componentName)
		if component == nil {
			return fmt.Errorf("unknown language: %s\nRun 'aicof list --available --type languages' to see available languages", componentName)
		}
		alreadyInstalled = config.HasLanguage(componentName)

	case "framework", "fw", "f":
		component = core.FindFramework(componentName)
		if component == nil {
			return fmt.Errorf("unknown framework: %s\nRun 'aicof list --available --type frameworks' to see available frameworks", componentName)
		}
		alreadyInstalled = config.HasFramework(componentName)

	case "workflow", "wf", "w":
		component = core.FindWorkflow(componentName)
		if component == nil {
			return fmt.Errorf("unknown workflow: %s\nRun 'aicof list --available --type workflows' to see available workflows", componentName)
		}
		alreadyInstalled = config.HasWorkflow(componentName)

	default:
		return fmt.Errorf("unknown component type: %s\nValid types: language, framework, workflow", componentType)
	}

	// Check if already installed
	if alreadyInstalled {
		ui.Warn("%s '%s' is already installed", componentType, componentName)
		return nil
	}

	// Download the component
	spinner := ui.NewSpinner(fmt.Sprintf("Downloading %s...", component.Name))
	spinner.Start()

	downloader, err := core.NewDownloader()
	if err != nil {
		spinner.Error("Failed to initialize")
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// Download the specific version
	cachePath, err := downloader.DownloadVersion(config.Version)
	if err != nil {
		spinner.Error("Download failed")
		return fmt.Errorf("failed to download: %w", err)
	}
	spinner.Stop()

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Copy the file
	if err := core.CopyFromCache(cachePath, cwd, component.Path); err != nil {
		return fmt.Errorf("failed to install %s: %w", component.Name, err)
	}

	// Update config
	switch componentType {
	case "language", "lang", "l":
		config.AddLanguage(componentName)
	case "framework", "fw", "f":
		config.AddFramework(componentName)
	case "workflow", "wf", "w":
		config.AddWorkflow(componentName)
	}

	if err := config.Save(cwd); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	ui.Success("Installed %s", component.Path)
	ui.Success("Updated aicof.yaml")

	return nil
}
