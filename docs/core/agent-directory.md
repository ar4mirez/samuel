---
title: The .agent Directory
description: Project-specific context that grows over time
---

# The .agent Directory

The `.agent/` directory stores project-specific context that grows organically with your project.

---

## Overview

While CLAUDE.md provides universal guardrails, `.agent/` stores context specific to YOUR project:

```
.agent/
├── README.md              # How to use .agent/
├── project.md             # Your tech stack (created when chosen)
├── patterns.md            # Coding patterns (created when emerge)
├── state.md               # Current work (for multi-session)
├── skills/                # Auto-load based on file type (21 language guide skills)
│   ├── typescript-guide/SKILL.md
│   ├── python-guide/SKILL.md
│   ├── go-guide/SKILL.md
│   ├── rust-guide/SKILL.md
│   ├── kotlin-guide/SKILL.md
│   └── ... (16 more)
├── framework-guides/      # Framework-specific patterns (33 frameworks)
│   ├── react.md
│   ├── django.md
│   ├── rails.md
│   ├── spring-boot-java.md
│   └── ... (29 more)
├── workflows/             # On-demand workflows
│   ├── create-prd.md
│   ├── generate-tasks.md
│   ├── initialize-project.md
│   ├── troubleshooting.md
│   └── generate-agents-md.md
├── tasks/                 # PRDs and task lists (created on demand)
│   └── NNNN-prd-feature-name.md
└── memory/                # Decision logs (created on demand)
    └── YYYY-MM-DD-topic.md
```

---

## File Types

### Pre-Created (Templates)

These files come with the template:

| File | Purpose | Loaded |
|------|---------|--------|
| `README.md` | How to use .agent/ | On-demand |
| `project.md.template` | Template for project.md | Reference |
| `state.md.template` | Template for state.md | Reference |
| `skills/<lang>-guide/SKILL.md` | Language-specific rules (21) | Auto-load |
| `framework-guides/*.md` | Framework-specific patterns (33) | On-demand |
| `workflows/*.md` | Structured workflows | On-demand |

### Created Over Time

These files are created as your project evolves:

| File | Purpose | When Created |
|------|---------|--------------|
| `project.md` | Tech stack, architecture | When stack chosen |
| `patterns.md` | Coding conventions | When patterns emerge |
| `state.md` | Current work status | For multi-session work |
| `tasks/*.md` | PRDs and task lists | COMPLEX mode |
| `memory/*.md` | Decision logs | Significant decisions |

---

## Progressive Growth

The directory grows naturally:

```
Day 1:
.agent/
├── README.md
├── skills/
├── workflows/
└── (templates)

Week 1:
.agent/
├── project.md          ← Created when tech stack chosen
├── skills/
└── workflows/

Month 1:
.agent/
├── project.md
├── patterns.md         ← Created when patterns emerge
├── skills/
├── workflows/
└── tasks/
    └── 0001-prd-auth.md  ← Created for complex feature

Ongoing:
.agent/
├── project.md
├── patterns.md
├── state.md            ← Created for long-running work
├── skills/
├── workflows/
├── tasks/
│   ├── 0001-prd-auth.md
│   └── 0002-prd-search.md
└── memory/
    └── 2025-01-15-cache-strategy.md  ← Created for key decisions
```

!!! tip "Don't Over-Document"

    Let files emerge naturally. Don't create `project.md` on day one - wait until you make architecture decisions.

---

## Key Files

### project.md

Documents your tech stack and architecture decisions:

```markdown
# Project Configuration

## Tech Stack
- **Language**: TypeScript 5.0
- **Runtime**: Node.js 20 LTS
- **Framework**: Express.js 4.18
- **Database**: PostgreSQL 15
- **ORM**: Prisma 5.0
- **Testing**: Vitest + Playwright

## Architecture
- Monolithic API (single service)
- Repository pattern for data access
- Middleware-based request handling

## Conventions
- ESM modules (import/export)
- Zod for runtime validation
- JWT for authentication

## External Services
- AWS S3 for file storage
- SendGrid for email
- Stripe for payments
```

**When to create**: After first major architecture decision.

---

### patterns.md

Documents coding patterns specific to your project:

```markdown
# Coding Patterns

## API Endpoints

All endpoints follow this pattern:

```typescript
// src/routes/users.ts
router.post('/', validateBody(CreateUserSchema), async (req, res) => {
  const user = await userService.create(req.validated);
  res.status(201).json(user);
});
```

## Error Handling

Custom errors with status codes:

```typescript
throw new AppError('User not found', 404);
```

## Database Queries

Always use transactions for writes:

```typescript
await prisma.$transaction(async (tx) => {
  // multiple operations
});
```
```

**When to create**: After 2-3 patterns emerge.

---

### state.md

Tracks current work for multi-session projects:

```markdown
# Current Work State

## Active Task
Building user authentication (PRD: 0001-prd-user-auth.md)

## Progress
- [x] Task 1.1: Database schema
- [x] Task 1.2: User model
- [ ] Task 1.3: Registration endpoint ← IN PROGRESS
- [ ] Task 1.4: Login endpoint

## Blockers
- Need clarification on OAuth providers (asked in Task 2.1)

## Next Steps
1. Complete registration endpoint
2. Add tests for registration
3. Start login endpoint

## Session Notes
- 2025-01-15: Started auth implementation
- 2025-01-16: Completed database schema
```

**When to create**: For work spanning multiple sessions.

---

### tasks/

PRDs and task breakdowns for complex features:

```
tasks/
├── 0001-prd-user-auth.md         # PRD document
├── tasks-0001-prd-user-auth.md   # Task breakdown
├── 0002-prd-search.md
└── tasks-0002-prd-search.md
```

**Naming**: `NNNN-prd-feature-name.md`

**When created**: COMPLEX mode with PRD workflow.

---

### memory/

Decision logs for significant choices:

```
memory/
├── 2025-01-10-database-choice.md
├── 2025-01-15-auth-strategy.md
└── 2025-01-20-caching-approach.md
```

**Format**:

```markdown
# Decision: Caching Strategy

**Date**: 2025-01-20
**Status**: Decided

## Context
API response times increasing as data grows.

## Options Considered
1. Redis cache layer
2. In-memory cache (node-cache)
3. Database query optimization

## Decision
Redis cache layer for shared state across instances.

## Consequences
- Need Redis infrastructure
- Cache invalidation complexity
- Improved response times (expected 50ms → 10ms)
```

**When to create**: Significant architectural decisions.

---

## Loading Protocol

AI follows this protocol when starting a session:

```
1. Load CLAUDE.md (always)
   ↓
2. Check for .agent/state.md
   - If exists: Load to resume work
   ↓
3. Check for .agent/project.md
   - If exists: Load for context
   ↓
4. During work: Auto-load language guide based on file extensions
   ↓
5. On-demand: Load workflows, patterns.md, memory/ as needed
```

---

## Language Guide Skills

Auto-loaded based on file extensions you're working with (21 languages):

| Language | Extensions | Guide |
|----------|------------|-------|
| TypeScript | `.ts`, `.tsx`, `.js`, `.jsx` | `skills/typescript-guide/SKILL.md` |
| Python | `.py` | `skills/python-guide/SKILL.md` |
| Go | `.go` | `skills/go-guide/SKILL.md` |
| Rust | `.rs` | `skills/rust-guide/SKILL.md` |
| Kotlin | `.kt`, `.kts` | `skills/kotlin-guide/SKILL.md` |
| Java | `.java` | `skills/java-guide/SKILL.md` |
| C# | `.cs` | `skills/csharp-guide/SKILL.md` |
| PHP | `.php` | `skills/php-guide/SKILL.md` |
| Swift | `.swift` | `skills/swift-guide/SKILL.md` |
| C/C++ | `.c`, `.cpp`, `.h`, `.hpp` | `skills/cpp-guide/SKILL.md` |
| Ruby | `.rb` | `skills/ruby-guide/SKILL.md` |
| *+ 10 more* | SQL, Shell, R, Dart, HTML/CSS, Lua, Assembly, CUDA, Solidity, Zig | |

**No manual selection needed** - AI detects automatically.

[:octicons-arrow-right-24: Language Guides](../languages/index.md)

---

## Framework Guides

On-demand framework-specific patterns (33 frameworks across 11 language families):

| Language | Frameworks |
|----------|------------|
| TypeScript/JS | React, Next.js, Express |
| Python | Django, FastAPI, Flask |
| Go | Gin, Echo, Fiber |
| Rust | Axum, Actix-web, Rocket |
| Kotlin | Spring Boot (Kotlin), Ktor, Android Compose |
| Java | Spring Boot, Quarkus, Micronaut |
| C# | ASP.NET Core, Blazor, Unity |
| PHP | Laravel, Symfony, WordPress |
| Swift | SwiftUI, UIKit, Vapor |
| Ruby | Rails, Sinatra, Hanami |
| Dart | Flutter, Shelf, Dart Frog |

**How to load**: Reference the framework or work in a project that uses it.

[:octicons-arrow-right-24: Framework Guides](../frameworks/index.md)

---

## Workflows

On-demand workflows for structured tasks:

| Workflow | When to Use |
|----------|-------------|
| `initialize-project.md` | New or existing project setup |
| `create-prd.md` | Plan complex features |
| `generate-tasks.md` | Break PRD into tasks |
| `troubleshooting.md` | Debug systematically |
| `generate-agents-md.md` | Cross-tool compatibility |

**How to invoke**:

```
@.agent/workflows/create-prd.md
```

[:octicons-arrow-right-24: Workflows](../workflows/index.md)

---

## Best Practices

### Do

- [x] Let files emerge naturally
- [x] Update project.md when stack changes
- [x] Document significant decisions in memory/
- [x] Keep state.md current during long work

### Don't

- [ ] Create project.md before making decisions
- [ ] Document every small choice
- [ ] Over-organize the directory
- [ ] Manually edit language guide skills (they're templates)

---

## Git Configuration

### What to Commit

```
✓ CLAUDE.md
✓ .agent/README.md
✓ .agent/skills/
✓ .agent/workflows/
✓ .agent/tasks/EXAMPLE-*.md (examples)
✓ .agent/project.md.template
✓ .agent/state.md.template
```

### What to Optionally Gitignore

```gitignore
# Generated files (optional)
.agent/project.md
.agent/patterns.md
.agent/state.md
.agent/tasks/*.md
!.agent/tasks/EXAMPLE-*.md
.agent/memory/*.md
!.agent/memory/.gitkeep
```

!!! note "Team Decision"

    Some teams commit project.md and patterns.md to share context. Others keep them local. Choose what works for your team.

---

## Related

<div class="grid cards" markdown>

-   :material-file-document:{ .lg .middle } **CLAUDE.md**

    ---

    The main instruction file.

    [:octicons-arrow-right-24: CLAUDE.md](claude-md.md)

-   :material-code-braces:{ .lg .middle } **Language Guides**

    ---

    Auto-loaded language rules (21 languages).

    [:octicons-arrow-right-24: Languages](../languages/index.md)

-   :material-layers:{ .lg .middle } **Framework Guides**

    ---

    Framework-specific patterns (33 frameworks).

    [:octicons-arrow-right-24: Frameworks](../frameworks/index.md)

-   :material-cog:{ .lg .middle } **Workflows**

    ---

    On-demand structured workflows.

    [:octicons-arrow-right-24: Workflows](../workflows/index.md)

</div>
