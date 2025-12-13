---
title: Generate AGENTS.md
description: Create cross-tool compatible instructions
---

# Generate AGENTS.md Workflow

Create a standalone AGENTS.md file for cross-tool compatibility.

---

## When to Use

- **Multi-tool teams** - Some use Claude Code, others use Cursor
- **CI/CD integration** - Simpler instructions for automated tools
- **Different content** - Want different rules for different tools

---

## Background: AGENTS.md Standard

[AGENTS.md](https://agents.md) is a universal standard for AI coding assistant instructions, adopted by 20,000+ repositories.

| Tool | Primary File | Fallback |
|------|--------------|----------|
| **Claude Code** | CLAUDE.md | AGENTS.md |
| **Cursor** | AGENTS.md | - |
| **OpenAI Codex** | AGENTS.md | - |
| **GitHub Copilot** | AGENTS.md | - |
| **Google Jules** | AGENTS.md | - |

---

## Option 1: Symlink (Recommended)

The simplest approach - both files stay in sync:

```bash
ln -s CLAUDE.md AGENTS.md
```

**Pros**:

- Zero maintenance
- Always in sync
- No duplication

**Cons**:

- Same content for all tools
- Symlinks not supported on some Windows setups

---

## Option 2: Generate Standalone

Use the workflow to generate a separate file:

```
@.agent/workflows/generate-agents-md.md

Generate AGENTS.md for cross-tool compatibility
```

### What AI Does

1. **Extracts Operations section** from CLAUDE.md
2. **Extracts Boundaries section** from CLAUDE.md
3. **Adds Quick Reference** for common tasks
4. **Creates standalone AGENTS.md**

### Output Structure

```markdown
# AGENTS.md

## Operations
- Setup Commands
- Testing Commands
- Build Commands
- Code Style Commands

## Boundaries
- Protected Files
- Never Commit
- Ask Before Modifying

## Quick Reference
- Task classification
- Common guardrails
- Code standards
```

---

## Comparison

### CLAUDE.md (Full)

~500 lines containing:

- Operations
- Boundaries
- Quick Reference
- **Core Guardrails (35+ rules)**
- **4D Methodology**
- **SDLC stages**
- **Context System**
- **Anti-Patterns**

### AGENTS.md (Operations-Only)

~200 lines containing:

- Operations
- Boundaries
- Quick Reference

!!! note "Why Smaller?"

    AGENTS.md follows the "operations first" pattern - commands and boundaries only. The full methodology stays in CLAUDE.md.

---

## When to Use Each

### Use Symlink When

- Single-tool team
- Want identical behavior across tools
- On Unix/Mac (symlinks work)

### Use Standalone When

- Multi-tool team with different needs
- Want simpler instructions for some tools
- Need Windows compatibility
- CI/CD tools that need minimal context

---

## Maintaining Both Files

If using standalone AGENTS.md:

### After CLAUDE.md Changes

```
@.agent/workflows/generate-agents-md.md

Regenerate AGENTS.md to match updated CLAUDE.md
```

### Keeping in Sync

Add to your PR checklist:

- [ ] If CLAUDE.md changed, regenerate AGENTS.md

---

## Example Generated AGENTS.md

```markdown
# AGENTS.md

> AI coding assistant instructions (operations-focused)
> For full methodology, see CLAUDE.md

## Operations

### Setup Commands
```bash
npm install          # Install dependencies
npm run dev          # Start development server
```

### Testing Commands
```bash
npm test             # Run tests
npm run test:cov     # With coverage
```

### Build Commands
```bash
npm run build        # Production build
npm run lint         # Lint check
```

## Boundaries

### Protected Files
- package-lock.json (dependency lock)
- .env files (environment config)
- Database migrations (once deployed)

### Never Commit
- Secrets, API keys
- node_modules/
- .env files

### Ask Before Modifying
- Authentication logic
- Public API contracts
- CI/CD configuration

## Quick Reference

### Task Classification
- ATOMIC: <5 files, clear scope
- FEATURE: 5-10 files
- COMPLEX: >10 files

### Code Standards
- Functions ≤50 lines
- Files ≤300 lines
- Tests >80% for business logic
- Conventional commits
```

---

## Related

<div class="grid cards" markdown>

-   :material-file-document:{ .lg .middle } **CLAUDE.md**

    ---

    The full instruction file.

    [:octicons-arrow-right-24: CLAUDE.md](../core/claude-md.md)

-   :material-tools:{ .lg .middle } **Cross-Tool Setup**

    ---

    Detailed cross-tool configuration.

    [:octicons-arrow-right-24: Cross-Tool](../reference/cross-tool.md)

</div>
