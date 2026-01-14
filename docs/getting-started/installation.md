---
title: Installation
description: Detailed installation instructions for AICoF
---

# Installation

Multiple ways to add AICoF to your project.

---

## Option 1: Direct Copy (Recommended)

The simplest approach - copy the files directly:

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
    ├── project.md.template      # Template for project.md
    ├── state.md.template        # Template for state.md
    ├── language-guides/         # Auto-load based on file type
    │   ├── typescript.md
    │   ├── python.md
    │   ├── go.md
    │   ├── rust.md
    │   └── kotlin.md
    ├── workflows/               # On-demand workflows
    │   ├── create-prd.md
    │   ├── generate-tasks.md
    │   ├── initialize-project.md
    │   ├── troubleshooting.md
    │   └── generate-agents-md.md
    ├── tasks/                   # PRDs and task lists (created on demand)
    └── memory/                  # Decision logs (created on demand)
```

---

## Option 2: Git Subtree

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

## Option 3: Download ZIP

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

After installation, verify the files are in place:

```bash
# Check CLAUDE.md exists
cat CLAUDE.md | head -20

# Check .agent directory structure
ls -la .agent/

# Check language guides
ls .agent/language-guides/

# Check workflows
ls .agent/workflows/
```

Expected output:

```
.agent/
├── README.md
├── language-guides/
│   ├── README.md
│   ├── go.md
│   ├── kotlin.md
│   ├── python.md
│   ├── rust.md
│   └── typescript.md
├── memory/
├── tasks/
└── workflows/
    ├── README.md
    ├── create-prd.md
    ├── generate-agents-md.md
    ├── generate-tasks.md
    ├── initialize-project.md
    └── troubleshooting.md
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

### If Using Direct Copy

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

### If Using Subtree

```bash
git subtree pull --prefix=.ai-template \
    https://github.com/ar4mirez/aicof.git main --squash

# Copy updated files
cp .ai-template/CLAUDE.md ./
cp -r .ai-template/.agent ./
```

---

## Troubleshooting

### Files Not Loading

If AI doesn't seem to follow guardrails:

1. Verify `CLAUDE.md` is in project root
2. Check file permissions: `chmod 644 CLAUDE.md`
3. Explicitly remind AI: "Follow guardrails in CLAUDE.md"

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

-   :material-play:{ .lg .middle } **Your First Task**

    ---

    Try the system with a real task.

    [:octicons-arrow-right-24: First Task](first-task.md)

</div>
