package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

// Skill validation constants per Agent Skills specification
const (
	MaxSkillNameLength     = 64
	MaxDescriptionLength   = 1024
	MaxCompatibilityLength = 500
)

// SkillMetadata represents SKILL.md frontmatter per Agent Skills spec
type SkillMetadata struct {
	Name          string            `yaml:"name"`
	Description   string            `yaml:"description"`
	License       string            `yaml:"license,omitempty"`
	Compatibility string            `yaml:"compatibility,omitempty"`
	AllowedTools  string            `yaml:"allowed-tools,omitempty"`
	Metadata      map[string]string `yaml:"metadata,omitempty"`
}

// SkillInfo contains parsed skill information
type SkillInfo struct {
	Path        string
	DirName     string
	Metadata    SkillMetadata
	Body        string
	HasScripts  bool
	HasRefs     bool
	HasAssets   bool
	Errors      []string
}

// toTitleCase converts a kebab-case name to Title Case
func toTitleCase(s string) string {
	words := strings.Split(s, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

// ValidateSkillName checks name format per Agent Skills specification
func ValidateSkillName(name string) []string {
	var errors []string

	if name == "" {
		errors = append(errors, "name is required")
		return errors
	}

	if len(name) > MaxSkillNameLength {
		errors = append(errors, fmt.Sprintf("name exceeds %d character limit (%d chars)", MaxSkillNameLength, len(name)))
	}

	// Check for uppercase characters
	if name != strings.ToLower(name) {
		errors = append(errors, "name must be lowercase")
	}

	// Check for valid characters
	for _, r := range name {
		if !unicode.IsLower(r) && !unicode.IsDigit(r) && r != '-' {
			errors = append(errors, "name may only contain lowercase letters, digits, and hyphens")
			break
		}
	}

	// Check start/end with hyphen
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		errors = append(errors, "name cannot start or end with a hyphen")
	}

	// Check for consecutive hyphens
	if strings.Contains(name, "--") {
		errors = append(errors, "name cannot contain consecutive hyphens")
	}

	return errors
}

// ValidateSkillDescription checks description per Agent Skills specification
func ValidateSkillDescription(description string) []string {
	var errors []string

	if strings.TrimSpace(description) == "" {
		errors = append(errors, "description is required")
		return errors
	}

	if len(description) > MaxDescriptionLength {
		errors = append(errors, fmt.Sprintf("description exceeds %d character limit (%d chars)", MaxDescriptionLength, len(description)))
	}

	return errors
}

// ValidateSkillCompatibility checks compatibility field per Agent Skills specification
func ValidateSkillCompatibility(compatibility string) []string {
	var errors []string

	if compatibility != "" && len(compatibility) > MaxCompatibilityLength {
		errors = append(errors, fmt.Sprintf("compatibility exceeds %d character limit (%d chars)", MaxCompatibilityLength, len(compatibility)))
	}

	return errors
}

// ValidateSkillMetadata validates the complete frontmatter
func ValidateSkillMetadata(meta SkillMetadata, dirName string) []string {
	var errors []string

	// Validate name
	errors = append(errors, ValidateSkillName(meta.Name)...)

	// Check name matches directory
	if meta.Name != "" && dirName != "" && meta.Name != dirName {
		errors = append(errors, fmt.Sprintf("skill name '%s' must match directory name '%s'", meta.Name, dirName))
	}

	// Validate description
	errors = append(errors, ValidateSkillDescription(meta.Description)...)

	// Validate compatibility (optional)
	errors = append(errors, ValidateSkillCompatibility(meta.Compatibility)...)

	return errors
}

// ParseSkillMD parses SKILL.md content and extracts frontmatter and body
func ParseSkillMD(content string) (*SkillMetadata, string, error) {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return nil, "", fmt.Errorf("empty SKILL.md file")
	}

	// Check for frontmatter delimiter
	if strings.TrimSpace(lines[0]) != "---" {
		return nil, "", fmt.Errorf("SKILL.md must start with YAML frontmatter (---)")
	}

	// Find closing delimiter
	var frontmatterEnd int
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			frontmatterEnd = i
			break
		}
	}

	if frontmatterEnd == 0 {
		return nil, "", fmt.Errorf("SKILL.md frontmatter not closed (missing ---)")
	}

	// Extract frontmatter
	frontmatterContent := strings.Join(lines[1:frontmatterEnd], "\n")

	var meta SkillMetadata
	if err := yaml.Unmarshal([]byte(frontmatterContent), &meta); err != nil {
		return nil, "", fmt.Errorf("invalid YAML frontmatter: %w", err)
	}

	// Extract body
	body := ""
	if frontmatterEnd+1 < len(lines) {
		body = strings.Join(lines[frontmatterEnd+1:], "\n")
	}

	return &meta, strings.TrimSpace(body), nil
}

// LoadSkillInfo loads and validates a skill from a directory
func LoadSkillInfo(skillDir string) (*SkillInfo, error) {
	info := &SkillInfo{
		Path:    skillDir,
		DirName: filepath.Base(skillDir),
	}

	// Check SKILL.md exists
	skillMDPath := filepath.Join(skillDir, "SKILL.md")
	content, err := os.ReadFile(skillMDPath)
	if err != nil {
		if os.IsNotExist(err) {
			info.Errors = append(info.Errors, "missing required file: SKILL.md")
			return info, nil
		}
		return nil, fmt.Errorf("failed to read SKILL.md: %w", err)
	}

	// Parse SKILL.md
	meta, body, err := ParseSkillMD(string(content))
	if err != nil {
		info.Errors = append(info.Errors, err.Error())
		return info, nil
	}

	info.Metadata = *meta
	info.Body = body

	// Validate metadata
	info.Errors = append(info.Errors, ValidateSkillMetadata(*meta, info.DirName)...)

	// Check optional directories
	info.HasScripts = dirExists(filepath.Join(skillDir, "scripts"))
	info.HasRefs = dirExists(filepath.Join(skillDir, "references"))
	info.HasAssets = dirExists(filepath.Join(skillDir, "assets"))

	return info, nil
}

// ScanSkillsDirectory scans a directory for skills and returns their info
func ScanSkillsDirectory(skillsDir string) ([]*SkillInfo, error) {
	var skills []*SkillInfo

	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return skills, nil
		}
		return nil, fmt.Errorf("failed to read skills directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip hidden directories and special files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		skillPath := filepath.Join(skillsDir, entry.Name())

		// Check if SKILL.md exists
		if _, err := os.Stat(filepath.Join(skillPath, "SKILL.md")); os.IsNotExist(err) {
			continue
		}

		info, err := LoadSkillInfo(skillPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load skill '%s': %w", entry.Name(), err)
		}

		skills = append(skills, info)
	}

	return skills, nil
}

// GenerateSkillsSection generates the "Available Skills" markdown section
func GenerateSkillsSection(skills []*SkillInfo) string {
	if len(skills) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("## Available Skills\n\n")
	sb.WriteString("Skills extend AI capabilities. Load a skill when task matches its description.\n\n")
	sb.WriteString("| Skill | Description |\n")
	sb.WriteString("|-------|-------------|\n")

	for _, skill := range skills {
		if len(skill.Errors) > 0 {
			continue // Skip invalid skills
		}
		// Truncate description for table display
		desc := skill.Metadata.Description
		desc = strings.ReplaceAll(desc, "\n", " ")
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		fmt.Fprintf(&sb, "| %s | %s |\n", skill.Metadata.Name, desc)
	}

	sb.WriteString("\n**To use a skill**: Read `.claude/skills/<skill-name>/SKILL.md`\n")

	return sb.String()
}

// GetSkillTemplate returns the template content for a new SKILL.md file
func GetSkillTemplate(name string) string {
	return fmt.Sprintf(`---
name: %s
description: |
  Brief description of what this skill does and when to use it.
  Include specific triggers and keywords that should activate this skill.
license: MIT
metadata:
  author: your-name
  version: "1.0"
---

# %s

## Purpose

Describe what capability this skill provides to AI agents.

## When to Use

- Scenario 1: When the user asks for...
- Scenario 2: When working with...

## Instructions

Step-by-step instructions for the AI agent:

1. First, analyze the request
2. Then, perform the action
3. Finally, verify the result

## Examples

### Example 1: Basic Usage

**Input**: User request example

**Output**:
`+"```"+`
Expected output
`+"```"+`

## Notes

Any additional context, warnings, or best practices.
`, name, toTitleCase(name))
}

// CreateSkillScaffold creates a new skill directory with template files
func CreateSkillScaffold(skillsDir, name string) error {
	skillPath := filepath.Join(skillsDir, name)

	// Check if skill already exists
	if _, err := os.Stat(skillPath); !os.IsNotExist(err) {
		return fmt.Errorf("skill '%s' already exists", name)
	}

	// Create skill directory
	if err := os.MkdirAll(skillPath, 0755); err != nil {
		return fmt.Errorf("failed to create skill directory: %w", err)
	}

	// Create SKILL.md
	skillMDPath := filepath.Join(skillPath, "SKILL.md")
	if err := os.WriteFile(skillMDPath, []byte(GetSkillTemplate(name)), 0644); err != nil {
		return fmt.Errorf("failed to create SKILL.md: %w", err)
	}

	// Create optional directories
	for _, dir := range []string{"scripts", "references", "assets"} {
		dirPath := filepath.Join(skillPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", dir, err)
		}

		// Create .gitkeep
		gitkeepPath := filepath.Join(dirPath, ".gitkeep")
		if err := os.WriteFile(gitkeepPath, []byte(""), 0644); err != nil {
			return fmt.Errorf("failed to create .gitkeep: %w", err)
		}
	}

	return nil
}

// UpdateCLAUDEMDSkillsSection updates the skills section in CLAUDE.md
func UpdateCLAUDEMDSkillsSection(claudeMDPath string, skills []*SkillInfo) error {
	content, err := os.ReadFile(claudeMDPath)
	if err != nil {
		return fmt.Errorf("failed to read CLAUDE.md: %w", err)
	}

	skillsSection := GenerateSkillsSection(skills)
	if skillsSection == "" {
		return nil // No skills to add
	}

	contentStr := string(content)

	// Look for skills marker comments
	startMarker := "<!-- SKILLS_START -->"
	endMarker := "<!-- SKILLS_END -->"

	startIdx := strings.Index(contentStr, startMarker)
	endIdx := strings.Index(contentStr, endMarker)

	var newContent string
	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		// Replace existing section
		newContent = contentStr[:startIdx] +
			startMarker + "\n" +
			skillsSection +
			contentStr[endIdx:]
	} else {
		// Skills section doesn't exist, don't add it automatically
		// The user can add the markers manually if they want auto-updates
		return nil
	}

	return os.WriteFile(claudeMDPath, []byte(newContent), 0644)
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// CountLines counts lines in a string (for token estimation)
func CountLines(s string) int {
	scanner := bufio.NewScanner(strings.NewReader(s))
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}
