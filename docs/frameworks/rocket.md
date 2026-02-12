# Rocket Framework Guide

> **Applies to**: Rocket 0.5+, Rust Web APIs, Type-Safe Web Applications
> **Use with**: `.claude/skills/rust-guide/SKILL.md`

---

## Overview

Rocket is a Rust web framework focused on ease of use, developer experience, and type safety. It uses Rust's type system to ensure correctness at compile time, with powerful features like request guards, fairings (middleware), and derive macros for minimal boilerplate.

### When to Use Rocket
- **Type-Safe APIs**: Maximum compile-time guarantees
- **Developer Experience**: Clean, expressive syntax
- **Rapid Prototyping**: Quick to get started with minimal boilerplate
- **Form Handling**: Excellent support for HTML forms and multipart uploads
- **Traditional Web Apps**: Both APIs and server-rendered HTML

### When NOT to Use Rocket
- **Async-Heavy Workloads**: Consider Axum for pure async performance
- **Minimal Dependencies**: Rocket has heavier compile times
- **Maximum Control**: Use Axum/Actix for fine-grained control

---

## Project Structure

```
myproject/
├── Cargo.toml
├── Rocket.toml                 # Rocket configuration
├── src/
│   ├── main.rs                 # Application entry point
│   ├── lib.rs                  # Library exports
│   ├── routes/
│   │   ├── mod.rs
│   │   ├── users.rs
│   │   └── health.rs
│   ├── models/
│   │   ├── mod.rs
│   │   └── user.rs
│   ├── guards/                 # Request guards
│   │   ├── mod.rs
│   │   └── auth.rs
│   ├── fairings/               # Middleware (fairings)
│   │   ├── mod.rs
│   │   └── logging.rs
│   ├── db/
│   │   ├── mod.rs
│   │   └── pool.rs
│   ├── services/
│   │   ├── mod.rs
│   │   └── user_service.rs
│   ├── error.rs
│   └── config.rs
├── tests/
│   └── api_tests.rs
├── migrations/
└── static/                     # Static files
```

---

## Dependencies

```toml
# Cargo.toml
[package]
name = "myproject"
version = "0.1.0"
edition = "2021"

[dependencies]
# Web framework
rocket = { version = "0.5", features = ["json", "secrets"] }

# Database
rocket_db_pools = { version = "0.1", features = ["sqlx_postgres"] }
sqlx = { version = "0.7", features = ["runtime-tokio", "postgres", "chrono", "uuid"] }

# Serialization
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"

# Validation
validator = { version = "0.16", features = ["derive"] }

# Authentication
jsonwebtoken = "9"
bcrypt = "0.15"

# Utilities
chrono = { version = "0.4", features = ["serde"] }
uuid = { version = "1", features = ["v4", "serde"] }
thiserror = "1.0"
dotenvy = "0.15"

# Async runtime (Rocket uses Tokio by default)
tokio = { version = "1", features = ["full"] }

[dev-dependencies]
rocket = { version = "0.5", features = ["json", "secrets"] }
```

---

## Application Entry Point

```rust
// src/main.rs
#[macro_use] extern crate rocket;

use rocket::{Build, Rocket};
use rocket_db_pools::Database;

mod config;
mod db;
mod error;
mod fairings;
mod guards;
mod models;
mod routes;
mod services;

use db::DbPool;
use fairings::RequestLogger;

#[launch]
fn rocket() -> Rocket<Build> {
    // Load environment variables
    dotenvy::dotenv().ok();

    rocket::build()
        // Attach database pool
        .attach(DbPool::init())
        // Attach custom fairings
        .attach(RequestLogger)
        // Mount routes
        .mount("/api", routes::api_routes())
        .mount("/health", routes::health_routes())
        // Register catchers for error handling
        .register("/", catchers![
            error::not_found,
            error::internal_error,
            error::unauthorized,
            error::bad_request,
        ])
}
```

---

## Configuration

```toml
# Rocket.toml
[default]
address = "0.0.0.0"
port = 8000
workers = 16
keep_alive = 5
log_level = "normal"
temp_dir = "/tmp"
limits = { form = "64 kB", json = "1 MiB" }

[default.databases.db_pool]
url = "postgres://user:password@localhost/mydb"
min_connections = 5
max_connections = 20
connect_timeout = 5
idle_timeout = 300

[debug]
log_level = "debug"

[release]
secret_key = "generate-a-256-bit-base64-key"
log_level = "critical"
```

```rust
// src/config.rs
use rocket::serde::Deserialize;

#[derive(Debug, Deserialize)]
#[serde(crate = "rocket::serde")]
pub struct AppConfig {
    pub jwt_secret: String,
    pub jwt_expiration_hours: i64,
}

impl Default for AppConfig {
    fn default() -> Self {
        Self {
            jwt_secret: std::env::var("JWT_SECRET")
                .unwrap_or_else(|_| "development-secret-key".to_string()),
            jwt_expiration_hours: std::env::var("JWT_EXPIRATION_HOURS")
                .unwrap_or_else(|_| "24".to_string())
                .parse()
                .unwrap_or(24),
        }
    }
}
```

---

## Error Handling

```rust
// src/error.rs
use rocket::http::Status;
use rocket::response::{self, Responder, Response};
use rocket::serde::json::Json;
use rocket::Request;
use serde::Serialize;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum AppError {
    #[error("Resource not found: {0}")]
    NotFound(String),

    #[error("Unauthorized: {0}")]
    Unauthorized(String),

    #[error("Validation error: {0}")]
    Validation(String),

    #[error("Database error: {0}")]
    Database(#[from] sqlx::Error),

    #[error("Internal server error")]
    Internal(String),
}

#[derive(Serialize)]
pub struct ErrorResponse {
    pub error: String,
    pub message: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub details: Option<Vec<String>>,
}

impl<'r> Responder<'r, 'static> for AppError {
    fn respond_to(self, request: &'r Request<'_>) -> response::Result<'static> {
        let (status, error_response) = match &self {
            AppError::NotFound(msg) => (
                Status::NotFound,
                ErrorResponse {
                    error: "NOT_FOUND".to_string(),
                    message: msg.clone(),
                    details: None,
                },
            ),
            AppError::Unauthorized(msg) => (
                Status::Unauthorized,
                ErrorResponse {
                    error: "UNAUTHORIZED".to_string(),
                    message: msg.clone(),
                    details: None,
                },
            ),
            AppError::Validation(msg) => (
                Status::BadRequest,
                ErrorResponse {
                    error: "VALIDATION_ERROR".to_string(),
                    message: msg.clone(),
                    details: None,
                },
            ),
            AppError::Database(e) => {
                eprintln!("Database error: {:?}", e);
                (
                    Status::InternalServerError,
                    ErrorResponse {
                        error: "DATABASE_ERROR".to_string(),
                        message: "A database error occurred".to_string(),
                        details: None,
                    },
                )
            }
            AppError::Internal(msg) => {
                eprintln!("Internal error: {}", msg);
                (
                    Status::InternalServerError,
                    ErrorResponse {
                        error: "INTERNAL_ERROR".to_string(),
                        message: "An internal error occurred".to_string(),
                        details: None,
                    },
                )
            }
        };

        Response::build_from(Json(error_response).respond_to(request)?)
            .status(status)
            .ok()
    }
}

// Error catchers
#[catch(404)]
pub fn not_found() -> Json<ErrorResponse> {
    Json(ErrorResponse {
        error: "NOT_FOUND".to_string(),
        message: "The requested resource was not found".to_string(),
        details: None,
    })
}

#[catch(500)]
pub fn internal_error() -> Json<ErrorResponse> {
    Json(ErrorResponse {
        error: "INTERNAL_ERROR".to_string(),
        message: "An internal server error occurred".to_string(),
        details: None,
    })
}

#[catch(401)]
pub fn unauthorized() -> Json<ErrorResponse> {
    Json(ErrorResponse {
        error: "UNAUTHORIZED".to_string(),
        message: "Authentication required".to_string(),
        details: None,
    })
}

#[catch(400)]
pub fn bad_request() -> Json<ErrorResponse> {
    Json(ErrorResponse {
        error: "BAD_REQUEST".to_string(),
        message: "Invalid request".to_string(),
        details: None,
    })
}
```

---

## Models

```rust
// src/models/user.rs
use chrono::{DateTime, Utc};
use rocket::serde::{Deserialize, Serialize};
use sqlx::FromRow;
use uuid::Uuid;
use validator::Validate;

#[derive(Debug, Clone, Serialize, FromRow)]
#[serde(crate = "rocket::serde")]
pub struct User {
    pub id: Uuid,
    pub email: String,
    #[serde(skip_serializing)]
    pub password_hash: String,
    pub name: String,
    pub role: String,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Deserialize, Validate)]
#[serde(crate = "rocket::serde")]
pub struct CreateUser {
    #[validate(email(message = "Invalid email format"))]
    pub email: String,

    #[validate(length(min = 8, message = "Password must be at least 8 characters"))]
    pub password: String,

    #[validate(length(min = 1, max = 100, message = "Name must be 1-100 characters"))]
    pub name: String,
}

#[derive(Debug, Deserialize, Validate)]
#[serde(crate = "rocket::serde")]
pub struct UpdateUser {
    #[validate(length(min = 1, max = 100, message = "Name must be 1-100 characters"))]
    pub name: Option<String>,

    #[validate(email(message = "Invalid email format"))]
    pub email: Option<String>,
}

#[derive(Debug, Serialize)]
#[serde(crate = "rocket::serde")]
pub struct UserResponse {
    pub id: Uuid,
    pub email: String,
    pub name: String,
    pub role: String,
    pub created_at: DateTime<Utc>,
}

impl From<User> for UserResponse {
    fn from(user: User) -> Self {
        Self {
            id: user.id,
            email: user.email,
            name: user.name,
            role: user.role,
            created_at: user.created_at,
        }
    }
}

#[derive(Debug, Deserialize, Validate)]
#[serde(crate = "rocket::serde")]
pub struct LoginRequest {
    #[validate(email)]
    pub email: String,
    pub password: String,
}

#[derive(Debug, Serialize)]
#[serde(crate = "rocket::serde")]
pub struct LoginResponse {
    pub token: String,
    pub user: UserResponse,
}
```

---

## Database Pool

```rust
// src/db/mod.rs
use rocket_db_pools::{sqlx, Database};

#[derive(Database)]
#[database("db_pool")]
pub struct DbPool(sqlx::PgPool);

// src/db/pool.rs
pub use super::DbPool;

// Helper to get a connection from the pool
use rocket_db_pools::Connection;

pub type DbConn = Connection<DbPool>;
```

---

## Request Guards

```rust
// src/guards/auth.rs
use crate::config::AppConfig;
use crate::error::AppError;
use jsonwebtoken::{decode, DecodingKey, Validation};
use rocket::http::Status;
use rocket::request::{FromRequest, Outcome, Request};
use rocket::serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize)]
pub struct Claims {
    pub sub: Uuid,
    pub email: String,
    pub role: String,
    pub exp: i64,
}

#[derive(Debug)]
pub struct AuthUser {
    pub user_id: Uuid,
    pub email: String,
    pub role: String,
}

#[rocket::async_trait]
impl<'r> FromRequest<'r> for AuthUser {
    type Error = AppError;

    async fn from_request(request: &'r Request<'_>) -> Outcome<Self, Self::Error> {
        // Get Authorization header
        let auth_header = request.headers().get_one("Authorization");

        let token = match auth_header {
            Some(header) if header.starts_with("Bearer ") => &header[7..],
            _ => {
                return Outcome::Error((
                    Status::Unauthorized,
                    AppError::Unauthorized("Missing or invalid Authorization header".to_string()),
                ));
            }
        };

        // Get config
        let config = AppConfig::default();

        // Verify token
        let token_data = match decode::<Claims>(
            token,
            &DecodingKey::from_secret(config.jwt_secret.as_bytes()),
            &Validation::default(),
        ) {
            Ok(data) => data,
            Err(e) => {
                return Outcome::Error((
                    Status::Unauthorized,
                    AppError::Unauthorized(format!("Invalid token: {}", e)),
                ));
            }
        };

        Outcome::Success(AuthUser {
            user_id: token_data.claims.sub,
            email: token_data.claims.email,
            role: token_data.claims.role,
        })
    }
}

// Admin-only guard
#[derive(Debug)]
pub struct AdminUser(pub AuthUser);

#[rocket::async_trait]
impl<'r> FromRequest<'r> for AdminUser {
    type Error = AppError;

    async fn from_request(request: &'r Request<'_>) -> Outcome<Self, Self::Error> {
        let auth_user = match AuthUser::from_request(request).await {
            Outcome::Success(user) => user,
            Outcome::Error(e) => return Outcome::Error(e),
            Outcome::Forward(s) => return Outcome::Forward(s),
        };

        if auth_user.role != "admin" {
            return Outcome::Error((
                Status::Forbidden,
                AppError::Unauthorized("Admin access required".to_string()),
            ));
        }

        Outcome::Success(AdminUser(auth_user))
    }
}

// Optional auth guard
#[derive(Debug)]
pub struct OptionalUser(pub Option<AuthUser>);

#[rocket::async_trait]
impl<'r> FromRequest<'r> for OptionalUser {
    type Error = std::convert::Infallible;

    async fn from_request(request: &'r Request<'_>) -> Outcome<Self, Self::Error> {
        match AuthUser::from_request(request).await {
            Outcome::Success(user) => Outcome::Success(OptionalUser(Some(user))),
            _ => Outcome::Success(OptionalUser(None)),
        }
    }
}
```

---

## Fairings (Middleware)

```rust
// src/fairings/logging.rs
use rocket::fairing::{Fairing, Info, Kind};
use rocket::{Data, Request, Response};
use std::time::Instant;

pub struct RequestLogger;

#[rocket::async_trait]
impl Fairing for RequestLogger {
    fn info(&self) -> Info {
        Info {
            name: "Request Logger",
            kind: Kind::Request | Kind::Response,
        }
    }

    async fn on_request(&self, request: &mut Request<'_>, _: &mut Data<'_>) {
        // Store request start time
        request.local_cache(|| Instant::now());
    }

    async fn on_response<'r>(&self, request: &'r Request<'_>, response: &mut Response<'r>) {
        let start_time = request.local_cache(|| Instant::now());
        let duration = start_time.elapsed();

        let method = request.method();
        let uri = request.uri();
        let status = response.status();

        println!(
            "{} {} {} - {:?}",
            method,
            uri,
            status.code,
            duration
        );
    }
}

// CORS Fairing
use rocket::http::Header;

pub struct Cors;

#[rocket::async_trait]
impl Fairing for Cors {
    fn info(&self) -> Info {
        Info {
            name: "CORS",
            kind: Kind::Response,
        }
    }

    async fn on_response<'r>(&self, _request: &'r Request<'_>, response: &mut Response<'r>) {
        response.set_header(Header::new("Access-Control-Allow-Origin", "*"));
        response.set_header(Header::new(
            "Access-Control-Allow-Methods",
            "GET, POST, PUT, DELETE, OPTIONS",
        ));
        response.set_header(Header::new(
            "Access-Control-Allow-Headers",
            "Content-Type, Authorization",
        ));
        response.set_header(Header::new("Access-Control-Allow-Credentials", "true"));
    }
}
```

---

## Services

```rust
// src/services/user_service.rs
use crate::config::AppConfig;
use crate::error::AppError;
use crate::guards::Claims;
use crate::models::user::{CreateUser, LoginRequest, LoginResponse, UpdateUser, User, UserResponse};
use bcrypt::{hash, verify, DEFAULT_COST};
use chrono::{Duration, Utc};
use jsonwebtoken::{encode, EncodingKey, Header};
use sqlx::PgPool;
use uuid::Uuid;
use validator::Validate;

pub struct UserService;

impl UserService {
    pub async fn create_user(pool: &PgPool, input: CreateUser) -> Result<UserResponse, AppError> {
        // Validate input
        input.validate().map_err(|e| AppError::Validation(e.to_string()))?;

        // Check if user exists
        let existing = sqlx::query_scalar::<_, i64>(
            "SELECT COUNT(*) FROM users WHERE email = $1"
        )
        .bind(&input.email)
        .fetch_one(pool)
        .await?;

        if existing > 0 {
            return Err(AppError::Validation("Email already registered".to_string()));
        }

        // Hash password
        let password_hash = hash(&input.password, DEFAULT_COST)
            .map_err(|e| AppError::Internal(e.to_string()))?;

        let now = Utc::now();
        let id = Uuid::new_v4();

        let user = sqlx::query_as::<_, User>(
            r#"
            INSERT INTO users (id, email, password_hash, name, role, created_at, updated_at)
            VALUES ($1, $2, $3, $4, 'user', $5, $5)
            RETURNING *
            "#,
        )
        .bind(id)
        .bind(&input.email)
        .bind(&password_hash)
        .bind(&input.name)
        .bind(now)
        .fetch_one(pool)
        .await?;

        Ok(user.into())
    }

    pub async fn login(pool: &PgPool, input: LoginRequest) -> Result<LoginResponse, AppError> {
        input.validate().map_err(|e| AppError::Validation(e.to_string()))?;

        let user = sqlx::query_as::<_, User>(
            "SELECT * FROM users WHERE email = $1"
        )
        .bind(&input.email)
        .fetch_optional(pool)
        .await?
        .ok_or_else(|| AppError::Unauthorized("Invalid credentials".to_string()))?;

        // Verify password
        let valid = verify(&input.password, &user.password_hash)
            .map_err(|e| AppError::Internal(e.to_string()))?;

        if !valid {
            return Err(AppError::Unauthorized("Invalid credentials".to_string()));
        }

        // Generate token
        let config = AppConfig::default();
        let expiration = Utc::now() + Duration::hours(config.jwt_expiration_hours);

        let claims = Claims {
            sub: user.id,
            email: user.email.clone(),
            role: user.role.clone(),
            exp: expiration.timestamp(),
        };

        let token = encode(
            &Header::default(),
            &claims,
            &EncodingKey::from_secret(config.jwt_secret.as_bytes()),
        )
        .map_err(|e| AppError::Internal(e.to_string()))?;

        Ok(LoginResponse {
            token,
            user: user.into(),
        })
    }

    pub async fn get_user(pool: &PgPool, user_id: Uuid) -> Result<UserResponse, AppError> {
        let user = sqlx::query_as::<_, User>(
            "SELECT * FROM users WHERE id = $1"
        )
        .bind(user_id)
        .fetch_optional(pool)
        .await?
        .ok_or_else(|| AppError::NotFound(format!("User {} not found", user_id)))?;

        Ok(user.into())
    }

    pub async fn update_user(
        pool: &PgPool,
        user_id: Uuid,
        input: UpdateUser,
    ) -> Result<UserResponse, AppError> {
        input.validate().map_err(|e| AppError::Validation(e.to_string()))?;

        let user = sqlx::query_as::<_, User>(
            r#"
            UPDATE users
            SET
                name = COALESCE($2, name),
                email = COALESCE($3, email),
                updated_at = $4
            WHERE id = $1
            RETURNING *
            "#,
        )
        .bind(user_id)
        .bind(&input.name)
        .bind(&input.email)
        .bind(Utc::now())
        .fetch_optional(pool)
        .await?
        .ok_or_else(|| AppError::NotFound(format!("User {} not found", user_id)))?;

        Ok(user.into())
    }

    pub async fn delete_user(pool: &PgPool, user_id: Uuid) -> Result<(), AppError> {
        let result = sqlx::query("DELETE FROM users WHERE id = $1")
            .bind(user_id)
            .execute(pool)
            .await?;

        if result.rows_affected() == 0 {
            return Err(AppError::NotFound(format!("User {} not found", user_id)));
        }

        Ok(())
    }

    pub async fn list_users(
        pool: &PgPool,
        limit: i64,
        offset: i64,
    ) -> Result<Vec<UserResponse>, AppError> {
        let users = sqlx::query_as::<_, User>(
            "SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2"
        )
        .bind(limit)
        .bind(offset)
        .fetch_all(pool)
        .await?;

        Ok(users.into_iter().map(|u| u.into()).collect())
    }
}
```

---

## Routes

```rust
// src/routes/mod.rs
use rocket::Route;

mod health;
mod users;

pub fn api_routes() -> Vec<Route> {
    routes![
        users::register,
        users::login,
        users::get_current_user,
        users::update_current_user,
        users::get_user,
        users::list_users,
        users::delete_user,
    ]
}

pub fn health_routes() -> Vec<Route> {
    routes![
        health::health_check,
        health::ready_check,
    ]
}

// src/routes/health.rs
use rocket::serde::json::Json;
use rocket::serde::Serialize;
use crate::db::DbConn;

#[derive(Serialize)]
#[serde(crate = "rocket::serde")]
pub struct HealthResponse {
    pub status: String,
    pub version: String,
}

#[get("/")]
pub fn health_check() -> Json<HealthResponse> {
    Json(HealthResponse {
        status: "ok".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
    })
}

#[get("/ready")]
pub async fn ready_check(mut db: DbConn) -> Result<Json<HealthResponse>, &'static str> {
    // Check database connection
    sqlx::query("SELECT 1")
        .execute(&mut **db)
        .await
        .map_err(|_| "Database not ready")?;

    Ok(Json(HealthResponse {
        status: "ready".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
    }))
}

// src/routes/users.rs
use crate::db::DbConn;
use crate::error::AppError;
use crate::guards::{AdminUser, AuthUser};
use crate::models::user::{CreateUser, LoginRequest, LoginResponse, UpdateUser, UserResponse};
use crate::services::UserService;
use rocket::serde::json::Json;
use uuid::Uuid;

#[post("/register", data = "<input>")]
pub async fn register(
    mut db: DbConn,
    input: Json<CreateUser>,
) -> Result<Json<UserResponse>, AppError> {
    let user = UserService::create_user(&mut **db, input.into_inner()).await?;
    Ok(Json(user))
}

#[post("/login", data = "<input>")]
pub async fn login(
    mut db: DbConn,
    input: Json<LoginRequest>,
) -> Result<Json<LoginResponse>, AppError> {
    let response = UserService::login(&mut **db, input.into_inner()).await?;
    Ok(Json(response))
}

#[get("/users/me")]
pub async fn get_current_user(
    mut db: DbConn,
    auth: AuthUser,
) -> Result<Json<UserResponse>, AppError> {
    let user = UserService::get_user(&mut **db, auth.user_id).await?;
    Ok(Json(user))
}

#[put("/users/me", data = "<input>")]
pub async fn update_current_user(
    mut db: DbConn,
    auth: AuthUser,
    input: Json<UpdateUser>,
) -> Result<Json<UserResponse>, AppError> {
    let user = UserService::update_user(&mut **db, auth.user_id, input.into_inner()).await?;
    Ok(Json(user))
}

#[get("/users/<id>")]
pub async fn get_user(
    mut db: DbConn,
    _auth: AuthUser,
    id: &str,
) -> Result<Json<UserResponse>, AppError> {
    let user_id = Uuid::parse_str(id)
        .map_err(|_| AppError::Validation("Invalid user ID".to_string()))?;
    let user = UserService::get_user(&mut **db, user_id).await?;
    Ok(Json(user))
}

#[get("/users?<limit>&<offset>")]
pub async fn list_users(
    mut db: DbConn,
    _auth: AdminUser,
    limit: Option<i64>,
    offset: Option<i64>,
) -> Result<Json<Vec<UserResponse>>, AppError> {
    let users = UserService::list_users(
        &mut **db,
        limit.unwrap_or(20),
        offset.unwrap_or(0),
    ).await?;
    Ok(Json(users))
}

#[delete("/users/<id>")]
pub async fn delete_user(
    mut db: DbConn,
    _auth: AdminUser,
    id: &str,
) -> Result<(), AppError> {
    let user_id = Uuid::parse_str(id)
        .map_err(|_| AppError::Validation("Invalid user ID".to_string()))?;
    UserService::delete_user(&mut **db, user_id).await?;
    Ok(())
}
```

---

## Testing

```rust
// tests/api_tests.rs
use rocket::http::{ContentType, Header, Status};
use rocket::local::asynchronous::Client;
use serde_json::{json, Value};

async fn get_client() -> Client {
    let rocket = myproject::rocket();
    Client::tracked(rocket).await.expect("valid rocket instance")
}

#[rocket::async_test]
async fn test_health_check() {
    let client = get_client().await;
    let response = client.get("/health").dispatch().await;

    assert_eq!(response.status(), Status::Ok);

    let body: Value = response.into_json().await.unwrap();
    assert_eq!(body["status"], "ok");
}

#[rocket::async_test]
async fn test_register_user() {
    let client = get_client().await;

    let response = client
        .post("/api/register")
        .header(ContentType::JSON)
        .body(json!({
            "email": "test@example.com",
            "password": "password123",
            "name": "Test User"
        }).to_string())
        .dispatch()
        .await;

    assert_eq!(response.status(), Status::Ok);

    let body: Value = response.into_json().await.unwrap();
    assert_eq!(body["email"], "test@example.com");
    assert_eq!(body["name"], "Test User");
}

#[rocket::async_test]
async fn test_register_invalid_email() {
    let client = get_client().await;

    let response = client
        .post("/api/register")
        .header(ContentType::JSON)
        .body(json!({
            "email": "invalid-email",
            "password": "password123",
            "name": "Test User"
        }).to_string())
        .dispatch()
        .await;

    assert_eq!(response.status(), Status::BadRequest);
}

#[rocket::async_test]
async fn test_login_and_access() {
    let client = get_client().await;

    // Register
    client
        .post("/api/register")
        .header(ContentType::JSON)
        .body(json!({
            "email": "login_test@example.com",
            "password": "password123",
            "name": "Login Test"
        }).to_string())
        .dispatch()
        .await;

    // Login
    let login_response = client
        .post("/api/login")
        .header(ContentType::JSON)
        .body(json!({
            "email": "login_test@example.com",
            "password": "password123"
        }).to_string())
        .dispatch()
        .await;

    assert_eq!(login_response.status(), Status::Ok);

    let login_body: Value = login_response.into_json().await.unwrap();
    let token = login_body["token"].as_str().unwrap();

    // Access protected route
    let me_response = client
        .get("/api/users/me")
        .header(Header::new("Authorization", format!("Bearer {}", token)))
        .dispatch()
        .await;

    assert_eq!(me_response.status(), Status::Ok);

    let me_body: Value = me_response.into_json().await.unwrap();
    assert_eq!(me_body["email"], "login_test@example.com");
}

#[rocket::async_test]
async fn test_unauthorized_access() {
    let client = get_client().await;

    let response = client.get("/api/users/me").dispatch().await;

    assert_eq!(response.status(), Status::Unauthorized);
}

#[rocket::async_test]
async fn test_invalid_token() {
    let client = get_client().await;

    let response = client
        .get("/api/users/me")
        .header(Header::new("Authorization", "Bearer invalid_token"))
        .dispatch()
        .await;

    assert_eq!(response.status(), Status::Unauthorized);
}
```

---

## Form Handling

```rust
// Rocket excels at form handling
use rocket::form::{Form, FromForm};
use rocket::fs::TempFile;

#[derive(FromForm)]
pub struct UploadForm<'r> {
    pub name: String,
    pub description: Option<String>,
    pub file: TempFile<'r>,
}

#[post("/upload", data = "<form>")]
pub async fn upload(
    mut form: Form<UploadForm<'_>>,
    auth: AuthUser,
) -> Result<Json<Value>, AppError> {
    // Save the file
    let filename = format!("uploads/{}_{}", auth.user_id, form.name);
    form.file.persist_to(&filename).await
        .map_err(|e| AppError::Internal(e.to_string()))?;

    Ok(Json(json!({
        "status": "uploaded",
        "filename": filename
    })))
}

// Multipart form with validation
#[derive(FromForm)]
pub struct ProfileForm<'r> {
    #[field(validate = len(1..100))]
    pub name: String,

    #[field(validate = contains('@'))]
    pub email: String,

    #[field(validate = len(..1_000_000))]  // Max 1MB
    pub avatar: Option<TempFile<'r>>,
}
```

---

## Commands

```bash
# Development
cargo run

# Release build
cargo build --release

# Run tests
cargo test

# Run specific test
cargo test test_health_check

# Format code
cargo fmt

# Lint
cargo clippy

# Check
cargo check

# Generate docs
cargo doc --open

# Run with logging
ROCKET_LOG_LEVEL=debug cargo run
```

---

## Best Practices

### DO
- ✓ Use request guards for authentication/authorization
- ✓ Implement `Responder` for custom error types
- ✓ Use fairings for cross-cutting concerns
- ✓ Leverage Rocket's type system for validation
- ✓ Use `rocket_db_pools` for database connections
- ✓ Configure via `Rocket.toml` for different environments
- ✓ Use derive macros to reduce boilerplate
- ✓ Test with `rocket::local::asynchronous::Client`

### DON'T
- ✗ Use unwrap/expect in production code
- ✗ Ignore validation errors
- ✗ Store secrets in `Rocket.toml` (use environment variables)
- ✗ Block async code with synchronous operations
- ✗ Skip error catchers for common HTTP errors
- ✗ Use global state when request-local state suffices

---

## Rocket vs Axum vs Actix-web Comparison

| Feature | Rocket | Axum | Actix-web |
|---------|--------|------|-----------|
| **Learning Curve** | Gentle | Moderate | Moderate |
| **Type Safety** | Excellent | Excellent | Good |
| **Performance** | Good | Excellent | Excellent |
| **Ecosystem** | Growing | Tower ecosystem | Mature |
| **Macros** | Heavy use | Minimal | Moderate |
| **Forms/Uploads** | Excellent | Good | Good |
| **Configuration** | Built-in TOML | Manual | Manual |
| **Compile Time** | Slower | Faster | Fast |

**Choose Rocket when**: Type safety and developer experience are top priorities, or you need excellent form handling.

---

## Migration Guide

### From Actix-web to Rocket

```rust
// Actix-web
#[get("/users/{id}")]
async fn get_user(path: web::Path<Uuid>) -> impl Responder {
    // ...
}

// Rocket
#[get("/users/<id>")]
async fn get_user(id: &str) -> Result<Json<UserResponse>, AppError> {
    let user_id = Uuid::parse_str(id)?;
    // ...
}
```

### From Axum to Rocket

```rust
// Axum
async fn get_user(
    State(pool): State<PgPool>,
    Path(id): Path<Uuid>,
) -> Result<Json<User>, AppError> {
    // ...
}

// Rocket
#[get("/users/<id>")]
async fn get_user(
    mut db: DbConn,
    id: &str,
) -> Result<Json<UserResponse>, AppError> {
    let user_id = Uuid::parse_str(id)?;
    // ...
}
```

---

## References

- [Rocket Documentation](https://rocket.rs/v0.5/)
- [Rocket Guide](https://rocket.rs/v0.5/guide/)
- [Rocket API Reference](https://api.rocket.rs/v0.5/)
- [rocket_db_pools](https://crates.io/crates/rocket_db_pools)
- [Rocket GitHub](https://github.com/SergioBenitez/Rocket)
- [Rocket Examples](https://github.com/SergioBenitez/Rocket/tree/v0.5/examples)
