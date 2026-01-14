# CLAUDE.md

AI-assisted development instructions. Opinionated guardrails for writing quality software.

> **AGENTS.md Compatible**: This file follows the [AGENTS.md](https://agents.md) standard structure.
> Operations section first, then methodology. For other AI tools, symlink: `ln -s CLAUDE.md AGENTS.md`

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
- **COMPLEX** (>10 files, new subsystem) → Use PRD workflow | .agent/workflows/create-prd.md

**Common Guardrails** (validate first):
✓ Function ≤50 lines | ✓ File ≤300 lines | ✓ Input validation | ✓ Parameterized queries
✓ Tests >80% (critical) | ✓ Conventional commits | ✓ No secrets in code

**Emergency Quick Links:**
- Security issue? → .agent/workflows/security-audit.md
- Tests failing? → .agent/workflows/troubleshooting.md
- Stuck >30 min? → .agent/workflows/troubleshooting.md
- Complex feature? → .agent/workflows/create-prd.md
- Code review? → .agent/workflows/code-review.md
- Language-specific? → .agent/language-guides/

**Workflows** (on-demand):
- Planning: initialize-project, create-prd, generate-tasks
- Quality: code-review, security-audit, testing-strategy
- Maintenance: cleanup-project, refactoring, dependency-update, update-framework
- Utility: troubleshooting, generate-agents-md, document-work

**Load Language Guide** (automatic based on file extensions):
- TypeScript/JavaScript → .agent/language-guides/typescript.md
- Python → .agent/language-guides/python.md
- Go → .agent/language-guides/go.md
- Rust → .agent/language-guides/rust.md
- Kotlin → .agent/language-guides/kotlin.md
- Java → .agent/language-guides/java.md
- C# → .agent/language-guides/csharp.md
- PHP → .agent/language-guides/php.md
- Swift → .agent/language-guides/swift.md
- C/C++ → .agent/language-guides/cpp.md
- Ruby → .agent/language-guides/ruby.md
- SQL → .agent/language-guides/sql.md
- Shell/Bash → .agent/language-guides/shell.md
- R → .agent/language-guides/r.md
- Dart/Flutter → .agent/language-guides/dart.md
- HTML/CSS → .agent/language-guides/html-css.md
- Lua → .agent/language-guides/lua.md
- Assembly → .agent/language-guides/assembly.md
- CUDA → .agent/language-guides/cuda.md
- Solidity → .agent/language-guides/solidity.md
- Zig → .agent/language-guides/zig.md

**Load Framework Guide** (when using specific framework):
- TypeScript/JavaScript:
  - React → .agent/framework-guides/react.md
  - Next.js → .agent/framework-guides/nextjs.md
  - Express → .agent/framework-guides/express.md
- Python:
  - Django → .agent/framework-guides/django.md
  - FastAPI → .agent/framework-guides/fastapi.md
  - Flask → .agent/framework-guides/flask.md
- Go:
  - Gin → .agent/framework-guides/gin.md
  - Echo → .agent/framework-guides/echo.md
  - Fiber → .agent/framework-guides/fiber.md
- Rust:
  - Axum → .agent/framework-guides/axum.md
  - Actix-web → .agent/framework-guides/actix-web.md
  - Rocket → .agent/framework-guides/rocket.md
- Kotlin:
  - Spring Boot (Kotlin) → .agent/framework-guides/spring-boot-kotlin.md
  - Ktor → .agent/framework-guides/ktor.md
  - Android Compose → .agent/framework-guides/android-compose.md
- Java:
  - Spring Boot → .agent/framework-guides/spring-boot-java.md
  - Quarkus → .agent/framework-guides/quarkus.md
  - Micronaut → .agent/framework-guides/micronaut.md
- C#:
  - ASP.NET Core → .agent/framework-guides/aspnet-core.md
  - Blazor → .agent/framework-guides/blazor.md
  - Unity → .agent/framework-guides/unity.md
- PHP:
  - Laravel → .agent/framework-guides/laravel.md
  - Symfony → .agent/framework-guides/symfony.md
  - WordPress → .agent/framework-guides/wordpress.md
- Swift:
  - SwiftUI → .agent/framework-guides/swiftui.md
  - UIKit → .agent/framework-guides/uikit.md
  - Vapor → .agent/framework-guides/vapor.md
- Ruby:
  - Rails → .agent/framework-guides/rails.md
  - Sinatra → .agent/framework-guides/sinatra.md
  - Hanami → .agent/framework-guides/hanami.md
- Dart:
  - Flutter → .agent/framework-guides/flutter.md
  - Shelf → .agent/framework-guides/shelf.md
  - Dart Frog → .agent/framework-guides/dart-frog.md

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
3. **Develop**: Create `.agent/project.md` with plan → Implement incrementally.
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
1. Use `.agent/workflows/create-prd.md` to define requirements
2. Use `.agent/workflows/generate-tasks.md` to break down implementation
3. Implement tasks step-by-step with verification checkpoints
4. See `.agent/workflows/README.md` for full documentation

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
**Complex**: Create `.agent/project.md` OR use PRD workflow

**Checkpoints:**
- [ ] Requirements clear and testable
- [ ] Scope defined (what's included, what's not)
- [ ] Dependencies identified
- [ ] Breaking changes flagged

### Stage 2: Implementation
**Always:**
- Write tests first (TDD) or alongside code
- Load language-specific guide: .agent/language-guides/{language}.md
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
- Update `.agent/patterns.md` if new pattern emerged
- Update `.agent/state.md` with progress
- Add `.agent/memory/YYYY-MM-DD-topic.md` for key decisions

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
**Compatibility**: AGENTS.md standard (symlink for other tools)
**Current**: ~500 lines (target: <600, Operations section adds ~150 lines)

### .agent/ Directory
**Loaded**: On-demand, when needed
**Purpose**: Project-specific context that grows over time

**Structure:**
```
.agent/
├── README.md              # How to use .agent/
├── project.md             # Tech stack, architecture (create when tech chosen)
├── patterns.md            # Coding patterns (create when patterns emerge)
├── state.md              # Current work (create for multi-session work)
├── language-guides/      # Language-specific guardrails (pre-created, 21 languages)
│   ├── typescript.md     # TypeScript/JavaScript, React, Node.js
│   ├── python.md         # Python, Django, FastAPI
│   ├── go.md             # Go, microservices
│   ├── rust.md           # Rust, systems programming
│   ├── kotlin.md         # Kotlin, Android, Spring Boot
│   ├── java.md           # Java, Spring Boot, enterprise
│   ├── csharp.md         # C#, .NET, ASP.NET Core
│   ├── php.md            # PHP, Laravel, Symfony
│   ├── swift.md          # Swift, iOS, macOS
│   ├── cpp.md            # C/C++, systems, embedded
│   ├── ruby.md           # Ruby, Rails
│   ├── sql.md            # SQL, PostgreSQL, MySQL
│   ├── shell.md          # Shell/Bash scripting
│   ├── r.md              # R, statistical computing
│   ├── dart.md           # Dart, Flutter
│   ├── html-css.md       # HTML5, CSS3, accessibility
│   ├── lua.md            # Lua, Love2D, Neovim
│   ├── assembly.md       # x86-64, ARM64, RISC-V
│   ├── cuda.md           # CUDA, GPU computing
│   ├── solidity.md       # Solidity, Ethereum, smart contracts
│   └── zig.md            # Zig, systems programming
├── framework-guides/     # Framework-specific templates (pre-created, 33 frameworks)
│   ├── react.md          # React 18+, hooks, state management
│   ├── nextjs.md         # Next.js 14+, App Router, RSC
│   ├── express.md        # Express.js, middleware, REST APIs
│   ├── django.md         # Django 5+, ORM, admin, DRF
│   ├── fastapi.md        # FastAPI, async, Pydantic, OpenAPI
│   ├── flask.md          # Flask, blueprints, extensions
│   ├── gin.md            # Gin, middleware, REST APIs
│   ├── echo.md           # Echo, middleware, routing
│   ├── fiber.md          # Fiber, Express-style, high-performance
│   ├── axum.md           # Axum, Tower, async Rust
│   ├── actix-web.md      # Actix-web, actors, high-performance
│   ├── rocket.md         # Rocket, type-safe, macros
│   ├── spring-boot-kotlin.md  # Spring Boot with Kotlin, coroutines
│   ├── ktor.md           # Ktor, coroutines, DSL
│   ├── android-compose.md # Jetpack Compose, Material 3
│   ├── spring-boot.md    # Spring Boot Java, JPA, Security
│   ├── quarkus.md        # Quarkus, GraalVM, reactive
│   ├── micronaut.md      # Micronaut, compile-time DI
│   ├── aspnet-core.md    # ASP.NET Core, Minimal APIs, EF Core
│   ├── blazor.md         # Blazor, WebAssembly, SignalR
│   ├── unity.md          # Unity, C# scripting, game dev
│   ├── laravel.md        # Laravel 11+, Eloquent, Blade
│   ├── symfony.md        # Symfony 7+, Doctrine, Twig
│   ├── wordpress.md      # WordPress, themes, plugins, REST API
│   ├── swiftui.md        # SwiftUI, declarative UI, Combine
│   ├── uikit.md          # UIKit, programmatic/storyboard
│   ├── vapor.md          # Vapor, Fluent, async Swift
│   ├── rails.md          # Rails 7+, ActiveRecord, Hotwire
│   ├── sinatra.md        # Sinatra, lightweight Ruby
│   ├── hanami.md         # Hanami 2+, clean architecture
│   ├── flutter.md        # Flutter, Riverpod, go_router
│   ├── shelf.md          # Shelf, middleware HTTP server
│   └── dart-frog.md      # Dart Frog, file-based routing
├── workflows/            # Structured workflows (pre-created, 11 workflows)
│   ├── initialize-project.md  # Project setup
│   ├── create-prd.md          # Requirements documents
│   ├── generate-tasks.md      # Task breakdown
│   ├── code-review.md         # Pre-commit quality review
│   ├── security-audit.md      # Security assessment
│   ├── testing-strategy.md    # Test planning & coverage
│   ├── cleanup-project.md     # Prune unused guides
│   ├── refactoring.md         # Technical debt remediation
│   ├── dependency-update.md   # Safe dependency updates
│   ├── troubleshooting.md     # Debugging workflow
│   └── generate-agents-md.md  # Cross-tool compatibility
├── tasks/                # PRDs and task lists (created during COMPLEX mode)
│   ├── NNNN-prd-feature-name.md
│   └── tasks-NNNN-prd-feature-name.md
└── memory/               # Decision logs (created as needed)
    └── YYYY-MM-DD-topic.md
```

**Loading Protocol:**
- **Session Start**: AI loads CLAUDE.md → Checks for state.md → Reads if exists
- **During Work**: Load language guide based on file extensions (automatic)
- **Complex Features**: Load workflows (PRD, task generation)
- **Reference Needed**: Load patterns.md, project.md, memory/ on-demand

**Progressive Growth:**
- Day 1: Only CLAUDE.md + templates ✓
- First code: Language guide loaded automatically
- Week 1: Create `.agent/project.md` when architecture decided
- Month 1: `.agent/patterns.md` populated with conventions
- Ongoing: `.agent/memory/` captures complex decisions

**For details on any .agent/ file, see:** `.agent/README.md`

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

**For new projects:** Use `.agent/workflows/initialize-project.md`

**For existing projects:** Use `.agent/workflows/initialize-project.md`

AI will ask questions, analyze codebase, and create `.agent/project.md` with findings.

---

## When Stuck

**See:** `.agent/workflows/troubleshooting.md`

**Quick recovery:**
1. STOP trying random solutions (>30 min = stuck)
2. Document what you've tried
3. Simplify & isolate (minimal reproduction)
4. Check fundamentals (dependencies, config, versions)
5. Ask user with clear problem statement
6. Record solution in `.agent/memory/`

---

## Success Criteria

### Per Task
- [ ] All guardrails validated
- [ ] Tests pass with coverage thresholds
- [ ] No security vulnerabilities introduced
- [ ] Documentation updated
- [ ] Commit follows conventions

### Per Session
- [ ] State documented in `.agent/state.md` (if multi-session)
- [ ] New patterns added to `.agent/patterns.md` (if emerged)
- [ ] Key decisions in `.agent/memory/` (if significant)
- [ ] No broken tests left behind
- [ ] Progress measurable

---

## Version & Changelog

**Current Version**: 1.7.0
**Last Updated**: 2025-01-14

### Changelog

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
- ✅ Added framework-guides/ directory under .agent/
- ✅ Updated Quick Reference with "Load Framework Guide" section
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
- ✅ Extracted initialization to .agent/workflows/initialize-project.md
- ✅ Extracted troubleshooting to .agent/workflows/troubleshooting.md
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

**Remember**: This file is your guardrails. The `.agent/` directory is your memory. Small atomic changes. Validate continuously. Document progressively.

**Cross-Tool Compatibility**: `ln -s CLAUDE.md AGENTS.md` for tools that read AGENTS.md
