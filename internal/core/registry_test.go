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

func TestFindSkill(t *testing.T) {
	tests := []struct {
		name      string
		wantFound bool
	}{
		{"go-guide", true},
		{"commit-message", true},
		{"react", true},
		{"auto", true},
		{"typescript-guide", true},
		{"nonexistent", false},
		{"", false},
		{"GO-GUIDE", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindSkill(tt.name)
			if (got != nil) != tt.wantFound {
				t.Errorf("FindSkill(%q) found = %v, want %v",
					tt.name, got != nil, tt.wantFound)
			}
			if got != nil && got.Name != tt.name {
				t.Errorf("FindSkill(%q).Name = %q", tt.name, got.Name)
			}
		})
	}
}

func TestGetAllSkillNames(t *testing.T) {
	names := GetAllSkillNames()
	if len(names) == 0 {
		t.Fatal("GetAllSkillNames() returned empty slice")
	}

	if len(names) != len(Skills) {
		t.Errorf("GetAllSkillNames() returned %d names, want %d",
			len(names), len(Skills))
	}

	expected := []string{
		"commit-message", "go-guide", "react", "auto", "sync-claude-md",
	}
	for _, exp := range expected {
		found := false
		for _, name := range names {
			if name == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetAllSkillNames() missing %q", exp)
		}
	}
}

func TestGetLanguageSkills(t *testing.T) {
	skills := GetLanguageSkills()
	if len(skills) == 0 {
		t.Fatal("GetLanguageSkills() returned empty slice")
	}

	for _, s := range skills {
		if s.Category != "language" {
			t.Errorf("GetLanguageSkills() returned skill %q with category %q",
				s.Name, s.Category)
		}
	}

	// Verify known language skills are present
	expected := map[string]bool{
		"go-guide":     false,
		"python-guide": false,
		"rust-guide":   false,
	}
	for _, s := range skills {
		if _, ok := expected[s.Name]; ok {
			expected[s.Name] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Errorf("GetLanguageSkills() missing expected skill %q", name)
		}
	}
}

func TestGetFrameworkSkills(t *testing.T) {
	skills := GetFrameworkSkills()
	if len(skills) == 0 {
		t.Fatal("GetFrameworkSkills() returned empty slice")
	}

	for _, s := range skills {
		if s.Category != "framework" {
			t.Errorf("GetFrameworkSkills() returned skill %q with category %q",
				s.Name, s.Category)
		}
	}

	// Verify known framework skills are present
	expected := map[string]bool{
		"react":  false,
		"django": false,
		"gin":    false,
	}
	for _, s := range skills {
		if _, ok := expected[s.Name]; ok {
			expected[s.Name] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Errorf("GetFrameworkSkills() missing expected skill %q", name)
		}
	}
}

func TestGetWorkflowSkills(t *testing.T) {
	skills := GetWorkflowSkills()
	if len(skills) == 0 {
		t.Fatal("GetWorkflowSkills() returned empty slice")
	}

	for _, s := range skills {
		if s.Category != "workflow" {
			t.Errorf("GetWorkflowSkills() returned skill %q with category %q",
				s.Name, s.Category)
		}
	}

	// Verify known workflow skills are present
	expected := map[string]bool{
		"create-prd":     false,
		"code-review":    false,
		"auto":           false,
		"troubleshooting": false,
	}
	for _, s := range skills {
		if _, ok := expected[s.Name]; ok {
			expected[s.Name] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Errorf("GetWorkflowSkills() missing expected skill %q", name)
		}
	}
}

func TestSkillCategoryCoverage(t *testing.T) {
	lang := GetLanguageSkills()
	fw := GetFrameworkSkills()
	wf := GetWorkflowSkills()

	categorized := len(lang) + len(fw) + len(wf)
	// commit-message has no category, so total should be categorized + uncategorized
	if categorized >= len(Skills) {
		// All or most skills should have a category
		return
	}
	uncategorized := len(Skills) - categorized
	if uncategorized > 5 {
		t.Errorf("Too many uncategorized skills: %d out of %d",
			uncategorized, len(Skills))
	}
}

func TestLanguageToSkillName(t *testing.T) {
	tests := []struct {
		langName string
		want     string
	}{
		{"go", "go-guide"},
		{"python", "python-guide"},
		{"typescript", "typescript-guide"},
		{"rust", "rust-guide"},
		{"", "-guide"},
	}

	for _, tt := range tests {
		t.Run(tt.langName, func(t *testing.T) {
			got := LanguageToSkillName(tt.langName)
			if got != tt.want {
				t.Errorf("LanguageToSkillName(%q) = %q, want %q",
					tt.langName, got, tt.want)
			}
		})
	}
}

func TestSkillToLanguageName(t *testing.T) {
	tests := []struct {
		skillName string
		want      string
	}{
		{"go-guide", "go"},
		{"python-guide", "python"},
		{"typescript-guide", "typescript"},
		{"rust-guide", "rust"},
		{"react", "react"},
		{"commit-message", "commit-message"},
		{"guide", "guide"},
		{"-guide", "-guide"},
		{"a-guide", "a"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.skillName, func(t *testing.T) {
			got := SkillToLanguageName(tt.skillName)
			if got != tt.want {
				t.Errorf("SkillToLanguageName(%q) = %q, want %q",
					tt.skillName, got, tt.want)
			}
		})
	}
}

func TestLanguageToSkillNameRoundTrip(t *testing.T) {
	for _, lang := range Languages {
		skillName := LanguageToSkillName(lang.Name)
		langName := SkillToLanguageName(skillName)
		if langName != lang.Name {
			t.Errorf("round-trip failed: %q -> %q -> %q",
				lang.Name, skillName, langName)
		}
	}
}

func TestFrameworkToSkillName(t *testing.T) {
	tests := []struct {
		fwName string
		want   string
	}{
		{"react", "react"},
		{"gin", "gin"},
		{"django", "django"},
	}

	for _, tt := range tests {
		t.Run(tt.fwName, func(t *testing.T) {
			got := FrameworkToSkillName(tt.fwName)
			if got != tt.want {
				t.Errorf("FrameworkToSkillName(%q) = %q, want %q",
					tt.fwName, got, tt.want)
			}
		})
	}
}

func TestWorkflowToSkillName(t *testing.T) {
	tests := []struct {
		wfName string
		want   string
	}{
		{"create-prd", "create-prd"},
		{"auto", "auto"},
		{"code-review", "code-review"},
	}

	for _, tt := range tests {
		t.Run(tt.wfName, func(t *testing.T) {
			got := WorkflowToSkillName(tt.wfName)
			if got != tt.want {
				t.Errorf("WorkflowToSkillName(%q) = %q, want %q",
					tt.wfName, got, tt.want)
			}
		})
	}
}

func TestSkillsHaveFields(t *testing.T) {
	for _, s := range Skills {
		if s.Name == "" {
			t.Error("Skill has empty name")
		}
		if s.Path == "" {
			t.Errorf("Skill %q has empty path", s.Name)
		}
		if s.Description == "" {
			t.Errorf("Skill %q has empty description", s.Name)
		}
	}
}
