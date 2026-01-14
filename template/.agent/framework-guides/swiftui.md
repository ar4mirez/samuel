# SwiftUI Framework Guide

> **Framework**: SwiftUI 5.0+ (iOS 17+, macOS 14+)
> **Language**: Swift 5.9+
> **Type**: Declarative UI Framework
> **Platform**: iOS, macOS, watchOS, tvOS, visionOS

---

## Overview

SwiftUI is Apple's modern declarative UI framework for building native apps across all Apple platforms. It uses a reactive data-driven approach with automatic UI updates.

**Use SwiftUI when:**
- Building new Apple platform apps
- Need cross-platform Apple development (iOS, macOS, watchOS, tvOS)
- Want declarative, reactive UI patterns
- Building with latest iOS features
- Rapid prototyping and iteration

**Consider alternatives when:**
- Supporting iOS < 14 (use UIKit)
- Need fine-grained UIKit control
- Complex custom drawing/animations
- Existing large UIKit codebase

---

## Project Structure

```
MyApp/
├── MyApp.xcodeproj
├── MyApp/
│   ├── App/
│   │   ├── MyApp.swift              # @main entry point
│   │   └── AppDelegate.swift        # Optional UIKit lifecycle
│   ├── Features/
│   │   ├── Authentication/
│   │   │   ├── Views/
│   │   │   │   ├── LoginView.swift
│   │   │   │   └── SignUpView.swift
│   │   │   ├── ViewModels/
│   │   │   │   └── AuthViewModel.swift
│   │   │   └── Models/
│   │   │       └── User.swift
│   │   ├── Home/
│   │   │   ├── Views/
│   │   │   ├── ViewModels/
│   │   │   └── Components/
│   │   └── Settings/
│   ├── Core/
│   │   ├── Network/
│   │   │   ├── APIClient.swift
│   │   │   └── Endpoints.swift
│   │   ├── Storage/
│   │   │   └── UserDefaults+Extensions.swift
│   │   ├── Extensions/
│   │   │   ├── View+Extensions.swift
│   │   │   └── Color+Extensions.swift
│   │   └── Utilities/
│   ├── Shared/
│   │   ├── Components/
│   │   │   ├── LoadingView.swift
│   │   │   ├── ErrorView.swift
│   │   │   └── PrimaryButton.swift
│   │   └── Modifiers/
│   │       └── CardModifier.swift
│   ├── Resources/
│   │   ├── Assets.xcassets
│   │   └── Localizable.strings
│   └── Preview Content/
├── MyAppTests/
├── MyAppUITests/
└── Package.swift                     # SPM dependencies
```

---

## App Entry Point

### Basic App Structure

```swift
import SwiftUI

@main
struct MyApp: App {
    // App-level state
    @StateObject private var appState = AppState()
    @StateObject private var authManager = AuthManager()

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(appState)
                .environmentObject(authManager)
        }
    }
}

// Content view with navigation
struct ContentView: View {
    @EnvironmentObject var authManager: AuthManager

    var body: some View {
        Group {
            if authManager.isAuthenticated {
                MainTabView()
            } else {
                AuthenticationView()
            }
        }
        .animation(.easeInOut, value: authManager.isAuthenticated)
    }
}
```

### App with Scene Phases

```swift
import SwiftUI

@main
struct MyApp: App {
    @Environment(\.scenePhase) private var scenePhase
    @StateObject private var dataController = DataController()

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environment(\.managedObjectContext, dataController.container.viewContext)
        }
        .onChange(of: scenePhase) { _, newPhase in
            switch newPhase {
            case .active:
                print("App became active")
            case .inactive:
                print("App became inactive")
            case .background:
                dataController.save()
            @unknown default:
                break
            }
        }
    }
}
```

---

## Views and Components

### Basic View Structure

```swift
import SwiftUI

struct UserProfileView: View {
    let user: User
    @State private var isEditing = false

    var body: some View {
        VStack(spacing: 16) {
            profileHeader
            statsSection
            actionButtons
        }
        .padding()
        .navigationTitle("Profile")
        .sheet(isPresented: $isEditing) {
            EditProfileView(user: user)
        }
    }

    // MARK: - Subviews

    private var profileHeader: some View {
        VStack(spacing: 8) {
            AsyncImage(url: URL(string: user.avatarURL)) { image in
                image
                    .resizable()
                    .aspectRatio(contentMode: .fill)
            } placeholder: {
                Circle()
                    .fill(Color.gray.opacity(0.3))
            }
            .frame(width: 100, height: 100)
            .clipShape(Circle())

            Text(user.name)
                .font(.title2)
                .fontWeight(.semibold)

            Text(user.email)
                .font(.subheadline)
                .foregroundStyle(.secondary)
        }
    }

    private var statsSection: some View {
        HStack(spacing: 32) {
            StatView(title: "Posts", value: user.postCount)
            StatView(title: "Followers", value: user.followerCount)
            StatView(title: "Following", value: user.followingCount)
        }
        .padding(.vertical)
    }

    private var actionButtons: some View {
        HStack(spacing: 16) {
            Button("Edit Profile") {
                isEditing = true
            }
            .buttonStyle(.borderedProminent)

            Button("Share") {
                // Share action
            }
            .buttonStyle(.bordered)
        }
    }
}

// Reusable component
struct StatView: View {
    let title: String
    let value: Int

    var body: some View {
        VStack(spacing: 4) {
            Text("\(value)")
                .font(.headline)
            Text(title)
                .font(.caption)
                .foregroundStyle(.secondary)
        }
    }
}
```

### Custom View Modifiers

```swift
import SwiftUI

// Card modifier
struct CardModifier: ViewModifier {
    func body(content: Content) -> some View {
        content
            .background(Color(.systemBackground))
            .cornerRadius(12)
            .shadow(color: .black.opacity(0.1), radius: 8, x: 0, y: 4)
    }
}

extension View {
    func cardStyle() -> some View {
        modifier(CardModifier())
    }
}

// Loading overlay modifier
struct LoadingModifier: ViewModifier {
    let isLoading: Bool

    func body(content: Content) -> some View {
        ZStack {
            content
                .disabled(isLoading)
                .blur(radius: isLoading ? 2 : 0)

            if isLoading {
                ProgressView()
                    .scaleEffect(1.5)
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
                    .background(Color.black.opacity(0.2))
            }
        }
    }
}

extension View {
    func loading(_ isLoading: Bool) -> some View {
        modifier(LoadingModifier(isLoading: isLoading))
    }
}

// Usage
struct ContentView: View {
    @State private var isLoading = false

    var body: some View {
        VStack {
            Text("Content")
                .padding()
                .cardStyle()
        }
        .loading(isLoading)
    }
}
```

---

## State Management

### @Observable (iOS 17+)

```swift
import SwiftUI
import Observation

@Observable
final class UserViewModel {
    var user: User?
    var isLoading = false
    var errorMessage: String?

    private let userService: UserServiceProtocol

    init(userService: UserServiceProtocol = UserService()) {
        self.userService = userService
    }

    func loadUser(id: String) async {
        isLoading = true
        errorMessage = nil

        do {
            user = try await userService.fetchUser(id: id)
        } catch {
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }

    func updateUser(name: String) async {
        guard var currentUser = user else { return }

        currentUser.name = name

        do {
            user = try await userService.updateUser(currentUser)
        } catch {
            errorMessage = error.localizedDescription
        }
    }
}

// View using @Observable
struct UserView: View {
    @State private var viewModel = UserViewModel()
    let userId: String

    var body: some View {
        Group {
            if viewModel.isLoading {
                ProgressView()
            } else if let user = viewModel.user {
                UserProfileView(user: user)
            } else if let error = viewModel.errorMessage {
                ErrorView(message: error) {
                    Task { await viewModel.loadUser(id: userId) }
                }
            }
        }
        .task {
            await viewModel.loadUser(id: userId)
        }
    }
}
```

### ObservableObject (iOS 13+)

```swift
import SwiftUI
import Combine

final class AuthViewModel: ObservableObject {
    @Published var email = ""
    @Published var password = ""
    @Published var isAuthenticated = false
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let authService: AuthServiceProtocol
    private var cancellables = Set<AnyCancellable>()

    init(authService: AuthServiceProtocol = AuthService()) {
        self.authService = authService
    }

    var isFormValid: Bool {
        !email.isEmpty && email.contains("@") && password.count >= 8
    }

    @MainActor
    func login() async {
        guard isFormValid else {
            errorMessage = "Please enter valid credentials"
            return
        }

        isLoading = true
        errorMessage = nil

        do {
            let user = try await authService.login(email: email, password: password)
            isAuthenticated = true
            print("Logged in: \(user.name)")
        } catch {
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }

    func logout() {
        authService.logout()
        isAuthenticated = false
        email = ""
        password = ""
    }
}

// View with ObservableObject
struct LoginView: View {
    @StateObject private var viewModel = AuthViewModel()

    var body: some View {
        VStack(spacing: 20) {
            TextField("Email", text: $viewModel.email)
                .textFieldStyle(.roundedBorder)
                .textContentType(.emailAddress)
                .autocapitalization(.none)

            SecureField("Password", text: $viewModel.password)
                .textFieldStyle(.roundedBorder)
                .textContentType(.password)

            if let error = viewModel.errorMessage {
                Text(error)
                    .foregroundStyle(.red)
                    .font(.caption)
            }

            Button("Login") {
                Task { await viewModel.login() }
            }
            .buttonStyle(.borderedProminent)
            .disabled(!viewModel.isFormValid || viewModel.isLoading)
        }
        .padding()
        .loading(viewModel.isLoading)
    }
}
```

### Environment and Dependency Injection

```swift
import SwiftUI

// Environment key for custom dependency
struct APIClientKey: EnvironmentKey {
    static let defaultValue: APIClientProtocol = APIClient()
}

extension EnvironmentValues {
    var apiClient: APIClientProtocol {
        get { self[APIClientKey.self] }
        set { self[APIClientKey.self] = newValue }
    }
}

// Using environment dependency
struct PostListView: View {
    @Environment(\.apiClient) private var apiClient
    @State private var posts: [Post] = []

    var body: some View {
        List(posts) { post in
            PostRowView(post: post)
        }
        .task {
            posts = try? await apiClient.fetchPosts() ?? []
        }
    }
}

// Injecting mock in previews
#Preview {
    PostListView()
        .environment(\.apiClient, MockAPIClient())
}
```

---

## Navigation

### NavigationStack (iOS 16+)

```swift
import SwiftUI

// Navigation path for type-safe navigation
struct MainView: View {
    @State private var navigationPath = NavigationPath()

    var body: some View {
        NavigationStack(path: $navigationPath) {
            HomeView(navigationPath: $navigationPath)
                .navigationDestination(for: User.self) { user in
                    UserDetailView(user: user)
                }
                .navigationDestination(for: Post.self) { post in
                    PostDetailView(post: post)
                }
                .navigationDestination(for: Route.self) { route in
                    destinationView(for: route)
                }
        }
    }

    @ViewBuilder
    private func destinationView(for route: Route) -> some View {
        switch route {
        case .settings:
            SettingsView()
        case .profile(let userId):
            ProfileView(userId: userId)
        case .notifications:
            NotificationsView()
        }
    }
}

// Route enum for navigation
enum Route: Hashable {
    case settings
    case profile(userId: String)
    case notifications
}

// Home view with navigation
struct HomeView: View {
    @Binding var navigationPath: NavigationPath
    @State private var users: [User] = []

    var body: some View {
        List(users) { user in
            NavigationLink(value: user) {
                UserRowView(user: user)
            }
        }
        .navigationTitle("Home")
        .toolbar {
            ToolbarItem(placement: .topBarTrailing) {
                Button {
                    navigationPath.append(Route.settings)
                } label: {
                    Image(systemName: "gear")
                }
            }
        }
    }
}
```

### TabView with Navigation

```swift
import SwiftUI

struct MainTabView: View {
    @State private var selectedTab = Tab.home

    enum Tab: String, CaseIterable {
        case home, search, profile

        var icon: String {
            switch self {
            case .home: return "house"
            case .search: return "magnifyingglass"
            case .profile: return "person"
            }
        }

        var title: String {
            rawValue.capitalized
        }
    }

    var body: some View {
        TabView(selection: $selectedTab) {
            NavigationStack {
                HomeView()
            }
            .tabItem {
                Label(Tab.home.title, systemImage: Tab.home.icon)
            }
            .tag(Tab.home)

            NavigationStack {
                SearchView()
            }
            .tabItem {
                Label(Tab.search.title, systemImage: Tab.search.icon)
            }
            .tag(Tab.search)

            NavigationStack {
                ProfileView()
            }
            .tabItem {
                Label(Tab.profile.title, systemImage: Tab.profile.icon)
            }
            .tag(Tab.profile)
        }
    }
}
```

---

## Lists and Data Display

### Modern List Patterns

```swift
import SwiftUI

struct PostListView: View {
    @State private var posts: [Post] = []
    @State private var isLoading = false
    @State private var searchText = ""

    var filteredPosts: [Post] {
        if searchText.isEmpty {
            return posts
        }
        return posts.filter { $0.title.localizedCaseInsensitiveContains(searchText) }
    }

    var body: some View {
        List {
            ForEach(filteredPosts) { post in
                PostRowView(post: post)
                    .swipeActions(edge: .trailing) {
                        Button(role: .destructive) {
                            deletePost(post)
                        } label: {
                            Label("Delete", systemImage: "trash")
                        }

                        Button {
                            archivePost(post)
                        } label: {
                            Label("Archive", systemImage: "archivebox")
                        }
                        .tint(.orange)
                    }
                    .swipeActions(edge: .leading) {
                        Button {
                            toggleFavorite(post)
                        } label: {
                            Label("Favorite", systemImage: post.isFavorite ? "star.fill" : "star")
                        }
                        .tint(.yellow)
                    }
            }
        }
        .listStyle(.plain)
        .searchable(text: $searchText, prompt: "Search posts")
        .refreshable {
            await loadPosts()
        }
        .overlay {
            if filteredPosts.isEmpty && !isLoading {
                ContentUnavailableView.search(text: searchText)
            }
        }
        .task {
            await loadPosts()
        }
    }

    private func loadPosts() async {
        isLoading = true
        // Fetch posts
        isLoading = false
    }

    private func deletePost(_ post: Post) {
        posts.removeAll { $0.id == post.id }
    }

    private func archivePost(_ post: Post) {
        // Archive logic
    }

    private func toggleFavorite(_ post: Post) {
        if let index = posts.firstIndex(where: { $0.id == post.id }) {
            posts[index].isFavorite.toggle()
        }
    }
}

// Custom row view
struct PostRowView: View {
    let post: Post

    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            HStack {
                Text(post.title)
                    .font(.headline)

                Spacer()

                if post.isFavorite {
                    Image(systemName: "star.fill")
                        .foregroundStyle(.yellow)
                }
            }

            Text(post.excerpt)
                .font(.subheadline)
                .foregroundStyle(.secondary)
                .lineLimit(2)

            HStack {
                Label("\(post.likeCount)", systemImage: "heart")
                Label("\(post.commentCount)", systemImage: "bubble.right")
                Spacer()
                Text(post.createdAt, style: .relative)
            }
            .font(.caption)
            .foregroundStyle(.tertiary)
        }
        .padding(.vertical, 4)
    }
}
```

### LazyVGrid and LazyHGrid

```swift
import SwiftUI

struct PhotoGridView: View {
    let photos: [Photo]

    private let columns = [
        GridItem(.adaptive(minimum: 100, maximum: 200), spacing: 4)
    ]

    var body: some View {
        ScrollView {
            LazyVGrid(columns: columns, spacing: 4) {
                ForEach(photos) { photo in
                    AsyncImage(url: photo.thumbnailURL) { image in
                        image
                            .resizable()
                            .aspectRatio(1, contentMode: .fill)
                    } placeholder: {
                        Rectangle()
                            .fill(Color.gray.opacity(0.3))
                    }
                    .aspectRatio(1, contentMode: .fit)
                    .clipped()
                }
            }
            .padding(4)
        }
    }
}

// Horizontal scroll with sections
struct CategoryView: View {
    let categories: [Category]

    var body: some View {
        ScrollView {
            LazyVStack(alignment: .leading, spacing: 24, pinnedViews: .sectionHeaders) {
                ForEach(categories) { category in
                    Section {
                        ScrollView(.horizontal, showsIndicators: false) {
                            LazyHStack(spacing: 12) {
                                ForEach(category.items) { item in
                                    ItemCardView(item: item)
                                        .frame(width: 150)
                                }
                            }
                            .padding(.horizontal)
                        }
                    } header: {
                        Text(category.name)
                            .font(.title2)
                            .fontWeight(.bold)
                            .padding(.horizontal)
                            .padding(.vertical, 8)
                            .frame(maxWidth: .infinity, alignment: .leading)
                            .background(.bar)
                    }
                }
            }
        }
    }
}
```

---

## Networking

### Async/Await API Client

```swift
import Foundation

protocol APIClientProtocol {
    func fetch<T: Decodable>(_ endpoint: Endpoint) async throws -> T
    func send<T: Decodable, U: Encodable>(_ endpoint: Endpoint, body: U) async throws -> T
}

final class APIClient: APIClientProtocol {
    private let session: URLSession
    private let decoder: JSONDecoder
    private let encoder: JSONEncoder
    private let baseURL: URL

    init(
        baseURL: URL = URL(string: "https://api.example.com")!,
        session: URLSession = .shared
    ) {
        self.baseURL = baseURL
        self.session = session

        self.decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        decoder.dateDecodingStrategy = .iso8601

        self.encoder = JSONEncoder()
        encoder.keyEncodingStrategy = .convertToSnakeCase
    }

    func fetch<T: Decodable>(_ endpoint: Endpoint) async throws -> T {
        let request = try buildRequest(for: endpoint)
        return try await execute(request)
    }

    func send<T: Decodable, U: Encodable>(_ endpoint: Endpoint, body: U) async throws -> T {
        var request = try buildRequest(for: endpoint)
        request.httpBody = try encoder.encode(body)
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        return try await execute(request)
    }

    private func buildRequest(for endpoint: Endpoint) throws -> URLRequest {
        guard let url = URL(string: endpoint.path, relativeTo: baseURL) else {
            throw APIError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = endpoint.method.rawValue
        request.timeoutInterval = 30

        for (key, value) in endpoint.headers {
            request.setValue(value, forHTTPHeaderField: key)
        }

        return request
    }

    private func execute<T: Decodable>(_ request: URLRequest) async throws -> T {
        let (data, response) = try await session.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }

        guard (200...299).contains(httpResponse.statusCode) else {
            throw APIError.httpError(statusCode: httpResponse.statusCode, data: data)
        }

        return try decoder.decode(T.self, from: data)
    }
}

// Endpoint definition
struct Endpoint {
    let path: String
    let method: HTTPMethod
    var headers: [String: String] = [:]

    enum HTTPMethod: String {
        case get = "GET"
        case post = "POST"
        case put = "PUT"
        case delete = "DELETE"
    }
}

extension Endpoint {
    static func users() -> Endpoint {
        Endpoint(path: "/users", method: .get)
    }

    static func user(id: String) -> Endpoint {
        Endpoint(path: "/users/\(id)", method: .get)
    }

    static func createUser() -> Endpoint {
        Endpoint(path: "/users", method: .post)
    }
}

// API errors
enum APIError: LocalizedError {
    case invalidURL
    case invalidResponse
    case httpError(statusCode: Int, data: Data)
    case decodingError(Error)

    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Invalid URL"
        case .invalidResponse:
            return "Invalid response from server"
        case .httpError(let statusCode, _):
            return "HTTP error: \(statusCode)"
        case .decodingError(let error):
            return "Decoding error: \(error.localizedDescription)"
        }
    }
}
```

---

## Forms and Input

### Form with Validation

```swift
import SwiftUI

struct RegistrationView: View {
    @State private var formData = RegistrationForm()
    @State private var isSubmitting = false
    @State private var showingAlert = false
    @State private var alertMessage = ""

    var body: some View {
        NavigationStack {
            Form {
                Section("Personal Information") {
                    TextField("Full Name", text: $formData.name)
                        .textContentType(.name)

                    TextField("Email", text: $formData.email)
                        .textContentType(.emailAddress)
                        .autocapitalization(.none)
                        .keyboardType(.emailAddress)

                    DatePicker("Date of Birth", selection: $formData.dateOfBirth, displayedComponents: .date)
                }

                Section("Account") {
                    SecureField("Password", text: $formData.password)
                        .textContentType(.newPassword)

                    SecureField("Confirm Password", text: $formData.confirmPassword)
                        .textContentType(.newPassword)

                    if !formData.passwordsMatch && !formData.confirmPassword.isEmpty {
                        Text("Passwords do not match")
                            .font(.caption)
                            .foregroundStyle(.red)
                    }
                }

                Section("Preferences") {
                    Toggle("Receive Newsletter", isOn: $formData.receiveNewsletter)

                    Picker("Notification Frequency", selection: $formData.notificationFrequency) {
                        ForEach(NotificationFrequency.allCases) { frequency in
                            Text(frequency.rawValue).tag(frequency)
                        }
                    }
                }

                Section {
                    Button("Create Account") {
                        Task { await submit() }
                    }
                    .frame(maxWidth: .infinity)
                    .disabled(!formData.isValid || isSubmitting)
                }
            }
            .navigationTitle("Register")
            .loading(isSubmitting)
            .alert("Registration", isPresented: $showingAlert) {
                Button("OK", role: .cancel) { }
            } message: {
                Text(alertMessage)
            }
        }
    }

    private func submit() async {
        isSubmitting = true
        defer { isSubmitting = false }

        do {
            try await registerUser(formData)
            alertMessage = "Account created successfully!"
        } catch {
            alertMessage = error.localizedDescription
        }

        showingAlert = true
    }
}

// Form data model
struct RegistrationForm {
    var name = ""
    var email = ""
    var dateOfBirth = Date()
    var password = ""
    var confirmPassword = ""
    var receiveNewsletter = false
    var notificationFrequency = NotificationFrequency.daily

    var isValidEmail: Bool {
        email.contains("@") && email.contains(".")
    }

    var isValidPassword: Bool {
        password.count >= 8
    }

    var passwordsMatch: Bool {
        password == confirmPassword
    }

    var isValid: Bool {
        !name.isEmpty &&
        isValidEmail &&
        isValidPassword &&
        passwordsMatch
    }
}

enum NotificationFrequency: String, CaseIterable, Identifiable {
    case daily = "Daily"
    case weekly = "Weekly"
    case monthly = "Monthly"
    case never = "Never"

    var id: String { rawValue }
}
```

---

## Testing

### View Model Testing

```swift
import XCTest
@testable import MyApp

final class UserViewModelTests: XCTestCase {
    var sut: UserViewModel!
    var mockService: MockUserService!

    override func setUp() {
        super.setUp()
        mockService = MockUserService()
        sut = UserViewModel(userService: mockService)
    }

    override func tearDown() {
        sut = nil
        mockService = nil
        super.tearDown()
    }

    func test_loadUser_success_updatesUser() async {
        // Given
        let expectedUser = User(id: "1", name: "John", email: "john@example.com")
        mockService.userToReturn = expectedUser

        // When
        await sut.loadUser(id: "1")

        // Then
        XCTAssertEqual(sut.user, expectedUser)
        XCTAssertFalse(sut.isLoading)
        XCTAssertNil(sut.errorMessage)
    }

    func test_loadUser_failure_setsErrorMessage() async {
        // Given
        mockService.errorToThrow = APIError.invalidResponse

        // When
        await sut.loadUser(id: "1")

        // Then
        XCTAssertNil(sut.user)
        XCTAssertFalse(sut.isLoading)
        XCTAssertNotNil(sut.errorMessage)
    }
}

// Mock service
final class MockUserService: UserServiceProtocol {
    var userToReturn: User?
    var errorToThrow: Error?

    func fetchUser(id: String) async throws -> User {
        if let error = errorToThrow {
            throw error
        }
        guard let user = userToReturn else {
            throw APIError.invalidResponse
        }
        return user
    }
}
```

### UI Testing

```swift
import XCTest

final class LoginUITests: XCTestCase {
    var app: XCUIApplication!

    override func setUp() {
        super.setUp()
        continueAfterFailure = false
        app = XCUIApplication()
        app.launchArguments = ["UI_TESTING"]
        app.launch()
    }

    func test_login_withValidCredentials_navigatesToHome() {
        // Enter credentials
        let emailField = app.textFields["Email"]
        emailField.tap()
        emailField.typeText("test@example.com")

        let passwordField = app.secureTextFields["Password"]
        passwordField.tap()
        passwordField.typeText("password123")

        // Tap login button
        app.buttons["Login"].tap()

        // Verify navigation to home
        XCTAssertTrue(app.navigationBars["Home"].waitForExistence(timeout: 5))
    }

    func test_login_withInvalidCredentials_showsError() {
        let emailField = app.textFields["Email"]
        emailField.tap()
        emailField.typeText("invalid")

        let loginButton = app.buttons["Login"]
        XCTAssertFalse(loginButton.isEnabled)
    }
}
```

---

## Preview Provider

### Preview with Mock Data

```swift
import SwiftUI

struct UserProfileView_Previews: PreviewProvider {
    static var previews: some View {
        Group {
            // Default preview
            UserProfileView(user: .preview)
                .previewDisplayName("Default")

            // Dark mode
            UserProfileView(user: .preview)
                .preferredColorScheme(.dark)
                .previewDisplayName("Dark Mode")

            // Different device
            UserProfileView(user: .preview)
                .previewDevice("iPhone SE (3rd generation)")
                .previewDisplayName("iPhone SE")

            // Large text
            UserProfileView(user: .preview)
                .environment(\.sizeCategory, .accessibilityExtraLarge)
                .previewDisplayName("Large Text")
        }
    }
}

// Preview data extension
extension User {
    static var preview: User {
        User(
            id: "preview",
            name: "John Doe",
            email: "john@example.com",
            avatarURL: "https://example.com/avatar.jpg",
            postCount: 42,
            followerCount: 1234,
            followingCount: 567
        )
    }
}

// iOS 17+ Preview macro
#Preview("User Profile") {
    NavigationStack {
        UserProfileView(user: .preview)
    }
}

#Preview("Dark Mode") {
    NavigationStack {
        UserProfileView(user: .preview)
    }
    .preferredColorScheme(.dark)
}
```

---

## Best Practices

### Performance
- ✓ Use `@State` for view-local state
- ✓ Use `@StateObject` for owned objects (created once)
- ✓ Use `@ObservedObject` for passed-in objects
- ✓ Avoid expensive computations in body
- ✓ Use `Equatable` views for optimization
- ✓ Prefer `LazyVStack`/`LazyHStack` for long lists

### Architecture
- ✓ Keep views small and focused
- ✓ Extract reusable components
- ✓ Use view models for business logic
- ✓ Inject dependencies via environment
- ✓ Use protocols for testability

### Common Pitfalls
- ❌ Modifying state during view update
- ❌ Creating `@StateObject` in computed property
- ❌ Using `@ObservedObject` for owned objects
- ❌ Heavy computation in `body`
- ❌ Not handling optional binding properly

---

## References

- [SwiftUI Documentation](https://developer.apple.com/documentation/swiftui)
- [Apple Human Interface Guidelines](https://developer.apple.com/design/human-interface-guidelines/)
- [SwiftUI by Example](https://www.hackingwithswift.com/quick-start/swiftui)
- [WWDC SwiftUI Sessions](https://developer.apple.com/videos/swiftui)
