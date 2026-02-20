---
title: Quick Start
description: Get started with Samuel in 60 seconds
---

# Quick Start Guide

Get up and running with Samuel (Artificial Intelligence Coding Framework) in under a minute.

---

## 60-Second Setup

=== "CLI (Recommended)"

    ```bash
    # 1. Install the CLI
    brew tap ar4mirez/tap && brew install samuel

    # 2. Initialize your project
    cd your-project
    samuel init

    # 3. Done! Start coding with AI guardrails
    ```

=== "Curl Install"

    ```bash
    # 1. Install via curl
    curl -sSL https://raw.githubusercontent.com/ar4mirez/samuel/main/install.sh | sh

    # 2. Initialize your project
    cd your-project
    samuel init
    ```

=== "Go Install"

    ```bash
    # 1. Install via Go
    go install github.com/ar4mirez/samuel/cmd/samuel@latest

    # 2. Initialize your project
    cd your-project
    samuel init
    ```

=== "Manual Setup"

    ```bash
    # 1. Copy to your project
    cp -r /path/to/samuel/{CLAUDE.md,.claude} ./

    # 2. (Optional) For cross-tool compatibility
    ln -s CLAUDE.md AGENTS.md
    ```

**The system works immediately:**

- [x] AI loads CLAUDE.md automatically (500 lines of guardrails + operations)
- [x] Language guides auto-load based on file extensions
- [x] 14 CLI commands for component management
- [x] Progressive - starts minimal, grows with your project
- [x] Cross-tool: Works with any AI assistant via AGENTS.md

---

## Discover Components Before Installing

One of Samuel's best features is **component discovery**. Before installing anything, explore what's available:

### Search for Components

```bash
# Find anything matching "react"
samuel search react

# Search only frameworks
samuel search --type fw api

# Fuzzy matching - finds "python" even with typos
samuel search pythn
```

### Preview Before Installing

```bash
# See component details
samuel info framework react

# Preview the actual content (first 20 lines)
samuel info lang typescript --preview 20

# See what's related
samuel info fw nextjs  # Shows: Related to typescript, react
```

### Add What You Need

```bash
# Add a framework
samuel add framework react

# Add a language guide
samuel add language typescript

# Add a workflow
samuel add workflow code-review

# Use short aliases
samuel add fw react
samuel add lang ts
samuel add wf security-audit
```

### Check What's Installed

```bash
# List installed components
samuel list

# See what's available (not installed)
samuel list --available

# Filter by type
samuel list --type frameworks
```

---

## Choose Your Path

### :material-plus-circle: Path 1: New Project with CLI

```bash
# Create and initialize a new project
mkdir my-project && cd my-project
samuel init

# Or initialize with specific components
samuel init --languages typescript,python --frameworks react
```

**The CLI will:**

1. Create `CLAUDE.md` with all guardrails
2. Create `.claude/` directory structure
3. Let you select languages, frameworks, workflows interactively
4. Set up cross-tool compatibility (AGENTS.md symlink)

---

### :material-folder-open: Path 2: Existing Project

```bash
cd existing-project
samuel init
```

**AI will:**

1. Detect existing files (won't overwrite without `--force`)
2. Create Samuel structure
3. Let you select relevant components

**Or use the workflow for deep analysis:**

```
@.claude/skills/initialize-project/SKILL.md

"This is an existing project - analyze the codebase"
```

---

### :material-lightning-bolt: Path 3: Jump Right In

Just start coding. AI follows guardrails automatically.

=== "Simple Tasks"

    ```
    "Fix the login button alignment"
    ```

    AI uses **ATOMIC mode** - single file, quick fix, tests, commit.

=== "Features"

    ```
    "Add user profile editing"
    ```

    AI uses **FEATURE mode** - breaks into subtasks, implements systematically.

=== "Complex Work"

    ```
    "Build real-time chat with WebSockets"
    ```

    AI suggests **COMPLEX mode** - offers PRD workflow for structured approach.

---

## CLI Command Reference

Here are the most useful commands to get started:

| Command | Description | Example |
|---------|-------------|---------|
| `init` | Initialize Samuel in a project | `samuel init` |
| `search` | Find components by keyword | `samuel search react` |
| `info` | Show component details | `samuel info fw react --preview` |
| `add` | Add a component | `samuel add lang typescript` |
| `remove` | Remove a component | `samuel remove fw react` |
| `list` | List components | `samuel list --available` |
| `config` | Manage configuration | `samuel config list` |
| `diff` | Compare versions | `samuel diff --installed` |
| `update` | Update framework | `samuel update --check` |
| `doctor` | Check installation health | `samuel doctor` |
| `version` | Show version info | `samuel version` |
| `skill` | Manage Agent Skills | `samuel skill list` |
| `auto` | Autonomous AI coding loop | `samuel auto status` |
| `sync` | Sync per-folder CLAUDE.md/AGENTS.md | `samuel sync --dry-run` |

**Type aliases**: `language` (lang, l), `framework` (fw, f), `workflow` (wf, w)

[:octicons-arrow-right-24: Full CLI Reference](../reference/cli.md)

---

## Verify Installation

After setup, verify everything is working:

```bash
# Check installation health
samuel doctor

# See what's installed
samuel list

# Check for updates
samuel version --check
```

**Expected output from `samuel doctor`:**

```
Samuel Health Check
==================

[OK] CLAUDE.md exists
[OK] .claude/ directory exists
[OK] Configuration valid
[OK] 3 languages installed
[OK] 2 frameworks installed
[OK] All workflows available

Status: Healthy
```

---

## What Happens Next?

Once initialized, the system is ready. Here's what to expect:

### Automatic Loading

When you start working with an AI assistant:

1. **CLAUDE.md loads automatically** - 500 lines of guardrails and operations
2. **Language guides auto-load** - Based on file extensions you're working with
3. **Guardrails are enforced** - Code quality, security, testing standards

### Progressive Growth

The `.claude/` directory grows with your project:

| Timeline | What Gets Created |
|----------|-------------------|
| **Day 1** | Only CLAUDE.md + templates |
| **Week 1** | `.claude/project.md` (tech stack) |
| **Month 1** | `.claude/patterns.md` (conventions) |
| **Ongoing** | `.claude/memory/` (decisions) |

!!! tip "Don't Over-Document"

    Let the documentation grow naturally. Don't create `project.md` on day one - wait until you make architecture decisions.

---

## Next Steps

<div class="grid cards" markdown>

-   :material-console:{ .lg .middle } **CLI Reference**

    ---

    Learn all 14 commands and their options.

    [:octicons-arrow-right-24: CLI Commands](../reference/cli.md)

-   :material-book:{ .lg .middle } **Learn the Methodology**

    ---

    Understand the 4D approach: Deconstruct, Diagnose, Develop, Deliver.

    [:octicons-arrow-right-24: 4D Methodology](../core/methodology.md)

-   :material-shield:{ .lg .middle } **Review Guardrails**

    ---

    See all 35+ rules that AI will follow.

    [:octicons-arrow-right-24: All Guardrails](../core/guardrails.md)

-   :material-cog:{ .lg .middle } **Try a Workflow**

    ---

    Use PRD workflow for a complex feature.

    [:octicons-arrow-right-24: Workflows](../workflows/index.md)

</div>
