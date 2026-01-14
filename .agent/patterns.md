# AICoF Codebase Patterns

> **Purpose**: Document coding patterns and conventions established in this project
>
> **Created**: 2026-01-14
> **Last Updated**: 2026-01-14

---

## Go CLI Patterns

### CLI Flag Detection Pattern

**When to use**: When CLI needs different behavior for interactive vs flag-provided inputs

**Example**:

```go
func runCommand(cmd *cobra.Command, args []string) error {
    // Get flag values
    flagA, _ := cmd.Flags().GetString("flag-a")
    flagB, _ := cmd.Flags().GetStringSlice("flag-b")

    // Track if user provided CLI flags (skip prompts if so)
    cliProvided := flagA != "" || len(flagB) > 0

    // Only prompt when flags weren't provided
    if !nonInteractive && !cliProvided {
        // Interactive prompts here
    }
}
```

**Why**: Users expect explicit CLI flags to be honored without additional prompts. This pattern provides good UX for both interactive and scripted usage.

**See also**: [init.go](packages/cli/internal/cmd/init.go#L49-L53)

---

### Promptui Confirm Pattern

**When to use**: When prompting for yes/no confirmation

**Example**:

```go
func Confirm(label string, defaultYes bool) (bool, error) {
    suffix := " [y/N]"
    defaultStr := "n"
    if defaultYes {
        suffix = " [Y/n]"
        defaultStr = "y"
    }

    prompt := promptui.Prompt{
        Label:   label + suffix,
        Default: defaultStr,  // NOT using IsConfirm: true
    }

    result, err := prompt.Run()
    if err != nil {
        return false, err
    }

    result = strings.ToLower(strings.TrimSpace(result))
    if result == "" {
        return defaultYes, nil
    }
    return result == "y" || result == "yes", nil
}
```

**Why**: The `IsConfirm: true` option in promptui has quirky behavior where Enter on empty input returns `ErrAbort`. Using a standard prompt with default value gives predictable behavior.

**See also**: [prompts.go](packages/cli/internal/ui/prompts.go#L120-L148)

---

### Repository Safety Check Pattern

**When to use**: When a CLI tool operates on repositories and needs to prevent accidental self-modification

**Example**:

```go
// isAICOFRepository checks if the target directory is the AICoF repository itself
func isAICOFRepository(targetDir string) bool {
    // Check for unique directory structure
    templateDir := filepath.Join(targetDir, "template")
    if info, err := os.Stat(templateDir); err == nil && info.IsDir() {
        claudeMD := filepath.Join(templateDir, "CLAUDE.md")
        if _, err := os.Stat(claudeMD); err == nil {
            return true
        }
    }

    // Check for source code directory
    cliDir := filepath.Join(targetDir, "packages", "cli")
    if info, err := os.Stat(cliDir); err == nil && info.IsDir() {
        goMod := filepath.Join(cliDir, "go.mod")
        if _, err := os.Stat(goMod); err == nil {
            return true
        }
    }

    return false
}
```

**Why**: Prevents users from accidentally overwriting the framework source when running CLI commands from within the repository.

**See also**: [init.go](packages/cli/internal/cmd/init.go#L349-L371)

---

### Template Prefix Pattern

**When to use**: When distributing files from a repository that has separate template and source directories

**Example**:

```go
// TemplatePrefix is the path prefix where template files are stored
const TemplatePrefix = "template/"

// GetSourcePath returns the source path for a destination path
func GetSourcePath(destPath string) string {
    return TemplatePrefix + destPath
}

// In extraction:
srcPath := filepath.Join(e.sourcePath, TemplatePrefix, path)
dstPath := filepath.Join(e.destPath, path)
```

**Why**: Allows clear separation between distributable template files and project-specific source code in the same repository.

**See also**: [registry.go](packages/cli/internal/core/registry.go#L12-L14), [extractor.go](packages/cli/internal/core/extractor.go#L51-L52)

---

### Error Display Pattern (Cobra)

**When to use**: When Cobra's default error handling is silenced for cleaner output

**Example**:

```go
// In root.go
var rootCmd = &cobra.Command{
    SilenceErrors: true,  // Prevent duplicate error messages
    SilenceUsage:  true,  // Don't show usage on every error
}

// In main.go - handle errors explicitly
func main() {
    if err := cmd.Execute(); err != nil {
        red := color.New(color.FgRed).SprintFunc()
        fmt.Fprintf(os.Stderr, "%s %s\n", red("Error:"), err.Error())
        os.Exit(1)
    }
}
```

**Why**: Cobra's default error handling can be noisy. This pattern gives clean, colored error output while maintaining proper exit codes.

**See also**: [main.go](packages/cli/cmd/aicof/main.go#L11-L18), [root.go](packages/cli/internal/cmd/root.go)

---

## Documentation Patterns

### Workflow Document Structure

**When to use**: When creating new workflow documents in `.agent/workflows/`

**Example**:

```markdown
# Workflow Name

Brief description.

---

## When to Use

| Trigger | Action |
|---------|--------|
| ...     | ...    |

---

## Prerequisites

- [ ] Prerequisite 1
- [ ] Prerequisite 2

---

## Process Overview

```
1. Step One
   └── Details
         ↓
2. Step Two
   └── Details
```

---

## Phase 1: First Phase

### AI Will Review
...

---

## Usage Examples

### Example 1: Scenario Name
...
```

**Why**: Consistent structure makes workflows easier to follow and maintain.

**See also**: [document-work.md](template/.agent/workflows/document-work.md)

---

## Testing Patterns

### CLI E2E Test Setup

**When to use**: Testing CLI commands that interact with external services

**Example**:

```go
func TestInitCommand(t *testing.T) {
    // Create temp directory
    tmpDir := t.TempDir()

    // Run command
    rootCmd := NewRootCmd()
    rootCmd.SetArgs([]string{"init", "--template", "minimal", tmpDir})

    err := rootCmd.Execute()
    require.NoError(t, err)

    // Verify files created
    _, err = os.Stat(filepath.Join(tmpDir, "CLAUDE.md"))
    require.NoError(t, err)
}
```

**Why**: Temp directories are automatically cleaned up, and isolated tests prevent interference.

**See also**: Testing PRD in [.agent/tasks/0001-prd-cli-testing.md](.agent/tasks/0001-prd-cli-testing.md)
