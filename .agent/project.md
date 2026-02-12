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

### CLI Tool (Go)

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
aicof/
├── cmd/aicof/                # CLI entry point (main.go)
├── internal/                 # Go implementation
│   ├── commands/             # Command implementations (8 commands)
│   ├── core/                 # Core packages (config, registry, extractor)
│   ├── github/               # GitHub API client
│   └── ui/                   # User interface helpers (prompts, spinner)
├── template/                 # Distributable template files
│   ├── CLAUDE.md             # Main AI instructions
│   ├── AI_INSTRUCTIONS.md    # Quick start guide
│   └── .agent/               # AI context directory
│       ├── skills/           # 21 language guides + 33 framework skills
│       └── workflows/        # 13 structured workflows
├── .agent/                   # Project-specific AI context (dogfooding)
│   ├── skills/               # Go guide (this is a Go project)
│   ├── workflows/            # All 13 workflows
│   ├── tasks/                # PRDs and task lists
│   └── memory/               # Decision logs
├── docs/                     # MkDocs documentation source
├── CLAUDE.md                 # Framework instructions (dogfooding)
├── AI_INSTRUCTIONS.md        # Quick start guide (dogfooding)
├── aicof.yaml                # Framework config (dogfooding)
├── go.mod                    # Go module definition
├── Makefile                  # Build targets
└── .goreleaser.yaml          # Release automation
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

### Decision: Go Project Structure at Root

**Date**: 2026-01-14
**Context**: Initial structure had Go code in `packages/cli/` (Node.js monorepo pattern), not idiomatic Go
**Options Considered**:

1. Keep `packages/cli/` structure - Cons: Confusing for Go developers, non-standard
2. Move Go code to root with standard layout - Pros: Idiomatic, cleaner imports

**Decision**: Standard Go layout at repository root (`cmd/`, `internal/`)
**Rationale**: Go projects typically have `go.mod` at root with `cmd/` and `internal/` directories
**Implementation**: Renamed `internal/cmd` to `internal/commands` to avoid confusion with root `cmd/`

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
make test             # Run all tests
go test -cover ./...  # With coverage
```

---

## Development Setup

```bash
# Clone repo
git clone https://github.com/ar4mirez/aicof.git
cd aicof

# Build CLI
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
