---
title: Framework Guides
description: Framework-specific patterns and best practices for 33 frameworks across 11 language families
---

# Framework Guides

Framework-specific patterns, best practices, and conventions that supplement the language guides.

---

## Overview

When working with a specific framework, AI loads the appropriate guide to provide deep, framework-specific guidance:

- **Project Structure** - Recommended directory layout and conventions
- **Core Patterns** - Framework idioms and best practices
- **Testing** - Framework-specific testing approaches
- **Common Pitfalls** - What to avoid
- **Security** - Framework-specific security considerations
- **Performance** - Optimization techniques

---

## Quick Reference

### By Language

| Language | Frameworks |
|----------|------------|
| **TypeScript/JS** | React, Next.js, Express |
| **Python** | Django, FastAPI, Flask |
| **Go** | Gin, Echo, Fiber |
| **Rust** | Axum, Actix-web, Rocket |
| **Kotlin** | Spring Boot (Kotlin), Ktor, Android Compose |
| **Java** | Spring Boot, Quarkus, Micronaut |
| **C#** | ASP.NET Core, Blazor, Unity |
| **PHP** | Laravel, Symfony, WordPress |
| **Swift** | SwiftUI, UIKit, Vapor |
| **Ruby** | Rails, Sinatra, Hanami |
| **Dart** | Flutter, Shelf, Dart Frog |

---

## TypeScript/JavaScript Frameworks

<div class="grid cards" markdown>

-   :simple-react:{ .lg .middle } **React**

    ---

    React 18+, hooks, state management, component patterns, testing with Vitest.

    [:octicons-arrow-right-24: React Guide](react.md)

-   :simple-nextdotjs:{ .lg .middle } **Next.js**

    ---

    Next.js 14+, App Router, React Server Components, API routes, middleware.

    [:octicons-arrow-right-24: Next.js Guide](nextjs.md)

-   :simple-express:{ .lg .middle } **Express**

    ---

    Express.js, middleware patterns, REST APIs, error handling, security.

    [:octicons-arrow-right-24: Express Guide](express.md)

</div>

---

## Python Frameworks

<div class="grid cards" markdown>

-   :simple-django:{ .lg .middle } **Django**

    ---

    Django 5+, ORM, admin, Django REST Framework, templates, security.

    [:octicons-arrow-right-24: Django Guide](django.md)

-   :simple-fastapi:{ .lg .middle } **FastAPI**

    ---

    FastAPI, async, Pydantic validation, OpenAPI, dependency injection.

    [:octicons-arrow-right-24: FastAPI Guide](fastapi.md)

-   :simple-flask:{ .lg .middle } **Flask**

    ---

    Flask, blueprints, extensions, Jinja2 templates, SQLAlchemy.

    [:octicons-arrow-right-24: Flask Guide](flask.md)

</div>

---

## Go Frameworks

<div class="grid cards" markdown>

-   :simple-go:{ .lg .middle } **Gin**

    ---

    Gin, middleware, routing, validation, JSON handling, high performance.

    [:octicons-arrow-right-24: Gin Guide](gin.md)

-   :simple-go:{ .lg .middle } **Echo**

    ---

    Echo, middleware, routing, context, data binding, WebSocket.

    [:octicons-arrow-right-24: Echo Guide](echo.md)

-   :simple-go:{ .lg .middle } **Fiber**

    ---

    Fiber, Express-inspired, high-performance, zero memory allocation.

    [:octicons-arrow-right-24: Fiber Guide](fiber.md)

</div>

---

## Rust Frameworks

<div class="grid cards" markdown>

-   :simple-rust:{ .lg .middle } **Axum**

    ---

    Axum, Tower ecosystem, async handlers, extractors, type-safe routing.

    [:octicons-arrow-right-24: Axum Guide](axum.md)

-   :simple-rust:{ .lg .middle } **Actix-web**

    ---

    Actix-web, actor model, high-performance, middleware, WebSocket.

    [:octicons-arrow-right-24: Actix-web Guide](actix-web.md)

-   :simple-rust:{ .lg .middle } **Rocket**

    ---

    Rocket, type-safe, macro-driven, fairings, managed state.

    [:octicons-arrow-right-24: Rocket Guide](rocket.md)

</div>

---

## Kotlin Frameworks

<div class="grid cards" markdown>

-   :simple-spring:{ .lg .middle } **Spring Boot (Kotlin)**

    ---

    Spring Boot with Kotlin, coroutines, WebFlux, null safety.

    [:octicons-arrow-right-24: Spring Boot Kotlin Guide](spring-boot-kotlin.md)

-   :simple-ktor:{ .lg .middle } **Ktor**

    ---

    Ktor, coroutines, DSL configuration, plugins, lightweight.

    [:octicons-arrow-right-24: Ktor Guide](ktor.md)

-   :simple-android:{ .lg .middle } **Android Compose**

    ---

    Jetpack Compose, Material 3, state management, navigation.

    [:octicons-arrow-right-24: Android Compose Guide](android-compose.md)

</div>

---

## Java Frameworks

<div class="grid cards" markdown>

-   :simple-spring:{ .lg .middle } **Spring Boot**

    ---

    Spring Boot, JPA, Spring Security, REST, WebFlux, reactive.

    [:octicons-arrow-right-24: Spring Boot Guide](spring-boot-java.md)

-   :simple-quarkus:{ .lg .middle } **Quarkus**

    ---

    Quarkus, GraalVM native compilation, reactive, CDI, fast startup.

    [:octicons-arrow-right-24: Quarkus Guide](quarkus.md)

-   :material-server:{ .lg .middle } **Micronaut**

    ---

    Micronaut, compile-time DI, GraalVM, cloud-native, low memory.

    [:octicons-arrow-right-24: Micronaut Guide](micronaut.md)

</div>

---

## C# / .NET Frameworks

<div class="grid cards" markdown>

-   :simple-dotnet:{ .lg .middle } **ASP.NET Core**

    ---

    ASP.NET Core 8+, Minimal APIs, Entity Framework Core, Blazor, SignalR.

    [:octicons-arrow-right-24: ASP.NET Core Guide](aspnet-core.md)

-   :simple-blazor:{ .lg .middle } **Blazor**

    ---

    Blazor WebAssembly/Server, components, SignalR, interop.

    [:octicons-arrow-right-24: Blazor Guide](blazor.md)

-   :simple-unity:{ .lg .middle } **Unity**

    ---

    Unity game development, C# scripting, MonoBehaviour, physics.

    [:octicons-arrow-right-24: Unity Guide](unity.md)

</div>

---

## PHP Frameworks

<div class="grid cards" markdown>

-   :simple-laravel:{ .lg .middle } **Laravel**

    ---

    Laravel 11+, Eloquent ORM, Blade, queues, Livewire, Inertia.

    [:octicons-arrow-right-24: Laravel Guide](laravel.md)

-   :simple-symfony:{ .lg .middle } **Symfony**

    ---

    Symfony 7+, Doctrine ORM, Twig, bundles, Flex.

    [:octicons-arrow-right-24: Symfony Guide](symfony.md)

-   :simple-wordpress:{ .lg .middle } **WordPress**

    ---

    WordPress themes, plugins, REST API, blocks, WooCommerce.

    [:octicons-arrow-right-24: WordPress Guide](wordpress.md)

</div>

---

## Swift Frameworks

<div class="grid cards" markdown>

-   :simple-swift:{ .lg .middle } **SwiftUI**

    ---

    SwiftUI, declarative UI, Combine, state management, animations.

    [:octicons-arrow-right-24: SwiftUI Guide](swiftui.md)

-   :simple-swift:{ .lg .middle } **UIKit**

    ---

    UIKit, programmatic/storyboard, Auto Layout, view controllers.

    [:octicons-arrow-right-24: UIKit Guide](uikit.md)

-   :simple-swift:{ .lg .middle } **Vapor**

    ---

    Vapor, Fluent ORM, async Swift, middleware, WebSocket.

    [:octicons-arrow-right-24: Vapor Guide](vapor.md)

</div>

---

## Ruby Frameworks

<div class="grid cards" markdown>

-   :simple-rubyonrails:{ .lg .middle } **Rails**

    ---

    Rails 7+, ActiveRecord, Hotwire, Turbo, Stimulus, Action Cable.

    [:octicons-arrow-right-24: Rails Guide](rails.md)

-   :simple-ruby:{ .lg .middle } **Sinatra**

    ---

    Sinatra, lightweight Ruby web apps, DSL routing, middleware.

    [:octicons-arrow-right-24: Sinatra Guide](sinatra.md)

-   :simple-ruby:{ .lg .middle } **Hanami**

    ---

    Hanami 2+, clean architecture, dry-rb ecosystem, modular.

    [:octicons-arrow-right-24: Hanami Guide](hanami.md)

</div>

---

## Dart Frameworks

<div class="grid cards" markdown>

-   :simple-flutter:{ .lg .middle } **Flutter**

    ---

    Flutter, Riverpod, go_router, Material 3, cross-platform.

    [:octicons-arrow-right-24: Flutter Guide](flutter.md)

-   :simple-dart:{ .lg .middle } **Shelf**

    ---

    Shelf, middleware-based HTTP server, composable handlers.

    [:octicons-arrow-right-24: Shelf Guide](shelf.md)

-   :simple-dart:{ .lg .middle } **Dart Frog**

    ---

    Dart Frog, file-based routing, server-side Dart, middleware.

    [:octicons-arrow-right-24: Dart Frog Guide](dart-frog.md)

</div>

---

## Relationship to Language Guides

Framework guides **supplement** language guides:

| Type | Purpose | Example |
|------|---------|---------|
| **Language Guide** | General language rules | TypeScript strict mode, Python type hints |
| **Framework Guide** | Framework-specific patterns | React hooks, Django ORM queries |

Both are loaded together when working with a framework.

---

## Contributing

Want to add a framework skill?

1. Copy an existing framework skill directory as a template
2. Adapt `SKILL.md` and `references/` for your framework
3. Add to `.agent/skills/<framework-name>/`
4. Add to `docs/frameworks/`
5. Update navigation in `mkdocs.yml`
6. Update this index
7. Submit a pull request
