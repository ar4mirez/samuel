# Security Audit Workflow

> **Purpose**: Proactive security assessment covering OWASP Top 10, dependency vulnerabilities, secrets detection, and security best practices.

---

## When to Use

| Trigger | Priority | Description |
|---------|----------|-------------|
| **Pre-Production** | Critical | Before any production deployment |
| **Monthly Review** | High | Regular security hygiene |
| **Auth Changes** | Critical | After adding/modifying authentication |
| **External Integration** | High | When adding third-party services |
| **Dependency Updates** | Medium | After major dependency changes |
| **Security Incident** | Critical | Post-incident review |

---

## Audit Scope

### Full Audit
Complete security review across all categories. Time: 2-4 hours.

### Focused Audit
Target specific area (e.g., authentication only). Time: 30-60 minutes.

### Quick Scan
Automated checks only (dependencies, secrets). Time: 5-10 minutes.

---

## Prerequisites

Before starting audit:

- [ ] Access to codebase and dependencies
- [ ] Access to environment configuration (sanitized)
- [ ] List of external services/APIs used
- [ ] Authentication flow documentation (if exists)
- [ ] Previous audit reports (if available)

---

## Audit Process

```
Phase 1: OWASP Top 10 Review
    ↓
Phase 2: Dependency Vulnerability Scan
    ↓
Phase 3: Secrets Detection
    ↓
Phase 4: Input Validation Audit
    ↓
Phase 5: Authentication & Authorization
    ↓
Phase 6: API Security
    ↓
Phase 7: Report & Remediation
```

---

## Phase 1: OWASP Top 10 Review

### A01:2021 - Broken Access Control

**Check**: Authorization on all sensitive operations

```
Questions:
- Are all endpoints protected with authorization checks?
- Is there role-based access control (RBAC)?
- Can users access other users' data?
- Are direct object references protected?
```

**Red Flags**:
```javascript
// FAIL: No authorization check
app.get('/api/users/:id', (req, res) => {
  return db.getUser(req.params.id)  // Any user can access any ID
})

// PASS: Authorization verified
app.get('/api/users/:id', authorize, (req, res) => {
  if (req.user.id !== req.params.id && !req.user.isAdmin) {
    return res.status(403).json({ error: 'Forbidden' })
  }
  return db.getUser(req.params.id)
})
```

### A02:2021 - Cryptographic Failures

**Check**: Sensitive data protection

```
Questions:
- Is sensitive data encrypted at rest?
- Is TLS enforced for all connections?
- Are passwords properly hashed (bcrypt, argon2)?
- Are encryption keys securely managed?
```

**Requirements**:
- Passwords: bcrypt (cost 10+) or argon2
- Tokens: cryptographically random (crypto.randomBytes)
- TLS: Version 1.2+ only
- No MD5/SHA1 for security purposes

### A03:2021 - Injection

**Check**: All user input handling

```
Questions:
- Are SQL queries parameterized?
- Is user input escaped in templates?
- Are OS commands avoided (or properly escaped)?
- Is JSON/XML parsing safe?
```

**Injection Types**:
| Type | Check | Mitigation |
|------|-------|------------|
| SQL | Parameterized queries | Use ORM or prepared statements |
| NoSQL | Query sanitization | Validate input types |
| XSS | Output encoding | Use templating engine escaping |
| Command | Avoid shell execution | Use language APIs directly |
| LDAP | Input sanitization | Escape special characters |

### A04:2021 - Insecure Design

**Check**: Security in architecture

```
Questions:
- Is there defense in depth?
- Are trust boundaries defined?
- Is least privilege applied?
- Are security requirements documented?
```

### A05:2021 - Security Misconfiguration

**Check**: Configuration security

```
Questions:
- Are default credentials changed?
- Are unnecessary features disabled?
- Are error messages generic (no stack traces)?
- Are security headers configured?
```

**Required Headers**:
```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Content-Security-Policy: default-src 'self'
Strict-Transport-Security: max-age=31536000
X-XSS-Protection: 0 (rely on CSP instead)
```

### A06:2021 - Vulnerable Components

**Check**: Dependency security (see Phase 2)

### A07:2021 - Authentication Failures

**Check**: Authentication implementation (see Phase 5)

### A08:2021 - Data Integrity Failures

**Check**: Software and data integrity

```
Questions:
- Are dependencies verified (checksums)?
- Is CI/CD pipeline secured?
- Are updates authenticated?
- Is deserialization safe?
```

### A09:2021 - Security Logging Failures

**Check**: Logging and monitoring

```
Questions:
- Are authentication events logged?
- Are access control failures logged?
- Are logs protected from tampering?
- Is there alerting for suspicious activity?
```

**Required Logging**:
- Login attempts (success/failure)
- Password changes
- Permission changes
- Access denied events
- Admin actions

### A10:2021 - Server-Side Request Forgery

**Check**: External request handling

```
Questions:
- Are user-provided URLs validated?
- Are internal services protected?
- Is URL schema restricted (http/https only)?
- Are redirects validated?
```

---

## Phase 2: Dependency Vulnerability Scan

### 2.1 Run Audit Commands

**Node.js**:
```bash
npm audit
npm audit --json > audit-report.json

# For detailed report
npm audit --audit-level=moderate
```

**Python**:
```bash
pip-audit
pip-audit --format=json > audit-report.json

# Or using safety
safety check --json > audit-report.json
```

**Go**:
```bash
govulncheck ./...
```

**Rust**:
```bash
cargo audit
cargo audit --json > audit-report.json
```

**Ruby**:
```bash
bundle audit check
```

### 2.2 Analyze Results

For each vulnerability:

| Severity | Action | Timeline |
|----------|--------|----------|
| Critical | Immediate fix or remove | Hours |
| High | Fix in current sprint | Days |
| Moderate | Schedule fix | Weeks |
| Low | Track for update | Next release |

### 2.3 Check for Updates

```bash
# Node.js
npm outdated

# Python
pip list --outdated

# Go
go list -u -m all

# Rust
cargo outdated
```

---

## Phase 3: Secrets Detection

### 3.1 Automated Scanning

**Using git-secrets**:
```bash
git secrets --scan
git secrets --scan-history
```

**Using gitleaks**:
```bash
gitleaks detect --source . --verbose
```

**Using truffleHog**:
```bash
trufflehog filesystem .
```

### 3.2 Manual Review

Check for patterns:

| Pattern | Example | Location |
|---------|---------|----------|
| API Keys | `sk_live_`, `AKIA` | Config files, code |
| Passwords | `password=`, `passwd` | Config, .env |
| Tokens | `token=`, `bearer` | Headers, config |
| Connection Strings | `mongodb://user:pass@` | Config |
| Private Keys | `-----BEGIN RSA` | Any file |
| AWS Credentials | `aws_access_key_id` | Config, code |

### 3.3 Environment Variable Audit

```
Check:
- All secrets in environment variables (not code)
- .env files in .gitignore
- No .env files committed historically
- Environment variables have secure defaults
```

---

## Phase 4: Input Validation Audit

### 4.1 Identify All Input Sources

| Source | Examples | Risk Level |
|--------|----------|------------|
| Request body | JSON, form data | High |
| URL parameters | `/users/:id` | High |
| Query strings | `?search=term` | High |
| Headers | `Authorization`, custom | Medium |
| Cookies | Session, preferences | Medium |
| File uploads | Images, documents | Critical |
| WebSocket messages | Real-time data | High |

### 4.2 Validation Checklist

For each input source:

- [ ] Schema validation (Zod, Pydantic, etc.)
- [ ] Type checking
- [ ] Length limits
- [ ] Format validation (email, URL, etc.)
- [ ] Allowlist validation (when applicable)
- [ ] Sanitization for output context

### 4.3 File Upload Security

```
Requirements:
- File type validation (magic bytes, not just extension)
- Size limits enforced
- Virus/malware scanning
- Secure storage location (outside web root)
- Randomized filenames
- No executable permissions
```

---

## Phase 5: Authentication & Authorization

### 5.1 Authentication Review

**Password Security**:
- [ ] Minimum length (12+ characters recommended)
- [ ] Complexity requirements (or passphrase)
- [ ] Hashing: bcrypt (cost 10+) or argon2
- [ ] No password in logs or error messages
- [ ] Rate limiting on login attempts
- [ ] Account lockout policy

**Session Security**:
- [ ] Secure session generation
- [ ] HttpOnly cookie flag
- [ ] Secure cookie flag (HTTPS)
- [ ] SameSite cookie attribute
- [ ] Session timeout
- [ ] Session invalidation on logout
- [ ] Regenerate session on privilege change

**Multi-Factor Authentication**:
- [ ] MFA available for sensitive accounts
- [ ] Recovery codes securely generated
- [ ] MFA bypass prevention

### 5.2 Authorization Review

**Access Control**:
- [ ] Authorization on every endpoint
- [ ] Role-based access control
- [ ] Least privilege principle
- [ ] Deny by default
- [ ] No client-side only checks

**Token Security** (JWT/OAuth):
- [ ] Strong signing algorithm (RS256, ES256)
- [ ] Token expiration
- [ ] Token refresh mechanism
- [ ] Token revocation capability
- [ ] No sensitive data in token payload

---

## Phase 6: API Security

### 6.1 Rate Limiting

```
Requirements:
- Rate limiting on all endpoints
- Stricter limits on authentication endpoints
- Graduated response (slow down, then block)
- Per-user and per-IP limits
```

### 6.2 CORS Configuration

```javascript
// Secure CORS configuration
{
  origin: ['https://app.example.com'],  // Not '*'
  methods: ['GET', 'POST', 'PUT', 'DELETE'],
  allowedHeaders: ['Content-Type', 'Authorization'],
  credentials: true,
  maxAge: 86400
}
```

### 6.3 API Versioning & Deprecation

```
Check:
- API versioning strategy
- Deprecated endpoints documented
- Sunset headers on deprecated endpoints
- Migration path for deprecated features
```

### 6.4 Error Handling

```
Requirements:
- Generic error messages to clients
- Detailed errors in logs only
- No stack traces in production
- Consistent error format
- No information leakage
```

---

## Phase 7: Report & Remediation

### 7.1 Audit Report Template

```markdown
# Security Audit Report

**Date**: 2025-01-15
**Auditor**: AI Assistant
**Scope**: Full Audit
**Duration**: 3 hours

## Executive Summary

| Severity | Count | Status |
|----------|-------|--------|
| Critical | 0 | - |
| High | 2 | Action Required |
| Medium | 5 | Scheduled |
| Low | 8 | Tracked |

**Overall Risk Level**: Medium

## Critical Findings

None identified.

## High Severity Findings

### H1: SQL Injection in User Search
**Location**: `src/api/users.ts:45`
**Description**: User search query built with string concatenation
**Impact**: Database compromise, data theft
**Remediation**: Use parameterized queries
**Timeline**: Immediate (24 hours)

### H2: Missing Rate Limiting on Login
**Location**: `src/auth/login.ts`
**Description**: No rate limiting on login endpoint
**Impact**: Brute force attacks possible
**Remediation**: Add rate limiting (100/15min)
**Timeline**: This sprint

## Medium Severity Findings

[Details...]

## Low Severity Findings

[Details...]

## Recommendations

1. Implement security headers middleware
2. Add automated dependency scanning to CI/CD
3. Enable audit logging for all admin actions
4. Schedule quarterly security reviews

## Appendix

### Tools Used
- npm audit
- gitleaks
- Manual code review

### Files Reviewed
- src/api/*.ts
- src/auth/*.ts
- src/middleware/*.ts
```

### 7.2 Remediation Priority Matrix

| Finding | Severity | Effort | Priority |
|---------|----------|--------|----------|
| SQL Injection | Critical | Low | Immediate |
| Missing Auth | High | Medium | Sprint 1 |
| Weak Hash | High | Low | Sprint 1 |
| Missing Headers | Medium | Low | Sprint 2 |
| Old Dependency | Low | Low | Backlog |

### 7.3 Follow-up Actions

- [ ] Create tickets for all findings
- [ ] Schedule remediation work
- [ ] Plan re-audit after fixes
- [ ] Update security documentation
- [ ] Brief team on findings

---

## Quick Scan Commands

### All-in-One Quick Scan

```bash
# Node.js project
npm audit && npx gitleaks detect

# Python project
pip-audit && gitleaks detect

# Go project
govulncheck ./... && gitleaks detect

# Rust project
cargo audit && gitleaks detect
```

---

## Checklist Summary

### OWASP Top 10
- [ ] A01: Broken Access Control
- [ ] A02: Cryptographic Failures
- [ ] A03: Injection
- [ ] A04: Insecure Design
- [ ] A05: Security Misconfiguration
- [ ] A06: Vulnerable Components
- [ ] A07: Authentication Failures
- [ ] A08: Data Integrity Failures
- [ ] A09: Logging Failures
- [ ] A10: SSRF

### Dependencies
- [ ] Vulnerability scan completed
- [ ] Critical/High issues addressed
- [ ] Outdated packages identified

### Secrets
- [ ] Automated scan completed
- [ ] No secrets in code
- [ ] Environment variables secured

### Authentication
- [ ] Password policy enforced
- [ ] Session security configured
- [ ] MFA available (if applicable)

### API Security
- [ ] Rate limiting enabled
- [ ] CORS configured
- [ ] Security headers set
- [ ] Error handling secure

---

## Related Workflows

- [code-review.md](code-review.md) - Includes security checks
- [dependency-update.md](dependency-update.md) - Safe dependency updates
- [troubleshooting.md](troubleshooting.md) - Security incident response
