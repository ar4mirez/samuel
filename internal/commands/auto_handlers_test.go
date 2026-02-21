package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ar4mirez/samuel/internal/core"
)

func TestDetectQualityChecks(t *testing.T) {
	tests := []struct {
		name       string
		files      []string
		wantChecks []string
	}{
		{
			name:       "go project",
			files:      []string{"go.mod"},
			wantChecks: []string{"go test ./...", "go vet ./...", "go build ./..."},
		},
		{
			name:       "node project",
			files:      []string{"package.json"},
			wantChecks: []string{"npm test", "npm run lint", "npm run build"},
		},
		{
			name:       "rust project",
			files:      []string{"Cargo.toml"},
			wantChecks: []string{"cargo test", "cargo clippy", "cargo build"},
		},
		{
			name:       "python project",
			files:      []string{"requirements.txt"},
			wantChecks: []string{"pytest", "ruff check ."},
		},
		{
			name:       "empty directory",
			files:      []string{},
			wantChecks: []string{},
		},
		{
			name:       "go takes priority over node",
			files:      []string{"go.mod", "package.json"},
			wantChecks: []string{"go test ./...", "go vet ./...", "go build ./..."},
		},
		{
			name:       "node takes priority over rust",
			files:      []string{"package.json", "Cargo.toml"},
			wantChecks: []string{"npm test", "npm run lint", "npm run build"},
		},
		{
			name:       "rust takes priority over python",
			files:      []string{"Cargo.toml", "requirements.txt"},
			wantChecks: []string{"cargo test", "cargo clippy", "cargo build"},
		},
		{
			name:       "unrecognized project file",
			files:      []string{"Makefile"},
			wantChecks: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, f := range tt.files {
				if err := os.WriteFile(filepath.Join(dir, f), []byte(""), 0644); err != nil {
					t.Fatalf("failed to create file %s: %v", f, err)
				}
			}

			got := detectQualityChecks(dir)

			if len(got) != len(tt.wantChecks) {
				t.Fatalf("detectQualityChecks() returned %d checks, want %d\ngot:  %v\nwant: %v",
					len(got), len(tt.wantChecks), got, tt.wantChecks)
			}
			for i, check := range got {
				if check != tt.wantChecks[i] {
					t.Errorf("check[%d] = %q, want %q", i, check, tt.wantChecks[i])
				}
			}
		})
	}
}

func TestCountTaskStatuses(t *testing.T) {
	tests := []struct {
		name  string
		tasks []core.AutoTask
		want  map[string]int
	}{
		{
			name:  "empty tasks",
			tasks: []core.AutoTask{},
			want: map[string]int{
				"pending": 0, "in_progress": 0, "completed": 0,
				"skipped": 0, "blocked": 0,
			},
		},
		{
			name: "all pending",
			tasks: []core.AutoTask{
				{ID: "1", Status: "pending"},
				{ID: "2", Status: "pending"},
				{ID: "3", Status: "pending"},
			},
			want: map[string]int{
				"pending": 3, "in_progress": 0, "completed": 0,
				"skipped": 0, "blocked": 0,
			},
		},
		{
			name: "all completed",
			tasks: []core.AutoTask{
				{ID: "1", Status: "completed"},
				{ID: "2", Status: "completed"},
			},
			want: map[string]int{
				"pending": 0, "in_progress": 0, "completed": 2,
				"skipped": 0, "blocked": 0,
			},
		},
		{
			name: "mixed statuses",
			tasks: []core.AutoTask{
				{ID: "1", Status: "completed"},
				{ID: "2", Status: "in_progress"},
				{ID: "3", Status: "pending"},
				{ID: "4", Status: "blocked"},
				{ID: "5", Status: "skipped"},
				{ID: "6", Status: "pending"},
			},
			want: map[string]int{
				"pending": 2, "in_progress": 1, "completed": 1,
				"skipped": 1, "blocked": 1,
			},
		},
		{
			name: "unknown status counted separately",
			tasks: []core.AutoTask{
				{ID: "1", Status: "completed"},
				{ID: "2", Status: "unknown"},
			},
			want: map[string]int{
				"pending": 0, "in_progress": 0, "completed": 1,
				"skipped": 0, "blocked": 0, "unknown": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prd := &core.AutoPRD{Tasks: tt.tasks}
			got := countTaskStatuses(prd)

			for status, wantCount := range tt.want {
				if got[status] != wantCount {
					t.Errorf("count[%q] = %d, want %d", status, got[status], wantCount)
				}
			}
			// Verify no unexpected keys
			for status, gotCount := range got {
				if _, exists := tt.want[status]; !exists {
					t.Errorf("unexpected status %q with count %d", status, gotCount)
				}
			}
		})
	}
}

func TestValidateSandbox(t *testing.T) {
	t.Run("none mode always succeeds", func(t *testing.T) {
		err := validateSandbox("none")
		if err != nil {
			t.Errorf("validateSandbox(\"none\") returned error: %v", err)
		}
	})

	t.Run("empty string passes", func(t *testing.T) {
		err := validateSandbox("")
		if err != nil {
			t.Errorf("validateSandbox(\"\") returned error: %v", err)
		}
	})

	t.Run("unrecognized mode passes", func(t *testing.T) {
		// Non-docker modes skip both docker checks entirely
		err := validateSandbox("local")
		if err != nil {
			t.Errorf("validateSandbox(\"local\") returned error: %v", err)
		}
	})

	t.Run("docker mode checks docker availability", func(t *testing.T) {
		err := validateSandbox(core.SandboxDocker)
		// Result depends on docker availability in environment;
		// just verify it doesn't panic and returns a valid result type
		if err != nil {
			// Expected when docker is not available
			t.Logf("docker unavailable (expected in CI): %v", err)
		}
	})

	t.Run("docker-sandbox mode checks docker sandbox availability", func(t *testing.T) {
		err := validateSandbox(core.SandboxDockerSandbox)
		// Result depends on docker sandbox availability in environment
		if err != nil {
			t.Logf("docker sandbox unavailable (expected in CI): %v", err)
		}
	})
}
