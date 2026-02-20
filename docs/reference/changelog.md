---
title: Changelog
description: Version history and updates
---

# Changelog

All notable changes to Samuel (Artificial Intelligence Coding Framework).

---

## [Unreleased]

### Added

- **7 Community Skills** from Anthropic (`github.com/anthropics/skills`):
  - `algorithmic-art` — Generative art with p5.js and seeded randomness
  - `doc-coauthoring` — Collaborative document writing workflow
  - `frontend-design` — Design-driven frontend interface creation
  - `mcp-builder` — MCP server development guide
  - `theme-factory` — Professional theme styling for artifacts
  - `web-artifacts-builder` — React/TypeScript/shadcn web application toolchain
  - `webapp-testing` — Playwright-based web application testing
- **`samuel auto pilot`** — Zero-setup autonomous mode that discovers and implements tasks without PRD setup
- **Go-native auto loop** — Autonomous loop now runs entirely in Go (replaces shell script)

### Changed

- Auto loop rewritten from shell script (`auto.sh`) to Go-native implementation
- `create-skill` workflow enhanced with Anthropic's skill-creator guidance

### Removed

- Shell-script auto loop (`auto.sh`)
- `DockerSandboxConfig` / `BuildDockerArgs` from docker.go (replaced by simpler sandbox configuration)

---

## [2.0.0] - 2026-02-12

### Renamed to Samuel

Breaking change: project renamed from AICoF to **Samuel**. Migrated from `.agent/` to `.claude/` native directory. Added autonomous AI coding loop.

### Added

- **Autonomous AI Coding Loop** (Ralph Wiggum methodology):
  - `samuel auto init` - Initialize autonomous loop from PRD
  - `samuel auto convert` - Convert markdown PRD/tasks to prd.json
  - `samuel auto status` - Show loop progress
  - `samuel auto start` - Begin/resume autonomous execution
  - `samuel auto task` - Manage tasks (list, complete, skip, reset, add)
- **Auto workflow skill** at `.claude/skills/auto/SKILL.md`
- **Per-folder CLAUDE.md** support for hierarchical instructions
- **Auto config keys**: `auto.enabled`, `auto.ai_tool`, `auto.max_iterations`, `auto.quality_checks`
- **AGENTS.md** now a real file copy (not symlink) for cross-tool compatibility
- **Homebrew formula** — `brew install samuel` (was cask)

### Changed

- **Project renamed**: `aicof` → `samuel` (binary, config, module path, repository)
- Config file: `aicof.yaml` → `samuel.yaml`
- Cache/config dirs: `~/.config/aicof/` → `~/.config/samuel/`
- Migrated from `.agent/` to `.claude/` (native Claude Code directory)
- Skills now live in `.claude/skills/` (native skill discovery)
- Homebrew distribution changed from cask to formula

### Breaking

- Binary renamed from `aicof` to `samuel`
- Config file renamed from `aicof.yaml` to `samuel.yaml`
- `.agent/` directory no longer used — migrate to `.claude/`
- Workflow files moved from `.claude/workflows/` to `.claude/skills/`
- Language guides moved from `.claude/language-guides/` to `.claude/skills/<name>-guide/`

---

## [1.8.0] - 2026-02-04

### Agent Skills Integration

Added Agent Skills support following the open standard for AI agent capabilities.

### Added

- **Agent Skills Management CLI**:
  - `samuel skill create <name>` - Scaffold new skills
  - `samuel skill validate [name]` - Validate against Agent Skills specification
  - `samuel skill list` - List installed skills
  - `samuel skill info <name>` - Show skill details
- **Create Skill workflow** at `.claude/skills/create-skill/SKILL.md`
- Skills compatible with 25+ agent products
- `installed.skills` config key

### Changed

- All components now follow Agent Skills standard (SKILL.md with YAML frontmatter)
- Language guides renamed to `<name>-guide/SKILL.md` format
- Framework guides moved to `<name>/SKILL.md` format
- Workflows moved to `<name>/SKILL.md` format

---

## [1.7.0] - 2026-01-14

### Rebrand and CLI Enhancements

Rebranded to Samuel (Artificial Intelligence Coding Framework) with expanded CLI capabilities.

### Added

- **Rebranded to Samuel** - clearer identity for the framework
- **AGENTS.md Compatible** - cross-tool support for 20+ AI coding tools

### Changed

- Updated branding throughout documentation and CLI
- Improved search scoring and fuzzy matching

---

## [1.6.0] - 2026-01-10

### CLI-First Workflow

Major expansion of the CLI tool with discovery-first workflow.

### Added

- **CLI Commands**: search, info, diff, config (bringing total to 11)
- Fuzzy search with relevance scoring
- Component preview before installing
- Version comparison and diffing

---

## [1.5.0] - 2025-12-15

### Framework Guides Release

Major expansion adding 33 framework-specific guides across 11 language families.

### Added

- **Framework Guides** (33 total):
  - **TypeScript/JavaScript**: React, Next.js, Express
  - **Python**: Django, FastAPI, Flask
  - **Go**: Gin, Echo, Fiber
  - **Rust**: Axum, Actix-web, Rocket
  - **Kotlin**: Spring Boot (Kotlin), Ktor, Android Compose
  - **Java**: Spring Boot, Quarkus, Micronaut
  - **C#**: ASP.NET Core, Blazor, Unity
  - **PHP**: Laravel, Symfony, WordPress
  - **Swift**: SwiftUI, UIKit, Vapor
  - **Ruby**: Rails, Sinatra, Hanami
  - **Dart**: Flutter, Shelf, Dart Frog

### Changed

- Added "Load Framework Guide" section to CLAUDE.md Quick Reference
- Added Framework Guides section to README.md
- Added frameworks section to documentation site
- Updated documentation site navigation

### Technical

- Total guides now 54 (21 language + 33 framework)
- Each framework guide includes: Setup, Project Structure, Patterns, Testing, Common Pitfalls, Security

---

## [1.4.0] - 2025-12-13

### Comprehensive Language Coverage Release

Major expansion adding 10 specialized domain language guides, bringing the total to 21 languages.

### Added

- **Specialized Domain Language Guides**:
  - SQL (PostgreSQL, MySQL, SQLite, query optimization)
  - Shell/Bash (scripting, automation, POSIX)
  - R (statistical computing, tidyverse, Shiny)
  - Dart (Flutter, mobile development)
  - HTML/CSS (web standards, accessibility, BEM)
  - Lua (scripting, Love2D, Neovim)
  - Assembly (x86-64, ARM64, RISC-V)
  - CUDA (GPU computing, parallel processing)
  - Solidity (Ethereum, smart contracts, DeFi)
  - Zig (systems programming, C interop)

### Changed

- Updated Quick Reference with all 21 language guides
- Updated `.claude/` directory structure documentation
- Updated documentation site with all language guides

### Technical

- Now supports 21 programming languages total
- All guides include: Core Principles, Guardrails, Validation, Testing, Tooling, Common Pitfalls

---

## [1.3.0] - 2025-12-13

### Extended Language Support Release

Major expansion adding 6 enterprise language guides covering top programming languages.

### Added

- **Enterprise Language Guides**:
  - Java (Spring Boot, Jakarta EE, JUnit 5)
  - C# (.NET 8, ASP.NET Core, xUnit)
  - PHP (Laravel, Symfony, PHPUnit)
  - Swift (iOS, macOS, SwiftUI)
  - C/C++ (C++20, CMake, GoogleTest)
  - Ruby (Rails 7, RSpec, Sidekiq)

### Changed

- Updated Quick Reference with all language guides
- Updated `.claude/` directory structure documentation
- Updated documentation site with enterprise language guides

### Technical

- Now supports 11 languages total (TypeScript, Python, Go, Rust, Kotlin + new)
- Each guide follows standard template: Core Principles → Guardrails → Validation → Testing → Tooling → Pitfalls

---

## [1.2.0] - 2025-12-13

### AGENTS.md Compatibility Release

Major update adding cross-tool compatibility through the AGENTS.md standard.

### Added

- **Operations Section** - Commands first, context second pattern
  - Setup Commands (Node, Python, Go, Rust)
  - Testing Commands with coverage targets
  - Build & Deploy Commands
  - Code Style Commands (format, lint)
  - Environment Variables template
- **Boundaries Section** - Prominent "do not touch" list
- **AGENTS.md Compatibility** - Works with 20+ AI coding tools
- **Kotlin Language Guide** - Full support for Kotlin development
- **AGENTS.md Generator Workflow** - Create standalone AGENTS.md

### Changed

- Restructured CLAUDE.md for "commands first, context second" pattern
- Added symlink instructions for cross-tool compatibility
- Consolidated Protected Boundaries into Boundaries section
- Updated documentation links

### Technical

- ~500 lines (up from 400, due to Operations section)
- Full AGENTS.md standard compliance
- Cross-tool verified with Claude Code, Cursor

---

## [1.1.0] - 2025-01-15

### Phase 1 Optimization Release

Focus on reducing token usage while adding critical guardrails.

### Added

- **Quick Reference Section** - Fast access to common information
- **Critical Guardrails**:
  - Dependency license checking
  - Database migration rollbacks
  - Semantic versioning for breaking changes
  - Smoke test validation for deployments
- **Language Guides**:
  - TypeScript/JavaScript (comprehensive)
  - Python (Django, FastAPI, Data Science)
  - Go (microservices, CLIs)
  - Rust (systems, WebAssembly)

### Changed

- Reduced CLAUDE.md from 490 → 400 lines (18% reduction)
- Clarified workflow requirements (MANDATORY vs RECOMMENDED)
- Extracted language-specific content to `.claude/language-guides/`
- Extracted initialization to `.claude/skills/initialize-project/SKILL.md`
- Extracted troubleshooting to `.claude/skills/troubleshooting/SKILL.md`

### Technical

- Token optimization for AI context windows
- Better separation of concerns
- Progressive loading protocol

---

## [1.0.0] - 2025-01-14

### Initial Release

First public release of Samuel (formerly AI Code Template).

### Added

- **CLAUDE.md** - Core instruction file (~490 lines)
  - 30+ guardrails for code quality, security, testing, git
  - 4D methodology (Deconstruct, Diagnose, Develop, Deliver)
  - 3 modes (ATOMIC, FEATURE, COMPLEX)
  - Software Development Lifecycle stages

- **.claude/ Directory Structure**
  - README.md for directory documentation
  - project.md.template for tech stack
  - state.md.template for work tracking
  - language-guides/ placeholder
  - workflows/ placeholder
  - tasks/ for PRDs
  - memory/ for decisions

- **Workflows**
  - create-prd.md - Product Requirements Documents
  - generate-tasks.md - Task breakdown
  - initialize-project.md - Project setup

### Philosophy

- Small, validated changes
- Quality enforced through guardrails
- Documentation grows organically
- Progressive disclosure of complexity

---

## Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR** (1.x.x → 2.0.0): Breaking changes to guardrails or methodology
- **MINOR** (1.0.x → 1.1.0): New features, language guides, workflows
- **PATCH** (1.0.0 → 1.0.1): Bug fixes, documentation improvements

---

## Upgrade Guide

### From 1.x to 2.0.0

**Breaking change**: Full migration from `.agent/` to `.claude/`.

1. Run `samuel update --version 2.0.0` (handles migration automatically)
2. Or manually:
   - Move `.agent/` contents to `.claude/skills/`
   - Delete `.agent/` directory
   - Update CLAUDE.md to v2.0.0 format
   - Copy AGENTS.md as a real file (not symlink)
3. Remove deprecated directories: `.agent/memory/`, `.agent/tasks/`, `.agent/rfd/`

### From 1.7.x to 1.8.x

1. New CLI commands available: `samuel skill create|validate|list|info`
2. Skills directory `.claude/skills/` now follows Agent Skills standard
3. Optional: Create custom skills with `samuel skill create <name>`

### From 1.4.x to 1.5.x

1. Copy new framework guides to `.claude/framework-guides/`:
   - All 33 framework guides (react.md, nextjs.md, express.md, django.md, etc.)
2. Update CLAUDE.md (includes "Load Framework Guide" section)
3. Copy new `docs/frameworks/` directory for documentation site

### From 1.3.x to 1.4.x

1. Copy new specialized language guides to `.claude/language-guides/`:
   - sql.md, shell.md, r.md, dart.md, html-css.md
   - lua.md, assembly.md, cuda.md, solidity.md, zig.md
2. Update CLAUDE.md (includes all 21 language guide references)

### From 1.2.x to 1.3.x

1. Copy new enterprise language guides to `.claude/language-guides/`:
   - java.md, csharp.md, php.md, swift.md, cpp.md, ruby.md
2. Update CLAUDE.md (includes all 11 language guide references)

### From 1.1.x to 1.2.x

1. Copy new CLAUDE.md (includes Operations section)
2. Add AGENTS.md symlink if needed: `ln -s CLAUDE.md AGENTS.md`
3. Copy new `.claude/skills/generate-agents-md/SKILL.md`
4. Copy new `.claude/language-guides/kotlin.md` (if using Kotlin)

### From 1.0.x to 1.1.x

1. Copy new CLAUDE.md
2. Copy new `.claude/language-guides/` directory
3. Copy new `.claude/skills/troubleshooting/SKILL.md`
4. Review Quick Reference section

---

## Future Plans

### Potential Features

- [x] ~~Additional language guides (Java, C#, PHP, Swift)~~ - Completed in v1.3.0
- [x] ~~Specialized language guides (SQL, Shell, Lua, etc.)~~ - Completed in v1.4.0
- [x] ~~Framework-specific templates (Next.js, Django, Rails)~~ - Completed in v1.5.0
- [x] ~~Agent Skills integration~~ - Completed in v1.8.0
- [x] ~~Autonomous AI coding loop~~ - Completed in v2.0.0
- [ ] IDE integrations (VS Code extension)
- [ ] Automated guardrail checking
- [ ] Web-based dashboard for auto loop monitoring
- [ ] Interactive documentation tutorials

### Community Contributions Welcome

- New language guides (Scala, Elixir, Haskell, Julia, etc.)
- Framework templates
- Integration examples
- Documentation improvements

---

## Links

- [GitHub Repository](https://github.com/ar4mirez/samuel)
- [Issues](https://github.com/ar4mirez/samuel/issues)
- [Discussions](https://github.com/ar4mirez/samuel/discussions)
