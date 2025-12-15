# Micronaut Framework Guide

> **Framework**: Micronaut 4.x
> **Language**: Java 21+
> **Type**: Cloud-Native Microservices Framework
> **Use Cases**: Microservices, Serverless, CLI Applications, IoT

---

## Overview

Micronaut is a modern, JVM-based framework designed for building modular, easily testable microservice and serverless applications. Unlike reflection-based frameworks, Micronaut uses compile-time dependency injection and AOP, resulting in minimal memory footprint and fast startup times.

### Key Features
- **Compile-time DI**: No reflection at runtime
- **GraalVM Native**: First-class native image support
- **Reactive**: Built-in reactive programming support
- **Cloud-Native**: Service discovery, distributed tracing, config management
- **Fast Startup**: Ideal for serverless and containers
- **Low Memory**: Minimal runtime overhead

---

## Project Structure

```
myapp/
├── src/
│   ├── main/
│   │   ├── java/
│   │   │   └── com/example/
│   │   │       ├── Application.java
│   │   │       ├── controller/
│   │   │       │   └── UserController.java
│   │   │       ├── service/
│   │   │       │   ├── UserService.java
│   │   │       │   └── impl/
│   │   │       │       └── UserServiceImpl.java
│   │   │       ├── repository/
│   │   │       │   └── UserRepository.java
│   │   │       ├── domain/
│   │   │       │   └── User.java
│   │   │       ├── dto/
│   │   │       │   ├── UserRequest.java
│   │   │       │   └── UserResponse.java
│   │   │       ├── mapper/
│   │   │       │   └── UserMapper.java
│   │   │       └── exception/
│   │   │           ├── ResourceNotFoundException.java
│   │   │           └── GlobalExceptionHandler.java
│   │   └── resources/
│   │       ├── application.yml
│   │       ├── application-dev.yml
│   │       ├── application-prod.yml
│   │       └── logback.xml
│   └── test/
│       └── java/
│           └── com/example/
│               ├── controller/
│               │   └── UserControllerTest.java
│               └── service/
│                   └── UserServiceTest.java
├── build.gradle
├── gradle.properties
├── settings.gradle
└── README.md
```

---

## Dependencies (build.gradle)

```groovy
plugins {
    id("io.micronaut.application") version "4.2.1"
    id("io.micronaut.aot") version "4.2.1"
    id("com.google.devtools.ksp") version "1.9.21-1.0.16"
}

version = "0.1"
group = "com.example"

repositories {
    mavenCentral()
}

dependencies {
    // Micronaut Core
    annotationProcessor("io.micronaut:micronaut-http-validation")
    annotationProcessor("io.micronaut.serde:micronaut-serde-processor")
    annotationProcessor("io.micronaut.validation:micronaut-validation-processor")

    implementation("io.micronaut:micronaut-http-client")
    implementation("io.micronaut.serde:micronaut-serde-jackson")
    implementation("io.micronaut.validation:micronaut-validation")

    // Database
    annotationProcessor("io.micronaut.data:micronaut-data-processor")
    implementation("io.micronaut.data:micronaut-data-jdbc")
    implementation("io.micronaut.sql:micronaut-jdbc-hikari")
    implementation("io.micronaut.flyway:micronaut-flyway")
    runtimeOnly("org.postgresql:postgresql")

    // Security
    annotationProcessor("io.micronaut.security:micronaut-security-annotations")
    implementation("io.micronaut.security:micronaut-security-jwt")

    // OpenAPI
    annotationProcessor("io.micronaut.openapi:micronaut-openapi")
    implementation("io.swagger.core.v3:swagger-annotations")

    // MapStruct
    annotationProcessor("org.mapstruct:mapstruct-processor:1.5.5.Final")
    implementation("org.mapstruct:mapstruct:1.5.5.Final")

    // Lombok
    compileOnly("org.projectlombok:lombok")
    annotationProcessor("org.projectlombok:lombok")
    annotationProcessor("org.projectlombok:lombok-mapstruct-binding:0.2.0")

    // Health & Metrics
    implementation("io.micronaut:micronaut-management")
    implementation("io.micronaut.micrometer:micronaut-micrometer-core")
    implementation("io.micronaut.micrometer:micronaut-micrometer-registry-prometheus")

    // Testing
    testImplementation("io.micronaut:micronaut-http-client")
    testImplementation("io.micronaut.test:micronaut-test-junit5")
    testImplementation("org.junit.jupiter:junit-jupiter-api")
    testImplementation("org.mockito:mockito-core")
    testImplementation("org.assertj:assertj-core")
    testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine")
    testImplementation("org.testcontainers:junit-jupiter")
    testImplementation("org.testcontainers:postgresql")
}

application {
    mainClass.set("com.example.Application")
}

java {
    sourceCompatibility = JavaVersion.toVersion("21")
    targetCompatibility = JavaVersion.toVersion("21")
}

micronaut {
    runtime("netty")
    testRuntime("junit5")
    processing {
        incremental(true)
        annotations("com.example.*")
    }
    aot {
        optimizeServiceLoading = false
        convertYamlToJava = false
        precomputeOperations = true
        cacheEnvironment = true
        optimizeClassLoading = true
        deduceEnvironment = true
        optimizeNetty = true
    }
}

graalvmNative.toolchainDetection = false
```

---

## Configuration (application.yml)

```yaml
micronaut:
  application:
    name: myapp
  server:
    port: 8080
    cors:
      enabled: true
      configurations:
        web:
          allowed-origins:
            - http://localhost:3000
          allowed-methods:
            - GET
            - POST
            - PUT
            - DELETE
          allowed-headers:
            - Content-Type
            - Authorization

  security:
    authentication: bearer
    token:
      jwt:
        signatures:
          secret:
            generator:
              secret: ${JWT_SECRET:changeme-in-production-use-256-bit-key}
              jws-algorithm: HS256
        generator:
          access-token:
            expiration: 3600  # 1 hour
          refresh-token:
            enabled: true
            secret: ${JWT_REFRESH_SECRET:changeme-refresh-secret}
    intercept-url-map:
      - pattern: /health/**
        http-method: GET
        access:
          - isAnonymous()
      - pattern: /api/auth/**
        http-method: POST
        access:
          - isAnonymous()
      - pattern: /swagger/**
        access:
          - isAnonymous()
      - pattern: /api/**
        access:
          - isAuthenticated()

  router:
    static-resources:
      swagger:
        paths: classpath:META-INF/swagger
        mapping: /swagger/**

datasources:
  default:
    url: jdbc:postgresql://localhost:5432/myapp
    username: ${DB_USERNAME:postgres}
    password: ${DB_PASSWORD:postgres}
    driver-class-name: org.postgresql.Driver
    dialect: POSTGRES
    schema-generate: NONE

flyway:
  datasources:
    default:
      enabled: true
      locations: classpath:db/migration

jackson:
  serialization:
    indent-output: true
    write-dates-as-timestamps: false
  deserialization:
    fail-on-unknown-properties: false

endpoints:
  health:
    enabled: true
    sensitive: false
    details-visible: AUTHENTICATED
  info:
    enabled: true
    sensitive: false
  prometheus:
    enabled: true
    sensitive: false

logger:
  levels:
    com.example: DEBUG
    io.micronaut.data: DEBUG

---
# Development profile
micronaut:
  environments:
    - dev

datasources:
  default:
    url: jdbc:postgresql://localhost:5432/myapp_dev

---
# Production profile
micronaut:
  environments:
    - prod

logger:
  levels:
    com.example: INFO
    io.micronaut.data: WARN
```

---

## Application Entry Point

```java
package com.example;

import io.micronaut.runtime.Micronaut;
import io.swagger.v3.oas.annotations.OpenAPIDefinition;
import io.swagger.v3.oas.annotations.info.Contact;
import io.swagger.v3.oas.annotations.info.Info;
import io.swagger.v3.oas.annotations.info.License;

@OpenAPIDefinition(
    info = @Info(
        title = "My API",
        version = "1.0",
        description = "API documentation for MyApp",
        license = @License(name = "MIT"),
        contact = @Contact(name = "Support", email = "support@example.com")
    )
)
public class Application {

    public static void main(String[] args) {
        Micronaut.run(Application.class, args);
    }
}
```

---

## Domain Entity

```java
package com.example.domain;

import io.micronaut.data.annotation.*;
import io.micronaut.serde.annotation.Serdeable;
import lombok.*;

import java.time.Instant;

@Serdeable
@MappedEntity("users")
@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class User {

    @Id
    @GeneratedValue(GeneratedValue.Type.AUTO)
    private Long id;

    @Column("email")
    private String email;

    @Column("password")
    private String password;

    @Column("first_name")
    private String firstName;

    @Column("last_name")
    private String lastName;

    @Column("role")
    @Builder.Default
    private String role = "USER";

    @Column("active")
    @Builder.Default
    private Boolean active = true;

    @DateCreated
    @Column("created_at")
    private Instant createdAt;

    @DateUpdated
    @Column("updated_at")
    private Instant updatedAt;
}
```

---

## DTOs

```java
package com.example.dto;

import io.micronaut.serde.annotation.Serdeable;
import jakarta.validation.constraints.*;

@Serdeable
public record UserRequest(
    @NotBlank(message = "Email is required")
    @Email(message = "Invalid email format")
    @Size(max = 255, message = "Email must not exceed 255 characters")
    String email,

    @NotBlank(message = "Password is required")
    @Size(min = 8, max = 100, message = "Password must be between 8 and 100 characters")
    String password,

    @NotBlank(message = "First name is required")
    @Size(max = 100, message = "First name must not exceed 100 characters")
    String firstName,

    @NotBlank(message = "Last name is required")
    @Size(max = 100, message = "Last name must not exceed 100 characters")
    String lastName,

    @Pattern(regexp = "^(USER|ADMIN)$", message = "Role must be USER or ADMIN")
    String role
) {
    public UserRequest {
        if (role == null) {
            role = "USER";
        }
    }
}
```

```java
package com.example.dto;

import io.micronaut.serde.annotation.Serdeable;
import java.time.Instant;

@Serdeable
public record UserResponse(
    Long id,
    String email,
    String firstName,
    String lastName,
    String role,
    Boolean active,
    Instant createdAt,
    Instant updatedAt
) {}
```

```java
package com.example.dto;

import io.micronaut.serde.annotation.Serdeable;
import java.util.List;

@Serdeable
public record PageResponse<T>(
    List<T> content,
    int page,
    int size,
    long totalElements,
    int totalPages,
    boolean first,
    boolean last
) {
    public static <T> PageResponse<T> of(List<T> content, int page, int size, long totalElements) {
        int totalPages = (int) Math.ceil((double) totalElements / size);
        return new PageResponse<>(
            content,
            page,
            size,
            totalElements,
            totalPages,
            page == 0,
            page >= totalPages - 1
        );
    }
}
```

---

## Mapper

```java
package com.example.mapper;

import com.example.domain.User;
import com.example.dto.UserRequest;
import com.example.dto.UserResponse;
import org.mapstruct.*;

import java.util.List;

@Mapper(componentModel = "jsr330")
public interface UserMapper {

    @Mapping(target = "id", ignore = true)
    @Mapping(target = "createdAt", ignore = true)
    @Mapping(target = "updatedAt", ignore = true)
    @Mapping(target = "active", constant = "true")
    User toEntity(UserRequest request);

    UserResponse toResponse(User user);

    List<UserResponse> toResponseList(List<User> users);

    @BeanMapping(nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.IGNORE)
    @Mapping(target = "id", ignore = true)
    @Mapping(target = "createdAt", ignore = true)
    @Mapping(target = "updatedAt", ignore = true)
    @Mapping(target = "password", ignore = true)
    void updateEntity(UserRequest request, @MappingTarget User user);
}
```

---

## Repository

```java
package com.example.repository;

import com.example.domain.User;
import io.micronaut.data.annotation.Query;
import io.micronaut.data.jdbc.annotation.JdbcRepository;
import io.micronaut.data.model.Page;
import io.micronaut.data.model.Pageable;
import io.micronaut.data.model.query.builder.sql.Dialect;
import io.micronaut.data.repository.PageableRepository;

import java.util.List;
import java.util.Optional;

@JdbcRepository(dialect = Dialect.POSTGRES)
public interface UserRepository extends PageableRepository<User, Long> {

    Optional<User> findByEmail(String email);

    boolean existsByEmail(String email);

    List<User> findByActiveTrue();

    Page<User> findByRole(String role, Pageable pageable);

    @Query("SELECT u FROM User u WHERE u.active = true ORDER BY u.createdAt DESC")
    List<User> findAllActiveOrderByCreatedAtDesc();

    @Query("UPDATE User u SET u.active = :active WHERE u.id = :id")
    void updateActiveStatus(Long id, boolean active);
}
```

---

## Service Layer

```java
package com.example.service;

import com.example.dto.PageResponse;
import com.example.dto.UserRequest;
import com.example.dto.UserResponse;

import java.util.List;

public interface UserService {

    UserResponse createUser(UserRequest request);

    UserResponse getUserById(Long id);

    UserResponse getUserByEmail(String email);

    PageResponse<UserResponse> getAllUsers(int page, int size);

    List<UserResponse> getActiveUsers();

    UserResponse updateUser(Long id, UserRequest request);

    void deleteUser(Long id);

    void deactivateUser(Long id);
}
```

```java
package com.example.service.impl;

import com.example.domain.User;
import com.example.dto.PageResponse;
import com.example.dto.UserRequest;
import com.example.dto.UserResponse;
import com.example.exception.DuplicateResourceException;
import com.example.exception.ResourceNotFoundException;
import com.example.mapper.UserMapper;
import com.example.repository.UserRepository;
import com.example.service.UserService;
import io.micronaut.data.model.Pageable;
import io.micronaut.transaction.annotation.Transactional;
import jakarta.inject.Singleton;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;

import java.util.List;

@Singleton
@RequiredArgsConstructor
@Slf4j
@Transactional
public class UserServiceImpl implements UserService {

    private final UserRepository userRepository;
    private final UserMapper userMapper;
    private final BCryptPasswordEncoder passwordEncoder;

    @Override
    public UserResponse createUser(UserRequest request) {
        log.debug("Creating user with email: {}", request.email());

        if (userRepository.existsByEmail(request.email())) {
            throw new DuplicateResourceException("User", "email", request.email());
        }

        User user = userMapper.toEntity(request);
        user.setPassword(passwordEncoder.encode(request.password()));

        User savedUser = userRepository.save(user);
        log.info("Created user with id: {}", savedUser.getId());

        return userMapper.toResponse(savedUser);
    }

    @Override
    @Transactional(readOnly = true)
    public UserResponse getUserById(Long id) {
        log.debug("Fetching user by id: {}", id);

        return userRepository.findById(id)
            .map(userMapper::toResponse)
            .orElseThrow(() -> new ResourceNotFoundException("User", "id", id));
    }

    @Override
    @Transactional(readOnly = true)
    public UserResponse getUserByEmail(String email) {
        log.debug("Fetching user by email: {}", email);

        return userRepository.findByEmail(email)
            .map(userMapper::toResponse)
            .orElseThrow(() -> new ResourceNotFoundException("User", "email", email));
    }

    @Override
    @Transactional(readOnly = true)
    public PageResponse<UserResponse> getAllUsers(int page, int size) {
        log.debug("Fetching users page: {}, size: {}", page, size);

        var pageable = Pageable.from(page, size);
        var userPage = userRepository.findAll(pageable);

        List<UserResponse> content = userMapper.toResponseList(userPage.getContent());
        return PageResponse.of(content, page, size, userPage.getTotalSize());
    }

    @Override
    @Transactional(readOnly = true)
    public List<UserResponse> getActiveUsers() {
        log.debug("Fetching all active users");
        return userMapper.toResponseList(userRepository.findByActiveTrue());
    }

    @Override
    public UserResponse updateUser(Long id, UserRequest request) {
        log.debug("Updating user with id: {}", id);

        User user = userRepository.findById(id)
            .orElseThrow(() -> new ResourceNotFoundException("User", "id", id));

        // Check email uniqueness if changed
        if (!user.getEmail().equals(request.email()) &&
            userRepository.existsByEmail(request.email())) {
            throw new DuplicateResourceException("User", "email", request.email());
        }

        userMapper.updateEntity(request, user);

        User updatedUser = userRepository.update(user);
        log.info("Updated user with id: {}", id);

        return userMapper.toResponse(updatedUser);
    }

    @Override
    public void deleteUser(Long id) {
        log.debug("Deleting user with id: {}", id);

        if (!userRepository.existsById(id)) {
            throw new ResourceNotFoundException("User", "id", id);
        }

        userRepository.deleteById(id);
        log.info("Deleted user with id: {}", id);
    }

    @Override
    public void deactivateUser(Long id) {
        log.debug("Deactivating user with id: {}", id);

        if (!userRepository.existsById(id)) {
            throw new ResourceNotFoundException("User", "id", id);
        }

        userRepository.updateActiveStatus(id, false);
        log.info("Deactivated user with id: {}", id);
    }
}
```

### Password Encoder Bean

```java
package com.example.config;

import io.micronaut.context.annotation.Factory;
import jakarta.inject.Singleton;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;

@Factory
public class SecurityBeans {

    @Singleton
    public BCryptPasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder(12);
    }
}
```

---

## Controller

```java
package com.example.controller;

import com.example.dto.PageResponse;
import com.example.dto.UserRequest;
import com.example.dto.UserResponse;
import com.example.service.UserService;
import io.micronaut.http.HttpResponse;
import io.micronaut.http.HttpStatus;
import io.micronaut.http.annotation.*;
import io.micronaut.scheduling.TaskExecutors;
import io.micronaut.scheduling.annotation.ExecuteOn;
import io.micronaut.security.annotation.Secured;
import io.micronaut.security.rules.SecurityRule;
import io.micronaut.validation.Validated;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;

import java.net.URI;
import java.util.List;

@Controller("/api/users")
@Validated
@RequiredArgsConstructor
@ExecuteOn(TaskExecutors.BLOCKING)
@Secured(SecurityRule.IS_AUTHENTICATED)
@Tag(name = "Users", description = "User management endpoints")
public class UserController {

    private final UserService userService;

    @Post
    @Status(HttpStatus.CREATED)
    @Operation(summary = "Create a new user", description = "Creates a new user account")
    @ApiResponse(responseCode = "201", description = "User created successfully",
        content = @Content(schema = @Schema(implementation = UserResponse.class)))
    @ApiResponse(responseCode = "400", description = "Invalid request")
    @ApiResponse(responseCode = "409", description = "User already exists")
    public HttpResponse<UserResponse> createUser(@Body @Valid UserRequest request) {
        UserResponse user = userService.createUser(request);
        return HttpResponse.created(user)
            .headers(headers -> headers.location(URI.create("/api/users/" + user.id())));
    }

    @Get("/{id}")
    @Operation(summary = "Get user by ID", description = "Retrieves a user by their ID")
    @ApiResponse(responseCode = "200", description = "User found",
        content = @Content(schema = @Schema(implementation = UserResponse.class)))
    @ApiResponse(responseCode = "404", description = "User not found")
    public UserResponse getUserById(
            @Parameter(description = "User ID") @PathVariable Long id) {
        return userService.getUserById(id);
    }

    @Get
    @Operation(summary = "Get all users", description = "Retrieves all users with pagination")
    @ApiResponse(responseCode = "200", description = "Users retrieved successfully")
    public PageResponse<UserResponse> getAllUsers(
            @Parameter(description = "Page number") @QueryValue(defaultValue = "0") int page,
            @Parameter(description = "Page size") @QueryValue(defaultValue = "20") int size) {
        return userService.getAllUsers(page, size);
    }

    @Get("/active")
    @Operation(summary = "Get active users", description = "Retrieves all active users")
    public List<UserResponse> getActiveUsers() {
        return userService.getActiveUsers();
    }

    @Put("/{id}")
    @Operation(summary = "Update user", description = "Updates an existing user")
    @ApiResponse(responseCode = "200", description = "User updated successfully")
    @ApiResponse(responseCode = "404", description = "User not found")
    @ApiResponse(responseCode = "409", description = "Email already in use")
    public UserResponse updateUser(
            @PathVariable Long id,
            @Body @Valid UserRequest request) {
        return userService.updateUser(id, request);
    }

    @Delete("/{id}")
    @Status(HttpStatus.NO_CONTENT)
    @Secured({"ROLE_ADMIN"})
    @Operation(summary = "Delete user", description = "Deletes a user (Admin only)")
    @ApiResponse(responseCode = "204", description = "User deleted successfully")
    @ApiResponse(responseCode = "404", description = "User not found")
    public void deleteUser(@PathVariable Long id) {
        userService.deleteUser(id);
    }

    @Put("/{id}/deactivate")
    @Status(HttpStatus.NO_CONTENT)
    @Secured({"ROLE_ADMIN"})
    @Operation(summary = "Deactivate user", description = "Deactivates a user (Admin only)")
    public void deactivateUser(@PathVariable Long id) {
        userService.deactivateUser(id);
    }
}
```

---

## Exception Handling

### Custom Exceptions

```java
package com.example.exception;

public class ResourceNotFoundException extends RuntimeException {

    private final String resourceName;
    private final String fieldName;
    private final Object fieldValue;

    public ResourceNotFoundException(String resourceName, String fieldName, Object fieldValue) {
        super(String.format("%s not found with %s: '%s'", resourceName, fieldName, fieldValue));
        this.resourceName = resourceName;
        this.fieldName = fieldName;
        this.fieldValue = fieldValue;
    }

    public String getResourceName() {
        return resourceName;
    }

    public String getFieldName() {
        return fieldName;
    }

    public Object getFieldValue() {
        return fieldValue;
    }
}
```

```java
package com.example.exception;

public class DuplicateResourceException extends RuntimeException {

    private final String resourceName;
    private final String fieldName;
    private final Object fieldValue;

    public DuplicateResourceException(String resourceName, String fieldName, Object fieldValue) {
        super(String.format("%s already exists with %s: '%s'", resourceName, fieldName, fieldValue));
        this.resourceName = resourceName;
        this.fieldName = fieldName;
        this.fieldValue = fieldValue;
    }

    public String getResourceName() {
        return resourceName;
    }

    public String getFieldName() {
        return fieldName;
    }

    public Object getFieldValue() {
        return fieldValue;
    }
}
```

### Error Response DTO

```java
package com.example.dto;

import io.micronaut.serde.annotation.Serdeable;

import java.time.Instant;
import java.util.List;
import java.util.Map;

@Serdeable
public record ErrorResponse(
    String type,
    String title,
    int status,
    String detail,
    String instance,
    Instant timestamp,
    Map<String, Object> properties,
    List<FieldError> errors
) {
    @Serdeable
    public record FieldError(
        String field,
        String message,
        Object rejectedValue
    ) {}

    public static ErrorResponse of(int status, String title, String detail, String instance) {
        return new ErrorResponse(
            "about:blank",
            title,
            status,
            detail,
            instance,
            Instant.now(),
            null,
            null
        );
    }

    public static ErrorResponse withValidationErrors(
            int status, String title, String detail, String instance, List<FieldError> errors) {
        return new ErrorResponse(
            "about:blank",
            title,
            status,
            detail,
            instance,
            Instant.now(),
            null,
            errors
        );
    }
}
```

### Global Exception Handler

```java
package com.example.exception;

import com.example.dto.ErrorResponse;
import io.micronaut.context.annotation.Requires;
import io.micronaut.http.HttpRequest;
import io.micronaut.http.HttpResponse;
import io.micronaut.http.HttpStatus;
import io.micronaut.http.annotation.Produces;
import io.micronaut.http.server.exceptions.ExceptionHandler;
import jakarta.inject.Singleton;
import jakarta.validation.ConstraintViolationException;
import lombok.extern.slf4j.Slf4j;

import java.util.stream.Collectors;

@Singleton
@Slf4j
@Produces
@Requires(classes = {ResourceNotFoundException.class, ExceptionHandler.class})
public class ResourceNotFoundExceptionHandler
    implements ExceptionHandler<ResourceNotFoundException, HttpResponse<ErrorResponse>> {

    @Override
    public HttpResponse<ErrorResponse> handle(HttpRequest request, ResourceNotFoundException exception) {
        log.warn("Resource not found: {}", exception.getMessage());

        ErrorResponse error = ErrorResponse.of(
            HttpStatus.NOT_FOUND.getCode(),
            "Resource Not Found",
            exception.getMessage(),
            request.getPath()
        );

        return HttpResponse.notFound(error);
    }
}

@Singleton
@Slf4j
@Produces
@Requires(classes = {DuplicateResourceException.class, ExceptionHandler.class})
class DuplicateResourceExceptionHandler
    implements ExceptionHandler<DuplicateResourceException, HttpResponse<ErrorResponse>> {

    @Override
    public HttpResponse<ErrorResponse> handle(HttpRequest request, DuplicateResourceException exception) {
        log.warn("Duplicate resource: {}", exception.getMessage());

        ErrorResponse error = ErrorResponse.of(
            HttpStatus.CONFLICT.getCode(),
            "Resource Already Exists",
            exception.getMessage(),
            request.getPath()
        );

        return HttpResponse.status(HttpStatus.CONFLICT).body(error);
    }
}

@Singleton
@Slf4j
@Produces
@Requires(classes = {ConstraintViolationException.class, ExceptionHandler.class})
class ValidationExceptionHandler
    implements ExceptionHandler<ConstraintViolationException, HttpResponse<ErrorResponse>> {

    @Override
    public HttpResponse<ErrorResponse> handle(HttpRequest request, ConstraintViolationException exception) {
        log.warn("Validation error: {}", exception.getMessage());

        var fieldErrors = exception.getConstraintViolations().stream()
            .map(violation -> new ErrorResponse.FieldError(
                violation.getPropertyPath().toString(),
                violation.getMessage(),
                violation.getInvalidValue()
            ))
            .collect(Collectors.toList());

        ErrorResponse error = ErrorResponse.withValidationErrors(
            HttpStatus.BAD_REQUEST.getCode(),
            "Validation Failed",
            "One or more fields have validation errors",
            request.getPath(),
            fieldErrors
        );

        return HttpResponse.badRequest(error);
    }
}
```

---

## Authentication

### Auth Controller

```java
package com.example.controller;

import com.example.dto.AuthRequest;
import com.example.dto.AuthResponse;
import com.example.service.AuthService;
import io.micronaut.http.annotation.Body;
import io.micronaut.http.annotation.Controller;
import io.micronaut.http.annotation.Post;
import io.micronaut.security.annotation.Secured;
import io.micronaut.security.rules.SecurityRule;
import io.micronaut.validation.Validated;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;

@Controller("/api/auth")
@Validated
@RequiredArgsConstructor
@Secured(SecurityRule.IS_ANONYMOUS)
public class AuthController {

    private final AuthService authService;

    @Post("/login")
    public AuthResponse login(@Body @Valid AuthRequest request) {
        return authService.authenticate(request);
    }

    @Post("/refresh")
    public AuthResponse refresh(@Body RefreshRequest request) {
        return authService.refreshToken(request.refreshToken());
    }
}
```

```java
package com.example.dto;

import io.micronaut.serde.annotation.Serdeable;
import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;

@Serdeable
public record AuthRequest(
    @NotBlank @Email String email,
    @NotBlank String password
) {}

@Serdeable
public record AuthResponse(
    String accessToken,
    String refreshToken,
    String tokenType,
    Long expiresIn
) {}

@Serdeable
public record RefreshRequest(
    @NotBlank String refreshToken
) {}
```

### Auth Service

```java
package com.example.service;

import com.example.domain.User;
import com.example.dto.AuthRequest;
import com.example.dto.AuthResponse;
import com.example.exception.AuthenticationException;
import com.example.repository.UserRepository;
import io.micronaut.security.token.jwt.generator.JwtTokenGenerator;
import io.micronaut.security.token.jwt.generator.claims.JwtClaims;
import io.micronaut.security.token.jwt.generator.claims.JwtClaimsSetAdapter;
import jakarta.inject.Singleton;
import lombok.RequiredArgsConstructor;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;

import java.util.Map;
import java.util.Optional;

@Singleton
@RequiredArgsConstructor
public class AuthService {

    private final UserRepository userRepository;
    private final BCryptPasswordEncoder passwordEncoder;
    private final JwtTokenGenerator jwtTokenGenerator;

    public AuthResponse authenticate(AuthRequest request) {
        User user = userRepository.findByEmail(request.email())
            .orElseThrow(() -> new AuthenticationException("Invalid credentials"));

        if (!passwordEncoder.matches(request.password(), user.getPassword())) {
            throw new AuthenticationException("Invalid credentials");
        }

        if (!user.getActive()) {
            throw new AuthenticationException("User account is deactivated");
        }

        return generateTokens(user);
    }

    public AuthResponse refreshToken(String refreshToken) {
        // Validate refresh token and generate new tokens
        // Implementation depends on your refresh token strategy
        throw new UnsupportedOperationException("Implement refresh token logic");
    }

    private AuthResponse generateTokens(User user) {
        Map<String, Object> claims = Map.of(
            "sub", user.getEmail(),
            "roles", user.getRole(),
            "userId", user.getId()
        );

        Optional<String> accessToken = jwtTokenGenerator.generateToken(claims);

        return new AuthResponse(
            accessToken.orElseThrow(() -> new RuntimeException("Failed to generate token")),
            generateRefreshToken(user),
            "Bearer",
            3600L
        );
    }

    private String generateRefreshToken(User user) {
        // Generate refresh token - implement based on your strategy
        return "refresh-token-placeholder";
    }
}
```

---

## Database Migration (Flyway)

```sql
-- src/main/resources/db/migration/V1__create_users_table.sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'USER',
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_active ON users(active);
```

---

## Testing

### Unit Test

```java
package com.example.service;

import com.example.domain.User;
import com.example.dto.PageResponse;
import com.example.dto.UserRequest;
import com.example.dto.UserResponse;
import com.example.exception.DuplicateResourceException;
import com.example.exception.ResourceNotFoundException;
import com.example.mapper.UserMapper;
import com.example.repository.UserRepository;
import com.example.service.impl.UserServiceImpl;
import io.micronaut.data.model.Page;
import io.micronaut.data.model.Pageable;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;

import java.time.Instant;
import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
@DisplayName("UserService")
class UserServiceTest {

    @Mock
    private UserRepository userRepository;

    @Mock
    private UserMapper userMapper;

    @Mock
    private BCryptPasswordEncoder passwordEncoder;

    private UserServiceImpl userService;

    @BeforeEach
    void setUp() {
        userService = new UserServiceImpl(userRepository, userMapper, passwordEncoder);
    }

    @Nested
    @DisplayName("createUser")
    class CreateUser {

        @Test
        @DisplayName("should create user when email is unique")
        void shouldCreateUserWhenEmailIsUnique() {
            // Given
            UserRequest request = new UserRequest(
                "test@example.com", "password123", "John", "Doe", "USER");
            User user = createUser(1L, "test@example.com");
            UserResponse response = createUserResponse(1L, "test@example.com");

            when(userRepository.existsByEmail(request.email())).thenReturn(false);
            when(userMapper.toEntity(request)).thenReturn(user);
            when(passwordEncoder.encode(request.password())).thenReturn("encodedPassword");
            when(userRepository.save(user)).thenReturn(user);
            when(userMapper.toResponse(user)).thenReturn(response);

            // When
            UserResponse result = userService.createUser(request);

            // Then
            assertThat(result).isNotNull();
            assertThat(result.email()).isEqualTo("test@example.com");
            verify(userRepository).save(user);
        }

        @Test
        @DisplayName("should throw DuplicateResourceException when email exists")
        void shouldThrowWhenEmailExists() {
            // Given
            UserRequest request = new UserRequest(
                "existing@example.com", "password123", "John", "Doe", "USER");

            when(userRepository.existsByEmail(request.email())).thenReturn(true);

            // When/Then
            assertThatThrownBy(() -> userService.createUser(request))
                .isInstanceOf(DuplicateResourceException.class)
                .hasMessageContaining("email");
        }
    }

    @Nested
    @DisplayName("getUserById")
    class GetUserById {

        @Test
        @DisplayName("should return user when found")
        void shouldReturnUserWhenFound() {
            // Given
            Long userId = 1L;
            User user = createUser(userId, "test@example.com");
            UserResponse response = createUserResponse(userId, "test@example.com");

            when(userRepository.findById(userId)).thenReturn(Optional.of(user));
            when(userMapper.toResponse(user)).thenReturn(response);

            // When
            UserResponse result = userService.getUserById(userId);

            // Then
            assertThat(result).isNotNull();
            assertThat(result.id()).isEqualTo(userId);
        }

        @Test
        @DisplayName("should throw ResourceNotFoundException when not found")
        void shouldThrowWhenNotFound() {
            // Given
            Long userId = 999L;
            when(userRepository.findById(userId)).thenReturn(Optional.empty());

            // When/Then
            assertThatThrownBy(() -> userService.getUserById(userId))
                .isInstanceOf(ResourceNotFoundException.class)
                .hasMessageContaining("User");
        }
    }

    private User createUser(Long id, String email) {
        return User.builder()
            .id(id)
            .email(email)
            .password("encodedPassword")
            .firstName("John")
            .lastName("Doe")
            .role("USER")
            .active(true)
            .createdAt(Instant.now())
            .updatedAt(Instant.now())
            .build();
    }

    private UserResponse createUserResponse(Long id, String email) {
        return new UserResponse(
            id, email, "John", "Doe", "USER", true, Instant.now(), Instant.now());
    }
}
```

### Controller Test

```java
package com.example.controller;

import com.example.dto.UserRequest;
import com.example.dto.UserResponse;
import com.example.service.UserService;
import io.micronaut.http.HttpRequest;
import io.micronaut.http.HttpResponse;
import io.micronaut.http.HttpStatus;
import io.micronaut.http.client.HttpClient;
import io.micronaut.http.client.annotation.Client;
import io.micronaut.http.client.exceptions.HttpClientResponseException;
import io.micronaut.security.authentication.UsernamePasswordCredentials;
import io.micronaut.security.token.jwt.render.BearerAccessRefreshToken;
import io.micronaut.test.annotation.MockBean;
import io.micronaut.test.extensions.junit5.annotation.MicronautTest;
import jakarta.inject.Inject;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.time.Instant;

import static org.assertj.core.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@MicronautTest
@DisplayName("UserController")
class UserControllerTest {

    @Inject
    @Client("/")
    HttpClient client;

    @Inject
    UserService userService;

    private String accessToken;

    @MockBean(UserService.class)
    UserService mockUserService() {
        return mock(UserService.class);
    }

    @BeforeEach
    void setUp() {
        // Get authentication token for tests
        UsernamePasswordCredentials credentials =
            new UsernamePasswordCredentials("test@example.com", "password");
        HttpRequest<?> request = HttpRequest.POST("/api/auth/login", credentials);

        // In real tests, you would authenticate properly
        // For now, we'll skip auth in tests or use @MockBean for security
    }

    @Test
    @DisplayName("GET /api/users/{id} should return user when found")
    void getUserById_shouldReturnUser() {
        // Given
        Long userId = 1L;
        UserResponse response = new UserResponse(
            userId, "test@example.com", "John", "Doe",
            "USER", true, Instant.now(), Instant.now());

        when(userService.getUserById(userId)).thenReturn(response);

        // When
        HttpRequest<?> request = HttpRequest.GET("/api/users/" + userId);
        // Add auth header: .bearerAuth(accessToken)

        // In a real test with proper auth:
        // HttpResponse<UserResponse> httpResponse = client.toBlocking()
        //     .exchange(request, UserResponse.class);

        // Then
        // assertThat(httpResponse.status()).isEqualTo(HttpStatus.OK);
        // assertThat(httpResponse.body().id()).isEqualTo(userId);
    }

    @Test
    @DisplayName("POST /api/users should create user")
    void createUser_shouldCreateUser() {
        // Given
        UserRequest request = new UserRequest(
            "new@example.com", "password123", "Jane", "Doe", "USER");
        UserResponse response = new UserResponse(
            1L, "new@example.com", "Jane", "Doe",
            "USER", true, Instant.now(), Instant.now());

        when(userService.createUser(any(UserRequest.class))).thenReturn(response);

        // When/Then - similar pattern to above
    }
}
```

### Integration Test with Testcontainers

```java
package com.example;

import io.micronaut.http.HttpRequest;
import io.micronaut.http.client.HttpClient;
import io.micronaut.http.client.annotation.Client;
import io.micronaut.test.extensions.junit5.annotation.MicronautTest;
import io.micronaut.test.support.TestPropertyProvider;
import jakarta.inject.Inject;
import org.junit.jupiter.api.*;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import java.util.Map;

@MicronautTest
@Testcontainers
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
@DisplayName("Integration Tests")
class IntegrationTest implements TestPropertyProvider {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:15-alpine")
        .withDatabaseName("testdb")
        .withUsername("test")
        .withPassword("test");

    @Inject
    @Client("/")
    HttpClient client;

    @Override
    public Map<String, String> getProperties() {
        postgres.start();
        return Map.of(
            "datasources.default.url", postgres.getJdbcUrl(),
            "datasources.default.username", postgres.getUsername(),
            "datasources.default.password", postgres.getPassword()
        );
    }

    @Test
    @DisplayName("Health endpoint should return UP")
    void healthEndpoint_shouldReturnUp() {
        HttpRequest<?> request = HttpRequest.GET("/health");

        var response = client.toBlocking().exchange(request, String.class);

        assertThat(response.status().getCode()).isEqualTo(200);
    }
}
```

---

## Build & Run Commands

```bash
# Build
./gradlew build

# Run in development
./gradlew run

# Run with specific profile
MICRONAUT_ENVIRONMENTS=dev ./gradlew run

# Run tests
./gradlew test

# Run with coverage
./gradlew test jacocoTestReport

# Build native image (requires GraalVM)
./gradlew nativeCompile

# Run native image
./build/native/nativeCompile/myapp

# Build Docker image
./gradlew dockerBuild

# Build native Docker image
./gradlew dockerBuildNative

# Generate OpenAPI spec
./gradlew generateOpenApi

# Clean build
./gradlew clean build
```

---

## Best Practices

### Do's
- ✓ Use constructor injection (compile-time verified)
- ✓ Use `@Serdeable` for all DTOs
- ✓ Use `@ExecuteOn(TaskExecutors.BLOCKING)` for blocking operations
- ✓ Use `@Transactional` for database operations
- ✓ Use validation annotations on DTOs
- ✓ Handle exceptions with proper HTTP status codes
- ✓ Use native compilation for production
- ✓ Use Micronaut Data for repositories
- ✓ Enable health checks and metrics
- ✓ Use configuration properties with sensible defaults

### Don'ts
- ✗ Don't use reflection-based libraries without GraalVM configuration
- ✗ Don't use `@Autowired` (Spring annotation)
- ✗ Don't block in reactive streams
- ✗ Don't expose sensitive data in error messages
- ✗ Don't skip validation
- ✗ Don't use mutable DTOs
- ✗ Don't ignore compile-time warnings

---

## Comparison with Other Frameworks

| Feature | Micronaut | Spring Boot | Quarkus |
|---------|-----------|-------------|---------|
| DI | Compile-time | Runtime (reflection) | Compile-time |
| Startup Time | ~100ms | ~2-5s | ~100ms |
| Memory | ~50MB | ~200MB+ | ~50MB |
| Native Image | Excellent | Good (Spring Native) | Excellent |
| Reactive | Built-in | WebFlux (optional) | Mutiny |
| Data Access | Micronaut Data | Spring Data | Panache |
| Configuration | Type-safe | @Value, @ConfigurationProperties | Type-safe |
| Testing | Built-in mocking | Spring Test | @QuarkusTest |
| Learning Curve | Moderate | Low (familiar) | Moderate |
| Ecosystem | Growing | Extensive | Growing |

---

## References

- [Micronaut Documentation](https://docs.micronaut.io/)
- [Micronaut Guides](https://guides.micronaut.io/)
- [Micronaut Data](https://micronaut-projects.github.io/micronaut-data/latest/guide/)
- [Micronaut Security](https://micronaut-projects.github.io/micronaut-security/latest/guide/)
- [GraalVM Native Image](https://www.graalvm.org/native-image/)
