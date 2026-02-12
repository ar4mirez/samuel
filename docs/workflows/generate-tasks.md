---
title: Generate Tasks
description: Break PRDs into actionable task lists
---

# Generate Tasks Workflow

Convert a Product Requirements Document (PRD) into an actionable, numbered task list.

---

## When to Use

After creating a PRD with the create-prd workflow:

```
@.claude/skills/generate-tasks/SKILL.md

Use the PRD at .claude/tasks/0001-prd-user-auth.md
```

---

## What AI Does

### Step 1: Analyze PRD

AI reads the PRD and identifies:

- Functional requirements
- Technical considerations
- Dependencies between components
- Testing requirements

### Step 2: Generate Task Breakdown

Creates a hierarchical task list:

- **Phases** - Major milestones
- **Tasks** - Individual work items
- **Subtasks** - Granular steps

### Step 3: Add Verification

Each task includes:

- Acceptance criteria
- Dependencies
- Verification steps

### Step 4: Save Task List

Saves to: `.claude/tasks/tasks-NNNN-prd-feature-name.md`

---

## Task Structure

### Format

```markdown
# Task List: User Authentication

**PRD**: 0001-prd-user-auth.md
**Created**: 2025-01-15
**Status**: Not Started

---

## Phase 1: Database & Models

### Task 1.1: Create User Schema
**Complexity**: Low
**Dependencies**: None
**Files**: src/db/schema/user.ts

**Steps**:
1. Create user table migration
2. Add email, password_hash, created_at fields
3. Add unique constraint on email
4. Run migration

**Verification**:
- [ ] Migration runs successfully
- [ ] Can create user in database
- [ ] Email uniqueness enforced

---

### Task 1.2: Create User Model
**Complexity**: Low
**Dependencies**: 1.1
**Files**: src/models/user.ts

**Steps**:
1. Create User interface
2. Create UserRepository class
3. Implement create, findByEmail, findById methods

**Verification**:
- [ ] Unit tests pass
- [ ] Can create and retrieve users
```

---

## Task Properties

### Complexity Levels

| Level | Description | Estimated Effort |
|-------|-------------|------------------|
| **Low** | Single file, clear implementation | ~30 minutes |
| **Medium** | Multiple files, some decisions | ~2 hours |
| **High** | Complex logic, integration | ~4+ hours |

### Dependencies

Tasks can depend on other tasks:

```markdown
**Dependencies**: 1.1, 1.2
```

AI ensures tasks are ordered correctly so dependencies complete first.

### Verification Steps

Every task has verification criteria:

```markdown
**Verification**:
- [ ] Unit tests pass
- [ ] Integration test passes
- [ ] Manual testing confirms behavior
```

---

## Example Task List

### User Authentication Feature

```markdown
# Task List: User Authentication

## Phase 1: Database & Models (Tasks 1.1-1.4)

### Task 1.1: Create Database Schema
- Create users table
- Create sessions table
- Verification: Migrations run

### Task 1.2: Create User Model
- Dependencies: 1.1
- Verification: Unit tests pass

### Task 1.3: Create Session Model
- Dependencies: 1.1
- Verification: Unit tests pass

### Task 1.4: Add Password Hashing
- Dependencies: 1.2
- Verification: Passwords hashed correctly

---

## Phase 2: API Endpoints (Tasks 2.1-2.5)

### Task 2.1: Registration Endpoint
- POST /api/auth/register
- Dependencies: 1.2, 1.4
- Verification: Can register new user

### Task 2.2: Login Endpoint
- POST /api/auth/login
- Dependencies: 1.2, 1.3, 1.4
- Verification: Can login and receive token

### Task 2.3: JWT Middleware
- Validate tokens on protected routes
- Dependencies: 1.3
- Verification: Protected routes require auth

### Task 2.4: Logout Endpoint
- POST /api/auth/logout
- Dependencies: 1.3
- Verification: Sessions invalidated

### Task 2.5: Current User Endpoint
- GET /api/auth/me
- Dependencies: 2.3
- Verification: Returns current user data

---

## Phase 3: OAuth Integration (Tasks 3.1-3.4)

### Task 3.1: Google OAuth Setup
- Configure Google OAuth
- Dependencies: 2.1
- Verification: Can initiate Google auth

...
```

---

## Working Through Tasks

### Starting a Task

```
Start task 1.1 from the user auth task list
```

AI will:

1. Read the task details
2. Check dependencies are complete
3. Implement the task
4. Run verification steps
5. Mark task complete

### Continuing to Next Task

```
Continue with task 1.2
```

AI will:

1. Verify 1.1 is complete
2. Implement 1.2
3. Run verification
4. Update task status

### Checking Progress

```
What's the status of the user auth tasks?
```

AI will:

1. Read the task list
2. Report completed/in-progress/pending
3. Suggest next task

---

## Task List Management

### Updating Status

Task statuses:

- **Not Started** - Work not begun
- **In Progress** - Currently working
- **Blocked** - Waiting on something
- **Complete** - Verified and done

### Adding Tasks

If new requirements emerge:

```
Add a new task to Phase 2 for email verification
```

AI will:

1. Add the task with proper numbering
2. Update dependencies if needed
3. Add verification steps

### Removing Tasks

If requirements change:

```
Remove task 3.4 - we're not doing social login
```

AI will:

1. Remove the task
2. Update any dependent tasks
3. Note the change

---

## Tips for Effective Tasks

### Keep Tasks Small

```markdown
# ❌ Too big
Task 1: Build entire authentication system

# ✅ Right size
Task 1.1: Create user database schema
Task 1.2: Create user model
Task 1.3: Implement password hashing
Task 1.4: Create registration endpoint
```

### Clear Verification

```markdown
# ❌ Vague
Verification: Works correctly

# ✅ Specific
Verification:
- [ ] POST /api/auth/register returns 201
- [ ] User appears in database
- [ ] Password is hashed (not plaintext)
- [ ] Duplicate email returns 409
```

### Explicit Dependencies

```markdown
# ❌ Implicit
Task 2.1: Login endpoint (assumes registration exists)

# ✅ Explicit
Task 2.1: Login endpoint
Dependencies: 1.2 (User model), 1.4 (Password hashing)
```

---

## Related

<div class="grid cards" markdown>

-   :material-file-document:{ .lg .middle } **Create PRD**

    ---

    Create the PRD before generating tasks.

    [:octicons-arrow-right-24: Create PRD](create-prd.md)

-   :material-cog:{ .lg .middle } **4D Methodology**

    ---

    How tasks fit into the development process.

    [:octicons-arrow-right-24: Methodology](../core/methodology.md)

</div>
