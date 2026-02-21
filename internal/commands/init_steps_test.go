package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ar4mirez/samuel/internal/core"
	"github.com/spf13/cobra"
)

// newInitCmd creates a fresh cobra.Command with init flags for testing.
func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "init", RunE: func(cmd *cobra.Command, args []string) error { return nil }}
	cmd.Flags().StringP("template", "t", "", "Template")
	cmd.Flags().StringSlice("languages", nil, "Languages")
	cmd.Flags().StringSlice("frameworks", nil, "Frameworks")
	cmd.Flags().BoolP("force", "f", false, "Force")
	cmd.Flags().Bool("non-interactive", false, "Non-interactive")
	return cmd
}

func TestParseInitFlags(t *testing.T) {
	t.Run("defaults_no_args_no_flags", func(t *testing.T) {
		cmd := newInitCmd()
		flags, err := parseInitFlags(cmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if flags.force {
			t.Error("force should be false by default")
		}
		if flags.nonInteractive {
			t.Error("nonInteractive should be false by default")
		}
		if flags.templateName != "" {
			t.Errorf("templateName = %q, want empty", flags.templateName)
		}
		if len(flags.languageFlags) != 0 {
			t.Errorf("languageFlags = %v, want empty", flags.languageFlags)
		}
		if len(flags.frameworkFlags) != 0 {
			t.Errorf("frameworkFlags = %v, want empty", flags.frameworkFlags)
		}
		if flags.cliProvided {
			t.Error("cliProvided should be false when no flags set")
		}
		if flags.createDir {
			t.Error("createDir should be false for '.' target")
		}
		// absTargetDir should be the absolute path of "."
		cwd, _ := os.Getwd()
		if flags.absTargetDir != cwd {
			t.Errorf("absTargetDir = %q, want %q", flags.absTargetDir, cwd)
		}
	})

	t.Run("target_arg_existing_dir", func(t *testing.T) {
		dir := t.TempDir()
		cmd := newInitCmd()
		flags, err := parseInitFlags(cmd, []string{dir})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if flags.absTargetDir != dir {
			t.Errorf("absTargetDir = %q, want %q", flags.absTargetDir, dir)
		}
		if flags.createDir {
			t.Error("createDir should be false for existing directory")
		}
	})

	t.Run("target_arg_nonexistent_dir", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "new-project")
		cmd := newInitCmd()
		flags, err := parseInitFlags(cmd, []string{dir})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if flags.absTargetDir != dir {
			t.Errorf("absTargetDir = %q, want %q", flags.absTargetDir, dir)
		}
		if !flags.createDir {
			t.Error("createDir should be true for nonexistent directory")
		}
	})

	t.Run("force_flag", func(t *testing.T) {
		cmd := newInitCmd()
		if err := cmd.Flags().Set("force", "true"); err != nil {
			t.Fatal(err)
		}
		flags, err := parseInitFlags(cmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !flags.force {
			t.Error("force should be true")
		}
	})

	t.Run("non_interactive_flag", func(t *testing.T) {
		cmd := newInitCmd()
		if err := cmd.Flags().Set("non-interactive", "true"); err != nil {
			t.Fatal(err)
		}
		flags, err := parseInitFlags(cmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !flags.nonInteractive {
			t.Error("nonInteractive should be true")
		}
	})

	t.Run("template_flag_sets_cli_provided", func(t *testing.T) {
		cmd := newInitCmd()
		if err := cmd.Flags().Set("template", "minimal"); err != nil {
			t.Fatal(err)
		}
		flags, err := parseInitFlags(cmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if flags.templateName != "minimal" {
			t.Errorf("templateName = %q, want %q", flags.templateName, "minimal")
		}
		if !flags.cliProvided {
			t.Error("cliProvided should be true when template is set")
		}
	})

	t.Run("languages_flag_sets_cli_provided", func(t *testing.T) {
		cmd := newInitCmd()
		if err := cmd.Flags().Set("languages", "go,python"); err != nil {
			t.Fatal(err)
		}
		flags, err := parseInitFlags(cmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(flags.languageFlags) != 2 {
			t.Fatalf("languageFlags = %v, want 2 elements", flags.languageFlags)
		}
		if !flags.cliProvided {
			t.Error("cliProvided should be true when languages are set")
		}
	})

	t.Run("frameworks_flag_sets_cli_provided", func(t *testing.T) {
		cmd := newInitCmd()
		if err := cmd.Flags().Set("frameworks", "gin"); err != nil {
			t.Fatal(err)
		}
		flags, err := parseInitFlags(cmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(flags.frameworkFlags) != 1 {
			t.Fatalf("frameworkFlags = %v, want 1 element", flags.frameworkFlags)
		}
		if !flags.cliProvided {
			t.Error("cliProvided should be true when frameworks are set")
		}
	})

	t.Run("dot_target_does_not_set_create_dir", func(t *testing.T) {
		cmd := newInitCmd()
		flags, err := parseInitFlags(cmd, []string{"."})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if flags.createDir {
			t.Error("createDir should be false for '.' target")
		}
	})
}

func TestValidateInitTarget(t *testing.T) {
	t.Run("valid_directory", func(t *testing.T) {
		dir := t.TempDir()
		flags := &initFlags{absTargetDir: dir}
		if err := validateInitTarget(flags); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("samuel_repository_rejected", func(t *testing.T) {
		dir := t.TempDir()
		templateDir := filepath.Join(dir, "template")
		if err := os.MkdirAll(templateDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(templateDir, "CLAUDE.md"), []byte("# test"), 0644); err != nil {
			t.Fatal(err)
		}
		flags := &initFlags{absTargetDir: dir}
		err := validateInitTarget(flags)
		if err == nil {
			t.Error("expected error for Samuel repository")
		}
	})

	t.Run("config_exists_without_force", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "samuel.yaml"), []byte("version: 1.0"), 0644); err != nil {
			t.Fatal(err)
		}
		flags := &initFlags{absTargetDir: dir, force: false}
		err := validateInitTarget(flags)
		if err == nil {
			t.Error("expected error when config exists without force")
		}
	})

	t.Run("config_exists_with_force", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "samuel.yaml"), []byte("version: 1.0"), 0644); err != nil {
			t.Fatal(err)
		}
		flags := &initFlags{absTargetDir: dir, force: true}
		if err := validateInitTarget(flags); err != nil {
			t.Errorf("unexpected error with force: %v", err)
		}
	})

	t.Run("alt_config_exists_without_force", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, ".samuel.yaml"), []byte("version: 1.0"), 0644); err != nil {
			t.Fatal(err)
		}
		flags := &initFlags{absTargetDir: dir, force: false}
		err := validateInitTarget(flags)
		if err == nil {
			t.Error("expected error when alt config exists without force")
		}
	})
}

func TestSelectComponents_NonInteractive(t *testing.T) {
	t.Run("template_resolves_to_components", func(t *testing.T) {
		flags := &initFlags{
			nonInteractive: true,
			templateName:   "starter",
			cliProvided:    true,
		}
		sel, err := selectComponents(flags)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sel.template == nil {
			t.Fatal("template should not be nil")
		}
		if sel.template.Name != "starter" {
			t.Errorf("template name = %q, want %q", sel.template.Name, "starter")
		}
		// starter has typescript, python, go
		if len(sel.languages) != 3 {
			t.Errorf("got %d languages, want 3", len(sel.languages))
		}
	})

	t.Run("unknown_template_returns_error", func(t *testing.T) {
		flags := &initFlags{
			nonInteractive: true,
			templateName:   "nonexistent",
			cliProvided:    true,
		}
		_, err := selectComponents(flags)
		if err == nil {
			t.Error("expected error for unknown template")
		}
	})

	t.Run("language_flags_override_template", func(t *testing.T) {
		flags := &initFlags{
			nonInteractive: true,
			templateName:   "full",
			languageFlags:  []string{"go", "rust"},
			cliProvided:    true,
		}
		sel, err := selectComponents(flags)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sel.languages) != 2 {
			t.Fatalf("got %d languages, want 2: %v", len(sel.languages), sel.languages)
		}
		expected := map[string]bool{"go": true, "rust": true}
		for _, lang := range sel.languages {
			if !expected[lang] {
				t.Errorf("unexpected language %q", lang)
			}
		}
	})

	t.Run("framework_flags_override_template", func(t *testing.T) {
		flags := &initFlags{
			nonInteractive: true,
			templateName:   "full",
			frameworkFlags: []string{"gin", "echo"},
			cliProvided:    true,
		}
		sel, err := selectComponents(flags)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sel.frameworks) != 2 {
			t.Fatalf("got %d frameworks, want 2: %v", len(sel.frameworks), sel.frameworks)
		}
		expected := map[string]bool{"gin": true, "echo": true}
		for _, fw := range sel.frameworks {
			if !expected[fw] {
				t.Errorf("unexpected framework %q", fw)
			}
		}
	})

	t.Run("defaults_to_starter_when_nothing_selected", func(t *testing.T) {
		flags := &initFlags{
			nonInteractive: true,
			cliProvided:    false,
		}
		sel, err := selectComponents(flags)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sel.template == nil {
			t.Fatal("template should not be nil")
		}
		if sel.template.Name != "starter" {
			t.Errorf("template name = %q, want %q", sel.template.Name, "starter")
		}
	})

	t.Run("minimal_template_empty_languages", func(t *testing.T) {
		flags := &initFlags{
			nonInteractive: true,
			templateName:   "minimal",
			cliProvided:    true,
		}
		sel, err := selectComponents(flags)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sel.languages) != 0 {
			t.Errorf("minimal template should have 0 languages, got %d", len(sel.languages))
		}
		if len(sel.frameworks) != 0 {
			t.Errorf("minimal template should have 0 frameworks, got %d", len(sel.frameworks))
		}
	})

	t.Run("language_only_flags_no_template", func(t *testing.T) {
		flags := &initFlags{
			nonInteractive: true,
			languageFlags:  []string{"python"},
			cliProvided:    true,
		}
		sel, err := selectComponents(flags)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// No template set, languages from flags
		if sel.template != nil {
			t.Errorf("template should be nil when only languages provided, got %q", sel.template.Name)
		}
		if len(sel.languages) != 1 || sel.languages[0] != "python" {
			t.Errorf("languages = %v, want [python]", sel.languages)
		}
	})
}

func TestUpdateSkillsAndAgentsMD(t *testing.T) {
	t.Run("creates_agents_md_from_claude_md", func(t *testing.T) {
		dir := t.TempDir()
		claudeContent := "# CLAUDE.md\n\nTest content\n"
		claudeMDPath := filepath.Join(dir, "CLAUDE.md")
		if err := os.WriteFile(claudeMDPath, []byte(claudeContent), 0644); err != nil {
			t.Fatal(err)
		}
		// No skills directory — should still copy CLAUDE.md to AGENTS.md
		updateSkillsAndAgentsMD(dir)

		agentsContent, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
		if err != nil {
			t.Fatalf("AGENTS.md not created: %v", err)
		}
		if string(agentsContent) != claudeContent {
			t.Errorf("AGENTS.md content = %q, want %q", string(agentsContent), claudeContent)
		}
	})

	t.Run("with_skills_directory", func(t *testing.T) {
		dir := t.TempDir()
		skillsDir := filepath.Join(dir, ".claude", "skills", "test-skill")
		if err := os.MkdirAll(skillsDir, 0755); err != nil {
			t.Fatal(err)
		}
		skillContent := `---
name: test-skill
description: A test skill
compatibility: ["claude"]
---
# Test Skill

Instructions here.
`
		if err := os.WriteFile(filepath.Join(skillsDir, "SKILL.md"), []byte(skillContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Create CLAUDE.md with skills markers
		claudeContent := `# CLAUDE.md

<!-- SKILLS_START -->
## Available Skills
<!-- SKILLS_END -->
`
		if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte(claudeContent), 0644); err != nil {
			t.Fatal(err)
		}

		skills := updateSkillsAndAgentsMD(dir)
		if len(skills) == 0 {
			t.Error("expected at least 1 skill to be found")
		}

		// Verify AGENTS.md was created
		if _, err := os.Stat(filepath.Join(dir, "AGENTS.md")); err != nil {
			t.Errorf("AGENTS.md should exist: %v", err)
		}
	})

	t.Run("without_claude_md", func(t *testing.T) {
		dir := t.TempDir()
		// No CLAUDE.md — should not create AGENTS.md
		skills := updateSkillsAndAgentsMD(dir)
		if len(skills) != 0 {
			t.Errorf("expected 0 skills, got %d", len(skills))
		}
		if _, err := os.Stat(filepath.Join(dir, "AGENTS.md")); err == nil {
			t.Error("AGENTS.md should not be created without CLAUDE.md")
		}
	})

	t.Run("returns_discovered_skills", func(t *testing.T) {
		dir := t.TempDir()
		skillsDir := filepath.Join(dir, ".claude", "skills")
		for _, name := range []string{"skill-a", "skill-b"} {
			sDir := filepath.Join(skillsDir, name)
			if err := os.MkdirAll(sDir, 0755); err != nil {
				t.Fatal(err)
			}
			content := "---\nname: " + name + "\ndescription: Test\ncompatibility: [\"claude\"]\n---\n# " + name + "\n"
			if err := os.WriteFile(filepath.Join(sDir, "SKILL.md"), []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# Test"), 0644); err != nil {
			t.Fatal(err)
		}
		skills := updateSkillsAndAgentsMD(dir)
		if len(skills) != 2 {
			t.Errorf("expected 2 skills, got %d", len(skills))
		}
	})
}

func TestInstallAndSetup_CreateDir(t *testing.T) {
	t.Run("creates_directory_when_flagged", func(t *testing.T) {
		parent := t.TempDir()
		newDir := filepath.Join(parent, "my-project")
		flags := &initFlags{
			absTargetDir: newDir,
			createDir:    true,
		}
		sel := &initSelections{
			languages:  []string{},
			frameworks: []string{},
		}

		// installAndSetup will fail at the extractor stage since there's
		// no cached download, but the directory creation happens first
		_ = installAndSetup(flags, sel, "1.0.0", filepath.Join(parent, "nonexistent-cache"))

		// The directory should have been created
		info, err := os.Stat(newDir)
		if err != nil {
			t.Fatalf("directory was not created: %v", err)
		}
		if !info.IsDir() {
			t.Error("created path is not a directory")
		}
	})
}

func TestInitSelections_Struct(t *testing.T) {
	t.Run("zero_value_is_usable", func(t *testing.T) {
		sel := &initSelections{}
		if sel.template != nil {
			t.Error("template should be nil by default")
		}
		if sel.languages != nil {
			t.Error("languages should be nil by default")
		}
		if sel.frameworks != nil {
			t.Error("frameworks should be nil by default")
		}
	})

	t.Run("populated_correctly", func(t *testing.T) {
		tmpl := core.FindTemplate("starter")
		sel := &initSelections{
			template:   tmpl,
			languages:  []string{"go", "python"},
			frameworks: []string{"gin"},
		}
		if sel.template.Name != "starter" {
			t.Errorf("template name = %q, want %q", sel.template.Name, "starter")
		}
		if len(sel.languages) != 2 {
			t.Errorf("languages count = %d, want 2", len(sel.languages))
		}
		if len(sel.frameworks) != 1 {
			t.Errorf("frameworks count = %d, want 1", len(sel.frameworks))
		}
	})
}
