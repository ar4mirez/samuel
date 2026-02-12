package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatProgressEntry(t *testing.T) {
	tests := []struct {
		name  string
		entry ProgressEntry
		check func(t *testing.T, result string)
	}{
		{
			name: "full entry",
			entry: ProgressEntry{
				Iteration: 3,
				TaskID:    "1.1",
				Type:      ProgressCompleted,
				Message:   "Created user schema",
			},
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "[iteration:3]") {
					t.Error("expected iteration tag")
				}
				if !strings.Contains(result, "[task:1.1]") {
					t.Error("expected task tag")
				}
				if !strings.Contains(result, "COMPLETED: Created user schema") {
					t.Error("expected type and message")
				}
			},
		},
		{
			name: "entry without task",
			entry: ProgressEntry{
				Iteration: 1,
				Type:      ProgressQualityCheck,
				Message:   "go test ./... PASSED",
			},
			check: func(t *testing.T, result string) {
				if strings.Contains(result, "[task:") {
					t.Error("should not contain task tag")
				}
				if !strings.Contains(result, "QUALITY_CHECK: go test ./... PASSED") {
					t.Error("expected quality check message")
				}
			},
		},
		{
			name: "entry without iteration",
			entry: ProgressEntry{
				Type:    ProgressLearning,
				Message: "Some insight",
			},
			check: func(t *testing.T, result string) {
				if strings.Contains(result, "[iteration:") {
					t.Error("should not contain iteration tag")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatProgressEntry(tt.entry)
			// All entries start with timestamp
			if !strings.HasPrefix(result, "[") {
				t.Error("expected timestamp prefix")
			}
			tt.check(t, result)
		})
	}
}

func TestAppendProgress(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.txt")

	entry1 := ProgressEntry{Iteration: 1, TaskID: "1.0", Type: ProgressStarted, Message: "Task 1"}
	entry2 := ProgressEntry{Iteration: 1, TaskID: "1.0", Type: ProgressCompleted, Message: "Done"}

	if err := AppendProgress(path, entry1); err != nil {
		t.Fatalf("first append failed: %v", err)
	}
	if err := AppendProgress(path, entry2); err != nil {
		t.Fatalf("second append failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestAppendProgress_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new-progress.txt")

	entry := ProgressEntry{Type: ProgressStarted, Message: "Begin"}
	if err := AppendProgress(path, entry); err != nil {
		t.Fatalf("append failed: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to be created")
	}
}

func TestReadProgressTail(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.txt")

	// Write 5 lines
	content := "line1\nline2\nline3\nline4\nline5\n"
	os.WriteFile(path, []byte(content), 0644)

	tests := []struct {
		name  string
		lines int
		want  int
	}{
		{"last 3", 3, 3},
		{"last 10 (more than exists)", 10, 5},
		{"all lines (0)", 0, 5},
		{"last 1", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadProgressTail(path, tt.lines)
			if err != nil {
				t.Fatalf("ReadProgressTail failed: %v", err)
			}
			if len(got) != tt.want {
				t.Errorf("expected %d lines, got %d", tt.want, len(got))
			}
		})
	}
}

func TestReadProgressTail_MissingFile(t *testing.T) {
	lines, err := ReadProgressTail("/nonexistent/progress.txt", 5)
	if err != nil {
		t.Errorf("expected nil error for missing file, got: %v", err)
	}
	if lines != nil {
		t.Errorf("expected nil lines, got %v", lines)
	}
}
