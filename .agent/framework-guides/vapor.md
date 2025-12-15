# Vapor Framework Guide

> **Framework**: Vapor 4.x
> **Language**: Swift 5.9+
> **Type**: Server-Side Web Framework
> **Platform**: Linux, macOS

---

## Overview

Vapor is a popular server-side Swift framework for building web applications, APIs, and backend services. It provides a type-safe, expressive API with async/await support.

**Use Vapor when:**
- Building Swift-native backend services
- Need type-safe API development
- Sharing code between iOS/macOS apps and backend
- Building with async/await patterns
- Want Swift package ecosystem integration

**Consider alternatives when:**
- Team lacks Swift experience
- Need extensive middleware ecosystem
- Require specific database support not available
- Need maximum raw performance (consider C++/Rust)

---

## Project Structure

```
MyVaporApp/
├── Package.swift
├── Sources/
│   └── App/
│       ├── Controllers/
│       │   ├── UserController.swift
│       │   └── AuthController.swift
│       ├── Models/
│       │   ├── User.swift
│       │   └── Post.swift
│       ├── DTOs/
│       │   ├── UserDTO.swift
│       │   └── CreateUserRequest.swift
│       ├── Migrations/
│       │   ├── CreateUser.swift
│       │   └── CreatePost.swift
│       ├── Middleware/
│       │   ├── AuthMiddleware.swift
│       │   └── ErrorMiddleware.swift
│       ├── Services/
│       │   ├── UserService.swift
│       │   └── EmailService.swift
│       ├── Extensions/
│       │   └── Request+Extensions.swift
│       ├── configure.swift
│       ├── routes.swift
│       └── entrypoint.swift
├── Tests/
│   └── AppTests/
│       ├── UserControllerTests.swift
│       └── AuthControllerTests.swift
├── Resources/
│   └── Views/
│       └── index.leaf
├── Public/
│   ├── css/
│   └── js/
└── docker-compose.yml
```

---

## Package.swift

```swift
// swift-tools-version:5.9
import PackageDescription

let package = Package(
    name: "MyVaporApp",
    platforms: [
        .macOS(.v13)
    ],
    dependencies: [
        // Vapor
        .package(url: "https://github.com/vapor/vapor.git", from: "4.89.0"),
        // Fluent ORM
        .package(url: "https://github.com/vapor/fluent.git", from: "4.9.0"),
        // PostgreSQL driver
        .package(url: "https://github.com/vapor/fluent-postgres-driver.git", from: "2.8.0"),
        // Redis
        .package(url: "https://github.com/vapor/redis.git", from: "4.10.0"),
        // JWT authentication
        .package(url: "https://github.com/vapor/jwt.git", from: "4.2.0"),
        // Leaf templating
        .package(url: "https://github.com/vapor/leaf.git", from: "4.3.0"),
    ],
    targets: [
        .executableTarget(
            name: "App",
            dependencies: [
                .product(name: "Vapor", package: "vapor"),
                .product(name: "Fluent", package: "fluent"),
                .product(name: "FluentPostgresDriver", package: "fluent-postgres-driver"),
                .product(name: "Redis", package: "redis"),
                .product(name: "JWT", package: "jwt"),
                .product(name: "Leaf", package: "leaf"),
            ]
        ),
        .testTarget(
            name: "AppTests",
            dependencies: [
                .target(name: "App"),
                .product(name: "XCTVapor", package: "vapor"),
            ]
        )
    ]
)
```

---

## Application Configuration

### entrypoint.swift

```swift
import Vapor
import Logging

@main
enum Entrypoint {
    static func main() async throws {
        var env = try Environment.detect()
        try LoggingSystem.bootstrap(from: &env)

        let app = try await Application.make(env)

        do {
            try await configure(app)
            try await app.execute()
        } catch {
            app.logger.report(error: error)
            try? await app.asyncShutdown()
            throw error
        }
    }
}
```

### configure.swift

```swift
import Vapor
import Fluent
import FluentPostgresDriver
import Redis
import JWT
import Leaf

func configure(_ app: Application) async throws {
    // MARK: - Environment

    let environment = app.environment

    // MARK: - Middleware

    app.middleware.use(FileMiddleware(publicDirectory: app.directory.publicDirectory))
    app.middleware.use(ErrorMiddleware.default(environment: environment))
    app.middleware.use(CORSMiddleware())

    // MARK: - Database

    if let databaseURL = Environment.get("DATABASE_URL") {
        try app.databases.use(.postgres(url: databaseURL), as: .psql)
    } else {
        app.databases.use(
            .postgres(
                hostname: Environment.get("DB_HOST") ?? "localhost",
                port: Environment.get("DB_PORT").flatMap(Int.init) ?? 5432,
                username: Environment.get("DB_USER") ?? "vapor",
                password: Environment.get("DB_PASSWORD") ?? "vapor",
                database: Environment.get("DB_NAME") ?? "vapor_dev"
            ),
            as: .psql
        )
    }

    // MARK: - Migrations

    app.migrations.add(CreateUser())
    app.migrations.add(CreatePost())
    app.migrations.add(CreateUserToken())

    if environment == .development {
        try await app.autoMigrate()
    }

    // MARK: - Redis

    if let redisURL = Environment.get("REDIS_URL") {
        app.redis.configuration = try RedisConfiguration(url: redisURL)
    } else {
        app.redis.configuration = try RedisConfiguration(
            hostname: Environment.get("REDIS_HOST") ?? "localhost",
            port: Environment.get("REDIS_PORT").flatMap(Int.init) ?? 6379
        )
    }

    // MARK: - JWT

    guard let jwtSecret = Environment.get("JWT_SECRET") else {
        fatalError("JWT_SECRET environment variable not set")
    }

    await app.jwt.keys.add(hmac: HMACKey(from: jwtSecret), digestAlgorithm: .sha256)

    // MARK: - Leaf

    app.views.use(.leaf)

    // MARK: - Routes

    try routes(app)
}
```

### routes.swift

```swift
import Vapor

func routes(_ app: Application) throws {
    // Health check
    app.get("health") { req -> HTTPStatus in
        return .ok
    }

    // API routes
    let api = app.grouped("api", "v1")

    // Public routes
    try api.register(collection: AuthController())

    // Protected routes
    let protected = api.grouped(JWTAuthMiddleware())
    try protected.register(collection: UserController())
    try protected.register(collection: PostController())

    // Web routes
    app.get { req async throws -> View in
        return try await req.view.render("index", ["title": "Welcome"])
    }
}
```

---

## Models with Fluent

### User Model

```swift
import Fluent
import Vapor

final class User: Model, Content, @unchecked Sendable {
    static let schema = "users"

    @ID(key: .id)
    var id: UUID?

    @Field(key: "email")
    var email: String

    @Field(key: "password_hash")
    var passwordHash: String

    @Field(key: "name")
    var name: String

    @Enum(key: "role")
    var role: Role

    @Field(key: "is_active")
    var isActive: Bool

    @Timestamp(key: "created_at", on: .create)
    var createdAt: Date?

    @Timestamp(key: "updated_at", on: .update)
    var updatedAt: Date?

    @Children(for: \.$user)
    var posts: [Post]

    @Children(for: \.$user)
    var tokens: [UserToken]

    init() {}

    init(
        id: UUID? = nil,
        email: String,
        passwordHash: String,
        name: String,
        role: Role = .user,
        isActive: Bool = true
    ) {
        self.id = id
        self.email = email
        self.passwordHash = passwordHash
        self.name = name
        self.role = role
        self.isActive = isActive
    }
}

// MARK: - Role Enum

extension User {
    enum Role: String, Codable, CaseIterable {
        case admin
        case user
        case guest
    }
}

// MARK: - Authentication

extension User: ModelAuthenticatable {
    static let usernameKey = \User.$email
    static let passwordHashKey = \User.$passwordHash

    func verify(password: String) throws -> Bool {
        try Bcrypt.verify(password, created: self.passwordHash)
    }
}

// MARK: - Convenience Methods

extension User {
    static func create(
        email: String,
        password: String,
        name: String,
        role: Role = .user,
        on database: Database
    ) async throws -> User {
        let passwordHash = try Bcrypt.hash(password)
        let user = User(email: email, passwordHash: passwordHash, name: name, role: role)
        try await user.save(on: database)
        return user
    }

    func generateToken() throws -> UserToken {
        try UserToken(
            value: [UInt8].random(count: 32).base64,
            userID: self.requireID()
        )
    }
}
```

### Post Model

```swift
import Fluent
import Vapor

final class Post: Model, Content, @unchecked Sendable {
    static let schema = "posts"

    @ID(key: .id)
    var id: UUID?

    @Field(key: "title")
    var title: String

    @Field(key: "content")
    var content: String

    @Field(key: "slug")
    var slug: String

    @Enum(key: "status")
    var status: Status

    @Parent(key: "user_id")
    var user: User

    @Timestamp(key: "published_at", on: .none)
    var publishedAt: Date?

    @Timestamp(key: "created_at", on: .create)
    var createdAt: Date?

    @Timestamp(key: "updated_at", on: .update)
    var updatedAt: Date?

    init() {}

    init(
        id: UUID? = nil,
        title: String,
        content: String,
        slug: String,
        status: Status = .draft,
        userID: User.IDValue
    ) {
        self.id = id
        self.title = title
        self.content = content
        self.slug = slug
        self.status = status
        self.$user.id = userID
    }

    enum Status: String, Codable, CaseIterable {
        case draft
        case published
        case archived
    }
}
```

### UserToken Model

```swift
import Fluent
import Vapor

final class UserToken: Model, Content, @unchecked Sendable {
    static let schema = "user_tokens"

    @ID(key: .id)
    var id: UUID?

    @Field(key: "value")
    var value: String

    @Parent(key: "user_id")
    var user: User

    @Timestamp(key: "created_at", on: .create)
    var createdAt: Date?

    @Timestamp(key: "expires_at", on: .none)
    var expiresAt: Date?

    init() {}

    init(id: UUID? = nil, value: String, userID: User.IDValue) {
        self.id = id
        self.value = value
        self.$user.id = userID
        self.expiresAt = Date().addingTimeInterval(60 * 60 * 24 * 7) // 7 days
    }
}

extension UserToken: ModelTokenAuthenticatable {
    static let valueKey = \UserToken.$value
    static let userKey = \UserToken.$user

    var isValid: Bool {
        guard let expiresAt = expiresAt else { return true }
        return expiresAt > Date()
    }
}
```

---

## Migrations

### Create User Migration

```swift
import Fluent

struct CreateUser: AsyncMigration {
    func prepare(on database: Database) async throws {
        let role = try await database.enum("user_role")
            .case("admin")
            .case("user")
            .case("guest")
            .create()

        try await database.schema("users")
            .id()
            .field("email", .string, .required)
            .field("password_hash", .string, .required)
            .field("name", .string, .required)
            .field("role", role, .required)
            .field("is_active", .bool, .required, .custom("DEFAULT true"))
            .field("created_at", .datetime)
            .field("updated_at", .datetime)
            .unique(on: "email")
            .create()
    }

    func revert(on database: Database) async throws {
        try await database.schema("users").delete()
        try await database.enum("user_role").delete()
    }
}
```

### Create Post Migration

```swift
import Fluent

struct CreatePost: AsyncMigration {
    func prepare(on database: Database) async throws {
        let status = try await database.enum("post_status")
            .case("draft")
            .case("published")
            .case("archived")
            .create()

        try await database.schema("posts")
            .id()
            .field("title", .string, .required)
            .field("content", .string, .required)
            .field("slug", .string, .required)
            .field("status", status, .required)
            .field("user_id", .uuid, .required, .references("users", "id", onDelete: .cascade))
            .field("published_at", .datetime)
            .field("created_at", .datetime)
            .field("updated_at", .datetime)
            .unique(on: "slug")
            .create()

        // Create index
        try await database.schema("posts")
            .index(on: "user_id")
            .update()
    }

    func revert(on database: Database) async throws {
        try await database.schema("posts").delete()
        try await database.enum("post_status").delete()
    }
}
```

---

## DTOs (Data Transfer Objects)

### User DTOs

```swift
import Vapor

// MARK: - Request DTOs

struct CreateUserRequest: Content, Validatable {
    let email: String
    let password: String
    let name: String

    static func validations(_ validations: inout Validations) {
        validations.add("email", as: String.self, is: .email)
        validations.add("password", as: String.self, is: .count(8...))
        validations.add("name", as: String.self, is: !.empty)
    }
}

struct UpdateUserRequest: Content, Validatable {
    let name: String?
    let email: String?

    static func validations(_ validations: inout Validations) {
        validations.add("email", as: String?.self, is: .nil || .email)
        validations.add("name", as: String?.self, is: .nil || !.empty)
    }
}

struct LoginRequest: Content, Validatable {
    let email: String
    let password: String

    static func validations(_ validations: inout Validations) {
        validations.add("email", as: String.self, is: .email)
        validations.add("password", as: String.self, is: !.empty)
    }
}

// MARK: - Response DTOs

struct UserResponse: Content {
    let id: UUID
    let email: String
    let name: String
    let role: User.Role
    let createdAt: Date?

    init(user: User) throws {
        self.id = try user.requireID()
        self.email = user.email
        self.name = user.name
        self.role = user.role
        self.createdAt = user.createdAt
    }
}

struct TokenResponse: Content {
    let accessToken: String
    let tokenType: String
    let expiresIn: Int

    init(token: String, expiresIn: Int = 3600) {
        self.accessToken = token
        self.tokenType = "Bearer"
        self.expiresIn = expiresIn
    }
}

struct PaginatedResponse<T: Content>: Content {
    let items: [T]
    let metadata: PageMetadata
}

struct PageMetadata: Content {
    let page: Int
    let perPage: Int
    let total: Int
    let totalPages: Int
}
```

---

## Controllers

### User Controller

```swift
import Vapor
import Fluent

struct UserController: RouteCollection {
    func boot(routes: RoutesBuilder) throws {
        let users = routes.grouped("users")

        users.get(use: index)
        users.get(":userID", use: show)
        users.put(":userID", use: update)
        users.delete(":userID", use: delete)

        // Current user routes
        users.get("me", use: me)
        users.put("me", use: updateMe)
    }

    // MARK: - Handlers

    /// GET /users
    @Sendable
    func index(req: Request) async throws -> PaginatedResponse<UserResponse> {
        let page = try req.query.decode(PageRequest.self)

        let users = try await User.query(on: req.db)
            .filter(\.$isActive == true)
            .sort(\.$createdAt, .descending)
            .paginate(PageRequest(page: page.page, per: page.per))

        let items = try users.items.map { try UserResponse(user: $0) }

        return PaginatedResponse(
            items: items,
            metadata: PageMetadata(
                page: page.page,
                perPage: page.per,
                total: users.metadata.total,
                totalPages: users.metadata.pageCount
            )
        )
    }

    /// GET /users/:userID
    @Sendable
    func show(req: Request) async throws -> UserResponse {
        guard let user = try await User.find(req.parameters.get("userID"), on: req.db) else {
            throw Abort(.notFound, reason: "User not found")
        }
        return try UserResponse(user: user)
    }

    /// PUT /users/:userID
    @Sendable
    func update(req: Request) async throws -> UserResponse {
        // Check admin permission
        let currentUser = try req.auth.require(User.self)
        guard currentUser.role == .admin else {
            throw Abort(.forbidden, reason: "Admin access required")
        }

        guard let user = try await User.find(req.parameters.get("userID"), on: req.db) else {
            throw Abort(.notFound, reason: "User not found")
        }

        try UpdateUserRequest.validate(content: req)
        let updateRequest = try req.content.decode(UpdateUserRequest.self)

        if let name = updateRequest.name {
            user.name = name
        }
        if let email = updateRequest.email {
            // Check email uniqueness
            if let existing = try await User.query(on: req.db)
                .filter(\.$email == email)
                .filter(\.$id != user.requireID())
                .first() {
                throw Abort(.conflict, reason: "Email already in use")
            }
            user.email = email
        }

        try await user.save(on: req.db)
        return try UserResponse(user: user)
    }

    /// DELETE /users/:userID
    @Sendable
    func delete(req: Request) async throws -> HTTPStatus {
        let currentUser = try req.auth.require(User.self)
        guard currentUser.role == .admin else {
            throw Abort(.forbidden, reason: "Admin access required")
        }

        guard let user = try await User.find(req.parameters.get("userID"), on: req.db) else {
            throw Abort(.notFound, reason: "User not found")
        }

        // Soft delete
        user.isActive = false
        try await user.save(on: req.db)

        return .noContent
    }

    /// GET /users/me
    @Sendable
    func me(req: Request) async throws -> UserResponse {
        let user = try req.auth.require(User.self)
        return try UserResponse(user: user)
    }

    /// PUT /users/me
    @Sendable
    func updateMe(req: Request) async throws -> UserResponse {
        let user = try req.auth.require(User.self)

        try UpdateUserRequest.validate(content: req)
        let updateRequest = try req.content.decode(UpdateUserRequest.self)

        if let name = updateRequest.name {
            user.name = name
        }
        if let email = updateRequest.email {
            if let existing = try await User.query(on: req.db)
                .filter(\.$email == email)
                .filter(\.$id != user.requireID())
                .first() {
                throw Abort(.conflict, reason: "Email already in use")
            }
            user.email = email
        }

        try await user.save(on: req.db)
        return try UserResponse(user: user)
    }
}

// Page request helper
struct PageRequest: Content {
    var page: Int
    var per: Int

    init(page: Int = 1, per: Int = 20) {
        self.page = max(1, page)
        self.per = min(100, max(1, per))
    }
}
```

### Auth Controller

```swift
import Vapor
import Fluent

struct AuthController: RouteCollection {
    func boot(routes: RoutesBuilder) throws {
        let auth = routes.grouped("auth")

        auth.post("register", use: register)
        auth.post("login", use: login)

        // Protected routes
        let protected = auth.grouped(JWTAuthMiddleware())
        protected.post("logout", use: logout)
        protected.post("refresh", use: refresh)
    }

    // MARK: - Handlers

    /// POST /auth/register
    @Sendable
    func register(req: Request) async throws -> UserResponse {
        try CreateUserRequest.validate(content: req)
        let createRequest = try req.content.decode(CreateUserRequest.self)

        // Check if email already exists
        if let _ = try await User.query(on: req.db)
            .filter(\.$email == createRequest.email)
            .first() {
            throw Abort(.conflict, reason: "Email already registered")
        }

        let user = try await User.create(
            email: createRequest.email,
            password: createRequest.password,
            name: createRequest.name,
            on: req.db
        )

        return try UserResponse(user: user)
    }

    /// POST /auth/login
    @Sendable
    func login(req: Request) async throws -> TokenResponse {
        try LoginRequest.validate(content: req)
        let loginRequest = try req.content.decode(LoginRequest.self)

        guard let user = try await User.query(on: req.db)
            .filter(\.$email == loginRequest.email)
            .filter(\.$isActive == true)
            .first() else {
            throw Abort(.unauthorized, reason: "Invalid credentials")
        }

        guard try user.verify(password: loginRequest.password) else {
            throw Abort(.unauthorized, reason: "Invalid credentials")
        }

        // Generate JWT token
        let payload = UserPayload(
            subject: .init(value: try user.requireID().uuidString),
            expiration: .init(value: Date().addingTimeInterval(3600))
        )

        let token = try await req.jwt.sign(payload)

        return TokenResponse(token: token, expiresIn: 3600)
    }

    /// POST /auth/logout
    @Sendable
    func logout(req: Request) async throws -> HTTPStatus {
        // For JWT, logout is typically client-side
        // Optionally blacklist token in Redis
        let user = try req.auth.require(User.self)

        // Delete all tokens for user (if using database tokens)
        try await UserToken.query(on: req.db)
            .filter(\.$user.$id == user.requireID())
            .delete()

        return .noContent
    }

    /// POST /auth/refresh
    @Sendable
    func refresh(req: Request) async throws -> TokenResponse {
        let user = try req.auth.require(User.self)

        let payload = UserPayload(
            subject: .init(value: try user.requireID().uuidString),
            expiration: .init(value: Date().addingTimeInterval(3600))
        )

        let token = try await req.jwt.sign(payload)

        return TokenResponse(token: token, expiresIn: 3600)
    }
}
```

---

## JWT Authentication

### JWT Payload

```swift
import JWT
import Vapor

struct UserPayload: JWTPayload {
    var subject: SubjectClaim
    var expiration: ExpirationClaim
    var isAdmin: Bool?

    func verify(using algorithm: some JWTAlgorithm) throws {
        try expiration.verifyNotExpired()
    }
}
```

### JWT Auth Middleware

```swift
import Vapor
import JWT

struct JWTAuthMiddleware: AsyncMiddleware {
    func respond(to request: Request, chainingTo next: any AsyncResponder) async throws -> Response {
        // Extract token from Authorization header
        guard let token = request.headers.bearerAuthorization?.token else {
            throw Abort(.unauthorized, reason: "Missing authorization token")
        }

        do {
            // Verify and decode JWT
            let payload = try await request.jwt.verify(token, as: UserPayload.self)

            // Get user ID from payload
            guard let userID = UUID(payload.subject.value) else {
                throw Abort(.unauthorized, reason: "Invalid token payload")
            }

            // Load user from database
            guard let user = try await User.find(userID, on: request.db),
                  user.isActive else {
                throw Abort(.unauthorized, reason: "User not found or inactive")
            }

            // Authenticate request
            request.auth.login(user)

            return try await next.respond(to: request)
        } catch let error as JWTError {
            throw Abort(.unauthorized, reason: "Invalid token: \(error.localizedDescription)")
        }
    }
}
```

---

## Error Handling

### Custom Error Middleware

```swift
import Vapor

struct AppErrorMiddleware: AsyncMiddleware {
    func respond(to request: Request, chainingTo next: any AsyncResponder) async throws -> Response {
        do {
            return try await next.respond(to: request)
        } catch let abort as AbortError {
            return try await handleAbortError(abort, for: request)
        } catch let error as ValidationError {
            return try await handleValidationError(error, for: request)
        } catch {
            return try await handleUnknownError(error, for: request)
        }
    }

    private func handleAbortError(_ error: AbortError, for request: Request) async throws -> Response {
        let response = ErrorResponse(
            error: true,
            reason: error.reason,
            code: error.status.code
        )

        return try await response.encodeResponse(status: error.status, for: request)
    }

    private func handleValidationError(_ error: ValidationError, for request: Request) async throws -> Response {
        let response = ErrorResponse(
            error: true,
            reason: "Validation failed",
            code: 400,
            details: error.description
        )

        return try await response.encodeResponse(status: .badRequest, for: request)
    }

    private func handleUnknownError(_ error: Error, for request: Request) async throws -> Response {
        request.logger.error("Unexpected error: \(error)")

        let response = ErrorResponse(
            error: true,
            reason: request.application.environment.isRelease
                ? "An internal error occurred"
                : error.localizedDescription,
            code: 500
        )

        return try await response.encodeResponse(status: .internalServerError, for: request)
    }
}

struct ErrorResponse: Content {
    let error: Bool
    let reason: String
    let code: UInt
    let details: String?

    init(error: Bool, reason: String, code: UInt, details: String? = nil) {
        self.error = error
        self.reason = reason
        self.code = code
        self.details = details
    }
}
```

---

## Services

### User Service

```swift
import Vapor
import Fluent

protocol UserServiceProtocol: Sendable {
    func findByID(_ id: UUID, on db: Database) async throws -> User?
    func findByEmail(_ email: String, on db: Database) async throws -> User?
    func create(_ request: CreateUserRequest, on db: Database) async throws -> User
    func update(_ user: User, with request: UpdateUserRequest, on db: Database) async throws -> User
    func delete(_ user: User, on db: Database) async throws
}

struct UserService: UserServiceProtocol {
    func findByID(_ id: UUID, on db: Database) async throws -> User? {
        try await User.find(id, on: db)
    }

    func findByEmail(_ email: String, on db: Database) async throws -> User? {
        try await User.query(on: db)
            .filter(\.$email == email)
            .first()
    }

    func create(_ request: CreateUserRequest, on db: Database) async throws -> User {
        let passwordHash = try Bcrypt.hash(request.password)
        let user = User(
            email: request.email,
            passwordHash: passwordHash,
            name: request.name
        )
        try await user.save(on: db)
        return user
    }

    func update(_ user: User, with request: UpdateUserRequest, on db: Database) async throws -> User {
        if let name = request.name {
            user.name = name
        }
        if let email = request.email {
            user.email = email
        }
        try await user.save(on: db)
        return user
    }

    func delete(_ user: User, on db: Database) async throws {
        user.isActive = false
        try await user.save(on: db)
    }
}

// Register in configure.swift
extension Application {
    var userService: UserServiceProtocol {
        UserService()
    }
}

extension Request {
    var userService: UserServiceProtocol {
        application.userService
    }
}
```

---

## Testing

### Controller Tests

```swift
import XCTVapor
@testable import App

final class UserControllerTests: XCTestCase {
    var app: Application!

    override func setUp() async throws {
        app = try await Application.make(.testing)
        try await configure(app)
        try await app.autoMigrate()
    }

    override func tearDown() async throws {
        try await app.autoRevert()
        try await app.asyncShutdown()
        app = nil
    }

    func testRegisterUser() async throws {
        let createRequest = CreateUserRequest(
            email: "test@example.com",
            password: "password123",
            name: "Test User"
        )

        try await app.test(.POST, "api/v1/auth/register") { req in
            try req.content.encode(createRequest)
        } afterResponse: { res in
            XCTAssertEqual(res.status, .ok)

            let user = try res.content.decode(UserResponse.self)
            XCTAssertEqual(user.email, "test@example.com")
            XCTAssertEqual(user.name, "Test User")
        }
    }

    func testLoginUser() async throws {
        // Create user first
        let user = try await User.create(
            email: "login@example.com",
            password: "password123",
            name: "Login User",
            on: app.db
        )

        let loginRequest = LoginRequest(
            email: "login@example.com",
            password: "password123"
        )

        try await app.test(.POST, "api/v1/auth/login") { req in
            try req.content.encode(loginRequest)
        } afterResponse: { res in
            XCTAssertEqual(res.status, .ok)

            let token = try res.content.decode(TokenResponse.self)
            XCTAssertFalse(token.accessToken.isEmpty)
            XCTAssertEqual(token.tokenType, "Bearer")
        }
    }

    func testGetUsersRequiresAuth() async throws {
        try await app.test(.GET, "api/v1/users") { res in
            XCTAssertEqual(res.status, .unauthorized)
        }
    }

    func testGetUsersWithAuth() async throws {
        // Create user and get token
        let user = try await User.create(
            email: "auth@example.com",
            password: "password123",
            name: "Auth User",
            on: app.db
        )

        let payload = UserPayload(
            subject: .init(value: try user.requireID().uuidString),
            expiration: .init(value: Date().addingTimeInterval(3600))
        )
        let token = try await app.jwt.keys.sign(payload)

        try await app.test(.GET, "api/v1/users") { req in
            req.headers.bearerAuthorization = BearerAuthorization(token: token)
        } afterResponse: { res in
            XCTAssertEqual(res.status, .ok)

            let response = try res.content.decode(PaginatedResponse<UserResponse>.self)
            XCTAssertFalse(response.items.isEmpty)
        }
    }
}
```

---

## Docker Configuration

### Dockerfile

```dockerfile
# Build stage
FROM swift:5.9-jammy as build

WORKDIR /app

# Copy dependencies first for caching
COPY Package.swift Package.resolved ./
RUN swift package resolve

# Copy source and build
COPY . .
RUN swift build -c release --static-swift-stdlib

# Production stage
FROM ubuntu:jammy

RUN apt-get update && apt-get install -y \
    libcurl4 \
    libxml2 \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /app/.build/release/App /app/App
COPY --from=build /app/Public /app/Public
COPY --from=build /app/Resources /app/Resources

ENV ENVIRONMENT=production

EXPOSE 8080

ENTRYPOINT ["./App"]
CMD ["serve", "--env", "production", "--hostname", "0.0.0.0", "--port", "8080"]
```

### docker-compose.yml

```yaml
version: "3.8"

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://vapor:vapor@db:5432/vapor
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - db
      - redis

  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: vapor
      POSTGRES_PASSWORD: vapor
      POSTGRES_DB: vapor
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

---

## Best Practices

### Performance
- ✓ Use async/await throughout
- ✓ Implement pagination for list endpoints
- ✓ Use eager loading to prevent N+1 queries
- ✓ Cache frequently accessed data in Redis
- ✓ Use database indexes for query optimization

### Security
- ✓ Always validate input with Validatable
- ✓ Use Bcrypt for password hashing
- ✓ Implement rate limiting
- ✓ Use HTTPS in production
- ✓ Sanitize user input
- ✓ Use parameterized queries (Fluent handles this)

### Architecture
- ✓ Use DTOs for request/response
- ✓ Keep controllers thin
- ✓ Extract business logic to services
- ✓ Use dependency injection
- ✓ Write comprehensive tests

---

## References

- [Vapor Documentation](https://docs.vapor.codes/)
- [Fluent Documentation](https://docs.vapor.codes/fluent/overview/)
- [Vapor Discord](https://discord.gg/vapor)
- [API Guidelines](https://docs.vapor.codes/getting-started/folder-structure/)
- [Swift Server Workgroup](https://www.swift.org/sswg/)
