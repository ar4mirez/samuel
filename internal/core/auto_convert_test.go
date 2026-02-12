package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTaskMarkdown(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantCount int
		wantErr   bool
		checks    func(t *testing.T, tasks []AutoTask)
	}{
		{
			name: "standard format with parent and children",
			content: `## High-Level Tasks

- [ ] 1.0 Database Schema & Migrations
  - [ ] 1.1 Create user schema [~2,000 tokens - Simple]
  - [ ] 1.2 Create session schema [~3,000 tokens - Medium]
- [ ] 2.0 Backend Authentication Service
  - [ ] 2.1 Implement password hashing [~5,000 tokens - Complex]
`,
			wantCount: 5,
			checks: func(t *testing.T, tasks []AutoTask) {
				// Parent task
				if tasks[0].ID != "1.0" {
					t.Errorf("expected first task ID 1.0, got %s", tasks[0].ID)
				}
				if tasks[0].ParentID != "" {
					t.Errorf("expected no parent for 1.0, got %s", tasks[0].ParentID)
				}
				// Child task
				if tasks[1].ID != "1.1" {
					t.Errorf("expected second task ID 1.1, got %s", tasks[1].ID)
				}
				if tasks[1].ParentID != "1.0" {
					t.Errorf("expected parent 1.0 for 1.1, got %s", tasks[1].ParentID)
				}
				if tasks[1].Complexity != TaskComplexitySimple {
					t.Errorf("expected simple complexity, got %s", tasks[1].Complexity)
				}
				// New parent group
				if tasks[3].ParentID != "" {
					t.Errorf("expected no parent for 2.0, got %s", tasks[3].ParentID)
				}
				if tasks[4].ParentID != "2.0" {
					t.Errorf("expected parent 2.0 for 2.1, got %s", tasks[4].ParentID)
				}
				if tasks[4].Complexity != TaskComplexityComplex {
					t.Errorf("expected complex complexity, got %s", tasks[4].Complexity)
				}
			},
		},
		{
			name: "checked items marked as completed",
			content: `- [x] 1.0 Completed task
- [ ] 2.0 Pending task
`,
			wantCount: 2,
			checks: func(t *testing.T, tasks []AutoTask) {
				if tasks[0].Status != TaskStatusCompleted {
					t.Errorf("expected completed, got %s", tasks[0].Status)
				}
				if tasks[1].Status != TaskStatusPending {
					t.Errorf("expected pending, got %s", tasks[1].Status)
				}
			},
		},
		{
			name: "uppercase X in checkbox",
			content: `- [X] 1.0 Completed task
`,
			wantCount: 1,
			checks: func(t *testing.T, tasks []AutoTask) {
				if tasks[0].Status != TaskStatusCompleted {
					t.Errorf("expected completed, got %s", tasks[0].Status)
				}
			},
		},
		{
			name: "missing complexity defaults to medium",
			content: `- [ ] 1.0 Task without complexity estimate
`,
			wantCount: 1,
			checks: func(t *testing.T, tasks []AutoTask) {
				if tasks[0].Complexity != TaskComplexityMedium {
					t.Errorf("expected medium complexity, got %s", tasks[0].Complexity)
				}
			},
		},
		{
			name: "skips non-task lines",
			content: `## Some Header

This is a description paragraph.

- [ ] 1.0 Actual task

Some more text.
`,
			wantCount: 1,
		},
		{
			name:    "no valid tasks returns error",
			content: "Just some text without any task lines",
			wantErr: true,
		},
		{
			name:    "empty content",
			content: "",
			wantErr: true,
		},
		{
			name: "task with comma in token count",
			content: `- [ ] 1.0 Large task [~15,000 tokens - Complex]
`,
			wantCount: 1,
			checks: func(t *testing.T, tasks []AutoTask) {
				if tasks[0].Title != "Large task" {
					t.Errorf("expected 'Large task', got '%s'", tasks[0].Title)
				}
				if tasks[0].Complexity != TaskComplexityComplex {
					t.Errorf("expected complex, got %s", tasks[0].Complexity)
				}
			},
		},
		{
			name: "dependencies set for child tasks",
			content: `- [ ] 1.0 Parent
  - [ ] 1.1 Child
`,
			wantCount: 2,
			checks: func(t *testing.T, tasks []AutoTask) {
				if len(tasks[1].DependsOn) != 1 || tasks[1].DependsOn[0] != "1.0" {
					t.Errorf("expected depends_on [1.0], got %v", tasks[1].DependsOn)
				}
			},
		},
		{
			name: "tab-indented child tasks",
			content: `- [ ] 1.0 Parent
	- [ ] 1.1 Tab indented child
`,
			wantCount: 2,
			checks: func(t *testing.T, tasks []AutoTask) {
				if tasks[1].ParentID != "1.0" {
					t.Errorf("expected parent 1.0, got %s", tasks[1].ParentID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := ParseTaskMarkdown(tt.content)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tasks) != tt.wantCount {
				t.Errorf("expected %d tasks, got %d", tt.wantCount, len(tasks))
			}
			if tt.checks != nil {
				tt.checks(t, tasks)
			}
		})
	}
}

func TestParseTaskLine(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		wantNil    bool
		wantID     string
		wantTitle  string
		wantStatus string
	}{
		{
			name:       "standard task line",
			line:       "- [ ] 1.0 Create database schema",
			wantID:     "1.0",
			wantTitle:  "Create database schema",
			wantStatus: TaskStatusPending,
		},
		{
			name:       "completed task",
			line:       "- [x] 2.3 Implement auth middleware",
			wantID:     "2.3",
			wantTitle:  "Implement auth middleware",
			wantStatus: TaskStatusCompleted,
		},
		{
			name:       "with complexity annotation",
			line:       "- [ ] 1.1 Create user table [~2,000 tokens - Simple]",
			wantID:     "1.1",
			wantTitle:  "Create user table",
			wantStatus: TaskStatusPending,
		},
		{
			name:    "non-task line",
			line:    "## Some Header",
			wantNil: true,
		},
		{
			name:    "empty line",
			line:    "",
			wantNil: true,
		},
		{
			name:    "bullet without checkbox",
			line:    "- Some regular list item",
			wantNil: true,
		},
		{
			name:       "indented child task",
			line:       "  - [ ] 3.1 Child task [~1,000 tokens - Simple]",
			wantID:     "3.1",
			wantTitle:  "Child task",
			wantStatus: TaskStatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, _ := parseTaskLine(tt.line)
			if tt.wantNil {
				if task != nil {
					t.Errorf("expected nil, got task %s", task.ID)
				}
				return
			}
			if task == nil {
				t.Fatal("expected task, got nil")
			}
			if task.ID != tt.wantID {
				t.Errorf("ID: expected %s, got %s", tt.wantID, task.ID)
			}
			if task.Title != tt.wantTitle {
				t.Errorf("Title: expected %s, got %s", tt.wantTitle, task.Title)
			}
			if task.Status != tt.wantStatus {
				t.Errorf("Status: expected %s, got %s", tt.wantStatus, task.Status)
			}
		})
	}
}

func TestConvertMarkdownToPRD(t *testing.T) {
	dir := t.TempDir()

	prdPath := filepath.Join(dir, "0001-prd-auth.md")
	if err := os.WriteFile(prdPath, []byte(`# User Authentication

## Introduction
This feature adds user authentication.

## Goals
1. Allow users to create accounts
`), 0644); err != nil {
		t.Fatalf("failed to write PRD file: %v", err)
	}

	tasksPath := filepath.Join(dir, "tasks-0001-prd-auth.md")
	if err := os.WriteFile(tasksPath, []byte(`## Task List

- [ ] 1.0 Database Setup
  - [ ] 1.1 Create user schema [~2,000 tokens - Simple]
  - [ ] 1.2 Run migrations [~1,000 tokens - Simple]
- [ ] 2.0 Auth Service
  - [ ] 2.1 Password hashing [~3,000 tokens - Medium]
`), 0644); err != nil {
		t.Fatalf("failed to write tasks file: %v", err)
	}

	prd, err := ConvertMarkdownToPRD(prdPath, tasksPath)
	if err != nil {
		t.Fatalf("ConvertMarkdownToPRD failed: %v", err)
	}

	if prd.Project.Name != "user-authentication" {
		t.Errorf("expected name 'user-authentication', got '%s'", prd.Project.Name)
	}
	if prd.Project.Description != "User Authentication" {
		t.Errorf("expected description 'User Authentication', got '%s'", prd.Project.Description)
	}
	if prd.Project.SourcePRD != prdPath {
		t.Errorf("expected source_prd %s, got %s", prdPath, prd.Project.SourcePRD)
	}
	if len(prd.Tasks) != 5 {
		t.Errorf("expected 5 tasks, got %d", len(prd.Tasks))
	}
	if prd.Progress.TotalTasks != 5 {
		t.Errorf("expected total_tasks 5, got %d", prd.Progress.TotalTasks)
	}
}

func TestConvertMarkdownToPRD_PRDOnly(t *testing.T) {
	dir := t.TempDir()

	prdPath := filepath.Join(dir, "0001-prd-feature.md")
	if err := os.WriteFile(prdPath, []byte("# My Feature\n\nSome description."), 0644); err != nil {
		t.Fatalf("failed to write PRD file: %v", err)
	}

	prd, err := ConvertMarkdownToPRD(prdPath, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if prd.Project.Name != "my-feature" {
		t.Errorf("expected 'my-feature', got '%s'", prd.Project.Name)
	}
	if len(prd.Tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(prd.Tasks))
	}
}

func TestConvertMarkdownToPRD_MissingPRD(t *testing.T) {
	_, err := ConvertMarkdownToPRD("/nonexistent/prd.md", "")
	if err == nil {
		t.Error("expected error for missing PRD")
	}
}

func TestConvertMarkdownToPRD_MissingTasks(t *testing.T) {
	dir := t.TempDir()
	prdPath := filepath.Join(dir, "prd.md")
	if err := os.WriteFile(prdPath, []byte("# Test"), 0644); err != nil {
		t.Fatalf("failed to write PRD file: %v", err)
	}

	_, err := ConvertMarkdownToPRD(prdPath, "/nonexistent/tasks.md")
	if err == nil {
		t.Error("expected error for missing tasks file")
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"User Authentication", "user-authentication"},
		{"My Feature (v2)", "my-feature-v2"},
		{"  Spaces  Everywhere  ", "spaces-everywhere"},
		{"UPPERCASE", "uppercase"},
		{"with--dashes", "with-dashes"},
		{"special!@#chars", "specialchars"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := slugify(tt.input)
			if got != tt.want {
				t.Errorf("slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExtractPRDMetadata(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantName string
		wantDesc string
	}{
		{
			name:     "standard H1",
			content:  "# User Authentication\n\nSome content",
			wantName: "user-authentication",
			wantDesc: "User Authentication",
		},
		{
			name:     "no heading",
			content:  "Just some content without heading",
			wantName: "unnamed-project",
			wantDesc: "Converted from PRD",
		},
		{
			name:     "heading after other content",
			content:  "---\nfrontmatter: true\n---\n\n# My Project\n\nContent",
			wantName: "my-project",
			wantDesc: "My Project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, desc := extractPRDMetadata(tt.content)
			if name != tt.wantName {
				t.Errorf("name: expected %q, got %q", tt.wantName, name)
			}
			if desc != tt.wantDesc {
				t.Errorf("desc: expected %q, got %q", tt.wantDesc, desc)
			}
		})
	}
}

func TestFindTasksFile(t *testing.T) {
	dir := t.TempDir()

	prdPath := filepath.Join(dir, "0001-prd-feature.md")
	if err := os.WriteFile(prdPath, []byte("# PRD"), 0644); err != nil {
		t.Fatalf("failed to write PRD file: %v", err)
	}

	// No tasks file yet
	if got := FindTasksFile(prdPath); got != "" {
		t.Errorf("expected empty, got %s", got)
	}

	// Create the tasks file
	tasksPath := filepath.Join(dir, "tasks-0001-prd-feature.md")
	if err := os.WriteFile(tasksPath, []byte("# Tasks"), 0644); err != nil {
		t.Fatalf("failed to write tasks file: %v", err)
	}

	if got := FindTasksFile(prdPath); got != tasksPath {
		t.Errorf("expected %s, got %s", tasksPath, got)
	}
}
