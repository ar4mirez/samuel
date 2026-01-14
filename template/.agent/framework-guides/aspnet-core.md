# ASP.NET Core Framework Guide

> **Framework**: ASP.NET Core 8.x (LTS)
> **Language**: C# 12
> **Type**: Web API, MVC, Minimal APIs
> **Use Cases**: Enterprise APIs, Microservices, Full-stack web apps

---

## Overview

ASP.NET Core is a cross-platform, high-performance framework for building modern, cloud-enabled web applications. It supports both traditional MVC and minimal API patterns with excellent tooling and enterprise-grade features.

---

## Project Structure

```
MyApp/
├── MyApp.Api/                    # Web API project
│   ├── Controllers/              # API controllers
│   │   └── UsersController.cs
│   ├── Endpoints/                # Minimal API endpoints (alternative)
│   │   └── UserEndpoints.cs
│   ├── Middleware/               # Custom middleware
│   │   └── ExceptionMiddleware.cs
│   ├── Filters/                  # Action filters
│   │   └── ValidationFilter.cs
│   ├── Extensions/               # Service extensions
│   │   └── ServiceCollectionExtensions.cs
│   ├── Program.cs                # Application entry point
│   ├── appsettings.json          # Configuration
│   ├── appsettings.Development.json
│   └── MyApp.Api.csproj
├── MyApp.Core/                   # Domain/business logic
│   ├── Entities/
│   │   └── User.cs
│   ├── Interfaces/
│   │   ├── IUserRepository.cs
│   │   └── IUserService.cs
│   ├── Services/
│   │   └── UserService.cs
│   ├── Exceptions/
│   │   └── DomainException.cs
│   └── MyApp.Core.csproj
├── MyApp.Infrastructure/         # Data access, external services
│   ├── Data/
│   │   ├── AppDbContext.cs
│   │   └── Configurations/
│   │       └── UserConfiguration.cs
│   ├── Repositories/
│   │   └── UserRepository.cs
│   └── MyApp.Infrastructure.csproj
├── MyApp.Contracts/              # DTOs, API contracts
│   ├── Requests/
│   │   └── UserRequest.cs
│   ├── Responses/
│   │   └── UserResponse.cs
│   └── MyApp.Contracts.csproj
├── tests/
│   ├── MyApp.UnitTests/
│   │   └── Services/
│   │       └── UserServiceTests.cs
│   └── MyApp.IntegrationTests/
│       └── Controllers/
│           └── UsersControllerTests.cs
├── MyApp.sln
├── Directory.Build.props
├── Directory.Packages.props      # Central package management
└── docker-compose.yml
```

---

## Project Configuration

### Directory.Build.props
```xml
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

### Directory.Packages.props
```xml
<Project>
  <PropertyGroup>
    <ManagePackageVersionsCentrally>true</ManagePackageVersionsCentrally>
  </PropertyGroup>

  <ItemGroup>
    <!-- ASP.NET Core -->
    <PackageVersion Include="Microsoft.AspNetCore.Authentication.JwtBearer" Version="8.0.0" />
    <PackageVersion Include="Microsoft.AspNetCore.OpenApi" Version="8.0.0" />
    <PackageVersion Include="Swashbuckle.AspNetCore" Version="6.5.0" />

    <!-- Entity Framework -->
    <PackageVersion Include="Microsoft.EntityFrameworkCore" Version="8.0.0" />
    <PackageVersion Include="Microsoft.EntityFrameworkCore.Design" Version="8.0.0" />
    <PackageVersion Include="Npgsql.EntityFrameworkCore.PostgreSQL" Version="8.0.0" />

    <!-- Validation & Mapping -->
    <PackageVersion Include="FluentValidation.AspNetCore" Version="11.3.0" />
    <PackageVersion Include="Mapster" Version="7.4.0" />
    <PackageVersion Include="Mapster.DependencyInjection" Version="1.0.1" />

    <!-- Logging & Monitoring -->
    <PackageVersion Include="Serilog.AspNetCore" Version="8.0.0" />
    <PackageVersion Include="Serilog.Sinks.Console" Version="5.0.1" />

    <!-- Testing -->
    <PackageVersion Include="xunit" Version="2.6.2" />
    <PackageVersion Include="xunit.runner.visualstudio" Version="2.5.4" />
    <PackageVersion Include="Moq" Version="4.20.70" />
    <PackageVersion Include="FluentAssertions" Version="6.12.0" />
    <PackageVersion Include="Microsoft.AspNetCore.Mvc.Testing" Version="8.0.0" />
    <PackageVersion Include="Testcontainers.PostgreSql" Version="3.6.0" />
    <PackageVersion Include="coverlet.collector" Version="6.0.0" />
  </ItemGroup>
</Project>
```

### MyApp.Api.csproj
```xml
<Project Sdk="Microsoft.NET.Sdk.Web">
  <ItemGroup>
    <PackageReference Include="Microsoft.AspNetCore.Authentication.JwtBearer" />
    <PackageReference Include="Microsoft.AspNetCore.OpenApi" />
    <PackageReference Include="Swashbuckle.AspNetCore" />
    <PackageReference Include="FluentValidation.AspNetCore" />
    <PackageReference Include="Mapster" />
    <PackageReference Include="Mapster.DependencyInjection" />
    <PackageReference Include="Serilog.AspNetCore" />
    <PackageReference Include="Serilog.Sinks.Console" />
  </ItemGroup>

  <ItemGroup>
    <ProjectReference Include="..\MyApp.Core\MyApp.Core.csproj" />
    <ProjectReference Include="..\MyApp.Infrastructure\MyApp.Infrastructure.csproj" />
    <ProjectReference Include="..\MyApp.Contracts\MyApp.Contracts.csproj" />
  </ItemGroup>
</Project>
```

---

## Application Configuration

### appsettings.json
```json
{
  "ConnectionStrings": {
    "DefaultConnection": "Host=localhost;Database=myapp;Username=postgres;Password=postgres"
  },
  "Jwt": {
    "Secret": "your-secret-key-at-least-32-characters-long",
    "Issuer": "myapp",
    "Audience": "myapp-users",
    "ExpirationMinutes": 60
  },
  "Logging": {
    "LogLevel": {
      "Default": "Information",
      "Microsoft.AspNetCore": "Warning",
      "Microsoft.EntityFrameworkCore": "Warning"
    }
  },
  "AllowedHosts": "*"
}
```

### Program.cs
```csharp
using FluentValidation;
using Mapster;
using MapsterMapper;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.EntityFrameworkCore;
using Microsoft.IdentityModel.Tokens;
using Microsoft.OpenApi.Models;
using MyApp.Api.Middleware;
using MyApp.Core.Interfaces;
using MyApp.Core.Services;
using MyApp.Infrastructure.Data;
using MyApp.Infrastructure.Repositories;
using Serilog;
using System.Reflection;
using System.Text;

var builder = WebApplication.CreateBuilder(args);

// Serilog
builder.Host.UseSerilog((context, config) =>
    config.ReadFrom.Configuration(context.Configuration));

// DbContext
builder.Services.AddDbContext<AppDbContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")));

// Repositories
builder.Services.AddScoped<IUserRepository, UserRepository>();

// Services
builder.Services.AddScoped<IUserService, UserService>();

// FluentValidation
builder.Services.AddValidatorsFromAssemblyContaining<Program>();

// Mapster
var config = TypeAdapterConfig.GlobalSettings;
config.Scan(Assembly.GetExecutingAssembly());
builder.Services.AddSingleton(config);
builder.Services.AddScoped<IMapper, ServiceMapper>();

// Authentication
builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
    .AddJwtBearer(options =>
    {
        options.TokenValidationParameters = new TokenValidationParameters
        {
            ValidateIssuer = true,
            ValidateAudience = true,
            ValidateLifetime = true,
            ValidateIssuerSigningKey = true,
            ValidIssuer = builder.Configuration["Jwt:Issuer"],
            ValidAudience = builder.Configuration["Jwt:Audience"],
            IssuerSigningKey = new SymmetricSecurityKey(
                Encoding.UTF8.GetBytes(builder.Configuration["Jwt:Secret"]!))
        };
    });

builder.Services.AddAuthorization();

// Controllers
builder.Services.AddControllers();

// Swagger
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen(options =>
{
    options.SwaggerDoc("v1", new OpenApiInfo
    {
        Title = "MyApp API",
        Version = "v1"
    });

    options.AddSecurityDefinition("Bearer", new OpenApiSecurityScheme
    {
        Description = "JWT Authorization header using the Bearer scheme",
        Name = "Authorization",
        In = ParameterLocation.Header,
        Type = SecuritySchemeType.Http,
        Scheme = "bearer"
    });

    options.AddSecurityRequirement(new OpenApiSecurityRequirement
    {
        {
            new OpenApiSecurityScheme
            {
                Reference = new OpenApiReference
                {
                    Type = ReferenceType.SecurityScheme,
                    Id = "Bearer"
                }
            },
            Array.Empty<string>()
        }
    });
});

// Health checks
builder.Services.AddHealthChecks()
    .AddDbContextCheck<AppDbContext>();

var app = builder.Build();

// Middleware pipeline
app.UseMiddleware<ExceptionMiddleware>();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseSerilogRequestLogging();

app.UseAuthentication();
app.UseAuthorization();

app.MapControllers();
app.MapHealthChecks("/health");

app.Run();

// Make Program class accessible for testing
public partial class Program { }
```

---

## Domain Layer

### Entities/User.cs
```csharp
namespace MyApp.Core.Entities;

public class User
{
    public long Id { get; set; }
    public required string Email { get; set; }
    public required string PasswordHash { get; set; }
    public required string FirstName { get; set; }
    public required string LastName { get; set; }
    public UserRole Role { get; set; } = UserRole.User;
    public bool Active { get; set; } = true;
    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
    public DateTime? UpdatedAt { get; set; }

    public string FullName => $"{FirstName} {LastName}";
}

public enum UserRole
{
    User,
    Admin
}
```

### Exceptions/DomainException.cs
```csharp
namespace MyApp.Core.Exceptions;

public abstract class DomainException : Exception
{
    public abstract int StatusCode { get; }

    protected DomainException(string message) : base(message) { }
}

public class NotFoundException : DomainException
{
    public override int StatusCode => 404;

    public NotFoundException(string entity, object id)
        : base($"{entity} with ID {id} was not found") { }
}

public class ConflictException : DomainException
{
    public override int StatusCode => 409;

    public ConflictException(string entity, string field, object value)
        : base($"{entity} with {field} '{value}' already exists") { }
}

public class ValidationException : DomainException
{
    public override int StatusCode => 400;
    public IDictionary<string, string[]> Errors { get; }

    public ValidationException(IDictionary<string, string[]> errors)
        : base("Validation failed")
    {
        Errors = errors;
    }
}
```

---

## Contracts Layer

### Requests/UserRequest.cs
```csharp
namespace MyApp.Contracts.Requests;

public record CreateUserRequest(
    string Email,
    string Password,
    string FirstName,
    string LastName
);

public record UpdateUserRequest(
    string? FirstName,
    string? LastName,
    bool? Active
);
```

### Responses/UserResponse.cs
```csharp
namespace MyApp.Contracts.Responses;

public record UserResponse(
    long Id,
    string Email,
    string FirstName,
    string LastName,
    string Role,
    bool Active,
    DateTime CreatedAt
);

public record PagedResponse<T>(
    IReadOnlyList<T> Items,
    int Page,
    int PageSize,
    long TotalCount,
    int TotalPages
);
```

---

## Validation

### Validators/CreateUserRequestValidator.cs
```csharp
using FluentValidation;
using MyApp.Contracts.Requests;

namespace MyApp.Api.Validators;

public class CreateUserRequestValidator : AbstractValidator<CreateUserRequest>
{
    public CreateUserRequestValidator()
    {
        RuleFor(x => x.Email)
            .NotEmpty().WithMessage("Email is required")
            .EmailAddress().WithMessage("Invalid email format")
            .MaximumLength(255).WithMessage("Email must not exceed 255 characters");

        RuleFor(x => x.Password)
            .NotEmpty().WithMessage("Password is required")
            .MinimumLength(8).WithMessage("Password must be at least 8 characters")
            .Matches("[A-Z]").WithMessage("Password must contain uppercase letter")
            .Matches("[a-z]").WithMessage("Password must contain lowercase letter")
            .Matches("[0-9]").WithMessage("Password must contain a digit")
            .Matches("[^a-zA-Z0-9]").WithMessage("Password must contain special character");

        RuleFor(x => x.FirstName)
            .NotEmpty().WithMessage("First name is required")
            .MaximumLength(100).WithMessage("First name must not exceed 100 characters");

        RuleFor(x => x.LastName)
            .NotEmpty().WithMessage("Last name is required")
            .MaximumLength(100).WithMessage("Last name must not exceed 100 characters");
    }
}
```

---

## Mapping

### Mappings/UserMappingConfig.cs
```csharp
using Mapster;
using MyApp.Contracts.Requests;
using MyApp.Contracts.Responses;
using MyApp.Core.Entities;

namespace MyApp.Api.Mappings;

public class UserMappingConfig : IRegister
{
    public void Register(TypeAdapterConfig config)
    {
        config.NewConfig<CreateUserRequest, User>()
            .Ignore(dest => dest.Id)
            .Ignore(dest => dest.PasswordHash)
            .Ignore(dest => dest.CreatedAt)
            .Ignore(dest => dest.UpdatedAt);

        config.NewConfig<User, UserResponse>()
            .Map(dest => dest.Role, src => src.Role.ToString());
    }
}
```

---

## Infrastructure Layer

### Data/AppDbContext.cs
```csharp
using Microsoft.EntityFrameworkCore;
using MyApp.Core.Entities;

namespace MyApp.Infrastructure.Data;

public class AppDbContext : DbContext
{
    public AppDbContext(DbContextOptions<AppDbContext> options) : base(options) { }

    public DbSet<User> Users => Set<User>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.ApplyConfigurationsFromAssembly(typeof(AppDbContext).Assembly);
    }

    public override Task<int> SaveChangesAsync(CancellationToken cancellationToken = default)
    {
        foreach (var entry in ChangeTracker.Entries<User>())
        {
            if (entry.State == EntityState.Modified)
            {
                entry.Entity.UpdatedAt = DateTime.UtcNow;
            }
        }

        return base.SaveChangesAsync(cancellationToken);
    }
}
```

### Data/Configurations/UserConfiguration.cs
```csharp
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using MyApp.Core.Entities;

namespace MyApp.Infrastructure.Data.Configurations;

public class UserConfiguration : IEntityTypeConfiguration<User>
{
    public void Configure(EntityTypeBuilder<User> builder)
    {
        builder.ToTable("users");

        builder.HasKey(u => u.Id);

        builder.Property(u => u.Id)
            .HasColumnName("id");

        builder.Property(u => u.Email)
            .HasColumnName("email")
            .HasMaxLength(255)
            .IsRequired();

        builder.Property(u => u.PasswordHash)
            .HasColumnName("password_hash")
            .HasMaxLength(255)
            .IsRequired();

        builder.Property(u => u.FirstName)
            .HasColumnName("first_name")
            .HasMaxLength(100)
            .IsRequired();

        builder.Property(u => u.LastName)
            .HasColumnName("last_name")
            .HasMaxLength(100)
            .IsRequired();

        builder.Property(u => u.Role)
            .HasColumnName("role")
            .HasConversion<string>()
            .HasMaxLength(50);

        builder.Property(u => u.Active)
            .HasColumnName("active")
            .HasDefaultValue(true);

        builder.Property(u => u.CreatedAt)
            .HasColumnName("created_at");

        builder.Property(u => u.UpdatedAt)
            .HasColumnName("updated_at");

        builder.HasIndex(u => u.Email)
            .IsUnique();

        builder.Ignore(u => u.FullName);
    }
}
```

### Repositories/UserRepository.cs
```csharp
using Microsoft.EntityFrameworkCore;
using MyApp.Core.Entities;
using MyApp.Core.Interfaces;
using MyApp.Infrastructure.Data;

namespace MyApp.Infrastructure.Repositories;

public class UserRepository : IUserRepository
{
    private readonly AppDbContext _context;

    public UserRepository(AppDbContext context)
    {
        _context = context;
    }

    public async Task<User?> GetByIdAsync(long id, CancellationToken ct = default)
    {
        return await _context.Users.FindAsync([id], ct);
    }

    public async Task<User?> GetByEmailAsync(string email, CancellationToken ct = default)
    {
        return await _context.Users
            .FirstOrDefaultAsync(u => u.Email == email, ct);
    }

    public async Task<(IReadOnlyList<User> Items, long TotalCount)> GetAllAsync(
        int page, int pageSize, CancellationToken ct = default)
    {
        var query = _context.Users.AsNoTracking();

        var totalCount = await query.LongCountAsync(ct);

        var items = await query
            .OrderByDescending(u => u.CreatedAt)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(ct);

        return (items, totalCount);
    }

    public async Task<bool> ExistsByEmailAsync(string email, CancellationToken ct = default)
    {
        return await _context.Users.AnyAsync(u => u.Email == email, ct);
    }

    public async Task<User> CreateAsync(User user, CancellationToken ct = default)
    {
        _context.Users.Add(user);
        await _context.SaveChangesAsync(ct);
        return user;
    }

    public async Task<User> UpdateAsync(User user, CancellationToken ct = default)
    {
        _context.Users.Update(user);
        await _context.SaveChangesAsync(ct);
        return user;
    }

    public async Task DeleteAsync(User user, CancellationToken ct = default)
    {
        _context.Users.Remove(user);
        await _context.SaveChangesAsync(ct);
    }
}
```

---

## Service Layer

### Interfaces/IUserService.cs
```csharp
using MyApp.Contracts.Requests;
using MyApp.Contracts.Responses;

namespace MyApp.Core.Interfaces;

public interface IUserService
{
    Task<UserResponse> GetByIdAsync(long id, CancellationToken ct = default);
    Task<PagedResponse<UserResponse>> GetAllAsync(int page, int pageSize, CancellationToken ct = default);
    Task<UserResponse> CreateAsync(CreateUserRequest request, CancellationToken ct = default);
    Task<UserResponse> UpdateAsync(long id, UpdateUserRequest request, CancellationToken ct = default);
    Task DeleteAsync(long id, CancellationToken ct = default);
}
```

### Services/UserService.cs
```csharp
using MapsterMapper;
using Microsoft.Extensions.Logging;
using MyApp.Contracts.Requests;
using MyApp.Contracts.Responses;
using MyApp.Core.Entities;
using MyApp.Core.Exceptions;
using MyApp.Core.Interfaces;
using BC = BCrypt.Net.BCrypt;

namespace MyApp.Core.Services;

public class UserService : IUserService
{
    private readonly IUserRepository _userRepository;
    private readonly IMapper _mapper;
    private readonly ILogger<UserService> _logger;

    public UserService(
        IUserRepository userRepository,
        IMapper mapper,
        ILogger<UserService> logger)
    {
        _userRepository = userRepository;
        _mapper = mapper;
        _logger = logger;
    }

    public async Task<UserResponse> GetByIdAsync(long id, CancellationToken ct = default)
    {
        var user = await _userRepository.GetByIdAsync(id, ct)
            ?? throw new NotFoundException(nameof(User), id);

        return _mapper.Map<UserResponse>(user);
    }

    public async Task<PagedResponse<UserResponse>> GetAllAsync(
        int page, int pageSize, CancellationToken ct = default)
    {
        var (items, totalCount) = await _userRepository.GetAllAsync(page, pageSize, ct);

        var responses = _mapper.Map<List<UserResponse>>(items);
        var totalPages = (int)Math.Ceiling(totalCount / (double)pageSize);

        return new PagedResponse<UserResponse>(responses, page, pageSize, totalCount, totalPages);
    }

    public async Task<UserResponse> CreateAsync(CreateUserRequest request, CancellationToken ct = default)
    {
        if (await _userRepository.ExistsByEmailAsync(request.Email, ct))
        {
            throw new ConflictException(nameof(User), "email", request.Email);
        }

        var user = _mapper.Map<User>(request);
        user.PasswordHash = BC.HashPassword(request.Password);

        var created = await _userRepository.CreateAsync(user, ct);

        _logger.LogInformation("User created with ID {UserId}", created.Id);

        return _mapper.Map<UserResponse>(created);
    }

    public async Task<UserResponse> UpdateAsync(
        long id, UpdateUserRequest request, CancellationToken ct = default)
    {
        var user = await _userRepository.GetByIdAsync(id, ct)
            ?? throw new NotFoundException(nameof(User), id);

        if (request.FirstName is not null)
            user.FirstName = request.FirstName;

        if (request.LastName is not null)
            user.LastName = request.LastName;

        if (request.Active.HasValue)
            user.Active = request.Active.Value;

        var updated = await _userRepository.UpdateAsync(user, ct);

        _logger.LogInformation("User updated with ID {UserId}", updated.Id);

        return _mapper.Map<UserResponse>(updated);
    }

    public async Task DeleteAsync(long id, CancellationToken ct = default)
    {
        var user = await _userRepository.GetByIdAsync(id, ct)
            ?? throw new NotFoundException(nameof(User), id);

        await _userRepository.DeleteAsync(user, ct);

        _logger.LogInformation("User deleted with ID {UserId}", id);
    }
}
```

---

## API Controller

### Controllers/UsersController.cs
```csharp
using FluentValidation;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using MyApp.Contracts.Requests;
using MyApp.Contracts.Responses;
using MyApp.Core.Interfaces;

namespace MyApp.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
[Produces("application/json")]
public class UsersController : ControllerBase
{
    private readonly IUserService _userService;
    private readonly IValidator<CreateUserRequest> _createValidator;

    public UsersController(
        IUserService userService,
        IValidator<CreateUserRequest> createValidator)
    {
        _userService = userService;
        _createValidator = createValidator;
    }

    /// <summary>
    /// Get all users with pagination
    /// </summary>
    [HttpGet]
    [Authorize]
    [ProducesResponseType(typeof(PagedResponse<UserResponse>), StatusCodes.Status200OK)]
    public async Task<ActionResult<PagedResponse<UserResponse>>> GetAll(
        [FromQuery] int page = 1,
        [FromQuery] int pageSize = 10,
        CancellationToken ct = default)
    {
        var result = await _userService.GetAllAsync(page, pageSize, ct);
        return Ok(result);
    }

    /// <summary>
    /// Get user by ID
    /// </summary>
    [HttpGet("{id:long}")]
    [Authorize]
    [ProducesResponseType(typeof(UserResponse), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<ActionResult<UserResponse>> GetById(long id, CancellationToken ct)
    {
        var user = await _userService.GetByIdAsync(id, ct);
        return Ok(user);
    }

    /// <summary>
    /// Create a new user
    /// </summary>
    [HttpPost]
    [ProducesResponseType(typeof(UserResponse), StatusCodes.Status201Created)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status409Conflict)]
    public async Task<ActionResult<UserResponse>> Create(
        [FromBody] CreateUserRequest request,
        CancellationToken ct)
    {
        var validation = await _createValidator.ValidateAsync(request, ct);
        if (!validation.IsValid)
        {
            return BadRequest(validation.Errors);
        }

        var user = await _userService.CreateAsync(request, ct);
        return CreatedAtAction(nameof(GetById), new { id = user.Id }, user);
    }

    /// <summary>
    /// Update an existing user
    /// </summary>
    [HttpPut("{id:long}")]
    [Authorize]
    [ProducesResponseType(typeof(UserResponse), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<ActionResult<UserResponse>> Update(
        long id,
        [FromBody] UpdateUserRequest request,
        CancellationToken ct)
    {
        var user = await _userService.UpdateAsync(id, request, ct);
        return Ok(user);
    }

    /// <summary>
    /// Delete a user
    /// </summary>
    [HttpDelete("{id:long}")]
    [Authorize(Roles = "Admin")]
    [ProducesResponseType(StatusCodes.Status204NoContent)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<IActionResult> Delete(long id, CancellationToken ct)
    {
        await _userService.DeleteAsync(id, ct);
        return NoContent();
    }
}
```

---

## Minimal API Alternative

### Endpoints/UserEndpoints.cs
```csharp
using FluentValidation;
using MyApp.Contracts.Requests;
using MyApp.Contracts.Responses;
using MyApp.Core.Interfaces;

namespace MyApp.Api.Endpoints;

public static class UserEndpoints
{
    public static IEndpointRouteBuilder MapUserEndpoints(this IEndpointRouteBuilder routes)
    {
        var group = routes.MapGroup("/api/users")
            .WithTags("Users")
            .WithOpenApi();

        group.MapGet("/", GetAll)
            .RequireAuthorization()
            .Produces<PagedResponse<UserResponse>>();

        group.MapGet("/{id:long}", GetById)
            .RequireAuthorization()
            .Produces<UserResponse>()
            .Produces(StatusCodes.Status404NotFound);

        group.MapPost("/", Create)
            .Produces<UserResponse>(StatusCodes.Status201Created)
            .Produces(StatusCodes.Status400BadRequest);

        group.MapPut("/{id:long}", Update)
            .RequireAuthorization()
            .Produces<UserResponse>()
            .Produces(StatusCodes.Status404NotFound);

        group.MapDelete("/{id:long}", Delete)
            .RequireAuthorization("Admin")
            .Produces(StatusCodes.Status204NoContent);

        return routes;
    }

    private static async Task<IResult> GetAll(
        IUserService service,
        int page = 1,
        int pageSize = 10,
        CancellationToken ct = default)
    {
        var result = await service.GetAllAsync(page, pageSize, ct);
        return Results.Ok(result);
    }

    private static async Task<IResult> GetById(
        long id,
        IUserService service,
        CancellationToken ct)
    {
        var user = await service.GetByIdAsync(id, ct);
        return Results.Ok(user);
    }

    private static async Task<IResult> Create(
        CreateUserRequest request,
        IUserService service,
        IValidator<CreateUserRequest> validator,
        CancellationToken ct)
    {
        var validation = await validator.ValidateAsync(request, ct);
        if (!validation.IsValid)
        {
            return Results.BadRequest(validation.Errors);
        }

        var user = await service.CreateAsync(request, ct);
        return Results.Created($"/api/users/{user.Id}", user);
    }

    private static async Task<IResult> Update(
        long id,
        UpdateUserRequest request,
        IUserService service,
        CancellationToken ct)
    {
        var user = await service.UpdateAsync(id, request, ct);
        return Results.Ok(user);
    }

    private static async Task<IResult> Delete(
        long id,
        IUserService service,
        CancellationToken ct)
    {
        await service.DeleteAsync(id, ct);
        return Results.NoContent();
    }
}
```

---

## Exception Handling Middleware

### Middleware/ExceptionMiddleware.cs
```csharp
using System.Text.Json;
using MyApp.Core.Exceptions;

namespace MyApp.Api.Middleware;

public class ExceptionMiddleware
{
    private readonly RequestDelegate _next;
    private readonly ILogger<ExceptionMiddleware> _logger;

    public ExceptionMiddleware(RequestDelegate next, ILogger<ExceptionMiddleware> logger)
    {
        _next = next;
        _logger = logger;
    }

    public async Task InvokeAsync(HttpContext context)
    {
        try
        {
            await _next(context);
        }
        catch (Exception ex)
        {
            await HandleExceptionAsync(context, ex);
        }
    }

    private async Task HandleExceptionAsync(HttpContext context, Exception exception)
    {
        var (statusCode, response) = exception switch
        {
            NotFoundException ex => (ex.StatusCode, new ErrorResponse(ex.Message)),
            ConflictException ex => (ex.StatusCode, new ErrorResponse(ex.Message)),
            ValidationException ex => (ex.StatusCode, new ValidationErrorResponse(ex.Message, ex.Errors)),
            _ => (500, new ErrorResponse("An unexpected error occurred"))
        };

        if (statusCode == 500)
        {
            _logger.LogError(exception, "Unhandled exception occurred");
        }

        context.Response.ContentType = "application/json";
        context.Response.StatusCode = statusCode;

        var json = JsonSerializer.Serialize(response, new JsonSerializerOptions
        {
            PropertyNamingPolicy = JsonNamingPolicy.CamelCase
        });

        await context.Response.WriteAsync(json);
    }
}

public record ErrorResponse(string Message);

public record ValidationErrorResponse(
    string Message,
    IDictionary<string, string[]> Errors
);
```

---

## Database Migrations

### Create Migration
```bash
# Install EF tools
dotnet tool install --global dotnet-ef

# Add migration
dotnet ef migrations add InitialCreate -p MyApp.Infrastructure -s MyApp.Api

# Apply migration
dotnet ef database update -p MyApp.Infrastructure -s MyApp.Api
```

### Example Migration
```csharp
using Microsoft.EntityFrameworkCore.Migrations;
using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;

namespace MyApp.Infrastructure.Migrations;

public partial class InitialCreate : Migration
{
    protected override void Up(MigrationBuilder migrationBuilder)
    {
        migrationBuilder.CreateTable(
            name: "users",
            columns: table => new
            {
                id = table.Column<long>(nullable: false)
                    .Annotation("Npgsql:ValueGenerationStrategy",
                        NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                email = table.Column<string>(maxLength: 255, nullable: false),
                password_hash = table.Column<string>(maxLength: 255, nullable: false),
                first_name = table.Column<string>(maxLength: 100, nullable: false),
                last_name = table.Column<string>(maxLength: 100, nullable: false),
                role = table.Column<string>(maxLength: 50, nullable: false),
                active = table.Column<bool>(nullable: false, defaultValue: true),
                created_at = table.Column<DateTime>(nullable: false),
                updated_at = table.Column<DateTime>(nullable: true)
            },
            constraints: table =>
            {
                table.PrimaryKey("PK_users", x => x.id);
            });

        migrationBuilder.CreateIndex(
            name: "IX_users_email",
            table: "users",
            column: "email",
            unique: true);
    }

    protected override void Down(MigrationBuilder migrationBuilder)
    {
        migrationBuilder.DropTable(name: "users");
    }
}
```

---

## Testing

### Unit Tests
```csharp
using FluentAssertions;
using MapsterMapper;
using Microsoft.Extensions.Logging;
using Moq;
using MyApp.Contracts.Requests;
using MyApp.Core.Entities;
using MyApp.Core.Exceptions;
using MyApp.Core.Interfaces;
using MyApp.Core.Services;
using Xunit;

namespace MyApp.UnitTests.Services;

public class UserServiceTests
{
    private readonly Mock<IUserRepository> _repositoryMock;
    private readonly Mock<IMapper> _mapperMock;
    private readonly Mock<ILogger<UserService>> _loggerMock;
    private readonly UserService _sut;

    public UserServiceTests()
    {
        _repositoryMock = new Mock<IUserRepository>();
        _mapperMock = new Mock<IMapper>();
        _loggerMock = new Mock<ILogger<UserService>>();
        _sut = new UserService(
            _repositoryMock.Object,
            _mapperMock.Object,
            _loggerMock.Object);
    }

    [Fact]
    public async Task GetByIdAsync_WhenUserExists_ReturnsUser()
    {
        // Arrange
        var user = new User
        {
            Id = 1,
            Email = "test@example.com",
            PasswordHash = "hash",
            FirstName = "John",
            LastName = "Doe"
        };

        _repositoryMock.Setup(r => r.GetByIdAsync(1, It.IsAny<CancellationToken>()))
            .ReturnsAsync(user);

        _mapperMock.Setup(m => m.Map<UserResponse>(user))
            .Returns(new UserResponse(1, "test@example.com", "John", "Doe", "User", true, DateTime.UtcNow));

        // Act
        var result = await _sut.GetByIdAsync(1);

        // Assert
        result.Should().NotBeNull();
        result.Email.Should().Be("test@example.com");
    }

    [Fact]
    public async Task GetByIdAsync_WhenUserNotFound_ThrowsNotFoundException()
    {
        // Arrange
        _repositoryMock.Setup(r => r.GetByIdAsync(999, It.IsAny<CancellationToken>()))
            .ReturnsAsync((User?)null);

        // Act
        var act = () => _sut.GetByIdAsync(999);

        // Assert
        await act.Should().ThrowAsync<NotFoundException>();
    }

    [Fact]
    public async Task CreateAsync_WhenEmailExists_ThrowsConflictException()
    {
        // Arrange
        var request = new CreateUserRequest(
            "existing@example.com",
            "Password123!",
            "John",
            "Doe");

        _repositoryMock.Setup(r => r.ExistsByEmailAsync(request.Email, It.IsAny<CancellationToken>()))
            .ReturnsAsync(true);

        // Act
        var act = () => _sut.CreateAsync(request);

        // Assert
        await act.Should().ThrowAsync<ConflictException>();
    }
}
```

### Integration Tests
```csharp
using System.Net;
using System.Net.Http.Json;
using FluentAssertions;
using Microsoft.AspNetCore.Mvc.Testing;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using MyApp.Contracts.Requests;
using MyApp.Contracts.Responses;
using MyApp.Infrastructure.Data;
using Testcontainers.PostgreSql;
using Xunit;

namespace MyApp.IntegrationTests.Controllers;

public class UsersControllerTests : IAsyncLifetime
{
    private readonly PostgreSqlContainer _postgres = new PostgreSqlBuilder()
        .WithImage("postgres:15-alpine")
        .Build();

    private WebApplicationFactory<Program> _factory = null!;
    private HttpClient _client = null!;

    public async Task InitializeAsync()
    {
        await _postgres.StartAsync();

        _factory = new WebApplicationFactory<Program>()
            .WithWebHostBuilder(builder =>
            {
                builder.ConfigureServices(services =>
                {
                    var descriptor = services.SingleOrDefault(
                        d => d.ServiceType == typeof(DbContextOptions<AppDbContext>));

                    if (descriptor != null)
                        services.Remove(descriptor);

                    services.AddDbContext<AppDbContext>(options =>
                        options.UseNpgsql(_postgres.GetConnectionString()));
                });
            });

        _client = _factory.CreateClient();

        // Apply migrations
        using var scope = _factory.Services.CreateScope();
        var db = scope.ServiceProvider.GetRequiredService<AppDbContext>();
        await db.Database.MigrateAsync();
    }

    public async Task DisposeAsync()
    {
        await _factory.DisposeAsync();
        await _postgres.DisposeAsync();
    }

    [Fact]
    public async Task Create_WithValidRequest_ReturnsCreated()
    {
        // Arrange
        var request = new CreateUserRequest(
            "test@example.com",
            "Password123!",
            "John",
            "Doe");

        // Act
        var response = await _client.PostAsJsonAsync("/api/users", request);

        // Assert
        response.StatusCode.Should().Be(HttpStatusCode.Created);

        var user = await response.Content.ReadFromJsonAsync<UserResponse>();
        user.Should().NotBeNull();
        user!.Email.Should().Be("test@example.com");
        user.FirstName.Should().Be("John");
    }

    [Fact]
    public async Task Create_WithInvalidEmail_ReturnsBadRequest()
    {
        // Arrange
        var request = new CreateUserRequest(
            "invalid-email",
            "Password123!",
            "John",
            "Doe");

        // Act
        var response = await _client.PostAsJsonAsync("/api/users", request);

        // Assert
        response.StatusCode.Should().Be(HttpStatusCode.BadRequest);
    }
}
```

---

## Build & Run Commands

```bash
# Restore dependencies
dotnet restore

# Build
dotnet build

# Run
dotnet run --project MyApp.Api

# Run with watch
dotnet watch run --project MyApp.Api

# Run tests
dotnet test

# Run tests with coverage
dotnet test --collect:"XPlat Code Coverage"

# Format code
dotnet format

# Add migration
dotnet ef migrations add MigrationName -p MyApp.Infrastructure -s MyApp.Api

# Apply migrations
dotnet ef database update -p MyApp.Infrastructure -s MyApp.Api

# Publish
dotnet publish -c Release -o ./publish

# Docker build
docker build -t myapp:latest .
```

---

## Dockerfile

```dockerfile
# Build stage
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

COPY ["MyApp.Api/MyApp.Api.csproj", "MyApp.Api/"]
COPY ["MyApp.Core/MyApp.Core.csproj", "MyApp.Core/"]
COPY ["MyApp.Infrastructure/MyApp.Infrastructure.csproj", "MyApp.Infrastructure/"]
COPY ["MyApp.Contracts/MyApp.Contracts.csproj", "MyApp.Contracts/"]
COPY ["Directory.Build.props", "."]
COPY ["Directory.Packages.props", "."]

RUN dotnet restore "MyApp.Api/MyApp.Api.csproj"

COPY . .
RUN dotnet publish "MyApp.Api/MyApp.Api.csproj" -c Release -o /app/publish

# Runtime stage
FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS runtime
WORKDIR /app
EXPOSE 8080

COPY --from=build /app/publish .

# Non-root user
USER app

ENTRYPOINT ["dotnet", "MyApp.Api.dll"]
```

---

## Best Practices

### DO
- ✓ Use Central Package Management for version consistency
- ✓ Use `CancellationToken` for async operations
- ✓ Use records for DTOs (immutability)
- ✓ Use FluentValidation for complex validation
- ✓ Use Clean Architecture separation
- ✓ Use EF Core migrations for schema changes
- ✓ Use health checks for monitoring
- ✓ Use Serilog for structured logging
- ✓ Use Testcontainers for integration tests

### DON'T
- ✗ Don't expose entities in API responses
- ✗ Don't use synchronous database calls
- ✗ Don't catch exceptions without proper handling
- ✗ Don't hardcode connection strings
- ✗ Don't skip validation
- ✗ Don't use magic strings

---

## Framework Comparison

| Feature | ASP.NET Core | Spring Boot | Express |
|---------|-------------|-------------|---------|
| Language | C# | Java/Kotlin | JavaScript/TypeScript |
| DI | Built-in | Built-in | Manual/Third-party |
| ORM | Entity Framework | JPA/Hibernate | Prisma/TypeORM |
| Performance | Excellent | Very Good | Good |
| Learning Curve | Moderate | Steep | Gentle |
| Enterprise | Excellent | Excellent | Good |
| Microservices | Excellent | Excellent | Good |

---

## References

- [ASP.NET Core Documentation](https://docs.microsoft.com/aspnet/core)
- [Entity Framework Core](https://docs.microsoft.com/ef/core)
- [FluentValidation](https://docs.fluentvalidation.net/)
- [Mapster](https://github.com/MapsterMapper/Mapster)
- [Serilog](https://serilog.net/)
- [xUnit](https://xunit.net/)
- [FluentAssertions](https://fluentassertions.com/)
