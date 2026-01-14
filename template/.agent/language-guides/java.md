# Java Guide

> **Applies to**: Java 17+, Spring Boot, Jakarta EE, Android (legacy), Microservices

---

## Core Principles

1. **Type Safety**: Strong static typing, compile-time checks
2. **Object-Oriented**: Classes, interfaces, inheritance (prefer composition)
3. **Platform Independence**: Write once, run anywhere (JVM)
4. **Modern Java**: Use records, sealed classes, pattern matching (Java 17+)
5. **Null Safety**: Use `Optional`, `@Nullable`/`@NonNull` annotations

---

## Language-Specific Guardrails

### Java Version & Setup
- ✓ Use Java 17+ LTS (21 LTS recommended for new projects)
- ✓ Use build tools: Maven (`pom.xml`) or Gradle (`build.gradle`)
- ✓ Pin dependency versions (avoid dynamic versions)
- ✓ Include Java version in `pom.xml` or `build.gradle`

### Code Style (Google Java Style)
- ✓ Run formatter before every commit (google-java-format, Spotless)
- ✓ Follow Google Java Style Guide or project conventions
- ✓ Line length: 100 characters
- ✓ Use `camelCase` for methods and variables
- ✓ Use `PascalCase` for classes and interfaces
- ✓ Use `SCREAMING_SNAKE_CASE` for constants
- ✓ 4-space indentation (not tabs)
- ✓ Braces on same line (K&R style)

### Modern Java Features (17+)
- ✓ Use `var` for local variables with obvious types
- ✓ Use records for immutable data classes
- ✓ Use sealed classes for restricted hierarchies
- ✓ Use pattern matching in `instanceof` and `switch`
- ✓ Use text blocks for multi-line strings
- ✓ Use `Optional` for nullable return values
- ✓ Use Stream API for collection processing

### Null Safety
- ✓ Never return `null` from methods (use `Optional`)
- ✓ Validate parameters with `Objects.requireNonNull()`
- ✓ Use `@Nullable`/`@NonNull` annotations (JSR-305 or JetBrains)
- ✓ Use `Optional.ofNullable()` for potentially null values
- ✓ Avoid `Optional` as method parameters or fields

### Exception Handling
- ✓ Use specific exception types (not generic `Exception`)
- ✓ Don't catch `Exception` or `Throwable` unless re-throwing
- ✓ Always close resources (try-with-resources)
- ✓ Don't swallow exceptions (at minimum, log them)
- ✓ Create custom exceptions for domain errors
- ✓ Document exceptions in Javadoc with `@throws`

### Collections & Streams
- ✓ Use immutable collections where possible (`List.of()`, `Set.of()`)
- ✓ Prefer Stream API over imperative loops
- ✓ Use `Collectors` for complex reductions
- ✓ Avoid side effects in stream operations
- ✓ Use parallel streams only for CPU-bound, large datasets

### Concurrency
- ✓ Prefer `ExecutorService` over raw threads
- ✓ Use `CompletableFuture` for async operations
- ✓ Use virtual threads (Java 21+) for I/O-bound tasks
- ✓ Prefer immutable objects for thread safety
- ✓ Use `ConcurrentHashMap`, `AtomicInteger`, etc. for mutable shared state
- ✓ Always handle `InterruptedException` properly

---

## Project Structure

### Maven Standard Layout
```
myproject/
├── pom.xml
├── src/
│   ├── main/
│   │   ├── java/
│   │   │   └── com/example/myproject/
│   │   │       ├── Application.java
│   │   │       ├── controller/
│   │   │       ├── service/
│   │   │       ├── repository/
│   │   │       ├── model/
│   │   │       ├── dto/
│   │   │       └── config/
│   │   └── resources/
│   │       ├── application.yml
│   │       └── db/migration/
│   └── test/
│       ├── java/
│       └── resources/
└── README.md
```

### Guardrails
- ✓ Follow Maven/Gradle standard directory layout
- ✓ Package by feature or layer (be consistent)
- ✓ One public class per file
- ✓ Tests mirror main source structure

---

## Validation & Input Handling

### Recommended Libraries
- **Bean Validation (JSR-380)**: Constraint annotations
- **Hibernate Validator**: Reference implementation
- **Jakarta Validation**: Modern standard

### Pattern
```java
import jakarta.validation.Valid;
import jakarta.validation.constraints.*;

public record UserCreate(
    @NotBlank @Email
    String email,

    @Positive
    int age,

    @NotBlank @Pattern(regexp = "^(admin|user|guest)$")
    String role
) {}

// Custom validator
@Documented
@Constraint(validatedBy = UniqueEmailValidator.class)
@Target({ElementType.FIELD})
@Retention(RetentionPolicy.RUNTIME)
public @interface UniqueEmail {
    String message() default "Email already exists";
    Class<?>[] groups() default {};
    Class<? extends Payload>[] payload() default {};
}

// Usage in controller
@PostMapping("/users")
public ResponseEntity<User> createUser(@Valid @RequestBody UserCreate request) {
    return ResponseEntity.ok(userService.create(request));
}
```

---

## Testing

### Frameworks
- **JUnit 5**: Standard testing framework
- **Mockito**: Mocking library
- **AssertJ**: Fluent assertions
- **Testcontainers**: Integration testing with containers
- **ArchUnit**: Architecture testing

### Guardrails
- ✓ Test files: `*Test.java` (unit), `*IT.java` (integration)
- ✓ Use descriptive test names: `shouldThrowExceptionWhenEmailInvalid()`
- ✓ Use `@DisplayName` for readable test descriptions
- ✓ Use `@Nested` for grouping related tests
- ✓ Use `@ParameterizedTest` for multiple test cases
- ✓ Coverage target: >80% for business logic
- ✓ Use Testcontainers for database integration tests

### Example
```java
import org.junit.jupiter.api.*;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;
import static org.assertj.core.api.Assertions.*;
import static org.mockito.Mockito.*;

@DisplayName("UserService")
class UserServiceTest {

    private UserRepository repository;
    private UserService service;

    @BeforeEach
    void setUp() {
        repository = mock(UserRepository.class);
        service = new UserService(repository);
    }

    @Nested
    @DisplayName("createUser")
    class CreateUser {

        @Test
        @DisplayName("should create user with valid data")
        void shouldCreateUserWithValidData() {
            var request = new UserCreate("test@example.com", 25, "user");
            var expected = new User(1L, "test@example.com", 25, "user");
            when(repository.save(any())).thenReturn(expected);

            var result = service.create(request);

            assertThat(result)
                .isNotNull()
                .satisfies(user -> {
                    assertThat(user.email()).isEqualTo("test@example.com");
                    assertThat(user.age()).isEqualTo(25);
                });
            verify(repository).save(any(User.class));
        }

        @ParameterizedTest
        @ValueSource(strings = {"", " ", "invalid"})
        @DisplayName("should throw exception for invalid email")
        void shouldThrowExceptionForInvalidEmail(String email) {
            var request = new UserCreate(email, 25, "user");

            assertThatThrownBy(() -> service.create(request))
                .isInstanceOf(ValidationException.class)
                .hasMessageContaining("email");
        }
    }
}
```

---

## Tooling

### Essential Tools
- **Checkstyle**: Style enforcement
- **SpotBugs**: Bug detection
- **PMD**: Code analysis
- **Spotless**: Code formatting
- **JaCoCo**: Code coverage
- **OWASP Dependency-Check**: Vulnerability scanning

### Configuration (Maven)
```xml
<!-- pom.xml -->
<properties>
    <java.version>21</java.version>
    <maven.compiler.source>${java.version}</maven.compiler.source>
    <maven.compiler.target>${java.version}</maven.compiler.target>
</properties>

<build>
    <plugins>
        <!-- Compiler -->
        <plugin>
            <groupId>org.apache.maven.plugins</groupId>
            <artifactId>maven-compiler-plugin</artifactId>
            <version>3.11.0</version>
            <configuration>
                <compilerArgs>
                    <arg>-Xlint:all</arg>
                </compilerArgs>
            </configuration>
        </plugin>

        <!-- Spotless formatter -->
        <plugin>
            <groupId>com.diffplug.spotless</groupId>
            <artifactId>spotless-maven-plugin</artifactId>
            <version>2.40.0</version>
            <configuration>
                <java>
                    <googleJavaFormat>
                        <version>1.18.1</version>
                    </googleJavaFormat>
                </java>
            </configuration>
        </plugin>

        <!-- JaCoCo coverage -->
        <plugin>
            <groupId>org.jacoco</groupId>
            <artifactId>jacoco-maven-plugin</artifactId>
            <version>0.8.11</version>
        </plugin>
    </plugins>
</build>
```

### Pre-Commit Commands
```bash
# Format
./mvnw spotless:apply

# Compile
./mvnw compile

# Test
./mvnw test

# Test with coverage
./mvnw verify

# Static analysis
./mvnw checkstyle:check pmd:check spotbugs:check

# Dependency vulnerability check
./mvnw dependency-check:check
```

---

## Common Pitfalls

### Don't Do This
```java
// Returning null
public User findById(Long id) {
    return repository.findById(id).orElse(null);
}

// Catching generic Exception
try {
    riskyOperation();
} catch (Exception e) {
    // Lost information
}

// Mutable return values
public List<User> getUsers() {
    return this.users; // Caller can modify internal state
}

// Raw types
List list = new ArrayList();

// Not closing resources
FileInputStream fis = new FileInputStream("file.txt");
// ... use fis
fis.close(); // May not be called if exception occurs
```

### Do This Instead
```java
// Use Optional
public Optional<User> findById(Long id) {
    return repository.findById(id);
}

// Catch specific exceptions
try {
    riskyOperation();
} catch (IOException e) {
    logger.error("I/O error during operation", e);
    throw new ServiceException("Failed to process", e);
}

// Return immutable copies
public List<User> getUsers() {
    return List.copyOf(this.users);
}

// Use generics
List<User> users = new ArrayList<>();

// Try-with-resources
try (var fis = new FileInputStream("file.txt")) {
    // ... use fis
} // Automatically closed
```

---

## Framework-Specific Patterns

### Spring Boot
```java
// REST Controller
@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;

    @GetMapping("/{id}")
    public ResponseEntity<User> getUser(@PathVariable Long id) {
        return userService.findById(id)
            .map(ResponseEntity::ok)
            .orElse(ResponseEntity.notFound().build());
    }

    @PostMapping
    public ResponseEntity<User> createUser(@Valid @RequestBody UserCreate request) {
        var user = userService.create(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(user);
    }
}

// Service with transactions
@Service
@RequiredArgsConstructor
public class UserService {

    private final UserRepository repository;

    @Transactional(readOnly = true)
    public Optional<User> findById(Long id) {
        return repository.findById(id);
    }

    @Transactional
    public User create(UserCreate request) {
        var user = User.builder()
            .email(request.email())
            .age(request.age())
            .role(request.role())
            .build();
        return repository.save(user);
    }
}

// Repository
public interface UserRepository extends JpaRepository<User, Long> {
    Optional<User> findByEmail(String email);

    @Query("SELECT u FROM User u WHERE u.role = :role")
    List<User> findByRole(@Param("role") String role);
}
```

### Records (Java 17+)
```java
// Immutable data class
public record User(
    Long id,
    String email,
    int age,
    String role
) {
    // Compact constructor for validation
    public User {
        Objects.requireNonNull(email, "Email is required");
        if (age < 0) {
            throw new IllegalArgumentException("Age must be non-negative");
        }
    }

    // Additional methods
    public boolean isAdmin() {
        return "admin".equals(role);
    }
}

// DTO records
public record UserCreate(String email, int age, String role) {}
public record UserResponse(Long id, String email, String role) {
    public static UserResponse from(User user) {
        return new UserResponse(user.id(), user.email(), user.role());
    }
}
```

---

## Performance Considerations

### Optimization Guardrails
- ✓ Use connection pooling (HikariCP)
- ✓ Enable JPA/Hibernate batch operations
- ✓ Use `@Transactional(readOnly = true)` for read operations
- ✓ Avoid N+1 queries (use fetch joins, entity graphs)
- ✓ Use pagination for large datasets
- ✓ Profile with JFR, VisualVM, or async-profiler before optimizing
- ✓ Use virtual threads (Java 21+) for I/O-bound workloads

### Example
```java
// Fetch join to avoid N+1
@Query("SELECT u FROM User u LEFT JOIN FETCH u.orders WHERE u.id = :id")
Optional<User> findByIdWithOrders(@Param("id") Long id);

// Pagination
Page<User> findByRole(String role, Pageable pageable);

// Batch operations
@Modifying
@Query("UPDATE User u SET u.status = :status WHERE u.lastLogin < :date")
int updateInactiveUsers(@Param("status") String status, @Param("date") LocalDate date);
```

---

## Security Best Practices

### Guardrails
- ✓ Never hardcode secrets (use environment variables, Vault)
- ✓ Use parameterized queries (JPA/Hibernate handles this)
- ✓ Validate all inputs with Bean Validation
- ✓ Hash passwords with BCrypt or Argon2
- ✓ Enable CSRF protection (Spring Security default)
- ✓ Use HTTPS in production
- ✓ Run OWASP Dependency-Check regularly
- ✓ Keep dependencies updated

### Example
```java
// Password hashing (Spring Security)
@Bean
public PasswordEncoder passwordEncoder() {
    return new BCryptPasswordEncoder(12);
}

// Security configuration
@Configuration
@EnableWebSecurity
public class SecurityConfig {

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        return http
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/api/public/**").permitAll()
                .requestMatchers("/api/admin/**").hasRole("ADMIN")
                .anyRequest().authenticated()
            )
            .csrf(csrf -> csrf
                .csrfTokenRepository(CookieCsrfTokenRepository.withHttpOnlyFalse())
            )
            .sessionManagement(session -> session
                .sessionCreationPolicy(SessionCreationPolicy.STATELESS)
            )
            .build();
    }
}
```

---

## References

- [Oracle Java Documentation](https://docs.oracle.com/en/java/)
- [Google Java Style Guide](https://google.github.io/styleguide/javaguide.html)
- [Effective Java (Book)](https://www.oreilly.com/library/view/effective-java/9780134686097/)
- [Spring Boot Documentation](https://docs.spring.io/spring-boot/docs/current/reference/html/)
- [JUnit 5 User Guide](https://junit.org/junit5/docs/current/user-guide/)
- [AssertJ Documentation](https://assertj.github.io/doc/)
- [Baeldung](https://www.baeldung.com/) (tutorials)
