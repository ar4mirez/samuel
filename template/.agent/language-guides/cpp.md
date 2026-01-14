# C/C++ Guide

> **Applies to**: C17/C23, C++17/C++20/C++23, Embedded Systems, Game Development, Systems Programming

---

## Core Principles

1. **Memory Safety**: RAII, smart pointers, no raw owning pointers
2. **Modern C++**: Use C++17/20/23 features, avoid legacy patterns
3. **Zero-Cost Abstractions**: High-level code with predictable performance
4. **Compile-Time Validation**: constexpr, concepts, static_assert
5. **Resource Management**: Deterministic destruction, exception safety

---

## Language-Specific Guardrails

### Version & Setup
- ✓ Use C++17 minimum (C++20/23 recommended for new projects)
- ✓ Use CMake for build configuration
- ✓ Use vcpkg, Conan, or FetchContent for dependencies
- ✓ Enable all warnings: `-Wall -Wextra -Wpedantic` (GCC/Clang)
- ✓ Treat warnings as errors in CI: `-Werror`

### Code Style (Google C++ Style / Core Guidelines)
- ✓ Follow C++ Core Guidelines or Google C++ Style
- ✓ Run clang-format before every commit
- ✓ Use `snake_case` for functions and variables
- ✓ Use `PascalCase` for types (classes, structs, enums)
- ✓ Use `SCREAMING_SNAKE_CASE` for macros and constants
- ✓ 2 or 4-space indentation (be consistent)
- ✓ Line length: 80-120 characters

### Memory Management (CRITICAL)
- ✓ Use RAII for all resource management
- ✓ Prefer `std::unique_ptr` for single ownership
- ✓ Use `std::shared_ptr` only when shared ownership is needed
- ✓ Never use `new`/`delete` directly (use smart pointers or containers)
- ✓ Use `std::make_unique`/`std::make_shared`
- ✓ Avoid raw owning pointers (non-owning raw pointers are OK)
- ✓ Use `std::span` (C++20) for non-owning array views

### Modern C++ Features
- ✓ Use `auto` for complex types, explicit types for clarity
- ✓ Use range-based for loops: `for (const auto& item : container)`
- ✓ Use structured bindings: `auto [key, value] = pair;`
- ✓ Use `std::optional` instead of nullable pointers
- ✓ Use `std::variant` instead of unions
- ✓ Use `constexpr` for compile-time computation
- ✓ Use concepts (C++20) for template constraints
- ✓ Use `std::string_view` for non-owning string references

### Error Handling
- ✓ Use exceptions for exceptional conditions
- ✓ Use `std::expected` (C++23) or `std::optional` for expected failures
- ✓ Create custom exception types derived from `std::exception`
- ✓ Ensure exception safety (basic guarantee minimum)
- ✓ Use `noexcept` for functions that won't throw
- ✓ RAII ensures cleanup even with exceptions

### Const Correctness
- ✓ Use `const` by default for variables
- ✓ Use `const` for method parameters that won't be modified
- ✓ Use `const` member functions for non-mutating operations
- ✓ Use `constexpr` for compile-time constants
- ✓ Prefer `std::string_view` over `const std::string&` for read-only

### Thread Safety
- ✓ Use `std::mutex` and `std::lock_guard` for synchronization
- ✓ Use `std::atomic` for simple shared state
- ✓ Use `std::jthread` (C++20) or properly manage `std::thread`
- ✓ Prefer immutable data for thread safety
- ✓ Use `std::async` or thread pools for task-based concurrency

---

## Project Structure

### CMake Project Layout
```
myproject/
├── CMakeLists.txt
├── cmake/                    # CMake modules
├── include/
│   └── myproject/           # Public headers
│       ├── myproject.hpp
│       └── types.hpp
├── src/                     # Implementation files
│   ├── myproject.cpp
│   └── internal/            # Private headers
├── tests/
│   ├── CMakeLists.txt
│   └── test_myproject.cpp
├── examples/
├── docs/
├── vcpkg.json               # Or conanfile.txt
└── README.md
```

### CMakeLists.txt Example
```cmake
cmake_minimum_required(VERSION 3.20)
project(myproject VERSION 1.0.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# Warnings
add_compile_options(-Wall -Wextra -Wpedantic)

# Library
add_library(myproject
    src/myproject.cpp
)

target_include_directories(myproject
    PUBLIC
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        $<INSTALL_INTERFACE:include>
)

# Tests
enable_testing()
add_subdirectory(tests)
```

---

## Smart Pointers & RAII

### Ownership Patterns
```cpp
#include <memory>
#include <vector>

// Unique ownership (most common)
auto create_user() -> std::unique_ptr<User> {
    return std::make_unique<User>("test@example.com", 25);
}

// Transfer ownership
void process_user(std::unique_ptr<User> user) {
    // Takes ownership, user destroyed at end of function
}

// Shared ownership (use sparingly)
auto create_shared_cache() -> std::shared_ptr<Cache> {
    return std::make_shared<Cache>();
}

// Non-owning observation (raw pointer is fine)
void observe_user(const User* user) {
    if (user) {
        // Read-only access
    }
}

// Non-owning reference (preferred over raw pointer)
void process_user_ref(const User& user) {
    // Guaranteed non-null
}

// RAII resource wrapper
class FileHandle {
public:
    explicit FileHandle(const std::string& path)
        : handle_(std::fopen(path.c_str(), "r")) {
        if (!handle_) {
            throw std::runtime_error("Failed to open file");
        }
    }

    ~FileHandle() {
        if (handle_) {
            std::fclose(handle_);
        }
    }

    // Non-copyable
    FileHandle(const FileHandle&) = delete;
    FileHandle& operator=(const FileHandle&) = delete;

    // Movable
    FileHandle(FileHandle&& other) noexcept : handle_(other.handle_) {
        other.handle_ = nullptr;
    }

    FileHandle& operator=(FileHandle&& other) noexcept {
        if (this != &other) {
            if (handle_) std::fclose(handle_);
            handle_ = other.handle_;
            other.handle_ = nullptr;
        }
        return *this;
    }

private:
    FILE* handle_;
};
```

---

## Error Handling

### Custom Exceptions
```cpp
#include <stdexcept>
#include <string>

class AppError : public std::runtime_error {
public:
    explicit AppError(const std::string& message)
        : std::runtime_error(message) {}
};

class NotFoundError : public AppError {
public:
    NotFoundError(const std::string& resource, const std::string& id)
        : AppError(resource + " with ID " + id + " not found")
        , resource_(resource)
        , id_(id) {}

    [[nodiscard]] const std::string& resource() const noexcept { return resource_; }
    [[nodiscard]] const std::string& id() const noexcept { return id_; }

private:
    std::string resource_;
    std::string id_;
};

// Usage
auto get_user(const std::string& id) -> User {
    auto user = db_.find_user(id);
    if (!user) {
        throw NotFoundError("User", id);
    }
    return *user;
}
```

### std::expected (C++23) or std::optional
```cpp
#include <optional>
#include <expected>  // C++23

// Using std::optional for absence
auto find_user(const std::string& id) -> std::optional<User> {
    auto it = users_.find(id);
    if (it == users_.end()) {
        return std::nullopt;
    }
    return it->second;
}

// Using std::expected (C++23) for errors
enum class ParseError {
    InvalidFormat,
    OutOfRange,
    Empty
};

auto parse_int(std::string_view str) -> std::expected<int, ParseError> {
    if (str.empty()) {
        return std::unexpected(ParseError::Empty);
    }

    try {
        return std::stoi(std::string(str));
    } catch (const std::invalid_argument&) {
        return std::unexpected(ParseError::InvalidFormat);
    } catch (const std::out_of_range&) {
        return std::unexpected(ParseError::OutOfRange);
    }
}

// Usage
auto result = parse_int("42");
if (result) {
    std::cout << *result << '\n';
} else {
    // Handle error
}
```

---

## Testing

### Frameworks
- **GoogleTest**: Industry standard
- **Catch2**: Header-only, modern syntax
- **doctest**: Lightweight, fast compilation

### Guardrails
- ✓ Test files: `*_test.cpp` or `test_*.cpp`
- ✓ Use descriptive test names
- ✓ Use fixtures for setup/teardown
- ✓ Test edge cases: empty, null, boundary values
- ✓ Use mocking sparingly (prefer dependency injection)
- ✓ Coverage target: >80% for business logic

### Example (GoogleTest)
```cpp
#include <gtest/gtest.h>
#include <gmock/gmock.h>
#include "myproject/user_service.hpp"

class UserServiceTest : public ::testing::Test {
protected:
    void SetUp() override {
        repository_ = std::make_unique<MockUserRepository>();
        service_ = std::make_unique<UserService>(repository_.get());
    }

    std::unique_ptr<MockUserRepository> repository_;
    std::unique_ptr<UserService> service_;
};

TEST_F(UserServiceTest, CreateUser_WithValidData_ReturnsUser) {
    // Arrange
    UserCreate request{"test@example.com", 25, "user"};
    User expected{1, "test@example.com", 25, "user"};

    EXPECT_CALL(*repository_, save(::testing::_))
        .WillOnce(::testing::Return(expected));

    // Act
    auto result = service_->create(request);

    // Assert
    EXPECT_EQ(result.email, "test@example.com");
    EXPECT_EQ(result.age, 25);
}

TEST_F(UserServiceTest, CreateUser_WithInvalidEmail_ThrowsException) {
    UserCreate request{"invalid", 25, "user"};

    EXPECT_THROW(service_->create(request), ValidationError);
}

// Parameterized test
class InvalidEmailTest : public UserServiceTest,
                         public ::testing::WithParamInterface<std::string> {};

TEST_P(InvalidEmailTest, RejectsInvalidEmails) {
    UserCreate request{GetParam(), 25, "user"};
    EXPECT_THROW(service_->create(request), ValidationError);
}

INSTANTIATE_TEST_SUITE_P(
    InvalidEmails,
    InvalidEmailTest,
    ::testing::Values("", " ", "invalid", "test@", "@example.com")
);
```

### Example (Catch2)
```cpp
#include <catch2/catch_test_macros.hpp>
#include <catch2/generators/catch_generators.hpp>
#include "myproject/user_service.hpp"

TEST_CASE("UserService creates users", "[user][service]") {
    auto repository = std::make_unique<InMemoryUserRepository>();
    UserService service(repository.get());

    SECTION("with valid data returns user") {
        UserCreate request{"test@example.com", 25, "user"};

        auto result = service.create(request);

        REQUIRE(result.email == "test@example.com");
        REQUIRE(result.age == 25);
    }

    SECTION("with invalid email throws exception") {
        auto email = GENERATE("", " ", "invalid", "test@");
        UserCreate request{email, 25, "user"};

        REQUIRE_THROWS_AS(service.create(request), ValidationError);
    }
}
```

---

## Tooling

### Essential Tools
- **clang-format**: Code formatting
- **clang-tidy**: Static analysis
- **cppcheck**: Additional static analysis
- **AddressSanitizer**: Memory error detection
- **Valgrind**: Memory leak detection
- **gcov/lcov**: Code coverage

### Configuration
```yaml
# .clang-format
BasedOnStyle: Google
IndentWidth: 4
ColumnLimit: 100
AllowShortFunctionsOnASingleLine: Inline
BreakBeforeBraces: Attach
PointerAlignment: Left
```

```yaml
# .clang-tidy
Checks: >
  -*,
  bugprone-*,
  clang-analyzer-*,
  cppcoreguidelines-*,
  modernize-*,
  performance-*,
  readability-*,
  -modernize-use-trailing-return-type,
  -readability-identifier-length

WarningsAsErrors: '*'
HeaderFilterRegex: '.*'

CheckOptions:
  - key: readability-identifier-naming.ClassCase
    value: CamelCase
  - key: readability-identifier-naming.FunctionCase
    value: lower_case
  - key: readability-identifier-naming.VariableCase
    value: lower_case
  - key: readability-identifier-naming.ConstantCase
    value: UPPER_CASE
```

### Pre-Commit Commands
```bash
# Format
clang-format -i src/*.cpp include/**/*.hpp

# Static analysis
clang-tidy src/*.cpp -- -I include

# Build
cmake -B build -DCMAKE_BUILD_TYPE=Release
cmake --build build

# Test
ctest --test-dir build --output-on-failure

# Build with sanitizers (debug)
cmake -B build-asan -DCMAKE_BUILD_TYPE=Debug \
    -DCMAKE_CXX_FLAGS="-fsanitize=address,undefined"
cmake --build build-asan
```

---

## Common Pitfalls

### Don't Do This
```cpp
// Raw owning pointer
User* create_user() {
    return new User();  // Who deletes this?
}

// Manual memory management
void process() {
    int* arr = new int[100];
    // ... code that might throw
    delete[] arr;  // May not be reached!
}

// Returning reference to local
const std::string& get_name() {
    std::string name = "test";
    return name;  // Dangling reference!
}

// Using C-style casts
auto ptr = (Derived*)base_ptr;

// Uninitialized variables
int count;
process(count);  // Undefined behavior
```

### Do This Instead
```cpp
// Smart pointer with clear ownership
auto create_user() -> std::unique_ptr<User> {
    return std::make_unique<User>();
}

// RAII with containers
void process() {
    std::vector<int> arr(100);
    // Automatically cleaned up, even if exception thrown
}

// Return by value (move semantics)
std::string get_name() {
    std::string name = "test";
    return name;  // Moved, not copied
}

// Use C++ casts
auto ptr = dynamic_cast<Derived*>(base_ptr);
if (!ptr) { /* handle failure */ }

// Initialize all variables
int count = 0;
process(count);
```

---

## Modern C++ Patterns

### Concepts (C++20)
```cpp
#include <concepts>

// Define a concept
template<typename T>
concept Printable = requires(T t) {
    { std::cout << t } -> std::same_as<std::ostream&>;
};

// Use concept as constraint
template<Printable T>
void print(const T& value) {
    std::cout << value << '\n';
}

// Concept for container
template<typename T>
concept Container = requires(T t) {
    { t.begin() } -> std::input_iterator;
    { t.end() } -> std::input_iterator;
    { t.size() } -> std::convertible_to<std::size_t>;
};

template<Container C>
void process_container(const C& container) {
    for (const auto& item : container) {
        // Process item
    }
}
```

### Ranges (C++20)
```cpp
#include <ranges>
#include <vector>
#include <algorithm>

void process_users(std::vector<User>& users) {
    // Filter and transform with ranges
    auto adult_emails = users
        | std::views::filter([](const User& u) { return u.age >= 18; })
        | std::views::transform([](const User& u) { return u.email; });

    for (const auto& email : adult_emails) {
        std::cout << email << '\n';
    }

    // Sort with ranges
    std::ranges::sort(users, {}, &User::age);
}
```

### std::variant and std::visit
```cpp
#include <variant>
#include <string>

using Value = std::variant<int, double, std::string>;

auto process_value(const Value& v) -> std::string {
    return std::visit([](auto&& arg) -> std::string {
        using T = std::decay_t<decltype(arg)>;
        if constexpr (std::is_same_v<T, int>) {
            return "int: " + std::to_string(arg);
        } else if constexpr (std::is_same_v<T, double>) {
            return "double: " + std::to_string(arg);
        } else {
            return "string: " + arg;
        }
    }, v);
}
```

---

## Performance Considerations

### Optimization Guardrails
- ✓ Prefer stack allocation over heap
- ✓ Use `reserve()` for vectors when size is known
- ✓ Use move semantics for expensive-to-copy objects
- ✓ Avoid unnecessary copies (use const references)
- ✓ Use `std::string_view` for read-only string operations
- ✓ Profile with perf, Valgrind, or platform profilers
- ✓ Use `[[likely]]`/`[[unlikely]]` hints for hot paths

### Example
```cpp
// Reserve capacity
std::vector<User> users;
users.reserve(1000);  // Avoid reallocations

// Move instead of copy
std::vector<std::string> collect_names(std::vector<User> users) {
    std::vector<std::string> names;
    names.reserve(users.size());

    for (auto& user : users) {
        names.push_back(std::move(user.name));  // Move, don't copy
    }
    return names;
}

// string_view for non-owning
void process(std::string_view text) {
    // No allocation, works with string literals too
}

// Pass by const reference for read-only
void analyze(const std::vector<Data>& data) {
    // No copy
}
```

---

## Security Best Practices

### Guardrails
- ✓ Use smart pointers (prevent memory leaks, use-after-free)
- ✓ Use bounds-checked containers (`at()` or ranges)
- ✓ Validate all external input
- ✓ Use secure string functions (`std::string`, not C strings)
- ✓ Enable AddressSanitizer in CI
- ✓ Use ASLR, stack canaries, and other compiler protections
- ✓ Keep dependencies updated

### Build Flags for Security
```cmake
# Security hardening
target_compile_options(myproject PRIVATE
    -fstack-protector-strong
    -D_FORTIFY_SOURCE=2
    -fPIE
)

target_link_options(myproject PRIVATE
    -pie
    -Wl,-z,relro,-z,now
)

# Sanitizers for development
if(CMAKE_BUILD_TYPE STREQUAL "Debug")
    target_compile_options(myproject PRIVATE
        -fsanitize=address,undefined
        -fno-omit-frame-pointer
    )
    target_link_options(myproject PRIVATE
        -fsanitize=address,undefined
    )
endif()
```

---

## C-Specific Guidelines

When working with C (C17/C23):

### Guardrails
- ✓ Use C17 or C23 standard
- ✓ Always check return values
- ✓ Use `const` for read-only parameters
- ✓ Use `restrict` for non-aliasing pointers
- ✓ Initialize all variables
- ✓ Use `static` for internal linkage
- ✓ Use `_Noreturn` / `[[noreturn]]` for functions that don't return

### Memory Safety in C
```c
#include <stdlib.h>
#include <string.h>

// Always check allocation
void* safe_malloc(size_t size) {
    void* ptr = malloc(size);
    if (!ptr && size > 0) {
        abort();  // Or handle error appropriately
    }
    return ptr;
}

// Use snprintf, not sprintf
char buffer[256];
int result = snprintf(buffer, sizeof(buffer), "User: %s", username);
if (result < 0 || (size_t)result >= sizeof(buffer)) {
    // Handle truncation or error
}

// Clean up pattern
typedef struct {
    char* data;
    size_t size;
} Buffer;

Buffer* buffer_create(size_t size) {
    Buffer* buf = safe_malloc(sizeof(Buffer));
    buf->data = safe_malloc(size);
    buf->size = size;
    return buf;
}

void buffer_destroy(Buffer* buf) {
    if (buf) {
        free(buf->data);
        free(buf);
    }
}
```

---

## References

- [C++ Core Guidelines](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines)
- [Google C++ Style Guide](https://google.github.io/styleguide/cppguide.html)
- [cppreference.com](https://en.cppreference.com/)
- [Effective Modern C++ (Book)](https://www.oreilly.com/library/view/effective-modern-c/9781491908419/)
- [C++ Weekly (YouTube)](https://www.youtube.com/c/lefticus1)
- [CppCon Talks](https://www.youtube.com/user/CppCon)
- [clang-tidy checks](https://clang.llvm.org/extra/clang-tidy/checks/list.html)
