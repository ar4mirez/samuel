# Android Jetpack Compose Framework Guide

> **Applies to**: Android SDK 24+, Jetpack Compose 1.5+, Kotlin 1.9+, Material Design 3

---

## Overview

Jetpack Compose is Android's modern toolkit for building native UI. It uses a declarative approach with Kotlin, eliminating the need for XML layouts and providing a more intuitive way to create user interfaces.

**Best For**: Modern Android apps, UI-heavy applications, new Android projects

**Key Features**:
- Declarative UI with Kotlin
- Less code, fewer bugs
- Powerful state management
- Material Design 3 built-in
- Interoperable with existing Views
- Strong tooling support in Android Studio

---

## Project Structure

```
app/
├── src/
│   └── main/
│       ├── java/com/example/myapp/
│       │   ├── MyApplication.kt           # Application class
│       │   ├── MainActivity.kt            # Main activity
│       │   ├── navigation/
│       │   │   └── AppNavigation.kt       # Navigation graph
│       │   ├── ui/
│       │   │   ├── theme/                 # Theme configuration
│       │   │   │   ├── Color.kt
│       │   │   │   ├── Type.kt
│       │   │   │   └── Theme.kt
│       │   │   ├── screens/               # Screen composables
│       │   │   │   ├── home/
│       │   │   │   │   ├── HomeScreen.kt
│       │   │   │   │   └── HomeViewModel.kt
│       │   │   │   ├── detail/
│       │   │   │   │   ├── DetailScreen.kt
│       │   │   │   │   └── DetailViewModel.kt
│       │   │   │   └── auth/
│       │   │   │       ├── LoginScreen.kt
│       │   │   │       └── AuthViewModel.kt
│       │   │   └── components/            # Reusable components
│       │   │       ├── AppBar.kt
│       │   │       ├── LoadingIndicator.kt
│       │   │       └── ErrorMessage.kt
│       │   ├── data/
│       │   │   ├── model/                 # Data models
│       │   │   │   └── User.kt
│       │   │   ├── remote/                # API services
│       │   │   │   ├── ApiService.kt
│       │   │   │   └── dto/
│       │   │   ├── local/                 # Local storage
│       │   │   │   ├── AppDatabase.kt
│       │   │   │   └── dao/
│       │   │   └── repository/            # Repositories
│       │   │       └── UserRepository.kt
│       │   └── di/                        # Dependency injection
│       │       ├── AppModule.kt
│       │       └── NetworkModule.kt
│       ├── res/
│       │   ├── values/
│       │   │   ├── strings.xml
│       │   │   └── themes.xml
│       │   └── drawable/
│       └── AndroidManifest.xml
├── build.gradle.kts
└── proguard-rules.pro
```

---

## Dependencies (build.gradle.kts - App Module)

```kotlin
plugins {
    id("com.android.application")
    id("org.jetbrains.kotlin.android")
    id("com.google.devtools.ksp")
    id("com.google.dagger.hilt.android")
    kotlin("plugin.serialization")
}

android {
    namespace = "com.example.myapp"
    compileSdk = 34

    defaultConfig {
        applicationId = "com.example.myapp"
        minSdk = 24
        targetSdk = 34
        versionCode = 1
        versionName = "1.0"

        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
        vectorDrawables {
            useSupportLibrary = true
        }
    }

    buildTypes {
        release {
            isMinifyEnabled = true
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
        }
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }

    kotlinOptions {
        jvmTarget = "17"
    }

    buildFeatures {
        compose = true
    }

    composeOptions {
        kotlinCompilerExtensionVersion = "1.5.8"
    }

    packaging {
        resources {
            excludes += "/META-INF/{AL2.0,LGPL2.1}"
        }
    }
}

dependencies {
    val composeBom = platform("androidx.compose:compose-bom:2024.01.00")
    implementation(composeBom)

    // Compose UI
    implementation("androidx.compose.ui:ui")
    implementation("androidx.compose.ui:ui-graphics")
    implementation("androidx.compose.ui:ui-tooling-preview")
    implementation("androidx.compose.material3:material3")
    implementation("androidx.compose.material:material-icons-extended")

    // Activity & Lifecycle
    implementation("androidx.activity:activity-compose:1.8.2")
    implementation("androidx.lifecycle:lifecycle-runtime-compose:2.7.0")
    implementation("androidx.lifecycle:lifecycle-viewmodel-compose:2.7.0")

    // Navigation
    implementation("androidx.navigation:navigation-compose:2.7.6")
    implementation("androidx.hilt:hilt-navigation-compose:1.1.0")

    // Dependency Injection - Hilt
    implementation("com.google.dagger:hilt-android:2.50")
    ksp("com.google.dagger:hilt-compiler:2.50")

    // Networking - Retrofit + OkHttp
    implementation("com.squareup.retrofit2:retrofit:2.9.0")
    implementation("com.squareup.okhttp3:okhttp:4.12.0")
    implementation("com.squareup.okhttp3:logging-interceptor:4.12.0")
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.6.2")
    implementation("com.jakewharton.retrofit:retrofit2-kotlinx-serialization-converter:1.0.0")

    // Local Database - Room
    implementation("androidx.room:room-runtime:2.6.1")
    implementation("androidx.room:room-ktx:2.6.1")
    ksp("androidx.room:room-compiler:2.6.1")

    // DataStore for preferences
    implementation("androidx.datastore:datastore-preferences:1.0.0")

    // Image Loading - Coil
    implementation("io.coil-kt:coil-compose:2.5.0")

    // Coroutines
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-android:1.7.3")

    // Testing
    testImplementation("junit:junit:4.13.2")
    testImplementation("org.jetbrains.kotlinx:kotlinx-coroutines-test:1.7.3")
    testImplementation("io.mockk:mockk:1.13.9")
    testImplementation("app.cash.turbine:turbine:1.0.0")

    androidTestImplementation(composeBom)
    androidTestImplementation("androidx.compose.ui:ui-test-junit4")
    androidTestImplementation("androidx.test.ext:junit:1.1.5")
    androidTestImplementation("androidx.test.espresso:espresso-core:3.5.1")
    androidTestImplementation("com.google.dagger:hilt-android-testing:2.50")

    debugImplementation("androidx.compose.ui:ui-tooling")
    debugImplementation("androidx.compose.ui:ui-test-manifest")
}
```

---

## Application Setup

### Application Class

```kotlin
// MyApplication.kt
package com.example.myapp

import android.app.Application
import dagger.hilt.android.HiltAndroidApp

@HiltAndroidApp
class MyApplication : Application()
```

### MainActivity

```kotlin
// MainActivity.kt
package com.example.myapp

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.ui.Modifier
import com.example.myapp.navigation.AppNavigation
import com.example.myapp.ui.theme.MyAppTheme
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            MyAppTheme {
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    AppNavigation()
                }
            }
        }
    }
}
```

---

## Theme Configuration

### Colors

```kotlin
// ui/theme/Color.kt
package com.example.myapp.ui.theme

import androidx.compose.ui.graphics.Color

// Light Theme
val Purple40 = Color(0xFF6650a4)
val PurpleGrey40 = Color(0xFF625b71)
val Pink40 = Color(0xFF7D5260)

// Dark Theme
val Purple80 = Color(0xFFD0BCFF)
val PurpleGrey80 = Color(0xFFCCC2DC)
val Pink80 = Color(0xFFEFB8C8)

// Custom Colors
val Primary = Color(0xFF1976D2)
val PrimaryVariant = Color(0xFF004BA0)
val Secondary = Color(0xFF03DAC6)
val Error = Color(0xFFB00020)
val Background = Color(0xFFFFFBFE)
val Surface = Color(0xFFFFFBFE)
val OnPrimary = Color.White
val OnSecondary = Color.Black
val OnBackground = Color(0xFF1C1B1F)
val OnSurface = Color(0xFF1C1B1F)
```

### Typography

```kotlin
// ui/theme/Type.kt
package com.example.myapp.ui.theme

import androidx.compose.material3.Typography
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.sp

val Typography = Typography(
    displayLarge = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.Normal,
        fontSize = 57.sp,
        lineHeight = 64.sp,
        letterSpacing = (-0.25).sp
    ),
    headlineLarge = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.Normal,
        fontSize = 32.sp,
        lineHeight = 40.sp,
        letterSpacing = 0.sp
    ),
    titleLarge = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.Normal,
        fontSize = 22.sp,
        lineHeight = 28.sp,
        letterSpacing = 0.sp
    ),
    bodyLarge = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.Normal,
        fontSize = 16.sp,
        lineHeight = 24.sp,
        letterSpacing = 0.5.sp
    ),
    labelSmall = TextStyle(
        fontFamily = FontFamily.Default,
        fontWeight = FontWeight.Medium,
        fontSize = 11.sp,
        lineHeight = 16.sp,
        letterSpacing = 0.5.sp
    )
)
```

### Theme

```kotlin
// ui/theme/Theme.kt
package com.example.myapp.ui.theme

import android.app.Activity
import android.os.Build
import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.runtime.SideEffect
import androidx.compose.ui.graphics.toArgb
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.platform.LocalView
import androidx.core.view.WindowCompat

private val DarkColorScheme = darkColorScheme(
    primary = Purple80,
    secondary = PurpleGrey80,
    tertiary = Pink80
)

private val LightColorScheme = lightColorScheme(
    primary = Purple40,
    secondary = PurpleGrey40,
    tertiary = Pink40
)

@Composable
fun MyAppTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    dynamicColor: Boolean = true,
    content: @Composable () -> Unit
) {
    val colorScheme = when {
        dynamicColor && Build.VERSION.SDK_INT >= Build.VERSION_CODES.S -> {
            val context = LocalContext.current
            if (darkTheme) dynamicDarkColorScheme(context) else dynamicLightColorScheme(context)
        }
        darkTheme -> DarkColorScheme
        else -> LightColorScheme
    }

    val view = LocalView.current
    if (!view.isInEditMode) {
        SideEffect {
            val window = (view.context as Activity).window
            window.statusBarColor = colorScheme.primary.toArgb()
            WindowCompat.getInsetsController(window, view).isAppearanceLightStatusBars = !darkTheme
        }
    }

    MaterialTheme(
        colorScheme = colorScheme,
        typography = Typography,
        content = content
    )
}
```

---

## Navigation

```kotlin
// navigation/AppNavigation.kt
package com.example.myapp.navigation

import androidx.compose.runtime.Composable
import androidx.navigation.NavType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import androidx.navigation.navArgument
import com.example.myapp.ui.screens.auth.LoginScreen
import com.example.myapp.ui.screens.detail.DetailScreen
import com.example.myapp.ui.screens.home.HomeScreen

sealed class Screen(val route: String) {
    object Login : Screen("login")
    object Home : Screen("home")
    object Detail : Screen("detail/{userId}") {
        fun createRoute(userId: String) = "detail/$userId"
    }
}

@Composable
fun AppNavigation() {
    val navController = rememberNavController()

    NavHost(
        navController = navController,
        startDestination = Screen.Home.route
    ) {
        composable(Screen.Login.route) {
            LoginScreen(
                onLoginSuccess = {
                    navController.navigate(Screen.Home.route) {
                        popUpTo(Screen.Login.route) { inclusive = true }
                    }
                }
            )
        }

        composable(Screen.Home.route) {
            HomeScreen(
                onUserClick = { userId ->
                    navController.navigate(Screen.Detail.createRoute(userId))
                },
                onLogout = {
                    navController.navigate(Screen.Login.route) {
                        popUpTo(Screen.Home.route) { inclusive = true }
                    }
                }
            )
        }

        composable(
            route = Screen.Detail.route,
            arguments = listOf(
                navArgument("userId") { type = NavType.StringType }
            )
        ) { backStackEntry ->
            val userId = backStackEntry.arguments?.getString("userId") ?: ""
            DetailScreen(
                userId = userId,
                onBackClick = { navController.popBackStack() }
            )
        }
    }
}
```

---

## State Management

### UI State

```kotlin
// ui/screens/home/HomeUiState.kt
package com.example.myapp.ui.screens.home

import com.example.myapp.data.model.User

sealed interface HomeUiState {
    object Loading : HomeUiState
    data class Success(val users: List<User>) : HomeUiState
    data class Error(val message: String) : HomeUiState
}
```

### ViewModel

```kotlin
// ui/screens/home/HomeViewModel.kt
package com.example.myapp.ui.screens.home

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.example.myapp.data.repository.UserRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.*
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class HomeViewModel @Inject constructor(
    private val userRepository: UserRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow<HomeUiState>(HomeUiState.Loading)
    val uiState: StateFlow<HomeUiState> = _uiState.asStateFlow()

    private val _searchQuery = MutableStateFlow("")
    val searchQuery: StateFlow<String> = _searchQuery.asStateFlow()

    init {
        loadUsers()
    }

    fun loadUsers() {
        viewModelScope.launch {
            _uiState.value = HomeUiState.Loading
            userRepository.getUsers()
                .catch { e ->
                    _uiState.value = HomeUiState.Error(e.message ?: "Unknown error")
                }
                .collect { users ->
                    _uiState.value = HomeUiState.Success(users)
                }
        }
    }

    fun onSearchQueryChange(query: String) {
        _searchQuery.value = query
    }

    fun deleteUser(userId: String) {
        viewModelScope.launch {
            try {
                userRepository.deleteUser(userId)
                loadUsers()
            } catch (e: Exception) {
                _uiState.value = HomeUiState.Error(e.message ?: "Failed to delete user")
            }
        }
    }

    fun refresh() {
        loadUsers()
    }
}
```

---

## Screen Composables

### Home Screen

```kotlin
// ui/screens/home/HomeScreen.kt
package com.example.myapp.ui.screens.home

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Delete
import androidx.compose.material.icons.filled.Logout
import androidx.compose.material.icons.filled.Refresh
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import com.example.myapp.data.model.User
import com.example.myapp.ui.components.ErrorMessage
import com.example.myapp.ui.components.LoadingIndicator

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HomeScreen(
    onUserClick: (String) -> Unit,
    onLogout: () -> Unit,
    viewModel: HomeViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsStateWithLifecycle()
    val searchQuery by viewModel.searchQuery.collectAsStateWithLifecycle()

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Users") },
                actions = {
                    IconButton(onClick = { viewModel.refresh() }) {
                        Icon(Icons.Default.Refresh, contentDescription = "Refresh")
                    }
                    IconButton(onClick = onLogout) {
                        Icon(Icons.Default.Logout, contentDescription = "Logout")
                    }
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            // Search Bar
            OutlinedTextField(
                value = searchQuery,
                onValueChange = viewModel::onSearchQueryChange,
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                placeholder = { Text("Search users...") },
                singleLine = true
            )

            // Content
            when (val state = uiState) {
                is HomeUiState.Loading -> {
                    LoadingIndicator()
                }
                is HomeUiState.Success -> {
                    UserList(
                        users = state.users.filter {
                            it.name.contains(searchQuery, ignoreCase = true) ||
                            it.email.contains(searchQuery, ignoreCase = true)
                        },
                        onUserClick = onUserClick,
                        onDeleteClick = viewModel::deleteUser
                    )
                }
                is HomeUiState.Error -> {
                    ErrorMessage(
                        message = state.message,
                        onRetry = viewModel::loadUsers
                    )
                }
            }
        }
    }
}

@Composable
private fun UserList(
    users: List<User>,
    onUserClick: (String) -> Unit,
    onDeleteClick: (String) -> Unit
) {
    if (users.isEmpty()) {
        Box(
            modifier = Modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            Text("No users found")
        }
    } else {
        LazyColumn(
            contentPadding = PaddingValues(16.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            items(
                items = users,
                key = { it.id }
            ) { user ->
                UserCard(
                    user = user,
                    onClick = { onUserClick(user.id) },
                    onDeleteClick = { onDeleteClick(user.id) }
                )
            }
        }
    }
}

@Composable
private fun UserCard(
    user: User,
    onClick: () -> Unit,
    onDeleteClick: () -> Unit
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = user.name,
                    style = MaterialTheme.typography.titleMedium
                )
                Text(
                    text = user.email,
                    style = MaterialTheme.typography.bodyMedium,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
            }
            IconButton(onClick = onDeleteClick) {
                Icon(
                    Icons.Default.Delete,
                    contentDescription = "Delete",
                    tint = MaterialTheme.colorScheme.error
                )
            }
        }
    }
}
```

### Login Screen

```kotlin
// ui/screens/auth/LoginScreen.kt
package com.example.myapp.ui.screens.auth

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Visibility
import androidx.compose.material.icons.filled.VisibilityOff
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle

@Composable
fun LoginScreen(
    onLoginSuccess: () -> Unit,
    viewModel: AuthViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsStateWithLifecycle()
    var email by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }
    var passwordVisible by remember { mutableStateOf(false) }

    LaunchedEffect(uiState) {
        if (uiState is AuthUiState.Success) {
            onLoginSuccess()
        }
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text(
            text = "Welcome Back",
            style = MaterialTheme.typography.headlineLarge
        )

        Spacer(modifier = Modifier.height(32.dp))

        OutlinedTextField(
            value = email,
            onValueChange = { email = it },
            modifier = Modifier.fillMaxWidth(),
            label = { Text("Email") },
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Email),
            singleLine = true,
            isError = uiState is AuthUiState.Error
        )

        Spacer(modifier = Modifier.height(16.dp))

        OutlinedTextField(
            value = password,
            onValueChange = { password = it },
            modifier = Modifier.fillMaxWidth(),
            label = { Text("Password") },
            visualTransformation = if (passwordVisible) {
                VisualTransformation.None
            } else {
                PasswordVisualTransformation()
            },
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
            singleLine = true,
            isError = uiState is AuthUiState.Error,
            trailingIcon = {
                IconButton(onClick = { passwordVisible = !passwordVisible }) {
                    Icon(
                        if (passwordVisible) Icons.Default.VisibilityOff else Icons.Default.Visibility,
                        contentDescription = "Toggle password visibility"
                    )
                }
            }
        )

        if (uiState is AuthUiState.Error) {
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                text = (uiState as AuthUiState.Error).message,
                color = MaterialTheme.colorScheme.error,
                style = MaterialTheme.typography.bodySmall
            )
        }

        Spacer(modifier = Modifier.height(24.dp))

        Button(
            onClick = { viewModel.login(email, password) },
            modifier = Modifier.fillMaxWidth(),
            enabled = email.isNotBlank() && password.isNotBlank() && uiState !is AuthUiState.Loading
        ) {
            if (uiState is AuthUiState.Loading) {
                CircularProgressIndicator(
                    modifier = Modifier.size(20.dp),
                    color = MaterialTheme.colorScheme.onPrimary
                )
            } else {
                Text("Login")
            }
        }
    }
}
```

---

## Reusable Components

### Loading Indicator

```kotlin
// ui/components/LoadingIndicator.kt
package com.example.myapp.ui.components

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier

@Composable
fun LoadingIndicator(
    modifier: Modifier = Modifier
) {
    Box(
        modifier = modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        CircularProgressIndicator()
    }
}
```

### Error Message

```kotlin
// ui/components/ErrorMessage.kt
package com.example.myapp.ui.components

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Error
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

@Composable
fun ErrorMessage(
    message: String,
    onRetry: (() -> Unit)? = null,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .fillMaxSize()
            .padding(16.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Icon(
            Icons.Default.Error,
            contentDescription = null,
            modifier = Modifier.size(48.dp),
            tint = MaterialTheme.colorScheme.error
        )
        Spacer(modifier = Modifier.height(16.dp))
        Text(
            text = message,
            style = MaterialTheme.typography.bodyLarge,
            color = MaterialTheme.colorScheme.error
        )
        onRetry?.let {
            Spacer(modifier = Modifier.height(16.dp))
            Button(onClick = it) {
                Text("Retry")
            }
        }
    }
}
```

---

## Data Layer

### Model

```kotlin
// data/model/User.kt
package com.example.myapp.data.model

import kotlinx.serialization.Serializable

@Serializable
data class User(
    val id: String,
    val email: String,
    val name: String,
    val avatarUrl: String? = null,
)
```

### API Service

```kotlin
// data/remote/ApiService.kt
package com.example.myapp.data.remote

import com.example.myapp.data.model.User
import com.example.myapp.data.remote.dto.LoginRequest
import com.example.myapp.data.remote.dto.LoginResponse
import retrofit2.http.*

interface ApiService {

    @GET("users")
    suspend fun getUsers(): List<User>

    @GET("users/{id}")
    suspend fun getUser(@Path("id") id: String): User

    @POST("users")
    suspend fun createUser(@Body user: User): User

    @DELETE("users/{id}")
    suspend fun deleteUser(@Path("id") id: String)

    @POST("auth/login")
    suspend fun login(@Body request: LoginRequest): LoginResponse
}
```

### DTOs

```kotlin
// data/remote/dto/AuthDto.kt
package com.example.myapp.data.remote.dto

import kotlinx.serialization.Serializable

@Serializable
data class LoginRequest(
    val email: String,
    val password: String,
)

@Serializable
data class LoginResponse(
    val token: String,
    val userId: String,
)
```

### Repository

```kotlin
// data/repository/UserRepository.kt
package com.example.myapp.data.repository

import com.example.myapp.data.model.User
import com.example.myapp.data.remote.ApiService
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class UserRepository @Inject constructor(
    private val apiService: ApiService
) {

    fun getUsers(): Flow<List<User>> = flow {
        val users = apiService.getUsers()
        emit(users)
    }

    suspend fun getUser(id: String): User {
        return apiService.getUser(id)
    }

    suspend fun createUser(user: User): User {
        return apiService.createUser(user)
    }

    suspend fun deleteUser(id: String) {
        apiService.deleteUser(id)
    }
}
```

---

## Dependency Injection (Hilt)

### App Module

```kotlin
// di/AppModule.kt
package com.example.myapp.di

import android.content.Context
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.preferencesDataStore
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "settings")

@Module
@InstallIn(SingletonComponent::class)
object AppModule {

    @Provides
    @Singleton
    fun provideDataStore(
        @ApplicationContext context: Context
    ): DataStore<Preferences> {
        return context.dataStore
    }
}
```

### Network Module

```kotlin
// di/NetworkModule.kt
package com.example.myapp.di

import com.example.myapp.data.remote.ApiService
import com.jakewharton.retrofit2.converter.kotlinx.serialization.asConverterFactory
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import kotlinx.serialization.json.Json
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import java.util.concurrent.TimeUnit
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object NetworkModule {

    @Provides
    @Singleton
    fun provideJson(): Json = Json {
        ignoreUnknownKeys = true
        isLenient = true
        encodeDefaults = true
    }

    @Provides
    @Singleton
    fun provideOkHttpClient(): OkHttpClient {
        return OkHttpClient.Builder()
            .connectTimeout(30, TimeUnit.SECONDS)
            .readTimeout(30, TimeUnit.SECONDS)
            .addInterceptor(HttpLoggingInterceptor().apply {
                level = HttpLoggingInterceptor.Level.BODY
            })
            .build()
    }

    @Provides
    @Singleton
    fun provideRetrofit(
        okHttpClient: OkHttpClient,
        json: Json
    ): Retrofit {
        return Retrofit.Builder()
            .baseUrl("https://api.example.com/")
            .client(okHttpClient)
            .addConverterFactory(json.asConverterFactory("application/json".toMediaType()))
            .build()
    }

    @Provides
    @Singleton
    fun provideApiService(retrofit: Retrofit): ApiService {
        return retrofit.create(ApiService::class.java)
    }
}
```

---

## Testing

### ViewModel Unit Test

```kotlin
// test/HomeViewModelTest.kt
package com.example.myapp.ui.screens.home

import app.cash.turbine.test
import com.example.myapp.data.model.User
import com.example.myapp.data.repository.UserRepository
import io.mockk.coEvery
import io.mockk.mockk
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.flow.flowOf
import kotlinx.coroutines.test.*
import org.junit.After
import org.junit.Before
import org.junit.Test
import kotlin.test.assertEquals
import kotlin.test.assertTrue

@OptIn(ExperimentalCoroutinesApi::class)
class HomeViewModelTest {

    private val testDispatcher = StandardTestDispatcher()
    private lateinit var userRepository: UserRepository
    private lateinit var viewModel: HomeViewModel

    @Before
    fun setup() {
        Dispatchers.setMain(testDispatcher)
        userRepository = mockk()
    }

    @After
    fun tearDown() {
        Dispatchers.resetMain()
    }

    @Test
    fun `loadUsers success updates state`() = runTest {
        val users = listOf(
            User(id = "1", email = "test@example.com", name = "Test User")
        )
        coEvery { userRepository.getUsers() } returns flowOf(users)

        viewModel = HomeViewModel(userRepository)
        advanceUntilIdle()

        viewModel.uiState.test {
            val state = awaitItem()
            assertTrue(state is HomeUiState.Success)
            assertEquals(users, (state as HomeUiState.Success).users)
        }
    }

    @Test
    fun `loadUsers error updates state with error`() = runTest {
        coEvery { userRepository.getUsers() } throws RuntimeException("Network error")

        viewModel = HomeViewModel(userRepository)
        advanceUntilIdle()

        viewModel.uiState.test {
            val state = awaitItem()
            assertTrue(state is HomeUiState.Error)
            assertEquals("Network error", (state as HomeUiState.Error).message)
        }
    }
}
```

### Compose UI Test

```kotlin
// androidTest/HomeScreenTest.kt
package com.example.myapp

import androidx.compose.ui.test.*
import androidx.compose.ui.test.junit4.createComposeRule
import com.example.myapp.data.model.User
import com.example.myapp.ui.screens.home.HomeScreen
import com.example.myapp.ui.screens.home.HomeUiState
import com.example.myapp.ui.screens.home.HomeViewModel
import com.example.myapp.ui.theme.MyAppTheme
import io.mockk.every
import io.mockk.mockk
import kotlinx.coroutines.flow.MutableStateFlow
import org.junit.Rule
import org.junit.Test

class HomeScreenTest {

    @get:Rule
    val composeTestRule = createComposeRule()

    @Test
    fun displays_loading_indicator_when_loading() {
        val viewModel = mockk<HomeViewModel>(relaxed = true)
        every { viewModel.uiState } returns MutableStateFlow(HomeUiState.Loading)
        every { viewModel.searchQuery } returns MutableStateFlow("")

        composeTestRule.setContent {
            MyAppTheme {
                HomeScreen(
                    onUserClick = {},
                    onLogout = {},
                    viewModel = viewModel
                )
            }
        }

        composeTestRule.onNode(hasTestTag("loading_indicator")).assertExists()
    }

    @Test
    fun displays_users_when_success() {
        val users = listOf(
            User(id = "1", email = "test@example.com", name = "Test User")
        )
        val viewModel = mockk<HomeViewModel>(relaxed = true)
        every { viewModel.uiState } returns MutableStateFlow(HomeUiState.Success(users))
        every { viewModel.searchQuery } returns MutableStateFlow("")

        composeTestRule.setContent {
            MyAppTheme {
                HomeScreen(
                    onUserClick = {},
                    onLogout = {},
                    viewModel = viewModel
                )
            }
        }

        composeTestRule.onNodeWithText("Test User").assertExists()
        composeTestRule.onNodeWithText("test@example.com").assertExists()
    }

    @Test
    fun displays_error_message_when_error() {
        val viewModel = mockk<HomeViewModel>(relaxed = true)
        every { viewModel.uiState } returns MutableStateFlow(HomeUiState.Error("Network error"))
        every { viewModel.searchQuery } returns MutableStateFlow("")

        composeTestRule.setContent {
            MyAppTheme {
                HomeScreen(
                    onUserClick = {},
                    onLogout = {},
                    viewModel = viewModel
                )
            }
        }

        composeTestRule.onNodeWithText("Network error").assertExists()
        composeTestRule.onNodeWithText("Retry").assertExists()
    }
}
```

---

## Commands

```bash
# Build debug APK
./gradlew assembleDebug

# Build release APK
./gradlew assembleRelease

# Install on connected device
./gradlew installDebug

# Run unit tests
./gradlew test

# Run instrumented tests
./gradlew connectedAndroidTest

# Lint check
./gradlew lint

# Clean build
./gradlew clean

# Generate signed bundle
./gradlew bundleRelease
```

---

## Best Practices

### Do's
- ✅ Use `collectAsStateWithLifecycle()` for Flow collection
- ✅ Use `remember` for expensive computations
- ✅ Use `LaunchedEffect` for side effects
- ✅ Keep composables small and focused
- ✅ Use `Modifier` as first parameter in composables
- ✅ Use sealed interfaces for UI state
- ✅ Use Hilt for dependency injection
- ✅ Preview composables with `@Preview`

### Don'ts
- ❌ Don't perform heavy operations in composables
- ❌ Don't use mutable state outside StateFlow/State
- ❌ Don't forget to handle configuration changes
- ❌ Don't use hardcoded strings (use resources)
- ❌ Don't skip error handling
- ❌ Don't create ViewModels manually (use Hilt)

---

## Comparison: Compose vs XML Views

| Feature | Jetpack Compose | XML Views |
|---------|-----------------|-----------|
| **Code** | Kotlin only | XML + Kotlin |
| **UI Updates** | Declarative recomposition | Imperative invalidation |
| **State** | Built-in state management | Manual binding |
| **Preview** | @Preview annotation | Layout preview |
| **Learning Curve** | Moderate (new paradigm) | Lower (established) |
| **Performance** | Optimized recomposition | View hierarchy traversal |
| **Animations** | Built-in animation APIs | Animator/Transition |
| **Testing** | Compose testing APIs | Espresso |
| **Interop** | Supports Views in Compose | Can embed Compose |

---

## When to Use Compose

**Choose Compose when**:
- Starting a new Android project
- Building UI-heavy applications
- Want less boilerplate code
- Team knows Kotlin well
- Need rapid UI iteration

**Consider Views when**:
- Maintaining existing XML-based apps
- Complex custom views already exist
- Team more familiar with Views
- Targeting older Android versions (API < 21)
