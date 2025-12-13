---
title: Go Guide
description: Go development guardrails and best practices
---

# Go Guide

> **Applies to**: Go 1.20+, Microservices, APIs, CLIs

---

## Core Principles

1. **Simplicity**: Prefer simple, readable code over clever solutions
2. **Concurrency**: Use goroutines and channels for concurrent operations
3. **Errors Are Values**: Explicit error handling, no exceptions
4. **Composition Over Inheritance**: Interfaces and struct embedding
5. **Standard Library First**: Rich stdlib, minimize dependencies

---

## Language-Specific Guardrails

### Go Version & Setup

```
✓ Use Go 1.20+ (1.21+ for improved performance)
✓ Use Go modules (go.mod) for dependency management
✓ Run go mod tidy before committing
✓ Pin major versions in go.mod
```

### Code Style (Effective Go)

```
✓ Run gofmt before every commit (auto-format)
✓ Run go vet to catch common mistakes
✓ Run golangci-lint for comprehensive linting
✓ Use goimports for import management
✓ Package names: lowercase, no underscores (userservice not user_service)
✓ Exported names: PascalCase (UserService)
✓ Unexported names: camelCase (userService)
```

### Error Handling

```
✓ Always check errors: if err != nil { return err }
✓ Return errors, don't panic (panic only for unrecoverable errors)
✓ Wrap errors with context: fmt.Errorf("failed to fetch user: %w", err)
✓ Use custom error types for domain errors
✓ Don't ignore errors with _ unless justified with comment
```

### Concurrency

```
✓ Use context.Context for cancellation and timeouts
✓ Always set timeout for HTTP requests
✓ Use sync.WaitGroup for goroutine coordination
✓ Close channels from sender side only
✓ Use select with default to avoid blocking
```

### Interfaces

```
✓ Accept interfaces, return structs
✓ Define interfaces where they're used (not where implemented)
✓ Keep interfaces small (1-3 methods ideal)
✓ Use io.Reader, io.Writer from stdlib when applicable
```

---

## Error Handling Patterns

### Basic Pattern

```go
func GetUser(id string) (*User, error) {
    user, err := db.FindUserByID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %s: %w", id, err)
    }
    return user, nil
}
```

### Custom Errors

```go
type NotFoundError struct {
    Resource string
    ID       string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}

// Usage
func GetUser(id string) (*User, error) {
    user, err := db.FindUserByID(id)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, &NotFoundError{Resource: "user", ID: id}
    }
    if err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }
    return user, nil
}

// Checking
user, err := GetUser("123")
if err != nil {
    var notFound *NotFoundError
    if errors.As(err, &notFound) {
        // Handle not found
    }
}
```

---

## Testing

### Guardrails

```
✓ Test files: *_test.go (same package)
✓ Test functions: func TestFunctionName(t *testing.T)
✓ Table-driven tests for multiple cases
✓ Use t.Helper() in test helpers
✓ Use subtests: t.Run("subtest name", func(t *testing.T) {...})
✓ Coverage target: >80% for business logic
✓ Benchmark critical paths: func BenchmarkFunction(b *testing.B)
```

### Table-Driven Tests

```go
func TestCalculate(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
        wantErr  bool
    }{
        {
            name:     "positive numbers",
            a:        2,
            b:        3,
            expected: 5,
            wantErr:  false,
        },
        {
            name:     "negative numbers",
            a:        -2,
            b:        -3,
            expected: -5,
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Calculate(tt.a, tt.b)

            if tt.wantErr {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                return
            }

            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }

            if result != tt.expected {
                t.Errorf("got %d, want %d", result, tt.expected)
            }
        })
    }
}
```

---

## Tooling

### Essential Commands

```bash
# Format code
go fmt ./...
gofmt -s -w .           # Simplified formatting

# Vet (detect suspicious constructs)
go vet ./...

# Test
go test ./...
go test -cover ./...    # With coverage
go test -race ./...     # Race detector

# Build
go build ./cmd/api

# Mod operations
go mod tidy             # Clean up dependencies
go mod vendor           # Vendor dependencies
go mod verify           # Verify dependencies

# Linting
golangci-lint run       # Comprehensive linting
```

### Configuration

```yaml
# .golangci.yml
linters:
  enable:
    - gofmt
    - govet
    - staticcheck
    - ineffassign
    - misspell
    - gosec        # Security
    - errcheck     # Unchecked errors
    - gocyclo      # Cyclomatic complexity
    - dupl         # Code duplication

linters-settings:
  gocyclo:
    min-complexity: 10
  dupl:
    threshold: 100
```

---

## Common Pitfalls

### Don't Do This

```go
// ❌ Ignoring errors
result, _ := doSomething()

// ❌ Not using context for cancellation
func LongRunningTask() {
    time.Sleep(10 * time.Minute)
}

// ❌ Goroutine leak (no way to stop)
go func() {
    for {
        doWork()
    }
}()

// ❌ Range loop variable capture
for _, item := range items {
    go func() {
        process(item) // Wrong: captures loop variable
    }()
}
```

### Do This Instead

```go
// ✅ Proper error handling
result, err := doSomething()
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// ✅ Use context for cancellation
func LongRunningTask(ctx context.Context) error {
    select {
    case <-time.After(10 * time.Minute):
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

// ✅ Goroutine with cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go func() {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            doWork()
        }
    }
}()

// ✅ Correct loop variable capture
for _, item := range items {
    item := item // Capture loop variable
    go func() {
        process(item)
    }()
}
```

---

## HTTP Server Patterns

### Basic Server with Graceful Shutdown

```go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/health", healthHandler)
    mux.HandleFunc("/users", usersHandler)

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Graceful shutdown
    go func() {
        sigint := make(chan os.Signal, 1)
        signal.Notify(sigint, os.Interrupt)
        <-sigint

        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        if err := srv.Shutdown(ctx); err != nil {
            log.Printf("HTTP server shutdown error: %v", err)
        }
    }()

    log.Println("Starting server on :8080")
    if err := srv.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatalf("HTTP server error: %v", err)
    }
}
```

### Middleware Pattern

```go
type Middleware func(http.Handler) http.Handler

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

// Usage
handler := LoggingMiddleware(AuthMiddleware(http.HandlerFunc(usersHandler)))
```

---

## Performance Considerations

### Optimization Guardrails

```
✓ Use sync.Pool for frequently allocated objects
✓ Avoid string concatenation in loops (use strings.Builder)
✓ Use buffered channels when appropriate
✓ Profile before optimizing: go test -bench, pprof
✓ Benchmark critical paths with testing.B
```

### Example

```go
// String building
var sb strings.Builder
for _, s := range items {
    sb.WriteString(s)
}
result := sb.String()

// Object pooling
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

buf := bufferPool.Get().(*bytes.Buffer)
defer bufferPool.Put(buf)
buf.Reset()
```

---

## Security Best Practices

### Guardrails

```
✓ Use parameterized queries (prevents SQL injection)
✓ Validate all user inputs
✓ Use crypto/rand for random numbers (not math/rand)
✓ Hash passwords with bcrypt or argon2
✓ Use HTTPS (TLS) for production
✓ Run gosec to detect security issues
```

### Example

```go
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

---

## References

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Go Proverbs](https://go-proverbs.github.io/)
- [golangci-lint](https://golangci-lint.run/)
