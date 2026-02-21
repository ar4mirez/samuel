package core

import "testing"

func TestPriorityRank(t *testing.T) {
	tests := []struct {
		priority string
		want     int
	}{
		{TaskPriorityCritical, 0},
		{TaskPriorityHigh, 1},
		{TaskPriorityMedium, 2},
		{TaskPriorityLow, 3},
		{"", 2},
		{"unknown", 2},
	}

	for _, tt := range tests {
		t.Run(tt.priority, func(t *testing.T) {
			got := priorityRank(tt.priority)
			if got != tt.want {
				t.Errorf("priorityRank(%q) = %d, want %d", tt.priority, got, tt.want)
			}
		})
	}
}

func TestGetAvailableTasks_UnmetDependencies(t *testing.T) {
	prd := NewAutoPRD("test", "desc")
	prd.Tasks = []AutoTask{
		{ID: "1", Title: "First", Status: TaskStatusPending},
		{ID: "2", Title: "Second", Status: TaskStatusPending, DependsOn: []string{"1"}},
		{ID: "3", Title: "Third", Status: TaskStatusPending, DependsOn: []string{"1", "2"}},
	}

	available := prd.getAvailableTasks()
	if len(available) != 1 {
		t.Fatalf("expected 1 available task, got %d", len(available))
	}
	if available[0].ID != "1" {
		t.Errorf("expected task 1, got %s", available[0].ID)
	}
}

func TestGetAvailableTasks_InProgressNotAvailable(t *testing.T) {
	prd := NewAutoPRD("test", "desc")
	prd.Tasks = []AutoTask{
		{ID: "1", Title: "In progress", Status: TaskStatusInProgress},
		{ID: "2", Title: "Blocked", Status: TaskStatusBlocked},
	}

	available := prd.getAvailableTasks()
	if len(available) != 0 {
		t.Errorf("expected 0 available tasks, got %d", len(available))
	}
}

func TestFindTask_NotFound(t *testing.T) {
	prd := NewAutoPRD("test", "desc")
	prd.Tasks = []AutoTask{
		{ID: "1", Title: "Only task", Status: TaskStatusPending},
	}

	if task := prd.findTask("999"); task != nil {
		t.Errorf("expected nil for non-existent task, got %v", task)
	}
}

func TestSkipTask_NotFound(t *testing.T) {
	prd := NewAutoPRD("test", "desc")
	if err := prd.SkipTask("nonexistent"); err == nil {
		t.Error("expected error for non-existent task")
	}
}

func TestResetTask_NotFound(t *testing.T) {
	prd := NewAutoPRD("test", "desc")
	if err := prd.ResetTask("nonexistent"); err == nil {
		t.Error("expected error for non-existent task")
	}
}

func TestAddTask_PreservesExplicitStatus(t *testing.T) {
	prd := NewAutoPRD("test", "desc")
	task := AutoTask{ID: "1", Title: "Blocked task", Status: TaskStatusBlocked}
	if err := prd.AddTask(task); err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}
	if prd.Tasks[0].Status != TaskStatusBlocked {
		t.Errorf("expected status %q preserved, got %q", TaskStatusBlocked, prd.Tasks[0].Status)
	}
}

func TestValidateTasks_EmptyID(t *testing.T) {
	tasks := []AutoTask{
		{ID: "", Title: "No ID", Status: TaskStatusPending},
	}
	errors := validateTasks(tasks)

	found := false
	for _, e := range errors {
		if e == "task missing ID" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'task missing ID' error, got %v", errors)
	}
}

func TestValidateTasks_EmptyIDSkipsFurtherChecks(t *testing.T) {
	// A task with empty ID should not produce "missing title" or "invalid status"
	tasks := []AutoTask{
		{ID: "", Title: "", Status: ""},
	}
	errors := validateTasks(tasks)
	if len(errors) != 1 {
		t.Errorf("expected exactly 1 error for empty ID task, got %d: %v", len(errors), errors)
	}
}

func TestRecalculateProgress_SkippedNotCounted(t *testing.T) {
	prd := NewAutoPRD("test", "desc")
	prd.Tasks = []AutoTask{
		{ID: "1", Status: TaskStatusCompleted},
		{ID: "2", Status: TaskStatusSkipped},
		{ID: "3", Status: TaskStatusBlocked},
	}
	prd.RecalculateProgress()

	if prd.Progress.CompletedTasks != 1 {
		t.Errorf("expected 1 completed, got %d", prd.Progress.CompletedTasks)
	}
	if prd.Progress.Status == LoopStatusCompleted {
		t.Error("should not be completed when tasks are skipped/blocked")
	}
}

func TestIsValidStatus(t *testing.T) {
	valid := []string{
		TaskStatusPending, TaskStatusInProgress,
		TaskStatusCompleted, TaskStatusSkipped, TaskStatusBlocked,
	}
	for _, s := range valid {
		if !isValidStatus(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}

	invalid := []string{"", "unknown", "done", "PENDING"}
	for _, s := range invalid {
		if isValidStatus(s) {
			t.Errorf("expected %q to be invalid", s)
		}
	}
}

func TestAllDependenciesMet(t *testing.T) {
	completed := map[string]bool{"1": true, "2": true}

	tests := []struct {
		name string
		deps []string
		want bool
	}{
		{"nil deps", nil, true},
		{"empty deps", []string{}, true},
		{"all met", []string{"1", "2"}, true},
		{"one unmet", []string{"1", "3"}, false},
		{"all unmet", []string{"3", "4"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := allDependenciesMet(tt.deps, completed)
			if got != tt.want {
				t.Errorf("allDependenciesMet(%v) = %v, want %v", tt.deps, got, tt.want)
			}
		})
	}
}
