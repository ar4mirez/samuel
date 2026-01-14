# Workflows Directory

This directory contains structured workflows for AI-assisted development, covering the entire software development lifecycle.

## Purpose

Workflows provide step-by-step guidance for AI to tackle complex tasks systematically, with built-in verification checkpoints and quality guardrails.

## Workflow Categories

| Category | Workflows | Purpose |
|----------|-----------|---------|
| **Planning** | create-prd, generate-tasks, initialize-project | Define and break down work |
| **Quality** | code-review, security-audit, testing-strategy | Validate and improve code |
| **Maintenance** | cleanup-project, refactoring, dependency-update, update-framework | Keep codebase healthy |
| **Utility** | troubleshooting, generate-agents-md, document-work | Support and compatibility |

---

## Available Workflows

### Planning Workflows

#### initialize-project.md
**Use for**: Setting up new or existing projects with AI-assisted development

**When to use**:
- Starting a new project
- Adding AI assistance to existing project
- Need to analyze and document project structure

**How to use**:
```
@.agent/workflows/initialize-project.md
```

**Output**: `.agent/project.md`, analysis of project structure

---

#### create-prd.md
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

#### generate-tasks.md
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

### Quality Workflows

#### code-review.md
**Use for**: Systematic validation against all guardrails before committing

**When to use**:
- Before any commit (automated mode)
- During pull request review (interactive mode)
- After FEATURE/COMPLEX mode completion

**How to use**:
```
@.agent/workflows/code-review.md

Review the changes in src/api/
```

**Output**: Code review report with pass/fail/warning status

---

#### security-audit.md
**Use for**: Proactive security assessment covering OWASP Top 10

**When to use**:
- Before production deployments
- Monthly security review
- After adding authentication/authorization
- When adding external integrations

**How to use**:
```
@.agent/workflows/security-audit.md

Audit the authentication module
```

**Output**: Security audit report with findings and remediation steps

---

#### testing-strategy.md
**Use for**: Planning and achieving test coverage targets

**When to use**:
- After COMPLEX features
- When coverage drops below thresholds
- During test debt sprints
- Establishing testing foundation

**How to use**:
```
@.agent/workflows/testing-strategy.md

Plan testing for the payment module
```

**Output**: Test strategy with prioritized test backlog

---

### Maintenance Workflows

#### cleanup-project.md
**Use for**: Pruning unused guides and reducing .agent/ bloat

**When to use**:
- After initialize-project.md completes
- Quarterly maintenance
- Before major releases
- When .agent/ exceeds 1MB

**How to use**:
```
@.agent/workflows/cleanup-project.md

Clean up this project's .agent/ directory
```

**Output**: Leaner .agent/ directory, active-guides.json manifest

---

#### refactoring.md
**Use for**: Structured approach to technical debt remediation

**When to use**:
- Functions >50 lines, files >300 lines
- Cyclomatic complexity >10
- Code duplication in 3+ places
- Quarterly debt review

**How to use**:
```
@.agent/workflows/refactoring.md

Refactor the order processing module
```

**Output**: Refactoring plan with incremental steps

---

#### dependency-update.md
**Use for**: Safe and systematic dependency updates

**When to use**:
- Security vulnerability in dependency
- Monthly maintenance
- Major version available
- Before releases

**How to use**:
```
@.agent/workflows/dependency-update.md

Update project dependencies
```

**Output**: Update plan with testing and rollback procedures

---

#### update-framework.md
**Use for**: Updating AICoF to latest version while preserving customizations

**When to use**:
- New version released
- Want new language/framework guides
- Monthly maintenance check
- Team version sync

**How to use**:
```
@.agent/workflows/update-framework.md

Update to the latest version of AICoF
```

**Output**: Updated CLAUDE.md and .agent/ files with customizations preserved

---

### Utility Workflows

#### troubleshooting.md
**Use for**: Systematic debugging when stuck

**When to use**:
- Stuck for >30 minutes
- Tests failing unexpectedly
- Unclear error messages
- Need structured debugging

**How to use**:
```
@.agent/workflows/troubleshooting.md

Debug: API returns 500 error
```

**Output**: Diagnostic analysis and solution recommendations

---

#### generate-agents-md.md
**Use for**: Creating AGENTS.md for cross-tool compatibility

**When to use**:
- Team uses multiple AI tools (Cursor, Codex, Copilot + Claude)
- Open-source project needs universal AI instructions
- Symlinks not suitable for deployment

**How to use**:
```
@.agent/workflows/generate-agents-md.md
```

**Output**: `./AGENTS.md`

**Alternative**: Use symlink instead: `ln -s CLAUDE.md AGENTS.md`

---

#### document-work.md
**Use for**: Capturing patterns, decisions, and learnings from recent development work

**When to use**:
- After completing significant features
- End of work session
- New pattern identified
- Before vacation/handoff

**How to use**:
```
@.agent/workflows/document-work.md

Document the work from today's session on the authentication feature.
We made decisions about JWT vs sessions and established a new error handling pattern.
```

**Output**: Updated patterns.md, new memory files, state.md updates

---

## Workflow Map

```
                    ┌─────────────────────────────────────────┐
                    │           SDLC Stage Map                │
                    └─────────────────────────────────────────┘

    Planning              Implementation           Validation            Maintenance
    ────────              ──────────────           ──────────            ───────────

┌──────────────┐                              ┌──────────────┐    ┌──────────────┐
│ initialize-  │─────────────────────────────▶│ cleanup-     │    │ dependency-  │
│ project      │                              │ project      │    │ update       │
└──────────────┘                              └──────────────┘    └──────────────┘
       │                                             │                   │
       ▼                                             ▼                   │
┌──────────────┐      ┌──────────────┐      ┌──────────────┐            │
│ create-prd   │─────▶│ generate-    │─────▶│ code-review  │◀───────────┘
│              │      │ tasks        │      │              │
└──────────────┘      └──────────────┘      └──────────────┘
                                                   │
                                                   ▼
                                            ┌──────────────┐
                                            │ security-    │
                                            │ audit        │
                                            └──────────────┘
                                                   │
                                                   ▼
                      ┌──────────────┐      ┌──────────────┐
                      │ refactoring  │◀─────│ testing-     │
                      │              │      │ strategy     │
                      └──────────────┘      └──────────────┘
                             │
                             ▼
┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│ update-      │      │ trouble-     │      │ document-    │
│ framework    │      │ shooting     │      │ work         │
└──────────────┘      └──────────────┘      └──────────────┘
 (maintenance)         (any stage)           (any stage)
```

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

## Detailed Workflow Examples

This section provides comprehensive examples for each workflow, showing real-world usage scenarios.

### Planning Workflows

#### initialize-project.md Example

**Scenario**: Setting up a new TypeScript/React project

**User Request:**
```
@.agent/workflows/initialize-project.md

I'm starting a new e-commerce dashboard project
```

**AI Conversation:**
```
AI: "I'll help you set up this project. Let me ask some discovery questions:

1. What tech stack?
   - Language: TypeScript
   - Framework: React with Next.js
   - Styling: Tailwind CSS

2. What architecture?
   - Jamstack with API routes

3. Testing approach?
   - Jest + React Testing Library
   - Playwright for E2E
   - Coverage target: 80%

4. Deployment target?
   - Vercel

5. Database?
   - PostgreSQL with Prisma ORM

6. Additional requirements?
   - Authentication: NextAuth.js
   - API: REST with tRPC
   - State: Zustand

Is this correct? Any modifications?"

User: "Yes, that's right"

AI: "Creating project structure..."
```

**Output:**
- `.agent/project.md` with tech stack documentation
- Recommended directory structure
- Initial configuration files

---

#### create-prd.md Example

**Scenario**: Defining a user authentication feature

**User Request:**
```
@.agent/workflows/create-prd.md

I want to build user authentication with email/password and social login
```

**AI Conversation:**
```
AI: "I'll create a PRD for the authentication feature. Let me ask some clarifying questions:

1. Goals: What's the primary objective?
   User: "Secure user access with multiple login options"

2. User Stories: Who are the users?
   User: "End users of the e-commerce platform"

3. OAuth Providers: Which social logins?
   User: "Google and GitHub"

4. Session Management: JWT or sessions?
   User: "JWT with refresh tokens"

5. Password Requirements?
   User: "Min 8 chars, 1 uppercase, 1 number"

6. Additional Features?
   User: "Password reset via email, remember me option"

Creating PRD..."
```

**Output:**
- `.agent/tasks/0001-prd-user-authentication.md`

---

#### generate-tasks.md Example

**Scenario**: Breaking down the authentication PRD

**User Request:**
```
@.agent/workflows/generate-tasks.md

Generate tasks for @.agent/tasks/0001-prd-user-authentication.md
```

**AI Response:**
```
AI: "Analyzing PRD... Generating high-level tasks:

Parent Tasks:
1.0 Database Schema & Models
2.0 Authentication API Endpoints
3.0 JWT Token Management
4.0 OAuth Integration (Google, GitHub)
5.0 Password Reset Flow
6.0 Frontend Auth Components
7.0 Testing & Security Validation

Review these parent tasks. Say 'Go' to generate detailed sub-tasks."

User: "Go"

AI: "Generating detailed sub-tasks...

1.0 Database Schema & Models
  1.1 Create User model with email, password hash, roles
  1.2 Create Session/RefreshToken model
  1.3 Create PasswordReset model
  1.4 Write database migrations
  1.5 Add model unit tests

2.0 Authentication API Endpoints
  2.1 POST /auth/register - user registration
  2.2 POST /auth/login - email/password login
  2.3 POST /auth/logout - session invalidation
  2.4 POST /auth/refresh - token refresh
  2.5 Add API integration tests

[...continues for all 35 sub-tasks]"
```

**Output:**
- `.agent/tasks/tasks-0001-prd-user-authentication.md`

---

### Quality Workflows

#### code-review.md Example

**Scenario**: Pre-commit review of API changes

**User Request:**
```
@.agent/workflows/code-review.md

Review the changes in src/api/auth/
```

**AI Response:**
```
AI: "Running code review...

## Code Review Report

### Files Reviewed:
- src/api/auth/login.ts
- src/api/auth/register.ts
- src/api/auth/middleware.ts

### Results:

✅ PASS: Code Quality
- All functions < 50 lines
- No cyclomatic complexity issues
- Types properly defined

⚠️ WARNING: Security
- Line 45 in login.ts: Consider rate limiting
- Line 23 in register.ts: Add email validation

✅ PASS: Testing
- Unit tests present
- Coverage: 85%

❌ FAIL: Input Validation
- Line 12 in register.ts: Missing Zod schema validation

### Recommendations:
1. Add rate limiting middleware to login endpoint
2. Add Zod schema for registration input
3. Consider adding brute-force protection

Fix required items before committing."
```

---

#### security-audit.md Example

**Scenario**: Quick security scan before deployment

**User Request:**
```
@.agent/workflows/security-audit.md

Quick scan of the authentication module before we deploy
```

**AI Response:**
```
AI: "Running quick security scan (scope: authentication module)...

## Security Audit Report

### Scope: src/api/auth/, src/lib/auth/

### OWASP Top 10 Check:

✅ A01 Broken Access Control: PASS
   - Role checks implemented correctly

✅ A02 Cryptographic Failures: PASS
   - bcrypt used for password hashing (cost factor 12)
   - JWT signed with RS256

⚠️ A03 Injection: WARNING
   - SQL queries use Prisma (parameterized) ✅
   - User input in error messages (line 67) - sanitize

✅ A07 Auth Failures: PASS
   - Account lockout after 5 attempts
   - Secure session management

### Dependencies:
✅ No known vulnerabilities in auth-related packages

### Secrets Check:
✅ No hardcoded secrets found
✅ Environment variables used correctly

### Summary:
- Critical: 0
- High: 0
- Medium: 1 (error message sanitization)
- Low: 0

Ready for deployment with minor fix."
```

---

#### testing-strategy.md Example

**Scenario**: Improving test coverage for payment module

**User Request:**
```
@.agent/workflows/testing-strategy.md

Plan testing for the payment module - we're at 45% coverage
```

**AI Response:**
```
AI: "Analyzing payment module test coverage...

## Current State:
- Coverage: 45% (target: 80%)
- Files: 12 source files, 4 test files
- Missing: Integration tests, edge cases

## Test Pyramid Plan:

### Unit Tests (Priority 1):
- [ ] PaymentProcessor.processPayment() - 8 test cases
- [ ] RefundHandler.initiateRefund() - 5 test cases
- [ ] PriceCalculator.calculateTotal() - 10 test cases

### Integration Tests (Priority 2):
- [ ] Stripe API integration - 4 test cases
- [ ] Database transaction rollback - 3 test cases
- [ ] Webhook handling - 6 test cases

### Edge Cases (Priority 3):
- [ ] Currency conversion edge cases
- [ ] Partial refund scenarios
- [ ] Network timeout handling

## Estimated Coverage After:
- Unit tests: +25% → 70%
- Integration tests: +10% → 80%
- Edge cases: +5% → 85%

Start with unit tests? Say 'Go' to begin implementation."
```

---

### Maintenance Workflows

#### cleanup-project.md Example

**Scenario**: Quarterly cleanup of .agent/ directory

**User Request:**
```
@.agent/workflows/cleanup-project.md

Clean up this project's .agent/ directory - it's getting bloated
```

**AI Response:**
```
AI: "Analyzing .agent/ directory...

## Cleanup Analysis:

### Current Size: 2.4 MB

### Active Languages Detected:
- TypeScript ✅ (keep typescript.md)
- Python ❌ (no .py files found)
- Go ❌ (no .go files found)

### Recommendations:

#### Remove (unused guides):
- .agent/language-guides/python.md (156 KB)
- .agent/language-guides/go.md (142 KB)
- .agent/language-guides/rust.md (178 KB)
- .agent/framework-guides/django.md (98 KB)
- .agent/framework-guides/flask.md (87 KB)

#### Archive (stale tasks):
- .agent/tasks/0001-prd-old-feature.md (6 months old)
- .agent/tasks/tasks-0001-*.md

#### Consolidate (old memory files):
- 12 memory files from Q1 2025 → quarterly-summary-2025-q1.md

### Projected Savings: 1.8 MB (75% reduction)

Proceed with cleanup? (y/n)"
```

---

#### refactoring.md Example

**Scenario**: Breaking up a large function

**User Request:**
```
@.agent/workflows/refactoring.md

Refactor the processOrder function - it's 120 lines and hard to test
```

**AI Response:**
```
AI: "Analyzing processOrder function...

## Refactoring Plan

### Current State:
- Location: src/services/orderService.ts:45-165
- Lines: 120 (limit: 50)
- Cyclomatic Complexity: 15 (limit: 10)
- Test Coverage: 30%

### Identified Responsibilities:
1. Validate order items (lines 48-72)
2. Calculate pricing (lines 73-95)
3. Check inventory (lines 96-115)
4. Create order record (lines 116-140)
5. Send notifications (lines 141-160)

### Refactoring Steps:

Step 1: Extract validateOrderItems()
- Create function with single responsibility
- Add unit tests
- Replace inline code

Step 2: Extract calculateOrderTotal()
- Pure function, easy to test
- Add unit tests

Step 3: Extract checkInventory()
- Async function
- Add integration tests

Step 4: Extract createOrderRecord()
- Database transaction
- Add tests

Step 5: Extract sendOrderNotifications()
- Side effects isolated
- Add tests with mocks

### After Refactoring:
- processOrder: ~20 lines (orchestration only)
- 5 focused helper functions
- Estimated coverage: 85%

Start with Step 1?"
```

---

#### dependency-update.md Example

**Scenario**: Security vulnerability in a dependency

**User Request:**
```
@.agent/workflows/dependency-update.md

We got a security alert for lodash - need to update urgently
```

**AI Response:**
```
AI: "Analyzing dependency vulnerability...

## Vulnerability Report

### Affected Package:
- lodash@4.17.19
- Severity: HIGH
- CVE: CVE-2021-23337
- Issue: Prototype pollution

### Update Plan:

#### Phase 1: Assessment
- Current version: 4.17.19
- Fixed version: 4.17.21
- Breaking changes: None

#### Phase 2: Update
```bash
npm update lodash
```

#### Phase 3: Test
- Run full test suite
- Check lodash usage in codebase
- Verify no runtime errors

### Files Using lodash:
- src/utils/helpers.ts (import { merge, cloneDeep })
- src/services/dataTransform.ts (import { groupBy })

### Risk Assessment:
- LOW: Minor version update
- No breaking changes expected
- Well-tested functions used

Proceed with update?"
```

---

#### update-framework.md Example

**Scenario**: Updating to latest AICoF version

**User Request:**
```
@.agent/workflows/update-framework.md

Update to the latest version of AICoF
```

**AI Response:**
```
AI: "Checking current installation...

## Version Check:
- Current: v1.5.0
- Latest: v1.6.0

## What's New in v1.6.0:
- New workflow: document-work.md
- New workflow: update-framework.md
- Detailed usage examples in README
- Repository name fixes

## Your Customizations Detected:
- CLAUDE.md: No modifications
- .agent/project.md: Project-specific (will preserve)
- .agent/patterns.md: Has custom patterns (will preserve)

## Update Strategy: Full Replace (recommended)

### Steps:
1. Backup project-specific files
2. Download v1.6.0
3. Replace template files
4. Restore your files

Proceed with update? (y/n)"
```

---

### Utility Workflows

#### troubleshooting.md Example

**Scenario**: API returning 500 errors

**User Request:**
```
@.agent/workflows/troubleshooting.md

Debug: API returns 500 error on /api/orders endpoint
```

**AI Response:**
```
AI: "Starting systematic debugging...

## Problem: 500 Error on /api/orders

### Phase 1: Reproduce
- Endpoint: GET /api/orders
- Error: 500 Internal Server Error
- Frequency: Consistent

### Phase 2: Isolate
Checking logs...

Found: TypeError: Cannot read property 'id' of undefined
Location: src/api/orders/index.ts:34

### Phase 3: Analyze

Code at line 34:
```typescript
const userId = req.user.id; // req.user is undefined
```

Root Cause: Auth middleware not running before this route.

### Phase 4: Fix Options

Option A: Add auth middleware
```typescript
router.get('/orders', authMiddleware, getOrders);
```

Option B: Check for auth in handler
```typescript
if (!req.user) return res.status(401).json({ error: 'Unauthorized' });
```

Recommendation: Option A (middleware approach)

### Phase 5: Verify
After fix, run:
```bash
npm test src/api/orders/
curl -H "Authorization: Bearer $TOKEN" localhost:3000/api/orders
```

Apply Option A fix?"
```

---

#### document-work.md Example

**Scenario**: End of sprint documentation

**User Request:**
```
@.agent/workflows/document-work.md

Document the work from this sprint - we implemented the payment integration
and made some architectural decisions
```

**AI Response:**
```
AI: "Analyzing recent commits (last 2 weeks)...

## Work Summary

### Commits Analyzed: 47
### Files Changed: 32
### Key Changes Identified:

## Patterns Discovered:

### 1. Payment Webhook Pattern
Used for: Stripe webhook handling
```typescript
// Verify signature → Parse event → Route to handler → Acknowledge
```
Adding to patterns.md...

### 2. Retry with Exponential Backoff
Used for: External API calls
```typescript
// attempt → wait(2^n seconds) → retry until maxAttempts
```
Adding to patterns.md...

## Decisions Made:

### Decision 1: Stripe over PayPal
Creating memory file: 2025-01-14-payment-provider-choice.md
- Context: Needed payment integration
- Options: Stripe, PayPal, Square
- Decision: Stripe (better API, webhooks, documentation)

### Decision 2: Webhook Queue Architecture
Creating memory file: 2025-01-14-webhook-processing.md
- Context: Handle high webhook volume
- Decision: Queue-based processing with Redis

## State Update:
Updating state.md with:
- Completed: Payment integration
- In Progress: Subscription billing
- Next: Invoice generation

Documentation complete!"
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
