# PRD: GitHub Release Integration & Documentation Automation

> **ID**: 0003
> **Status**: Implemented
> **Created**: 2026-01-14
> **Author**: Claude Code (AI)

---

## Introduction

This feature implements complete GitHub release automation for the AICoF CLI, enabling automated binary builds, distribution, and documentation updates when version tags are pushed.

**Goal**: Enable seamless release workflow where pushing a version tag automatically builds, tests, and releases the CLI with proper documentation.

---

## Goals

1. Automate release process on version tag push (`v*`)
2. Run pre-release validation (tests, lint) before building
3. Build cross-platform binaries (Linux, macOS, Windows)
4. Create GitHub Release with downloadable assets
5. Generate and publish Homebrew formula
6. Provide curl-based installer for quick installation
7. Ensure documentation is always updated with releases
8. Establish clear release checklist for maintainers

---

## User Stories

**US-001**: As a maintainer, I want releases to be automated so that I only need to push a tag to create a full release.

**US-002**: As a user, I want to install via `curl | sh` so that I can quickly get started without managing downloads.

**US-003**: As a macOS user, I want to install via Homebrew so that I can manage updates alongside other tools.

**US-004**: As a maintainer, I want CI to run on PRs so that code quality is validated before merging.

---

## Functional Requirements

### CI Pipeline (FR-001 to FR-004)

- FR-001: Run tests on every push to main and PRs
- FR-002: Run linter (golangci-lint) on every push and PR
- FR-003: Build for multiple platforms to catch cross-platform issues
- FR-004: Only run CI for Go-related file changes

### Release Pipeline (FR-005 to FR-010)

- FR-005: Trigger on version tags matching `v*` pattern
- FR-006: Run full test suite before building
- FR-007: Build binaries for 5 platform combinations
- FR-008: Create GitHub Release with all artifacts
- FR-009: Generate checksums.txt for verification
- FR-010: Update Homebrew formula in homebrew-tap repository

### Install Script (FR-011 to FR-015)

- FR-011: Detect OS (Linux, macOS, Windows)
- FR-012: Detect architecture (amd64, arm64)
- FR-013: Download correct binary from latest release
- FR-014: Verify checksum before installation
- FR-015: Install to /usr/local/bin (or custom location)

---

## Non-Goals

- ❌ GPG signing of binaries (deferred to future)
- ❌ SBOM generation (deferred to future)
- ❌ Scoop package manager support (Windows)
- ❌ APT/YUM package repositories
- ❌ Automated version bumping
- ❌ Release notes generation from commits (manual CHANGELOG)

---

## Technical Considerations

### Existing Infrastructure

- GoReleaser already configured (needs ldflags path fix)
- Documentation workflow exists (triggers on release)
- Makefile has release targets

### Files Created/Modified

| File | Action | Purpose |
|------|--------|---------|
| `.goreleaser.yaml` | Modified | Fixed ldflags paths |
| `.github/workflows/ci.yml` | Created | CI pipeline |
| `.github/workflows/release.yml` | Created | Release automation |
| `install.sh` | Created | Curl installer |
| `.github/RELEASE_CHECKLIST.md` | Created | Maintainer guide |

### External Dependencies

1. **Homebrew tap repository**: `ar4mirez/homebrew-tap` must be created
2. **GitHub secret**: `HOMEBREW_TAP_GITHUB_TOKEN` for formula updates

---

## Guardrails Affected

### Security

- ✓ Checksum verification in install script
- ✓ No secrets exposed in logs
- ✓ HTTPS-only downloads

### Testing

- ✓ Tests run before release build
- ✓ CI validates all PRs

### Code Quality

- ✓ Linter runs on all PRs
- ✓ Multi-platform build validation

---

## Success Metrics

### Technical

- Release workflow completes in <10 minutes
- All 5 platform binaries built successfully
- Install script works on Linux, macOS, Windows (WSL)
- Zero failed releases due to automation issues

### User Experience

- Time from tag push to release: <10 minutes
- Single command installation works
- Documentation auto-updates on release

---

## Implementation Summary

### Phase 1: Fix Existing Issues ✅

- Fixed `.goreleaser.yaml` ldflags paths from `internal/cmd` to `internal/commands`

### Phase 2: CI Workflow ✅

- Created `.github/workflows/ci.yml`
- Runs tests, lint, and multi-platform build verification
- Triggers on pushes and PRs to main

### Phase 3: Release Workflow ✅

- Created `.github/workflows/release.yml`
- Triggers on `v*` tags
- Runs tests before GoReleaser
- Includes Homebrew formula generation

### Phase 4: Install Script ✅

- Created `install.sh`
- Cross-platform detection (OS + arch)
- Checksum verification
- Colorful output with ASCII art

### Phase 5: Documentation ✅

- Created `.github/RELEASE_CHECKLIST.md`
- Comprehensive pre/post release checklist
- Versioning guidelines
- Troubleshooting guide

---

## Open Items

### Required External Setup

1. [ ] Create `ar4mirez/homebrew-tap` repository on GitHub
2. [ ] Create PAT with `repo` scope for Homebrew tap
3. [ ] Add `HOMEBREW_TAP_GITHUB_TOKEN` secret to aicof repository

### Future Enhancements

- GPG signing of binaries
- SBOM generation
- Scoop package for Windows
- Automated changelog generation

---

## References

- GoReleaser documentation: https://goreleaser.com
- GitHub Actions: https://docs.github.com/en/actions
- Homebrew: https://docs.brew.sh
