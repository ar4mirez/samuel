# Axum Framework Guide

> **Applies to**: Axum 0.7+, Rust Web APIs, Microservices
> **Complements**: `.agent/skills/rust-guide/SKILL.md`

---

## Overview

Axum is a modern, ergonomic web framework built on top of Tokio and Tower. It provides a modular architecture with compile-time extraction of request data and first-class support for async/await.

### Key Features
- Built on Tokio and Tower ecosystem
- Type-safe extractors for request data
- Compile-time route checking
- Minimal boilerplate with maximum flexibility
- First-class middleware support via Tower
- WebSocket and SSE support built-in
- Excellent error handling with custom error types

---

## Project Structure

```
myproject/
├── Cargo.toml
├── src/
│   ├── main.rs
│   ├── lib.rs
│   ├── config.rs
│   ├── routes/
│   │   ├── mod.rs
│   │   ├── users.rs
│   │   └── health.rs
│   ├── handlers/
│   │   ├── mod.rs
│   │   └── users.rs
│   ├── models/
│   │   ├── mod.rs
│   │   └── user.rs
│   ├── services/
│   │   ├── mod.rs
│   │   └── user_service.rs
│   ├── repositories/
│   │   ├── mod.rs
│   │   └── user_repository.rs
│   ├── extractors/
│   │   ├── mod.rs
│   │   └── auth.rs
│   ├── middleware/
│   │   ├── mod.rs
│   │   └── logging.rs
│   └── errors/
│       ├── mod.rs
│       └── app_error.rs
├── tests/
│   └── integration_tests.rs
└── migrations/
```

---

## Dependencies (Cargo.toml)

```toml
[package]
name = "myproject"
version = "0.1.0"
edition = "2021"

[dependencies]
# Web framework
axum = { version = "0.7", features = ["macros"] }
axum-extra = { version = "0.9", features = ["typed-header"] }
tokio = { version = "1.0", features = ["full"] }
tower = { version = "0.4", features = ["full"] }
tower-http = { version = "0.5", features = ["cors", "trace", "compression-gzip"] }

# Serialization
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"

# Database
sqlx = { version = "0.7", features = ["runtime-tokio", "postgres", "uuid", "chrono"] }

# Validation
validator = { version = "0.16", features = ["derive"] }

# Authentication
jsonwebtoken = "9.0"

# Error handling
thiserror = "1.0"
anyhow = "1.0"

# Logging
tracing = "0.1"
tracing-subscriber = { version = "0.3", features = ["env-filter"] }

# Configuration
config = "0.14"
dotenvy = "0.15"

# Utilities
uuid = { version = "1.0", features = ["v4", "serde"] }
chrono = { version = "0.4", features = ["serde"] }

[dev-dependencies]
axum-test = "14.0"
mockall = "0.12"
```

---

## Application Entry Point

```rust
// src/main.rs
use std::net::SocketAddr;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

mod config;
mod errors;
mod extractors;
mod handlers;
mod middleware;
mod models;
mod repositories;
mod routes;
mod services;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // Load environment variables
    dotenvy::dotenv().ok();

    // Initialize tracing
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| "myproject=debug,tower_http=debug,axum=trace".into()),
        )
        .with(tracing_subscriber::fmt::layer())
        .init();

    // Load configuration
    let config = config::Config::load()?;

    // Create database pool
    let pool = sqlx::PgPool::connect(&config.database_url).await?;

    // Run migrations
    sqlx::migrate!("./migrations").run(&pool).await?;

    // Build application state
    let state = AppState::new(pool, config.clone());

    // Build router
    let app = routes::create_router(state);

    // Start server
    let addr = SocketAddr::from(([0, 0, 0, 0], config.port));
    tracing::info!("Starting server on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await?;
    axum::serve(listener, app)
        .with_graceful_shutdown(shutdown_signal())
        .await?;

    Ok(())
}

async fn shutdown_signal() {
    tokio::signal::ctrl_c()
        .await
        .expect("Failed to install CTRL+C signal handler");
    tracing::info!("Shutdown signal received");
}

// Shared application state
#[derive(Clone)]
pub struct AppState {
    pub pool: sqlx::PgPool,
    pub config: config::Config,
}

impl AppState {
    pub fn new(pool: sqlx::PgPool, config: config::Config) -> Self {
        Self { pool, config }
    }
}
```

---

## Configuration

```rust
// src/config.rs
use serde::Deserialize;

#[derive(Clone, Deserialize)]
pub struct Config {
    pub port: u16,
    pub database_url: String,
    pub jwt_secret: String,
    pub jwt_expiration_hours: i64,
}

impl Config {
    pub fn load() -> anyhow::Result<Self> {
        let config = config::Config::builder()
            .add_source(config::Environment::default())
            .build()?;

        Ok(config.try_deserialize()?)
    }
}
```

---

## Router Setup

```rust
// src/routes/mod.rs
use axum::{
    middleware,
    routing::{get, post, put, delete},
    Router,
};
use tower_http::{
    compression::CompressionLayer,
    cors::{Any, CorsLayer},
    trace::TraceLayer,
};

use crate::{handlers, middleware as app_middleware, AppState};

pub mod health;
pub mod users;

pub fn create_router(state: AppState) -> Router {
    let cors = CorsLayer::new()
        .allow_origin(Any)
        .allow_methods(Any)
        .allow_headers(Any);

    let public_routes = Router::new()
        .route("/health", get(handlers::health::health_check))
        .route("/api/v1/auth/register", post(handlers::users::register))
        .route("/api/v1/auth/login", post(handlers::users::login));

    let protected_routes = Router::new()
        .route("/api/v1/users", get(handlers::users::list_users))
        .route("/api/v1/users/:id", get(handlers::users::get_user))
        .route("/api/v1/users/:id", put(handlers::users::update_user))
        .route("/api/v1/users/:id", delete(handlers::users::delete_user))
        .route("/api/v1/users/me", get(handlers::users::get_current_user))
        .layer(middleware::from_fn_with_state(
            state.clone(),
            app_middleware::auth::auth_middleware,
        ));

    Router::new()
        .merge(public_routes)
        .merge(protected_routes)
        .layer(TraceLayer::new_for_http())
        .layer(CompressionLayer::new())
        .layer(cors)
        .with_state(state)
}
```

---

## Error Handling

```rust
// src/errors/app_error.rs
use axum::{
    http::StatusCode,
    response::{IntoResponse, Response},
    Json,
};
use serde_json::json;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum AppError {
    #[error("Not found: {0}")]
    NotFound(String),

    #[error("Bad request: {0}")]
    BadRequest(String),

    #[error("Unauthorized: {0}")]
    Unauthorized(String),

    #[error("Forbidden: {0}")]
    Forbidden(String),

    #[error("Conflict: {0}")]
    Conflict(String),

    #[error("Validation error: {0}")]
    Validation(String),

    #[error("Internal server error")]
    Internal(#[from] anyhow::Error),

    #[error("Database error")]
    Database(#[from] sqlx::Error),
}

impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let (status, error_message) = match &self {
            AppError::NotFound(msg) => (StatusCode::NOT_FOUND, msg.clone()),
            AppError::BadRequest(msg) => (StatusCode::BAD_REQUEST, msg.clone()),
            AppError::Unauthorized(msg) => (StatusCode::UNAUTHORIZED, msg.clone()),
            AppError::Forbidden(msg) => (StatusCode::FORBIDDEN, msg.clone()),
            AppError::Conflict(msg) => (StatusCode::CONFLICT, msg.clone()),
            AppError::Validation(msg) => (StatusCode::UNPROCESSABLE_ENTITY, msg.clone()),
            AppError::Internal(err) => {
                tracing::error!("Internal error: {:?}", err);
                (StatusCode::INTERNAL_SERVER_ERROR, "Internal server error".to_string())
            }
            AppError::Database(err) => {
                tracing::error!("Database error: {:?}", err);
                (StatusCode::INTERNAL_SERVER_ERROR, "Database error".to_string())
            }
        };

        let body = Json(json!({
            "success": false,
            "error": {
                "code": status.as_u16(),
                "message": error_message
            }
        }));

        (status, body).into_response()
    }
}

pub type AppResult<T> = Result<T, AppError>;
```

---

## Models

```rust
// src/models/user.rs
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use uuid::Uuid;
use validator::Validate;

#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct User {
    pub id: Uuid,
    pub email: String,
    #[serde(skip_serializing)]
    pub password_hash: String,
    pub name: String,
    pub role: String,
    pub is_active: bool,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Deserialize, Validate)]
pub struct CreateUserDto {
    #[validate(email(message = "Invalid email format"))]
    pub email: String,

    #[validate(length(min = 8, message = "Password must be at least 8 characters"))]
    pub password: String,

    #[validate(length(min = 1, max = 100, message = "Name must be 1-100 characters"))]
    pub name: String,
}

#[derive(Debug, Deserialize, Validate)]
pub struct UpdateUserDto {
    #[validate(length(min = 1, max = 100, message = "Name must be 1-100 characters"))]
    pub name: Option<String>,

    pub is_active: Option<bool>,
}

#[derive(Debug, Deserialize, Validate)]
pub struct LoginDto {
    #[validate(email(message = "Invalid email format"))]
    pub email: String,

    #[validate(length(min = 1, message = "Password is required"))]
    pub password: String,
}

#[derive(Debug, Serialize)]
pub struct UserResponse {
    pub id: Uuid,
    pub email: String,
    pub name: String,
    pub role: String,
    pub is_active: bool,
    pub created_at: DateTime<Utc>,
}

impl From<User> for UserResponse {
    fn from(user: User) -> Self {
        Self {
            id: user.id,
            email: user.email,
            name: user.name,
            role: user.role,
            is_active: user.is_active,
            created_at: user.created_at,
        }
    }
}

#[derive(Debug, Serialize)]
pub struct AuthResponse {
    pub user: UserResponse,
    pub token: String,
}
```

---

## Extractors

```rust
// src/extractors/auth.rs
use axum::{
    async_trait,
    extract::FromRequestParts,
    http::{header::AUTHORIZATION, request::Parts, StatusCode},
    RequestPartsExt,
};
use axum_extra::TypedHeader;
use headers::{authorization::Bearer, Authorization};
use jsonwebtoken::{decode, encode, DecodingKey, EncodingKey, Header, Validation};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

use crate::{errors::AppError, AppState};

#[derive(Debug, Serialize, Deserialize)]
pub struct Claims {
    pub sub: Uuid,
    pub email: String,
    pub role: String,
    pub exp: i64,
    pub iat: i64,
}

impl Claims {
    pub fn new(user_id: Uuid, email: String, role: String, expiration_hours: i64) -> Self {
        let now = chrono::Utc::now();
        Self {
            sub: user_id,
            email,
            role,
            exp: (now + chrono::Duration::hours(expiration_hours)).timestamp(),
            iat: now.timestamp(),
        }
    }
}

pub fn create_token(claims: &Claims, secret: &str) -> Result<String, AppError> {
    encode(
        &Header::default(),
        claims,
        &EncodingKey::from_secret(secret.as_bytes()),
    )
    .map_err(|e| AppError::Internal(anyhow::anyhow!("Failed to create token: {}", e)))
}

pub fn verify_token(token: &str, secret: &str) -> Result<Claims, AppError> {
    decode::<Claims>(
        token,
        &DecodingKey::from_secret(secret.as_bytes()),
        &Validation::default(),
    )
    .map(|data| data.claims)
    .map_err(|_| AppError::Unauthorized("Invalid token".to_string()))
}

// Custom extractor for authenticated user
pub struct AuthUser {
    pub user_id: Uuid,
    pub email: String,
    pub role: String,
}

#[async_trait]
impl FromRequestParts<AppState> for AuthUser {
    type Rejection = AppError;

    async fn from_request_parts(
        parts: &mut Parts,
        state: &AppState,
    ) -> Result<Self, Self::Rejection> {
        let TypedHeader(Authorization(bearer)) = parts
            .extract::<TypedHeader<Authorization<Bearer>>>()
            .await
            .map_err(|_| AppError::Unauthorized("Missing authorization header".to_string()))?;

        let claims = verify_token(bearer.token(), &state.config.jwt_secret)?;

        Ok(AuthUser {
            user_id: claims.sub,
            email: claims.email,
            role: claims.role,
        })
    }
}
```

---

## Middleware

```rust
// src/middleware/auth.rs
use axum::{
    extract::State,
    http::{header::AUTHORIZATION, Request},
    middleware::Next,
    response::Response,
};

use crate::{errors::AppError, extractors::auth::verify_token, AppState};

pub async fn auth_middleware<B>(
    State(state): State<AppState>,
    request: Request<B>,
    next: Next<B>,
) -> Result<Response, AppError> {
    let auth_header = request
        .headers()
        .get(AUTHORIZATION)
        .and_then(|header| header.to_str().ok())
        .ok_or_else(|| AppError::Unauthorized("Missing authorization header".to_string()))?;

    let token = auth_header
        .strip_prefix("Bearer ")
        .ok_or_else(|| AppError::Unauthorized("Invalid authorization format".to_string()))?;

    verify_token(token, &state.config.jwt_secret)?;

    Ok(next.run(request).await)
}
```

---

## Repository

```rust
// src/repositories/user_repository.rs
use sqlx::PgPool;
use uuid::Uuid;

use crate::{errors::AppError, models::user::User};

pub struct UserRepository {
    pool: PgPool,
}

impl UserRepository {
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }

    pub async fn find_by_id(&self, id: Uuid) -> Result<Option<User>, AppError> {
        let user = sqlx::query_as::<_, User>(
            r#"
            SELECT id, email, password_hash, name, role, is_active, created_at, updated_at
            FROM users
            WHERE id = $1
            "#,
        )
        .bind(id)
        .fetch_optional(&self.pool)
        .await?;

        Ok(user)
    }

    pub async fn find_by_email(&self, email: &str) -> Result<Option<User>, AppError> {
        let user = sqlx::query_as::<_, User>(
            r#"
            SELECT id, email, password_hash, name, role, is_active, created_at, updated_at
            FROM users
            WHERE email = $1
            "#,
        )
        .bind(email)
        .fetch_optional(&self.pool)
        .await?;

        Ok(user)
    }

    pub async fn find_all(&self, limit: i64, offset: i64) -> Result<Vec<User>, AppError> {
        let users = sqlx::query_as::<_, User>(
            r#"
            SELECT id, email, password_hash, name, role, is_active, created_at, updated_at
            FROM users
            ORDER BY created_at DESC
            LIMIT $1 OFFSET $2
            "#,
        )
        .bind(limit)
        .bind(offset)
        .fetch_all(&self.pool)
        .await?;

        Ok(users)
    }

    pub async fn create(
        &self,
        email: &str,
        password_hash: &str,
        name: &str,
    ) -> Result<User, AppError> {
        let user = sqlx::query_as::<_, User>(
            r#"
            INSERT INTO users (email, password_hash, name)
            VALUES ($1, $2, $3)
            RETURNING id, email, password_hash, name, role, is_active, created_at, updated_at
            "#,
        )
        .bind(email)
        .bind(password_hash)
        .bind(name)
        .fetch_one(&self.pool)
        .await?;

        Ok(user)
    }

    pub async fn update(&self, id: Uuid, name: Option<&str>, is_active: Option<bool>) -> Result<User, AppError> {
        let user = sqlx::query_as::<_, User>(
            r#"
            UPDATE users
            SET
                name = COALESCE($2, name),
                is_active = COALESCE($3, is_active),
                updated_at = NOW()
            WHERE id = $1
            RETURNING id, email, password_hash, name, role, is_active, created_at, updated_at
            "#,
        )
        .bind(id)
        .bind(name)
        .bind(is_active)
        .fetch_one(&self.pool)
        .await?;

        Ok(user)
    }

    pub async fn delete(&self, id: Uuid) -> Result<(), AppError> {
        sqlx::query("DELETE FROM users WHERE id = $1")
            .bind(id)
            .execute(&self.pool)
            .await?;

        Ok(())
    }
}
```

---

## Service Layer

```rust
// src/services/user_service.rs
use sqlx::PgPool;
use uuid::Uuid;
use validator::Validate;

use crate::{
    config::Config,
    errors::AppError,
    extractors::auth::{create_token, Claims},
    models::user::{AuthResponse, CreateUserDto, LoginDto, UpdateUserDto, User, UserResponse},
    repositories::user_repository::UserRepository,
};

pub struct UserService {
    repository: UserRepository,
    config: Config,
}

impl UserService {
    pub fn new(pool: PgPool, config: Config) -> Self {
        Self {
            repository: UserRepository::new(pool),
            config,
        }
    }

    pub async fn register(&self, dto: CreateUserDto) -> Result<AuthResponse, AppError> {
        dto.validate()
            .map_err(|e| AppError::Validation(e.to_string()))?;

        // Check if email already exists
        if self.repository.find_by_email(&dto.email).await?.is_some() {
            return Err(AppError::Conflict("Email already registered".to_string()));
        }

        // Hash password
        let password_hash = hash_password(&dto.password)?;

        // Create user
        let user = self
            .repository
            .create(&dto.email, &password_hash, &dto.name)
            .await?;

        // Generate token
        let claims = Claims::new(
            user.id,
            user.email.clone(),
            user.role.clone(),
            self.config.jwt_expiration_hours,
        );
        let token = create_token(&claims, &self.config.jwt_secret)?;

        Ok(AuthResponse {
            user: user.into(),
            token,
        })
    }

    pub async fn login(&self, dto: LoginDto) -> Result<AuthResponse, AppError> {
        dto.validate()
            .map_err(|e| AppError::Validation(e.to_string()))?;

        let user = self
            .repository
            .find_by_email(&dto.email)
            .await?
            .ok_or_else(|| AppError::Unauthorized("Invalid credentials".to_string()))?;

        if !user.is_active {
            return Err(AppError::Forbidden("Account is deactivated".to_string()));
        }

        if !verify_password(&dto.password, &user.password_hash)? {
            return Err(AppError::Unauthorized("Invalid credentials".to_string()));
        }

        let claims = Claims::new(
            user.id,
            user.email.clone(),
            user.role.clone(),
            self.config.jwt_expiration_hours,
        );
        let token = create_token(&claims, &self.config.jwt_secret)?;

        Ok(AuthResponse {
            user: user.into(),
            token,
        })
    }

    pub async fn get_user(&self, id: Uuid) -> Result<UserResponse, AppError> {
        let user = self
            .repository
            .find_by_id(id)
            .await?
            .ok_or_else(|| AppError::NotFound(format!("User {} not found", id)))?;

        Ok(user.into())
    }

    pub async fn list_users(&self, page: i64, per_page: i64) -> Result<Vec<UserResponse>, AppError> {
        let offset = (page - 1) * per_page;
        let users = self.repository.find_all(per_page, offset).await?;
        Ok(users.into_iter().map(|u| u.into()).collect())
    }

    pub async fn update_user(&self, id: Uuid, dto: UpdateUserDto) -> Result<UserResponse, AppError> {
        dto.validate()
            .map_err(|e| AppError::Validation(e.to_string()))?;

        // Check if user exists
        self.repository
            .find_by_id(id)
            .await?
            .ok_or_else(|| AppError::NotFound(format!("User {} not found", id)))?;

        let user = self
            .repository
            .update(id, dto.name.as_deref(), dto.is_active)
            .await?;

        Ok(user.into())
    }

    pub async fn delete_user(&self, id: Uuid) -> Result<(), AppError> {
        self.repository
            .find_by_id(id)
            .await?
            .ok_or_else(|| AppError::NotFound(format!("User {} not found", id)))?;

        self.repository.delete(id).await?;
        Ok(())
    }
}

fn hash_password(password: &str) -> Result<String, AppError> {
    bcrypt::hash(password, bcrypt::DEFAULT_COST)
        .map_err(|e| AppError::Internal(anyhow::anyhow!("Failed to hash password: {}", e)))
}

fn verify_password(password: &str, hash: &str) -> Result<bool, AppError> {
    bcrypt::verify(password, hash)
        .map_err(|e| AppError::Internal(anyhow::anyhow!("Failed to verify password: {}", e)))
}
```

---

## Handlers

```rust
// src/handlers/users.rs
use axum::{
    extract::{Path, Query, State},
    http::StatusCode,
    Json,
};
use serde::Deserialize;
use uuid::Uuid;

use crate::{
    errors::{AppError, AppResult},
    extractors::auth::AuthUser,
    models::user::{AuthResponse, CreateUserDto, LoginDto, UpdateUserDto, UserResponse},
    services::user_service::UserService,
    AppState,
};

#[derive(Deserialize)]
pub struct PaginationQuery {
    #[serde(default = "default_page")]
    pub page: i64,
    #[serde(default = "default_per_page")]
    pub per_page: i64,
}

fn default_page() -> i64 { 1 }
fn default_per_page() -> i64 { 20 }

pub async fn register(
    State(state): State<AppState>,
    Json(dto): Json<CreateUserDto>,
) -> AppResult<(StatusCode, Json<AuthResponse>)> {
    let service = UserService::new(state.pool, state.config);
    let response = service.register(dto).await?;
    Ok((StatusCode::CREATED, Json(response)))
}

pub async fn login(
    State(state): State<AppState>,
    Json(dto): Json<LoginDto>,
) -> AppResult<Json<AuthResponse>> {
    let service = UserService::new(state.pool, state.config);
    let response = service.login(dto).await?;
    Ok(Json(response))
}

pub async fn get_current_user(
    State(state): State<AppState>,
    auth_user: AuthUser,
) -> AppResult<Json<UserResponse>> {
    let service = UserService::new(state.pool, state.config);
    let user = service.get_user(auth_user.user_id).await?;
    Ok(Json(user))
}

pub async fn list_users(
    State(state): State<AppState>,
    Query(pagination): Query<PaginationQuery>,
    _auth_user: AuthUser,
) -> AppResult<Json<Vec<UserResponse>>> {
    let service = UserService::new(state.pool, state.config);
    let users = service.list_users(pagination.page, pagination.per_page).await?;
    Ok(Json(users))
}

pub async fn get_user(
    State(state): State<AppState>,
    Path(id): Path<Uuid>,
    _auth_user: AuthUser,
) -> AppResult<Json<UserResponse>> {
    let service = UserService::new(state.pool, state.config);
    let user = service.get_user(id).await?;
    Ok(Json(user))
}

pub async fn update_user(
    State(state): State<AppState>,
    Path(id): Path<Uuid>,
    _auth_user: AuthUser,
    Json(dto): Json<UpdateUserDto>,
) -> AppResult<Json<UserResponse>> {
    let service = UserService::new(state.pool, state.config);
    let user = service.update_user(id, dto).await?;
    Ok(Json(user))
}

pub async fn delete_user(
    State(state): State<AppState>,
    Path(id): Path<Uuid>,
    _auth_user: AuthUser,
) -> AppResult<StatusCode> {
    let service = UserService::new(state.pool, state.config);
    service.delete_user(id).await?;
    Ok(StatusCode::NO_CONTENT)
}

// Health check handler
pub mod health {
    use axum::Json;
    use serde_json::{json, Value};

    pub async fn health_check() -> Json<Value> {
        Json(json!({
            "status": "healthy",
            "timestamp": chrono::Utc::now().to_rfc3339()
        }))
    }
}
```

---

## Testing

```rust
// tests/integration_tests.rs
use axum::{
    body::Body,
    http::{Request, StatusCode},
};
use serde_json::json;
use tower::ServiceExt;

use myproject::{routes::create_router, AppState};

async fn setup_test_app() -> axum::Router {
    // Create test database pool
    let pool = sqlx::PgPool::connect("postgres://test:test@localhost/test_db")
        .await
        .unwrap();

    let config = myproject::config::Config {
        port: 3000,
        database_url: String::new(),
        jwt_secret: "test_secret".to_string(),
        jwt_expiration_hours: 24,
    };

    let state = AppState::new(pool, config);
    create_router(state)
}

#[tokio::test]
async fn test_health_check() {
    let app = setup_test_app().await;

    let response = app
        .oneshot(Request::builder().uri("/health").body(Body::empty()).unwrap())
        .await
        .unwrap();

    assert_eq!(response.status(), StatusCode::OK);
}

#[tokio::test]
async fn test_register_user() {
    let app = setup_test_app().await;

    let body = json!({
        "email": "test@example.com",
        "password": "password123",
        "name": "Test User"
    });

    let response = app
        .oneshot(
            Request::builder()
                .method("POST")
                .uri("/api/v1/auth/register")
                .header("Content-Type", "application/json")
                .body(Body::from(serde_json::to_string(&body).unwrap()))
                .unwrap(),
        )
        .await
        .unwrap();

    assert_eq!(response.status(), StatusCode::CREATED);
}

#[tokio::test]
async fn test_login_invalid_credentials() {
    let app = setup_test_app().await;

    let body = json!({
        "email": "nonexistent@example.com",
        "password": "wrongpassword"
    });

    let response = app
        .oneshot(
            Request::builder()
                .method("POST")
                .uri("/api/v1/auth/login")
                .header("Content-Type", "application/json")
                .body(Body::from(serde_json::to_string(&body).unwrap()))
                .unwrap(),
        )
        .await
        .unwrap();

    assert_eq!(response.status(), StatusCode::UNAUTHORIZED);
}
```

---

## Database Migration

```sql
-- migrations/001_create_users.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
```

---

## Commands

```bash
# Development
cargo run

# Watch mode (with cargo-watch)
cargo watch -x run

# Build release
cargo build --release

# Run tests
cargo test

# Run with specific log level
RUST_LOG=debug cargo run

# Check code
cargo check
cargo clippy

# Format code
cargo fmt

# Run migrations
sqlx migrate run

# Create new migration
sqlx migrate add create_users
```

---

## Best Practices

### Extractors
- Use custom extractors for request validation
- Implement `FromRequestParts` for reusable extraction logic
- Use `State` for shared application state
- Use `Path`, `Query`, `Json` for typed request data

### Error Handling
- Implement `IntoResponse` for custom error types
- Use `thiserror` for error definitions
- Return appropriate HTTP status codes
- Log errors with context using `tracing`

### Performance
- Use connection pooling with SQLx
- Enable compression with `tower-http`
- Use async/await throughout
- Consider using `tower::buffer` for backpressure

### Testing
- Use `axum-test` crate for integration testing
- Mock repositories for unit tests
- Test error cases and edge cases
- Use separate test database

---

## Comparison: Axum vs Actix-web vs Rocket

| Feature | Axum | Actix-web | Rocket |
|---------|------|-----------|--------|
| Async Runtime | Tokio | Actix (Tokio) | Tokio |
| Type Safety | Compile-time extractors | Runtime | Compile-time |
| Middleware | Tower | Custom | Fairings |
| Performance | Very Fast | Fastest | Fast |
| Ecosystem | Tower/Hyper | Actix | Rocket |
| Learning Curve | Moderate | Moderate | Gentle |
| Maturity | Newer | Mature | Mature |

---

## References

- [Axum Documentation](https://docs.rs/axum)
- [Tower Documentation](https://docs.rs/tower)
- [Tokio Documentation](https://tokio.rs)
- [SQLx Documentation](https://docs.rs/sqlx)
- [Axum Examples](https://github.com/tokio-rs/axum/tree/main/examples)
