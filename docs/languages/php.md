# PHP Guide

> **Applies to**: PHP 8.1+, Laravel, Symfony, WordPress, APIs

---

## Core Principles

1. **Modern PHP**: Use PHP 8.1+ features (attributes, enums, readonly)
2. **Type Safety**: Strict types, typed properties, return types
3. **PSR Standards**: Follow PSR-1, PSR-4, PSR-12 coding standards
4. **Composer**: Dependency management and autoloading
5. **Framework Best Practices**: Follow Laravel/Symfony conventions

---

## Language-Specific Guardrails

### PHP Version & Setup
- ✓ Use PHP 8.1+ (8.3 recommended for new projects)
- ✓ Enable `declare(strict_types=1)` in all files
- ✓ Use Composer for dependency management
- ✓ Pin dependency versions in `composer.json`
- ✓ Include PHP version in `composer.json` `require` section

### Code Style (PSR-12)
- ✓ Follow PSR-12 Extended Coding Style
- ✓ Run `php-cs-fixer` or `phpcs` before every commit
- ✓ Use `PascalCase` for classes
- ✓ Use `camelCase` for methods and variables
- ✓ Use `SCREAMING_SNAKE_CASE` for constants
- ✓ 4-space indentation (not tabs)
- ✓ Line length: 120 characters max
- ✓ One class per file, file name matches class name

### Type Safety (PHP 8+)
- ✓ Declare strict types: `declare(strict_types=1);`
- ✓ Type all method parameters and return values
- ✓ Use union types: `string|int`
- ✓ Use nullable types: `?string` or `string|null`
- ✓ Use typed properties in classes
- ✓ Use `mixed` only when truly necessary

### Modern PHP Features (8.1+)
- ✓ Use enums for fixed sets of values
- ✓ Use readonly properties and classes
- ✓ Use constructor property promotion
- ✓ Use named arguments for clarity
- ✓ Use attributes instead of docblock annotations where possible
- ✓ Use match expressions instead of switch
- ✓ Use nullsafe operator: `$obj?->method()`
- ✓ Use first-class callable syntax: `$fn = strlen(...)`

### Error Handling
- ✓ Use exceptions for error handling (not return codes)
- ✓ Create custom exception classes for domain errors
- ✓ Catch specific exceptions, not generic `Exception`
- ✓ Always provide meaningful exception messages
- ✓ Log exceptions with context
- ✓ Don't suppress errors with `@`

### Security
- ✓ Never trust user input (validate and sanitize)
- ✓ Use prepared statements for database queries
- ✓ Use `htmlspecialchars()` or templating engine escaping
- ✓ Use `password_hash()` and `password_verify()` for passwords
- ✓ Use CSRF tokens for forms
- ✓ Set proper session configuration
- ✓ Keep dependencies updated

---

## Project Structure

### Laravel Standard
```
myproject/
├── app/
│   ├── Console/          # Artisan commands
│   ├── Exceptions/       # Custom exceptions
│   ├── Http/
│   │   ├── Controllers/
│   │   ├── Middleware/
│   │   └── Requests/     # Form requests (validation)
│   ├── Models/           # Eloquent models
│   ├── Providers/        # Service providers
│   └── Services/         # Business logic
├── config/               # Configuration files
├── database/
│   ├── factories/
│   ├── migrations/
│   └── seeders/
├── resources/
│   └── views/            # Blade templates
├── routes/
│   └── api.php
├── tests/
│   ├── Feature/
│   └── Unit/
├── composer.json
└── phpunit.xml
```

### PSR-4 Autoloading
```json
{
    "autoload": {
        "psr-4": {
            "App\\": "app/",
            "Database\\": "database/"
        }
    },
    "autoload-dev": {
        "psr-4": {
            "Tests\\": "tests/"
        }
    }
}
```

---

## Validation & Input Handling

### Laravel Validation
```php
<?php

declare(strict_types=1);

namespace App\Http\Requests;

use Illuminate\Foundation\Http\FormRequest;

final class UserCreateRequest extends FormRequest
{
    public function authorize(): bool
    {
        return true;
    }

    /**
     * @return array<string, array<int, string>>
     */
    public function rules(): array
    {
        return [
            'email' => ['required', 'email', 'max:255', 'unique:users'],
            'age' => ['required', 'integer', 'min:1', 'max:150'],
            'role' => ['required', 'string', 'in:admin,user,guest'],
            'password' => ['required', 'string', 'min:8', 'confirmed'],
        ];
    }

    /**
     * @return array<string, string>
     */
    public function messages(): array
    {
        return [
            'email.unique' => 'This email is already registered.',
            'role.in' => 'Role must be admin, user, or guest.',
        ];
    }
}
```

### Symfony Validation
```php
<?php

declare(strict_types=1);

namespace App\Dto;

use Symfony\Component\Validator\Constraints as Assert;

final readonly class UserCreate
{
    public function __construct(
        #[Assert\NotBlank]
        #[Assert\Email]
        #[Assert\Length(max: 255)]
        public string $email,

        #[Assert\Positive]
        #[Assert\LessThan(150)]
        public int $age,

        #[Assert\Choice(choices: ['admin', 'user', 'guest'])]
        public string $role,
    ) {}
}
```

---

## Testing

### Frameworks
- **PHPUnit**: Standard testing framework
- **Pest**: Modern, elegant syntax (Laravel-friendly)
- **Mockery**: Mocking library
- **Laravel Dusk**: Browser testing

### Guardrails
- ✓ Test files: `*Test.php` in `tests/` directory
- ✓ Test methods: `test_*` or `@test` annotation
- ✓ Use descriptive names: `test_user_creation_fails_with_invalid_email()`
- ✓ Use data providers for parameterized tests
- ✓ Mock external services
- ✓ Coverage target: >80% for business logic
- ✓ Use database transactions for integration tests

### Example (PHPUnit)
```php
<?php

declare(strict_types=1);

namespace Tests\Unit;

use App\Services\UserService;
use App\Repositories\UserRepository;
use App\Exceptions\ValidationException;
use PHPUnit\Framework\TestCase;
use PHPUnit\Framework\Attributes\Test;
use PHPUnit\Framework\Attributes\DataProvider;
use Mockery;

final class UserServiceTest extends TestCase
{
    private UserService $service;
    private UserRepository $repository;

    protected function setUp(): void
    {
        parent::setUp();
        $this->repository = Mockery::mock(UserRepository::class);
        $this->service = new UserService($this->repository);
    }

    protected function tearDown(): void
    {
        Mockery::close();
        parent::tearDown();
    }

    #[Test]
    public function creates_user_with_valid_data(): void
    {
        $data = [
            'email' => 'test@example.com',
            'age' => 25,
            'role' => 'user',
        ];

        $this->repository
            ->shouldReceive('create')
            ->once()
            ->with($data)
            ->andReturn(new User(1, ...$data));

        $user = $this->service->create($data);

        $this->assertSame('test@example.com', $user->email);
        $this->assertSame(25, $user->age);
    }

    #[Test]
    #[DataProvider('invalidEmailProvider')]
    public function throws_exception_for_invalid_email(string $email): void
    {
        $this->expectException(ValidationException::class);
        $this->expectExceptionMessage('Invalid email');

        $this->service->create([
            'email' => $email,
            'age' => 25,
            'role' => 'user',
        ]);
    }

    /**
     * @return array<string, array{string}>
     */
    public static function invalidEmailProvider(): array
    {
        return [
            'empty' => [''],
            'no at sign' => ['invalid'],
            'no domain' => ['test@'],
        ];
    }
}
```

### Example (Pest - Laravel)
```php
<?php

declare(strict_types=1);

use App\Models\User;
use function Pest\Laravel\{postJson, assertDatabaseHas};

test('user can be created with valid data', function () {
    $response = postJson('/api/users', [
        'email' => 'test@example.com',
        'age' => 25,
        'role' => 'user',
        'password' => 'password123',
        'password_confirmation' => 'password123',
    ]);

    $response
        ->assertCreated()
        ->assertJsonPath('email', 'test@example.com');

    assertDatabaseHas('users', ['email' => 'test@example.com']);
});

test('user creation fails with invalid email', function (string $email) {
    $response = postJson('/api/users', [
        'email' => $email,
        'age' => 25,
        'role' => 'user',
    ]);

    $response->assertUnprocessable()
        ->assertJsonValidationErrors(['email']);
})->with([
    'empty' => '',
    'invalid format' => 'not-an-email',
    'missing domain' => 'test@',
]);
```

---

## Tooling

### Essential Tools
- **PHP-CS-Fixer**: Code formatting
- **PHPStan** / **Psalm**: Static analysis
- **PHPUnit**: Testing
- **Composer**: Dependency management
- **Laravel Pint**: Laravel code style (wrapper for PHP-CS-Fixer)

### Configuration
```php
<?php
// .php-cs-fixer.php

$finder = PhpCsFixer\Finder::create()
    ->in(__DIR__ . '/app')
    ->in(__DIR__ . '/tests');

return (new PhpCsFixer\Config())
    ->setRules([
        '@PSR12' => true,
        '@PHP81Migration' => true,
        'strict_param' => true,
        'declare_strict_types' => true,
        'array_syntax' => ['syntax' => 'short'],
        'ordered_imports' => ['sort_algorithm' => 'alpha'],
        'no_unused_imports' => true,
        'final_class' => true,
        'void_return' => true,
    ])
    ->setRiskyAllowed(true)
    ->setFinder($finder);
```

```yaml
# phpstan.neon
parameters:
    level: 8
    paths:
        - app
        - tests
    checkMissingIterableValueType: true
    checkGenericClassInNonGenericObjectType: true
```

```json
// composer.json scripts
{
    "scripts": {
        "lint": "php-cs-fixer fix --dry-run --diff",
        "fix": "php-cs-fixer fix",
        "analyse": "phpstan analyse",
        "test": "phpunit",
        "test:coverage": "phpunit --coverage-html coverage"
    }
}
```

### Pre-Commit Commands
```bash
# Format
composer fix
# Or Laravel
./vendor/bin/pint

# Static analysis
composer analyse

# Test
composer test

# Full check
composer lint && composer analyse && composer test
```

---

## Common Pitfalls

### Don't Do This
```php
<?php
// No strict types
// Missing type declarations
function process($data) {
    return $data;
}

// Using @ to suppress errors
$result = @file_get_contents($url);

// SQL injection
$query = "SELECT * FROM users WHERE id = " . $_GET['id'];

// XSS vulnerability
echo $_GET['name'];

// Loose comparison
if ($value == '0') { } // true for '', 0, null, false, []

// Not validating input
$user = User::create($_POST);
```

### Do This Instead
```php
<?php

declare(strict_types=1);

// Proper type declarations
function process(array $data): array
{
    return $data;
}

// Proper error handling
$result = file_get_contents($url);
if ($result === false) {
    throw new RuntimeException('Failed to fetch URL');
}

// Prepared statements
$stmt = $pdo->prepare('SELECT * FROM users WHERE id = ?');
$stmt->execute([$id]);

// Escaped output
echo htmlspecialchars($name, ENT_QUOTES, 'UTF-8');
// Or in Blade: {{ $name }}

// Strict comparison
if ($value === '0') { }

// Validated input
$validated = $request->validated();
$user = User::create($validated);
```

---

## Framework-Specific Patterns

### Laravel Controller
```php
<?php

declare(strict_types=1);

namespace App\Http\Controllers;

use App\Http\Requests\UserCreateRequest;
use App\Http\Resources\UserResource;
use App\Models\User;
use App\Services\UserService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Resources\Json\AnonymousResourceCollection;

final class UserController extends Controller
{
    public function __construct(
        private readonly UserService $userService,
    ) {}

    public function index(): AnonymousResourceCollection
    {
        $users = User::query()
            ->with('profile')
            ->paginate(15);

        return UserResource::collection($users);
    }

    public function store(UserCreateRequest $request): JsonResponse
    {
        $user = $this->userService->create($request->validated());

        return UserResource::make($user)
            ->response()
            ->setStatusCode(201);
    }

    public function show(User $user): UserResource
    {
        return UserResource::make($user->load('profile'));
    }

    public function destroy(User $user): JsonResponse
    {
        $user->delete();

        return response()->json(null, 204);
    }
}
```

### Laravel Service
```php
<?php

declare(strict_types=1);

namespace App\Services;

use App\Models\User;
use App\Events\UserCreated;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Hash;

final class UserService
{
    /**
     * @param array{email: string, age: int, role: string, password: string} $data
     */
    public function create(array $data): User
    {
        return DB::transaction(function () use ($data): User {
            $user = User::create([
                'email' => $data['email'],
                'age' => $data['age'],
                'role' => $data['role'],
                'password' => Hash::make($data['password']),
            ]);

            event(new UserCreated($user));

            return $user;
        });
    }
}
```

### Symfony Controller
```php
<?php

declare(strict_types=1);

namespace App\Controller;

use App\Dto\UserCreate;
use App\Entity\User;
use App\Repository\UserRepository;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\Attribute\MapRequestPayload;
use Symfony\Component\Routing\Attribute\Route;

#[Route('/api/users')]
final class UserController extends AbstractController
{
    public function __construct(
        private readonly UserRepository $userRepository,
    ) {}

    #[Route('', methods: ['GET'])]
    public function index(): JsonResponse
    {
        $users = $this->userRepository->findAll();

        return $this->json($users);
    }

    #[Route('', methods: ['POST'])]
    public function create(#[MapRequestPayload] UserCreate $dto): JsonResponse
    {
        $user = new User(
            email: $dto->email,
            age: $dto->age,
            role: $dto->role,
        );

        $this->userRepository->save($user, flush: true);

        return $this->json($user, Response::HTTP_CREATED);
    }

    #[Route('/{id}', methods: ['GET'])]
    public function show(User $user): JsonResponse
    {
        return $this->json($user);
    }
}
```

### Enums (PHP 8.1+)
```php
<?php

declare(strict_types=1);

namespace App\Enums;

enum UserRole: string
{
    case Admin = 'admin';
    case User = 'user';
    case Guest = 'guest';

    public function label(): string
    {
        return match ($this) {
            self::Admin => 'Administrator',
            self::User => 'Regular User',
            self::Guest => 'Guest User',
        };
    }

    public function permissions(): array
    {
        return match ($this) {
            self::Admin => ['read', 'write', 'delete', 'admin'],
            self::User => ['read', 'write'],
            self::Guest => ['read'],
        };
    }
}

// Usage
$role = UserRole::from('admin');
echo $role->label(); // "Administrator"
```

---

## Performance Considerations

### Optimization Guardrails
- ✓ Use eager loading to avoid N+1 queries
- ✓ Use database indexes for frequently queried columns
- ✓ Use caching (Redis, Memcached) for expensive operations
- ✓ Use queues for long-running tasks
- ✓ Use pagination for large datasets
- ✓ Profile with Xdebug, Blackfire, or Laravel Telescope
- ✓ Enable OPcache in production

### Example
```php
<?php

// Eager loading
$users = User::with(['profile', 'orders'])->get();

// Query caching (Laravel)
$users = Cache::remember('users.active', 3600, function () {
    return User::where('active', true)->get();
});

// Chunking for large datasets
User::query()
    ->where('created_at', '<', now()->subYear())
    ->chunkById(1000, function ($users) {
        foreach ($users as $user) {
            // Process user
        }
    });

// Queue job
dispatch(new SendWelcomeEmail($user));
```

---

## Security Best Practices

### Guardrails
- ✓ Always use `declare(strict_types=1)`
- ✓ Use prepared statements (Eloquent/Doctrine do this automatically)
- ✓ Escape output in templates
- ✓ Use `password_hash()` with `PASSWORD_DEFAULT` or `PASSWORD_ARGON2ID`
- ✓ Validate and sanitize all user input
- ✓ Use CSRF protection
- ✓ Set secure session settings
- ✓ Keep dependencies updated (`composer audit`)
- ✓ Use HTTPS in production

### Example
```php
<?php

declare(strict_types=1);

// Password hashing
$hash = password_hash($password, PASSWORD_ARGON2ID, [
    'memory_cost' => 65536,
    'time_cost' => 4,
    'threads' => 3,
]);

// Password verification
if (password_verify($input, $hash)) {
    // Valid password
}

// Session security (php.ini or runtime)
ini_set('session.cookie_httponly', '1');
ini_set('session.cookie_secure', '1');
ini_set('session.cookie_samesite', 'Strict');
ini_set('session.use_strict_mode', '1');

// Input sanitization
$email = filter_var($input, FILTER_SANITIZE_EMAIL);
$int = filter_var($input, FILTER_VALIDATE_INT);

// XSS prevention
$safe = htmlspecialchars($userInput, ENT_QUOTES | ENT_HTML5, 'UTF-8');
```

---

## References

- [PHP Documentation](https://www.php.net/docs.php)
- [PHP-FIG PSR Standards](https://www.php-fig.org/psr/)
- [Laravel Documentation](https://laravel.com/docs)
- [Symfony Documentation](https://symfony.com/doc/current/index.html)
- [PHPUnit Documentation](https://phpunit.de/documentation.html)
- [PHPStan Documentation](https://phpstan.org/user-guide/getting-started)
- [PHP: The Right Way](https://phptherightway.com/)
