# Blazor Framework Guide

> **Framework**: Blazor (.NET 8)
> **Language**: C# 12
> **Type**: SPA Framework (WebAssembly & Server)
> **Use Cases**: Interactive web apps, dashboards, line-of-business apps

---

## Overview

Blazor is a .NET framework for building interactive web applications using C# instead of JavaScript. It supports two hosting models: Blazor WebAssembly (runs in browser) and Blazor Server (runs on server with SignalR). .NET 8 introduces Blazor United with hybrid rendering modes.

---

## Hosting Models

| Model | Description | Best For |
|-------|-------------|----------|
| **Blazor WebAssembly** | Runs entirely in browser via WebAssembly | Offline-capable apps, PWAs |
| **Blazor Server** | Runs on server, UI updates via SignalR | Low-latency, SEO-important |
| **Blazor United (.NET 8)** | Hybrid static + interactive rendering | Best of both worlds |

---

## Project Structure

```
MyApp/
├── MyApp.Client/                    # Blazor WebAssembly project
│   ├── Pages/                       # Routable components
│   │   ├── Home.razor
│   │   ├── Users/
│   │   │   ├── UserList.razor
│   │   │   ├── UserDetail.razor
│   │   │   └── UserForm.razor
│   │   └── Counter.razor
│   ├── Components/                  # Reusable components
│   │   ├── Layout/
│   │   │   ├── MainLayout.razor
│   │   │   └── NavMenu.razor
│   │   ├── Common/
│   │   │   ├── LoadingSpinner.razor
│   │   │   ├── ErrorBoundary.razor
│   │   │   └── Modal.razor
│   │   └── Users/
│   │       └── UserCard.razor
│   ├── Services/                    # Client-side services
│   │   ├── IUserService.cs
│   │   └── UserService.cs
│   ├── State/                       # State management
│   │   └── AppState.cs
│   ├── wwwroot/                     # Static assets
│   │   ├── css/
│   │   └── index.html               # WASM only
│   ├── _Imports.razor               # Global imports
│   ├── App.razor                    # Root component
│   ├── Program.cs
│   └── MyApp.Client.csproj
├── MyApp.Server/                    # API/Server host
│   ├── Controllers/
│   │   └── UsersController.cs
│   ├── Program.cs
│   └── MyApp.Server.csproj
├── MyApp.Shared/                    # Shared models/DTOs
│   ├── Models/
│   │   └── User.cs
│   ├── DTOs/
│   │   ├── UserRequest.cs
│   │   └── UserResponse.cs
│   └── MyApp.Shared.csproj
├── tests/
│   ├── MyApp.Client.Tests/
│   │   └── Components/
│   │       └── UserCardTests.cs
│   └── MyApp.Server.Tests/
└── MyApp.sln
```

---

## Project Configuration

### MyApp.Client.csproj (WebAssembly)
```xml
<Project Sdk="Microsoft.NET.Sdk.BlazorWebAssembly">

  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
    <Nullable>enable</Nullable>
    <ImplicitUsings>enable</ImplicitUsings>
  </PropertyGroup>

  <ItemGroup>
    <PackageReference Include="Microsoft.AspNetCore.Components.WebAssembly" Version="8.0.0" />
    <PackageReference Include="Microsoft.AspNetCore.Components.WebAssembly.Authentication" Version="8.0.0" />
    <PackageReference Include="Microsoft.Extensions.Http" Version="8.0.0" />
    <PackageReference Include="Blazored.LocalStorage" Version="4.4.0" />
    <PackageReference Include="Blazored.Toast" Version="4.1.0" />
  </ItemGroup>

  <ItemGroup>
    <ProjectReference Include="..\MyApp.Shared\MyApp.Shared.csproj" />
  </ItemGroup>

</Project>
```

### Program.cs (WebAssembly)
```csharp
using Microsoft.AspNetCore.Components.Web;
using Microsoft.AspNetCore.Components.WebAssembly.Hosting;
using MyApp.Client;
using MyApp.Client.Services;
using MyApp.Client.State;
using Blazored.LocalStorage;
using Blazored.Toast;

var builder = WebAssemblyHostBuilder.CreateDefault(args);
builder.RootComponents.Add<App>("#app");
builder.RootComponents.Add<HeadOutlet>("head::after");

// HTTP Client
builder.Services.AddScoped(sp =>
    new HttpClient { BaseAddress = new Uri(builder.HostEnvironment.BaseAddress) });

// Typed HTTP clients
builder.Services.AddHttpClient<IUserService, UserService>(client =>
    client.BaseAddress = new Uri(builder.HostEnvironment.BaseAddress));

// State management
builder.Services.AddScoped<AppState>();

// Third-party services
builder.Services.AddBlazoredLocalStorage();
builder.Services.AddBlazoredToast();

// Authentication
builder.Services.AddAuthorizationCore();

await builder.Build().RunAsync();
```

### _Imports.razor
```razor
@using System.Net.Http
@using System.Net.Http.Json
@using Microsoft.AspNetCore.Components.Forms
@using Microsoft.AspNetCore.Components.Routing
@using Microsoft.AspNetCore.Components.Web
@using Microsoft.AspNetCore.Components.Web.Virtualization
@using Microsoft.AspNetCore.Components.WebAssembly.Http
@using Microsoft.AspNetCore.Authorization
@using Microsoft.JSInterop
@using MyApp.Client
@using MyApp.Client.Components
@using MyApp.Client.Components.Layout
@using MyApp.Client.Components.Common
@using MyApp.Client.Services
@using MyApp.Client.State
@using MyApp.Shared.Models
@using MyApp.Shared.DTOs
```

---

## Components

### App.razor
```razor
<CascadingAuthenticationState>
    <Router AppAssembly="@typeof(App).Assembly">
        <Found Context="routeData">
            <AuthorizeRouteView RouteData="@routeData"
                                DefaultLayout="@typeof(MainLayout)">
                <NotAuthorized>
                    <RedirectToLogin />
                </NotAuthorized>
            </AuthorizeRouteView>
            <FocusOnNavigate RouteData="@routeData" Selector="h1" />
        </Found>
        <NotFound>
            <PageTitle>Not Found</PageTitle>
            <LayoutView Layout="@typeof(MainLayout)">
                <div class="alert alert-warning">
                    <h4>Page not found</h4>
                    <p>Sorry, the page you requested could not be found.</p>
                </div>
            </LayoutView>
        </NotFound>
    </Router>
</CascadingAuthenticationState>
```

### MainLayout.razor
```razor
@inherits LayoutComponentBase

<div class="page">
    <div class="sidebar">
        <NavMenu />
    </div>

    <main>
        <div class="top-row px-4">
            <AuthorizeView>
                <Authorized>
                    <span>Hello, @context.User.Identity?.Name!</span>
                    <a href="authentication/logout">Log out</a>
                </Authorized>
                <NotAuthorized>
                    <a href="authentication/login">Log in</a>
                </NotAuthorized>
            </AuthorizeView>
        </div>

        <article class="content px-4">
            <ErrorBoundary @ref="errorBoundary">
                <ChildContent>
                    @Body
                </ChildContent>
                <ErrorContent Context="ex">
                    <div class="alert alert-danger">
                        <h4>An error occurred</h4>
                        <p>@ex.Message</p>
                        <button class="btn btn-primary" @onclick="Recover">
                            Try again
                        </button>
                    </div>
                </ErrorContent>
            </ErrorBoundary>
        </article>
    </main>
</div>

@code {
    private ErrorBoundary? errorBoundary;

    private void Recover()
    {
        errorBoundary?.Recover();
    }
}
```

### NavMenu.razor
```razor
<div class="nav-menu">
    <nav class="flex-column">
        <div class="nav-item px-3">
            <NavLink class="nav-link" href="" Match="NavLinkMatch.All">
                <span class="oi oi-home" aria-hidden="true"></span> Home
            </NavLink>
        </div>
        <div class="nav-item px-3">
            <NavLink class="nav-link" href="users">
                <span class="oi oi-people" aria-hidden="true"></span> Users
            </NavLink>
        </div>
        <AuthorizeView Roles="Admin">
            <div class="nav-item px-3">
                <NavLink class="nav-link" href="admin">
                    <span class="oi oi-cog" aria-hidden="true"></span> Admin
                </NavLink>
            </div>
        </AuthorizeView>
    </nav>
</div>
```

---

## Pages

### Pages/Users/UserList.razor
```razor
@page "/users"
@attribute [Authorize]
@inject IUserService UserService
@inject NavigationManager Navigation

<PageTitle>Users</PageTitle>

<div class="d-flex justify-content-between align-items-center mb-4">
    <h1>Users</h1>
    <button class="btn btn-primary" @onclick="NavigateToCreate">
        <span class="oi oi-plus"></span> Add User
    </button>
</div>

@if (loading)
{
    <LoadingSpinner />
}
else if (error is not null)
{
    <div class="alert alert-danger">
        <h4>Error loading users</h4>
        <p>@error</p>
        <button class="btn btn-outline-danger" @onclick="LoadUsers">
            Retry
        </button>
    </div>
}
else if (users is null || users.Count == 0)
{
    <div class="alert alert-info">
        No users found. Click "Add User" to create one.
    </div>
}
else
{
    <div class="row">
        @foreach (var user in users)
        {
            <div class="col-md-4 mb-3">
                <UserCard User="user"
                          OnEdit="() => NavigateToEdit(user.Id)"
                          OnDelete="() => DeleteUser(user.Id)" />
            </div>
        }
    </div>

    <nav>
        <ul class="pagination">
            <li class="page-item @(currentPage == 1 ? "disabled" : "")">
                <button class="page-link" @onclick="PreviousPage">Previous</button>
            </li>
            @for (int i = 1; i <= totalPages; i++)
            {
                var pageNum = i;
                <li class="page-item @(currentPage == pageNum ? "active" : "")">
                    <button class="page-link" @onclick="() => GoToPage(pageNum)">
                        @pageNum
                    </button>
                </li>
            }
            <li class="page-item @(currentPage == totalPages ? "disabled" : "")">
                <button class="page-link" @onclick="NextPage">Next</button>
            </li>
        </ul>
    </nav>
}

@code {
    private List<UserResponse>? users;
    private bool loading = true;
    private string? error;
    private int currentPage = 1;
    private int totalPages = 1;
    private const int PageSize = 10;

    protected override async Task OnInitializedAsync()
    {
        await LoadUsers();
    }

    private async Task LoadUsers()
    {
        loading = true;
        error = null;

        try
        {
            var result = await UserService.GetUsersAsync(currentPage, PageSize);
            users = result.Items.ToList();
            totalPages = result.TotalPages;
        }
        catch (Exception ex)
        {
            error = ex.Message;
        }
        finally
        {
            loading = false;
        }
    }

    private void NavigateToCreate() => Navigation.NavigateTo("/users/new");

    private void NavigateToEdit(long id) => Navigation.NavigateTo($"/users/{id}/edit");

    private async Task DeleteUser(long id)
    {
        if (await UserService.DeleteUserAsync(id))
        {
            await LoadUsers();
        }
    }

    private async Task GoToPage(int page)
    {
        currentPage = page;
        await LoadUsers();
    }

    private async Task PreviousPage() => await GoToPage(currentPage - 1);

    private async Task NextPage() => await GoToPage(currentPage + 1);
}
```

### Pages/Users/UserForm.razor
```razor
@page "/users/new"
@page "/users/{Id:long}/edit"
@attribute [Authorize]
@inject IUserService UserService
@inject NavigationManager Navigation
@inject IToastService ToastService

<PageTitle>@(IsEdit ? "Edit User" : "New User")</PageTitle>

<h1>@(IsEdit ? "Edit User" : "New User")</h1>

@if (loading)
{
    <LoadingSpinner />
}
else
{
    <div class="row">
        <div class="col-md-6">
            <EditForm Model="model" OnValidSubmit="HandleSubmit">
                <DataAnnotationsValidator />
                <ValidationSummary class="alert alert-danger" />

                <div class="mb-3">
                    <label for="email" class="form-label">Email</label>
                    <InputText id="email" class="form-control" @bind-Value="model.Email"
                               disabled="@IsEdit" />
                    <ValidationMessage For="@(() => model.Email)" class="text-danger" />
                </div>

                @if (!IsEdit)
                {
                    <div class="mb-3">
                        <label for="password" class="form-label">Password</label>
                        <InputText id="password" type="password" class="form-control"
                                   @bind-Value="model.Password" />
                        <ValidationMessage For="@(() => model.Password)" class="text-danger" />
                    </div>
                }

                <div class="mb-3">
                    <label for="firstName" class="form-label">First Name</label>
                    <InputText id="firstName" class="form-control" @bind-Value="model.FirstName" />
                    <ValidationMessage For="@(() => model.FirstName)" class="text-danger" />
                </div>

                <div class="mb-3">
                    <label for="lastName" class="form-label">Last Name</label>
                    <InputText id="lastName" class="form-control" @bind-Value="model.LastName" />
                    <ValidationMessage For="@(() => model.LastName)" class="text-danger" />
                </div>

                @if (IsEdit)
                {
                    <div class="mb-3 form-check">
                        <InputCheckbox id="active" class="form-check-input" @bind-Value="model.Active" />
                        <label for="active" class="form-check-label">Active</label>
                    </div>
                }

                <div class="d-flex gap-2">
                    <button type="submit" class="btn btn-primary" disabled="@submitting">
                        @if (submitting)
                        {
                            <span class="spinner-border spinner-border-sm"></span>
                        }
                        @(IsEdit ? "Update" : "Create")
                    </button>
                    <button type="button" class="btn btn-secondary" @onclick="Cancel">
                        Cancel
                    </button>
                </div>
            </EditForm>
        </div>
    </div>
}

@code {
    [Parameter]
    public long? Id { get; set; }

    private bool IsEdit => Id.HasValue;

    private UserFormModel model = new();
    private bool loading;
    private bool submitting;

    protected override async Task OnInitializedAsync()
    {
        if (IsEdit)
        {
            loading = true;
            try
            {
                var user = await UserService.GetUserAsync(Id!.Value);
                model = new UserFormModel
                {
                    Email = user.Email,
                    FirstName = user.FirstName,
                    LastName = user.LastName,
                    Active = user.Active
                };
            }
            catch (Exception ex)
            {
                ToastService.ShowError($"Failed to load user: {ex.Message}");
                Navigation.NavigateTo("/users");
            }
            finally
            {
                loading = false;
            }
        }
    }

    private async Task HandleSubmit()
    {
        submitting = true;

        try
        {
            if (IsEdit)
            {
                await UserService.UpdateUserAsync(Id!.Value, new UpdateUserRequest(
                    model.FirstName,
                    model.LastName,
                    model.Active));
                ToastService.ShowSuccess("User updated successfully");
            }
            else
            {
                await UserService.CreateUserAsync(new CreateUserRequest(
                    model.Email,
                    model.Password,
                    model.FirstName,
                    model.LastName));
                ToastService.ShowSuccess("User created successfully");
            }

            Navigation.NavigateTo("/users");
        }
        catch (Exception ex)
        {
            ToastService.ShowError(ex.Message);
        }
        finally
        {
            submitting = false;
        }
    }

    private void Cancel() => Navigation.NavigateTo("/users");

    private class UserFormModel
    {
        [Required]
        [EmailAddress]
        public string Email { get; set; } = string.Empty;

        [Required(ErrorMessage = "Password is required for new users")]
        [MinLength(8)]
        public string Password { get; set; } = string.Empty;

        [Required]
        [MaxLength(100)]
        public string FirstName { get; set; } = string.Empty;

        [Required]
        [MaxLength(100)]
        public string LastName { get; set; } = string.Empty;

        public bool Active { get; set; } = true;
    }
}
```

---

## Reusable Components

### Components/Users/UserCard.razor
```razor
<div class="card h-100">
    <div class="card-body">
        <div class="d-flex align-items-center mb-3">
            <div class="avatar me-3">
                @User.FirstName[0]@User.LastName[0]
            </div>
            <div>
                <h5 class="card-title mb-0">@User.FirstName @User.LastName</h5>
                <small class="text-muted">@User.Email</small>
            </div>
        </div>

        <div class="mb-2">
            <span class="badge @(User.Active ? "bg-success" : "bg-secondary")">
                @(User.Active ? "Active" : "Inactive")
            </span>
            <span class="badge bg-info">@User.Role</span>
        </div>

        <p class="card-text text-muted small">
            Created: @User.CreatedAt.ToString("MMM dd, yyyy")
        </p>
    </div>
    <div class="card-footer bg-transparent">
        <div class="btn-group btn-group-sm">
            <button class="btn btn-outline-primary" @onclick="OnEdit">
                <span class="oi oi-pencil"></span> Edit
            </button>
            <button class="btn btn-outline-danger" @onclick="ConfirmDelete">
                <span class="oi oi-trash"></span> Delete
            </button>
        </div>
    </div>
</div>

@if (showConfirmation)
{
    <Modal Title="Confirm Delete" OnClose="CancelDelete">
        <BodyContent>
            <p>Are you sure you want to delete @User.FirstName @User.LastName?</p>
        </BodyContent>
        <FooterContent>
            <button class="btn btn-secondary" @onclick="CancelDelete">Cancel</button>
            <button class="btn btn-danger" @onclick="HandleDelete">Delete</button>
        </FooterContent>
    </Modal>
}

@code {
    [Parameter, EditorRequired]
    public UserResponse User { get; set; } = default!;

    [Parameter]
    public EventCallback OnEdit { get; set; }

    [Parameter]
    public EventCallback OnDelete { get; set; }

    private bool showConfirmation;

    private void ConfirmDelete() => showConfirmation = true;

    private void CancelDelete() => showConfirmation = false;

    private async Task HandleDelete()
    {
        showConfirmation = false;
        await OnDelete.InvokeAsync();
    }
}
```

### Components/Common/Modal.razor
```razor
<div class="modal fade show" style="display: block;" tabindex="-1">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title">@Title</h5>
                <button type="button" class="btn-close" @onclick="Close"></button>
            </div>
            <div class="modal-body">
                @BodyContent
            </div>
            @if (FooterContent is not null)
            {
                <div class="modal-footer">
                    @FooterContent
                </div>
            }
        </div>
    </div>
</div>
<div class="modal-backdrop fade show"></div>

@code {
    [Parameter]
    public string Title { get; set; } = "Modal";

    [Parameter]
    public RenderFragment? BodyContent { get; set; }

    [Parameter]
    public RenderFragment? FooterContent { get; set; }

    [Parameter]
    public EventCallback OnClose { get; set; }

    private async Task Close() => await OnClose.InvokeAsync();
}
```

### Components/Common/LoadingSpinner.razor
```razor
<div class="d-flex justify-content-center align-items-center p-4">
    <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
    </div>
    @if (!string.IsNullOrEmpty(Message))
    {
        <span class="ms-2">@Message</span>
    }
</div>

@code {
    [Parameter]
    public string? Message { get; set; }
}
```

---

## Services

### Services/IUserService.cs
```csharp
using MyApp.Shared.DTOs;

namespace MyApp.Client.Services;

public interface IUserService
{
    Task<PagedResponse<UserResponse>> GetUsersAsync(int page = 1, int pageSize = 10);
    Task<UserResponse> GetUserAsync(long id);
    Task<UserResponse> CreateUserAsync(CreateUserRequest request);
    Task<UserResponse> UpdateUserAsync(long id, UpdateUserRequest request);
    Task<bool> DeleteUserAsync(long id);
}
```

### Services/UserService.cs
```csharp
using System.Net.Http.Json;
using MyApp.Shared.DTOs;

namespace MyApp.Client.Services;

public class UserService : IUserService
{
    private readonly HttpClient _httpClient;

    public UserService(HttpClient httpClient)
    {
        _httpClient = httpClient;
    }

    public async Task<PagedResponse<UserResponse>> GetUsersAsync(int page = 1, int pageSize = 10)
    {
        var response = await _httpClient.GetFromJsonAsync<PagedResponse<UserResponse>>(
            $"api/users?page={page}&pageSize={pageSize}");

        return response ?? new PagedResponse<UserResponse>([], 0, 0, 0, 0);
    }

    public async Task<UserResponse> GetUserAsync(long id)
    {
        var response = await _httpClient.GetAsync($"api/users/{id}");
        response.EnsureSuccessStatusCode();
        return (await response.Content.ReadFromJsonAsync<UserResponse>())!;
    }

    public async Task<UserResponse> CreateUserAsync(CreateUserRequest request)
    {
        var response = await _httpClient.PostAsJsonAsync("api/users", request);

        if (!response.IsSuccessStatusCode)
        {
            var error = await response.Content.ReadAsStringAsync();
            throw new HttpRequestException(error);
        }

        return (await response.Content.ReadFromJsonAsync<UserResponse>())!;
    }

    public async Task<UserResponse> UpdateUserAsync(long id, UpdateUserRequest request)
    {
        var response = await _httpClient.PutAsJsonAsync($"api/users/{id}", request);

        if (!response.IsSuccessStatusCode)
        {
            var error = await response.Content.ReadAsStringAsync();
            throw new HttpRequestException(error);
        }

        return (await response.Content.ReadFromJsonAsync<UserResponse>())!;
    }

    public async Task<bool> DeleteUserAsync(long id)
    {
        var response = await _httpClient.DeleteAsync($"api/users/{id}");
        return response.IsSuccessStatusCode;
    }
}
```

---

## State Management

### State/AppState.cs
```csharp
namespace MyApp.Client.State;

public class AppState
{
    private UserResponse? _currentUser;
    private bool _isDarkMode;

    public UserResponse? CurrentUser
    {
        get => _currentUser;
        set
        {
            _currentUser = value;
            NotifyStateChanged();
        }
    }

    public bool IsDarkMode
    {
        get => _isDarkMode;
        set
        {
            _isDarkMode = value;
            NotifyStateChanged();
        }
    }

    public event Action? OnChange;

    private void NotifyStateChanged() => OnChange?.Invoke();
}
```

### Using State in Components
```razor
@inject AppState State
@implements IDisposable

<div class="@(State.IsDarkMode ? "dark-theme" : "light-theme")">
    @if (State.CurrentUser is not null)
    {
        <span>Welcome, @State.CurrentUser.FirstName!</span>
    }
</div>

@code {
    protected override void OnInitialized()
    {
        State.OnChange += StateHasChanged;
    }

    public void Dispose()
    {
        State.OnChange -= StateHasChanged;
    }
}
```

---

## JavaScript Interop

### wwwroot/js/interop.js
```javascript
window.blazorInterop = {
    showAlert: function (message) {
        alert(message);
    },

    getWindowWidth: function () {
        return window.innerWidth;
    },

    setLocalStorage: function (key, value) {
        localStorage.setItem(key, value);
    },

    getLocalStorage: function (key) {
        return localStorage.getItem(key);
    },

    downloadFile: function (filename, contentType, content) {
        const blob = new Blob([content], { type: contentType });
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = filename;
        link.click();
        URL.revokeObjectURL(url);
    }
};
```

### Services/JsInteropService.cs
```csharp
using Microsoft.JSInterop;

namespace MyApp.Client.Services;

public class JsInteropService
{
    private readonly IJSRuntime _jsRuntime;

    public JsInteropService(IJSRuntime jsRuntime)
    {
        _jsRuntime = jsRuntime;
    }

    public async Task ShowAlertAsync(string message)
    {
        await _jsRuntime.InvokeVoidAsync("blazorInterop.showAlert", message);
    }

    public async Task<int> GetWindowWidthAsync()
    {
        return await _jsRuntime.InvokeAsync<int>("blazorInterop.getWindowWidth");
    }

    public async Task DownloadFileAsync(string filename, string contentType, byte[] content)
    {
        await _jsRuntime.InvokeVoidAsync("blazorInterop.downloadFile",
            filename, contentType, content);
    }
}
```

---

## Testing

### Component Unit Tests with bUnit
```csharp
using Bunit;
using FluentAssertions;
using Microsoft.Extensions.DependencyInjection;
using Moq;
using MyApp.Client.Components.Users;
using MyApp.Client.Services;
using MyApp.Shared.DTOs;
using Xunit;

namespace MyApp.Client.Tests.Components;

public class UserCardTests : TestContext
{
    [Fact]
    public void UserCard_DisplaysUserInfo()
    {
        // Arrange
        var user = new UserResponse(
            Id: 1,
            Email: "test@example.com",
            FirstName: "John",
            LastName: "Doe",
            Role: "User",
            Active: true,
            CreatedAt: DateTime.UtcNow);

        // Act
        var cut = RenderComponent<UserCard>(parameters => parameters
            .Add(p => p.User, user));

        // Assert
        cut.Find(".card-title").TextContent.Should().Contain("John Doe");
        cut.Find("small.text-muted").TextContent.Should().Contain("test@example.com");
        cut.Find(".badge.bg-success").TextContent.Should().Contain("Active");
    }

    [Fact]
    public void UserCard_EditButton_InvokesCallback()
    {
        // Arrange
        var user = new UserResponse(1, "test@example.com", "John", "Doe", "User", true, DateTime.UtcNow);
        var editClicked = false;

        var cut = RenderComponent<UserCard>(parameters => parameters
            .Add(p => p.User, user)
            .Add(p => p.OnEdit, EventCallback.Factory.Create(this, () => editClicked = true)));

        // Act
        cut.Find("button.btn-outline-primary").Click();

        // Assert
        editClicked.Should().BeTrue();
    }

    [Fact]
    public void UserCard_DeleteButton_ShowsConfirmation()
    {
        // Arrange
        var user = new UserResponse(1, "test@example.com", "John", "Doe", "User", true, DateTime.UtcNow);

        var cut = RenderComponent<UserCard>(parameters => parameters
            .Add(p => p.User, user));

        // Act
        cut.Find("button.btn-outline-danger").Click();

        // Assert
        cut.Find(".modal").Should().NotBeNull();
        cut.Find(".modal-body").TextContent.Should().Contain("Are you sure");
    }
}
```

### Page Tests
```csharp
using Bunit;
using FluentAssertions;
using Microsoft.Extensions.DependencyInjection;
using Moq;
using MyApp.Client.Pages.Users;
using MyApp.Client.Services;
using MyApp.Shared.DTOs;
using Xunit;

namespace MyApp.Client.Tests.Pages;

public class UserListTests : TestContext
{
    private readonly Mock<IUserService> _userServiceMock;

    public UserListTests()
    {
        _userServiceMock = new Mock<IUserService>();
        Services.AddSingleton(_userServiceMock.Object);

        // Add NavigationManager
        Services.AddSingleton(new MockNavigationManager());
    }

    [Fact]
    public async Task UserList_LoadsAndDisplaysUsers()
    {
        // Arrange
        var users = new List<UserResponse>
        {
            new(1, "user1@example.com", "John", "Doe", "User", true, DateTime.UtcNow),
            new(2, "user2@example.com", "Jane", "Smith", "Admin", true, DateTime.UtcNow)
        };

        _userServiceMock.Setup(s => s.GetUsersAsync(1, 10))
            .ReturnsAsync(new PagedResponse<UserResponse>(users, 1, 10, 2, 1));

        // Act
        var cut = RenderComponent<UserList>();
        await cut.InvokeAsync(() => Task.Delay(100)); // Wait for async load

        // Assert
        cut.FindAll(".card").Count.Should().Be(2);
    }

    [Fact]
    public async Task UserList_ShowsLoading_WhenFetching()
    {
        // Arrange
        _userServiceMock.Setup(s => s.GetUsersAsync(1, 10))
            .Returns(async () =>
            {
                await Task.Delay(1000);
                return new PagedResponse<UserResponse>([], 1, 10, 0, 0);
            });

        // Act
        var cut = RenderComponent<UserList>();

        // Assert
        cut.Find(".spinner-border").Should().NotBeNull();
    }
}
```

---

## .NET 8 Blazor United (SSR + Interactive)

### Interactive Render Modes
```razor
@* Static SSR (default) *@
<StaticComponent />

@* Interactive Server *@
<InteractiveComponent @rendermode="InteractiveServer" />

@* Interactive WebAssembly *@
<InteractiveComponent @rendermode="InteractiveWebAssembly" />

@* Interactive Auto (starts Server, switches to WASM) *@
<InteractiveComponent @rendermode="InteractiveAuto" />
```

### Streaming Rendering
```razor
@page "/users"
@attribute [StreamRendering]

@if (users is null)
{
    <p>Loading users...</p>
}
else
{
    @foreach (var user in users)
    {
        <UserCard User="user" />
    }
}

@code {
    private List<UserResponse>? users;

    protected override async Task OnInitializedAsync()
    {
        // This will stream the content as data loads
        users = await UserService.GetUsersAsync();
    }
}
```

---

## Build & Run Commands

```bash
# Create new Blazor WebAssembly app
dotnet new blazorwasm -o MyApp.Client

# Create Blazor Server app
dotnet new blazorserver -o MyApp.Server

# Create Blazor Web App (.NET 8 - United)
dotnet new blazor -o MyApp

# Run
dotnet run --project MyApp.Client

# Run with watch
dotnet watch run --project MyApp.Client

# Build
dotnet build

# Publish (WebAssembly)
dotnet publish -c Release

# Run tests
dotnet test

# Add bUnit for testing
dotnet add package bunit
```

---

## Best Practices

### DO
- ✓ Use `@key` directive for list rendering performance
- ✓ Use `EventCallback` instead of `Action` for component events
- ✓ Implement `IDisposable` for cleanup (event handlers, timers)
- ✓ Use `@rendermode` appropriately for interactive components
- ✓ Use `StateHasChanged()` only when necessary
- ✓ Use virtualization for large lists (`<Virtualize>`)
- ✓ Use `CascadingValue` for shared state
- ✓ Use typed HttpClient for API calls

### DON'T
- ✗ Don't call JavaScript in `OnInitialized` (use `OnAfterRender`)
- ✗ Don't mutate parameters directly
- ✗ Don't use synchronous HTTP calls
- ✗ Don't ignore component lifecycle
- ✗ Don't overuse JavaScript interop
- ✗ Don't forget to handle loading/error states

---

## Framework Comparison

| Feature | Blazor | React | Angular | Vue |
|---------|--------|-------|---------|-----|
| Language | C# | JavaScript/TypeScript | TypeScript | JavaScript/TypeScript |
| Rendering | WASM/Server | Virtual DOM | Real DOM | Virtual DOM |
| State | Built-in/Fluxor | Redux/Context | Services/NgRx | Vuex/Pinia |
| Learning Curve | Moderate (.NET devs) | Moderate | Steep | Gentle |
| Performance | Very Good | Excellent | Good | Excellent |
| Bundle Size | Larger (WASM) | Small | Medium | Small |
| SEO | Server/SSR | SSR needed | SSR needed | SSR needed |

---

## References

- [Blazor Documentation](https://docs.microsoft.com/aspnet/core/blazor)
- [Blazor University](https://blazor-university.com/)
- [bUnit Testing](https://bunit.dev/)
- [Blazored Libraries](https://github.com/Blazored)
- [Awesome Blazor](https://github.com/AdrienTorris/awesome-blazor)
