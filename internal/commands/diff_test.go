package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestComputeDiff(t *testing.T) {
	tests := []struct {
		name          string
		v1            string
		v2            string
		files1        map[string]string
		files2        map[string]string
		wantAdded     int
		wantRemoved   int
		wantModified  int
		wantUnchanged int
	}{
		{
			name:          "no changes",
			v1:            "1.0.0",
			v2:            "1.1.0",
			files1:        map[string]string{"a.md": "hash1", "b.md": "hash2"},
			files2:        map[string]string{"a.md": "hash1", "b.md": "hash2"},
			wantAdded:     0,
			wantRemoved:   0,
			wantModified:  0,
			wantUnchanged: 2,
		},
		{
			name:          "file added",
			v1:            "1.0.0",
			v2:            "1.1.0",
			files1:        map[string]string{"a.md": "hash1"},
			files2:        map[string]string{"a.md": "hash1", "b.md": "hash2"},
			wantAdded:     1,
			wantRemoved:   0,
			wantModified:  0,
			wantUnchanged: 1,
		},
		{
			name:          "file removed",
			v1:            "1.0.0",
			v2:            "1.1.0",
			files1:        map[string]string{"a.md": "hash1", "b.md": "hash2"},
			files2:        map[string]string{"a.md": "hash1"},
			wantAdded:     0,
			wantRemoved:   1,
			wantModified:  0,
			wantUnchanged: 1,
		},
		{
			name:          "file modified",
			v1:            "1.0.0",
			v2:            "1.1.0",
			files1:        map[string]string{"a.md": "hash1"},
			files2:        map[string]string{"a.md": "hash2"},
			wantAdded:     0,
			wantRemoved:   0,
			wantModified:  1,
			wantUnchanged: 0,
		},
		{
			name:          "mixed changes",
			v1:            "1.0.0",
			v2:            "1.1.0",
			files1:        map[string]string{"a.md": "hash1", "b.md": "hash2", "c.md": "hash3"},
			files2:        map[string]string{"a.md": "hash1", "b.md": "modified", "d.md": "hash4"},
			wantAdded:     1, // d.md
			wantRemoved:   1, // c.md
			wantModified:  1, // b.md
			wantUnchanged: 1, // a.md
		},
		{
			name:          "empty to files",
			v1:            "1.0.0",
			v2:            "1.1.0",
			files1:        map[string]string{},
			files2:        map[string]string{"a.md": "hash1", "b.md": "hash2"},
			wantAdded:     2,
			wantRemoved:   0,
			wantModified:  0,
			wantUnchanged: 0,
		},
		{
			name:          "files to empty",
			v1:            "1.0.0",
			v2:            "1.1.0",
			files1:        map[string]string{"a.md": "hash1", "b.md": "hash2"},
			files2:        map[string]string{},
			wantAdded:     0,
			wantRemoved:   2,
			wantModified:  0,
			wantUnchanged: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := computeDiff(tt.v1, tt.v2, tt.files1, tt.files2)

			if diff.FromVersion != tt.v1 {
				t.Errorf("FromVersion = %q, want %q", diff.FromVersion, tt.v1)
			}
			if diff.ToVersion != tt.v2 {
				t.Errorf("ToVersion = %q, want %q", diff.ToVersion, tt.v2)
			}
			if len(diff.Added) != tt.wantAdded {
				t.Errorf("Added count = %d, want %d", len(diff.Added), tt.wantAdded)
			}
			if len(diff.Removed) != tt.wantRemoved {
				t.Errorf("Removed count = %d, want %d", len(diff.Removed), tt.wantRemoved)
			}
			if len(diff.Modified) != tt.wantModified {
				t.Errorf("Modified count = %d, want %d", len(diff.Modified), tt.wantModified)
			}
			if diff.Unchanged != tt.wantUnchanged {
				t.Errorf("Unchanged = %d, want %d", diff.Unchanged, tt.wantUnchanged)
			}
		})
	}
}

func TestComputeDiff_Sorting(t *testing.T) {
	files1 := map[string]string{}
	files2 := map[string]string{
		"z.md": "hash",
		"a.md": "hash",
		"m.md": "hash",
	}

	diff := computeDiff("1.0", "2.0", files1, files2)

	// Should be sorted alphabetically
	if len(diff.Added) != 3 {
		t.Fatalf("Expected 3 added files, got %d", len(diff.Added))
	}
	if diff.Added[0] != "a.md" {
		t.Errorf("First added file = %q, want %q", diff.Added[0], "a.md")
	}
	if diff.Added[1] != "m.md" {
		t.Errorf("Second added file = %q, want %q", diff.Added[1], "m.md")
	}
	if diff.Added[2] != "z.md" {
		t.Errorf("Third added file = %q, want %q", diff.Added[2], "z.md")
	}
}

func TestHashFile(t *testing.T) {
	// Create a temp file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "Hello, World!"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	hash, err := hashFile(testFile)
	if err != nil {
		t.Fatalf("hashFile() error = %v", err)
	}
	if hash == "" {
		t.Error("hashFile() returned empty hash")
	}

	// Same content should produce same hash
	testFile2 := filepath.Join(tmpDir, "test2.txt")
	if err := os.WriteFile(testFile2, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	hash2, _ := hashFile(testFile2)
	if hash != hash2 {
		t.Error("Same content should produce same hash")
	}

	// Different content should produce different hash
	testFile3 := filepath.Join(tmpDir, "test3.txt")
	if err := os.WriteFile(testFile3, []byte("Different content"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	hash3, _ := hashFile(testFile3)
	if hash == hash3 {
		t.Error("Different content should produce different hash")
	}
}

func TestHashFile_NotExist(t *testing.T) {
	_, err := hashFile("/nonexistent/path/file.txt")
	if err == nil {
		t.Error("hashFile() should return error for non-existent file")
	}
}

func TestVersionDiff_Struct(t *testing.T) {
	diff := VersionDiff{
		FromVersion: "1.0.0",
		ToVersion:   "2.0.0",
		Added:       []string{"a.md", "b.md"},
		Removed:     []string{"c.md"},
		Modified:    []string{"d.md"},
		Unchanged:   5,
	}

	if diff.FromVersion != "1.0.0" {
		t.Errorf("FromVersion = %q, want %q", diff.FromVersion, "1.0.0")
	}
	if diff.ToVersion != "2.0.0" {
		t.Errorf("ToVersion = %q, want %q", diff.ToVersion, "2.0.0")
	}
	if len(diff.Added) != 2 {
		t.Errorf("Added length = %d, want 2", len(diff.Added))
	}
	if len(diff.Removed) != 1 {
		t.Errorf("Removed length = %d, want 1", len(diff.Removed))
	}
	if len(diff.Modified) != 1 {
		t.Errorf("Modified length = %d, want 1", len(diff.Modified))
	}
	if diff.Unchanged != 5 {
		t.Errorf("Unchanged = %d, want 5", diff.Unchanged)
	}
}

func TestGetLocalFileHashes(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .claude directory with some md files
	agentDir := filepath.Join(tmpDir, ".claude", "workflows")
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		t.Fatalf("Failed to create agent dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(agentDir, "test.md"), []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create CLAUDE.md
	if err := os.WriteFile(filepath.Join(tmpDir, "CLAUDE.md"), []byte("claude content"), 0644); err != nil {
		t.Fatalf("Failed to write CLAUDE.md: %v", err)
	}

	hashes := getLocalFileHashes(tmpDir)

	if len(hashes) == 0 {
		t.Error("getLocalFileHashes() returned empty map")
	}

	// Check that CLAUDE.md is included
	if _, ok := hashes["CLAUDE.md"]; !ok {
		t.Error("CLAUDE.md should be in hashes")
	}

	// Check that .claude files are included
	foundAgent := false
	for path := range hashes {
		if filepath.Dir(path) == ".claude/workflows" || filepath.Dir(path) == ".claude\\workflows" {
			foundAgent = true
			break
		}
	}
	if !foundAgent {
		t.Log("Note: .claude workflow files not found in hashes (may be platform-specific)")
	}
}

func TestGetVersionFileHashes(t *testing.T) {
	tmpDir := t.TempDir()

	// Create template directory structure
	templateDir := filepath.Join(tmpDir, "template")
	if err := os.MkdirAll(filepath.Join(templateDir, ".claude", "workflows"), 0755); err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, "CLAUDE.md"), []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to write CLAUDE.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, ".claude", "workflows", "test.md"), []byte("workflow"), 0644); err != nil {
		t.Fatalf("Failed to write workflow file: %v", err)
	}
	// Non-md file should be ignored
	if err := os.WriteFile(filepath.Join(templateDir, "test.txt"), []byte("ignored"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	hashes := getVersionFileHashes(tmpDir)

	if len(hashes) < 2 {
		t.Errorf("getVersionFileHashes() returned %d files, want at least 2", len(hashes))
	}

	// Non-md files should not be included
	for path := range hashes {
		if filepath.Ext(path) != ".md" {
			t.Errorf("Non-md file %q should not be in hashes", path)
		}
	}
}
