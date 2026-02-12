# Spring Boot (Kotlin) Framework Guide

> **Applies to**: Spring Boot 3.x, Kotlin 1.9+, REST APIs, Microservices
> **Use with**: `.agent/skills/kotlin-guide/SKILL.md`

---

## Overview

Spring Boot with Kotlin combines the power of Spring's ecosystem with Kotlin's modern language features. Spring Boot 3.x offers first-class Kotlin support including coroutines, null safety, and DSL-style configuration.

### When to Use Spring Boot (Kotlin)
- **Enterprise Applications**: Proven enterprise-grade framework
- **Microservices**: Excellent with Spring Cloud
- **Existing Spring Ecosystem**: Leverage vast Spring libraries
- **Team Familiarity**: Teams with Java/Spring background
- **Complex Business Logic**: Rich feature set for complex domains

### When NOT to Use
- **Lightweight Services**: Consider Ktor for simpler services
- **Mobile Backend**: Ktor is more lightweight
- **Minimal Dependencies**: Spring Boot has a larger footprint

---

## Project Structure

```
myproject/
├── build.gradle.kts
├── settings.gradle.kts
├── src/
│   ├── main/
│   │   ├── kotlin/
│   │   │   └── com/example/myproject/
│   │   │       ├── MyProjectApplication.kt
│   │   │       ├── config/
│   │   │       │   ├── SecurityConfig.kt
│   │   │       │   └── WebConfig.kt
│   │   │       ├── controller/
│   │   │       │   └── UserController.kt
│   │   │       ├── service/
│   │   │       │   └── UserService.kt
│   │   │       ├── repository/
│   │   │       │   └── UserRepository.kt
│   │   │       ├── model/
│   │   │       │   ├── entity/
│   │   │       │   │   └── User.kt
│   │   │       │   └── dto/
│   │   │       │       └── UserDto.kt
│   │   │       ├── exception/
│   │   │       │   ├── GlobalExceptionHandler.kt
│   │   │       │   └── Exceptions.kt
│   │   │       └── security/
│   │   │           └── JwtService.kt
│   │   └── resources/
│   │       ├── application.yml
│   │       ├── application-dev.yml
│   │       └── db/migration/
│   └── test/
│       └── kotlin/
│           └── com/example/myproject/
│               ├── controller/
│               ├── service/
│               └── integration/
└── README.md
```

---

## Dependencies

```kotlin
// build.gradle.kts
import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    id("org.springframework.boot") version "3.2.0"
    id("io.spring.dependency-management") version "1.1.4"
    kotlin("jvm") version "1.9.21"
    kotlin("plugin.spring") version "1.9.21"
    kotlin("plugin.jpa") version "1.9.21"
}

group = "com.example"
version = "0.0.1-SNAPSHOT"

java {
    sourceCompatibility = JavaVersion.VERSION_21
}

repositories {
    mavenCentral()
}

dependencies {
    // Spring Boot
    implementation("org.springframework.boot:spring-boot-starter-web")
    implementation("org.springframework.boot:spring-boot-starter-validation")
    implementation("org.springframework.boot:spring-boot-starter-security")
    implementation("org.springframework.boot:spring-boot-starter-data-jpa")

    // Kotlin
    implementation("com.fasterxml.jackson.module:jackson-module-kotlin")
    implementation("org.jetbrains.kotlin:kotlin-reflect")

    // Coroutines (optional but recommended)
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-reactor")

    // Database
    runtimeOnly("org.postgresql:postgresql")

    // JWT
    implementation("io.jsonwebtoken:jjwt-api:0.12.3")
    runtimeOnly("io.jsonwebtoken:jjwt-impl:0.12.3")
    runtimeOnly("io.jsonwebtoken:jjwt-jackson:0.12.3")

    // Development
    developmentOnly("org.springframework.boot:spring-boot-devtools")

    // Testing
    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testImplementation("org.springframework.security:spring-security-test")
    testImplementation("io.mockk:mockk:1.13.8")
    testImplementation("com.ninja-squad:springmockk:4.0.2")
}

tasks.withType<KotlinCompile> {
    kotlinOptions {
        freeCompilerArgs += "-Xjsr305=strict"
        jvmTarget = "21"
    }
}

tasks.withType<Test> {
    useJUnitPlatform()
}
```

---

## Application Entry Point

```kotlin
// src/main/kotlin/com/example/myproject/MyProjectApplication.kt
package com.example.myproject

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
class MyProjectApplication

fun main(args: Array<String>) {
    runApplication<MyProjectApplication>(*args)
}
```

---

## Configuration

```yaml
# src/main/resources/application.yml
spring:
  application:
    name: myproject
  datasource:
    url: jdbc:postgresql://localhost:5432/mydb
    username: ${DB_USERNAME:postgres}
    password: ${DB_PASSWORD:postgres}
    driver-class-name: org.postgresql.Driver
  jpa:
    hibernate:
      ddl-auto: validate
    show-sql: false
    properties:
      hibernate:
        dialect: org.hibernate.dialect.PostgreSQLDialect
        format_sql: true

server:
  port: 8080
  error:
    include-message: always
    include-binding-errors: always

jwt:
  secret: ${JWT_SECRET:your-256-bit-secret-key-here}
  expiration: 86400000  # 24 hours in milliseconds

logging:
  level:
    com.example.myproject: DEBUG
    org.springframework.security: DEBUG
```

```kotlin
// src/main/kotlin/com/example/myproject/config/JwtProperties.kt
package com.example.myproject.config

import org.springframework.boot.context.properties.ConfigurationProperties

@ConfigurationProperties(prefix = "jwt")
data class JwtProperties(
    val secret: String,
    val expiration: Long,
)
```

---

## Entity and DTOs

```kotlin
// src/main/kotlin/com/example/myproject/model/entity/User.kt
package com.example.myproject.model.entity

import jakarta.persistence.*
import org.hibernate.annotations.CreationTimestamp
import org.hibernate.annotations.UpdateTimestamp
import java.time.Instant
import java.util.*

@Entity
@Table(name = "users")
data class User(
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    val id: UUID? = null,

    @Column(unique = true, nullable = false)
    val email: String,

    @Column(nullable = false)
    val passwordHash: String,

    @Column(nullable = false)
    val name: String,

    @Column(nullable = false)
    val role: String = "USER",

    @CreationTimestamp
    @Column(updatable = false)
    val createdAt: Instant? = null,

    @UpdateTimestamp
    val updatedAt: Instant? = null,
)

// src/main/kotlin/com/example/myproject/model/dto/UserDto.kt
package com.example.myproject.model.dto

import com.example.myproject.model.entity.User
import jakarta.validation.constraints.Email
import jakarta.validation.constraints.NotBlank
import jakarta.validation.constraints.Size
import java.time.Instant
import java.util.*

data class CreateUserRequest(
    @field:NotBlank(message = "Email is required")
    @field:Email(message = "Invalid email format")
    val email: String,

    @field:NotBlank(message = "Password is required")
    @field:Size(min = 8, message = "Password must be at least 8 characters")
    val password: String,

    @field:NotBlank(message = "Name is required")
    @field:Size(min = 1, max = 100, message = "Name must be between 1 and 100 characters")
    val name: String,
)

data class UpdateUserRequest(
    @field:Size(min = 1, max = 100, message = "Name must be between 1 and 100 characters")
    val name: String? = null,

    @field:Email(message = "Invalid email format")
    val email: String? = null,
)

data class LoginRequest(
    @field:NotBlank(message = "Email is required")
    @field:Email(message = "Invalid email format")
    val email: String,

    @field:NotBlank(message = "Password is required")
    val password: String,
)

data class LoginResponse(
    val token: String,
    val user: UserResponse,
)

data class UserResponse(
    val id: UUID,
    val email: String,
    val name: String,
    val role: String,
    val createdAt: Instant,
) {
    companion object {
        fun from(user: User): UserResponse = UserResponse(
            id = user.id!!,
            email = user.email,
            name = user.name,
            role = user.role,
            createdAt = user.createdAt!!,
        )
    }
}
```

---

## Repository

```kotlin
// src/main/kotlin/com/example/myproject/repository/UserRepository.kt
package com.example.myproject.repository

import com.example.myproject.model.entity.User
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository
import java.util.*

@Repository
interface UserRepository : JpaRepository<User, UUID> {
    fun findByEmail(email: String): User?
    fun existsByEmail(email: String): Boolean
}
```

---

## Service Layer

```kotlin
// src/main/kotlin/com/example/myproject/service/UserService.kt
package com.example.myproject.service

import com.example.myproject.exception.ConflictException
import com.example.myproject.exception.NotFoundException
import com.example.myproject.exception.UnauthorizedException
import com.example.myproject.model.dto.*
import com.example.myproject.model.entity.User
import com.example.myproject.repository.UserRepository
import com.example.myproject.security.JwtService
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional
import java.util.*

@Service
class UserService(
    private val userRepository: UserRepository,
    private val passwordEncoder: PasswordEncoder,
    private val jwtService: JwtService,
) {
    @Transactional
    fun createUser(request: CreateUserRequest): UserResponse {
        if (userRepository.existsByEmail(request.email)) {
            throw ConflictException("Email already registered")
        }

        val user = User(
            email = request.email,
            passwordHash = passwordEncoder.encode(request.password),
            name = request.name,
        )

        val savedUser = userRepository.save(user)
        return UserResponse.from(savedUser)
    }

    fun login(request: LoginRequest): LoginResponse {
        val user = userRepository.findByEmail(request.email)
            ?: throw UnauthorizedException("Invalid credentials")

        if (!passwordEncoder.matches(request.password, user.passwordHash)) {
            throw UnauthorizedException("Invalid credentials")
        }

        val token = jwtService.generateToken(user)
        return LoginResponse(token = token, user = UserResponse.from(user))
    }

    @Transactional(readOnly = true)
    fun getUserById(id: UUID): UserResponse {
        val user = userRepository.findById(id)
            .orElseThrow { NotFoundException("User not found: $id") }
        return UserResponse.from(user)
    }

    @Transactional(readOnly = true)
    fun getAllUsers(pageable: Pageable): Page<UserResponse> {
        return userRepository.findAll(pageable).map { UserResponse.from(it) }
    }

    @Transactional
    fun updateUser(id: UUID, request: UpdateUserRequest): UserResponse {
        val user = userRepository.findById(id)
            .orElseThrow { NotFoundException("User not found: $id") }

        val updatedUser = user.copy(
            name = request.name ?: user.name,
            email = request.email ?: user.email,
        )

        return UserResponse.from(userRepository.save(updatedUser))
    }

    @Transactional
    fun deleteUser(id: UUID) {
        if (!userRepository.existsById(id)) {
            throw NotFoundException("User not found: $id")
        }
        userRepository.deleteById(id)
    }
}
```

---

## Controller

```kotlin
// src/main/kotlin/com/example/myproject/controller/UserController.kt
package com.example.myproject.controller

import com.example.myproject.model.dto.*
import com.example.myproject.security.CurrentUser
import com.example.myproject.service.UserService
import jakarta.validation.Valid
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.data.web.PageableDefault
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.security.access.prepost.PreAuthorize
import org.springframework.web.bind.annotation.*
import java.util.*

@RestController
@RequestMapping("/api")
class UserController(
    private val userService: UserService,
) {
    @PostMapping("/register")
    fun register(@Valid @RequestBody request: CreateUserRequest): ResponseEntity<UserResponse> {
        val user = userService.createUser(request)
        return ResponseEntity.status(HttpStatus.CREATED).body(user)
    }

    @PostMapping("/login")
    fun login(@Valid @RequestBody request: LoginRequest): ResponseEntity<LoginResponse> {
        val response = userService.login(request)
        return ResponseEntity.ok(response)
    }

    @GetMapping("/users/me")
    fun getCurrentUser(@CurrentUser userId: UUID): ResponseEntity<UserResponse> {
        val user = userService.getUserById(userId)
        return ResponseEntity.ok(user)
    }

    @PutMapping("/users/me")
    fun updateCurrentUser(
        @CurrentUser userId: UUID,
        @Valid @RequestBody request: UpdateUserRequest,
    ): ResponseEntity<UserResponse> {
        val user = userService.updateUser(userId, request)
        return ResponseEntity.ok(user)
    }

    @GetMapping("/users/{id}")
    @PreAuthorize("hasRole('ADMIN') or @userSecurity.isOwner(#id, authentication)")
    fun getUser(@PathVariable id: UUID): ResponseEntity<UserResponse> {
        val user = userService.getUserById(id)
        return ResponseEntity.ok(user)
    }

    @GetMapping("/users")
    @PreAuthorize("hasRole('ADMIN')")
    fun getAllUsers(
        @PageableDefault(size = 20) pageable: Pageable,
    ): ResponseEntity<Page<UserResponse>> {
        val users = userService.getAllUsers(pageable)
        return ResponseEntity.ok(users)
    }

    @DeleteMapping("/users/{id}")
    @PreAuthorize("hasRole('ADMIN')")
    fun deleteUser(@PathVariable id: UUID): ResponseEntity<Unit> {
        userService.deleteUser(id)
        return ResponseEntity.noContent().build()
    }
}

// Health controller
@RestController
@RequestMapping("/health")
class HealthController {
    @GetMapping
    fun health() = mapOf(
        "status" to "ok",
        "timestamp" to System.currentTimeMillis(),
    )
}
```

---

## Security Configuration

```kotlin
// src/main/kotlin/com/example/myproject/config/SecurityConfig.kt
package com.example.myproject.config

import com.example.myproject.security.JwtAuthenticationFilter
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.security.config.annotation.method.configuration.EnableMethodSecurity
import org.springframework.security.config.annotation.web.builders.HttpSecurity
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity
import org.springframework.security.config.http.SessionCreationPolicy
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.security.web.SecurityFilterChain
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter

@Configuration
@EnableWebSecurity
@EnableMethodSecurity
class SecurityConfig(
    private val jwtAuthenticationFilter: JwtAuthenticationFilter,
) {
    @Bean
    fun securityFilterChain(http: HttpSecurity): SecurityFilterChain {
        return http
            .csrf { it.disable() }
            .sessionManagement { it.sessionCreationPolicy(SessionCreationPolicy.STATELESS) }
            .authorizeHttpRequests { auth ->
                auth
                    .requestMatchers("/api/register", "/api/login").permitAll()
                    .requestMatchers("/health/**").permitAll()
                    .requestMatchers("/actuator/**").permitAll()
                    .anyRequest().authenticated()
            }
            .addFilterBefore(jwtAuthenticationFilter, UsernamePasswordAuthenticationFilter::class.java)
            .build()
    }

    @Bean
    fun passwordEncoder(): PasswordEncoder = BCryptPasswordEncoder(12)
}
```

```kotlin
// src/main/kotlin/com/example/myproject/security/JwtService.kt
package com.example.myproject.security

import com.example.myproject.config.JwtProperties
import com.example.myproject.model.entity.User
import io.jsonwebtoken.Claims
import io.jsonwebtoken.Jwts
import io.jsonwebtoken.security.Keys
import org.springframework.stereotype.Service
import java.util.*
import javax.crypto.SecretKey

@Service
class JwtService(
    private val jwtProperties: JwtProperties,
) {
    private val secretKey: SecretKey = Keys.hmacShaKeyFor(jwtProperties.secret.toByteArray())

    fun generateToken(user: User): String {
        val now = Date()
        val expiration = Date(now.time + jwtProperties.expiration)

        return Jwts.builder()
            .subject(user.id.toString())
            .claim("email", user.email)
            .claim("role", user.role)
            .issuedAt(now)
            .expiration(expiration)
            .signWith(secretKey)
            .compact()
    }

    fun validateToken(token: String): Boolean {
        return try {
            val claims = extractAllClaims(token)
            !claims.expiration.before(Date())
        } catch (e: Exception) {
            false
        }
    }

    fun extractUserId(token: String): UUID {
        val subject = extractAllClaims(token).subject
        return UUID.fromString(subject)
    }

    fun extractRole(token: String): String {
        return extractAllClaims(token)["role"] as String
    }

    private fun extractAllClaims(token: String): Claims {
        return Jwts.parser()
            .verifyWith(secretKey)
            .build()
            .parseSignedClaims(token)
            .payload
    }
}

// src/main/kotlin/com/example/myproject/security/JwtAuthenticationFilter.kt
package com.example.myproject.security

import jakarta.servlet.FilterChain
import jakarta.servlet.http.HttpServletRequest
import jakarta.servlet.http.HttpServletResponse
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken
import org.springframework.security.core.authority.SimpleGrantedAuthority
import org.springframework.security.core.context.SecurityContextHolder
import org.springframework.stereotype.Component
import org.springframework.web.filter.OncePerRequestFilter

@Component
class JwtAuthenticationFilter(
    private val jwtService: JwtService,
) : OncePerRequestFilter() {

    override fun doFilterInternal(
        request: HttpServletRequest,
        response: HttpServletResponse,
        filterChain: FilterChain,
    ) {
        val authHeader = request.getHeader("Authorization")

        if (authHeader == null || !authHeader.startsWith("Bearer ")) {
            filterChain.doFilter(request, response)
            return
        }

        val token = authHeader.substring(7)

        if (jwtService.validateToken(token)) {
            val userId = jwtService.extractUserId(token)
            val role = jwtService.extractRole(token)

            val authorities = listOf(SimpleGrantedAuthority("ROLE_$role"))
            val authentication = UsernamePasswordAuthenticationToken(userId, null, authorities)

            SecurityContextHolder.getContext().authentication = authentication
        }

        filterChain.doFilter(request, response)
    }
}

// src/main/kotlin/com/example/myproject/security/CurrentUser.kt
package com.example.myproject.security

import org.springframework.security.core.annotation.AuthenticationPrincipal

@Target(AnnotationTarget.VALUE_PARAMETER)
@Retention(AnnotationRetention.RUNTIME)
@AuthenticationPrincipal
annotation class CurrentUser
```

---

## Exception Handling

```kotlin
// src/main/kotlin/com/example/myproject/exception/Exceptions.kt
package com.example.myproject.exception

class NotFoundException(message: String) : RuntimeException(message)
class UnauthorizedException(message: String) : RuntimeException(message)
class ConflictException(message: String) : RuntimeException(message)
class ValidationException(message: String) : RuntimeException(message)

// src/main/kotlin/com/example/myproject/exception/GlobalExceptionHandler.kt
package com.example.myproject.exception

import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.security.access.AccessDeniedException
import org.springframework.web.bind.MethodArgumentNotValidException
import org.springframework.web.bind.annotation.ExceptionHandler
import org.springframework.web.bind.annotation.RestControllerAdvice
import java.time.Instant

data class ErrorResponse(
    val timestamp: Instant = Instant.now(),
    val status: Int,
    val error: String,
    val message: String,
    val details: List<String>? = null,
)

@RestControllerAdvice
class GlobalExceptionHandler {

    @ExceptionHandler(NotFoundException::class)
    fun handleNotFound(ex: NotFoundException): ResponseEntity<ErrorResponse> {
        return ResponseEntity
            .status(HttpStatus.NOT_FOUND)
            .body(ErrorResponse(
                status = 404,
                error = "NOT_FOUND",
                message = ex.message ?: "Resource not found",
            ))
    }

    @ExceptionHandler(UnauthorizedException::class)
    fun handleUnauthorized(ex: UnauthorizedException): ResponseEntity<ErrorResponse> {
        return ResponseEntity
            .status(HttpStatus.UNAUTHORIZED)
            .body(ErrorResponse(
                status = 401,
                error = "UNAUTHORIZED",
                message = ex.message ?: "Unauthorized",
            ))
    }

    @ExceptionHandler(ConflictException::class)
    fun handleConflict(ex: ConflictException): ResponseEntity<ErrorResponse> {
        return ResponseEntity
            .status(HttpStatus.CONFLICT)
            .body(ErrorResponse(
                status = 409,
                error = "CONFLICT",
                message = ex.message ?: "Resource conflict",
            ))
    }

    @ExceptionHandler(AccessDeniedException::class)
    fun handleAccessDenied(ex: AccessDeniedException): ResponseEntity<ErrorResponse> {
        return ResponseEntity
            .status(HttpStatus.FORBIDDEN)
            .body(ErrorResponse(
                status = 403,
                error = "FORBIDDEN",
                message = "Access denied",
            ))
    }

    @ExceptionHandler(MethodArgumentNotValidException::class)
    fun handleValidation(ex: MethodArgumentNotValidException): ResponseEntity<ErrorResponse> {
        val details = ex.bindingResult.fieldErrors.map { "${it.field}: ${it.defaultMessage}" }
        return ResponseEntity
            .status(HttpStatus.BAD_REQUEST)
            .body(ErrorResponse(
                status = 400,
                error = "VALIDATION_ERROR",
                message = "Validation failed",
                details = details,
            ))
    }

    @ExceptionHandler(Exception::class)
    fun handleGeneric(ex: Exception): ResponseEntity<ErrorResponse> {
        ex.printStackTrace()
        return ResponseEntity
            .status(HttpStatus.INTERNAL_SERVER_ERROR)
            .body(ErrorResponse(
                status = 500,
                error = "INTERNAL_ERROR",
                message = "An unexpected error occurred",
            ))
    }
}
```

---

## Testing

```kotlin
// src/test/kotlin/com/example/myproject/service/UserServiceTest.kt
package com.example.myproject.service

import com.example.myproject.exception.ConflictException
import com.example.myproject.exception.UnauthorizedException
import com.example.myproject.model.dto.CreateUserRequest
import com.example.myproject.model.dto.LoginRequest
import com.example.myproject.model.entity.User
import com.example.myproject.repository.UserRepository
import com.example.myproject.security.JwtService
import io.mockk.*
import io.mockk.impl.annotations.InjectMockKs
import io.mockk.impl.annotations.MockK
import io.mockk.junit5.MockKExtension
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.assertThrows
import org.junit.jupiter.api.extension.ExtendWith
import org.springframework.security.crypto.password.PasswordEncoder
import java.time.Instant
import java.util.*
import kotlin.test.assertEquals
import kotlin.test.assertNotNull

@ExtendWith(MockKExtension::class)
class UserServiceTest {

    @MockK
    lateinit var userRepository: UserRepository

    @MockK
    lateinit var passwordEncoder: PasswordEncoder

    @MockK
    lateinit var jwtService: JwtService

    @InjectMockKs
    lateinit var userService: UserService

    private val testUser = User(
        id = UUID.randomUUID(),
        email = "test@example.com",
        passwordHash = "hashedPassword",
        name = "Test User",
        role = "USER",
        createdAt = Instant.now(),
        updatedAt = Instant.now(),
    )

    @BeforeEach
    fun setUp() {
        clearAllMocks()
    }

    @Test
    fun `createUser should create user successfully`() {
        val request = CreateUserRequest(
            email = "new@example.com",
            password = "password123",
            name = "New User",
        )

        every { userRepository.existsByEmail(request.email) } returns false
        every { passwordEncoder.encode(request.password) } returns "hashedPassword"
        every { userRepository.save(any()) } returns testUser.copy(
            email = request.email,
            name = request.name,
        )

        val result = userService.createUser(request)

        assertNotNull(result)
        verify { userRepository.save(any()) }
    }

    @Test
    fun `createUser should throw ConflictException when email exists`() {
        val request = CreateUserRequest(
            email = "existing@example.com",
            password = "password123",
            name = "User",
        )

        every { userRepository.existsByEmail(request.email) } returns true

        assertThrows<ConflictException> {
            userService.createUser(request)
        }
    }

    @Test
    fun `login should return token for valid credentials`() {
        val request = LoginRequest(
            email = "test@example.com",
            password = "password123",
        )

        every { userRepository.findByEmail(request.email) } returns testUser
        every { passwordEncoder.matches(request.password, testUser.passwordHash) } returns true
        every { jwtService.generateToken(testUser) } returns "jwt-token"

        val result = userService.login(request)

        assertEquals("jwt-token", result.token)
        assertNotNull(result.user)
    }

    @Test
    fun `login should throw UnauthorizedException for invalid password`() {
        val request = LoginRequest(
            email = "test@example.com",
            password = "wrongpassword",
        )

        every { userRepository.findByEmail(request.email) } returns testUser
        every { passwordEncoder.matches(request.password, testUser.passwordHash) } returns false

        assertThrows<UnauthorizedException> {
            userService.login(request)
        }
    }
}

// src/test/kotlin/com/example/myproject/controller/UserControllerIntegrationTest.kt
package com.example.myproject.controller

import com.example.myproject.model.dto.CreateUserRequest
import com.example.myproject.model.dto.LoginRequest
import com.fasterxml.jackson.databind.ObjectMapper
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.http.MediaType
import org.springframework.test.context.ActiveProfiles
import org.springframework.test.web.servlet.MockMvc
import org.springframework.test.web.servlet.post
import org.springframework.test.web.servlet.get
import org.springframework.transaction.annotation.Transactional

@SpringBootTest
@AutoConfigureMockMvc
@ActiveProfiles("test")
@Transactional
class UserControllerIntegrationTest {

    @Autowired
    lateinit var mockMvc: MockMvc

    @Autowired
    lateinit var objectMapper: ObjectMapper

    @Test
    fun `register should create new user`() {
        val request = CreateUserRequest(
            email = "newuser@example.com",
            password = "password123",
            name = "New User",
        )

        mockMvc.post("/api/register") {
            contentType = MediaType.APPLICATION_JSON
            content = objectMapper.writeValueAsString(request)
        }.andExpect {
            status { isCreated() }
            jsonPath("$.email") { value("newuser@example.com") }
            jsonPath("$.name") { value("New User") }
        }
    }

    @Test
    fun `register should fail with invalid email`() {
        val request = CreateUserRequest(
            email = "invalid-email",
            password = "password123",
            name = "User",
        )

        mockMvc.post("/api/register") {
            contentType = MediaType.APPLICATION_JSON
            content = objectMapper.writeValueAsString(request)
        }.andExpect {
            status { isBadRequest() }
        }
    }

    @Test
    fun `getCurrentUser should require authentication`() {
        mockMvc.get("/api/users/me")
            .andExpect {
                status { isUnauthorized() }
            }
    }
}
```

---

## Commands

```bash
# Development
./gradlew bootRun

# Build
./gradlew build

# Test
./gradlew test

# Test with coverage
./gradlew test jacocoTestReport

# Format (with ktlint)
./gradlew ktlintFormat

# Lint
./gradlew ktlintCheck

# Build JAR
./gradlew bootJar

# Run JAR
java -jar build/libs/myproject-0.0.1-SNAPSHOT.jar

# Docker build
docker build -t myproject .
```

---

## Best Practices

### DO
- ✓ Use data classes for DTOs and entities
- ✓ Use Kotlin null safety features
- ✓ Use `@Transactional(readOnly = true)` for read operations
- ✓ Use constructor injection (default in Kotlin)
- ✓ Use coroutines for async operations where beneficial
- ✓ Use `@Valid` for request validation
- ✓ Use Spring profiles for environment configuration

### DON'T
- ✗ Use `lateinit var` for dependencies (use constructor injection)
- ✗ Ignore null safety in Spring Data repositories
- ✗ Use `!!` operator without proper null checks
- ✗ Mix reactive and blocking code
- ✗ Store secrets in configuration files

---

## References

- [Spring Boot Documentation](https://docs.spring.io/spring-boot/docs/current/reference/html/)
- [Spring Kotlin Support](https://docs.spring.io/spring-framework/reference/languages/kotlin.html)
- [Spring Security](https://docs.spring.io/spring-security/reference/index.html)
- [Spring Data JPA](https://docs.spring.io/spring-data/jpa/docs/current/reference/html/)
- [Kotlin + Spring Boot Guide](https://spring.io/guides/tutorials/spring-boot-kotlin/)
