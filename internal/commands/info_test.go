package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ar4mirez/samuel/internal/core"
)

func TestFindComponent(t *testing.T) {
	tests := []struct {
		componentType string
		name          string
		wantFound     bool
	}{
		{"language", "go", true},
		{"language", "python", true},
		{"language", "nonexistent", false},
		{"framework", "react", true},
		{"framework", "nextjs", true},
		{"framework", "nonexistent", false},
		{"workflow", "create-prd", true},
		{"workflow", "nonexistent", false},
		{"invalid", "go", false},
	}

	for _, tt := range tests {
		t.Run(tt.componentType+"/"+tt.name, func(t *testing.T) {
			got := findComponent(tt.componentType, tt.name)
			if (got != nil) != tt.wantFound {
				t.Errorf("findComponent(%q, %q) found = %v, want %v",
					tt.componentType, tt.name, got != nil, tt.wantFound)
			}
		})
	}
}

func TestCheckInstallStatus(t *testing.T) {
	config := &core.Config{
		Installed: core.InstalledItems{
			Languages:  []string{"go", "python"},
			Frameworks: []string{"react"},
			Workflows:  []string{"all"},
		},
	}

	tests := []struct {
		componentType string
		name          string
		want          bool
	}{
		{"language", "go", true},
		{"language", "rust", false},
		{"framework", "react", true},
		{"framework", "vue", false},
		{"workflow", "create-prd", true}, // "all" includes everything
		{"invalid", "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.componentType+"/"+tt.name, func(t *testing.T) {
			got := checkInstallStatus(config, tt.componentType, tt.name)
			if got != tt.want {
				t.Errorf("checkInstallStatus(config, %q, %q) = %v, want %v",
					tt.componentType, tt.name, got, tt.want)
			}
		})
	}
}

func TestCheckInstallStatus_NilConfig(t *testing.T) {
	got := checkInstallStatus(nil, "language", "go")
	if got {
		t.Error("checkInstallStatus(nil, ...) = true, want false")
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		size int64
		want string
	}{
		{0, "0 bytes"},
		{1, "1 bytes"},
		{500, "500 bytes"},
		{1023, "1023 bytes"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{10240, "10.0 KB"},
		{1048576, "1.0 MB"},
		{1572864, "1.5 MB"},
		{10485760, "10.0 MB"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatFileSize(tt.size)
			if got != tt.want {
				t.Errorf("formatFileSize(%d) = %q, want %q", tt.size, got, tt.want)
			}
		})
	}
}

func TestGetFrameworksForLanguage(t *testing.T) {
	tests := []struct {
		lang      string
		wantCount int
		contains  string
	}{
		{"typescript", 3, "react"},
		{"python", 3, "django"},
		{"go", 3, "gin"},
		{"rust", 3, "axum"},
		{"nonexistent", 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			got := getFrameworksForLanguage(tt.lang)
			if len(got) != tt.wantCount {
				t.Errorf("getFrameworksForLanguage(%q) returned %d items, want %d",
					tt.lang, len(got), tt.wantCount)
			}
			if tt.contains != "" {
				found := false
				for _, r := range got {
					if r.Name == tt.contains {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("getFrameworksForLanguage(%q) should contain %q", tt.lang, tt.contains)
				}
			}
		})
	}
}

func TestGetLanguageForFramework(t *testing.T) {
	tests := []struct {
		framework string
		wantLang  string
		wantCount int
	}{
		{"react", "typescript", 1},
		{"nextjs", "typescript", 1},
		{"django", "python", 1},
		{"gin", "go", 1},
		{"axum", "rust", 1},
		{"rails", "ruby", 1},
		{"flutter", "dart", 1},
		{"nonexistent", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.framework, func(t *testing.T) {
			got := getLanguageForFramework(tt.framework)
			if len(got) != tt.wantCount {
				t.Errorf("getLanguageForFramework(%q) returned %d items, want %d",
					tt.framework, len(got), tt.wantCount)
			}
			if tt.wantCount > 0 && got[0].Name != tt.wantLang {
				t.Errorf("getLanguageForFramework(%q)[0].Name = %q, want %q",
					tt.framework, got[0].Name, tt.wantLang)
			}
		})
	}
}

func TestGetRelatedComponents(t *testing.T) {
	tests := []struct {
		name          string
		componentType string
		componentName string
		wantCount     int
	}{
		{"language with frameworks", "language", "python", 3},
		{"framework with language", "framework", "react", 1},
		{"workflow has no related", "workflow", "create-prd", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			component := &core.Component{Name: tt.componentName}
			got := getRelatedComponents(component, tt.componentType)
			if len(got) != tt.wantCount {
				t.Errorf("getRelatedComponents(%q) returned %d items, want %d",
					tt.componentName, len(got), tt.wantCount)
			}
		})
	}
}

func TestIsInstalled(t *testing.T) {
	config := &core.Config{
		Installed: core.InstalledItems{
			Languages:  []string{"go"},
			Frameworks: []string{"react"},
			Workflows:  []string{"create-prd"},
		},
	}

	tests := []struct {
		componentType string
		name          string
		want          bool
	}{
		{"language", "go", true},
		{"language", "python", false},
		{"framework", "react", true},
		{"framework", "vue", false},
		{"workflow", "create-prd", true},
		{"workflow", "security-audit", false},
		{"invalid", "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.componentType+"/"+tt.name, func(t *testing.T) {
			got := isInstalled(config, tt.componentType, tt.name)
			if got != tt.want {
				t.Errorf("isInstalled(config, %q, %q) = %v, want %v",
					tt.componentType, tt.name, got, tt.want)
			}
		})
	}
}

func TestGetFilePreview(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")

	content := `Line 1
Line 2
Line 3
Line 4
Line 5`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	tests := []struct {
		lines     int
		wantLines int
	}{
		{3, 3},
		{5, 5},
		{10, 5}, // File only has 5 lines
		{0, 0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			preview, err := getFilePreview(testFile, tt.lines)
			if err != nil {
				t.Fatalf("getFilePreview() error = %v", err)
			}
			// Count lines in output
			if tt.wantLines == 0 {
				if preview != "" {
					t.Errorf("getFilePreview(%d) should return empty, got %q", tt.lines, preview)
				}
				return
			}
			// Check that preview is not empty for non-zero lines
			if preview == "" && tt.lines > 0 {
				t.Errorf("getFilePreview(%d) returned empty preview", tt.lines)
			}
		})
	}
}

func TestGetFilePreview_NotExist(t *testing.T) {
	_, err := getFilePreview("/nonexistent/file.md", 5)
	if err == nil {
		t.Error("getFilePreview() should return error for non-existent file")
	}
}

func TestRelatedComponent_Struct(t *testing.T) {
	r := RelatedComponent{
		Name:        "react",
		Type:        "framework",
		Description: "React framework",
	}

	if r.Name != "react" {
		t.Errorf("Name = %q, want %q", r.Name, "react")
	}
	if r.Type != "framework" {
		t.Errorf("Type = %q, want %q", r.Type, "framework")
	}
	if r.Description != "React framework" {
		t.Errorf("Description = %q, want %q", r.Description, "React framework")
	}
}

// Test framework language mappings are complete
func TestLanguageFrameworkMappings(t *testing.T) {
	// These languages should have framework mappings
	languages := []string{
		"typescript", "python", "go", "rust", "kotlin",
		"java", "csharp", "php", "swift", "ruby", "dart",
	}

	for _, lang := range languages {
		t.Run(lang, func(t *testing.T) {
			frameworks := getFrameworksForLanguage(lang)
			if len(frameworks) == 0 {
				t.Errorf("Language %q has no framework mappings", lang)
			}
		})
	}
}

// Test framework to language mappings
func TestFrameworkLanguageMappings(t *testing.T) {
	frameworks := []string{
		"react", "nextjs", "express",
		"django", "fastapi", "flask",
		"gin", "echo", "fiber",
		"axum", "actix-web", "rocket",
	}

	for _, fw := range frameworks {
		t.Run(fw, func(t *testing.T) {
			langs := getLanguageForFramework(fw)
			if len(langs) == 0 {
				t.Errorf("Framework %q has no language mapping", fw)
			}
		})
	}
}
