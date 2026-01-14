# Workflow: Generate Task List from PRD

> **Purpose**: Break down Product Requirements Document into actionable, verifiable tasks
>
> **Use When**: PRD created, ready to plan implementation
> **4D Phase**: Develop (planning implementation steps)

---

## Goal

Guide AI in creating detailed, step-by-step task list in Markdown format based on existing PRD. The task list should guide a developer through implementation with built-in quality checkpoints.

## Prerequisites

✅ **Before using this workflow:**
- PRD exists in `.agent/tasks/NNNN-prd-feature-name.md`
- PRD reviewed and approved by user
- `.agent/project.md` exists (or will be created during this process)
- Current codebase reviewed for context

## Process

### Phase 1: Analysis & High-Level Planning

#### Step 1: Receive PRD Reference
User points AI to specific PRD file.

**Example**:
```
"Generate tasks for @.agent/tasks/0001-prd-user-authentication.md using @.agent/workflows/generate-tasks.md"
```

#### Step 2: Read PRD
AI reads and analyzes:
- Functional requirements
- User stories
- Technical considerations
- Guardrails affected
- Success metrics

#### Step 3: Assess Current State
**CRITICAL**: Review existing codebase before planning.

**AI must:**

1. **Check `.agent/project.md`** for:
   - Tech stack and framework choices
   - Existing architecture patterns
   - Testing framework and conventions
   - Deployment configuration

2. **Scan codebase** for:
   - Existing directory structure
   - Similar features (auth, API routes, etc.)
   - Naming conventions
   - Test file patterns
   - Component organization

3. **Identify reusable components**:
   - Existing utilities that can be leveraged
   - UI components from design system
   - Database connection patterns
   - Error handling patterns
   - Validation utilities

4. **Check guardrails compliance**:
   - Current file lengths (don't add to files approaching limits)
   - Existing test coverage patterns
   - Security patterns in use
   - Performance considerations

**Output**: Brief summary for user review before generating tasks.

**Example**:
```markdown
## Current State Assessment

**Existing Infrastructure:**
- Express.js backend with TypeScript
- Prisma ORM with PostgreSQL
- Jest for testing (78% coverage currently)
- React frontend with Tailwind CSS

**Reusable Components:**
- `lib/validation.ts` - Input validation utilities (use Zod)
- `components/Form.tsx` - Form component (reuse for auth forms)
- `lib/errors.ts` - Centralized error handling (extend for auth errors)

**Patterns to Follow:**
- API routes in `src/api/<feature>/` structure
- Tests alongside code (`*.test.ts`)
- Repository pattern for data access
- Environment variables via `.env` with fallbacks

**Guardrails Note:**
- `src/api/user/user.controller.ts` is 285 lines (near 300 limit - don't modify)
- Create new auth module instead of extending user module
```

#### Step 4: Generate Parent Tasks
Based on PRD and current state, create high-level tasks.

**Guidelines:**
- Typically 4-7 parent tasks
- Each parent = logical phase of implementation
- Order by dependency (database → backend → frontend → tests)
- Include setup, implementation, testing, documentation

**Save to**: `.agent/tasks/tasks-NNNN-prd-feature-name.md`

**Present to user** in this format:

```markdown
## High-Level Tasks

- [ ] 1.0 Database Schema & Migrations
- [ ] 2.0 Backend Authentication Service
- [ ] 3.0 API Routes & Middleware
- [ ] 4.0 Frontend Auth Components
- [ ] 5.0 OAuth Integration
- [ ] 6.0 Testing & Validation
- [ ] 7.0 Documentation & Deployment

**Relevant Files** (identified):
- `src/db/schemas/user.schema.ts` - Create
- `src/auth/auth.service.ts` - Create
- `src/api/auth/auth.controller.ts` - Create
... (full list in task file)
```

**Then inform user**:
```
"High-level tasks generated. Review the task structure above.
Ready to generate detailed sub-tasks? Respond with 'Go' to proceed."
```

#### Step 5: Wait for User Confirmation
**PAUSE** - Do not proceed until user responds with "Go" or equivalent.

This checkpoint ensures high-level plan aligns with expectations before diving into details.

---

### Phase 2: Detailed Task Breakdown

#### Step 6: Generate Sub-Tasks
Once user confirms "Go", break down each parent task into actionable sub-tasks.

**Sub-task Guidelines:**

1. **Granularity**: Each sub-task = 30-60 min of work (3,000-6,000 tokens)
2. **Atomicity**: Complete in one session, one commit
3. **Testability**: Each sub-task includes verification step
4. **Dependency**: Clear which sub-tasks depend on others
5. **Guardrails**: Each task validates against relevant guardrails

**Sub-task Naming Convention:**
- Use present tense verbs: "Create", "Implement", "Add", "Update", "Test"
- Be specific: "Create user schema" not "Setup database"
- Include success criterion: "Create user schema (with email, password_hash, created_at fields)"

**Example Parent → Sub-tasks**:
```markdown
- [ ] 1.0 Database Schema & Migrations
  - [ ] 1.1 Create user schema with Prisma (email, password_hash, oauth_provider, created_at, updated_at)
  - [ ] 1.2 Create session schema with Prisma (user_id FK, token, expires_at, created_at)
  - [ ] 1.3 Generate and run Prisma migration for user and session tables
  - [ ] 1.4 Verify schema in PostgreSQL (check indexes, constraints)
  - [ ] 1.5 Update .env.example with DATABASE_URL placeholder
```

#### Step 7: Add Guardrail Validation Per Task
For each sub-task, identify which guardrails must be validated.

**Format**: Add guardrail checklist as comment under relevant tasks.

**Example**:
```markdown
- [ ] 2.3 Implement password hashing service (bcrypt, cost factor 12)
  <!-- Guardrails:
    ✓ All environment variables have secure defaults (BCRYPT_ROUNDS=12)
    ✓ All exported functions have type signatures and JSDoc
    ✓ Function ≤50 lines (hash, verify, validate strength as separate functions)
    ✓ Edge cases tested (null, empty, too short, invalid chars)
  -->
```

#### Step 8: Add Complexity Estimates
For each sub-task, estimate complexity in tokens (not time).

**Complexity Levels:**
- **Simple**: <2,000 tokens (single function, straightforward logic)
- **Medium**: 2,000-5,000 tokens (multiple functions, some complexity)
- **Complex**: 5,000-10,000 tokens (multiple files, integration, edge cases)

**Example**:
```markdown
- [ ] 2.3 Implement password hashing service [~3,000 tokens - Medium]
```

#### Step 9: Identify Relevant Files
List all files that will be created or modified, including tests.

**Format**:
```markdown
## Relevant Files

### Create (New Files)
- `src/db/schemas/user.schema.ts` - User table schema with Prisma
- `src/db/schemas/user.schema.test.ts` - Schema validation tests
- `src/auth/auth.service.ts` - Core authentication business logic
- `src/auth/auth.service.test.ts` - Service unit tests (>95% coverage target)
- `src/auth/auth.middleware.ts` - JWT validation middleware
- `src/auth/auth.middleware.test.ts` - Middleware tests
- `src/auth/auth.types.ts` - TypeScript types and interfaces
- `src/api/auth/auth.controller.ts` - Route handlers (login, register, logout)
- `src/api/auth/auth.controller.test.ts` - Controller integration tests
- `src/api/auth/auth.routes.ts` - Express route definitions

### Modify (Existing Files)
- `src/app.ts` - Add auth routes to Express app
- `src/middleware/index.ts` - Export auth middleware
- `.env.example` - Add JWT_SECRET, BCRYPT_ROUNDS placeholders
- `package.json` - Add dependencies (bcrypt, jsonwebtoken, passport)

### Reference (Read but don't modify)
- `lib/validation.ts` - Use existing Zod validators
- `lib/errors.ts` - Use existing error classes
- `.agent/project.md` - Check for architecture patterns
```

#### Step 10: Add Implementation Notes
Provide guidance for each task section.

**Include:**
- Testing requirements specific to this feature
- Performance considerations
- Security reminders
- Integration points
- Common pitfalls to avoid

**Example**:
```markdown
## Implementation Notes

### Testing Requirements
- All auth logic must have >95% coverage (business-critical)
- Test files alongside code (e.g., `auth.service.ts` + `auth.service.test.ts`)
- Use `npm test src/auth` to run auth module tests only
- Integration tests must test actual JWT validation, not mocks

### Security Checklist (verify before each commit)
- [ ] All passwords hashed (never store plain text)
- [ ] All JWT secrets from environment variables
- [ ] All inputs validated with Zod schemas
- [ ] All SQL queries parameterized (use Prisma)
- [ ] Rate limiting on login/register endpoints
- [ ] Authentication events logged

### Performance Targets
- Login/register: <200ms (p95)
- JWT validation: <50ms (p95)
- Password hashing: 100-300ms (bcrypt cost 12)

### Common Pitfalls
- ❌ Don't validate JWT in database query (use middleware)
- ❌ Don't store tokens in localStorage (use httpOnly cookies)
- ❌ Don't return sensitive data in error messages
- ❌ Don't skip rate limiting (prevents brute force)
```

#### Step 11: Generate Final Output
Combine all sections into final task list structure.

---

## Output Format

The generated task list MUST follow this structure:

```markdown
# Task List: [Feature Name]

> **Source PRD**: `.agent/tasks/NNNN-prd-feature-name.md`
> **Generated**: YYYY-MM-DD
> **Status**: Not Started / In Progress / Completed

---

## Current State Assessment

[Summary from Step 3]

---

## Relevant Files

### Create (New Files)
- `path/to/file.ts` - Description [Complexity: Simple/Medium/Complex]
- `path/to/file.test.ts` - Unit tests for file.ts [Complexity: Medium]

### Modify (Existing Files)
- `path/to/existing.ts` - Changes needed [Complexity: Simple]

### Reference (Read Only)
- `path/to/reference.ts` - Existing pattern to follow

---

## Implementation Notes

### Testing Requirements
[From Step 10]

### Security Checklist
[From Step 10]

### Performance Targets
[From Step 10]

### Common Pitfalls
[From Step 10]

---

## Tasks

- [ ] 1.0 Parent Task Title
  - [ ] 1.1 Sub-task description [~2,000 tokens - Simple]
    <!-- Guardrails: ✓ Specific guardrails for this task -->
  - [ ] 1.2 Sub-task description [~4,000 tokens - Medium]
    <!-- Guardrails: ✓ Specific guardrails for this task -->

- [ ] 2.0 Parent Task Title
  - [ ] 2.1 Sub-task description [~3,000 tokens - Medium]
  - [ ] 2.2 Sub-task description [~5,000 tokens - Complex]

- [ ] 3.0 Testing & Validation
  - [ ] 3.1 Verify all unit tests passing
  - [ ] 3.2 Verify coverage >95% for auth module
  - [ ] 3.3 Verify all guardrails validated
  - [ ] 3.4 Run security audit (npm audit)
  - [ ] 3.5 Performance testing (load testing if applicable)

- [ ] 4.0 Documentation & Deployment
  - [ ] 4.1 Update API documentation
  - [ ] 4.2 Update .agent/project.md with auth architecture
  - [ ] 4.3 Add to .agent/patterns.md if new patterns emerged
  - [ ] 4.4 Create deployment checklist
  - [ ] 4.5 Update README with setup instructions

---

## Progress Tracking

**Total Tasks**: X parent, Y sub-tasks
**Completed**: 0/Y (0%)
**In Progress**: None
**Blocked**: None

**Last Updated**: YYYY-MM-DD HH:MM

---

## Success Criteria

Before marking this task list complete, verify:
- [ ] All functional requirements from PRD implemented
- [ ] All guardrails validated (especially security & testing)
- [ ] Test coverage meets targets (>95% for auth)
- [ ] Performance targets met (<200ms auth responses)
- [ ] Documentation updated (.agent/ files, API docs, README)
- [ ] No security vulnerabilities (npm audit clean)
- [ ] User acceptance testing passed (if applicable)

---

## Next Steps

1. Review this task list with user
2. User confirms or requests changes
3. Start with task 1.1 (implement atomically, one commit per task)
4. After each task: Run tests, validate guardrails, commit
5. Update progress tracking after each task
6. Mark parent task complete only when all sub-tasks done
```

---

## Interaction Model

### Two-Phase Approach

**Phase 1**: Generate parent tasks → Present to user → Wait for "Go"
**Phase 2**: Generate detailed sub-tasks → Save file → Inform user

This ensures high-level plan approved before diving into details.

### During Implementation

User will say: "Start on task 1.1" or "Continue with next task"

AI should:
1. Implement the specific sub-task
2. Write tests alongside code
3. Validate against guardrails in task comment
4. Run tests
5. Verify success criteria
6. Commit with message: `feat(auth): task 1.1 - create user schema`
7. Mark task complete in task list file
8. Ask: "Task 1.1 complete. Review changes? Or continue with 1.2?"

---

## Guardrails Integration

For each sub-task, identify which of CLAUDE.md's 30+ guardrails apply.

**Common guardrails per task type:**

### Database Schema Tasks
- ✓ All environment variables have secure defaults
- ✓ Parameterized queries (Prisma enforces this)

### Service/Logic Tasks
- ✓ No function exceeds 50 lines
- ✓ Cyclomatic complexity ≤10 per function
- ✓ All exported functions have type signatures and JSDoc
- ✓ All user inputs validated before processing

### API/Route Tasks
- ✓ All API boundaries have input validation (Zod schemas)
- ✓ All async operations have timeout/cancellation
- ✓ API responses <200ms for simple queries

### Test Tasks
- ✓ Coverage targets: >95% for auth (business-critical)
- ✓ All public APIs have unit tests
- ✓ Edge cases tested (null, empty, boundary values)
- ✓ No test interdependencies

### Security Tasks (Auth-related)
- ✓ All user inputs validated before processing
- ✓ All database queries parameterized
- ✓ All environment variables have secure defaults
- ✓ Dependencies checked for vulnerabilities

---

## Output Location

**File**: `.agent/tasks/tasks-NNNN-prd-feature-name.md`
**Format**: Markdown
**Naming**: Match PRD filename (e.g., PRD is `0001-prd-user-auth.md` → tasks are `tasks-0001-prd-user-auth.md`)

---

## Final Instructions

### For AI Assistant:

**Phase 1 (High-Level)**:
1. ✅ Read PRD thoroughly
2. ✅ Assess current codebase (scan for existing patterns)
3. ✅ Generate 4-7 parent tasks
4. ✅ Identify all relevant files
5. ✅ Present to user
6. ❌ **WAIT for "Go" before Phase 2**

**Phase 2 (Detailed)**:
1. ✅ Break each parent into 3-8 sub-tasks
2. ✅ Add guardrail validation per task
3. ✅ Add complexity estimates
4. ✅ Add implementation notes
5. ✅ Save task list file
6. ✅ Inform user: "Task list ready at [path]. Ready to start task 1.1?"

**During Implementation**:
1. ✅ Implement one sub-task at a time
2. ✅ Validate guardrails for that specific task
3. ✅ Write tests alongside code
4. ✅ Commit after each task (atomic commits)
5. ✅ Update task list (mark completed)
6. ✅ Ask user to review before continuing

### Target Audience

**Primary**: Junior developer with AI assistance
**Assumption**: Developer knows fundamentals but needs clear guidance on:
- Where to create files
- What patterns to follow
- Which guardrails to validate
- How to test properly

---

**Remember**: Good task breakdown = clear roadmap. Each sub-task should be:
- **Atomic**: Complete in one session
- **Testable**: Has clear success criterion
- **Guarded**: Validates specific guardrails
- **Estimated**: Complexity known upfront
- **Dependent**: Dependencies clear (order matters)
