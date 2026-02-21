package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateRemovePath(t *testing.T) {
	tests := []struct {
		name          string
		projectDir    string
		componentPath string
		wantErr       bool
		wantPath      string
	}{
		{
			name:          "valid_relative_path",
			projectDir:    "/project",
			componentPath: ".claude/skills/go-guide",
			wantErr:       false,
			wantPath:      filepath.Join("/project", ".claude/skills/go-guide"),
		},
		{
			name:          "valid_nested_path",
			projectDir:    "/project",
			componentPath: ".claude/skills/typescript-guide",
			wantErr:       false,
			wantPath:      filepath.Join("/project", ".claude/skills/typescript-guide"),
		},
		{
			name:          "traversal_parent_escape",
			projectDir:    "/project",
			componentPath: "../../etc/passwd",
			wantErr:       true,
		},
		{
			name:          "traversal_deep_escape",
			projectDir:    "/project",
			componentPath: ".claude/../../..",
			wantErr:       true,
		},
		{
			name:          "traversal_dot_dot_only",
			projectDir:    "/project",
			componentPath: "..",
			wantErr:       true,
		},
		{
			name:          "current_dir_resolves_to_base",
			projectDir:    "/project",
			componentPath: ".",
			wantErr:       false,
			wantPath:      "/project",
		},
		{
			name:          "valid_single_file",
			projectDir:    "/project",
			componentPath: "CLAUDE.md",
			wantErr:       false,
			wantPath:      filepath.Join("/project", "CLAUDE.md"),
		},
		{
			name:          "path_with_redundant_separators",
			projectDir:    "/project",
			componentPath: ".claude///skills//go-guide",
			wantErr:       false,
			wantPath:      filepath.Join("/project", ".claude/skills/go-guide"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateRemovePath(tt.projectDir, tt.componentPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRemovePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.wantPath {
				t.Errorf("validateRemovePath() = %q, want %q", got, tt.wantPath)
			}
			if tt.wantErr && got != "" {
				t.Errorf("validateRemovePath() returned non-empty path %q on error", got)
			}
		})
	}
}

func TestValidateRemovePath_WithTempDir(t *testing.T) {
	projectDir := t.TempDir()

	t.Run("valid_path_within_project", func(t *testing.T) {
		subDir := filepath.Join(projectDir, ".claude", "skills")
		if err := os.MkdirAll(subDir, 0755); err != nil {
			t.Fatal(err)
		}

		got, err := validateRemovePath(projectDir, ".claude/skills/test-guide")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := filepath.Join(projectDir, ".claude/skills/test-guide")
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("traversal_rejected", func(t *testing.T) {
		_, err := validateRemovePath(projectDir, "../../etc/important_file")
		if err == nil {
			t.Fatal("expected error for path traversal, got nil")
		}
		if !strings.Contains(err.Error(), "path traversal detected") {
			t.Errorf("error should mention path traversal, got: %v", err)
		}
	})
}
