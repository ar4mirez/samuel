# CLI Interactive Mode Fix

**Date**: 2026-01-14
**Status**: Decided
**Affects**: packages/cli/internal/cmd/init.go, packages/cli/internal/ui/prompts.go

## Context

Users reported that running `aicof init --template full test-project` would hang indefinitely, only proceeding after pressing Ctrl+C. The issue occurred even when command-line flags explicitly specified all options.

## Problem Analysis

Three interactive prompts were triggering even when CLI flags were provided:

1. **Confirmation prompt** (line 188-195) - Always triggered in non-interactive mode
2. **Language selection prompt** (line 127-150) - Triggered when template wasn't "full"
3. **Framework selection prompt** (line 152-177) - Same issue

Additionally, the `promptui.Prompt` with `IsConfirm: true` had quirky behavior where Enter on empty input returned `ErrAbort` instead of using the default.

## Options Considered

### Option A: Add more `--non-interactive` checks
- **Pros**: Simple, explicit
- **Cons**: Requires users to always add `--non-interactive` flag

### Option B: Track CLI-provided flags and skip prompts
- **Pros**: Better UX - if user provided flags, assume they know what they want
- **Cons**: Slightly more complex logic

## Decision

Option B: Track whether user provided CLI flags (`--template`, `--languages`, `--frameworks`) and skip interactive prompts when any are provided.

## Implementation

### Changes to init.go

```go
// Track if user provided CLI flags (skip confirmation prompt if so)
cliProvided := templateFlag != "" || len(languageFlags) > 0 || len(frameworkFlags) > 0

// Working variable for template name (may be set interactively)
templateName := templateFlag
```

Added `!cliProvided` to all interactive prompt conditions:
- Line 128: Language selection
- Line 153: Framework selection
- Line 195: Confirmation prompt

### Changes to prompts.go

Fixed `Confirm` function to not use `IsConfirm: true`:

```go
prompt := promptui.Prompt{
    Label:   label + suffix,
    Default: defaultStr,  // "y" or "n" instead of ""
}
```

## Consequences

- Users can now use `aicof init --template full <dir>` without any prompts
- Interactive mode still works fully when no flags provided
- Ctrl+C properly cancels operations
- Enter on confirmation now correctly uses default value

## Testing

Verified all three templates work without hanging:
```bash
aicof init --template minimal test-dir  # Works
aicof init --template starter test-dir  # Works
aicof init --template full test-dir     # Works
aicof init --languages go,rust test-dir # Works
```

## References

- Related commits: fix in progress (uncommitted)
- Previous issue: CLI hanging for 10 minutes during init
