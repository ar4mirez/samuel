---
title: RFDs (Requests for Discussion)
description: Technical decision documents for AICoF projects
---

# RFDs

Requests for Discussion (RFDs) document technical decisions, explore options, and build consensus before implementation.

## What is an RFD?

RFDs are inspired by [Oxide Computer Company's RFD system](https://oxide.computer/blog/rfd-1-requests-for-discussion), which itself draws from the IETF's RFC tradition.

> "Ideas should be timely rather than polished."

**Core principles:**

- **Share early, refine together** - Don't polish in isolation
- **Document options** - Show what was considered, not just what was chosen
- **Build consensus** - Get input before committing to implementation
- **Preserve context** - Future readers need to understand "why"

## RFD vs PRD

| Aspect | RFD | PRD |
|--------|-----|-----|
| **Purpose** | Explore options, build consensus | Define implementation |
| **Focus** | "Why" and "What options" | "What" and "How" |
| **Timing** | Early ideation | Ready to implement |
| **Outcome** | Decision documented | Task list generated |

**Typical flow:**

```
Idea → RFD (explore options) → Decision → PRD (plan implementation) → Tasks → Code
```

## RFD States

| State | Description | Location |
|-------|-------------|----------|
| **Prediscussion** | Very early, author still forming idea | `.agent/rfd/` (private) |
| **Ideation** | Ready for informal feedback | `.agent/rfd/` (private) |
| **Discussion** | Open for team discussion | `docs/rfd/` (public) |
| **Published** | Decision made, documented | `docs/rfd/` (public) |
| **Committed** | Implementation complete | `docs/rfd/` (public) |
| **Abandoned** | Rejected or superseded | `docs/rfd/` (public) |

```
Prediscussion → Ideation → Discussion → Published → Committed
                    ↓           ↓
                    └─────→ Abandoned
```

## Creating an RFD

Use the workflow:

```
@.agent/workflows/create-rfd.md

Explore options for [topic]
```

Or ask Claude to suggest one during a discussion:

```
"Should we use Redis or PostgreSQL for sessions?"
```

Claude will suggest creating an RFD when it detects you're exploring multiple approaches.

## Active RFDs

RFDs currently in **Discussion** state, seeking feedback:

| RFD | Title | Authors | Updated |
|-----|-------|---------|---------|
| *No active RFDs* | | | |

## Committed RFDs

RFDs with decisions made and implementation complete:

| RFD | Title | Authors | Updated |
|-----|-------|---------|---------|
| [0001](0001.md) | Progressive Disclosure Architecture for AI System Instructions | ar4mirez | 2026-01-15 |
| [0002](0002.md) | Idiomatic Go Project Layout for CLI Tools | ar4mirez | 2026-01-15 |
| [0003](0003.md) | Composable CLI Commands Over Overloaded Flags | ar4mirez | 2026-01-15 |
| [0004](0004.md) | Smart Interactive Mode via Flag Detection | ar4mirez | 2026-01-15 |

## All RFDs

| RFD | Title | State | Labels |
|-----|-------|-------|--------|
| [0001](0001.md) | Progressive Disclosure Architecture | Committed | architecture, documentation, token-efficiency |
| [0002](0002.md) | Idiomatic Go Project Layout | Committed | go, project-structure, cli |
| [0003](0003.md) | Composable CLI Commands | Committed | cli, ux, architecture |
| [0004](0004.md) | Smart Interactive Mode | Committed | cli, ux, pattern |

## RFD Index

The master index of all RFDs is maintained in `rfd-index.yaml` at the project root. This file tracks:

- RFD numbers and titles
- Current state
- Authors and labels
- File paths
- Related PRDs

## When to Use RFDs

**Use an RFD when:**

- Exploring multiple valid approaches
- Need team input before deciding
- Making architectural decisions
- Proposing significant changes
- Decision affects multiple team members

**Skip RFDs for:**

- Obvious decisions with one option
- Small, easily reversible changes
- Implementation details (use PRD instead)
- Urgent fixes (document post-hoc if needed)

## Learn More

- [Create RFD Workflow](../workflows/create-rfd.md) - Full workflow documentation
- [Oxide's RFD 1](https://rfd.shared.oxide.computer/rfd/0001) - Original inspiration
- [Create PRD Workflow](../workflows/create-prd.md) - For implementation planning after decisions
