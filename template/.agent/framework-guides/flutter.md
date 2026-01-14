# Flutter Framework Guide

> **Framework**: Flutter 3.x
> **Language**: Dart 3.x
> **Type**: Cross-Platform UI Framework
> **Use Cases**: Mobile apps (iOS/Android), Web, Desktop

---

## Overview

Flutter is Google's UI toolkit for building natively compiled applications for mobile, web, and desktop from a single codebase. It uses Dart and provides its own rendering engine for consistent cross-platform UI.

### Key Features
- Hot reload for fast development
- Widget-based declarative UI
- Single codebase for multiple platforms
- Rich set of Material and Cupertino widgets
- Strong typing with Dart null safety

---

## Project Structure

```
myapp/
├── lib/
│   ├── main.dart                 # App entry point
│   ├── app.dart                  # App widget
│   ├── features/
│   │   ├── auth/
│   │   │   ├── data/
│   │   │   │   ├── models/
│   │   │   │   │   └── user_model.dart
│   │   │   │   ├── repositories/
│   │   │   │   │   └── auth_repository.dart
│   │   │   │   └── datasources/
│   │   │   │       └── auth_remote_datasource.dart
│   │   │   ├── domain/
│   │   │   │   ├── entities/
│   │   │   │   │   └── user.dart
│   │   │   │   ├── repositories/
│   │   │   │   │   └── auth_repository.dart
│   │   │   │   └── usecases/
│   │   │   │       └── login_usecase.dart
│   │   │   └── presentation/
│   │   │       ├── screens/
│   │   │       │   └── login_screen.dart
│   │   │       ├── widgets/
│   │   │       │   └── login_form.dart
│   │   │       └── providers/
│   │   │           └── auth_provider.dart
│   │   └── home/
│   │       └── ...
│   ├── core/
│   │   ├── constants/
│   │   │   └── app_constants.dart
│   │   ├── errors/
│   │   │   └── failures.dart
│   │   ├── network/
│   │   │   └── api_client.dart
│   │   ├── theme/
│   │   │   └── app_theme.dart
│   │   ├── utils/
│   │   │   └── validators.dart
│   │   └── widgets/
│   │       └── common_widgets.dart
│   └── l10n/                     # Localization
│       └── app_en.arb
├── test/
│   ├── unit/
│   ├── widget/
│   └── integration/
├── pubspec.yaml
├── analysis_options.yaml
└── README.md
```

---

## Application Setup

### Main Entry Point

```dart
// lib/main.dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'app.dart';
import 'core/di/injection.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize dependencies
  await configureDependencies();

  runApp(
    const ProviderScope(
      child: MyApp(),
    ),
  );
}
```

### App Widget

```dart
// lib/app.dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'core/router/app_router.dart';
import 'core/theme/app_theme.dart';

class MyApp extends ConsumerWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(routerProvider);

    return MaterialApp.router(
      title: 'My App',
      theme: AppTheme.light,
      darkTheme: AppTheme.dark,
      themeMode: ThemeMode.system,
      routerConfig: router,
      debugShowCheckedModeBanner: false,
    );
  }
}
```

### Theme Configuration

```dart
// lib/core/theme/app_theme.dart
import 'package:flutter/material.dart';

class AppTheme {
  AppTheme._();

  static const _primaryColor = Color(0xFF6750A4);
  static const _secondaryColor = Color(0xFF625B71);

  static ThemeData get light {
    return ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: _primaryColor,
        brightness: Brightness.light,
      ),
      appBarTheme: const AppBarTheme(
        centerTitle: true,
        elevation: 0,
      ),
      inputDecorationTheme: InputDecorationTheme(
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        filled: true,
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          minimumSize: const Size(double.infinity, 48),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
        ),
      ),
    );
  }

  static ThemeData get dark {
    return ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: _primaryColor,
        brightness: Brightness.dark,
      ),
      appBarTheme: const AppBarTheme(
        centerTitle: true,
        elevation: 0,
      ),
    );
  }
}
```

---

## Routing (go_router)

### Router Configuration

```dart
// lib/core/router/app_router.dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../features/auth/presentation/screens/login_screen.dart';
import '../../features/auth/presentation/providers/auth_provider.dart';
import '../../features/home/presentation/screens/home_screen.dart';
import '../../features/profile/presentation/screens/profile_screen.dart';

final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authStateProvider);

  return GoRouter(
    initialLocation: '/',
    redirect: (context, state) {
      final isLoggedIn = authState.valueOrNull != null;
      final isLoggingIn = state.matchedLocation == '/login';

      if (!isLoggedIn && !isLoggingIn) {
        return '/login';
      }
      if (isLoggedIn && isLoggingIn) {
        return '/';
      }
      return null;
    },
    routes: [
      GoRoute(
        path: '/login',
        name: 'login',
        builder: (context, state) => const LoginScreen(),
      ),
      ShellRoute(
        builder: (context, state, child) => ScaffoldWithNavBar(child: child),
        routes: [
          GoRoute(
            path: '/',
            name: 'home',
            builder: (context, state) => const HomeScreen(),
          ),
          GoRoute(
            path: '/profile',
            name: 'profile',
            builder: (context, state) => const ProfileScreen(),
          ),
          GoRoute(
            path: '/users/:id',
            name: 'user-detail',
            builder: (context, state) {
              final id = state.pathParameters['id']!;
              return UserDetailScreen(userId: id);
            },
          ),
        ],
      ),
    ],
    errorBuilder: (context, state) => ErrorScreen(error: state.error),
  );
});

class ScaffoldWithNavBar extends StatelessWidget {
  final Widget child;

  const ScaffoldWithNavBar({super.key, required this.child});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: child,
      bottomNavigationBar: NavigationBar(
        selectedIndex: _calculateSelectedIndex(context),
        onDestinationSelected: (index) => _onItemTapped(index, context),
        destinations: const [
          NavigationDestination(
            icon: Icon(Icons.home_outlined),
            selectedIcon: Icon(Icons.home),
            label: 'Home',
          ),
          NavigationDestination(
            icon: Icon(Icons.person_outline),
            selectedIcon: Icon(Icons.person),
            label: 'Profile',
          ),
        ],
      ),
    );
  }

  int _calculateSelectedIndex(BuildContext context) {
    final location = GoRouterState.of(context).matchedLocation;
    if (location.startsWith('/profile')) return 1;
    return 0;
  }

  void _onItemTapped(int index, BuildContext context) {
    switch (index) {
      case 0:
        context.goNamed('home');
        break;
      case 1:
        context.goNamed('profile');
        break;
    }
  }
}
```

---

## State Management (Riverpod)

### Providers

```dart
// lib/features/auth/presentation/providers/auth_provider.dart
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../domain/entities/user.dart';
import '../../domain/repositories/auth_repository.dart';

// Repository provider
final authRepositoryProvider = Provider<AuthRepository>((ref) {
  return AuthRepositoryImpl(
    remoteDataSource: ref.watch(authRemoteDataSourceProvider),
    localDataSource: ref.watch(authLocalDataSourceProvider),
  );
});

// Auth state
final authStateProvider = StreamProvider<User?>((ref) {
  final repository = ref.watch(authRepositoryProvider);
  return repository.authStateChanges;
});

// Current user
final currentUserProvider = Provider<User?>((ref) {
  return ref.watch(authStateProvider).valueOrNull;
});

// Login notifier
final loginProvider = AsyncNotifierProvider<LoginNotifier, void>(
  LoginNotifier.new,
);

class LoginNotifier extends AsyncNotifier<void> {
  @override
  Future<void> build() async {}

  Future<void> login(String email, String password) async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(() async {
      final repository = ref.read(authRepositoryProvider);
      await repository.signInWithEmailAndPassword(email, password);
    });
  }

  Future<void> loginWithGoogle() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(() async {
      final repository = ref.read(authRepositoryProvider);
      await repository.signInWithGoogle();
    });
  }

  Future<void> logout() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(() async {
      final repository = ref.read(authRepositoryProvider);
      await repository.signOut();
    });
  }
}
```

### State with AsyncNotifier

```dart
// lib/features/home/presentation/providers/users_provider.dart
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../domain/entities/user.dart';
import '../../domain/repositories/user_repository.dart';

final usersProvider = AsyncNotifierProvider<UsersNotifier, List<User>>(
  UsersNotifier.new,
);

class UsersNotifier extends AsyncNotifier<List<User>> {
  @override
  Future<List<User>> build() async {
    return _fetchUsers();
  }

  Future<List<User>> _fetchUsers() async {
    final repository = ref.read(userRepositoryProvider);
    return repository.getUsers();
  }

  Future<void> refresh() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(_fetchUsers);
  }

  Future<void> addUser(User user) async {
    final repository = ref.read(userRepositoryProvider);
    await repository.createUser(user);

    // Update state optimistically or refetch
    state = AsyncData([...state.value ?? [], user]);
  }

  Future<void> deleteUser(String id) async {
    final repository = ref.read(userRepositoryProvider);
    await repository.deleteUser(id);

    state = AsyncData(
      state.value?.where((u) => u.id != id).toList() ?? [],
    );
  }
}

// Filtered users
final searchQueryProvider = StateProvider<String>((ref) => '');

final filteredUsersProvider = Provider<AsyncValue<List<User>>>((ref) {
  final users = ref.watch(usersProvider);
  final query = ref.watch(searchQueryProvider).toLowerCase();

  return users.whenData((userList) {
    if (query.isEmpty) return userList;
    return userList.where((user) {
      return user.name.toLowerCase().contains(query) ||
          user.email.toLowerCase().contains(query);
    }).toList();
  });
});
```

---

## Widgets & Screens

### Screen Widget

```dart
// lib/features/home/presentation/screens/home_screen.dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/users_provider.dart';
import '../widgets/user_list_tile.dart';

class HomeScreen extends ConsumerWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final usersAsync = ref.watch(filteredUsersProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Users'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => ref.read(usersProvider.notifier).refresh(),
          ),
        ],
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: SearchBar(
              hintText: 'Search users...',
              onChanged: (value) {
                ref.read(searchQueryProvider.notifier).state = value;
              },
            ),
          ),
          Expanded(
            child: usersAsync.when(
              data: (users) => users.isEmpty
                  ? const Center(child: Text('No users found'))
                  : ListView.builder(
                      itemCount: users.length,
                      itemBuilder: (context, index) {
                        return UserListTile(user: users[index]);
                      },
                    ),
              loading: () => const Center(child: CircularProgressIndicator()),
              error: (error, stack) => Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text('Error: $error'),
                    const SizedBox(height: 16),
                    ElevatedButton(
                      onPressed: () => ref.refresh(usersProvider),
                      child: const Text('Retry'),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showAddUserDialog(context, ref),
        child: const Icon(Icons.add),
      ),
    );
  }

  void _showAddUserDialog(BuildContext context, WidgetRef ref) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (context) => const AddUserBottomSheet(),
    );
  }
}
```

### Reusable Widgets

```dart
// lib/features/home/presentation/widgets/user_list_tile.dart
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../domain/entities/user.dart';

class UserListTile extends StatelessWidget {
  final User user;

  const UserListTile({super.key, required this.user});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return ListTile(
      leading: CircleAvatar(
        backgroundImage: user.avatarUrl != null
            ? NetworkImage(user.avatarUrl!)
            : null,
        child: user.avatarUrl == null
            ? Text(user.name[0].toUpperCase())
            : null,
      ),
      title: Text(user.name),
      subtitle: Text(user.email),
      trailing: const Icon(Icons.chevron_right),
      onTap: () => context.pushNamed('user-detail', pathParameters: {'id': user.id}),
    );
  }
}
```

### Form Widget

```dart
// lib/features/auth/presentation/widgets/login_form.dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/utils/validators.dart';
import '../providers/auth_provider.dart';

class LoginForm extends ConsumerStatefulWidget {
  const LoginForm({super.key});

  @override
  ConsumerState<LoginForm> createState() => _LoginFormState();
}

class _LoginFormState extends ConsumerState<LoginForm> {
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _obscurePassword = true;

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    await ref.read(loginProvider.notifier).login(
      _emailController.text.trim(),
      _passwordController.text,
    );
  }

  @override
  Widget build(BuildContext context) {
    final loginState = ref.watch(loginProvider);

    ref.listen(loginProvider, (previous, next) {
      next.whenOrNull(
        error: (error, _) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text(error.toString())),
          );
        },
      );
    });

    return Form(
      key: _formKey,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          TextFormField(
            controller: _emailController,
            decoration: const InputDecoration(
              labelText: 'Email',
              prefixIcon: Icon(Icons.email_outlined),
            ),
            keyboardType: TextInputType.emailAddress,
            textInputAction: TextInputAction.next,
            validator: Validators.email,
            enabled: !loginState.isLoading,
          ),
          const SizedBox(height: 16),
          TextFormField(
            controller: _passwordController,
            decoration: InputDecoration(
              labelText: 'Password',
              prefixIcon: const Icon(Icons.lock_outline),
              suffixIcon: IconButton(
                icon: Icon(
                  _obscurePassword
                      ? Icons.visibility_outlined
                      : Icons.visibility_off_outlined,
                ),
                onPressed: () {
                  setState(() => _obscurePassword = !_obscurePassword);
                },
              ),
            ),
            obscureText: _obscurePassword,
            textInputAction: TextInputAction.done,
            validator: Validators.password,
            enabled: !loginState.isLoading,
            onFieldSubmitted: (_) => _submit(),
          ),
          const SizedBox(height: 24),
          ElevatedButton(
            onPressed: loginState.isLoading ? null : _submit,
            child: loginState.isLoading
                ? const SizedBox(
                    height: 20,
                    width: 20,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Text('Sign In'),
          ),
          const SizedBox(height: 16),
          OutlinedButton.icon(
            onPressed: loginState.isLoading
                ? null
                : () => ref.read(loginProvider.notifier).loginWithGoogle(),
            icon: Image.asset('assets/google_logo.png', height: 20),
            label: const Text('Continue with Google'),
          ),
        ],
      ),
    );
  }
}
```

---

## Data Layer

### Models with Freezed

```dart
// lib/features/auth/data/models/user_model.dart
import 'package:freezed_annotation/freezed_annotation.dart';

import '../../domain/entities/user.dart';

part 'user_model.freezed.dart';
part 'user_model.g.dart';

@freezed
class UserModel with _$UserModel {
  const UserModel._();

  const factory UserModel({
    required String id,
    required String email,
    required String name,
    @JsonKey(name: 'avatar_url') String? avatarUrl,
    @JsonKey(name: 'created_at') required DateTime createdAt,
  }) = _UserModel;

  factory UserModel.fromJson(Map<String, dynamic> json) =>
      _$UserModelFromJson(json);

  User toEntity() => User(
        id: id,
        email: email,
        name: name,
        avatarUrl: avatarUrl,
        createdAt: createdAt,
      );

  factory UserModel.fromEntity(User user) => UserModel(
        id: user.id,
        email: user.email,
        name: user.name,
        avatarUrl: user.avatarUrl,
        createdAt: user.createdAt,
      );
}
```

### Repository Implementation

```dart
// lib/features/auth/data/repositories/auth_repository_impl.dart
import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../../domain/entities/user.dart';
import '../../domain/repositories/auth_repository.dart';
import '../datasources/auth_remote_datasource.dart';
import '../datasources/auth_local_datasource.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthRemoteDataSource remoteDataSource;
  final AuthLocalDataSource localDataSource;

  AuthRepositoryImpl({
    required this.remoteDataSource,
    required this.localDataSource,
  });

  @override
  Stream<User?> get authStateChanges {
    return remoteDataSource.authStateChanges.map((model) => model?.toEntity());
  }

  @override
  Future<Either<Failure, User>> signInWithEmailAndPassword(
    String email,
    String password,
  ) async {
    try {
      final userModel = await remoteDataSource.signInWithEmailAndPassword(
        email,
        password,
      );
      await localDataSource.cacheUser(userModel);
      return Right(userModel.toEntity());
    } on AuthException catch (e) {
      return Left(AuthFailure(e.message));
    } catch (e) {
      return Left(ServerFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, User>> signInWithGoogle() async {
    try {
      final userModel = await remoteDataSource.signInWithGoogle();
      await localDataSource.cacheUser(userModel);
      return Right(userModel.toEntity());
    } on AuthException catch (e) {
      return Left(AuthFailure(e.message));
    } catch (e) {
      return Left(ServerFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> signOut() async {
    try {
      await remoteDataSource.signOut();
      await localDataSource.clearCache();
      return const Right(null);
    } catch (e) {
      return Left(ServerFailure(e.toString()));
    }
  }
}
```

### API Client

```dart
// lib/core/network/api_client.dart
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../constants/app_constants.dart';
import '../errors/exceptions.dart';

final dioProvider = Provider<Dio>((ref) {
  final dio = Dio(
    BaseOptions(
      baseUrl: AppConstants.apiBaseUrl,
      connectTimeout: const Duration(seconds: 30),
      receiveTimeout: const Duration(seconds: 30),
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    ),
  );

  dio.interceptors.addAll([
    AuthInterceptor(ref),
    LogInterceptor(
      requestBody: true,
      responseBody: true,
    ),
  ]);

  return dio;
});

class AuthInterceptor extends Interceptor {
  final Ref ref;

  AuthInterceptor(this.ref);

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    final token = await ref.read(tokenProvider.future);
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    if (err.response?.statusCode == 401) {
      ref.read(authProvider.notifier).logout();
    }
    handler.next(err);
  }
}

final apiClientProvider = Provider<ApiClient>((ref) {
  return ApiClient(ref.watch(dioProvider));
});

class ApiClient {
  final Dio _dio;

  ApiClient(this._dio);

  Future<T> get<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
    required T Function(dynamic) parser,
  }) async {
    try {
      final response = await _dio.get(path, queryParameters: queryParameters);
      return parser(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<T> post<T>(
    String path, {
    dynamic data,
    required T Function(dynamic) parser,
  }) async {
    try {
      final response = await _dio.post(path, data: data);
      return parser(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  AppException _handleError(DioException error) {
    switch (error.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.receiveTimeout:
      case DioExceptionType.sendTimeout:
        return NetworkException('Connection timeout');
      case DioExceptionType.badResponse:
        final statusCode = error.response?.statusCode;
        final message = error.response?.data?['message'] ?? 'Unknown error';
        if (statusCode == 401) {
          return UnauthorizedException(message);
        } else if (statusCode == 404) {
          return NotFoundException(message);
        } else if (statusCode == 422) {
          return ValidationException(message);
        }
        return ServerException(message);
      default:
        return NetworkException('Network error');
    }
  }
}
```

---

## Testing

### Widget Tests

```dart
// test/widget/login_form_test.dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';

import 'package:myapp/features/auth/presentation/widgets/login_form.dart';
import 'package:myapp/features/auth/presentation/providers/auth_provider.dart';

class MockAuthRepository extends Mock implements AuthRepository {}

void main() {
  late MockAuthRepository mockRepository;

  setUp(() {
    mockRepository = MockAuthRepository();
  });

  Widget createWidget() {
    return ProviderScope(
      overrides: [
        authRepositoryProvider.overrideWithValue(mockRepository),
      ],
      child: const MaterialApp(
        home: Scaffold(body: LoginForm()),
      ),
    );
  }

  group('LoginForm', () {
    testWidgets('renders email and password fields', (tester) async {
      await tester.pumpWidget(createWidget());

      expect(find.byType(TextFormField), findsNWidgets(2));
      expect(find.text('Email'), findsOneWidget);
      expect(find.text('Password'), findsOneWidget);
    });

    testWidgets('shows validation errors for empty fields', (tester) async {
      await tester.pumpWidget(createWidget());

      await tester.tap(find.text('Sign In'));
      await tester.pump();

      expect(find.text('Email is required'), findsOneWidget);
      expect(find.text('Password is required'), findsOneWidget);
    });

    testWidgets('calls login with valid credentials', (tester) async {
      when(() => mockRepository.signInWithEmailAndPassword(any(), any()))
          .thenAnswer((_) async => Right(testUser));

      await tester.pumpWidget(createWidget());

      await tester.enterText(
        find.widgetWithText(TextFormField, 'Email'),
        'test@example.com',
      );
      await tester.enterText(
        find.widgetWithText(TextFormField, 'Password'),
        'password123',
      );
      await tester.tap(find.text('Sign In'));
      await tester.pump();

      verify(() => mockRepository.signInWithEmailAndPassword(
            'test@example.com',
            'password123',
          )).called(1);
    });

    testWidgets('shows loading indicator during login', (tester) async {
      when(() => mockRepository.signInWithEmailAndPassword(any(), any()))
          .thenAnswer((_) async {
        await Future.delayed(const Duration(seconds: 1));
        return Right(testUser);
      });

      await tester.pumpWidget(createWidget());

      await tester.enterText(
        find.widgetWithText(TextFormField, 'Email'),
        'test@example.com',
      );
      await tester.enterText(
        find.widgetWithText(TextFormField, 'Password'),
        'password123',
      );
      await tester.tap(find.text('Sign In'));
      await tester.pump();

      expect(find.byType(CircularProgressIndicator), findsOneWidget);
    });
  });
}
```

### Provider Tests

```dart
// test/unit/auth_provider_test.dart
import 'package:dartz/dartz.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';

import 'package:myapp/features/auth/domain/entities/user.dart';
import 'package:myapp/features/auth/domain/repositories/auth_repository.dart';
import 'package:myapp/features/auth/presentation/providers/auth_provider.dart';

class MockAuthRepository extends Mock implements AuthRepository {}

void main() {
  late MockAuthRepository mockRepository;
  late ProviderContainer container;

  final testUser = User(
    id: '1',
    email: 'test@example.com',
    name: 'Test User',
    createdAt: DateTime.now(),
  );

  setUp(() {
    mockRepository = MockAuthRepository();
    container = ProviderContainer(
      overrides: [
        authRepositoryProvider.overrideWithValue(mockRepository),
      ],
    );
  });

  tearDown(() {
    container.dispose();
  });

  group('LoginNotifier', () {
    test('login success updates state', () async {
      when(() => mockRepository.signInWithEmailAndPassword(any(), any()))
          .thenAnswer((_) async => Right(testUser));

      final notifier = container.read(loginProvider.notifier);
      await notifier.login('test@example.com', 'password');

      final state = container.read(loginProvider);
      expect(state.hasError, false);
      expect(state.isLoading, false);
    });

    test('login failure sets error state', () async {
      when(() => mockRepository.signInWithEmailAndPassword(any(), any()))
          .thenAnswer((_) async => Left(AuthFailure('Invalid credentials')));

      final notifier = container.read(loginProvider.notifier);
      await notifier.login('test@example.com', 'wrong');

      final state = container.read(loginProvider);
      expect(state.hasError, true);
    });
  });
}
```

### Integration Tests

```dart
// integration_test/app_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';

import 'package:myapp/main.dart' as app;

void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  group('end-to-end test', () {
    testWidgets('login flow', (tester) async {
      app.main();
      await tester.pumpAndSettle();

      // Should be on login screen
      expect(find.text('Sign In'), findsOneWidget);

      // Enter credentials
      await tester.enterText(
        find.byKey(const Key('email_field')),
        'test@example.com',
      );
      await tester.enterText(
        find.byKey(const Key('password_field')),
        'password123',
      );

      // Submit
      await tester.tap(find.text('Sign In'));
      await tester.pumpAndSettle();

      // Should be on home screen
      expect(find.text('Home'), findsOneWidget);
    });
  });
}
```

---

## Configuration

### pubspec.yaml

```yaml
name: myapp
description: My Flutter application
publish_to: 'none'
version: 1.0.0+1

environment:
  sdk: '>=3.0.0 <4.0.0'

dependencies:
  flutter:
    sdk: flutter

  # State management
  flutter_riverpod: ^2.4.0
  riverpod_annotation: ^2.3.0

  # Navigation
  go_router: ^12.0.0

  # Networking
  dio: ^5.3.0

  # Data classes
  freezed_annotation: ^2.4.0
  json_annotation: ^4.8.0

  # Functional programming
  dartz: ^0.10.1

  # Local storage
  shared_preferences: ^2.2.0
  flutter_secure_storage: ^9.0.0

  # Utils
  intl: ^0.18.0
  logger: ^2.0.0

dev_dependencies:
  flutter_test:
    sdk: flutter
  integration_test:
    sdk: flutter

  # Code generation
  build_runner: ^2.4.0
  freezed: ^2.4.0
  json_serializable: ^6.7.0
  riverpod_generator: ^2.3.0

  # Testing
  mocktail: ^1.0.0

  # Linting
  flutter_lints: ^3.0.0

flutter:
  uses-material-design: true
  assets:
    - assets/images/
    - assets/icons/
```

### analysis_options.yaml

```yaml
include: package:flutter_lints/flutter.yaml

analyzer:
  language:
    strict-casts: true
    strict-inference: true
    strict-raw-types: true
  errors:
    missing_required_param: error
    missing_return: error
  exclude:
    - "**/*.g.dart"
    - "**/*.freezed.dart"

linter:
  rules:
    - always_declare_return_types
    - avoid_dynamic_calls
    - avoid_print
    - avoid_type_to_string
    - cancel_subscriptions
    - close_sinks
    - prefer_const_constructors
    - prefer_const_declarations
    - prefer_final_fields
    - prefer_final_locals
    - require_trailing_commas
    - unawaited_futures
    - unnecessary_await_in_return
    - use_key_in_widget_constructors
```

---

## Commands

```bash
# Create new project
flutter create myapp

# Run app
flutter run
flutter run -d chrome  # Web
flutter run -d macos   # macOS

# Build
flutter build apk
flutter build ios
flutter build web
flutter build macos

# Code generation
dart run build_runner build --delete-conflicting-outputs
dart run build_runner watch

# Tests
flutter test
flutter test --coverage
flutter test integration_test/

# Analyze
flutter analyze
dart fix --apply

# Format
dart format .
```

---

## References

- [Flutter Documentation](https://docs.flutter.dev/)
- [Dart Documentation](https://dart.dev/guides)
- [Riverpod Documentation](https://riverpod.dev/)
- [go_router Documentation](https://pub.dev/packages/go_router)
- [Freezed Package](https://pub.dev/packages/freezed)
- [Flutter Testing](https://docs.flutter.dev/testing)
