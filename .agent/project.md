# Project: AICoF (Artificial Intelligence Coding Framework)

> **Purpose**: Document tech stack, architecture, and key design decisions
>
> **Created**: 2026-01-14
> **Last Updated**: 2026-01-14

---

## Tech Stack

### Languages & Runtimes
- **Primary Language**: Go 1.21+ (CLI tool)
- **Documentation**: Python 3.x (MkDocs)
- **Template Content**: Markdown

### CLI Tool (`packages/cli/`)
- **Framework**: Cobra (CLI framework)
- **Configuration**: gopkg.in/yaml.v3
- **HTTP Client**: net/http (standard library)
- **Archive Handling**: archive/tar, compress/gzip (standard library)
- **UI**: fatih/color, manifoldco/promptui, schollz/progressbar

### Documentation Site
- **Generator**: MkDocs with Material theme
- **Hosting**: GitHub Pages

---

## Architecture

### Repository Structure
```
ai-code-template/
├── template/                  # Distributable template files
│   ├── CLAUDE.md             # Main AI instructions
│   ├── AI_INSTRUCTIONS.md    # Alternative format
│   └── .agent/               # AI context directory
│       ├── language-guides/  # 21 language-specific guides
│       ├── framework-guides/ # 33 framework-specific guides
│       └── workflows/        # 13 structured workflows
├── .agent/                   # Project-specific AI context (this project)
│   ├── tasks/               # PRDs and task lists
│   └── memory/              # Decision logs
├── packages/cli/             # Go CLI tool
│   ├── cmd/aicof/           # Entry point
│   └── internal/            # Implementation
│       ├── cmd/             # Command implementations
│       ├── core/            # Core packages
│       ├── github/          # GitHub API client
│       └── ui/              # User interface helpers
└── docs/                     # MkDocs documentation source
```

### CLI Architecture
- **Commands**: init, update, add, remove, list, doctor, version
- **Core Flow**: User command → GitHub download → Tar extraction → Local files
- **Caching**: Downloaded versions cached in ~/.cache/aicof/
- **Configuration**: Project config in aicof.yaml

---

## Key Design Decisions

### Decision: Repository Restructure (template/ directory)
**Date**: 2026-01-14
**Context**: `.agent/` directory served dual purpose - both template files for distribution AND project-specific context for this repo's development.
**Options Considered**:
1. Single `.agent/` with complex exclusions - Cons: Confusing, error-prone
2. Separate `template/` directory - Pros: Clear separation, explicit paths

**Decision**: Separate `template/` directory
**Rationale**: Clear separation of distributable templates from project-specific context
**Trade-offs**: CLI paths need `template/` prefix, slight increase in path complexity

### Decision: Go for CLI
**Date**: 2026-01-14
**Context**: Needed standalone binary without runtime dependencies
**Decision**: Go with Cobra framework
**Rationale**: Single binary distribution, cross-platform, excellent CLI ecosystem

### Decision: GitHub Archive Downloads (no git)
**Date**: 2026-01-14
**Context**: Don't want to require git on user machines
**Decision**: Use GitHub's tar.gz archive API
**Rationale**: Works without git, smaller downloads, simpler implementation

---

## Testing Strategy

### Test Types
- **Unit**: Core packages (config, registry, extractor, github client)
- **E2E**: All 7 CLI commands with real GitHub API

### Coverage Targets
- Core business logic: 80%+
- Overall: 60%+

### Test Commands
```bash
cd packages/cli
make test           # Run all tests
go test -cover ./...  # With coverage
```

---

## Development Setup

```bash
# Clone repo
git clone https://github.com/ar4mirez/aicof.git
cd aicof

# Build CLI
cd packages/cli
make build
./bin/aicof version

# Run docs locally
pip install -r requirements-docs.txt
mkdocs serve
```

---

## Known Issues

### No GitHub Releases Yet
**Description**: Repository has no releases, so CLI falls back to `dev` version from main branch
**Workaround**: First release will create proper version tags
**Tracking**: Expected to resolve with first official release

---

## External Resources

- **GitHub Repository**: https://github.com/ar4mirez/aicof
- **Documentation**: https://ar4mirez.github.io/aicof/
