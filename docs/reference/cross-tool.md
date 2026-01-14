---
title: Cross-Tool Compatibility
description: Using AICoF with multiple AI coding assistants
---

# Cross-Tool Compatibility

AICoF works with any AI coding assistant through the AGENTS.md standard.

---

## Supported Tools

| Tool | Primary File | Fallback | Status |
|------|--------------|----------|--------|
| **Claude Code** | CLAUDE.md | AGENTS.md | Full support |
| **Cursor** | .cursorrules, AGENTS.md | - | Full support |
| **OpenAI Codex** | AGENTS.md | - | Full support |
| **GitHub Copilot** | AGENTS.md | - | Partial |
| **Google Jules** | AGENTS.md | - | Full support |
| **Codeium** | AGENTS.md | - | Full support |
| **Amazon Q** | AGENTS.md | - | Partial |

---

## Setup Options

### Option 1: Symlink (Recommended)

Creates a link so both files are always identical:

=== "macOS/Linux"

    ```bash
    ln -s CLAUDE.md AGENTS.md
    ```

=== "Windows (Admin PowerShell)"

    ```powershell
    New-Item -ItemType SymbolicLink -Path AGENTS.md -Target CLAUDE.md
    ```

**Pros**:

- Zero maintenance
- Always in sync
- No duplication

**Cons**:

- Windows requires admin privileges
- Some tools may not follow symlinks

### Option 2: Generate Standalone

Create a separate AGENTS.md with operations-only content:

```
@.agent/workflows/generate-agents-md.md
```

**Pros**:

- Works everywhere
- Can customize per-tool
- Simpler content for some tools

**Cons**:

- Manual updates needed
- Potential drift from CLAUDE.md

### Option 3: Copy

Simply copy CLAUDE.md:

```bash
cp CLAUDE.md AGENTS.md
```

**Pros**:

- Works everywhere
- Simple

**Cons**:

- Manual sync needed
- Easy to forget updates

---

## Tool-Specific Setup

### Claude Code

Claude Code reads CLAUDE.md automatically. No additional setup needed.

```
your-project/
├── CLAUDE.md       ← Loaded automatically
└── .agent/         ← Language guides, workflows
```

### Cursor

Cursor reads `.cursorrules` or `AGENTS.md`:

**Option A**: Symlink CLAUDE.md

```bash
ln -s CLAUDE.md .cursorrules
# or
ln -s CLAUDE.md AGENTS.md
```

**Option B**: Create .cursorrules

```
@.agent/workflows/generate-agents-md.md

Generate .cursorrules for Cursor
```

### GitHub Copilot

Copilot can read AGENTS.md for context:

```bash
ln -s CLAUDE.md AGENTS.md
```

!!! note "Limited Support"

    Copilot has limited instruction-following compared to Claude Code. Guardrails may not be strictly enforced.

### OpenAI Codex / ChatGPT

Use AGENTS.md:

```bash
ln -s CLAUDE.md AGENTS.md
```

---

## Multi-Tool Teams

When team members use different tools:

### Recommended Setup

1. Keep CLAUDE.md as source of truth
2. Create symlink: `ln -s CLAUDE.md AGENTS.md`
3. Add both to git

```bash
git add CLAUDE.md AGENTS.md
git commit -m "chore: add AI coding instructions"
```

### Team Workflow

1. **Edit only CLAUDE.md** - Never edit AGENTS.md directly
2. **Symlink stays updated** - Changes propagate automatically
3. **All tools see same rules** - Consistent behavior

### For Windows Users

If symlinks don't work:

```
@.agent/workflows/generate-agents-md.md

Generate standalone AGENTS.md
```

Add to `.gitattributes`:

```gitattributes
AGENTS.md merge=ours
```

This prevents merge conflicts when regenerating.

---

## File Comparison

### CLAUDE.md (Full - ~500 lines)

Contains everything:

- Operations (commands)
- Boundaries (protected files)
- Quick Reference
- Core Guardrails (35+ rules)
- 4D Methodology
- SDLC stages
- Context System
- Anti-Patterns

### AGENTS.md (Operations - ~200 lines)

Contains essentials:

- Operations (commands)
- Boundaries (protected files)
- Quick Reference

---

## Maintaining Sync

### With Symlink

Automatic - no action needed.

### With Standalone AGENTS.md

Add to your workflow:

1. **PR Checklist**:
   - [ ] If CLAUDE.md changed, regenerate AGENTS.md

2. **CI Check** (optional):
   ```yaml
   - name: Check AGENTS.md is current
     run: |
       # Generate expected AGENTS.md
       # Compare with committed version
       # Fail if different
   ```

3. **Pre-commit hook** (optional):
   ```bash
   #!/bin/bash
   if git diff --cached --name-only | grep -q "CLAUDE.md"; then
     echo "CLAUDE.md changed - regenerate AGENTS.md"
   fi
   ```

---

## Troubleshooting

### Symlink Not Working

**Windows**: Requires admin privileges or Developer Mode enabled.

**Git**: Ensure `core.symlinks=true`:

```bash
git config core.symlinks true
```

### Tool Ignoring Instructions

1. Check file location (must be in project root)
2. Check file name (exact match required)
3. Restart the tool
4. Explicitly reference: "Follow instructions in CLAUDE.md"

### Different Behavior Between Tools

Each tool interprets instructions differently. For consistent results:

1. Keep rules specific and testable
2. Use guardrails that can be verified
3. Test with each tool your team uses

---

## Best Practices

### Do

- [x] Use symlink when possible
- [x] Keep CLAUDE.md as source of truth
- [x] Test with all tools your team uses
- [x] Document which tools are supported

### Don't

- [ ] Edit AGENTS.md directly (edit CLAUDE.md)
- [ ] Maintain two different instruction sets
- [ ] Assume all tools behave identically
- [ ] Ignore tool-specific limitations

---

## Related

<div class="grid cards" markdown>

-   :material-file-document:{ .lg .middle } **CLAUDE.md**

    ---

    The source of truth for instructions.

    [:octicons-arrow-right-24: CLAUDE.md](../core/claude-md.md)

-   :material-cog:{ .lg .middle } **Generate AGENTS.md**

    ---

    Workflow for generating standalone file.

    [:octicons-arrow-right-24: Generate AGENTS.md](../workflows/generate-agents-md.md)

</div>
