package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ar4mirez/samuel/internal/core"
)

func TestCheckCLAUDEMD(t *testing.T) {
	t.Run("missing_file", func(t *testing.T) {
		dir := t.TempDir()
		result := checkCLAUDEMD(dir)
		if result.passed {
			t.Error("expected check to fail when CLAUDE.md is missing")
		}
		if result.name != "CLAUDE.md" {
			t.Errorf("expected name CLAUDE.md, got %q", result.name)
		}
		if !result.fixable {
			t.Error("expected check to be fixable")
		}
	})

	t.Run("present_without_version", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# My Project"), 0644); err != nil {
			t.Fatal(err)
		}
		result := checkCLAUDEMD(dir)
		if !result.passed {
			t.Error("expected check to pass when CLAUDE.md exists")
		}
		if result.message != "Present" {
			t.Errorf("expected message %q, got %q", "Present", result.message)
		}
	})

	t.Run("present_with_bold_version", func(t *testing.T) {
		dir := t.TempDir()
		content := "# Project\n**Current Version**: 1.2.3\n"
		if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		result := checkCLAUDEMD(dir)
		if !result.passed {
			t.Error("expected check to pass")
		}
		if result.message != "Present (v1.2.3)" {
			t.Errorf("expected message %q, got %q", "Present (v1.2.3)", result.message)
		}
	})

	t.Run("present_with_plain_version", func(t *testing.T) {
		dir := t.TempDir()
		content := "Current Version: 4.5.6\n"
		if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		result := checkCLAUDEMD(dir)
		if !result.passed {
			t.Error("expected check to pass")
		}
		if result.message != "Present (v4.5.6)" {
			t.Errorf("expected message %q, got %q", "Present (v4.5.6)", result.message)
		}
	})
}

func TestCheckAGENTSMD(t *testing.T) {
	t.Run("missing_file", func(t *testing.T) {
		dir := t.TempDir()
		result := checkAGENTSMD(dir)
		if result.passed {
			t.Error("expected check to fail when AGENTS.md is missing")
		}
		if !result.fixable {
			t.Error("expected check to be fixable")
		}
	})

	t.Run("present", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte("# Agents"), 0644); err != nil {
			t.Fatal(err)
		}
		result := checkAGENTSMD(dir)
		if !result.passed {
			t.Error("expected check to pass when AGENTS.md exists")
		}
		if result.message != "Present" {
			t.Errorf("expected message %q, got %q", "Present", result.message)
		}
	})
}

func TestCheckDirectoryStructure(t *testing.T) {
	t.Run("all_present", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, ".claude", "skills"), 0755); err != nil {
			t.Fatal(err)
		}
		result, missingDirs := checkDirectoryStructure(dir)
		if !result.passed {
			t.Error("expected check to pass when directories exist")
		}
		if len(missingDirs) != 0 {
			t.Errorf("expected no missing dirs, got %v", missingDirs)
		}
	})

	t.Run("all_missing", func(t *testing.T) {
		dir := t.TempDir()
		result, missingDirs := checkDirectoryStructure(dir)
		if result.passed {
			t.Error("expected check to fail when directories are missing")
		}
		if len(missingDirs) != 2 {
			t.Errorf("expected 2 missing dirs, got %d: %v", len(missingDirs), missingDirs)
		}
		if !result.fixable {
			t.Error("expected check to be fixable")
		}
	})

	t.Run("only_claude_dir", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.Mkdir(filepath.Join(dir, ".claude"), 0755); err != nil {
			t.Fatal(err)
		}
		result, missingDirs := checkDirectoryStructure(dir)
		if result.passed {
			t.Error("expected check to fail when skills dir is missing")
		}
		if len(missingDirs) != 1 {
			t.Errorf("expected 1 missing dir, got %d: %v", len(missingDirs), missingDirs)
		}
		if missingDirs[0] != ".claude/skills" {
			t.Errorf("expected missing dir .claude/skills, got %q", missingDirs[0])
		}
	})
}

func TestCheckInstalledSkills(t *testing.T) {
	t.Run("all_present", func(t *testing.T) {
		dir := t.TempDir()
		// Create a fake component path with SKILL.md
		skillPath := filepath.Join(dir, ".claude", "skills", "go-guide")
		if err := os.MkdirAll(skillPath, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(skillPath, "SKILL.md"), []byte("# Go"), 0644); err != nil {
			t.Fatal(err)
		}

		finder := func(name string) *core.Component {
			return &core.Component{
				Name: name,
				Path: ".claude/skills/go-guide",
			}
		}
		result := checkInstalledSkills(dir, []string{"go"}, "Language guides", finder)
		if !result.passed {
			t.Error("expected check to pass when skill files exist")
		}
	})

	t.Run("some_missing", func(t *testing.T) {
		dir := t.TempDir()
		finder := func(name string) *core.Component {
			return &core.Component{
				Name: name,
				Path: ".claude/skills/" + name,
			}
		}
		result := checkInstalledSkills(dir, []string{"go", "python"}, "Language guides", finder)
		if result.passed {
			t.Error("expected check to fail when skill files are missing")
		}
		if !result.fixable {
			t.Error("expected check to be fixable")
		}
	})

	t.Run("unknown_component", func(t *testing.T) {
		dir := t.TempDir()
		finder := func(name string) *core.Component {
			return nil // component not found in registry
		}
		result := checkInstalledSkills(dir, []string{"unknown"}, "Language guides", finder)
		if !result.passed {
			t.Error("expected check to pass when component is not in registry (skipped)")
		}
	})

	t.Run("empty_names", func(t *testing.T) {
		dir := t.TempDir()
		finder := func(name string) *core.Component {
			return nil
		}
		result := checkInstalledSkills(dir, []string{}, "Language guides", finder)
		if !result.passed {
			t.Error("expected check to pass with empty names list")
		}
		if result.message != "All 0 installed language guides present" {
			t.Errorf("unexpected message: %q", result.message)
		}
	})
}

func TestCheckInstalledComponents(t *testing.T) {
	t.Run("with_known_components", func(t *testing.T) {
		dir := t.TempDir()

		// Use a real language from the registry
		goComp := core.FindLanguage("go")
		if goComp == nil {
			t.Skip("go language not in registry")
		}

		// Create the SKILL.md so it passes
		skillPath := filepath.Join(dir, goComp.Path)
		if err := os.MkdirAll(skillPath, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(skillPath, "SKILL.md"), []byte("# Go"), 0644); err != nil {
			t.Fatal(err)
		}

		config := &core.Config{
			Installed: core.InstalledItems{
				Languages: []string{"go"},
			},
		}

		results := checkInstalledComponents(dir, config)
		if len(results) != 3 {
			t.Errorf("expected 3 results (languages, frameworks, workflows), got %d", len(results))
		}
		// Languages check should pass since we created the SKILL.md
		if !results[0].passed {
			t.Errorf("expected languages check to pass, got: %s", results[0].message)
		}
	})

	t.Run("with_workflow_all", func(t *testing.T) {
		dir := t.TempDir()
		config := &core.Config{
			Installed: core.InstalledItems{
				Workflows: []string{"all"},
			},
		}

		results := checkInstalledComponents(dir, config)
		if len(results) != 3 {
			t.Errorf("expected 3 results, got %d", len(results))
		}
		// Workflows check should fail since no SKILL.md files exist
		workflowResult := results[2]
		if workflowResult.passed {
			t.Error("expected workflows check to fail when files are missing")
		}
	})
}

func TestCheckSkillsIntegrity(t *testing.T) {
	t.Run("no_skills_directory", func(t *testing.T) {
		dir := t.TempDir()
		results := checkSkillsIntegrity(dir)
		if results != nil {
			t.Errorf("expected nil results when skills dir doesn't exist, got %v", results)
		}
	})

	t.Run("empty_skills_directory", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, ".claude", "skills"), 0755); err != nil {
			t.Fatal(err)
		}
		results := checkSkillsIntegrity(dir)
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if !results[0].passed {
			t.Error("expected check to pass for empty skills dir")
		}
		if results[0].message != "No skills installed" {
			t.Errorf("unexpected message: %q", results[0].message)
		}
	})

	t.Run("valid_skill", func(t *testing.T) {
		dir := t.TempDir()
		skillDir := filepath.Join(dir, ".claude", "skills", "test-skill")
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			t.Fatal(err)
		}
		// Write a minimal valid SKILL.md
		content := `---
name: test-skill
description: A test skill
version: 1.0.0
---
# Test Skill
Body content here.`
		if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		results := checkSkillsIntegrity(dir)
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if !results[0].passed {
			t.Errorf("expected check to pass for valid skill, got: %s", results[0].message)
		}
	})

	t.Run("invalid_skill", func(t *testing.T) {
		dir := t.TempDir()
		skillDir := filepath.Join(dir, ".claude", "skills", "bad-skill")
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			t.Fatal(err)
		}
		// Write an invalid SKILL.md (missing required metadata)
		if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# No metadata"), 0644); err != nil {
			t.Fatal(err)
		}
		results := checkSkillsIntegrity(dir)
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].passed {
			t.Error("expected check to fail for invalid skill")
		}
	})

	t.Run("mixed_valid_and_invalid", func(t *testing.T) {
		dir := t.TempDir()
		skillsDir := filepath.Join(dir, ".claude", "skills")

		// Valid skill
		validDir := filepath.Join(skillsDir, "valid-skill")
		if err := os.MkdirAll(validDir, 0755); err != nil {
			t.Fatal(err)
		}
		validContent := `---
name: valid-skill
description: A valid skill
version: 1.0.0
---
# Valid Skill
Body.`
		if err := os.WriteFile(filepath.Join(validDir, "SKILL.md"), []byte(validContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Invalid skill
		invalidDir := filepath.Join(skillsDir, "invalid-skill")
		if err := os.MkdirAll(invalidDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(invalidDir, "SKILL.md"), []byte("no metadata"), 0644); err != nil {
			t.Fatal(err)
		}

		results := checkSkillsIntegrity(dir)
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].passed {
			t.Error("expected check to fail when some skills are invalid")
		}
	})
}

func TestCheckAutoHealth(t *testing.T) {
	t.Run("valid_prd", func(t *testing.T) {
		dir := t.TempDir()
		autoDir := filepath.Join(dir, ".claude", "auto")
		if err := os.MkdirAll(autoDir, 0755); err != nil {
			t.Fatal(err)
		}

		prd := &core.AutoPRD{
			Version: "1.0",
			Project: core.AutoProject{
				Name:        "test",
				Description: "test project",
				CreatedAt:   "2026-01-01T00:00:00Z",
				UpdatedAt:   "2026-01-01T00:00:00Z",
			},
			Config: core.AutoConfig{
				MaxIterations: 10,
				AITool:        "claude",
			},
			Tasks: []core.AutoTask{
				{ID: "1", Title: "Task 1", Status: "completed"},
				{ID: "2", Title: "Task 2", Status: "pending"},
			},
		}
		data, err := json.MarshalIndent(prd, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(autoDir, "prd.json"), data, 0644); err != nil {
			t.Fatal(err)
		}

		results := checkAutoHealth(dir)
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if !results[0].passed {
			t.Errorf("expected check to pass for valid prd, got: %s", results[0].message)
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		dir := t.TempDir()
		autoDir := filepath.Join(dir, ".claude", "auto")
		if err := os.MkdirAll(autoDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(autoDir, "prd.json"), []byte("{invalid"), 0644); err != nil {
			t.Fatal(err)
		}

		results := checkAutoHealth(dir)
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].passed {
			t.Error("expected check to fail for invalid JSON")
		}
	})

	t.Run("missing_prd", func(t *testing.T) {
		dir := t.TempDir()
		results := checkAutoHealth(dir)
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].passed {
			t.Error("expected check to fail when prd.json is missing")
		}
	})

	t.Run("validation_errors", func(t *testing.T) {
		dir := t.TempDir()
		autoDir := filepath.Join(dir, ".claude", "auto")
		if err := os.MkdirAll(autoDir, 0755); err != nil {
			t.Fatal(err)
		}

		// PRD with missing required fields
		prd := &core.AutoPRD{
			// Version intentionally empty — triggers validation error
			Project: core.AutoProject{
				Name: "test",
			},
		}
		data, err := json.MarshalIndent(prd, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(autoDir, "prd.json"), data, 0644); err != nil {
			t.Fatal(err)
		}

		results := checkAutoHealth(dir)
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].passed {
			t.Error("expected check to fail when PRD has validation errors")
		}
	})
}

func TestCheckLocalModifications(t *testing.T) {
	t.Run("claude_md_present", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# Modified"), 0644); err != nil {
			t.Fatal(err)
		}
		config := &core.Config{}
		results := checkLocalModifications(dir, config)
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if !results[0].passed {
			t.Error("expected check to pass when CLAUDE.md exists")
		}
	})

	t.Run("claude_md_absent", func(t *testing.T) {
		dir := t.TempDir()
		config := &core.Config{}
		results := checkLocalModifications(dir, config)
		if results != nil {
			t.Errorf("expected nil results when CLAUDE.md doesn't exist, got %v", results)
		}
	})
}

func TestCheckModification(t *testing.T) {
	t.Run("existing_file", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, "test.txt")
		if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
		if !checkModification(filePath) {
			t.Error("expected true for existing file")
		}
	})

	t.Run("nonexistent_file", func(t *testing.T) {
		if checkModification("/nonexistent/path/file.txt") {
			t.Error("expected false for nonexistent file")
		}
	})
}

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "bold_version",
			content: "**Current Version**: 1.2.3",
			want:    "1.2.3",
		},
		{
			name:    "plain_version",
			content: "Current Version: 4.5.6",
			want:    "4.5.6",
		},
		{
			name:    "no_version",
			content: "# Some content\nNo version here",
			want:    "",
		},
		{
			name:    "empty_content",
			content: "",
			want:    "",
		},
		{
			name:    "version_in_multiline",
			content: "# Header\n\n**Current Version**: 10.20.30\n\nOther stuff",
			want:    "10.20.30",
		},
		{
			name:    "bold_preferred_over_plain",
			content: "**Current Version**: 1.0.0\nCurrent Version: 2.0.0",
			want:    "1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractVersion(tt.content)
			if got != tt.want {
				t.Errorf("extractVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}
