# Dependency Update Workflow

> **Purpose**: Safe and systematic dependency updates with vulnerability management, license checking, and rollback planning.

---

## When to Use

| Trigger | Priority | Description |
|---------|----------|-------------|
| **Security Vulnerability** | Critical | Known CVE in dependency |
| **Monthly Maintenance** | High | Regular update cycle |
| **Major Version** | Medium | New major version available |
| **Pre-Release** | High | Before production deployments |
| **Breaking Bug** | Critical | Bug in current dependency |

---

## Update Strategy

### Update Types

| Type | Risk | Frequency | Testing |
|------|------|-----------|---------|
| **Patch** (x.x.1) | Low | Weekly/Auto | Basic |
| **Minor** (x.1.0) | Low-Medium | Monthly | Standard |
| **Major** (1.0.0) | High | Quarterly | Comprehensive |

### Semantic Versioning

```
MAJOR.MINOR.PATCH
  │     │     │
  │     │     └── Bug fixes (backward compatible)
  │     └──────── New features (backward compatible)
  └────────────── Breaking changes
```

---

## Prerequisites

Before starting:

- [ ] All tests passing
- [ ] Clean git working directory
- [ ] Recent backup/checkpoint
- [ ] Time for testing and potential rollback
- [ ] Access to changelogs/release notes

---

## Update Process

```
Phase 1: Audit Dependencies
    ↓
Phase 2: Check Vulnerabilities
    ↓
Phase 3: Check License Compatibility
    ↓
Phase 4: Plan Updates
    ↓
Phase 5: Execute Updates
    ↓
Phase 6: Test & Validate
    ↓
Phase 7: Document & Deploy
```

---

## Phase 1: Audit Dependencies

### 1.1 List Outdated Dependencies

**Node.js**:
```bash
npm outdated
# or
yarn outdated
# or
pnpm outdated
```

**Python**:
```bash
pip list --outdated
# or with poetry
poetry show --outdated
```

**Go**:
```bash
go list -u -m all
```

**Rust**:
```bash
cargo outdated
```

**Ruby**:
```bash
bundle outdated
```

### 1.2 Analyze Output

Create update inventory:

| Package | Current | Latest | Type | Risk |
|---------|---------|--------|------|------|
| react | 18.2.0 | 18.3.0 | Minor | Low |
| next | 14.0.0 | 15.0.0 | Major | High |
| lodash | 4.17.20 | 4.17.21 | Patch | Low |
| typescript | 5.0.0 | 5.4.0 | Minor | Medium |

### 1.3 Direct vs Transitive

Identify dependency type:

- **Direct**: Listed in package.json/requirements.txt
- **Transitive**: Dependencies of dependencies

Focus on direct dependencies first.

---

## Phase 2: Check Vulnerabilities

### 2.1 Run Security Audit

**Node.js**:
```bash
npm audit
npm audit --json > audit-report.json
```

**Python**:
```bash
pip-audit
# or
safety check
```

**Go**:
```bash
govulncheck ./...
```

**Rust**:
```bash
cargo audit
```

**Ruby**:
```bash
bundle audit check
```

### 2.2 Analyze Vulnerabilities

| Severity | Action | Timeline |
|----------|--------|----------|
| Critical | Immediate fix | Hours |
| High | Fix this sprint | Days |
| Moderate | Schedule fix | Weeks |
| Low | Track for next cycle | Month |

### 2.3 Vulnerability Report

```markdown
## Vulnerability Summary

### Critical (0)
None

### High (2)
1. **axios@0.21.0** - SSRF vulnerability
   - CVE: CVE-2021-3749
   - Fix: Upgrade to 0.21.1+
   - Impact: Server-side request forgery

2. **lodash@4.17.15** - Prototype pollution
   - CVE: CVE-2021-23337
   - Fix: Upgrade to 4.17.21
   - Impact: Remote code execution

### Moderate (3)
[Details...]
```

---

## Phase 3: Check License Compatibility

### 3.1 Audit Licenses

**Node.js**:
```bash
npx license-checker --summary
npx license-checker --onlyAllow "MIT;ISC;Apache-2.0;BSD-2-Clause;BSD-3-Clause"
```

**Python**:
```bash
pip-licenses
pip-licenses --allow-only="MIT;Apache-2.0;BSD"
```

### 3.2 License Compatibility Matrix

| Your License | Compatible With |
|-------------|-----------------|
| MIT | Most licenses |
| Apache-2.0 | MIT, BSD, Apache |
| GPL-3.0 | GPL-3.0 only |
| Proprietary | MIT, BSD, Apache (check terms) |

### 3.3 Red Flag Licenses

| License | Risk | Action |
|---------|------|--------|
| GPL-3.0 | High | Requires source disclosure |
| AGPL-3.0 | High | Network copyleft |
| SSPL | High | Service restrictions |
| Unlicensed | High | No rights granted |
| WTFPL | Medium | Ambiguous |

---

## Phase 4: Plan Updates

### 4.1 Prioritize Updates

**Priority Order**:
1. Security vulnerabilities (Critical/High)
2. Security vulnerabilities (Moderate)
3. Patch updates (low risk)
4. Minor updates (medium risk)
5. Major updates (high risk)

### 4.2 Batch vs Individual

| Approach | When | Risk |
|----------|------|------|
| **Individual** | Major updates, risky deps | Lower |
| **Batched** | Patches, minor updates | Medium |
| **All at once** | Fresh project, comprehensive testing | Higher |

### 4.3 Update Plan

```markdown
## Update Plan - 2025-01-15

### Batch 1: Security Fixes (Immediate)
- [ ] axios: 0.21.0 → 0.21.4 (security)
- [ ] lodash: 4.17.15 → 4.17.21 (security)

### Batch 2: Patch Updates (This Week)
- [ ] react: 18.2.0 → 18.2.1
- [ ] typescript: 5.0.0 → 5.0.4

### Batch 3: Minor Updates (This Sprint)
- [ ] jest: 29.5.0 → 29.7.0
- [ ] eslint: 8.45.0 → 8.56.0

### Batch 4: Major Updates (Next Sprint)
- [ ] next: 14.0.0 → 15.0.0
  - Breaking changes: [link to changelog]
  - Migration guide: [link]
  - Estimated effort: 4 hours
```

---

## Phase 5: Execute Updates

### 5.1 Create Update Branch

```bash
git checkout -b chore/dependency-updates-2025-01
```

### 5.2 Update One at a Time (Risky Deps)

For major updates or risky dependencies:

```bash
# Node.js
npm install package@version

# Python
pip install package==version

# Go
go get package@version

# Rust
cargo update -p package --precise version
```

### 5.3 Batch Update (Safe Deps)

For patches and minor updates:

```bash
# Node.js - Update all patches
npm update

# Python - Update to latest
pip install -U package1 package2 package3

# Go - Update all
go get -u ./...

# Rust - Update all
cargo update
```

### 5.4 Lock File Update

After updates, verify lock files are updated:

- `package-lock.json` / `yarn.lock` / `pnpm-lock.yaml`
- `requirements.txt` / `poetry.lock` / `Pipfile.lock`
- `go.sum`
- `Cargo.lock`
- `Gemfile.lock`

### 5.5 Commit Strategy

**Individual commits for major updates**:
```bash
git add package.json package-lock.json
git commit -m "chore(deps): upgrade next from 14.0.0 to 15.0.0

Breaking changes addressed:
- Updated App Router usage
- Migrated deprecated APIs

Refs: #123"
```

**Single commit for batched updates**:
```bash
git add package.json package-lock.json
git commit -m "chore(deps): update dependencies - 2025-01

Security fixes:
- axios: 0.21.0 → 0.21.4 (CVE-2021-3749)
- lodash: 4.17.15 → 4.17.21 (CVE-2021-23337)

Patch updates:
- react: 18.2.0 → 18.2.1
- typescript: 5.0.0 → 5.0.4"
```

---

## Phase 6: Test & Validate

### 6.1 Test Suite

Run full test suite:

```bash
# Node.js
npm test
npm run test:e2e

# Python
pytest
pytest --cov

# Go
go test ./...

# Rust
cargo test
```

### 6.2 Type Check

```bash
# TypeScript
npm run typecheck

# Python
mypy .

# Rust
cargo check
```

### 6.3 Lint Check

```bash
# Node.js
npm run lint

# Python
ruff check .

# Go
golangci-lint run

# Rust
cargo clippy
```

### 6.4 Build Check

```bash
# Node.js
npm run build

# Go
go build ./...

# Rust
cargo build --release
```

### 6.5 Manual Testing

For major updates:

- [ ] Critical paths work (login, checkout, etc.)
- [ ] No console errors
- [ ] No visual regressions
- [ ] Performance acceptable

---

## Phase 7: Document & Deploy

### 7.1 Update Documentation

If APIs changed:

- [ ] Update README if needed
- [ ] Update API documentation
- [ ] Update migration notes

### 7.2 Create Pull Request

```markdown
## Dependency Updates - January 2025

### Security Fixes
- axios: 0.21.0 → 0.21.4 (CVE-2021-3749)
- lodash: 4.17.15 → 4.17.21 (CVE-2021-23337)

### Updates
| Package | From | To | Type |
|---------|------|-----|------|
| react | 18.2.0 | 18.2.1 | Patch |
| next | 14.0.0 | 15.0.0 | Major |
| typescript | 5.0.0 | 5.0.4 | Patch |

### Breaking Changes
- Next.js 15: [Changes addressed]

### Testing
- [x] Unit tests pass
- [x] Integration tests pass
- [x] E2E tests pass
- [x] Manual testing complete

### Rollback
If issues arise, revert this PR and pin versions:
```bash
git revert <commit-hash>
```
```

### 7.3 Deploy Strategy

| Environment | Strategy |
|-------------|----------|
| Development | Merge immediately |
| Staging | Deploy after PR merge |
| Production | Deploy after staging validation |

---

## Rollback Procedures

### If Tests Fail

```bash
# Reset to before updates
git checkout package.json package-lock.json
npm install
```

### If Production Issues

```bash
# Revert the commit
git revert <update-commit-hash>
npm install
# Deploy revert
```

### Pin Problematic Dependency

```json
// package.json
{
  "dependencies": {
    "problematic-package": "1.2.3"  // Pin to working version
  },
  "resolutions": {
    "problematic-package": "1.2.3"  // Force transitive deps
  }
}
```

---

## Automation

### Dependabot (GitHub)

```yaml
# .github/dependabot.yml
version: 2
updates:
  - package-ecosystem: "npm"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
    groups:
      patch-updates:
        patterns:
          - "*"
        update-types:
          - "patch"
```

### Renovate

```json
// renovate.json
{
  "extends": ["config:base"],
  "packageRules": [
    {
      "matchUpdateTypes": ["patch"],
      "automerge": true
    },
    {
      "matchUpdateTypes": ["major"],
      "labels": ["breaking-change"]
    }
  ]
}
```

---

## Quick Reference

### Commands by Language

| Task | Node.js | Python | Go | Rust |
|------|---------|--------|----|------|
| List outdated | `npm outdated` | `pip list --outdated` | `go list -u -m all` | `cargo outdated` |
| Security audit | `npm audit` | `pip-audit` | `govulncheck ./...` | `cargo audit` |
| Update all | `npm update` | `pip install -U` | `go get -u ./...` | `cargo update` |
| Update one | `npm install pkg@ver` | `pip install pkg==ver` | `go get pkg@ver` | `cargo update -p pkg` |

---

## Checklist

### Pre-Update
- [ ] Tests passing
- [ ] Clean git state
- [ ] Outdated list generated
- [ ] Vulnerabilities checked
- [ ] Licenses checked
- [ ] Update plan created

### During Update
- [ ] Branch created
- [ ] Updates applied
- [ ] Lock files updated
- [ ] Commits atomic and descriptive

### Post-Update
- [ ] All tests pass
- [ ] Type checks pass
- [ ] Lint passes
- [ ] Build succeeds
- [ ] Manual testing done
- [ ] PR created
- [ ] Rollback plan ready

---

## Related Workflows

- [security-audit.md](security-audit.md) - Includes vulnerability scanning
- [code-review.md](code-review.md) - Review updated code
- [troubleshooting.md](troubleshooting.md) - If updates cause issues
