# AICoF - Artificial Intelligence Coding Framework

> **Build smarter, faster, and more scalable software**
> Cross-tool compatible • Opinionated guardrails • Tech-stack agnostic • Token-optimized

[![Version](https://img.shields.io/badge/version-1.7.0-blue.svg)](CLAUDE.md)
[![AGENTS.md](https://img.shields.io/badge/AGENTS.md-compatible-brightgreen.svg)](https://agents.md)
[![Status](https://img.shields.io/badge/status-production%20ready-brightgreen.svg)](CLAUDE.md)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

---

## What's New in v1.7.0

- **Rebranded to AICoF** - Artificial Intelligence Coding Framework with clearer identity
- **33 Framework Guides** - Comprehensive framework-specific guidance across 11 language families
- **21 Language Guides** - All major programming languages covered
- **AGENTS.md Compatible** - Works with Claude Code, Cursor, Codex, Copilot, and 20+ other AI tools

---

## Quick Start (60 Seconds)

### Option 1: Using the CLI (Recommended)

```bash
# Install CLI
brew tap ar4mirez/tap
brew install aicof

# Or via curl
curl -sSL https://raw.githubusercontent.com/ar4mirez/aicof/main/install.sh | sh

# Initialize your project
aicof init my-project
cd my-project

# Explore available components
aicof search react           # Find components by keyword
aicof info framework react   # View component details
aicof list --available       # List all components

# Start coding with AI - guardrails apply automatically!
```

### Option 2: Manual Copy

```bash
# 1. Copy template files to your project
cp -r /path/to/aicof/template/{CLAUDE.md,AI_INSTRUCTIONS.md,.agent} ./

# 2. (Optional) For cross-tool compatibility
ln -s CLAUDE.md AGENTS.md

# 3. Start coding with AI - guardrails apply automatically!
```

**The system works immediately:**

- AI loads [CLAUDE.md](template/CLAUDE.md) automatically (500 lines of guardrails + operations)
- Language guides auto-load based on file extensions
- Workflows available when you need them
- Progressive - starts minimal, grows with your project

**[Read the full Quick Start Guide](template/AI_INSTRUCTIONS.md)**

---

## Cross-Tool Compatibility (AGENTS.md)

This system follows the [AGENTS.md](https://agents.md) standard - the universal format for AI coding assistants adopted by 20,000+ repositories.

### How It Works

| Tool | Primary File | Fallback |
|------|--------------|----------|
| **Claude Code** | CLAUDE.md | AGENTS.md |
| **Cursor** | AGENTS.md | CLAUDE.md |
| **OpenAI Codex** | AGENTS.md | CLAUDE.md |
| **GitHub Copilot** | AGENTS.md | CLAUDE.md |
| **Google Jules** | AGENTS.md | CLAUDE.md |

### Setup for Cross-Tool Teams

**Option 1: Symlink (recommended)**
```bash
ln -s CLAUDE.md AGENTS.md
```

**Option 2: Generate standalone AGENTS.md**
```
@.agent/workflows/generate-agents-md.md
```

**Why both files?**
- `CLAUDE.md` = Full methodology (guardrails + 4D + workflows)
- `AGENTS.md` = Operations only (commands, boundaries, style)

---

## What Is AICoF?

AICoF (Artificial Intelligence Coding Framework) is an **opinionated AI development framework** designed for professional software teams.

### Key Features

| Feature | Description |
|---------|-------------|
| **35+ Guardrails** | Testable rules, not vague suggestions |
| **21 Language Guides** | All major languages with auto-loading support |
| **33 Framework Guides** | Framework-specific patterns and best practices |
| **13 Workflows** | PRD, tasks, init, troubleshooting, code-review, and more |
| **3 Modes** | ATOMIC/FEATURE/COMPLEX (scales from bugs to architecture) |
| **4D Methodology** | Deconstruct → Diagnose → Develop → Deliver |
| **Cross-Tool** | Works with any AI coding assistant |

### Philosophy

> Small, validated changes. Quality enforced. Documentation grows organically.

---

## CLI Commands

The `aicof` CLI manages framework installation, updates, and component discovery.

### Core Commands

| Command | Description | Example |
|---------|-------------|---------|
| `init [project]` | Initialize AICoF in a project | `aicof init my-app` |
| `update` | Update to latest framework version | `aicof update` |
| `doctor` | Check installation health | `aicof doctor` |
| `version` | Show CLI and framework versions | `aicof version` |

### Component Management

| Command | Description | Example |
|---------|-------------|---------|
| `add <type> <name>` | Add a component | `aicof add framework react` |
| `remove <type> <name>` | Remove a component | `aicof remove language rust` |
| `list [--available]` | List installed/available components | `aicof list --available` |

**Type aliases**: `language` (lang, l), `framework` (fw, f), `workflow` (wf, w)

### Discovery Commands

| Command | Description | Example |
|---------|-------------|---------|
| `search <query>` | Search components by keyword | `aicof search api` |
| `info <type> <name>` | Show component details | `aicof info fw nextjs` |
| `diff [v1] [v2]` | Compare versions | `aicof diff v1.6.0 v1.7.0` |

### Configuration

| Command | Description | Example |
|---------|-------------|---------|
| `config list` | Show all config values | `aicof config list` |
| `config get <key>` | Get a config value | `aicof config get version` |
| `config set <key> <value>` | Set a config value | `aicof config set registry https://...` |

**Valid keys**: `version`, `registry`, `installed.languages`, `installed.frameworks`, `installed.workflows`

### Command Examples

```bash
# Initialize and customize
aicof init my-project
aicof add lang typescript
aicof add fw react nextjs
aicof add wf code-review security-audit

# Discover components
aicof search python           # Fuzzy search across all types
aicof search --type fw api    # Search only frameworks
aicof info lang go --preview 20  # Preview first 20 lines

# Manage updates
aicof diff                    # Compare installed vs latest
aicof diff v1.6.0 v1.7.0      # Compare specific versions
aicof update                  # Apply updates

# Troubleshoot
aicof doctor                  # Check installation health
aicof config list             # View current configuration
```

---

## Professional Prompt Examples

### Bug Fixes (ATOMIC Mode)

```
Fix the null pointer exception in UserService.getProfile() when user.email is undefined
```

```
The checkout button is disabled after form validation passes - debug and fix
```

```
Memory leak in useWebSocket hook - component doesn't cleanup on unmount
```

### Feature Development (FEATURE Mode)

```
Add pagination to the /api/users endpoint with cursor-based navigation
```

```
Implement email verification flow: send verification link, validate token, update user status
```

```
Create a reusable DataTable component with sorting, filtering, and column resizing
```

### Complex Features (COMPLEX Mode)

```
@.agent/workflows/create-prd.md
Build a real-time notification system with WebSocket support, push notifications,
and user preference management
```

```
@.agent/workflows/create-prd.md
Implement multi-tenant architecture: tenant isolation, data partitioning,
tenant-specific configurations
```

```
@.agent/workflows/create-prd.md
Add comprehensive audit logging: user actions, data changes, security events,
with retention policies and export functionality
```

### Code Review & Analysis

```
Review the authentication module for security vulnerabilities and suggest improvements
```

```
Analyze the database query patterns in OrderService and identify N+1 query issues
```

```
Evaluate the error handling strategy in the API layer and propose a consistent approach
```

### Refactoring

```
Refactor UserController to use dependency injection and separate business logic into a service layer
```

```
Extract common validation logic from all form components into a reusable validation hook
```

```
Migrate the codebase from callbacks to async/await while maintaining backward compatibility
```

### Architecture & Planning

```
@.agent/workflows/create-prd.md
Design a caching strategy for the product catalog: cache invalidation,
distributed caching, cache warming
```

```
Analyze the current monolith and propose a microservices decomposition plan
with clear service boundaries
```

```
Create a database migration strategy for adding multi-region support
without downtime
```

### Debugging & Troubleshooting

```
@.agent/workflows/troubleshooting.md
Production error: "Connection pool exhausted" appearing intermittently under load
```

```
Performance degradation: API response times increased from 50ms to 500ms
after last deployment
```

```
Investigate why the CI pipeline is flaky - tests pass locally but fail
randomly in GitHub Actions
```

### Testing

```
Generate comprehensive unit tests for the PaymentService with edge cases
for failed transactions, refunds, and currency conversion
```

```
Create integration tests for the OAuth flow covering Google, GitHub,
and email/password authentication
```

```
Add E2E tests for the checkout flow using Playwright, including
error scenarios and payment failures
```

### Documentation

```
Generate API documentation for the /api/v2/orders endpoint including
request/response examples and error codes
```

```
Create a technical design document for the upcoming search feature
with architecture diagrams
```

### DevOps & Infrastructure

```
Create a Dockerfile for the Node.js API with multi-stage build,
non-root user, and health checks
```

```
Write GitHub Actions workflow for: lint, test, build, deploy to staging
on PR, deploy to production on merge
```

---

## How It Works

### The 3 Modes

**ATOMIC** (<5 files, clear scope)
```
"Fix the login button alignment"
```
- Direct implementation
- Quick validation
- One commit
- Examples: Bug fixes, styling, simple features

**FEATURE** (5-10 files)
```
"Add user profile editing with avatar upload"
```
- Break into 3-5 subtasks
- Implement sequentially
- Integration testing
- Examples: New component, API endpoint, refactoring

**COMPLEX** (>10 files, new subsystem)
```
@.agent/workflows/create-prd.md
"Build user authentication with OAuth"
```
- Create PRD (Product Requirements Document)
- Generate task breakdown
- Step-by-step implementation
- Examples: Authentication, payments, analytics

**AI auto-detects which mode to use.**

### The Guardrails (35+ Rules)

**Code Quality:**
- Functions ≤50 lines
- Files ≤300 lines
- Complexity ≤10 per function
- All exports have types/docs

**Security (CRITICAL):**
- All inputs validated
- Parameterized queries only
- No secrets in code
- Dependencies checked for vulnerabilities + licenses

**Testing (CRITICAL):**
- >80% coverage for business logic
- >60% overall coverage
- Tests for all public APIs
- Regression tests for bugs

**Git:**
- Conventional commits (`feat:`, `fix:`, etc.)
- One logical change per commit
- All tests pass before push
- PRs required (no direct commits to main)

**[See all guardrails in CLAUDE.md](CLAUDE.md)**

### The 4D Methodology

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ DECONSTRUCT │ ──▶ │  DIAGNOSE   │ ──▶ │   DEVELOP   │ ──▶ │   DELIVER   │
│             │     │             │     │             │     │             │
│ Break down  │     │ Identify    │     │ Implement   │     │ Validate    │
│ the task    │     │ risks &     │     │ with tests  │     │ & commit    │
│             │     │ dependencies│     │             │     │             │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
```

---

## Use Cases

### Starting a New Project

```
@.agent/workflows/initialize-project.md
"Initialize a new TypeScript API with Express, PostgreSQL, and Jest"
```

AI will:
1. Ask clarifying questions (architecture, deployment target)
2. Create directory structure
3. Generate config files (tsconfig, package.json, etc.)
4. Create `.agent/project.md` documenting decisions
5. Set up testing framework

### Onboarding to Existing Project

```
@.agent/workflows/initialize-project.md
"This is an existing project - analyze the codebase and document patterns"
```

AI will:
1. Scan tech stack (package.json, requirements.txt, etc.)
2. Analyze code patterns and conventions
3. Review recent commits
4. Create `.agent/project.md` with findings
5. Extract patterns to `.agent/patterns.md`

### Building Complex Features

```
@.agent/workflows/create-prd.md
"Build a subscription billing system with Stripe integration"
```

AI will:
1. Ask clarifying questions (plans, trial periods, webhooks)
2. Create comprehensive PRD
3. Generate task breakdown (20-50 subtasks)
4. Implement step-by-step with verification
5. Update documentation

### Cross-Tool Team Workflow

```
@.agent/workflows/generate-agents-md.md
"Generate AGENTS.md for our team using Cursor and Claude Code"
```

AI will:
1. Extract Operations section from CLAUDE.md
2. Generate universal AGENTS.md
3. Ensure both tools work consistently

---

## Documentation

| Document | Purpose | When to Read |
|----------|---------|--------------|
| [AI_INSTRUCTIONS.md](template/AI_INSTRUCTIONS.md) | **Quick Start Guide** | Read this first |
| [CLAUDE.md](template/CLAUDE.md) | Core guardrails & methodology | AI loads automatically |
| [.agent/README.md](template/.agent/README.md) | .agent/ folder structure | When customizing |

### Language Guides (Auto-Load)

| Language | Files | Guide |
|----------|-------|-------|
| TypeScript/JavaScript | `.ts`, `.tsx`, `.js`, `.jsx` | [SKILL.md](template/.agent/skills/typescript-guide/SKILL.md) |
| Python | `.py` | [SKILL.md](template/.agent/skills/python-guide/SKILL.md) |
| Go | `.go` | [SKILL.md](template/.agent/skills/go-guide/SKILL.md) |
| Rust | `.rs` | [SKILL.md](template/.agent/skills/rust-guide/SKILL.md) |
| Kotlin | `.kt`, `.kts` | [SKILL.md](template/.agent/skills/kotlin-guide/SKILL.md) |
| Java | `.java` | [SKILL.md](template/.agent/skills/java-guide/SKILL.md) |
| C# | `.cs` | [SKILL.md](template/.agent/skills/csharp-guide/SKILL.md) |
| PHP | `.php` | [SKILL.md](template/.agent/skills/php-guide/SKILL.md) |
| Swift | `.swift` | [SKILL.md](template/.agent/skills/swift-guide/SKILL.md) |
| C/C++ | `.c`, `.cpp`, `.h`, `.hpp` | [SKILL.md](template/.agent/skills/cpp-guide/SKILL.md) |
| Ruby | `.rb` | [SKILL.md](template/.agent/skills/ruby-guide/SKILL.md) |
| SQL | `.sql` | [SKILL.md](template/.agent/skills/sql-guide/SKILL.md) |
| Shell/Bash | `.sh`, `.bash` | [SKILL.md](template/.agent/skills/shell-guide/SKILL.md) |
| R | `.r`, `.R` | [SKILL.md](template/.agent/skills/r-guide/SKILL.md) |
| Dart/Flutter | `.dart` | [SKILL.md](template/.agent/skills/dart-guide/SKILL.md) |
| HTML/CSS | `.html`, `.css`, `.scss` | [SKILL.md](template/.agent/skills/html-css-guide/SKILL.md) |
| Lua | `.lua` | [SKILL.md](template/.agent/skills/lua-guide/SKILL.md) |
| Assembly | `.asm`, `.s` | [SKILL.md](template/.agent/skills/assembly-guide/SKILL.md) |
| CUDA | `.cu`, `.cuh` | [SKILL.md](template/.agent/skills/cuda-guide/SKILL.md) |
| Solidity | `.sol` | [SKILL.md](template/.agent/skills/solidity-guide/SKILL.md) |
| Zig | `.zig` | [SKILL.md](template/.agent/skills/zig-guide/SKILL.md) |

### Framework Guides (On-Demand)

| Language | Frameworks |
|----------|------------|
| TypeScript/JS | [React](template/.agent/framework-guides/react.md), [Next.js](template/.agent/framework-guides/nextjs.md), [Express](template/.agent/framework-guides/express.md) |
| Python | [Django](template/.agent/framework-guides/django.md), [FastAPI](template/.agent/framework-guides/fastapi.md), [Flask](template/.agent/framework-guides/flask.md) |
| Go | [Gin](template/.agent/framework-guides/gin.md), [Echo](template/.agent/framework-guides/echo.md), [Fiber](template/.agent/framework-guides/fiber.md) |
| Rust | [Axum](template/.agent/framework-guides/axum.md), [Actix-web](template/.agent/framework-guides/actix-web.md), [Rocket](template/.agent/framework-guides/rocket.md) |
| Kotlin | [Spring Boot](template/.agent/framework-guides/spring-boot-kotlin.md), [Ktor](template/.agent/framework-guides/ktor.md), [Android Compose](template/.agent/framework-guides/android-compose.md) |
| Java | [Spring Boot](template/.agent/framework-guides/spring-boot-java.md), [Quarkus](template/.agent/framework-guides/quarkus.md), [Micronaut](template/.agent/framework-guides/micronaut.md) |
| C# | [ASP.NET Core](template/.agent/framework-guides/aspnet-core.md), [Blazor](template/.agent/framework-guides/blazor.md), [Unity](template/.agent/framework-guides/unity.md) |
| PHP | [Laravel](template/.agent/framework-guides/laravel.md), [Symfony](template/.agent/framework-guides/symfony.md), [WordPress](template/.agent/framework-guides/wordpress.md) |
| Swift | [SwiftUI](template/.agent/framework-guides/swiftui.md), [UIKit](template/.agent/framework-guides/uikit.md), [Vapor](template/.agent/framework-guides/vapor.md) |
| Ruby | [Rails](template/.agent/framework-guides/rails.md), [Sinatra](template/.agent/framework-guides/sinatra.md), [Hanami](template/.agent/framework-guides/hanami.md) |
| Dart | [Flutter](template/.agent/framework-guides/flutter.md), [Shelf](template/.agent/framework-guides/shelf.md), [Dart Frog](template/.agent/framework-guides/dart-frog.md) |

### Workflows (On-Demand)

| Workflow | Purpose |
|----------|---------|
| [Initialize Project](template/.agent/workflows/initialize-project.md) | Setup new/existing projects |
| [Create PRD](template/.agent/workflows/create-prd.md) | Plan complex features |
| [Generate Tasks](template/.agent/workflows/generate-tasks.md) | Break PRDs into tasks |
| [Code Review](template/.agent/workflows/code-review.md) | Pre-commit quality review |
| [Security Audit](template/.agent/workflows/security-audit.md) | Security assessment |
| [Testing Strategy](template/.agent/workflows/testing-strategy.md) | Test planning & coverage |
| [Refactoring](template/.agent/workflows/refactoring.md) | Technical debt remediation |
| [Dependency Update](template/.agent/workflows/dependency-update.md) | Safe dependency updates |
| [Troubleshooting](template/.agent/workflows/troubleshooting.md) | Debug systematically |
| [Cleanup Project](template/.agent/workflows/cleanup-project.md) | Prune unused guides |
| [Document Work](template/.agent/workflows/document-work.md) | Capture patterns & decisions |
| [Update Framework](template/.agent/workflows/update-framework.md) | Update AICoF safely |
| [Generate AGENTS.md](template/.agent/workflows/generate-agents-md.md) | Cross-tool compatibility |

---

## Repository Structure

```text
aicof/
├── template/                    # Distributable template files
│   ├── CLAUDE.md               # Main AI instructions (copy to your project)
│   ├── AI_INSTRUCTIONS.md      # Quick start guide
│   └── .agent/                 # AI context directory
│       ├── skills/             # 21 language-specific guides
│       ├── framework-guides/   # 33 framework-specific guides
│       └── workflows/          # 13 structured workflows
├── cmd/aicof/                  # CLI entry point
├── internal/                   # CLI implementation (commands, core, ui)
├── docs/                       # Documentation website source
└── .agent/                     # This project's own AI context
    ├── tasks/                  # PRDs and task lists
    └── memory/                 # Decision logs
```

**Why this structure?**

- `template/` contains files distributed to users via the CLI
- `packages/cli/` contains the CLI tool that manages installations
- `.agent/` at root is for developing AICoF itself (dogfooding)

---

## System Stats

| Metric | Value |
|--------|-------|
| **Version** | 1.7.0 |
| **Status** | Production Ready |
| **AGENTS.md** | Compatible |
| **Total Files** | 67 markdown files |
| **CLAUDE.md** | ~500 lines |
| **Language Guides** | 21 (all major programming languages) |
| **Framework Guides** | 33 (across 11 language families) |
| **Workflows** | 13 (PRD, tasks, init, troubleshoot, code-review, etc.) |
| **Guardrails** | 35+ testable rules |

---

## Learning Path

### Week 1: Learn the Basics

- [ ] Install AICoF via CLI or copy template files to your project
- [ ] Write 5 features using ATOMIC mode
- [ ] Review guardrails in [CLAUDE.md](template/CLAUDE.md)
- [ ] Check which language guide applies to you

### Week 2: Try Complex Features

- [ ] Use PRD workflow for a medium feature
- [ ] Generate task breakdown
- [ ] Implement step-by-step
- [ ] Notice how `.agent/project.md` grows

### Week 3: Customize & Extend

- [ ] Add project-specific patterns to `.agent/patterns.md`
- [ ] Create first decision log in `.agent/memory/`
- [ ] (Multi-tool) Set up AGENTS.md for your team
- [ ] Experiment with different modes

---

## FAQ

### Q: Do I need to use Claude Code?
**A:** No! While designed for Claude Code, the system works with any AI coding assistant. Use the AGENTS.md symlink for other tools like Cursor or Codex.

### Q: Which file should I edit - CLAUDE.md or AGENTS.md?
**A:** Edit CLAUDE.md (the source of truth). If using symlink, AGENTS.md updates automatically. If using standalone AGENTS.md, regenerate it after CLAUDE.md changes.

### Q: My language isn't covered - what do I do?
**A:** The core guardrails in CLAUDE.md are language-agnostic (90% still applies). Consider contributing a new language guide!

### Q: Is the PRD workflow required?
**A:** No! Only for complex features (>10 files, new subsystems). Most work is ATOMIC or FEATURE mode.

### Q: Can I customize the guardrails?
**A:** Yes! Edit CLAUDE.md for your team. Common customizations:
- File length limits (300 → 500)
- Coverage targets (80% → 90%)
- Commit format (add your types)

---

## Contributing

Contributions welcome! Ideas:
- Additional language guides (Scala, Elixir, Haskell, OCaml, Julia)
- Additional framework guides for existing languages
- Integration examples with other AI tools
- Real-world case studies

**To contribute:**
1. Fork this repository
2. Create a feature branch
3. Make your changes following the guardrails
4. Submit a pull request

---

## License

MIT License - See [LICENSE](LICENSE) file for details.

---

## Acknowledgments

Built with:
- **Claude Code** - Anthropic's AI coding assistant
- **AGENTS.md Standard** - Universal AI agent instructions
- **4D Methodology** - Systematic problem-solving approach
- **Community feedback** - Continuous improvement

---

## Support

- **Documentation**: [AI_INSTRUCTIONS.md](AI_INSTRUCTIONS.md)
- **Issues**: [GitHub Issues](https://github.com/ar4mirez/aicof/issues)
- **Discussions**: [GitHub Discussions](https://github.com/ar4mirez/aicof/discussions)

---

**Happy coding with AI!**

*Works with Claude Code, Cursor, Codex, Copilot, and any AI assistant that reads AGENTS.md*
