package commands

import (
	"fmt"
	"strings"

	"github.com/ar4mirez/samuel/internal/core"
	"github.com/ar4mirez/samuel/internal/ui"
)

// displaySkillMetadata renders the metadata section of a skill info display.
func displaySkillMetadata(info *core.SkillInfo) {
	ui.Section("Metadata")
	ui.TableRow("Name", info.Metadata.Name)

	desc := strings.TrimSpace(info.Metadata.Description)
	if strings.Contains(desc, "\n") {
		ui.Print("  Description:")
		for _, line := range strings.Split(desc, "\n") {
			ui.Print("    %s", strings.TrimSpace(line))
		}
	} else {
		ui.TableRow("Description", desc)
	}

	if info.Metadata.License != "" {
		ui.TableRow("License", info.Metadata.License)
	}

	if info.Metadata.Compatibility != "" {
		ui.TableRow("Compatibility", info.Metadata.Compatibility)
	}

	if len(info.Metadata.Metadata) > 0 {
		ui.Print("  Custom metadata:")
		for k, v := range info.Metadata.Metadata {
			ui.Print("    %s: %s", k, v)
		}
	}
}

// displaySkillStructure renders the structure and stats section of a skill info display.
func displaySkillStructure(info *core.SkillInfo) {
	ui.Section("Structure")
	ui.TableRow("Path", info.Path)

	dirs := []string{}
	if info.HasScripts {
		dirs = append(dirs, "scripts/")
	}
	if info.HasRefs {
		dirs = append(dirs, "references/")
	}
	if info.HasAssets {
		dirs = append(dirs, "assets/")
	}
	if len(dirs) > 0 {
		ui.TableRow("Directories", strings.Join(dirs, ", "))
	}

	if info.Body != "" {
		lines := core.CountLines(info.Body)
		ui.TableRow("Body lines", fmt.Sprintf("%d", lines))
		if lines > 500 {
			ui.WarnItem(1, "Consider splitting content >500 lines")
		}
	}
}
