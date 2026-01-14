# Shelf Framework Guide

> **Framework**: Shelf 1.x
> **Language**: Dart 3.x
> **Type**: Middleware-Based HTTP Server
> **Use Cases**: REST APIs, Microservices, Backend Services

---

## Overview

Shelf is a web server middleware framework for Dart. It provides a simple, composable model for building HTTP servers using middleware pipelines. Shelf is the foundation for many Dart backend frameworks.

### Key Features
- Middleware-based architecture
- Composable request/response handlers
- Streaming support
- WebSocket support (via shelf_web_socket)
- Simple and lightweight

---

## Project Structure

```
myapp/
├── bin/
│   └── server.dart               # Entry point
├── lib/
│   ├── src/
│   │   ├── app.dart              # Application setup
│   │   ├── config/
│   │   │   └── config.dart       # Configuration
│   │   ├── handlers/
│   │   │   ├── users_handler.dart
│   │   │   └── health_handler.dart
│   │   ├── middleware/
│   │   │   ├── auth_middleware.dart
│   │   │   ├── cors_middleware.dart
│   │   │   └── logging_middleware.dart
│   │   ├── models/
│   │   │   └── user.dart
│   │   ├── repositories/
│   │   │   └── user_repository.dart
│   │   └── services/
│   │       └── user_service.dart
│   └── myapp.dart                # Library export
├── test/
│   └── handlers/
│       └── users_handler_test.dart
├── pubspec.yaml
├── analysis_options.yaml
└── Dockerfile
```

---

## Application Setup

### Entry Point

```dart
// bin/server.dart
import 'dart:io';

import 'package:shelf/shelf_io.dart' as shelf_io;

import 'package:myapp/myapp.dart';

Future<void> main() async {
  final config = Config.fromEnvironment();
  final app = Application(config);

  final handler = await app.createHandler();

  final server = await shelf_io.serve(
    handler,
    InternetAddress.anyIPv4,
    config.port,
  );

  print('Server running on http://${server.address.host}:${server.port}');

  // Graceful shutdown
  ProcessSignal.sigint.watch().listen((_) async {
    print('Shutting down...');
    await app.close();
    await server.close();
    exit(0);
  });
}
```

### Application Class

```dart
// lib/src/app.dart
import 'package:shelf/shelf.dart';
import 'package:shelf_router/shelf_router.dart';

import 'config/config.dart';
import 'handlers/health_handler.dart';
import 'handlers/users_handler.dart';
import 'middleware/auth_middleware.dart';
import 'middleware/cors_middleware.dart';
import 'middleware/logging_middleware.dart';
import 'repositories/user_repository.dart';
import 'services/user_service.dart';

class Application {
  final Config config;

  late final UserRepository _userRepository;
  late final UserService _userService;

  Application(this.config);

  Future<Handler> createHandler() async {
    // Initialize dependencies
    _userRepository = UserRepository();
    _userService = UserService(_userRepository);

    // Create handlers
    final healthHandler = HealthHandler();
    final usersHandler = UsersHandler(_userService);

    // Build router
    final router = Router()
      ..mount('/health', healthHandler.router.call)
      ..mount('/api/v1/users', usersHandler.router.call);

    // Build middleware pipeline
    final pipeline = const Pipeline()
        .addMiddleware(loggingMiddleware())
        .addMiddleware(corsMiddleware())
        .addMiddleware(handleErrors())
        .addHandler(router.call);

    return pipeline;
  }

  Future<void> close() async {
    await _userRepository.close();
  }
}

// Error handling middleware
Middleware handleErrors() {
  return (Handler innerHandler) {
    return (Request request) async {
      try {
        return await innerHandler(request);
      } on NotFoundException catch (e) {
        return Response.notFound(
          jsonEncode({'error': e.message}),
          headers: {'Content-Type': 'application/json'},
        );
      } on ValidationException catch (e) {
        return Response(
          422,
          body: jsonEncode({'error': e.message, 'errors': e.errors}),
          headers: {'Content-Type': 'application/json'},
        );
      } on UnauthorizedException {
        return Response.forbidden(
          jsonEncode({'error': 'Unauthorized'}),
          headers: {'Content-Type': 'application/json'},
        );
      } catch (e, stack) {
        print('Error: $e\n$stack');
        return Response.internalServerError(
          body: jsonEncode({'error': 'Internal server error'}),
          headers: {'Content-Type': 'application/json'},
        );
      }
    };
  };
}
```

### Configuration

```dart
// lib/src/config/config.dart
import 'dart:io';

class Config {
  final int port;
  final String databaseUrl;
  final String jwtSecret;
  final bool isDevelopment;

  const Config({
    required this.port,
    required this.databaseUrl,
    required this.jwtSecret,
    required this.isDevelopment,
  });

  factory Config.fromEnvironment() {
    return Config(
      port: int.parse(Platform.environment['PORT'] ?? '8080'),
      databaseUrl: Platform.environment['DATABASE_URL'] ??
          'postgres://localhost/myapp',
      jwtSecret: Platform.environment['JWT_SECRET'] ??
          'development-secret-key',
      isDevelopment: Platform.environment['DART_ENV'] != 'production',
    );
  }
}
```

---

## Handlers (Controllers)

### Basic Handler

```dart
// lib/src/handlers/health_handler.dart
import 'dart:convert';

import 'package:shelf/shelf.dart';
import 'package:shelf_router/shelf_router.dart';

class HealthHandler {
  Router get router {
    final router = Router();

    router.get('/', _health);
    router.get('/ready', _ready);
    router.get('/live', _live);

    return router;
  }

  Response _health(Request request) {
    return Response.ok(
      jsonEncode({
        'status': 'healthy',
        'timestamp': DateTime.now().toIso8601String(),
      }),
      headers: {'Content-Type': 'application/json'},
    );
  }

  Response _ready(Request request) {
    // Check database connection, etc.
    return Response.ok(
      jsonEncode({'ready': true}),
      headers: {'Content-Type': 'application/json'},
    );
  }

  Response _live(Request request) {
    return Response.ok(
      jsonEncode({'live': true}),
      headers: {'Content-Type': 'application/json'},
    );
  }
}
```

### CRUD Handler

```dart
// lib/src/handlers/users_handler.dart
import 'dart:convert';

import 'package:shelf/shelf.dart';
import 'package:shelf_router/shelf_router.dart';

import '../models/user.dart';
import '../services/user_service.dart';
import '../middleware/auth_middleware.dart';

class UsersHandler {
  final UserService _userService;

  UsersHandler(this._userService);

  Router get router {
    final router = Router();

    // Public routes
    router.post('/register', _register);
    router.post('/login', _login);

    // Protected routes (add auth middleware individually or to group)
    router.get('/', _authMiddleware(_getAll));
    router.get('/<id>', _authMiddleware(_getById));
    router.put('/<id>', _authMiddleware(_update));
    router.delete('/<id>', _authMiddleware(_delete));

    return router;
  }

  Handler _authMiddleware(Handler handler) {
    return const Pipeline()
        .addMiddleware(authMiddleware())
        .addHandler(handler);
  }

  Future<Response> _register(Request request) async {
    final body = await request.readAsString();
    final json = jsonDecode(body) as Map<String, dynamic>;

    final user = await _userService.register(
      email: json['email'] as String,
      password: json['password'] as String,
      name: json['name'] as String,
    );

    return Response(
      201,
      body: jsonEncode(user.toJson()),
      headers: {'Content-Type': 'application/json'},
    );
  }

  Future<Response> _login(Request request) async {
    final body = await request.readAsString();
    final json = jsonDecode(body) as Map<String, dynamic>;

    final result = await _userService.login(
      email: json['email'] as String,
      password: json['password'] as String,
    );

    return Response.ok(
      jsonEncode({
        'token': result.token,
        'user': result.user.toJson(),
      }),
      headers: {'Content-Type': 'application/json'},
    );
  }

  Future<Response> _getAll(Request request) async {
    final page = int.tryParse(request.url.queryParameters['page'] ?? '1') ?? 1;
    final limit = int.tryParse(request.url.queryParameters['limit'] ?? '20') ?? 20;

    final users = await _userService.getAll(page: page, limit: limit);

    return Response.ok(
      jsonEncode({
        'data': users.map((u) => u.toJson()).toList(),
        'page': page,
        'limit': limit,
      }),
      headers: {'Content-Type': 'application/json'},
    );
  }

  Future<Response> _getById(Request request, String id) async {
    final user = await _userService.getById(id);

    if (user == null) {
      throw NotFoundException('User not found');
    }

    return Response.ok(
      jsonEncode(user.toJson()),
      headers: {'Content-Type': 'application/json'},
    );
  }

  Future<Response> _update(Request request, String id) async {
    final currentUser = request.context['user'] as User;

    // Check authorization
    if (currentUser.id != id && !currentUser.isAdmin) {
      throw UnauthorizedException();
    }

    final body = await request.readAsString();
    final json = jsonDecode(body) as Map<String, dynamic>;

    final user = await _userService.update(id, json);

    return Response.ok(
      jsonEncode(user.toJson()),
      headers: {'Content-Type': 'application/json'},
    );
  }

  Future<Response> _delete(Request request, String id) async {
    final currentUser = request.context['user'] as User;

    if (currentUser.id != id && !currentUser.isAdmin) {
      throw UnauthorizedException();
    }

    await _userService.delete(id);

    return Response(204);
  }
}
```

---

## Middleware

### Logging Middleware

```dart
// lib/src/middleware/logging_middleware.dart
import 'package:shelf/shelf.dart';

Middleware loggingMiddleware() {
  return (Handler innerHandler) {
    return (Request request) async {
      final stopwatch = Stopwatch()..start();

      print('[${DateTime.now()}] ${request.method} ${request.requestedUri}');

      final response = await innerHandler(request);

      stopwatch.stop();
      print(
        '[${DateTime.now()}] ${request.method} ${request.requestedUri} '
        '${response.statusCode} ${stopwatch.elapsedMilliseconds}ms',
      );

      return response;
    };
  };
}
```

### CORS Middleware

```dart
// lib/src/middleware/cors_middleware.dart
import 'package:shelf/shelf.dart';

Middleware corsMiddleware({
  List<String> allowedOrigins = const ['*'],
  List<String> allowedMethods = const ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  List<String> allowedHeaders = const ['Content-Type', 'Authorization'],
}) {
  return (Handler innerHandler) {
    return (Request request) async {
      // Handle preflight requests
      if (request.method == 'OPTIONS') {
        return Response.ok(
          '',
          headers: _corsHeaders(allowedOrigins, allowedMethods, allowedHeaders),
        );
      }

      final response = await innerHandler(request);

      return response.change(
        headers: {
          ...response.headers,
          ..._corsHeaders(allowedOrigins, allowedMethods, allowedHeaders),
        },
      );
    };
  };
}

Map<String, String> _corsHeaders(
  List<String> origins,
  List<String> methods,
  List<String> headers,
) {
  return {
    'Access-Control-Allow-Origin': origins.join(', '),
    'Access-Control-Allow-Methods': methods.join(', '),
    'Access-Control-Allow-Headers': headers.join(', '),
    'Access-Control-Max-Age': '86400',
  };
}
```

### Authentication Middleware

```dart
// lib/src/middleware/auth_middleware.dart
import 'dart:convert';

import 'package:shelf/shelf.dart';
import 'package:dart_jsonwebtoken/dart_jsonwebtoken.dart';

import '../models/user.dart';
import '../config/config.dart';

Middleware authMiddleware() {
  return (Handler innerHandler) {
    return (Request request) async {
      final authHeader = request.headers['Authorization'];

      if (authHeader == null || !authHeader.startsWith('Bearer ')) {
        throw UnauthorizedException();
      }

      final token = authHeader.substring(7);

      try {
        final config = request.context['config'] as Config;
        final jwt = JWT.verify(token, SecretKey(config.jwtSecret));

        final payload = jwt.payload as Map<String, dynamic>;
        final user = User.fromJson(payload['user'] as Map<String, dynamic>);

        // Add user to request context
        final updatedRequest = request.change(
          context: {...request.context, 'user': user},
        );

        return await innerHandler(updatedRequest);
      } on JWTExpiredException {
        throw UnauthorizedException('Token expired');
      } on JWTException {
        throw UnauthorizedException('Invalid token');
      }
    };
  };
}

// Optional auth - doesn't fail if no token
Middleware optionalAuthMiddleware() {
  return (Handler innerHandler) {
    return (Request request) async {
      final authHeader = request.headers['Authorization'];

      if (authHeader != null && authHeader.startsWith('Bearer ')) {
        final token = authHeader.substring(7);

        try {
          final config = request.context['config'] as Config;
          final jwt = JWT.verify(token, SecretKey(config.jwtSecret));

          final payload = jwt.payload as Map<String, dynamic>;
          final user = User.fromJson(payload['user'] as Map<String, dynamic>);

          return await innerHandler(
            request.change(context: {...request.context, 'user': user}),
          );
        } catch (_) {
          // Ignore invalid tokens
        }
      }

      return await innerHandler(request);
    };
  };
}

// Exceptions
class UnauthorizedException implements Exception {
  final String message;
  UnauthorizedException([this.message = 'Unauthorized']);
}

class NotFoundException implements Exception {
  final String message;
  NotFoundException(this.message);
}

class ValidationException implements Exception {
  final String message;
  final Map<String, List<String>> errors;
  ValidationException(this.message, [this.errors = const {}]);
}
```

### Rate Limiting Middleware

```dart
// lib/src/middleware/rate_limit_middleware.dart
import 'dart:async';

import 'package:shelf/shelf.dart';

class RateLimiter {
  final int maxRequests;
  final Duration window;
  final Map<String, List<DateTime>> _requests = {};

  RateLimiter({
    this.maxRequests = 100,
    this.window = const Duration(minutes: 1),
  });

  bool isAllowed(String key) {
    final now = DateTime.now();
    final windowStart = now.subtract(window);

    _requests[key] = (_requests[key] ?? [])
        .where((time) => time.isAfter(windowStart))
        .toList();

    if (_requests[key]!.length >= maxRequests) {
      return false;
    }

    _requests[key]!.add(now);
    return true;
  }

  int remaining(String key) {
    final now = DateTime.now();
    final windowStart = now.subtract(window);

    final count = (_requests[key] ?? [])
        .where((time) => time.isAfter(windowStart))
        .length;

    return (maxRequests - count).clamp(0, maxRequests);
  }
}

Middleware rateLimitMiddleware({
  int maxRequests = 100,
  Duration window = const Duration(minutes: 1),
}) {
  final limiter = RateLimiter(maxRequests: maxRequests, window: window);

  return (Handler innerHandler) {
    return (Request request) async {
      final clientIp = request.headers['X-Forwarded-For'] ??
          request.headers['X-Real-IP'] ??
          'unknown';

      if (!limiter.isAllowed(clientIp)) {
        return Response(
          429,
          body: jsonEncode({'error': 'Too many requests'}),
          headers: {
            'Content-Type': 'application/json',
            'Retry-After': '60',
            'X-RateLimit-Remaining': '0',
          },
        );
      }

      final response = await innerHandler(request);

      return response.change(
        headers: {
          ...response.headers,
          'X-RateLimit-Remaining': '${limiter.remaining(clientIp)}',
          'X-RateLimit-Limit': '$maxRequests',
        },
      );
    };
  };
}
```

---

## Models

### User Model

```dart
// lib/src/models/user.dart
import 'package:freezed_annotation/freezed_annotation.dart';

part 'user.freezed.dart';
part 'user.g.dart';

@freezed
class User with _$User {
  const User._();

  const factory User({
    required String id,
    required String email,
    required String name,
    @Default(false) bool isAdmin,
    @JsonKey(includeToJson: false) String? passwordHash,
    required DateTime createdAt,
    DateTime? updatedAt,
  }) = _User;

  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);

  Map<String, dynamic> toPublicJson() {
    return {
      'id': id,
      'email': email,
      'name': name,
      'isAdmin': isAdmin,
      'createdAt': createdAt.toIso8601String(),
      if (updatedAt != null) 'updatedAt': updatedAt!.toIso8601String(),
    };
  }
}

@freezed
class LoginResult with _$LoginResult {
  const factory LoginResult({
    required String token,
    required User user,
  }) = _LoginResult;
}
```

---

## Services

### User Service

```dart
// lib/src/services/user_service.dart
import 'package:dart_jsonwebtoken/dart_jsonwebtoken.dart';
import 'package:bcrypt/bcrypt.dart';

import '../config/config.dart';
import '../models/user.dart';
import '../repositories/user_repository.dart';
import '../middleware/auth_middleware.dart';

class UserService {
  final UserRepository _repository;
  final Config _config;

  UserService(this._repository, this._config);

  Future<User> register({
    required String email,
    required String password,
    required String name,
  }) async {
    // Validate
    _validateEmail(email);
    _validatePassword(password);

    // Check if email exists
    final existing = await _repository.findByEmail(email);
    if (existing != null) {
      throw ValidationException(
        'Email already exists',
        {'email': ['Email is already registered']},
      );
    }

    // Hash password
    final passwordHash = BCrypt.hashpw(password, BCrypt.gensalt());

    // Create user
    final user = await _repository.create(
      email: email.toLowerCase().trim(),
      passwordHash: passwordHash,
      name: name.trim(),
    );

    return user;
  }

  Future<LoginResult> login({
    required String email,
    required String password,
  }) async {
    final user = await _repository.findByEmail(email.toLowerCase().trim());

    if (user == null || user.passwordHash == null) {
      throw UnauthorizedException('Invalid credentials');
    }

    final isValid = BCrypt.checkpw(password, user.passwordHash!);
    if (!isValid) {
      throw UnauthorizedException('Invalid credentials');
    }

    // Generate JWT
    final jwt = JWT(
      {
        'user': user.toPublicJson(),
        'iat': DateTime.now().millisecondsSinceEpoch ~/ 1000,
      },
      issuer: 'myapp',
      subject: user.id,
    );

    final token = jwt.sign(
      SecretKey(_config.jwtSecret),
      expiresIn: const Duration(days: 7),
    );

    return LoginResult(token: token, user: user);
  }

  Future<List<User>> getAll({int page = 1, int limit = 20}) async {
    return _repository.findAll(
      offset: (page - 1) * limit,
      limit: limit,
    );
  }

  Future<User?> getById(String id) async {
    return _repository.findById(id);
  }

  Future<User> update(String id, Map<String, dynamic> data) async {
    final user = await _repository.findById(id);
    if (user == null) {
      throw NotFoundException('User not found');
    }

    return _repository.update(
      id,
      name: data['name'] as String?,
      email: data['email'] as String?,
    );
  }

  Future<void> delete(String id) async {
    final user = await _repository.findById(id);
    if (user == null) {
      throw NotFoundException('User not found');
    }

    await _repository.delete(id);
  }

  void _validateEmail(String email) {
    final emailRegex = RegExp(r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$');
    if (!emailRegex.hasMatch(email)) {
      throw ValidationException(
        'Invalid email',
        {'email': ['Please enter a valid email address']},
      );
    }
  }

  void _validatePassword(String password) {
    if (password.length < 8) {
      throw ValidationException(
        'Invalid password',
        {'password': ['Password must be at least 8 characters']},
      );
    }
  }
}
```

---

## Repositories

### User Repository

```dart
// lib/src/repositories/user_repository.dart
import '../models/user.dart';

abstract class UserRepository {
  Future<User> create({
    required String email,
    required String passwordHash,
    required String name,
  });

  Future<User?> findById(String id);
  Future<User?> findByEmail(String email);
  Future<List<User>> findAll({int offset = 0, int limit = 20});
  Future<User> update(String id, {String? name, String? email});
  Future<void> delete(String id);
  Future<void> close();
}

// In-memory implementation for example
class InMemoryUserRepository implements UserRepository {
  final Map<String, User> _users = {};
  int _idCounter = 0;

  @override
  Future<User> create({
    required String email,
    required String passwordHash,
    required String name,
  }) async {
    final id = (++_idCounter).toString();
    final user = User(
      id: id,
      email: email,
      name: name,
      passwordHash: passwordHash,
      createdAt: DateTime.now(),
    );
    _users[id] = user;
    return user;
  }

  @override
  Future<User?> findById(String id) async {
    return _users[id];
  }

  @override
  Future<User?> findByEmail(String email) async {
    return _users.values.where((u) => u.email == email).firstOrNull;
  }

  @override
  Future<List<User>> findAll({int offset = 0, int limit = 20}) async {
    return _users.values.skip(offset).take(limit).toList();
  }

  @override
  Future<User> update(String id, {String? name, String? email}) async {
    final user = _users[id];
    if (user == null) throw Exception('User not found');

    final updated = user.copyWith(
      name: name ?? user.name,
      email: email ?? user.email,
      updatedAt: DateTime.now(),
    );
    _users[id] = updated;
    return updated;
  }

  @override
  Future<void> delete(String id) async {
    _users.remove(id);
  }

  @override
  Future<void> close() async {}
}
```

---

## WebSocket Support

```dart
// lib/src/handlers/websocket_handler.dart
import 'dart:convert';

import 'package:shelf/shelf.dart';
import 'package:shelf_web_socket/shelf_web_socket.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class WebSocketHandler {
  final Set<WebSocketChannel> _clients = {};

  Handler get handler => webSocketHandler(_onConnection);

  void _onConnection(WebSocketChannel webSocket) {
    _clients.add(webSocket);
    print('Client connected. Total: ${_clients.length}');

    webSocket.stream.listen(
      (message) => _handleMessage(webSocket, message),
      onDone: () {
        _clients.remove(webSocket);
        print('Client disconnected. Total: ${_clients.length}');
      },
      onError: (error) {
        print('WebSocket error: $error');
        _clients.remove(webSocket);
      },
    );
  }

  void _handleMessage(WebSocketChannel sender, dynamic message) {
    try {
      final data = jsonDecode(message as String) as Map<String, dynamic>;
      final type = data['type'] as String;

      switch (type) {
        case 'ping':
          sender.sink.add(jsonEncode({'type': 'pong'}));
          break;
        case 'broadcast':
          _broadcast(data['payload']);
          break;
        default:
          sender.sink.add(jsonEncode({
            'type': 'error',
            'message': 'Unknown message type',
          }));
      }
    } catch (e) {
      sender.sink.add(jsonEncode({
        'type': 'error',
        'message': 'Invalid message format',
      }));
    }
  }

  void _broadcast(dynamic payload) {
    final message = jsonEncode({
      'type': 'message',
      'payload': payload,
    });

    for (final client in _clients) {
      client.sink.add(message);
    }
  }
}
```

---

## Testing

### Handler Tests

```dart
// test/handlers/users_handler_test.dart
import 'dart:convert';

import 'package:shelf/shelf.dart';
import 'package:test/test.dart';
import 'package:mocktail/mocktail.dart';

import 'package:myapp/src/handlers/users_handler.dart';
import 'package:myapp/src/services/user_service.dart';
import 'package:myapp/src/models/user.dart';

class MockUserService extends Mock implements UserService {}

void main() {
  late MockUserService mockService;
  late UsersHandler handler;

  setUp(() {
    mockService = MockUserService();
    handler = UsersHandler(mockService);
  });

  group('UsersHandler', () {
    group('POST /register', () {
      test('returns 201 with created user', () async {
        final user = User(
          id: '1',
          email: 'test@example.com',
          name: 'Test User',
          createdAt: DateTime.now(),
        );

        when(() => mockService.register(
              email: any(named: 'email'),
              password: any(named: 'password'),
              name: any(named: 'name'),
            )).thenAnswer((_) async => user);

        final request = Request(
          'POST',
          Uri.parse('http://localhost/register'),
          body: jsonEncode({
            'email': 'test@example.com',
            'password': 'password123',
            'name': 'Test User',
          }),
        );

        final response = await handler.router.call(request);

        expect(response.statusCode, equals(201));

        final body = jsonDecode(await response.readAsString());
        expect(body['email'], equals('test@example.com'));
      });

      test('returns 422 with validation errors', () async {
        when(() => mockService.register(
              email: any(named: 'email'),
              password: any(named: 'password'),
              name: any(named: 'name'),
            )).thenThrow(ValidationException(
          'Invalid email',
          {'email': ['Please enter a valid email']},
        ));

        final request = Request(
          'POST',
          Uri.parse('http://localhost/register'),
          body: jsonEncode({
            'email': 'invalid',
            'password': 'password123',
            'name': 'Test',
          }),
        );

        // Need to wrap with error handler
        final handlerWithErrors = const Pipeline()
            .addMiddleware(handleErrors())
            .addHandler(handler.router.call);

        final response = await handlerWithErrors(request);

        expect(response.statusCode, equals(422));
      });
    });

    group('GET /', () {
      test('returns paginated users', () async {
        final users = [
          User(id: '1', email: 'a@test.com', name: 'A', createdAt: DateTime.now()),
          User(id: '2', email: 'b@test.com', name: 'B', createdAt: DateTime.now()),
        ];

        when(() => mockService.getAll(page: any(named: 'page'), limit: any(named: 'limit')))
            .thenAnswer((_) async => users);

        // Add authenticated user to context
        final request = Request(
          'GET',
          Uri.parse('http://localhost/?page=1&limit=20'),
        ).change(context: {
          'user': User(id: '1', email: 'admin@test.com', name: 'Admin', createdAt: DateTime.now()),
        });

        final response = await handler.router.call(request);

        expect(response.statusCode, equals(200));

        final body = jsonDecode(await response.readAsString());
        expect(body['data'], hasLength(2));
      });
    });
  });
}
```

---

## Configuration

### pubspec.yaml

```yaml
name: myapp
description: Shelf backend application
version: 1.0.0

environment:
  sdk: '>=3.0.0 <4.0.0'

dependencies:
  shelf: ^1.4.0
  shelf_router: ^1.1.0
  shelf_web_socket: ^1.0.0

  # Data
  freezed_annotation: ^2.4.0
  json_annotation: ^4.8.0

  # Auth
  dart_jsonwebtoken: ^2.12.0
  bcrypt: ^1.1.0

  # Utils
  args: ^2.4.0

dev_dependencies:
  test: ^1.24.0
  mocktail: ^1.0.0

  # Code generation
  build_runner: ^2.4.0
  freezed: ^2.4.0
  json_serializable: ^6.7.0

executables:
  server: server
```

### Dockerfile

```dockerfile
# Build stage
FROM dart:stable AS build

WORKDIR /app

# Copy dependencies
COPY pubspec.* ./
RUN dart pub get

# Copy source
COPY . .

# Generate code
RUN dart run build_runner build --delete-conflicting-outputs

# Build
RUN dart compile exe bin/server.dart -o bin/server

# Runtime stage
FROM scratch

COPY --from=build /runtime/ /
COPY --from=build /app/bin/server /app/bin/server

EXPOSE 8080

CMD ["/app/bin/server"]
```

---

## Commands

```bash
# Run development server
dart run bin/server.dart

# Run with hot reload (using dart_dev or similar)
dart pub global activate webdev
webdev serve

# Run tests
dart test

# Generate code
dart run build_runner build

# Build executable
dart compile exe bin/server.dart -o server

# Format
dart format .

# Analyze
dart analyze
```

---

## References

- [Shelf Documentation](https://pub.dev/packages/shelf)
- [shelf_router](https://pub.dev/packages/shelf_router)
- [shelf_web_socket](https://pub.dev/packages/shelf_web_socket)
- [Dart Server Tutorial](https://dart.dev/tutorials/server/httpserver)
