# Task List: Go Project Restructure & Framework Self-Initialization

**PRD**: 0002-prd-go-project-restructure.md
**Created**: 2026-01-14
**Status**: Not Started

---

## Phase 1: Go Project Restructure

### Task 1.1: Create Root Directory Structure
**Complexity**: Low
**Dependencies**: None
**Files**: cmd/, internal/

**Steps**:
1. Create `cmd/aicof/` directory at repository root
2. Create `internal/commands/` directory at repository root
3. Create `internal/core/` directory at repository root
4. Create `internal/github/` directory at repository root
5. Create `internal/ui/` directory at repository root

**Verification**:
- [ ] Directory structure exists at root level
- [ ] No errors creating directories

---

### Task 1.2: Move Go Module Files
**Complexity**: Low
**Dependencies**: None
**Files**: go.mod, go.sum

**Steps**:
1. `git mv packages/cli/go.mod go.mod`
2. `git mv packages/cli/go.sum go.sum`
3. Verify module path is `github.com/ar4mirez/aicof` (no change needed)

**Verification**:
- [ ] `go.mod` exists at repository root
- [ ] `go.sum` exists at repository root
- [ ] Module path unchanged

---

### Task 1.3: Move Entry Point
**Complexity**: Low
**Dependencies**: 1.1
**Files**: cmd/aicof/main.go

**Steps**:
1. `git mv packages/cli/cmd/aicof/main.go cmd/aicof/main.go`
2. Update import path: `github.com/ar4mirez/aicof/internal/cmd` → `github.com/ar4mirez/aicof/internal/commands`

**Verification**:
- [ ] `cmd/aicof/main.go` exists at root
- [ ] Import path updated to `internal/commands`

---

### Task 1.4: Move and Rename Command Package
**Complexity**: Medium
**Dependencies**: 1.1
**Files**: internal/commands/*.go (8 files)

**Steps**:
1. Move all files from `packages/cli/internal/cmd/` to `internal/commands/`:
   - `git mv packages/cli/internal/cmd/root.go internal/commands/root.go`
   - `git mv packages/cli/internal/cmd/init.go internal/commands/init.go`
   - `git mv packages/cli/internal/cmd/update.go internal/commands/update.go`
   - `git mv packages/cli/internal/cmd/add.go internal/commands/add.go`
   - `git mv packages/cli/internal/cmd/remove.go internal/commands/remove.go`
   - `git mv packages/cli/internal/cmd/list.go internal/commands/list.go`
   - `git mv packages/cli/internal/cmd/doctor.go internal/commands/doctor.go`
   - `git mv packages/cli/internal/cmd/version.go internal/commands/version.go`
2. Update package declaration in each file: `package cmd` → `package commands`

**Verification**:
- [ ] All 8 files moved to `internal/commands/`
- [ ] Package declaration is `package commands` in all files
- [ ] No files remain in `packages/cli/internal/cmd/`

---

### Task 1.5: Move Core Package
**Complexity**: Low
**Dependencies**: 1.1
**Files**: internal/core/*.go (4 files)

**Steps**:
1. Move all files from `packages/cli/internal/core/` to `internal/core/`:
   - `git mv packages/cli/internal/core/config.go internal/core/config.go`
   - `git mv packages/cli/internal/core/downloader.go internal/core/downloader.go`
   - `git mv packages/cli/internal/core/extractor.go internal/core/extractor.go`
   - `git mv packages/cli/internal/core/registry.go internal/core/registry.go`

**Verification**:
- [ ] All 4 files moved to `internal/core/`
- [ ] No files remain in `packages/cli/internal/core/`

---

### Task 1.6: Move GitHub Package
**Complexity**: Low
**Dependencies**: 1.1
**Files**: internal/github/client.go

**Steps**:
1. `git mv packages/cli/internal/github/client.go internal/github/client.go`

**Verification**:
- [ ] File moved to `internal/github/`
- [ ] No files remain in `packages/cli/internal/github/`

---

### Task 1.7: Move UI Package
**Complexity**: Low
**Dependencies**: 1.1
**Files**: internal/ui/*.go (3 files)

**Steps**:
1. Move all files from `packages/cli/internal/ui/` to `internal/ui/`:
   - `git mv packages/cli/internal/ui/output.go internal/ui/output.go`
   - `git mv packages/cli/internal/ui/prompts.go internal/ui/prompts.go`
   - `git mv packages/cli/internal/ui/spinner.go internal/ui/spinner.go`

**Verification**:
- [ ] All 3 files moved to `internal/ui/`
- [ ] No files remain in `packages/cli/internal/ui/`

---

### Task 1.8: Update All Import Paths
**Complexity**: Medium
**Dependencies**: 1.3, 1.4, 1.5, 1.6, 1.7
**Files**: Multiple Go files

**Steps**:
1. In `cmd/aicof/main.go`:
   - Change `github.com/ar4mirez/aicof/internal/cmd` → `github.com/ar4mirez/aicof/internal/commands`
2. In all `internal/commands/*.go` files:
   - Update any imports from `internal/cmd` → `internal/commands`
3. Run `go mod tidy` to clean up

**Verification**:
- [ ] `go build ./cmd/aicof` succeeds
- [ ] No import errors
- [ ] `go mod tidy` completes without errors

---

### Task 1.9: Move Build Configuration
**Complexity**: Medium
**Dependencies**: 1.2
**Files**: Makefile, .goreleaser.yaml

**Steps**:
1. `git mv packages/cli/Makefile Makefile`
2. `git mv packages/cli/.goreleaser.yaml .goreleaser.yaml`
3. Update Makefile paths (should work as-is since paths are relative)
4. Update `.goreleaser.yaml`:
   - Ensure `main: ./cmd/aicof` is correct

**Verification**:
- [ ] `make build` works from repository root
- [ ] Binary created at `bin/aicof`
- [ ] `.goreleaser.yaml` has correct paths

---

### Task 1.10: Update .gitignore
**Complexity**: Low
**Dependencies**: 1.9
**Files**: .gitignore

**Steps**:
1. Update `.gitignore` to remove `packages/cli/bin/` entry
2. Add `bin/` at root level if not already present

**Verification**:
- [ ] `bin/` is ignored at root level
- [ ] No reference to `packages/cli/bin/`

---

### Task 1.11: Delete packages/ Directory
**Complexity**: Low
**Dependencies**: 1.3, 1.4, 1.5, 1.6, 1.7, 1.9
**Files**: packages/ (directory)

**Steps**:
1. Verify all files have been moved out of `packages/cli/`
2. Delete empty directories: `rm -rf packages/`
3. Stage deletion: `git add -A`

**Verification**:
- [ ] `packages/` directory no longer exists
- [ ] All Go code is at root level
- [ ] No orphaned files

---

### Task 1.12: Verify Build and Commands
**Complexity**: Medium
**Dependencies**: 1.8, 1.9, 1.11
**Files**: None (verification only)

**Steps**:
1. Run `go build ./cmd/aicof`
2. Run `./bin/aicof version`
3. Test all 7 commands:
   - `./bin/aicof init --help`
   - `./bin/aicof update --help`
   - `./bin/aicof add --help`
   - `./bin/aicof remove --help`
   - `./bin/aicof list --help`
   - `./bin/aicof doctor --help`
   - `./bin/aicof version`
4. Test actual init: `./bin/aicof init --template minimal /tmp/test-restructure`

**Verification**:
- [ ] Build succeeds from root
- [ ] All 7 commands respond correctly
- [ ] Init command works end-to-end
- [ ] No runtime errors

---

## Phase 2: Framework Self-Initialization

### Task 2.1: Copy CLAUDE.md to Root
**Complexity**: Low
**Dependencies**: None
**Files**: CLAUDE.md

**Steps**:
1. `cp template/CLAUDE.md CLAUDE.md`
2. Verify content is identical

**Verification**:
- [ ] `CLAUDE.md` exists at repository root
- [ ] Content matches `template/CLAUDE.md`

---

### Task 2.2: Copy AI_INSTRUCTIONS.md to Root
**Complexity**: Low
**Dependencies**: None
**Files**: AI_INSTRUCTIONS.md

**Steps**:
1. `cp template/AI_INSTRUCTIONS.md AI_INSTRUCTIONS.md`
2. Verify content is identical

**Verification**:
- [ ] `AI_INSTRUCTIONS.md` exists at repository root
- [ ] Content matches `template/AI_INSTRUCTIONS.md`

---

### Task 2.3: Add Go Language Guide
**Complexity**: Low
**Dependencies**: None
**Files**: .agent/language-guides/go.md

**Steps**:
1. Create directory: `mkdir -p .agent/language-guides`
2. Copy guide: `cp template/.agent/language-guides/go.md .agent/language-guides/go.md`

**Verification**:
- [ ] `.agent/language-guides/go.md` exists
- [ ] Content matches template version

---

### Task 2.4: Add All Workflows
**Complexity**: Low
**Dependencies**: None
**Files**: .agent/workflows/*.md (13 files)

**Steps**:
1. Create directory: `mkdir -p .agent/workflows`
2. Copy all workflows: `cp template/.agent/workflows/*.md .agent/workflows/`
3. Verify all 13 workflows copied:
   - initialize-project.md
   - create-prd.md
   - generate-tasks.md
   - code-review.md
   - security-audit.md
   - testing-strategy.md
   - cleanup-project.md
   - refactoring.md
   - dependency-update.md
   - update-framework.md
   - troubleshooting.md
   - generate-agents-md.md
   - document-work.md

**Verification**:
- [ ] `.agent/workflows/` directory exists
- [ ] All 13 workflow files present
- [ ] Content matches template versions

---

### Task 2.5: Update .agent/README.md
**Complexity**: Low
**Dependencies**: 2.3, 2.4
**Files**: .agent/README.md

**Steps**:
1. Update `.agent/README.md` to reflect that this project now has:
   - Go language guide
   - All workflows
   - Project-specific files (project.md, patterns.md, state.md)

**Verification**:
- [ ] README accurately describes available guides and workflows
- [ ] Structure section is up to date

---

## Phase 3: Documentation & Cleanup

### Task 3.1: Update .agent/project.md
**Complexity**: Medium
**Dependencies**: 1.11
**Files**: .agent/project.md

**Steps**:
1. Update "Repository Structure" section to show new layout
2. Update "Development Setup" section with new paths
3. Update "CLI Architecture" section if needed
4. Remove references to `packages/cli/`

**Verification**:
- [ ] Structure diagram matches actual directory layout
- [ ] Build commands are correct for new structure
- [ ] No references to old `packages/cli/` path

---

### Task 3.2: Merge CLI README into Root README
**Complexity**: Medium
**Dependencies**: 1.11
**Files**: README.md

**Steps**:
1. Read `packages/cli/README.md` content (before deletion)
2. Integrate relevant sections into root `README.md`:
   - CLI installation instructions
   - CLI usage examples
   - Build instructions
3. Update build paths to reflect new structure
4. Add section noting that this project uses AICoF framework

**Verification**:
- [ ] Root README has CLI documentation
- [ ] Build instructions work with new structure
- [ ] Installation instructions are accurate
- [ ] No references to old structure

---

### Task 3.3: Update GitHub Actions
**Complexity**: Medium
**Dependencies**: 1.11
**Files**: .github/workflows/*.yml

**Steps**:
1. Review `.github/workflows/` for any Go-related workflows
2. Update paths from `packages/cli/` to root
3. Update working directory references
4. Update any `go build` or `go test` commands

**Verification**:
- [ ] CI workflows reference correct paths
- [ ] `go build ./cmd/aicof` used instead of nested paths
- [ ] Workflows can find go.mod at root

---

### Task 3.4: Update .agent/state.md
**Complexity**: Low
**Dependencies**: All previous tasks
**Files**: .agent/state.md

**Steps**:
1. Update state.md to reflect completed restructure
2. Mark PRD 0002 tasks as complete
3. Note any follow-up items

**Verification**:
- [ ] State reflects current project status
- [ ] Completed tasks marked done

---

### Task 3.5: Final Verification
**Complexity**: Medium
**Dependencies**: All previous tasks
**Files**: None (verification only)

**Steps**:
1. Run full build: `go build ./cmd/aicof`
2. Run all tests: `go test ./...`
3. Test CLI end-to-end:
   ```bash
   ./bin/aicof version
   ./bin/aicof init --template minimal /tmp/final-test
   ./bin/aicof doctor
   ```
4. Verify `go install github.com/ar4mirez/aicof/cmd/aicof@latest` would work (check module path)
5. Verify CLAUDE.md is found by AI tools (open in new session)
6. Verify `.agent/language-guides/go.md` exists and is readable
7. Verify all 13 workflows in `.agent/workflows/`

**Verification**:
- [ ] Build succeeds
- [ ] Tests pass (if any exist)
- [ ] All CLI commands work
- [ ] CLAUDE.md at root
- [ ] Go guide available
- [ ] All workflows available
- [ ] No `packages/` directory

---

### Task 3.6: Commit and Push
**Complexity**: Low
**Dependencies**: 3.5
**Files**: All changed files

**Steps**:
1. Stage all changes: `git add -A`
2. Review staged changes: `git status`
3. Commit with descriptive message:
   ```
   refactor: restructure to standard Go project layout

   - Move go.mod and Go code from packages/cli/ to repository root
   - Rename internal/cmd to internal/commands
   - Add CLAUDE.md and AI_INSTRUCTIONS.md at root
   - Add Go language guide to .agent/language-guides/
   - Add all 13 workflows to .agent/workflows/
   - Update documentation for new structure
   - Delete packages/ directory

   This restructure:
   - Follows standard Go project layout
   - Enables `go install github.com/ar4mirez/aicof/cmd/aicof@latest`
   - Self-initializes with AICoF framework (dogfooding)

   BREAKING: Contributors must update their local setup
   ```
4. Push to remote: `git push`

**Verification**:
- [ ] Commit created successfully
- [ ] Push succeeds
- [ ] Remote has new structure

---

## Summary

| Phase | Tasks | Estimated Complexity |
|-------|-------|---------------------|
| Phase 1: Go Restructure | 1.1 - 1.12 | Medium (file moves, import updates) |
| Phase 2: Framework Init | 2.1 - 2.5 | Low (file copies) |
| Phase 3: Documentation | 3.1 - 3.6 | Medium (content updates) |

**Total Tasks**: 23
**Critical Path**: 1.1 → 1.2-1.7 (parallel) → 1.8 → 1.9 → 1.11 → 1.12 → 3.5 → 3.6

**Risk Areas**:
- Task 1.8 (import path updates) - must be thorough
- Task 1.12 (verification) - catch issues before commit
- Task 3.3 (CI updates) - ensure builds don't break
