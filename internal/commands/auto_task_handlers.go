package commands

import (
	"fmt"
	"os"

	"github.com/ar4mirez/samuel/internal/core"
	"github.com/ar4mirez/samuel/internal/ui"
	"github.com/spf13/cobra"
)

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
