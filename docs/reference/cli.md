---
title: CLI Command Reference
description: Complete reference for all AICoF CLI commands
---

# CLI Command Reference

The AICoF CLI provides 11 commands for managing and discovering components. This page documents all commands, flags, and usage examples.

---

## Installation

```bash
# macOS/Linux (Homebrew)
brew tap ar4mirez/tap && brew install aicof

# macOS/Linux (Curl)
curl -sSL https://raw.githubusercontent.com/ar4mirez/aicof/main/install.sh | sh

# Go
go install github.com/ar4mirez/aicof/cmd/aicof@latest
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
aicof --verbose init
aicof --no-color list
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
aicof add language typescript
aicof add lang typescript
aicof add l typescript
```

---

## Commands

### init

Initialize AICoF in a new or existing project.

**Usage:**

```bash
aicof init [project-name] [flags]
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
aicof init

# Create new project directory
aicof init my-project

# Non-interactive with specific components
aicof init --languages go,python --frameworks gin --non-interactive

# Force overwrite existing files
aicof init --force
```

---

### search

Search for components by keyword with fuzzy matching.

**Usage:**

```bash
aicof search <query> [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--type` | `-t` | Filter by type (lang/fw/wf) |
| `--limit` | `-l` | Maximum results (default: 10) |

**Examples:**

```bash
# Search across all component types
aicof search react

# Search only frameworks
aicof search --type fw api

# Search only languages
aicof search -t lang script

# Fuzzy matching - finds "python" even with typo
aicof search pythn

# Limit results
aicof search web --limit 5
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
aicof info <type> <name> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--preview <lines>` | Preview first N lines of the guide content |
| `--no-related` | Hide related components section |

**Examples:**

```bash
# Show framework details
aicof info framework react

# Show language with content preview
aicof info lang typescript --preview 20

# Using aliases
aicof info fw nextjs
aicof info l go --preview 30

# Hide related components
aicof info fw react --no-related
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
aicof add <type> <name> [flags]
```

**Examples:**

```bash
# Add a language guide
aicof add language typescript
aicof add lang python
aicof add l go

# Add a framework guide
aicof add framework react
aicof add fw django
aicof add f rails

# Add a workflow
aicof add workflow code-review
aicof add wf security-audit
aicof add w testing-strategy
```

---

### remove

Remove a language guide, framework guide, or workflow.

**Usage:**

```bash
aicof remove <type> <name> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--force` | Remove without confirmation prompt |

**Examples:**

```bash
# Remove with confirmation
aicof remove framework react

# Remove without confirmation
aicof remove fw react --force

# Remove a workflow
aicof remove wf code-review
```

---

### list

List installed or available components.

**Usage:**

```bash
aicof list [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--available` | Show available (not installed) components |
| `--type <type>` | Filter by type (languages/frameworks/workflows) |

**Examples:**

```bash
# List all installed components
aicof list

# List available (not installed) components
aicof list --available

# List only installed frameworks
aicof list --type frameworks

# List available languages
aicof list --available --type languages
```

---

### config

Manage AICoF configuration.

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

**Examples:**

```bash
# Show all configuration
aicof config list

# Get specific value
aicof config get version
aicof config get installed.languages

# Set values
aicof config set registry https://github.com/ar4mirez/aicof
aicof config set installed.languages go,rust,python
```

---

### diff

Compare versions to see what changed.

**Usage:**

```bash
aicof diff [version1] [version2] [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--installed` | Compare installed version vs latest |
| `--components` | Show only component changes (languages, frameworks, workflows) |

**Examples:**

```bash
# Compare installed vs latest
aicof diff --installed

# Compare two specific versions
aicof diff 1.6.0 1.7.0

# Compare local vs specific version
aicof diff 1.7.0

# Show only component changes
aicof diff --installed --components
```

**Output shows:**

- Added files
- Removed files
- Modified files
- Summary statistics

---

### update

Update AICoF to the latest version.

**Usage:**

```bash
aicof update [flags]
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
aicof update --check

# Preview changes before updating
aicof update --diff

# Update to latest
aicof update

# Update to specific version
aicof update --version 1.7.0

# Force update without prompts
aicof update --force
```

---

### doctor

Check installation health and diagnose issues.

**Usage:**

```bash
aicof doctor [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--fix` | Attempt to automatically fix issues |

**Examples:**

```bash
# Run health check
aicof doctor

# Auto-fix issues
aicof doctor --fix
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
aicof version [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--check` | Check for newer versions available |

**Examples:**

```bash
# Show version
aicof version

# Check for updates
aicof version --check
```

---

## Common Workflows

### Setting Up a New Project

```bash
# Create and enter project directory
mkdir my-project && cd my-project

# Initialize with defaults
aicof init

# Or with specific components
aicof init --languages typescript,python --frameworks react,fastapi
```

### Discovering Components

```bash
# Search for what you need
aicof search api

# Get details
aicof info fw fastapi --preview 20

# Add it
aicof add fw fastapi
```

### Keeping Up to Date

```bash
# Check what's new
aicof diff --installed

# Preview changes
aicof update --check

# Update
aicof update
```

### Troubleshooting

```bash
# Check health
aicof doctor

# Fix issues
aicof doctor --fix

# Verbose output for debugging
aicof --verbose doctor
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
