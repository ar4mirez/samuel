---
title: The .claude Directory
description: Native Claude Code skills directory for AI-assisted development
---

# The .claude Directory

The `.claude/` directory is Claude Code's native project directory. Samuel uses it to store skills that extend AI capabilities.

---

## Overview

While CLAUDE.md provides universal guardrails, `.claude/` stores skills specific to YOUR project:

```
.claude/
├── skills/                # Agent Skills (language guides, frameworks, workflows)
│   ├── README.md
│   ├── typescript-guide/  # 21 language guide skills
│   │   ├── SKILL.md
│   │   └── references/
│   ├── react/             # 33 framework skills
│   │   ├── SKILL.md
│   │   └── references/
│   ├── create-prd/        # 16 workflow skills
│   │   └── SKILL.md
│   ├── auto/              # Autonomous loop skill
│   │   └── SKILL.md
│   └── ... (more skills)
├── auto/                  # Autonomous loop state (generated)
│   ├── prd.json           # Machine-readable task state
│   ├── progress.md       # Append-only learnings journal
│   ├── prompt.md          # Iteration prompt template
│   └── auto.sh            # Loop orchestration script
└── settings.local.json    # Claude Code local settings
```

---

## What's Inside

### Skills

Skills are capability modules following the [Agent Skills](https://agentskills.io) standard:

| Type | Count | Loaded | Example |
|------|-------|--------|---------|
| Language guides | 21 | Auto-load by file extension | `skills/go-guide/SKILL.md` |
| Framework skills | 33 | On-demand | `skills/react/SKILL.md` |
| Workflow skills | 16 | On-demand | `skills/create-prd/SKILL.md` |

### Per-Folder CLAUDE.md Files

Samuel creates stub `CLAUDE.md` files in existing project directories during `samuel init`. These are loaded automatically when AI works in that directory. Customize them with folder-specific instructions, conventions, and constraints.

---

## Language Guide Skills

Auto-loaded based on file extensions you're working with (21 languages):

| Language | Extensions | Guide |
|----------|------------|-------|
| TypeScript | `.ts`, `.tsx`, `.js`, `.jsx` | `skills/typescript-guide/SKILL.md` |
| Python | `.py` | `skills/python-guide/SKILL.md` |
| Go | `.go` | `skills/go-guide/SKILL.md` |
| Rust | `.rs` | `skills/rust-guide/SKILL.md` |
| Kotlin | `.kt`, `.kts` | `skills/kotlin-guide/SKILL.md` |
| Java | `.java` | `skills/java-guide/SKILL.md` |
| C# | `.cs` | `skills/csharp-guide/SKILL.md` |
| PHP | `.php` | `skills/php-guide/SKILL.md` |
| Swift | `.swift` | `skills/swift-guide/SKILL.md` |
| C/C++ | `.c`, `.cpp`, `.h`, `.hpp` | `skills/cpp-guide/SKILL.md` |
| Ruby | `.rb` | `skills/ruby-guide/SKILL.md` |
| *+ 10 more* | SQL, Shell, R, Dart, HTML/CSS, Lua, Assembly, CUDA, Solidity, Zig | |

**No manual selection needed** - AI detects automatically.

[:octicons-arrow-right-24: Language Guides](../languages/index.md)

---

## Framework Skills

On-demand framework-specific patterns (33 framework skills across 11 language families):

| Language | Frameworks |
|----------|------------|
| TypeScript/JS | React, Next.js, Express |
| Python | Django, FastAPI, Flask |
| Go | Gin, Echo, Fiber |
| Rust | Axum, Actix-web, Rocket |
| Kotlin | Spring Boot (Kotlin), Ktor, Android Compose |
| Java | Spring Boot, Quarkus, Micronaut |
| C# | ASP.NET Core, Blazor, Unity |
| PHP | Laravel, Symfony, WordPress |
| Swift | SwiftUI, UIKit, Vapor |
| Ruby | Rails, Sinatra, Hanami |
| Dart | Flutter, Shelf, Dart Frog |

**How to load**: Reference the framework or work in a project that uses it.

[:octicons-arrow-right-24: Framework Guides](../frameworks/index.md)

---

## Workflows

On-demand workflows for structured tasks:

| Workflow | When to Use |
|----------|-------------|
| `initialize-project` | New or existing project setup |
| `create-prd` | Plan complex features |
| `generate-tasks` | Break PRD into tasks |
| `troubleshooting` | Debug systematically |
| `generate-agents-md` | Cross-tool compatibility |
| `create-skill` | Create portable Agent Skills |
| `auto` | Autonomous AI coding loop |

**How to invoke**:

```
@.claude/skills/create-prd/SKILL.md
```

[:octicons-arrow-right-24: Workflows](../workflows/index.md)

---

## Loading Protocol

AI follows this protocol when starting a session:

```
1. Load CLAUDE.md (always)
   ↓
2. Load per-folder CLAUDE.md (when working in a subdirectory)
   ↓
3. Auto-load language guide based on file extensions
   ↓
4. On-demand: Load workflow or framework skills as needed
```

---

## Git Configuration

### What to Commit

```
CLAUDE.md
AGENTS.md
.claude/skills/               # Language guides, frameworks, and workflows
Per-folder CLAUDE.md files     # Folder-specific instructions
```

### What to Gitignore

```gitignore
# Claude Code local settings
.claude/settings.local.json
```

---

## Best Practices

### Do

- [x] Customize per-folder CLAUDE.md files with specific instructions
- [x] Create custom skills for recurring project tasks
- [x] Use `samuel skill create` to scaffold new skills

### Don't

- [ ] Manually edit language guide skills (they're templates, use `samuel update`)
- [ ] Over-organize the directory
- [ ] Put secrets or credentials in CLAUDE.md files

---

## Related

<div class="grid cards" markdown>

-   :material-file-document:{ .lg .middle } **CLAUDE.md**

    ---

    The main instruction file.

    [:octicons-arrow-right-24: CLAUDE.md](claude-md.md)

-   :material-code-braces:{ .lg .middle } **Language Guides**

    ---

    Auto-loaded language rules (21 languages).

    [:octicons-arrow-right-24: Languages](../languages/index.md)

-   :material-layers:{ .lg .middle } **Framework Guides**

    ---

    Framework-specific patterns (33 frameworks).

    [:octicons-arrow-right-24: Frameworks](../frameworks/index.md)

-   :material-cog:{ .lg .middle } **Workflows**

    ---

    On-demand structured workflows.

    [:octicons-arrow-right-24: Workflows](../workflows/index.md)

</div>
