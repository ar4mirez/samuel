# Gin Framework Guide

> **Applies to**: Gin 1.9+, REST APIs, Microservices, Web Applications
> **Language Guide**: @.claude/skills/go-guide/SKILL.md

---

## Overview

Gin is a high-performance HTTP web framework written in Go featuring a martini-like API with performance up to 40x faster. It's the most popular Go web framework, ideal for building REST APIs and microservices.

**Use Gin when:**
- Building high-performance REST APIs
- You need a mature, well-documented framework
- Middleware ecosystem is important
- You want a balanced approach (not too minimal, not too heavy)

**Consider alternatives when:**
- You need maximum minimalism (use standard library)
- You want built-in WebSocket support (use Fiber)
- You prefer a different API style (use Echo)

---

## Project Structure

```
myproject/
├── cmd/
│   └── api/
│       └── main.go           # Entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration
│   ├── handler/
│   │   ├── handler.go        # Handler interface
│   │   ├── user.go           # User handlers
│   │   └── auth.go           # Auth handlers
│   ├── middleware/
│   │   ├── auth.go           # JWT middleware
│   │   ├── cors.go           # CORS middleware
│   │   └── logger.go         # Request logging
│   ├── model/
│   │   ├── user.go           # User model
│   │   └── response.go       # Response models
│   ├── repository/
│   │   ├── repository.go     # Repository interface
│   │   └── user.go           # User repository
│   ├── service/
│   │   ├── service.go        # Service interface
│   │   └── user.go           # User service
│   └── router/
│       └── router.go         # Route definitions
├── pkg/
│   ├── validator/
│   │   └── validator.go      # Custom validators
│   └── response/
│       └── response.go       # Response helpers
├── migrations/
├── .env.example
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Application Setup

### cmd/api/main.go
```go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"myproject/internal/config"
	"myproject/internal/handler"
	"myproject/internal/middleware"
	"myproject/internal/repository"
	"myproject/internal/router"
	"myproject/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := initDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize layers
	repos := repository.NewRepositories(db)
	services := service.NewServices(repos, cfg)
	handlers := handler.NewHandlers(services)

	// Setup router
	r := router.Setup(handlers, middleware.NewMiddleware(cfg))

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func initDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
```

---

## Configuration

### internal/config/config.go
```go
package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        string
	DatabaseURL string

	JWTSecret          string
	JWTAccessDuration  time.Duration
	JWTRefreshDuration time.Duration

	CORSOrigins []string
}

func Load() (*Config, error) {
	// Load .env file in development
	_ = godotenv.Load()

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/myapp?sslmode=disable"),

		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		JWTAccessDuration:  parseDuration(getEnv("JWT_ACCESS_DURATION", "15m")),
		JWTRefreshDuration: parseDuration(getEnv("JWT_REFRESH_DURATION", "168h")),

		CORSOrigins: []string{getEnv("CORS_ORIGINS", "*")},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 15 * time.Minute
	}
	return d
}
```

---

## Models

### internal/model/user.go
```go
package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	IsAdmin   bool           `json:"is_admin" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// Request DTOs
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required,min=1,max=100"`
	LastName  string `json:"last_name" binding:"required,min=1,max=100"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name" binding:"omitempty,min=1,max=100"`
	LastName  *string `json:"last_name" binding:"omitempty,min=1,max=100"`
	IsActive  *bool   `json:"is_active"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Response DTOs
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	FullName  string    `json:"full_name"`
	IsActive  bool      `json:"is_active"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		FullName:  u.FullName(),
		IsActive:  u.IsActive,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}
```

### internal/model/response.go
```go
package model

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    *PageMeta   `json:"meta"`
}

type PageMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Success: true,
		Data:    data,
	}
}

func NewErrorResponse(message string) *Response {
	return &Response{
		Success: false,
		Error:   message,
	}
}

func NewPaginatedResponse(data interface{}, page, perPage int, total int64) *PaginatedResponse {
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Success: true,
		Data:    data,
		Meta: &PageMeta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
```

---

## Repository Layer

### internal/repository/repository.go
```go
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

### internal/repository/user.go
```go
package repository

import (
	"context"

	"myproject/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetAll(ctx context.Context, page, perPage int) ([]model.User, int64, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
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

func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAll(ctx context.Context, page, perPage int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	offset := (page - 1) * perPage

	err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Offset(offset).
		Limit(perPage).
		Order("created_at DESC").
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}
```

---

## Service Layer

### internal/service/service.go
```go
package service

import (
	"myproject/internal/config"
	"myproject/internal/repository"
)

type Services struct {
	User UserService
	Auth AuthService
}

func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	authService := NewAuthService(repos.User, cfg)
	return &Services{
		User: NewUserService(repos.User),
		Auth: authService,
	}
}
```

### internal/service/user.go
```go
package service

import (
	"context"
	"errors"

	"myproject/internal/model"
	"myproject/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserService interface {
	Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetAll(ctx context.Context, page, perPage int) ([]model.User, int64, error)
	Update(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.User, error)
	Delete(ctx context.Context, id uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// Check if user exists
	existing, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, ErrUserAlreadyExists
	}

	user := &model.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetAll(ctx context.Context, page, perPage int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return s.userRepo.GetAll(ctx, page, perPage)
}

func (s *userService) Update(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	_, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.userRepo.Delete(ctx, id)
}
```

### internal/service/auth.go
```go
package service

import (
	"context"
	"errors"
	"time"

	"myproject/internal/config"
	"myproject/internal/model"
	"myproject/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(ctx context.Context, req *model.LoginRequest) (*model.TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.TokenResponse, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type Claims struct {
	UserID  uint `json:"user_id"`
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (*model.TokenResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.CheckPassword(req.Password) {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	return s.generateTokens(user)
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*model.TokenResponse, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	return s.generateTokens(user)
}

func (s *authService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (s *authService) generateTokens(user *model.User) (*model.TokenResponse, error) {
	now := time.Now()
	accessExp := now.Add(s.cfg.JWTAccessDuration)
	refreshExp := now.Add(s.cfg.JWTRefreshDuration)

	// Access token
	accessClaims := &Claims{
		UserID:  user.ID,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshClaims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &model.TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.cfg.JWTAccessDuration.Seconds()),
	}, nil
}
```

---

## Handlers

### internal/handler/handler.go
```go
package handler

import "myproject/internal/service"

type Handlers struct {
	User *UserHandler
	Auth *AuthHandler
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		User: NewUserHandler(services.User),
		Auth: NewAuthHandler(services.Auth),
	}
}
```

### internal/handler/user.go
```go
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"myproject/internal/model"
	"myproject/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUsers godoc
// @Summary Get all users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} model.PaginatedResponse
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	users, total, err := h.userService.GetAll(c.Request.Context(), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	c.JSON(http.StatusOK, model.NewPaginatedResponse(responses, page, perPage, total))
}

// GetUser godoc
// @Summary Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.Response
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid user ID"))
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, model.NewErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(user.ToResponse()))
}

// CreateUser godoc
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body model.CreateUserRequest true "User data"
// @Success 201 {object} model.Response
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
		return
	}

	user, err := h.userService.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, model.NewErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(user.ToResponse()))
}

// UpdateUser godoc
// @Summary Update user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body model.UpdateUserRequest true "User data"
// @Success 200 {object} model.Response
// @Router /users/{id} [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid user ID"))
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
		return
	}

	user, err := h.userService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, model.NewErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(user.ToResponse()))
}

// DeleteUser godoc
// @Summary Delete user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid user ID"))
		return
	}

	err = h.userService.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, model.NewErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}

// GetCurrentUser godoc
// @Summary Get current authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} model.Response
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("unauthorized"))
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(user.ToResponse()))
}
```

### internal/handler/auth.go
```go
package handler

import (
	"net/http"

	"myproject/internal/model"
	"myproject/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.TokenResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
		return
	}

	tokens, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Refresh godoc
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} model.TokenResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
		return
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, tokens)
}
```

---

## Middleware

### internal/middleware/auth.go
```go
package middleware

import (
	"net/http"
	"strings"

	"myproject/internal/model"
	"myproject/internal/service"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse("missing authorization header"))
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse("invalid authorization header"))
			return
		}

		claims, err := m.authService.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse("invalid token"))
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("is_admin", claims.IsAdmin)
		c.Next()
	}
}

func (m *Middleware) AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			c.AbortWithStatusJSON(http.StatusForbidden, model.NewErrorResponse("admin access required"))
			return
		}
		c.Next()
	}
}

type Middleware struct {
	authService service.AuthService
}

func NewMiddleware(cfg interface{}) *Middleware {
	return &Middleware{}
}

func (m *Middleware) SetAuthService(authService service.AuthService) {
	m.authService = authService
}
```

### internal/middleware/logger.go
```go
package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if query != "" {
			path = path + "?" + query
		}

		log.Printf("[GIN] %3d | %13v | %15s | %-7s %s",
			status,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
```

---

## Router

### internal/router/router.go
```go
package router

import (
	"myproject/internal/handler"
	"myproject/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(handlers *handler.Handlers, mw *middleware.Middleware) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.Auth.Login)
			auth.POST("/refresh", handlers.Auth.Refresh)
		}

		// User routes
		users := v1.Group("/users")
		{
			// Public
			users.POST("", handlers.User.CreateUser)

			// Protected
			users.Use(mw.Auth())
			users.GET("", handlers.User.GetUsers)
			users.GET("/me", handlers.User.GetCurrentUser)
			users.GET("/:id", handlers.User.GetUser)
			users.PATCH("/:id", handlers.User.UpdateUser)

			// Admin only
			users.DELETE("/:id", mw.AdminOnly(), handlers.User.DeleteUser)
		}
	}

	return r
}
```

---

## Testing

### internal/handler/user_test.go
```go
package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"myproject/internal/handler"
	"myproject/internal/model"
	"myproject/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of UserService
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

func (m *MockUserService) GetAll(ctx context.Context, page, perPage int) ([]model.User, int64, error) {
	args := m.Called(ctx, page, perPage)
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

func setupTestRouter(h *handler.UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/users", h.CreateUser)
	r.GET("/users/:id", h.GetUser)
	r.GET("/users", h.GetUsers)
	return r
}

func TestCreateUser(t *testing.T) {
	mockService := new(MockUserService)
	h := handler.NewUserHandler(mockService)
	router := setupTestRouter(h)

	t.Run("successful creation", func(t *testing.T) {
		expectedUser := &model.User{
			ID:        1,
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
		}

		mockService.On("Create", mock.Anything, mock.Anything).Return(expectedUser, nil).Once()

		body := model.CreateUserRequest{
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("duplicate email", func(t *testing.T) {
		mockService.On("Create", mock.Anything, mock.Anything).
			Return(nil, service.ErrUserAlreadyExists).Once()

		body := model.CreateUserRequest{
			Email:     "existing@example.com",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	mockService := new(MockUserService)
	h := handler.NewUserHandler(mockService)
	router := setupTestRouter(h)

	t.Run("user found", func(t *testing.T) {
		expectedUser := &model.User{
			ID:        1,
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
		}

		mockService.On("GetByID", mock.Anything, uint(1)).Return(expectedUser, nil).Once()

		req, _ := http.NewRequest("GET", "/users/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockService.On("GetByID", mock.Anything, uint(999)).
			Return(nil, service.ErrUserNotFound).Once()

		req, _ := http.NewRequest("GET", "/users/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}
```

---

## Commands Reference

```bash
# Initialize project
go mod init myproject

# Install dependencies
go mod tidy

# Run development server
go run cmd/api/main.go

# Build binary
go build -o bin/api cmd/api/main.go

# Run tests
go test ./...
go test -v -cover ./...

# Run with race detection
go test -race ./...

# Lint
golangci-lint run

# Generate Swagger docs
swag init -g cmd/api/main.go

# Database migrations (using golang-migrate)
migrate -path migrations -database "$DATABASE_URL" up
migrate -path migrations -database "$DATABASE_URL" down
```

---

## Dependencies

```go
// go.mod
module myproject

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/joho/godotenv v1.5.1
    golang.org/x/crypto v0.14.0
    gorm.io/driver/postgres v1.5.4
    gorm.io/gorm v1.25.5
)

require (
    github.com/stretchr/testify v1.8.4 // testing
    github.com/swaggo/gin-swagger v1.6.0 // swagger
    github.com/swaggo/swag v1.16.2 // swagger
)
```

---

## Best Practices

### Gin-Specific Guidelines
- ✓ Use application factory pattern for testability
- ✓ Group routes with versioning (/api/v1)
- ✓ Use middleware for cross-cutting concerns
- ✓ Use gin.Context for request-scoped data
- ✓ Use binding tags for validation
- ✓ Return consistent JSON response structure
- ✓ Use proper HTTP status codes

### Performance Guidelines
- ✓ Use gin.ReleaseMode in production
- ✓ Configure proper timeouts
- ✓ Use connection pooling for database
- ✓ Implement graceful shutdown
- ✓ Use pagination for list endpoints

---

## References

- [Gin Documentation](https://gin-gonic.com/docs/)
- [Gin GitHub](https://github.com/gin-gonic/gin)
- [GORM Documentation](https://gorm.io/docs/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
