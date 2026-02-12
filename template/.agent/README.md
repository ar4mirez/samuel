# .agent/ Directory

This directory contains project-specific AI context that grows organically as the project evolves.

## AGENTS.md Compatibility

This system is compatible with the [AGENTS.md](https://agents.md) standard (v1.0, July 2025).

**How it works:**
- `CLAUDE.md` contains Operations section (AGENTS.md compatible) + full methodology
- For other AI tools (Cursor, Codex, Copilot), create symlink: `ln -s CLAUDE.md AGENTS.md`
- Or use `@.agent/workflows/generate-agents-md.md` to generate standalone AGENTS.md

**File Priority:**
| Tool | Primary | Fallback |
|------|---------|----------|
| Claude Code | CLAUDE.md | AGENTS.md |
| Cursor, Codex, etc. | AGENTS.md | CLAUDE.md |

## Philosophy: Progressive Growth

**Day 1**: This folder may be empty except templates. That's OK!
**Week 1**: AI creates `project.md` when architecture decisions are made
**Month 1**: `patterns.md` populated with discovered conventions
**Ongoing**: `memory/` captures significant decisions and learnings

## Structure

```
.agent/
├── README.md              # This file
├── project.md             # Tech stack, architecture, key decisions
├── patterns.md            # Coding patterns and conventions
├── state.md              # Current work, tasks, blockers
├── skills/               # Agent Skills (language guides, utilities)
│   ├── README.md         # Skills documentation
│   ├── go-guide/         # Go language guide skill
│   │   ├── SKILL.md      # Core guardrails and patterns
│   │   └── references/   # Detailed patterns, pitfalls, security
│   ├── typescript-guide/  # TypeScript language guide skill
│   ├── python-guide/      # Python language guide skill
│   └── ... (18 more)     # rust, kotlin, java, csharp, php, swift, cpp, ruby, sql, shell, r, dart, html-css, lua, assembly, cuda, solidity, zig
├── framework-guides/     # Framework-specific patterns (33 frameworks)
│   ├── README.md         # Index of all framework guides
│   ├── react.md          # React
│   ├── django.md         # Django
│   ├── rails.md          # Rails
│   └── ... (30 more)     # Next.js, Express, FastAPI, Flask, Gin, Echo, Fiber, Axum, Actix-web, Rocket, Spring Boot (Kotlin), Ktor, Android Compose, Spring Boot (Java), Quarkus, Micronaut, ASP.NET Core, Blazor, Unity, Laravel, Symfony, WordPress, SwiftUI, UIKit, Vapor, Sinatra, Hanami, Flutter, Shelf, Dart Frog
├── workflows/            # Structured workflows for complex features
│   ├── README.md         # Workflow documentation
│   ├── create-prd.md     # Product Requirements Document workflow
│   ├── generate-tasks.md # Task breakdown workflow
│   ├── initialize-project.md # Project initialization workflow
│   ├── troubleshooting.md # Debugging workflow
│   └── generate-agents-md.md # AGENTS.md generator for cross-tool compatibility
├── tasks/                # PRDs and task lists for complex features
│   ├── NNNN-prd-feature-name.md
│   └── tasks-NNNN-prd-feature-name.md
└── memory/               # Decision log (one file per topic)
    └── YYYY-MM-DD-topic.md
```

## When Files Are Created

### project.md
**Trigger**: First architecture decision or tech stack choice
**Created by**: AI (with user confirmation)
**Contains**:
- Tech stack and versions
- Architecture patterns
- Framework choices
- Key design decisions

**Example creation moment**:
- User: "Let's use React with TypeScript"
- AI: Creates `project.md` documenting this choice

### patterns.md
**Trigger**: Same pattern used 3+ times
**Created by**: AI (observes repetition)
**Contains**:
- API error handling patterns
- Database query patterns
- Component structure conventions
- Testing patterns

**Example creation moment**:
- AI notices consistent use of error boundary pattern
- AI: "I've noticed we use this error handling pattern repeatedly. Should I document it in `.agent/patterns.md`?"

### state.md
**Trigger**: Multi-session work begins
**Created by**: User request or AI suggestion
**Contains**:
- Current sprint/tasks
- Work in progress
- Blockers and dependencies
- Next steps

**Example creation moment**:
- User: "We'll be working on this over the next few weeks"
- AI: Creates `state.md` to track progress across sessions

### memory/YYYY-MM-DD-topic.md
**Trigger**: Significant decision or complex problem solved
**Created by**: AI (after key moments)
**Contains**:
- Context: What was the problem/decision?
- Analysis: What options were considered?
- Decision: What was chosen and why?
- Outcome: What happened as a result?

**Example creation moment**:
- After debugging complex authentication issue
- After choosing between competing architectures
- After discovering non-obvious solution

### workflows/ (Complex Feature Development)
**Trigger**: COMPLEX mode scenario (>10 files, new subsystem, unclear scope)
**Created by**: Already exists (templates provided)
**Contains**:
- `create-prd.md` - Workflow for generating Product Requirements Documents
- `generate-tasks.md` - Workflow for breaking PRDs into task lists
- `README.md` - Workflow documentation and usage guide

**When to use**:
- Building major feature (user authentication, payment system, etc.)
- Requirements unclear or ambiguous
- Multiple stakeholders need alignment
- Feature takes >1 week to implement

**Example workflow**:
1. User: "Build user authentication with OAuth"
2. AI detects COMPLEX mode → Suggests PRD workflow
3. AI uses `@.agent/workflows/create-prd.md` → Creates `.agent/tasks/0001-prd-user-auth.md`
4. User reviews and approves PRD
5. AI uses `@.agent/workflows/generate-tasks.md` → Creates `./agent/tasks/tasks-0001-prd-user-auth.md`
6. AI implements tasks step-by-step with verification

**Skip workflows for**:
- ATOMIC tasks (single file, bug fixes)
- Small features (<5 files, clear scope)
- Refactoring with defined boundaries

### tasks/ (PRDs and Task Lists)
**Trigger**: Complex feature development using workflows
**Created by**: AI during workflow execution
**Contains**:
- `NNNN-prd-feature-name.md` - Product Requirements Documents
- `tasks-NNNN-prd-feature-name.md` - Task breakdowns for implementation

**Naming convention**:
- PRDs: 0001, 0002, 0003... (sequential numbering)
- Tasks: Match PRD number (tasks-0001, tasks-0002, etc.)

**Example**:
```
.agent/tasks/
├── 0001-prd-user-authentication.md      # Requirements for auth feature
├── tasks-0001-prd-user-authentication.md  # Task breakdown (42 subtasks)
├── 0002-prd-payment-integration.md       # Requirements for payments
└── tasks-0002-prd-payment-integration.md  # Task breakdown (28 subtasks)
```

## Loading Protocol

### Session Start
1. AI automatically loads CLAUDE.md
2. AI checks for `state.md` → Reads if exists
3. AI asks about architecture if unclear → May read `project.md`

### During Work
- Need pattern reference → Read `patterns.md`
- Need architecture context → Read `project.md`
- Research past decision → Search `memory/`

### Session End
- AI updates `state.md` with progress (if exists)
- AI adds to `patterns.md` if new pattern emerged
- AI creates `memory/YYYY-MM-DD-topic.md` for key decisions

## Templates

This directory includes `.template` files:
- `project.md.template` - Starting point for project documentation
- `state.md.template` - Starting point for work tracking

**To use**: Copy template, remove `.template` suffix, fill in your content.

**AI can do this**: Just ask "Create project.md with our tech stack" and AI will use the template.

## Best Practices

### DO
✓ Let files grow naturally (don't force structure upfront)
✓ Document decisions at the time they're made
✓ Keep memory/ files focused (one topic per file)
✓ Update state.md regularly during active development
✓ Review and consolidate patterns monthly

### DON'T
❌ Create all files on Day 1 (most will be empty/wrong)
❌ Document obvious things (file structure, basic syntax)
❌ Let documentation lag behind code (update together)
❌ Write novels (keep entries concise, scannable)
❌ Duplicate CLAUDE.md content here

## File Size Guidelines

- **project.md**: 50-150 lines (concise overview)
- **patterns.md**: 100-300 lines (grows over time)
- **state.md**: 30-100 lines (current focus only)
- **memory/*.md**: 20-100 lines each (focused on one decision)

If files exceed these ranges, consider:
- Splitting patterns.md by category
- Archiving old state.md entries
- Breaking complex memory entries into multiple files

## Maintenance

### Weekly
- [ ] Update state.md with current priorities
- [ ] Add newly-discovered patterns to patterns.md
- [ ] Archive completed tasks from state.md

### Monthly
- [ ] Review patterns.md for outdated entries
- [ ] Consolidate related memory/ files if applicable
- [ ] Update project.md if tech stack changed

### Per Phase/Release
- [ ] Archive state.md → memory/YYYY-MM-DD-phase-N-summary.md
- [ ] Review all .agent/ files for accuracy
- [ ] Remove obsolete patterns

## Integration with CLAUDE.md

**CLAUDE.md** = Rules and guardrails (universal)
**.agent/** = Context and memory (project-specific)

Think of it as:
- CLAUDE.md: "How we write software" (the process)
- .agent/: "What we're building" (the product)

AI loads CLAUDE.md always, loads .agent/ files on-demand when context is needed.

---

**Remember**: Start minimal. Grow organically. Document progressively.
