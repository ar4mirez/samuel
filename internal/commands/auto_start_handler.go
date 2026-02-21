package commands

import (
	"fmt"
	"os"

	"github.com/ar4mirez/samuel/internal/core"
	"github.com/ar4mirez/samuel/internal/ui"
	"github.com/spf13/cobra"
)

func runAutoStart(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	prdPath := core.GetAutoPRDPath(cwd)
	prd, err := core.LoadAutoPRD(prdPath)
	if err != nil {
		return fmt.Errorf("failed to load prd.json. Run 'samuel auto init' first: %w", err)
	}

	// Resolve sandbox mode: CLI flag overrides prd.json config
	sandbox := prd.Config.Sandbox
	if flagSandbox, _ := cmd.Flags().GetString("sandbox"); flagSandbox != "" {
		sandbox = flagSandbox
	}
	sandboxImage := prd.Config.SandboxImage
	if flagImage, _ := cmd.Flags().GetString("sandbox-image"); flagImage != "" {
		sandboxImage = flagImage
	}
	sandboxTemplate := prd.Config.SandboxTemplate
	if flagTpl, _ := cmd.Flags().GetString("sandbox-template"); flagTpl != "" {
		sandboxTemplate = flagTpl
	}

	if !core.IsValidSandboxMode(sandbox) {
		return fmt.Errorf("unsupported sandbox mode: %s (supported: %v)", sandbox, core.GetSupportedSandboxModes())
	}

	if err := validateSandbox(sandbox); err != nil {
		return err
	}

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		return printStartDryRun(prd, cwd, sandbox, sandboxImage, sandboxTemplate)
	}

	skipConfirm, _ := cmd.Flags().GetBool("yes")
	if !skipConfirm {
		confirmed, confirmErr := ui.Confirm("Start autonomous loop?", false)
		if confirmErr != nil || !confirmed {
			ui.Info("Cancelled")
			return nil
		}
	}

	cfg := core.NewLoopConfig(cwd, prd)
	cfg.Sandbox = sandbox
	cfg.SandboxImage = sandboxImage
	cfg.SandboxTpl = sandboxTemplate

	if iterOverride, _ := cmd.Flags().GetInt("iterations"); iterOverride > 0 {
		cfg.MaxIterations = iterOverride
	}

	cfg.OnIterStart = func(iter int, iterType string) {
		ui.Info("[iteration:%d] Starting iteration %d of %d", iter, iter, cfg.MaxIterations)
	}
	cfg.OnIterEnd = func(iter int, err error) {
		if err != nil {
			ui.Warn("[iteration:%d] Agent exited with error: %v", iter, err)
		} else {
			ui.Info("[iteration:%d] Iteration %d complete.", iter, iter)
		}
	}

	ui.Info("Starting auto loop...")
	ui.Print("  AI Tool:  %s", cfg.AITool)
	ui.Print("  Sandbox:  %s", sandbox)
	ui.Print("")

	if err := core.RunAutoLoop(cfg); err != nil {
		return fmt.Errorf("auto loop exited with error: %w", err)
	}

	printLoopSummary(prdPath)
	return nil
}

func printStartDryRun(prd *core.AutoPRD, cwd, sandbox, sandboxImage, sandboxTemplate string) error {
	ui.Header("Dry Run - Auto Loop")
	ui.Print("  AI Tool:    %s", prd.Config.AITool)
	ui.Print("  Iterations: %d", prd.Config.MaxIterations)
	ui.Print("  Sandbox:    %s", sandbox)
	if sandbox == core.SandboxDocker {
		image := sandboxImage
		if image == "" {
			image = core.DefaultSandboxImage
		}
		ui.Print("  Image:      %s", image)
	}
	if sandbox == core.SandboxDockerSandbox {
		ui.Print("  Workspace:  %s (same path inside VM)", cwd)
		if sandboxTemplate != "" {
			ui.Print("  Template:   %s", sandboxTemplate)
		}
		ui.Print("  Note:       API keys read from shell config (~/.bashrc, ~/.zshrc)")
	}
	ui.Print("  Tasks:      %d pending", countTaskStatuses(prd)["pending"])
	ui.Print("")
	ui.Print("  Quality checks:")
	for _, check := range prd.Config.QualityChecks {
		ui.Print("    - %s", check)
	}
	ui.Print("")
	ui.Info("Run without --dry-run to execute")
	return nil
}

func printLoopSummary(prdPath string) {
	finalPRD, err := core.LoadAutoPRD(prdPath)
	if err != nil {
		ui.Warn("Failed to load prd.json for summary: %v", err)
	}
	if finalPRD == nil {
		return
	}
	finalPRD.RecalculateProgress()
	remaining := finalPRD.Progress.TotalTasks - finalPRD.Progress.CompletedTasks
	if remaining == 0 {
		ui.Success("All tasks completed!")
	} else {
		ui.Info("Loop finished. Remaining tasks: %d", remaining)
		ui.Info("Run 'samuel auto status' for details.")
	}
}
