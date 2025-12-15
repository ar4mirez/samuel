# Dart Frog Framework Guide

> **Framework**: Dart Frog 1.x
> **Language**: Dart 3.x
> **Type**: Backend Framework
> **Use Cases**: REST APIs, Full-Stack Dart, Serverless Functions

---

## Overview

Dart Frog is a fast, minimalistic backend framework for Dart. It's built on top of Shelf and provides file-based routing, middleware support, and a great developer experience with hot reload.

### Key Features
- File-based routing (like Next.js)
- Built-in middleware support
- Dependency injection
- Hot reload for development
- Easy deployment (Docker, Cloud Run, etc.)
- Full-stack Dart (pairs with Flutter)

---

## Project Structure

```
myapp/
├── routes/
│   ├── _middleware.dart          # Global middleware
│   ├── index.dart                # GET /
│   ├── health.dart               # GET /health
│   ├── api/
│   │   ├── _middleware.dart      # API middleware
│   │   ├── v1/
│   │   │   ├── _middleware.dart  # V1 middleware
│   │   │   ├── users/
│   │   │   │   ├── index.dart    # /api/v1/users
│   │   │   │   └── [id].dart     # /api/v1/users/:id
│   │   │   └── posts/
│   │   │       ├── index.dart
│   │   │       └── [id].dart
│   │   └── auth/
│   │       ├── login.dart
│   │       └── register.dart
│   └── ws.dart                   # WebSocket endpoint
├── lib/
│   ├── src/
│   │   ├── models/
│   │   │   └── user.dart
│   │   ├── repositories/
│   │   │   └── user_repository.dart
│   │   ├── services/
│   │   │   └── user_service.dart
│   │   └── middleware/
│   │       └── auth_provider.dart
│   └── myapp.dart
├── test/
│   └── routes/
│       └── api/
│           └── v1/
│               └── users/
│                   └── index_test.dart
├── pubspec.yaml
└── Dockerfile
```

---

## Route Handlers

### Basic Route

```dart
// routes/index.dart
import 'package:dart_frog/dart_frog.dart';

Response onRequest(RequestContext context) {
  return Response.json({
    'message': 'Welcome to My API',
    'version': '1.0.0',
  });
}
```

### Health Check

```dart
// routes/health.dart
import 'package:dart_frog/dart_frog.dart';

Response onRequest(RequestContext context) {
  return Response.json({
    'status': 'healthy',
    'timestamp': DateTime.now().toIso8601String(),
  });
}
```

### HTTP Method Handling

```dart
// routes/api/v1/users/index.dart
import 'dart:io';

import 'package:dart_frog/dart_frog.dart';

import 'package:myapp/myapp.dart';

Future<Response> onRequest(RequestContext context) async {
  return switch (context.request.method) {
    HttpMethod.get => _getUsers(context),
    HttpMethod.post => _createUser(context),
    _ => Future.value(Response(statusCode: HttpStatus.methodNotAllowed)),
  };
}

Future<Response> _getUsers(RequestContext context) async {
  final userService = context.read<UserService>();

  final page = int.tryParse(
    context.request.uri.queryParameters['page'] ?? '1',
  ) ?? 1;
  final limit = int.tryParse(
    context.request.uri.queryParameters['limit'] ?? '20',
  ) ?? 20;

  final users = await userService.getAll(page: page, limit: limit);

  return Response.json({
    'data': users.map((u) => u.toJson()).toList(),
    'page': page,
    'limit': limit,
  });
}

Future<Response> _createUser(RequestContext context) async {
  final userService = context.read<UserService>();

  final body = await context.request.json() as Map<String, dynamic>;

  try {
    final user = await userService.create(
      email: body['email'] as String,
      name: body['name'] as String,
      password: body['password'] as String,
    );

    return Response.json(
      user.toJson(),
      statusCode: HttpStatus.created,
    );
  } on ValidationException catch (e) {
    return Response.json(
      {'error': e.message, 'errors': e.errors},
      statusCode: HttpStatus.unprocessableEntity,
    );
  }
}
```

### Dynamic Routes

```dart
// routes/api/v1/users/[id].dart
import 'dart:io';

import 'package:dart_frog/dart_frog.dart';

import 'package:myapp/myapp.dart';

Future<Response> onRequest(RequestContext context, String id) async {
  return switch (context.request.method) {
    HttpMethod.get => _getUser(context, id),
    HttpMethod.put => _updateUser(context, id),
    HttpMethod.delete => _deleteUser(context, id),
    _ => Future.value(Response(statusCode: HttpStatus.methodNotAllowed)),
  };
}

Future<Response> _getUser(RequestContext context, String id) async {
  final userService = context.read<UserService>();

  final user = await userService.getById(id);

  if (user == null) {
    return Response.json(
      {'error': 'User not found'},
      statusCode: HttpStatus.notFound,
    );
  }

  return Response.json(user.toJson());
}

Future<Response> _updateUser(RequestContext context, String id) async {
  final userService = context.read<UserService>();
  final currentUser = context.read<User?>();

  // Authorization check
  if (currentUser == null || (currentUser.id != id && !currentUser.isAdmin)) {
    return Response.json(
      {'error': 'Forbidden'},
      statusCode: HttpStatus.forbidden,
    );
  }

  final body = await context.request.json() as Map<String, dynamic>;

  try {
    final user = await userService.update(id, body);
    return Response.json(user.toJson());
  } on NotFoundException {
    return Response.json(
      {'error': 'User not found'},
      statusCode: HttpStatus.notFound,
    );
  }
}

Future<Response> _deleteUser(RequestContext context, String id) async {
  final userService = context.read<UserService>();
  final currentUser = context.read<User?>();

  if (currentUser == null || (currentUser.id != id && !currentUser.isAdmin)) {
    return Response.json(
      {'error': 'Forbidden'},
      statusCode: HttpStatus.forbidden,
    );
  }

  try {
    await userService.delete(id);
    return Response(statusCode: HttpStatus.noContent);
  } on NotFoundException {
    return Response.json(
      {'error': 'User not found'},
      statusCode: HttpStatus.notFound,
    );
  }
}
```

### Nested Dynamic Routes

```dart
// routes/api/v1/users/[userId]/posts/[postId].dart
import 'package:dart_frog/dart_frog.dart';

Future<Response> onRequest(
  RequestContext context,
  String userId,
  String postId,
) async {
  final postService = context.read<PostService>();

  final post = await postService.getByUserAndId(userId, postId);

  if (post == null) {
    return Response.json(
      {'error': 'Post not found'},
      statusCode: 404,
    );
  }

  return Response.json(post.toJson());
}
```

---

## Middleware

### Global Middleware

```dart
// routes/_middleware.dart
import 'package:dart_frog/dart_frog.dart';

import 'package:myapp/myapp.dart';

Handler middleware(Handler handler) {
  return handler
      .use(requestLogger())
      .use(corsMiddleware())
      .use(errorHandler())
      .use(provider<Config>((_) => Config.fromEnvironment()))
      .use(provider<UserRepository>((_) => UserRepositoryImpl()))
      .use(
        provider<UserService>((context) => UserService(
          context.read<UserRepository>(),
          context.read<Config>(),
        )),
      );
}
```

### API Middleware

```dart
// routes/api/_middleware.dart
import 'package:dart_frog/dart_frog.dart';

Handler middleware(Handler handler) {
  return handler.use(jsonContentType());
}

Middleware jsonContentType() {
  return (handler) {
    return (context) async {
      final response = await handler(context);
      return response.copyWith(
        headers: {
          ...response.headers,
          'Content-Type': 'application/json',
        },
      );
    };
  };
}
```

### Authentication Middleware

```dart
// routes/api/v1/_middleware.dart
import 'package:dart_frog/dart_frog.dart';

import 'package:myapp/myapp.dart';

Handler middleware(Handler handler) {
  return handler.use(authProvider());
}

// lib/src/middleware/auth_provider.dart
import 'package:dart_frog/dart_frog.dart';
import 'package:dart_jsonwebtoken/dart_jsonwebtoken.dart';

Middleware authProvider() {
  return provider<User?>((context) {
    final authHeader = context.request.headers['Authorization'];

    if (authHeader == null || !authHeader.startsWith('Bearer ')) {
      return null;
    }

    final token = authHeader.substring(7);

    try {
      final config = context.read<Config>();
      final jwt = JWT.verify(token, SecretKey(config.jwtSecret));

      final payload = jwt.payload as Map<String, dynamic>;
      return User.fromJson(payload['user'] as Map<String, dynamic>);
    } catch (_) {
      return null;
    }
  });
}
```

### Protected Route Middleware

```dart
// routes/api/v1/users/_middleware.dart
import 'dart:io';

import 'package:dart_frog/dart_frog.dart';

import 'package:myapp/myapp.dart';

Handler middleware(Handler handler) {
  return handler.use(requireAuth());
}

Middleware requireAuth() {
  return (handler) {
    return (context) {
      final user = context.read<User?>();

      if (user == null) {
        return Response.json(
          {'error': 'Unauthorized'},
          statusCode: HttpStatus.unauthorized,
        );
      }

      return handler(context);
    };
  };
}
```

### Request Logger

```dart
// lib/src/middleware/request_logger.dart
import 'package:dart_frog/dart_frog.dart';

Middleware requestLogger() {
  return (handler) {
    return (context) async {
      final stopwatch = Stopwatch()..start();
      final request = context.request;

      print('[${DateTime.now()}] ${request.method.value} ${request.uri}');

      final response = await handler(context);

      stopwatch.stop();
      print(
        '[${DateTime.now()}] ${request.method.value} ${request.uri} '
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
import 'package:dart_frog/dart_frog.dart';

Middleware corsMiddleware({
  List<String> allowedOrigins = const ['*'],
  List<String> allowedMethods = const ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  List<String> allowedHeaders = const ['Content-Type', 'Authorization'],
}) {
  final corsHeaders = {
    'Access-Control-Allow-Origin': allowedOrigins.join(', '),
    'Access-Control-Allow-Methods': allowedMethods.join(', '),
    'Access-Control-Allow-Headers': allowedHeaders.join(', '),
    'Access-Control-Max-Age': '86400',
  };

  return (handler) {
    return (context) async {
      // Handle preflight
      if (context.request.method == HttpMethod.options) {
        return Response(headers: corsHeaders);
      }

      final response = await handler(context);

      return response.copyWith(
        headers: {...response.headers, ...corsHeaders},
      );
    };
  };
}
```

### Error Handler

```dart
// lib/src/middleware/error_handler.dart
import 'dart:io';

import 'package:dart_frog/dart_frog.dart';

import '../exceptions.dart';

Middleware errorHandler() {
  return (handler) {
    return (context) async {
      try {
        return await handler(context);
      } on ValidationException catch (e) {
        return Response.json(
          {'error': e.message, 'errors': e.errors},
          statusCode: HttpStatus.unprocessableEntity,
        );
      } on NotFoundException catch (e) {
        return Response.json(
          {'error': e.message},
          statusCode: HttpStatus.notFound,
        );
      } on UnauthorizedException catch (e) {
        return Response.json(
          {'error': e.message},
          statusCode: HttpStatus.unauthorized,
        );
      } on ForbiddenException catch (e) {
        return Response.json(
          {'error': e.message},
          statusCode: HttpStatus.forbidden,
        );
      } catch (e, stack) {
        print('Error: $e\n$stack');
        return Response.json(
          {'error': 'Internal server error'},
          statusCode: HttpStatus.internalServerError,
        );
      }
    };
  };
}
```

---

## Dependency Injection

### Provider Pattern

```dart
// routes/_middleware.dart
import 'package:dart_frog/dart_frog.dart';

Handler middleware(Handler handler) {
  return handler
      // Configuration
      .use(provider<Config>((_) => Config.fromEnvironment()))

      // Database connection
      .use(provider<Database>((context) {
        final config = context.read<Config>();
        return Database(config.databaseUrl);
      }))

      // Repositories
      .use(provider<UserRepository>((context) {
        return UserRepositoryImpl(context.read<Database>());
      }))
      .use(provider<PostRepository>((context) {
        return PostRepositoryImpl(context.read<Database>());
      }))

      // Services
      .use(provider<UserService>((context) {
        return UserService(
          context.read<UserRepository>(),
          context.read<Config>(),
        );
      }))
      .use(provider<PostService>((context) {
        return PostService(
          context.read<PostRepository>(),
          context.read<UserRepository>(),
        );
      }));
}
```

### Using Dependencies in Routes

```dart
// routes/api/v1/users/index.dart
import 'package:dart_frog/dart_frog.dart';

import 'package:myapp/myapp.dart';

Future<Response> onRequest(RequestContext context) async {
  // Read dependencies from context
  final userService = context.read<UserService>();
  final config = context.read<Config>();
  final currentUser = context.read<User?>(); // nullable if optional auth

  // Use dependencies
  final users = await userService.getAll();

  return Response.json({
    'data': users.map((u) => u.toJson()).toList(),
    'environment': config.isDevelopment ? 'development' : 'production',
  });
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

  Map<String, dynamic> toPublicJson() => {
        'id': id,
        'email': email,
        'name': name,
        'isAdmin': isAdmin,
        'createdAt': createdAt.toIso8601String(),
        if (updatedAt != null) 'updatedAt': updatedAt!.toIso8601String(),
      };
}
```

### Exceptions

```dart
// lib/src/exceptions.dart
class AppException implements Exception {
  final String message;
  const AppException(this.message);

  @override
  String toString() => message;
}

class ValidationException extends AppException {
  final Map<String, List<String>> errors;

  const ValidationException(super.message, [this.errors = const {}]);
}

class NotFoundException extends AppException {
  const NotFoundException(super.message);
}

class UnauthorizedException extends AppException {
  const UnauthorizedException([super.message = 'Unauthorized']);
}

class ForbiddenException extends AppException {
  const ForbiddenException([super.message = 'Forbidden']);
}
```

---

## Authentication Routes

### Login

```dart
// routes/api/v1/auth/login.dart
import 'dart:io';

import 'package:dart_frog/dart_frog.dart';

import 'package:myapp/myapp.dart';

Future<Response> onRequest(RequestContext context) async {
  if (context.request.method != HttpMethod.post) {
    return Response(statusCode: HttpStatus.methodNotAllowed);
  }

  final userService = context.read<UserService>();
  final body = await context.request.json() as Map<String, dynamic>;

  try {
    final result = await userService.login(
      email: body['email'] as String,
      password: body['password'] as String,
    );

    return Response.json({
      'token': result.token,
      'user': result.user.toPublicJson(),
    });
  } on UnauthorizedException {
    return Response.json(
      {'error': 'Invalid credentials'},
      statusCode: HttpStatus.unauthorized,
    );
  }
}
```

### Register

```dart
// routes/api/v1/auth/register.dart
import 'dart:io';

import 'package:dart_frog/dart_frog.dart';

import 'package:myapp/myapp.dart';

Future<Response> onRequest(RequestContext context) async {
  if (context.request.method != HttpMethod.post) {
    return Response(statusCode: HttpStatus.methodNotAllowed);
  }

  final userService = context.read<UserService>();
  final body = await context.request.json() as Map<String, dynamic>;

  try {
    final user = await userService.register(
      email: body['email'] as String,
      password: body['password'] as String,
      name: body['name'] as String,
    );

    return Response.json(
      user.toPublicJson(),
      statusCode: HttpStatus.created,
    );
  } on ValidationException catch (e) {
    return Response.json(
      {'error': e.message, 'errors': e.errors},
      statusCode: HttpStatus.unprocessableEntity,
    );
  }
}
```

### Current User

```dart
// routes/api/v1/auth/me.dart
import 'dart:io';

import 'package:dart_frog/dart_frog.dart';

import 'package:myapp/myapp.dart';

Response onRequest(RequestContext context) {
  if (context.request.method != HttpMethod.get) {
    return Response(statusCode: HttpStatus.methodNotAllowed);
  }

  final user = context.read<User?>();

  if (user == null) {
    return Response.json(
      {'error': 'Unauthorized'},
      statusCode: HttpStatus.unauthorized,
    );
  }

  return Response.json(user.toPublicJson());
}
```

---

## WebSocket Support

```dart
// routes/ws.dart
import 'dart:convert';

import 'package:dart_frog/dart_frog.dart';
import 'package:dart_frog_web_socket/dart_frog_web_socket.dart';

Future<Response> onRequest(RequestContext context) async {
  final handler = webSocketHandler((channel, protocol) {
    // Handle new connection
    print('Client connected');

    // Listen for messages
    channel.stream.listen(
      (message) {
        try {
          final data = jsonDecode(message as String) as Map<String, dynamic>;
          _handleMessage(channel, data);
        } catch (e) {
          channel.sink.add(jsonEncode({
            'type': 'error',
            'message': 'Invalid message format',
          }));
        }
      },
      onDone: () => print('Client disconnected'),
      onError: (error) => print('WebSocket error: $error'),
    );
  });

  return handler(context);
}

void _handleMessage(WebSocketChannel channel, Map<String, dynamic> data) {
  final type = data['type'] as String?;

  switch (type) {
    case 'ping':
      channel.sink.add(jsonEncode({'type': 'pong'}));
      break;
    case 'echo':
      channel.sink.add(jsonEncode({
        'type': 'echo',
        'data': data['data'],
      }));
      break;
    default:
      channel.sink.add(jsonEncode({
        'type': 'error',
        'message': 'Unknown message type: $type',
      }));
  }
}
```

---

## Testing

### Route Tests

```dart
// test/routes/api/v1/users/index_test.dart
import 'dart:io';

import 'package:dart_frog/dart_frog.dart';
import 'package:mocktail/mocktail.dart';
import 'package:test/test.dart';

import 'package:myapp/myapp.dart';
import '../../../../../routes/api/v1/users/index.dart' as route;

class MockRequestContext extends Mock implements RequestContext {}

class MockUserService extends Mock implements UserService {}

void main() {
  late MockRequestContext context;
  late MockUserService mockUserService;

  setUp(() {
    context = MockRequestContext();
    mockUserService = MockUserService();

    when(() => context.read<UserService>()).thenReturn(mockUserService);
  });

  group('GET /api/v1/users', () {
    test('returns list of users', () async {
      final users = [
        User(id: '1', email: 'a@test.com', name: 'A', createdAt: DateTime.now()),
        User(id: '2', email: 'b@test.com', name: 'B', createdAt: DateTime.now()),
      ];

      when(() => mockUserService.getAll(page: any(named: 'page'), limit: any(named: 'limit')))
          .thenAnswer((_) async => users);

      when(() => context.request).thenReturn(
        Request.get(Uri.parse('http://localhost/api/v1/users?page=1&limit=20')),
      );

      final response = await route.onRequest(context);

      expect(response.statusCode, equals(HttpStatus.ok));

      final body = await response.json() as Map<String, dynamic>;
      expect(body['data'], hasLength(2));
    });
  });

  group('POST /api/v1/users', () {
    test('creates user with valid data', () async {
      final user = User(
        id: '1',
        email: 'test@example.com',
        name: 'Test',
        createdAt: DateTime.now(),
      );

      when(() => mockUserService.create(
            email: any(named: 'email'),
            name: any(named: 'name'),
            password: any(named: 'password'),
          )).thenAnswer((_) async => user);

      when(() => context.request).thenReturn(
        Request.post(
          Uri.parse('http://localhost/api/v1/users'),
          body: '{"email":"test@example.com","name":"Test","password":"password123"}',
        ),
      );

      final response = await route.onRequest(context);

      expect(response.statusCode, equals(HttpStatus.created));
    });

    test('returns 422 with validation errors', () async {
      when(() => mockUserService.create(
            email: any(named: 'email'),
            name: any(named: 'name'),
            password: any(named: 'password'),
          )).thenThrow(
        ValidationException('Validation failed', {
          'email': ['Invalid email format'],
        }),
      );

      when(() => context.request).thenReturn(
        Request.post(
          Uri.parse('http://localhost/api/v1/users'),
          body: '{"email":"invalid","name":"Test","password":"password123"}',
        ),
      );

      final response = await route.onRequest(context);

      expect(response.statusCode, equals(HttpStatus.unprocessableEntity));
    });
  });
}
```

### Middleware Tests

```dart
// test/middleware/auth_provider_test.dart
import 'package:dart_frog/dart_frog.dart';
import 'package:dart_jsonwebtoken/dart_jsonwebtoken.dart';
import 'package:mocktail/mocktail.dart';
import 'package:test/test.dart';

import 'package:myapp/myapp.dart';

class MockRequestContext extends Mock implements RequestContext {}

void main() {
  group('authProvider', () {
    test('returns null when no auth header', () {
      final context = MockRequestContext();
      final request = Request.get(Uri.parse('http://localhost/'));

      when(() => context.request).thenReturn(request);
      when(() => context.read<Config>()).thenReturn(
        Config(port: 8080, databaseUrl: '', jwtSecret: 'secret', isDevelopment: true),
      );

      final middleware = authProvider();
      final handler = middleware((ctx) => Response());

      // The provider returns User? - we need to test differently
      // This is a simplified example
    });

    test('returns user when valid token', () {
      final context = MockRequestContext();
      final config = Config(
        port: 8080,
        databaseUrl: '',
        jwtSecret: 'secret',
        isDevelopment: true,
      );

      final user = User(
        id: '1',
        email: 'test@test.com',
        name: 'Test',
        createdAt: DateTime.now(),
      );

      final jwt = JWT({'user': user.toJson()});
      final token = jwt.sign(SecretKey(config.jwtSecret));

      final request = Request.get(
        Uri.parse('http://localhost/'),
        headers: {'Authorization': 'Bearer $token'},
      );

      when(() => context.request).thenReturn(request);
      when(() => context.read<Config>()).thenReturn(config);

      // Test the provider logic
    });
  });
}
```

---

## Configuration

### pubspec.yaml

```yaml
name: myapp
description: Dart Frog backend application
version: 1.0.0
publish_to: none

environment:
  sdk: '>=3.0.0 <4.0.0'

dependencies:
  dart_frog: ^1.0.0
  dart_frog_web_socket: ^1.0.0

  # Data
  freezed_annotation: ^2.4.0
  json_annotation: ^4.8.0

  # Auth
  dart_jsonwebtoken: ^2.12.0
  bcrypt: ^1.1.0

dev_dependencies:
  dart_frog_test: ^1.0.0
  test: ^1.24.0
  mocktail: ^1.0.0
  very_good_analysis: ^5.0.0

  # Code generation
  build_runner: ^2.4.0
  freezed: ^2.4.0
  json_serializable: ^6.7.0
```

### analysis_options.yaml

```yaml
include: package:very_good_analysis/analysis_options.yaml

analyzer:
  exclude:
    - "**/*.g.dart"
    - "**/*.freezed.dart"
  errors:
    invalid_annotation_target: ignore

linter:
  rules:
    public_member_api_docs: false
```

### Dockerfile

```dockerfile
# Build stage
FROM dart:stable AS build

WORKDIR /app

# Install dart_frog
RUN dart pub global activate dart_frog_cli

# Copy dependencies
COPY pubspec.* ./
RUN dart pub get

# Copy source
COPY . .

# Generate code
RUN dart run build_runner build --delete-conflicting-outputs

# Build for production
RUN dart pub global run dart_frog_cli:dart_frog build

# Runtime stage
FROM scratch

COPY --from=build /runtime/ /
COPY --from=build /app/build/bin/server /app/bin/server

EXPOSE 8080

CMD ["/app/bin/server"]
```

### docker-compose.yaml

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - DATABASE_URL=postgres://postgres:password@db:5432/myapp
      - JWT_SECRET=${JWT_SECRET}
      - DART_ENV=production
    depends_on:
      - db

  db:
    image: postgres:15
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=myapp
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

volumes:
  postgres_data:
```

---

## Commands

```bash
# Create new project
dart_frog create myapp

# Run development server (with hot reload)
dart_frog dev

# Build for production
dart_frog build

# Run production build
./build/bin/server

# Generate new route
dart_frog new route /api/v1/users

# Generate new middleware
dart_frog new middleware auth

# Run tests
dart test

# Generate code
dart run build_runner build

# Format
dart format .

# Analyze
dart analyze
```

---

## Deployment

### Cloud Run

```bash
# Build and push
gcloud builds submit --tag gcr.io/PROJECT_ID/myapp

# Deploy
gcloud run deploy myapp \
  --image gcr.io/PROJECT_ID/myapp \
  --platform managed \
  --allow-unauthenticated \
  --set-env-vars "JWT_SECRET=your-secret"
```

### Fly.io

```toml
# fly.toml
app = "myapp"
primary_region = "iad"

[build]
  dockerfile = "Dockerfile"

[env]
  PORT = "8080"

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = true
  auto_start_machines = true
  min_machines_running = 0
```

```bash
fly launch
fly deploy
fly secrets set JWT_SECRET=your-secret
```

---

## References

- [Dart Frog Documentation](https://dartfrog.vgv.dev/)
- [Dart Frog GitHub](https://github.com/VeryGoodOpenSource/dart_frog)
- [Very Good Ventures Blog](https://verygood.ventures/blog)
- [Dart Server Tutorial](https://dart.dev/tutorials/server/httpserver)
