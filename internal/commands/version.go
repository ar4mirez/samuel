package commands

import (
	"fmt"
	"os"

	"github.com/ar4mirez/samuel/internal/core"
	"github.com/ar4mirez/samuel/internal/ui"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show CLI and framework version information",
	Long: `Display version information for the Samuel CLI and installed framework.

Examples:
  samuel version              # Show version info
  samuel version --check      # Check for updates`,
	RunE: runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().Bool("check", false, "Check for available updates")
}

func runVersion(cmd *cobra.Command, args []string) error {
	checkUpdate, _ := cmd.Flags().GetBool("check")

	// Show CLI version
	ui.Bold("Samuel CLI")
	ui.TableRow("Version", Version)
	ui.TableRow("Commit", Commit)
	ui.TableRow("Built", BuildDate)

	// Try to load local config for framework version
	config, err := core.LoadConfig()
	if err == nil {
		fmt.Println()
		ui.Bold("Installed Framework")
		ui.TableRow("Version", config.Version)
		ui.TableRow("Languages", fmt.Sprintf("%d installed", len(config.Installed.Languages)))
		ui.TableRow("Frameworks", fmt.Sprintf("%d installed", len(config.Installed.Frameworks)))

		workflowCount := len(config.Installed.Workflows)
		if len(config.Installed.Workflows) == 1 && config.Installed.Workflows[0] == "all" {
			workflowCount = len(core.Workflows)
		}
		ui.TableRow("Workflows", fmt.Sprintf("%d installed", workflowCount))
	} else if !os.IsNotExist(err) {
		ui.Warn("Could not load framework config: %v", err)
	} else {
		fmt.Println()
		ui.Dim("No Samuel framework installed in current directory")
	}

	// Check for updates if requested
	if checkUpdate {
		fmt.Println()
		ui.Info("Checking for updates...")

		downloader, err := core.NewDownloader()
		if err != nil {
			return fmt.Errorf("failed to initialize downloader: %w", err)
		}

		// Check CLI updates
		cliInfo, err := downloader.CheckForUpdates(Version)
		if err != nil {
			ui.Warn("Could not check for CLI updates: %v", err)
		} else {
			if cliInfo.UpdateNeeded {
				ui.Success("New CLI version available: %s → %s", Version, cliInfo.Latest)
				ui.Info("Update with: samuel self-update")
			} else {
				ui.Success("CLI is up to date")
			}
		}

		// Check framework updates if installed
		if config != nil {
			fwInfo, err := downloader.CheckForUpdates(config.Version)
			if err != nil {
				ui.Warn("Could not check for framework updates: %v", err)
			} else {
				if fwInfo.UpdateNeeded {
					ui.Success("New framework version available: %s → %s", config.Version, fwInfo.Latest)
					ui.Info("Update with: samuel update")
				} else {
					ui.Success("Framework is up to date")
				}
			}
		}
	}

	return nil
}
