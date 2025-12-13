---
title: Language Guides
description: Language-specific guardrails that auto-load based on file extensions
---

# Language Guides

Language-specific rules and best practices that auto-load based on the files you're working with.

---

## Overview

When you work on a file, AI automatically loads the appropriate language guide:

| Language | Extensions | Guide |
|----------|------------|-------|
| **TypeScript** | `.ts`, `.tsx`, `.js`, `.jsx` | [typescript.md](typescript.md) |
| **Python** | `.py` | [python.md](python.md) |
| **Go** | `.go` | [go.md](go.md) |
| **Rust** | `.rs` | [rust.md](rust.md) |
| **Kotlin** | `.kt`, `.kts` | [kotlin.md](kotlin.md) |

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
7. **Framework Patterns** - React, Django, Gin, etc.
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

### Validation Libraries

| Language | Recommended |
|----------|-------------|
| TypeScript | Zod |
| Python | Pydantic |
| Go | go-playground/validator |
| Rust | serde + custom validators |
| Kotlin | kotlinx.serialization + init blocks |

### Testing Frameworks

| Language | Recommended |
|----------|-------------|
| TypeScript | Vitest or Jest |
| Python | pytest |
| Go | built-in testing |
| Rust | built-in testing |
| Kotlin | Kotest + MockK |

### Linting/Formatting

| Language | Format | Lint |
|----------|--------|------|
| TypeScript | Prettier | ESLint |
| Python | Black | Ruff |
| Go | gofmt | golangci-lint |
| Rust | rustfmt | Clippy |
| Kotlin | ktlint | detekt |

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

---

## Adding a New Language

Want to add Java, C#, Swift, or another language?

1. Copy an existing guide as a template
2. Adapt the sections for your language
3. Add to `.agent/language-guides/`
4. Update this index

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

## Language Guides

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
