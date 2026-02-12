package commands

import (
	"github.com/spf13/cobra"
)

var autoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Autonomous AI coding loop (Ralph Wiggum methodology)",
	Long: `Manage autonomous AI coding loops using the Ralph Wiggum methodology.

The auto command enables unattended AI-driven development: an AI agent
independently selects, implements, and commits tasks from a structured
task list (prd.json), running across multiple fresh context windows.

Subcommands:
  init      Initialize autonomous loop for a project
  convert   Convert markdown PRD/tasks to prd.json
  status    Show loop progress and current state
  start     Begin or resume the autonomous loop
  task      Manage individual tasks (list, complete, skip, reset, add)

Workflow:
  1. aicof auto init --prd .claude/tasks/0001-prd-feature.md
  2. Review .claude/auto/prd.json and prompt.md
  3. aicof auto start
  4. aicof auto status  (check progress periodically)

Examples:
  aicof auto init --prd .claude/tasks/0001-prd-auth.md
  aicof auto init --ai-tool amp --max-iterations 100
  aicof auto convert .claude/tasks/0001-prd-auth.md
  aicof auto status
  aicof auto start --iterations 20
  aicof auto task list
  aicof auto task complete 1.1`,
}

var autoInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize autonomous loop for a project",
	Long: `Initialize the autonomous loop directory structure and configuration.

Creates .claude/auto/ with:
  - prd.json      Machine-readable task state
  - progress.txt  Append-only learnings journal
  - prompt.md     Iteration prompt template
  - auto.sh       Loop orchestration script

If --prd is provided, converts the PRD and associated task file to prd.json.

Examples:
  aicof auto init
  aicof auto init --prd .claude/tasks/0001-prd-auth.md
  aicof auto init --ai-tool amp --max-iterations 100`,
	RunE: runAutoInit,
}

var autoConvertCmd = &cobra.Command{
	Use:   "convert <prd-path>",
	Short: "Convert markdown PRD/tasks to prd.json",
	Long: `Convert a markdown PRD and optional task list into prd.json format.

Automatically looks for a matching tasks file using the convention:
  PRD: .claude/tasks/0001-prd-feature.md
  Tasks: .claude/tasks/tasks-0001-prd-feature.md

Examples:
  aicof auto convert .claude/tasks/0001-prd-auth.md`,
	Args: cobra.ExactArgs(1),
	RunE: runAutoConvert,
}

var autoStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show autonomous loop status",
	Long: `Display the current state of the autonomous loop including
task progress, iteration count, and recent activity.

Examples:
  aicof auto status`,
	RunE: runAutoStatus,
}

var autoStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Begin or resume the autonomous loop",
	Long: `Start the autonomous AI coding loop by executing auto.sh.

The loop will iterate until all tasks are completed or the maximum
iteration count is reached.

Examples:
  aicof auto start
  aicof auto start --iterations 20
  aicof auto start --dry-run
  aicof auto start --yes`,
	RunE: runAutoStart,
}

var autoTaskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage individual tasks in prd.json",
	Long: `Manually manage tasks in the autonomous loop task list.

Subcommands:
  list      List all tasks with status
  complete  Mark a task as completed
  skip      Mark a task as skipped
  reset     Reset a task to pending
  add       Add a new task

Examples:
  aicof auto task list
  aicof auto task complete 1.1
  aicof auto task skip 2.3
  aicof auto task reset 1.1
  aicof auto task add "3.0" "New parent task"`,
}

var autoTaskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks with status",
	RunE:  runAutoTaskList,
}

var autoTaskCompleteCmd = &cobra.Command{
	Use:   "complete <task-id>",
	Short: "Mark a task as completed",
	Args:  cobra.ExactArgs(1),
	RunE:  runAutoTaskComplete,
}

var autoTaskSkipCmd = &cobra.Command{
	Use:   "skip <task-id>",
	Short: "Mark a task as skipped",
	Args:  cobra.ExactArgs(1),
	RunE:  runAutoTaskSkip,
}

var autoTaskResetCmd = &cobra.Command{
	Use:   "reset <task-id>",
	Short: "Reset a task to pending",
	Args:  cobra.ExactArgs(1),
	RunE:  runAutoTaskReset,
}

var autoTaskAddCmd = &cobra.Command{
	Use:   "add <task-id> <title>",
	Short: "Add a new task",
	Args:  cobra.ExactArgs(2),
	RunE:  runAutoTaskAdd,
}

func init() {
	rootCmd.AddCommand(autoCmd)
	autoCmd.AddCommand(autoInitCmd)
	autoCmd.AddCommand(autoConvertCmd)
	autoCmd.AddCommand(autoStatusCmd)
	autoCmd.AddCommand(autoStartCmd)
	autoCmd.AddCommand(autoTaskCmd)
	autoTaskCmd.AddCommand(autoTaskListCmd)
	autoTaskCmd.AddCommand(autoTaskCompleteCmd)
	autoTaskCmd.AddCommand(autoTaskSkipCmd)
	autoTaskCmd.AddCommand(autoTaskResetCmd)
	autoTaskCmd.AddCommand(autoTaskAddCmd)

	// init flags
	autoInitCmd.Flags().String("prd", "", "Path to PRD markdown file to convert")
	autoInitCmd.Flags().String("ai-tool", "claude", "AI tool to use (claude, amp, cursor, codex)")
	autoInitCmd.Flags().Int("max-iterations", 50, "Maximum loop iterations")

	// start flags
	autoStartCmd.Flags().Int("iterations", 0, "Override max iterations for this run")
	autoStartCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	autoStartCmd.Flags().Bool("dry-run", false, "Show what would happen without executing")
}
