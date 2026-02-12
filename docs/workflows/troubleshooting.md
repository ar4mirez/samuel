---
title: Troubleshooting
description: Systematic debugging when stuck
---

# Troubleshooting Workflow

A structured approach to debugging when you've been stuck for more than 30 minutes.

---

## When to Use

- **Stuck >30 minutes** on the same issue
- **Tests failing** unexpectedly
- **Build broken** and unclear why
- **Security issue** discovered
- **Performance degraded** significantly

---

## How to Invoke

```
@.claude/skills/troubleshooting/SKILL.md

Been stuck on "Connection pool exhausted" error for 45 minutes
```

---

## The Process

### Step 1: STOP

**Do not** keep trying random solutions.

!!! danger "Stop Immediately"

    After 30 minutes of the same issue, random attempts waste time and may make things worse.

### Step 2: Document

What have you tried?

```markdown
## Attempted Solutions
1. Restarted the server - still fails
2. Increased pool size to 50 - same error
3. Added connection timeout - no change
```

### Step 3: Simplify

Can you reproduce in isolation?

- Create minimal test case
- Remove complexity
- Isolate the failing component

### Step 4: Check Fundamentals

Common culprits:

- [ ] Dependencies installed and correct version?
- [ ] Configuration correct?
- [ ] Environment variables set?
- [ ] File permissions correct?
- [ ] Right version running? (restart server, clear cache)

### Step 5: Search

- Google the exact error message
- Check GitHub issues for dependencies
- Search Stack Overflow
- Review `.claude/memory/` for similar problems

### Step 6: Ask for Help

Present a clear problem statement:

```markdown
## Problem
Connection pool exhausted after ~100 requests

## What I'm Trying To Do
Handle concurrent API requests

## What Happens Instead
Error: "Connection pool exhausted" after 100 requests

## What I've Tried
1. Increased pool size (10 → 50)
2. Added connection timeout (30s)
3. Restarted server

## Environment
- Node.js 20
- PostgreSQL 16
- pg library 8.11
```

### Step 7: Document Solution

Once resolved, create `.claude/memory/YYYY-MM-DD-issue-name.md`:

```markdown
# Issue: Connection Pool Exhausted

**Date**: 2025-01-15
**Resolved**: Yes

## Problem
Connection pool exhausted after ~100 concurrent requests.

## Root Cause
Connections not being released after query completion.
Missing `finally` block to release connection.

## Solution
```javascript
const client = await pool.connect();
try {
  const result = await client.query(sql);
  return result;
} finally {
  client.release(); // Always release!
}
```

## Prevention
- Add linter rule for unreleased connections
- Add connection monitoring
- Review similar code for same pattern
```

---

## Common Issues

### Tests Breaking Unexpectedly

**Diagnosis Steps**:

1. When did it break?
   ```bash
   git diff HEAD~1
   ```

2. Isolate the failing test:
   ```bash
   npm test -- --grep "specific test"
   ```

3. Check for flakiness (run 10 times)

4. Check for test interdependence

**Recovery**:

```bash
# Revert to last working state
git reset --hard HEAD~1

# Re-apply changes incrementally
# Test after each change
```

---

### Build Failing

**Common Causes**:

- [ ] Dependency version mismatch
- [ ] Missing environment variable
- [ ] TypeScript/compiler error
- [ ] Circular dependency

**Quick Fixes**:

```bash
# Clear and reinstall dependencies
rm -rf node_modules package-lock.json
npm install

# Clear build cache
rm -rf dist .next .cache
npm run build
```

---

### Security Issue Found

!!! danger "PRIORITY: CRITICAL"

    Security issues take priority over all other work.

**Immediate Actions**:

1. **DO NOT** commit vulnerable code
2. **Fix immediately** - Security takes priority
3. **Add regression test** - Prevent reintroduction
4. **Review similar code** - Same pattern elsewhere?
5. **Document** in `.claude/memory/`

**Common Security Fixes**:

=== "SQL Injection"

    ```javascript
    // ❌ Vulnerable
    db.query(`SELECT * FROM users WHERE id = ${userId}`);

    // ✅ Fixed
    db.query('SELECT * FROM users WHERE id = ?', [userId]);
    ```

=== "XSS"

    ```javascript
    // ❌ Vulnerable
    element.innerHTML = userInput;

    // ✅ Fixed
    element.textContent = userInput;
    ```

=== "Hardcoded Secrets"

    ```javascript
    // ❌ Vulnerable
    const API_KEY = "sk_live_abc123";

    // ✅ Fixed
    const API_KEY = process.env.API_KEY;
    ```

---

### Performance Degraded

**Diagnosis**:

1. **Measure** - Don't guess
   ```bash
   # Node.js
   node --prof app.js

   # Python
   python -m cProfile script.py
   ```

2. **Identify bottleneck**
   - N+1 queries?
   - Large data in memory?
   - Missing index?
   - Blocking I/O?

3. **Fix the biggest issue first**

**Common Fixes**:

- Add database indexes
- Use pagination/streaming
- Add caching
- Fix N+1 queries with eager loading

---

## Red Flags

Stop and reassess if you see:

!!! warning "Code Red Flags"

    - Same error after 3 different fixes
    - Solution getting more complex
    - "It just works now" without understanding why
    - Touching >10 files for a "simple" fix
    - Breaking tests to make new code work

!!! warning "Process Red Flags"

    - Skipping tests "I'll add them later"
    - Committing commented-out code
    - Ignoring linter errors
    - Using `any` to "make TypeScript happy"
    - Copying code without understanding it

**When you see red flags**:

1. STOP adding code
2. Revert to last working state
3. Apply COMPLEX mode methodology
4. Ask for guidance

---

## Recovery Checklist

After resolving any major issue:

- [ ] Tests passing
- [ ] Guardrails validated
- [ ] Root cause understood (not just symptom fixed)
- [ ] Similar code reviewed for same issue
- [ ] Solution documented in `.claude/memory/`
- [ ] Prevention added (test, linter rule, etc.)

---

## Getting Help

### From AI

Provide clear context:

```
@.claude/skills/troubleshooting/SKILL.md

Problem: API returns 500 on user creation
Tried: Checking logs, validating input, restarting server
Error: "Cannot read property 'id' of undefined"
```

### From Documentation

- [Guardrails](../core/guardrails.md) - Check if you're violating rules
- [Language Guides](../languages/index.md) - Language-specific debugging

### From Humans

Prepare:

1. Minimal reproduction case
2. Error messages (full text)
3. What you've tried
4. Environment details

---

## Related

<div class="grid cards" markdown>

-   :material-shield:{ .lg .middle } **Guardrails**

    ---

    Rules that help prevent issues.

    [:octicons-arrow-right-24: Guardrails](../core/guardrails.md)

-   :material-cog:{ .lg .middle } **Methodology**

    ---

    Structured approach to development.

    [:octicons-arrow-right-24: 4D Methodology](../core/methodology.md)

</div>
