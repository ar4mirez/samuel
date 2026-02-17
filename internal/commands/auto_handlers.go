package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ar4mirez/samuel/internal/core"
	"github.com/ar4mirez/samuel/internal/ui"
	"github.com/spf13/cobra"
)

func runAutoInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	if !core.ConfigExists(cwd) {
		return fmt.Errorf("no Samuel installation found. Run 'samuel init' first")
	}

	aiTool, _ := cmd.Flags().GetString("ai-tool")
	maxIter, _ := cmd.Flags().GetInt("max-iterations")
	prdPath, _ := cmd.Flags().GetString("prd")
	sandbox, _ := cmd.Flags().GetString("sandbox")
	sandboxImage, _ := cmd.Flags().GetString("sandbox-image")
	sandboxTemplate, _ := cmd.Flags().GetString("sandbox-template")

	if !core.IsValidAITool(aiTool) {
		return fmt.Errorf("unsupported AI tool: %s (supported: %v)", aiTool, core.GetSupportedAITools())
	}

	if !core.IsValidSandboxMode(sandbox) {
		return fmt.Errorf("unsupported sandbox mode: %s (supported: %v)", sandbox, core.GetSupportedSandboxModes())
	}

	return initAutoDir(cwd, prdPath, aiTool, maxIter, sandbox, sandboxImage, sandboxTemplate)
}

func initAutoDir(cwd, prdPath, aiTool string, maxIter int, sandbox, sandboxImage, sandboxTemplate string) error {
	autoDir := core.GetAutoDir(cwd)
	if err := os.MkdirAll(autoDir, 0755); err != nil {
		return fmt.Errorf("failed to create auto directory: %w", err)
	}

	config := core.AutoConfig{
		MaxIterations:   maxIter,
		QualityChecks:   detectQualityChecks(cwd),
		AITool:          aiTool,
		PromptFile:      filepath.Join(core.AutoDir, core.AutoPromptFile),
		Sandbox:         sandbox,
		SandboxImage:    sandboxImage,
		SandboxTemplate: sandboxTemplate,
	}

	if err := writeAutoFiles(autoDir, config); err != nil {
		return err
	}

	if prdPath != "" {
		if err := convertAndSavePRD(cwd, prdPath); err != nil {
			return err
		}
	} else {
		prd := core.NewAutoPRD("my-project", "Autonomous loop project")
		prd.Config = config
		prdFile := filepath.Join(autoDir, core.AutoPRDFile)
		if err := prd.Save(prdFile); err != nil {
			return fmt.Errorf("failed to save prd.json: %w", err)
		}
	}

	printInitSummary(autoDir, prdPath)
	return nil
}

func writeAutoFiles(autoDir string, config core.AutoConfig) error {
	promptContent := core.GeneratePromptFile(config)
	promptPath := filepath.Join(autoDir, core.AutoPromptFile)
	if err := os.WriteFile(promptPath, []byte(promptContent), 0644); err != nil {
		return fmt.Errorf("failed to write prompt.md: %w", err)
	}

	scriptContent := core.GenerateAutoScript(config)
	scriptPath := filepath.Join(autoDir, core.AutoScriptFile)
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("failed to write auto.sh: %w", err)
	}

	progressPath := filepath.Join(autoDir, core.AutoProgressFile)
	if _, err := os.Stat(progressPath); os.IsNotExist(err) {
		if err := os.WriteFile(progressPath, []byte(""), 0644); err != nil {
			return fmt.Errorf("failed to create progress.md: %w", err)
		}
	}

	return nil
}

func printInitSummary(autoDir, prdPath string) {
	ui.Success("Auto loop initialized at %s/", autoDir)
	ui.Print("")
	ui.Print("  Files created:")
	ui.Print("    %s", filepath.Join(core.AutoDir, core.AutoPRDFile))
	ui.Print("    %s", filepath.Join(core.AutoDir, core.AutoProgressFile))
	ui.Print("    %s", filepath.Join(core.AutoDir, core.AutoPromptFile))
	ui.Print("    %s", filepath.Join(core.AutoDir, core.AutoScriptFile))
	ui.Print("")

	if prdPath != "" {
		ui.Info("PRD converted from: %s", prdPath)
	} else {
		ui.Info("No PRD provided. Add tasks with 'samuel auto task add'")
		ui.Info("Or convert a PRD with 'samuel auto convert <prd-path>'")
	}

	ui.Print("")
	ui.Info("Next steps:")
	ui.Print("  1. Review %s", filepath.Join(core.AutoDir, core.AutoPRDFile))
	ui.Print("  2. Review %s", filepath.Join(core.AutoDir, core.AutoPromptFile))
	ui.Print("  3. Run 'samuel auto start' to begin the loop")
}

func detectQualityChecks(cwd string) []string {
	checks := []string{}
	if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
		return []string{"go test ./...", "go vet ./...", "go build ./..."}
	}
	if _, err := os.Stat(filepath.Join(cwd, "package.json")); err == nil {
		return []string{"npm test", "npm run lint", "npm run build"}
	}
	if _, err := os.Stat(filepath.Join(cwd, "Cargo.toml")); err == nil {
		return []string{"cargo test", "cargo clippy", "cargo build"}
	}
	if _, err := os.Stat(filepath.Join(cwd, "requirements.txt")); err == nil {
		return []string{"pytest", "ruff check ."}
	}
	return checks
}

func runAutoConvert(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	return convertAndSavePRD(cwd, args[0])
}

func convertAndSavePRD(cwd, prdPath string) error {
	tasksPath := core.FindTasksFile(prdPath)

	spinner := ui.NewSpinner("Converting PRD to prd.json")
	spinner.Start()

	prd, err := core.ConvertMarkdownToPRD(prdPath, tasksPath)
	if err != nil {
		spinner.Error("Conversion failed")
		return fmt.Errorf("failed to convert PRD: %w", err)
	}

	prdFile := core.GetAutoPRDPath(cwd)
	if err := prd.Save(prdFile); err != nil {
		spinner.Error("Save failed")
		return fmt.Errorf("failed to save prd.json: %w", err)
	}

	spinner.Success("Converted successfully")
	ui.Print("")
	ui.Print("  Project: %s", prd.Project.Name)
	ui.Print("  Tasks:   %d", prd.Progress.TotalTasks)
	if tasksPath != "" {
		ui.Print("  Source:  %s + %s", prdPath, tasksPath)
	} else {
		ui.Print("  Source:  %s (no task file found)", prdPath)
	}
	ui.Print("  Output:  %s", prdFile)
	return nil
}

func runAutoStatus(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	prdPath := core.GetAutoPRDPath(cwd)
	prd, err := core.LoadAutoPRD(prdPath)
	if err != nil {
		return fmt.Errorf("no auto loop found. Run 'samuel auto init' first")
	}

	prd.RecalculateProgress()
	printStatus(prd)
	return nil
}

func printStatus(prd *core.AutoPRD) {
	ui.Header("Auto Loop Status")

	ui.TableRow("Project", prd.Project.Name)
	ui.TableRow("Status", prd.Progress.Status)

	pct := 0
	if prd.Progress.TotalTasks > 0 {
		pct = (prd.Progress.CompletedTasks * 100) / prd.Progress.TotalTasks
	}
	ui.TableRow("Progress", fmt.Sprintf("%d/%d tasks (%d%%)",
		prd.Progress.CompletedTasks, prd.Progress.TotalTasks, pct))
	ui.TableRow("AI Tool", prd.Config.AITool)
	ui.TableRow("Sandbox", prd.Config.Sandbox)
	if prd.Config.Sandbox == core.SandboxDocker && prd.Config.SandboxImage != "" {
		ui.TableRow("Sandbox Image", prd.Config.SandboxImage)
	}
	if prd.Config.Sandbox == core.SandboxDockerSandbox && prd.Config.SandboxTemplate != "" {
		ui.TableRow("Sandbox Template", prd.Config.SandboxTemplate)
	}
	ui.TableRow("Max Iterations", fmt.Sprintf("%d", prd.Config.MaxIterations))

	if prd.Progress.TotalIterationsRun > 0 {
		ui.TableRow("Iterations Run", fmt.Sprintf("%d", prd.Progress.TotalIterationsRun))
	}
	if prd.Progress.LastIterationAt != "" {
		ui.TableRow("Last Iteration", prd.Progress.LastIterationAt)
	}

	// Count by status
	counts := countTaskStatuses(prd)
	ui.Print("")
	ui.Print("  Pending: %d  Completed: %d  Blocked: %d  Skipped: %d",
		counts["pending"], counts["completed"], counts["blocked"], counts["skipped"])

	next := prd.GetNextTask()
	if next != nil {
		ui.Print("")
		ui.Info("Next task: %s %s", next.ID, next.Title)
	}
}

func countTaskStatuses(prd *core.AutoPRD) map[string]int {
	counts := map[string]int{
		"pending": 0, "in_progress": 0, "completed": 0, "skipped": 0, "blocked": 0,
	}
	for _, t := range prd.Tasks {
		counts[t.Status]++
	}
	return counts
}

func runAutoStart(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	scriptPath := filepath.Join(core.GetAutoDir(cwd), core.AutoScriptFile)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("auto.sh not found. Run 'samuel auto init' first")
	}

	prdPath := core.GetAutoPRDPath(cwd)
	prd, err := core.LoadAutoPRD(prdPath)
	if err != nil {
		return fmt.Errorf("failed to load prd.json: %w", err)
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

	if sandbox == core.SandboxDocker {
		if err := core.CheckDockerAvailable(); err != nil {
			return fmt.Errorf("docker sandbox unavailable: %w", err)
		}
	}
	if sandbox == core.SandboxDockerSandbox {
		if err := core.CheckDockerSandboxAvailable(); err != nil {
			return fmt.Errorf("docker sandbox unavailable: %w", err)
		}
	}

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		return printDryRun(prd, cwd, scriptPath, sandbox, sandboxImage, sandboxTemplate)
	}

	skipConfirm, _ := cmd.Flags().GetBool("yes")
	if !skipConfirm {
		confirmed, err := ui.Confirm("Start autonomous loop?", false)
		if err != nil || !confirmed {
			ui.Info("Cancelled")
			return nil
		}
	}

	iterOverride, _ := cmd.Flags().GetInt("iterations")

	switch sandbox {
	case core.SandboxDocker:
		return executeAutoScriptInDocker(cwd, scriptPath, iterOverride, sandboxImage)
	case core.SandboxDockerSandbox:
		return executeAutoInDockerSandbox(cwd, prd, iterOverride, sandboxTemplate)
	default:
		return executeAutoScript(scriptPath, iterOverride)
	}
}

func printDryRun(prd *core.AutoPRD, cwd, scriptPath, sandbox, sandboxImage, sandboxTemplate string) error {
	ui.Header("Dry Run - Auto Loop")
	ui.Print("  Script:     %s", scriptPath)
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

func executeAutoScript(scriptPath string, iterOverride int) error {
	args := []string{scriptPath}
	if iterOverride > 0 {
		args = append(args, fmt.Sprintf("%d", iterOverride))
	}

	execCmd := exec.Command("bash", args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	ui.Info("Starting auto loop...")
	ui.Print("")

	if err := execCmd.Run(); err != nil {
		return fmt.Errorf("auto loop exited with error: %w", err)
	}
	return nil
}

func executeAutoScriptInDocker(cwd, scriptPath string, iterOverride int, image string) error {
	relScript, err := filepath.Rel(cwd, scriptPath)
	if err != nil {
		return fmt.Errorf("failed to compute relative script path: %w", err)
	}

	dockerConfig := core.DockerSandboxConfig{
		Image:        image,
		WorkDir:      cwd,
		ScriptPath:   relScript,
		IterOverride: iterOverride,
	}

	args := core.BuildDockerArgs(dockerConfig)
	execCmd := exec.Command("docker", args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	resolvedImage := image
	if resolvedImage == "" {
		resolvedImage = core.DefaultSandboxImage
	}

	ui.Info("Starting auto loop in Docker sandbox...")
	ui.Print("  Image: %s", resolvedImage)
	ui.Print("  Mount: %s -> %s", cwd, core.DockerContainerMount)
	ui.Print("")

	if err := execCmd.Run(); err != nil {
		return fmt.Errorf("docker auto loop exited with error: %w", err)
	}
	return nil
}

func executeAutoInDockerSandbox(
	cwd string,
	prd *core.AutoPRD,
	iterOverride int,
	template string,
) error {
	maxIter := prd.Config.MaxIterations
	if iterOverride > 0 {
		maxIter = iterOverride
	}

	prdPath := core.GetAutoPRDPath(cwd)
	promptPath := filepath.Join(cwd, prd.Config.PromptFile)

	ui.Info("Starting auto loop in Docker Sandbox (microVM)...")
	ui.Print("  Agent:     %s", prd.Config.AITool)
	ui.Print("  Workspace: %s", cwd)
	if template != "" {
		ui.Print("  Template:  %s", template)
	}
	ui.Print("")

	pauseSeconds := 2
	if val := os.Getenv("PAUSE_SECONDS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			pauseSeconds = parsed
		}
	}

	var i int
	for i = 1; i <= maxIter; i++ {
		currentPRD, err := core.LoadAutoPRD(prdPath)
		if err != nil {
			return fmt.Errorf("iteration %d: failed to reload prd.json: %w", i, err)
		}

		if currentPRD.GetNextTask() == nil {
			ui.Success("All tasks completed after %d iterations!", i-1)
			return nil
		}

		ui.Info("[iteration:%d] Starting iteration %d of %d", i, i, maxIter)

		agentArgs, err := core.GetAgentArgs(prd.Config.AITool, promptPath)
		if err != nil {
			return fmt.Errorf("iteration %d: %w", i, err)
		}
		config := core.DockerSandboxRunConfig{
			Agent:     prd.Config.AITool,
			WorkDir:   cwd,
			Template:  template,
			AgentArgs: agentArgs,
		}

		args := core.BuildDockerSandboxArgs(config)
		execCmd := exec.Command("docker", args...)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		execCmd.Stdin = os.Stdin

		if err := execCmd.Run(); err != nil {
			ui.Warn("[iteration:%d] Agent exited with error: %v. Continuing...", i, err)
		}

		time.Sleep(time.Duration(pauseSeconds) * time.Second)
		ui.Info("[iteration:%d] Iteration %d complete.", i, i)
	}

	// Final summary
	finalPRD, _ := core.LoadAutoPRD(prdPath)
	if finalPRD != nil {
		finalPRD.RecalculateProgress()
		remaining := finalPRD.Progress.TotalTasks - finalPRD.Progress.CompletedTasks
		ui.Info("Loop finished after %d iterations. Remaining tasks: %d", i-1, remaining)
	}

	return nil
}

func runAutoTaskList(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	prd, err := core.LoadAutoPRD(core.GetAutoPRDPath(cwd))
	if err != nil {
		return fmt.Errorf("no auto loop found. Run 'samuel auto init' first")
	}

	ui.Header("Tasks")
	for _, t := range prd.Tasks {
		icon := taskStatusIcon(t.Status)
		indent := 0
		if t.ParentID != "" {
			indent = 1
		}
		ui.ListItem(indent, "%s %s %s", icon, t.ID, t.Title)
	}

	ui.Print("")
	prd.RecalculateProgress()
	ui.Print("Total: %d  Completed: %d  Pending: %d",
		prd.Progress.TotalTasks, prd.Progress.CompletedTasks,
		prd.Progress.TotalTasks-prd.Progress.CompletedTasks)
	return nil
}

func taskStatusIcon(status string) string {
	switch status {
	case core.TaskStatusCompleted:
		return "[x]"
	case core.TaskStatusSkipped:
		return "[-]"
	case core.TaskStatusBlocked:
		return "[!]"
	case core.TaskStatusInProgress:
		return "[>]"
	default:
		return "[ ]"
	}
}

func runAutoTaskComplete(cmd *cobra.Command, args []string) error {
	return updateTaskStatus(args[0], func(prd *core.AutoPRD, id string) error {
		return prd.CompleteTask(id, "", 0)
	}, "completed")
}

func runAutoTaskSkip(cmd *cobra.Command, args []string) error {
	return updateTaskStatus(args[0], func(prd *core.AutoPRD, id string) error {
		return prd.SkipTask(id)
	}, "skipped")
}

func runAutoTaskReset(cmd *cobra.Command, args []string) error {
	return updateTaskStatus(args[0], func(prd *core.AutoPRD, id string) error {
		return prd.ResetTask(id)
	}, "reset to pending")
}

func updateTaskStatus(id string, fn func(*core.AutoPRD, string) error, label string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	prdPath := core.GetAutoPRDPath(cwd)
	prd, err := core.LoadAutoPRD(prdPath)
	if err != nil {
		return fmt.Errorf("no auto loop found. Run 'samuel auto init' first")
	}

	if err := fn(prd, id); err != nil {
		return err
	}

	if err := prd.Save(prdPath); err != nil {
		return fmt.Errorf("failed to save prd.json: %w", err)
	}

	ui.Success("Task %s %s", id, label)
	return nil
}

func runAutoTaskAdd(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	prdPath := core.GetAutoPRDPath(cwd)
	prd, err := core.LoadAutoPRD(prdPath)
	if err != nil {
		return fmt.Errorf("no auto loop found. Run 'samuel auto init' first")
	}

	task := core.AutoTask{
		ID:       args[0],
		Title:    args[1],
		Status:   core.TaskStatusPending,
		Priority: core.TaskPriorityMedium,
	}

	if err := prd.AddTask(task); err != nil {
		return err
	}

	if err := prd.Save(prdPath); err != nil {
		return fmt.Errorf("failed to save prd.json: %w", err)
	}

	ui.Success("Task %s added: %s", task.ID, task.Title)
	return nil
}
