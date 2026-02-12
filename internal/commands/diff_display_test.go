package commands

import (
	"testing"
)

func TestCategorizeFiles(t *testing.T) {
	files := []string{
		".agent/skills/go-guide/SKILL.md",
		".agent/skills/python-guide/SKILL.md",
		".agent/skills/python-guide/references/patterns.md",
		".agent/skills/react/SKILL.md",
		".agent/skills/nextjs/SKILL.md",
		".agent/skills/django/SKILL.md",
		".agent/workflows/create-prd.md",
		"CLAUDE.md",
		"README.md",
	}

	langs, fws, wfs := categorizeFiles(files)

	if len(langs) != 2 {
		t.Errorf("categorizeFiles() langs = %d, want 2; got %v", len(langs), langs)
	}
	if len(fws) != 3 {
		t.Errorf("categorizeFiles() fws = %d, want 3; got %v", len(fws), fws)
	}
	if len(wfs) != 1 {
		t.Errorf("categorizeFiles() wfs = %d, want 1", len(wfs))
	}

	// Check extracted names use skill directory names
	foundGo := false
	for _, l := range langs {
		if l == "go-guide" {
			foundGo = true
		}
	}
	if !foundGo {
		t.Errorf("categorizeFiles() should extract 'go-guide' from skill path, got %v", langs)
	}

	// Check framework names are extracted correctly
	foundReact := false
	for _, f := range fws {
		if f == "react" {
			foundReact = true
		}
	}
	if !foundReact {
		t.Errorf("categorizeFiles() should extract 'react' from framework skill path, got %v", fws)
	}
}

func TestCategorizeFiles_Empty(t *testing.T) {
	langs, fws, wfs := categorizeFiles([]string{})

	if len(langs) != 0 || len(fws) != 0 || len(wfs) != 0 {
		t.Error("categorizeFiles([]) should return empty slices")
	}
}

func TestCategorizeFiles_NoComponents(t *testing.T) {
	files := []string{
		"CLAUDE.md",
		"README.md",
		"other.txt",
	}

	langs, fws, wfs := categorizeFiles(files)

	if len(langs) != 0 || len(fws) != 0 || len(wfs) != 0 {
		t.Error("categorizeFiles() with no component files should return empty slices")
	}
}

func TestCategorizeFiles_FrameworkSkillDedup(t *testing.T) {
	// Multiple files from same framework skill should be deduplicated
	files := []string{
		".agent/skills/react/SKILL.md",
		".agent/skills/react/references/patterns.md",
	}

	_, fws, _ := categorizeFiles(files)

	if len(fws) != 1 {
		t.Errorf("categorizeFiles() should deduplicate framework skill files, got %d: %v", len(fws), fws)
	}
}

func TestCategorizeOtherFiles(t *testing.T) {
	added := []string{
		".agent/skills/go-guide/SKILL.md",
		"CLAUDE.md",
		"new-file.md",
	}
	modified := []string{
		".agent/skills/react/SKILL.md",
		"README.md",
	}
	removed := []string{
		".agent/workflows/old.md",
		"deleted.md",
	}

	addedOther, modifiedOther, removedOther := categorizeOtherFiles(added, modified, removed)

	if len(addedOther) != 2 { // CLAUDE.md and new-file.md
		t.Errorf("categorizeOtherFiles() addedOther = %d, want 2; got %v", len(addedOther), addedOther)
	}
	if len(modifiedOther) != 1 { // README.md
		t.Errorf("categorizeOtherFiles() modifiedOther = %d, want 1", len(modifiedOther))
	}
	if len(removedOther) != 1 { // deleted.md
		t.Errorf("categorizeOtherFiles() removedOther = %d, want 1", len(removedOther))
	}
}

func TestCategorizeOtherFiles_AllComponents(t *testing.T) {
	added := []string{".agent/skills/go-guide/SKILL.md"}
	modified := []string{".agent/skills/react/SKILL.md"}
	removed := []string{".agent/workflows/old.md"}

	addedOther, modifiedOther, removedOther := categorizeOtherFiles(added, modified, removed)

	if len(addedOther) != 0 || len(modifiedOther) != 0 || len(removedOther) != 0 {
		t.Error("categorizeOtherFiles() with only component files should return empty slices")
	}
}

func TestExtractComponentName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{".agent/workflows/create-prd.md", "create-prd"},
		{"CLAUDE.md", "CLAUDE"},
		{"path/to/file.md", "file"},
		{"simple.md", "simple"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := extractComponentName(tt.path)
			if got != tt.want {
				t.Errorf("extractComponentName(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestDisplayFileDiff_NoDifferences(t *testing.T) {
	diff := &VersionDiff{
		FromVersion: "1.0.0",
		ToVersion:   "1.0.0",
		Added:       []string{},
		Removed:     []string{},
		Modified:    []string{},
		Unchanged:   10,
	}

	// This test just verifies it doesn't panic
	// The actual output goes to stdout
	displayFileDiff(diff)
}

func TestDisplayFileDiff_WithChanges(t *testing.T) {
	diff := &VersionDiff{
		FromVersion: "1.0.0",
		ToVersion:   "2.0.0",
		Added:       []string{"new.md"},
		Removed:     []string{"old.md"},
		Modified:    []string{"changed.md"},
		Unchanged:   5,
	}

	// This test just verifies it doesn't panic
	displayFileDiff(diff)
}

func TestDisplayComponentDiff(t *testing.T) {
	diff := &VersionDiff{
		FromVersion: "1.0.0",
		ToVersion:   "2.0.0",
		Added: []string{
			".agent/skills/rust-guide/SKILL.md",
			".agent/skills/axum/SKILL.md",
		},
		Removed: []string{
			".agent/workflows/old.md",
		},
		Modified: []string{
			"CLAUDE.md",
		},
		Unchanged: 50,
	}

	// This test just verifies it doesn't panic
	displayComponentDiff(diff)
}

func TestExtractSkillDirName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{".agent/skills/go-guide/SKILL.md", "go-guide"},
		{".agent/skills/commit-message/SKILL.md", "commit-message"},
		{".agent/skills/go-guide/references/patterns.md", "go-guide"},
		{".agent/skills/react/SKILL.md", "react"},
		{".agent/skills/spring-boot-kotlin/SKILL.md", "spring-boot-kotlin"},
		{"CLAUDE.md", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := extractSkillDirName(tt.path)
			if got != tt.want {
				t.Errorf("extractSkillDirName(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsFrameworkSkill(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"react", true},
		{"django", true},
		{"gin", true},
		{"go-guide", false},
		{"commit-message", false},
		{"nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isFrameworkSkill(tt.name)
			if got != tt.want {
				t.Errorf("isFrameworkSkill(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestDisplayComponentChanges(t *testing.T) {
	added := []string{"new1", "new2"}
	modified := []string{"mod1"}
	removed := []string{"rem1", "rem2", "rem3"}

	// This test just verifies it doesn't panic
	displayComponentChanges(added, modified, removed)
}

func TestDisplayComponentChanges_Empty(t *testing.T) {
	// This test just verifies it doesn't panic with empty slices
	displayComponentChanges(nil, nil, nil)
	displayComponentChanges([]string{}, []string{}, []string{})
}
