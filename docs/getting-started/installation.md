---
title: Installation
description: Detailed installation instructions for AICoF
---

# Installation

Multiple ways to install AICoF. The CLI is recommended for the best experience.

---

## Option 1: CLI Installation (Recommended)

The AICoF CLI provides the easiest way to install and manage the framework.

### macOS / Linux (Homebrew)

```bash
brew tap ar4mirez/tap && brew install aicof
```

### macOS / Linux (Curl)

```bash
curl -sSL https://raw.githubusercontent.com/ar4mirez/aicof/main/install.sh | sh
```

### Go Install

```bash
go install github.com/ar4mirez/aicof/cmd/aicof@latest
```

### Verify CLI Installation

```bash
# Check version
aicof version

# Run health check
aicof doctor
```

### Initialize Your Project

```bash
cd your-project
aicof init
```

The CLI will interactively guide you through setup:

1. Create `CLAUDE.md` with all guardrails
2. Create `.agent/` directory structure
3. Select languages, frameworks, and workflows
4. Set up cross-tool compatibility (AGENTS.md symlink)

---

## Option 2: Direct Copy (Manual)

For environments where the CLI isn't available:

```bash
# Clone the repository
git clone https://github.com/ar4mirez/aicof.git

# Copy to your project
cp -r aicof/{CLAUDE.md,.agent} /path/to/your-project/

# Clean up
rm -rf aicof
```

### What Gets Copied

```
your-project/
├── CLAUDE.md                    # Core guardrails (~500 lines)
└── .agent/                      # Project context directory
    ├── README.md                # How to use .agent/
    ├── skills/                  # Language guides + framework skills (54 total)
    │   ├── <lang>-guide/        # 21 language guide skills
    │   │   ├── SKILL.md
    │   │   └── references/
    │   └── <framework>/         # 33 framework skills
    │       ├── SKILL.md
    │       └── references/
    ├── workflows/               # 13 workflows
    ├── tasks/                   # PRDs and task lists (created on demand)
    └── memory/                  # Decision logs (created on demand)
```

---

## Option 3: Git Subtree

Keep AICoF as a subtree for easier updates:

```bash
# Add as subtree (first time)
git subtree add --prefix=.ai-template \
    https://github.com/ar4mirez/aicof.git main --squash

# Copy files to root
cp .ai-template/CLAUDE.md ./
cp -r .ai-template/.agent ./

# Update later
git subtree pull --prefix=.ai-template \
    https://github.com/ar4mirez/aicof.git main --squash
```

!!! note "Subtree Benefits"

    - Easy updates with `git subtree pull`
    - Full history preserved
    - No submodule complexity

---

## Option 4: Download ZIP

For projects without git or one-time use:

1. Go to [GitHub Releases](https://github.com/ar4mirez/aicof/releases)
2. Download the latest release ZIP
3. Extract `CLAUDE.md` and `.agent/` to your project

---

## Cross-Tool Setup

If your team uses multiple AI tools (Claude Code, Cursor, Codex, etc.):

### Option A: Symlink (Recommended)

```bash
ln -s CLAUDE.md AGENTS.md
```

Both files stay in sync automatically.

### Option B: Generate Standalone

```
@.agent/workflows/generate-agents-md.md
```

Creates a separate `AGENTS.md` with operations-only content.

!!! info "When to Use Standalone"

    - Team members use tools that don't support CLAUDE.md
    - You want different content for different tools
    - CI/CD tools need simpler instructions

---

## Verify Installation

### With CLI (Recommended)

```bash
# Comprehensive health check
aicof doctor

# List installed components
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

### Manual Verification

```bash
# Check CLAUDE.md exists
cat CLAUDE.md | head -20

# Check .agent directory structure
ls -la .agent/

# Check language guides
ls .agent/skills/

# Check workflows
ls .agent/workflows/
```

---

## Post-Installation: Add Components

After installation, add the components you need:

### Using CLI

```bash
# Search for what you need
aicof search react

# Preview before installing
aicof info framework react --preview 20

# Add components
aicof add language typescript
aicof add framework react
aicof add workflow code-review

# Use aliases for speed
aicof add lang go
aicof add fw nextjs
aicof add wf security-audit
```

### Check What's Available

```bash
# List all available components
aicof list --available

# Filter by type
aicof list --available --type frameworks
```

---

## Git Configuration

### Add to .gitignore (Optional)

You may want to ignore generated files:

```gitignore
# AICoF - generated files
.agent/project.md
.agent/patterns.md
.agent/state.md
.agent/tasks/*.md
!.agent/tasks/EXAMPLE-*.md
.agent/memory/*.md
!.agent/memory/.gitkeep
```

### Commit the Template Files

```bash
# Add template files
git add CLAUDE.md .agent/

# Commit
git commit -m "chore: add AICoF for AI-assisted development"
```

---

## Update to Latest Version

### Update with CLI

```bash
# Check what's changed
aicof diff --installed

# Preview changes before updating
aicof update --check

# Update to latest
aicof update

# Or update to specific version
aicof update --version 1.7.0
```

### Manual Update

```bash
# Clone latest version
git clone https://github.com/ar4mirez/aicof.git temp-update

# Backup your customizations
cp CLAUDE.md CLAUDE.md.backup
cp -r .agent .agent.backup

# Copy new files
cp temp-update/CLAUDE.md ./
cp -r temp-update/.agent ./

# Review and merge your customizations
# (manually compare .backup files with new ones)

# Clean up
rm -rf temp-update
rm CLAUDE.md.backup
rm -rf .agent.backup
```

---

## Troubleshooting

### CLI Not Found

If `aicof` command is not found after installation:

```bash
# Check if it's in PATH
which aicof

# Add to PATH (bash/zsh)
export PATH="$PATH:/usr/local/bin"

# Or reinstall
curl -sSL https://raw.githubusercontent.com/ar4mirez/aicof/main/install.sh | sh
```

### Files Not Loading

If AI doesn't seem to follow guardrails:

1. Verify `CLAUDE.md` is in project root
2. Check file permissions: `chmod 644 CLAUDE.md`
3. Explicitly remind AI: "Follow guardrails in CLAUDE.md"
4. Run `aicof doctor` to check for issues

### Symlink Issues on Windows

Windows requires administrator privileges for symlinks. Use Option B (Generate Standalone) instead:

```
@.agent/workflows/generate-agents-md.md
```

### Large Repository

For large repos, you might want to exclude documentation from searches:

```gitattributes
# .gitattributes
.agent/** linguist-documentation
CLAUDE.md linguist-documentation
```

---

## Next Steps

<div class="grid cards" markdown>

-   :material-rocket-launch:{ .lg .middle } **Quick Start**

    ---

    Get coding in 60 seconds.

    [:octicons-arrow-right-24: Quick Start](quick-start.md)

-   :material-console:{ .lg .middle } **CLI Reference**

    ---

    Learn all 11 CLI commands.

    [:octicons-arrow-right-24: CLI Commands](../reference/cli.md)

-   :material-play:{ .lg .middle } **Your First Task**

    ---

    Try the system with a real task.

    [:octicons-arrow-right-24: First Task](first-task.md)

</div>
