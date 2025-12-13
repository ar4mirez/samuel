---
title: Changelog
description: Version history and updates
---

# Changelog

All notable changes to AI Code Template.

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
- Updated `.agent/` directory structure documentation
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
- Updated `.agent/` directory structure documentation
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
- Extracted language-specific content to `.agent/language-guides/`
- Extracted initialization to `.agent/workflows/initialize-project.md`
- Extracted troubleshooting to `.agent/workflows/troubleshooting.md`

### Technical

- Token optimization for AI context windows
- Better separation of concerns
- Progressive loading protocol

---

## [1.0.0] - 2025-01-14

### Initial Release

First public release of AI Code Template.

### Added

- **CLAUDE.md** - Core instruction file (~490 lines)
  - 30+ guardrails for code quality, security, testing, git
  - 4D methodology (Deconstruct, Diagnose, Develop, Deliver)
  - 3 modes (ATOMIC, FEATURE, COMPLEX)
  - Software Development Lifecycle stages

- **.agent/ Directory Structure**
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

### From 1.3.x to 1.4.x

1. Copy new specialized language guides to `.agent/language-guides/`:
   - sql.md, shell.md, r.md, dart.md, html-css.md
   - lua.md, assembly.md, cuda.md, solidity.md, zig.md
2. Update CLAUDE.md (includes all 21 language guide references)

### From 1.2.x to 1.3.x

1. Copy new enterprise language guides to `.agent/language-guides/`:
   - java.md, csharp.md, php.md, swift.md, cpp.md, ruby.md
2. Update CLAUDE.md (includes all 11 language guide references)

### From 1.1.x to 1.2.x

1. Copy new CLAUDE.md (includes Operations section)
2. Add AGENTS.md symlink if needed: `ln -s CLAUDE.md AGENTS.md`
3. Copy new `.agent/workflows/generate-agents-md.md`
4. Copy new `.agent/language-guides/kotlin.md` (if using Kotlin)

### From 1.0.x to 1.1.x

1. Copy new CLAUDE.md
2. Copy new `.agent/language-guides/` directory
3. Copy new `.agent/workflows/troubleshooting.md`
4. Review Quick Reference section

---

## Future Plans

### Potential Features

- [x] ~~Additional language guides (Java, C#, PHP, Swift)~~ - Completed in v1.3.0
- [x] ~~Specialized language guides (SQL, Shell, Lua, etc.)~~ - Completed in v1.4.0
- [ ] Framework-specific templates (Next.js, Django, Rails)
- [ ] IDE integrations (VS Code extension)
- [ ] Automated guardrail checking
- [ ] Interactive documentation tutorials
- [ ] Project scaffold generators

### Community Contributions Welcome

- New language guides (Scala, Elixir, Haskell, Julia, etc.)
- Framework templates
- Integration examples
- Documentation improvements

---

## Links

- [GitHub Repository](https://github.com/ar4mirez/ai-code-template)
- [Issues](https://github.com/ar4mirez/ai-code-template/issues)
- [Discussions](https://github.com/ar4mirez/ai-code-template/discussions)
