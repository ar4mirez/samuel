# Swift Guide

> **Applies to**: Swift 5.9+, iOS, macOS, SwiftUI, UIKit, Server-Side Swift

---

## Core Principles

1. **Safety First**: Optionals, strong typing, memory safety
2. **Protocol-Oriented**: Prefer protocols over inheritance
3. **Value Types**: Prefer structs over classes
4. **Modern Concurrency**: async/await, actors, structured concurrency
5. **Clarity Over Brevity**: Clear, expressive code

---

## Language-Specific Guardrails

### Swift Version & Setup
- ✓ Use Swift 5.9+ (latest stable)
- ✓ Use Swift Package Manager (SPM) for dependencies
- ✓ Enable strict concurrency checking
- ✓ Target minimum iOS version appropriate for features used

### Code Style (Swift API Design Guidelines)
- ✓ Follow [Swift API Design Guidelines](https://swift.org/documentation/api-design-guidelines/)
- ✓ Run SwiftLint before every commit
- ✓ Use `camelCase` for functions, variables, properties
- ✓ Use `PascalCase` for types (structs, classes, enums, protocols)
- ✓ 4-space indentation
- ✓ Line length: 120 characters
- ✓ Name methods based on their side effects (mutating vs non-mutating)

### Optionals
- ✓ Prefer `guard let` for early exits
- ✓ Use `if let` for optional binding in simple cases
- ✓ Avoid force unwrapping (`!`) except in tests or known-safe situations
- ✓ Use nil-coalescing: `value ?? defaultValue`
- ✓ Use optional chaining: `object?.property?.method()`
- ✓ Use `map` and `flatMap` for optional transformations

### Value Types vs Reference Types
- ✓ Prefer `struct` over `class` (value semantics, thread-safe)
- ✓ Use `class` when identity matters or for UIKit/AppKit subclassing
- ✓ Use `actor` for shared mutable state in concurrent contexts
- ✓ Make properties `let` by default, use `var` only when needed
- ✓ Use `final` on classes that won't be subclassed

### Modern Concurrency (Swift 5.5+)
- ✓ Use `async`/`await` over completion handlers
- ✓ Use `Task` for structured concurrency
- ✓ Use `actor` for thread-safe mutable state
- ✓ Use `@MainActor` for UI updates
- ✓ Handle `Task` cancellation properly
- ✓ Use `AsyncSequence` for streams of values
- ✓ Avoid `DispatchQueue` unless necessary for legacy code

### Error Handling
- ✓ Use `throws` for functions that can fail
- ✓ Use `do-catch` for error handling
- ✓ Create custom `Error` types for domain errors
- ✓ Use `Result` for async operations without async/await
- ✓ Avoid `try?` unless you truly don't need error details
- ✓ Use `try!` only when failure is impossible

---

## Project Structure

### iOS App Structure
```
MyApp/
├── MyApp.xcodeproj
├── MyApp/
│   ├── App/
│   │   ├── MyApp.swift           # @main entry point
│   │   └── AppDelegate.swift     # If using UIKit lifecycle
│   ├── Features/
│   │   ├── Authentication/
│   │   │   ├── Views/
│   │   │   ├── ViewModels/
│   │   │   └── Models/
│   │   └── Home/
│   ├── Core/
│   │   ├── Network/
│   │   ├── Storage/
│   │   └── Extensions/
│   ├── Resources/
│   │   ├── Assets.xcassets
│   │   └── Localizable.strings
│   └── Supporting Files/
│       └── Info.plist
├── MyAppTests/
├── MyAppUITests/
└── Package.swift                  # SPM dependencies
```

### Swift Package Structure
```
MyPackage/
├── Package.swift
├── Sources/
│   └── MyPackage/
│       ├── MyPackage.swift
│       └── Internal/
├── Tests/
│   └── MyPackageTests/
└── README.md
```

---

## SwiftUI Patterns

### View with ViewModel
```swift
import SwiftUI

struct UserListView: View {
    @StateObject private var viewModel = UserListViewModel()

    var body: some View {
        NavigationStack {
            content
                .navigationTitle("Users")
                .task {
                    await viewModel.loadUsers()
                }
                .refreshable {
                    await viewModel.loadUsers()
                }
        }
    }

    @ViewBuilder
    private var content: some View {
        switch viewModel.state {
        case .loading:
            ProgressView()
        case .loaded(let users):
            userList(users)
        case .error(let message):
            errorView(message)
        }
    }

    private func userList(_ users: [User]) -> some View {
        List(users) { user in
            NavigationLink(value: user) {
                UserRowView(user: user)
            }
        }
        .navigationDestination(for: User.self) { user in
            UserDetailView(user: user)
        }
    }

    private func errorView(_ message: String) -> some View {
        ContentUnavailableView(
            "Error",
            systemImage: "exclamationmark.triangle",
            description: Text(message)
        )
    }
}
```

### ViewModel with @Observable (iOS 17+)
```swift
import Foundation

@Observable
final class UserListViewModel {
    private(set) var state: ViewState = .loading

    private let userService: UserServiceProtocol

    init(userService: UserServiceProtocol = UserService()) {
        self.userService = userService
    }

    @MainActor
    func loadUsers() async {
        state = .loading

        do {
            let users = try await userService.fetchUsers()
            state = .loaded(users)
        } catch {
            state = .error(error.localizedDescription)
        }
    }
}

enum ViewState {
    case loading
    case loaded([User])
    case error(String)
}
```

### ViewModel with ObservableObject (iOS 13+)
```swift
import Foundation
import Combine

@MainActor
final class UserListViewModel: ObservableObject {
    @Published private(set) var state: ViewState = .loading

    private let userService: UserServiceProtocol

    init(userService: UserServiceProtocol = UserService()) {
        self.userService = userService
    }

    func loadUsers() async {
        state = .loading

        do {
            let users = try await userService.fetchUsers()
            state = .loaded(users)
        } catch {
            state = .error(error.localizedDescription)
        }
    }
}
```

---

## Testing

### Frameworks
- **XCTest**: Built-in testing framework
- **Swift Testing** (Xcode 16+): Modern testing framework
- **Quick/Nimble**: BDD-style testing

### Guardrails
- ✓ Test files: `*Tests.swift`
- ✓ Test methods: `func test*()` (XCTest) or `@Test` (Swift Testing)
- ✓ Use descriptive names: `test_createUser_withValidData_returnsUser()`
- ✓ Use protocols for dependency injection (testability)
- ✓ Mock external dependencies
- ✓ Coverage target: >80% for business logic
- ✓ Test async code with expectations or async/await

### Example (XCTest)
```swift
import XCTest
@testable import MyApp

final class UserServiceTests: XCTestCase {
    private var sut: UserService!
    private var mockAPIClient: MockAPIClient!

    override func setUp() {
        super.setUp()
        mockAPIClient = MockAPIClient()
        sut = UserService(apiClient: mockAPIClient)
    }

    override func tearDown() {
        sut = nil
        mockAPIClient = nil
        super.tearDown()
    }

    func test_fetchUsers_withSuccessResponse_returnsUsers() async throws {
        // Given
        let expectedUsers = [User(id: "1", email: "test@example.com")]
        mockAPIClient.result = .success(expectedUsers)

        // When
        let users = try await sut.fetchUsers()

        // Then
        XCTAssertEqual(users.count, 1)
        XCTAssertEqual(users.first?.email, "test@example.com")
    }

    func test_fetchUsers_withNetworkError_throwsError() async {
        // Given
        mockAPIClient.result = .failure(NetworkError.connectionFailed)

        // When/Then
        do {
            _ = try await sut.fetchUsers()
            XCTFail("Expected error to be thrown")
        } catch {
            XCTAssertTrue(error is NetworkError)
        }
    }
}

// Mock
final class MockAPIClient: APIClientProtocol {
    var result: Result<[User], Error> = .success([])

    func fetch<T: Decodable>(_ endpoint: Endpoint) async throws -> T {
        switch result {
        case .success(let value):
            return value as! T
        case .failure(let error):
            throw error
        }
    }
}
```

### Example (Swift Testing - Xcode 16+)
```swift
import Testing
@testable import MyApp

@Suite("UserService Tests")
struct UserServiceTests {
    let mockAPIClient = MockAPIClient()
    let sut: UserService

    init() {
        sut = UserService(apiClient: mockAPIClient)
    }

    @Test("Fetches users successfully")
    func fetchUsersSuccess() async throws {
        mockAPIClient.result = .success([User(id: "1", email: "test@example.com")])

        let users = try await sut.fetchUsers()

        #expect(users.count == 1)
        #expect(users.first?.email == "test@example.com")
    }

    @Test("Throws error on network failure")
    func fetchUsersNetworkError() async {
        mockAPIClient.result = .failure(NetworkError.connectionFailed)

        await #expect(throws: NetworkError.self) {
            try await sut.fetchUsers()
        }
    }

    @Test("Validates email format", arguments: ["", "invalid", "test@"])
    func invalidEmailValidation(email: String) {
        #expect(throws: ValidationError.self) {
            try User(email: email).validate()
        }
    }
}
```

---

## Tooling

### Essential Tools
- **SwiftLint**: Code style enforcement
- **SwiftFormat**: Code formatting
- **swift-testing**: Modern testing (Xcode 16+)
- **xcbeautify**: Readable Xcode build output

### Configuration
```yaml
# .swiftlint.yml
disabled_rules:
  - trailing_whitespace

opt_in_rules:
  - empty_count
  - explicit_init
  - fatal_error_message
  - first_where
  - force_unwrapping
  - implicitly_unwrapped_optional
  - private_action
  - private_outlet
  - redundant_nil_coalescing

line_length:
  warning: 120
  error: 150

type_body_length:
  warning: 300
  error: 400

file_length:
  warning: 400
  error: 500

function_body_length:
  warning: 40
  error: 60

cyclomatic_complexity:
  warning: 10
  error: 15

nesting:
  type_level: 2
  function_level: 3

custom_rules:
  no_print:
    name: "No print statements"
    regex: "print\\("
    message: "Use Logger instead of print"
    severity: warning
```

### Pre-Commit Commands
```bash
# Lint
swiftlint lint --strict

# Format
swiftformat .

# Build
xcodebuild build -scheme MyApp -destination 'platform=iOS Simulator,name=iPhone 15'

# Test
xcodebuild test -scheme MyApp -destination 'platform=iOS Simulator,name=iPhone 15'

# Swift Package
swift build
swift test
```

---

## Common Pitfalls

### Don't Do This
```swift
// Force unwrapping without justification
let name = optionalName!

// Retain cycles in closures
class ViewModel {
    var onComplete: (() -> Void)?

    func setup() {
        onComplete = {
            self.doSomething() // Retain cycle!
        }
    }
}

// Blocking main thread
DispatchQueue.main.sync { } // Can deadlock

// Ignoring errors
let data = try? loadData() // Error details lost

// Massive view controllers/views
struct ContentView: View {
    var body: some View {
        // 500+ lines of view code
    }
}
```

### Do This Instead
```swift
// Safe unwrapping
guard let name = optionalName else {
    return
}

// Capture list to prevent retain cycle
class ViewModel {
    var onComplete: (() -> Void)?

    func setup() {
        onComplete = { [weak self] in
            self?.doSomething()
        }
    }
}

// Async/await instead of blocking
@MainActor
func updateUI() async {
    // UI updates here
}

// Proper error handling
do {
    let data = try loadData()
} catch {
    logger.error("Failed to load: \(error)")
}

// Extract subviews
struct ContentView: View {
    var body: some View {
        VStack {
            HeaderView()
            ContentListView()
            FooterView()
        }
    }
}
```

---

## Networking

### Modern Async Networking
```swift
import Foundation

protocol APIClientProtocol {
    func fetch<T: Decodable>(_ endpoint: Endpoint) async throws -> T
}

final class APIClient: APIClientProtocol {
    private let session: URLSession
    private let decoder: JSONDecoder

    init(session: URLSession = .shared) {
        self.session = session
        self.decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        decoder.dateDecodingStrategy = .iso8601
    }

    func fetch<T: Decodable>(_ endpoint: Endpoint) async throws -> T {
        let request = try endpoint.urlRequest()

        let (data, response) = try await session.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse else {
            throw NetworkError.invalidResponse
        }

        guard (200...299).contains(httpResponse.statusCode) else {
            throw NetworkError.httpError(statusCode: httpResponse.statusCode)
        }

        return try decoder.decode(T.self, from: data)
    }
}

enum NetworkError: Error {
    case invalidURL
    case invalidResponse
    case httpError(statusCode: Int)
    case connectionFailed
}

struct Endpoint {
    let path: String
    let method: HTTPMethod
    let headers: [String: String]
    let body: Data?

    func urlRequest() throws -> URLRequest {
        guard let url = URL(string: "https://api.example.com" + path) else {
            throw NetworkError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = method.rawValue
        request.allHTTPHeaderFields = headers
        request.httpBody = body
        return request
    }
}

enum HTTPMethod: String {
    case get = "GET"
    case post = "POST"
    case put = "PUT"
    case delete = "DELETE"
}
```

---

## Concurrency with Actors

### Actor for Thread-Safe State
```swift
actor UserCache {
    private var cache: [String: User] = [:]

    func get(_ id: String) -> User? {
        cache[id]
    }

    func set(_ user: User) {
        cache[user.id] = user
    }

    func clear() {
        cache.removeAll()
    }
}

// Usage
let cache = UserCache()

Task {
    await cache.set(user)
    let cached = await cache.get("123")
}
```

### MainActor for UI Updates
```swift
@MainActor
final class UserViewModel: ObservableObject {
    @Published private(set) var users: [User] = []

    func loadUsers() async {
        // This runs on main actor, safe for @Published
        let fetchedUsers = await userService.fetchUsers()
        users = fetchedUsers
    }
}
```

---

## Performance Considerations

### Optimization Guardrails
- ✓ Use `lazy` for expensive computed properties
- ✓ Use value types (structs) for thread safety and copy-on-write
- ✓ Avoid unnecessary allocations in hot paths
- ✓ Use `@inlinable` for performance-critical generic functions
- ✓ Profile with Instruments before optimizing
- ✓ Use `ContiguousArray` for performance-critical array operations
- ✓ Prefer `Set` or `Dictionary` for lookups over `Array`

### Example
```swift
// Lazy initialization
struct DataProcessor {
    lazy var expensiveResource: Resource = {
        Resource.load()
    }()
}

// Value type with copy-on-write
struct LargeData {
    private var storage: Storage

    private class Storage {
        var data: [Int]
        init(data: [Int]) { self.data = data }
    }

    mutating func append(_ value: Int) {
        if !isKnownUniquelyReferenced(&storage) {
            storage = Storage(data: storage.data)
        }
        storage.data.append(value)
    }
}

// Efficient collection operations
let ids = Set(users.map(\.id)) // O(1) lookups
let userById = Dictionary(uniqueKeysWithValues: users.map { ($0.id, $0) })
```

---

## Security Best Practices

### Guardrails
- ✓ Use Keychain for sensitive data (tokens, passwords)
- ✓ Use App Transport Security (ATS) - HTTPS only
- ✓ Validate server certificates for SSL pinning (when required)
- ✓ Use `CryptoKit` for cryptographic operations
- ✓ Never hardcode secrets in source code
- ✓ Use `Data Protection` for file encryption
- ✓ Validate and sanitize user input

### Example
```swift
import Security
import CryptoKit

// Keychain storage
final class KeychainManager {
    enum KeychainError: Error {
        case duplicateItem
        case itemNotFound
        case unexpectedStatus(OSStatus)
    }

    func save(_ data: Data, for key: String) throws {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key,
            kSecValueData as String: data,
            kSecAttrAccessible as String: kSecAttrAccessibleWhenUnlockedThisDeviceOnly
        ]

        let status = SecItemAdd(query as CFDictionary, nil)

        guard status == errSecSuccess else {
            if status == errSecDuplicateItem {
                throw KeychainError.duplicateItem
            }
            throw KeychainError.unexpectedStatus(status)
        }
    }

    func retrieve(for key: String) throws -> Data {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key,
            kSecReturnData as String: true
        ]

        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)

        guard status == errSecSuccess, let data = result as? Data else {
            throw KeychainError.itemNotFound
        }

        return data
    }
}

// Hashing with CryptoKit
func hashPassword(_ password: String) -> String {
    let data = Data(password.utf8)
    let hash = SHA256.hash(data: data)
    return hash.compactMap { String(format: "%02x", $0) }.joined()
}
```

---

## References

- [Swift Documentation](https://docs.swift.org/swift-book/)
- [Swift API Design Guidelines](https://swift.org/documentation/api-design-guidelines/)
- [Apple Human Interface Guidelines](https://developer.apple.com/design/human-interface-guidelines/)
- [SwiftUI Documentation](https://developer.apple.com/documentation/swiftui)
- [Swift Concurrency](https://docs.swift.org/swift-book/documentation/the-swift-programming-language/concurrency/)
- [SwiftLint](https://github.com/realm/SwiftLint)
- [Point-Free (Swift videos)](https://www.pointfree.co/)
