package commands

import (
	"fmt"

	"github.com/ar4mirez/samuel/internal/core"
	"github.com/ar4mirez/samuel/internal/ui"
)

type pilotStats struct {
	discoveryCount int
	implCount      int
}

func printPilotDryRun(autoCfg core.AutoConfig, pilotCfg *core.PilotConfig, cwd string) error {
	ui.Header("Dry Run - Pilot Mode")
	ui.Print("  AI Tool:           %s", autoCfg.AITool)
	ui.Print("  Max Iterations:    %d", autoCfg.MaxIterations)
	ui.Print("  Discover Interval: %d", pilotCfg.DiscoverInterval)
	ui.Print("  Max Tasks/Disc:    %d", pilotCfg.MaxDiscoveryTasks)
	ui.Print("  Sandbox:           %s", autoCfg.Sandbox)
	if pilotCfg.Focus != "" {
		ui.Print("  Focus:             %s", pilotCfg.Focus)
	}
	ui.Print("  Project:           %s", cwd)
	ui.Print("")
	ui.Print("  Quality checks:")
	for _, check := range autoCfg.QualityChecks {
		ui.Print("    - %s", check)
	}
	ui.Print("")
	ui.Print("  Loop plan:")
	ui.Print("    1. Discovery: analyze project, generate tasks")
	ui.Print("    2. Implementation: pick and implement top task")
	ui.Print("    3. Repeat until max iterations or no more work")
	ui.Print("")
	ui.Info("Run without --dry-run to execute")
	return nil
}

func printPilotSummary(prdPath string, stats pilotStats) {
	ui.Print("")
	ui.Header("Pilot Summary")

	finalPRD, _ := core.LoadAutoPRD(prdPath)
	if finalPRD == nil {
		ui.Print("  Could not load final state.")
		return
	}

	finalPRD.RecalculateProgress()
	total := finalPRD.Progress.TotalTasks
	completed := finalPRD.Progress.CompletedTasks

	ui.TableRow("Discovery iterations", fmt.Sprintf("%d", stats.discoveryCount))
	ui.TableRow("Impl iterations", fmt.Sprintf("%d", stats.implCount))
	ui.TableRow("Total iterations", fmt.Sprintf("%d", stats.discoveryCount+stats.implCount))
	ui.TableRow("Tasks generated", fmt.Sprintf("%d", total))
	ui.TableRow("Tasks completed", fmt.Sprintf("%d", completed))

	if total > 0 {
		remaining := total - completed
		if remaining == 0 {
			ui.Success("All tasks completed!")
		} else {
			ui.Info("Remaining tasks: %d", remaining)
			ui.Info("Run 'samuel auto start' to continue, or 'samuel auto status' for details.")
		}
	}
}
