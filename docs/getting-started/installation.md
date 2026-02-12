---
title: Installation
description: Detailed installation instructions for Samuel
---

# Installation

Multiple ways to install Samuel. The CLI is recommended for the best experience.

---

## Option 1: CLI Installation (Recommended)

The Samuel CLI provides the easiest way to install and manage the framework.

### macOS / Linux (Homebrew)

```bash
brew tap ar4mirez/tap && brew install samuel
```

### macOS / Linux (Curl)

```bash
curl -sSL https://raw.githubusercontent.com/ar4mirez/samuel/main/install.sh | sh
```

### Go Install

```bash
go install github.com/ar4mirez/samuel/cmd/samuel@latest
```

### Verify CLI Installation

```bash
# Check version
samuel version

# Run health check
samuel doctor
```

### Initialize Your Project

```bash
cd your-project
samuel init
```

The CLI will interactively guide you through setup:

1. Create `CLAUDE.md` with all guardrails
2. Create `.claude/` directory structure
3. Select languages, frameworks, and workflows
4. Set up cross-tool compatibility (AGENTS.md symlink)

---

## Option 2: Direct Copy (Manual)

For environments where the CLI isn't available:

```bash
# Clone the repository
git clone https://github.com/ar4mirez/samuel.git

# Copy to your project
cp -r samuel/{CLAUDE.md,.claude} /path/to/your-project/

# Clean up
rm -rf samuel
```

### What Gets Copied

```
your-project/
├── CLAUDE.md                    # Core guardrails (~500 lines)
└── .claude/                      # Project context directory
    ├── README.md                # How to use .claude/
    ├── skills/                  # Language guides + framework skills (54 total)
    │   ├── <lang>-guide/        # 21 language guide skills
    │   │   ├── SKILL.md
    │   │   └── references/
    │   └── <framework>/         # 33 framework skills
    │       ├── SKILL.md
    │       └── references/
    │   ├── create-prd/          # 15 workflow skills
    │   │   └── SKILL.md
    │   └── ...                  # initialize-project, generate-tasks, code-review, etc.
    ├── tasks/                   # PRDs and task lists (created on demand)
    └── memory/                  # Decision logs (created on demand)
```

---

## Option 3: Git Subtree

Keep Samuel as a subtree for easier updates:

```bash
# Add as subtree (first time)
git subtree add --prefix=.ai-template \
    https://github.com/ar4mirez/samuel.git main --squash

# Copy files to root
cp .ai-template/CLAUDE.md ./
cp -r .ai-template/.claude ./

# Update later
git subtree pull --prefix=.ai-template \
    https://github.com/ar4mirez/samuel.git main --squash
```

!!! note "Subtree Benefits"

    - Easy updates with `git subtree pull`
    - Full history preserved
    - No submodule complexity

---

## Option 4: Download ZIP

For projects without git or one-time use:

1. Go to [GitHub Releases](https://github.com/ar4mirez/samuel/releases)
2. Download the latest release ZIP
3. Extract `CLAUDE.md` and `.claude/` to your project

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
@.claude/skills/generate-agents-md/SKILL.md
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
samuel doctor

# List installed components
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

### Manual Verification

```bash
# Check CLAUDE.md exists
cat CLAUDE.md | head -20

# Check .claude directory structure
ls -la .claude/

# Check language guides
ls .claude/skills/

# Check skills (includes workflow skills)
ls .claude/skills/
```

---

## Post-Installation: Add Components

After installation, add the components you need:

### Using CLI

```bash
# Search for what you need
samuel search react

# Preview before installing
samuel info framework react --preview 20

# Add components
samuel add language typescript
samuel add framework react
samuel add workflow code-review

# Use aliases for speed
samuel add lang go
samuel add fw nextjs
samuel add wf security-audit
```

### Check What's Available

```bash
# List all available components
samuel list --available

# Filter by type
samuel list --available --type frameworks
```

---

## Git Configuration

### Add to .gitignore (Optional)

You may want to ignore generated files:

```gitignore
# Samuel - generated files
.claude/project.md
.claude/patterns.md
.claude/state.md
.claude/tasks/*.md
!.claude/tasks/EXAMPLE-*.md
.claude/memory/*.md
!.claude/memory/.gitkeep
```

### Commit the Template Files

```bash
# Add template files
git add CLAUDE.md .claude/

# Commit
git commit -m "chore: add Samuel for AI-assisted development"
```

---

## Update to Latest Version

### Update with CLI

```bash
# Check what's changed
samuel diff --installed

# Preview changes before updating
samuel update --check

# Update to latest
samuel update

# Or update to specific version
samuel update --version 1.7.0
```

### Manual Update

```bash
# Clone latest version
git clone https://github.com/ar4mirez/samuel.git temp-update

# Backup your customizations
cp CLAUDE.md CLAUDE.md.backup
cp -r .claude .claude.backup

# Copy new files
cp temp-update/CLAUDE.md ./
cp -r temp-update/.claude ./

# Review and merge your customizations
# (manually compare .backup files with new ones)

# Clean up
rm -rf temp-update
rm CLAUDE.md.backup
rm -rf .claude.backup
```

---

## Troubleshooting

### CLI Not Found

If `samuel` command is not found after installation:

```bash
# Check if it's in PATH
which samuel

# Add to PATH (bash/zsh)
export PATH="$PATH:/usr/local/bin"

# Or reinstall
curl -sSL https://raw.githubusercontent.com/ar4mirez/samuel/main/install.sh | sh
```

### Files Not Loading

If AI doesn't seem to follow guardrails:

1. Verify `CLAUDE.md` is in project root
2. Check file permissions: `chmod 644 CLAUDE.md`
3. Explicitly remind AI: "Follow guardrails in CLAUDE.md"
4. Run `samuel doctor` to check for issues

### Symlink Issues on Windows

Windows requires administrator privileges for symlinks. Use Option B (Generate Standalone) instead:

```
@.claude/skills/generate-agents-md/SKILL.md
```

### Large Repository

For large repos, you might want to exclude documentation from searches:

```gitattributes
# .gitattributes
.claude/** linguist-documentation
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
