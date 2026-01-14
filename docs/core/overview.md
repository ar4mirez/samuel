---
title: Core System Overview
description: Understanding the AICoF architecture
---

# Core System Overview

AICoF (Artificial Intelligence Coding Framework) is built on a simple but powerful architecture designed for progressive adoption and cross-tool compatibility.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Your Project                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────┐    ┌──────────────────────────────────┐   │
│  │  CLAUDE.md  │    │           .agent/                │   │
│  │             │    │                                  │   │
│  │  • Commands │    │  ├── language-guides/ (auto-load)│   │
│  │  • Guards   │    │  ├── workflows/ (on-demand)      │   │
│  │  • Methods  │    │  ├── tasks/ (generated)          │   │
│  │             │    │  ├── memory/ (decisions)         │   │
│  └─────────────┘    │  └── project.md (your stack)     │   │
│        ↑            └──────────────────────────────────┘   │
│        │                           ↑                        │
│        └───────────────────────────┘                        │
│                    References                               │
└─────────────────────────────────────────────────────────────┘
```

---

## Core Components

### CLAUDE.md - The Brain

The main instruction file loaded by AI assistants. Contains:

| Section | Purpose |
|---------|---------|
| **Operations** | Commands for setup, testing, building |
| **Boundaries** | Files/actions AI should not modify |
| **Quick Reference** | Task classification, emergency links |
| **Guardrails** | 35+ testable rules |
| **4D Methodology** | ATOMIC/FEATURE/COMPLEX modes |
| **SDLC** | Software Development Lifecycle stages |
| **Context System** | How `.agent/` directory works |
| **Anti-Patterns** | What to avoid |

**Size**: ~500 lines (optimized for token efficiency)

[:octicons-arrow-right-24: Learn about CLAUDE.md](claude-md.md)

---

### .agent/ Directory - The Memory

Project-specific context that grows over time:

```
.agent/
├── README.md              # How to use .agent/
├── project.md             # Your tech stack (created when chosen)
├── patterns.md            # Coding patterns (created when emerge)
├── state.md               # Current work (for multi-session)
├── language-guides/       # Auto-load based on file type
├── workflows/             # On-demand workflows
├── tasks/                 # PRDs and task lists
└── memory/                # Decision logs
```

**Loading Protocol**:

1. **Session Start**: AI loads CLAUDE.md → checks for state.md
2. **During Work**: Language guide loads based on file extensions
3. **Complex Features**: Workflows loaded on-demand
4. **Reference Needed**: patterns.md, project.md loaded as needed

[:octicons-arrow-right-24: Learn about .agent directory](agent-directory.md)

---

## The 3 Modes

AI auto-detects which mode based on task complexity:

### ATOMIC Mode

**For**: Bug fixes, small features (<5 files)

```mermaid
graph LR
    A[Task] --> B[Deconstruct]
    B --> C[Diagnose]
    C --> D[Develop]
    D --> E[Deliver]
    E --> F[Done ✓]
```

- Direct implementation
- Quick validation
- One commit

### FEATURE Mode

**For**: New components, API endpoints (5-10 files)

```mermaid
graph TD
    A[Task] --> B[Break into 3-5 subtasks]
    B --> C1[Subtask 1]
    B --> C2[Subtask 2]
    B --> C3[Subtask 3]
    C1 --> D[Integration Test]
    C2 --> D
    C3 --> D
    D --> E[Documentation]
    E --> F[Done ✓]
```

- Break into subtasks
- Sequential implementation
- Multiple commits

### COMPLEX Mode

**For**: New subsystems, architecture (>10 files)

```mermaid
graph TD
    A[Task] --> B[Create PRD]
    B --> C[Generate Tasks]
    C --> D[Task 1]
    D --> E[Task 2]
    E --> F[Task N]
    F --> G[Integration]
    G --> H[Documentation]
    H --> I[Done ✓]
```

- PRD workflow (optional but recommended)
- Full task breakdown
- Staged delivery

[:octicons-arrow-right-24: Learn the Methodology](methodology.md)

---

## The 4D Methodology

Every task follows four phases:

| Phase | Purpose | Key Question |
|-------|---------|--------------|
| **Deconstruct** | Break down the task | What's the minimal change? |
| **Diagnose** | Identify risks | Will this break anything? |
| **Develop** | Implement with tests | Does it meet guardrails? |
| **Deliver** | Validate and commit | Is it ready for production? |

[:octicons-arrow-right-24: Deep dive into 4D](methodology.md)

---

## Guardrails System

35+ testable rules across categories:

<div class="grid cards" markdown>

-   :material-code-braces:{ .lg .middle } **Code Quality**

    ---

    - Functions ≤50 lines
    - Files ≤300 lines
    - Complexity ≤10
    - All exports typed

-   :material-shield-lock:{ .lg .middle } **Security**

    ---

    - Input validation
    - Parameterized queries
    - No hardcoded secrets
    - Dependency auditing

-   :material-test-tube:{ .lg .middle } **Testing**

    ---

    - >80% business logic
    - >60% overall
    - Regression tests
    - No flaky tests

-   :material-git:{ .lg .middle } **Git**

    ---

    - Conventional commits
    - Atomic changes
    - No direct to main
    - All tests pass

</div>

[:octicons-arrow-right-24: All Guardrails](guardrails.md)

---

## Language Guides

Auto-loaded based on file extensions:

| Language | Extensions | Guide |
|----------|------------|-------|
| TypeScript | `.ts`, `.tsx`, `.js`, `.jsx` | [typescript.md](../languages/typescript.md) |
| Python | `.py` | [python.md](../languages/python.md) |
| Go | `.go` | [go.md](../languages/go.md) |
| Rust | `.rs` | [rust.md](../languages/rust.md) |
| Kotlin | `.kt`, `.kts` | [kotlin.md](../languages/kotlin.md) |

**No manual selection needed** - AI detects and loads automatically.

---

## Workflows

On-demand workflows for complex tasks:

| Workflow | When to Use |
|----------|-------------|
| [Initialize Project](../workflows/initialize-project.md) | New or existing project setup |
| [Create PRD](../workflows/create-prd.md) | Plan complex features |
| [Generate Tasks](../workflows/generate-tasks.md) | Break PRD into tasks |
| [Troubleshooting](../workflows/troubleshooting.md) | Debug systematically |
| [Generate AGENTS.md](../workflows/generate-agents-md.md) | Cross-tool compatibility |

---

## Progressive Growth

The system grows with your project:

```
Day 1:     CLAUDE.md + .agent/ templates only
           ↓
Week 1:    .agent/project.md created (tech stack)
           ↓
Month 1:   .agent/patterns.md populated (conventions)
           ↓
Ongoing:   .agent/memory/ captures decisions
```

!!! tip "Don't Over-Document"

    Let documentation emerge naturally. Don't create files preemptively.

---

## Cross-Tool Compatibility

Works with any AI coding assistant:

| Tool | Primary File | Fallback |
|------|--------------|----------|
| **Claude Code** | CLAUDE.md | AGENTS.md |
| **Cursor** | AGENTS.md | CLAUDE.md |
| **OpenAI Codex** | AGENTS.md | - |
| **GitHub Copilot** | AGENTS.md | - |
| **Google Jules** | AGENTS.md | - |

Setup: `ln -s CLAUDE.md AGENTS.md`

[:octicons-arrow-right-24: Cross-Tool Setup](../reference/cross-tool.md)

---

## Next Steps

<div class="grid cards" markdown>

-   :material-file-document:{ .lg .middle } **CLAUDE.md Deep Dive**

    ---

    Understand the main instruction file.

    [:octicons-arrow-right-24: CLAUDE.md](claude-md.md)

-   :material-shield:{ .lg .middle } **All Guardrails**

    ---

    Review the 35+ testable rules.

    [:octicons-arrow-right-24: Guardrails](guardrails.md)

</div>
