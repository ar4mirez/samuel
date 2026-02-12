# Create Skill Workflow

> **Purpose**: Create reusable AI agent capability modules following the [Agent Skills](https://agentskills.io) open standard. Skills are portable across 25+ agent products including Claude Code, Cursor, and VS Code.

---

## When to Use

| Trigger | Description |
|---------|-------------|
| **New Capability** | Creating a reusable AI agent capability |
| **Scaffold Skill** | Setting up the directory structure for a new skill |
| **Convert Guide** | Converting an existing guide or workflow to Agent Skills format |
| **Cross-Tool Sharing** | Building a capability module portable across AI tools |

---

## Process

### Step 1: Define the Skill

Before creating the skill, clarify its purpose:

1. **What capability** does this skill provide?
    - Example: "Process PDF files", "Generate API documentation", "Manage database migrations"

2. **When should it activate?**
    - What user requests should trigger this skill?
    - What keywords or contexts indicate this skill is needed?

3. **What resources does it need?**
    - Scripts for complex operations?
    - Reference documentation?
    - Templates or assets?

4. **What's the scope?**
    - Keep skills focused on one capability
    - Split large skills into multiple smaller ones

### Step 2: Scaffold the Skill

```bash
samuel skill create <skill-name>
```

**Name requirements:**

- Lowercase alphanumeric and hyphens only
- No consecutive hyphens (`--`)
- Cannot start or end with hyphen
- Maximum 64 characters

This creates:

```text
.claude/skills/<skill-name>/
├── SKILL.md           # Pre-filled template
├── scripts/           # For executable code
├── references/        # For additional docs
└── assets/            # For templates, data
```

### Step 3: Write SKILL.md

Edit the generated SKILL.md with required YAML frontmatter:

```yaml
---
name: skill-name
description: |
  What this skill does and when to use it.
  Include specific triggers and keywords.
license: MIT
metadata:
  author: your-name
  version: "1.0"
---
```

**Description best practices:**

- Describe both *what* and *when*
- Include keywords that trigger activation
- Be specific ("Process PDF files" not "Helps with documents")
- Maximum 1024 characters

Then write the body content with:

1. **Purpose**: What capability this provides
2. **When to Use**: Scenarios that trigger this skill
3. **Instructions**: Step-by-step procedure
4. **Examples**: Input/output pairs
5. **Notes**: Warnings, edge cases, best practices

Keep the body **under 500 lines** — use reference files for detailed content.

### Step 4: Add Resources (Optional)

| Directory | Purpose | Example |
|-----------|---------|---------|
| `scripts/` | Executable code for deterministic operations | `process.py`, `migrate.sh` |
| `references/` | Documentation AI loads on-demand | `api-guide.md`, `schema.md` |
| `assets/` | Templates and data files | `template.html`, `config.json` |

### Step 5: Validate

```bash
# Validate specific skill
samuel skill validate <skill-name>

# Validate all skills
samuel skill validate
```

**Validation checks:**

- SKILL.md exists with valid YAML frontmatter
- Name matches directory name
- Name format is correct (lowercase, hyphens, max 64 chars)
- Description is present and under 1024 characters
- Compatibility field under 500 characters (if present)

### Step 6: Test

1. Load the skill in your AI agent
2. Try scenarios from "When to Use"
3. Verify instructions are followed correctly
4. Check that examples produce expected output
5. Test edge cases

---

## Best Practices

### Concise is Key

The context window is a shared resource. Only include what the AI doesn't already know:

```markdown
<!-- Good (50 tokens) -->
## Extract PDF Text
Use pdfplumber:
import pdfplumber
with pdfplumber.open("file.pdf") as pdf:
    text = pdf.pages[0].extract_text()

<!-- Bad (150 tokens) -->
## Extract PDF Text
PDF files are a common format... [unnecessary explanation]
```

### Set Appropriate Freedom

Match specificity to task fragility:

| Freedom Level | When to Use | Example |
|--------------|-------------|---------|
| High | Multiple valid approaches | Code review process |
| Medium | Preferred pattern exists | Report generation |
| Low | Fragile/critical operations | Database migrations |

### Use Progressive Disclosure

1. **Metadata** (~100 tokens): Always loaded by the agent
2. **SKILL.md body** (<5000 tokens): Loaded on activation
3. **References/Scripts**: Loaded on-demand

Keep SKILL.md lean; move details to reference files.

---

## Checklist

Before finalizing your skill:

- [ ] Name follows conventions (lowercase, hyphens, max 64 chars)
- [ ] Description is specific and under 1024 chars
- [ ] SKILL.md body is under 500 lines
- [ ] Instructions are clear and step-by-step
- [ ] Examples show input/output pairs
- [ ] Validation passes (`samuel skill validate`)
- [ ] Tested with real scenarios
- [ ] Scripts handle errors gracefully (if applicable)

---

## See Also

- [CLI Reference: skill](../reference/cli.md#skill) - Full command reference
- [Agent Skills Specification](https://agentskills.io/specification) - The open standard
- [The .claude Directory](../core/agent-directory.md) - Where skills live
