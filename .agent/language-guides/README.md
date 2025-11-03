# Language-Specific Guides

This directory contains language and framework-specific guardrails and best practices.

## When to Use

AI automatically loads the relevant language guide when:
1. Tech stack is defined in `.agent/project.md`
2. File extensions detected (`.ts`, `.py`, `.go`, `.rs`, `.kt`, `.kts`)
3. User explicitly requests: `@.agent/language-guides/typescript.md`

## Available Guides

- **[typescript.md](typescript.md)** - TypeScript & JavaScript (React, Node.js, Next.js)
- **[python.md](python.md)** - Python (Django, FastAPI, Flask)
- **[go.md](go.md)** - Go microservices and APIs
- **[rust.md](rust.md)** - Rust systems programming
- **[kotlin.md](kotlin.md)** - Kotlin (Android, Spring Boot, Ktor)

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
