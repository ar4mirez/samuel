package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewAutoPRD(t *testing.T) {
	prd := NewAutoPRD("test-project", "A test project")

	if prd.Version != AutoSchemaVer {
		t.Errorf("expected version %s, got %s", AutoSchemaVer, prd.Version)
	}
	if prd.Project.Name != "test-project" {
		t.Errorf("expected project name 'test-project', got %s", prd.Project.Name)
	}
	if prd.Project.Description != "A test project" {
		t.Errorf("expected description 'A test project', got %s", prd.Project.Description)
	}
	if prd.Config.MaxIterations != 50 {
		t.Errorf("expected max_iterations 50, got %d", prd.Config.MaxIterations)
	}
	if prd.Config.AITool != "claude" {
		t.Errorf("expected ai_tool 'claude', got %s", prd.Config.AITool)
	}
	if prd.Progress.Status != LoopStatusNotStarted {
		t.Errorf("expected status 'not_started', got %s", prd.Progress.Status)
	}
	if len(prd.Tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(prd.Tasks))
	}
}

func TestAutoPRD_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	prdPath := filepath.Join(dir, "prd.json")

	prd := NewAutoPRD("test", "desc")
	prd.Tasks = []AutoTask{
		{ID: "1.0", Title: "Parent task", Status: TaskStatusPending, Priority: TaskPriorityHigh},
		{ID: "1.1", Title: "Sub task", Status: TaskStatusPending, ParentID: "1.0", DependsOn: []string{"1.0"}},
	}

	if err := prd.Save(prdPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := LoadAutoPRD(prdPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Project.Name != "test" {
		t.Errorf("expected name 'test', got %s", loaded.Project.Name)
	}
	if len(loaded.Tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(loaded.Tasks))
	}
	if loaded.Progress.TotalTasks != 2 {
		t.Errorf("expected total_tasks 2, got %d", loaded.Progress.TotalTasks)
	}
}

func TestAutoPRD_SaveCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	prdPath := filepath.Join(dir, "nested", "deep", "prd.json")

	prd := NewAutoPRD("test", "desc")
	if err := prd.Save(prdPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, err := os.Stat(prdPath); os.IsNotExist(err) {
		t.Error("expected prd.json to exist")
	}
}

func TestAutoPRD_SaveAtomicity(t *testing.T) {
	dir := t.TempDir()
	prdPath := filepath.Join(dir, "prd.json")

	prd := NewAutoPRD("test", "desc")
	if err := prd.Save(prdPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify no .tmp file left behind
	tmpPath := prdPath + ".tmp"
	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		t.Error("expected temp file to be cleaned up")
	}
}

func TestLoadAutoPRD_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	prdPath := filepath.Join(dir, "prd.json")

	os.WriteFile(prdPath, []byte("{invalid json"), 0644)

	_, err := LoadAutoPRD(prdPath)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadAutoPRD_NotFound(t *testing.T) {
	_, err := LoadAutoPRD("/nonexistent/prd.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestAutoPRD_GetNextTask(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []AutoTask
		wantID   string
		wantNil  bool
	}{
		{
			name:    "empty tasks",
			tasks:   []AutoTask{},
			wantNil: true,
		},
		{
			name: "single pending task",
			tasks: []AutoTask{
				{ID: "1.0", Title: "Task 1", Status: TaskStatusPending},
			},
			wantID: "1.0",
		},
		{
			name: "skips completed tasks",
			tasks: []AutoTask{
				{ID: "1.0", Title: "Done", Status: TaskStatusCompleted},
				{ID: "2.0", Title: "Pending", Status: TaskStatusPending},
			},
			wantID: "2.0",
		},
		{
			name: "respects priority ordering",
			tasks: []AutoTask{
				{ID: "1.0", Title: "Low", Status: TaskStatusPending, Priority: TaskPriorityLow},
				{ID: "2.0", Title: "Critical", Status: TaskStatusPending, Priority: TaskPriorityCritical},
				{ID: "3.0", Title: "High", Status: TaskStatusPending, Priority: TaskPriorityHigh},
			},
			wantID: "2.0",
		},
		{
			name: "same priority uses ID order",
			tasks: []AutoTask{
				{ID: "3.0", Title: "Third", Status: TaskStatusPending, Priority: TaskPriorityMedium},
				{ID: "1.0", Title: "First", Status: TaskStatusPending, Priority: TaskPriorityMedium},
				{ID: "2.0", Title: "Second", Status: TaskStatusPending, Priority: TaskPriorityMedium},
			},
			wantID: "1.0",
		},
		{
			name: "respects dependencies",
			tasks: []AutoTask{
				{ID: "1.0", Title: "First", Status: TaskStatusPending, Priority: TaskPriorityCritical},
				{ID: "2.0", Title: "Second", Status: TaskStatusPending, Priority: TaskPriorityCritical, DependsOn: []string{"1.0"}},
			},
			wantID: "1.0",
		},
		{
			name: "dependency met via completed",
			tasks: []AutoTask{
				{ID: "1.0", Title: "Done", Status: TaskStatusCompleted},
				{ID: "2.0", Title: "Ready", Status: TaskStatusPending, DependsOn: []string{"1.0"}},
			},
			wantID: "2.0",
		},
		{
			name: "dependency met via skipped",
			tasks: []AutoTask{
				{ID: "1.0", Title: "Skipped", Status: TaskStatusSkipped},
				{ID: "2.0", Title: "Ready", Status: TaskStatusPending, DependsOn: []string{"1.0"}},
			},
			wantID: "2.0",
		},
		{
			name: "all completed returns nil",
			tasks: []AutoTask{
				{ID: "1.0", Title: "Done", Status: TaskStatusCompleted},
				{ID: "2.0", Title: "Done", Status: TaskStatusCompleted},
			},
			wantNil: true,
		},
		{
			name: "blocked tasks not returned",
			tasks: []AutoTask{
				{ID: "1.0", Title: "Blocked", Status: TaskStatusBlocked},
			},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prd := NewAutoPRD("test", "desc")
			prd.Tasks = tt.tasks

			got := prd.GetNextTask()
			if tt.wantNil {
				if got != nil {
					t.Errorf("expected nil, got task %s", got.ID)
				}
				return
			}
			if got == nil {
				t.Fatal("expected task, got nil")
			}
			if got.ID != tt.wantID {
				t.Errorf("expected task %s, got %s", tt.wantID, got.ID)
			}
		})
	}
}

func TestAutoPRD_CompleteTask(t *testing.T) {
	prd := NewAutoPRD("test", "desc")
	prd.Tasks = []AutoTask{
		{ID: "1.0", Title: "Task", Status: TaskStatusPending},
	}

	if err := prd.CompleteTask("1.0", "abc123", 3); err != nil {
		t.Fatalf("CompleteTask failed: %v", err)
	}

	task := prd.findTask("1.0")
	if task.Status != TaskStatusCompleted {
		t.Errorf("expected completed, got %s", task.Status)
	}
	if task.CommitSHA != "abc123" {
		t.Errorf("expected commit abc123, got %s", task.CommitSHA)
	}
	if task.Iteration != 3 {
		t.Errorf("expected iteration 3, got %d", task.Iteration)
	}
	if task.CompletedAt == "" {
		t.Error("expected completed_at to be set")
	}
}

func TestAutoPRD_CompleteTask_NotFound(t *testing.T) {
	prd := NewAutoPRD("test", "desc")
	if err := prd.CompleteTask("nonexistent", "", 0); err == nil {
		t.Error("expected error for missing task")
	}
}

func TestAutoPRD_SkipTask(t *testing.T) {
	prd := NewAutoPRD("test", "desc")
	prd.Tasks = []AutoTask{
		{ID: "1.0", Title: "Task", Status: TaskStatusPending},
	}

	if err := prd.SkipTask("1.0"); err != nil {
		t.Fatalf("SkipTask failed: %v", err)
	}

	if prd.Tasks[0].Status != TaskStatusSkipped {
		t.Errorf("expected skipped, got %s", prd.Tasks[0].Status)
	}
}

func TestAutoPRD_ResetTask(t *testing.T) {
	prd := NewAutoPRD("test", "desc")
	prd.Tasks = []AutoTask{
		{ID: "1.0", Title: "Task", Status: TaskStatusCompleted, CommitSHA: "abc", Iteration: 2, CompletedAt: "2026-01-01"},
	}

	if err := prd.ResetTask("1.0"); err != nil {
		t.Fatalf("ResetTask failed: %v", err)
	}

	task := prd.Tasks[0]
	if task.Status != TaskStatusPending {
		t.Errorf("expected pending, got %s", task.Status)
	}
	if task.CommitSHA != "" {
		t.Error("expected commit_sha to be cleared")
	}
	if task.Iteration != 0 {
		t.Error("expected iteration to be cleared")
	}
	if task.CompletedAt != "" {
		t.Error("expected completed_at to be cleared")
	}
}

func TestAutoPRD_AddTask(t *testing.T) {
	prd := NewAutoPRD("test", "desc")

	task := AutoTask{ID: "1.0", Title: "New task"}
	if err := prd.AddTask(task); err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	if len(prd.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(prd.Tasks))
	}
	if prd.Tasks[0].Status != TaskStatusPending {
		t.Errorf("expected pending status, got %s", prd.Tasks[0].Status)
	}
}

func TestAutoPRD_AddTask_DuplicateID(t *testing.T) {
	prd := NewAutoPRD("test", "desc")
	prd.Tasks = []AutoTask{{ID: "1.0", Title: "Existing", Status: TaskStatusPending}}

	err := prd.AddTask(AutoTask{ID: "1.0", Title: "Duplicate"})
	if err == nil {
		t.Error("expected error for duplicate ID")
	}
}

func TestAutoPRD_AddTask_EmptyID(t *testing.T) {
	prd := NewAutoPRD("test", "desc")
	err := prd.AddTask(AutoTask{Title: "No ID"})
	if err == nil {
		t.Error("expected error for empty ID")
	}
}

func TestAutoPRD_RecalculateProgress(t *testing.T) {
	tests := []struct {
		name           string
		tasks          []AutoTask
		wantTotal      int
		wantCompleted  int
		wantStatus     string
	}{
		{
			name:          "no tasks",
			tasks:         []AutoTask{},
			wantTotal:     0,
			wantCompleted: 0,
			wantStatus:    LoopStatusNotStarted,
		},
		{
			name: "mixed statuses",
			tasks: []AutoTask{
				{ID: "1.0", Status: TaskStatusCompleted},
				{ID: "2.0", Status: TaskStatusPending},
				{ID: "3.0", Status: TaskStatusSkipped},
			},
			wantTotal:     3,
			wantCompleted: 1,
			wantStatus:    LoopStatusNotStarted,
		},
		{
			name: "all completed",
			tasks: []AutoTask{
				{ID: "1.0", Status: TaskStatusCompleted},
				{ID: "2.0", Status: TaskStatusCompleted},
			},
			wantTotal:     2,
			wantCompleted: 2,
			wantStatus:    LoopStatusCompleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prd := NewAutoPRD("test", "desc")
			prd.Tasks = tt.tasks
			prd.RecalculateProgress()

			if prd.Progress.TotalTasks != tt.wantTotal {
				t.Errorf("total: expected %d, got %d", tt.wantTotal, prd.Progress.TotalTasks)
			}
			if prd.Progress.CompletedTasks != tt.wantCompleted {
				t.Errorf("completed: expected %d, got %d", tt.wantCompleted, prd.Progress.CompletedTasks)
			}
			if prd.Progress.Status != tt.wantStatus {
				t.Errorf("status: expected %s, got %s", tt.wantStatus, prd.Progress.Status)
			}
		})
	}
}

func TestValidateAutoPRD(t *testing.T) {
	tests := []struct {
		name      string
		prd       *AutoPRD
		wantCount int
	}{
		{
			name:      "valid prd",
			prd:       NewAutoPRD("test", "desc"),
			wantCount: 0,
		},
		{
			name: "missing version",
			prd: &AutoPRD{
				Project: AutoProject{Name: "test"},
				Tasks:   []AutoTask{},
			},
			wantCount: 1,
		},
		{
			name: "missing project name",
			prd: &AutoPRD{
				Version: "1.0",
				Tasks:   []AutoTask{},
			},
			wantCount: 1,
		},
		{
			name: "duplicate task IDs",
			prd: &AutoPRD{
				Version: "1.0",
				Project: AutoProject{Name: "test"},
				Tasks: []AutoTask{
					{ID: "1.0", Title: "A", Status: TaskStatusPending},
					{ID: "1.0", Title: "B", Status: TaskStatusPending},
				},
			},
			wantCount: 1,
		},
		{
			name: "invalid dependency reference",
			prd: &AutoPRD{
				Version: "1.0",
				Project: AutoProject{Name: "test"},
				Tasks: []AutoTask{
					{ID: "1.0", Title: "A", Status: TaskStatusPending, DependsOn: []string{"9.9"}},
				},
			},
			wantCount: 1,
		},
		{
			name: "task missing title",
			prd: &AutoPRD{
				Version: "1.0",
				Project: AutoProject{Name: "test"},
				Tasks: []AutoTask{
					{ID: "1.0", Status: TaskStatusPending},
				},
			},
			wantCount: 1,
		},
		{
			name: "invalid status",
			prd: &AutoPRD{
				Version: "1.0",
				Project: AutoProject{Name: "test"},
				Tasks: []AutoTask{
					{ID: "1.0", Title: "A", Status: "invalid"},
				},
			},
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateAutoPRD(tt.prd)
			if len(errors) != tt.wantCount {
				t.Errorf("expected %d errors, got %d: %v", tt.wantCount, len(errors), errors)
			}
		})
	}
}

func TestAutoPRD_JSONRoundTrip(t *testing.T) {
	prd := NewAutoPRD("test-project", "A test project")
	prd.Tasks = []AutoTask{
		{
			ID:            "1.0",
			Title:         "Setup database",
			Status:        TaskStatusPending,
			Priority:      TaskPriorityCritical,
			Complexity:    TaskComplexityComplex,
			FilesToCreate: []string{"db/schema.sql"},
			Guardrails:    []string{"parameterized queries"},
		},
		{
			ID:        "1.1",
			Title:     "Create user table",
			Status:    TaskStatusPending,
			ParentID:  "1.0",
			DependsOn: []string{"1.0"},
		},
	}

	data, err := json.MarshalIndent(prd, "", "  ")
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var loaded AutoPRD
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if loaded.Tasks[0].Priority != TaskPriorityCritical {
		t.Errorf("expected critical priority, got %s", loaded.Tasks[0].Priority)
	}
	if len(loaded.Tasks[0].FilesToCreate) != 1 {
		t.Errorf("expected 1 file_to_create, got %d", len(loaded.Tasks[0].FilesToCreate))
	}
	if loaded.Tasks[1].DependsOn[0] != "1.0" {
		t.Errorf("expected depends_on [1.0], got %v", loaded.Tasks[1].DependsOn)
	}
}

func TestGetAutoPRDPath(t *testing.T) {
	got := GetAutoPRDPath("/project")
	want := filepath.Join("/project", ".claude/auto", "prd.json")
	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
}

func TestGetAutoDir(t *testing.T) {
	got := GetAutoDir("/project")
	want := filepath.Join("/project", ".claude/auto")
	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
}
