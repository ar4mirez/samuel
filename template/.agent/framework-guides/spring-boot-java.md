# Spring Boot (Java) Framework Guide

> **Framework**: Spring Boot 3.x with Java 17+
> **Type**: Full-stack Web Framework
> **Use Cases**: Enterprise applications, REST APIs, microservices, batch processing

---

## Overview

Spring Boot is the industry-standard Java framework for building production-ready applications with minimal configuration. It provides auto-configuration, embedded servers, and a comprehensive ecosystem.

**Key Features:**
- Auto-configuration and convention over configuration
- Embedded Tomcat, Jetty, or Undertow
- Production-ready features (metrics, health checks, externalized config)
- Comprehensive security with Spring Security
- Excellent testing support
- Massive ecosystem and community

---

## Project Structure

```
myproject/
├── src/
│   ├── main/
│   │   ├── java/
│   │   │   └── com/example/myproject/
│   │   │       ├── MyProjectApplication.java
│   │   │       ├── config/
│   │   │       │   ├── SecurityConfig.java
│   │   │       │   └── WebConfig.java
│   │   │       ├── controller/
│   │   │       │   └── UserController.java
│   │   │       ├── service/
│   │   │       │   ├── UserService.java
│   │   │       │   └── impl/
│   │   │       │       └── UserServiceImpl.java
│   │   │       ├── repository/
│   │   │       │   └── UserRepository.java
│   │   │       ├── model/
│   │   │       │   ├── entity/
│   │   │       │   │   └── User.java
│   │   │       │   └── dto/
│   │   │       │       ├── UserRequest.java
│   │   │       │       └── UserResponse.java
│   │   │       ├── exception/
│   │   │       │   ├── GlobalExceptionHandler.java
│   │   │       │   └── ResourceNotFoundException.java
│   │   │       └── mapper/
│   │   │           └── UserMapper.java
│   │   └── resources/
│   │       ├── application.yml
│   │       ├── application-dev.yml
│   │       ├── application-prod.yml
│   │       └── db/
│   │           └── migration/
│   │               └── V1__create_users_table.sql
│   └── test/
│       └── java/
│           └── com/example/myproject/
│               ├── controller/
│               │   └── UserControllerTest.java
│               ├── service/
│               │   └── UserServiceTest.java
│               └── integration/
│                   └── UserIntegrationTest.java
├── pom.xml
├── Dockerfile
└── README.md
```

---

## Dependencies (pom.xml)

```xml
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0
         https://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <parent>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-parent</artifactId>
        <version>3.2.0</version>
        <relativePath/>
    </parent>

    <groupId>com.example</groupId>
    <artifactId>myproject</artifactId>
    <version>0.0.1-SNAPSHOT</version>
    <name>myproject</name>
    <description>Spring Boot Application</description>

    <properties>
        <java.version>21</java.version>
        <mapstruct.version>1.5.5.Final</mapstruct.version>
    </properties>

    <dependencies>
        <!-- Web -->
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
        </dependency>

        <!-- Data JPA -->
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-data-jpa</artifactId>
        </dependency>

        <!-- Validation -->
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-validation</artifactId>
        </dependency>

        <!-- Security -->
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-security</artifactId>
        </dependency>

        <!-- Actuator (production monitoring) -->
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-actuator</artifactId>
        </dependency>

        <!-- Database -->
        <dependency>
            <groupId>org.postgresql</groupId>
            <artifactId>postgresql</artifactId>
            <scope>runtime</scope>
        </dependency>

        <!-- Flyway migrations -->
        <dependency>
            <groupId>org.flywaydb</groupId>
            <artifactId>flyway-core</artifactId>
        </dependency>

        <!-- MapStruct -->
        <dependency>
            <groupId>org.mapstruct</groupId>
            <artifactId>mapstruct</artifactId>
            <version>${mapstruct.version}</version>
        </dependency>

        <!-- Lombok -->
        <dependency>
            <groupId>org.projectlombok</groupId>
            <artifactId>lombok</artifactId>
            <optional>true</optional>
        </dependency>

        <!-- OpenAPI/Swagger -->
        <dependency>
            <groupId>org.springdoc</groupId>
            <artifactId>springdoc-openapi-starter-webmvc-ui</artifactId>
            <version>2.3.0</version>
        </dependency>

        <!-- Testing -->
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-test</artifactId>
            <scope>test</scope>
        </dependency>
        <dependency>
            <groupId>org.springframework.security</groupId>
            <artifactId>spring-security-test</artifactId>
            <scope>test</scope>
        </dependency>
        <dependency>
            <groupId>org.testcontainers</groupId>
            <artifactId>postgresql</artifactId>
            <scope>test</scope>
        </dependency>
        <dependency>
            <groupId>org.testcontainers</groupId>
            <artifactId>junit-jupiter</artifactId>
            <scope>test</scope>
        </dependency>
        <dependency>
            <groupId>com.h2database</groupId>
            <artifactId>h2</artifactId>
            <scope>test</scope>
        </dependency>
    </dependencies>

    <build>
        <plugins>
            <plugin>
                <groupId>org.springframework.boot</groupId>
                <artifactId>spring-boot-maven-plugin</artifactId>
                <configuration>
                    <excludes>
                        <exclude>
                            <groupId>org.projectlombok</groupId>
                            <artifactId>lombok</artifactId>
                        </exclude>
                    </excludes>
                </configuration>
            </plugin>
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-compiler-plugin</artifactId>
                <configuration>
                    <annotationProcessorPaths>
                        <path>
                            <groupId>org.projectlombok</groupId>
                            <artifactId>lombok</artifactId>
                            <version>${lombok.version}</version>
                        </path>
                        <path>
                            <groupId>org.mapstruct</groupId>
                            <artifactId>mapstruct-processor</artifactId>
                            <version>${mapstruct.version}</version>
                        </path>
                        <path>
                            <groupId>org.projectlombok</groupId>
                            <artifactId>lombok-mapstruct-binding</artifactId>
                            <version>0.2.0</version>
                        </path>
                    </annotationProcessorPaths>
                </configuration>
            </plugin>
        </plugins>
    </build>
</project>
```

---

## Configuration

### application.yml

```yaml
spring:
  application:
    name: myproject

  profiles:
    active: dev

  datasource:
    url: jdbc:postgresql://localhost:5432/myproject
    username: ${DB_USERNAME:postgres}
    password: ${DB_PASSWORD:postgres}
    hikari:
      maximum-pool-size: 10
      minimum-idle: 5
      connection-timeout: 30000
      idle-timeout: 600000
      max-lifetime: 1800000

  jpa:
    hibernate:
      ddl-auto: validate
    open-in-view: false
    properties:
      hibernate:
        format_sql: true
        jdbc:
          batch_size: 25
        order_inserts: true
        order_updates: true

  flyway:
    enabled: true
    locations: classpath:db/migration

server:
  port: 8080
  error:
    include-message: always
    include-binding-errors: always

management:
  endpoints:
    web:
      exposure:
        include: health,info,metrics,prometheus
  endpoint:
    health:
      show-details: when-authorized

logging:
  level:
    root: INFO
    com.example.myproject: DEBUG
    org.springframework.security: DEBUG

springdoc:
  api-docs:
    path: /api-docs
  swagger-ui:
    path: /swagger-ui.html
```

### application-dev.yml

```yaml
spring:
  jpa:
    show-sql: true

logging:
  level:
    org.hibernate.SQL: DEBUG
    org.hibernate.orm.jdbc.bind: TRACE
```

### application-prod.yml

```yaml
spring:
  jpa:
    show-sql: false

server:
  error:
    include-message: never
    include-stacktrace: never

logging:
  level:
    root: WARN
    com.example.myproject: INFO
```

---

## Application Entry Point

```java
package com.example.myproject;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class MyProjectApplication {

    public static void main(String[] args) {
        SpringApplication.run(MyProjectApplication.class, args);
    }
}
```

---

## Entity Model

```java
package com.example.myproject.model.entity;

import jakarta.persistence.*;
import lombok.*;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;

import java.time.LocalDateTime;

@Entity
@Table(name = "users")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false, unique = true)
    private String email;

    @Column(nullable = false)
    private String password;

    @Column(nullable = false)
    private String name;

    @Column(nullable = false)
    @Enumerated(EnumType.STRING)
    @Builder.Default
    private Role role = Role.USER;

    @Column(nullable = false)
    @Builder.Default
    private boolean active = true;

    @CreationTimestamp
    @Column(updatable = false)
    private LocalDateTime createdAt;

    @UpdateTimestamp
    private LocalDateTime updatedAt;

    public enum Role {
        USER, ADMIN
    }
}
```

---

## DTOs with Validation

### Request DTO

```java
package com.example.myproject.model.dto;

import jakarta.validation.constraints.*;
import lombok.Builder;

@Builder
public record UserRequest(
    @NotBlank(message = "Email is required")
    @Email(message = "Invalid email format")
    String email,

    @NotBlank(message = "Password is required")
    @Size(min = 8, max = 100, message = "Password must be between 8 and 100 characters")
    @Pattern(
        regexp = "^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])(?=.*[@#$%^&+=!]).*$",
        message = "Password must contain at least one digit, lowercase, uppercase, and special character"
    )
    String password,

    @NotBlank(message = "Name is required")
    @Size(min = 2, max = 100, message = "Name must be between 2 and 100 characters")
    String name,

    String role
) {}
```

### Response DTO

```java
package com.example.myproject.model.dto;

import lombok.Builder;
import java.time.LocalDateTime;

@Builder
public record UserResponse(
    Long id,
    String email,
    String name,
    String role,
    boolean active,
    LocalDateTime createdAt
) {}
```

---

## Mapper

```java
package com.example.myproject.mapper;

import com.example.myproject.model.dto.UserRequest;
import com.example.myproject.model.dto.UserResponse;
import com.example.myproject.model.entity.User;
import org.mapstruct.*;

import java.util.List;

@Mapper(componentModel = "spring")
public interface UserMapper {

    @Mapping(target = "id", ignore = true)
    @Mapping(target = "createdAt", ignore = true)
    @Mapping(target = "updatedAt", ignore = true)
    @Mapping(target = "active", constant = "true")
    @Mapping(target = "role", expression = "java(mapRole(request.role()))")
    User toEntity(UserRequest request);

    UserResponse toResponse(User user);

    List<UserResponse> toResponseList(List<User> users);

    @BeanMapping(nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.IGNORE)
    @Mapping(target = "id", ignore = true)
    @Mapping(target = "createdAt", ignore = true)
    @Mapping(target = "updatedAt", ignore = true)
    @Mapping(target = "password", ignore = true)
    void updateEntity(UserRequest request, @MappingTarget User user);

    default User.Role mapRole(String role) {
        if (role == null || role.isBlank()) {
            return User.Role.USER;
        }
        return User.Role.valueOf(role.toUpperCase());
    }
}
```

---

## Repository

```java
package com.example.myproject.repository;

import com.example.myproject.model.entity.User;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface UserRepository extends JpaRepository<User, Long> {

    Optional<User> findByEmail(String email);

    boolean existsByEmail(String email);

    List<User> findByActiveTrue();

    Page<User> findByRole(User.Role role, Pageable pageable);

    @Query("SELECT u FROM User u WHERE u.active = true AND u.role = :role")
    List<User> findActiveUsersByRole(@Param("role") User.Role role);

    @Query(value = "SELECT * FROM users WHERE email LIKE %:domain%", nativeQuery = true)
    List<User> findByEmailDomain(@Param("domain") String domain);
}
```

---

## Service Layer

### Service Interface

```java
package com.example.myproject.service;

import com.example.myproject.model.dto.UserRequest;
import com.example.myproject.model.dto.UserResponse;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;

import java.util.List;

public interface UserService {

    UserResponse createUser(UserRequest request);

    UserResponse getUserById(Long id);

    UserResponse getUserByEmail(String email);

    Page<UserResponse> getAllUsers(Pageable pageable);

    List<UserResponse> getActiveUsers();

    UserResponse updateUser(Long id, UserRequest request);

    void deleteUser(Long id);

    void deactivateUser(Long id);
}
```

### Service Implementation

```java
package com.example.myproject.service.impl;

import com.example.myproject.exception.ResourceNotFoundException;
import com.example.myproject.exception.DuplicateResourceException;
import com.example.myproject.mapper.UserMapper;
import com.example.myproject.model.dto.UserRequest;
import com.example.myproject.model.dto.UserResponse;
import com.example.myproject.model.entity.User;
import com.example.myproject.repository.UserRepository;
import com.example.myproject.service.UserService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Service
@RequiredArgsConstructor
@Slf4j
@Transactional(readOnly = true)
public class UserServiceImpl implements UserService {

    private final UserRepository userRepository;
    private final UserMapper userMapper;
    private final PasswordEncoder passwordEncoder;

    @Override
    @Transactional
    public UserResponse createUser(UserRequest request) {
        log.info("Creating user with email: {}", request.email());

        if (userRepository.existsByEmail(request.email())) {
            throw new DuplicateResourceException("User", "email", request.email());
        }

        User user = userMapper.toEntity(request);
        user.setPassword(passwordEncoder.encode(request.password()));

        User savedUser = userRepository.save(user);
        log.info("User created with id: {}", savedUser.getId());

        return userMapper.toResponse(savedUser);
    }

    @Override
    public UserResponse getUserById(Long id) {
        log.debug("Fetching user by id: {}", id);

        return userRepository.findById(id)
            .map(userMapper::toResponse)
            .orElseThrow(() -> new ResourceNotFoundException("User", "id", id));
    }

    @Override
    public UserResponse getUserByEmail(String email) {
        log.debug("Fetching user by email: {}", email);

        return userRepository.findByEmail(email)
            .map(userMapper::toResponse)
            .orElseThrow(() -> new ResourceNotFoundException("User", "email", email));
    }

    @Override
    public Page<UserResponse> getAllUsers(Pageable pageable) {
        log.debug("Fetching all users with pagination");

        return userRepository.findAll(pageable)
            .map(userMapper::toResponse);
    }

    @Override
    public List<UserResponse> getActiveUsers() {
        log.debug("Fetching all active users");

        return userMapper.toResponseList(userRepository.findByActiveTrue());
    }

    @Override
    @Transactional
    public UserResponse updateUser(Long id, UserRequest request) {
        log.info("Updating user with id: {}", id);

        User user = userRepository.findById(id)
            .orElseThrow(() -> new ResourceNotFoundException("User", "id", id));

        // Check email uniqueness if changed
        if (!user.getEmail().equals(request.email()) &&
            userRepository.existsByEmail(request.email())) {
            throw new DuplicateResourceException("User", "email", request.email());
        }

        userMapper.updateEntity(request, user);
        User updatedUser = userRepository.save(user);

        return userMapper.toResponse(updatedUser);
    }

    @Override
    @Transactional
    public void deleteUser(Long id) {
        log.info("Deleting user with id: {}", id);

        if (!userRepository.existsById(id)) {
            throw new ResourceNotFoundException("User", "id", id);
        }

        userRepository.deleteById(id);
    }

    @Override
    @Transactional
    public void deactivateUser(Long id) {
        log.info("Deactivating user with id: {}", id);

        User user = userRepository.findById(id)
            .orElseThrow(() -> new ResourceNotFoundException("User", "id", id));

        user.setActive(false);
        userRepository.save(user);
    }
}
```

---

## Controller

```java
package com.example.myproject.controller;

import com.example.myproject.model.dto.UserRequest;
import com.example.myproject.model.dto.UserResponse;
import com.example.myproject.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.web.PageableDefault;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/v1/users")
@RequiredArgsConstructor
@Tag(name = "Users", description = "User management APIs")
public class UserController {

    private final UserService userService;

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    @Operation(summary = "Create a new user")
    @ApiResponse(responseCode = "201", description = "User created successfully")
    @ApiResponse(responseCode = "400", description = "Invalid input")
    @ApiResponse(responseCode = "409", description = "Email already exists")
    public ResponseEntity<UserResponse> createUser(@Valid @RequestBody UserRequest request) {
        UserResponse response = userService.createUser(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }

    @GetMapping("/{id}")
    @Operation(summary = "Get user by ID")
    @ApiResponse(responseCode = "200", description = "User found")
    @ApiResponse(responseCode = "404", description = "User not found")
    public ResponseEntity<UserResponse> getUserById(
            @Parameter(description = "User ID") @PathVariable Long id) {
        return ResponseEntity.ok(userService.getUserById(id));
    }

    @GetMapping
    @Operation(summary = "Get all users with pagination")
    public ResponseEntity<Page<UserResponse>> getAllUsers(
            @PageableDefault(size = 20, sort = "id") Pageable pageable) {
        return ResponseEntity.ok(userService.getAllUsers(pageable));
    }

    @GetMapping("/active")
    @Operation(summary = "Get all active users")
    public ResponseEntity<List<UserResponse>> getActiveUsers() {
        return ResponseEntity.ok(userService.getActiveUsers());
    }

    @PutMapping("/{id}")
    @Operation(summary = "Update user")
    @ApiResponse(responseCode = "200", description = "User updated successfully")
    @ApiResponse(responseCode = "404", description = "User not found")
    public ResponseEntity<UserResponse> updateUser(
            @PathVariable Long id,
            @Valid @RequestBody UserRequest request) {
        return ResponseEntity.ok(userService.updateUser(id, request));
    }

    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    @PreAuthorize("hasRole('ADMIN')")
    @Operation(summary = "Delete user (Admin only)")
    @ApiResponse(responseCode = "204", description = "User deleted successfully")
    @ApiResponse(responseCode = "404", description = "User not found")
    public ResponseEntity<Void> deleteUser(@PathVariable Long id) {
        userService.deleteUser(id);
        return ResponseEntity.noContent().build();
    }

    @PatchMapping("/{id}/deactivate")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    @Operation(summary = "Deactivate user")
    public ResponseEntity<Void> deactivateUser(@PathVariable Long id) {
        userService.deactivateUser(id);
        return ResponseEntity.noContent().build();
    }
}
```

---

## Exception Handling

### Custom Exceptions

```java
package com.example.myproject.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.NOT_FOUND)
public class ResourceNotFoundException extends RuntimeException {

    public ResourceNotFoundException(String resource, String field, Object value) {
        super(String.format("%s not found with %s: '%s'", resource, field, value));
    }
}
```

```java
package com.example.myproject.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.CONFLICT)
public class DuplicateResourceException extends RuntimeException {

    public DuplicateResourceException(String resource, String field, Object value) {
        super(String.format("%s already exists with %s: '%s'", resource, field, value));
    }
}
```

### Global Exception Handler

```java
package com.example.myproject.exception;

import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ProblemDetail;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import java.net.URI;
import java.time.Instant;
import java.util.HashMap;
import java.util.Map;

@RestControllerAdvice
@Slf4j
public class GlobalExceptionHandler {

    @ExceptionHandler(ResourceNotFoundException.class)
    public ProblemDetail handleResourceNotFound(ResourceNotFoundException ex) {
        log.warn("Resource not found: {}", ex.getMessage());

        ProblemDetail problem = ProblemDetail.forStatusAndDetail(
            HttpStatus.NOT_FOUND,
            ex.getMessage()
        );
        problem.setTitle("Resource Not Found");
        problem.setType(URI.create("https://api.example.com/errors/not-found"));
        problem.setProperty("timestamp", Instant.now());

        return problem;
    }

    @ExceptionHandler(DuplicateResourceException.class)
    public ProblemDetail handleDuplicateResource(DuplicateResourceException ex) {
        log.warn("Duplicate resource: {}", ex.getMessage());

        ProblemDetail problem = ProblemDetail.forStatusAndDetail(
            HttpStatus.CONFLICT,
            ex.getMessage()
        );
        problem.setTitle("Duplicate Resource");
        problem.setType(URI.create("https://api.example.com/errors/conflict"));
        problem.setProperty("timestamp", Instant.now());

        return problem;
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ProblemDetail handleValidationErrors(MethodArgumentNotValidException ex) {
        log.warn("Validation failed: {}", ex.getMessage());

        Map<String, String> errors = new HashMap<>();
        ex.getBindingResult().getAllErrors().forEach(error -> {
            String fieldName = ((FieldError) error).getField();
            String errorMessage = error.getDefaultMessage();
            errors.put(fieldName, errorMessage);
        });

        ProblemDetail problem = ProblemDetail.forStatusAndDetail(
            HttpStatus.BAD_REQUEST,
            "Validation failed"
        );
        problem.setTitle("Validation Error");
        problem.setType(URI.create("https://api.example.com/errors/validation"));
        problem.setProperty("timestamp", Instant.now());
        problem.setProperty("errors", errors);

        return problem;
    }

    @ExceptionHandler(Exception.class)
    public ProblemDetail handleGenericException(Exception ex) {
        log.error("Unexpected error occurred", ex);

        ProblemDetail problem = ProblemDetail.forStatusAndDetail(
            HttpStatus.INTERNAL_SERVER_ERROR,
            "An unexpected error occurred"
        );
        problem.setTitle("Internal Server Error");
        problem.setType(URI.create("https://api.example.com/errors/internal"));
        problem.setProperty("timestamp", Instant.now());

        return problem;
    }
}
```

---

## Security Configuration

```java
package com.example.myproject.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.annotation.method.configuration.EnableMethodSecurity;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.web.SecurityFilterChain;

@Configuration
@EnableWebSecurity
@EnableMethodSecurity
public class SecurityConfig {

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        return http
            .csrf(AbstractHttpConfigurer::disable)
            .sessionManagement(session ->
                session.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/api-docs/**", "/swagger-ui/**", "/swagger-ui.html").permitAll()
                .requestMatchers("/actuator/health", "/actuator/info").permitAll()
                .requestMatchers(HttpMethod.POST, "/api/v1/users").permitAll()
                .requestMatchers(HttpMethod.GET, "/api/v1/users/**").authenticated()
                .requestMatchers(HttpMethod.DELETE, "/api/v1/users/**").hasRole("ADMIN")
                .anyRequest().authenticated()
            )
            .httpBasic(httpBasic -> {})
            .build();
    }

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder(12);
    }
}
```

---

## Database Migration (Flyway)

### V1__create_users_table.sql

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'USER',
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_active ON users(active);
```

---

## Testing

### Unit Tests

```java
package com.example.myproject.service;

import com.example.myproject.exception.DuplicateResourceException;
import com.example.myproject.exception.ResourceNotFoundException;
import com.example.myproject.mapper.UserMapper;
import com.example.myproject.model.dto.UserRequest;
import com.example.myproject.model.dto.UserResponse;
import com.example.myproject.model.entity.User;
import com.example.myproject.repository.UserRepository;
import com.example.myproject.service.impl.UserServiceImpl;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.security.crypto.password.PasswordEncoder;

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
    private PasswordEncoder passwordEncoder;

    @InjectMocks
    private UserServiceImpl userService;

    private User user;
    private UserRequest request;
    private UserResponse response;

    @BeforeEach
    void setUp() {
        user = User.builder()
            .id(1L)
            .email("test@example.com")
            .password("encodedPassword")
            .name("Test User")
            .role(User.Role.USER)
            .active(true)
            .build();

        request = UserRequest.builder()
            .email("test@example.com")
            .password("Password123!")
            .name("Test User")
            .build();

        response = UserResponse.builder()
            .id(1L)
            .email("test@example.com")
            .name("Test User")
            .role("USER")
            .active(true)
            .build();
    }

    @Nested
    @DisplayName("createUser")
    class CreateUser {

        @Test
        @DisplayName("should create user with valid data")
        void shouldCreateUserWithValidData() {
            // Given
            when(userRepository.existsByEmail(request.email())).thenReturn(false);
            when(userMapper.toEntity(request)).thenReturn(user);
            when(passwordEncoder.encode(request.password())).thenReturn("encodedPassword");
            when(userRepository.save(any(User.class))).thenReturn(user);
            when(userMapper.toResponse(user)).thenReturn(response);

            // When
            UserResponse result = userService.createUser(request);

            // Then
            assertThat(result).isNotNull();
            assertThat(result.email()).isEqualTo("test@example.com");
            verify(userRepository).save(any(User.class));
        }

        @Test
        @DisplayName("should throw exception when email already exists")
        void shouldThrowExceptionWhenEmailExists() {
            // Given
            when(userRepository.existsByEmail(request.email())).thenReturn(true);

            // When/Then
            assertThatThrownBy(() -> userService.createUser(request))
                .isInstanceOf(DuplicateResourceException.class)
                .hasMessageContaining("email");

            verify(userRepository, never()).save(any());
        }
    }

    @Nested
    @DisplayName("getUserById")
    class GetUserById {

        @Test
        @DisplayName("should return user when found")
        void shouldReturnUserWhenFound() {
            // Given
            when(userRepository.findById(1L)).thenReturn(Optional.of(user));
            when(userMapper.toResponse(user)).thenReturn(response);

            // When
            UserResponse result = userService.getUserById(1L);

            // Then
            assertThat(result).isNotNull();
            assertThat(result.id()).isEqualTo(1L);
        }

        @Test
        @DisplayName("should throw exception when not found")
        void shouldThrowExceptionWhenNotFound() {
            // Given
            when(userRepository.findById(999L)).thenReturn(Optional.empty());

            // When/Then
            assertThatThrownBy(() -> userService.getUserById(999L))
                .isInstanceOf(ResourceNotFoundException.class)
                .hasMessageContaining("999");
        }
    }
}
```

### Integration Tests

```java
package com.example.myproject.integration;

import com.example.myproject.model.dto.UserRequest;
import com.example.myproject.model.dto.UserResponse;
import com.example.myproject.model.entity.User;
import com.example.myproject.repository.UserRepository;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.security.test.context.support.WithMockUser;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.springframework.test.web.servlet.MockMvc;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import static org.hamcrest.Matchers.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@SpringBootTest
@AutoConfigureMockMvc
@Testcontainers
@DisplayName("User API Integration Tests")
class UserIntegrationTest {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:15-alpine");

    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
    }

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ObjectMapper objectMapper;

    @Autowired
    private UserRepository userRepository;

    @BeforeEach
    void setUp() {
        userRepository.deleteAll();
    }

    @Test
    @DisplayName("POST /api/v1/users - should create user")
    void shouldCreateUser() throws Exception {
        UserRequest request = UserRequest.builder()
            .email("test@example.com")
            .password("Password123!")
            .name("Test User")
            .build();

        mockMvc.perform(post("/api/v1/users")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(request)))
            .andExpect(status().isCreated())
            .andExpect(jsonPath("$.id").exists())
            .andExpect(jsonPath("$.email").value("test@example.com"))
            .andExpect(jsonPath("$.name").value("Test User"));
    }

    @Test
    @DisplayName("POST /api/v1/users - should return 400 for invalid email")
    void shouldReturnBadRequestForInvalidEmail() throws Exception {
        UserRequest request = UserRequest.builder()
            .email("invalid-email")
            .password("Password123!")
            .name("Test User")
            .build();

        mockMvc.perform(post("/api/v1/users")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(request)))
            .andExpect(status().isBadRequest())
            .andExpect(jsonPath("$.errors.email").exists());
    }

    @Test
    @WithMockUser
    @DisplayName("GET /api/v1/users/{id} - should return user")
    void shouldReturnUser() throws Exception {
        // Create a user first
        User user = userRepository.save(User.builder()
            .email("test@example.com")
            .password("encoded")
            .name("Test User")
            .build());

        mockMvc.perform(get("/api/v1/users/" + user.getId()))
            .andExpect(status().isOk())
            .andExpect(jsonPath("$.id").value(user.getId()))
            .andExpect(jsonPath("$.email").value("test@example.com"));
    }

    @Test
    @WithMockUser
    @DisplayName("GET /api/v1/users/{id} - should return 404 for non-existent user")
    void shouldReturnNotFoundForNonExistentUser() throws Exception {
        mockMvc.perform(get("/api/v1/users/999"))
            .andExpect(status().isNotFound())
            .andExpect(jsonPath("$.title").value("Resource Not Found"));
    }

    @Test
    @WithMockUser(roles = "ADMIN")
    @DisplayName("DELETE /api/v1/users/{id} - admin should delete user")
    void adminShouldDeleteUser() throws Exception {
        User user = userRepository.save(User.builder()
            .email("test@example.com")
            .password("encoded")
            .name("Test User")
            .build());

        mockMvc.perform(delete("/api/v1/users/" + user.getId()))
            .andExpect(status().isNoContent());

        assertThat(userRepository.findById(user.getId())).isEmpty();
    }

    @Test
    @WithMockUser(roles = "USER")
    @DisplayName("DELETE /api/v1/users/{id} - regular user should get forbidden")
    void regularUserShouldNotDeleteUser() throws Exception {
        User user = userRepository.save(User.builder()
            .email("test@example.com")
            .password("encoded")
            .name("Test User")
            .build());

        mockMvc.perform(delete("/api/v1/users/" + user.getId()))
            .andExpect(status().isForbidden());
    }
}
```

---

## Commands

```bash
# Create project
mvn archetype:generate -DgroupId=com.example -DartifactId=myproject \
    -DarchetypeArtifactId=maven-archetype-quickstart -DinteractiveMode=false

# Or use Spring Initializr
# https://start.spring.io/

# Build
./mvnw clean package

# Run
./mvnw spring-boot:run

# Run with profile
./mvnw spring-boot:run -Dspring-boot.run.profiles=dev

# Test
./mvnw test

# Test with coverage
./mvnw verify

# Format (with spotless)
./mvnw spotless:apply

# Static analysis (with checkstyle)
./mvnw checkstyle:check

# Build Docker image
./mvnw spring-boot:build-image

# Generate native executable (GraalVM)
./mvnw -Pnative native:compile
```

---

## Best Practices

### Do's
- ✓ Use constructor injection (with `@RequiredArgsConstructor`)
- ✓ Use `@Transactional(readOnly = true)` for read operations
- ✓ Use DTOs for API requests/responses (never expose entities)
- ✓ Use MapStruct for entity-DTO mapping
- ✓ Use pagination for list endpoints
- ✓ Use ProblemDetail (RFC 7807) for error responses
- ✓ Use profiles for environment-specific configuration
- ✓ Use Flyway/Liquibase for database migrations
- ✓ Write unit tests for services, integration tests for controllers
- ✓ Use Testcontainers for realistic integration tests

### Don'ts
- ✗ Don't use `@Autowired` on fields (use constructor injection)
- ✗ Don't expose JPA entities directly in APIs
- ✗ Don't use `spring.jpa.hibernate.ddl-auto=update` in production
- ✗ Don't catch generic `Exception` in controllers
- ✗ Don't put business logic in controllers
- ✗ Don't use `@Transactional` at class level carelessly
- ✗ Don't hardcode configuration values

---

## Comparison: Spring Boot (Java) vs Spring Boot (Kotlin)

| Aspect | Java | Kotlin |
|--------|------|--------|
| Null Safety | Optional, annotations | Built-in null types |
| Data Classes | Records (Java 16+) | data class |
| Boilerplate | More (use Lombok) | Less (built-in) |
| Coroutines | Project Reactor | Native coroutines |
| DSL Support | Limited | Excellent (Router DSL) |
| Learning Curve | Standard | Additional Kotlin learning |
| Ecosystem | Mature | Growing |

---

## References

- [Spring Boot Documentation](https://docs.spring.io/spring-boot/docs/current/reference/html/)
- [Spring Data JPA](https://docs.spring.io/spring-data/jpa/docs/current/reference/html/)
- [Spring Security](https://docs.spring.io/spring-security/reference/)
- [Baeldung Spring Tutorials](https://www.baeldung.com/spring-tutorial)
- [MapStruct Documentation](https://mapstruct.org/documentation/stable/reference/html/)
