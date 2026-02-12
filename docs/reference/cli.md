---
title: CLI Command Reference
description: Complete reference for all Samuel CLI commands
---

# CLI Command Reference

The Samuel CLI provides 13 commands for managing and discovering components. This page documents all commands, flags, and usage examples.

---

## Installation

```bash
# macOS/Linux (Homebrew)
brew tap ar4mirez/tap && brew install samuel

# macOS/Linux (Curl)
curl -sSL https://raw.githubusercontent.com/ar4mirez/samuel/main/install.sh | sh

# Go
go install github.com/ar4mirez/samuel/cmd/samuel@latest
```

---

## Global Flags

These flags work with all commands:

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Enable verbose output for debugging |
| `--no-color` | | Disable colored output |
| `--help` | `-h` | Show help for any command |

**Example:**

```bash
samuel --verbose init
samuel --no-color list
```

---

## Type Aliases

When specifying component types, you can use these aliases:

| Full Name | Aliases |
|-----------|---------|
| `language` | `lang`, `l` |
| `framework` | `fw`, `f` |
| `workflow` | `wf`, `w` |

**Example:**

```bash
# These are equivalent
samuel add language typescript
samuel add lang typescript
samuel add l typescript
```

---

## Commands

### init

Initialize Samuel in a new or existing project.

**Usage:**

```bash
samuel init [project-name] [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--template <name>` | Use a specific template (minimal, full) |
| `--languages <list>` | Pre-select languages (comma-separated) |
| `--frameworks <list>` | Pre-select frameworks (comma-separated) |
| `--workflows <list>` | Pre-select workflows (comma-separated) |
| `--force` | Overwrite existing files without prompting |
| `--non-interactive` | Skip all prompts, use defaults or flags |

**Examples:**

```bash
# Interactive setup in current directory
samuel init

# Create new project directory
samuel init my-project

# Non-interactive with specific components
samuel init --languages go,python --frameworks gin --non-interactive

# Force overwrite existing files
samuel init --force
```

---

### search

Search for components by keyword with fuzzy matching.

**Usage:**

```bash
samuel search <query> [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--type` | `-t` | Filter by type (lang/fw/wf) |
| `--limit` | `-l` | Maximum results (default: 10) |

**Examples:**

```bash
# Search across all component types
samuel search react

# Search only frameworks
samuel search --type fw api

# Search only languages
samuel search -t lang script

# Fuzzy matching - finds "python" even with typo
samuel search pythn

# Limit results
samuel search web --limit 5
```

**Search Scoring:**

Results are ranked by relevance:

| Match Type | Score | Example |
|------------|-------|---------|
| Exact match | 100 | "react" → React |
| Prefix match | 80 | "type" → TypeScript |
| Contains | 60 | "script" → TypeScript |
| Description match | 40 | "web framework" → Express |
| Fuzzy match (≤2 edits) | 25 | "pythn" → Python |

---

### info

Show detailed information about a component.

**Usage:**

```bash
samuel info <type> <name> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--preview <lines>` | Preview first N lines of the guide content |
| `--no-related` | Hide related components section |

**Examples:**

```bash
# Show framework details
samuel info framework react

# Show language with content preview
samuel info lang typescript --preview 20

# Using aliases
samuel info fw nextjs
samuel info l go --preview 30

# Hide related components
samuel info fw react --no-related
```

**Output includes:**

- Component name and type
- Description
- File path and size
- Installation status
- Related components (frameworks ↔ languages)
- Optional content preview

---

### add

Add a language guide, framework guide, or workflow.

**Usage:**

```bash
samuel add <type> <name> [flags]
```

**Examples:**

```bash
# Add a language guide
samuel add language typescript
samuel add lang python
samuel add l go

# Add a framework guide
samuel add framework react
samuel add fw django
samuel add f rails

# Add a workflow
samuel add workflow code-review
samuel add wf security-audit
samuel add w testing-strategy
```

---

### remove

Remove a language guide, framework guide, or workflow.

**Usage:**

```bash
samuel remove <type> <name> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--force` | Remove without confirmation prompt |

**Examples:**

```bash
# Remove with confirmation
samuel remove framework react

# Remove without confirmation
samuel remove fw react --force

# Remove a workflow
samuel remove wf code-review
```

---

### list

List installed or available components.

**Usage:**

```bash
samuel list [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--available` | Show available (not installed) components |
| `--type <type>` | Filter by type (languages/frameworks/workflows) |

**Examples:**

```bash
# List all installed components
samuel list

# List available (not installed) components
samuel list --available

# List only installed frameworks
samuel list --type frameworks

# List available languages
samuel list --available --type languages
```

---

### config

Manage Samuel configuration.

**Subcommands:**

| Subcommand | Description |
|------------|-------------|
| `config list` | Show all configuration values |
| `config get <key>` | Get a specific configuration value |
| `config set <key> <value>` | Set a configuration value |

**Valid Configuration Keys:**

| Key | Description |
|-----|-------------|
| `version` | Installed framework version |
| `registry` | GitHub repository URL for updates |
| `installed.languages` | Comma-separated list of installed languages |
| `installed.frameworks` | Comma-separated list of installed frameworks |
| `installed.workflows` | Comma-separated list of installed workflows |
| `installed.skills` | Comma-separated list of installed skills |
| `auto.enabled` | Enable autonomous AI coding loop |
| `auto.ai_tool` | AI tool for auto loop (claude, amp, codex) |
| `auto.max_iterations` | Maximum loop iterations (default: 50) |
| `auto.quality_checks` | Quality check commands for auto loop |

**Examples:**

```bash
# Show all configuration
samuel config list

# Get specific value
samuel config get version
samuel config get installed.languages

# Set values
samuel config set registry https://github.com/ar4mirez/samuel
samuel config set installed.languages go,rust,python
```

---

### diff

Compare versions to see what changed.

**Usage:**

```bash
samuel diff [version1] [version2] [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--installed` | Compare installed version vs latest |
| `--components` | Show only component changes (languages, frameworks, workflows) |

**Examples:**

```bash
# Compare installed vs latest
samuel diff --installed

# Compare two specific versions
samuel diff 1.6.0 1.7.0

# Compare local vs specific version
samuel diff 1.7.0

# Show only component changes
samuel diff --installed --components
```

**Output shows:**

- Added files
- Removed files
- Modified files
- Summary statistics

---

### update

Update Samuel to the latest version.

**Usage:**

```bash
samuel update [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--check` | Check for updates without applying |
| `--diff` | Show changes before updating |
| `--force` | Update without confirmation |
| `--version <v>` | Update to a specific version |

**Examples:**

```bash
# Check for updates
samuel update --check

# Preview changes before updating
samuel update --diff

# Update to latest
samuel update

# Update to specific version
samuel update --version 1.7.0

# Force update without prompts
samuel update --force
```

---

### doctor

Check installation health and diagnose issues.

**Usage:**

```bash
samuel doctor [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--fix` | Attempt to automatically fix issues |

**Examples:**

```bash
# Run health check
samuel doctor

# Auto-fix issues
samuel doctor --fix
```

**Checks performed:**

- CLAUDE.md exists and is readable
- .claude/ directory exists with correct structure
- Configuration file is valid
- Installed components are accessible
- No orphaned or corrupted files

---

### version

Show version information.

**Usage:**

```bash
samuel version [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--check` | Check for newer versions available |

**Examples:**

```bash
# Show version
samuel version

# Check for updates
samuel version --check
```

---

### skill

Manage Agent Skills — portable capability modules for AI agents.

**Subcommands:**

| Subcommand | Description |
|------------|-------------|
| `skill create <name>` | Create a new skill scaffold |
| `skill validate [name]` | Validate skill(s) against the Agent Skills spec |
| `skill list` | List installed skills |
| `skill info <name>` | Show detailed information about a skill |

**Examples:**

```bash
# Create a new skill
samuel skill create database-ops

# Validate all skills
samuel skill validate

# Validate a specific skill
samuel skill validate database-ops

# List installed skills
samuel skill list

# Show skill details
samuel skill info database-ops
```

**Skill name requirements:**

- Lowercase alphanumeric and hyphens only
- No consecutive hyphens (`--`)
- Cannot start or end with hyphen
- Maximum 64 characters

**Created files:**

```text
.claude/skills/<name>/
├── SKILL.md           # Skill definition with YAML frontmatter
├── scripts/           # Executable code
├── references/        # Additional documentation
└── assets/            # Templates and data files
```

---

### auto

Autonomous AI coding loop using the Ralph Wiggum methodology.

**Subcommands:**

| Subcommand | Description |
|------------|-------------|
| `auto init` | Initialize autonomous loop for a project |
| `auto convert <prd-path>` | Convert markdown PRD/tasks to prd.json |
| `auto status` | Show loop progress and current state |
| `auto start` | Begin or resume the autonomous loop |
| `auto task list` | List all tasks with status |
| `auto task complete <id>` | Mark a task as completed |
| `auto task skip <id>` | Mark a task as skipped |
| `auto task reset <id>` | Reset a task to pending |
| `auto task add <id> <title>` | Add a new task |

**init flags:**

| Flag | Description |
|------|-------------|
| `--prd <path>` | Path to PRD markdown file to convert |
| `--ai-tool <name>` | AI tool to use: claude, amp, cursor, codex (default: claude) |
| `--max-iterations <n>` | Maximum loop iterations (default: 50) |

**start flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--iterations <n>` | | Override max iterations for this run |
| `--yes` | `-y` | Skip confirmation prompt |
| `--dry-run` | | Show what would happen without executing |

**Examples:**

```bash
# Initialize from a PRD
samuel auto init --prd .claude/tasks/0001-prd-auth.md

# Initialize with custom settings
samuel auto init --ai-tool amp --max-iterations 100

# Convert a PRD to prd.json
samuel auto convert .claude/tasks/0001-prd-auth.md

# Check loop status
samuel auto status

# Start the loop
samuel auto start

# Start with custom iterations, skip confirmation
samuel auto start --iterations 20 --yes

# Dry run
samuel auto start --dry-run

# Manage tasks
samuel auto task list
samuel auto task complete 1.1
samuel auto task skip 2.3
samuel auto task reset 1.1
samuel auto task add "3.0" "New parent task"
```

**Generated files:**

```text
.claude/auto/
├── prd.json        # Machine-readable task state
├── progress.md    # Append-only learnings journal
├── prompt.md       # Iteration prompt template
└── auto.sh         # Loop orchestration script
```

[:octicons-arrow-right-24: Auto Workflow Guide](../workflows/auto.md)

---

## Common Workflows

### Setting Up a New Project

```bash
# Create and enter project directory
mkdir my-project && cd my-project

# Initialize with defaults
samuel init

# Or with specific components
samuel init --languages typescript,python --frameworks react,fastapi
```

### Discovering Components

```bash
# Search for what you need
samuel search api

# Get details
samuel info fw fastapi --preview 20

# Add it
samuel add fw fastapi
```

### Keeping Up to Date

```bash
# Check what's new
samuel diff --installed

# Preview changes
samuel update --check

# Update
samuel update
```

### Troubleshooting

```bash
# Check health
samuel doctor

# Fix issues
samuel doctor --fix

# Verbose output for debugging
samuel --verbose doctor
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | Component not found |
| 4 | Configuration error |

---

## Environment Variables

| Variable | Description |
|----------|-------------|
| `AICOF_NO_COLOR` | Disable colored output (same as `--no-color`) |
| `AICOF_VERBOSE` | Enable verbose output (same as `--verbose`) |

---

## See Also

- [Quick Start](../getting-started/quick-start.md) - Get started in 60 seconds
- [Installation](../getting-started/installation.md) - Detailed installation options
- [Workflows](../workflows/index.md) - Available workflows
