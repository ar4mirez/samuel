# AICoF Codebase Patterns

> **Purpose**: Document coding patterns and conventions established in this project
>
> **Created**: 2026-01-14
> **Last Updated**: 2026-01-14

---

## Go CLI Patterns

### Standard Go Project Layout

**When to use**: When structuring a Go CLI application repository

**Example**:

```text
project-root/
├── cmd/
│   └── appname/
│       └── main.go          # Entry point, minimal code
├── internal/
│   ├── commands/            # CLI command implementations (not cmd/)
│   ├── core/                # Business logic
│   └── ui/                  # User interface helpers
├── go.mod                   # At repository root
├── go.sum
├── Makefile
└── .goreleaser.yaml
```

**Why**: This is the idiomatic Go project layout. Key conventions:

- `go.mod` at repository root
- `cmd/` contains entry points (main packages)
- `internal/` contains private packages (not importable by external code)
- Avoid naming internal packages `cmd` to prevent confusion with root `cmd/`
- Entry point should be minimal - just parse args and call into internal packages

**See also**: [Go project layout](https://github.com/golang-standards/project-layout), [memory/2026-01-14-go-project-restructure.md](.agent/memory/2026-01-14-go-project-restructure.md)

---

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

**See also**: [init.go](internal/commands/init.go#L49-L53)

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

**See also**: [prompts.go](internal/ui/prompts.go#L120-L148)

---

### Repository Safety Check Pattern

**When to use**: When a CLI tool operates on repositories and needs to prevent accidental self-modification

**Example**:

```go
// isAICOFRepository checks if the target directory is the AICoF repository itself
func isAICOFRepository(targetDir string) bool {
    // Check for unique directory structure (template/ with CLAUDE.md)
    templateDir := filepath.Join(targetDir, "template")
    if info, err := os.Stat(templateDir); err == nil && info.IsDir() {
        claudeMD := filepath.Join(templateDir, "CLAUDE.md")
        if _, err := os.Stat(claudeMD); err == nil {
            return true
        }
    }

    // Check for Go source code at root (cmd/aicof with main.go)
    cmdDir := filepath.Join(targetDir, "cmd", "aicof")
    if info, err := os.Stat(cmdDir); err == nil && info.IsDir() {
        mainGo := filepath.Join(cmdDir, "main.go")
        if _, err := os.Stat(mainGo); err == nil {
            return true
        }
    }

    return false
}
```

**Why**: Prevents users from accidentally overwriting the framework source when running CLI commands from within the repository.

**See also**: [init.go](internal/commands/init.go#L349-L371)

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

**See also**: [registry.go](internal/core/registry.go#L12-L14), [extractor.go](internal/core/extractor.go#L51-L52)

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

**See also**: [main.go](cmd/aicof/main.go#L11-L18), [root.go](internal/commands/root.go)

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

---

### Table-Driven Unit Tests

**When to use**: When testing functions with multiple input/output scenarios

**Example**:

```go
func TestMatchScore(t *testing.T) {
    tests := []struct {
        name        string
        query       string
        compName    string
        description string
        wantScore   int
    }{
        {"exact match", "react", "react", "React framework", 100},
        {"prefix match", "type", "typescript", "TypeScript", 80},
        {"contains match", "script", "typescript", "TypeScript", 60},
        {"no match", "xyz123", "react", "React framework", 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := matchScore(tt.query, tt.compName, tt.description)
            if got != tt.wantScore {
                t.Errorf("matchScore(%q, %q, %q) = %d, want %d",
                    tt.query, tt.compName, tt.description, got, tt.wantScore)
            }
        })
    }
}
```

**Why**: Table-driven tests provide:

- Clear test case visibility
- Easy to add new cases
- Consistent error messages
- Subtests with `t.Run()` for granular failure reporting

**See also**: [search_test.go](internal/commands/search_test.go), [config_test.go](internal/core/config_test.go)

---

### Coverage-Focused Testing Strategy

**When to use**: When implementing tests for existing code

**Strategy**:

1. **Test helper functions first** - Business logic in helper functions is most valuable to test
2. **Skip command runners** - `run*` functions require integration tests
3. **Skip display functions** - UI output functions have low testing value
4. **Target 80%+ for business logic** - Even if overall coverage is lower

**Example coverage breakdown**:

```text
Helper functions (test these):
  matchScore          100.0%
  levenshteinDistance 100.0%
  computeDiff         100.0%
  formatFileSize      100.0%

Command runners (skip for unit tests):
  runSearch           0.0%  (integration test)
  runInfo             0.0%  (integration test)

Display functions (low value):
  displayResults      0.0%  (UI only)
```

**Why**: Focusing on business logic gives maximum value. Command runners need real CLI execution (integration tests), and display functions just format output.

**See also**: Testing strategy in [.agent/tasks/0001-prd-cli-testing.md](.agent/tasks/0001-prd-cli-testing.md)

---

### Fuzzy Search Pattern

**When to use**: When implementing search functionality with typo tolerance

**Example**:

```go
func matchScore(query, name, description string) int {
    queryLower := strings.ToLower(query)
    nameLower := strings.ToLower(name)

    // Priority order (highest to lowest score)
    if nameLower == queryLower {
        return 100  // Exact match
    }
    if strings.HasPrefix(nameLower, queryLower) {
        return 80   // Prefix match
    }
    if strings.Contains(nameLower, queryLower) {
        return 60   // Contains match
    }
    if strings.Contains(descLower, queryLower) {
        return 40   // Description match
    }

    // Fuzzy match using Levenshtein distance
    dist := levenshteinDistance(queryLower, nameLower)
    if dist <= 2 && dist < len(nameLower)/2 {
        return 30 - dist*5  // Penalize more edits
    }
    return 0
}
```

**Why**: This scoring approach:

- Exact matches always win
- Prefix matches beat substring matches
- Description matches are lower priority
- Fuzzy matching catches typos but with penalty

**See also**: [search.go](internal/commands/search.go#L169-L205)

---

### Related Components Pattern

**When to use**: When showing relationships between components

**Example**:

```go
// Language → Frameworks mapping
var languageFrameworks = map[string][]string{
    "typescript": {"react", "nextjs", "express"},
    "python":     {"django", "fastapi", "flask"},
    "go":         {"gin", "echo", "fiber"},
}

// Framework → Language mapping (reverse lookup)
func getLanguageForFramework(framework string) []RelatedComponent {
    for lang, frameworks := range languageFrameworks {
        for _, fw := range frameworks {
            if fw == framework {
                return []RelatedComponent{{Name: lang, Type: "language"}}
            }
        }
    }
    return nil
}
```

**Why**: Bidirectional mappings help users discover related components when browsing `info` output.

**See also**: [info.go](internal/commands/info.go#L177-L240)
