package core

import (
	"testing"
)

func TestFindLanguage(t *testing.T) {
	tests := []struct {
		name      string
		wantFound bool
	}{
		{"go", true},
		{"python", true},
		{"typescript", true},
		{"rust", true},
		{"nonexistent", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindLanguage(tt.name)
			if (got != nil) != tt.wantFound {
				t.Errorf("FindLanguage(%q) found = %v, want %v", tt.name, got != nil, tt.wantFound)
			}
			if got != nil && got.Name != tt.name {
				t.Errorf("FindLanguage(%q).Name = %q", tt.name, got.Name)
			}
		})
	}
}

func TestFindFramework(t *testing.T) {
	tests := []struct {
		name      string
		wantFound bool
	}{
		{"react", true},
		{"nextjs", true},
		{"django", true},
		{"gin", true},
		{"axum", true},
		{"nonexistent", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindFramework(tt.name)
			if (got != nil) != tt.wantFound {
				t.Errorf("FindFramework(%q) found = %v, want %v", tt.name, got != nil, tt.wantFound)
			}
			if got != nil && got.Name != tt.name {
				t.Errorf("FindFramework(%q).Name = %q", tt.name, got.Name)
			}
		})
	}
}

func TestFindWorkflow(t *testing.T) {
	tests := []struct {
		name      string
		wantFound bool
	}{
		{"create-prd", true},
		{"code-review", true},
		{"security-audit", true},
		{"nonexistent", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindWorkflow(tt.name)
			if (got != nil) != tt.wantFound {
				t.Errorf("FindWorkflow(%q) found = %v, want %v", tt.name, got != nil, tt.wantFound)
			}
			if got != nil && got.Name != tt.name {
				t.Errorf("FindWorkflow(%q).Name = %q", tt.name, got.Name)
			}
		})
	}
}

func TestFindTemplate(t *testing.T) {
	tests := []struct {
		name      string
		wantFound bool
	}{
		{"full", true},
		{"minimal", true},
		{"nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindTemplate(tt.name)
			if (got != nil) != tt.wantFound {
				t.Errorf("FindTemplate(%q) found = %v, want %v", tt.name, got != nil, tt.wantFound)
			}
			if got != nil && got.Name != tt.name {
				t.Errorf("FindTemplate(%q).Name = %q", tt.name, got.Name)
			}
		})
	}
}

func TestGetAllLanguageNames(t *testing.T) {
	names := GetAllLanguageNames()
	if len(names) == 0 {
		t.Error("GetAllLanguageNames() returned empty slice")
	}

	// Check that some expected languages are included
	expected := []string{"go", "python", "typescript", "rust"}
	for _, exp := range expected {
		found := false
		for _, name := range names {
			if name == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetAllLanguageNames() missing %q", exp)
		}
	}
}

func TestGetAllFrameworkNames(t *testing.T) {
	names := GetAllFrameworkNames()
	if len(names) == 0 {
		t.Error("GetAllFrameworkNames() returned empty slice")
	}

	// Check that some expected frameworks are included
	expected := []string{"react", "django", "gin", "axum"}
	for _, exp := range expected {
		found := false
		for _, name := range names {
			if name == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetAllFrameworkNames() missing %q", exp)
		}
	}
}

func TestGetAllWorkflowNames(t *testing.T) {
	names := GetAllWorkflowNames()
	if len(names) == 0 {
		t.Error("GetAllWorkflowNames() returned empty slice")
	}

	// Check that some expected workflows are included
	expected := []string{"create-prd", "code-review", "security-audit"}
	for _, exp := range expected {
		found := false
		for _, name := range names {
			if name == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetAllWorkflowNames() missing %q", exp)
		}
	}
}

func TestGetComponentPaths(t *testing.T) {
	tests := []struct {
		name       string
		languages  []string
		frameworks []string
		workflows  []string
		wantMin    int
	}{
		{
			name:       "empty",
			languages:  []string{},
			frameworks: []string{},
			workflows:  []string{},
			wantMin:    len(CoreFiles), // At least core files
		},
		{
			name:       "with languages",
			languages:  []string{"go", "python"},
			frameworks: []string{},
			workflows:  []string{},
			wantMin:    len(CoreFiles) + 2,
		},
		{
			name:       "with frameworks",
			languages:  []string{},
			frameworks: []string{"react", "gin"},
			workflows:  []string{},
			wantMin:    len(CoreFiles) + 2,
		},
		{
			name:       "all workflows",
			languages:  []string{},
			frameworks: []string{},
			workflows:  []string{"all"},
			wantMin:    len(CoreFiles) + len(Workflows),
		},
		{
			name:       "specific workflows",
			languages:  []string{},
			frameworks: []string{},
			workflows:  []string{"create-prd", "code-review"},
			wantMin:    len(CoreFiles) + 2,
		},
		{
			name:       "mixed",
			languages:  []string{"go"},
			frameworks: []string{"gin"},
			workflows:  []string{"create-prd"},
			wantMin:    len(CoreFiles) + 3,
		},
		{
			name:       "nonexistent ignored",
			languages:  []string{"nonexistent"},
			frameworks: []string{"nonexistent"},
			workflows:  []string{"nonexistent"},
			wantMin:    len(CoreFiles), // Non-existent should be ignored
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := GetComponentPaths(tt.languages, tt.frameworks, tt.workflows)
			if len(paths) < tt.wantMin {
				t.Errorf("GetComponentPaths() returned %d paths, want at least %d", len(paths), tt.wantMin)
			}
		})
	}
}

func TestGetSourcePath(t *testing.T) {
	tests := []struct {
		destPath string
		want     string
	}{
		{"CLAUDE.md", "template/CLAUDE.md"},
		{".claude/skills/create-prd/SKILL.md", "template/.claude/skills/create-prd/SKILL.md"},
		{"", "template/"},
	}

	for _, tt := range tests {
		t.Run(tt.destPath, func(t *testing.T) {
			got := GetSourcePath(tt.destPath)
			if got != tt.want {
				t.Errorf("GetSourcePath(%q) = %q, want %q", tt.destPath, got, tt.want)
			}
		})
	}
}

func TestGetSourcePaths(t *testing.T) {
	destPaths := []string{"CLAUDE.md", ".claude/skills/create-prd/SKILL.md"}
	srcPaths := GetSourcePaths(destPaths)

	if len(srcPaths) != len(destPaths) {
		t.Errorf("GetSourcePaths() returned %d paths, want %d", len(srcPaths), len(destPaths))
	}

	for i, src := range srcPaths {
		expected := TemplatePrefix + destPaths[i]
		if src != expected {
			t.Errorf("GetSourcePaths()[%d] = %q, want %q", i, src, expected)
		}
	}
}

func TestLanguagesHavePaths(t *testing.T) {
	for _, lang := range Languages {
		if lang.Name == "" {
			t.Error("Language has empty name")
		}
		if lang.Path == "" {
			t.Errorf("Language %q has empty path", lang.Name)
		}
		if lang.Description == "" {
			t.Errorf("Language %q has empty description", lang.Name)
		}
	}
}

func TestFrameworksHavePaths(t *testing.T) {
	for _, fw := range Frameworks {
		if fw.Name == "" {
			t.Error("Framework has empty name")
		}
		if fw.Path == "" {
			t.Errorf("Framework %q has empty path", fw.Name)
		}
		if fw.Description == "" {
			t.Errorf("Framework %q has empty description", fw.Name)
		}
	}
}

func TestWorkflowsHavePaths(t *testing.T) {
	for _, wf := range Workflows {
		if wf.Name == "" {
			t.Error("Workflow has empty name")
		}
		if wf.Path == "" {
			t.Errorf("Workflow %q has empty path", wf.Name)
		}
		if wf.Description == "" {
			t.Errorf("Workflow %q has empty description", wf.Name)
		}
	}
}

func TestCoreFilesNotEmpty(t *testing.T) {
	if len(CoreFiles) == 0 {
		t.Error("CoreFiles is empty")
	}
}

func TestTemplatesHaveFields(t *testing.T) {
	for _, tmpl := range Templates {
		if tmpl.Name == "" {
			t.Error("Template has empty name")
		}
		if tmpl.Description == "" {
			t.Errorf("Template %q has empty description", tmpl.Name)
		}
	}
}

func TestConstants(t *testing.T) {
	if DefaultRegistry == "" {
		t.Error("DefaultRegistry is empty")
	}
	if DefaultOwner == "" {
		t.Error("DefaultOwner is empty")
	}
	if DefaultRepo == "" {
		t.Error("DefaultRepo is empty")
	}
	if TemplatePrefix == "" {
		t.Error("TemplatePrefix is empty")
	}
}
