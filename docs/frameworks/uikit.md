# UIKit Framework Guide

> **Framework**: UIKit (iOS 12+)
> **Language**: Swift 5.0+
> **Type**: Imperative UI Framework
> **Platform**: iOS, tvOS, Mac Catalyst

---

## Overview

UIKit is Apple's traditional imperative UI framework for building iOS applications. It provides fine-grained control over UI elements and is the foundation of iOS app development.

**Use UIKit when:**
- Supporting older iOS versions (< iOS 14)
- Need fine-grained control over animations
- Working with complex custom UI
- Integrating with existing UIKit codebases
- Advanced collection view layouts
- Complex gesture handling

**Consider alternatives when:**
- Building new apps targeting iOS 15+ (consider SwiftUI)
- Simple CRUD applications
- Need rapid prototyping

---

## Project Structure

```
MyApp/
├── MyApp.xcodeproj
├── MyApp/
│   ├── Application/
│   │   ├── AppDelegate.swift
│   │   ├── SceneDelegate.swift
│   │   └── AppCoordinator.swift
│   ├── Features/
│   │   ├── Authentication/
│   │   │   ├── Controllers/
│   │   │   │   ├── LoginViewController.swift
│   │   │   │   └── SignUpViewController.swift
│   │   │   ├── ViewModels/
│   │   │   │   └── AuthViewModel.swift
│   │   │   ├── Views/
│   │   │   │   └── LoginFormView.swift
│   │   │   └── Coordinator/
│   │   │       └── AuthCoordinator.swift
│   │   ├── Home/
│   │   │   ├── Controllers/
│   │   │   ├── ViewModels/
│   │   │   ├── Views/
│   │   │   └── Cells/
│   │   └── Profile/
│   ├── Core/
│   │   ├── Network/
│   │   ├── Storage/
│   │   ├── Extensions/
│   │   └── Utilities/
│   ├── Shared/
│   │   ├── Views/
│   │   ├── Cells/
│   │   └── Protocols/
│   └── Resources/
│       ├── Assets.xcassets
│       ├── Storyboards/
│       └── Localizable.strings
├── MyAppTests/
└── MyAppUITests/
```

---

## Application Lifecycle

### SceneDelegate (iOS 13+)

```swift
import UIKit

class SceneDelegate: UIResponder, UIWindowSceneDelegate {
    var window: UIWindow?
    var appCoordinator: AppCoordinator?

    func scene(
        _ scene: UIScene,
        willConnectTo session: UISceneSession,
        options connectionOptions: UIScene.ConnectionOptions
    ) {
        guard let windowScene = (scene as? UIWindowScene) else { return }

        let window = UIWindow(windowScene: windowScene)
        self.window = window

        // Setup coordinator
        let navigationController = UINavigationController()
        appCoordinator = AppCoordinator(navigationController: navigationController)

        window.rootViewController = navigationController
        window.makeKeyAndVisible()

        appCoordinator?.start()
    }

    func sceneDidEnterBackground(_ scene: UIScene) {
        // Save state
        CoreDataManager.shared.saveContext()
    }
}
```

### AppDelegate

```swift
import UIKit

@main
class AppDelegate: UIResponder, UIApplicationDelegate {
    func application(
        _ application: UIApplication,
        didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?
    ) -> Bool {
        // Configure appearance
        configureAppearance()

        // Configure third-party SDKs
        configureAnalytics()

        return true
    }

    private func configureAppearance() {
        // Navigation bar appearance
        let navBarAppearance = UINavigationBarAppearance()
        navBarAppearance.configureWithOpaqueBackground()
        navBarAppearance.backgroundColor = .systemBackground
        navBarAppearance.titleTextAttributes = [.foregroundColor: UIColor.label]

        UINavigationBar.appearance().standardAppearance = navBarAppearance
        UINavigationBar.appearance().scrollEdgeAppearance = navBarAppearance

        // Tab bar appearance
        let tabBarAppearance = UITabBarAppearance()
        tabBarAppearance.configureWithOpaqueBackground()
        UITabBar.appearance().standardAppearance = tabBarAppearance
        UITabBar.appearance().scrollEdgeAppearance = tabBarAppearance
    }

    // MARK: - UISceneSession Lifecycle

    func application(
        _ application: UIApplication,
        configurationForConnecting connectingSceneSession: UISceneSession,
        options: UIScene.ConnectionOptions
    ) -> UISceneConfiguration {
        return UISceneConfiguration(name: "Default Configuration", sessionRole: connectingSceneSession.role)
    }
}
```

---

## View Controllers

### Base View Controller

```swift
import UIKit

class BaseViewController: UIViewController {
    // MARK: - Properties

    private lazy var loadingView: LoadingView = {
        let view = LoadingView()
        view.translatesAutoresizingMaskIntoConstraints = false
        return view
    }()

    // MARK: - Lifecycle

    override func viewDidLoad() {
        super.viewDidLoad()
        setupUI()
        setupConstraints()
        setupBindings()
    }

    // MARK: - Setup Methods (Override in subclasses)

    func setupUI() {
        view.backgroundColor = .systemBackground
    }

    func setupConstraints() {
        // Override in subclasses
    }

    func setupBindings() {
        // Override in subclasses
    }

    // MARK: - Loading State

    func showLoading() {
        view.addSubview(loadingView)
        NSLayoutConstraint.activate([
            loadingView.topAnchor.constraint(equalTo: view.topAnchor),
            loadingView.leadingAnchor.constraint(equalTo: view.leadingAnchor),
            loadingView.trailingAnchor.constraint(equalTo: view.trailingAnchor),
            loadingView.bottomAnchor.constraint(equalTo: view.bottomAnchor)
        ])
        loadingView.startAnimating()
    }

    func hideLoading() {
        loadingView.stopAnimating()
        loadingView.removeFromSuperview()
    }

    // MARK: - Error Handling

    func showError(_ error: Error) {
        let alert = UIAlertController(
            title: "Error",
            message: error.localizedDescription,
            preferredStyle: .alert
        )
        alert.addAction(UIAlertAction(title: "OK", style: .default))
        present(alert, animated: true)
    }
}
```

### Feature View Controller

```swift
import UIKit
import Combine

final class UserListViewController: BaseViewController {
    // MARK: - Properties

    private let viewModel: UserListViewModel
    private var cancellables = Set<AnyCancellable>()

    weak var coordinator: UserCoordinator?

    // MARK: - UI Elements

    private lazy var tableView: UITableView = {
        let table = UITableView(frame: .zero, style: .plain)
        table.translatesAutoresizingMaskIntoConstraints = false
        table.register(UserCell.self, forCellReuseIdentifier: UserCell.identifier)
        table.delegate = self
        table.dataSource = self
        table.rowHeight = UITableView.automaticDimension
        table.estimatedRowHeight = 80
        table.refreshControl = refreshControl
        return table
    }()

    private lazy var refreshControl: UIRefreshControl = {
        let control = UIRefreshControl()
        control.addTarget(self, action: #selector(handleRefresh), for: .valueChanged)
        return control
    }()

    private lazy var emptyStateView: EmptyStateView = {
        let view = EmptyStateView()
        view.translatesAutoresizingMaskIntoConstraints = false
        view.configure(
            image: UIImage(systemName: "person.3"),
            title: "No Users",
            message: "There are no users to display"
        )
        view.isHidden = true
        return view
    }()

    // MARK: - Initialization

    init(viewModel: UserListViewModel = UserListViewModel()) {
        self.viewModel = viewModel
        super.init(nibName: nil, bundle: nil)
    }

    required init?(coder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Setup

    override func setupUI() {
        super.setupUI()
        title = "Users"

        view.addSubview(tableView)
        view.addSubview(emptyStateView)

        navigationItem.rightBarButtonItem = UIBarButtonItem(
            barButtonSystemItem: .add,
            target: self,
            action: #selector(addUserTapped)
        )
    }

    override func setupConstraints() {
        NSLayoutConstraint.activate([
            tableView.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor),
            tableView.leadingAnchor.constraint(equalTo: view.leadingAnchor),
            tableView.trailingAnchor.constraint(equalTo: view.trailingAnchor),
            tableView.bottomAnchor.constraint(equalTo: view.bottomAnchor),

            emptyStateView.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            emptyStateView.centerYAnchor.constraint(equalTo: view.centerYAnchor),
            emptyStateView.leadingAnchor.constraint(greaterThanOrEqualTo: view.leadingAnchor, constant: 32),
            emptyStateView.trailingAnchor.constraint(lessThanOrEqualTo: view.trailingAnchor, constant: -32)
        ])
    }

    override func setupBindings() {
        viewModel.$users
            .receive(on: DispatchQueue.main)
            .sink { [weak self] users in
                self?.tableView.reloadData()
                self?.emptyStateView.isHidden = !users.isEmpty
            }
            .store(in: &cancellables)

        viewModel.$isLoading
            .receive(on: DispatchQueue.main)
            .sink { [weak self] isLoading in
                if isLoading {
                    self?.showLoading()
                } else {
                    self?.hideLoading()
                    self?.refreshControl.endRefreshing()
                }
            }
            .store(in: &cancellables)

        viewModel.$error
            .compactMap { $0 }
            .receive(on: DispatchQueue.main)
            .sink { [weak self] error in
                self?.showError(error)
            }
            .store(in: &cancellables)
    }

    // MARK: - Lifecycle

    override func viewDidLoad() {
        super.viewDidLoad()
        viewModel.loadUsers()
    }

    // MARK: - Actions

    @objc private func handleRefresh() {
        viewModel.loadUsers()
    }

    @objc private func addUserTapped() {
        coordinator?.showAddUser()
    }
}

// MARK: - UITableViewDataSource

extension UserListViewController: UITableViewDataSource {
    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return viewModel.users.count
    }

    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        guard let cell = tableView.dequeueReusableCell(
            withIdentifier: UserCell.identifier,
            for: indexPath
        ) as? UserCell else {
            return UITableViewCell()
        }

        let user = viewModel.users[indexPath.row]
        cell.configure(with: user)
        return cell
    }
}

// MARK: - UITableViewDelegate

extension UserListViewController: UITableViewDelegate {
    func tableView(_ tableView: UITableView, didSelectRowAt indexPath: IndexPath) {
        tableView.deselectRow(at: indexPath, animated: true)
        let user = viewModel.users[indexPath.row]
        coordinator?.showUserDetail(user)
    }

    func tableView(
        _ tableView: UITableView,
        trailingSwipeActionsConfigurationForRowAt indexPath: IndexPath
    ) -> UISwipeActionsConfiguration? {
        let deleteAction = UIContextualAction(style: .destructive, title: "Delete") { [weak self] _, _, completion in
            self?.viewModel.deleteUser(at: indexPath.row)
            completion(true)
        }
        deleteAction.image = UIImage(systemName: "trash")

        return UISwipeActionsConfiguration(actions: [deleteAction])
    }
}
```

---

## MVVM Architecture

### View Model with Combine

```swift
import Foundation
import Combine

final class UserListViewModel {
    // MARK: - Published Properties

    @Published private(set) var users: [User] = []
    @Published private(set) var isLoading = false
    @Published private(set) var error: Error?

    // MARK: - Dependencies

    private let userService: UserServiceProtocol
    private var cancellables = Set<AnyCancellable>()

    // MARK: - Initialization

    init(userService: UserServiceProtocol = UserService()) {
        self.userService = userService
    }

    // MARK: - Public Methods

    func loadUsers() {
        isLoading = true
        error = nil

        Task { @MainActor in
            do {
                users = try await userService.fetchUsers()
            } catch {
                self.error = error
            }
            isLoading = false
        }
    }

    func deleteUser(at index: Int) {
        let user = users[index]
        users.remove(at: index)

        Task {
            do {
                try await userService.deleteUser(id: user.id)
            } catch {
                // Restore on failure
                await MainActor.run {
                    users.insert(user, at: index)
                    self.error = error
                }
            }
        }
    }

    func searchUsers(query: String) {
        if query.isEmpty {
            loadUsers()
            return
        }

        Task { @MainActor in
            do {
                users = try await userService.searchUsers(query: query)
            } catch {
                self.error = error
            }
        }
    }
}
```

### View Model with Closures (Legacy)

```swift
import Foundation

final class UserDetailViewModel {
    // MARK: - Output Closures

    var onUserLoaded: ((User) -> Void)?
    var onLoadingStateChanged: ((Bool) -> Void)?
    var onError: ((Error) -> Void)?

    // MARK: - Properties

    private(set) var user: User?
    private let userService: UserServiceProtocol
    private let userId: String

    // MARK: - Initialization

    init(userId: String, userService: UserServiceProtocol = UserService()) {
        self.userId = userId
        self.userService = userService
    }

    // MARK: - Public Methods

    func loadUser() {
        onLoadingStateChanged?(true)

        Task { @MainActor in
            do {
                let user = try await userService.fetchUser(id: userId)
                self.user = user
                onUserLoaded?(user)
            } catch {
                onError?(error)
            }
            onLoadingStateChanged?(false)
        }
    }
}
```

---

## Coordinator Pattern

### Protocol Definition

```swift
import UIKit

protocol Coordinator: AnyObject {
    var childCoordinators: [Coordinator] { get set }
    var navigationController: UINavigationController { get set }

    func start()
}

extension Coordinator {
    func addChild(_ coordinator: Coordinator) {
        childCoordinators.append(coordinator)
    }

    func removeChild(_ coordinator: Coordinator) {
        childCoordinators.removeAll { $0 === coordinator }
    }
}
```

### App Coordinator

```swift
import UIKit

final class AppCoordinator: Coordinator {
    var childCoordinators: [Coordinator] = []
    var navigationController: UINavigationController

    private let authManager: AuthManager

    init(navigationController: UINavigationController, authManager: AuthManager = .shared) {
        self.navigationController = navigationController
        self.authManager = authManager
    }

    func start() {
        if authManager.isAuthenticated {
            showMainFlow()
        } else {
            showAuthFlow()
        }
    }

    private func showAuthFlow() {
        let coordinator = AuthCoordinator(navigationController: navigationController)
        coordinator.delegate = self
        addChild(coordinator)
        coordinator.start()
    }

    private func showMainFlow() {
        let coordinator = MainTabCoordinator(navigationController: navigationController)
        addChild(coordinator)
        coordinator.start()
    }
}

extension AppCoordinator: AuthCoordinatorDelegate {
    func authCoordinatorDidFinish(_ coordinator: AuthCoordinator) {
        removeChild(coordinator)
        showMainFlow()
    }
}
```

### Feature Coordinator

```swift
import UIKit

protocol UserCoordinatorDelegate: AnyObject {
    func userCoordinatorDidFinish(_ coordinator: UserCoordinator)
}

final class UserCoordinator: Coordinator {
    var childCoordinators: [Coordinator] = []
    var navigationController: UINavigationController

    weak var delegate: UserCoordinatorDelegate?

    init(navigationController: UINavigationController) {
        self.navigationController = navigationController
    }

    func start() {
        let viewModel = UserListViewModel()
        let viewController = UserListViewController(viewModel: viewModel)
        viewController.coordinator = self
        navigationController.pushViewController(viewController, animated: true)
    }

    func showUserDetail(_ user: User) {
        let viewModel = UserDetailViewModel(userId: user.id)
        let viewController = UserDetailViewController(viewModel: viewModel)
        viewController.coordinator = self
        navigationController.pushViewController(viewController, animated: true)
    }

    func showAddUser() {
        let viewModel = AddUserViewModel()
        let viewController = AddUserViewController(viewModel: viewModel)
        viewController.delegate = self

        let navController = UINavigationController(rootViewController: viewController)
        navigationController.present(navController, animated: true)
    }

    func showEditUser(_ user: User) {
        let viewModel = EditUserViewModel(user: user)
        let viewController = EditUserViewController(viewModel: viewModel)
        viewController.delegate = self
        navigationController.pushViewController(viewController, animated: true)
    }
}

extension UserCoordinator: AddUserViewControllerDelegate {
    func addUserViewControllerDidAddUser(_ controller: AddUserViewController, user: User) {
        controller.dismiss(animated: true)
        // Notify or refresh
    }

    func addUserViewControllerDidCancel(_ controller: AddUserViewController) {
        controller.dismiss(animated: true)
    }
}
```

---

## Custom Views

### Programmatic View

```swift
import UIKit

final class UserCell: UITableViewCell {
    static let identifier = "UserCell"

    // MARK: - UI Elements

    private let avatarImageView: UIImageView = {
        let imageView = UIImageView()
        imageView.translatesAutoresizingMaskIntoConstraints = false
        imageView.contentMode = .scaleAspectFill
        imageView.clipsToBounds = true
        imageView.backgroundColor = .systemGray5
        imageView.layer.cornerRadius = 24
        return imageView
    }()

    private let nameLabel: UILabel = {
        let label = UILabel()
        label.translatesAutoresizingMaskIntoConstraints = false
        label.font = .preferredFont(forTextStyle: .headline)
        label.textColor = .label
        return label
    }()

    private let emailLabel: UILabel = {
        let label = UILabel()
        label.translatesAutoresizingMaskIntoConstraints = false
        label.font = .preferredFont(forTextStyle: .subheadline)
        label.textColor = .secondaryLabel
        return label
    }()

    private let stackView: UIStackView = {
        let stack = UIStackView()
        stack.translatesAutoresizingMaskIntoConstraints = false
        stack.axis = .vertical
        stack.spacing = 4
        return stack
    }()

    // MARK: - Initialization

    override init(style: UITableViewCell.CellStyle, reuseIdentifier: String?) {
        super.init(style: style, reuseIdentifier: reuseIdentifier)
        setupUI()
        setupConstraints()
    }

    required init?(coder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Setup

    private func setupUI() {
        contentView.addSubview(avatarImageView)
        contentView.addSubview(stackView)

        stackView.addArrangedSubview(nameLabel)
        stackView.addArrangedSubview(emailLabel)

        accessoryType = .disclosureIndicator
    }

    private func setupConstraints() {
        NSLayoutConstraint.activate([
            avatarImageView.leadingAnchor.constraint(equalTo: contentView.leadingAnchor, constant: 16),
            avatarImageView.centerYAnchor.constraint(equalTo: contentView.centerYAnchor),
            avatarImageView.widthAnchor.constraint(equalToConstant: 48),
            avatarImageView.heightAnchor.constraint(equalToConstant: 48),
            avatarImageView.topAnchor.constraint(greaterThanOrEqualTo: contentView.topAnchor, constant: 12),
            avatarImageView.bottomAnchor.constraint(lessThanOrEqualTo: contentView.bottomAnchor, constant: -12),

            stackView.leadingAnchor.constraint(equalTo: avatarImageView.trailingAnchor, constant: 12),
            stackView.trailingAnchor.constraint(equalTo: contentView.trailingAnchor, constant: -16),
            stackView.centerYAnchor.constraint(equalTo: contentView.centerYAnchor)
        ])
    }

    // MARK: - Configuration

    func configure(with user: User) {
        nameLabel.text = user.name
        emailLabel.text = user.email

        // Load avatar image
        if let url = URL(string: user.avatarURL) {
            ImageLoader.shared.load(url: url) { [weak self] image in
                self?.avatarImageView.image = image
            }
        }
    }

    // MARK: - Reuse

    override func prepareForReuse() {
        super.prepareForReuse()
        avatarImageView.image = nil
        nameLabel.text = nil
        emailLabel.text = nil
    }
}
```

### Reusable Component

```swift
import UIKit

final class PrimaryButton: UIButton {
    // MARK: - Properties

    private var originalBackgroundColor: UIColor?

    override var isEnabled: Bool {
        didSet {
            alpha = isEnabled ? 1.0 : 0.5
        }
    }

    override var isHighlighted: Bool {
        didSet {
            UIView.animate(withDuration: 0.1) {
                self.transform = self.isHighlighted ? CGAffineTransform(scaleX: 0.95, y: 0.95) : .identity
            }
        }
    }

    // MARK: - Initialization

    override init(frame: CGRect) {
        super.init(frame: frame)
        setupUI()
    }

    required init?(coder: NSCoder) {
        super.init(coder: coder)
        setupUI()
    }

    convenience init(title: String) {
        self.init(frame: .zero)
        setTitle(title, for: .normal)
    }

    // MARK: - Setup

    private func setupUI() {
        backgroundColor = .systemBlue
        setTitleColor(.white, for: .normal)
        titleLabel?.font = .preferredFont(forTextStyle: .headline)

        layer.cornerRadius = 12

        contentEdgeInsets = UIEdgeInsets(top: 14, left: 24, bottom: 14, right: 24)

        originalBackgroundColor = backgroundColor
    }

    // MARK: - Loading State

    private var activityIndicator: UIActivityIndicatorView?
    private var savedTitle: String?

    func showLoading() {
        savedTitle = title(for: .normal)
        setTitle(nil, for: .normal)
        isEnabled = false

        let indicator = UIActivityIndicatorView(style: .medium)
        indicator.color = .white
        indicator.translatesAutoresizingMaskIntoConstraints = false
        addSubview(indicator)

        NSLayoutConstraint.activate([
            indicator.centerXAnchor.constraint(equalTo: centerXAnchor),
            indicator.centerYAnchor.constraint(equalTo: centerYAnchor)
        ])

        indicator.startAnimating()
        activityIndicator = indicator
    }

    func hideLoading() {
        activityIndicator?.stopAnimating()
        activityIndicator?.removeFromSuperview()
        activityIndicator = nil

        setTitle(savedTitle, for: .normal)
        isEnabled = true
    }
}
```

---

## Collection Views

### Modern Collection View with Diffable Data Source

```swift
import UIKit

final class PhotoGridViewController: UIViewController {
    // MARK: - Types

    enum Section {
        case main
    }

    typealias DataSource = UICollectionViewDiffableDataSource<Section, Photo>
    typealias Snapshot = NSDiffableDataSourceSnapshot<Section, Photo>

    // MARK: - Properties

    private var collectionView: UICollectionView!
    private var dataSource: DataSource!
    private let viewModel: PhotoGridViewModel

    // MARK: - Initialization

    init(viewModel: PhotoGridViewModel = PhotoGridViewModel()) {
        self.viewModel = viewModel
        super.init(nibName: nil, bundle: nil)
    }

    required init?(coder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Lifecycle

    override func viewDidLoad() {
        super.viewDidLoad()
        setupCollectionView()
        setupDataSource()
        setupBindings()
        viewModel.loadPhotos()
    }

    // MARK: - Setup

    private func setupCollectionView() {
        let layout = createLayout()
        collectionView = UICollectionView(frame: view.bounds, collectionViewLayout: layout)
        collectionView.autoresizingMask = [.flexibleWidth, .flexibleHeight]
        collectionView.delegate = self

        view.addSubview(collectionView)
    }

    private func createLayout() -> UICollectionViewLayout {
        let itemSize = NSCollectionLayoutSize(
            widthDimension: .fractionalWidth(1/3),
            heightDimension: .fractionalWidth(1/3)
        )
        let item = NSCollectionLayoutItem(layoutSize: itemSize)
        item.contentInsets = NSDirectionalEdgeInsets(top: 2, leading: 2, bottom: 2, trailing: 2)

        let groupSize = NSCollectionLayoutSize(
            widthDimension: .fractionalWidth(1.0),
            heightDimension: .fractionalWidth(1/3)
        )
        let group = NSCollectionLayoutGroup.horizontal(layoutSize: groupSize, subitems: [item])

        let section = NSCollectionLayoutSection(group: group)

        return UICollectionViewCompositionalLayout(section: section)
    }

    private func setupDataSource() {
        let cellRegistration = UICollectionView.CellRegistration<PhotoCell, Photo> { cell, indexPath, photo in
            cell.configure(with: photo)
        }

        dataSource = DataSource(collectionView: collectionView) { collectionView, indexPath, photo in
            return collectionView.dequeueConfiguredReusableCell(
                using: cellRegistration,
                for: indexPath,
                item: photo
            )
        }
    }

    private func setupBindings() {
        viewModel.onPhotosLoaded = { [weak self] photos in
            self?.applySnapshot(photos: photos)
        }
    }

    private func applySnapshot(photos: [Photo], animating: Bool = true) {
        var snapshot = Snapshot()
        snapshot.appendSections([.main])
        snapshot.appendItems(photos)
        dataSource.apply(snapshot, animatingDifferences: animating)
    }
}

// MARK: - UICollectionViewDelegate

extension PhotoGridViewController: UICollectionViewDelegate {
    func collectionView(_ collectionView: UICollectionView, didSelectItemAt indexPath: IndexPath) {
        guard let photo = dataSource.itemIdentifier(for: indexPath) else { return }
        // Handle selection
    }
}
```

### Compositional Layout with Multiple Sections

```swift
import UIKit

final class HomeViewController: UIViewController {
    // MARK: - Types

    enum Section: Int, CaseIterable {
        case featured
        case categories
        case recent

        var title: String {
            switch self {
            case .featured: return "Featured"
            case .categories: return "Categories"
            case .recent: return "Recent"
            }
        }
    }

    // MARK: - Layout

    private func createLayout() -> UICollectionViewLayout {
        UICollectionViewCompositionalLayout { [weak self] sectionIndex, environment in
            guard let section = Section(rawValue: sectionIndex) else { return nil }

            switch section {
            case .featured:
                return self?.createFeaturedSection()
            case .categories:
                return self?.createCategoriesSection()
            case .recent:
                return self?.createRecentSection()
            }
        }
    }

    private func createFeaturedSection() -> NSCollectionLayoutSection {
        let itemSize = NSCollectionLayoutSize(
            widthDimension: .fractionalWidth(1.0),
            heightDimension: .fractionalHeight(1.0)
        )
        let item = NSCollectionLayoutItem(layoutSize: itemSize)

        let groupSize = NSCollectionLayoutSize(
            widthDimension: .fractionalWidth(0.9),
            heightDimension: .absolute(250)
        )
        let group = NSCollectionLayoutGroup.horizontal(layoutSize: groupSize, subitems: [item])

        let section = NSCollectionLayoutSection(group: group)
        section.orthogonalScrollingBehavior = .groupPagingCentered
        section.contentInsets = NSDirectionalEdgeInsets(top: 16, leading: 0, bottom: 16, trailing: 0)
        section.interGroupSpacing = 16

        section.boundarySupplementaryItems = [createSectionHeader()]

        return section
    }

    private func createCategoriesSection() -> NSCollectionLayoutSection {
        let itemSize = NSCollectionLayoutSize(
            widthDimension: .fractionalWidth(1.0),
            heightDimension: .fractionalHeight(1.0)
        )
        let item = NSCollectionLayoutItem(layoutSize: itemSize)

        let groupSize = NSCollectionLayoutSize(
            widthDimension: .absolute(100),
            heightDimension: .absolute(100)
        )
        let group = NSCollectionLayoutGroup.horizontal(layoutSize: groupSize, subitems: [item])

        let section = NSCollectionLayoutSection(group: group)
        section.orthogonalScrollingBehavior = .continuous
        section.contentInsets = NSDirectionalEdgeInsets(top: 8, leading: 16, bottom: 16, trailing: 16)
        section.interGroupSpacing = 12

        section.boundarySupplementaryItems = [createSectionHeader()]

        return section
    }

    private func createRecentSection() -> NSCollectionLayoutSection {
        let itemSize = NSCollectionLayoutSize(
            widthDimension: .fractionalWidth(1.0),
            heightDimension: .estimated(100)
        )
        let item = NSCollectionLayoutItem(layoutSize: itemSize)

        let groupSize = NSCollectionLayoutSize(
            widthDimension: .fractionalWidth(1.0),
            heightDimension: .estimated(100)
        )
        let group = NSCollectionLayoutGroup.vertical(layoutSize: groupSize, subitems: [item])

        let section = NSCollectionLayoutSection(group: group)
        section.contentInsets = NSDirectionalEdgeInsets(top: 8, leading: 16, bottom: 16, trailing: 16)
        section.interGroupSpacing = 12

        section.boundarySupplementaryItems = [createSectionHeader()]

        return section
    }

    private func createSectionHeader() -> NSCollectionLayoutBoundarySupplementaryItem {
        let headerSize = NSCollectionLayoutSize(
            widthDimension: .fractionalWidth(1.0),
            heightDimension: .estimated(44)
        )
        return NSCollectionLayoutBoundarySupplementaryItem(
            layoutSize: headerSize,
            elementKind: UICollectionView.elementKindSectionHeader,
            alignment: .top
        )
    }
}
```

---

## Networking

### URL Session Wrapper

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

    init(baseURL: URL, session: URLSession = .shared) {
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

        // Add auth token if available
        if let token = AuthManager.shared.accessToken {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        return request
    }

    private func execute<T: Decodable>(_ request: URLRequest) async throws -> T {
        let (data, response) = try await session.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }

        guard (200...299).contains(httpResponse.statusCode) else {
            throw APIError.httpError(statusCode: httpResponse.statusCode)
        }

        do {
            return try decoder.decode(T.self, from: data)
        } catch {
            throw APIError.decodingError(error)
        }
    }
}
```

---

## Testing

### View Controller Testing

```swift
import XCTest
@testable import MyApp

final class UserListViewControllerTests: XCTestCase {
    var sut: UserListViewController!
    var mockViewModel: MockUserListViewModel!
    var mockCoordinator: MockUserCoordinator!

    override func setUp() {
        super.setUp()
        mockViewModel = MockUserListViewModel()
        mockCoordinator = MockUserCoordinator()
        sut = UserListViewController(viewModel: mockViewModel)
        sut.coordinator = mockCoordinator
        sut.loadViewIfNeeded()
    }

    override func tearDown() {
        sut = nil
        mockViewModel = nil
        mockCoordinator = nil
        super.tearDown()
    }

    func test_viewDidLoad_loadsUsers() {
        // When
        sut.viewDidLoad()

        // Then
        XCTAssertTrue(mockViewModel.loadUsersCalled)
    }

    func test_tableView_numberOfRows_matchesUsersCount() {
        // Given
        mockViewModel.users = [User.stub(), User.stub()]

        // When
        let rowCount = sut.tableView(sut.tableView, numberOfRowsInSection: 0)

        // Then
        XCTAssertEqual(rowCount, 2)
    }

    func test_didSelectRow_navigatesToUserDetail() {
        // Given
        let user = User.stub()
        mockViewModel.users = [user]
        let indexPath = IndexPath(row: 0, section: 0)

        // When
        sut.tableView(sut.tableView, didSelectRowAt: indexPath)

        // Then
        XCTAssertEqual(mockCoordinator.shownUser?.id, user.id)
    }
}

// Mocks
final class MockUserListViewModel: UserListViewModel {
    var loadUsersCalled = false

    override func loadUsers() {
        loadUsersCalled = true
    }
}

final class MockUserCoordinator: UserCoordinator {
    var shownUser: User?

    override func showUserDetail(_ user: User) {
        shownUser = user
    }
}
```

---

## Best Practices

### Memory Management
- ✓ Use `[weak self]` in closures
- ✓ Break retain cycles with `weak` delegates
- ✓ Cancel network tasks in `deinit`
- ✓ Use Instruments to detect leaks

### Architecture
- ✓ Use MVVM + Coordinator pattern
- ✓ Keep view controllers thin
- ✓ Extract reusable views
- ✓ Use dependency injection

### Performance
- ✓ Reuse cells properly
- ✓ Load images asynchronously
- ✓ Use diffable data sources
- ✓ Profile with Instruments

---

## References

- [UIKit Documentation](https://developer.apple.com/documentation/uikit)
- [Modern Collection Views](https://developer.apple.com/videos/play/wwdc2020/10026/)
- [Diffable Data Sources](https://developer.apple.com/videos/play/wwdc2019/220/)
- [Human Interface Guidelines](https://developer.apple.com/design/human-interface-guidelines/)
