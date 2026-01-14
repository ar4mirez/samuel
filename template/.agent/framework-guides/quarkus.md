# Quarkus Framework Guide

> **Framework**: Quarkus 3.x with Java 17+
> **Type**: Supersonic Subatomic Java Framework
> **Use Cases**: Cloud-native microservices, serverless, Kubernetes, GraalVM native

---

## Overview

Quarkus is a Kubernetes-native Java framework designed for GraalVM and HotSpot, providing fast startup times, low memory footprint, and developer joy through live reload.

**Key Features:**
- Native compilation with GraalVM (millisecond startup)
- Live coding with Dev Services
- Kubernetes-native design
- Unified reactive and imperative programming
- Extensive extension ecosystem
- Dev UI for development productivity

---

## Project Structure

```
myproject/
├── src/
│   ├── main/
│   │   ├── java/
│   │   │   └── com/example/
│   │   │       ├── resource/
│   │   │       │   └── UserResource.java
│   │   │       ├── service/
│   │   │       │   └── UserService.java
│   │   │       ├── repository/
│   │   │       │   └── UserRepository.java
│   │   │       ├── entity/
│   │   │       │   └── User.java
│   │   │       ├── dto/
│   │   │       │   ├── UserRequest.java
│   │   │       │   └── UserResponse.java
│   │   │       ├── mapper/
│   │   │       │   └── UserMapper.java
│   │   │       └── exception/
│   │   │           └── ExceptionMappers.java
│   │   ├── resources/
│   │   │   ├── application.properties
│   │   │   └── db/
│   │   │       └── migration/
│   │   │           └── V1__create_users.sql
│   │   └── docker/
│   │       ├── Dockerfile.jvm
│   │       └── Dockerfile.native
│   └── test/
│       └── java/
│           └── com/example/
│               ├── resource/
│               │   └── UserResourceTest.java
│               └── service/
│                   └── UserServiceTest.java
├── pom.xml
└── README.md
```

---

## Dependencies (pom.xml)

```xml
<?xml version="1.0"?>
<project xsi:schemaLocation="http://maven.apache.org/POM/4.0.0
         https://maven.apache.org/xsd/maven-4.0.0.xsd"
         xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
    <modelVersion>4.0.0</modelVersion>

    <groupId>com.example</groupId>
    <artifactId>myproject</artifactId>
    <version>1.0.0-SNAPSHOT</version>

    <properties>
        <compiler-plugin.version>3.11.0</compiler-plugin.version>
        <maven.compiler.release>21</maven.compiler.release>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
        <quarkus.platform.artifact-id>quarkus-bom</quarkus.platform.artifact-id>
        <quarkus.platform.group-id>io.quarkus.platform</quarkus.platform.group-id>
        <quarkus.platform.version>3.6.0</quarkus.platform.version>
        <surefire-plugin.version>3.1.2</surefire-plugin.version>
    </properties>

    <dependencyManagement>
        <dependencies>
            <dependency>
                <groupId>${quarkus.platform.group-id}</groupId>
                <artifactId>${quarkus.platform.artifact-id}</artifactId>
                <version>${quarkus.platform.version}</version>
                <type>pom</type>
                <scope>import</scope>
            </dependency>
        </dependencies>
    </dependencyManagement>

    <dependencies>
        <!-- RESTEasy Reactive -->
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-resteasy-reactive</artifactId>
        </dependency>
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-resteasy-reactive-jackson</artifactId>
        </dependency>

        <!-- Hibernate Reactive with Panache -->
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-hibernate-reactive-panache</artifactId>
        </dependency>
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-reactive-pg-client</artifactId>
        </dependency>

        <!-- Validation -->
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-hibernate-validator</artifactId>
        </dependency>

        <!-- Security -->
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-smallrye-jwt</artifactId>
        </dependency>
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-smallrye-jwt-build</artifactId>
        </dependency>

        <!-- Health & Metrics -->
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-smallrye-health</artifactId>
        </dependency>
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-micrometer-registry-prometheus</artifactId>
        </dependency>

        <!-- OpenAPI -->
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-smallrye-openapi</artifactId>
        </dependency>

        <!-- Flyway -->
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-flyway</artifactId>
        </dependency>
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-jdbc-postgresql</artifactId>
        </dependency>

        <!-- MapStruct -->
        <dependency>
            <groupId>org.mapstruct</groupId>
            <artifactId>mapstruct</artifactId>
            <version>1.5.5.Final</version>
        </dependency>

        <!-- Testing -->
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-junit5</artifactId>
            <scope>test</scope>
        </dependency>
        <dependency>
            <groupId>io.rest-assured</groupId>
            <artifactId>rest-assured</artifactId>
            <scope>test</scope>
        </dependency>
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-test-security</artifactId>
            <scope>test</scope>
        </dependency>
        <dependency>
            <groupId>io.quarkus</groupId>
            <artifactId>quarkus-junit5-mockito</artifactId>
            <scope>test</scope>
        </dependency>
        <dependency>
            <groupId>org.testcontainers</groupId>
            <artifactId>postgresql</artifactId>
            <scope>test</scope>
        </dependency>
    </dependencies>

    <build>
        <plugins>
            <plugin>
                <groupId>${quarkus.platform.group-id}</groupId>
                <artifactId>quarkus-maven-plugin</artifactId>
                <version>${quarkus.platform.version}</version>
                <extensions>true</extensions>
                <executions>
                    <execution>
                        <goals>
                            <goal>build</goal>
                            <goal>generate-code</goal>
                            <goal>generate-code-tests</goal>
                        </goals>
                    </execution>
                </executions>
            </plugin>
            <plugin>
                <artifactId>maven-compiler-plugin</artifactId>
                <version>${compiler-plugin.version}</version>
                <configuration>
                    <annotationProcessorPaths>
                        <path>
                            <groupId>org.mapstruct</groupId>
                            <artifactId>mapstruct-processor</artifactId>
                            <version>1.5.5.Final</version>
                        </path>
                    </annotationProcessorPaths>
                </configuration>
            </plugin>
            <plugin>
                <artifactId>maven-surefire-plugin</artifactId>
                <version>${surefire-plugin.version}</version>
                <configuration>
                    <systemPropertyVariables>
                        <java.util.logging.manager>
                            org.jboss.logmanager.LogManager
                        </java.util.logging.manager>
                    </systemPropertyVariables>
                </configuration>
            </plugin>
        </plugins>
    </build>

    <profiles>
        <profile>
            <id>native</id>
            <activation>
                <property>
                    <name>native</name>
                </property>
            </activation>
            <properties>
                <quarkus.package.type>native</quarkus.package.type>
            </properties>
        </profile>
    </profiles>
</project>
```

---

## Configuration

### application.properties

```properties
# Application
quarkus.application.name=myproject
quarkus.http.port=8080

# Database - Reactive
quarkus.datasource.db-kind=postgresql
quarkus.datasource.username=${DB_USERNAME:postgres}
quarkus.datasource.password=${DB_PASSWORD:postgres}
quarkus.datasource.reactive.url=vertx-reactive:postgresql://localhost:5432/myproject

# Database - JDBC (for Flyway)
quarkus.datasource.jdbc.url=jdbc:postgresql://localhost:5432/myproject

# Hibernate
quarkus.hibernate-orm.database.generation=validate
quarkus.hibernate-orm.log.sql=true

# Flyway
quarkus.flyway.migrate-at-start=true
quarkus.flyway.locations=db/migration

# JWT Security
mp.jwt.verify.publickey.location=publicKey.pem
mp.jwt.verify.issuer=https://example.com
smallrye.jwt.sign.key.location=privateKey.pem

# OpenAPI
quarkus.smallrye-openapi.path=/api-docs
quarkus.swagger-ui.always-include=true
quarkus.swagger-ui.path=/swagger-ui

# Health
quarkus.smallrye-health.root-path=/health

# Dev Services (auto-starts containers in dev mode)
quarkus.devservices.enabled=true
%dev.quarkus.datasource.devservices.enabled=true
%dev.quarkus.datasource.devservices.image-name=postgres:15-alpine

# Logging
quarkus.log.level=INFO
quarkus.log.category."com.example".level=DEBUG

# Production profile
%prod.quarkus.datasource.reactive.url=vertx-reactive:postgresql://${DB_HOST}:${DB_PORT}/${DB_NAME}
%prod.quarkus.log.console.json=true
```

---

## Entity with Panache

```java
package com.example.entity;

import io.quarkus.hibernate.reactive.panache.PanacheEntity;
import io.smallrye.mutiny.Uni;
import jakarta.persistence.*;
import java.time.LocalDateTime;
import java.util.List;

@Entity
@Table(name = "users")
public class User extends PanacheEntity {

    @Column(nullable = false, unique = true)
    public String email;

    @Column(nullable = false)
    public String password;

    @Column(nullable = false)
    public String name;

    @Column(nullable = false)
    @Enumerated(EnumType.STRING)
    public Role role = Role.USER;

    @Column(nullable = false)
    public boolean active = true;

    @Column(name = "created_at", updatable = false)
    public LocalDateTime createdAt;

    @Column(name = "updated_at")
    public LocalDateTime updatedAt;

    @PrePersist
    void onCreate() {
        createdAt = LocalDateTime.now();
        updatedAt = LocalDateTime.now();
    }

    @PreUpdate
    void onUpdate() {
        updatedAt = LocalDateTime.now();
    }

    public enum Role {
        USER, ADMIN
    }

    // Panache query methods
    public static Uni<User> findByEmail(String email) {
        return find("email", email).firstResult();
    }

    public static Uni<List<User>> findActive() {
        return list("active", true);
    }

    public static Uni<List<User>> findByRole(Role role) {
        return list("role", role);
    }

    public static Uni<Boolean> existsByEmail(String email) {
        return count("email", email).map(count -> count > 0);
    }
}
```

---

## DTOs with Validation

### Request DTO

```java
package com.example.dto;

import jakarta.validation.constraints.*;

public record UserRequest(
    @NotBlank(message = "Email is required")
    @Email(message = "Invalid email format")
    String email,

    @NotBlank(message = "Password is required")
    @Size(min = 8, max = 100, message = "Password must be between 8 and 100 characters")
    String password,

    @NotBlank(message = "Name is required")
    @Size(min = 2, max = 100, message = "Name must be between 2 and 100 characters")
    String name,

    String role
) {}
```

### Response DTO

```java
package com.example.dto;

import java.time.LocalDateTime;

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
package com.example.mapper;

import com.example.dto.UserRequest;
import com.example.dto.UserResponse;
import com.example.entity.User;
import org.mapstruct.*;

import java.util.List;

@Mapper(componentModel = "cdi")
public interface UserMapper {

    @Mapping(target = "id", ignore = true)
    @Mapping(target = "createdAt", ignore = true)
    @Mapping(target = "updatedAt", ignore = true)
    @Mapping(target = "active", constant = "true")
    @Mapping(target = "role", expression = "java(mapRole(request.role()))")
    User toEntity(UserRequest request);

    @Mapping(target = "role", source = "role")
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

    default String mapRole(User.Role role) {
        return role.name();
    }
}
```

---

## Repository (Optional with Panache)

```java
package com.example.repository;

import com.example.entity.User;
import io.quarkus.hibernate.reactive.panache.PanacheRepository;
import io.smallrye.mutiny.Uni;
import jakarta.enterprise.context.ApplicationScoped;

import java.util.List;

@ApplicationScoped
public class UserRepository implements PanacheRepository<User> {

    public Uni<User> findByEmail(String email) {
        return find("email", email).firstResult();
    }

    public Uni<List<User>> findActive() {
        return list("active", true);
    }

    public Uni<List<User>> findByRole(User.Role role) {
        return list("role", role);
    }

    public Uni<Boolean> existsByEmail(String email) {
        return count("email", email).map(count -> count > 0);
    }
}
```

---

## Service Layer

```java
package com.example.service;

import com.example.dto.UserRequest;
import com.example.dto.UserResponse;
import com.example.entity.User;
import com.example.exception.DuplicateResourceException;
import com.example.exception.ResourceNotFoundException;
import com.example.mapper.UserMapper;
import com.example.repository.UserRepository;
import io.quarkus.elytron.security.common.BcryptUtil;
import io.quarkus.hibernate.reactive.panache.common.WithTransaction;
import io.smallrye.mutiny.Uni;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import org.jboss.logging.Logger;

import java.util.List;

@ApplicationScoped
public class UserService {

    private static final Logger LOG = Logger.getLogger(UserService.class);

    @Inject
    UserRepository userRepository;

    @Inject
    UserMapper userMapper;

    @WithTransaction
    public Uni<UserResponse> createUser(UserRequest request) {
        LOG.infof("Creating user with email: %s", request.email());

        return userRepository.existsByEmail(request.email())
            .flatMap(exists -> {
                if (exists) {
                    return Uni.createFrom().failure(
                        new DuplicateResourceException("User", "email", request.email())
                    );
                }

                User user = userMapper.toEntity(request);
                user.password = BcryptUtil.bcryptHash(request.password());

                return userRepository.persist(user)
                    .map(userMapper::toResponse);
            });
    }

    public Uni<UserResponse> getUserById(Long id) {
        LOG.debugf("Fetching user by id: %d", id);

        return userRepository.findById(id)
            .onItem().ifNull().failWith(() ->
                new ResourceNotFoundException("User", "id", id))
            .map(userMapper::toResponse);
    }

    public Uni<UserResponse> getUserByEmail(String email) {
        LOG.debugf("Fetching user by email: %s", email);

        return userRepository.findByEmail(email)
            .onItem().ifNull().failWith(() ->
                new ResourceNotFoundException("User", "email", email))
            .map(userMapper::toResponse);
    }

    public Uni<List<UserResponse>> getAllUsers() {
        LOG.debug("Fetching all users");

        return userRepository.listAll()
            .map(userMapper::toResponseList);
    }

    public Uni<List<UserResponse>> getActiveUsers() {
        LOG.debug("Fetching active users");

        return userRepository.findActive()
            .map(userMapper::toResponseList);
    }

    @WithTransaction
    public Uni<UserResponse> updateUser(Long id, UserRequest request) {
        LOG.infof("Updating user with id: %d", id);

        return userRepository.findById(id)
            .onItem().ifNull().failWith(() ->
                new ResourceNotFoundException("User", "id", id))
            .flatMap(user -> {
                // Check email uniqueness if changed
                if (!user.email.equals(request.email())) {
                    return userRepository.existsByEmail(request.email())
                        .flatMap(exists -> {
                            if (exists) {
                                return Uni.createFrom().failure(
                                    new DuplicateResourceException("User", "email", request.email())
                                );
                            }
                            userMapper.updateEntity(request, user);
                            return userRepository.persist(user);
                        });
                }
                userMapper.updateEntity(request, user);
                return userRepository.persist(user);
            })
            .map(userMapper::toResponse);
    }

    @WithTransaction
    public Uni<Void> deleteUser(Long id) {
        LOG.infof("Deleting user with id: %d", id);

        return userRepository.findById(id)
            .onItem().ifNull().failWith(() ->
                new ResourceNotFoundException("User", "id", id))
            .flatMap(user -> userRepository.delete(user))
            .replaceWithVoid();
    }

    @WithTransaction
    public Uni<Void> deactivateUser(Long id) {
        LOG.infof("Deactivating user with id: %d", id);

        return userRepository.findById(id)
            .onItem().ifNull().failWith(() ->
                new ResourceNotFoundException("User", "id", id))
            .flatMap(user -> {
                user.active = false;
                return userRepository.persist(user);
            })
            .replaceWithVoid();
    }
}
```

---

## Resource (Controller)

```java
package com.example.resource;

import com.example.dto.UserRequest;
import com.example.dto.UserResponse;
import com.example.service.UserService;
import io.smallrye.mutiny.Uni;
import jakarta.annotation.security.RolesAllowed;
import jakarta.inject.Inject;
import jakarta.validation.Valid;
import jakarta.ws.rs.*;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import org.eclipse.microprofile.openapi.annotations.Operation;
import org.eclipse.microprofile.openapi.annotations.enums.SecuritySchemeType;
import org.eclipse.microprofile.openapi.annotations.responses.APIResponse;
import org.eclipse.microprofile.openapi.annotations.security.SecurityRequirement;
import org.eclipse.microprofile.openapi.annotations.security.SecurityScheme;
import org.eclipse.microprofile.openapi.annotations.tags.Tag;

import java.net.URI;
import java.util.List;

@Path("/api/v1/users")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
@Tag(name = "Users", description = "User management APIs")
@SecurityScheme(
    securitySchemeName = "jwt",
    type = SecuritySchemeType.HTTP,
    scheme = "bearer",
    bearerFormat = "JWT"
)
public class UserResource {

    @Inject
    UserService userService;

    @POST
    @Operation(summary = "Create a new user")
    @APIResponse(responseCode = "201", description = "User created successfully")
    @APIResponse(responseCode = "400", description = "Invalid input")
    @APIResponse(responseCode = "409", description = "Email already exists")
    public Uni<Response> createUser(@Valid UserRequest request) {
        return userService.createUser(request)
            .map(user -> Response
                .created(URI.create("/api/v1/users/" + user.id()))
                .entity(user)
                .build());
    }

    @GET
    @Path("/{id}")
    @RolesAllowed({"USER", "ADMIN"})
    @SecurityRequirement(name = "jwt")
    @Operation(summary = "Get user by ID")
    @APIResponse(responseCode = "200", description = "User found")
    @APIResponse(responseCode = "404", description = "User not found")
    public Uni<UserResponse> getUserById(@PathParam("id") Long id) {
        return userService.getUserById(id);
    }

    @GET
    @RolesAllowed({"USER", "ADMIN"})
    @SecurityRequirement(name = "jwt")
    @Operation(summary = "Get all users")
    public Uni<List<UserResponse>> getAllUsers() {
        return userService.getAllUsers();
    }

    @GET
    @Path("/active")
    @RolesAllowed({"USER", "ADMIN"})
    @SecurityRequirement(name = "jwt")
    @Operation(summary = "Get all active users")
    public Uni<List<UserResponse>> getActiveUsers() {
        return userService.getActiveUsers();
    }

    @PUT
    @Path("/{id}")
    @RolesAllowed({"USER", "ADMIN"})
    @SecurityRequirement(name = "jwt")
    @Operation(summary = "Update user")
    @APIResponse(responseCode = "200", description = "User updated successfully")
    @APIResponse(responseCode = "404", description = "User not found")
    public Uni<UserResponse> updateUser(
            @PathParam("id") Long id,
            @Valid UserRequest request) {
        return userService.updateUser(id, request);
    }

    @DELETE
    @Path("/{id}")
    @RolesAllowed("ADMIN")
    @SecurityRequirement(name = "jwt")
    @Operation(summary = "Delete user (Admin only)")
    @APIResponse(responseCode = "204", description = "User deleted successfully")
    @APIResponse(responseCode = "404", description = "User not found")
    public Uni<Response> deleteUser(@PathParam("id") Long id) {
        return userService.deleteUser(id)
            .map(ignored -> Response.noContent().build());
    }

    @PATCH
    @Path("/{id}/deactivate")
    @RolesAllowed({"USER", "ADMIN"})
    @SecurityRequirement(name = "jwt")
    @Operation(summary = "Deactivate user")
    public Uni<Response> deactivateUser(@PathParam("id") Long id) {
        return userService.deactivateUser(id)
            .map(ignored -> Response.noContent().build());
    }
}
```

---

## Exception Handling

### Custom Exceptions

```java
package com.example.exception;

public class ResourceNotFoundException extends RuntimeException {

    public ResourceNotFoundException(String resource, String field, Object value) {
        super(String.format("%s not found with %s: '%s'", resource, field, value));
    }
}
```

```java
package com.example.exception;

public class DuplicateResourceException extends RuntimeException {

    public DuplicateResourceException(String resource, String field, Object value) {
        super(String.format("%s already exists with %s: '%s'", resource, field, value));
    }
}
```

### Exception Mappers

```java
package com.example.exception;

import jakarta.validation.ConstraintViolationException;
import jakarta.ws.rs.core.Response;
import jakarta.ws.rs.ext.ExceptionMapper;
import jakarta.ws.rs.ext.Provider;
import org.jboss.logging.Logger;

import java.time.Instant;
import java.util.HashMap;
import java.util.Map;

@Provider
public class ExceptionMappers {

    private static final Logger LOG = Logger.getLogger(ExceptionMappers.class);

    @Provider
    public static class ResourceNotFoundMapper implements ExceptionMapper<ResourceNotFoundException> {
        @Override
        public Response toResponse(ResourceNotFoundException e) {
            LOG.warn("Resource not found: " + e.getMessage());

            Map<String, Object> error = new HashMap<>();
            error.put("title", "Resource Not Found");
            error.put("status", 404);
            error.put("detail", e.getMessage());
            error.put("timestamp", Instant.now().toString());

            return Response.status(Response.Status.NOT_FOUND)
                .entity(error)
                .build();
        }
    }

    @Provider
    public static class DuplicateResourceMapper implements ExceptionMapper<DuplicateResourceException> {
        @Override
        public Response toResponse(DuplicateResourceException e) {
            LOG.warn("Duplicate resource: " + e.getMessage());

            Map<String, Object> error = new HashMap<>();
            error.put("title", "Duplicate Resource");
            error.put("status", 409);
            error.put("detail", e.getMessage());
            error.put("timestamp", Instant.now().toString());

            return Response.status(Response.Status.CONFLICT)
                .entity(error)
                .build();
        }
    }

    @Provider
    public static class ValidationExceptionMapper implements ExceptionMapper<ConstraintViolationException> {
        @Override
        public Response toResponse(ConstraintViolationException e) {
            LOG.warn("Validation failed: " + e.getMessage());

            Map<String, String> errors = new HashMap<>();
            e.getConstraintViolations().forEach(violation -> {
                String field = violation.getPropertyPath().toString();
                errors.put(field, violation.getMessage());
            });

            Map<String, Object> error = new HashMap<>();
            error.put("title", "Validation Error");
            error.put("status", 400);
            error.put("detail", "Validation failed");
            error.put("errors", errors);
            error.put("timestamp", Instant.now().toString());

            return Response.status(Response.Status.BAD_REQUEST)
                .entity(error)
                .build();
        }
    }
}
```

---

## JWT Authentication

### Auth Resource

```java
package com.example.resource;

import com.example.dto.LoginRequest;
import com.example.dto.TokenResponse;
import com.example.entity.User;
import com.example.exception.AuthenticationException;
import io.quarkus.elytron.security.common.BcryptUtil;
import io.smallrye.jwt.build.Jwt;
import io.smallrye.mutiny.Uni;
import jakarta.inject.Inject;
import jakarta.validation.Valid;
import jakarta.ws.rs.*;
import jakarta.ws.rs.core.MediaType;
import org.eclipse.microprofile.config.inject.ConfigProperty;

import java.time.Duration;
import java.util.Set;

@Path("/api/v1/auth")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
public class AuthResource {

    @ConfigProperty(name = "mp.jwt.verify.issuer")
    String issuer;

    @POST
    @Path("/login")
    public Uni<TokenResponse> login(@Valid LoginRequest request) {
        return User.findByEmail(request.email())
            .onItem().ifNull().failWith(() ->
                new AuthenticationException("Invalid credentials"))
            .flatMap(user -> {
                if (!BcryptUtil.matches(request.password(), user.password)) {
                    return Uni.createFrom().failure(
                        new AuthenticationException("Invalid credentials")
                    );
                }

                String token = Jwt.issuer(issuer)
                    .upn(user.email)
                    .groups(Set.of(user.role.name()))
                    .claim("userId", user.id)
                    .expiresIn(Duration.ofHours(24))
                    .sign();

                return Uni.createFrom().item(new TokenResponse(token, "Bearer", 86400L));
            });
    }
}
```

---

## Health Checks

```java
package com.example.health;

import com.example.repository.UserRepository;
import io.smallrye.health.api.Wellness;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import org.eclipse.microprofile.health.*;

@Wellness
@ApplicationScoped
public class DatabaseHealthCheck implements HealthCheck {

    @Inject
    UserRepository userRepository;

    @Override
    public HealthCheckResponse call() {
        try {
            long count = userRepository.count().await().indefinitely();
            return HealthCheckResponse.up("Database connection")
                .withData("users_count", count)
                .build();
        } catch (Exception e) {
            return HealthCheckResponse.down("Database connection")
                .withData("error", e.getMessage())
                .build();
        }
    }
}
```

---

## Testing

### Resource Tests

```java
package com.example.resource;

import com.example.dto.UserRequest;
import io.quarkus.test.junit.QuarkusTest;
import io.quarkus.test.security.TestSecurity;
import io.restassured.http.ContentType;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.*;

@QuarkusTest
@DisplayName("UserResource")
class UserResourceTest {

    @Test
    @DisplayName("POST /api/v1/users - should create user")
    void shouldCreateUser() {
        UserRequest request = new UserRequest(
            "test@example.com",
            "Password123!",
            "Test User",
            null
        );

        given()
            .contentType(ContentType.JSON)
            .body(request)
        .when()
            .post("/api/v1/users")
        .then()
            .statusCode(201)
            .body("id", notNullValue())
            .body("email", equalTo("test@example.com"))
            .body("name", equalTo("Test User"));
    }

    @Test
    @DisplayName("POST /api/v1/users - should return 400 for invalid email")
    void shouldReturnBadRequestForInvalidEmail() {
        UserRequest request = new UserRequest(
            "invalid-email",
            "Password123!",
            "Test User",
            null
        );

        given()
            .contentType(ContentType.JSON)
            .body(request)
        .when()
            .post("/api/v1/users")
        .then()
            .statusCode(400)
            .body("errors.email", notNullValue());
    }

    @Test
    @TestSecurity(user = "test", roles = "USER")
    @DisplayName("GET /api/v1/users/{id} - should return user")
    void shouldReturnUser() {
        // Assuming user with id 1 exists
        given()
        .when()
            .get("/api/v1/users/1")
        .then()
            .statusCode(200)
            .body("id", equalTo(1));
    }

    @Test
    @DisplayName("GET /api/v1/users/{id} - should return 401 without auth")
    void shouldReturnUnauthorizedWithoutAuth() {
        given()
        .when()
            .get("/api/v1/users/1")
        .then()
            .statusCode(401);
    }

    @Test
    @TestSecurity(user = "admin", roles = "ADMIN")
    @DisplayName("DELETE /api/v1/users/{id} - admin should delete user")
    void adminShouldDeleteUser() {
        // Create a user first, then delete
        UserRequest request = new UserRequest(
            "delete@example.com",
            "Password123!",
            "Delete User",
            null
        );

        Integer id = given()
            .contentType(ContentType.JSON)
            .body(request)
        .when()
            .post("/api/v1/users")
        .then()
            .statusCode(201)
            .extract().path("id");

        given()
        .when()
            .delete("/api/v1/users/" + id)
        .then()
            .statusCode(204);
    }

    @Test
    @TestSecurity(user = "user", roles = "USER")
    @DisplayName("DELETE /api/v1/users/{id} - regular user should get forbidden")
    void regularUserShouldNotDeleteUser() {
        given()
        .when()
            .delete("/api/v1/users/1")
        .then()
            .statusCode(403);
    }
}
```

### Service Tests

```java
package com.example.service;

import com.example.dto.UserRequest;
import com.example.dto.UserResponse;
import com.example.entity.User;
import com.example.exception.DuplicateResourceException;
import com.example.exception.ResourceNotFoundException;
import com.example.mapper.UserMapper;
import com.example.repository.UserRepository;
import io.quarkus.test.InjectMock;
import io.quarkus.test.junit.QuarkusTest;
import io.smallrye.mutiny.Uni;
import jakarta.inject.Inject;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@QuarkusTest
@DisplayName("UserService")
class UserServiceTest {

    @Inject
    UserService userService;

    @InjectMock
    UserRepository userRepository;

    @InjectMock
    UserMapper userMapper;

    private User user;
    private UserRequest request;
    private UserResponse response;

    @BeforeEach
    void setUp() {
        user = new User();
        user.id = 1L;
        user.email = "test@example.com";
        user.password = "encodedPassword";
        user.name = "Test User";
        user.role = User.Role.USER;
        user.active = true;

        request = new UserRequest(
            "test@example.com",
            "Password123!",
            "Test User",
            null
        );

        response = new UserResponse(
            1L,
            "test@example.com",
            "Test User",
            "USER",
            true,
            null
        );
    }

    @Test
    @DisplayName("createUser should create user with valid data")
    void shouldCreateUserWithValidData() {
        // Given
        when(userRepository.existsByEmail(request.email())).thenReturn(Uni.createFrom().item(false));
        when(userMapper.toEntity(request)).thenReturn(user);
        when(userRepository.persist(any(User.class))).thenReturn(Uni.createFrom().item(user));
        when(userMapper.toResponse(user)).thenReturn(response);

        // When
        UserResponse result = userService.createUser(request).await().indefinitely();

        // Then
        assertThat(result).isNotNull();
        assertThat(result.email()).isEqualTo("test@example.com");
    }

    @Test
    @DisplayName("createUser should throw exception when email exists")
    void shouldThrowExceptionWhenEmailExists() {
        // Given
        when(userRepository.existsByEmail(request.email())).thenReturn(Uni.createFrom().item(true));

        // When/Then
        assertThatThrownBy(() -> userService.createUser(request).await().indefinitely())
            .isInstanceOf(DuplicateResourceException.class)
            .hasMessageContaining("email");
    }

    @Test
    @DisplayName("getUserById should return user when found")
    void shouldReturnUserWhenFound() {
        // Given
        when(userRepository.findById(1L)).thenReturn(Uni.createFrom().item(user));
        when(userMapper.toResponse(user)).thenReturn(response);

        // When
        UserResponse result = userService.getUserById(1L).await().indefinitely();

        // Then
        assertThat(result).isNotNull();
        assertThat(result.id()).isEqualTo(1L);
    }

    @Test
    @DisplayName("getUserById should throw exception when not found")
    void shouldThrowExceptionWhenNotFound() {
        // Given
        when(userRepository.findById(999L)).thenReturn(Uni.createFrom().nullItem());

        // When/Then
        assertThatThrownBy(() -> userService.getUserById(999L).await().indefinitely())
            .isInstanceOf(ResourceNotFoundException.class)
            .hasMessageContaining("999");
    }
}
```

---

## Native Compilation

### Dockerfile.native

```dockerfile
FROM quay.io/quarkus/ubi-quarkus-mandrel-builder-image:jdk-21 AS build
COPY --chown=quarkus:quarkus mvnw /code/mvnw
COPY --chown=quarkus:quarkus .mvn /code/.mvn
COPY --chown=quarkus:quarkus pom.xml /code/
USER quarkus
WORKDIR /code
RUN ./mvnw -B org.apache.maven.plugins:maven-dependency-plugin:3.1.2:go-offline
COPY src /code/src
RUN ./mvnw package -Pnative -DskipTests

FROM quay.io/quarkus/quarkus-micro-image:2.0
WORKDIR /work/
COPY --from=build /code/target/*-runner /work/application
RUN chmod 775 /work
EXPOSE 8080
CMD ["./application", "-Dquarkus.http.host=0.0.0.0"]
```

---

## Commands

```bash
# Create project
mvn io.quarkus.platform:quarkus-maven-plugin:3.6.0:create \
    -DprojectGroupId=com.example \
    -DprojectArtifactId=myproject \
    -Dextensions="resteasy-reactive-jackson,hibernate-reactive-panache,reactive-pg-client"

# Dev mode (live reload)
./mvnw quarkus:dev

# Build
./mvnw package

# Build native
./mvnw package -Pnative

# Build native in container (no GraalVM needed locally)
./mvnw package -Pnative -Dquarkus.native.container-build=true

# Test
./mvnw test

# Test with native
./mvnw verify -Pnative

# List extensions
./mvnw quarkus:list-extensions

# Add extension
./mvnw quarkus:add-extension -Dextensions="openapi"

# Build container image
./mvnw package -Dquarkus.container-image.build=true

# Run Kubernetes
./mvnw package -Dquarkus.kubernetes.deploy=true
```

---

## Best Practices

### Do's
- ✓ Use Dev Services for local development
- ✓ Use Panache for simpler data access
- ✓ Use reactive programming for I/O-bound operations
- ✓ Use CDI for dependency injection
- ✓ Configure for native compilation early
- ✓ Use health checks and metrics
- ✓ Use `@WithTransaction` for transactional operations
- ✓ Test with `@QuarkusTest` and `@TestSecurity`

### Don'ts
- ✗ Don't use reflection without registration for native
- ✗ Don't use blocking operations in reactive pipelines
- ✗ Don't ignore native compilation compatibility
- ✗ Don't use synchronous APIs with reactive datasources
- ✗ Don't hardcode configuration (use ConfigProperty)

---

## Comparison: Quarkus vs Spring Boot

| Aspect | Quarkus | Spring Boot |
|--------|---------|-------------|
| Startup Time | ~10ms (native), ~1s (JVM) | ~3-5s |
| Memory | ~10MB (native), ~100MB (JVM) | ~200-400MB |
| Native Compilation | First-class support | Spring Native (less mature) |
| Dev Experience | Live reload, Dev Services | Spring DevTools |
| Reactive | Unified (RESTEasy Reactive + Mutiny) | WebFlux (separate stack) |
| Ecosystem | Growing | Massive |
| Learning Curve | Moderate | Well-known |

---

## References

- [Quarkus Documentation](https://quarkus.io/guides/)
- [Quarkus Extensions](https://quarkus.io/extensions/)
- [Panache Guide](https://quarkus.io/guides/hibernate-orm-panache)
- [RESTEasy Reactive Guide](https://quarkus.io/guides/resteasy-reactive)
- [SmallRye JWT](https://quarkus.io/guides/security-jwt)
- [Native Compilation](https://quarkus.io/guides/building-native-image)
