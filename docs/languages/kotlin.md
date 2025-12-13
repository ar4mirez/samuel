---
title: Kotlin Guide
description: Kotlin development guardrails and best practices
---

# Kotlin Guide

> **Applies to**: Kotlin 1.9+, Android, Spring Boot, Ktor, Multiplatform

---

## Core Principles

1. **Null Safety First**: Type system eliminates NullPointerExceptions
2. **Conciseness & Readability**: Reduce boilerplate, expressive code
3. **Seamless Java Interoperability**: Works with existing Java codebases
4. **Coroutines for Concurrency**: Structured async programming
5. **Immutability by Default**: Prefer `val` over `var`, immutable collections

---

## Language-Specific Guardrails

### Kotlin Version & Setup

```
✓ Use Kotlin 1.9+ (2.0+ recommended for K2 compiler)
✓ Use Gradle Kotlin DSL for build configuration
✓ Pin dependency versions in version catalogs (libs.versions.toml)
✓ Include Kotlin version in build.gradle.kts
```

### Null Safety

```
✓ Prefer non-nullable types by default
✓ Use safe call operator ?. instead of explicit null checks
✓ Use Elvis operator ?: for default values
✓ Avoid !! operator without explicit justification and comment
✓ Use lateinit only when initialization truly deferred (DI, Android)
✓ Use by lazy for expensive one-time initialization
✓ Use requireNotNull() or checkNotNull() with descriptive messages
```

### Code Style (ktlint)

```
✓ Run ./gradlew ktlintFormat before every commit
✓ Follow official Kotlin coding conventions
✓ Line length: 100 characters (ktlint default)
✓ Use camelCase for functions and properties
✓ Use PascalCase for classes and interfaces
✓ Use SCREAMING_SNAKE_CASE for constants
✓ 4-space indentation (not tabs)
✓ Use trailing commas in multiline declarations
```

### Immutability

```
✓ Prefer val over var (immutable by default)
✓ Use data class for DTOs (automatic equals/hashCode/toString/copy)
✓ Use immutable collections: List, Set, Map instead of MutableList, etc.
✓ Use .copy() for data class modifications
✓ Mark properties private unless they need to be public
```

### Coroutines

```
✓ Always use structured concurrency (never GlobalScope)
✓ Use suspend functions for async operations
✓ Make suspend functions main-safe (use withContext(Dispatchers.IO))
✓ Use coroutineScope or supervisorScope for concurrent work
✓ Always handle CancellationException properly (don't catch it)
✓ Set timeouts for long-running operations (withTimeout)
✓ Use Flow for streams of data
```

---

## Validation & Input Handling

### Pattern

```kotlin
import kotlinx.serialization.Serializable
import kotlinx.serialization.decodeFromString
import kotlinx.serialization.json.Json

@Serializable
data class UserCreate(
    val email: String,
    val age: Int,
    val role: String
) {
    init {
        require(email.contains("@")) { "Invalid email format" }
        require(age > 0) { "Age must be positive" }
        require(role in listOf("admin", "user", "guest")) { "Invalid role" }
    }
}

// Usage
fun createUser(json: String): User {
    val validated = Json.decodeFromString<UserCreate>(json)
    return User(email = validated.email, age = validated.age, role = validated.role)
}
```

---

## Testing

### Frameworks

| Framework | Use Case |
|-----------|----------|
| **Kotest** | Kotlin-idiomatic, multiple test styles (recommended) |
| **JUnit 5** | Industry standard (verbose but mature) |
| **MockK** | Mocking library (mocks final classes) |
| **Turbine** | Flow testing library |

### Guardrails

```
✓ Test files: *Test.kt or descriptive names
✓ Use descriptive test names with backticks
✓ Prefer Kotest's shouldBe over JUnit's assertEquals
✓ Use MockK for mocking (not Mockito for Kotlin)
✓ Use coEvery, coVerify for suspend functions
✓ Test coroutines with runTest (from kotlinx-coroutines-test)
✓ Coverage target: >80% for business logic
```

### Example (Kotest with MockK)

```kotlin
import io.kotest.core.spec.style.StringSpec
import io.kotest.matchers.shouldBe
import io.mockk.*
import kotlinx.coroutines.test.runTest

class UserServiceTest : StringSpec({
    val repository = mockk<UserRepository>()
    val service = UserService(repository)

    afterTest {
        clearAllMocks()
    }

    "should return user when id is valid" {
        val user = User(id = "123", email = "test@example.com")
        coEvery { repository.findById("123") } returns user

        val result = service.getUser("123")

        result shouldBe user
        coVerify(exactly = 1) { repository.findById("123") }
    }

    "should throw exception when email is invalid" {
        val exception = shouldThrow<ValidationException> {
            service.createUser(email = "invalid", age = 25)
        }
        exception.message shouldBe "Invalid email format"
    }
})
```

---

## Tooling

### Essential Tools

| Tool | Purpose |
|------|---------|
| **ktlint** | Code formatter (opinionated) |
| **detekt** | Static analysis (code smells, complexity) |
| **Gradle Kotlin DSL** | Build configuration |
| **KSP** | Kotlin Symbol Processing (faster than KAPT) |

### Configuration

```kotlin
// build.gradle.kts
plugins {
    kotlin("jvm") version "1.9.22"
    kotlin("plugin.serialization") version "1.9.22"
    id("org.jlleitschuh.gradle.ktlint") version "12.1.0"
    id("io.gitlab.arturbosch.detekt") version "1.23.4"
}

kotlin {
    jvmToolchain(17)
    compilerOptions {
        freeCompilerArgs.add("-Xjsr305=strict")
    }
}
```

### Pre-Commit Commands

```bash
# Format
./gradlew ktlintFormat

# Lint
./gradlew ktlintCheck

# Static analysis
./gradlew detekt

# Test
./gradlew test

# Build
./gradlew build
```

---

## Common Pitfalls

### Don't Do This

```kotlin
// ❌ Using !!
val value = nullableValue!!

// ❌ Using var when val works
var user = User("John")

// ❌ GlobalScope for coroutines
GlobalScope.launch {
    doWork()
}

// ❌ Mutable collections as public properties
class UserManager {
    val users = mutableListOf<User>()
}
```

### Do This Instead

```kotlin
// ✅ Safe call or Elvis operator
val value = nullableValue ?: defaultValue

// ✅ Immutable by default
val user = User("John")

// ✅ Structured concurrency
suspend fun doWork() = coroutineScope {
    launch {
        // Work here
    }
}

// ✅ Immutable public API
class UserManager {
    private val _users = mutableListOf<User>()
    val users: List<User> get() = _users.toList()
}
```

---

## Framework-Specific Patterns

### Android + Jetpack Compose

```kotlin
// ViewModel with StateFlow
class UserViewModel(
    private val repository: UserRepository
) : ViewModel() {
    private val _uiState = MutableStateFlow<UiState>(UiState.Loading)
    val uiState: StateFlow<UiState> = _uiState.asStateFlow()

    init {
        loadUsers()
    }

    private fun loadUsers() {
        viewModelScope.launch {
            _uiState.value = UiState.Loading
            try {
                val users = repository.getUsers()
                _uiState.value = UiState.Success(users)
            } catch (e: Exception) {
                _uiState.value = UiState.Error(e.message ?: "Unknown error")
            }
        }
    }
}

sealed interface UiState {
    object Loading : UiState
    data class Success(val users: List<User>) : UiState
    data class Error(val message: String) : UiState
}

// Composable UI
@Composable
fun UserScreen(viewModel: UserViewModel = hiltViewModel()) {
    val uiState by viewModel.uiState.collectAsStateWithLifecycle()

    when (val state = uiState) {
        is UiState.Loading -> LoadingIndicator()
        is UiState.Success -> UserList(users = state.users)
        is UiState.Error -> ErrorMessage(message = state.message)
    }
}
```

### Spring Boot + Kotlin

```kotlin
// REST Controller
@RestController
@RequestMapping("/api/users")
class UserController(private val service: UserService) {

    @GetMapping("/{id}")
    suspend fun getUser(@PathVariable id: String): ResponseEntity<User> {
        val user = service.findById(id) ?: return ResponseEntity.notFound().build()
        return ResponseEntity.ok(user)
    }

    @PostMapping
    suspend fun createUser(@Valid @RequestBody request: UserCreateRequest): User {
        return service.create(request)
    }
}

// Data class with validation
data class UserCreateRequest(
    @field:Email(message = "Invalid email format")
    val email: String,

    @field:Positive(message = "Age must be positive")
    val age: Int,

    @field:NotBlank
    val name: String
)
```

---

## Performance Considerations

### Optimization Guardrails

```
✓ Use sequences for large collections with multiple operations
✓ Use inline functions for lambdas passed to higher-order functions
✓ Use by lazy for expensive one-time initialization
✓ Avoid creating unnecessary objects in loops
✓ Use buildList, buildMap for collection builders
✓ Profile before optimizing (Android Profiler, JProfiler)
```

### Example

```kotlin
// Use sequences for large datasets
val result = hugeList
    .asSequence()
    .filter { it.isActive }
    .map { it.name }
    .take(10)
    .toList()

// Lazy initialization
val expensiveResource by lazy {
    // Computed only once, when first accessed
    loadExpensiveResource()
}
```

---

## Security Best Practices

### Guardrails

```
✓ Never hardcode secrets (use environment variables or secure storage)
✓ Validate all user inputs with data class init blocks or Bean Validation
✓ Use parameterized queries (prevent SQL injection)
✓ Hash passwords with bcrypt or Argon2
✓ Enable CSRF protection (Spring Security)
✓ Run ./gradlew dependencyCheckAnalyze for vulnerability scanning
✓ Use ProGuard/R8 for Android (obfuscation)
```

### Example

```kotlin
// Password hashing (Spring Security)
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder

class PasswordService {
    private val encoder = BCryptPasswordEncoder(12)

    fun hashPassword(plainPassword: String): String {
        return encoder.encode(plainPassword)
    }

    fun verifyPassword(plainPassword: String, hashedPassword: String): Boolean {
        return encoder.matches(plainPassword, hashedPassword)
    }
}
```

---

## References

- [Official Kotlin Documentation](https://kotlinlang.org/docs/home.html)
- [Kotlin Coding Conventions](https://kotlinlang.org/docs/coding-conventions.html)
- [ktlint](https://pinterest.github.io/ktlint/latest/)
- [detekt](https://detekt.dev/)
- [Kotest](https://kotest.io/)
- [MockK](https://mockk.io/)
- [Kotlin Coroutines Guide](https://kotlinlang.org/docs/coroutines-guide.html)
- [Android Kotlin Guide](https://developer.android.com/kotlin)
