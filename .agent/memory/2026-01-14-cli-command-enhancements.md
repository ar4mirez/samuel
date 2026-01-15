# CLI Command Enhancements Decision

**Date**: 2026-01-14
**Status**: Decided
**Affects**: CLI tool, user experience, discoverability

## Context

The AICoF CLI had 7 commands (`init`, `update`, `add`, `remove`, `list`, `doctor`, `version`). User feedback and analysis identified gaps:

1. **Discoverability** - No way to search for components by keyword
2. **Preview** - Can't see component details before installing
3. **Configuration** - Must manually edit `aicof.yaml`
4. **Version comparison** - Can't see what changed before updating

## Options Considered

### Option A: Add 4 New Commands (Chosen)

- `search` - Fuzzy search across all components
- `info` - Show component details with preview
- `config` - Manage configuration via CLI
- `diff` - Compare versions before update

**Pros:**

- Each command has single responsibility
- Follows Unix philosophy (do one thing well)
- Easy to discover via `--help`
- Consistent with existing command patterns

**Cons:**

- More commands to learn
- More code to maintain

### Option B: Extend Existing Commands

Extend `list` with search, extend `add` with preview, etc.

**Pros:**

- Fewer top-level commands

**Cons:**

- Overloaded commands with too many flags
- Harder to discover features
- Complex `--help` output

### Option C: Interactive TUI Mode

Single interactive command with menus.

**Pros:**

- Discoverable via browsing

**Cons:**

- Doesn't work in scripts/CI
- Requires additional TUI library
- Diverges from existing CLI patterns

## Decision

**Option A: Add 4 New Commands**

Rationale:

1. Maintains consistency with existing commands
2. Each command is scriptable and composable
3. Follows established CLI patterns (git, npm, docker)
4. Easier to test individual commands
5. Clear separation of concerns

## Implementation Details

### Command Specifications

| Command | Usage | Description |
|---------|-------|-------------|
| `search` | `aicof search <query> [--type]` | Fuzzy search with Levenshtein distance ≤2 |
| `info` | `aicof info <type> <name> [--preview]` | Show details, file size, related components |
| `config` | `aicof config list\|get\|set` | View/modify aicof.yaml values |
| `diff` | `aicof diff [--installed\|v1 v2]` | Compare versions, show added/modified/removed |

### Testing Strategy

- Unit tests for helper functions (80%+ coverage)
- Integration tests for command runners (future)
- 7 test files created with 77+ test functions

### Files Added

```text
internal/commands/
├── search.go, search_test.go
├── info.go, info_test.go
├── config_cmd.go, config_cmd_test.go
├── diff.go, diff_test.go, diff_display.go, diff_display_test.go

internal/core/
├── config_test.go (enhanced config.go)
├── registry_test.go
```

## Consequences

### Positive

- Users can discover components without reading docs
- Safe previews before installing/updating
- Configuration changes don't require manual YAML editing
- Version comparison reduces update anxiety

### Negative

- 4 more commands to maintain
- Increased test surface area
- Documentation needs updating

### Technical Debt

- Integration tests not yet implemented for command runners
- README documentation needs updating

## References

- PRD: `.agent/tasks/0004-prd-cli-command-enhancements.md`
- Patterns: `.agent/patterns.md` (Fuzzy Search Pattern, Testing Patterns)
