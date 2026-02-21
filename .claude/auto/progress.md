# Auto Pilot Progress

## Discovery Log

[2026-02-21T15:51:24Z] [discovery] FOUND: Initial discovery iteration — analyzed full codebase

### Test Coverage Gaps
- `internal/github/` — 0% coverage (client.go: 279 LOC, handles GitHub API)
- `internal/ui/` — 0% coverage (3 files, 418 LOC total)
- `internal/core/skill.go` — 431 LOC, no tests (skill parsing, validation, indexing)
- `internal/core/downloader.go` — 262 LOC, no tests (archive download/extraction)
- `internal/core/auto_tasks.go` — 213 LOC, no tests (task state management)
- `internal/core/auto_prompt.go` — 104 LOC, no tests (prompt generation)
- `internal/commands/` — 76% of files untested (16/21 source files)

### Code Quality Violations
- `runDoctor()` in doctor.go: 362 lines (limit: 50)
- `runInit()` in init.go: 277 lines (limit: 50)
- 5 non-test files exceed 300-line limit: doctor.go (468), init.go (434), skill.go (431), config.go (425), extractor.go (378)

### Error Handling Issues
- `init.go:250` — error from `ScanSkillsDirectory` silently discarded
- `doctor.go:95` — error from `os.ReadFile` silently discarded
- `search.go:57` — error from `LoadConfig` silently discarded

### Dead Code
- `doctor.go:398` — `_ = extractor` silences unused variable warning

### Positive Findings
- All tests pass (`go test ./...`)
- No `go vet` warnings
- No TODO/FIXME/HACK markers in code
- No unused imports
- Good error wrapping patterns with `fmt.Errorf(%w)`
- Comprehensive table-driven tests where tests exist

## Iteration Log

[2026-02-21T16:00:00Z] [iteration:1] [task:1] COMPLETED: Fixed silent error handling in init.go, doctor.go, and search.go
- init.go: `ScanSkillsDirectory` error now logged via `ui.Warn`
- doctor.go: `os.ReadFile` error now logged via `ui.Warn`
- search.go: `LoadConfig` error logged via `ui.Warn` only for non-file-not-found errors (since missing config is expected when project isn't initialized)
- LEARNING: For search.go, `LoadConfig` returns `os.ErrNotExist` when no config file exists — this is a normal case (user hasn't run `samuel init` yet), so only warn on other error types
- Commit: c3c5562

[2026-02-21T16:30:00Z] [iteration:2] [task:3] COMPLETED: Added unit tests for internal/github/client.go
- Created client_test.go with 28 test cases covering all exported functions
- Used `redirectTransport` pattern (custom `http.RoundTripper`) to redirect HTTP requests to `httptest.NewServer` — zero changes to production code
- Coverage: 89.7% (0% → 89.7%)
- Table-driven tests for: GetLatestRelease, GetLatestVersionOrBranch, GetTags, DownloadArchive, DownloadBranchArchive, DownloadFile, CheckForUpdates
- Also tested: NewClient, URL construction (GetArchiveURL, GetBranchArchiveURL), request headers, constants
- LEARNING: `CheckForUpdates` has a nil pointer dereference if `GetLatestRelease` returns nil (404 case). The function accesses `release.TagName` without nil check. Not fixed in this task — should be a separate bug fix.
- LEARNING: The `redirectTransport` pattern is reusable for testing any HTTP client code without modifying production code. The transport intercepts `RoundTrip` and rewrites `URL.Host` to point at the test server.

[2026-02-21T17:00:00Z] [iteration:3] [task:4] COMPLETED: Added unit tests for internal/core/skill.go
- Created skill_test.go with 50 test cases covering all 14 functions
- Coverage: 95.8% (0% → 95.8%) — 11 functions at 100%, ParseSkillMD at 95%, LoadSkillInfo at 94.7%, ScanSkillsDirectory at 89.5%, CreateSkillScaffold at 75%
- Used t.TempDir() for all filesystem tests — clean, no manual cleanup needed
- Table-driven tests for: toTitleCase, ValidateSkillName, ValidateSkillDescription, ValidateSkillCompatibility, ValidateSkillMetadata, ParseSkillMD, CountLines
- Subtests for filesystem-dependent functions: LoadSkillInfo, ScanSkillsDirectory, CreateSkillScaffold, UpdateCLAUDEMDSkillsSection, dirExists
- LEARNING: ScanSkillsDirectory correctly skips hidden dirs (prefixed with `.`) and non-directory entries — tested both cases
- LEARNING: ValidateSkillName with uppercase chars triggers TWO errors (lowercase check + invalid chars check) because uppercase letters fail both conditions
- LEARNING: ParseSkillMD splitting on empty string returns `[""]` (1 element), not `[]` (0 elements), so the `len(lines) == 0` guard in the function is effectively dead code
- Commit: 16771cf

[2026-02-21T17:30:00Z] [iteration:4] [task:5] COMPLETED: Added unit tests for internal/core/auto_tasks.go
- Created auto_tasks_test.go with 12 test functions (including table-driven tests)
- Coverage: 100% on all 13 functions (previously ~90% from tests scattered in auto_test.go)
- Filled gaps: priorityRank default case, SkipTask/ResetTask not-found paths, validateTasks empty-ID edge case
- Table-driven tests for: priorityRank (6 cases), allDependenciesMet (5 cases), isValidStatus (9 cases)
- Subtests for: getAvailableTasks (unmet deps, in-progress exclusion), findTask not-found, AddTask status preservation
- LEARNING: auto_test.go already had substantial tests covering auto_tasks.go functions, but with small gaps. The new file fills those gaps to achieve 100%
- LEARNING: validateTasks uses `continue` after detecting empty ID, so further checks (title, status) are skipped for that task — only 1 error produced
- Commit: 3c73231

### Tasks Generated (Initial Discovery): 10
| ID | Priority | Title |
|----|----------|-------|
| 1  | high     | Fix silent error handling in init.go, doctor.go, and search.go |
| 2  | medium   | Remove unused extractor variable in doctor.go |
| 3  | high     | Add unit tests for internal/github/client.go |
| 4  | high     | Add unit tests for internal/core/skill.go |
| 5  | high     | Add unit tests for internal/core/auto_tasks.go |
| 6  | medium   | Refactor runDoctor() into smaller helper functions |
| 7  | medium   | Refactor runInit() into smaller helper functions |
| 8  | high     | Add unit tests for internal/core/downloader.go |
| 9  | medium   | Add unit tests for internal/core/auto_prompt.go |
| 10 | low      | Reduce file size of internal/core/skill.go below 300-line limit |

[2026-02-21T18:00:00Z] [iteration:5] [task:11] COMPLETED: Fixed symlink traversal vulnerability in tar extraction
- Added `validateSymlinkTarget()` function that rejects: (1) absolute symlink targets, (2) relative targets that resolve outside the destination directory
- Validation uses `filepath.Join(filepath.Dir(symlinkPath), linkTarget)` to resolve relative targets against the symlink's parent directory, then checks containment with `strings.HasPrefix`
- Created `downloader_test.go` with 12 test cases covering: symlink traversal attacks, absolute symlink rejection, valid symlinks, basic extraction, and path traversal in file names
- The existing path traversal check (line 140) already protects file *names* — this fix closes the gap for symlink *targets*
- LEARNING: Symlink traversal is a two-step attack: (1) create a symlink pointing outside dest, (2) write files through the symlink. Validating `target` (the symlink location) alone is insufficient — `Linkname` (where it points) must also be validated
- LEARNING: `filepath.Clean` resolves `../` components, so `filepath.Join(dir, "../../etc") → filepath.Clean(result)` gives the actual resolved path for containment checking
- Commit: 1cf9484

---

[2026-02-21T17:45:00Z] [discovery] FOUND: Second discovery iteration — deep security & quality analysis

### Security Vulnerabilities (NEW)
- **CRITICAL**: `downloader.go:124-133` — Symlink targets not validated in tar extraction; malicious archives can escape destination via symlink traversal
- **CRITICAL**: `client.go:216-228` — Nil pointer dereference in `CheckForUpdates` when `GetLatestRelease` returns `(nil, nil)` (no releases case)
- **HIGH**: `extractor.go:216-248` — `ReadFile`, `WriteFile`, `RemoveFile`, `BackupFile` accept paths without traversal validation; `../../etc/passwd` would escape destination
- **HIGH**: `downloader.go:117-121` — Unbounded `io.Copy` in tar extraction (decompression bomb risk); `client.go:204-206` — unbounded `io.ReadAll` in `DownloadFile` (OOM risk)

### Additional Error Handling Issues (NEW)
- `diff.go:201` — `_ = filepath.Walk(agentDir, ...)` silently discards walk errors
- `diff.go:223` — `_ = filepath.Walk(templatePath, ...)` silently discards walk errors
- `diff.go:192` — `_ = filepath.Glob(...)` discards error from malformed patterns

### Test Coverage Update
- Overall coverage: **37.5%** (below 60% target)
- `internal/core` at **65.9%**, `internal/github` at **89.7%**
- `internal/ui/` still at **0%** — `output.go` (137 LOC) is most testable
- `internal/core/extractor.go` at **~15%** (13 of 17 functions untested)
- `internal/core/registry.go` at **~85%** but 6 skill lookup functions at 0%

### Code Quality Violations (NEW)
- `runUpdate()` in update.go: **220 lines** (4.4x the 50-line limit) — worst new violation
- `runRemove()`: 98 lines, `runAdd()`: 97 lines, `executePilotLoop()`: 95 lines
- `sync.go`: **431 lines** (exceeds 300-line file limit)
- `commands/skill.go`: **375 lines**, `search.go`: **337 lines** (exceed limits)
- 45+ raw `fmt.Print*` calls in commands bypassing `ui` package abstraction
- 8 uses of deprecated `filepath.Walk` (should be `filepath.WalkDir` since Go 1.16)

### Positive Findings
- `go vet ./...` still clean
- No TODO/FIXME/HACK markers in code
- Good test quality where tests exist (table-driven, t.TempDir patterns)

[2026-02-21T19:00:00Z] [iteration:6] [task:12] COMPLETED: Fixed nil pointer dereference in CheckForUpdates
- Added nil guard: `if release == nil { return nil, fmt.Errorf(...) }` before accessing `release.TagName`
- Added regression test `no_releases_returns_error` to TestCheckForUpdates table-driven tests
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: `GetLatestRelease` returns `(nil, nil)` on 404 — this is by design (documented in godoc). All callers must check for nil. `GetLatestVersionOrBranch` already had the nil check; only `CheckForUpdates` was missing it.
- Commit: e484979

### Tasks Generated (Second Discovery): 10
| ID | Priority | Title |
|----|----------|-------|
| 11 | critical | Fix symlink traversal vulnerability in tar extraction |
| 12 | critical | Fix nil pointer dereference in CheckForUpdates |
| 13 | high     | Add path traversal validation to extractor.go file operations |
| 14 | high     | Add size limits to tar extraction and file downloads |
| 15 | medium   | Fix silently discarded filepath.Walk errors in diff.go |
| 16 | high     | Add unit tests for internal/core/extractor.go |
| 17 | medium   | Add unit tests for internal/ui/output.go |
| 18 | medium   | Refactor runUpdate() into smaller helper functions |
| 19 | medium   | Add unit tests for registry.go skill lookup functions |
| 20 | low      | Reduce file size of internal/core/sync.go below 300-line limit |

[2026-02-21T19:30:00Z] [iteration:7] [task:8] COMPLETED: Added unit tests for internal/core/downloader.go filesystem functions
- Added 22 test cases: 6 extractTarGz edge cases, 3 copyFile, 3 copyDir, 3 ClearCache, 2 GetCacheSize
- Coverage: validateSymlinkTarget 100%, copyDir 90%, copyFile 92.3%, ClearCache 80%, GetCacheSize 87.5%, extractTarGz 75.8%
- Thin wrapper methods (NewDownloader, DownloadVersion, GetLatestVersion, DownloadFile, CheckForUpdates) left at 0% — they delegate to github.Client which is tested at 89.7%
- extractTarGz edge cases: invalid gzip input, empty archive, nested dirs without explicit entries, file permissions preserved, directory path traversal
- LEARNING: `extractTarGz` auto-creates parent directories for regular files via `os.MkdirAll(filepath.Dir(target))`, so archives don't need explicit directory entries — tested this behavior explicitly
- LEARNING: The `contains` helper in the test file uses `bytes.Contains` instead of `strings.Contains` — works identically but could be simplified. Not worth changing since it's existing test code.
- Commit: e07b37c

[2026-02-21T20:00:00Z] [iteration:8] [task:13] COMPLETED: Added path traversal validation to extractor.go file operations
- Added `validateContainedPath(baseDir, relativePath)` helper using `filepath.Clean` + `strings.HasPrefix` containment check
- Applied to 6 functions: `ReadFile`, `WriteFile`, `RemoveFile`, `BackupFile`, `FileExists`, `ValidateExtraction`
- `BackupFile` validates both source path (against `destPath`) and destination path (against `backupDir`)
- Added 11 regression tests: `TestValidateContainedPath` (8 table-driven cases), plus individual tests for ReadFile, WriteFile, RemoveFile, FileExists, BackupFile, ValidateExtraction with traversal paths
- Also added 2 positive tests (ReadFile_ValidPath, WriteFile_ValidPath) to verify normal operation isn't broken
- LEARNING: On Unix, `filepath.Join("/base", "/absolute")` does NOT replace the base — it produces `/base/absolute`. So absolute paths as the second arg to `filepath.Join` don't actually escape on Unix (they're treated as relative). The containment check naturally passes for this case.
- LEARNING: The `validateContainedPath` pattern is the same one used in `downloader.go` for symlink validation (task 11) — the project now has consistent path containment checks across both tar extraction and file operations.
- Commit: 8b69661

---

[2026-02-21T21:00:00Z] [discovery] FOUND: Third discovery iteration — TOCTOU security, error handling, coverage, and file size analysis

### Security Vulnerabilities (NEW)
- **CRITICAL**: `auto_loop.go:116` — TOCTOU in AITool validation. `cfg.AITool` from PRD file is passed to `exec.Command` without re-validation. PRD is re-read from disk each iteration (line 74), so file modification between CLI parse and loop execution bypasses the `IsValidAITool()` check in `parseAutoFlags()`.
- **HIGH**: `extractor.go:319-338` — `RestoreBackup` computes `dstPath = filepath.Join(e.destPath, relPath)` without calling `validateContainedPath`. Task 13 covered ReadFile/WriteFile/RemoveFile/BackupFile but missed RestoreBackup and CopyFromCache.

### Error Handling Issues (NEW)
- **HIGH**: 4 `LoadAutoPRD` errors silently discarded: `auto_pilot.go:196`, `auto_pilot.go:225`, `auto_start_handler.go:124`, `auto_pilot_output.go:45`
- **MEDIUM**: 3 `filepath.Rel` errors silently discarded in `extractor.go` (lines 87, 104, 117) — empty `relPath` on failure causes incorrect file path construction

### Test Coverage Update
- Overall coverage: `internal/commands` **19.4%**, `internal/core` **73.0%**, `internal/github` **89.9%**, `internal/ui` **0.0%**
- `auto_loop.go`: 5 of 9 functions at 0% (InvokeAgent, invokeAgentLocal, invokeAgentDocker, invokeAgentDockerSandbox, buildDockerRunArgs)
- `config.go`: GetGlobalConfigPath (0%), GetCachePath (0%), EnsureCacheDir (0%), GetValue (45%)

### Code Quality Update
- 19 functions exceed 50-line limit (top 3: runDoctor 366, runInit 281, runUpdate 219)
- 9 non-test files exceed 300-line limit: doctor.go (471), registry.go (460), init.go (437), sync.go (431), skill.go (431), config.go (425), extractor.go (413), skill.go (375), search.go (337)
- 8 deprecated `filepath.Walk` calls (should be `filepath.WalkDir` since Go 1.16)
- No panic() calls, no TODOs/FIXMEs, go vet clean

### Positive Findings
- `go test ./...` all pass, `go vet ./...` clean, `go build ./...` clean
- No new security suppression comments (#nosec)
- HTTP client correctly uses 30s timeout (not http.DefaultClient)
- No hardcoded credentials or tokens
- Good test patterns established (table-driven, t.TempDir, redirectTransport)

### Tasks Generated (Third Discovery): 10
| ID | Priority | Title |
|----|----------|-------|
| 21 | critical | Fix TOCTOU vulnerability in AITool validation for auto loop |
| 22 | high     | Add path containment validation to RestoreBackup and CopyFromCache |
| 23 | high     | Fix silently discarded LoadAutoPRD errors in auto pilot commands |
| 24 | medium   | Replace deprecated filepath.Walk with filepath.WalkDir across codebase |
| 25 | medium   | Fix silently discarded filepath.Rel errors in extractor.go |
| 26 | high     | Add unit tests for auto_loop.go agent invocation functions |
| 27 | medium   | Add unit tests for config.go path and value functions |
| 28 | low      | Reduce file size of internal/core/config.go below 300-line limit |
| 29 | low      | Reduce file size of internal/core/registry.go below 300-line limit |
| 30 | low      | Reduce file size of internal/core/extractor.go below 300-line limit |

[2026-02-21T21:30:00Z] [iteration:9] [task:21] COMPLETED: Fixed TOCTOU vulnerability in AITool validation
- Added `IsValidAITool(cfg.AITool)` check at the top of `InvokeAgent()` before dispatching to local/docker/docker-sandbox execution paths
- The validation is a single point of defense that covers all three execution modes
- Added 9 regression tests: 5 invalid tool cases (path traversal, shell injection, empty, arbitrary binary, unknown tool) and 4 valid tool acceptance tests
- LEARNING: The described TOCTOU was partially inaccurate — `cfg.AITool` is set once in `NewLoopConfig` and not updated when PRD is re-read in the loop. However, in `auto_start_handler.go`, `NewLoopConfig` reads `AITool` from `prd.Config.AITool` without any validation (only sandbox mode is validated). The fix at `InvokeAgent` level provides defense-in-depth regardless of how `cfg.AITool` is populated.
- LEARNING: `exec.Command` with Go does NOT invoke a shell, so shell metacharacters like `;` in the tool name won't cause injection. But an arbitrary binary path (e.g., `/bin/sh`) would still be executed. The allow-list validation prevents both attack vectors.
- Commit: 4db5316

[2026-02-21T22:00:00Z] [iteration:10] [task:14] COMPLETED: Added size limits to tar extraction and file downloads
- Added `MaxExtractedFileSize` (100MB) as a package-level var in `downloader.go` — applied via `io.LimitReader(tarReader, MaxExtractedFileSize+1)` in `extractTarGz`, checked after copy
- Added `MaxDownloadFileSize` (10MB) as a package-level var in `client.go` — applied via `io.LimitReader(resp.Body, MaxDownloadFileSize+1)` in `DownloadFile`, checked after read
- Both use the `LimitReader(r, max+1)` + `n > max` pattern: read up to max+1 bytes, then check if we hit the boundary. Files exactly at the limit succeed; files exceeding it fail.
- Used `var` instead of `const` for the limits so tests can temporarily override them (test isolation with save/restore pattern)
- Added 4 regression tests: `TestExtractTarGz_FileSizeLimit` (oversized rejected), `TestExtractTarGz_FileSizeAtLimit` (exact size OK), `TestDownloadFile_SizeLimit/oversized_file_rejected`, `TestDownloadFile_SizeLimit/file_at_limit_succeeds`
- LEARNING: The `io.LimitReader(r, n+1)` + `bytesRead > n` pattern is preferable to `io.LimitReader(r, n)` because the latter silently truncates without error — you can't distinguish between a file exactly `n` bytes and one larger than `n`. The +1 approach lets you detect the overflow.
- LEARNING: Using `var` for size limits instead of `const` is the right tradeoff — it enables test isolation while keeping the production default immutable in practice. The save/restore `defer` pattern prevents test pollution.
- Commit: 050f9f6

[2026-02-21T22:30:00Z] [iteration:11] [task:16] COMPLETED: Added unit tests for internal/core/extractor.go
- Added 22 new test cases covering all previously untested functions
- Created `createTemplateFile` helper to DRY up test setup for template directory structures
- Test coverage for extractor.go: 10 functions at 100%, 9 functions at 73-85% (up from ~15% with 13 of 17 functions at 0%)
- Tests cover: shouldSkip (9 table-driven cases), NewExtractor/GetSourcePath/GetDestPath, Extract (single file, skip existing, force overwrite, source not found, directory, multiple files), ExtractAll (with files, skips .git/node_modules, no template dir), RestoreBackup (nested files, empty dir), RemoveFile/FileExists/BackupFile/ValidateExtraction (positive-path tests)
- Overall `internal/core` package coverage: 80.4%
- LEARNING: `shouldSkip` treats `.github` and `.gitignore` as skip targets because it uses `strings.HasPrefix(path, ".git")` — this is intentional for the template directory context where only `.claude/` directories should exist
- LEARNING: `ExtractAll` handles missing template directory gracefully via `filepath.SkipAll` in the Walk callback — tested with `TestExtractAll_NoTemplateDir`
- Commit: fbc447e

[2026-02-21T23:00:00Z] [iteration:12] [task:22] COMPLETED: Added path containment validation to RestoreBackup and CopyFromCache
- `RestoreBackup`: replaced `dstPath := filepath.Join(e.destPath, relPath)` with `dstPath, err := validateContainedPath(e.destPath, relPath)` — now rejects backup entries that would write outside the destination directory
- `CopyFromCache`: replaced `dstPath := filepath.Join(destPath, filePath)` with `dstPath, err := validateContainedPath(destPath, filePath)` — now rejects file paths that escape the destination
- Added 3 regression tests: `TestRestoreBackup_PathTraversal` (validates normal operation still works with validation in place), `TestCopyFromCache_PathTraversal` (file traversal rejected), `TestCopyFromCache_PathTraversal_Directory` (directory traversal rejected)
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: `RestoreBackup` uses `filepath.Walk` which generates paths from the actual filesystem tree, so `filepath.Rel` will produce clean relative paths for real files. The traversal risk is more theoretical (crafted backup directory or symlink following), but the validation provides defense-in-depth consistent with ReadFile/WriteFile/RemoveFile/BackupFile.
- LEARNING: `CopyFromCache` is a standalone function (not a method on Extractor) — `validateContainedPath` works equally well for standalone functions since it only needs a base dir and relative path.

---

[2026-02-21T23:30:00Z] [discovery] FOUND: Fourth discovery iteration — race condition, test coverage, error handling

### Security Issues (NEW)
- **HIGH**: `internal/ui/spinner.go:35-50` — Data race between Start() goroutine reading `s.bar` and Stop() setting `s.bar = nil`. No synchronization protects the field. Could cause nil pointer dereference under concurrent access.

### Test Coverage Update
- Overall coverage: `cmd/samuel` **0%**, `internal/commands` **19.4%**, `internal/core` **80.4%**, `internal/github` **89.4%**, `internal/ui` **0.0%**
- `internal/commands/` has 16 of 21 source files with 0% test coverage
- Highest-value untested command files: `auto_handlers.go` (pure functions: detectQualityChecks, countTaskStatuses, validateSandbox), `auto_task_handlers.go` (taskStatusIcon), `init.go` (expandLanguages, expandFrameworks, isSamuelRepository), `doctor.go` (extractVersion)
- All tests passing, `go vet` clean, `go build` clean

### Error Handling Issues (NEW)
- **MEDIUM**: `sync.go:221` — `rel, _ := filepath.Rel(opts.RootDir, path)` silently discards error in SyncFolderCLAUDEMDs — could produce incorrect relative paths
- **MEDIUM**: `diff.go:192,206,233` — Three `filepath.Rel` calls silently discard errors — incorrect hash map keys could cause wrong diff results
- **MEDIUM**: `list.go:122` — `config, _ := core.LoadConfig()` discards non-file-not-found errors (corrupt YAML, permission denied)

### Code Quality Update
- 9 non-test files exceed 300-line limit (no change from previous discovery)
- 2 command files exceed 300-line limit that don't have pending refactor tasks: `skill.go` (375 lines), `search.go` (337 lines)
- No new panic() calls, no hardcoded credentials, no TODO/FIXME markers
- File permissions consistently 0755 (dirs) and 0644 (files)

### Positive Findings
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- Security posture improved significantly since previous discoveries (path traversal, symlink, size limits, TOCTOU all fixed)
- Good test quality patterns maintained (table-driven, t.TempDir, redirectTransport)
- No new environment variable leaks or temp file issues
- HTTP client properly configured with 30s timeout and size limits
- All exec.Command calls validated against allowlist

### Tasks Generated (Fourth Discovery): 10
| ID | Priority | Title |
|----|----------|-------|
| 31 | high     | Fix race condition in ui/spinner.go Start/Stop |
| 32 | high     | Add unit tests for auto_handlers.go pure functions |
| 33 | medium   | Add unit tests for auto_task_handlers.go functions |
| 34 | high     | Add unit tests for init.go helper functions |
| 35 | medium   | Add unit tests for doctor.go helper functions |
| 36 | medium   | Fix silently discarded filepath.Rel error in sync.go |
| 37 | medium   | Fix silently discarded filepath.Rel errors in diff.go |
| 38 | low      | Reduce file size of commands/skill.go below 300-line limit |
| 39 | low      | Reduce file size of commands/search.go below 300-line limit |
| 40 | medium   | Fix silently discarded LoadConfig error in list.go |

[2026-02-21T23:45:00Z] [iteration:13] [task:23] COMPLETED: Fixed silently discarded LoadAutoPRD errors in auto pilot commands
- Fixed 4 locations: auto_pilot.go (2), auto_start_handler.go (1), auto_pilot_output.go (1)
- All 4 now log `ui.Warn` when `LoadAutoPRD` returns a non-nil error, preserving the error context
- Existing nil-check behavior kept intact — the code still degrades gracefully when PRD is nil
- Mid-loop locations (auto_pilot.go:196, :225): warn + continue with existing nil handling. These affect discovery task counting and empty-discovery loop termination.
- Summary locations (printLoopSummary, printPilotSummary): warn + early return on nil (unchanged behavior, just with error logging)
- LEARNING: All 4 call sites already had nil guards (`if finalPRD == nil`), so the behavior is safe even without the error. The issue was purely about error visibility — a corrupt prd.json or permission error would be silently swallowed, making debugging difficult.
- Commit: d7339f9

[2026-02-22T00:00:00Z] [iteration:14] [task:26] COMPLETED: Added unit tests for auto_loop.go agent invocation functions
- Added 15 new test cases: buildDockerRunArgs (3 tests), InvokeAgent dispatch (4 modes), invokeAgentLocal (2), invokeAgentDocker (2), invokeAgentDockerSandbox (2), RunAutoLoop consecutive failures (1), RunAutoLoop callbacks (1)
- Coverage improvement: InvokeAgent 100%, invokeAgentLocal 100%, invokeAgentDockerSandbox 100%, buildDockerRunArgs 100%, invokeAgentDocker 86.7%, RunAutoLoop 90%
- Overall `internal/core` package coverage: 83.9% (up from ~80.4%)
- Key technique: For dispatch tests, use claude tool with non-existent prompt file to fail fast at GetAgentArgs stage — avoids slow docker container pulls while still verifying the switch dispatches to the correct code path
- For docker invoke tests, use `nonexistent-image-test:0.0.0` to fail fast on image pull instead of pulling real images
- LEARNING: Docker sandbox (`docker sandbox run`) reuses existing sandboxes by name. If a sandbox was previously created in the same test run, subsequent calls to the same sandbox name are fast. But initial creation pulls images and is very slow (23+ seconds). Tests should avoid triggering real sandbox creation.
- LEARNING: `invokeAgentDocker` has an unreachable error path: `filepath.Rel(projectDir, promptPath)` on Unix only fails when the paths are on different volumes (not possible on Unix). The 86.7% coverage gap is this error path. This is acceptable — the error check is defense-in-depth for potential Windows support.
- Commit: 1cf5627

[2026-02-22T00:30:00Z] [iteration:15] [task:31] COMPLETED: Fixed race condition in ui/spinner.go Start/Stop
- Replaced unsynchronized `s.bar` nil check with `done chan struct{}` + `sync.Once` pattern
- `Start()` goroutine now uses `select` on `done` channel + `time.Ticker` instead of reading a shared pointer
- `Stop()` uses `sync.Once` to close the `done` channel exactly once — safe for concurrent and repeated calls
- Added 7 regression tests: basic start/stop, idempotent stop, stop-without-start, concurrent stop (10 goroutines), rapid start/stop, Success(), Error()
- All tests pass with `-race` flag — no data races detected
- LEARNING: The `select` on a closed channel approach (`select { case <-s.done: return; default: close(s.done) }`) is NOT safe for concurrent callers — two goroutines can both enter the `default` branch before either closes. `sync.Once` is the correct primitive for one-time channel close.
- LEARNING: `time.NewTicker` is preferred over `time.Sleep` in goroutine loops because the ticker can be deferred-stopped, ensuring cleanup when the goroutine exits via the `done` channel. The `select` pattern with ticker also avoids busy-waiting during the sleep interval when shutdown is requested.
- Commit: 0819b18

[2026-02-22T01:00:00Z] [iteration:16] [task:32] COMPLETED: Added unit tests for auto_handlers.go pure functions
- Created auto_handlers_test.go with 19 test cases covering 3 functions
- `detectQualityChecks`: 9 table-driven tests (go.mod, package.json, Cargo.toml, requirements.txt, empty dir, priority ordering for go>node>rust>python, unrecognized file)
- `countTaskStatuses`: 5 table-driven tests (empty, all pending, all completed, mixed statuses, unknown status key)
- `validateSandbox`: 5 subtests (none, empty, unrecognized mode, docker, docker-sandbox)
- LEARNING: `detectQualityChecks` uses early-return pattern — the first matching project file wins. Priority order: go.mod > package.json > Cargo.toml > requirements.txt. Multiple project files present = first match takes priority.
- LEARNING: `countTaskStatuses` initializes 5 known status keys but unknown statuses (e.g., "unknown") also get counted via map increment — the function doesn't validate against known statuses. The `unknown_status_counted_separately` test verifies this behavior.
- LEARNING: `validateSandbox` only checks docker availability for "docker" and "docker-sandbox" modes. Any other string (including empty) passes through without error since neither `if` branch matches.
- Commit: 66435a9

[2026-02-22T01:30:00Z] [iteration:17] [task:34] COMPLETED: Added unit tests for init.go helper functions
- Created init_test.go with 44 test cases covering 4 functions
- `expandLanguages`: 18 table-driven tests (all 10 aliases, comma-separated, multiple flags, unknown filtered, uppercase normalized, whitespace trimmed)
- `expandFrameworks`: 10 table-driven tests (2 aliases, comma-separated, multiple flags, unknown filtered, uppercase normalized, whitespace trimmed)
- `isSamuelRepository`: 6 subtests using t.TempDir() (empty dir, template dir without CLAUDE.md, template dir with CLAUDE.md, packages/cli with go.mod, packages/cli without go.mod, nonexistent dir)
- `getRelevantFrameworks`: 10 subtests (3 single-language mapping tests, multi-language dedup, unknown/empty/nil input, Component type validation, all 11 mapped languages verified, registry consistency check)
- LEARNING: `expandLanguages` maps both "js" and "javascript" to "typescript" — the alias "ts" also maps there. Three aliases converge to one language.
- LEARNING: `isSamuelRepository` has two independent detection paths (template/CLAUDE.md OR packages/cli/go.mod). Either match triggers detection. The "packages/cli" check is for a legacy or alternative repo structure.
- LEARNING: `getRelevantFrameworks` uses a `seen` map to deduplicate — if two languages share a framework, it only appears once. The mapping covers 11 languages × 3 frameworks each = 33 framework mappings.
- Commit: 45c720b

[2026-02-22T02:00:00Z] [iteration:18] [task:41] COMPLETED: Added path traversal validation to remove.go component deletion
- Added `validateRemovePath(projectDir, componentPath)` helper with `filepath.Clean` + `strings.HasPrefix` containment check
- Applied before `os.Remove` call in `runRemove` — rejects component paths that resolve outside the project directory
- Added 10 regression tests: 8 table-driven cases (valid paths, traversal attacks, edge cases) + 2 TempDir-based tests
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: `component.Path` actually comes from the hard-coded registry (`core.FindLanguage/Framework/Workflow`), not directly from `samuel.yaml` user input. The risk is more theoretical (crafted registry entry or future code change), but the validation provides defense-in-depth consistent with `validateContainedPath` in `extractor.go`.
- LEARNING: The `validateRemovePath` pattern is identical to `validateContainedPath` in extractor.go. In a future cleanup iteration, these could be consolidated into a shared utility function in the `core` package by exporting `ValidateContainedPath`.
- Commit: ec2d38b

---

[2026-02-21T01:30:00Z] [discovery] FOUND: Fifth discovery iteration — security, coverage, error handling

### Security Vulnerabilities (NEW)
- **HIGH**: `remove.go:108` — `filePath := filepath.Join(cwd, component.Path)` uses `component.Path` from `samuel.yaml` without path traversal validation. A tampered config could cause `os.Remove` to delete arbitrary files outside the project directory. No `validateContainedPath` equivalent is applied.
- **MEDIUM**: `skill.go:148` — User-provided skill name from CLI used directly in `filepath.Join(skillsDir, name)`. While `ValidateSkillName` restricts format, the path is constructed *before* validation in the create command flow (line 134). A name like `../../etc` could escape the skills directory.

### Test Coverage Update
- Overall: `cmd/samuel` **0%**, `internal/commands` **20.3%**, `internal/core` **83.9%**, `internal/github` **89.4%**, `internal/ui` **17.4%**
- `internal/commands/` still has 15 of 21 source files with 0% test coverage
- New testable files identified: `config_cmd.go` (2 pure functions: `formatConfigValue`, `isValidConfigKey`), `diff.go` (pure function: `computeDiff`), `auto_pilot.go` (flag parsing functions)
- `internal/ui` improved from 0% to 17.4% (spinner tests added in task 31), but `output.go` and `prompts.go` still untested

### Error Handling Issues (NEW)
- **MEDIUM**: `info.go:62` — `config, _ := core.LoadConfig()` silently discards error (same pattern as task 40 for list.go)
- **LOW**: `config_cmd.go:168` — `oldValue, _ := config.GetValue(key)` silently discards error (display-only, non-critical)

### Code Quality Update
- 19 functions exceed 50-line limit (unchanged from previous discovery)
- 9 non-test files exceed 300-line limit (unchanged)
- `auto_pilot.go` at 296 lines — just under 300-line limit but `executePilotLoop` is 100 lines (2x function limit)
- `go vet ./...` clean, no panic() calls, no TODO/FIXME markers
- No hardcoded secrets or credentials found
- All `exec.Command` calls protected by AITool allowlist
- HTTP client properly uses custom client with 30s timeout (no `http.DefaultClient`)

### Positive Findings
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- Security posture strong: path traversal, symlink, size limits, TOCTOU all fixed in previous iterations
- `internal/core` at 83.9% exceeds 80% business logic target
- `internal/github` at 89.4% well above target
- Established test patterns (table-driven, t.TempDir, redirectTransport) consistently applied
- Docker env var forwarding uses explicit allowlist (not forwarding all vars)
- Decompression bomb and download size protections in place

[2026-02-22T02:30:00Z] [iteration:19] [task:2] COMPLETED: Removed unused extractor variable in doctor.go
- Removed `extractor := core.NewExtractor(cachePath, cwd)` on line 390 and `_ = extractor` on line 401
- The code uses `core.CopyFromCache(cachePath, cwd, path)` directly — the extractor was never needed
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: The `core` import is still needed (13 other usages in the file), so no import cleanup was required
- Commit: 0d77327

[2026-02-22T03:00:00Z] [iteration:20] [task:6] COMPLETED: Refactored runDoctor() into smaller helper functions
- Split 363-line `runDoctor()` into a 43-line orchestrator that delegates to focused functions
- Created `doctor_checks.go` (291 lines) with individual check functions: `checkConfigFile`, `checkCLAUDEMD`, `checkAGENTSMD`, `checkDirectoryStructure`, `checkInstalledComponents`, `checkSkillsIntegrity`, `checkAutoHealth`, `checkLocalModifications`
- Kept `doctor.go` (175 lines) with command registration, `runDoctor`, `printCheckResults`, `printCheckSummary`, `performAutoFix`, `restoreMissingComponents`
- Introduced `checkInstalledSkills` generic helper to deduplicate the language/framework/workflow check pattern (originally 3 near-identical code blocks)
- All functions under 50 lines, both files under 300 lines
- `checkModification` signature simplified: removed unused `_ string` parameter (was `config.Version` but never used)
- LEARNING: The three component checks (language, framework, workflow) follow an identical pattern: migrate legacy config, iterate names, check SKILL.md existence. A generic function with a `finder func(string) *core.Component` parameter cleanly abstracts this. The only special case is workflows expanding `"all"` to `GetAllWorkflowNames()`, handled in `checkInstalledComponents` before calling the generic helper.
- LEARNING: Splitting into two files (orchestrator + checks) keeps both under the 300-line limit while the original 470-line single file far exceeded it. The logical grouping (doctor.go = command setup + display + fix, doctor_checks.go = individual checks) is a natural split.
- Commit: b51ee55

### Tasks Generated (Fifth Discovery): 10
| ID | Priority | Title |
|----|----------|-------|
| 41 | high     | Add path traversal validation to remove.go component deletion |
| 42 | medium   | Add skill name validation against path traversal in skill.go |
| 43 | medium   | Fix silently discarded LoadConfig error in info.go |
| 44 | medium   | Add unit tests for config_cmd.go pure functions |
| 45 | low      | Fix silently discarded GetValue error in config_cmd.go |
| 46 | medium   | Add unit tests for auto_pilot.go parsePilotFlags and parseAutoFlags |
| 47 | medium   | Add unit tests for diff.go pure functions |
| 48 | medium   | Add unit tests for list.go helper functions |
| 49 | low      | Reduce file size of auto_pilot.go / refactor executePilotLoop |
| 50 | medium   | Add unit tests for sync.go GenerateFolderCLAUDEMD function |

---

[2026-02-22T04:00:00Z] [discovery] FOUND: Sixth discovery iteration — function size violations, error handling, test coverage

### Code Quality Violations (NEW)
- **runAdd()** in add.go: **97 lines** (1.9x the 50-line limit) — no existing refactor task
- **runRemove()** in remove.go: **101 lines** (2x the 50-line limit) — no existing refactor task
- **DownloadVersion()** in downloader.go: **65 lines** (1.3x the 50-line limit)
- **extractTarGz()** in downloader.go: **69 lines** (1.4x the 50-line limit)
- **runVersion()** in version.go: **68 lines** (1.4x the 50-line limit)
- downloader.go at 296 lines — borderline 300-line limit with two oversized functions

### Error Handling Issues (NEW)
- **MEDIUM**: `downloader.go:49` — `os.RemoveAll(cacheDest)` error silently discarded when clearing dev cache. Failed removal + subsequent download = corrupted state
- **MEDIUM**: `update.go:147-155` — Two `os.ReadFile` errors silently swallowed with `continue`. Permission or I/O errors cause files to be treated as "unchanged", silently skipped during update

### Performance Issue (NEW)
- **LOW**: `update.go:157` — `string(localContent) != string(cacheContent)` allocates two string copies for comparison. `bytes.Equal` does the same comparison with zero allocations

### Code Duplication (NEW)
- Three nearly identical path containment validation functions: `validateContainedPath` (extractor.go), `validateRemovePath` (remove.go), `validateSymlinkTarget` (downloader.go). Should be consolidated into a single exported `ValidateContainedPath` in core package

### Test Coverage Update
- Overall: `cmd/samuel` **0%**, `internal/commands` **22.6%**, `internal/core` **83.9%**, `internal/github` **89.4%**, `internal/ui` **17.4%**
- `internal/commands` improved from 19.4% → 22.6% (tasks 32, 34, 41 completed since last discovery)
- Files with 0% coverage still not covered by existing tasks: `add.go` (131 LOC), `version.go` (99 LOC), `remove.go` (partial — only validateRemovePath tested)
- 91 total silently discarded errors across 17 files (many are Cobra flag parsing, lower risk)
- 33 `fmt.Println` calls in commands/ bypassing ui package (consistency issue, not prioritized)

### Positive Findings
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- `go vet` clean — no issues
- No panic() calls, no TODO/FIXME/HACK markers
- No hardcoded credentials or secrets
- Security posture strong (path traversal, symlink, size limits, TOCTOU all fixed)
- `internal/core` at 83.9% exceeds 80% business logic coverage target
- `internal/github` at 89.4% well above target
- Linear commit history, clean branch structure
- All 20 completed tasks have passing quality checks

### Tasks Generated (Sixth Discovery): 10
| ID | Priority | Title |
|----|----------|-------|
| 51 | medium   | Fix silently discarded os.RemoveAll error in downloader.go dev cache clear |
| 52 | medium   | Fix silently discarded ReadFile errors in update.go file comparison |
| 53 | low      | Use bytes.Equal instead of string conversion for file comparison in update.go |
| 54 | medium   | Refactor runAdd() into smaller helper functions |
| 55 | medium   | Refactor runRemove() into smaller helper functions |
| 56 | medium   | Add unit tests for add.go command logic |
| 57 | medium   | Add unit tests for version.go runVersion function |
| 58 | low      | Refactor DownloadVersion() and extractTarGz() in downloader.go below 50-line limit |
| 59 | low      | Consolidate path containment validation into shared core utility |
| 60 | medium   | Add unit tests for remove.go validateRemovePath and removal logic |

[2026-02-22T04:30:00Z] [iteration:21] [task:7] COMPLETED: Refactored runInit() into smaller helper functions
- Split 277-line `runInit()` into a 28-line orchestrator that delegates to focused functions
- Created `init_steps.go` (284 lines) with step functions: parseInitFlags, validateInitTarget, selectComponents, selectTemplateInteractive, selectLanguagesInteractive, selectFrameworksInteractive, displayAndConfirm, downloadFramework, installAndSetup, updateSkillsAndAgentsMD
- Kept `init.go` (222 lines) with command registration, runInit orchestrator, and existing helpers (expandLanguages, expandFrameworks, isSamuelRepository, getRelevantFrameworks), plus reportInitResults and saveInitConfig
- Introduced `initFlags` and `initSelections` structs to reduce parameter passing between functions
- All functions under 50 lines, both files under 300 lines
- LEARNING: The `selectComponents` function was initially 60 lines even after extracting interactive selection helpers. Extracting `selectTemplateInteractive` (template options construction + prompt) brought it under 50. The key insight: the condition check + option construction + prompt + error handling pattern (13 lines) is a reusable extraction target.
- LEARNING: `ScanSkillsDirectory` returns `[]*core.SkillInfo` (pointer slice), not `[]core.SkillInfo`. When refactoring, type mismatches in return values are caught immediately by the compiler — always `go build` after extraction.
- LEARNING: Similar to the doctor.go refactor (task 6), splitting into two files (orchestrator + steps) keeps both under the 300-line limit. The logical grouping (init.go = command setup + helpers, init_steps.go = initialization steps) is a natural split.
- Commit: f9017e3

[2026-02-22T05:00:00Z] [iteration:22] [task:9] COMPLETED: Added unit tests for internal/core/auto_prompt.go
- Created auto_prompt_test.go with 12 test cases covering both exported functions
- Coverage: GetDefaultPromptTemplate 100%, GeneratePromptFile 100%
- Table-driven tests for GeneratePromptFile: 8 cases covering basic config, quality checks, pilot mode, combined features, empty/zero edge cases
- Separate tests for: template idempotency, template composition (starts with default template), quality checks ordering
- LEARNING: `auto_prompt.go` is a pure function module — no filesystem, no HTTP, no state. Both functions are deterministic string builders. The `GetDefaultPromptTemplate` function uses Go string concatenation with backtick-quoted strings to embed backtick characters (e.g., `` "`CLAUDE.md`" ``), which makes the template content itself contain markdown code formatting.
- LEARNING: `GeneratePromptFile` appends to the base template using `strings.Builder`. The pilot mode section and quality checks section are independently conditional — testing all 4 combinations (neither, only checks, only pilot, both) provides full branch coverage.
- Commit: 3bc2809

[2026-02-22T05:30:00Z] [iteration:23] [task:15] COMPLETED: Fixed silently discarded filepath.Walk and filepath.Glob errors in diff.go
- `filepath.Glob` error (line 190): now checked — logs warning via `ui.Warn` with the pattern that failed, then `continue` to skip to next pattern
- `filepath.Walk` for `.agent` directory (line 201): return value now checked — non-`os.IsNotExist` errors logged via `ui.Warn`
- `filepath.Walk` for template directory (line 223): return value now checked — non-`os.IsNotExist` errors logged via `ui.Warn`
- All three use `!os.IsNotExist(err)` guard because missing directories are expected cases (`.agent` dir may not exist, template path may not exist for certain versions)
- LEARNING: `filepath.Walk` with a callback that returns nil for all errors will itself always return nil. The Walk return value reflects only what the callback returns. So `_ = filepath.Walk(...)` is technically safe when the callback swallows errors, but handling the return value is still good practice for defense-in-depth (e.g., if the callback is later changed to return errors).
- LEARNING: `filepath.Glob` only returns errors for syntactically invalid patterns (per Go docs). Since the patterns here are hardcoded constants (`"CLAUDE.md"`, `"AGENTS.md"`, `".claude/**/*.md"`), the error path should never trigger in practice. The fix is still worthwhile for correctness and to avoid the `_ =` anti-pattern.
- Commit: cacd3f2

[2026-02-22T06:00:00Z] [iteration:24] [task:17] COMPLETED: Added unit tests for internal/ui/output.go
- Created output_test.go with 42 test cases covering all 14 exported functions + DisableColors
- Coverage: ui package improved from 17.4% to 49.6%
- `captureStdout` helper redirects both `os.Stdout` AND `color.Output` to a pipe — necessary because `color.Println()` methods write to `color.Output` (set at package init), not the current `os.Stdout`
- `captureStderr` helper for `Error()` which writes to `os.Stderr`
- Table-driven tests for: Success (3), Error (2), ListItem (4), SuccessItem (3), WarnItem (2), ErrorItem (2)
- Single assertion tests for: DisableColors, Warn, Info, Print, Bold, Dim, Header, Section, TableRow, ColoredTableRow
- All tests disable colors via `color.NoColor = true` + defer restore for deterministic output assertions
- LEARNING: `fatih/color` has two output paths: `Fprintf(writer, ...)` uses the explicit writer, but `Println(...)` uses `c.Output()` which defaults to `color.Output` (package-level var, set to `os.Stdout` at init). Replacing `os.Stdout` alone doesn't capture `Println` output — must also redirect `color.Output`.
- LEARNING: `Header` function uses both `fmt.Println()` (writes to `os.Stdout`) and `boldColor.Println(title)` (writes to `color.Output`). The two different output targets caused the initial test failure — the title line went to the original stdout while the blank lines went to the pipe.
- Commit: 73d44fa

---

[2026-02-22T06:30:00Z] [discovery] FOUND: Seventh discovery iteration — file close errors, Docker image validation, new file coverage gaps, code quality

### Data Integrity Issues (NEW)
- **HIGH**: `downloader.go:168,171` — `file.Close()` after `io.Copy` without error check. On some filesystems, `Close()` flushes buffered data; failed close = silently corrupted extracted files
- **HIGH**: `downloader.go:255` — `defer dstFile.Close()` in `copyDir` without error check after `io.Copy`. Same corruption risk for directory copies
- **HIGH**: `extractor.go:415` — `defer dst.Close()` in `copyFile` without error check after `io.Copy`. Same pattern, same risk

### Security Issues (NEW)
- **MEDIUM**: `auto_loop.go:146` — `SandboxImage` from `prd.json` passed directly to `docker run` without format validation. A tampered `prd.json` could specify a malicious Docker image name
- **MEDIUM**: `config_cmd.go:162-175` — `config set registry` accepts any string without URL scheme validation. A malicious registry URL could redirect downloads

### Test Coverage Update
- Overall: `cmd/samuel` **0%**, `internal/commands` **22.2%**, `internal/core` **85.6%**, `internal/github` **89.4%**, `internal/ui` **49.6%**
- `doctor_checks.go` (291 LOC) — created by task 6 refactor, **0% coverage** (new file, no tests)
- `init_steps.go` (284 LOC) — created by task 7 refactor, **0% coverage** (new file, no tests)
- `commands/skill.go` (375 LOC) — **0% coverage**, largest untested command file
- `commands/update.go` (264 LOC) — **0% coverage**, contains 219-line runUpdate function
- `commands/auto_start_handler.go` (139 LOC) — **0% coverage**
- 15 of 21 source files in `internal/commands/` have no corresponding test file

### Code Quality Update
- 17 functions exceed 50-line limit (top new violations not covered by prior tasks):
  - `runSkillInfo()` 92 lines (skill.go:284)
  - `runAutoStart()` 80 lines (auto_start_handler.go:12)
  - `listInstalled()` 73 lines (list.go:43)
  - `runVersion()` 71 lines (version.go:28)
  - `runSync()` 67 lines (sync.go:38)
  - `runSkillValidate()` 67 lines (skill.go:162)
- 7 non-test files exceed 300-line limit (registry.go 460, sync.go 431, skill.go 431, config.go 425, extractor.go 419, commands/skill.go 375, search.go 337)
- 8 deprecated `filepath.Walk` calls remain (6 in extractor.go, 2 in downloader.go)
- 7 magic numbers without named constants in skill.go, prompts.go, spinner.go

### Positive Findings
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- `go vet` clean — no issues
- `internal/core` at 85.6% — exceeds 80% business logic target
- `internal/github` at 89.4% — well above target
- No panic() calls, no TODO/FIXME/HACK markers, no hardcoded secrets
- Security posture strong: path traversal, symlink, size limits, TOCTOU, race condition all fixed
- HTTP client properly configured (30s timeout, size limits, no DefaultClient)
- All exec.Command calls validated against allowlist
- Cobra argument counts properly enforced across all commands

### Tasks Generated (Seventh Discovery): 10
| ID | Priority | Title |
|----|----------|-------|
| 61 | high     | Fix unchecked file.Close() errors after io.Copy in downloader.go and extractor.go |
| 62 | medium   | Add Docker sandbox image name validation in auto_loop.go |
| 63 | medium   | Add unit tests for doctor_checks.go check functions |
| 64 | medium   | Add unit tests for init_steps.go step functions |
| 65 | medium   | Add unit tests for commands/update.go helper logic |
| 66 | medium   | Refactor runSkillInfo() in commands/skill.go below 50-line limit |
| 67 | medium   | Refactor runAutoStart() in auto_start_handler.go below 50-line limit |
| 68 | medium   | Add URL scheme validation for registry config value in config_cmd.go |
| 69 | low      | Replace deprecated filepath.Walk with filepath.WalkDir in core package |
| 70 | low      | Extract magic numbers to named constants in skill.go and ui/prompts.go |

[2026-02-22T07:00:00Z] [iteration:25] [task:61] COMPLETED: Fixed unchecked file.Close() errors after io.Copy in downloader.go and extractor.go
- `extractTarGz`: Normal-path `file.Close()` now checked — returns `fmt.Errorf("failed to close file %q: %w", ...)`. Error-path close (line 168) left unchecked since the primary io.Copy error is more relevant.
- `copyFile` (downloader.go): Added named return `(err error)` and replaced `defer dstFile.Close()` with deferred closure that checks `cerr := dstFile.Close(); cerr != nil && err == nil` pattern.
- `copySingleFile` (extractor.go): Same named return + deferred error check pattern applied.
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: The `defer func() { if cerr := f.Close(); cerr != nil && err == nil { err = cerr } }()` pattern requires a named return value. The `err == nil` guard ensures the Close error doesn't mask a more important io.Copy error — if io.Copy already failed, the Close error is secondary.
- LEARNING: In `copySingleFile`, the `if err := os.MkdirAll(...)` block uses short variable declaration (`:=`) inside the `if` scope, so it does NOT shadow the named return `err`. After the `if` block exits, subsequent uses of `err` still refer to the named return. This is a subtle but important Go scoping distinction.
- Commit: a53222d

[2026-02-22T07:30:00Z] [iteration:26] [task:62] COMPLETED: Added Docker sandbox image name validation in auto_loop.go
- Added `IsValidSandboxImage()` function in `docker.go` using a regex pattern matching Docker image reference format: `[registry/]name[:tag][@digest]`
- Rejects: empty strings, absolute paths (`/bin/sh`), relative path escapes (`../`), shell metacharacters (`;`, `$()`, backticks, `|`, `&`, spaces)
- Accepts: standard images (`node:lts`), registry-prefixed (`ghcr.io/owner/image:latest`), digest-pinned (`node:lts@sha256:...`)
- Added validation call in `invokeAgentDocker` before image is passed to `buildDockerRunArgs`
- Added 20 test cases in `TestIsValidSandboxImage` (9 valid, 11 invalid) and 4 regression tests in `TestInvokeAgentDocker_RejectsInvalidImage`
- LEARNING: Docker image names always start with an alphanumeric character — names starting with `.` or `/` are filesystem paths, not image references. The regex `^[a-zA-Z0-9]` naturally rejects both absolute and relative path attacks.
- LEARNING: `exec.Command("docker", "run", image)` does NOT invoke a shell, so shell metacharacters in `image` won't cause command injection. However, a malicious image name like `evil-registry.com/backdoor:latest` is still dangerous because Docker will pull and run it. The validation ensures the format is valid but cannot prevent pulling from untrusted registries — that's a policy decision outside scope.
- Commit: 8f60188

[2026-02-22T08:00:00Z] [iteration:27] [task:63] COMPLETED: Added unit tests for doctor_checks.go check functions
- Created doctor_checks_test.go with 34 test cases covering all 10 functions + extractVersion
- Functions tested: checkCLAUDEMD (4), checkAGENTSMD (2), checkDirectoryStructure (3), checkInstalledSkills (4), checkInstalledComponents (2), checkSkillsIntegrity (5), checkAutoHealth (4), checkLocalModifications (2), checkModification (2), extractVersion (6 table-driven)
- Used t.TempDir() for all filesystem tests — crafted directory structures with SKILL.md, CLAUDE.md, AGENTS.md, prd.json as needed
- checkConfigFile not directly tested (reads from cwd via core.LoadConfig which is hard to redirect without changing cwd)
- checkAutoHealth tests create real prd.json files with valid/invalid JSON/structure to test all branches
- checkSkillsIntegrity tests create valid and invalid SKILL.md files (with/without YAML frontmatter) to exercise the skill validation pipeline
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: `checkInstalledSkills` with a nil-returning finder (unknown component) silently skips that component — it only checks SKILL.md existence when `finder(name) != nil`. This means unknown components in the config are not flagged as errors, which is intentional (forward compatibility with newer registry versions).
- LEARNING: `extractVersion` tries the bold pattern (`**Current Version**: X.Y.Z`) first, then falls back to plain. If both match, bold wins. This is tested explicitly in the `bold_preferred_over_plain` case.
- LEARNING: `checkAutoHealth` depends on `core.ValidateAutoPRD` which checks for empty Version, empty Project.Name, empty Project.Description, etc. A PRD with only `Version` missing triggers validation errors — useful for testing the validation-errors branch without constructing a fully invalid PRD.

[2026-02-22T08:30:00Z] [iteration:28] [task:64] COMPLETED: Added unit tests for init_steps.go step functions
- Created init_steps_test.go with 28 test cases covering 5 functions + 2 structs
- `parseInitFlags`: 9 tests (defaults, target args existing/nonexistent, force, non-interactive, template/languages/frameworks flags, dot target, cliProvided detection)
- `validateInitTarget`: 5 tests (valid dir, samuel repo rejected, config exists with/without force, alt config)
- `selectComponents` (non-interactive paths): 7 tests (template resolution, unknown template error, language/framework flag overrides, defaults to starter, minimal template, language-only flags without template)
- `updateSkillsAndAgentsMD`: 4 tests (CLAUDE.md → AGENTS.md copy, skills directory scanning, no CLAUDE.md, multiple skills returned)
- `installAndSetup`: 1 test (directory creation when createDir flag set)
- `initSelections` struct: 2 tests (zero value, populated)
- Used `newInitCmd()` helper to create fresh cobra.Command with correct flags for each test
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: `selectComponents` has many interactive paths (template selection, language selection, framework selection) that require UI prompts. The non-interactive + cliProvided code paths are fully testable without mocking. Setting `nonInteractive=true` or `cliProvided=true` bypasses all `ui.Select`/`ui.MultiSelect` calls.
- LEARNING: `parseInitFlags` uses cobra's `cmd.Flags().GetBool/GetString/GetStringSlice` which silently return zero values if the flag wasn't registered. The test helper `newInitCmd()` must register all 5 flags to match the real command setup, otherwise tests would pass with wrong zero values.
- LEARNING: `installAndSetup` can be partially tested — directory creation happens before the extractor runs, so even with a nonexistent cache path, we can verify the directory was created. The extractor errors are handled gracefully (reported via `ui.Error` calls).
- Commit: 7a9c814

---

[2026-02-22T09:00:00Z] [discovery] FOUND: Eighth discovery iteration — error handling, test coverage, code quality, performance

### Error Handling Issues (NEW)
- **HIGH**: `info.go:62` and `list.go:122` — `config, _ := core.LoadConfig()` silently discards errors (permission denied, corrupt YAML, I/O errors). Same pattern previously fixed in search.go (task 1). Non-ErrNotExist errors should be warned.
- **MEDIUM**: `downloader.go:49` — `os.RemoveAll(cacheDest)` error silently discarded when clearing dev branch cache. Failed removal + new download = corrupted mixed state.
- **MEDIUM**: `update.go:147-155` — Two `os.ReadFile` calls silently `continue` on error. Permission/I/O errors cause files to be treated as "unchanged" and skipped during update.
- **MEDIUM**: `sync.go:221` — `rel, _ := filepath.Rel(opts.RootDir, path)` discards error. Empty `rel` on failure causes incorrect path construction.

### Test Coverage Update
- Overall: `cmd/samuel` **0%**, `internal/commands` **22.6%**, `internal/core` **85.5%**, `internal/github` **89.4%**, `internal/ui` **49.6%**
- `internal/commands/` still has 13 of 23 source files with 0% test coverage
- Highest-value untested files: `skill.go` (375 LOC, 0%), `list.go` (164 LOC, 0%), `config_cmd.go` (216 LOC, minimal 62-line tests)
- `internal/core` at 85.5% exceeds the 80% business logic target
- `internal/github` at 89.4% well above target
- `internal/ui` improved from 0% to 49.6% across previous iterations

### Code Quality Issues (NEW/PERSISTING)
- `runUpdate()` at 219 lines (4.4x the 50-line function limit) — largest function in codebase, no existing refactor task
- `update.go:157` uses `string(localContent) != string(cacheContent)` — allocates 2 unnecessary string copies; `bytes.Equal` is zero-allocation
- 3 duplicate path containment validation functions: `validateContainedPath` (extractor.go), `validateRemovePath` (remove.go), `validateSymlinkTarget` (downloader.go) — should consolidate

### Positive Findings
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- `go vet` clean — no issues
- No TODO/FIXME/HACK markers, no hardcoded secrets, no panic() calls
- Security posture strong: path traversal, symlink, size limits, TOCTOU, race condition, Docker image validation all fixed
- HTTP client properly configured (30s timeout, size limits, no DefaultClient)
- All exec.Command calls validated against allowlist
- 28 iterations completed with 4 tasks finished in this batch (61-64)

### Tasks Generated (Eighth Discovery): 10
| ID | Priority | Title |
|----|----------|-------|
| 71 | high     | Fix silently discarded LoadConfig errors in info.go and list.go |
| 72 | high     | Add unit tests for commands/skill.go skill management functions |
| 73 | medium   | Fix silently discarded os.RemoveAll error in downloader.go dev cache clear |
| 74 | medium   | Fix silently discarded ReadFile errors in update.go file comparison |
| 75 | medium   | Fix silently discarded filepath.Rel error in sync.go SyncFolderCLAUDEMDs |
| 76 | medium   | Add unit tests for commands/list.go helper functions |
| 77 | medium   | Add unit tests for commands/config_cmd.go pure functions |
| 78 | medium   | Refactor runUpdate() in commands/update.go below 50-line limit |
| 79 | low      | Use bytes.Equal instead of string conversion in update.go |
| 80 | low      | Consolidate path containment validation into shared core utility |

[2026-02-22T09:30:00Z] [iteration:29] [task:71] COMPLETED: Fixed silently discarded LoadConfig errors in info.go and list.go
- `info.go:62`: Changed `config, _ := core.LoadConfig()` to check error — warns via `ui.Warn` for non-`os.ErrNotExist` errors
- `list.go:122` (`listAvailable`): Same pattern applied — warns on non-ErrNotExist errors while keeping config nil-safe
- Both match the pattern established in search.go (task 1): `if configErr != nil && !os.IsNotExist(configErr) { ui.Warn(...) }`
- Note: `listInstalled` (list.go:44) already properly handles LoadConfig errors with a full error return — only `listAvailable` needed fixing
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: The `config, _ := core.LoadConfig()` anti-pattern appeared in 4 command files total: search.go (fixed task 1), info.go (fixed this task), list.go (fixed this task), and the `os` import was already present in both files so no import changes were needed.
- LEARNING: `listInstalled` vs `listAvailable` have different error handling requirements: `listInstalled` needs a valid config (returns error on failure), while `listAvailable` uses config only for marking installed items (config is optional/best-effort). The warn-and-continue pattern is correct for `listAvailable`.
- Commit: 8550e2c

[2026-02-22T10:00:00Z] [iteration:31] [task:72] COMPLETED: Added unit tests for commands/skill.go skill management functions
- Created skill_test.go with 25 test cases covering all 4 command handlers
- `runSkillCreate`: 5 tests (invalid name, consecutive hyphens, no config, successful creation, duplicate skill)
- `runSkillValidate`: 7 tests (no skills dir, empty skills dir, specific not found, specific valid, specific invalid, all valid, mixed valid/invalid)
- `runSkillList`: 5 tests (no skills dir, empty skills dir, valid skills, skills with errors, long description truncation)
- `runSkillInfo`: 8 tests (not found, valid with metadata, optional directories, multiline description, custom metadata, validation errors, long body >500 lines warning, missing SKILL.md)
- Used `setupSkillTestDir` helper that creates t.TempDir() with samuel.yaml + .claude/skills/ and os.Chdir with defer restore
- Used `createSkillDir` helper to DRY up SKILL.md file creation in test directories
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: `runSkillInfo` does NOT return an error when the skill has validation errors — it displays the errors via `ui.ErrorItem` and returns nil. This contrasts with `runSkillValidate` which returns `fmt.Errorf` when any skill is invalid. The difference is intentional: info is display-only, validate is a check command.
- LEARNING: `ScanSkillsDirectory` only includes directories that contain a SKILL.md file. A directory without SKILL.md is silently skipped — it doesn't appear in the scan results and doesn't trigger errors. But `LoadSkillInfo` (used by info/validate with specific name) returns info with `Errors: ["missing required file: SKILL.md"]` for such directories.
- LEARNING: The `os.Chdir(t.TempDir()) + defer restore` pattern is necessary because all 4 command handlers use `os.Getwd()` to find the skills directory. There's no way to inject the working directory without this approach.
- Commit: 1875291

[2026-02-22T10:30:00Z] [iteration:32] [task:65] COMPLETED: Added unit tests for commands/update.go helper logic
- Extracted `categorizeFileChanges` from `runUpdate` into a standalone testable function with `fileChanges` struct
- Created update_test.go with 16 test cases covering 3 functions
- `fileExists`: 4 tests (existing file, nonexistent, directory, empty path)
- `categorizeFileChanges`: 10 tests (empty paths, new file, unchanged, modified, removed in new version, mixed categories, nested paths, unreadable local/cache files, empty files)
- `runUpdate`: 2 tests (no config → error, corrupt config → error)
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: `runUpdate` is heavily coupled to `core.NewDownloader()` which calls `EnsureCacheDir()` and `github.NewClient()`. The only testable early-exit paths that don't hit the downloader are: (1) no config file, (2) corrupt config file. The "already up to date" path still requires `NewDownloader()` to succeed first.
- LEARNING: Extracting `categorizeFileChanges` as a standalone function is the right approach for this type of orchestration function. The function is pure (depends only on filesystem state) and makes the core business logic testable without mocking HTTP clients or downloaders. This is the same pattern used in other tasks (e.g., task 6 extracting doctor checks).
- LEARNING: The `fileChanges` struct with named fields is clearer than returning 3 separate slices. It also makes it easy to add new categories (e.g., "deleted files") in the future without changing the function signature.
- Commit: 690365a

---

[2026-02-22T11:00:00Z] [discovery] FOUND: Ninth discovery iteration — test coverage, function size violations, refactoring opportunities

### Test Coverage Update
- Overall: `cmd/samuel` **0%**, `internal/commands` **38.9%** (up from 22.6%), `internal/core` **85.5%**, `internal/github` **89.4%**, `internal/ui` **49.6%**
- `internal/commands` improved significantly from recent task completions (63, 64, 65, 72)
- `internal/core` at 85.5% exceeds the 80% business logic target
- `internal/github` at 89.4% well above target

### Untested Command Files (0% Coverage)
- `auto_task_handlers.go` (127 LOC) — 7 functions: taskStatusIcon (pure), updateTaskStatus, runAutoTaskAdd, etc.
- `add.go` (130 LOC) — runAdd at 97 lines, component validation and config update logic
- `auto_pilot.go` (296 LOC) — parsePilotFlags and parseAutoFlags are pure functions, executePilotLoop is 100 lines
- `auto_pilot_output.go` (73 LOC) — printPilotDryRun and printPilotSummary display functions
- `auto_start_handler.go` (139 LOC) — runAutoStart at 80 lines (task 67 tracks refactor)
- `version.go` (98 LOC) — runVersion at 71 lines, version display and update checking
- `sync.go` (113 LOC) — runSync at 67 lines, relPath pure helper untested
- `auto.go` (184 LOC) — init() only, command registration (not testable)
- `doctor.go` (175 LOC) — orchestrator + display functions (checks tested in doctor_checks_test.go)

### Function Size Violations (NEW — not covered by existing tasks)
- `runAdd()` in add.go: **97 lines** (1.94x limit) — mixes validation, download, extraction, config update
- `executePilotLoop()` in auto_pilot.go: **100 lines** (2x limit) — interleaved discovery/implementation loop with state tracking
- `runVersion()` in version.go: **71 lines** (1.42x limit) — mixes version display and update checking
- `runSync()` in sync.go: **67 lines** (1.34x limit) — mixes config loading, sync execution, result display

### Code Quality Issues (PERSISTING)
- 14 pending tasks remain from previous discoveries (tasks 66-80)
- `internal/commands` at 38.9% is below the 60% overall coverage target
- 9 non-test files still exceed 300-line limit
- 8 deprecated `filepath.Walk` calls remain (task 69 pending)
- 3 duplicate path containment validation functions (task 80 pending)

### Positive Findings
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- `go vet` clean — no issues
- No panic() calls, no TODO/FIXME/HACK markers, no hardcoded secrets
- Security posture strong: all critical and high-priority security tasks completed
- `internal/core` at 85.5% exceeds 80% business logic coverage target
- 32 iterations completed, 6 tasks completed in current batch (61-65, 71-72)
- Established test patterns consistently applied (table-driven, t.TempDir, redirectTransport, stdout capture)

### Tasks Generated (Ninth Discovery): 10
| ID | Priority | Title |
|----|----------|-------|
| 81 | medium   | Add unit tests for auto_task_handlers.go task management functions |
| 82 | medium   | Refactor runAdd() in commands/add.go below 50-line limit |
| 83 | medium   | Add unit tests for commands/add.go component addition logic |
| 84 | medium   | Refactor executePilotLoop() in auto_pilot.go below 50-line limit |
| 85 | medium   | Add unit tests for auto_pilot.go flag parsing functions |
| 86 | low      | Refactor runVersion() in commands/version.go below 50-line limit |
| 87 | low      | Add unit tests for commands/version.go version display logic |
| 88 | low      | Add unit tests for commands/sync.go relPath helper and runSync error paths |
| 89 | low      | Add unit tests for auto_pilot_output.go display functions |
| 90 | low      | Refactor runSync() in commands/sync.go below 50-line limit |

[2026-02-22T12:00:00Z] [iteration:33] [task:66] COMPLETED: Refactored runSkillInfo() in commands/skill.go below 50-line limit
- Split 92-line `runSkillInfo()` into a 37-line orchestrator by extracting two display helpers
- Created `skill_info.go` (68 lines) with: `displaySkillMetadata(info)` (metadata section, ~28 lines), `displaySkillStructure(info)` (structure + stats section, ~25 lines)
- Kept validation display inline in `runSkillInfo()` (~10 lines) as it's short enough
- skill.go reduced from 376 to 320 lines (closer to 300-line limit)
- All existing tests pass — 8 `TestRunSkillInfo` subtests verify the display helpers work correctly through `runSkillInfo`
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: The split follows the same pattern as doctor.go (task 6) and init.go (task 7): orchestrator stays in the main file, extracted helpers go to a dedicated file. The naming convention `<command>_info.go` is descriptive for display-related helpers.
- LEARNING: Validation display (~10 lines) was left inline in `runSkillInfo()` because extracting a 10-line function with no logic beyond conditional formatting adds indirection without meaningful benefit. The 37-line result is well under the 50-line limit.
- Commit: 3a077ad

[2026-02-22T12:30:00Z] [iteration:34] [task:67] COMPLETED: Refactored runAutoStart() in auto_start_handler.go below 50-line limit
- Split 80-line `runAutoStart()` into a 50-line orchestrator by extracting two helpers
- `resolveSandboxFlags(cmd, prd)` (15 lines): extracts sandbox, image, template from CLI flags with PRD defaults
- `buildLoopConfig(cmd, cwd, prd, sandbox, image, template)` (23 lines): creates LoopConfig with CLI overrides, iteration limit, and UI callbacks
- File grew from 140 to 153 lines — well under the 300-line limit
- All existing tests pass, all quality checks pass
- LEARNING: `core.NewLoopConfig()` returns a value type `core.LoopConfig` (not a pointer), and `core.RunAutoLoop()` takes a value type. When extracting `buildLoopConfig`, the return type must be `core.LoopConfig` not `*core.LoopConfig` — the compiler catches this immediately.
- LEARNING: For small files (140 lines), keeping extracted helpers in the same file is cleaner than creating a separate file. The doctor.go/init.go split pattern (task 6/7) was necessary because those files were 400+ lines. Here, a single file at 153 lines is readable and well-organized.
- Commit: ee69280

[2026-02-22T13:00:00Z] [iteration:35] [task:68] COMPLETED: Added URL scheme validation for registry config value in config_cmd.go
- Added `validateRegistryURL(value)` function using `net/url.Parse` + scheme/host checks
- Rejects: HTTP URLs, FTP URLs, non-URL strings, empty strings, URLs without hosts (e.g., `https://`)
- Accepts: valid HTTPS URLs with paths, ports, etc.
- Validation runs in `runConfigSet` before config loading — fail-fast on invalid registry values
- Added 9 table-driven tests covering valid HTTPS URLs, HTTP rejection, empty scheme, FTP, empty string, plain text, scheme-only
- File stayed at 239 lines — well under 300-line limit
- LEARNING: Go's `url.Parse` is very permissive — `"not-a-url"` parses successfully with an empty scheme and the string as the path. Checking `u.Scheme != "https"` catches this case naturally since the scheme will be empty.
- LEARNING: `url.Parse("https://")` succeeds with `Host: ""` — the host check is necessary to reject scheme-only URLs that `url.Parse` considers valid.
- Commit: 57a1aa7

[2026-02-22T13:30:00Z] [iteration:36] [task:73] COMPLETED: Fixed silently discarded os.RemoveAll error in downloader.go dev cache clear
- Changed `os.RemoveAll(cacheDest)` to check and return error with context: `fmt.Errorf("failed to clear dev cache: %w", err)`
- A failed removal followed by a new download could produce corrupted mixed state (old + new files)
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: This is a one-line fix pattern that recurs across the codebase (tasks 1, 15, 23, 71 all fixed similar `_ =` error suppression). The consistent fix is: check the error, log/warn for non-critical paths, return error for critical paths. Here, cache clearing is critical (corrupted state risk), so returning the error is correct.
- Commit: 781b2ba

[2026-02-22T14:00:00Z] [iteration:37] [task:74] COMPLETED: Fixed silently discarded ReadFile errors in update.go file comparison
- Added `ui.Warn("Skipping %s: failed to read local file: %v", ...)` and `ui.Warn("Skipping %s: failed to read cached file: %v", ...)` in `categorizeFileChanges`
- Previously, `os.ReadFile` errors silently caused files to be skipped via `continue` — permission denied or I/O errors made files appear "unchanged"
- The fix was a previous iteration's incomplete work — the code changes were made but not committed
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- Existing tests already covered this path: `TestCategorizeFileChanges/unreadable_local_file_skipped` and `unreadable_cache_file_skipped` (from task 65) — tests output the warning messages as expected
- LEARNING: When a previous iteration marks a task as `in_progress` but doesn't commit, the next iteration should check `git diff` to see if the work was already done. In this case, the code changes were correct and just needed quality checks + commit.
- LEARNING: Task 65 already created the test cases for this exact scenario (unreadable files) when it extracted `categorizeFileChanges`, but the `ui.Warn` calls weren't in the code yet. The tests still passed because they only verified the file was skipped (not in any category), not that a warning was logged. The tests now additionally exercise the warning path.
- Commit: 9d4c70d

[2026-02-22T14:30:00Z] [iteration:38] [task:75] COMPLETED: Fixed silently discarded filepath.Rel error in sync.go SyncFolderCLAUDEMDs
- Changed `rel, _ := filepath.Rel(opts.RootDir, path)` to check the error
- On failure, appends to `result.Errors` with context and returns `filepath.SkipDir` to skip the directory
- Follows the existing `result.Errors` pattern used by `AnalyzeFolder` errors (line 230) rather than `ui.Warn`, since the `core` package does not import `ui`
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: The `core` package has no dependency on `ui` — error reporting in `core` uses `result.Errors` (append to a slice) rather than direct console output. This is the correct pattern for a library package. The calling code (in `commands/sync.go`) can decide how to present errors to the user.
- LEARNING: `filepath.Rel` only fails on Unix when the two paths are on different volumes (which doesn't happen on Unix — volumes are a Windows concept). So this error path is primarily for Windows cross-volume scenarios. The fix is still correct for defense-in-depth and avoids the `_ =` anti-pattern.
- Commit: cd347a9

[2026-02-22T15:00:00Z] [iteration:39] [task:76] COMPLETED: Added unit tests for commands/list.go helper functions
- Created list_test.go with 22 test cases covering 3 functions
- `listInstalled`: 12 tests (no config, corrupt config, empty config, languages only, frameworks only, workflows all, specific workflows, unknown component names, type filter languages/frameworks/workflows, mixed components)
- `listAvailable`: 7 tests (no config, with config marks installed, all workflows installed, type filter languages/frameworks/workflows, corrupt config warns and continues)
- `runList`: 3 tests (default calls listInstalled, --available flag calls listAvailable, type filter passed through)
- Used `setupListTestDir` helper with `core.Config.Save(dir)` + `os.Chdir` pattern for config setup
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: `listInstalled` returns a proper error on corrupt config (non-ErrNotExist), while `listAvailable` warns and continues with nil config. This asymmetry is intentional: `listInstalled` requires a valid config to show installed items, while `listAvailable` uses config only as an optional enhancement to mark installed items.
- LEARNING: The `listInstalled` workflows section has special handling: when `config.Installed.Workflows` is `["all"]`, the count shows the full `len(core.Workflows)` and iterates over all workflow entries from the registry instead of the config list. This is the same "all" expansion pattern used in doctor_checks.go.
- LEARNING: Unknown component names (not in registry) are displayed without descriptions — `FindLanguage/FindFramework/FindWorkflow` returns nil, and the fallback path just prints the name. This is forward-compatible with newer registry versions.
- Commit: 6acf65c

[2026-02-22T15:30:00Z] [iteration:40] [task:77] COMPLETED: Added unit tests for commands/config_cmd.go pure functions and handlers
- Expanded config_cmd_test.go from 98 lines (3 test functions) to 270 lines (7 test functions)
- `isValidConfigKey`: expanded from 9 to 16 table-driven cases — added all 5 missing auto.* keys (auto.enabled, auto.ai_tool, auto.max_iterations, auto.quality_checks) and installed.skills, plus edge cases (partial prefixes "installed.", "auto.", leading/trailing spaces)
- `isValidConfigKey_MatchesValidConfigKeys`: new meta-test ensuring every key in `core.ValidConfigKeys` passes validation — prevents desync between the two
- `formatConfigValue`: expanded from 7 to 12 cases — added bool true/false, float, URL string, string with spaces
- `validateRegistryURL`: expanded from 9 to 11 cases — added bare host, SSH scheme rejection
- `runConfigList`: 3 tests (no config graceful, valid config display, corrupt config error)
- `runConfigGet`: 4 tests (invalid key error, no config graceful, valid key, installed.languages list)
- `runConfigSet`: 6 tests (invalid key, invalid registry URL, no config graceful, set version + verify reload, set registry + verify reload, set installed.languages + verify reload)
- Created `setupConfigTestDir` helper with `t.TempDir()` + `os.Chdir` + cleanup for config-dependent tests
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: The existing tests only covered 5 of 10 valid config keys. The `auto.*` and `installed.skills` keys were added to `ValidConfigKeys` after the initial tests were written. The `MatchesValidConfigKeys` meta-test prevents this drift in the future.
- LEARNING: `runConfigSet` with `installed.languages` uses `config.SetValue` which splits comma-separated strings into a slice internally. The test verifies the round-trip: set "go,rust,python" → reload → verify 3 items in the Languages slice.
- Commit: a9d0b2a

---

[2026-02-22T16:00:00Z] [discovery] FOUND: Tenth discovery iteration — error handling, test coverage, security, code quality

### Test Coverage Update
- Overall: `cmd/samuel` **0%**, `internal/commands` **45.4%** (up from 38.9%), `internal/core` **85.3%**, `internal/github` **89.4%**, `internal/ui` **49.6%**
- `internal/commands` improved from 38.9% → 45.4% — still below 60% target
- Overall coverage: **61.2%** — marginally above 60% minimum
- `go vet ./...` clean — zero warnings
- 11 source files in `internal/commands/` still have no corresponding test file

### Error Handling Issues (NEW — not covered by existing tasks)
- **HIGH**: `diff.go:196,210,239` — 3 instances of `relPath, _ := filepath.Rel(...)` silently discarded. Empty relPath on failure causes incorrect hash map keys and wrong diff results. Was task 37 (Fourth Discovery) but never completed.
- **HIGH**: `extractor.go:87,104` — 2 instances of `relPath, _ := filepath.Rel(e.destPath, dstPath)` silently discarded. Empty relPath causes incorrect backup/restore path construction. Was task 25 (Third Discovery) but never completed.
- **MEDIUM**: `auto_progress.go:51` — `defer f.Close()` on write handle opened with `os.O_APPEND|os.O_WRONLY`. Failed Close() means silently lost progress data.

### Security Concern (NEW)
- **MEDIUM**: `skill.go:148` — `filepath.Join(skillsDir, name)` uses CLI-provided name. While `ValidateSkillName` runs before this at line 138, adding a `validateContainedPath` check provides defense-in-depth against any future validation bypass.

### Test Coverage Gaps (NEW — not covered by existing tasks)
- `registry.go` has 6 skill lookup functions at 0%: FindSkill, GetAllSkillNames, GetLanguageSkills, SkillToLanguageName, GetFrameworkSkills, GetWorkflowSkills. These are pure functions operating on static data — trivial to test with table-driven tests.
- `diff.go` (305 LOC) has 0% coverage. computeDiff, collectHashes, hashFile are testable with t.TempDir() setups.
- `doctor.go` (175 LOC) has 0% coverage on orchestrator side. printCheckResults, printCheckSummary, performAutoFix are testable with crafted CheckResult slices.
- `auto_start_handler.go` (152 LOC) has 0% coverage. resolveSandboxFlags and buildLoopConfig (extracted by task 67) are testable pure functions.

### Code Quality Issues (NEW — not covered by existing tasks)
- `runRemove()` in remove.go: **101 lines** (2x the 50-line limit) — mixes component resolution, file deletion, config update. No existing refactor task.
- Search scoring magic numbers in search.go: 100, 80, 60, 40, 30, 5 used without named constants. Not covered by task 70 (which covers skill.go/prompts.go/spinner.go only).

### Positive Findings
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- `go vet` clean — no issues
- No TODO/FIXME/HACK markers in production code
- No hardcoded credentials or secrets
- Security posture strong: path traversal, symlink, size limits, TOCTOU, race condition, Docker image validation all fixed
- `internal/core` at 85.3% exceeds 80% business logic coverage target
- `internal/github` at 89.4% well above target
- 40 iterations completed, 14 tasks completed in this batch (61-77)
- 16 functions still exceed 50-line limit, 8 files still exceed 300-line limit

### Tasks Generated (Tenth Discovery): 10
| ID  | Priority | Title |
|-----|----------|-------|
| 91  | high     | Fix silently discarded filepath.Rel errors in diff.go |
| 92  | high     | Fix silently discarded filepath.Rel errors in extractor.go |
| 93  | high     | Add unit tests for registry.go skill lookup functions |
| 94  | medium   | Fix unchecked Close() error on write handle in auto_progress.go |
| 95  | medium   | Add skill name path traversal validation in skill.go |
| 96  | medium   | Add unit tests for diff.go comparison functions |
| 97  | medium   | Add unit tests for doctor.go display and orchestrator functions |
| 98  | medium   | Refactor runRemove() in commands/remove.go below 50-line limit |
| 99  | medium   | Add unit tests for auto_start_handler.go pure functions |
| 100 | low      | Extract search scoring magic numbers to named constants in search.go |

[2026-02-22T16:30:00Z] [iteration:41] [task:91] COMPLETED: Fixed silently discarded filepath.Rel errors in diff.go
- Fixed 3 instances of `relPath, _ := filepath.Rel(...)` in `getLocalFileHashes` (2 locations) and `getVersionFileHashes` (1 location)
- Glob loop (line 196): checks error, logs via `ui.Warn`, `continue` to skip the match
- Walk callback for .agent dir (line 210): checks error, logs via `ui.Warn`, `return nil` to skip the entry
- Walk callback for template dir (line 239): checks error, logs via `ui.Warn`, `return nil` to skip the entry
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: This task was previously identified in both the Fourth Discovery (task 37) and Tenth Discovery (task 91). The earlier task 37 fixed only the `filepath.Walk` return values and `filepath.Glob` errors (committed as cacd3f2 in iteration 23), but the `filepath.Rel` errors within the callbacks remained. Task 91 completes the job by fixing the 3 `filepath.Rel` calls that were silently discarding errors inside the loop/callback bodies.
- LEARNING: In Walk callbacks, using `relErr` (not `err`) avoids shadowing the outer `err` parameter from the callback signature `func(path string, info os.FileInfo, err error)`.
- Commit: 50a1b18

[2026-02-22T17:00:00Z] [iteration:42] [task:92] COMPLETED: Fixed silently discarded filepath.Rel errors in extractor.go
- Fixed 3 instances of `relPath, _ := filepath.Rel(e.destPath, dstPath)` in extractFile (2 locations) and extractDir (1 location)
- extractFile skip path (line 87): checks error, returns `fmt.Errorf` with context
- extractFile success path (line 104): checks error, returns `fmt.Errorf` with context
- extractDir (line 117): checks error, returns `fmt.Errorf` with context
- Core package uses `return fmt.Errorf(...)` (not `ui.Warn`) since core doesn't depend on ui
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: The task description mentioned 2 instances at lines 87 and 104, but there was a 3rd instance at line 117 in `extractDir`. Always grep for the full pattern when fixing a class of issues.
- LEARNING: In the core package, error propagation via `return fmt.Errorf(...)` is preferred over `ui.Warn` because core is a dependency of commands, not the other way around. The UI layer should handle display; core should return errors.
- Commit: cb5acdc

[2026-02-22T16:30:00Z] [iteration:43] [task:93] COMPLETED: Added unit tests for registry.go skill lookup functions
- Added 12 new test functions covering FindSkill, GetAllSkillNames, GetLanguageSkills, GetFrameworkSkills, GetWorkflowSkills, SkillToLanguageName, LanguageToSkillName, FrameworkToSkillName, WorkflowToSkillName, round-trip conversion, category coverage, and Skills data integrity
- 287 lines of new test code added to registry_test.go
- All pure functions — no mocking needed, table-driven tests throughout
- LEARNING: SkillToLanguageName uses `len > 6` (not `>=`), meaning it requires at least 1 character before the "-guide" suffix. Input "-guide" (exactly 6 chars) is returned unchanged — this is intentional design.
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- Commit: 093a3df

[2026-02-22T17:30:00Z] [iteration:44] [task:78] COMPLETED: Refactored runUpdate() in commands/update.go below 50-line limit
- Split 184-line `runUpdate()` into a 44-line orchestrator that delegates to 5 focused helpers
- `downloadTargetVersion` (48 lines): downloader init, version resolution, version display, check-only and up-to-date short-circuits, download
- `displayChangeDiff` (27 lines): formats and prints file change summary for `--diff` flag
- `applyUpdate` (36 lines): orchestrates backup, extraction, reporting, and config save
- `backupModifiedFiles` (20 lines): creates timestamped backup directory and copies modified files
- `reportUpdateResults` (22 lines): displays new/preserved file counts and instructions for handling modifications
- `categorizeFileChanges` (37 lines) and `fileChanges` struct were already extracted by task 65 — retained as-is
- File went from 278 to exactly 300 lines — at the limit due to added function signatures and godocs
- All existing tests pass (16 tests across TestFileExists, TestCategorizeFileChanges, TestRunUpdate)
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- LEARNING: `downloadTargetVersion` returns `("", targetVersion, nil)` for two distinct "no-op" cases: (1) already up-to-date and (2) check-only mode. The caller checks `cachePath == ""` to detect both, since both mean "don't proceed with extraction." The empty-string sentinel is a clean Go idiom for optional return values.
- LEARNING: The original function had `var backupDir string` at the top of a 60-line block. Extracting `backupModifiedFiles` returns both `(backupDir, error)`, making the control flow clearer: backup is only attempted when there are modified files AND force is off.
- LEARNING: File size is the binding constraint, not function size. The original 278-line file grew to 304 after adding function signatures and godocs for 5 new functions. Trimming multi-line godocs to single lines (e.g., "backs up modified files, extracts updated files, reports results, and saves the updated config version" → "backs up modified files, extracts updates, and saves the config") saved 4 lines to hit exactly 300.
- Commit: a31c0fd

---

[2026-02-22T18:00:00Z] [discovery] FOUND: Eleventh discovery iteration — input validation, function size violations, test coverage, magic numbers

### Input Validation Issues (NEW)
- **HIGH**: `search.go:75` — `sortAndLimitResults(results, limit)` with negative `limit` from `--limit` flag causes runtime panic: `results[:limit]` with limit=-1 → "slice bounds out of range [:-1]". Zero limit silently returns no results. No validation on the flag value.
- **MEDIUM**: `auto_loop.go:30-42` — `PAUSE_SECONDS` and `MAX_CONSECUTIVE_FAILURES` env vars parsed via `strconv.Atoi` without range validation. Negative `PAUSE_SECONDS` → zero sleep (rapid-fire iterations). Zero `MAX_CONSECUTIVE_FAILURES` → loop exits on first failure (no fault tolerance).

### Function Size Violations (NEW — not covered by existing tasks)
- `extractTarGz()` in downloader.go: **70 lines** (1.4x limit) — security-critical tar extraction code
- `DownloadVersion()` in downloader.go: **66 lines** (1.32x limit) — cache + download + extraction orchestration
- `MultiSelect()` in prompts.go: **75 lines** (1.5x limit) — interactive selection loop with display construction
- `listInstalled()` in list.go: **72 lines** (1.44x limit) — three near-identical component display blocks
- `runSkillValidate()` in skill.go: **66 lines** (1.32x limit) — mixes loading, validation, and display

### Test Coverage Update
- Overall: `cmd/samuel` **0%**, `internal/commands` **~45%**, `internal/core` **~85%**, `internal/github` **89.4%**, `internal/ui` **~50%**
- `auto_handlers.go` (280 LOC) — only detectQualityChecks/countTaskStatuses/validateSandbox tested; initAutoDir/writeAutoFiles/printInitSummary untested
- 22 pending tasks remain (tasks 69-100), 18 completed in current batch
- 12 source files in `internal/commands/` still have no corresponding test file

### Magic Numbers (NEW — not covered by task 70 or 100)
- `auto_loop.go:30` — `pauseSecs := 2` and `maxConsecFails := 3` (loop defaults)
- `auto_loop.go:31,38` — `"PAUSE_SECONDS"` and `"MAX_CONSECUTIVE_FAILURES"` env var names inline
- `github/client.go:51` — `Timeout: 30 * time.Second` (HTTP client timeout)

### Positive Findings
- All quality checks pass: `go test ./...`, `go vet ./...`, `go build ./...`
- `go vet` clean — zero warnings
- No TODO/FIXME/HACK markers, no hardcoded secrets, no panic() calls
- Security posture strong: path traversal, symlink, size limits, TOCTOU, race condition, Docker image validation all fixed
- `internal/core` at ~85% exceeds 80% business logic coverage target
- `internal/github` at 89.4% well above target
- 44 iterations completed, 18 tasks completed in current batch
- Consistent test patterns maintained across all new tests

### Tasks Generated (Eleventh Discovery): 10
| ID  | Priority | Title |
|-----|----------|-------|
| 101 | high     | Fix potential panic from negative search limit flag in search.go |
| 102 | medium   | Add environment variable range validation in auto_loop.go |
| 103 | medium   | Add unit tests for auto_handlers.go init helper functions |
| 104 | medium   | Refactor extractTarGz() in downloader.go below 50-line limit |
| 105 | medium   | Refactor DownloadVersion() in downloader.go below 50-line limit |
| 106 | low      | Refactor MultiSelect() in prompts.go below 50-line limit |
| 107 | low      | Refactor listInstalled() in list.go below 50-line limit |
| 108 | low      | Refactor runSkillValidate() in skill.go below 50-line limit |
| 109 | low      | Extract magic numbers in auto_loop.go to named constants |
| 110 | low      | Extract HTTP client timeout to named constant in github/client.go |

[2026-02-21T23:00:00Z] [iteration:45] [task:101] COMPLETED: Fixed potential panic from negative search limit flag in search.go
- Added `defaultSearchLimit` constant (value: 20) replacing magic number in flag definition
- Added validation in `sortAndLimitResults()`: if `limit <= 0`, clamp to `defaultSearchLimit`
- Added 2 test cases to `TestSortAndLimitResults`: negative limit (-1) and zero limit (0) both fall back to default
- Added `TestDefaultSearchLimit` to verify constant value
- LEARNING: Go's `results[:limit]` with negative limit panics with "slice bounds out of range". The condition `len(results) > limit` is always true when limit < 0, so the panic path was guaranteed for any non-empty result set.
- Commit: 1207ffa

[2026-02-22T16:00:00Z] [iteration:46] [task:81] COMPLETED: Added unit tests for auto_task_handlers.go task management functions
- Created auto_task_handlers_test.go with 17 tests covering all exported and internal functions
- `taskStatusIcon`: table-driven tests for all 5 status constants + unknown + empty string
- `updateTaskStatus`: tested complete/skip/reset flows, task-not-found, missing PRD, corrupt JSON, and progress recalculation on save
- `runAutoTaskAdd`: success path, duplicate ID rejection, missing PRD, and JSON roundtrip verification
- `runAutoTaskList`: success with mixed statuses (including subtasks with ParentID), missing PRD
- `runAutoTaskComplete/Skip/Reset`: end-to-end through cobra handler wrappers
- LEARNING: Functions using `os.Getwd()` for PRD path resolution require `os.Chdir()` to temp dir in tests, with cleanup via `t.Cleanup(func() { os.Chdir(origDir) })`. The `setupTestPRD` helper creates `.claude/auto/prd.json` in a temp dir for reuse across tests.
- Commit: a4d9fe1
