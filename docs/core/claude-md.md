---
title: CLAUDE.md
description: The core instruction file for AI-assisted development
---

# CLAUDE.md

The main instruction file that AI assistants load automatically. Contains all guardrails, operations, and methodology in ~500 optimized lines.

---

## What is CLAUDE.md?

CLAUDE.md is the "brain" of the AI Code Template system. When an AI assistant starts working on your project, it loads this file to understand:

- **What commands to run** (Operations)
- **What to avoid** (Boundaries)
- **How to write code** (Guardrails)
- **How to approach tasks** (Methodology)

---

## File Structure

```markdown
# CLAUDE.md

## Operations           ← Commands first (AGENTS.md compatible)
  - Setup Commands
  - Testing Commands
  - Build & Deploy
  - Code Style

## Boundaries           ← Protected files and actions

## Quick Reference      ← Task classification, emergency links

## Core Guardrails      ← 35+ testable rules

## 4D Methodology       ← ATOMIC/FEATURE/COMPLEX modes

## Software Dev Lifecycle ← Planning → Implementation → Validation → Delivery

## Context System       ← How .agent/ directory works

## Anti-Patterns        ← What to avoid

## Initialization       ← How to start
## When Stuck          ← Recovery procedures
## Success Criteria    ← How to measure success
## Version & Changelog ← Updates
```

---

## Key Sections

### Operations

The first section - designed for quick reference and AGENTS.md compatibility:

```bash
# Example: Testing Commands
npm test             # Run all tests
npm run test:watch   # Watch mode
npm run test:cov     # With coverage (target: >80% business logic)

pytest               # Python: Run all tests
go test ./...        # Go: All tests
cargo test           # Rust: All tests
```

!!! info "Commands First"

    Operations section comes first because most AI interactions need quick command access. This follows the AGENTS.md standard.

---

### Boundaries

Critical section defining what AI should NOT modify:

**Protected Files**:

- `package-lock.json`, `yarn.lock`, `Cargo.lock` (dependency locks)
- `.env`, `.env.local` (environment configs)
- Database migrations (once in production)
- CI/CD configurations

**Never Commit**:

- Secrets, API keys, credentials
- `node_modules/`, `venv/`, `target/`
- Personal IDE configs

**Ask Before Modifying**:

- Authentication/authorization logic
- Database schemas (after deployment)
- Public API contracts

---

### Guardrails

35+ testable rules organized by category:

=== "Code Quality"

    ```
    ✓ Functions ≤50 lines
    ✓ Files ≤300 lines (components: 200, tests: 300)
    ✓ Cyclomatic complexity ≤10
    ✓ All exports have types and docs
    ✓ No magic numbers
    ✓ No commented-out code
    ✓ No TODO without issue reference
    ```

=== "Security"

    ```
    ✓ All inputs validated
    ✓ Parameterized queries only
    ✓ No hardcoded secrets
    ✓ File paths validated
    ✓ Async operations have timeouts
    ✓ Dependencies checked for vulnerabilities
    ✓ Dependencies checked for licenses
    ```

=== "Testing"

    ```
    ✓ >80% coverage for business logic
    ✓ >60% overall coverage
    ✓ All public APIs have tests
    ✓ All bugs have regression tests
    ✓ No test interdependencies
    ✓ Integration tests for external services
    ```

=== "Git"

    ```
    ✓ Conventional commits (type(scope): description)
    ✓ One logical change per commit
    ✓ All tests pass before push
    ✓ No direct commits to main
    ✓ Breaking changes = major version bump
    ```

[:octicons-arrow-right-24: Full Guardrails Reference](guardrails.md)

---

### 4D Methodology

The approach for every task:

```mermaid
graph LR
    D1[DECONSTRUCT] --> D2[DIAGNOSE]
    D2 --> D3[DEVELOP]
    D3 --> D4[DELIVER]
```

| Phase | ATOMIC | FEATURE | COMPLEX |
|-------|--------|---------|---------|
| **Deconstruct** | What's minimal change? | Break into 3-5 subtasks | Full decomposition |
| **Diagnose** | Will this break anything? | Integration points? | Risks, dependencies |
| **Develop** | Make change + tests | Implement sequentially | PRD → Tasks → Implement |
| **Deliver** | Validate → Commit | Integration test → Doc | Staged rollout |

[:octicons-arrow-right-24: Methodology Deep Dive](methodology.md)

---

## Cross-Tool Compatibility

CLAUDE.md follows the [AGENTS.md](https://agents.md) standard structure:

- **Operations first** - Commands, setup, boundaries
- **Context second** - Methodology, detailed guidelines

This means:

1. Claude Code reads CLAUDE.md natively
2. Other tools (Cursor, Codex) read AGENTS.md
3. Both files can contain the same content
4. Use symlink: `ln -s CLAUDE.md AGENTS.md`

---

## Customization

### Safe to Customize

- Line length limits (300 → 500)
- Coverage targets (80% → 90%)
- Commit format (add team-specific types)
- Add new guardrails specific to your project

### Customize with Caution

- Security guardrails (only make stricter)
- Testing requirements (only increase)
- Methodology (understand before changing)

### Don't Customize

- Fundamental principles (atomicity, validation)
- Security basics (input validation, parameterized queries)
- Git hygiene (conventional commits, atomic changes)

---

## Example: How AI Uses CLAUDE.md

**User prompt**: "Add a new API endpoint for user search"

**AI process** (following CLAUDE.md):

1. **Mode Detection**: 5-10 files → FEATURE mode
2. **Deconstruct**: Route, controller, service, tests, docs
3. **Diagnose**: Check existing patterns, auth requirements
4. **Develop** (following guardrails):
   - Input validation (Zod schema)
   - Parameterized query
   - Function ≤50 lines
   - Tests with >80% coverage
5. **Deliver**:
   - Run tests
   - Check all guardrails
   - Commit: `feat(api): add user search endpoint`

---

## Token Optimization

CLAUDE.md is optimized for token efficiency:

| Version | Lines | Notes |
|---------|-------|-------|
| v1.0.0 | 490 | Initial release |
| v1.1.0 | 400 | 18% reduction |
| v1.2.0 | 500 | Added Operations (AGENTS.md compatible) |

**Target**: <600 lines for the main file.

Language guides, workflows, and detailed documentation live in `.agent/` to avoid loading unnecessary context.

---

## Viewing the Full File

The complete CLAUDE.md is in your project root. Key line numbers:

| Section | Lines |
|---------|-------|
| Operations | 1-80 |
| Boundaries | 81-110 |
| Quick Reference | 112-135 |
| Core Guardrails | 137-195 |
| 4D Methodology | 197-260 |
| SDLC | 262-330 |
| Context System | 332-380 |
| Anti-Patterns | 382-420 |

---

## Related

<div class="grid cards" markdown>

-   :material-shield:{ .lg .middle } **Guardrails**

    ---

    Deep dive into all 35+ rules.

    [:octicons-arrow-right-24: View Guardrails](guardrails.md)

-   :material-cog:{ .lg .middle } **4D Methodology**

    ---

    Understanding the approach.

    [:octicons-arrow-right-24: Learn Methodology](methodology.md)

-   :material-folder:{ .lg .middle } **.agent Directory**

    ---

    Project context structure.

    [:octicons-arrow-right-24: .agent Structure](agent-directory.md)

-   :material-tools:{ .lg .middle } **Cross-Tool**

    ---

    AGENTS.md compatibility.

    [:octicons-arrow-right-24: Cross-Tool Setup](../reference/cross-tool.md)

</div>
