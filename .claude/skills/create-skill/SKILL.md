---
name: create-skill
description: |
  Agent Skill creation workflow. Use when creating new reusable AI agent skills,
  scaffolding skill directories, or converting existing guides into the
  portable Agent Skills standard format.
license: MIT
metadata:
  author: aicof
  version: "1.0"
  category: workflow
---

# Agent Skill Creation

[Agent Skills](https://agentskills.io) are capability modules that give AI agents new abilities. Unlike workflows (which guide how to approach tasks), skills define what an agent can do. Skills are portable across 25+ agent products including Claude Code, Cursor, and VS Code.

## When to Use

Use this skill when you need to:
- Create a new reusable AI agent capability
- Scaffold a skill directory structure
- Convert an existing guide or workflow into the Agent Skills format
- Build a capability module for AI agents

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

Run the AICoF CLI command:

```bash
aicof skill create <skill-name>
```

**Name Requirements**:
- Lowercase alphanumeric and hyphens only
- No consecutive hyphens (`--`)
- Cannot start or end with hyphen
- Maximum 64 characters

This creates:
```
.claude/skills/<skill-name>/
├── SKILL.md           # Pre-filled template
├── scripts/           # For executable code
├── references/        # For additional docs
└── assets/            # For templates, data
```

### Step 3: Write SKILL.md

Edit the generated SKILL.md with:

#### Required Frontmatter

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

**Description Best Practices**:
- Describe both *what* and *when*
- Include keywords that trigger activation
- Be specific, not vague ("Process PDF files" not "Helps with documents")
- Maximum 1024 characters

#### Body Content

Write clear instructions that tell the AI agent how to use this skill:

1. **Purpose**: What capability this provides
2. **When to Use**: Scenarios that trigger this skill
3. **Instructions**: Step-by-step procedure
4. **Examples**: Input/output pairs
5. **Notes**: Warnings, edge cases, best practices

**Keep under 500 lines** - use reference files for detailed content.

### Step 4: Add Resources (Optional)

#### Scripts (`scripts/`)

Add executable code for deterministic operations:

```python
# scripts/process.py
def process_data(input_path):
    # Implementation
    pass
```

**Script Best Practices**:
- Make scripts self-contained
- Include helpful error messages
- Handle edge cases
- Document parameters

#### References (`references/`)

Add documentation that AI loads on-demand:

```markdown
# references/api-guide.md

## API Endpoints
...
```

Use references for:
- Detailed technical docs
- Domain-specific knowledge
- Large examples
- Configuration guides

#### Assets (`assets/`)

Add templates and data files:

```
assets/
├── template.html
├── config.json
└── icons/
```

### Step 5: Validate

Run validation to check against Agent Skills spec:

```bash
aicof skill validate <skill-name>
```

**Validation Checks**:
- SKILL.md exists with valid YAML frontmatter
- Name matches directory name
- Name format is correct
- Description is present
- No fields exceed length limits

Fix any errors before proceeding.

### Step 6: Test

Test the skill with real tasks:

1. Load the skill in your AI agent
2. Try scenarios from "When to Use"
3. Verify instructions are followed correctly
4. Check that examples produce expected output
5. Test edge cases

### Step 7: Document

If the skill is significant, record in `.claude/memory/`:

```markdown
# .claude/memory/YYYY-MM-DD-skill-name.md

## Skill Created: <skill-name>

**Purpose**: Brief description

**Key Decisions**:
- Why this approach
- Alternatives considered

**Testing Notes**:
- What worked
- What needed adjustment
```

---

## Best Practices

### Concise is Key

The context window is a shared resource. Only include what the AI doesn't already know:

**Good** (50 tokens):
```markdown
## Extract PDF Text

Use pdfplumber:
```python
import pdfplumber
with pdfplumber.open("file.pdf") as pdf:
    text = pdf.pages[0].extract_text()
```
```

**Bad** (150 tokens):
```markdown
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

1. **Metadata** (~100 tokens): Always loaded
2. **SKILL.md body** (<5000 tokens): Loaded on activation
3. **References/Scripts**: Loaded on-demand

Keep SKILL.md lean; move details to reference files.

### Provide Examples

Examples teach better than explanations:

```markdown
### Example: User Request

**Input**: "Help me merge these PDFs"

**Output**:
```python
# Code that handles the request
```
```

---

## Checklist

Before finalizing your skill:

- [ ] Name follows conventions (lowercase, hyphens, max 64 chars)
- [ ] Description is specific and under 1024 chars
- [ ] SKILL.md body is under 500 lines
- [ ] Instructions are clear and step-by-step
- [ ] Examples show input/output pairs
- [ ] Validation passes (`aicof skill validate`)
- [ ] Tested with real scenarios
- [ ] Scripts handle errors gracefully (if applicable)

---

## Related

- **Agent Skills Specification**: https://agentskills.io/specification
- **Example Skills**: https://github.com/anthropics/skills
- **Skills README**: `.claude/skills/README.md`
