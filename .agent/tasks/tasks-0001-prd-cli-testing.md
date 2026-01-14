# Task List: AICoF CLI Testing Suite

> **Source PRD**: `.agent/tasks/0001-prd-cli-testing.md`
> **Generated**: 2026-01-14
> **Status**: Not Started

---

## Current State Assessment

**Repository Structure Issue:**
The current `.agent/` directory serves dual purposes:
1. Template files distributed to users
2. This project's own AI context (tasks, memory)

**Resolution Required:**
Move distributable files to `template/` directory. This affects:
- CLI's download/extraction logic (must extract from `template/`)
- Documentation paths
- User-facing instructions

**Existing Infrastructure:**
- Go CLI in `packages/cli/` (17 source files)
- Cobra CLI framework with 7 commands
- GitHub API client with branch fallback
- No existing tests

**Patterns to Follow:**
- Go standard `testing` package
- Table-driven tests
- `testify/assert` for assertions
- `httptest` for mocked HTTP

---

## Relevant Files

### Phase 0: Repository Restructure (Create)

- `template/CLAUDE.md` - Move from root
- `template/AI_INSTRUCTIONS.md` - Move from root
- `template/.agent/` - Move entire directory structure
- `template/.agent/language-guides/` - All 21 language guides
- `template/.agent/framework-guides/` - All 33 framework guides
- `template/.agent/workflows/` - All 13 workflows

### Phase 0: Repository Restructure (Modify)

- `packages/cli/internal/core/extractor.go` - Update paths to use `template/` prefix
- `packages/cli/internal/core/registry.go` - Update all component paths
- `README.md` - Update installation instructions
- `mkdocs.yml` - Update documentation paths if needed

### Phase 1: Unit Tests (Create)

- `packages/cli/internal/core/config_test.go` - Config parsing tests
- `packages/cli/internal/core/registry_test.go` - Registry lookup tests
- `packages/cli/internal/core/extractor_test.go` - File extraction tests
- `packages/cli/internal/github/client_test.go` - GitHub API tests (mocked)
- `packages/cli/testdata/` - Test fixtures directory

### Phase 2: E2E Tests (Create)

- `packages/cli/e2e/e2e_test.go` - E2E test suite setup
- `packages/cli/e2e/init_test.go` - Init command tests
- `packages/cli/e2e/update_test.go` - Update command tests
- `packages/cli/e2e/list_test.go` - List command tests
- `packages/cli/e2e/add_remove_test.go` - Add/Remove tests
- `packages/cli/e2e/doctor_test.go` - Doctor command tests
- `packages/cli/e2e/version_test.go` - Version command tests

### Reference (Read Only)

- `packages/cli/internal/cmd/*.go` - Command implementations
- `packages/cli/internal/core/*.go` - Core packages
- `CLAUDE.md` - Guardrails for testing requirements

---

## Implementation Notes

### Testing Requirements

- Unit tests: >80% coverage for `internal/core/`
- E2E tests: All 7 commands covered
- Tests must be independent (run in any order)
- Unit tests use mocked HTTP (fast, offline)
- E2E tests use real GitHub API (realistic)

### Test Execution

```bash
# Unit tests only (fast)
make test

# E2E tests only (requires network)
go test -v ./e2e/... -tags=e2e

# All tests with coverage
make test-coverage
```

### Common Pitfalls

- ❌ Don't hardcode paths - use `filepath.Join()`
- ❌ Don't leave temp directories - use `t.Cleanup()`
- ❌ Don't share state between tests
- ❌ Don't skip error assertions
- ✅ Use table-driven tests for multiple cases
- ✅ Mock external dependencies in unit tests

---

## Tasks

### Phase 0: Repository Restructure (PREREQUISITE)

- [ ] 0.0 Repository Restructure
  - [ ] 0.1 Create `template/` directory structure [~1,000 tokens - Simple]
  - [ ] 0.2 Move CLAUDE.md and AI_INSTRUCTIONS.md to `template/` [~500 tokens - Simple]
  - [ ] 0.3 Move `.agent/` contents to `template/.agent/` (exclude tasks/, memory/) [~1,000 tokens - Simple]
  - [ ] 0.4 Create project-specific `.agent/` with project.md, tasks/, memory/ [~500 tokens - Simple]
  - [ ] 0.5 Update CLI registry.go paths to use `template/` prefix [~2,000 tokens - Medium]
  - [ ] 0.6 Update CLI extractor.go to handle template/ structure [~1,500 tokens - Medium]
  - [ ] 0.7 Update README.md with new structure explanation [~1,000 tokens - Simple]
  - [ ] 0.8 Test CLI still works after restructure [~500 tokens - Simple]
    <!-- Guardrails:
      ✓ No breaking changes to user-facing behavior
      ✓ CLI downloads and extracts correctly
      ✓ All existing functionality preserved
    -->

### Phase 1: Unit Tests (Mocked)

- [ ] 1.0 Test Infrastructure Setup
  - [ ] 1.1 Add testify dependency to go.mod [~500 tokens - Simple]
  - [ ] 1.2 Create testdata/ directory with fixtures [~1,000 tokens - Simple]
  - [ ] 1.3 Create test helper functions (setupTestDir, etc.) [~1,500 tokens - Medium]
    <!-- Guardrails:
      ✓ Helper functions are reusable
      ✓ Temp directories cleaned up
    -->

- [ ] 2.0 Config Package Tests
  - [ ] 2.1 Test LoadConfig with valid YAML [~1,500 tokens - Medium]
  - [ ] 2.2 Test LoadConfig with missing file [~1,000 tokens - Simple]
  - [ ] 2.3 Test LoadConfig with malformed YAML [~1,000 tokens - Simple]
  - [ ] 2.4 Test Config.Save preserves format [~1,000 tokens - Simple]
  - [ ] 2.5 Test Add/Remove component methods [~1,500 tokens - Medium]
  - [ ] 2.6 Test Has* check methods [~1,000 tokens - Simple]
    <!-- Guardrails:
      ✓ Edge cases tested (null, empty)
      ✓ All public functions tested
    -->

- [ ] 3.0 Registry Package Tests
  - [ ] 3.1 Test FindLanguage returns correct component [~1,000 tokens - Simple]
  - [ ] 3.2 Test FindFramework returns correct component [~1,000 tokens - Simple]
  - [ ] 3.3 Test FindWorkflow returns correct component [~1,000 tokens - Simple]
  - [ ] 3.4 Test FindTemplate returns correct template [~1,000 tokens - Simple]
  - [ ] 3.5 Test GetComponentPaths returns all paths [~1,500 tokens - Medium]
  - [ ] 3.6 Test unknown component returns nil [~500 tokens - Simple]
    <!-- Guardrails:
      ✓ All 21 languages in registry
      ✓ All 33 frameworks in registry
      ✓ All 13 workflows in registry
    -->

- [ ] 4.0 Extractor Package Tests
  - [ ] 4.1 Test Extract single file [~2,000 tokens - Medium]
  - [ ] 4.2 Test Extract directory recursively [~2,000 tokens - Medium]
  - [ ] 4.3 Test Extract skips existing (force=false) [~1,500 tokens - Medium]
  - [ ] 4.4 Test Extract overwrites (force=true) [~1,500 tokens - Medium]
  - [ ] 4.5 Test Extract handles missing source [~1,000 tokens - Simple]
  - [ ] 4.6 Test BackupFile and RestoreBackup [~2,000 tokens - Medium]
    <!-- Guardrails:
      ✓ File paths validated (no traversal)
      ✓ Permissions preserved
    -->

- [ ] 5.0 GitHub Client Tests (Mocked)
  - [ ] 5.1 Create mock HTTP server helper [~2,000 tokens - Medium]
  - [ ] 5.2 Test GetLatestRelease parses JSON [~1,500 tokens - Medium]
  - [ ] 5.3 Test GetLatestRelease handles 404 [~1,000 tokens - Simple]
  - [ ] 5.4 Test GetLatestVersionOrBranch fallback [~1,500 tokens - Medium]
  - [ ] 5.5 Test URL construction methods [~1,000 tokens - Simple]
  - [ ] 5.6 Test DownloadArchive with mock [~2,000 tokens - Medium]
    <!-- Guardrails:
      ✓ HTTP client timeout tested
      ✓ Error responses handled
    -->

### Phase 2: E2E Tests (Real Network)

- [ ] 6.0 E2E Test Infrastructure
  - [ ] 6.1 Create e2e/ directory and test suite setup [~2,000 tokens - Medium]
  - [ ] 6.2 Create CLI runner helper (exec command, capture output) [~2,000 tokens - Medium]
  - [ ] 6.3 Create temp project setup/teardown helpers [~1,500 tokens - Medium]
    <!-- Guardrails:
      ✓ Tests isolated (each gets fresh directory)
      ✓ Cleanup on failure
    -->

- [ ] 7.0 Version Command E2E
  - [ ] 7.1 Test `aicof version` shows CLI info [~1,000 tokens - Simple]
  - [ ] 7.2 Test `aicof version` without installation [~1,000 tokens - Simple]
  - [ ] 7.3 Test `aicof version --check` [~1,500 tokens - Medium]

- [ ] 8.0 Init Command E2E
  - [ ] 8.1 Test `aicof init .` creates files [~2,500 tokens - Medium]
  - [ ] 8.2 Test `aicof init my-project` creates directory [~2,000 tokens - Medium]
  - [ ] 8.3 Test `aicof init --template minimal` [~2,000 tokens - Medium]
  - [ ] 8.4 Test `aicof init --template full` [~2,000 tokens - Medium]
  - [ ] 8.5 Test `aicof init --non-interactive` [~1,500 tokens - Medium]
  - [ ] 8.6 Test `aicof init --force` overwrites [~2,000 tokens - Medium]
  - [ ] 8.7 Test init fails when already initialized [~1,500 tokens - Medium]
    <!-- Guardrails:
      ✓ Creates aicof.yaml
      ✓ Creates CLAUDE.md
      ✓ Creates .agent/ structure
    -->

- [ ] 9.0 List Command E2E
  - [ ] 9.1 Test `aicof list` shows installed [~1,500 tokens - Medium]
  - [ ] 9.2 Test `aicof list --available` shows all [~1,500 tokens - Medium]
  - [ ] 9.3 Test `aicof list --type languages` filters [~1,500 tokens - Medium]
  - [ ] 9.4 Test `aicof list` without installation [~1,000 tokens - Simple]

- [ ] 10.0 Add/Remove Commands E2E
  - [ ] 10.1 Test `aicof add language go` [~2,000 tokens - Medium]
  - [ ] 10.2 Test `aicof add framework react` [~2,000 tokens - Medium]
  - [ ] 10.3 Test `aicof remove language go` [~2,000 tokens - Medium]
  - [ ] 10.4 Test add unknown component fails [~1,000 tokens - Simple]
  - [ ] 10.5 Test add when not initialized fails [~1,000 tokens - Simple]
  - [ ] 10.6 Test add already installed warns [~1,000 tokens - Simple]

- [ ] 11.0 Update Command E2E
  - [ ] 11.1 Test `aicof update --check` [~2,000 tokens - Medium]
  - [ ] 11.2 Test `aicof update --diff` [~2,000 tokens - Medium]
  - [ ] 11.3 Test `aicof update` applies changes [~3,000 tokens - Complex]
  - [ ] 11.4 Test update preserves modified files [~2,500 tokens - Medium]
  - [ ] 11.5 Test `aicof update --force` overwrites [~2,000 tokens - Medium]

- [ ] 12.0 Doctor Command E2E
  - [ ] 12.1 Test `aicof doctor` on healthy install [~1,500 tokens - Medium]
  - [ ] 12.2 Test doctor detects missing config [~1,500 tokens - Medium]
  - [ ] 12.3 Test doctor detects missing files [~1,500 tokens - Medium]
  - [ ] 12.4 Test `aicof doctor --fix` repairs [~2,500 tokens - Medium]

### Phase 3: Error Handling Tests

- [ ] 13.0 Error Scenarios
  - [ ] 13.1 Test network timeout handling [~2,000 tokens - Medium]
  - [ ] 13.2 Test invalid GitHub response [~1,500 tokens - Medium]
  - [ ] 13.3 Test invalid command arguments [~1,500 tokens - Medium]
  - [ ] 13.4 Test permission denied errors [~1,500 tokens - Medium]

### Phase 4: CI Integration

- [ ] 14.0 CI/CD Setup
  - [ ] 14.1 Update Makefile test targets [~1,000 tokens - Simple]
  - [ ] 14.2 Add test coverage reporting [~1,500 tokens - Medium]
  - [ ] 14.3 Create GitHub Actions workflow for tests [~2,000 tokens - Medium]
  - [ ] 14.4 Add coverage badge to README [~500 tokens - Simple]

---

## Progress Tracking

**Total Tasks**: 14 parent, 65 sub-tasks
**Completed**: 0/65 (0%)
**In Progress**: None
**Blocked**: None

**Last Updated**: 2026-01-14

---

## Success Criteria

Before marking this task list complete, verify:

- [ ] Repository restructured (template/ contains distributable files)
- [ ] CLI works correctly with new structure
- [ ] Unit test coverage >80% for internal/core/
- [ ] All 7 CLI commands have E2E tests
- [ ] Error scenarios tested
- [ ] Tests run in CI pipeline
- [ ] Documentation updated

---

## High-Level Summary

**Phase 0** (Prerequisite): Restructure repo - move templates to `template/`
**Phase 1**: Unit tests with mocked HTTP (~12,000 tokens)
**Phase 2**: E2E tests with real network (~35,000 tokens)
**Phase 3**: Error handling tests (~6,500 tokens)
**Phase 4**: CI integration (~5,000 tokens)

**Total Estimate**: ~58,500 tokens of test code

---

## Next Steps

1. **Review this task list** - confirm structure is correct
2. **Start with Phase 0** - repo restructure is prerequisite
3. Then proceed with Phase 1 unit tests (fast feedback)
4. Then Phase 2 E2E tests (comprehensive validation)

Ready to start with **Task 0.1** (create template/ directory structure)?
