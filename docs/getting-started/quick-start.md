---
title: Quick Start
description: Get started with AICoF in 60 seconds
---

# Quick Start Guide

Get up and running with AICoF (Artificial Intelligence Coding Framework) in under a minute.

---

## 60-Second Setup

=== "CLI (Recommended)"

    ```bash
    # 1. Install the CLI
    brew tap ar4mirez/tap && brew install aicof

    # 2. Initialize your project
    cd your-project
    aicof init

    # 3. Done! Start coding with AI guardrails
    ```

=== "Curl Install"

    ```bash
    # 1. Install via curl
    curl -sSL https://raw.githubusercontent.com/ar4mirez/aicof/main/install.sh | sh

    # 2. Initialize your project
    cd your-project
    aicof init
    ```

=== "Go Install"

    ```bash
    # 1. Install via Go
    go install github.com/ar4mirez/aicof/cmd/aicof@latest

    # 2. Initialize your project
    cd your-project
    aicof init
    ```

=== "Manual Setup"

    ```bash
    # 1. Copy to your project
    cp -r /path/to/aicof/{CLAUDE.md,.agent} ./

    # 2. (Optional) For cross-tool compatibility
    ln -s CLAUDE.md AGENTS.md
    ```

**The system works immediately:**

- [x] AI loads CLAUDE.md automatically (500 lines of guardrails + operations)
- [x] Language guides auto-load based on file extensions
- [x] 11 CLI commands for component management
- [x] Progressive - starts minimal, grows with your project
- [x] Cross-tool: Works with any AI assistant via AGENTS.md

---

## Discover Components Before Installing

One of AICoF's best features is **component discovery**. Before installing anything, explore what's available:

### Search for Components

```bash
# Find anything matching "react"
aicof search react

# Search only frameworks
aicof search --type fw api

# Fuzzy matching - finds "python" even with typos
aicof search pythn
```

### Preview Before Installing

```bash
# See component details
aicof info framework react

# Preview the actual content (first 20 lines)
aicof info lang typescript --preview 20

# See what's related
aicof info fw nextjs  # Shows: Related to typescript, react
```

### Add What You Need

```bash
# Add a framework
aicof add framework react

# Add a language guide
aicof add language typescript

# Add a workflow
aicof add workflow code-review

# Use short aliases
aicof add fw react
aicof add lang ts
aicof add wf security-audit
```

### Check What's Installed

```bash
# List installed components
aicof list

# See what's available (not installed)
aicof list --available

# Filter by type
aicof list --type frameworks
```

---

## Choose Your Path

### :material-plus-circle: Path 1: New Project with CLI

```bash
# Create and initialize a new project
mkdir my-project && cd my-project
aicof init

# Or initialize with specific components
aicof init --languages typescript,python --frameworks react
```

**The CLI will:**

1. Create `CLAUDE.md` with all guardrails
2. Create `.agent/` directory structure
3. Let you select languages, frameworks, workflows interactively
4. Set up cross-tool compatibility (AGENTS.md symlink)

---

### :material-folder-open: Path 2: Existing Project

```bash
cd existing-project
aicof init
```

**AI will:**

1. Detect existing files (won't overwrite without `--force`)
2. Create AICoF structure
3. Let you select relevant components

**Or use the workflow for deep analysis:**

```
@.agent/workflows/initialize-project.md

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
| `init` | Initialize AICoF in a project | `aicof init` |
| `search` | Find components by keyword | `aicof search react` |
| `info` | Show component details | `aicof info fw react --preview` |
| `add` | Add a component | `aicof add lang typescript` |
| `remove` | Remove a component | `aicof remove fw react` |
| `list` | List components | `aicof list --available` |
| `config` | Manage configuration | `aicof config list` |
| `diff` | Compare versions | `aicof diff --installed` |
| `update` | Update framework | `aicof update --check` |
| `doctor` | Check installation health | `aicof doctor` |
| `version` | Show version info | `aicof version` |

**Type aliases**: `language` (lang, l), `framework` (fw, f), `workflow` (wf, w)

[:octicons-arrow-right-24: Full CLI Reference](../reference/cli.md)

---

## Verify Installation

After setup, verify everything is working:

```bash
# Check installation health
aicof doctor

# See what's installed
aicof list

# Check for updates
aicof version --check
```

**Expected output from `aicof doctor`:**

```
AICoF Health Check
==================

[OK] CLAUDE.md exists
[OK] .agent/ directory exists
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

The `.agent/` directory grows with your project:

| Timeline | What Gets Created |
|----------|-------------------|
| **Day 1** | Only CLAUDE.md + templates |
| **Week 1** | `.agent/project.md` (tech stack) |
| **Month 1** | `.agent/patterns.md` (conventions) |
| **Ongoing** | `.agent/memory/` (decisions) |

!!! tip "Don't Over-Document"

    Let the documentation grow naturally. Don't create `project.md` on day one - wait until you make architecture decisions.

---

## Next Steps

<div class="grid cards" markdown>

-   :material-console:{ .lg .middle } **CLI Reference**

    ---

    Learn all 11 commands and their options.

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
