package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ar4mirez/aicof/internal/core"
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
	seen := make(map[string]bool) // Deduplicate skill directories
	for _, f := range files {
		if skillName := extractSkillDirName(f); skillName != "" {
			if !seen[skillName] {
				if strings.HasSuffix(skillName, "-guide") {
					langs = append(langs, skillName)
				} else if isFrameworkSkill(skillName) {
					fws = append(fws, skillName)
				}
				seen[skillName] = true
			}
		} else if strings.Contains(f, "workflows/") {
			name := extractComponentName(f)
			wfs = append(wfs, name)
		}
	}
	return
}

// isFrameworkSkill checks if a skill name corresponds to a known framework
func isFrameworkSkill(skillName string) bool {
	return core.FindFramework(skillName) != nil
}

func categorizeOtherFiles(added, modified, removed []string) (addedOther, modifiedOther, removedOther []string) {
	isComponent := func(f string) bool {
		if skillName := extractSkillDirName(f); skillName != "" {
			if strings.HasSuffix(skillName, "-guide") || isFrameworkSkill(skillName) {
				return true
			}
		}
		return strings.Contains(f, "workflows/")
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

// extractSkillDirName extracts the skill directory name from a path like
// ".agent/skills/go-guide/SKILL.md" -> "go-guide".
// Returns empty string if the path is not inside .agent/skills/.
func extractSkillDirName(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for i, part := range parts {
		if part == "skills" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
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
