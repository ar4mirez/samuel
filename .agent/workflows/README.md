# Workflows Directory

This directory contains structured workflows for complex feature development using AI assistance.

## Purpose

Workflows provide step-by-step guidance for AI to tackle large features systematically, with built-in verification checkpoints and quality guardrails.

## Available Workflows

### 1. generate-agents-md.md (NEW)
**Use for**: Creating AGENTS.md for cross-tool compatibility

**When to use**:
- Team uses multiple AI tools (Cursor, Codex, Copilot + Claude)
- Open-source project needs universal AI instructions
- Symlinks not suitable for deployment

**How to use**:
```
Generate AGENTS.md for this project using @.agent/workflows/generate-agents-md.md
```

**Output**: `./AGENTS.md`

**Alternative**: Use symlink instead: `ln -s CLAUDE.md AGENTS.md`

---

### 2. create-prd.md
**Use for**: Defining complex features with Product Requirements Documents

**When to use**:
- Feature affects >10 files
- Building new subsystem
- Requirements unclear/ambiguous
- Feature takes >1 week
- Breaking changes to architecture

**How to use**:
```
@.agent/workflows/create-prd.md

I want to build: [describe feature]
Reference files: @file1.ts @file2.tsx
```

**Output**: `.agent/tasks/NNNN-prd-feature-name.md`

---

### 3. generate-tasks.md
**Use for**: Breaking PRDs into actionable task lists

**When to use**:
- After PRD created and reviewed
- Ready to plan implementation
- Need structured roadmap

**How to use**:
```
Generate tasks for @.agent/tasks/0001-prd-feature-name.md using @.agent/workflows/generate-tasks.md
```

**Output**: `.agent/tasks/tasks-NNNN-prd-feature-name.md`

---

## Workflow Process

### Full Feature Development Flow

```
1. Define Feature
   ↓
   Use create-prd.md
   ↓
2. PRD Created (.agent/tasks/NNNN-prd-feature.md)
   ↓
   Review & Approve
   ↓
3. Generate Tasks
   ↓
   Use generate-tasks.md
   ↓
4. Task List Created (.agent/tasks/tasks-NNNN-prd-feature.md)
   ↓
   Review High-Level Tasks → User says "Go"
   ↓
   Detailed Sub-Tasks Generated
   ↓
5. Implement Step-by-Step
   ↓
   "Start on task 1.1"
   AI implements → Tests → Validates guardrails → Commits
   ↓
   "Continue with 1.2"
   Repeat until all tasks complete
   ↓
6. Feature Complete
   ↓
   Update .agent/project.md
   Update .agent/patterns.md (if new patterns)
   Create .agent/memory/YYYY-MM-DD-feature.md (if significant)
```

---

## Integration with CLAUDE.md

### 4D Methodology Mapping

**ATOMIC Mode** (default):
- Skip workflows entirely
- Direct implementation for single-file changes

**FEATURE Mode** (>5 files):
- May use task generation for complex features
- PRD optional (can generate tasks directly from description)

**COMPLEX Mode** (architecture changes):
- **MUST use create-prd.md** for scope definition
- **MUST use generate-tasks.md** for task breakdown
- Full workflow with checkpoints

### Escalation Triggers

AI automatically suggests workflows when:
- User request mentions "new system", "subsystem", "major feature"
- Analysis shows >10 files affected
- Scope unclear after initial questions
- Breaking changes to existing architecture

**AI message**:
```
"This appears to be a COMPLEX mode task (affects 15+ files, new auth subsystem).

I recommend using the PRD workflow:
1. @.agent/workflows/create-prd.md to define requirements
2. @.agent/workflows/generate-tasks.md to break down implementation

This ensures nothing is missed and provides clear checkpoints.

Proceed with PRD workflow? Or do you prefer a different approach?"
```

---

## Benefits

### 1. Structured Development
- Clear scope definition (PRD)
- Systematic breakdown (task list)
- Verifiable progress (check off tasks)

### 2. Quality Assurance
- Guardrails validated at each step
- Tests required for each task
- Security/performance checkpoints built-in

### 3. Reduced Risk
- User approval at key milestones (PRD, high-level tasks, each subtask)
- Atomic commits (easy rollback)
- Clear dependencies (no surprises)

### 4. Better Documentation
- PRD captures "why" decisions
- Task list captures "what" was done
- Commits capture "how" implementation worked

### 5. Manageable Complexity
- Large features broken into small, reviewable changes
- AI stays focused (one task at a time)
- Clear progress tracking (X/Y tasks complete)

---

## File Organization

After using workflows, your `.agent/tasks/` directory might look like:

```
.agent/tasks/
├── 0001-prd-user-authentication.md       # PRD for auth feature
├── tasks-0001-prd-user-authentication.md # Task breakdown
├── 0002-prd-dashboard-analytics.md       # PRD for dashboard
├── tasks-0002-prd-dashboard-analytics.md # Task breakdown
└── 0003-prd-payment-integration.md       # PRD for payments
    (task list generated after review)
```

**Naming convention**:
- PRDs: `NNNN-prd-feature-name.md` (0001, 0002, etc.)
- Tasks: `tasks-NNNN-prd-feature-name.md` (matches PRD number)

---

## Customization

### Adapting create-prd.md

**Modify these sections based on project needs**:
- Clarifying questions (add domain-specific questions)
- PRD structure (add/remove sections for your industry)
- Guardrails affected (add project-specific rules)

**Example**: For e-commerce project, add:
- "Payment processing requirements"
- "Inventory management implications"
- "Tax calculation considerations"

### Adapting generate-tasks.md

**Modify these sections based on workflow**:
- Parent task structure (different phases for your stack)
- Testing requirements (your coverage targets)
- Deployment steps (your CI/CD process)

**Example**: For serverless project, replace:
- "Database migrations" → "DynamoDB table definitions"
- "API routes" → "Lambda function handlers"
- "Deployment" → "SAM template updates"

---

## Tips for Success

### DO:
✅ Use workflows for genuinely complex features (>10 files, unclear scope)
✅ Review and approve PRD before generating tasks
✅ Review high-level tasks before AI generates sub-tasks
✅ Implement one sub-task at a time (don't skip ahead)
✅ Validate guardrails after each task
✅ Update task list as you progress

### DON'T:
❌ Use workflows for simple features (overhead not worth it)
❌ Skip PRD review (garbage in, garbage out)
❌ Generate sub-tasks without user "Go" (might go wrong direction)
❌ Implement multiple tasks in one commit (defeats atomic approach)
❌ Skip testing to "move faster" (creates technical debt)

---

## Examples

### Example 1: User Authentication (COMPLEX)
```
User: "Add user authentication with email/password and OAuth"

AI: "This is a COMPLEX feature (10+ files, new subsystem, security-critical).
     I recommend PRD workflow. Proceed?"

User: "Yes"

AI: Uses create-prd.md
    Asks clarifying questions (OAuth providers? Session management? Password reset?)
    Creates .agent/tasks/0001-prd-user-authentication.md

User: Reviews PRD, approves

AI: Uses generate-tasks.md
    Generates 7 parent tasks
    User says "Go"
    Generates 42 sub-tasks total
    Saves .agent/tasks/tasks-0001-prd-user-authentication.md

User: "Start on task 1.1"

AI: Implements task 1.1 (Create user schema)
    Writes test
    Validates guardrails
    Commits
    "Task 1.1 complete. Review? Or continue with 1.2?"

[Repeat for all 42 tasks]
```

### Example 2: Bug Fix (ATOMIC)
```
User: "Fix the login button styling"

AI: "This is ATOMIC mode (single file, clear scope).
     No workflows needed.
     Implementing directly..."

[AI fixes bug, tests, commits]
```

### Example 3: Refactoring (FEATURE)
```
User: "Refactor the payment module to use Stripe instead of PayPal"

AI: "This is FEATURE mode (affects 6 files, known scope).
     PRD optional. Want task breakdown or proceed directly?"

User: "Task breakdown please"

AI: Uses generate-tasks.md (can work without PRD for clear refactorings)
    Generates 4 parent tasks, 18 sub-tasks
    Saves task list

[Proceed with implementation]
```

---

## Troubleshooting

### Issue: AI not following task order
**Solution**: Be explicit: "Start on task 1.1" (don't just say "start")

### Issue: Tasks too granular (hundreds of sub-tasks)
**Solution**: Review parent tasks before "Go". Ask AI to consolidate if too fine-grained.

### Issue: Tasks too coarse (each is days of work)
**Solution**: Ask AI to break down specific parent task: "Break down task 3.0 into more sub-tasks"

### Issue: Guardrails not being validated
**Solution**: Remind AI: "Validate all guardrails from task comment before marking complete"

### Issue: Tests skipped
**Solution**: Check task list includes test tasks. Add if missing: "Add test tasks for each implementation task"

---

## Maintenance

### When to Update Workflows

**Update create-prd.md when**:
- New guardrails added to CLAUDE.md
- Project adopts new tech stack
- Industry-specific requirements emerge
- PRD sections not useful (remove them)

**Update generate-tasks.md when**:
- Testing approach changes (new framework, different targets)
- Deployment process changes (new CI/CD pipeline)
- File organization patterns shift
- New complexity estimation approach

### Review Frequency
- **Monthly**: Check if workflows still match project needs
- **Per phase**: After major version/release, update examples
- **When stuck**: If workflows not helping, identify what's missing

---

## Integration with Other .agent/ Files

### Workflows reference:
- **`.agent/project.md`**: Tech stack, architecture patterns
- **`.agent/patterns.md`**: Existing coding patterns to follow

### Workflows update:
- **`.agent/project.md`**: After PRD for new subsystem
- **`.agent/patterns.md`**: If new patterns emerge during implementation
- **`.agent/state.md`**: Track progress through task list
- **`.agent/memory/`**: Document key decisions during complex features

---

**Remember**: Workflows are tools, not mandates. Use them when they add value (complex features), skip them when they don't (simple changes). The goal is structured development, not bureaucracy.
