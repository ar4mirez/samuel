# PRD: Go Project Restructure & Framework Self-Initialization

**PRD Number**: 0002
**Created**: 2026-01-14
**Status**: Draft
**Author**: Claude Code (AI)

---

## Introduction

This PRD addresses two related needs:

1. **Go Project Restructure**: The current `packages/cli/` structure doesn't follow Go conventions
2. **Framework Self-Initialization**: The AICoF repository itself should properly use its own framework

**Problem Statement**:
- The Go CLI is nested in `packages/cli/` (Node.js monorepo pattern, not Go)
- The project's `.agent/` directory is incomplete - missing language guides, workflows, etc.
- As the framework that helps others build better, we should "eat our own dog food"

**Goal**: Transform this into a proper Go project that demonstrates best practices AND fully uses the AICoF framework it provides.

---

## Goals

1. Follow the [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
2. Move Go code to root-level directories (`cmd/`, `internal/`)
3. Place `go.mod` at repository root
4. Initialize the project with the AICoF framework (language guides, workflows)
5. Have CLAUDE.md at repository root (not just in template/)
6. Demonstrate proper framework usage as a reference implementation
7. Simplify build and development workflow

---

## User Stories

**US-001**: As a Go developer, I want to clone the repo and see familiar Go project structure (`cmd/`, `internal/`, `go.mod` at root) so I can navigate efficiently.

**US-002**: As someone learning AICoF, I want to see how the framework is used in a real Go project, so I can understand best practices.

**US-003**: As a contributor, I want access to the Go language guide and relevant workflows directly in `.agent/` so AI assistants can help me effectively.

**US-004**: As a user installing via `go install`, I want `go install github.com/ar4mirez/aicof/cmd/aicof@latest` to work correctly.

---

## Current Structure (Problems)

```
aicof/
├── packages/
│   └── cli/                    # ❌ Non-standard Go nesting
│       ├── go.mod              # ❌ Module not at root
│       ├── go.sum
│       ├── Makefile
│       ├── .goreleaser.yaml
│       ├── cmd/aicof/main.go
│       └── internal/
│           ├── cmd/            # ❌ Confusing name (internal/cmd vs cmd/)
│           ├── core/
│           ├── github/
│           └── ui/
├── template/                   # Template files for distribution
│   ├── CLAUDE.md              # Main AI instructions
│   ├── AI_INSTRUCTIONS.md
│   └── .agent/                # Full framework
│       ├── language-guides/   # 21 guides
│       ├── framework-guides/  # 33 guides
│       └── workflows/         # 13 workflows
├── .agent/                    # ❌ Incomplete - missing guides/workflows
│   ├── project.md
│   ├── patterns.md
│   ├── state.md
│   ├── tasks/
│   └── memory/
├── docs/                      # MkDocs documentation
└── (no CLAUDE.md at root)     # ❌ Should have framework at root too
```

**Issues**:
1. `packages/cli/` is Node.js pattern, not idiomatic Go
2. `go.mod` nested means `go install github.com/ar4mirez/aicof@latest` won't work
3. `internal/cmd/` naming conflicts with root `cmd/`
4. No `CLAUDE.md` at root - AI tools won't find it
5. `.agent/` missing language-guides/, workflows/ - incomplete framework
6. Project doesn't demonstrate proper framework usage

---

## Proposed Structure

```
aicof/
├── cmd/
│   └── aicof/
│       └── main.go             # CLI entry point
├── internal/
│   ├── commands/               # Renamed from cmd to avoid confusion
│   │   ├── root.go
│   │   ├── init.go
│   │   ├── update.go
│   │   ├── add.go
│   │   ├── remove.go
│   │   ├── list.go
│   │   ├── doctor.go
│   │   └── version.go
│   ├── core/
│   │   ├── config.go
│   │   ├── downloader.go
│   │   ├── extractor.go
│   │   └── registry.go
│   ├── github/
│   │   └── client.go
│   └── ui/
│       ├── output.go
│       ├── prompts.go
│       └── spinner.go
├── template/                   # Distributable templates (unchanged)
│   ├── CLAUDE.md
│   ├── AI_INSTRUCTIONS.md
│   └── .agent/
│       ├── language-guides/
│       ├── framework-guides/
│       └── workflows/
├── .agent/                     # ✅ COMPLETE framework for this project
│   ├── README.md
│   ├── project.md
│   ├── patterns.md
│   ├── state.md
│   ├── language-guides/        # ✅ Go guide for this project
│   │   └── go.md
│   ├── workflows/              # ✅ All workflows available
│   │   ├── create-prd.md
│   │   ├── generate-tasks.md
│   │   ├── code-review.md
│   │   └── ... (all 13)
│   ├── tasks/
│   └── memory/
├── docs/                       # MkDocs documentation (unchanged)
├── CLAUDE.md                   # ✅ Framework at root for AI tools
├── AI_INSTRUCTIONS.md          # ✅ Alternative format at root
├── go.mod                      # ✅ At repository root
├── go.sum
├── Makefile                    # ✅ At repository root
├── .goreleaser.yaml            # ✅ At repository root
├── README.md
├── CHANGELOG.md
└── LICENSE
```

---

## Functional Requirements

### Part A: Go Project Restructure

#### FR-001: Move Go Module to Root
The `go.mod` file must be at the repository root with module path `github.com/ar4mirez/aicof`.

#### FR-002: Standard cmd/ Directory
The `cmd/aicof/main.go` entry point must be at root level `cmd/aicof/main.go`.

#### FR-003: Rename internal/cmd to internal/commands
Rename `internal/cmd/` to `internal/commands/` to avoid confusion with root `cmd/`.

#### FR-004: Update All Import Paths
All import paths must be updated:
- `github.com/ar4mirez/aicof/internal/cmd` → `github.com/ar4mirez/aicof/internal/commands`

#### FR-005: Move Build Configuration
- `packages/cli/Makefile` → `Makefile` (root)
- `packages/cli/.goreleaser.yaml` → `.goreleaser.yaml` (root)
- Update paths in both files

#### FR-006: Delete packages/ Directory
After migration, remove the empty `packages/` directory.

### Part B: Framework Self-Initialization

#### FR-007: Add CLAUDE.md at Root
Copy `template/CLAUDE.md` to repository root. This is the same file - not a symlink - to ensure AI tools find it.

#### FR-008: Add AI_INSTRUCTIONS.md at Root
Copy `template/AI_INSTRUCTIONS.md` to repository root.

#### FR-009: Add Go Language Guide
Copy `template/.agent/language-guides/go.md` to `.agent/language-guides/go.md`.
This project is Go-only, so we only need the Go guide.

#### FR-010: Add All Workflows
Copy all workflows from `template/.agent/workflows/` to `.agent/workflows/`.
Workflows are language-agnostic and useful for any project.

#### FR-011: Keep Existing .agent/ Content
Preserve existing `.agent/` content:
- `project.md` (with updates for new structure)
- `patterns.md`
- `state.md`
- `tasks/` (PRDs)
- `memory/` (decisions)

### Part C: Documentation Updates

#### FR-012: Update README.md
- Merge `packages/cli/README.md` content into root README
- Update build instructions for new structure
- Document that this project uses its own framework

#### FR-013: Update .agent/project.md
Update architecture section to reflect new structure.

#### FR-014: Update CI/CD
Update GitHub Actions workflows for new paths.

---

## Non-Goals

- Changing the template/ structure - keep as-is
- Changing the docs/ structure - keep as-is
- Adding framework guides (this is CLI-only, no web frameworks)
- Changing the module name - keep `github.com/ar4mirez/aicof`
- Adding new CLI features - separate effort

---

## Technical Considerations

### Import Path Changes

**Before**:
```go
import "github.com/ar4mirez/aicof/internal/cmd"
```

**After**:
```go
import "github.com/ar4mirez/aicof/internal/commands"
```

### File Operations Summary

| Operation | From | To |
|-----------|------|-----|
| Move | `packages/cli/go.mod` | `go.mod` |
| Move | `packages/cli/go.sum` | `go.sum` |
| Move | `packages/cli/Makefile` | `Makefile` |
| Move | `packages/cli/.goreleaser.yaml` | `.goreleaser.yaml` |
| Move | `packages/cli/cmd/aicof/main.go` | `cmd/aicof/main.go` |
| Move+Rename | `packages/cli/internal/cmd/*.go` | `internal/commands/*.go` |
| Move | `packages/cli/internal/core/*.go` | `internal/core/*.go` |
| Move | `packages/cli/internal/github/*.go` | `internal/github/*.go` |
| Move | `packages/cli/internal/ui/*.go` | `internal/ui/*.go` |
| Copy | `template/CLAUDE.md` | `CLAUDE.md` |
| Copy | `template/AI_INSTRUCTIONS.md` | `AI_INSTRUCTIONS.md` |
| Copy | `template/.agent/language-guides/go.md` | `.agent/language-guides/go.md` |
| Copy | `template/.agent/workflows/*.md` | `.agent/workflows/*.md` |
| Delete | `packages/` | (remove directory) |
| Merge+Delete | `packages/cli/README.md` | (merge into root README.md) |

### Makefile Updates

Root-level `Makefile`:
```makefile
.PHONY: build clean test install

BINARY_NAME=aicof
BUILD_DIR=bin

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/aicof

clean:
	rm -rf $(BUILD_DIR)

test:
	go test -v ./...

install:
	go install ./cmd/aicof
```

### GoReleaser Updates

Root-level `.goreleaser.yaml`:
```yaml
builds:
  - main: ./cmd/aicof
    binary: aicof
    # ... rest unchanged
```

---

## Why Self-Initialize?

1. **Dogfooding**: We should use what we build. If AICoF helps others write better code, it should help us too.

2. **Reference Implementation**: Developers learning AICoF can look at this repo to see proper usage.

3. **AI Assistance**: With CLAUDE.md at root and Go guide in `.agent/`, AI tools can assist contributors effectively.

4. **Consistency**: Having workflows available means we can use create-prd, code-review, etc. on this project.

5. **Completeness**: A framework project without using its own framework is incomplete.

---

## Guardrails Affected

### Code Quality
- ✓ All import paths must compile after changes
- ✓ No circular dependencies introduced
- ✓ Build must succeed from repository root

### Testing
- ✓ All existing functionality must work after restructure
- ✓ Manual testing of all 7 commands required
- ✓ `go install` must work

### Git Hygiene
- ✓ Use `git mv` for file moves to preserve history
- ✓ Atomic commits or logical sequence
- ✓ Clear commit messages explaining restructure

---

## Success Metrics

1. `go build ./cmd/aicof` works from repository root
2. `go install github.com/ar4mirez/aicof/cmd/aicof@latest` works
3. All 7 CLI commands function correctly
4. GoReleaser builds succeed
5. GitHub Actions CI passes
6. `CLAUDE.md` exists at root and is found by AI tools
7. `.agent/language-guides/go.md` exists for Go development
8. `.agent/workflows/` contains all 13 workflows
9. No `packages/` directory remains
10. Project README documents framework usage

---

## Implementation Phases

### Phase 1: Go Restructure
1. Move `go.mod`, `go.sum` to root
2. Move `cmd/` to root
3. Move and rename `internal/cmd/` → `internal/commands/`
4. Move remaining `internal/` packages
5. Update all import paths
6. Move Makefile and .goreleaser.yaml
7. Delete packages/ directory
8. Verify build and all commands work

### Phase 2: Framework Self-Initialization
1. Copy CLAUDE.md to root
2. Copy AI_INSTRUCTIONS.md to root
3. Create `.agent/language-guides/` and add go.md
4. Create `.agent/workflows/` and copy all workflows
5. Update .agent/project.md for new structure

### Phase 3: Documentation & Cleanup
1. Merge packages/cli/README.md into root README
2. Update any docs/ references
3. Update GitHub Actions
4. Final testing
5. Commit and push

---

## Open Questions

1. **Q**: Should CLAUDE.md at root be a copy or symlink to template/CLAUDE.md?
   **A**: Copy - symlinks can cause issues with some tools and git hosting

2. **Q**: Should we add all 21 language guides or just Go?
   **A**: Just Go - this is a Go-only project. Other guides would be noise.

3. **Q**: Should we add any framework guides?
   **A**: No - this is a CLI tool, no web frameworks involved

---

## References

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Go Modules Reference](https://go.dev/ref/mod)
- [Effective Go](https://go.dev/doc/effective_go)
- AICoF Framework: `template/CLAUDE.md`
