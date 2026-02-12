---
title: FAQ
description: Frequently asked questions about AICoF
---

# Frequently Asked Questions

Common questions about AICoF (Artificial Intelligence Coding Framework).

---

## General

### Do I need to use Claude Code?

**No.** While designed for Claude Code, the system works with any AI coding assistant that reads instruction files. Use the AGENTS.md symlink for other tools:

```bash
ln -s CLAUDE.md AGENTS.md
```

### Which file should I edit - CLAUDE.md or AGENTS.md?

**Edit CLAUDE.md** (the source of truth). If using symlink, AGENTS.md updates automatically. If using standalone AGENTS.md, regenerate it after CLAUDE.md changes.

### Is this free to use?

**Yes.** AICoF is open source under the MIT license. Use it for personal or commercial projects.

---

## Setup

### What's the minimum setup required?

Just copy two items to your project:

```bash
cp CLAUDE.md /your-project/
cp -r .agent/ /your-project/
```

That's it. The system works immediately.

### Do I need to configure anything?

**No configuration required.** The system works out of the box. AI loads CLAUDE.md automatically and language guides auto-load based on file extensions.

### Can I use this with an existing project?

**Yes.** Use the initialize workflow:

```
@.agent/workflows/initialize-project.md

This is an existing project - analyze the codebase
```

---

## Languages

### My language isn't covered - what do I do?

The core guardrails in CLAUDE.md are **language-agnostic**. 90% of the rules still apply:

- Functions ≤50 lines ✓
- Files ≤300 lines ✓
- Input validation ✓
- Test coverage >80% ✓
- Conventional commits ✓

You can still use the system effectively. Consider contributing a language guide!

### How do I add a new language guide?

1. Copy an existing guide from `.agent/skills/<language>-guide/`
2. Create a new skill directory (e.g., `java-guide/SKILL.md`)
3. Adapt the content for your language
4. AI will auto-load based on file extensions

### Which language guide should I use?

**You don't choose - AI auto-loads based on file extensions:**

| Extensions | Guide |
|------------|-------|
| `.ts`, `.tsx`, `.js`, `.jsx` | TypeScript |
| `.py` | Python |
| `.go` | Go |
| `.rs` | Rust |
| `.kt`, `.kts` | Kotlin |

---

## Guardrails

### Can I customize the guardrails?

**Yes.** Edit CLAUDE.md for your team. Common customizations:

- File length limits: `300 → 500`
- Coverage targets: `80% → 90%`
- Commit format: Add your own types

### What if AI doesn't follow guardrails?

Explicitly remind it:

```
Follow the guardrails in CLAUDE.md. Specifically:
- Functions must be ≤50 lines
- All inputs must be validated
- Add tests for the new code
```

### Are guardrails enforced automatically?

AI follows guardrails when it loads CLAUDE.md. For strict enforcement, add:

- Linter rules
- Pre-commit hooks
- CI checks

---

## Workflows

### Is the PRD workflow required?

**No.** PRD is only for complex features (>10 files, new subsystems). Most work uses ATOMIC or FEATURE mode without PRD.

### When should I use the PRD workflow?

Use PRD when:

- Feature affects >10 files
- Building a new subsystem
- Requirements are unclear
- Multiple stakeholders need alignment

### Can I skip steps in a workflow?

**Generally no.** Workflows are designed as complete processes. Skipping steps often leads to issues later. If a workflow doesn't fit, use a simpler mode instead.

---

## Files & Structure

### Do I need to read all 500 lines of CLAUDE.md?

**No.** The Quick Reference section (lines 7-30) covers 80% of daily use. Read specific sections when you need guidance.

### What's the .agent/ directory for?

Project-specific context that grows over time:

- `project.md` - Your tech stack
- `patterns.md` - Your conventions
- `skills/<lang>-guide/` - Language-specific rules
- `workflows/` - Structured processes
- `memory/` - Decision logs

### Should I commit .agent/ files to git?

**Template files: Yes.** Language guides, workflows, README.

**Generated files: Team decision.** Some teams commit `project.md` and `patterns.md` to share context. Others keep them local.

---

## Methodology

### What's the difference between ATOMIC, FEATURE, and COMPLEX?

| Mode | Files | Approach |
|------|-------|----------|
| ATOMIC | <5 | Direct implementation |
| FEATURE | 5-10 | Break into subtasks |
| COMPLEX | >10 | PRD + task breakdown |

AI auto-detects which mode based on your request.

### What's the 4D methodology?

Four phases for every task:

1. **Deconstruct** - Break down the task
2. **Diagnose** - Identify risks
3. **Develop** - Implement with tests
4. **Deliver** - Validate and commit

### How do I know which mode AI is using?

AI usually indicates the mode:

- "This is a small fix, I'll implement directly" → ATOMIC
- "Let me break this into subtasks" → FEATURE
- "This is complex, want to create a PRD?" → COMPLEX

---

## Troubleshooting

### AI is not loading CLAUDE.md

1. Verify file is in project root
2. Check file name is exactly `CLAUDE.md`
3. Restart your AI tool
4. Explicitly reference: "Load instructions from CLAUDE.md"

### AI is ignoring my language guide

1. Check file extension matches
2. Explicitly load: "Load TypeScript guide from `.agent/skills/typescript-guide/SKILL.md`"

### AI is not following commits convention

Remind it:

```
Use conventional commits:
- feat(scope): description
- fix(scope): description
```

### I'm stuck and troubleshooting isn't helping

1. Use the troubleshooting workflow: `@.agent/workflows/troubleshooting.md`
2. Document what you've tried
3. Ask with a clear problem statement
4. Consider asking a human for fresh perspective

---

## Contributing

### How do I contribute a language guide?

1. Fork the repository
2. Copy an existing guide as template
3. Adapt for your language
4. Submit a pull request

### How do I report an issue?

Open an issue on GitHub:

- [GitHub Issues](https://github.com/ar4mirez/aicof/issues)

### How do I suggest improvements?

Open a discussion:

- [GitHub Discussions](https://github.com/ar4mirez/aicof/discussions)

---

## Still Have Questions?

- **Documentation**: Browse this site
- **GitHub Issues**: [Report bugs](https://github.com/ar4mirez/aicof/issues)
- **Discussions**: [Ask questions](https://github.com/ar4mirez/aicof/discussions)
