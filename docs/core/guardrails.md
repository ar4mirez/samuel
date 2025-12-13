---
title: Guardrails
description: 35+ testable rules for AI-assisted development
---

# Guardrails

35+ specific, testable rules that AI follows to ensure consistent, high-quality code.

---

## Overview

Guardrails are **not suggestions** - they're enforceable rules. AI validates each one during the Develop and Deliver phases.

```
✓ = MUST do (enforced)
```

---

## Code Quality

### Function Size

```
✓ No function exceeds 50 lines
```

**Why**: Long functions are hard to test, understand, and maintain.

**How to fix**: Extract helper functions for distinct operations.

```typescript
// ❌ Bad: 80-line function
function processOrder(order: Order) {
  // 80 lines of code...
}

// ✅ Good: Split into focused functions
function processOrder(order: Order) {
  validateOrder(order);
  calculateTotals(order);
  applyDiscounts(order);
  submitOrder(order);
}
```

---

### File Size

```
✓ No file exceeds 300 lines
  - Components: 200 lines
  - Tests: 300 lines
  - Utilities: 150 lines
```

**Why**: Large files indicate multiple responsibilities.

**How to fix**: Split into multiple files by concern.

---

### Complexity

```
✓ Cyclomatic complexity ≤10 per function
```

**Why**: High complexity = more paths = more bugs = harder testing.

**How to measure**: Most linters calculate this automatically.

**How to fix**: Extract conditions, use early returns, simplify logic.

```typescript
// ❌ Bad: Complexity 12
function getDiscount(user, cart, season, dayOfWeek) {
  if (user.isPremium) {
    if (cart.total > 100) {
      if (season === 'holiday') {
        // ...nested conditions continue
      }
    }
  }
}

// ✅ Good: Complexity 4
function getDiscount(user, cart, season, dayOfWeek) {
  if (!user.isPremium) return 0;
  if (cart.total <= 100) return 0;

  return calculateSeasonalDiscount(season, dayOfWeek);
}
```

---

### Type Signatures

```
✓ All exported functions have type signatures and documentation
```

**Why**: Self-documenting code, better IDE support, fewer bugs.

```typescript
// ❌ Bad: No types
export function formatCurrency(amount) {
  return `$${amount.toFixed(2)}`;
}

// ✅ Good: Typed and documented
/**
 * Formats a number as USD currency.
 * @param amount - The amount to format
 * @returns Formatted string like "$10.00"
 */
export function formatCurrency(amount: number): string {
  return `$${amount.toFixed(2)}`;
}
```

---

### No Magic Numbers

```
✓ No magic numbers (use named constants)
```

**Why**: Self-documenting, easy to change, prevents errors.

```typescript
// ❌ Bad
if (retries > 3) { ... }
setTimeout(fn, 86400000);

// ✅ Good
const MAX_RETRIES = 3;
const ONE_DAY_MS = 24 * 60 * 60 * 1000;

if (retries > MAX_RETRIES) { ... }
setTimeout(fn, ONE_DAY_MS);
```

---

### No Dead Code

```
✓ No commented-out code in commits
✓ No TODO without issue/ticket reference
✓ No dead code (unused imports, variables, functions)
```

**Why**: Confusing, clutters codebase, suggests incomplete work.

```typescript
// ❌ Bad
// function oldImplementation() { ... }
const unusedVariable = 42;
// TODO: fix this later

// ✅ Good
// TODO(#123): Implement caching for performance
```

---

## Security (CRITICAL)

!!! danger "Security guardrails are non-negotiable"

    These rules protect against OWASP Top 10 vulnerabilities.

### Input Validation

```
✓ All user inputs validated before processing
✓ All API boundaries have input validation
```

**Why**: Prevents injection attacks, data corruption.

```typescript
// ❌ Bad: No validation
app.post('/users', (req, res) => {
  const user = createUser(req.body);
});

// ✅ Good: Schema validation
const CreateUserSchema = z.object({
  email: z.string().email(),
  age: z.number().int().positive(),
});

app.post('/users', (req, res) => {
  const validated = CreateUserSchema.parse(req.body);
  const user = createUser(validated);
});
```

---

### Parameterized Queries

```
✓ All database queries use parameterized statements
```

**Why**: Prevents SQL injection.

```typescript
// ❌ DANGEROUS: SQL injection vulnerability
db.query(`SELECT * FROM users WHERE email = '${email}'`);

// ✅ Safe: Parameterized query
db.query('SELECT * FROM users WHERE email = ?', [email]);
```

---

### No Hardcoded Secrets

```
✓ All environment variables have secure defaults
✓ Never hardcode secrets
```

**Why**: Secrets in code = secrets in version history forever.

```typescript
// ❌ DANGEROUS: Hardcoded secret
const API_KEY = 'sk_live_abc123xyz';

// ✅ Safe: Environment variable
const API_KEY = process.env.API_KEY;
if (!API_KEY) throw new Error('API_KEY required');
```

---

### File Path Validation

```
✓ All file operations validate paths
```

**Why**: Prevents directory traversal attacks.

```typescript
// ❌ DANGEROUS: Path traversal vulnerability
const file = fs.readFileSync(`./uploads/${req.params.filename}`);

// ✅ Safe: Validate path
const safePath = path.join('./uploads', path.basename(req.params.filename));
if (!safePath.startsWith('./uploads/')) {
  throw new Error('Invalid path');
}
```

---

### Dependency Security

```
✓ Dependencies checked for known vulnerabilities before adding
✓ Dependencies checked for license compatibility
```

**How**:

```bash
# npm
npm audit

# Python
pip-audit

# Go
govulncheck ./...

# Rust
cargo audit
```

---

### Async Safety

```
✓ All async operations have timeout/cancellation mechanisms
```

**Why**: Prevents hanging requests, resource exhaustion.

```typescript
// ❌ Bad: No timeout
const response = await fetch(url);

// ✅ Good: With timeout
const controller = new AbortController();
const timeout = setTimeout(() => controller.abort(), 5000);
const response = await fetch(url, { signal: controller.signal });
clearTimeout(timeout);
```

---

### Migration Rollback

```
✓ All database migrations include rollback (down) function
```

**Why**: Safe deployments, quick recovery from issues.

```sql
-- Migration: 001_add_users_table.sql

-- UP
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL
);

-- DOWN
DROP TABLE users;
```

---

## Testing (CRITICAL)

!!! warning "Testing guardrails ensure reliability"

    Untested code is broken code waiting to happen.

### Coverage Targets

```
✓ >80% coverage for business logic
✓ >60% overall coverage
```

**Why**: Confidence in changes, catches regressions.

**How to measure**:

```bash
npm run test:cov
pytest --cov=src
go test -cover ./...
cargo tarpaulin
```

---

### Test Requirements

```
✓ All public APIs have unit tests
✓ All bug fixes include regression tests
✓ All edge cases explicitly tested (null, empty, boundary values)
```

**Why**: Prevents regressions, documents expected behavior.

```typescript
describe('formatCurrency', () => {
  it('formats positive numbers', () => { ... });
  it('formats zero', () => { ... });
  it('formats negative numbers', () => { ... });
  it('handles very large numbers', () => { ... });
  it('handles decimal precision', () => { ... });
});
```

---

### Test Quality

```
✓ Test names describe behavior
✓ No test interdependencies (tests run in any order)
```

**Why**: Self-documenting, reliable, maintainable.

```typescript
// ❌ Bad: Vague name
it('test1', () => { ... });

// ✅ Good: Describes behavior
it('should return 401 when token is expired', () => { ... });
```

---

### Integration Testing

```
✓ Integration tests for external service interactions
✓ All deployments include smoke test validation
```

**Why**: Unit tests don't catch integration issues.

---

## Git & Commits

### Conventional Commits

```
✓ Commit messages: type(scope): description
✓ Types: feat, fix, docs, refactor, test, chore, perf, ci
```

**Format**:

```
type(scope): short description

- Detail 1
- Detail 2

Refs: #issue-number
```

**Examples**:

```
feat(auth): add OAuth login support
fix(profile): handle null email gracefully
refactor(api): extract validation middleware
test(users): add integration tests for search
docs(readme): update installation instructions
chore(deps): upgrade React to v18
perf(query): add index for user email lookup
ci(actions): add automated deployment
```

---

### Atomic Commits

```
✓ One logical change per commit
✓ All commits must pass tests before pushing
```

**Why**: Easy to review, easy to revert, clean history.

```bash
# ❌ Bad: Multiple concerns
git commit -m "fix bug and add feature and update docs"

# ✅ Good: Atomic commits
git commit -m "fix(auth): validate email format"
git commit -m "feat(auth): add password strength indicator"
git commit -m "docs(auth): update login flow diagram"
```

---

### Branch Protection

```
✓ Branch naming: type/short-description
✓ No commits directly to main/master
✓ Breaking API changes require major version bump
```

**Examples**:

```
feat/user-auth
fix/login-button
refactor/api-middleware
```

---

## Performance

```
✓ No N+1 queries (batch database operations)
✓ Large datasets use pagination/streaming
✓ Expensive computations memoized/cached
✓ Frontend bundles < 200KB initial load
✓ API responses < 200ms for simple queries, < 1s for complex
```

**N+1 Query Fix**:

```typescript
// ❌ Bad: N+1 queries
const users = await User.findAll();
for (const user of users) {
  user.posts = await Post.findAll({ where: { userId: user.id } });
}

// ✅ Good: Eager loading
const users = await User.findAll({
  include: [Post],
});
```

---

## Quick Reference

### Per-Task Checklist

- [ ] All guardrails validated
- [ ] Tests pass with coverage thresholds
- [ ] No security vulnerabilities
- [ ] Documentation updated
- [ ] Commit follows conventions

### Validation Commands

```bash
# TypeScript
npm test && npm run lint && npm run typecheck

# Python
pytest && black --check . && mypy .

# Go
go test ./... && golangci-lint run

# Rust
cargo test && cargo clippy && cargo fmt --check
```

---

## Related

<div class="grid cards" markdown>

-   :material-cog:{ .lg .middle } **4D Methodology**

    ---

    When guardrails are applied.

    [:octicons-arrow-right-24: Methodology](methodology.md)

-   :material-code-braces:{ .lg .middle } **Language Guides**

    ---

    Language-specific rules.

    [:octicons-arrow-right-24: Languages](../languages/index.md)

-   :material-bug:{ .lg .middle } **Troubleshooting**

    ---

    When guardrails fail.

    [:octicons-arrow-right-24: Troubleshooting](../workflows/troubleshooting.md)

</div>
