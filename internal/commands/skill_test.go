package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// validSkillMD returns a well-formed SKILL.md content with the given name.
func validSkillMD(name, description string) string {
	return "---\nname: " + name + "\ndescription: " + description + "\n---\n\nBody content here.\n"
}

// setupSkillTestDir creates a temp dir with samuel.yaml and .claude/skills/.
// Returns the temp dir path and a cleanup function that restores the original cwd.
func setupSkillTestDir(t *testing.T) (string, func()) {
	t.Helper()
	dir := t.TempDir()

	// Create samuel.yaml so ConfigExists returns true
	if err := os.WriteFile(filepath.Join(dir, "samuel.yaml"), []byte("version: \"1.0.0\"\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create .claude/skills directory
	skillsDir := filepath.Join(dir, ".claude", "skills")
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		t.Fatal(err)
	}

	oldDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	return dir, func() { _ = os.Chdir(oldDir) }
}

// createSkillDir creates a skill directory with SKILL.md under the given skills dir.
func createSkillDir(t *testing.T, skillsDir, name, content string) {
	t.Helper()
	skillDir := filepath.Join(skillsDir, name)
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	skillMD := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillMD, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// --- runSkillCreate tests ---

func TestRunSkillCreate(t *testing.T) {
	t.Run("invalid_name_returns_error", func(t *testing.T) {
		_, cleanup := setupSkillTestDir(t)
		defer cleanup()

		cmd := &cobra.Command{}
		err := runSkillCreate(cmd, []string{"INVALID_NAME"})
		if err == nil {
			t.Fatal("expected error for invalid skill name")
		}
		if err.Error() != "skill name validation failed" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("name_with_consecutive_hyphens_rejected", func(t *testing.T) {
		_, cleanup := setupSkillTestDir(t)
		defer cleanup()

		cmd := &cobra.Command{}
		err := runSkillCreate(cmd, []string{"bad--name"})
		if err == nil {
			t.Fatal("expected error for consecutive hyphens")
		}
	})

	t.Run("no_config_returns_error", func(t *testing.T) {
		dir := t.TempDir()
		oldDir, _ := os.Getwd()
		if err := os.Chdir(dir); err != nil {
			t.Fatalf("failed to chdir: %v", err)
		}
		defer func() { _ = os.Chdir(oldDir) }()

		cmd := &cobra.Command{}
		err := runSkillCreate(cmd, []string{"my-skill"})
		if err == nil {
			t.Fatal("expected error when no samuel config")
		}
		expected := "no Samuel installation found. Run 'samuel init' first"
		if err.Error() != expected {
			t.Errorf("got %q, want %q", err.Error(), expected)
		}
	})

	t.Run("successful_creation", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		cmd := &cobra.Command{}
		err := runSkillCreate(cmd, []string{"my-new-skill"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify scaffold was created
		skillDir := filepath.Join(dir, ".claude", "skills", "my-new-skill")
		paths := []string{
			filepath.Join(skillDir, "SKILL.md"),
			filepath.Join(skillDir, "scripts", ".gitkeep"),
			filepath.Join(skillDir, "references", ".gitkeep"),
			filepath.Join(skillDir, "assets", ".gitkeep"),
		}
		for _, p := range paths {
			if _, err := os.Stat(p); os.IsNotExist(err) {
				t.Errorf("expected file to exist: %s", p)
			}
		}
	})

	t.Run("duplicate_skill_returns_error", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		// Create the skill first
		skillsDir := filepath.Join(dir, ".claude", "skills")
		createSkillDir(t, skillsDir, "existing-skill", validSkillMD("existing-skill", "A skill"))

		cmd := &cobra.Command{}
		err := runSkillCreate(cmd, []string{"existing-skill"})
		if err == nil {
			t.Fatal("expected error for duplicate skill")
		}
	})
}

// --- runSkillValidate tests ---

func TestRunSkillValidate(t *testing.T) {
	t.Run("no_skills_directory", func(t *testing.T) {
		dir := t.TempDir()
		// No .claude/skills directory
		if err := os.WriteFile(filepath.Join(dir, "samuel.yaml"), []byte("version: \"1.0.0\"\n"), 0644); err != nil {
			t.Fatal(err)
		}
		oldDir, _ := os.Getwd()
		if err := os.Chdir(dir); err != nil {
			t.Fatalf("failed to chdir: %v", err)
		}
		defer func() { _ = os.Chdir(oldDir) }()

		cmd := &cobra.Command{}
		err := runSkillValidate(cmd, []string{})
		if err != nil {
			t.Fatalf("expected nil error for missing skills dir, got: %v", err)
		}
	})

	t.Run("empty_skills_directory", func(t *testing.T) {
		_, cleanup := setupSkillTestDir(t)
		defer cleanup()

		cmd := &cobra.Command{}
		err := runSkillValidate(cmd, []string{})
		if err != nil {
			t.Fatalf("expected nil error for empty skills dir, got: %v", err)
		}
	})

	t.Run("specific_skill_not_found", func(t *testing.T) {
		_, cleanup := setupSkillTestDir(t)
		defer cleanup()

		cmd := &cobra.Command{}
		err := runSkillValidate(cmd, []string{"nonexistent"})
		if err == nil {
			t.Fatal("expected error for nonexistent skill")
		}
		if err.Error() != "skill 'nonexistent' not found" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("validate_specific_valid_skill", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		skillsDir := filepath.Join(dir, ".claude", "skills")
		createSkillDir(t, skillsDir, "good-skill", validSkillMD("good-skill", "A good skill"))

		cmd := &cobra.Command{}
		err := runSkillValidate(cmd, []string{"good-skill"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("validate_specific_invalid_skill", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		skillsDir := filepath.Join(dir, ".claude", "skills")
		// Name mismatch: dir name is "bad-skill" but SKILL.md name is "wrong-name"
		createSkillDir(t, skillsDir, "bad-skill", validSkillMD("wrong-name", "A skill"))

		cmd := &cobra.Command{}
		err := runSkillValidate(cmd, []string{"bad-skill"})
		if err == nil {
			t.Fatal("expected error for invalid skill")
		}
	})

	t.Run("validate_all_valid", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		skillsDir := filepath.Join(dir, ".claude", "skills")
		createSkillDir(t, skillsDir, "skill-a", validSkillMD("skill-a", "First skill"))
		createSkillDir(t, skillsDir, "skill-b", validSkillMD("skill-b", "Second skill"))

		cmd := &cobra.Command{}
		err := runSkillValidate(cmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("validate_all_mixed_valid_invalid", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		skillsDir := filepath.Join(dir, ".claude", "skills")
		createSkillDir(t, skillsDir, "valid-skill", validSkillMD("valid-skill", "Good"))
		// Missing description triggers validation error
		createSkillDir(t, skillsDir, "invalid-skill", "---\nname: invalid-skill\ndescription: \"\"\n---\n")

		cmd := &cobra.Command{}
		err := runSkillValidate(cmd, []string{})
		if err == nil {
			t.Fatal("expected error when some skills are invalid")
		}
	})
}

// --- runSkillList tests ---

func TestRunSkillList(t *testing.T) {
	t.Run("no_skills_directory", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "samuel.yaml"), []byte("version: \"1.0.0\"\n"), 0644); err != nil {
			t.Fatal(err)
		}
		oldDir, _ := os.Getwd()
		if err := os.Chdir(dir); err != nil {
			t.Fatalf("failed to chdir: %v", err)
		}
		defer func() { _ = os.Chdir(oldDir) }()

		cmd := &cobra.Command{}
		err := runSkillList(cmd, []string{})
		if err != nil {
			t.Fatalf("expected nil error for missing skills dir, got: %v", err)
		}
	})

	t.Run("empty_skills_directory", func(t *testing.T) {
		_, cleanup := setupSkillTestDir(t)
		defer cleanup()

		cmd := &cobra.Command{}
		err := runSkillList(cmd, []string{})
		if err != nil {
			t.Fatalf("expected nil error for empty skills dir, got: %v", err)
		}
	})

	t.Run("lists_valid_skills", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		skillsDir := filepath.Join(dir, ".claude", "skills")
		createSkillDir(t, skillsDir, "alpha", validSkillMD("alpha", "Alpha skill"))
		createSkillDir(t, skillsDir, "beta", validSkillMD("beta", "Beta skill"))

		cmd := &cobra.Command{}
		err := runSkillList(cmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("lists_skills_with_errors", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		skillsDir := filepath.Join(dir, ".claude", "skills")
		// Name mismatch causes validation error
		createSkillDir(t, skillsDir, "misnamed", validSkillMD("wrong-name", "Bad skill"))

		cmd := &cobra.Command{}
		err := runSkillList(cmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("long_description_truncated", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		skillsDir := filepath.Join(dir, ".claude", "skills")
		longDesc := "This is a very long description that exceeds sixty characters in total length to trigger truncation"
		createSkillDir(t, skillsDir, "verbose-skill", validSkillMD("verbose-skill", longDesc))

		cmd := &cobra.Command{}
		err := runSkillList(cmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// --- runSkillInfo tests ---

func TestRunSkillInfo(t *testing.T) {
	t.Run("skill_not_found", func(t *testing.T) {
		_, cleanup := setupSkillTestDir(t)
		defer cleanup()

		cmd := &cobra.Command{}
		err := runSkillInfo(cmd, []string{"nonexistent"})
		if err == nil {
			t.Fatal("expected error for nonexistent skill")
		}
		if err.Error() != "skill 'nonexistent' not found" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("valid_skill_with_metadata", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		skillsDir := filepath.Join(dir, ".claude", "skills")
		content := "---\nname: test-skill\ndescription: A test skill for unit testing\nlicense: MIT\ncompatibility: claude\n---\n\nThis is the body of the skill.\nIt has multiple lines.\n"
		createSkillDir(t, skillsDir, "test-skill", content)

		cmd := &cobra.Command{}
		err := runSkillInfo(cmd, []string{"test-skill"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("skill_with_optional_directories", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		skillsDir := filepath.Join(dir, ".claude", "skills")
		createSkillDir(t, skillsDir, "full-skill", validSkillMD("full-skill", "Full skill"))

		// Create optional directories
		skillDir := filepath.Join(skillsDir, "full-skill")
		for _, subdir := range []string{"scripts", "references", "assets"} {
			if err := os.MkdirAll(filepath.Join(skillDir, subdir), 0755); err != nil {
				t.Fatal(err)
			}
		}

		cmd := &cobra.Command{}
		err := runSkillInfo(cmd, []string{"full-skill"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("skill_with_multiline_description", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		skillsDir := filepath.Join(dir, ".claude", "skills")
		content := "---\nname: multi-desc\ndescription: |\n  Line one of description.\n  Line two of description.\n---\n\nBody.\n"
		createSkillDir(t, skillsDir, "multi-desc", content)

		cmd := &cobra.Command{}
		err := runSkillInfo(cmd, []string{"multi-desc"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("skill_with_custom_metadata", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		skillsDir := filepath.Join(dir, ".claude", "skills")
		content := "---\nname: meta-skill\ndescription: Skill with custom metadata\nmetadata:\n  author: tester\n  version: \"1.0\"\n---\n\nBody.\n"
		createSkillDir(t, skillsDir, "meta-skill", content)

		cmd := &cobra.Command{}
		err := runSkillInfo(cmd, []string{"meta-skill"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("skill_with_validation_errors", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		skillsDir := filepath.Join(dir, ".claude", "skills")
		// Name mismatch between dir name and SKILL.md name
		createSkillDir(t, skillsDir, "error-skill", validSkillMD("wrong-name", "Mismatched"))

		cmd := &cobra.Command{}
		// runSkillInfo does not return error for validation errors, it just displays them
		err := runSkillInfo(cmd, []string{"error-skill"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("skill_with_long_body", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		skillsDir := filepath.Join(dir, ".claude", "skills")
		// Create a body that exceeds 500 lines to trigger the warning
		body := ""
		for i := 0; i < 510; i++ {
			body += "Line of content.\n"
		}
		content := "---\nname: long-skill\ndescription: Skill with long body\n---\n\n" + body
		createSkillDir(t, skillsDir, "long-skill", content)

		cmd := &cobra.Command{}
		err := runSkillInfo(cmd, []string{"long-skill"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("skill_without_skill_md", func(t *testing.T) {
		dir, cleanup := setupSkillTestDir(t)
		defer cleanup()

		// Create skill directory without SKILL.md
		skillDir := filepath.Join(dir, ".claude", "skills", "empty-skill")
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			t.Fatal(err)
		}

		cmd := &cobra.Command{}
		// LoadSkillInfo returns info with Errors populated, not a Go error
		err := runSkillInfo(cmd, []string{"empty-skill"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
