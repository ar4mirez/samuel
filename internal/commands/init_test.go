package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ar4mirez/samuel/internal/core"
)

func TestExpandLanguages(t *testing.T) {
	tests := []struct {
		name     string
		flags    []string
		expected []string
	}{
		{
			name:     "direct language name",
			flags:    []string{"go"},
			expected: []string{"go"},
		},
		{
			name:     "alias ts to typescript",
			flags:    []string{"ts"},
			expected: []string{"typescript"},
		},
		{
			name:     "alias js to typescript",
			flags:    []string{"js"},
			expected: []string{"typescript"},
		},
		{
			name:     "alias javascript to typescript",
			flags:    []string{"javascript"},
			expected: []string{"typescript"},
		},
		{
			name:     "alias py to python",
			flags:    []string{"py"},
			expected: []string{"python"},
		},
		{
			name:     "alias cs to csharp",
			flags:    []string{"cs"},
			expected: []string{"csharp"},
		},
		{
			name:     "alias c++ to cpp",
			flags:    []string{"c++"},
			expected: []string{"cpp"},
		},
		{
			name:     "alias c to cpp",
			flags:    []string{"c"},
			expected: []string{"cpp"},
		},
		{
			name:     "alias rb to ruby",
			flags:    []string{"rb"},
			expected: []string{"ruby"},
		},
		{
			name:     "alias sh to shell",
			flags:    []string{"sh"},
			expected: []string{"shell"},
		},
		{
			name:     "alias bash to shell",
			flags:    []string{"bash"},
			expected: []string{"shell"},
		},
		{
			name:     "comma-separated values",
			flags:    []string{"go,py,ts"},
			expected: []string{"go", "python", "typescript"},
		},
		{
			name:     "multiple flags",
			flags:    []string{"go", "rust"},
			expected: []string{"go", "rust"},
		},
		{
			name:     "unknown language filtered out",
			flags:    []string{"go", "nonexistent"},
			expected: []string{"go"},
		},
		{
			name:     "all unknown returns nil",
			flags:    []string{"nonexistent"},
			expected: nil,
		},
		{
			name:     "empty flags",
			flags:    []string{},
			expected: nil,
		},
		{
			name:     "uppercase normalized",
			flags:    []string{"GO", "Python"},
			expected: []string{"go", "python"},
		},
		{
			name:     "whitespace trimmed",
			flags:    []string{" go , rust "},
			expected: []string{"go", "rust"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandLanguages(tt.flags)
			if len(result) != len(tt.expected) {
				t.Fatalf("got %d results %v, want %d %v",
					len(result), result, len(tt.expected), tt.expected)
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("result[%d] = %q, want %q", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestExpandFrameworks(t *testing.T) {
	tests := []struct {
		name     string
		flags    []string
		expected []string
	}{
		{
			name:     "direct framework name",
			flags:    []string{"react"},
			expected: []string{"react"},
		},
		{
			name:     "alias next to nextjs",
			flags:    []string{"next"},
			expected: []string{"nextjs"},
		},
		{
			name:     "alias spring to spring-boot-java",
			flags:    []string{"spring"},
			expected: []string{"spring-boot-java"},
		},
		{
			name:     "comma-separated values",
			flags:    []string{"react,django"},
			expected: []string{"react", "django"},
		},
		{
			name:     "multiple flags",
			flags:    []string{"gin", "echo"},
			expected: []string{"gin", "echo"},
		},
		{
			name:     "unknown framework filtered out",
			flags:    []string{"react", "nonexistent"},
			expected: []string{"react"},
		},
		{
			name:     "all unknown returns nil",
			flags:    []string{"nonexistent"},
			expected: nil,
		},
		{
			name:     "empty flags",
			flags:    []string{},
			expected: nil,
		},
		{
			name:     "uppercase normalized",
			flags:    []string{"React", "DJANGO"},
			expected: []string{"react", "django"},
		},
		{
			name:     "whitespace trimmed",
			flags:    []string{" gin , echo "},
			expected: []string{"gin", "echo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandFrameworks(tt.flags)
			if len(result) != len(tt.expected) {
				t.Fatalf("got %d results %v, want %d %v",
					len(result), result, len(tt.expected), tt.expected)
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("result[%d] = %q, want %q", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestIsSamuelRepository(t *testing.T) {
	t.Run("empty directory is not samuel repo", func(t *testing.T) {
		dir := t.TempDir()
		if isSamuelRepository(dir) {
			t.Error("empty directory should not be detected as samuel repo")
		}
	})

	t.Run("template dir without CLAUDE.md", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, "template"), 0755); err != nil {
			t.Fatal(err)
		}
		if isSamuelRepository(dir) {
			t.Error("template dir without CLAUDE.md should not match")
		}
	})

	t.Run("template dir with CLAUDE.md", func(t *testing.T) {
		dir := t.TempDir()
		templateDir := filepath.Join(dir, "template")
		if err := os.MkdirAll(templateDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(templateDir, "CLAUDE.md"), []byte("# test"), 0644); err != nil {
			t.Fatal(err)
		}
		if !isSamuelRepository(dir) {
			t.Error("template/CLAUDE.md should be detected as samuel repo")
		}
	})

	t.Run("packages/cli with go.mod", func(t *testing.T) {
		dir := t.TempDir()
		cliDir := filepath.Join(dir, "packages", "cli")
		if err := os.MkdirAll(cliDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(cliDir, "go.mod"), []byte("module samuel"), 0644); err != nil {
			t.Fatal(err)
		}
		if !isSamuelRepository(dir) {
			t.Error("packages/cli/go.mod should be detected as samuel repo")
		}
	})

	t.Run("packages/cli without go.mod", func(t *testing.T) {
		dir := t.TempDir()
		cliDir := filepath.Join(dir, "packages", "cli")
		if err := os.MkdirAll(cliDir, 0755); err != nil {
			t.Fatal(err)
		}
		if isSamuelRepository(dir) {
			t.Error("packages/cli without go.mod should not match")
		}
	})

	t.Run("nonexistent directory", func(t *testing.T) {
		if isSamuelRepository("/nonexistent/path/12345") {
			t.Error("nonexistent directory should not match")
		}
	})
}

func TestGetRelevantFrameworks(t *testing.T) {
	t.Run("typescript returns react, nextjs, express", func(t *testing.T) {
		result := getRelevantFrameworks([]string{"typescript"})
		expected := map[string]bool{"react": true, "nextjs": true, "express": true}

		if len(result) != len(expected) {
			t.Fatalf("got %d frameworks, want %d", len(result), len(expected))
		}
		for _, fw := range result {
			if !expected[fw.Name] {
				t.Errorf("unexpected framework %q", fw.Name)
			}
		}
	})

	t.Run("python returns django, fastapi, flask", func(t *testing.T) {
		result := getRelevantFrameworks([]string{"python"})
		expected := map[string]bool{"django": true, "fastapi": true, "flask": true}

		if len(result) != len(expected) {
			t.Fatalf("got %d frameworks, want %d", len(result), len(expected))
		}
		for _, fw := range result {
			if !expected[fw.Name] {
				t.Errorf("unexpected framework %q", fw.Name)
			}
		}
	})

	t.Run("go returns gin, echo, fiber", func(t *testing.T) {
		result := getRelevantFrameworks([]string{"go"})
		expected := map[string]bool{"gin": true, "echo": true, "fiber": true}

		if len(result) != len(expected) {
			t.Fatalf("got %d frameworks, want %d", len(result), len(expected))
		}
		for _, fw := range result {
			if !expected[fw.Name] {
				t.Errorf("unexpected framework %q", fw.Name)
			}
		}
	})

	t.Run("multiple languages deduplicates", func(t *testing.T) {
		result := getRelevantFrameworks([]string{"typescript", "python"})
		expected := map[string]bool{
			"react": true, "nextjs": true, "express": true,
			"django": true, "fastapi": true, "flask": true,
		}

		if len(result) != len(expected) {
			t.Fatalf("got %d frameworks, want %d", len(result), len(expected))
		}
		for _, fw := range result {
			if !expected[fw.Name] {
				t.Errorf("unexpected framework %q", fw.Name)
			}
		}
	})

	t.Run("unknown language returns empty", func(t *testing.T) {
		result := getRelevantFrameworks([]string{"nonexistent"})
		if len(result) != 0 {
			t.Errorf("got %d frameworks, want 0", len(result))
		}
	})

	t.Run("empty languages returns empty", func(t *testing.T) {
		result := getRelevantFrameworks([]string{})
		if len(result) != 0 {
			t.Errorf("got %d frameworks, want 0", len(result))
		}
	})

	t.Run("nil languages returns empty", func(t *testing.T) {
		result := getRelevantFrameworks(nil)
		if len(result) != 0 {
			t.Errorf("got %d frameworks, want 0", len(result))
		}
	})

	t.Run("result has correct Component type", func(t *testing.T) {
		result := getRelevantFrameworks([]string{"go"})
		for _, fw := range result {
			if fw.Category != "framework" {
				t.Errorf("framework %q has category %q, want %q",
					fw.Name, fw.Category, "framework")
			}
			if fw.Path == "" {
				t.Errorf("framework %q has empty path", fw.Name)
			}
		}
	})

	t.Run("all mapped languages have frameworks", func(t *testing.T) {
		mapped := []string{
			"typescript", "python", "go", "rust", "kotlin",
			"java", "csharp", "php", "swift", "ruby", "dart",
		}
		for _, lang := range mapped {
			result := getRelevantFrameworks([]string{lang})
			if len(result) == 0 {
				t.Errorf("language %q should have frameworks", lang)
			}
		}
	})

	t.Run("returned components match registry", func(t *testing.T) {
		result := getRelevantFrameworks([]string{"go"})
		for _, fw := range result {
			found := core.FindFramework(fw.Name)
			if found == nil {
				t.Errorf("framework %q from getRelevantFrameworks not in registry", fw.Name)
			}
		}
	})
}
