# Laravel Framework Guide

> **Framework**: Laravel 10+
> **Language**: PHP 8.1+
> **Type**: Full-Stack Web Framework
> **Use Cases**: Web Applications, APIs, E-commerce, SaaS

---

## Overview

Laravel is a PHP web application framework with expressive, elegant syntax. It provides structure and tools for building modern web applications with features like routing, authentication, queuing, and an ORM (Eloquent).

---

## Project Structure

```
myapp/
├── app/
│   ├── Console/              # Artisan commands
│   │   └── Commands/
│   ├── Exceptions/           # Exception handlers
│   ├── Http/
│   │   ├── Controllers/      # Request handlers
│   │   ├── Middleware/       # HTTP middleware
│   │   └── Requests/         # Form request validation
│   ├── Models/               # Eloquent models
│   ├── Providers/            # Service providers
│   └── Services/             # Business logic (custom)
├── bootstrap/
│   └── app.php
├── config/                   # Configuration files
│   ├── app.php
│   ├── database.php
│   └── ...
├── database/
│   ├── factories/            # Model factories
│   ├── migrations/           # Database migrations
│   └── seeders/             # Database seeders
├── public/
│   └── index.php            # Entry point
├── resources/
│   ├── css/
│   ├── js/
│   └── views/               # Blade templates
├── routes/
│   ├── api.php              # API routes
│   ├── web.php              # Web routes
│   └── console.php          # Console routes
├── storage/
│   ├── app/
│   ├── framework/
│   └── logs/
├── tests/
│   ├── Feature/
│   └── Unit/
├── .env                     # Environment variables
├── .env.example
├── artisan                  # CLI tool
├── composer.json
└── phpunit.xml
```

---

## Dependencies

### composer.json
```json
{
    "name": "mycompany/myapp",
    "type": "project",
    "require": {
        "php": "^8.1",
        "laravel/framework": "^10.0",
        "laravel/sanctum": "^3.2",
        "laravel/tinker": "^2.8"
    },
    "require-dev": {
        "fakerphp/faker": "^1.9.1",
        "laravel/pint": "^1.0",
        "laravel/sail": "^1.18",
        "mockery/mockery": "^1.4.4",
        "nunomaduro/collision": "^7.0",
        "phpunit/phpunit": "^10.1",
        "spatie/laravel-ignition": "^2.0"
    },
    "autoload": {
        "psr-4": {
            "App\\": "app/",
            "Database\\Factories\\": "database/factories/",
            "Database\\Seeders\\": "database/seeders/"
        }
    },
    "autoload-dev": {
        "psr-4": {
            "Tests\\": "tests/"
        }
    },
    "scripts": {
        "post-autoload-dump": [
            "Illuminate\\Foundation\\ComposerScripts::postAutoloadDump",
            "@php artisan package:discover --ansi"
        ],
        "test": "php artisan test",
        "lint": "vendor/bin/pint --test",
        "format": "vendor/bin/pint"
    }
}
```

---

## Core Patterns

### Model (Eloquent ORM)
```php
<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasMany;
use Illuminate\Database\Eloquent\SoftDeletes;

class User extends Model
{
    use HasFactory, SoftDeletes;

    /**
     * The attributes that are mass assignable.
     *
     * @var array<int, string>
     */
    protected $fillable = [
        'name',
        'email',
        'password',
        'role',
    ];

    /**
     * The attributes that should be hidden for serialization.
     *
     * @var array<int, string>
     */
    protected $hidden = [
        'password',
        'remember_token',
    ];

    /**
     * The attributes that should be cast.
     *
     * @var array<string, string>
     */
    protected $casts = [
        'email_verified_at' => 'datetime',
        'password' => 'hashed',
        'is_active' => 'boolean',
    ];

    // Relationships
    public function posts(): HasMany
    {
        return $this->hasMany(Post::class);
    }

    public function organization(): BelongsTo
    {
        return $this->belongsTo(Organization::class);
    }

    // Scopes
    public function scopeActive($query)
    {
        return $query->where('is_active', true);
    }

    public function scopeWithRole($query, string $role)
    {
        return $query->where('role', $role);
    }

    // Accessors
    public function getFullNameAttribute(): string
    {
        return "{$this->first_name} {$this->last_name}";
    }

    // Mutators
    public function setEmailAttribute(string $value): void
    {
        $this->attributes['email'] = strtolower($value);
    }
}
```

### Controller (Resource Controller)
```php
<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Http\Requests\StoreUserRequest;
use App\Http\Requests\UpdateUserRequest;
use App\Http\Resources\UserResource;
use App\Models\User;
use App\Services\UserService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Resources\Json\AnonymousResourceCollection;

class UserController extends Controller
{
    public function __construct(
        private readonly UserService $userService,
    ) {}

    /**
     * Display a listing of users.
     */
    public function index(): AnonymousResourceCollection
    {
        $users = User::query()
            ->with('organization')
            ->active()
            ->paginate(15);

        return UserResource::collection($users);
    }

    /**
     * Store a newly created user.
     */
    public function store(StoreUserRequest $request): JsonResponse
    {
        $user = $this->userService->create($request->validated());

        return UserResource::make($user)
            ->response()
            ->setStatusCode(201);
    }

    /**
     * Display the specified user.
     */
    public function show(User $user): UserResource
    {
        return UserResource::make($user->load('posts', 'organization'));
    }

    /**
     * Update the specified user.
     */
    public function update(UpdateUserRequest $request, User $user): UserResource
    {
        $user = $this->userService->update($user, $request->validated());

        return UserResource::make($user);
    }

    /**
     * Remove the specified user.
     */
    public function destroy(User $user): JsonResponse
    {
        $this->userService->delete($user);

        return response()->json(null, 204);
    }
}
```

### Form Request Validation
```php
<?php

declare(strict_types=1);

namespace App\Http\Requests;

use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;
use Illuminate\Validation\Rules\Password;

class StoreUserRequest extends FormRequest
{
    /**
     * Determine if the user is authorized to make this request.
     */
    public function authorize(): bool
    {
        return true;
    }

    /**
     * Get the validation rules that apply to the request.
     *
     * @return array<string, \Illuminate\Contracts\Validation\ValidationRule|array|string>
     */
    public function rules(): array
    {
        return [
            'name' => ['required', 'string', 'max:255'],
            'email' => [
                'required',
                'string',
                'email',
                'max:255',
                Rule::unique('users', 'email'),
            ],
            'password' => [
                'required',
                'confirmed',
                Password::min(8)
                    ->letters()
                    ->mixedCase()
                    ->numbers()
                    ->symbols(),
            ],
            'role' => ['required', 'string', Rule::in(['admin', 'user', 'guest'])],
            'organization_id' => ['nullable', 'exists:organizations,id'],
        ];
    }

    /**
     * Get custom messages for validator errors.
     *
     * @return array<string, string>
     */
    public function messages(): array
    {
        return [
            'email.unique' => 'This email address is already registered.',
            'role.in' => 'Role must be admin, user, or guest.',
        ];
    }
}
```

### Service Layer
```php
<?php

declare(strict_types=1);

namespace App\Services;

use App\Events\UserCreated;
use App\Models\User;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Hash;

class UserService
{
    /**
     * Create a new user.
     *
     * @param array<string, mixed> $data
     */
    public function create(array $data): User
    {
        return DB::transaction(function () use ($data): User {
            $user = User::create([
                'name' => $data['name'],
                'email' => $data['email'],
                'password' => Hash::make($data['password']),
                'role' => $data['role'],
                'organization_id' => $data['organization_id'] ?? null,
            ]);

            event(new UserCreated($user));

            return $user;
        });
    }

    /**
     * Update an existing user.
     *
     * @param array<string, mixed> $data
     */
    public function update(User $user, array $data): User
    {
        $updateData = collect($data)
            ->except(['password'])
            ->toArray();

        if (isset($data['password'])) {
            $updateData['password'] = Hash::make($data['password']);
        }

        $user->update($updateData);

        return $user->fresh();
    }

    /**
     * Delete a user.
     */
    public function delete(User $user): bool
    {
        return $user->delete();
    }
}
```

### API Resource
```php
<?php

declare(strict_types=1);

namespace App\Http\Resources;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class UserResource extends JsonResource
{
    /**
     * Transform the resource into an array.
     *
     * @return array<string, mixed>
     */
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'email' => $this->email,
            'role' => $this->role,
            'is_active' => $this->is_active,
            'created_at' => $this->created_at->toIso8601String(),
            'updated_at' => $this->updated_at->toIso8601String(),

            // Conditional relationships
            'organization' => OrganizationResource::make($this->whenLoaded('organization')),
            'posts' => PostResource::collection($this->whenLoaded('posts')),

            // Conditional attributes
            'posts_count' => $this->when(
                $this->posts_count !== null,
                $this->posts_count
            ),
        ];
    }
}
```

### Migration
```php
<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('users', function (Blueprint $table) {
            $table->id();
            $table->foreignId('organization_id')
                ->nullable()
                ->constrained()
                ->nullOnDelete();
            $table->string('name');
            $table->string('email')->unique();
            $table->timestamp('email_verified_at')->nullable();
            $table->string('password');
            $table->string('role')->default('user');
            $table->boolean('is_active')->default(true);
            $table->rememberToken();
            $table->timestamps();
            $table->softDeletes();

            $table->index(['role', 'is_active']);
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('users');
    }
};
```

### Routes
```php
<?php
// routes/api.php

declare(strict_types=1);

use App\Http\Controllers\Api\AuthController;
use App\Http\Controllers\Api\UserController;
use Illuminate\Support\Facades\Route;

Route::prefix('v1')->group(function () {
    // Public routes
    Route::post('/login', [AuthController::class, 'login']);
    Route::post('/register', [AuthController::class, 'register']);

    // Protected routes
    Route::middleware('auth:sanctum')->group(function () {
        Route::post('/logout', [AuthController::class, 'logout']);
        Route::get('/me', [AuthController::class, 'me']);

        Route::apiResource('users', UserController::class);

        // Admin-only routes
        Route::middleware('can:admin')->group(function () {
            Route::delete('/users/{user}/force', [UserController::class, 'forceDelete']);
        });
    });
});
```

### Middleware
```php
<?php

declare(strict_types=1);

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class EnsureUserIsAdmin
{
    /**
     * Handle an incoming request.
     */
    public function handle(Request $request, Closure $next): Response
    {
        if ($request->user()?->role !== 'admin') {
            abort(403, 'Access denied. Admin role required.');
        }

        return $next($request);
    }
}
```

### Event and Listener
```php
<?php
// app/Events/UserCreated.php

declare(strict_types=1);

namespace App\Events;

use App\Models\User;
use Illuminate\Broadcasting\InteractsWithSockets;
use Illuminate\Foundation\Events\Dispatchable;
use Illuminate\Queue\SerializesModels;

class UserCreated
{
    use Dispatchable, InteractsWithSockets, SerializesModels;

    public function __construct(
        public readonly User $user,
    ) {}
}
```

```php
<?php
// app/Listeners/SendWelcomeEmail.php

declare(strict_types=1);

namespace App\Listeners;

use App\Events\UserCreated;
use App\Mail\WelcomeEmail;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Support\Facades\Mail;

class SendWelcomeEmail implements ShouldQueue
{
    public function handle(UserCreated $event): void
    {
        Mail::to($event->user->email)
            ->send(new WelcomeEmail($event->user));
    }
}
```

### Job (Queue)
```php
<?php

declare(strict_types=1);

namespace App\Jobs;

use App\Models\User;
use App\Services\ReportService;
use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldBeUnique;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;

class GenerateUserReport implements ShouldQueue, ShouldBeUnique
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    public int $tries = 3;
    public int $backoff = 60;

    public function __construct(
        public readonly User $user,
        public readonly string $reportType,
    ) {}

    public function handle(ReportService $reportService): void
    {
        $reportService->generate($this->user, $this->reportType);
    }

    public function uniqueId(): string
    {
        return $this->user->id . '-' . $this->reportType;
    }

    public function failed(\Throwable $exception): void
    {
        // Handle job failure
        logger()->error('Report generation failed', [
            'user_id' => $this->user->id,
            'error' => $exception->getMessage(),
        ]);
    }
}
```

---

## Testing

### Feature Test
```php
<?php

declare(strict_types=1);

namespace Tests\Feature;

use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class UserControllerTest extends TestCase
{
    use RefreshDatabase;

    public function test_can_list_users(): void
    {
        $admin = User::factory()->create(['role' => 'admin']);
        User::factory()->count(5)->create();

        $response = $this->actingAs($admin)
            ->getJson('/api/v1/users');

        $response->assertOk()
            ->assertJsonCount(6, 'data')
            ->assertJsonStructure([
                'data' => [
                    '*' => ['id', 'name', 'email', 'role'],
                ],
                'meta' => ['current_page', 'total'],
            ]);
    }

    public function test_can_create_user(): void
    {
        $admin = User::factory()->create(['role' => 'admin']);

        $response = $this->actingAs($admin)
            ->postJson('/api/v1/users', [
                'name' => 'John Doe',
                'email' => 'john@example.com',
                'password' => 'Password123!',
                'password_confirmation' => 'Password123!',
                'role' => 'user',
            ]);

        $response->assertCreated()
            ->assertJsonPath('data.email', 'john@example.com');

        $this->assertDatabaseHas('users', [
            'email' => 'john@example.com',
        ]);
    }

    public function test_cannot_create_user_with_invalid_email(): void
    {
        $admin = User::factory()->create(['role' => 'admin']);

        $response = $this->actingAs($admin)
            ->postJson('/api/v1/users', [
                'name' => 'John Doe',
                'email' => 'invalid-email',
                'password' => 'Password123!',
                'password_confirmation' => 'Password123!',
                'role' => 'user',
            ]);

        $response->assertUnprocessable()
            ->assertJsonValidationErrors(['email']);
    }

    public function test_guest_cannot_access_users(): void
    {
        $response = $this->getJson('/api/v1/users');

        $response->assertUnauthorized();
    }
}
```

### Unit Test
```php
<?php

declare(strict_types=1);

namespace Tests\Unit;

use App\Models\User;
use App\Services\UserService;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class UserServiceTest extends TestCase
{
    use RefreshDatabase;

    private UserService $service;

    protected function setUp(): void
    {
        parent::setUp();
        $this->service = new UserService();
    }

    public function test_creates_user_with_hashed_password(): void
    {
        $user = $this->service->create([
            'name' => 'John Doe',
            'email' => 'john@example.com',
            'password' => 'plain-password',
            'role' => 'user',
        ]);

        $this->assertNotEquals('plain-password', $user->password);
        $this->assertTrue(password_verify('plain-password', $user->password));
    }

    public function test_updates_user_without_changing_password(): void
    {
        $user = User::factory()->create();
        $originalPassword = $user->password;

        $this->service->update($user, [
            'name' => 'Updated Name',
        ]);

        $this->assertEquals('Updated Name', $user->fresh()->name);
        $this->assertEquals($originalPassword, $user->fresh()->password);
    }
}
```

### Model Factory
```php
<?php

declare(strict_types=1);

namespace Database\Factories;

use App\Models\Organization;
use Illuminate\Database\Eloquent\Factories\Factory;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Str;

/**
 * @extends \Illuminate\Database\Eloquent\Factories\Factory<\App\Models\User>
 */
class UserFactory extends Factory
{
    public function definition(): array
    {
        return [
            'name' => fake()->name(),
            'email' => fake()->unique()->safeEmail(),
            'email_verified_at' => now(),
            'password' => Hash::make('password'),
            'role' => 'user',
            'is_active' => true,
            'remember_token' => Str::random(10),
        ];
    }

    public function admin(): static
    {
        return $this->state(fn (array $attributes) => [
            'role' => 'admin',
        ]);
    }

    public function inactive(): static
    {
        return $this->state(fn (array $attributes) => [
            'is_active' => false,
        ]);
    }

    public function withOrganization(): static
    {
        return $this->state(fn (array $attributes) => [
            'organization_id' => Organization::factory(),
        ]);
    }
}
```

---

## Commands

```bash
# Installation
composer create-project laravel/laravel myapp
cd myapp

# Development server
php artisan serve

# Artisan commands
php artisan make:model User -mfc      # Model with migration, factory, controller
php artisan make:controller Api/UserController --api
php artisan make:request StoreUserRequest
php artisan make:resource UserResource
php artisan make:middleware EnsureAdmin
php artisan make:event UserCreated
php artisan make:listener SendWelcomeEmail
php artisan make:job GenerateReport
php artisan make:mail WelcomeEmail
php artisan make:test UserControllerTest

# Database
php artisan migrate
php artisan migrate:fresh --seed
php artisan db:seed
php artisan make:migration create_posts_table

# Cache
php artisan cache:clear
php artisan config:cache
php artisan route:cache
php artisan view:cache

# Queue
php artisan queue:work
php artisan queue:failed

# Testing
php artisan test
php artisan test --filter=UserControllerTest
php artisan test --coverage

# Code quality
./vendor/bin/pint           # Format
./vendor/bin/pint --test    # Check formatting
./vendor/bin/phpstan        # Static analysis (if installed)

# Production
php artisan optimize
php artisan config:cache
php artisan route:cache
php artisan view:cache
```

---

## Best Practices

### Do
- ✓ Use Form Requests for validation
- ✓ Use API Resources for response transformation
- ✓ Use Service classes for business logic
- ✓ Use Jobs for long-running tasks
- ✓ Use Events/Listeners for decoupling
- ✓ Use Eloquent relationships with eager loading
- ✓ Use database transactions for multi-step operations
- ✓ Use model factories for testing
- ✓ Use Laravel Pint for code formatting

### Don't
- ✗ Don't put business logic in controllers
- ✗ Don't use `DB::raw()` without escaping
- ✗ Don't store sensitive data in `.env` committed to git
- ✗ Don't use `*` in Eloquent select queries in production
- ✗ Don't forget to add indexes for frequently queried columns
- ✗ Don't ignore N+1 query problems (use `with()`)
- ✗ Don't use synchronous operations for emails/notifications

---

## Framework Comparison

| Feature | Laravel | Symfony | WordPress |
|---------|---------|---------|-----------|
| Learning Curve | Moderate | Steep | Easy |
| Performance | Good | Excellent | Moderate |
| ORM | Eloquent | Doctrine | wpdb |
| Templating | Blade | Twig | PHP |
| Best For | Full-stack apps | Enterprise | CMS/Blogs |
| CLI Tool | Artisan | Console | WP-CLI |
| Auth | Built-in | Security | Built-in |

---

## References

- [Laravel Documentation](https://laravel.com/docs)
- [Laracasts](https://laracasts.com/)
- [Laravel News](https://laravel-news.com/)
- [Laravel Pint](https://laravel.com/docs/pint)
- [Laravel Sanctum](https://laravel.com/docs/sanctum)
- [Laravel Horizon](https://laravel.com/docs/horizon)
