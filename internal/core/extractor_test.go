package core

import (
	"os"
	"path/filepath"
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
