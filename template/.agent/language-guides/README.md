# Language-Specific Guides

This directory contains language and framework-specific guardrails and best practices.

## When to Use

AI automatically loads the relevant language guide when:
1. Tech stack is defined in `.agent/project.md`
2. File extensions detected (`.ts`, `.tsx`, `.js`, `.jsx`, `.py`, `.go`, `.rs`, `.kt`, `.kts`, `.java`, `.cs`, `.php`, `.swift`, `.c`, `.cpp`, `.h`, `.hpp`, `.rb`, `.sql`, `.sh`, `.bash`, `.r`, `.R`, `.dart`, `.html`, `.css`, `.scss`, `.lua`, `.asm`, `.s`, `.cu`, `.cuh`, `.sol`, `.zig`)
3. User explicitly requests: `@.agent/language-guides/typescript.md`

## Available Guides

### Core Languages
- **[typescript.md](typescript.md)** - TypeScript & JavaScript (React, Node.js, Next.js)
- **[python.md](python.md)** - Python (Django, FastAPI, Flask)
- **[go.md](go.md)** - Go microservices and APIs
- **[rust.md](rust.md)** - Rust systems programming
- **[kotlin.md](kotlin.md)** - Kotlin (Android, Spring Boot, Ktor)

### Enterprise & Web Languages
- **[java.md](java.md)** - Java (Spring Boot, Jakarta EE, enterprise)
- **[csharp.md](csharp.md)** - C# (.NET 8, ASP.NET Core, xUnit)
- **[php.md](php.md)** - PHP (Laravel, Symfony, PHPUnit)
- **[ruby.md](ruby.md)** - Ruby (Rails 7, RSpec, Sidekiq)

### Systems & Native Languages
- **[swift.md](swift.md)** - Swift (iOS, macOS, SwiftUI)
- **[cpp.md](cpp.md)** - C/C++ (C++20, CMake, GoogleTest)

### Data & Scripting Languages
- **[sql.md](sql.md)** - SQL (PostgreSQL, MySQL, query optimization)
- **[shell.md](shell.md)** - Shell/Bash (scripting, automation, POSIX)
- **[r.md](r.md)** - R (statistical computing, tidyverse, Shiny)

### Mobile & Web Frontend
- **[dart.md](dart.md)** - Dart (Flutter, mobile development)
- **[html-css.md](html-css.md)** - HTML/CSS (web standards, accessibility, BEM)

### Specialized Languages
- **[lua.md](lua.md)** - Lua (scripting, Love2D, Neovim)
- **[assembly.md](assembly.md)** - Assembly (x86-64, ARM64, RISC-V)
- **[cuda.md](cuda.md)** - CUDA (GPU computing, parallel processing)
- **[solidity.md](solidity.md)** - Solidity (Ethereum, smart contracts, DeFi)
- **[zig.md](zig.md)** - Zig (systems programming, C interop)

## How AI Should Use These

### Priority Order
1. **CLAUDE.md guardrails** (universal, always apply)
2. **Language-specific guardrails** (from this directory)
3. **Project-specific patterns** (from `.agent/patterns.md`)

### Loading Strategy
```
IF language detected:
  Load @.agent/language-guides/{language}.md
  Apply language-specific guardrails
  Override universal guardrails where conflicts exist
```

## Adding New Languages

To add a new language guide:

1. Create `{language}.md` in this directory
2. Follow the template structure (see typescript.md as reference)
3. Include:
   - Language-specific guardrails
   - Framework recommendations
   - Testing approaches
   - Common pitfalls
   - Tooling (linters, formatters)
4. Update this README with link

## Template Structure

```markdown
# {Language} Guide

## Core Principles
[Language philosophy, idioms]

## Guardrails
[Specific to this language]

## Frameworks
[Popular frameworks and their patterns]

## Testing
[Testing approach for this language]

## Tooling
[Linters, formatters, build tools]

## Common Pitfalls
[Language-specific mistakes to avoid]
```
