# Go Project Restructure: packages/cli/ to Root

**Date**: 2026-01-14
**Status**: Decided
**Affects**: All Go code, build configuration, import paths, CI/CD

## Context

The AICoF CLI was originally placed in `packages/cli/` following a Node.js monorepo pattern. This was not idiomatic Go and caused confusion:

- Go developers expect `go.mod` at repository root
- Standard Go layout uses `cmd/` and `internal/` at root
- The `packages/` pattern is from JavaScript/Node.js ecosystem
- Import paths were unnecessarily long: `github.com/ar4mirez/aicof/packages/cli/internal/cmd`

## Options Considered

### Option A: Keep packages/cli/ Structure

- **Pros**: No migration effort, existing code works
- **Cons**: Confusing for Go developers, non-idiomatic, longer import paths, harder to maintain

### Option B: Move Go Code to Repository Root

- **Pros**: Idiomatic Go layout, cleaner imports, standard project structure
- **Cons**: Requires file moves, import path updates, potential git history impact

## Decision

Move Go code to repository root with standard Go project layout:

```text
aicof/
├── cmd/aicof/main.go         # Entry point
├── internal/
│   ├── commands/             # Renamed from internal/cmd
│   ├── core/
│   ├── github/
│   └── ui/
├── go.mod                    # At root
├── Makefile
└── .goreleaser.yaml
```

Key implementation details:

1. **Package rename**: `internal/cmd` → `internal/commands` to avoid confusion with root `cmd/` directory
2. **Import paths**: Updated from `github.com/ar4mirez/aicof/packages/cli/internal/cmd` to `github.com/ar4mirez/aicof/internal/commands`
3. **Build config**: Makefile LDFLAGS updated to new paths
4. **History preservation**: Used `git mv` to preserve file history

## Consequences

### Positive

- Standard Go project structure that Go developers recognize
- Cleaner import paths
- Build tools work without specifying subdirectory
- Easier contribution from Go developers

### Negative

- Old `packages/cli/` paths in documentation needed updating
- Any external references to old paths are broken
- Patterns.md file paths needed updating

### Technical Debt

- Some documentation may still reference old paths (patterns.md was updated)
- CI/CD workflows may need path updates if added later

### Future Work Enabled

- Easier to add more Go packages at root level
- Standard tooling (go build ./..., go test ./...) works from root
- GoReleaser configuration is simplified

## Implementation Notes

The restructure was done in a single commit to keep the change atomic:

```bash
# Key commands used
git mv packages/cli/go.mod go.mod
git mv packages/cli/cmd cmd
git mv packages/cli/internal internal
# Then updated all package declarations and import paths
```

The package was renamed from `cmd` to `commands` because Go convention uses `cmd/` for entry points, and having `internal/cmd` would be confusing alongside the root `cmd/` directory.

## References

- PRD: `.agent/tasks/0002-prd-go-project-restructure.md`
- Task list: `.agent/tasks/tasks-0002-prd-go-project-restructure.md`
- Commit: `d7bd278 refactor: restructure to idiomatic Go project layout`
