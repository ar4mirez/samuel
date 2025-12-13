---
title: Changelog
description: Version history and updates
---

# Changelog

All notable changes to AI Code Template.

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

- [ ] Additional language guides (Java, C#, PHP, Swift)
- [ ] Framework-specific templates (Next.js, Django, Rails)
- [ ] IDE integrations (VS Code extension)
- [ ] Automated guardrail checking

### Community Contributions Welcome

- New language guides
- Framework templates
- Integration examples
- Documentation improvements

---

## Links

- [GitHub Repository](https://github.com/ar4mirez/ai-code-template)
- [Issues](https://github.com/ar4mirez/ai-code-template/issues)
- [Discussions](https://github.com/ar4mirez/ai-code-template/discussions)
