---
title: Initialize Project
description: Set up new projects or analyze existing codebases
---

# Initialize Project Workflow

Set up a new project or analyze an existing codebase to create project documentation.

---

## When to Use

- **New Project**: Starting from scratch
- **Existing Project**: Onboarding to a codebase you haven't documented yet
- **Team Alignment**: Creating shared project documentation

---

## How to Invoke

### New Project

```
@.claude/skills/initialize-project/SKILL.md

Initialize a new TypeScript API with Express, PostgreSQL, and Jest
```

### Existing Project

```
@.claude/skills/initialize-project/SKILL.md

This is an existing project - analyze the codebase and document patterns
```

---

## What AI Does

### For New Projects

1. **Asks Clarifying Questions**:
   - Tech stack (language, framework, versions)?
   - Architecture (monolith, microservices, serverless)?
   - Database (PostgreSQL, MongoDB, none)?
   - Testing approach (unit, integration, e2e)?
   - Deployment target (AWS, Vercel, Docker)?

2. **Creates Structure**:
   - Directory layout (best practices for stack)
   - Configuration files (tsconfig, package.json, etc.)
   - `.gitignore`, `.env.example`
   - Basic README

3. **Documents in .claude/**:
   - Creates `.claude/project.md` with decisions
   - Notes any patterns established

### For Existing Projects

1. **Scans Codebase**:
   - Package files (package.json, requirements.txt, go.mod)
   - Configuration files
   - Directory structure
   - README and documentation

2. **Analyzes Patterns**:
   - Code conventions
   - Testing patterns
   - API patterns
   - Error handling patterns

3. **Reviews History**:
   - Recent commits
   - Commit message format
   - Active contributors

4. **Documents Findings**:
   - Creates `.claude/project.md` with tech stack
   - Creates `.claude/patterns.md` with conventions
   - Identifies gaps and suggests improvements

---

## Questions AI Will Ask

### Tech Stack

- What language/framework are you using?
- Which version?
- Any specific libraries required?

### Architecture

- Monolith or microservices?
- Frontend/backend separation?
- API style (REST, GraphQL)?

### Database

- Which database?
- ORM preference?
- Migration strategy?

### Testing

- Testing framework preference?
- Coverage targets?
- E2E testing needed?

### Deployment

- Where will this deploy?
- CI/CD requirements?
- Containerization needed?

---

## Output: project.md

Example output for a TypeScript project:

```markdown
# Project Configuration

## Tech Stack
- **Language**: TypeScript 5.3
- **Runtime**: Node.js 20 LTS
- **Framework**: Express.js 4.18
- **Database**: PostgreSQL 16
- **ORM**: Prisma 5.7
- **Testing**: Vitest + Supertest
- **Linting**: ESLint + Prettier

## Architecture
- Monolithic API
- Repository pattern for data access
- Middleware-based request handling
- JWT authentication

## Directory Structure
```
src/
├── controllers/     # Route handlers
├── services/        # Business logic
├── repositories/    # Data access
├── middleware/      # Express middleware
├── types/           # TypeScript types
├── utils/           # Shared utilities
└── index.ts         # Entry point
```

## Conventions
- ESM modules (import/export)
- Zod for request validation
- Conventional commits
- Tests alongside source files

## External Services
- None configured yet

## Commands
- `npm run dev` - Development server
- `npm test` - Run tests
- `npm run build` - Production build
- `npm run lint` - Lint check
```

---

## Output: patterns.md

Example output when patterns are found:

```markdown
# Coding Patterns

## API Endpoints

Routes follow this pattern:

```typescript
// POST /api/users
router.post('/',
  validateBody(CreateUserSchema),
  async (req, res) => {
    const user = await userService.create(req.validated);
    res.status(201).json(user);
  }
);
```

## Error Handling

Centralized error handler in middleware:

```typescript
// Throw AppError in services
throw new AppError('User not found', 404);

// Caught by error middleware
app.use(errorHandler);
```

## Database Queries

Always use transactions for writes:

```typescript
await prisma.$transaction(async (tx) => {
  await tx.user.create({ data: userData });
  await tx.audit.create({ data: auditData });
});
```
```

---

## After Initialization

Once the workflow completes:

1. **Review the generated files** - Confirm accuracy
2. **Start coding** - Guardrails now apply
3. **Let documentation grow** - Add patterns as they emerge

!!! tip "Don't Over-Document"

    The initialization creates a foundation. Let `patterns.md` grow naturally as you work.

---

## Example Sessions

### New TypeScript API

```
User: @.claude/skills/initialize-project/SKILL.md

Initialize a new REST API for a task management app

AI: I'll help you set up a new project. Let me ask a few questions:

1. **Language**: TypeScript (confirmed from your request)
2. **Framework preference**:
   a) Express.js (lightweight, flexible)
   b) Fastify (performance-focused)
   c) NestJS (full-featured, Angular-style)

3. **Database**:
   a) PostgreSQL
   b) MongoDB
   c) SQLite (for simplicity)
   d) None (in-memory for now)

User: Express and PostgreSQL

AI: Great choices! A few more questions:

4. **ORM**:
   a) Prisma (type-safe, modern)
   b) TypeORM
   c) Raw SQL with pg

5. **Testing**:
   a) Vitest (fast, modern)
   b) Jest (established)

User: Prisma and Vitest

AI: Perfect. I'll create:
- Directory structure
- package.json with dependencies
- tsconfig.json (strict mode)
- Prisma schema
- .env.example
- .gitignore
- .claude/project.md

[Creates files...]

Done! Your project is ready. Run:
npm install
npm run dev

I've documented everything in .claude/project.md
```

### Existing Python Project

```
User: @.claude/skills/initialize-project/SKILL.md

This is an existing Django project - analyze it

AI: I'll analyze the codebase. Let me scan...

**Tech Stack Detected**:
- Python 3.11
- Django 4.2
- PostgreSQL (from settings.py)
- pytest + pytest-django

**Directory Structure**:
- apps/ (Django apps)
- config/ (settings)
- templates/
- static/

**Patterns Found**:
- Class-based views
- DRF serializers for API
- Custom user model
- Celery for background tasks

**Conventions**:
- Black formatting (pyproject.toml)
- isort for imports
- Type hints (some files)

**Observations**:
- Good test coverage in apps/users/
- Missing tests in apps/orders/
- No type hints in older code

I'll create .claude/project.md and .claude/patterns.md with these findings.

Would you like me to also create recommendations for improvements?
```

---

## Related

<div class="grid cards" markdown>

-   :material-file-document:{ .lg .middle } **Create PRD**

    ---

    After initialization, plan your first feature.

    [:octicons-arrow-right-24: Create PRD](create-prd.md)

-   :material-shield:{ .lg .middle } **Guardrails**

    ---

    Rules that apply to your new project.

    [:octicons-arrow-right-24: Guardrails](../core/guardrails.md)

</div>
