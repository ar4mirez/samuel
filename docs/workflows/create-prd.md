---
title: Create PRD
description: Plan complex features with a structured Product Requirements Document
---

# Create PRD Workflow

Generate a detailed Product Requirements Document (PRD) for complex features before implementation.

---

## When to Use

Use the PRD workflow when:

- [x] Feature affects **>10 files**
- [x] Building a **new subsystem** (auth, payments, etc.)
- [x] Requirements are **unclear or ambiguous**
- [x] Multiple stakeholders need **alignment**
- [x] Feature will take **>1 week** to implement

Skip PRD for:

- [ ] Bug fixes
- [ ] Small features (<5 files)
- [ ] Clear, well-defined tasks

---

## How to Invoke

```
@.claude/skills/create-prd/SKILL.md

Build a user authentication system with email/password and OAuth
```

---

## What AI Does

### Step 1: Check Context

AI reviews:

- `.claude/project.md` for tech stack
- Existing code patterns
- Related functionality

### Step 2: Ask Clarifying Questions

AI asks about:

- **Problem/Goal**: What problem does this solve?
- **Users**: Who uses this feature?
- **Core Functionality**: What are the must-haves?
- **User Stories**: 2-3 scenarios
- **Scope**: What's NOT included?
- **Data**: What data is involved?
- **Security**: Any sensitive data?
- **Performance**: Scale requirements?

### Step 3: Generate PRD

Creates a comprehensive document covering:

1. Introduction/Overview
2. Goals
3. User Stories
4. Functional Requirements
5. Non-Goals
6. Technical Considerations
7. Design Considerations
8. Guardrails Affected
9. Success Metrics
10. Implementation Estimate
11. Open Questions

### Step 4: Save PRD

Saves to: `.claude/tasks/NNNN-prd-feature-name.md`

---

## PRD Sections

### 1. Introduction

Brief description, problem statement, and goal.

```markdown
## Introduction

This feature adds user authentication to the platform,
allowing users to create accounts and access personalized content.

**Goal**: Enable secure authentication with email/password and OAuth.
```

### 2. Goals

Specific, measurable objectives.

```markdown
## Goals

1. Allow users to create accounts with email and password
2. Enable login via Google and GitHub OAuth
3. Implement secure session management with JWT
4. Achieve <200ms authentication response time
5. Support password reset flow
```

### 3. User Stories

User narratives describing feature usage.

```markdown
## User Stories

**US-001**: As a new user, I want to create an account with my email
so that I can access personalized features.

**US-002**: As a returning user, I want to log in quickly with Google
so that I don't have to remember another password.

**US-003**: As a user who forgot their password, I want to reset it
via email so that I can regain access.
```

### 4. Functional Requirements

Specific functionalities the feature must have.

```markdown
## Functional Requirements

### Authentication
FR-001: The system must allow users to register with email and password
FR-002: The system must validate email format and password strength
FR-003: The system must support OAuth login via Google and GitHub
FR-004: The system must generate and validate JWT tokens

### Security
FR-005: The system must hash passwords using bcrypt (cost factor: 12)
FR-006: The system must implement rate limiting (5 failed = 15min lockout)
FR-007: The system must log all authentication events
```

### 5. Non-Goals

What this feature will NOT include.

```markdown
## Non-Goals

- Two-factor authentication (2FA) - Deferred to v2
- Single Sign-On (SSO) integration - Future enhancement
- Biometric authentication - Not in scope
```

### 6. Technical Considerations

Implementation suggestions based on project context.

```markdown
## Technical Considerations

### Tech Stack Integration
- Use existing Express.js framework
- Add users and sessions tables to PostgreSQL
- Use Prisma for database access
- Use Zod for input validation

### File Organization
- src/auth/ - Main auth module
  - auth.controller.ts
  - auth.service.ts
  - auth.middleware.ts
  - auth.types.ts
```

### 7. Guardrails Affected

Which guardrails are critical for this feature.

```markdown
## Guardrails Affected

### Security (CRITICAL)
- ✓ All user inputs validated
- ✓ All database queries parameterized
- ✓ All passwords hashed with bcrypt
- ✓ All tokens have expiration

### Testing (CRITICAL)
- ✓ >95% coverage for auth logic
- ✓ All edge cases tested
```

### 8. Success Metrics

How success is measured.

```markdown
## Success Metrics

- Authentication response time < 200ms (p95)
- Test coverage >95% for auth module
- Zero critical security vulnerabilities
- 80% of users complete registration flow
```

---

## After PRD Creation

1. **Review the PRD** - Confirm requirements are accurate
2. **Clarify open questions** - Resolve any ambiguities
3. **Generate tasks** - Use generate-tasks workflow

```
@.claude/skills/generate-tasks/SKILL.md

Use the PRD at .claude/tasks/0001-prd-user-auth.md
```

!!! warning "Don't Skip to Implementation"

    After PRD creation, generate tasks before implementing. This ensures proper breakdown and verification steps.

---

## Example PRD

See a complete example: [Example PRD for API Rate Limiting](https://github.com/ar4mirez/samuel/blob/main/.claude/tasks/EXAMPLE-0001-prd-api-rate-limiting.md)

---

## Tips for Better PRDs

### Be Specific

```markdown
# ❌ Vague
"Users should be able to log in"

# ✅ Specific
"FR-001: The system must authenticate users with email/password
and return a JWT token valid for 1 hour"
```

### Include Non-Goals

Prevents scope creep:

```markdown
## Non-Goals
- Social login (Google, GitHub) - Phase 2
- Remember me functionality - Phase 2
- Account deletion - Separate feature
```

### Reference Guardrails

```markdown
## Guardrails Affected
- ✓ All inputs validated (Zod schemas for all endpoints)
- ✓ Parameterized queries (Prisma handles this)
- ✓ >80% test coverage (target 95% for auth)
```

---

## Related

<div class="grid cards" markdown>

-   :material-format-list-numbered:{ .lg .middle } **Generate Tasks**

    ---

    Break the PRD into actionable tasks.

    [:octicons-arrow-right-24: Generate Tasks](generate-tasks.md)

-   :material-shield:{ .lg .middle } **Guardrails**

    ---

    Rules to reference in your PRD.

    [:octicons-arrow-right-24: Guardrails](../core/guardrails.md)

</div>
