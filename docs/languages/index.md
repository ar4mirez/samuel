---
title: Language Guides
description: Language-specific guardrails that auto-load based on file extensions
---

# Language Guides

Language-specific rules and best practices that auto-load based on the files you're working with.

---

## Overview

When you work on a file, AI automatically loads the appropriate language guide:

### Core Languages

| Language | Extensions | Guide |
|----------|------------|-------|
| **TypeScript** | `.ts`, `.tsx`, `.js`, `.jsx` | [typescript.md](typescript.md) |
| **Python** | `.py` | [python.md](python.md) |
| **Go** | `.go` | [go.md](go.md) |
| **Rust** | `.rs` | [rust.md](rust.md) |
| **Kotlin** | `.kt`, `.kts` | [kotlin.md](kotlin.md) |

### Enterprise Languages

| Language | Extensions | Guide |
|----------|------------|-------|
| **Java** | `.java` | [java.md](java.md) |
| **C#** | `.cs` | [csharp.md](csharp.md) |
| **PHP** | `.php` | [php.md](php.md) |
| **Swift** | `.swift` | [swift.md](swift.md) |
| **C/C++** | `.c`, `.cpp`, `.h`, `.hpp` | [cpp.md](cpp.md) |
| **Ruby** | `.rb` | [ruby.md](ruby.md) |

### Specialized Domains

| Language | Extensions | Guide |
|----------|------------|-------|
| **SQL** | `.sql` | [sql.md](sql.md) |
| **Shell/Bash** | `.sh`, `.bash` | [shell.md](shell.md) |
| **R** | `.r`, `.R` | [r.md](r.md) |
| **Dart** | `.dart` | [dart.md](dart.md) |
| **HTML/CSS** | `.html`, `.css`, `.scss` | [html-css.md](html-css.md) |
| **Lua** | `.lua` | [lua.md](lua.md) |
| **Assembly** | `.asm`, `.s` | [assembly.md](assembly.md) |
| **CUDA** | `.cu`, `.cuh` | [cuda.md](cuda.md) |
| **Solidity** | `.sol` | [solidity.md](solidity.md) |
| **Zig** | `.zig` | [zig.md](zig.md) |

**No manual selection needed** - it just works!

---

## What's in a Language Guide?

Each guide contains:

1. **Core Principles** - Language philosophy
2. **Language-Specific Guardrails** - Rules beyond the universal ones
3. **Validation & Input Handling** - Recommended libraries
4. **Testing** - Frameworks and patterns
5. **Tooling** - Linters, formatters, configs
6. **Common Pitfalls** - What to avoid
7. **Framework Patterns** - React, Django, Rails, Spring, etc.
8. **Performance Considerations** - Language-specific optimizations
9. **Security Best Practices** - Language-specific security

---

## Quick Comparison

### Type Safety

| Language | Approach |
|----------|----------|
| TypeScript | Strict mode, no `any`, explicit types |
| Python | Type hints + mypy strict mode |
| Go | Built-in static types |
| Rust | Built-in static types, ownership |
| Kotlin | Null safety, smart casts |
| Java | Strong static typing, Optional |
| C# | Nullable reference types, strong typing |
| Swift | Optionals, strong static typing |
| C/C++ | Static types, no null safety (use smart pointers) |

### Validation Libraries

| Language | Recommended |
|----------|-------------|
| TypeScript | Zod |
| Python | Pydantic |
| Go | go-playground/validator |
| Rust | serde + custom validators |
| Kotlin | kotlinx.serialization + init blocks |
| Java | Bean Validation (JSR-380) |
| C# | FluentValidation |
| PHP | Laravel Validation, Symfony Validator |
| Ruby | ActiveModel::Validations |

### Testing Frameworks

| Language | Recommended |
|----------|-------------|
| TypeScript | Vitest or Jest |
| Python | pytest |
| Go | built-in testing |
| Rust | built-in testing |
| Kotlin | Kotest + MockK |
| Java | JUnit 5 + Mockito |
| C# | xUnit + Moq |
| PHP | PHPUnit or Pest |
| Ruby | RSpec |
| Swift | XCTest |
| C/C++ | GoogleTest |

### Linting/Formatting

| Language | Format | Lint |
|----------|--------|------|
| TypeScript | Prettier | ESLint |
| Python | Black | Ruff |
| Go | gofmt | golangci-lint |
| Rust | rustfmt | Clippy |
| Kotlin | ktlint | detekt |
| Java | google-java-format | Checkstyle, SpotBugs |
| C# | dotnet format | Roslyn Analyzers |
| PHP | PHP-CS-Fixer | PHPStan |
| Ruby | RuboCop | RuboCop |
| Swift | SwiftFormat | SwiftLint |
| C/C++ | clang-format | clang-tidy |

---

## Universal vs Language-Specific

### Universal (from CLAUDE.md)

These apply to ALL languages:

```
✓ Functions ≤50 lines
✓ Files ≤300 lines
✓ All inputs validated
✓ Tests >80% for business logic
✓ Conventional commits
```

### Language-Specific (from guides)

These are language-specific additions:

=== "TypeScript"

    ```
    ✓ strict: true in tsconfig
    ✓ No any without justification
    ✓ Prefer const over let
    ✓ async/await over raw promises
    ```

=== "Python"

    ```
    ✓ Type hints on all functions
    ✓ Use Pydantic for validation
    ✓ Follow PEP 8 (Black)
    ✓ Use async for I/O-bound ops
    ```

=== "Go"

    ```
    ✓ Always check errors
    ✓ Use context for cancellation
    ✓ Accept interfaces, return structs
    ✓ Table-driven tests
    ```

=== "Rust"

    ```
    ✓ Handle all Result<T,E>
    ✓ Prefer borrowing over cloning
    ✓ Minimize unsafe blocks
    ✓ Use thiserror for libraries
    ```

=== "Kotlin"

    ```
    ✓ Prefer val over var
    ✓ Use data classes for DTOs
    ✓ Structured concurrency
    ✓ Safe call operator ?.
    ```

=== "Java"

    ```
    ✓ Use Optional for nullable returns
    ✓ Records for immutable data
    ✓ Stream API over loops
    ✓ Try-with-resources
    ```

=== "C#"

    ```
    ✓ Enable nullable reference types
    ✓ Use records for DTOs
    ✓ async/await for I/O
    ✓ LINQ for collections
    ```

=== "PHP"

    ```
    ✓ declare(strict_types=1)
    ✓ Type hints everywhere
    ✓ Use enums (PHP 8.1+)
    ✓ Prepared statements only
    ```

=== "Ruby"

    ```
    ✓ Use frozen string literals
    ✓ Prefer symbols over strings
    ✓ Use guard clauses
    ✓ Follow Ruby Style Guide
    ```

=== "Swift"

    ```
    ✓ Use guard for early exits
    ✓ Prefer structs over classes
    ✓ async/await for concurrency
    ✓ Protocol-oriented design
    ```

---

## All Language Guides

### Core Languages

<div class="grid cards" markdown>

-   :simple-typescript:{ .lg .middle } **TypeScript**

    ---

    TypeScript, JavaScript, React, Node.js, Next.js, Vue, Angular.

    [:octicons-arrow-right-24: TypeScript Guide](typescript.md)

-   :simple-python:{ .lg .middle } **Python**

    ---

    Python 3.9+, Django, FastAPI, Flask, Data Science.

    [:octicons-arrow-right-24: Python Guide](python.md)

-   :simple-go:{ .lg .middle } **Go**

    ---

    Go 1.20+, Microservices, APIs, CLIs.

    [:octicons-arrow-right-24: Go Guide](go.md)

-   :simple-rust:{ .lg .middle } **Rust**

    ---

    Rust 1.70+, Systems Programming, WebAssembly, CLIs.

    [:octicons-arrow-right-24: Rust Guide](rust.md)

-   :simple-kotlin:{ .lg .middle } **Kotlin**

    ---

    Kotlin 1.9+, Android, Spring Boot, Ktor, Multiplatform.

    [:octicons-arrow-right-24: Kotlin Guide](kotlin.md)

</div>

### Enterprise Languages

<div class="grid cards" markdown>

-   :fontawesome-brands-java:{ .lg .middle } **Java**

    ---

    Java 17+, Spring Boot, Jakarta EE, Microservices.

    [:octicons-arrow-right-24: Java Guide](java.md)

-   :simple-csharp:{ .lg .middle } **C#**

    ---

    C# 11+, .NET 7+, ASP.NET Core, Unity, MAUI.

    [:octicons-arrow-right-24: C# Guide](csharp.md)

-   :simple-php:{ .lg .middle } **PHP**

    ---

    PHP 8.1+, Laravel, Symfony, WordPress.

    [:octicons-arrow-right-24: PHP Guide](php.md)

-   :simple-swift:{ .lg .middle } **Swift**

    ---

    Swift 5.9+, iOS, macOS, SwiftUI, Server-Side Swift.

    [:octicons-arrow-right-24: Swift Guide](swift.md)

-   :simple-cplusplus:{ .lg .middle } **C/C++**

    ---

    C17/C23, C++17/20/23, Systems Programming, Embedded.

    [:octicons-arrow-right-24: C/C++ Guide](cpp.md)

-   :simple-ruby:{ .lg .middle } **Ruby**

    ---

    Ruby 3.0+, Rails 7+, Sinatra, RSpec.

    [:octicons-arrow-right-24: Ruby Guide](ruby.md)

</div>

### Specialized Domains

<div class="grid cards" markdown>

-   :simple-postgresql:{ .lg .middle } **SQL**

    ---

    PostgreSQL, MySQL, SQLite, SQL Server, Query Optimization.

    [:octicons-arrow-right-24: SQL Guide](sql.md)

-   :simple-gnubash:{ .lg .middle } **Shell/Bash**

    ---

    Bash 4+, POSIX sh, Zsh, CI/CD Pipelines, Automation.

    [:octicons-arrow-right-24: Shell Guide](shell.md)

-   :simple-r:{ .lg .middle } **R**

    ---

    R 4.0+, Tidyverse, Shiny, Statistical Computing.

    [:octicons-arrow-right-24: R Guide](r.md)

-   :simple-dart:{ .lg .middle } **Dart**

    ---

    Dart 3.0+, Flutter, Mobile Development.

    [:octicons-arrow-right-24: Dart Guide](dart.md)

-   :simple-html5:{ .lg .middle } **HTML/CSS**

    ---

    HTML5, CSS3, Sass/SCSS, Accessibility, Responsive Design.

    [:octicons-arrow-right-24: HTML/CSS Guide](html-css.md)

-   :simple-lua:{ .lg .middle } **Lua**

    ---

    Lua 5.4+, Love2D, Neovim, Game Development.

    [:octicons-arrow-right-24: Lua Guide](lua.md)

-   :material-chip:{ .lg .middle } **Assembly**

    ---

    x86-64, ARM64, RISC-V, OS Development.

    [:octicons-arrow-right-24: Assembly Guide](assembly.md)

-   :simple-nvidia:{ .lg .middle } **CUDA**

    ---

    CUDA 11+, GPU Computing, Deep Learning, Scientific Computing.

    [:octicons-arrow-right-24: CUDA Guide](cuda.md)

-   :simple-solidity:{ .lg .middle } **Solidity**

    ---

    Solidity 0.8+, Ethereum, Smart Contracts, DeFi.

    [:octicons-arrow-right-24: Solidity Guide](solidity.md)

-   :simple-zig:{ .lg .middle } **Zig**

    ---

    Zig 0.11+, Systems Programming, C Interop, Game Development.

    [:octicons-arrow-right-24: Zig Guide](zig.md)

</div>

---

## Contributing a Language Guide

Want to contribute a guide for a language not yet covered?

1. Copy an existing guide as a template
2. Adapt the sections for your language
3. Add to `.agent/language-guides/` AND `docs/languages/`
4. Update `mkdocs.yml` navigation
5. Update this index
6. Submit a pull request

**Template structure**:

```markdown
# [Language] Guide

> **Applies to**: [Language version]+, [frameworks]

## Core Principles
## Language-Specific Guardrails
## Validation & Input Handling
## Testing
## Tooling
## Common Pitfalls
## Framework-Specific Patterns
## Performance Considerations
## Security Best Practices
## References
```

---

## My Language Isn't Listed

The core guardrails in CLAUDE.md are **language-agnostic**. 90% of the rules still apply:

- Functions ≤50 lines ✓
- Files ≤300 lines ✓
- Input validation ✓
- Parameterized queries ✓
- Test coverage >80% ✓
- Conventional commits ✓

You can still use the system effectively without a language-specific guide!

Consider contributing a guide for your language to help others.
