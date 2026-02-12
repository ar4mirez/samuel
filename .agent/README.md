# .agent/ - Project-Specific AI Context

This directory contains project-specific context for the AICoF framework development.

## Structure

```
.agent/
├── README.md              # This file
├── project.md             # Project architecture and tech stack
├── patterns.md            # Coding patterns and conventions
├── state.md               # Current work state
├── skills/                # Agent Skills (language guides, workflows, utilities)
│   ├── go-guide/          # Go language guide (this is a Go project)
│   │   └── SKILL.md
│   ├── initialize-project/ # Workflow skills (15 total)
│   │   └── SKILL.md
│   ├── create-prd/
│   │   └── SKILL.md
│   ├── generate-tasks/
│   │   └── SKILL.md
│   └── ...                # code-review, security-audit, testing-strategy,
│                          # cleanup-project, refactoring, dependency-update,
│                          # troubleshooting, generate-agents-md, document-work,
│                          # update-framework, create-rfd, create-skill, commit-message
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
- Workflow skills (in `.agent/skills/`) are available on-demand for planning, reviews, etc.
- State and memory files track ongoing work and decisions
