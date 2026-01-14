# AICoF CLI

Command-line interface for managing the [AICoF (Artificial Intelligence Coding Framework)](https://github.com/ar4mirez/aicof).

## Features

- **Initialize projects** with AI coding guardrails
- **Update framework** versions with smart diffing
- **Manage components** (language guides, framework guides, workflows)
- **Health checks** to verify installation integrity
- **No runtime dependencies** - single binary distribution

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap ar4mirez/tap
brew install aicof
```

### Curl (macOS/Linux)

```bash
curl -sSL https://raw.githubusercontent.com/ar4mirez/aicof/main/install.sh | sh
```

### Direct Download

Download the latest release from [GitHub Releases](https://github.com/ar4mirez/aicof/releases).

### From Source

```bash
git clone https://github.com/ar4mirez/aicof.git
cd aicof/packages/cli
make install
```

## Quick Start

```bash
# Initialize a new project
aicof init my-project

# Or initialize in current directory
aicof init .

# Check installation health
aicof doctor

# Update to latest framework version
aicof update

# List installed components
aicof list
```

## Commands

### `aicof init [project-name]`

Initialize AICoF framework in a new or existing project.

```bash
aicof init my-project              # Create new project
aicof init .                       # Initialize in current directory
aicof init --template minimal      # Use minimal template (just CLAUDE.md)
aicof init --template starter      # Use starter template (TypeScript, Python, Go)
aicof init --languages ts,py,go    # Select specific languages
aicof init --non-interactive       # Skip prompts
```

**Templates:**
- `full` - All 21 language guides, 33 framework guides, 13 workflows
- `starter` - Core files + TypeScript, Python, Go
- `minimal` - Just CLAUDE.md and workflows

### `aicof update`

Update the framework to the latest version.

```bash
aicof update                # Update to latest
aicof update --check        # Check for updates without applying
aicof update --diff         # Show what will change
aicof update --force        # Overwrite local modifications
aicof update --version 1.8.0  # Update to specific version
```

### `aicof add <type> <name>`

Add a component to your project.

```bash
aicof add language rust     # Add Rust language guide
aicof add framework django  # Add Django framework guide
aicof add workflow security-audit  # Add security audit workflow
```

### `aicof remove <type> <name>`

Remove a component from your project.

```bash
aicof remove language rust
aicof remove framework django
```

### `aicof list`

List installed or available components.

```bash
aicof list                    # List installed components
aicof list --available        # List all available components
aicof list --type languages   # Filter by type
```

### `aicof doctor`

Check installation health and fix issues.

```bash
aicof doctor        # Run health check
aicof doctor --fix  # Auto-fix issues
```

### `aicof version`

Show version information.

```bash
aicof version          # Show versions
aicof version --check  # Check for CLI updates
```

## Configuration

AICoF creates an `aicof.yaml` file in your project root:

```yaml
version: "1.7.0"
installed:
  languages:
    - typescript
    - python
    - go
  frameworks:
    - nextjs
    - fastapi
  workflows:
    - all
```

This file tracks:
- Installed framework version
- Which components are installed
- Allows selective updates

## Available Components

### Languages (21)

TypeScript, Python, Go, Rust, Kotlin, Java, C#, PHP, Swift, C/C++, Ruby, SQL, Shell/Bash, R, Dart, HTML/CSS, Lua, Assembly, CUDA, Solidity, Zig

### Frameworks (33)

- **TypeScript/JavaScript**: React, Next.js, Express
- **Python**: Django, FastAPI, Flask
- **Go**: Gin, Echo, Fiber
- **Rust**: Axum, Actix-web, Rocket
- **Kotlin**: Spring Boot, Ktor, Android Compose
- **Java**: Spring Boot, Quarkus, Micronaut
- **C#**: ASP.NET Core, Blazor, Unity
- **PHP**: Laravel, Symfony, WordPress
- **Swift**: SwiftUI, UIKit, Vapor
- **Ruby**: Rails, Sinatra, Hanami
- **Dart**: Flutter, Shelf, Dart Frog

### Workflows (13)

- `initialize-project` - Project setup
- `create-prd` - Requirements documents
- `generate-tasks` - Task breakdown
- `code-review` - Pre-commit quality review
- `security-audit` - Security assessment
- `testing-strategy` - Test planning
- `cleanup-project` - Prune unused guides
- `refactoring` - Technical debt remediation
- `dependency-update` - Safe dependency updates
- `update-framework` - AICoF version updates
- `troubleshooting` - Debugging workflow
- `generate-agents-md` - Cross-tool compatibility
- `document-work` - Capture patterns

## Development

### Building from Source

```bash
cd packages/cli

# Install dependencies
make deps

# Build
make build

# Run tests
make test

# Build for all platforms
make build-all
```

### Running Locally

```bash
# Run directly
make run ARGS="version"

# Or build and run
make build
./bin/aicof version
```

## License

MIT License - see [LICENSE](../../LICENSE)
