package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyFromCache_SingleFile(t *testing.T) {
	// Create temp cache structure: cachePath/template/somefile.txt
	cacheDir := t.TempDir()
	templateDir := filepath.Join(cacheDir, TemplatePrefix)
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, "somefile.txt"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	destDir := t.TempDir()

	err := CopyFromCache(cacheDir, destDir, "somefile.txt")
	if err != nil {
		t.Fatalf("CopyFromCache file: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(destDir, "somefile.txt"))
	if err != nil {
		t.Fatalf("failed to read copied file: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("got %q, want %q", string(data), "hello")
	}
}

func TestCopyFromCache_Directory(t *testing.T) {
	// Create temp cache structure with a skill directory
	cacheDir := t.TempDir()
	skillDir := filepath.Join(cacheDir, TemplatePrefix, ".claude", "skills", "flask")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# Flask Skill"), 0644); err != nil {
		t.Fatal(err)
	}

	destDir := t.TempDir()

	err := CopyFromCache(cacheDir, destDir, filepath.Join(".claude", "skills", "flask"))
	if err != nil {
		t.Fatalf("CopyFromCache directory: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(destDir, ".claude", "skills", "flask", "SKILL.md"))
	if err != nil {
		t.Fatalf("failed to read copied file: %v", err)
	}
	if string(data) != "# Flask Skill" {
		t.Errorf("got %q, want %q", string(data), "# Flask Skill")
	}
}

func TestCopyFromCache_DirectoryWithReferences(t *testing.T) {
	// Create temp cache structure with references/ subdirectory
	cacheDir := t.TempDir()
	skillDir := filepath.Join(cacheDir, TemplatePrefix, ".claude", "skills", "flask")
	refsDir := filepath.Join(skillDir, "references")
	if err := os.MkdirAll(refsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# Flask"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(refsDir, "patterns.md"), []byte("# Patterns"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(refsDir, "security.md"), []byte("# Security"), 0644); err != nil {
		t.Fatal(err)
	}

	destDir := t.TempDir()

	err := CopyFromCache(cacheDir, destDir, filepath.Join(".claude", "skills", "flask"))
	if err != nil {
		t.Fatalf("CopyFromCache with references: %v", err)
	}

	// Verify SKILL.md
	data, err := os.ReadFile(filepath.Join(destDir, ".claude", "skills", "flask", "SKILL.md"))
	if err != nil {
		t.Fatalf("failed to read SKILL.md: %v", err)
	}
	if string(data) != "# Flask" {
		t.Errorf("SKILL.md: got %q, want %q", string(data), "# Flask")
	}

	// Verify references/patterns.md
	data, err = os.ReadFile(filepath.Join(destDir, ".claude", "skills", "flask", "references", "patterns.md"))
	if err != nil {
		t.Fatalf("failed to read references/patterns.md: %v", err)
	}
	if string(data) != "# Patterns" {
		t.Errorf("patterns.md: got %q, want %q", string(data), "# Patterns")
	}

	// Verify references/security.md
	data, err = os.ReadFile(filepath.Join(destDir, ".claude", "skills", "flask", "references", "security.md"))
	if err != nil {
		t.Fatalf("failed to read references/security.md: %v", err)
	}
	if string(data) != "# Security" {
		t.Errorf("security.md: got %q, want %q", string(data), "# Security")
	}
}

func TestCopyFromCache_NestedSubdirectories(t *testing.T) {
	// Test deeply nested structure: skill/references/examples/basic.md
	cacheDir := t.TempDir()
	deepDir := filepath.Join(cacheDir, TemplatePrefix, ".claude", "skills", "test-skill", "references", "examples")
	if err := os.MkdirAll(deepDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(deepDir, "basic.md"), []byte("example"), 0644); err != nil {
		t.Fatal(err)
	}

	destDir := t.TempDir()

	err := CopyFromCache(cacheDir, destDir, filepath.Join(".claude", "skills", "test-skill"))
	if err != nil {
		t.Fatalf("CopyFromCache nested dirs: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(destDir, ".claude", "skills", "test-skill", "references", "examples", "basic.md"))
	if err != nil {
		t.Fatalf("failed to read nested file: %v", err)
	}
	if string(data) != "example" {
		t.Errorf("got %q, want %q", string(data), "example")
	}
}

func TestCopyFromCache_SourceNotFound(t *testing.T) {
	cacheDir := t.TempDir()
	templateDir := filepath.Join(cacheDir, TemplatePrefix)
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatal(err)
	}

	destDir := t.TempDir()

	err := CopyFromCache(cacheDir, destDir, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent source, got nil")
	}
}

func TestCopySingleFile(t *testing.T) {
	srcDir := t.TempDir()
	srcFile := filepath.Join(srcDir, "test.txt")
	if err := os.WriteFile(srcFile, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	destDir := t.TempDir()
	dstFile := filepath.Join(destDir, "sub", "test.txt")

	err := copySingleFile(srcFile, dstFile)
	if err != nil {
		t.Fatalf("copySingleFile: %v", err)
	}

	data, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}
	if string(data) != "content" {
		t.Errorf("got %q, want %q", string(data), "content")
	}
}

func TestCopyDirRecursive(t *testing.T) {
	srcDir := t.TempDir()
	subDir := filepath.Join(srcDir, "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "a.txt"), []byte("a"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "b.txt"), []byte("b"), 0644); err != nil {
		t.Fatal(err)
	}

	destDir := filepath.Join(t.TempDir(), "dest")

	err := copyDirRecursive(srcDir, destDir)
	if err != nil {
		t.Fatalf("copyDirRecursive: %v", err)
	}

	// Verify root file
	data, err := os.ReadFile(filepath.Join(destDir, "a.txt"))
	if err != nil {
		t.Fatalf("failed to read a.txt: %v", err)
	}
	if string(data) != "a" {
		t.Errorf("a.txt: got %q, want %q", string(data), "a")
	}

	// Verify sub file
	data, err = os.ReadFile(filepath.Join(destDir, "sub", "b.txt"))
	if err != nil {
		t.Fatalf("failed to read sub/b.txt: %v", err)
	}
	if string(data) != "b" {
		t.Errorf("b.txt: got %q, want %q", string(data), "b")
	}
}

func TestCopyFromCache_MultipleSubdirectories(t *testing.T) {
	// Simulate a skill with references/, scripts/, and assets/
	cacheDir := t.TempDir()
	skillBase := filepath.Join(cacheDir, TemplatePrefix, ".claude", "skills", "full-skill")
	dirs := []string{
		filepath.Join(skillBase, "references"),
		filepath.Join(skillBase, "scripts"),
		filepath.Join(skillBase, "assets"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatal(err)
		}
	}

	files := map[string]string{
		filepath.Join(skillBase, "SKILL.md"):              "# Skill",
		filepath.Join(skillBase, "references", "ref.md"):  "ref",
		filepath.Join(skillBase, "scripts", "run.sh"):     "#!/bin/bash",
		filepath.Join(skillBase, "assets", "diagram.svg"): "<svg/>",
	}
	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	destDir := t.TempDir()

	err := CopyFromCache(cacheDir, destDir, filepath.Join(".claude", "skills", "full-skill"))
	if err != nil {
		t.Fatalf("CopyFromCache multi-subdir: %v", err)
	}

	// Verify all files were copied
	checks := map[string]string{
		filepath.Join(destDir, ".claude", "skills", "full-skill", "SKILL.md"):              "# Skill",
		filepath.Join(destDir, ".claude", "skills", "full-skill", "references", "ref.md"):  "ref",
		filepath.Join(destDir, ".claude", "skills", "full-skill", "scripts", "run.sh"):     "#!/bin/bash",
		filepath.Join(destDir, ".claude", "skills", "full-skill", "assets", "diagram.svg"): "<svg/>",
	}
	for path, want := range checks {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("missing file %s: %v", path, err)
			continue
		}
		if string(data) != want {
			t.Errorf("%s: got %q, want %q", path, string(data), want)
		}
	}
}

func TestValidateContainedPath(t *testing.T) {
	base := "/safe/base"
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid simple file", "file.txt", false},
		{"valid nested file", "sub/dir/file.txt", false},
		{"traversal with dotdot", "../../etc/passwd", true},
		{"traversal mixed", "sub/../../etc/passwd", true},
		{"absolute path joined", "/etc/passwd", false}, // filepath.Join treats as relative on Unix
		{"dotdot only", "..", true},
		{"current dir", ".", false},
		{"trailing dotdot", "subdir/..", false}, // resolves to base itself
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateContainedPath(base, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateContainedPath(%q, %q) error = %v, wantErr %v",
					base, tt.path, err, tt.wantErr)
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), "path traversal") {
				t.Errorf("expected 'path traversal' in error, got: %v", err)
			}
		})
	}
}

func TestReadFile_PathTraversal(t *testing.T) {
	destDir := t.TempDir()
	ext := NewExtractor("", destDir)

	// Create a file outside destDir
	outsideDir := t.TempDir()
	outsideFile := filepath.Join(outsideDir, "secret.txt")
	if err := os.WriteFile(outsideFile, []byte("secret"), 0644); err != nil {
		t.Fatal(err)
	}

	// Attempt traversal
	relPath, _ := filepath.Rel(destDir, outsideFile)
	_, err := ext.ReadFile(relPath)
	if err == nil {
		t.Error("ReadFile should reject path traversal, got nil error")
	}
	if err != nil && !strings.Contains(err.Error(), "path traversal") {
		t.Errorf("expected path traversal error, got: %v", err)
	}
}

func TestWriteFile_PathTraversal(t *testing.T) {
	destDir := t.TempDir()
	ext := NewExtractor("", destDir)

	err := ext.WriteFile("../../evil.txt", []byte("malicious"))
	if err == nil {
		t.Error("WriteFile should reject path traversal, got nil error")
	}
}

func TestRemoveFile_PathTraversal(t *testing.T) {
	destDir := t.TempDir()
	ext := NewExtractor("", destDir)

	err := ext.RemoveFile("../../important.txt")
	if err == nil {
		t.Error("RemoveFile should reject path traversal, got nil error")
	}
}

func TestFileExists_PathTraversal(t *testing.T) {
	destDir := t.TempDir()
	ext := NewExtractor("", destDir)

	// Create a file outside destDir to confirm it's not detected
	outsideDir := t.TempDir()
	outsideFile := filepath.Join(outsideDir, "exists.txt")
	if err := os.WriteFile(outsideFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	relPath, _ := filepath.Rel(destDir, outsideFile)
	if ext.FileExists(relPath) {
		t.Error("FileExists should return false for path traversal")
	}
}

func TestBackupFile_PathTraversal(t *testing.T) {
	destDir := t.TempDir()
	backupDir := t.TempDir()
	ext := NewExtractor("", destDir)

	// Create a legitimate file in destDir
	if err := os.WriteFile(filepath.Join(destDir, "legit.txt"), []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	// Attempt source traversal (reading from outside destDir)
	err := ext.BackupFile("../../etc/passwd", backupDir)
	if err == nil {
		t.Error("BackupFile should reject source path traversal")
	}

	// Attempt destination traversal (writing outside backupDir)
	err = ext.BackupFile("legit.txt", backupDir)
	if err != nil {
		t.Errorf("BackupFile of legit file should succeed: %v", err)
	}
}

func TestValidateExtraction_PathTraversal(t *testing.T) {
	destDir := t.TempDir()
	ext := NewExtractor("", destDir)

	// Create a file inside destDir
	if err := os.WriteFile(filepath.Join(destDir, "valid.txt"), []byte("ok"), 0644); err != nil {
		t.Fatal(err)
	}

	missing := ext.ValidateExtraction([]string{"valid.txt", "../../etc/passwd"})
	if len(missing) != 1 {
		t.Fatalf("expected 1 missing, got %d: %v", len(missing), missing)
	}
	if missing[0] != "../../etc/passwd" {
		t.Errorf("expected traversal path in missing, got: %s", missing[0])
	}
}

func TestReadFile_ValidPath(t *testing.T) {
	destDir := t.TempDir()
	ext := NewExtractor("", destDir)

	// Create a valid file
	subDir := filepath.Join(destDir, "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "file.txt"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	data, err := ext.ReadFile("sub/file.txt")
	if err != nil {
		t.Fatalf("ReadFile valid path: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("got %q, want %q", string(data), "hello")
	}
}

func TestWriteFile_ValidPath(t *testing.T) {
	destDir := t.TempDir()
	ext := NewExtractor("", destDir)

	err := ext.WriteFile("sub/new.txt", []byte("content"))
	if err != nil {
		t.Fatalf("WriteFile valid path: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(destDir, "sub", "new.txt"))
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if string(data) != "content" {
		t.Errorf("got %q, want %q", string(data), "content")
	}
}
