# .agent/ - Project-Specific AI Context

This directory contains project-specific context for the AICoF framework development.

## Structure

```
.agent/
├── README.md           # This file
├── tasks/              # PRDs and task lists for this project
│   ├── 0001-prd-*.md   # PRD documents
│   └── tasks-*.md      # Task breakdowns
└── memory/             # Decision logs and learnings
    └── YYYY-MM-DD-*.md # Dated memory entries
```

## Note

The distributable template files (language guides, framework guides, workflows) are located in `template/.agent/`.

This separation allows:
1. This project to use AICoF for its own development
2. The template files to be distributed separately to users via the CLI
