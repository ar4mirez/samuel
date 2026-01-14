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
- ✓ Use Kotlin 1.9+ (2.0+ recommended for K2 compiler)
- ✓ Use Gradle Kotlin DSL for build configuration
- ✓ Pin dependency versions in version catalogs (`libs.versions.toml`)
- ✓ Include Kotlin version in `build.gradle.kts`

### Null Safety
- ✓ Prefer non-nullable types by default
- ✓ Use safe call operator `?.` instead of explicit null checks
- ✓ Use Elvis operator `?:` for default values
- ✓ Avoid `!!` operator without explicit justification and comment
- ✓ Use `lateinit` only when initialization truly deferred (DI, Android)
- ✓ Use `by lazy` for expensive one-time initialization
- ✓ Use `requireNotNull()` or `checkNotNull()` with descriptive messages

### Code Style (ktlint)
- ✓ Run `./gradlew ktlintFormat` before every commit
- ✓ Follow official Kotlin coding conventions
- ✓ Line length: 100 characters (ktlint default)
- ✓ Use `camelCase` for functions and properties
- ✓ Use `PascalCase` for classes and interfaces
- ✓ Use `SCREAMING_SNAKE_CASE` for constants
- ✓ 4-space indentation (not tabs)
- ✓ Use trailing commas in multiline declarations

### Immutability
- ✓ Prefer `val` over `var` (immutable by default)
- ✓ Use `data class` for DTOs (automatic equals/hashCode/toString/copy)
- ✓ Use immutable collections: `List`, `Set`, `Map` instead of `MutableList`, etc.
- ✓ Use `.copy()` for data class modifications
- ✓ Mark properties `private` unless they need to be public

### Functions & Classes
- ✓ Prefer expression bodies for single-expression functions
- ✓ Use named arguments for functions with multiple parameters
- ✓ Use default parameter values instead of overloads
- ✓ Functions ≤50 lines (extract helper functions)
- ✓ Classes ≤300 lines (split into smaller classes)
- ✓ Use sealed classes/interfaces for restricted hierarchies
- ✓ Use `object` for singletons
- ✓ Use `companion object` for factory methods
- ✓ Prefer composition over inheritance

### Coroutines
- ✓ Always use structured concurrency (never `GlobalScope`)
- ✓ Use `suspend` functions for async operations
- ✓ Make suspend functions main-safe (use `withContext(Dispatchers.IO)`)
- ✓ Use `coroutineScope` or `supervisorScope` for concurrent work
- ✓ Always handle `CancellationException` properly (don't catch it)
- ✓ Set timeouts for long-running operations (`withTimeout`)
- ✓ Use `Flow` for streams of data (not Channels for simple cases)

### Collections & Sequences
- ✓ Use sequences for large collections with multiple chained operations
- ✓ Prefer collection operators (map, filter) over imperative loops
- ✓ Use `mapNotNull`, `filterNotNull` to handle nullability
- ✓ Use `groupBy`, `partition` for categorization
- ✓ Understand lazy (sequences) vs eager (collections) evaluation

### Scope Functions
- ✓ Use `let` for null checks and transformations (`?.let { }`)
- ✓ Use `apply` for object configuration (returns receiver)
- ✓ Use `run` for computing results (returns lambda result)
- ✓ Use `also` for side effects like logging (returns receiver)
- ✓ Use `with` for calling multiple methods on object (non-extension)

### When Expression
- ✓ Prefer `when` over `if-else` chains
- ✓ Use exhaustive `when` with sealed classes (no else needed)
- ✓ Cover all enum cases in `when` (compiler enforced)
- ✓ Use `when` as expression (returns value) when possible

---

## Validation & Input Handling

### Recommended Libraries
- **kotlinx.serialization**: Type-safe JSON serialization
- **Bean Validation (JSR-380)**: Constraint validation (Spring Boot)
- **Arrow**: Functional validation and error handling

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
- **JUnit 5**: Industry standard (verbose but mature)
- **Kotest**: Kotlin-idiomatic, multiple test styles
- **MockK**: Mocking library (mocks final classes)
- **Turbine**: Flow testing library

### Guardrails
- ✓ Test files: `*Test.kt` or descriptive names
- ✓ Use descriptive test names with backticks: `` `should return user when id is valid` ``
- ✓ Prefer Kotest's `shouldBe` over JUnit's `assertEquals`
- ✓ Use MockK for mocking (not Mockito for Kotlin)
- ✓ Use `coEvery`, `coVerify` for suspend functions
- ✓ Test coroutines with `runTest` (from kotlinx-coroutines-test)
- ✓ Coverage target: >80% for business logic

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

    "should handle concurrent requests" {
        runTest {
            val users = listOf(
                User("1", "user1@example.com"),
                User("2", "user2@example.com")
            )
            coEvery { repository.findById(any()) } returnsMany users

            val results = service.getUsersConcurrently(listOf("1", "2"))

            results.size shouldBe 2
        }
    }
})
```

---

## Tooling

### Essential Tools
- **ktlint**: Code formatter (opinionated)
- **detekt**: Static analysis (code smells, complexity)
- **Gradle Kotlin DSL**: Build configuration
- **KSP**: Kotlin Symbol Processing (faster than KAPT)
- **Kotlin compiler**: Type checking

### Configuration Files
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

ktlint {
    version.set("1.1.1")
    android.set(false)
}

detekt {
    config.setFrom("$projectDir/detekt.yml")
    buildUponDefaultConfig = true
}
```

```yaml
# .editorconfig (for ktlint)
[*.{kt,kts}]
indent_size = 4
max_line_length = 100
insert_final_newline = true
ij_kotlin_allow_trailing_comma = true
ij_kotlin_allow_trailing_comma_on_call_site = true
```

```yaml
# detekt.yml
complexity:
  ComplexMethod:
    threshold: 10
  LongMethod:
    threshold: 50
  TooManyFunctions:
    thresholdInFiles: 15

style:
  MagicNumber:
    active: true
  MaxLineLength:
    maxLineLength: 100
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

### ❌ Don't Do This
```kotlin
// Using !!
val value = nullableValue!!

// Using var when val works
var user = User("John")

// Not handling nullable types
fun process(value: String?) {
    println(value.length) // Compilation error
}

// GlobalScope for coroutines
GlobalScope.launch {
    doWork()
}

// Mutable collections as public properties
class UserManager {
    val users = mutableListOf<User>()
}

// Imperative loops
val names = mutableListOf<String>()
for (user in users) {
    if (user.age > 18) {
        names.add(user.name)
    }
}
```

### ✅ Do This Instead
```kotlin
// Safe call or Elvis operator
val value = nullableValue ?: defaultValue

// Immutable by default
val user = User("John")

// Proper null handling
fun process(value: String?) {
    value?.let { println(it.length) }
}

// Structured concurrency
suspend fun doWork() = coroutineScope {
    launch {
        // Work here
    }
}

// Immutable public API
class UserManager {
    private val _users = mutableListOf<User>()
    val users: List<User> get() = _users.toList()
}

// Functional operators
val names = users
    .filter { it.age > 18 }
    .map { it.name }
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

// Service with coroutines
@Service
class UserService(private val repository: UserRepository) {
    suspend fun findById(id: String): User? = withContext(Dispatchers.IO) {
        repository.findById(id).orElse(null)
    }

    suspend fun create(request: UserCreateRequest): User = withContext(Dispatchers.IO) {
        val user = User(email = request.email, age = request.age, name = request.name)
        repository.save(user)
    }
}
```

### Ktor Server
```kotlin
fun Application.module() {
    install(ContentNegotiation) {
        json(Json {
            prettyPrint = true
            isLenient = true
        })
    }

    install(CallLogging)

    routing {
        route("/api/users") {
            get("/{id}") {
                val id = call.parameters["id"] ?: return@get call.respond(
                    HttpStatusCode.BadRequest,
                    "Missing id"
                )

                val user = userService.findById(id)
                if (user != null) {
                    call.respond(user)
                } else {
                    call.respond(HttpStatusCode.NotFound, "User not found")
                }
            }

            post {
                val request = call.receive<UserCreateRequest>()
                val user = userService.create(request)
                call.respond(HttpStatusCode.Created, user)
            }
        }
    }
}
```

---

## Project Structure

### Android
```
app/
├── src/
│   ├── main/
│   │   ├── java/com/example/
│   │   │   ├── data/          # Repositories, data sources
│   │   │   ├── domain/        # Use cases, models
│   │   │   ├── ui/            # Composables, ViewModels
│   │   │   │   ├── screens/
│   │   │   │   └── components/
│   │   │   └── di/            # Hilt modules
│   │   └── res/
│   └── test/
├── build.gradle.kts
└── proguard-rules.pro
```

### Spring Boot
```
src/
├── main/
│   ├── kotlin/com/example/
│   │   ├── controller/        # REST controllers
│   │   ├── service/           # Business logic
│   │   ├── repository/        # Data access
│   │   ├── model/             # Domain models
│   │   ├── dto/               # Data transfer objects
│   │   ├── config/            # Configuration classes
│   │   └── Application.kt     # Main entry point
│   └── resources/
│       ├── application.yml
│       └── db/migration/      # Flyway migrations
└── test/
```

---

## Performance Considerations

### Optimization Guardrails
- ✓ Use sequences for large collections with multiple operations
- ✓ Use inline functions for lambdas passed to higher-order functions
- ✓ Use `by lazy` for expensive one-time initialization
- ✓ Avoid creating unnecessary objects in loops
- ✓ Use `buildList`, `buildMap` for collection builders
- ✓ Profile before optimizing (Android Profiler, JProfiler, YourKit)

### Example
```kotlin
// Use sequences for large datasets
val result = hugeList
    .asSequence()
    .filter { it.isActive }
    .map { it.name }
    .take(10)
    .toList()

// Inline functions for performance-critical code
inline fun <T> measureTime(block: () -> T): Pair<T, Long> {
    val start = System.currentTimeMillis()
    val result = block()
    val time = System.currentTimeMillis() - start
    return result to time
}

// Lazy initialization
val expensiveResource by lazy {
    // Computed only once, when first accessed
    loadExpensiveResource()
}
```

---

## Security Best Practices

### Guardrails
- ✓ Never hardcode secrets (use environment variables or secure storage)
- ✓ Validate all user inputs with data class `init` blocks or Bean Validation
- ✓ Use parameterized queries (prevent SQL injection)
- ✓ Hash passwords with bcrypt or Argon2
- ✓ Enable CSRF protection (Spring Security)
- ✓ Use HTTPS/TLS for production
- ✓ Run `./gradlew dependencyCheckAnalyze` for vulnerability scanning
- ✓ Use ProGuard/R8 for Android (obfuscation)

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

// Input validation
@Serializable
data class LoginRequest(
    val email: String,
    val password: String
) {
    init {
        require(email.matches(Regex("^[A-Za-z0-9+_.-]+@(.+)$"))) {
            "Invalid email format"
        }
        require(password.length >= 8) {
            "Password must be at least 8 characters"
        }
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
- [Spring Boot Kotlin Support](https://docs.spring.io/spring-boot/docs/current/reference/html/features.html#features.kotlin)
- [Ktor Documentation](https://ktor.io/docs/)
