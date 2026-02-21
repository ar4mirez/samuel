package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestToTitleCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello-world", "Hello World"},
		{"go-guide", "Go Guide"},
		{"single", "Single"},
		{"a-b-c", "A B C"},
		{"", ""},
		{"already", "Already"},
		{"create-prd", "Create Prd"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toTitleCase(tt.input)
			if got != tt.want {
				t.Errorf("toTitleCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateSkillName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantCount int
		wantMsg   string // substring to match in first error
	}{
		{"valid_simple", "my-skill", 0, ""},
		{"valid_digits", "skill-123", 0, ""},
		{"valid_single", "skill", 0, ""},
		{"empty", "", 1, "name is required"},
		{"uppercase", "My-Skill", 2, "name must be lowercase"},
		{"too_long", strings.Repeat("a", 65), 1, "exceeds 64 character limit"},
		{"starts_with_hyphen", "-skill", 1, "cannot start or end with a hyphen"},
		{"ends_with_hyphen", "skill-", 1, "cannot start or end with a hyphen"},
		{"consecutive_hyphens", "my--skill", 1, "cannot contain consecutive hyphens"},
		{"invalid_chars", "my_skill", 1, "only contain lowercase letters"},
		{"spaces", "my skill", 1, "only contain lowercase letters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateSkillName(tt.input)
			if len(errs) != tt.wantCount {
				t.Errorf("ValidateSkillName(%q) returned %d errors, want %d: %v",
					tt.input, len(errs), tt.wantCount, errs)
				return
			}
			if tt.wantMsg != "" && len(errs) > 0 {
				if !strings.Contains(errs[0], tt.wantMsg) {
					t.Errorf("ValidateSkillName(%q) error = %q, want containing %q",
						tt.input, errs[0], tt.wantMsg)
				}
			}
		})
	}
}

func TestValidateSkillDescription(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantCount int
		wantMsg   string
	}{
		{"valid", "A useful skill description.", 0, ""},
		{"empty", "", 1, "description is required"},
		{"whitespace_only", "   ", 1, "description is required"},
		{"too_long", strings.Repeat("x", 1025), 1, "exceeds 1024 character limit"},
		{"at_limit", strings.Repeat("x", 1024), 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateSkillDescription(tt.input)
			if len(errs) != tt.wantCount {
				t.Errorf("ValidateSkillDescription(%q) returned %d errors, want %d: %v",
					tt.input, len(errs), tt.wantCount, errs)
				return
			}
			if tt.wantMsg != "" && len(errs) > 0 {
				if !strings.Contains(errs[0], tt.wantMsg) {
					t.Errorf("error = %q, want containing %q", errs[0], tt.wantMsg)
				}
			}
		})
	}
}

func TestValidateSkillCompatibility(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantCount int
	}{
		{"empty_is_valid", "", 0},
		{"valid", "Works with Claude, Cursor, Copilot", 0},
		{"at_limit", strings.Repeat("x", 500), 0},
		{"too_long", strings.Repeat("x", 501), 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateSkillCompatibility(tt.input)
			if len(errs) != tt.wantCount {
				t.Errorf("ValidateSkillCompatibility returned %d errors, want %d: %v",
					len(errs), tt.wantCount, errs)
			}
		})
	}
}

func TestValidateSkillMetadata(t *testing.T) {
	tests := []struct {
		name      string
		meta      SkillMetadata
		dirName   string
		wantCount int
		wantMsg   string
	}{
		{
			name: "valid_matching_dir",
			meta: SkillMetadata{
				Name:        "my-skill",
				Description: "A test skill.",
			},
			dirName:   "my-skill",
			wantCount: 0,
		},
		{
			name: "name_dir_mismatch",
			meta: SkillMetadata{
				Name:        "my-skill",
				Description: "A test skill.",
			},
			dirName:   "other-name",
			wantCount: 1,
			wantMsg:   "must match directory name",
		},
		{
			name: "empty_dir_skips_match",
			meta: SkillMetadata{
				Name:        "my-skill",
				Description: "A test skill.",
			},
			dirName:   "",
			wantCount: 0,
		},
		{
			name: "multiple_errors",
			meta: SkillMetadata{
				Name:        "",
				Description: "",
			},
			dirName:   "test",
			wantCount: 2, // name required + description required
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateSkillMetadata(tt.meta, tt.dirName)
			if len(errs) != tt.wantCount {
				t.Errorf("ValidateSkillMetadata returned %d errors, want %d: %v",
					len(errs), tt.wantCount, errs)
				return
			}
			if tt.wantMsg != "" && len(errs) > 0 {
				found := false
				for _, e := range errs {
					if strings.Contains(e, tt.wantMsg) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error containing %q, got %v", tt.wantMsg, errs)
				}
			}
		})
	}
}

func TestParseSkillMD(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantName string
		wantBody string
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid_full",
			content: `---
name: my-skill
description: A skill description.
license: MIT
---

# My Skill

Body content here.`,
			wantName: "my-skill",
			wantBody: "# My Skill\n\nBody content here.",
		},
		{
			name: "valid_no_body",
			content: `---
name: minimal
description: Minimal skill.
---`,
			wantName: "minimal",
			wantBody: "",
		},
		{
			name:    "empty_content",
			content: "",
			wantErr: true,
			errMsg:  "must start with YAML frontmatter",
		},
		{
			name:    "no_frontmatter",
			content: "# Just a heading",
			wantErr: true,
			errMsg:  "must start with YAML frontmatter",
		},
		{
			name: "unclosed_frontmatter",
			content: `---
name: broken
description: No closing delimiter.
`,
			wantErr: true,
			errMsg:  "frontmatter not closed",
		},
		{
			name: "invalid_yaml",
			content: `---
name: [invalid
---`,
			wantErr: true,
			errMsg:  "invalid YAML frontmatter",
		},
		{
			name: "with_metadata",
			content: `---
name: test-skill
description: Test.
metadata:
  author: tester
  version: "2.0"
---

Content.`,
			wantName: "test-skill",
			wantBody: "Content.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta, body, err := ParseSkillMD(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSkillMD error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("error = %q, want containing %q", err.Error(), tt.errMsg)
				}
				return
			}
			if meta.Name != tt.wantName {
				t.Errorf("name = %q, want %q", meta.Name, tt.wantName)
			}
			if body != tt.wantBody {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
		})
	}
}

func TestLoadSkillInfo(t *testing.T) {
	t.Run("valid_skill", func(t *testing.T) {
		dir := t.TempDir()
		skillDir := filepath.Join(dir, "my-skill")
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			t.Fatal(err)
		}

		content := `---
name: my-skill
description: A test skill.
---

Body.`
		if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		info, err := LoadSkillInfo(skillDir)
		if err != nil {
			t.Fatalf("LoadSkillInfo error: %v", err)
		}
		if info.DirName != "my-skill" {
			t.Errorf("DirName = %q, want %q", info.DirName, "my-skill")
		}
		if info.Metadata.Name != "my-skill" {
			t.Errorf("Name = %q, want %q", info.Metadata.Name, "my-skill")
		}
		if len(info.Errors) != 0 {
			t.Errorf("unexpected errors: %v", info.Errors)
		}
		if info.Body != "Body." {
			t.Errorf("Body = %q, want %q", info.Body, "Body.")
		}
	})

	t.Run("missing_skill_md", func(t *testing.T) {
		dir := t.TempDir()
		skillDir := filepath.Join(dir, "empty-skill")
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			t.Fatal(err)
		}

		info, err := LoadSkillInfo(skillDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(info.Errors) == 0 {
			t.Error("expected errors for missing SKILL.md")
		}
		if !strings.Contains(info.Errors[0], "missing required file") {
			t.Errorf("error = %q, want containing 'missing required file'", info.Errors[0])
		}
	})

	t.Run("invalid_frontmatter", func(t *testing.T) {
		dir := t.TempDir()
		skillDir := filepath.Join(dir, "bad-skill")
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			t.Fatal(err)
		}

		content := "# No frontmatter"
		if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		info, err := LoadSkillInfo(skillDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(info.Errors) == 0 {
			t.Error("expected errors for invalid frontmatter")
		}
	})

	t.Run("name_dir_mismatch", func(t *testing.T) {
		dir := t.TempDir()
		skillDir := filepath.Join(dir, "actual-dir")
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			t.Fatal(err)
		}

		content := `---
name: different-name
description: Mismatched name.
---`
		if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		info, err := LoadSkillInfo(skillDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(info.Errors) == 0 {
			t.Error("expected error for name/dir mismatch")
		}
	})

	t.Run("with_optional_dirs", func(t *testing.T) {
		dir := t.TempDir()
		skillDir := filepath.Join(dir, "full-skill")
		for _, sub := range []string{"scripts", "references", "assets"} {
			if err := os.MkdirAll(filepath.Join(skillDir, sub), 0755); err != nil {
				t.Fatal(err)
			}
		}

		content := `---
name: full-skill
description: Has all optional dirs.
---`
		if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		info, err := LoadSkillInfo(skillDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !info.HasScripts {
			t.Error("HasScripts should be true")
		}
		if !info.HasRefs {
			t.Error("HasRefs should be true")
		}
		if !info.HasAssets {
			t.Error("HasAssets should be true")
		}
	})
}

func TestScanSkillsDirectory(t *testing.T) {
	t.Run("multiple_skills", func(t *testing.T) {
		dir := t.TempDir()

		// Create two valid skills
		for _, name := range []string{"skill-a", "skill-b"} {
			skillDir := filepath.Join(dir, name)
			if err := os.MkdirAll(skillDir, 0755); err != nil {
				t.Fatal(err)
			}
			content := "---\nname: " + name + "\ndescription: Test.\n---\n"
			if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
		}

		// Create a directory without SKILL.md (should be skipped)
		if err := os.MkdirAll(filepath.Join(dir, "no-skill"), 0755); err != nil {
			t.Fatal(err)
		}

		// Create a hidden directory (should be skipped)
		hiddenDir := filepath.Join(dir, ".hidden")
		if err := os.MkdirAll(hiddenDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(hiddenDir, "SKILL.md"), []byte("---\nname: hidden\ndescription: x\n---\n"), 0644); err != nil {
			t.Fatal(err)
		}

		// Create a regular file (should be skipped)
		if err := os.WriteFile(filepath.Join(dir, "not-a-dir.txt"), []byte("file"), 0644); err != nil {
			t.Fatal(err)
		}

		skills, err := ScanSkillsDirectory(dir)
		if err != nil {
			t.Fatalf("ScanSkillsDirectory error: %v", err)
		}
		if len(skills) != 2 {
			t.Errorf("got %d skills, want 2", len(skills))
		}
	})

	t.Run("nonexistent_dir", func(t *testing.T) {
		skills, err := ScanSkillsDirectory("/nonexistent/path")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(skills) != 0 {
			t.Errorf("expected empty slice, got %d skills", len(skills))
		}
	})

	t.Run("empty_dir", func(t *testing.T) {
		dir := t.TempDir()
		skills, err := ScanSkillsDirectory(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(skills) != 0 {
			t.Errorf("expected empty slice, got %d skills", len(skills))
		}
	})
}

func TestGenerateSkillsSection(t *testing.T) {
	t.Run("empty_skills", func(t *testing.T) {
		result := GenerateSkillsSection(nil)
		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})

	t.Run("with_valid_skills", func(t *testing.T) {
		skills := []*SkillInfo{
			{
				Metadata: SkillMetadata{
					Name:        "go-guide",
					Description: "Go language patterns.",
				},
			},
			{
				Metadata: SkillMetadata{
					Name:        "react",
					Description: "React framework guide.",
				},
			},
		}

		result := GenerateSkillsSection(skills)
		if !strings.Contains(result, "go-guide") {
			t.Error("expected output to contain 'go-guide'")
		}
		if !strings.Contains(result, "react") {
			t.Error("expected output to contain 'react'")
		}
		if !strings.Contains(result, "## Available Skills") {
			t.Error("expected header")
		}
		if !strings.Contains(result, "| Skill | Description |") {
			t.Error("expected table header")
		}
	})

	t.Run("skips_invalid_skills", func(t *testing.T) {
		skills := []*SkillInfo{
			{
				Metadata: SkillMetadata{Name: "valid", Description: "OK."},
			},
			{
				Metadata: SkillMetadata{Name: "invalid", Description: "Bad."},
				Errors:   []string{"some error"},
			},
		}

		result := GenerateSkillsSection(skills)
		if !strings.Contains(result, "valid") {
			t.Error("expected valid skill in output")
		}
		if strings.Contains(result, "| invalid |") {
			t.Error("invalid skill should be skipped")
		}
	})

	t.Run("truncates_long_description", func(t *testing.T) {
		longDesc := strings.Repeat("x", 100)
		skills := []*SkillInfo{
			{
				Metadata: SkillMetadata{
					Name:        "long-desc",
					Description: longDesc,
				},
			},
		}

		result := GenerateSkillsSection(skills)
		if strings.Contains(result, longDesc) {
			t.Error("long description should be truncated")
		}
		if !strings.Contains(result, "...") {
			t.Error("truncated description should end with '...'")
		}
	})
}

func TestGetSkillTemplate(t *testing.T) {
	template := GetSkillTemplate("my-skill")

	if !strings.Contains(template, "name: my-skill") {
		t.Error("template should contain the skill name in frontmatter")
	}
	if !strings.Contains(template, "# My Skill") {
		t.Error("template should contain title-cased heading")
	}
	if !strings.Contains(template, "---") {
		t.Error("template should have frontmatter delimiters")
	}
	if !strings.Contains(template, "description:") {
		t.Error("template should have description field")
	}
}

func TestCreateSkillScaffold(t *testing.T) {
	t.Run("creates_structure", func(t *testing.T) {
		dir := t.TempDir()

		err := CreateSkillScaffold(dir, "test-skill")
		if err != nil {
			t.Fatalf("CreateSkillScaffold error: %v", err)
		}

		// Check SKILL.md exists
		skillMD := filepath.Join(dir, "test-skill", "SKILL.md")
		if _, err := os.Stat(skillMD); os.IsNotExist(err) {
			t.Error("SKILL.md should exist")
		}

		// Check optional directories
		for _, sub := range []string{"scripts", "references", "assets"} {
			subDir := filepath.Join(dir, "test-skill", sub)
			if _, err := os.Stat(subDir); os.IsNotExist(err) {
				t.Errorf("%s directory should exist", sub)
			}
			gitkeep := filepath.Join(subDir, ".gitkeep")
			if _, err := os.Stat(gitkeep); os.IsNotExist(err) {
				t.Errorf("%s/.gitkeep should exist", sub)
			}
		}

		// Verify SKILL.md content
		content, err := os.ReadFile(skillMD)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(content), "name: test-skill") {
			t.Error("SKILL.md should contain skill name")
		}
	})

	t.Run("already_exists", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, "existing"), 0755); err != nil {
			t.Fatal(err)
		}

		err := CreateSkillScaffold(dir, "existing")
		if err == nil {
			t.Error("expected error for existing skill")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("error = %q, want containing 'already exists'", err)
		}
	})
}

func TestUpdateCLAUDEMDSkillsSection(t *testing.T) {
	t.Run("replaces_existing_section", func(t *testing.T) {
		dir := t.TempDir()
		claudeMD := filepath.Join(dir, "CLAUDE.md")

		original := `# Project

Some content.

<!-- SKILLS_START -->
## Old Skills
Old content.
<!-- SKILLS_END -->

More content.`
		if err := os.WriteFile(claudeMD, []byte(original), 0644); err != nil {
			t.Fatal(err)
		}

		skills := []*SkillInfo{
			{
				Metadata: SkillMetadata{
					Name:        "new-skill",
					Description: "New skill description.",
				},
			},
		}

		err := UpdateCLAUDEMDSkillsSection(claudeMD, skills)
		if err != nil {
			t.Fatalf("UpdateCLAUDEMDSkillsSection error: %v", err)
		}

		content, err := os.ReadFile(claudeMD)
		if err != nil {
			t.Fatal(err)
		}
		contentStr := string(content)

		if !strings.Contains(contentStr, "new-skill") {
			t.Error("updated content should contain new skill")
		}
		if strings.Contains(contentStr, "Old Skills") {
			t.Error("old skills section should be replaced")
		}
		if !strings.Contains(contentStr, "More content.") {
			t.Error("content after markers should be preserved")
		}
		if !strings.Contains(contentStr, "Some content.") {
			t.Error("content before markers should be preserved")
		}
	})

	t.Run("no_markers_does_nothing", func(t *testing.T) {
		dir := t.TempDir()
		claudeMD := filepath.Join(dir, "CLAUDE.md")

		original := "# Project\n\nNo markers here."
		if err := os.WriteFile(claudeMD, []byte(original), 0644); err != nil {
			t.Fatal(err)
		}

		skills := []*SkillInfo{
			{
				Metadata: SkillMetadata{
					Name:        "skill",
					Description: "Desc.",
				},
			},
		}

		err := UpdateCLAUDEMDSkillsSection(claudeMD, skills)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		content, err := os.ReadFile(claudeMD)
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != original {
			t.Error("content should remain unchanged without markers")
		}
	})

	t.Run("empty_skills_does_nothing", func(t *testing.T) {
		dir := t.TempDir()
		claudeMD := filepath.Join(dir, "CLAUDE.md")

		original := "<!-- SKILLS_START -->\n<!-- SKILLS_END -->"
		if err := os.WriteFile(claudeMD, []byte(original), 0644); err != nil {
			t.Fatal(err)
		}

		err := UpdateCLAUDEMDSkillsSection(claudeMD, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		content, err := os.ReadFile(claudeMD)
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != original {
			t.Error("content should remain unchanged with empty skills")
		}
	})

	t.Run("missing_file", func(t *testing.T) {
		err := UpdateCLAUDEMDSkillsSection("/nonexistent/CLAUDE.md", []*SkillInfo{
			{Metadata: SkillMetadata{Name: "x", Description: "y"}},
		})
		if err == nil {
			t.Error("expected error for missing file")
		}
	})
}

func TestDirExists(t *testing.T) {
	t.Run("existing_dir", func(t *testing.T) {
		dir := t.TempDir()
		if !dirExists(dir) {
			t.Error("dirExists should return true for existing directory")
		}
	})

	t.Run("nonexistent", func(t *testing.T) {
		if dirExists("/nonexistent/path") {
			t.Error("dirExists should return false for nonexistent path")
		}
	})

	t.Run("file_not_dir", func(t *testing.T) {
		dir := t.TempDir()
		f := filepath.Join(dir, "file.txt")
		if err := os.WriteFile(f, []byte("hi"), 0644); err != nil {
			t.Fatal(err)
		}
		if dirExists(f) {
			t.Error("dirExists should return false for a file")
		}
	})
}

func TestCountLines(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"empty", "", 0},
		{"single_line", "hello", 1},
		{"two_lines", "hello\nworld", 2},
		{"trailing_newline", "hello\n", 1},
		{"multiple", "a\nb\nc\nd", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountLines(tt.input)
			if got != tt.want {
				t.Errorf("CountLines(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
