---
title: Rust Guide
description: Rust development guardrails and best practices
---

# Rust Guide

> **Applies to**: Rust 1.70+, Systems Programming, WebAssembly, CLIs

---

## Core Principles

1. **Memory Safety**: No null, no dangling pointers, no data races
2. **Ownership**: Clear ownership, borrowing, and lifetimes
3. **Zero-Cost Abstractions**: High-level code with low-level performance
4. **Fearless Concurrency**: Compiler prevents data races
5. **Explicit Over Implicit**: Types, errors, and behavior are explicit

---

## Language-Specific Guardrails

### Rust Version & Setup

```
✓ Use Rust 1.70+ (stable channel)
✓ Use Cargo.toml for dependency management
✓ Pin dependency versions or use ~ for compatible updates
✓ Run cargo update periodically for security patches
```

### Code Style (Rustfmt)

```
✓ Run cargo fmt before every commit
✓ Run cargo clippy for linting
✓ Follow Rust API Guidelines
✓ Use snake_case for functions, variables, modules
✓ Use PascalCase for types, traits, enums
✓ Use SCREAMING_SNAKE_CASE for constants
✓ Prefer exhaustive match over if let chains
```

### Error Handling

```
✓ Use Result<T, E> for recoverable errors
✓ Use custom error types (not just String)
✓ Implement std::error::Error trait for custom errors
✓ Use ? operator for error propagation
✓ Only use panic! or unwrap() when truly impossible to fail
✓ Use expect() with message over unwrap() when confident
```

### Ownership & Borrowing

```
✓ Prefer borrowing (&T) over owned (T) in function parameters
✓ Use &mut T only when mutation needed
✓ Avoid clone() unless necessary (understand cost)
✓ Use Cow<'a, T> for clone-on-write scenarios
✓ Lifetime annotations only when compiler can't infer
```

### Safety

```
✓ Minimize unsafe blocks (document invariants)
✓ Use #[forbid(unsafe_code)] in crates that shouldn't have unsafe
✓ Audit all unsafe in code reviews
✓ Prefer safe abstractions (Vec<T> over raw pointers)
```

---

## Error Handling Patterns

### Using `thiserror` (Libraries)

```rust
use thiserror::Error;

#[derive(Error, Debug)]
pub enum DataError {
    #[error("User {id} not found")]
    NotFound { id: String },

    #[error("Invalid email format: {0}")]
    InvalidEmail(String),

    #[error(transparent)]
    Database(#[from] sqlx::Error),
}
```

### Using `anyhow` (Applications)

```rust
use anyhow::{Context, Result};

fn read_config() -> Result<Config> {
    let contents = std::fs::read_to_string("config.toml")
        .context("Failed to read config file")?;

    toml::from_str(&contents)
        .context("Failed to parse config")
}
```

---

## Testing

### Guardrails

```
✓ Tests in same file: #[cfg(test)] mod tests { ... }
✓ Integration tests in tests/ directory
✓ Public API in lib.rs, private in modules
✓ Coverage target: >80% for libraries, >60% for applications
✓ Use cargo tarpaulin or cargo llvm-cov for coverage
```

### Example

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_addition() {
        assert_eq!(add(2, 2), 4);
    }

    #[test]
    #[should_panic(expected = "divide by zero")]
    fn test_divide_by_zero() {
        divide(10, 0);
    }

    #[test]
    fn test_result() -> Result<(), Box<dyn std::error::Error>> {
        let result = fallible_operation()?;
        assert_eq!(result, 42);
        Ok(())
    }
}
```

---

## Tooling

### Essential Commands

```bash
# Format
cargo fmt

# Lint
cargo clippy -- -D warnings

# Test
cargo test
cargo test --all-features

# Build
cargo build --release

# Check (faster than build)
cargo check

# Documentation
cargo doc --open

# Audit dependencies
cargo audit

# Benchmarks
cargo bench
```

### Configuration

```toml
# Cargo.toml
[package]
name = "myproject"
version = "0.1.0"
edition = "2021"
rust-version = "1.70"  # MSRV

[dependencies]
serde = { version = "1.0", features = ["derive"] }
tokio = { version = "1.0", features = ["full"] }

[profile.release]
lto = true           # Link-time optimization
codegen-units = 1    # Better optimization
opt-level = 3        # Maximum optimization
strip = true         # Strip symbols
```

---

## Common Pitfalls

### Don't Do This

```rust
// ❌ Unnecessary clone
fn process(data: String) {
    println!("{}", data.clone());
}

// ❌ Unwrap without justification
let value = risky_operation().unwrap();

// ❌ Ignoring Result
let _ = File::create("file.txt");

// ❌ String for errors
fn do_something() -> Result<(), String> {
    Err("something went wrong".to_string())
}
```

### Do This Instead

```rust
// ✅ Borrow instead of clone
fn process(data: &str) {
    println!("{}", data);
}

// ✅ Proper error handling
let value = risky_operation()
    .expect("Operation is guaranteed to succeed because...");

// ✅ Handle Result
File::create("file.txt").context("Failed to create file")?;

// ✅ Custom error type
fn do_something() -> Result<(), AppError> {
    Err(AppError::OperationFailed)
}
```

---

## Async Patterns (Tokio)

### Basic Async Function

```rust
use tokio;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let result = fetch_data().await?;
    println!("Result: {}", result);
    Ok(())
}

async fn fetch_data() -> Result<String, reqwest::Error> {
    let response = reqwest::get("https://api.example.com/data")
        .await?
        .text()
        .await?;

    Ok(response)
}
```

### Concurrent Tasks

```rust
use tokio::task;

async fn process_concurrently() -> Result<(), AppError> {
    let task1 = task::spawn(async { fetch_users().await });
    let task2 = task::spawn(async { fetch_posts().await });

    let (users, posts) = tokio::try_join!(task1, task2)?;

    Ok(())
}
```

### Timeouts

```rust
use tokio::time::{timeout, Duration};

async fn with_timeout() -> Result<Data, AppError> {
    let result = timeout(
        Duration::from_secs(5),
        fetch_data()
    ).await
    .map_err(|_| AppError::Timeout)?;

    result
}
```

---

## Web Server (Axum)

```rust
use axum::{
    routing::{get, post},
    Json, Router,
    extract::{Path, State},
    http::StatusCode,
};
use serde::{Deserialize, Serialize};
use std::sync::Arc;

#[derive(Clone)]
struct AppState {
    db: Arc<Database>,
}

#[tokio::main]
async fn main() {
    let state = AppState {
        db: Arc::new(Database::new()),
    };

    let app = Router::new()
        .route("/health", get(health))
        .route("/users", post(create_user))
        .route("/users/:id", get(get_user))
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000")
        .await
        .unwrap();

    axum::serve(listener, app).await.unwrap();
}

async fn health() -> &'static str {
    "OK"
}

async fn get_user(
    State(state): State<AppState>,
    Path(id): Path<u64>,
) -> Result<Json<User>, StatusCode> {
    let user = state.db.get_user(id)
        .await
        .map_err(|_| StatusCode::NOT_FOUND)?;

    Ok(Json(user))
}
```

---

## Performance Considerations

### Optimization Guardrails

```
✓ Use Vec::with_capacity() when size known
✓ Use iterators over manual loops (zero-cost abstractions)
✓ Profile with cargo flamegraph before optimizing
✓ Benchmark with criterion crate
✓ Use &str over String when possible
✓ Consider Arc<T> vs Rc<T> for sharing (Arc is thread-safe)
```

### Example

```rust
// Pre-allocate capacity
let mut vec = Vec::with_capacity(1000);

// Iterator chains (zero-cost)
let sum: i32 = numbers
    .iter()
    .filter(|&&x| x > 0)
    .map(|x| x * 2)
    .sum();

// String slices over owned strings
fn process(text: &str) { // Not String
    // ...
}
```

---

## Security Best Practices

### Guardrails

```
✓ Use #[forbid(unsafe_code)] unless unsafe truly needed
✓ Validate all external inputs (serde, custom validators)
✓ Use secrecy crate for sensitive data (API keys, passwords)
✓ Hash passwords with argon2 or bcrypt
✓ Use ring or rustls for cryptography
✓ Run cargo audit regularly
✓ Use cargo deny to check licenses and security
```

### Example

```rust
use secrecy::{Secret, ExposeSecret};

struct Config {
    api_key: Secret<String>,
}

fn use_api_key(config: &Config) {
    let key = config.api_key.expose_secret();
    // Use key (not logged, not displayed by Debug)
}
```

---

## References

- [The Rust Programming Language (Book)](https://doc.rust-lang.org/book/)
- [Rust by Example](https://doc.rust-lang.org/rust-by-example/)
- [Rust API Guidelines](https://rust-lang.github.io/api-guidelines/)
- [The Async Book](https://rust-lang.github.io/async-book/)
- [Clippy Lints](https://rust-lang.github.io/rust-clippy/)
