# PRD: AICoF CLI End-to-End Testing Suite

> **PRD Number**: 0001
> **Feature**: Comprehensive testing for aicof CLI
> **Status**: Draft
> **Created**: 2026-01-14

---

## Introduction

This PRD defines the testing strategy for the AICoF CLI tool (`aicof`). The CLI was recently implemented with 7 commands (init, update, add, remove, list, doctor, version) and needs comprehensive testing to ensure reliability before release.

**Goal**: Achieve >80% test coverage for business logic with both unit tests (fast, mocked) and end-to-end tests (real network) covering all commands and error scenarios.

---

## Goals

1. Create unit tests for all core packages with >80% coverage
2. Create E2E tests for all 7 CLI commands
3. Test both happy paths and error handling scenarios
4. Unit tests use mocked HTTP responses (fast, offline)
5. E2E tests use real GitHub API (realistic validation)
6. Ensure tests are maintainable and well-documented
7. Integrate tests into Makefile (`make test`)

---

## User Stories

**US-001**: As a developer, I want unit tests for core packages so that I can refactor code confidently without breaking functionality.

**US-002**: As a maintainer, I want E2E tests for all commands so that I can verify the CLI works correctly before releasing.

**US-003**: As a CI pipeline, I want fast unit tests that don't require network access so that builds complete quickly.

**US-004**: As a release manager, I want E2E tests that validate real GitHub downloads so that I know the CLI will work for users.

---

## Functional Requirements

### Unit Tests (Mocked Network)

**FR-001**: Test `internal/core/config.go`
- Parse valid YAML config
- Handle missing config file
- Handle malformed YAML
- Add/remove components correctly
- Save config preserves format

**FR-002**: Test `internal/core/registry.go`
- FindLanguage returns correct component
- FindFramework returns correct component
- FindWorkflow returns correct component
- FindTemplate returns correct template
- GetComponentPaths returns all expected paths
- Unknown component returns nil

**FR-003**: Test `internal/core/extractor.go`
- Extract single file correctly
- Extract directory recursively
- Skip existing files when force=false
- Overwrite existing files when force=true
- Handle missing source files
- Validate extraction results
- Backup and restore files

**FR-004**: Test `internal/core/downloader.go`
- Download version to cache (mocked)
- Return cached version if exists
- Clear dev cache on each download
- Handle download errors

**FR-005**: Test `internal/github/client.go`
- Parse latest release JSON
- Handle 404 (no releases)
- Construct archive URLs correctly
- Construct branch URLs correctly
- GetLatestVersionOrBranch fallback logic

### E2E Tests (Real Network)

**FR-006**: Test `aicof version`
- Shows CLI version info
- Shows "no framework" when not initialized
- `--check` flag works

**FR-007**: Test `aicof init`
- Creates project directory
- Installs CLAUDE.md
- Installs selected language guides
- Installs all workflows
- Creates aicof.yaml config
- `--template minimal` works
- `--template full` works
- `--non-interactive` works
- `--force` overwrites existing files
- Error when already initialized without --force

**FR-008**: Test `aicof list`
- Shows installed components
- `--available` shows all components
- `--type languages` filters correctly
- Works without installation (shows warning)

**FR-009**: Test `aicof add`
- Adds language guide
- Adds framework guide
- Updates aicof.yaml
- Error for unknown component
- Error when not initialized
- Warns when already installed

**FR-010**: Test `aicof remove`
- Removes language guide file
- Updates aicof.yaml
- Error when not installed
- `--force` skips confirmation

**FR-011**: Test `aicof update`
- `--check` shows available updates
- `--diff` shows changes
- Updates files correctly
- Preserves modified files
- Creates backup of modified files
- `--force` overwrites modified files

**FR-012**: Test `aicof doctor`
- Reports healthy installation
- Detects missing config
- Detects missing CLAUDE.md
- Detects missing directories
- Detects missing components
- `--fix` repairs issues

### Error Handling Tests

**FR-013**: Test error scenarios
- Network timeout handling
- Invalid GitHub response
- Corrupted archive
- Disk full simulation (if possible)
- Permission denied errors
- Invalid command arguments

---

## Non-Goals

- ❌ Performance benchmarking (not in scope for initial testing)
- ❌ Fuzzing or security testing (separate effort)
- ❌ Windows-specific testing (CI handles cross-platform)
- ❌ Interactive prompt testing (difficult to automate)
- ❌ Load testing (CLI is single-user)

---

## Technical Considerations

### Tech Stack

- **Testing Framework**: Go standard `testing` package
- **Assertions**: `testify/assert` for cleaner assertions
- **HTTP Mocking**: `httptest` for unit tests
- **Temp Directories**: `os.MkdirTemp` for isolated tests
- **Table-Driven Tests**: Go idiom for multiple test cases

### Test File Organization

```
packages/cli/
├── internal/
│   ├── core/
│   │   ├── config.go
│   │   ├── config_test.go        # Unit tests
│   │   ├── registry.go
│   │   ├── registry_test.go      # Unit tests
│   │   ├── extractor.go
│   │   ├── extractor_test.go     # Unit tests
│   │   ├── downloader.go
│   │   └── downloader_test.go    # Unit tests (mocked)
│   ├── github/
│   │   ├── client.go
│   │   └── client_test.go        # Unit tests (mocked)
│   └── cmd/
│       └── *_test.go             # Command unit tests
├── e2e/
│   ├── e2e_test.go               # E2E test suite
│   ├── init_test.go              # Init command E2E
│   ├── update_test.go            # Update command E2E
│   ├── add_remove_test.go        # Add/Remove E2E
│   ├── list_test.go              # List command E2E
│   ├── doctor_test.go            # Doctor command E2E
│   └── version_test.go           # Version command E2E
└── testdata/
    ├── valid_config.yaml         # Test fixtures
    ├── invalid_config.yaml
    └── mock_responses/
        ├── latest_release.json
        └── tags.json
```

### Test Helpers

```go
// Helper to create temp directory with cleanup
func setupTestDir(t *testing.T) string {
    dir, err := os.MkdirTemp("", "aicof-test-*")
    require.NoError(t, err)
    t.Cleanup(func() { os.RemoveAll(dir) })
    return dir
}

// Helper to run CLI command and capture output
func runCLI(t *testing.T, args ...string) (stdout, stderr string, exitCode int)

// Helper to create mock HTTP server
func mockGitHubServer(t *testing.T) *httptest.Server
```

### Dependencies to Add

```go
require (
    github.com/stretchr/testify v1.8.4  // Already indirect, make direct
)
```

---

## Guardrails Affected

### Testing (CRITICAL)
- ✓ Coverage targets: >80% for business logic
- ✓ All public APIs have unit tests
- ✓ Edge cases explicitly tested (null, empty, boundary values)
- ✓ Test names describe behavior
- ✓ No test interdependencies (tests run in any order)

### Code Quality
- ✓ No file exceeds 300 lines (split test files by command)
- ✓ Test helper functions are reusable

### Security
- ✓ Tests don't expose secrets
- ✓ Temp directories cleaned up after tests

---

## Success Metrics

### Technical Metrics
- Test coverage >80% for `internal/core/`
- Test coverage >60% for `internal/cmd/`
- All E2E tests pass with real GitHub API
- Tests complete in <60 seconds (unit) + <5 minutes (E2E)
- Zero flaky tests

### Quality Metrics
- All 7 commands have E2E tests
- Error scenarios covered for each command
- Tests are readable and maintainable

---

## Implementation Estimate

### Complexity Analysis

| Component | Estimated Effort | Type |
|-----------|-----------------|------|
| config_test.go | ~500 lines | ATOMIC |
| registry_test.go | ~300 lines | ATOMIC |
| extractor_test.go | ~600 lines | FEATURE |
| downloader_test.go | ~400 lines | FEATURE |
| client_test.go | ~400 lines | FEATURE |
| E2E test suite | ~800 lines | FEATURE |
| Test helpers | ~200 lines | ATOMIC |

**Total**: ~3,200 lines of test code

### Recommended Approach

1. Start with unit tests for `core/` packages (fast feedback)
2. Add mocked GitHub client tests
3. Create E2E test infrastructure
4. Add E2E tests command by command
5. Run full test suite and fix issues

---

## Open Questions

1. **CI Integration**: Should E2E tests run on every PR or only on release?
   - **Recommendation**: Unit tests on every PR, E2E on release branches

2. **Test Data**: Should we use a separate test repository or mock data?
   - **Recommendation**: Mocks for unit tests, real repo for E2E

3. **Coverage Enforcement**: Should we fail CI if coverage drops?
   - **Recommendation**: Yes, enforce >80% for `internal/core/`

---

## Next Steps

1. Generate task list using `@.agent/workflows/generate-tasks.md`
2. Implement tests in order of priority
3. Integrate into CI pipeline
4. Document test running instructions in README

---

**PRD Created**: 2026-01-14
**Author**: Claude (AI Assistant)
