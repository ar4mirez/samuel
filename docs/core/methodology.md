---
title: The 4D Methodology
description: A systematic approach to AI-assisted development
---

# The 4D Methodology

A structured approach that scales from simple bug fixes to complex system architecture.

---

## Overview

Every task, regardless of size, follows four phases:

```mermaid
graph LR
    A[DECONSTRUCT] --> B[DIAGNOSE]
    B --> C[DEVELOP]
    C --> D[DELIVER]

    style A fill:#e1f5fe
    style B fill:#fff3e0
    style C fill:#e8f5e9
    style D fill:#fce4ec
```

| Phase | Purpose | Key Output |
|-------|---------|------------|
| **Deconstruct** | Break down the task | Clear scope, subtasks |
| **Diagnose** | Identify risks | Dependencies, integration points |
| **Develop** | Implement with tests | Working code, passing tests |
| **Deliver** | Validate and commit | Production-ready code |

---

## The 3 Modes

AI auto-detects which mode based on task complexity:

### ATOMIC Mode

**Triggers**: <5 files, clear scope, single concern

**Examples**:

- Fix a bug
- Add a button
- Update styling
- Fix typo
- Add validation

**Process**:

```
Deconstruct: What's the minimal change?
     ↓
Diagnose: Will this break anything? Check dependencies.
     ↓
Develop: Make the change with tests.
     ↓
Deliver: Validate guardrails → Commit
```

**Timeline**: Single session, one commit

---

### FEATURE Mode

**Triggers**: 5-10 files, multiple concerns, integration needed

**Examples**:

- New API endpoint
- New component
- Refactoring
- Adding a service
- Multi-step form

**Process**:

```
Deconstruct: Break into 3-5 subtasks (each atomic)
     ↓
Diagnose: Identify integration points and dependencies
     ↓
Develop: Implement subtasks sequentially with tests
     ↓
Deliver: Integration test → Documentation → Review → Commit
```

**Timeline**: Hours to days, multiple commits

---

### COMPLEX Mode

**Triggers**: >10 files, new subsystem, architectural change, unclear requirements

**Examples**:

- Authentication system
- Payment integration
- Real-time features
- Database migration
- Microservice extraction

**Process**:

```
Deconstruct: Full decomposition into phases/milestones
     ↓
Diagnose: Analyze risks, dependencies, migration paths
     ↓
Develop: Create PRD → Generate tasks → Implement incrementally
     ↓
Deliver: Staged rollout → Documentation → Retrospective
```

**Timeline**: Days to weeks, many commits, possibly multiple PRs

!!! tip "PRD Workflow"

    For COMPLEX mode, use the PRD workflow:

    ```
    @.claude/skills/create-prd/SKILL.md
    ```

    This creates a structured requirements document before implementation.

---

## Escalation Triggers

AI automatically escalates when:

| Condition | Action |
|-----------|--------|
| Task affects >5 files | → FEATURE mode |
| Task affects >10 files | → COMPLEX mode (consider PRD) |
| Task affects >15 files OR new subsystem | → COMPLEX mode (PRD MANDATORY) |
| Requirements unclear | → Ask user for clarification first |

---

## Phase Details

### Phase 1: Deconstruct

**Goal**: Understand what needs to be done

=== "ATOMIC"

    - What's the minimal change?
    - Where does this change live?
    - Is the scope clear?

=== "FEATURE"

    - What are the subtasks?
    - What order should they be done?
    - What can be done in parallel?

=== "COMPLEX"

    - What are the phases/milestones?
    - What's the MVP?
    - What can be deferred?

**Output**: Clear scope definition, task breakdown (if needed)

---

### Phase 2: Diagnose

**Goal**: Identify what could go wrong

=== "ATOMIC"

    - Will this break existing functionality?
    - Are there dependent files?
    - Any edge cases?

=== "FEATURE"

    - Integration points with existing code?
    - Database changes needed?
    - API contracts affected?

=== "COMPLEX"

    - Migration risks?
    - Performance implications?
    - Security considerations?
    - Rollback strategy?

**Output**: Risk assessment, mitigation plan

**Checkpoints**:

- [ ] Requirements clear and testable
- [ ] Scope defined (included/excluded)
- [ ] Dependencies identified
- [ ] Breaking changes flagged

---

### Phase 3: Develop

**Goal**: Implement the solution

**Always**:

- Write tests first (TDD) or alongside code
- Follow language-specific guide
- Validate against guardrails continuously
- Keep changes focused

**Guardrails Applied**:

```
✓ Functions ≤50 lines
✓ Files ≤300 lines
✓ All inputs validated
✓ Parameterized queries
✓ Tests written
✓ No magic numbers
```

**Output**: Working code with tests

**Checkpoints**:

- [ ] Code follows all guardrails
- [ ] Tests written and passing
- [ ] No linter errors
- [ ] Types correct

---

### Phase 4: Deliver

**Goal**: Ensure code is production-ready

**Automated Validation**:

```bash
# Run tests
npm test           # or pytest, go test, cargo test

# Check coverage
npm run test:cov

# Lint
npm run lint

# Build
npm run build
```

**Manual Validation**:

- [ ] Edge cases considered
- [ ] Error handling implemented
- [ ] Performance acceptable
- [ ] Security reviewed

**Commit**:

```bash
git add <files>
git commit -m "type(scope): description

- Detail 1
- Detail 2

Refs: #issue-number"
```

---

## Example Walkthrough

### ATOMIC Example: Fix Login Button

**Task**: "The login button is disabled after form validation passes"

**Deconstruct**:

- Single component issue
- Likely state management bug
- Clear scope

**Diagnose**:

- Check button disabled condition
- Check form validation state
- No other components affected

**Develop**:

```typescript
// Found: button disabled when formValid is null (not false)
// Fix: Explicit boolean check
disabled={formValid !== true}
```

- Add test for edge case

**Deliver**:

```bash
npm test  # Passes
git commit -m "fix(auth): enable login button when form is valid

- Fix null check in disabled condition
- Add test for form validation state edge case

Fixes: #123"
```

---

### FEATURE Example: Add Avatar Upload

**Task**: "Add avatar upload to user profile"

**Deconstruct**:

1. File input component
2. Preview component
3. Validation logic (type, size)
4. Upload API call
5. Error handling
6. Tests

**Diagnose**:

- Uses existing `<FileInput>` component? No → create new
- API endpoint exists? No → need backend change
- Storage: S3 or local? → Check with user

**Develop** (sequential):

```
Subtask 1: Create AvatarUpload component
           → Test component renders
           → Commit: feat(profile): add avatar upload component

Subtask 2: Add file validation
           → Test validation (type, size)
           → Commit: feat(profile): add avatar file validation

Subtask 3: Add preview functionality
           → Test preview shows
           → Commit: feat(profile): add avatar preview

Subtask 4: Connect to API
           → Test upload success/failure
           → Commit: feat(profile): connect avatar upload to API
```

**Deliver**:

- Integration test: Full upload flow
- Documentation: Update user guide
- PR: Ready for review

---

### COMPLEX Example: Authentication System

**Task**: "Build user authentication with OAuth"

**Deconstruct** (PRD workflow):

```
@.claude/skills/create-prd/SKILL.md
```

AI asks clarifying questions, creates PRD with:

- User stories
- Functional requirements
- Non-goals
- Technical considerations
- Security requirements

**Diagnose** (from PRD):

- OAuth providers: Google, GitHub
- Session storage: JWT + Redis
- Rate limiting: 5 attempts = 15min lockout
- Migration: Add users table

**Develop** (task breakdown):

```
Phase 1: Email/Password Auth (Tasks 1.1-1.8)
Phase 2: OAuth Integration (Tasks 2.1-2.6)
Phase 3: Session Management (Tasks 3.1-3.4)
Phase 4: Password Reset (Tasks 4.1-4.3)
```

**Deliver**:

- Phase-by-phase deployment
- Feature flags for gradual rollout
- Documentation for all flows
- Security audit

---

## Anti-Patterns

!!! danger "Don't Do This"

    - **Skip Diagnose**: "It's a simple change" → causes regressions
    - **Big Bang Commits**: All changes in one commit → hard to review/revert
    - **Skip Tests**: "I'll add them later" → technical debt
    - **Ignore Guardrails**: "Just this once" → quality degrades

!!! success "Do This Instead"

    - Always diagnose, even for "simple" changes
    - Commit after each logical change
    - Tests alongside code
    - Guardrails are non-negotiable

---

## Related

<div class="grid cards" markdown>

-   :material-shield:{ .lg .middle } **Guardrails**

    ---

    Rules applied during Develop phase.

    [:octicons-arrow-right-24: View Guardrails](guardrails.md)

-   :material-file-document:{ .lg .middle } **PRD Workflow**

    ---

    Structured planning for COMPLEX mode.

    [:octicons-arrow-right-24: Create PRD](../workflows/create-prd.md)

-   :material-list-status:{ .lg .middle } **Task Generation**

    ---

    Break PRDs into actionable tasks.

    [:octicons-arrow-right-24: Generate Tasks](../workflows/generate-tasks.md)

</div>
