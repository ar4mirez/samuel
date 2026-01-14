# Framework Guides

Framework-specific best practices and patterns that supplement the language guides.

---

## Overview

When working with a specific framework, AI loads the appropriate guide to provide:

- **Project Structure** - Recommended directory layout
- **Core Patterns** - Framework idioms and conventions
- **Testing** - Framework-specific testing approaches
- **Common Pitfalls** - What to avoid
- **Security** - Framework-specific security considerations

---

## Available Framework Guides

### TypeScript/JavaScript

| Framework | Description | Guide |
|-----------|-------------|-------|
| **React** | React 18+, hooks, state management, component patterns | [react.md](react.md) |
| **Next.js** | Next.js 14+, App Router, RSC, API routes | [nextjs.md](nextjs.md) |
| **Express** | Express.js, middleware, REST APIs, error handling | [express.md](express.md) |

### Python

| Framework | Description | Guide |
|-----------|-------------|-------|
| **Django** | Django 5+, ORM, admin, DRF, templates | [django.md](django.md) |
| **FastAPI** | FastAPI, async, Pydantic, OpenAPI, dependency injection | [fastapi.md](fastapi.md) |
| **Flask** | Flask, blueprints, extensions, Jinja2 | [flask.md](flask.md) |

### Go

| Framework | Description | Guide |
|-----------|-------------|-------|
| **Gin** | Gin, middleware, routing, validation | [gin.md](gin.md) |
| **Echo** | Echo, middleware, routing, context | [echo.md](echo.md) |
| **Fiber** | Fiber, Express-style, high-performance | [fiber.md](fiber.md) |

### Rust

| Framework | Description | Guide |
|-----------|-------------|-------|
| **Axum** | Axum, Tower, async handlers, extractors | [axum.md](axum.md) |
| **Actix-web** | Actix-web, actors, high-performance | [actix-web.md](actix-web.md) |
| **Rocket** | Rocket, type-safe, macros, fairings | [rocket.md](rocket.md) |

### Kotlin

| Framework | Description | Guide |
|-----------|-------------|-------|
| **Spring Boot (Kotlin)** | Spring Boot with Kotlin, coroutines, WebFlux | [spring-boot-kotlin.md](spring-boot-kotlin.md) |
| **Ktor** | Ktor, coroutines, DSL, plugins | [ktor.md](ktor.md) |
| **Android Compose** | Jetpack Compose, Material 3, state management | [android-compose.md](android-compose.md) |

### Java

| Framework | Description | Guide |
|-----------|-------------|-------|
| **Spring Boot** | Spring Boot, JPA, Security, REST, WebFlux | [spring-boot-java.md](spring-boot-java.md) |
| **Quarkus** | Quarkus, GraalVM native, reactive, CDI | [quarkus.md](quarkus.md) |
| **Micronaut** | Micronaut, compile-time DI, GraalVM | [micronaut.md](micronaut.md) |

### C# / .NET

| Framework | Description | Guide |
|-----------|-------------|-------|
| **ASP.NET Core** | ASP.NET Core, Minimal APIs, EF Core, Blazor | [aspnet-core.md](aspnet-core.md) |
| **Blazor** | Blazor WebAssembly/Server, SignalR, components | [blazor.md](blazor.md) |
| **Unity** | Unity, C# scripting, game development patterns | [unity.md](unity.md) |

### PHP

| Framework | Description | Guide |
|-----------|-------------|-------|
| **Laravel** | Laravel 11+, Eloquent, Blade, queues, Livewire | [laravel.md](laravel.md) |
| **Symfony** | Symfony 7+, Doctrine, Twig, bundles | [symfony.md](symfony.md) |
| **WordPress** | WordPress themes, plugins, REST API, blocks | [wordpress.md](wordpress.md) |

### Swift

| Framework | Description | Guide |
|-----------|-------------|-------|
| **SwiftUI** | SwiftUI, declarative UI, Combine, state | [swiftui.md](swiftui.md) |
| **UIKit** | UIKit, programmatic/storyboard, Auto Layout | [uikit.md](uikit.md) |
| **Vapor** | Vapor, Fluent ORM, async Swift server | [vapor.md](vapor.md) |

### Ruby

| Framework | Description | Guide |
|-----------|-------------|-------|
| **Rails** | Rails 7+, ActiveRecord, Hotwire, Turbo | [rails.md](rails.md) |
| **Sinatra** | Sinatra, lightweight Ruby web apps | [sinatra.md](sinatra.md) |
| **Hanami** | Hanami 2+, clean architecture, dry-rb | [hanami.md](hanami.md) |

### Dart

| Framework | Description | Guide |
|-----------|-------------|-------|
| **Flutter** | Flutter, Riverpod, go_router, Material 3 | [flutter.md](flutter.md) |
| **Shelf** | Shelf, middleware HTTP server | [shelf.md](shelf.md) |
| **Dart Frog** | Dart Frog, file-based routing, server-side | [dart-frog.md](dart-frog.md) |

---

## Loading Protocol

Framework guides are loaded **on-demand** when:

1. You reference the framework by name (e.g., "using React", "in Django")
2. You're working in a project that uses the framework (detected via config files)
3. You explicitly request it

**Example triggers:**
- `package.json` with `react` dependency → React guide available
- `requirements.txt` with `django` → Django guide available
- `Cargo.toml` with `axum` → Axum guide available

---

## Guide Structure

Each framework guide follows this structure:

```markdown
# [Framework] Guide

> **Applies to**: [Framework version]+

## Quick Start
## Project Structure
## Core Patterns
## State Management (if applicable)
## Routing
## Data Fetching / Database
## Testing
## Common Pitfalls
## Security Best Practices
## Performance
## References
```

---

## Relationship to Language Guides

Framework guides **supplement** language guides, not replace them:

- **Language Guide**: General language rules (TypeScript strict mode, Python type hints)
- **Framework Guide**: Framework-specific patterns (React hooks, Django ORM)

Both are loaded together when working with a framework.

---

## Contributing

To add a new framework guide:

1. Copy an existing guide as a template
2. Adapt sections for your framework
3. Add to this README
4. Update CLAUDE.md "Load Framework Guide" section
5. Submit a pull request

---

## Statistics

- **Total Frameworks**: 33
- **Language Families**: 11
- **Average Guide Size**: ~30KB
