# PRD: CLI Command Enhancements

> **ID**: 0004
> **Status**: Implemented
> **Created**: 2026-01-14
> **Author**: Claude Code (AI)

---

## Introduction

This feature adds 4 new CLI commands (`search`, `info`, `config`, `diff`) and improves documentation for all existing commands to enhance discoverability and usability of the AICoF CLI.

**Goal**: Make it easier for users to discover, explore, and manage AICoF components.

---

## Goals

1. Enable component search with fuzzy matching for typos
2. Provide detailed component information before installation
3. Allow configuration management without manual YAML editing
4. Show version differences before updating
5. Document all CLI commands comprehensively

---

## User Stories

**US-001**: As a user, I want to search for components by keyword so that I can find relevant guides without knowing exact names.

**US-002**: As a user, I want to see component details (description, related components, preview) before installing so that I can make informed decisions.

**US-003**: As a user, I want to manage configuration via CLI so that I don't need to manually edit YAML files.

**US-004**: As a user, I want to see what changed between versions before updating so that I can understand the impact of updates.

---

## Functional Requirements

### Search Command (FR-001 to FR-005)

- FR-001: Search across languages, frameworks, and workflows
- FR-002: Support fuzzy matching with Levenshtein distance ≤ 2
- FR-003: Filter by component type with `--type` flag
- FR-004: Show installation status for each result
- FR-005: Sort results by relevance score

### Info Command (FR-006 to FR-010)

- FR-006: Display component description and file path
- FR-007: Show installation status and file metadata
- FR-008: List related components (frameworks for languages, language for frameworks)
- FR-009: Preview first N lines with `--preview` flag
- FR-010: Support type aliases (lang/l, fw/f, wf/w)

### Config Command (FR-011 to FR-015)

- FR-011: `config list` shows all configuration values
- FR-012: `config get <key>` retrieves specific value
- FR-013: `config set <key> <value>` updates configuration
- FR-014: Validate configuration keys before setting
- FR-015: Persist changes to aicof.yaml

### Diff Command (FR-016 to FR-020)

- FR-016: Compare installed files with latest version (default)
- FR-017: Compare two specific versions with arguments
- FR-018: Show added, modified, and removed files
- FR-019: Component-level view with `--components` flag
- FR-020: Download versions to cache if not present

---

## Non-Goals

- Real-time version checking (deferred to future)
- Automatic update suggestions (deferred to future)
- Web-based component browser
- Component dependency management

---

## Technical Considerations

### Implementation

| File | Lines | Purpose |
|------|-------|---------|
| `internal/commands/search.go` | ~270 | Search with fuzzy matching |
| `internal/commands/info.go` | ~275 | Component details and preview |
| `internal/commands/config_cmd.go` | ~200 | Configuration management |
| `internal/commands/diff.go` | ~420 | Version comparison |
| `internal/core/config.go` | +80 | GetValue/SetValue methods |

### Algorithms

- **Fuzzy Matching**: Levenshtein distance algorithm for typo tolerance
- **Score Ranking**: Exact=100, Prefix=80, Contains=60, Description=40, Fuzzy=30
- **File Comparison**: MD5 hashing for content comparison

### Type Aliases

All commands support shorthand aliases:
- `language` → `lang`, `l`
- `framework` → `fw`, `f`
- `workflow` → `wf`, `w`

---

## Implementation Summary

### Phase 1: Config Command ✅

Created `internal/commands/config_cmd.go`:
- Subcommands: list, get, set
- Added GetValue/SetValue/GetAllValues to config.go
- Valid keys validation

### Phase 2: Search Command ✅

Created `internal/commands/search.go`:
- Levenshtein distance algorithm for fuzzy matching
- Multi-score ranking system
- Type filtering with `--type` flag
- Result grouping by component type

### Phase 3: Info Command ✅

Created `internal/commands/info.go`:
- Component detail display
- Related components mapping (language ↔ framework)
- File preview with `--preview` flag
- Installation status with file metadata

### Phase 4: Diff Command ✅

Created `internal/commands/diff.go`:
- MD5-based file hashing for comparison
- Two modes: installed vs latest, version vs version
- File-level and component-level views
- Summary statistics

### Phase 5: Documentation ✅

Updated `README.md`:
- Added comprehensive CLI Commands section
- Core Commands, Component Management, Discovery, Configuration tables
- Command examples with practical workflows
- Updated workflow count to 13
- Expanded workflows documentation table

---

## Files Created/Modified

| File | Action | Purpose |
|------|--------|---------|
| `internal/commands/config_cmd.go` | Created | Config command |
| `internal/commands/search.go` | Created | Search command |
| `internal/commands/info.go` | Created | Info command |
| `internal/commands/diff.go` | Created | Diff command |
| `internal/core/config.go` | Modified | Added value methods |
| `README.md` | Modified | CLI documentation |

---

## Success Metrics

### Technical

- All 4 commands pass linting
- Build succeeds for all platforms
- Help text is clear and complete

### User Experience

- Search finds relevant results with typos
- Info shows useful context before installation
- Config changes persist correctly
- Diff accurately shows version changes

---

## References

- Plan file: `/Users/ar4mirez/.claude/plans/lazy-watching-feigenbaum.md`
- GitHub Release PRD: `.agent/tasks/0003-prd-github-release-integration.md`
