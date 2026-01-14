# Dart Guide

> **Applies to**: Dart 3.0+, Flutter, Backend Services, CLI Tools

---

## Core Principles

1. **Type Safety**: Sound null safety, strong typing
2. **Async First**: Futures, Streams, async/await
3. **Flutter Patterns**: Widget composition, state management
4. **Immutability**: Prefer immutable data structures
5. **Code Generation**: Use build_runner for JSON, routing

---

## Language-Specific Guardrails

### Dart Version & Setup
- ✓ Use Dart 3.0+ (sound null safety required)
- ✓ Use `pubspec.yaml` for dependencies
- ✓ Pin SDK version constraints
- ✓ Use `flutter pub get` or `dart pub get`

### Code Style (dart format)
- ✓ Run `dart format` before every commit
- ✓ Use `dart analyze` for static analysis
- ✓ Follow Effective Dart guidelines
- ✓ Use `lowerCamelCase` for variables, functions
- ✓ Use `UpperCamelCase` for classes, enums, types
- ✓ Use `lowercase_with_underscores` for files, libraries
- ✓ 2-space indentation
- ✓ Max line length: 80 characters

### Null Safety
- ✓ Use non-nullable types by default
- ✓ Use `?` for nullable types: `String?`
- ✓ Use `!` only when certain (prefer null checks)
- ✓ Use `??` for default values
- ✓ Use `?.` for null-aware access
- ✓ Use `late` sparingly (prefer initialization)

### Classes & Objects
- ✓ Prefer `final` fields (immutability)
- ✓ Use named constructors for clarity
- ✓ Use factory constructors appropriately
- ✓ Implement `==` and `hashCode` for value objects
- ✓ Use `@immutable` annotation for immutable classes

### Error Handling
- ✓ Use typed exceptions
- ✓ Catch specific exceptions
- ✓ Don't catch `Error` (programming errors)
- ✓ Use `rethrow` to preserve stack traces
- ✓ Document thrown exceptions

---

## Project Structure

### Flutter App
```
myapp/
├── lib/
│   ├── main.dart
│   ├── app.dart
│   ├── features/
│   │   ├── auth/
│   │   │   ├── data/
│   │   │   │   ├── models/
│   │   │   │   └── repositories/
│   │   │   ├── domain/
│   │   │   └── presentation/
│   │   │       ├── screens/
│   │   │       ├── widgets/
│   │   │       └── providers/
│   │   └── home/
│   ├── core/
│   │   ├── constants/
│   │   ├── extensions/
│   │   ├── utils/
│   │   └── theme/
│   └── shared/
│       └── widgets/
├── test/
│   ├── unit/
│   ├── widget/
│   └── integration/
├── pubspec.yaml
├── analysis_options.yaml
└── README.md
```

### Dart Package
```
mypackage/
├── lib/
│   ├── mypackage.dart       # Main export file
│   └── src/                 # Implementation
│       ├── models/
│       └── services/
├── test/
├── example/
├── pubspec.yaml
├── analysis_options.yaml
├── CHANGELOG.md
└── README.md
```

---

## Data Classes & Models

### With Freezed (Recommended)
```dart
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:json_annotation/json_annotation.dart';

part 'user.freezed.dart';
part 'user.g.dart';

@freezed
class User with _$User {
  const factory User({
    required String id,
    required String email,
    required int age,
    @Default('user') String role,
    DateTime? createdAt,
  }) = _User;

  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);
}

// Usage
final user = User(id: '1', email: 'test@example.com', age: 25);
final updatedUser = user.copyWith(age: 26);
final json = user.toJson();
```

### Manual Implementation
```dart
import 'package:meta/meta.dart';

@immutable
class User {
  final String id;
  final String email;
  final int age;
  final String role;
  final DateTime? createdAt;

  const User({
    required this.id,
    required this.email,
    required this.age,
    this.role = 'user',
    this.createdAt,
  });

  User copyWith({
    String? id,
    String? email,
    int? age,
    String? role,
    DateTime? createdAt,
  }) {
    return User(
      id: id ?? this.id,
      email: email ?? this.email,
      age: age ?? this.age,
      role: role ?? this.role,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] as String,
      email: json['email'] as String,
      age: json['age'] as int,
      role: json['role'] as String? ?? 'user',
      createdAt: json['created_at'] != null
          ? DateTime.parse(json['created_at'] as String)
          : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email': email,
      'age': age,
      'role': role,
      if (createdAt != null) 'created_at': createdAt!.toIso8601String(),
    };
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is User &&
        other.id == id &&
        other.email == email &&
        other.age == age &&
        other.role == role &&
        other.createdAt == createdAt;
  }

  @override
  int get hashCode => Object.hash(id, email, age, role, createdAt);

  @override
  String toString() => 'User(id: $id, email: $email, age: $age, role: $role)';
}
```

---

## Async Programming

### Futures
```dart
// Async function
Future<User> fetchUser(String id) async {
  final response = await http.get(Uri.parse('$baseUrl/users/$id'));

  if (response.statusCode == 200) {
    return User.fromJson(jsonDecode(response.body));
  } else {
    throw UserNotFoundException(id);
  }
}

// Error handling
Future<void> loadData() async {
  try {
    final user = await fetchUser('123');
    print(user);
  } on UserNotFoundException catch (e) {
    print('User not found: ${e.id}');
  } on SocketException {
    print('Network error');
  } catch (e, stackTrace) {
    print('Unexpected error: $e');
    print(stackTrace);
  }
}

// Parallel execution
Future<void> loadAll() async {
  final results = await Future.wait([
    fetchUsers(),
    fetchProducts(),
    fetchOrders(),
  ]);

  final users = results[0] as List<User>;
  final products = results[1] as List<Product>;
  final orders = results[2] as List<Order>;
}
```

### Streams
```dart
// Create stream
Stream<int> countStream(int max) async* {
  for (int i = 1; i <= max; i++) {
    await Future.delayed(const Duration(seconds: 1));
    yield i;
  }
}

// Stream transformations
final doubled = countStream(5).map((n) => n * 2);
final evens = countStream(10).where((n) => n.isEven);

// Listen to stream
final subscription = countStream(5).listen(
  (value) => print('Value: $value'),
  onError: (error) => print('Error: $error'),
  onDone: () => print('Done'),
);

// Cancel subscription
await subscription.cancel();

// StreamController
class UserService {
  final _userController = StreamController<User>.broadcast();

  Stream<User> get userStream => _userController.stream;

  void updateUser(User user) {
    _userController.add(user);
  }

  void dispose() {
    _userController.close();
  }
}
```

---

## Flutter Widgets

### Stateless Widget
```dart
class UserCard extends StatelessWidget {
  final User user;
  final VoidCallback? onTap;

  const UserCard({
    super.key,
    required this.user,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: ListTile(
        leading: CircleAvatar(
          child: Text(user.email[0].toUpperCase()),
        ),
        title: Text(user.email),
        subtitle: Text('Age: ${user.age}'),
        trailing: const Icon(Icons.chevron_right),
        onTap: onTap,
      ),
    );
  }
}
```

### Stateful Widget
```dart
class CounterWidget extends StatefulWidget {
  final int initialValue;

  const CounterWidget({
    super.key,
    this.initialValue = 0,
  });

  @override
  State<CounterWidget> createState() => _CounterWidgetState();
}

class _CounterWidgetState extends State<CounterWidget> {
  late int _count;

  @override
  void initState() {
    super.initState();
    _count = widget.initialValue;
  }

  void _increment() {
    setState(() {
      _count++;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Text(
          'Count: $_count',
          style: Theme.of(context).textTheme.headlineMedium,
        ),
        const SizedBox(height: 16),
        ElevatedButton(
          onPressed: _increment,
          child: const Text('Increment'),
        ),
      ],
    );
  }
}
```

### State Management with Riverpod
```dart
import 'package:flutter_riverpod/flutter_riverpod.dart';

// Simple provider
final counterProvider = StateProvider<int>((ref) => 0);

// Async provider
final userProvider = FutureProvider.family<User, String>((ref, id) async {
  final repository = ref.watch(userRepositoryProvider);
  return repository.getUser(id);
});

// Notifier provider
class UserListNotifier extends AsyncNotifier<List<User>> {
  @override
  Future<List<User>> build() async {
    return ref.watch(userRepositoryProvider).getUsers();
  }

  Future<void> addUser(User user) async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(() async {
      await ref.read(userRepositoryProvider).createUser(user);
      return ref.read(userRepositoryProvider).getUsers();
    });
  }
}

final userListProvider =
    AsyncNotifierProvider<UserListNotifier, List<User>>(UserListNotifier.new);

// Usage in widget
class UserListScreen extends ConsumerWidget {
  const UserListScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final usersAsync = ref.watch(userListProvider);

    return usersAsync.when(
      data: (users) => ListView.builder(
        itemCount: users.length,
        itemBuilder: (context, index) => UserCard(user: users[index]),
      ),
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, stack) => Center(child: Text('Error: $error')),
    );
  }
}
```

---

## Testing

### Unit Tests
```dart
import 'package:test/test.dart';
import 'package:mocktail/mocktail.dart';

class MockUserRepository extends Mock implements UserRepository {}

void main() {
  group('UserService', () {
    late UserService service;
    late MockUserRepository mockRepository;

    setUp(() {
      mockRepository = MockUserRepository();
      service = UserService(repository: mockRepository);
    });

    test('getUser returns user when found', () async {
      // Arrange
      final user = User(id: '1', email: 'test@example.com', age: 25);
      when(() => mockRepository.getUser('1')).thenAnswer((_) async => user);

      // Act
      final result = await service.getUser('1');

      // Assert
      expect(result, equals(user));
      verify(() => mockRepository.getUser('1')).called(1);
    });

    test('getUser throws when not found', () async {
      when(() => mockRepository.getUser('999'))
          .thenThrow(UserNotFoundException('999'));

      expect(
        () => service.getUser('999'),
        throwsA(isA<UserNotFoundException>()),
      );
    });
  });
}
```

### Widget Tests
```dart
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('UserCard', () {
    testWidgets('displays user information', (tester) async {
      final user = User(id: '1', email: 'test@example.com', age: 25);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: UserCard(user: user),
          ),
        ),
      );

      expect(find.text('test@example.com'), findsOneWidget);
      expect(find.text('Age: 25'), findsOneWidget);
    });

    testWidgets('calls onTap when tapped', (tester) async {
      var tapped = false;
      final user = User(id: '1', email: 'test@example.com', age: 25);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: UserCard(
              user: user,
              onTap: () => tapped = true,
            ),
          ),
        ),
      );

      await tester.tap(find.byType(UserCard));
      expect(tapped, isTrue);
    });
  });
}
```

---

## Error Handling

### Custom Exceptions
```dart
sealed class AppException implements Exception {
  const AppException(this.message);
  final String message;

  @override
  String toString() => message;
}

class NetworkException extends AppException {
  const NetworkException([super.message = 'Network error occurred']);
}

class ValidationException extends AppException {
  const ValidationException(super.message, {this.field});
  final String? field;
}

class UserNotFoundException extends AppException {
  const UserNotFoundException(this.userId)
      : super('User not found');
  final String userId;
}

// Usage
Future<User> getUser(String id) async {
  try {
    final response = await api.get('/users/$id');
    return User.fromJson(response.data);
  } on DioException catch (e) {
    if (e.response?.statusCode == 404) {
      throw UserNotFoundException(id);
    }
    throw NetworkException(e.message ?? 'Network error');
  }
}
```

### Result Type Pattern
```dart
sealed class Result<T> {
  const Result();
}

class Success<T> extends Result<T> {
  const Success(this.value);
  final T value;
}

class Failure<T> extends Result<T> {
  const Failure(this.error, [this.stackTrace]);
  final Object error;
  final StackTrace? stackTrace;
}

// Usage
Future<Result<User>> getUser(String id) async {
  try {
    final user = await repository.getUser(id);
    return Success(user);
  } catch (e, st) {
    return Failure(e, st);
  }
}

// Handle result
final result = await getUser('123');
switch (result) {
  case Success(:final value):
    print('User: ${value.email}');
  case Failure(:final error):
    print('Error: $error');
}
```

---

## Configuration

### analysis_options.yaml
```yaml
include: package:flutter_lints/flutter.yaml
# Or for Dart: package:lints/recommended.yaml

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
```

---

## Common Pitfalls

### Avoid These
```dart
// Using dynamic
dynamic user = fetchUser();

// Not awaiting futures
fetchData();  // Fire and forget bug

// Using ! without checking
String name = nullableString!;

// Mutable state in providers
final usersProvider = Provider((ref) => []);

// Building widgets in initState
@override
void initState() {
  super.initState();
  // Don't use context here!
}
```

### Do This Instead
```dart
// Use proper types
final User user = await fetchUser();

// Await or handle futures
await fetchData();
// Or intentionally ignore
unawaited(logAnalytics());

// Check null first
if (nullableString != null) {
  String name = nullableString;
}

// Use proper state management
final usersProvider = StateProvider<List<User>>((ref) => []);

// Use didChangeDependencies for context
@override
void didChangeDependencies() {
  super.didChangeDependencies();
  // Safe to use context here
}
```

---

## References

- [Dart Documentation](https://dart.dev/guides)
- [Effective Dart](https://dart.dev/guides/language/effective-dart)
- [Flutter Documentation](https://docs.flutter.dev/)
- [Riverpod Documentation](https://riverpod.dev/)
- [Freezed Package](https://pub.dev/packages/freezed)
- [Flutter Testing](https://docs.flutter.dev/testing)
