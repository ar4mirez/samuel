# Release Checklist

This checklist ensures all documentation and version references are updated before creating a new release.

## Pre-Release Checklist

### 1. Version Updates

- [ ] **CHANGELOG.md** - Add release notes under `## [x.y.z] - YYYY-MM-DD`
- [ ] **README.md** - Update version badge (line ~6)
- [ ] **template/CLAUDE.md** - Update version in "Version & Changelog" section
- [ ] **docs/reference/changelog.md** - Mirror root CHANGELOG updates

### 2. Code Quality

- [ ] All tests pass: `make test`
- [ ] Linter passes: `make lint` (if available) or `golangci-lint run`
- [ ] Build succeeds: `make build`
- [ ] GoReleaser dry-run passes: `make release-dry`

### 3. Documentation

- [ ] README reflects current features
- [ ] Installation instructions are accurate
- [ ] CLI command documentation is up-to-date

### 4. Git Preparation

- [ ] All changes committed
- [ ] Branch is up-to-date with main
- [ ] No uncommitted changes: `git status`

## Creating the Release

### Option A: Automated (Recommended)

1. Create and push a version tag:
   ```bash
   git tag -a v1.x.y -m "Release v1.x.y"
   git push origin v1.x.y
   ```

2. The release workflow will automatically:
   - Run tests
   - Build binaries for all platforms
   - Create GitHub release with assets
   - Update Homebrew formula (if configured)
   - Trigger documentation deployment

### Option B: Manual

1. Run GoReleaser locally:
   ```bash
   export GITHUB_TOKEN=your_token
   make release
   ```

## Post-Release Verification

- [ ] GitHub Release page shows all binaries
- [ ] Checksums file is included
- [ ] Release notes are formatted correctly
- [ ] Documentation site updated (https://ar4mirez.github.io/aicof/)
- [ ] Homebrew formula updated (if applicable)
- [ ] Install script works: `curl -sSL https://raw.githubusercontent.com/ar4mirez/aicof/main/install.sh | sh`

## Versioning Guidelines

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (x.0.0): Breaking changes to CLI commands or config format
- **MINOR** (0.x.0): New features, new guides, new workflows
- **PATCH** (0.0.x): Bug fixes, documentation updates

### Pre-release Versions

For testing releases before official release:
- Alpha: `v1.0.0-alpha.1`
- Beta: `v1.0.0-beta.1`
- Release candidate: `v1.0.0-rc.1`

## Hotfix Process

For urgent fixes to a released version:

1. Create hotfix branch from tag: `git checkout -b hotfix/v1.x.y v1.x.y`
2. Apply fix and test
3. Update CHANGELOG with patch version
4. Tag and release: `git tag -a v1.x.z -m "Hotfix release"`
5. Merge back to main if applicable

## Required Secrets

The release workflow requires these GitHub secrets:

| Secret | Purpose | Required |
|--------|---------|----------|
| `GITHUB_TOKEN` | Create releases, upload assets | Auto-provided |
| `HOMEBREW_TAP_GITHUB_TOKEN` | Update Homebrew formula | Optional |

## Troubleshooting

### GoReleaser fails

1. Check configuration: `goreleaser check`
2. Run dry-run: `goreleaser release --snapshot --clean`
3. Check ldflags paths match package structure

### Homebrew formula not updated

1. Verify `HOMEBREW_TAP_GITHUB_TOKEN` is set
2. Check homebrew-tap repository exists
3. Review GoReleaser logs for Homebrew step

### Documentation not deployed

1. Check docs workflow triggered
2. Verify mkdocs build succeeds locally
3. Check GitHub Pages settings
