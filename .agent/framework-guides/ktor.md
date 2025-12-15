# Ktor Framework Guide

> **Applies to**: Ktor 2.x, Kotlin 1.9+, Microservices, REST APIs, WebSockets

---

## Overview

Ktor is a lightweight, asynchronous framework built by JetBrains specifically for Kotlin. It leverages coroutines for non-blocking I/O and provides a flexible plugin system for extensibility.

**Best For**: Microservices, lightweight APIs, real-time applications, Kotlin-native projects

**Key Features**:
- Native Kotlin coroutines support
- Lightweight and modular (only include what you need)
- Type-safe routing with DSL
- Built-in testing support
- Multiple server engines (Netty, Jetty, CIO, Tomcat)
- WebSocket support out of the box
- Multiplatform client

---

## Project Structure

```
myapp/
├── src/
│   └── main/
│       ├── kotlin/
│       │   └── com/example/
│       │       ├── Application.kt           # Entry point
│       │       ├── plugins/                 # Ktor plugins configuration
│       │       │   ├── Routing.kt
│       │       │   ├── Serialization.kt
│       │       │   ├── Security.kt
│       │       │   ├── StatusPages.kt
│       │       │   └── Databases.kt
│       │       ├── routes/                  # Route definitions
│       │       │   ├── UserRoutes.kt
│       │       │   └── AuthRoutes.kt
│       │       ├── models/                  # Data models
│       │       │   ├── User.kt
│       │       │   └── Requests.kt
│       │       ├── services/                # Business logic
│       │       │   └── UserService.kt
│       │       ├── repositories/            # Data access
│       │       │   └── UserRepository.kt
│       │       └── utils/                   # Utilities
│       │           └── JwtUtils.kt
│       └── resources/
│           ├── application.conf             # HOCON configuration
│           └── logback.xml
├── src/
│   └── test/
│       └── kotlin/
│           └── com/example/
│               └── ApplicationTest.kt
├── build.gradle.kts
└── gradle.properties
```

---

## Dependencies (build.gradle.kts)

```kotlin
plugins {
    kotlin("jvm") version "1.9.22"
    kotlin("plugin.serialization") version "1.9.22"
    id("io.ktor.plugin") version "2.3.7"
}

group = "com.example"
version = "0.0.1"

application {
    mainClass.set("com.example.ApplicationKt")

    val isDevelopment: Boolean = project.ext.has("development")
    applicationDefaultJvmArgs = listOf("-Dio.ktor.development=$isDevelopment")
}

repositories {
    mavenCentral()
}

dependencies {
    // Ktor Server
    implementation("io.ktor:ktor-server-core-jvm")
    implementation("io.ktor:ktor-server-netty-jvm")
    implementation("io.ktor:ktor-server-host-common-jvm")

    // Content Negotiation & Serialization
    implementation("io.ktor:ktor-server-content-negotiation-jvm")
    implementation("io.ktor:ktor-serialization-kotlinx-json-jvm")

    // Authentication
    implementation("io.ktor:ktor-server-auth-jvm")
    implementation("io.ktor:ktor-server-auth-jwt-jvm")

    // Validation
    implementation("io.ktor:ktor-server-request-validation")

    // Status Pages (Error Handling)
    implementation("io.ktor:ktor-server-status-pages-jvm")

    // CORS
    implementation("io.ktor:ktor-server-cors-jvm")

    // Call Logging
    implementation("io.ktor:ktor-server-call-logging-jvm")

    // Database
    implementation("org.jetbrains.exposed:exposed-core:0.45.0")
    implementation("org.jetbrains.exposed:exposed-dao:0.45.0")
    implementation("org.jetbrains.exposed:exposed-jdbc:0.45.0")
    implementation("org.jetbrains.exposed:exposed-java-time:0.45.0")
    implementation("com.zaxxer:HikariCP:5.1.0")
    implementation("org.postgresql:postgresql:42.7.1")

    // Password Hashing
    implementation("at.favre.lib:bcrypt:0.10.2")

    // Logging
    implementation("ch.qos.logback:logback-classic:1.4.14")

    // Testing
    testImplementation("io.ktor:ktor-server-tests-jvm")
    testImplementation("io.ktor:ktor-client-content-negotiation")
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit:1.9.22")
    testImplementation("io.mockk:mockk:1.13.9")
    testImplementation("com.h2database:h2:2.2.224")
}

ktor {
    fatJar {
        archiveFileName.set("app.jar")
    }
}
```

---

## Application Entry Point

```kotlin
// src/main/kotlin/com/example/Application.kt
package com.example

import com.example.plugins.*
import io.ktor.server.application.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*

fun main() {
    embeddedServer(Netty, port = 8080, host = "0.0.0.0", module = Application::module)
        .start(wait = true)
}

fun Application.module() {
    configureSerialization()
    configureDatabases()
    configureSecurity()
    configureStatusPages()
    configureRouting()
}
```

---

## Configuration (application.conf)

```hocon
# src/main/resources/application.conf
ktor {
    deployment {
        port = 8080
        port = ${?PORT}
    }
    application {
        modules = [ com.example.ApplicationKt.module ]
    }
}

database {
    url = "jdbc:postgresql://localhost:5432/myapp"
    url = ${?DATABASE_URL}
    driver = "org.postgresql.Driver"
    user = "postgres"
    user = ${?DATABASE_USER}
    password = "password"
    password = ${?DATABASE_PASSWORD}
    maxPoolSize = 10
}

jwt {
    secret = "your-256-bit-secret-key-here-change-in-production"
    secret = ${?JWT_SECRET}
    issuer = "myapp"
    audience = "myapp-users"
    realm = "myapp"
    expirationMs = 3600000
}
```

---

## Plugin Configuration

### Serialization Plugin

```kotlin
// src/main/kotlin/com/example/plugins/Serialization.kt
package com.example.plugins

import io.ktor.serialization.kotlinx.json.*
import io.ktor.server.application.*
import io.ktor.server.plugins.contentnegotiation.*
import kotlinx.serialization.json.Json

fun Application.configureSerialization() {
    install(ContentNegotiation) {
        json(Json {
            prettyPrint = true
            isLenient = true
            ignoreUnknownKeys = true
            encodeDefaults = true
        })
    }
}
```

### Database Plugin (Exposed)

```kotlin
// src/main/kotlin/com/example/plugins/Databases.kt
package com.example.plugins

import com.zaxxer.hikari.HikariConfig
import com.zaxxer.hikari.HikariDataSource
import io.ktor.server.application.*
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.SchemaUtils
import org.jetbrains.exposed.sql.transactions.transaction
import com.example.repositories.Users

fun Application.configureDatabases() {
    val config = environment.config

    val hikariConfig = HikariConfig().apply {
        jdbcUrl = config.property("database.url").getString()
        driverClassName = config.property("database.driver").getString()
        username = config.property("database.user").getString()
        password = config.property("database.password").getString()
        maximumPoolSize = config.property("database.maxPoolSize").getString().toInt()
        isAutoCommit = false
        transactionIsolation = "TRANSACTION_REPEATABLE_READ"
        validate()
    }

    val dataSource = HikariDataSource(hikariConfig)
    Database.connect(dataSource)

    // Create tables
    transaction {
        SchemaUtils.create(Users)
    }
}
```

### Security Plugin (JWT)

```kotlin
// src/main/kotlin/com/example/plugins/Security.kt
package com.example.plugins

import com.auth0.jwt.JWT
import com.auth0.jwt.algorithms.Algorithm
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.auth.*
import io.ktor.server.auth.jwt.*
import io.ktor.server.response.*

fun Application.configureSecurity() {
    val config = environment.config
    val jwtSecret = config.property("jwt.secret").getString()
    val jwtIssuer = config.property("jwt.issuer").getString()
    val jwtAudience = config.property("jwt.audience").getString()
    val jwtRealm = config.property("jwt.realm").getString()

    install(Authentication) {
        jwt("auth-jwt") {
            realm = jwtRealm
            verifier(
                JWT
                    .require(Algorithm.HMAC256(jwtSecret))
                    .withAudience(jwtAudience)
                    .withIssuer(jwtIssuer)
                    .build()
            )
            validate { credential ->
                val userId = credential.payload.getClaim("userId").asString()
                val email = credential.payload.getClaim("email").asString()

                if (userId != null && email != null) {
                    JWTPrincipal(credential.payload)
                } else {
                    null
                }
            }
            challenge { _, _ ->
                call.respond(
                    HttpStatusCode.Unauthorized,
                    mapOf("error" to "Token is not valid or has expired")
                )
            }
        }
    }
}
```

### Status Pages (Error Handling)

```kotlin
// src/main/kotlin/com/example/plugins/StatusPages.kt
package com.example.plugins

import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.plugins.*
import io.ktor.server.plugins.requestvalidation.*
import io.ktor.server.plugins.statuspages.*
import io.ktor.server.response.*
import kotlinx.serialization.Serializable

@Serializable
data class ErrorResponse(
    val status: Int,
    val error: String,
    val message: String,
)

// Custom Exceptions
class NotFoundException(message: String) : RuntimeException(message)
class UnauthorizedException(message: String) : RuntimeException(message)
class ConflictException(message: String) : RuntimeException(message)
class BadRequestException(message: String) : RuntimeException(message)

fun Application.configureStatusPages() {
    install(StatusPages) {
        exception<NotFoundException> { call, cause ->
            call.respond(
                HttpStatusCode.NotFound,
                ErrorResponse(404, "NOT_FOUND", cause.message ?: "Resource not found")
            )
        }

        exception<UnauthorizedException> { call, cause ->
            call.respond(
                HttpStatusCode.Unauthorized,
                ErrorResponse(401, "UNAUTHORIZED", cause.message ?: "Unauthorized")
            )
        }

        exception<ConflictException> { call, cause ->
            call.respond(
                HttpStatusCode.Conflict,
                ErrorResponse(409, "CONFLICT", cause.message ?: "Resource conflict")
            )
        }

        exception<BadRequestException> { call, cause ->
            call.respond(
                HttpStatusCode.BadRequest,
                ErrorResponse(400, "BAD_REQUEST", cause.message ?: "Invalid request")
            )
        }

        exception<RequestValidationException> { call, cause ->
            call.respond(
                HttpStatusCode.BadRequest,
                ErrorResponse(
                    400,
                    "VALIDATION_ERROR",
                    cause.reasons.joinToString("; ")
                )
            )
        }

        exception<Throwable> { call, cause ->
            call.application.environment.log.error("Unhandled exception", cause)
            call.respond(
                HttpStatusCode.InternalServerError,
                ErrorResponse(500, "INTERNAL_ERROR", "An unexpected error occurred")
            )
        }
    }
}
```

---

## Models

### User Entity (Exposed)

```kotlin
// src/main/kotlin/com/example/models/User.kt
package com.example.models

import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.id.LongIdTable
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.javatime.datetime
import java.time.LocalDateTime

// Database Table
object Users : LongIdTable("users") {
    val email = varchar("email", 255).uniqueIndex()
    val passwordHash = varchar("password_hash", 255)
    val name = varchar("name", 255)
    val role = varchar("role", 50).default("user")
    val createdAt = datetime("created_at").default(LocalDateTime.now())
    val updatedAt = datetime("updated_at").default(LocalDateTime.now())
}

// Domain Model
data class User(
    val id: Long,
    val email: String,
    val passwordHash: String,
    val name: String,
    val role: String,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime,
) {
    companion object {
        fun fromRow(row: ResultRow): User = User(
            id = row[Users.id].value,
            email = row[Users.email],
            passwordHash = row[Users.passwordHash],
            name = row[Users.name],
            role = row[Users.role],
            createdAt = row[Users.createdAt],
            updatedAt = row[Users.updatedAt],
        )
    }
}

// Response DTO
@Serializable
data class UserResponse(
    val id: Long,
    val email: String,
    val name: String,
    val role: String,
    val createdAt: String,
) {
    companion object {
        fun from(user: User): UserResponse = UserResponse(
            id = user.id,
            email = user.email,
            name = user.name,
            role = user.role,
            createdAt = user.createdAt.toString(),
        )
    }
}
```

### Request DTOs

```kotlin
// src/main/kotlin/com/example/models/Requests.kt
package com.example.models

import kotlinx.serialization.Serializable

@Serializable
data class CreateUserRequest(
    val email: String,
    val password: String,
    val name: String,
)

@Serializable
data class UpdateUserRequest(
    val name: String? = null,
    val email: String? = null,
)

@Serializable
data class LoginRequest(
    val email: String,
    val password: String,
)

@Serializable
data class LoginResponse(
    val token: String,
    val user: UserResponse,
)
```

---

## Request Validation

```kotlin
// src/main/kotlin/com/example/plugins/Validation.kt
package com.example.plugins

import com.example.models.CreateUserRequest
import com.example.models.LoginRequest
import io.ktor.server.application.*
import io.ktor.server.plugins.requestvalidation.*

fun Application.configureValidation() {
    install(RequestValidation) {
        validate<CreateUserRequest> { request ->
            val errors = mutableListOf<String>()

            if (!request.email.contains("@")) {
                errors.add("Invalid email format")
            }
            if (request.password.length < 8) {
                errors.add("Password must be at least 8 characters")
            }
            if (request.name.isBlank()) {
                errors.add("Name is required")
            }

            if (errors.isNotEmpty()) {
                ValidationResult.Invalid(errors)
            } else {
                ValidationResult.Valid
            }
        }

        validate<LoginRequest> { request ->
            if (request.email.isBlank() || request.password.isBlank()) {
                ValidationResult.Invalid("Email and password are required")
            } else {
                ValidationResult.Valid
            }
        }
    }
}
```

---

## Repository

```kotlin
// src/main/kotlin/com/example/repositories/UserRepository.kt
package com.example.repositories

import com.example.models.User
import com.example.models.Users
import org.jetbrains.exposed.sql.*
import org.jetbrains.exposed.sql.SqlExpressionBuilder.eq
import org.jetbrains.exposed.sql.transactions.experimental.newSuspendedTransaction
import java.time.LocalDateTime

class UserRepository {

    private suspend fun <T> dbQuery(block: suspend () -> T): T =
        newSuspendedTransaction { block() }

    suspend fun findById(id: Long): User? = dbQuery {
        Users.select { Users.id eq id }
            .map(User::fromRow)
            .singleOrNull()
    }

    suspend fun findByEmail(email: String): User? = dbQuery {
        Users.select { Users.email eq email }
            .map(User::fromRow)
            .singleOrNull()
    }

    suspend fun findAll(limit: Int = 20, offset: Long = 0): List<User> = dbQuery {
        Users.selectAll()
            .orderBy(Users.createdAt, SortOrder.DESC)
            .limit(limit, offset)
            .map(User::fromRow)
    }

    suspend fun create(
        email: String,
        passwordHash: String,
        name: String,
        role: String = "user",
    ): User = dbQuery {
        val id = Users.insertAndGetId {
            it[Users.email] = email
            it[Users.passwordHash] = passwordHash
            it[Users.name] = name
            it[Users.role] = role
            it[createdAt] = LocalDateTime.now()
            it[updatedAt] = LocalDateTime.now()
        }

        Users.select { Users.id eq id }
            .map(User::fromRow)
            .single()
    }

    suspend fun update(id: Long, name: String?, email: String?): User? = dbQuery {
        val updated = Users.update({ Users.id eq id }) {
            name?.let { n -> it[Users.name] = n }
            email?.let { e -> it[Users.email] = e }
            it[updatedAt] = LocalDateTime.now()
        }

        if (updated > 0) {
            Users.select { Users.id eq id }
                .map(User::fromRow)
                .singleOrNull()
        } else {
            null
        }
    }

    suspend fun delete(id: Long): Boolean = dbQuery {
        Users.deleteWhere { Users.id eq id } > 0
    }

    suspend fun existsByEmail(email: String): Boolean = dbQuery {
        Users.select { Users.email eq email }.count() > 0
    }
}
```

---

## Service

```kotlin
// src/main/kotlin/com/example/services/UserService.kt
package com.example.services

import at.favre.lib.crypto.bcrypt.BCrypt
import com.example.models.*
import com.example.plugins.ConflictException
import com.example.plugins.NotFoundException
import com.example.plugins.UnauthorizedException
import com.example.repositories.UserRepository
import com.example.utils.JwtUtils

class UserService(
    private val userRepository: UserRepository,
    private val jwtUtils: JwtUtils,
) {

    suspend fun createUser(request: CreateUserRequest): UserResponse {
        if (userRepository.existsByEmail(request.email)) {
            throw ConflictException("Email already registered")
        }

        val passwordHash = BCrypt.withDefaults()
            .hashToString(12, request.password.toCharArray())

        val user = userRepository.create(
            email = request.email,
            passwordHash = passwordHash,
            name = request.name,
        )

        return UserResponse.from(user)
    }

    suspend fun getUserById(id: Long): UserResponse {
        val user = userRepository.findById(id)
            ?: throw NotFoundException("User not found")
        return UserResponse.from(user)
    }

    suspend fun getAllUsers(limit: Int, offset: Long): List<UserResponse> {
        return userRepository.findAll(limit, offset)
            .map(UserResponse::from)
    }

    suspend fun updateUser(id: Long, request: UpdateUserRequest): UserResponse {
        // Check email uniqueness if updating email
        if (request.email != null) {
            val existing = userRepository.findByEmail(request.email)
            if (existing != null && existing.id != id) {
                throw ConflictException("Email already in use")
            }
        }

        val user = userRepository.update(id, request.name, request.email)
            ?: throw NotFoundException("User not found")

        return UserResponse.from(user)
    }

    suspend fun deleteUser(id: Long) {
        if (!userRepository.delete(id)) {
            throw NotFoundException("User not found")
        }
    }

    suspend fun login(request: LoginRequest): LoginResponse {
        val user = userRepository.findByEmail(request.email)
            ?: throw UnauthorizedException("Invalid credentials")

        val passwordValid = BCrypt.verifyer()
            .verify(request.password.toCharArray(), user.passwordHash)
            .verified

        if (!passwordValid) {
            throw UnauthorizedException("Invalid credentials")
        }

        val token = jwtUtils.generateToken(user)

        return LoginResponse(
            token = token,
            user = UserResponse.from(user),
        )
    }
}
```

---

## JWT Utilities

```kotlin
// src/main/kotlin/com/example/utils/JwtUtils.kt
package com.example.utils

import com.auth0.jwt.JWT
import com.auth0.jwt.algorithms.Algorithm
import com.example.models.User
import java.util.*

class JwtUtils(
    private val secret: String,
    private val issuer: String,
    private val audience: String,
    private val expirationMs: Long,
) {

    fun generateToken(user: User): String {
        return JWT.create()
            .withAudience(audience)
            .withIssuer(issuer)
            .withClaim("userId", user.id.toString())
            .withClaim("email", user.email)
            .withClaim("role", user.role)
            .withExpiresAt(Date(System.currentTimeMillis() + expirationMs))
            .sign(Algorithm.HMAC256(secret))
    }
}
```

---

## Routing

### Main Routing Configuration

```kotlin
// src/main/kotlin/com/example/plugins/Routing.kt
package com.example.plugins

import com.example.repositories.UserRepository
import com.example.routes.authRoutes
import com.example.routes.userRoutes
import com.example.services.UserService
import com.example.utils.JwtUtils
import io.ktor.server.application.*
import io.ktor.server.plugins.callloging.*
import io.ktor.server.plugins.cors.routing.*
import io.ktor.server.request.*
import io.ktor.server.routing.*
import org.slf4j.event.Level

fun Application.configureRouting() {
    // CORS
    install(CORS) {
        anyHost()
        allowHeader("Content-Type")
        allowHeader("Authorization")
        allowMethod(io.ktor.http.HttpMethod.Options)
        allowMethod(io.ktor.http.HttpMethod.Put)
        allowMethod(io.ktor.http.HttpMethod.Delete)
    }

    // Call Logging
    install(CallLogging) {
        level = Level.INFO
        filter { call -> call.request.path().startsWith("/api") }
        format { call ->
            val status = call.response.status()
            val method = call.request.httpMethod.value
            val path = call.request.path()
            "$method $path - $status"
        }
    }

    // Initialize dependencies
    val config = environment.config
    val userRepository = UserRepository()
    val jwtUtils = JwtUtils(
        secret = config.property("jwt.secret").getString(),
        issuer = config.property("jwt.issuer").getString(),
        audience = config.property("jwt.audience").getString(),
        expirationMs = config.property("jwt.expirationMs").getString().toLong(),
    )
    val userService = UserService(userRepository, jwtUtils)

    routing {
        route("/api") {
            authRoutes(userService)
            userRoutes(userService)
        }
    }
}
```

### User Routes

```kotlin
// src/main/kotlin/com/example/routes/UserRoutes.kt
package com.example.routes

import com.example.models.CreateUserRequest
import com.example.models.UpdateUserRequest
import com.example.services.UserService
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.auth.*
import io.ktor.server.auth.jwt.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

fun Route.userRoutes(userService: UserService) {
    route("/users") {
        // Public: Create user
        post {
            val request = call.receive<CreateUserRequest>()
            val user = userService.createUser(request)
            call.respond(HttpStatusCode.Created, user)
        }

        // Protected routes
        authenticate("auth-jwt") {
            // Get all users
            get {
                val limit = call.request.queryParameters["limit"]?.toIntOrNull() ?: 20
                val offset = call.request.queryParameters["offset"]?.toLongOrNull() ?: 0
                val users = userService.getAllUsers(limit, offset)
                call.respond(users)
            }

            // Get current user
            get("/me") {
                val principal = call.principal<JWTPrincipal>()
                val userId = principal!!.payload.getClaim("userId").asString().toLong()
                val user = userService.getUserById(userId)
                call.respond(user)
            }

            // Get user by ID
            get("/{id}") {
                val id = call.parameters["id"]?.toLongOrNull()
                    ?: return@get call.respond(HttpStatusCode.BadRequest, "Invalid ID")
                val user = userService.getUserById(id)
                call.respond(user)
            }

            // Update user
            put("/{id}") {
                val id = call.parameters["id"]?.toLongOrNull()
                    ?: return@put call.respond(HttpStatusCode.BadRequest, "Invalid ID")
                val request = call.receive<UpdateUserRequest>()
                val user = userService.updateUser(id, request)
                call.respond(user)
            }

            // Delete user
            delete("/{id}") {
                val id = call.parameters["id"]?.toLongOrNull()
                    ?: return@delete call.respond(HttpStatusCode.BadRequest, "Invalid ID")
                userService.deleteUser(id)
                call.respond(HttpStatusCode.NoContent)
            }
        }
    }
}
```

### Auth Routes

```kotlin
// src/main/kotlin/com/example/routes/AuthRoutes.kt
package com.example.routes

import com.example.models.LoginRequest
import com.example.services.UserService
import io.ktor.server.application.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

fun Route.authRoutes(userService: UserService) {
    route("/auth") {
        post("/login") {
            val request = call.receive<LoginRequest>()
            val response = userService.login(request)
            call.respond(response)
        }
    }
}
```

---

## WebSocket Support

```kotlin
// src/main/kotlin/com/example/routes/WebSocketRoutes.kt
package com.example.routes

import io.ktor.server.routing.*
import io.ktor.server.websocket.*
import io.ktor.websocket.*
import kotlinx.coroutines.channels.consumeEach
import java.util.concurrent.ConcurrentHashMap

class ConnectionManager {
    private val connections = ConcurrentHashMap<String, WebSocketSession>()

    fun addConnection(userId: String, session: WebSocketSession) {
        connections[userId] = session
    }

    fun removeConnection(userId: String) {
        connections.remove(userId)
    }

    suspend fun broadcast(message: String) {
        connections.values.forEach { session ->
            session.send(Frame.Text(message))
        }
    }

    suspend fun sendTo(userId: String, message: String) {
        connections[userId]?.send(Frame.Text(message))
    }
}

fun Route.webSocketRoutes(connectionManager: ConnectionManager) {
    webSocket("/ws/{userId}") {
        val userId = call.parameters["userId"] ?: return@webSocket close(
            CloseReason(CloseReason.Codes.VIOLATED_POLICY, "No user ID")
        )

        connectionManager.addConnection(userId, this)

        try {
            incoming.consumeEach { frame ->
                when (frame) {
                    is Frame.Text -> {
                        val text = frame.readText()
                        connectionManager.broadcast("$userId: $text")
                    }
                    else -> {}
                }
            }
        } finally {
            connectionManager.removeConnection(userId)
        }
    }
}
```

---

## Testing

### Application Test Setup

```kotlin
// src/test/kotlin/com/example/ApplicationTest.kt
package com.example

import com.example.models.CreateUserRequest
import com.example.models.LoginRequest
import com.example.models.LoginResponse
import com.example.models.UserResponse
import io.ktor.client.call.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.request.*
import io.ktor.http.*
import io.ktor.serialization.kotlinx.json.*
import io.ktor.server.testing.*
import kotlin.test.*

class ApplicationTest {

    @Test
    fun `create user returns 201`() = testApplication {
        application {
            module()
        }

        val client = createClient {
            install(ContentNegotiation) {
                json()
            }
        }

        val response = client.post("/api/users") {
            contentType(ContentType.Application.Json)
            setBody(CreateUserRequest(
                email = "test@example.com",
                password = "password123",
                name = "Test User",
            ))
        }

        assertEquals(HttpStatusCode.Created, response.status)
        val user = response.body<UserResponse>()
        assertEquals("test@example.com", user.email)
        assertEquals("Test User", user.name)
    }

    @Test
    fun `login returns token`() = testApplication {
        application {
            module()
        }

        val client = createClient {
            install(ContentNegotiation) {
                json()
            }
        }

        // First create user
        client.post("/api/users") {
            contentType(ContentType.Application.Json)
            setBody(CreateUserRequest(
                email = "login@example.com",
                password = "password123",
                name = "Login User",
            ))
        }

        // Then login
        val response = client.post("/api/auth/login") {
            contentType(ContentType.Application.Json)
            setBody(LoginRequest(
                email = "login@example.com",
                password = "password123",
            ))
        }

        assertEquals(HttpStatusCode.OK, response.status)
        val loginResponse = response.body<LoginResponse>()
        assertNotNull(loginResponse.token)
    }

    @Test
    fun `get users requires authentication`() = testApplication {
        application {
            module()
        }

        val response = client.get("/api/users")

        assertEquals(HttpStatusCode.Unauthorized, response.status)
    }

    @Test
    fun `get users with valid token returns users`() = testApplication {
        application {
            module()
        }

        val client = createClient {
            install(ContentNegotiation) {
                json()
            }
        }

        // Create and login
        client.post("/api/users") {
            contentType(ContentType.Application.Json)
            setBody(CreateUserRequest(
                email = "auth@example.com",
                password = "password123",
                name = "Auth User",
            ))
        }

        val loginResponse = client.post("/api/auth/login") {
            contentType(ContentType.Application.Json)
            setBody(LoginRequest(
                email = "auth@example.com",
                password = "password123",
            ))
        }.body<LoginResponse>()

        // Use token
        val response = client.get("/api/users") {
            header(HttpHeaders.Authorization, "Bearer ${loginResponse.token}")
        }

        assertEquals(HttpStatusCode.OK, response.status)
        val users = response.body<List<UserResponse>>()
        assertTrue(users.isNotEmpty())
    }
}
```

### Service Unit Tests

```kotlin
// src/test/kotlin/com/example/services/UserServiceTest.kt
package com.example.services

import com.example.models.CreateUserRequest
import com.example.models.User
import com.example.plugins.ConflictException
import com.example.plugins.NotFoundException
import com.example.repositories.UserRepository
import com.example.utils.JwtUtils
import io.mockk.*
import kotlinx.coroutines.runBlocking
import java.time.LocalDateTime
import kotlin.test.*

class UserServiceTest {

    private lateinit var userRepository: UserRepository
    private lateinit var jwtUtils: JwtUtils
    private lateinit var userService: UserService

    @BeforeTest
    fun setup() {
        userRepository = mockk()
        jwtUtils = mockk()
        userService = UserService(userRepository, jwtUtils)
    }

    @AfterTest
    fun teardown() {
        clearAllMocks()
    }

    @Test
    fun `createUser with valid data returns user`() = runBlocking {
        val request = CreateUserRequest(
            email = "test@example.com",
            password = "password123",
            name = "Test User",
        )

        val user = User(
            id = 1,
            email = "test@example.com",
            passwordHash = "hashed",
            name = "Test User",
            role = "user",
            createdAt = LocalDateTime.now(),
            updatedAt = LocalDateTime.now(),
        )

        coEvery { userRepository.existsByEmail(any()) } returns false
        coEvery { userRepository.create(any(), any(), any(), any()) } returns user

        val result = userService.createUser(request)

        assertEquals("test@example.com", result.email)
        assertEquals("Test User", result.name)
        coVerify { userRepository.create(any(), any(), any(), any()) }
    }

    @Test
    fun `createUser with existing email throws ConflictException`(): Unit = runBlocking {
        val request = CreateUserRequest(
            email = "existing@example.com",
            password = "password123",
            name = "Test User",
        )

        coEvery { userRepository.existsByEmail(any()) } returns true

        assertFailsWith<ConflictException> {
            userService.createUser(request)
        }
    }

    @Test
    fun `getUserById with non-existent id throws NotFoundException`(): Unit = runBlocking {
        coEvery { userRepository.findById(any()) } returns null

        assertFailsWith<NotFoundException> {
            userService.getUserById(999)
        }
    }
}
```

---

## Commands

```bash
# Development
./gradlew run

# Build
./gradlew build

# Build fat JAR
./gradlew buildFatJar

# Run fat JAR
java -jar build/libs/app.jar

# Test
./gradlew test

# Format (if using ktlint)
./gradlew ktlintFormat

# Clean
./gradlew clean
```

---

## Best Practices

### Do's
- ✅ Use dependency injection pattern (pass services to routes)
- ✅ Use `suspend` functions for async database operations
- ✅ Configure plugins in separate files
- ✅ Use typed routes with Route extension functions
- ✅ Use `kotlinx.serialization` for JSON
- ✅ Use Exposed's `newSuspendedTransaction` for coroutine-safe DB access
- ✅ Handle errors with StatusPages plugin
- ✅ Use application.conf for configuration

### Don'ts
- ❌ Don't block coroutines with synchronous calls
- ❌ Don't use global state without proper synchronization
- ❌ Don't expose internal exceptions to clients
- ❌ Don't hardcode configuration values
- ❌ Don't skip request validation
- ❌ Don't mix business logic into routes

---

## Comparison: Ktor vs Spring Boot (Kotlin)

| Feature | Ktor | Spring Boot (Kotlin) |
|---------|------|---------------------|
| **Weight** | Lightweight, modular | Full-featured, heavier |
| **Learning Curve** | Lower (Kotlin-native) | Moderate (Spring ecosystem) |
| **Coroutines** | Native, first-class | Supported (WebFlux) |
| **Configuration** | HOCON/programmatic | YAML/properties, annotations |
| **ORM** | Exposed (recommended) | JPA/Hibernate |
| **Testing** | Built-in testApplication | Spring Test |
| **Best For** | Microservices, lightweight APIs | Enterprise, full-stack apps |
| **Startup Time** | Very fast | Slower (reflection) |
| **Memory** | Lower footprint | Higher footprint |
| **Ecosystem** | Growing | Mature, extensive |

---

## When to Use Ktor

**Choose Ktor when**:
- Building lightweight microservices
- Kotlin is your primary language
- You need fast startup times
- You want fine-grained control over dependencies
- Building real-time applications (WebSockets)
- Creating serverless functions

**Consider alternatives when**:
- You need extensive enterprise integrations
- Team is more familiar with Spring
- Project requires complex security configurations
- You need extensive third-party library support
