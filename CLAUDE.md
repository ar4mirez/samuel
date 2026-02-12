# CLAUDE.md

AI-assisted development instructions. Opinionated guardrails for writing quality software.

> **AGENTS.md Compatible**: This file follows the [AGENTS.md](https://agents.md) standard structure.
> Operations section first, then methodology. For other AI tools, copy: `cp CLAUDE.md AGENTS.md`

---

## Operations

> **Purpose**: Immediate, executable instructions for AI agents. Commands first, context second.
> This section is designed to be compatible with the AGENTS.md standard.

### Setup Commands

```bash
# Clone and install (customize per project)
git clone <repo-url> && cd <project>

# Node.js/TypeScript
npm install          # or: yarn install | pnpm install
npm run dev          # Start development server

# Python
python -m venv venv && source venv/bin/activate
pip install -r requirements.txt
python manage.py runserver  # Django
uvicorn main:app --reload   # FastAPI

# Go
go mod download
go run ./cmd/api

# Rust
cargo build
cargo run
```

### Testing Commands

```bash
# Node.js/TypeScript
npm test             # Run all tests
npm run test:watch   # Watch mode
npm run test:cov     # With coverage (target: >80% business logic)

# Python
pytest               # Run all tests
pytest -v            # Verbose
pytest --cov=src     # With coverage

# Go
go test ./...        # All tests
go test -cover ./... # With coverage
go test -race ./...  # Race detection

# Rust
cargo test           # All tests
cargo test -- --nocapture  # With output
cargo tarpaulin      # Coverage
```

### Build & Deploy Commands

```bash
# Node.js/TypeScript
npm run build        # Production build
npm run lint         # Lint check
npm run lint:fix     # Auto-fix lint issues
npm run typecheck    # TypeScript validation

# Python
black .              # Format
isort .              # Sort imports
mypy .               # Type check
ruff check .         # Lint

# Go
go build ./...       # Build all
go vet ./...         # Static analysis
golangci-lint run    # Comprehensive lint

# Rust
cargo build --release  # Production build
cargo clippy           # Lint
cargo fmt              # Format
```

### Code Style Commands

```bash
# Format before commit (all languages)
npm run format       # Prettier (JS/TS)
black . && isort .   # Python
go fmt ./...         # Go
cargo fmt            # Rust

# Lint before commit
npm run lint         # ESLint (JS/TS)
ruff check .         # Python
golangci-lint run    # Go
cargo clippy         # Rust
```

### Environment Variables

```bash
# Copy example env file (never commit .env)
cp .env.example .env

# Required variables (set in .env)
DATABASE_URL=        # Database connection string
API_KEY=             # External service API key
SECRET_KEY=          # Application secret (generate random)
NODE_ENV=            # development | production | test
```

---

## Boundaries (Do Not Touch)

> **Critical**: AI agents must NOT modify these without explicit permission.

### Protected Files
- `package-lock.json`, `yarn.lock`, `Cargo.lock`, `go.sum` (dependency locks)
- `.env`, `.env.local`, `.env.production` (environment configs)
- Database migration files (once applied to production)
- `tsconfig.json`, `Cargo.toml`, `go.mod` (without testing build)
- `.github/workflows/`, CI/CD configurations

### Never Commit
- Secrets, API keys, credentials, tokens
- `.env` files (commit `.env.example` instead)
- `node_modules/`, `venv/`, `target/`, build artifacts
- Personal IDE configs (`.vscode/settings.json`, `.idea/`)
- Large binaries, datasets (use Git LFS)

### Ask Before Modifying
- Authentication/authorization logic
- Database schemas (after production deployment)
- Public API contracts (breaks consumers)
- Build/deploy processes
- Major dependency versions

---

## Quick Reference

**Task Classification:**
- **ATOMIC** (<5 files, clear scope) → Implement directly
- **FEATURE** (5-10 files) → Break into subtasks
- **COMPLEX** (>10 files, new subsystem) → Use PRD workflow | .claude/skills/create-prd/SKILL.md

**Common Guardrails** (validate first):
✓ Function ≤50 lines | ✓ File ≤300 lines | ✓ Input validation | ✓ Parameterized queries
✓ Tests >80% (critical) | ✓ Conventional commits | ✓ No secrets in code

**Emergency Quick Links:**
- Security issue? → .claude/skills/security-audit/SKILL.md
- Tests failing? → .claude/skills/troubleshooting/SKILL.md
- Stuck >30 min? → .claude/skills/troubleshooting/SKILL.md
- Complex feature? → .claude/skills/create-prd/SKILL.md
- Code review? → .claude/skills/code-review/SKILL.md
- Language-specific? → .claude/skills/{lang}-guide/SKILL.md

**Workflows** (on-demand):

- Planning: initialize-project, create-rfd, create-prd, generate-tasks
- Quality: code-review, security-audit, testing-strategy
- Maintenance: cleanup-project, refactoring, dependency-update, update-framework
- Utility: troubleshooting, generate-agents-md, document-work, create-skill

**Skills** (capability modules - [Agent Skills](https://agentskills.io) standard):
- Create: `aicof skill create <name>` or `.claude/skills/create-skill/SKILL.md`
- Validate: `aicof skill validate`
- List: `aicof skill list`
- Load: `.claude/skills/<skill-name>/SKILL.md` when task matches description
- Spec: https://agentskills.io/specification

<!-- SKILLS_START -->
## Available Skills

Skills extend AI capabilities. Load a skill when task matches its description.

| Skill | Description |
|-------|-------------|
| commit-message | Generate descriptive commit messages by analyzing git diffs. Use when the us... |

**To use a skill**: Read `.claude/skills/<skill-name>/SKILL.md`
<!-- SKILLS_END -->

**RFD vs PRD** (when exploring options):

- **RFD** = "Why" (explore options, build consensus) → .claude/skills/create-rfd/SKILL.md
- **PRD** = "What" (define implementation) → .claude/skills/create-prd/SKILL.md
- Flow: Idea → RFD (explore) → Decision → PRD (plan) → Tasks → Code

**Load Language Guide** (automatic based on file extensions — language guides are Agent Skills):
- TypeScript/JavaScript → .claude/skills/typescript-guide/SKILL.md
- Python → .claude/skills/python-guide/SKILL.md
- Go → .claude/skills/go-guide/SKILL.md
- Rust → .claude/skills/rust-guide/SKILL.md
- Kotlin → .claude/skills/kotlin-guide/SKILL.md
- Java → .claude/skills/java-guide/SKILL.md
- C# → .claude/skills/csharp-guide/SKILL.md
- PHP → .claude/skills/php-guide/SKILL.md
- Swift → .claude/skills/swift-guide/SKILL.md
- C/C++ → .claude/skills/cpp-guide/SKILL.md
- Ruby → .claude/skills/ruby-guide/SKILL.md
- SQL → .claude/skills/sql-guide/SKILL.md
- Shell/Bash → .claude/skills/shell-guide/SKILL.md
- R → .claude/skills/r-guide/SKILL.md
- Dart/Flutter → .claude/skills/dart-guide/SKILL.md
- HTML/CSS → .claude/skills/html-css-guide/SKILL.md
- Lua → .claude/skills/lua-guide/SKILL.md
- Assembly → .claude/skills/assembly-guide/SKILL.md
- CUDA → .claude/skills/cuda-guide/SKILL.md
- Solidity → .claude/skills/solidity-guide/SKILL.md
- Zig → .claude/skills/zig-guide/SKILL.md

**Load Framework Skill** (when using specific framework — framework guides are Agent Skills):
- TypeScript/JavaScript:
  - React → .claude/skills/react/SKILL.md
  - Next.js → .claude/skills/nextjs/SKILL.md
  - Express → .claude/skills/express/SKILL.md
- Python:
  - Django → .claude/skills/django/SKILL.md
  - FastAPI → .claude/skills/fastapi/SKILL.md
  - Flask → .claude/skills/flask/SKILL.md
- Go:
  - Gin → .claude/skills/gin/SKILL.md
  - Echo → .claude/skills/echo/SKILL.md
  - Fiber → .claude/skills/fiber/SKILL.md
- Rust:
  - Axum → .claude/skills/axum/SKILL.md
  - Actix-web → .claude/skills/actix-web/SKILL.md
  - Rocket → .claude/skills/rocket/SKILL.md
- Kotlin:
  - Spring Boot (Kotlin) → .claude/skills/spring-boot-kotlin/SKILL.md
  - Ktor → .claude/skills/ktor/SKILL.md
  - Android Compose → .claude/skills/android-compose/SKILL.md
- Java:
  - Spring Boot → .claude/skills/spring-boot-java/SKILL.md
  - Quarkus → .claude/skills/quarkus/SKILL.md
  - Micronaut → .claude/skills/micronaut/SKILL.md
- C#:
  - ASP.NET Core → .claude/skills/aspnet-core/SKILL.md
  - Blazor → .claude/skills/blazor/SKILL.md
  - Unity → .claude/skills/unity/SKILL.md
- PHP:
  - Laravel → .claude/skills/laravel/SKILL.md
  - Symfony → .claude/skills/symfony/SKILL.md
  - WordPress → .claude/skills/wordpress/SKILL.md
- Swift:
  - SwiftUI → .claude/skills/swiftui/SKILL.md
  - UIKit → .claude/skills/uikit/SKILL.md
  - Vapor → .claude/skills/vapor/SKILL.md
- Ruby:
  - Rails → .claude/skills/rails/SKILL.md
  - Sinatra → .claude/skills/sinatra/SKILL.md
  - Hanami → .claude/skills/hanami/SKILL.md
- Dart:
  - Flutter → .claude/skills/flutter/SKILL.md
  - Shelf → .claude/skills/shelf/SKILL.md
  - Dart Frog → .claude/skills/dart-frog/SKILL.md

---

## Core Guardrails (ALWAYS ENFORCE)

### Code Quality
- ✓ No function exceeds 50 lines (split with helper functions)
- ✓ No file exceeds 300 lines (components: 200, tests: 300, utils: 150)
- ✓ Cyclomatic complexity ≤ 10 per function
- ✓ All exported functions have type signatures and documentation
- ✓ No magic numbers (use named constants)
- ✓ No commented-out code in commits (use git history)
- ✓ No `TODO` without issue/ticket reference
- ✓ No dead code (unused imports, variables, functions)

### Security (CRITICAL)
- ✓ All user inputs validated before processing
- ✓ All API boundaries have input validation (prefer schema validators: Zod, Pydantic, etc.)
- ✓ All database queries use parameterized statements (no string concatenation)
- ✓ All environment variables have secure defaults (never hardcode secrets)
- ✓ All file operations validate paths (prevent directory traversal)
- ✓ All async operations have timeout/cancellation mechanisms
- ✓ Dependencies checked for known vulnerabilities before adding
- ✓ Dependencies checked for license compatibility before adding
- ✓ All database migrations include rollback (down) function

### Testing (CRITICAL)
- ✓ Coverage targets: >80% for business logic, >60% overall
- ✓ All public APIs have unit tests
- ✓ All bug fixes include regression tests
- ✓ All edge cases explicitly tested (null, empty, boundary values)
- ✓ Test names describe behavior: `test_user_login_fails_with_invalid_password`
- ✓ No test interdependencies (tests run in any order)
- ✓ Integration tests for external service interactions
- ✓ All deployments include smoke test validation

### Git & Commits
- ✓ Commit messages: `type(scope): description` (conventional commits)
- ✓ Types: feat, fix, docs, refactor, test, chore, perf, ci
- ✓ One logical change per commit (atomic commits)
- ✓ All commits must pass tests before pushing
- ✓ Branch naming: `type/short-description` (e.g., `feat/user-auth`)
- ✓ No commits directly to main/master (use PRs)
- ✓ Breaking API changes require major version bump (Semantic Versioning)

### Performance
- ✓ No N+1 queries (batch database operations)
- ✓ Large datasets use pagination/streaming (not full loads)
- ✓ Expensive computations memoized/cached when appropriate
- ✓ Frontend bundles < 200KB initial load (code-split when needed)
- ✓ API responses < 200ms for simple queries, < 1s for complex

---

## 4D Methodology (Enhanced)

Apply appropriate mode based on task complexity:

### ATOMIC Mode (Default)
For single-file changes, bug fixes, small features:

1. **Deconstruct**: What's the minimal change needed?
2. **Diagnose**: Will this break anything? Check dependencies.
3. **Develop**: Make the change with tests.
4. **Deliver**: Validate (run tests, check guardrails) → Commit.

### FEATURE Mode
For multi-file features, new components, API endpoints:

1. **Deconstruct**: Break into 3-5 subtasks (each atomic).
2. **Diagnose**: Identify integration points and dependencies.
3. **Develop**: Implement subtasks sequentially with tests.
4. **Deliver**: Integration test → Documentation → Review → Commit.

### COMPLEX Mode
For architecture changes, major refactors, new systems:

1. **Deconstruct**: Full decomposition into phases/milestones.
2. **Diagnose**: Analyze risks, dependencies, migration paths.
3. **Develop**: Define plan in CLAUDE.md or per-folder CLAUDE.md → Implement incrementally.
4. **Deliver**: Staged rollout → Documentation → Retrospective.

**Structured Workflow (for COMPLEX tasks):**

**MANDATORY when:**
- New subsystem (authentication, payments, real-time features)
- Breaking architectural changes
- Affects >15 files
- Unclear requirements (needs scope definition)

**RECOMMENDED when:**
- Affects 10-15 files
- Multiple stakeholders
- Complex domain logic

**OPTIONAL when:**
- Affects <10 files with clear scope
- Well-defined refactoring

**Workflow:**
1. Use `.claude/skills/create-prd/SKILL.md` to define requirements
2. Use `.claude/skills/generate-tasks/SKILL.md` to break down implementation
3. Implement tasks step-by-step with verification checkpoints

**Escalation Triggers:**
- Task affects >5 files → FEATURE mode
- Task affects >10 files → COMPLEX mode (consider PRD workflow)
- Task affects >15 files OR new subsystem → COMPLEX mode (PRD workflow MANDATORY)
- Task unclear/ambiguous → Ask user for clarification first

---

## Software Development Lifecycle

### Stage 1: Planning
**Atomic**: Read existing code → Identify change location
**Feature**: Review related code → Sketch interfaces/contracts
**Complex**: Use PRD workflow or document plan in CLAUDE.md

**Checkpoints:**
- [ ] Requirements clear and testable
- [ ] Scope defined (what's included, what's not)
- [ ] Dependencies identified
- [ ] Breaking changes flagged

### Stage 2: Implementation
**Always:**
- Write tests first (TDD) or alongside code
- Load language-specific guide: .claude/skills/{language}-guide/SKILL.md
- Follow language/framework idioms
- Validate against guardrails continuously
- Keep changes focused (resist scope creep)

**Checkpoints:**
- [ ] Code follows all guardrails (review each ✓ above)
- [ ] Tests written and passing
- [ ] No linter errors/warnings
- [ ] Types correct (strict mode)

### Stage 3: Validation
**Automated:**
- Run full test suite
- Check coverage thresholds
- Run linter/formatter
- Verify build succeeds

**Manual:**
- [ ] Edge cases considered
- [ ] Error handling implemented
- [ ] Performance acceptable
- [ ] Security implications reviewed

### Stage 4: Documentation
**Code Level:**
- Function/class docstrings (what, why, params, returns)
- Inline comments for complex logic (why, not what)
- API documentation updated

**Project Level:**
- Update CLAUDE.md or per-folder CLAUDE.md if new patterns emerged
- Key decisions are captured via Claude Code's built-in memory system

### Stage 5: Commit
**Pre-Commit Checklist:**
- [ ] All tests pass
- [ ] Coverage thresholds met
- [ ] No linter errors
- [ ] All guardrails validated
- [ ] Commit message follows convention
- [ ] No sensitive data (secrets, credentials, PII)

**Commit:**
```bash
git add <files>
git commit -m "type(scope): description

- Detail 1
- Detail 2

Refs: #issue-number"
```

---

## Context System

### CLAUDE.md (This File)
**Loaded**: Always (every conversation)
**Purpose**: Operations + Guardrails + Methodology + SDLC workflow
**Compatibility**: AGENTS.md standard (copy for other tools: `cp CLAUDE.md AGENTS.md`)

### Per-Folder CLAUDE.md Files
**Loaded**: Automatically when working in that directory
**Purpose**: Directory-specific instructions, overrides, or context
**Example**: `src/api/CLAUDE.md` with API-specific conventions

### .claude/skills/ Directory
**Loaded**: On-demand, when task matches a skill description
**Purpose**: Agent Skills - capability modules following the [Agent Skills](https://agentskills.io) standard

**Structure:**
```
.claude/
├── skills/               # Agent Skills - capability modules
│   ├── README.md              # How to create and use skills
│   ├── commit-message/        # Generate commit messages
│   │   └── SKILL.md
│   ├── initialize-project/    # Project setup workflow
│   │   └── SKILL.md
│   ├── create-rfd/            # Request for Discussion (explore options)
│   │   └── SKILL.md
│   ├── create-prd/            # Requirements documents
│   │   └── SKILL.md
│   ├── generate-tasks/        # Task breakdown
│   │   └── SKILL.md
│   ├── code-review/           # Pre-commit quality review
│   │   └── SKILL.md
│   ├── security-audit/        # Security assessment
│   │   └── SKILL.md
│   ├── testing-strategy/      # Test planning & coverage
│   │   └── SKILL.md
│   ├── cleanup-project/       # Prune unused guides
│   │   └── SKILL.md
│   ├── refactoring/           # Technical debt remediation
│   │   └── SKILL.md
│   ├── dependency-update/     # Safe dependency updates
│   │   └── SKILL.md
│   ├── troubleshooting/       # Debugging workflow
│   │   └── SKILL.md
│   ├── generate-agents-md/    # Cross-tool compatibility
│   │   └── SKILL.md
│   ├── document-work/         # Capture patterns and decisions
│   │   └── SKILL.md
│   ├── update-framework/      # Update AICoF while preserving customizations
│   │   └── SKILL.md
│   ├── create-skill/          # Create Agent Skills
│   │   └── SKILL.md
│   ├── go-guide/              # Language guide skills (21 languages)
│   │   ├── SKILL.md           # Core guardrails and patterns
│   │   └── references/        # Detailed patterns, pitfalls, security
│   ├── python-guide/
│   │   └── SKILL.md
│   ├── typescript-guide/
│   │   └── SKILL.md
│   ├── ...                    # All 21 language guide skills
│   ├── react/                 # Framework skills (33 frameworks)
│   │   └── SKILL.md
│   ├── nextjs/
│   │   └── SKILL.md
│   ├── django/
│   │   └── SKILL.md
│   └── ...                    # All 33 framework skills
└── settings.local.json   # Claude Code local settings
```

**Loading Protocol:**
- **Session Start**: Claude Code loads CLAUDE.md automatically, plus any per-folder CLAUDE.md files
- **During Work**: Load language guide skill from `.claude/skills/<lang>-guide/SKILL.md` based on file extensions
- **Complex Features**: Load workflow skills (PRD, task generation)
- **Capability Needed**: Load skill from `.claude/skills/` when task matches skill description
- **Memory & State**: Claude Code's built-in memory system handles session state, patterns, and decisions

---

## Anti-Patterns (Avoid These)

### Code
- ❌ Premature optimization (measure first, optimize after)
- ❌ Over-engineering (YAGNI: You Aren't Gonna Need It)
- ❌ Copy-paste code (extract to shared function/component)
- ❌ Ignoring errors (every error needs handling)
- ❌ Mutable globals (use dependency injection, immutable state)

### Testing
- ❌ Testing implementation details (test behavior, not internals)
- ❌ Flaky tests (non-deterministic results indicate bad design)
- ❌ Test interdependencies (each test should be isolated)
- ❌ No assertions (tests must verify something)

### Process
- ❌ Committing directly to main (use feature branches + PRs)
- ❌ Batch commits (commit after each logical change)
- ❌ Skipping tests because "it's a small change" (small bugs exist too)
- ❌ Not reading error messages fully (they often tell you the fix)

### Documentation
- ❌ Outdated docs (worse than no docs)
- ❌ Obvious comments (`i++; // increment i`)
- ❌ No comments on complex logic (future you needs context)
- ❌ Docs separate from code (keep close: docstrings, inline, README)

---

## Initialization

**For new projects:** Use `.claude/skills/initialize-project/SKILL.md`

**For existing projects:** Use `.claude/skills/initialize-project/SKILL.md`

AI will ask questions, analyze codebase, and configure CLAUDE.md with findings.

---

## When Stuck

**See:** `.claude/skills/troubleshooting/SKILL.md`

**Quick recovery:**
1. STOP trying random solutions (>30 min = stuck)
2. Document what you've tried
3. Simplify & isolate (minimal reproduction)
4. Check fundamentals (dependencies, config, versions)
5. Ask user with clear problem statement
6. Record solution for future reference

---

## Success Criteria

### Per Task
- [ ] All guardrails validated
- [ ] Tests pass with coverage thresholds
- [ ] No security vulnerabilities introduced
- [ ] Documentation updated
- [ ] Commit follows conventions

### Per Session
- [ ] No broken tests left behind
- [ ] Progress measurable
- [ ] Key decisions captured for future reference

---

## Version & Changelog

**Current Version**: 2.0.0
**Last Updated**: 2026-02-11

### Changelog

**v2.0.0 (2026-02-11) - Migration from .agent/ to .claude/**
- BREAKING: Migrated all skills from `.agent/skills/` to `.claude/skills/`
- BREAKING: Removed `.agent/` directory structure (memory/, tasks/, rfd/, state.md, patterns.md)
- Replaced `.agent/` context system with Claude Code native conventions:
  - CLAUDE.md (root) - always loaded
  - Per-folder CLAUDE.md files - loaded when working in that directory
  - `.claude/skills/` - Agent Skills directory (native Claude Code skills)
- Session state, patterns, and decisions now use Claude Code's built-in memory system
- Updated AGENTS.md compatibility from symlink to real copy (`cp CLAUDE.md AGENTS.md`)
- Simplified Context System section
- Simplified Success Criteria (removed .agent/-dependent per-session items)
- Updated SDLC documentation stage to use CLAUDE.md instead of .agent/ files

**v1.8.0 (2025-02-04) - Agent Skills Integration**
- ✅ Added Agent Skills support (open standard for AI agent capabilities)
- ✅ New CLI commands: `aicof skill create`, `aicof skill validate`, `aicof skill list`, `aicof skill info`
- ✅ Added `.claude/skills/` directory for project skills
- ✅ Added `commit-message` example skill demonstrating the pattern
- ✅ Added `create-skill.md` workflow for creating new skills
- ✅ Updated Quick Reference with Skills section
- ✅ Extended `aicof doctor` to validate skills
- ✅ Extended `aicof search` to include skills
- ✅ Now supports 15 workflows total (was 14)
- ✅ Skills compatible with 25+ agent products (Claude Code, Cursor, VS Code, etc.)

**v1.7.0 (2025-01-14) - Rebrand to AICoF**
- ✅ Rebranded from "AI Code Template" to "AICoF (Artificial Intelligence Coding Framework)"
- ✅ Renamed GitHub repository from `ai-code-template` to `aicof`
- ✅ Updated all documentation, configuration, and references
- ✅ New brand positioning: "Build smarter, faster, and more scalable software"

**v1.6.0 (2025-01-14) - Workflow Enhancements & Documentation**
- ✅ Added 2 new workflows:
  - document-work.md: Capture patterns, decisions, and learnings from recent work
  - update-framework.md: Update AICoF while preserving customizations
- ✅ Added comprehensive usage examples for all 13 workflows in README
- ✅ Fixed repository name references
- ✅ Updated workflow map and documentation
- ✅ Now includes 13 workflows total (was 11)

**v1.5.0 (2025-12-15) - Framework-Specific Templates**
- ✅ Added 33 framework-specific guide templates across 11 languages:
  - TypeScript/JavaScript: React, Next.js, Express
  - Python: Django, FastAPI, Flask
  - Go: Gin, Echo, Fiber
  - Rust: Axum, Actix-web, Rocket
  - Kotlin: Spring Boot (Kotlin), Ktor, Android Compose
  - Java: Spring Boot, Quarkus, Micronaut
  - C#: ASP.NET Core, Blazor, Unity
  - PHP: Laravel, Symfony, WordPress
  - Swift: SwiftUI, UIKit, Vapor
  - Ruby: Rails, Sinatra, Hanami
  - Dart: Flutter, Shelf, Dart Frog
- ✅ Added framework skills under .claude/skills/ (migrated from framework-guides/)
- ✅ Updated Quick Reference with "Load Framework Skill" section
- ✅ Updated .agent/ directory structure documentation

**v1.4.0 (2025-12-13) - Comprehensive Language Coverage**
- ✅ Added 10 more language guides for specialized domains:
  - SQL (PostgreSQL, MySQL, query optimization)
  - Shell/Bash (scripting, automation, POSIX)
  - R (statistical computing, tidyverse, Shiny)
  - Dart (Flutter, mobile development)
  - HTML/CSS (web standards, accessibility, BEM)
  - Lua (scripting, Love2D, Neovim)
  - Assembly (x86-64, ARM64, RISC-V)
  - CUDA (GPU computing, parallel processing)
  - Solidity (Ethereum, smart contracts, DeFi)
  - Zig (systems programming, C interop)
- ✅ Now supports 21 languages total
- ✅ Updated Quick Reference with all language guides
- ✅ Updated .agent/ directory structure documentation

**v1.3.0 (2025-12-13) - Extended Language Support**
- ✅ Added 6 new language guides covering top programming languages:
  - Java (Spring Boot, enterprise, JUnit 5)
  - C# (.NET 8, ASP.NET Core, xUnit)
  - PHP (Laravel, Symfony, PHPUnit)
  - Swift (iOS, macOS, SwiftUI)
  - C/C++ (C++20, CMake, GoogleTest)
  - Ruby (Rails 7, RSpec, Sidekiq)
- ✅ Now supports 11 languages total (TypeScript, Python, Go, Rust, Kotlin + new)
- ✅ Updated Quick Reference with all language guides
- ✅ Updated .agent/ directory structure documentation

**v1.2.0 (2025-12-13) - AGENTS.md Compatibility**
- ✅ Added Operations section (AGENTS.md compatible)
  - Setup Commands (Node, Python, Go, Rust)
  - Testing Commands with coverage targets
  - Build & Deploy Commands
  - Code Style Commands (format, lint)
  - Environment Variables template
- ✅ Added Boundaries section (prominent "do not touch" list)
- ✅ Restructured for "commands first, context second" pattern
- ✅ Added symlink instructions for cross-tool compatibility
- ✅ Added Kotlin to language guides list
- ✅ Consolidated Protected Boundaries into Boundaries section
- ✅ Added AGENTS.md generator workflow

**v1.1.0 (2025-01-15) - Phase 1 Optimization**
- ✅ Reduced from 490 → 400 lines (18% reduction)
- ✅ Added Quick Reference section
- ✅ Added critical missing guardrails:
  - Dependency license checking
  - Database migration rollbacks
  - Semantic versioning for breaking changes
  - Smoke test validation
- ✅ Clarified workflow requirements (MANDATORY vs RECOMMENDED)
- ✅ Extracted language-specific guides to .agent/language-guides/
- ✅ Extracted initialization to .claude/skills/initialize-project/SKILL.md
- ✅ Extracted troubleshooting to .claude/skills/troubleshooting/SKILL.md
- ✅ Created comprehensive language guides (TypeScript, Python, Go, Rust)

**v1.0.0 (2025-01-14) - Initial Release**
- Initial CLAUDE.md system with 30+ guardrails
- 4D methodology (ATOMIC/FEATURE/COMPLEX)
- Integrated ai-dev-tasks workflows
- Progressive .agent/ directory structure

### Update This File When
- Adding/removing guardrails
- Changing methodology
- Project-wide quality standards shift
- New language guides added
- AGENTS.md standard evolves

---

**Remember**: This file is your guardrails. The `.claude/skills/` directory extends your capabilities. Small atomic changes. Validate continuously. Document progressively.

**Cross-Tool Compatibility**: `cp CLAUDE.md AGENTS.md` for tools that read AGENTS.md
