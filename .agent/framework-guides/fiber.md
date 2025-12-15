# Fiber Framework Guide

> **Applies to**: Fiber v2.50+, Go 1.21+, High-Performance REST APIs, Microservices
> **Load when**: Working with Fiber projects, Express-like Go APIs, high-throughput services

---

## Overview

Fiber is an Express-inspired web framework built on top of Fasthttp, the fastest HTTP engine for Go. It's designed for ease of use with zero memory allocation and performance in mind.

**Key Features**:
- Express-like API (familiar to Node.js developers)
- Built on Fasthttp (10x faster than net/http)
- Zero memory allocation in hot paths
- Built-in middleware collection
- WebSocket support
- Rate limiting
- Template engines

**When to use Fiber**:
- High-throughput APIs requiring maximum performance
- Teams familiar with Express.js
- Real-time applications with WebSockets
- Microservices requiring low latency

---

## Project Structure

```
myapi/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── handler/
│   │   ├── handler.go           # Handler container
│   │   ├── user_handler.go      # User handlers
│   │   └── health_handler.go    # Health check handlers
│   ├── middleware/
│   │   ├── auth.go              # JWT authentication
│   │   ├── logger.go            # Request logging
│   │   └── recover.go           # Panic recovery
│   ├── model/
│   │   ├── user.go              # User model
│   │   └── dto.go               # Data transfer objects
│   ├── repository/
│   │   ├── repository.go        # Repository interface
│   │   └── user_repository.go   # User repository
│   ├── service/
│   │   ├── service.go           # Service container
│   │   ├── user_service.go      # User business logic
│   │   └── auth_service.go      # Authentication service
│   └── router/
│       └── router.go            # Route definitions
├── pkg/
│   ├── validator/
│   │   └── validator.go         # Custom validator
│   └── response/
│       └── response.go          # Response helpers
├── go.mod
├── go.sum
└── README.md
```

---

## Core Patterns

### Application Setup

```go
// cmd/api/main.go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/helmet"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "github.com/gofiber/fiber/v2/middleware/requestid"

    "myapi/internal/config"
    "myapi/internal/handler"
    "myapi/internal/middleware"
    "myapi/internal/repository"
    "myapi/internal/router"
    "myapi/internal/service"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize database
    db, err := config.NewDatabase(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Initialize Fiber app with custom config
    app := fiber.New(fiber.Config{
        AppName:       cfg.AppName,
        ServerHeader:  cfg.AppName,
        ReadTimeout:   cfg.ReadTimeout,
        WriteTimeout:  cfg.WriteTimeout,
        IdleTimeout:   cfg.IdleTimeout,
        BodyLimit:     cfg.BodyLimit,
        Prefork:       cfg.Prefork,  // Enable prefork for multi-core
        ErrorHandler:  customErrorHandler,
    })

    // Global middleware
    app.Use(requestid.New())
    app.Use(middleware.Logger())
    app.Use(middleware.Recover())
    app.Use(cors.New(cors.Config{
        AllowOrigins:     cfg.CORSAllowOrigins,
        AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
        AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
        AllowCredentials: true,
    }))
    app.Use(helmet.New())
    app.Use(limiter.New(limiter.Config{
        Max:        cfg.RateLimitMax,
        Expiration: cfg.RateLimitExpiration,
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.IP()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
                "error": "Rate limit exceeded",
            })
        },
    }))

    // Initialize layers
    repos := repository.NewRepositories(db)
    services := service.NewServices(repos, cfg)
    handlers := handler.NewHandlers(services)

    // Setup routes
    router.Setup(app, handlers, services.Auth)

    // Graceful shutdown
    go func() {
        if err := app.Listen(":" + cfg.Port); err != nil {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    // Shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := app.ShutdownWithContext(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    // Close database
    sqlDB, _ := db.DB()
    sqlDB.Close()

    log.Println("Server exited")
}

// Custom error handler
func customErrorHandler(c *fiber.Ctx, err error) error {
    code := fiber.StatusInternalServerError
    message := "Internal Server Error"

    if e, ok := err.(*fiber.Error); ok {
        code = e.Code
        message = e.Message
    }

    return c.Status(code).JSON(fiber.Map{
        "error":   message,
        "code":    code,
        "request": c.Locals("requestid"),
    })
}
```

### Configuration

```go
// internal/config/config.go
package config

import (
    "os"
    "strconv"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type Config struct {
    // App
    AppName     string
    Environment string
    Port        string
    Debug       bool

    // Server
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    IdleTimeout  time.Duration
    BodyLimit    int
    Prefork      bool

    // Database
    DatabaseURL string

    // JWT
    JWTSecret          string
    JWTExpiration      time.Duration
    JWTRefreshDuration time.Duration

    // CORS
    CORSAllowOrigins string

    // Rate limiting
    RateLimitMax        int
    RateLimitExpiration time.Duration
}

func Load() (*Config, error) {
    return &Config{
        // App
        AppName:     getEnv("APP_NAME", "MyAPI"),
        Environment: getEnv("ENVIRONMENT", "development"),
        Port:        getEnv("PORT", "3000"),
        Debug:       getEnvBool("DEBUG", true),

        // Server
        ReadTimeout:  getEnvDuration("READ_TIMEOUT", 10*time.Second),
        WriteTimeout: getEnvDuration("WRITE_TIMEOUT", 10*time.Second),
        IdleTimeout:  getEnvDuration("IDLE_TIMEOUT", 120*time.Second),
        BodyLimit:    getEnvInt("BODY_LIMIT", 4*1024*1024), // 4MB
        Prefork:      getEnvBool("PREFORK", false),

        // Database
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/myapi?sslmode=disable"),

        // JWT
        JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
        JWTExpiration:      getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
        JWTRefreshDuration: getEnvDuration("JWT_REFRESH_DURATION", 7*24*time.Hour),

        // CORS
        CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "*"),

        // Rate limiting
        RateLimitMax:        getEnvInt("RATE_LIMIT_MAX", 100),
        RateLimitExpiration: getEnvDuration("RATE_LIMIT_EXPIRATION", time.Minute),
    }, nil
}

func NewDatabase(cfg *Config) (*gorm.DB, error) {
    logLevel := logger.Silent
    if cfg.Debug {
        logLevel = logger.Info
    }

    db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
        Logger: logger.Default.LogMode(logLevel),
    })
    if err != nil {
        return nil, err
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }

    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(25)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)

    return db, nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if i, err := strconv.Atoi(value); err == nil {
            return i
        }
    }
    return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        if b, err := strconv.ParseBool(value); err == nil {
            return b
        }
    }
    return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if d, err := time.ParseDuration(value); err == nil {
            return d
        }
    }
    return defaultValue
}
```

### Models and DTOs

```go
// internal/model/user.go
package model

import (
    "time"

    "gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    Email     string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
    Password  string         `gorm:"not null" json:"-"`
    Name      string         `gorm:"not null;size:100" json:"name"`
    Role      string         `gorm:"not null;default:user" json:"role"`
    Active    bool           `gorm:"not null;default:true" json:"active"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
    return "users"
}
```

```go
// internal/model/dto.go
package model

// Request DTOs
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email,max=255"`
    Password string `json:"password" validate:"required,min=8,max=72"`
    Name     string `json:"name" validate:"required,min=2,max=100"`
}

type UpdateUserRequest struct {
    Email string `json:"email" validate:"omitempty,email,max=255"`
    Name  string `json:"name" validate:"omitempty,min=2,max=100"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

// Response DTOs
type UserResponse struct {
    ID        uint   `json:"id"`
    Email     string `json:"email"`
    Name      string `json:"name"`
    Role      string `json:"role"`
    Active    bool   `json:"active"`
    CreatedAt string `json:"created_at"`
}

type AuthResponse struct {
    AccessToken  string       `json:"access_token"`
    RefreshToken string       `json:"refresh_token,omitempty"`
    ExpiresIn    int64        `json:"expires_in"`
    User         UserResponse `json:"user"`
}

type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    TotalItems int64       `json:"total_items"`
    TotalPages int         `json:"total_pages"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
    return UserResponse{
        ID:        u.ID,
        Email:     u.Email,
        Name:      u.Name,
        Role:      u.Role,
        Active:    u.Active,
        CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
    }
}
```

### Validator

```go
// pkg/validator/validator.go
package validator

import (
    "reflect"
    "strings"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
)

type CustomValidator struct {
    validator *validator.Validate
}

func New() *CustomValidator {
    v := validator.New()

    // Use JSON tag names in errors
    v.RegisterTagNameFunc(func(fld reflect.StructField) string {
        name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
        if name == "-" {
            return ""
        }
        return name
    })

    return &CustomValidator{validator: v}
}

func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

// ValidationError represents validation error details
type ValidationError struct {
    Field   string `json:"field"`
    Tag     string `json:"tag"`
    Message string `json:"message"`
}

// ParseValidationErrors extracts validation errors
func ParseValidationErrors(err error) []ValidationError {
    var errors []ValidationError

    if validationErrs, ok := err.(validator.ValidationErrors); ok {
        for _, e := range validationErrs {
            errors = append(errors, ValidationError{
                Field:   e.Field(),
                Tag:     e.Tag(),
                Message: getErrorMessage(e),
            })
        }
    }

    return errors
}

func getErrorMessage(e validator.FieldError) string {
    switch e.Tag() {
    case "required":
        return e.Field() + " is required"
    case "email":
        return e.Field() + " must be a valid email"
    case "min":
        return e.Field() + " must be at least " + e.Param() + " characters"
    case "max":
        return e.Field() + " must be at most " + e.Param() + " characters"
    default:
        return e.Field() + " is invalid"
    }
}

// ValidateBody parses and validates request body
func ValidateBody[T any](c *fiber.Ctx, v *CustomValidator) (*T, error) {
    var body T
    if err := c.BodyParser(&body); err != nil {
        return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
    }

    if err := v.Validate(&body); err != nil {
        return nil, err
    }

    return &body, nil
}
```

### Response Helpers

```go
// pkg/response/response.go
package response

import (
    "github.com/gofiber/fiber/v2"
    "myapi/pkg/validator"
)

// Success returns a success response
func Success(c *fiber.Ctx, data interface{}) error {
    return c.JSON(fiber.Map{
        "success": true,
        "data":    data,
    })
}

// Created returns a 201 created response
func Created(c *fiber.Ctx, data interface{}) error {
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "data":    data,
    })
}

// NoContent returns a 204 no content response
func NoContent(c *fiber.Ctx) error {
    return c.SendStatus(fiber.StatusNoContent)
}

// Error returns an error response
func Error(c *fiber.Ctx, code int, message string) error {
    return c.Status(code).JSON(fiber.Map{
        "success": false,
        "error":   message,
    })
}

// ValidationError returns a validation error response
func ValidationError(c *fiber.Ctx, err error) error {
    errors := validator.ParseValidationErrors(err)
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "success": false,
        "error":   "Validation failed",
        "details": errors,
    })
}

// Paginated returns a paginated response
func Paginated(c *fiber.Ctx, data interface{}, page, pageSize int, total int64) error {
    totalPages := int(total) / pageSize
    if int(total)%pageSize > 0 {
        totalPages++
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    data,
        "meta": fiber.Map{
            "page":        page,
            "page_size":   pageSize,
            "total_items": total,
            "total_pages": totalPages,
        },
    })
}
```

### Repository

```go
// internal/repository/repository.go
package repository

import "gorm.io/gorm"

type Repositories struct {
    User UserRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
    return &Repositories{
        User: NewUserRepository(db),
    }
}
```

```go
// internal/repository/user_repository.go
package repository

import (
    "context"

    "gorm.io/gorm"

    "myapi/internal/model"
)

type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    FindByID(ctx context.Context, id uint) (*model.User, error)
    FindByEmail(ctx context.Context, email string) (*model.User, error)
    FindAll(ctx context.Context, page, pageSize int) ([]model.User, int64, error)
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, id uint) error
    ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
    var user model.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    var user model.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
    var users []model.User
    var total int64

    r.db.WithContext(ctx).Model(&model.User{}).Count(&total)

    offset := (page - 1) * pageSize
    err := r.db.WithContext(ctx).
        Order("created_at DESC").
        Offset(offset).
        Limit(pageSize).
        Find(&users).Error

    return users, total, err
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
    return count > 0, err
}
```

### Service Layer

```go
// internal/service/service.go
package service

import (
    "myapi/internal/config"
    "myapi/internal/repository"
)

type Services struct {
    User UserService
    Auth AuthService
}

func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
    authService := NewAuthService(cfg)

    return &Services{
        User: NewUserService(repos.User, authService),
        Auth: authService,
    }
}
```

```go
// internal/service/auth_service.go
package service

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"

    "myapi/internal/config"
    "myapi/internal/model"
)

var (
    ErrInvalidToken     = errors.New("invalid token")
    ErrExpiredToken     = errors.New("token has expired")
    ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService interface {
    GenerateTokens(user *model.User) (*model.AuthResponse, error)
    ValidateToken(tokenString string) (*Claims, error)
    HashPassword(password string) (string, error)
    CheckPassword(hashedPassword, password string) bool
}

type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

type authService struct {
    secret          []byte
    expiration      time.Duration
    refreshDuration time.Duration
}

func NewAuthService(cfg *config.Config) AuthService {
    return &authService{
        secret:          []byte(cfg.JWTSecret),
        expiration:      cfg.JWTExpiration,
        refreshDuration: cfg.JWTRefreshDuration,
    }
}

func (s *authService) GenerateTokens(user *model.User) (*model.AuthResponse, error) {
    now := time.Now()
    expiresAt := now.Add(s.expiration)

    claims := &Claims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(now),
            Subject:   user.Email,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    accessToken, err := token.SignedString(s.secret)
    if err != nil {
        return nil, err
    }

    // Generate refresh token with longer expiration
    refreshClaims := &Claims{
        UserID: user.ID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshDuration)),
            IssuedAt:  jwt.NewNumericDate(now),
        },
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshTokenString, err := refreshToken.SignedString(s.secret)
    if err != nil {
        return nil, err
    }

    return &model.AuthResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshTokenString,
        ExpiresIn:    int64(s.expiration.Seconds()),
        User:         user.ToResponse(),
    }, nil
}

func (s *authService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrInvalidToken
        }
        return s.secret, nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrExpiredToken
        }
        return nil, ErrInvalidToken
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, ErrInvalidToken
    }

    return claims, nil
}

func (s *authService) HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func (s *authService) CheckPassword(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}
```

```go
// internal/service/user_service.go
package service

import (
    "context"
    "errors"

    "gorm.io/gorm"

    "myapi/internal/model"
    "myapi/internal/repository"
)

var (
    ErrUserNotFound      = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user already exists")
)

type UserService interface {
    Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
    GetByID(ctx context.Context, id uint) (*model.User, error)
    GetAll(ctx context.Context, page, pageSize int) ([]model.User, int64, error)
    Update(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.User, error)
    Delete(ctx context.Context, id uint) error
    Authenticate(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error)
}

type userService struct {
    repo repository.UserRepository
    auth AuthService
}

func NewUserService(repo repository.UserRepository, auth AuthService) UserService {
    return &userService{
        repo: repo,
        auth: auth,
    }
}

func (s *userService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
    exists, err := s.repo.ExistsByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, ErrUserAlreadyExists
    }

    hashedPassword, err := s.auth.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }

    user := &model.User{
        Email:    req.Email,
        Password: hashedPassword,
        Name:     req.Name,
        Role:     "user",
        Active:   true,
    }

    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

func (s *userService) GetByID(ctx context.Context, id uint) (*model.User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    return user, nil
}

func (s *userService) GetAll(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }
    return s.repo.FindAll(ctx, page, pageSize)
}

func (s *userService) Update(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }

    if req.Email != "" && req.Email != user.Email {
        exists, err := s.repo.ExistsByEmail(ctx, req.Email)
        if err != nil {
            return nil, err
        }
        if exists {
            return nil, ErrUserAlreadyExists
        }
        user.Email = req.Email
    }

    if req.Name != "" {
        user.Name = req.Name
    }

    if err := s.repo.Update(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
    _, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrUserNotFound
        }
        return err
    }

    return s.repo.Delete(ctx, id)
}

func (s *userService) Authenticate(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
    user, err := s.repo.FindByEmail(ctx, req.Email)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrInvalidCredentials
        }
        return nil, err
    }

    if !s.auth.CheckPassword(user.Password, req.Password) {
        return nil, ErrInvalidCredentials
    }

    if !user.Active {
        return nil, errors.New("account is deactivated")
    }

    return s.auth.GenerateTokens(user)
}
```

### Handlers

```go
// internal/handler/handler.go
package handler

import (
    "myapi/internal/service"
    "myapi/pkg/validator"
)

type Handlers struct {
    User   *UserHandler
    Health *HealthHandler
    validator *validator.CustomValidator
}

func NewHandlers(services *service.Services) *Handlers {
    v := validator.New()

    return &Handlers{
        User:      NewUserHandler(services.User, v),
        Health:    NewHealthHandler(),
        validator: v,
    }
}
```

```go
// internal/handler/user_handler.go
package handler

import (
    "errors"
    "strconv"

    "github.com/gofiber/fiber/v2"

    "myapi/internal/model"
    "myapi/internal/service"
    "myapi/pkg/response"
    "myapi/pkg/validator"
)

type UserHandler struct {
    service   service.UserService
    validator *validator.CustomValidator
}

func NewUserHandler(service service.UserService, v *validator.CustomValidator) *UserHandler {
    return &UserHandler{
        service:   service,
        validator: v,
    }
}

// Create godoc
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body model.CreateUserRequest true "User data"
// @Success 201 {object} model.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /api/v1/users [post]
func (h *UserHandler) Create(c *fiber.Ctx) error {
    req, err := validator.ValidateBody[model.CreateUserRequest](c, h.validator)
    if err != nil {
        if _, ok := err.(*fiber.Error); ok {
            return err
        }
        return response.ValidationError(c, err)
    }

    user, err := h.service.Create(c.Context(), req)
    if err != nil {
        if errors.Is(err, service.ErrUserAlreadyExists) {
            return response.Error(c, fiber.StatusConflict, "User already exists")
        }
        return err
    }

    return response.Created(c, user.ToResponse())
}

// GetByID godoc
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.UserResponse
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid user ID")
    }

    user, err := h.service.GetByID(c.Context(), uint(id))
    if err != nil {
        if errors.Is(err, service.ErrUserNotFound) {
            return response.Error(c, fiber.StatusNotFound, "User not found")
        }
        return err
    }

    return response.Success(c, user.ToResponse())
}

// GetAll godoc
// @Summary List all users
// @Tags users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} model.PaginatedResponse
// @Router /api/v1/users [get]
func (h *UserHandler) GetAll(c *fiber.Ctx) error {
    page := c.QueryInt("page", 1)
    pageSize := c.QueryInt("page_size", 20)

    users, total, err := h.service.GetAll(c.Context(), page, pageSize)
    if err != nil {
        return err
    }

    responses := make([]model.UserResponse, len(users))
    for i, user := range users {
        responses[i] = user.ToResponse()
    }

    return response.Paginated(c, responses, page, pageSize, total)
}

// Update godoc
// @Summary Update user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body model.UpdateUserRequest true "User data"
// @Success 200 {object} model.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid user ID")
    }

    req, err := validator.ValidateBody[model.UpdateUserRequest](c, h.validator)
    if err != nil {
        if _, ok := err.(*fiber.Error); ok {
            return err
        }
        return response.ValidationError(c, err)
    }

    user, err := h.service.Update(c.Context(), uint(id), req)
    if err != nil {
        if errors.Is(err, service.ErrUserNotFound) {
            return response.Error(c, fiber.StatusNotFound, "User not found")
        }
        if errors.Is(err, service.ErrUserAlreadyExists) {
            return response.Error(c, fiber.StatusConflict, "Email already in use")
        }
        return err
    }

    return response.Success(c, user.ToResponse())
}

// Delete godoc
// @Summary Delete user
// @Tags users
// @Param id path int true "User ID"
// @Success 204
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, "Invalid user ID")
    }

    if err := h.service.Delete(c.Context(), uint(id)); err != nil {
        if errors.Is(err, service.ErrUserNotFound) {
            return response.Error(c, fiber.StatusNotFound, "User not found")
        }
        return err
    }

    return response.NoContent(c)
}

// Login godoc
// @Summary Authenticate user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Credentials"
// @Success 200 {object} model.AuthResponse
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
    req, err := validator.ValidateBody[model.LoginRequest](c, h.validator)
    if err != nil {
        if _, ok := err.(*fiber.Error); ok {
            return err
        }
        return response.ValidationError(c, err)
    }

    authResponse, err := h.service.Authenticate(c.Context(), req)
    if err != nil {
        if errors.Is(err, service.ErrInvalidCredentials) {
            return response.Error(c, fiber.StatusUnauthorized, "Invalid credentials")
        }
        return err
    }

    return response.Success(c, authResponse)
}

// GetProfile godoc
// @Summary Get current user profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.UserResponse
// @Router /api/v1/auth/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)

    user, err := h.service.GetByID(c.Context(), userID)
    if err != nil {
        return err
    }

    return response.Success(c, user.ToResponse())
}
```

```go
// internal/handler/health_handler.go
package handler

import (
    "github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
    return &HealthHandler{}
}

func (h *HealthHandler) Health(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status":  "healthy",
        "service": "myapi",
    })
}

func (h *HealthHandler) Ready(c *fiber.Ctx) error {
    // Check database connection, external services, etc.
    return c.JSON(fiber.Map{
        "status": "ready",
    })
}
```

### Middleware

```go
// internal/middleware/auth.go
package middleware

import (
    "strings"

    "github.com/gofiber/fiber/v2"

    "myapi/internal/service"
)

func Auth(authService service.AuthService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format")
        }

        claims, err := authService.ValidateToken(parts[1])
        if err != nil {
            if err == service.ErrExpiredToken {
                return fiber.NewError(fiber.StatusUnauthorized, "Token has expired")
            }
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
        }

        // Store user info in context
        c.Locals("userID", claims.UserID)
        c.Locals("userEmail", claims.Email)
        c.Locals("userRole", claims.Role)

        return c.Next()
    }
}

// RequireRole middleware checks if user has required role
func RequireRole(roles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userRole, ok := c.Locals("userRole").(string)
        if !ok {
            return fiber.NewError(fiber.StatusForbidden, "Access denied")
        }

        for _, role := range roles {
            if userRole == role {
                return c.Next()
            }
        }

        return fiber.NewError(fiber.StatusForbidden, "Insufficient permissions")
    }
}
```

```go
// internal/middleware/logger.go
package middleware

import (
    "log"
    "time"

    "github.com/gofiber/fiber/v2"
)

func Logger() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()

        // Process request
        err := c.Next()

        // Log request
        log.Printf(
            "[%s] %s %s %d %s",
            c.Method(),
            c.Path(),
            c.IP(),
            c.Response().StatusCode(),
            time.Since(start),
        )

        return err
    }
}
```

```go
// internal/middleware/recover.go
package middleware

import (
    "log"
    "runtime/debug"

    "github.com/gofiber/fiber/v2"
)

func Recover() fiber.Handler {
    return func(c *fiber.Ctx) error {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Panic recovered: %v\n%s", r, debug.Stack())
                c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                    "error": "Internal server error",
                })
            }
        }()

        return c.Next()
    }
}
```

### Router

```go
// internal/router/router.go
package router

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/swagger"

    "myapi/internal/handler"
    "myapi/internal/middleware"
    "myapi/internal/service"
)

func Setup(app *fiber.App, h *handler.Handlers, authService service.AuthService) {
    // Health checks
    app.Get("/health", h.Health.Health)
    app.Get("/ready", h.Health.Ready)

    // Swagger docs
    app.Get("/swagger/*", swagger.HandlerDefault)

    // API v1
    v1 := app.Group("/api/v1")

    // Auth routes (public)
    auth := v1.Group("/auth")
    auth.Post("/login", h.User.Login)
    auth.Post("/register", h.User.Create)

    // Protected auth routes
    authProtected := auth.Group("", middleware.Auth(authService))
    authProtected.Get("/profile", h.User.GetProfile)

    // User routes (protected)
    users := v1.Group("/users", middleware.Auth(authService))
    users.Get("/", h.User.GetAll)
    users.Get("/:id", h.User.GetByID)
    users.Put("/:id", h.User.Update)
    users.Delete("/:id", middleware.RequireRole("admin"), h.User.Delete)
}
```

---

## Testing

### Handler Tests

```go
// internal/handler/user_handler_test.go
package handler

import (
    "bytes"
    "encoding/json"
    "net/http/httptest"
    "testing"

    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "myapi/internal/model"
    "myapi/internal/service"
    "myapi/pkg/validator"
)

type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetAll(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
    args := m.Called(ctx, page, pageSize)
    return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) Update(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.User, error) {
    args := m.Called(ctx, id, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, id uint) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *MockUserService) Authenticate(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.AuthResponse), args.Error(1)
}

func setupTestApp(mockService *MockUserService) *fiber.App {
    app := fiber.New()
    v := validator.New()
    handler := NewUserHandler(mockService, v)

    app.Post("/users", handler.Create)
    app.Get("/users/:id", handler.GetByID)
    app.Get("/users", handler.GetAll)

    return app
}

func TestUserHandler_Create(t *testing.T) {
    mockService := new(MockUserService)
    app := setupTestApp(mockService)

    t.Run("success", func(t *testing.T) {
        user := &model.User{
            ID:    1,
            Email: "test@example.com",
            Name:  "Test User",
            Role:  "user",
        }

        mockService.On("Create", mock.Anything, mock.MatchedBy(func(req *model.CreateUserRequest) bool {
            return req.Email == "test@example.com"
        })).Return(user, nil).Once()

        body, _ := json.Marshal(map[string]string{
            "email":    "test@example.com",
            "password": "password123",
            "name":     "Test User",
        })

        req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")

        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
        mockService.AssertExpectations(t)
    })

    t.Run("validation error", func(t *testing.T) {
        body, _ := json.Marshal(map[string]string{
            "email": "invalid",
        })

        req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")

        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
    })

    t.Run("user already exists", func(t *testing.T) {
        mockService.On("Create", mock.Anything, mock.Anything).
            Return(nil, service.ErrUserAlreadyExists).Once()

        body, _ := json.Marshal(map[string]string{
            "email":    "existing@example.com",
            "password": "password123",
            "name":     "Test User",
        })

        req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")

        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
        mockService.AssertExpectations(t)
    })
}

func TestUserHandler_GetByID(t *testing.T) {
    mockService := new(MockUserService)
    app := setupTestApp(mockService)

    t.Run("success", func(t *testing.T) {
        user := &model.User{
            ID:    1,
            Email: "test@example.com",
            Name:  "Test User",
        }

        mockService.On("GetByID", mock.Anything, uint(1)).Return(user, nil).Once()

        req := httptest.NewRequest("GET", "/users/1", nil)
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusOK, resp.StatusCode)
        mockService.AssertExpectations(t)
    })

    t.Run("not found", func(t *testing.T) {
        mockService.On("GetByID", mock.Anything, uint(999)).
            Return(nil, service.ErrUserNotFound).Once()

        req := httptest.NewRequest("GET", "/users/999", nil)
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
        mockService.AssertExpectations(t)
    })
}
```

---

## Commands

```bash
# Run development server
go run cmd/api/main.go

# Build
go build -o bin/api cmd/api/main.go

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...

# Generate swagger docs (with swag)
swag init -g cmd/api/main.go

# Database migrations (with golang-migrate)
migrate -path migrations -database "$DATABASE_URL" up
migrate -path migrations -database "$DATABASE_URL" down 1

# Lint
golangci-lint run

# Build for production
CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/api cmd/api/main.go
```

---

## Dependencies

```bash
# Core
go get github.com/gofiber/fiber/v2
go get github.com/gofiber/swagger

# Database
go get gorm.io/gorm
go get gorm.io/driver/postgres

# Validation
go get github.com/go-playground/validator/v10

# JWT
go get github.com/golang-jwt/jwt/v5

# Password hashing
go get golang.org/x/crypto/bcrypt

# Testing
go get github.com/stretchr/testify

# Documentation
go get github.com/swaggo/swag/cmd/swag
```

---

## Best Practices

### Performance
- ✓ Enable Prefork for multi-core utilization
- ✓ Use fasthttp's zero-allocation patterns
- ✓ Configure appropriate timeouts
- ✓ Use connection pooling for database
- ✓ Implement caching for frequently accessed data
- ✓ Use `fiber.Ctx.Context()` for request context

### Security
- ✓ Use Helmet middleware for security headers
- ✓ Implement rate limiting
- ✓ Validate all inputs
- ✓ Use CORS middleware appropriately
- ✓ Hash passwords with bcrypt
- ✓ Use HTTPS in production

### Error Handling
- ✓ Use custom error handler
- ✓ Return consistent error responses
- ✓ Log errors with context
- ✓ Don't expose internal errors to clients

### Code Organization
- ✓ Follow clean architecture
- ✓ Use dependency injection
- ✓ Keep handlers thin
- ✓ Business logic in services
- ✓ Data access in repositories

---

## Fiber vs Gin vs Echo

| Feature | Fiber | Gin | Echo |
|---------|-------|-----|------|
| Performance | Fastest (fasthttp) | Fast (net/http) | Fast (net/http) |
| API Style | Express-like | Gin-specific | Echo-specific |
| Memory | Zero alloc | Low alloc | Low alloc |
| Learning Curve | Easy (Express devs) | Moderate | Moderate |
| Ecosystem | Growing | Mature | Mature |
| WebSockets | Built-in | Plugin | Plugin |

---

## References

- [Fiber Documentation](https://docs.gofiber.io/)
- [Fiber GitHub](https://github.com/gofiber/fiber)
- [Fasthttp](https://github.com/valyala/fasthttp)
- [Fiber Recipes](https://github.com/gofiber/recipes)
- [Fiber Awesome List](https://github.com/gofiber/awesome-fiber)
