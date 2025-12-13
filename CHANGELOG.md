# Changelog

All notable changes to the AI Claude Code system will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

<!-- Add unreleased changes here -->

## [1.2.0] - 2025-12-13

### Added

#### AGENTS.md Compatibility
- **Operations section** in CLAUDE.md (AGENTS.md standard compatible)
  - Setup Commands for Node.js, Python, Go, Rust
  - Testing Commands with coverage targets
  - Build & Deploy Commands
  - Code Style Commands (format, lint)
  - Environment Variables template
- **Boundaries section** - prominent "do not touch" list for AI agents
- **Cross-tool symlink instructions** - `ln -s CLAUDE.md AGENTS.md`
- **generate-agents-md.md workflow** - generates standalone AGENTS.md for cross-tool teams

#### Documentation Site
- **MkDocs Material documentation site** with GitHub Pages
  - Light/dark mode toggle with system preference detection
  - Full-text search with highlighting
  - Responsive design for mobile and desktop
  - Code syntax highlighting with copy button
  - Mermaid diagram support
- **Comprehensive documentation structure**
  - Getting Started: quick-start, installation, first-task guides
  - Core System: overview, CLAUDE.md, methodology, guardrails, .agent directory
  - Language Guides: TypeScript, Python, Go, Rust, Kotlin
  - Workflows: initialize, create-prd, generate-tasks, troubleshooting, generate-agents-md
  - Reference: cross-tool compatibility, FAQ, changelog, contributing
- **GitHub Actions workflow** for automatic deployment to GitHub Pages
- **requirements-docs.txt** for documentation dependencies

#### Documentation
- **Professional prompt examples** in README.md (25+ real-world prompts)
  - Bug fixes (ATOMIC mode)
  - Feature development (FEATURE mode)
  - Complex features (COMPLEX mode)
  - Code review & analysis
  - Refactoring
  - Architecture & planning
  - Debugging & troubleshooting
  - Testing
  - Documentation
  - DevOps & infrastructure
- **4D Methodology visual diagram** in README.md
- **Cross-tool compatibility section** with tool priority table
- **Use cases section** with detailed examples

#### Language Guides
- **Kotlin language guide** (.agent/language-guides/kotlin.md)
  - Null safety patterns
  - Coroutines and structured concurrency
  - Data classes and sealed classes
  - Android and Spring Boot patterns
  - Kotest and MockK testing

### Changed

- **CLAUDE.md structure** - Operations section now at top (commands first, context second)
- **README.md** - Complete rewrite with professional focus
- **AI_INSTRUCTIONS.md** - Updated for v1.2.0 with cross-tool info
- **.agent/README.md** - Added AGENTS.md compatibility section
- **.agent/workflows/README.md** - Added generate-agents-md workflow

### Removed

- **Duplicate "Protected Boundaries" section** - consolidated into "Boundaries" at top of CLAUDE.md

## [1.1.0] - 2025-01-15

### Added

#### Guardrails
- Dependency license checking before adding packages
- Database migration rollback (down) function requirement
- Semantic versioning for breaking API changes
- Smoke test validation for deployments

#### Language Guides
- **TypeScript/JavaScript guide** (.agent/language-guides/typescript.md)
  - Strict mode configuration
  - Zod validation patterns
  - React hooks and functional components
  - Async/await best practices
- **Python guide** (.agent/language-guides/python.md)
  - Type hints (PEP 484)
  - Pydantic validation
  - Django and FastAPI patterns
  - pytest testing
- **Go guide** (.agent/language-guides/go.md)
  - Error handling patterns
  - Concurrency with goroutines
  - Table-driven tests
  - Interface design
- **Rust guide** (.agent/language-guides/rust.md)
  - Ownership and borrowing
  - Result<T,E> error handling
  - Async patterns with Tokio
  - Zero-cost abstractions

#### Workflows
- **initialize-project.md** - New/existing project setup workflow
- **troubleshooting.md** - Systematic debugging guide

#### Documentation
- **Quick Reference section** in CLAUDE.md
- **Line number references** for quick navigation

### Changed

- **CLAUDE.md size** - Reduced from 490 → 400 lines (18% reduction)
- **Workflow requirements** - Clarified MANDATORY vs RECOMMENDED vs OPTIONAL
- **Language-specific content** - Extracted to .agent/language-guides/
- **Initialization content** - Extracted to .agent/workflows/initialize-project.md
- **Troubleshooting content** - Extracted to .agent/workflows/troubleshooting.md

### Improved

- Token efficiency - ~1,600 tokens base context (vs. 4,000 before)
- Auto-loading - Language guides load based on file extensions
- Progressive loading - .agent/ files loaded on-demand

## [1.0.0] - 2025-01-14

### Added

#### Core System
- **CLAUDE.md** - Core guardrails file with 30+ rules
- **4D Methodology** - Deconstruct → Diagnose → Develop → Deliver
- **3 Modes** - ATOMIC, FEATURE, COMPLEX task classification
- **5 SDLC Stages** - Planning, Implementation, Validation, Documentation, Commit

#### Guardrails
- **Code Quality** (8 rules)
  - Function length ≤50 lines
  - File length ≤300 lines
  - Cyclomatic complexity ≤10
  - Type signatures and documentation
  - No magic numbers
  - No commented-out code
  - No TODO without issue reference
  - No dead code
- **Security** (6 rules)
  - Input validation
  - Parameterized queries
  - No hardcoded secrets
  - Path validation
  - Timeout/cancellation for async
  - Dependency vulnerability checking
- **Testing** (8 rules)
  - Coverage targets (80% business logic, 60% overall)
  - Public API tests
  - Regression tests for bugs
  - Edge case testing
  - Descriptive test names
  - No test interdependencies
  - Integration tests for external services
  - Smoke test validation
- **Git** (7 rules)
  - Conventional commits
  - Atomic commits
  - Tests pass before push
  - Branch naming conventions
  - No direct commits to main
  - PR required
  - Semantic versioning
- **Performance** (5 rules)
  - No N+1 queries
  - Pagination for large datasets
  - Memoization/caching
  - Bundle size limits
  - API response time limits

#### Workflows
- **create-prd.md** - Product Requirements Document workflow
- **generate-tasks.md** - Task breakdown from PRDs

#### Context System
- **.agent/ directory structure**
  - project.md - Tech stack and architecture
  - patterns.md - Coding patterns and conventions
  - state.md - Current work tracking
  - memory/ - Decision logs
  - tasks/ - PRDs and task lists
- **Progressive growth philosophy** - Start minimal, grow organically
- **On-demand loading** - Load context only when needed

#### Documentation
- **README.md** - Project overview and quick start
- **AI_INSTRUCTIONS.md** - Quick start guide
- **.agent/README.md** - Context system documentation
- **.agent/workflows/README.md** - Workflow documentation
- **Example PRD** - .agent/tasks/EXAMPLE-0001-prd-api-rate-limiting.md

### Technical Details

- **Token optimized** - Designed for efficient context usage
- **Tech-stack agnostic** - Works with any language/framework
- **AI-first design** - Structured for AI coding assistants

---

## Version History Summary

| Version | Date | Highlights |
|---------|------|------------|
| 1.2.0 | 2025-12-13 | AGENTS.md compatibility, MkDocs documentation site, Kotlin guide |
| 1.1.0 | 2025-01-15 | 4 language guides, 2 new workflows, 18% size reduction |
| 1.0.0 | 2025-01-14 | Initial release with 30+ guardrails, 4D methodology |

---

## Migration Guide

### From 1.1.0 to 1.2.0

No breaking changes. To enable cross-tool compatibility:

```bash
# Option 1: Symlink (recommended)
ln -s CLAUDE.md AGENTS.md

# Option 2: Generate standalone AGENTS.md
# Use @.agent/workflows/generate-agents-md.md
```

### From 1.0.0 to 1.1.0

No breaking changes. New language guides and workflows are additive.

---

## Links

- [CLAUDE.md](CLAUDE.md) - Core guardrails
- [README.md](README.md) - Project overview
- [AI_INSTRUCTIONS.md](AI_INSTRUCTIONS.md) - Quick start guide
- [AGENTS.md Standard](https://agents.md) - Cross-tool compatibility standard

[Unreleased]: https://github.com/ar4mirez/ai-claude-code/compare/v1.2.0...HEAD
[1.2.0]: https://github.com/ar4mirez/ai-claude-code/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/ar4mirez/ai-claude-code/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/ar4mirez/ai-claude-code/releases/tag/v1.0.0
