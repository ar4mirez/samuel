# .agent/ - Project-Specific AI Context

This directory contains project-specific context for the AICoF framework development.

## Structure

```
.agent/
├── README.md              # This file
├── project.md             # Project architecture and tech stack
├── patterns.md            # Coding patterns and conventions
├── state.md               # Current work state
├── skills/                # Language-specific guardrails
│   └── go-guide/          # Go language guide (this is a Go project)
│       └── SKILL.md
├── framework-guides/      # Framework-specific templates (empty - no framework used)
├── workflows/             # All 13 workflows
│   ├── initialize-project.md
│   ├── create-prd.md
│   ├── generate-tasks.md
│   ├── code-review.md
│   ├── security-audit.md
│   ├── testing-strategy.md
│   ├── cleanup-project.md
│   ├── refactoring.md
│   ├── dependency-update.md
│   ├── troubleshooting.md
│   ├── generate-agents-md.md
│   ├── document-work.md
│   └── update-framework.md
├── tasks/                 # PRDs and task lists
│   ├── 0001-prd-*.md      # PRD documents
│   └── tasks-*.md         # Task breakdowns
└── memory/                # Decision logs and learnings
    └── YYYY-MM-DD-*.md    # Dated memory entries
```

## Dogfooding

This project uses the AICoF framework for its own development. The framework files are:

- Root `CLAUDE.md` - Main AI instructions
- Root `AI_INSTRUCTIONS.md` - Quick start guide
- `.agent/` directory - Project context

The distributable template files are located in `template/.agent/` for distribution via the CLI.

## Usage

- Language guide (`go-guide/SKILL.md`) is automatically loaded when working on Go files
- Workflows are available on-demand for planning, reviews, etc.
- State and memory files track ongoing work and decisions
