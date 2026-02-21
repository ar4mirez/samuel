package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ar4mirez/samuel/internal/core"
	"github.com/ar4mirez/samuel/internal/ui"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Samuel framework to the latest version",
	Long: `Update the installed Samuel framework to the latest version.

This command will:
1. Check for available updates
2. Download the new version
3. Apply updates while preserving local modifications
4. Create backups of modified files

Examples:
  samuel update              # Update to latest version
  samuel update --check      # Check for updates without applying
  samuel update --diff       # Show what will change
  samuel update --force      # Overwrite local modifications`,
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().Bool("check", false, "Check for updates without applying")
	updateCmd.Flags().Bool("diff", false, "Show what files will change")
	updateCmd.Flags().BoolP("force", "f", false, "Overwrite local modifications")
	updateCmd.Flags().String("version", "", "Update to specific version")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	checkOnly, _ := cmd.Flags().GetBool("check")
	showDiff, _ := cmd.Flags().GetBool("diff")
	force, _ := cmd.Flags().GetBool("force")
	targetVersion, _ := cmd.Flags().GetString("version")

	// Load existing config
	config, err := core.LoadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no Samuel installation found. Run 'samuel init' first")
		}
		return fmt.Errorf("failed to load config: %w", err)
	}

	currentVersion := config.Version

	// Initialize downloader
	downloader, err := core.NewDownloader()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// Get target version (latest or specified)
	if targetVersion == "" {
		spinner := ui.NewSpinner("Checking for updates...")
		spinner.Start()

		latest, err := downloader.GetLatestVersion()
		if err != nil {
			spinner.Error("Failed to check for updates")
			return fmt.Errorf("failed to get latest version: %w", err)
		}
		spinner.Stop()
		targetVersion = latest
	}

	ui.Bold("Samuel Update")
	ui.TableRow("Current version", currentVersion)
	ui.TableRow("Target version", targetVersion)

	// Check if update needed
	if currentVersion == targetVersion && !force {
		fmt.Println()
		ui.Success("Already up to date!")
		return nil
	}

	// Check-only mode
	if checkOnly {
		if currentVersion != targetVersion {
			fmt.Println()
			ui.Success("Update available: %s → %s", currentVersion, targetVersion)
			ui.Info("Run 'samuel update' to apply")
		}
		return nil
	}

	// Download new version
	spinner := ui.NewSpinner("Downloading...")
	spinner.Start()

	cachePath, err := downloader.DownloadVersion(targetVersion)
	if err != nil {
		spinner.Error("Download failed")
		return fmt.Errorf("failed to download: %w", err)
	}
	spinner.Success(fmt.Sprintf("Downloaded v%s", targetVersion))

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Compute what needs to be updated
	paths := core.GetComponentPaths(
		config.Installed.Languages,
		config.Installed.Frameworks,
		config.Installed.Workflows,
	)

	// Check for local modifications
	extractor := core.NewExtractor(cachePath, cwd)
	changes := categorizeFileChanges(paths, cwd, cachePath)
	modifiedFiles := changes.modifiedFiles
	newFiles := changes.newFiles
	unchangedFiles := changes.unchangedFiles

	// Show diff if requested
	if showDiff {
		fmt.Println()
		ui.Section("Changes")

		if len(newFiles) > 0 {
			ui.ListItem(1, "%d new files:", len(newFiles))
			for _, f := range newFiles {
				ui.SuccessItem(2, "%s", f)
			}
		}

		if len(modifiedFiles) > 0 {
			ui.ListItem(1, "%d files with local modifications:", len(modifiedFiles))
			for _, f := range modifiedFiles {
				ui.WarnItem(2, "%s", f)
			}
		}

		if len(unchangedFiles) > 0 {
			ui.ListItem(1, "%d files to update:", len(unchangedFiles))
		}

		fmt.Println()
		if !force {
			ui.Info("Modified files will be preserved. Use --force to overwrite.")
		}
		return nil
	}

	// Create backup if there are modified files
	var backupDir string
	if len(modifiedFiles) > 0 && !force {
		backupDir = filepath.Join(cwd, fmt.Sprintf(".samuel-backup-%s", time.Now().Format("20060102-150405")))
		if err := os.MkdirAll(backupDir, 0755); err != nil {
			return fmt.Errorf("failed to create backup directory: %w", err)
		}

		for _, f := range modifiedFiles {
			if err := extractor.BackupFile(f, backupDir); err != nil {
				ui.Warn("Failed to backup %s: %v", f, err)
			}
		}
		ui.Success("Backed up %d modified files to %s", len(modifiedFiles), backupDir)
	}

	// Determine which files to update
	var filesToUpdate []string
	filesToUpdate = append(filesToUpdate, newFiles...)
	filesToUpdate = append(filesToUpdate, unchangedFiles...)

	if force {
		filesToUpdate = append(filesToUpdate, modifiedFiles...)
	}

	// Apply updates
	result, err := extractor.Extract(filesToUpdate, true)
	if err != nil {
		return fmt.Errorf("failed to apply updates: %w", err)
	}

	// Report results
	ui.Success("Updated %d files", len(result.FilesCreated))

	if len(newFiles) > 0 {
		ui.Success("Added %d new files", len(newFiles))
	}

	if len(modifiedFiles) > 0 && !force {
		ui.Warn("Preserved %d locally modified files", len(modifiedFiles))
		if backupDir != "" {
			ui.Info("Backups saved to: %s", backupDir)
		}
	}

	// Update config version
	config.Version = targetVersion
	if err := config.Save(cwd); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	ui.Success("Updated samuel.yaml to v%s", targetVersion)

	// Show modified files warning
	if len(modifiedFiles) > 0 && !force {
		fmt.Println()
		ui.Bold("Modified files preserved:")
		for _, f := range modifiedFiles {
			ui.WarnItem(1, "%s", f)
		}
		ui.Info("\nTo see changes: diff -u %s/<file> <file>", backupDir)
		ui.Info("To accept new version: cp %s/<file> <file>", backupDir)
	}

	return nil
}

// fileChanges holds the categorized file lists from comparing local vs cached files.
type fileChanges struct {
	newFiles       []string
	modifiedFiles  []string
	unchangedFiles []string
}

// categorizeFileChanges compares component paths between the local project and
// the cache, categorizing each file as new, modified, or unchanged.
func categorizeFileChanges(paths []string, cwd, cachePath string) fileChanges {
	var changes fileChanges

	for _, path := range paths {
		localPath := filepath.Join(cwd, path)
		cacheSrcPath := filepath.Join(cachePath, path)

		if !fileExists(cacheSrcPath) {
			continue
		}

		if !fileExists(localPath) {
			changes.newFiles = append(changes.newFiles, path)
			continue
		}

		localContent, err := os.ReadFile(localPath)
		if err != nil {
			ui.Warn("Skipping %s: failed to read local file: %v", path, err)
			continue
		}

		cacheContent, err := os.ReadFile(cacheSrcPath)
		if err != nil {
			ui.Warn("Skipping %s: failed to read cached file: %v", path, err)
			continue
		}

		if string(localContent) != string(cacheContent) {
			changes.modifiedFiles = append(changes.modifiedFiles, path)
		} else {
			changes.unchangedFiles = append(changes.unchangedFiles, path)
		}
	}

	return changes
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
