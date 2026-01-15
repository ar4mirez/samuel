package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ar4mirez/aicof/internal/ui"
)

func displayFileDiff(diff *VersionDiff) {
	ui.Bold("AICoF Version Diff")
	fmt.Println()
	ui.Print("Comparing: %s → %s", diff.FromVersion, diff.ToVersion)
	fmt.Println()

	if len(diff.Added) == 0 && len(diff.Modified) == 0 && len(diff.Removed) == 0 {
		ui.Success("No differences found")
		return
	}

	if len(diff.Added) > 0 {
		ui.Section(fmt.Sprintf("Added (%d files)", len(diff.Added)))
		for _, path := range diff.Added {
			fmt.Printf("  + %s\n", path)
		}
	}

	if len(diff.Modified) > 0 {
		ui.Section(fmt.Sprintf("Modified (%d files)", len(diff.Modified)))
		for _, path := range diff.Modified {
			fmt.Printf("  ~ %s\n", path)
		}
	}

	if len(diff.Removed) > 0 {
		ui.Section(fmt.Sprintf("Removed (%d files)", len(diff.Removed)))
		for _, path := range diff.Removed {
			fmt.Printf("  - %s\n", path)
		}
	}

	fmt.Println()
	ui.Dim("Summary: %d added, %d modified, %d removed, %d unchanged",
		len(diff.Added), len(diff.Modified), len(diff.Removed), diff.Unchanged)

	if len(diff.Added) > 0 || len(diff.Modified) > 0 {
		fmt.Println()
		ui.Info("Run 'aicof update' to apply these changes")
	}
}

func displayComponentDiff(diff *VersionDiff) {
	ui.Bold("AICoF Component Changes")
	fmt.Println()
	ui.Print("Comparing: %s → %s", diff.FromVersion, diff.ToVersion)
	fmt.Println()

	// Categorize changes by component type
	addedLangs, addedFws, addedWfs := categorizeFiles(diff.Added)
	modifiedLangs, modifiedFws, modifiedWfs := categorizeFiles(diff.Modified)
	removedLangs, removedFws, removedWfs := categorizeFiles(diff.Removed)

	// Languages
	if len(addedLangs) > 0 || len(modifiedLangs) > 0 || len(removedLangs) > 0 {
		ui.Section("Languages")
		displayComponentChanges(addedLangs, modifiedLangs, removedLangs)
	}

	// Frameworks
	if len(addedFws) > 0 || len(modifiedFws) > 0 || len(removedFws) > 0 {
		ui.Section("Frameworks")
		displayComponentChanges(addedFws, modifiedFws, removedFws)
	}

	// Workflows
	if len(addedWfs) > 0 || len(modifiedWfs) > 0 || len(removedWfs) > 0 {
		ui.Section("Workflows")
		displayComponentChanges(addedWfs, modifiedWfs, removedWfs)
	}

	// Other files
	addedOther, modifiedOther, removedOther := categorizeOtherFiles(diff.Added, diff.Modified, diff.Removed)
	if len(addedOther) > 0 || len(modifiedOther) > 0 || len(removedOther) > 0 {
		ui.Section("Other Files")
		displayComponentChanges(addedOther, modifiedOther, removedOther)
	}
}

func categorizeFiles(files []string) (langs, fws, wfs []string) {
	for _, f := range files {
		name := extractComponentName(f)
		if strings.Contains(f, "language-guides/") {
			langs = append(langs, name)
		} else if strings.Contains(f, "framework-guides/") {
			fws = append(fws, name)
		} else if strings.Contains(f, "workflows/") {
			wfs = append(wfs, name)
		}
	}
	return
}

func categorizeOtherFiles(added, modified, removed []string) (addedOther, modifiedOther, removedOther []string) {
	isComponent := func(f string) bool {
		return strings.Contains(f, "language-guides/") ||
			strings.Contains(f, "framework-guides/") ||
			strings.Contains(f, "workflows/")
	}

	for _, f := range added {
		if !isComponent(f) {
			addedOther = append(addedOther, f)
		}
	}
	for _, f := range modified {
		if !isComponent(f) {
			modifiedOther = append(modifiedOther, f)
		}
	}
	for _, f := range removed {
		if !isComponent(f) {
			removedOther = append(removedOther, f)
		}
	}
	return
}

func extractComponentName(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, ".md")
}

func displayComponentChanges(added, modified, removed []string) {
	for _, name := range added {
		fmt.Printf("  + %s (new)\n", name)
	}
	for _, name := range modified {
		fmt.Printf("  ~ %s (updated)\n", name)
	}
	for _, name := range removed {
		fmt.Printf("  - %s (removed)\n", name)
	}
}
