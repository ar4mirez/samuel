# C# Guide

> **Applies to**: C# 11+, .NET 7+, ASP.NET Core, Unity, MAUI, Blazor

---

## Core Principles

1. **Type Safety**: Strong static typing with nullability annotations
2. **Modern C#**: Records, pattern matching, async/await, LINQ
3. **Nullable Reference Types**: Enable and treat warnings as errors
4. **Async by Default**: Use async/await for I/O operations
5. **Immutability**: Prefer immutable types, use `readonly` and records

---

## Language-Specific Guardrails

### .NET Version & Setup
- ✓ Use .NET 7+ (8 LTS recommended for new projects)
- ✓ Use C# 11+ language features
- ✓ Enable nullable reference types: `<Nullable>enable</Nullable>`
- ✓ Treat warnings as errors in CI: `<TreatWarningsAsErrors>true</TreatWarningsAsErrors>`
- ✓ Use Central Package Management for version consistency

### Code Style (EditorConfig)
- ✓ Use `.editorconfig` for consistent style
- ✓ Run `dotnet format` before every commit
- ✓ Follow Microsoft C# Coding Conventions
- ✓ Use `PascalCase` for public members, types, namespaces
- ✓ Use `camelCase` for local variables and parameters
- ✓ Use `_camelCase` for private fields
- ✓ Use `SCREAMING_SNAKE_CASE` for constants
- ✓ File-scoped namespaces (reduces indentation)

### Nullable Reference Types
- ✓ Enable nullable context project-wide
- ✓ Annotate nullability explicitly (`string?` vs `string`)
- ✓ Use null-conditional operator: `obj?.Property`
- ✓ Use null-coalescing operator: `value ?? defaultValue`
- ✓ Use null-coalescing assignment: `value ??= defaultValue`
- ✓ Validate parameters with `ArgumentNullException.ThrowIfNull()`

### Modern C# Features
- ✓ Use records for immutable data types
- ✓ Use `init` properties for immutable objects
- ✓ Use pattern matching in `switch` expressions
- ✓ Use `with` expressions for record copying
- ✓ Use raw string literals for multi-line strings
- ✓ Use primary constructors (C# 12) for simple classes
- ✓ Use collection expressions: `[1, 2, 3]`

### Async/Await
- ✓ Use `async`/`await` for I/O-bound operations
- ✓ Suffix async methods with `Async`: `GetUserAsync()`
- ✓ Return `Task` or `Task<T>` (not `void` except event handlers)
- ✓ Use `ConfigureAwait(false)` in library code
- ✓ Use `CancellationToken` for cancellable operations
- ✓ Never use `.Result` or `.Wait()` (causes deadlocks)
- ✓ Use `ValueTask<T>` for hot paths that often complete synchronously

### LINQ
- ✓ Prefer LINQ method syntax over query syntax
- ✓ Use meaningful variable names in lambdas
- ✓ Avoid side effects in LINQ expressions
- ✓ Use `ToList()` or `ToArray()` to materialize when needed
- ✓ Be aware of deferred execution

### Exception Handling
- ✓ Catch specific exception types
- ✓ Use `when` clause for filtered catches
- ✓ Throw with `throw;` to preserve stack trace
- ✓ Create custom exceptions for domain errors
- ✓ Use `ArgumentException`, `InvalidOperationException` appropriately
- ✓ Document exceptions with `<exception>` XML docs

---

## Project Structure

### Standard Layout
```
MySolution/
├── MySolution.sln
├── Directory.Build.props          # Shared build properties
├── Directory.Packages.props       # Central package management
├── src/
│   ├── MyProject.Api/            # Web API
│   │   ├── Controllers/
│   │   ├── Middleware/
│   │   └── Program.cs
│   ├── MyProject.Core/           # Domain/business logic
│   │   ├── Entities/
│   │   ├── Interfaces/
│   │   └── Services/
│   └── MyProject.Infrastructure/ # Data access, external services
│       ├── Data/
│       └── Repositories/
├── tests/
│   ├── MyProject.UnitTests/
│   └── MyProject.IntegrationTests/
└── README.md
```

### Guardrails
- ✓ Separate concerns into projects (Clean Architecture)
- ✓ Core/Domain has no dependencies on Infrastructure
- ✓ Use `internal` for implementation details
- ✓ One class per file (matching file name)

---

## Validation & Input Handling

### Recommended Libraries
- **FluentValidation**: Fluent validation rules
- **DataAnnotations**: Built-in attribute validation
- **System.ComponentModel.DataAnnotations**: Standard annotations

### Pattern (FluentValidation)
```csharp
using FluentValidation;

public record UserCreate(string Email, int Age, string Role);

public class UserCreateValidator : AbstractValidator<UserCreate>
{
    public UserCreateValidator()
    {
        RuleFor(x => x.Email)
            .NotEmpty()
            .EmailAddress()
            .MaximumLength(255);

        RuleFor(x => x.Age)
            .GreaterThan(0)
            .LessThan(150);

        RuleFor(x => x.Role)
            .NotEmpty()
            .Must(role => new[] { "admin", "user", "guest" }.Contains(role))
            .WithMessage("Role must be admin, user, or guest");
    }
}

// Usage in ASP.NET Core
public class UsersController : ControllerBase
{
    [HttpPost]
    public async Task<ActionResult<User>> Create(
        [FromBody] UserCreate request,
        [FromServices] IValidator<UserCreate> validator)
    {
        var result = await validator.ValidateAsync(request);
        if (!result.IsValid)
            return BadRequest(result.Errors);

        var user = await _userService.CreateAsync(request);
        return CreatedAtAction(nameof(Get), new { id = user.Id }, user);
    }
}
```

---

## Testing

### Frameworks
- **xUnit**: Modern, extensible (recommended)
- **NUnit**: Feature-rich, mature
- **MSTest**: Microsoft's framework
- **Moq** / **NSubstitute**: Mocking libraries
- **FluentAssertions**: Readable assertions
- **Testcontainers**: Integration testing

### Guardrails
- ✓ Test files: `*Tests.cs` (unit), `*IntegrationTests.cs`
- ✓ Use descriptive test names: `CreateUser_WithValidData_ReturnsUser()`
- ✓ Use `[Theory]` with `[InlineData]` for parameterized tests
- ✓ Use `Arrange-Act-Assert` pattern
- ✓ Mock external dependencies only
- ✓ Coverage target: >80% for business logic

### Example (xUnit + FluentAssertions + Moq)
```csharp
using FluentAssertions;
using Moq;
using Xunit;

public class UserServiceTests
{
    private readonly Mock<IUserRepository> _repositoryMock;
    private readonly UserService _sut;

    public UserServiceTests()
    {
        _repositoryMock = new Mock<IUserRepository>();
        _sut = new UserService(_repositoryMock.Object);
    }

    [Fact]
    public async Task CreateAsync_WithValidData_ReturnsUser()
    {
        // Arrange
        var request = new UserCreate("test@example.com", 25, "user");
        var expected = new User(1, "test@example.com", 25, "user");
        _repositoryMock
            .Setup(r => r.AddAsync(It.IsAny<User>()))
            .ReturnsAsync(expected);

        // Act
        var result = await _sut.CreateAsync(request);

        // Assert
        result.Should().NotBeNull();
        result.Email.Should().Be("test@example.com");
        result.Age.Should().Be(25);
        _repositoryMock.Verify(r => r.AddAsync(It.IsAny<User>()), Times.Once);
    }

    [Theory]
    [InlineData("")]
    [InlineData(" ")]
    [InlineData("invalid")]
    public async Task CreateAsync_WithInvalidEmail_ThrowsValidationException(string email)
    {
        // Arrange
        var request = new UserCreate(email, 25, "user");

        // Act
        var act = () => _sut.CreateAsync(request);

        // Assert
        await act.Should()
            .ThrowAsync<ValidationException>()
            .WithMessage("*email*");
    }
}
```

---

## Tooling

### Essential Tools
- **dotnet format**: Code formatting
- **Roslyn Analyzers**: Built-in code analysis
- **StyleCop.Analyzers**: Style enforcement
- **SonarAnalyzer.CSharp**: Code quality
- **coverlet**: Code coverage
- **dotnet-outdated**: Dependency updates

### Configuration
```xml
<!-- Directory.Build.props -->
<Project>
  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
    <Nullable>enable</Nullable>
    <ImplicitUsings>enable</ImplicitUsings>
    <TreatWarningsAsErrors>true</TreatWarningsAsErrors>
    <AnalysisLevel>latest-recommended</AnalysisLevel>
  </PropertyGroup>

  <ItemGroup>
    <PackageReference Include="StyleCop.Analyzers" Version="1.2.0-beta.556">
      <PrivateAssets>all</PrivateAssets>
      <IncludeAssets>runtime; build; native; contentfiles; analyzers</IncludeAssets>
    </PackageReference>
  </ItemGroup>
</Project>
```

```ini
# .editorconfig
root = true

[*.cs]
indent_style = space
indent_size = 4
charset = utf-8
trim_trailing_whitespace = true
insert_final_newline = true

# Naming conventions
dotnet_naming_rule.private_fields_should_be_camel_case.symbols = private_fields
dotnet_naming_rule.private_fields_should_be_camel_case.style = camel_case_underscore
dotnet_naming_rule.private_fields_should_be_camel_case.severity = suggestion

dotnet_naming_symbols.private_fields.applicable_kinds = field
dotnet_naming_symbols.private_fields.applicable_accessibilities = private

dotnet_naming_style.camel_case_underscore.capitalization = camel_case
dotnet_naming_style.camel_case_underscore.required_prefix = _

# Prefer expression body
csharp_style_expression_bodied_methods = when_on_single_line:suggestion
csharp_style_expression_bodied_properties = true:suggestion

# Prefer file-scoped namespaces
csharp_style_namespace_declarations = file_scoped:warning
```

### Pre-Commit Commands
```bash
# Format
dotnet format

# Build
dotnet build --no-restore

# Test
dotnet test --no-build

# Test with coverage
dotnet test --collect:"XPlat Code Coverage"

# Publish
dotnet publish -c Release
```

---

## Common Pitfalls

### Don't Do This
```csharp
// Ignoring nullable warnings
string? name = GetName();
Console.WriteLine(name.Length); // NullReferenceException

// Using .Result or .Wait()
var result = GetDataAsync().Result; // Deadlock risk

// Mutable DTOs
public class UserDto
{
    public string Email { get; set; } = "";
}

// Catching Exception
try { DoWork(); }
catch (Exception) { } // Swallowed

// Not disposing resources
var stream = File.OpenRead("file.txt");
// stream never disposed
```

### Do This Instead
```csharp
// Handle nullable properly
string? name = GetName();
if (name is not null)
{
    Console.WriteLine(name.Length);
}
// Or
Console.WriteLine(name?.Length ?? 0);

// Use async/await
var result = await GetDataAsync();

// Immutable records
public record UserDto(string Email);

// Catch specific, handle properly
try { DoWork(); }
catch (IOException ex)
{
    _logger.LogError(ex, "Failed to process file");
    throw;
}

// Using statement
await using var stream = File.OpenRead("file.txt");
```

---

## Framework-Specific Patterns

### ASP.NET Core Minimal API
```csharp
var builder = WebApplication.CreateBuilder(args);

// Services
builder.Services.AddScoped<IUserService, UserService>();
builder.Services.AddScoped<IValidator<UserCreate>, UserCreateValidator>();

var app = builder.Build();

// Endpoints
app.MapGet("/users/{id}", async (int id, IUserService service) =>
{
    var user = await service.GetByIdAsync(id);
    return user is not null ? Results.Ok(user) : Results.NotFound();
});

app.MapPost("/users", async (UserCreate request, IUserService service, IValidator<UserCreate> validator) =>
{
    var validation = await validator.ValidateAsync(request);
    if (!validation.IsValid)
        return Results.BadRequest(validation.Errors);

    var user = await service.CreateAsync(request);
    return Results.Created($"/users/{user.Id}", user);
});

app.Run();
```

### Records and Pattern Matching
```csharp
// Immutable record
public record User(int Id, string Email, int Age, string Role)
{
    public bool IsAdmin => Role == "admin";
}

// Record with validation
public record UserCreate(string Email, int Age, string Role)
{
    public UserCreate
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(Email);
        ArgumentOutOfRangeException.ThrowIfNegativeOrZero(Age);
    }
}

// Pattern matching
public string GetUserCategory(User user) => user switch
{
    { Age: < 18 } => "Minor",
    { Role: "admin" } => "Administrator",
    { Age: >= 65 } => "Senior",
    _ => "Regular"
};

// List patterns (C# 11)
int[] numbers = [1, 2, 3, 4, 5];
var result = numbers switch
{
    [1, 2, ..] => "Starts with 1, 2",
    [.., 4, 5] => "Ends with 4, 5",
    [] => "Empty",
    _ => "Other"
};
```

### Entity Framework Core
```csharp
public class AppDbContext : DbContext
{
    public DbSet<User> Users => Set<User>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<User>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Email).HasMaxLength(255).IsRequired();
            entity.HasIndex(e => e.Email).IsUnique();
        });
    }
}

// Repository pattern
public class UserRepository : IUserRepository
{
    private readonly AppDbContext _context;

    public UserRepository(AppDbContext context) => _context = context;

    public async Task<User?> GetByIdAsync(int id, CancellationToken ct = default)
        => await _context.Users.FindAsync([id], ct);

    public async Task<User> AddAsync(User user, CancellationToken ct = default)
    {
        _context.Users.Add(user);
        await _context.SaveChangesAsync(ct);
        return user;
    }

    public async Task<IReadOnlyList<User>> GetByRoleAsync(string role, CancellationToken ct = default)
        => await _context.Users
            .Where(u => u.Role == role)
            .AsNoTracking()
            .ToListAsync(ct);
}
```

---

## Performance Considerations

### Optimization Guardrails
- ✓ Use `Span<T>` and `Memory<T>` for high-performance scenarios
- ✓ Use `ArrayPool<T>.Shared` for temporary arrays
- ✓ Use `StringBuilder` for string concatenation in loops
- ✓ Use `AsNoTracking()` for read-only EF Core queries
- ✓ Use compiled queries for frequently executed EF queries
- ✓ Profile with BenchmarkDotNet, dotTrace, or PerfView
- ✓ Use `ValueTask<T>` for hot async paths

### Example
```csharp
// Span for zero-allocation parsing
public static int CountDigits(ReadOnlySpan<char> text)
{
    int count = 0;
    foreach (var c in text)
    {
        if (char.IsDigit(c)) count++;
    }
    return count;
}

// ArrayPool for temporary buffers
var buffer = ArrayPool<byte>.Shared.Rent(1024);
try
{
    // Use buffer
}
finally
{
    ArrayPool<byte>.Shared.Return(buffer);
}

// Compiled EF query
private static readonly Func<AppDbContext, string, Task<User?>> GetUserByEmail =
    EF.CompileAsyncQuery((AppDbContext ctx, string email) =>
        ctx.Users.FirstOrDefault(u => u.Email == email));
```

---

## Security Best Practices

### Guardrails
- ✓ Never hardcode secrets (use User Secrets, Key Vault, env vars)
- ✓ Use parameterized queries (EF Core does this automatically)
- ✓ Validate all inputs with FluentValidation or DataAnnotations
- ✓ Hash passwords with ASP.NET Core Identity or BCrypt
- ✓ Enable HTTPS redirection and HSTS
- ✓ Use anti-forgery tokens for forms
- ✓ Run `dotnet list package --vulnerable` regularly
- ✓ Use `[Authorize]` attributes appropriately

### Example
```csharp
// User secrets (development)
// dotnet user-secrets set "ApiKey" "secret-value"
var apiKey = builder.Configuration["ApiKey"];

// Password hashing (Identity)
var hasher = new PasswordHasher<User>();
var hashed = hasher.HashPassword(user, plainPassword);
var result = hasher.VerifyHashedPassword(user, hashed, providedPassword);

// Security headers
app.UseHttpsRedirection();
app.UseHsts();
app.Use(async (context, next) =>
{
    context.Response.Headers.Add("X-Content-Type-Options", "nosniff");
    context.Response.Headers.Add("X-Frame-Options", "DENY");
    await next();
});

// Authorization
[Authorize(Roles = "Admin")]
[HttpDelete("{id}")]
public async Task<IActionResult> Delete(int id) { ... }
```

---

## References

- [Microsoft C# Documentation](https://docs.microsoft.com/en-us/dotnet/csharp/)
- [C# Coding Conventions](https://docs.microsoft.com/en-us/dotnet/csharp/fundamentals/coding-style/coding-conventions)
- [ASP.NET Core Documentation](https://docs.microsoft.com/en-us/aspnet/core/)
- [Entity Framework Core](https://docs.microsoft.com/en-us/ef/core/)
- [FluentValidation](https://docs.fluentvalidation.net/)
- [xUnit Documentation](https://xunit.net/)
- [FluentAssertions](https://fluentassertions.com/)
