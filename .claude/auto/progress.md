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
