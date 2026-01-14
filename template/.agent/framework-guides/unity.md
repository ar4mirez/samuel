# Unity Framework Guide

> **Framework**: Unity 2022 LTS+ with C#
> **Type**: Game Engine / Interactive Applications
> **Use Cases**: Games, VR/AR, Simulations, Interactive Media

---

## Overview

Unity is a cross-platform game engine for creating 2D/3D games, VR/AR experiences, simulations, and interactive applications. It uses C# as its primary scripting language.

### When to Use Unity
- ✅ 2D and 3D game development
- ✅ VR/AR applications (Meta Quest, HoloLens, etc.)
- ✅ Interactive simulations and training
- ✅ Cross-platform deployment (PC, Mobile, Console, Web)
- ✅ Rapid prototyping with visual editor

### When NOT to Use Unity
- ❌ Web applications (use Blazor, ASP.NET)
- ❌ Enterprise business software
- ❌ Simple mobile apps without 3D/game elements
- ❌ High-performance AAA graphics (consider Unreal)

---

## Project Structure

```
MyGame/
├── Assets/
│   ├── _Project/                    # Project-specific assets
│   │   ├── Art/
│   │   │   ├── Materials/
│   │   │   ├── Models/
│   │   │   ├── Sprites/
│   │   │   └── Textures/
│   │   ├── Audio/
│   │   │   ├── Music/
│   │   │   └── SFX/
│   │   ├── Prefabs/
│   │   │   ├── Characters/
│   │   │   ├── Environment/
│   │   │   └── UI/
│   │   ├── Scenes/
│   │   │   ├── MainMenu.unity
│   │   │   ├── GameLevel.unity
│   │   │   └── Loading.unity
│   │   ├── Scripts/
│   │   │   ├── Core/               # Core systems
│   │   │   │   ├── GameManager.cs
│   │   │   │   ├── SceneLoader.cs
│   │   │   │   └── ServiceLocator.cs
│   │   │   ├── Player/
│   │   │   │   ├── PlayerController.cs
│   │   │   │   ├── PlayerInput.cs
│   │   │   │   └── PlayerHealth.cs
│   │   │   ├── Enemies/
│   │   │   ├── UI/
│   │   │   │   ├── HUDController.cs
│   │   │   │   └── MenuController.cs
│   │   │   ├── Systems/
│   │   │   │   ├── AudioManager.cs
│   │   │   │   ├── SaveSystem.cs
│   │   │   │   └── ObjectPool.cs
│   │   │   └── Data/
│   │   │       ├── GameSettings.cs
│   │   │       └── PlayerData.cs
│   │   ├── ScriptableObjects/
│   │   │   ├── Items/
│   │   │   ├── Enemies/
│   │   │   └── Settings/
│   │   └── Settings/
│   │       ├── InputActions.inputactions
│   │       └── GameSettings.asset
│   ├── Plugins/                     # Third-party plugins
│   └── Resources/                   # Runtime-loaded assets
├── Packages/
│   └── manifest.json
├── ProjectSettings/
│   ├── ProjectSettings.asset
│   ├── InputManager.asset
│   └── TagManager.asset
├── Tests/
│   ├── EditMode/
│   └── PlayMode/
└── .gitignore
```

---

## Dependencies (Package Manager)

### manifest.json
```json
{
  "dependencies": {
    "com.unity.2d.sprite": "1.0.0",
    "com.unity.2d.tilemap": "1.0.0",
    "com.unity.cinemachine": "2.9.7",
    "com.unity.inputsystem": "1.7.0",
    "com.unity.textmeshpro": "3.0.6",
    "com.unity.addressables": "1.21.19",
    "com.unity.localization": "1.4.5",
    "com.unity.ai.navigation": "1.1.5",
    "com.unity.test-framework": "1.3.9",
    "com.unity.ide.rider": "3.0.27"
  }
}
```

---

## Core Components

### MonoBehaviour Lifecycle
```csharp
using UnityEngine;

namespace MyGame.Core
{
    /// <summary>
    /// Demonstrates the Unity MonoBehaviour lifecycle.
    /// Methods are called in the order shown.
    /// </summary>
    public class LifecycleExample : MonoBehaviour
    {
        // === INITIALIZATION ===

        // Called when script instance is loaded (even if disabled)
        private void Awake()
        {
            // Initialize references, don't access other objects
            Debug.Log("Awake");
        }

        // Called before first Update (only if enabled)
        private void Start()
        {
            // Safe to access other objects initialized in Awake
            Debug.Log("Start");
        }

        // Called when object becomes enabled
        private void OnEnable()
        {
            // Subscribe to events
            Debug.Log("OnEnable");
        }

        // === UPDATE LOOP ===

        // Called every frame
        private void Update()
        {
            // Game logic, input handling
            // Time.deltaTime for frame-independent movement
        }

        // Called every fixed timestep (physics)
        private void FixedUpdate()
        {
            // Physics calculations, rigidbody movement
            // Time.fixedDeltaTime is constant
        }

        // Called after all Update calls
        private void LateUpdate()
        {
            // Camera follow, post-processing
        }

        // === CLEANUP ===

        // Called when object becomes disabled
        private void OnDisable()
        {
            // Unsubscribe from events
            Debug.Log("OnDisable");
        }

        // Called when object is destroyed
        private void OnDestroy()
        {
            // Final cleanup
            Debug.Log("OnDestroy");
        }
    }
}
```

### Game Manager (Singleton Pattern)
```csharp
using UnityEngine;
using UnityEngine.SceneManagement;

namespace MyGame.Core
{
    public enum GameState
    {
        MainMenu,
        Playing,
        Paused,
        GameOver
    }

    public class GameManager : MonoBehaviour
    {
        public static GameManager Instance { get; private set; }

        [Header("Settings")]
        [SerializeField] private GameSettings settings;

        public GameState CurrentState { get; private set; } = GameState.MainMenu;
        public int Score { get; private set; }
        public int HighScore { get; private set; }

        public event System.Action<GameState> OnGameStateChanged;
        public event System.Action<int> OnScoreChanged;

        private void Awake()
        {
            if (Instance != null && Instance != this)
            {
                Destroy(gameObject);
                return;
            }

            Instance = this;
            DontDestroyOnLoad(gameObject);

            LoadHighScore();
        }

        public void SetState(GameState newState)
        {
            if (CurrentState == newState) return;

            CurrentState = newState;
            OnGameStateChanged?.Invoke(newState);

            switch (newState)
            {
                case GameState.Playing:
                    Time.timeScale = 1f;
                    break;
                case GameState.Paused:
                    Time.timeScale = 0f;
                    break;
                case GameState.GameOver:
                    Time.timeScale = 0f;
                    SaveHighScore();
                    break;
            }
        }

        public void AddScore(int points)
        {
            Score += points;
            OnScoreChanged?.Invoke(Score);

            if (Score > HighScore)
            {
                HighScore = Score;
            }
        }

        public void ResetGame()
        {
            Score = 0;
            OnScoreChanged?.Invoke(Score);
        }

        public void LoadScene(string sceneName)
        {
            SceneManager.LoadSceneAsync(sceneName);
        }

        public void QuitGame()
        {
            SaveHighScore();

            #if UNITY_EDITOR
            UnityEditor.EditorApplication.isPlaying = false;
            #else
            Application.Quit();
            #endif
        }

        private void LoadHighScore()
        {
            HighScore = PlayerPrefs.GetInt("HighScore", 0);
        }

        private void SaveHighScore()
        {
            PlayerPrefs.SetInt("HighScore", HighScore);
            PlayerPrefs.Save();
        }
    }
}
```

### Service Locator Pattern
```csharp
using System;
using System.Collections.Generic;
using UnityEngine;

namespace MyGame.Core
{
    /// <summary>
    /// Service Locator for dependency management.
    /// Alternative to Singleton pattern.
    /// </summary>
    public static class ServiceLocator
    {
        private static readonly Dictionary<Type, object> Services = new();

        public static void Register<T>(T service) where T : class
        {
            var type = typeof(T);
            if (Services.ContainsKey(type))
            {
                Debug.LogWarning($"Service {type.Name} already registered. Replacing.");
            }
            Services[type] = service;
        }

        public static T Get<T>() where T : class
        {
            var type = typeof(T);
            if (Services.TryGetValue(type, out var service))
            {
                return service as T;
            }

            Debug.LogError($"Service {type.Name} not found!");
            return null;
        }

        public static bool TryGet<T>(out T service) where T : class
        {
            var type = typeof(T);
            if (Services.TryGetValue(type, out var obj))
            {
                service = obj as T;
                return true;
            }

            service = null;
            return false;
        }

        public static void Unregister<T>() where T : class
        {
            Services.Remove(typeof(T));
        }

        public static void Clear()
        {
            Services.Clear();
        }
    }

    // Usage: Register services in bootstrapper
    public class GameBootstrapper : MonoBehaviour
    {
        [SerializeField] private AudioManager audioManager;
        [SerializeField] private SaveSystem saveSystem;

        private void Awake()
        {
            ServiceLocator.Register<IAudioManager>(audioManager);
            ServiceLocator.Register<ISaveSystem>(saveSystem);
        }

        private void OnDestroy()
        {
            ServiceLocator.Clear();
        }
    }
}
```

---

## Player Controller

### Input System Setup
```csharp
using UnityEngine;
using UnityEngine.InputSystem;

namespace MyGame.Player
{
    [RequireComponent(typeof(CharacterController))]
    public class PlayerController : MonoBehaviour
    {
        [Header("Movement")]
        [SerializeField] private float moveSpeed = 5f;
        [SerializeField] private float sprintMultiplier = 1.5f;
        [SerializeField] private float jumpHeight = 2f;
        [SerializeField] private float gravity = -15f;

        [Header("Look")]
        [SerializeField] private float lookSensitivity = 1f;
        [SerializeField] private float maxLookAngle = 85f;
        [SerializeField] private Transform cameraTransform;

        private CharacterController controller;
        private PlayerInput playerInput;

        private Vector2 moveInput;
        private Vector2 lookInput;
        private bool sprintInput;
        private bool jumpInput;

        private Vector3 velocity;
        private float xRotation;
        private bool isGrounded;

        private void Awake()
        {
            controller = GetComponent<CharacterController>();
            playerInput = GetComponent<PlayerInput>();

            Cursor.lockState = CursorLockMode.Locked;
            Cursor.visible = false;
        }

        private void OnEnable()
        {
            // Subscribe to Input System events
            var actionMap = playerInput.actions;
            actionMap["Move"].performed += OnMove;
            actionMap["Move"].canceled += OnMove;
            actionMap["Look"].performed += OnLook;
            actionMap["Look"].canceled += OnLook;
            actionMap["Jump"].performed += OnJump;
            actionMap["Sprint"].performed += OnSprint;
            actionMap["Sprint"].canceled += OnSprint;
        }

        private void OnDisable()
        {
            var actionMap = playerInput.actions;
            actionMap["Move"].performed -= OnMove;
            actionMap["Move"].canceled -= OnMove;
            actionMap["Look"].performed -= OnLook;
            actionMap["Look"].canceled -= OnLook;
            actionMap["Jump"].performed -= OnJump;
            actionMap["Sprint"].performed -= OnSprint;
            actionMap["Sprint"].canceled -= OnSprint;
        }

        private void Update()
        {
            HandleMovement();
            HandleLook();
        }

        private void HandleMovement()
        {
            isGrounded = controller.isGrounded;

            if (isGrounded && velocity.y < 0)
            {
                velocity.y = -2f; // Small downward force to keep grounded
            }

            // Calculate move direction relative to camera
            var forward = transform.forward;
            var right = transform.right;
            forward.y = 0f;
            right.y = 0f;
            forward.Normalize();
            right.Normalize();

            var currentSpeed = sprintInput ? moveSpeed * sprintMultiplier : moveSpeed;
            var moveDirection = (forward * moveInput.y + right * moveInput.x) * currentSpeed;

            // Apply gravity
            velocity.y += gravity * Time.deltaTime;

            // Jump
            if (jumpInput && isGrounded)
            {
                velocity.y = Mathf.Sqrt(jumpHeight * -2f * gravity);
                jumpInput = false;
            }

            // Move
            controller.Move((moveDirection + velocity) * Time.deltaTime);
        }

        private void HandleLook()
        {
            // Horizontal rotation (body)
            transform.Rotate(Vector3.up * lookInput.x * lookSensitivity);

            // Vertical rotation (camera only)
            xRotation -= lookInput.y * lookSensitivity;
            xRotation = Mathf.Clamp(xRotation, -maxLookAngle, maxLookAngle);
            cameraTransform.localRotation = Quaternion.Euler(xRotation, 0f, 0f);
        }

        // Input callbacks
        private void OnMove(InputAction.CallbackContext context) =>
            moveInput = context.ReadValue<Vector2>();

        private void OnLook(InputAction.CallbackContext context) =>
            lookInput = context.ReadValue<Vector2>();

        private void OnJump(InputAction.CallbackContext context) =>
            jumpInput = true;

        private void OnSprint(InputAction.CallbackContext context) =>
            sprintInput = context.ReadValueAsButton();
    }
}
```

### Player Health System
```csharp
using UnityEngine;
using UnityEngine.Events;

namespace MyGame.Player
{
    public class PlayerHealth : MonoBehaviour, IDamageable
    {
        [Header("Health")]
        [SerializeField] private int maxHealth = 100;
        [SerializeField] private float invincibilityDuration = 0.5f;

        [Header("Events")]
        public UnityEvent<int, int> OnHealthChanged; // current, max
        public UnityEvent OnDeath;
        public UnityEvent OnDamaged;

        public int CurrentHealth { get; private set; }
        public bool IsAlive => CurrentHealth > 0;
        public bool IsInvincible { get; private set; }

        private void Start()
        {
            CurrentHealth = maxHealth;
            OnHealthChanged?.Invoke(CurrentHealth, maxHealth);
        }

        public void TakeDamage(int damage)
        {
            if (!IsAlive || IsInvincible) return;

            CurrentHealth = Mathf.Max(0, CurrentHealth - damage);
            OnHealthChanged?.Invoke(CurrentHealth, maxHealth);
            OnDamaged?.Invoke();

            if (CurrentHealth <= 0)
            {
                Die();
            }
            else
            {
                StartCoroutine(InvincibilityCoroutine());
            }
        }

        public void Heal(int amount)
        {
            if (!IsAlive) return;

            CurrentHealth = Mathf.Min(maxHealth, CurrentHealth + amount);
            OnHealthChanged?.Invoke(CurrentHealth, maxHealth);
        }

        private void Die()
        {
            OnDeath?.Invoke();
            GameManager.Instance?.SetState(GameState.GameOver);
        }

        private System.Collections.IEnumerator InvincibilityCoroutine()
        {
            IsInvincible = true;
            yield return new WaitForSeconds(invincibilityDuration);
            IsInvincible = false;
        }
    }

    public interface IDamageable
    {
        void TakeDamage(int damage);
    }
}
```

---

## ScriptableObjects

### Item Data
```csharp
using UnityEngine;

namespace MyGame.Data
{
    public enum ItemType
    {
        Consumable,
        Equipment,
        Quest,
        Material
    }

    public enum Rarity
    {
        Common,
        Uncommon,
        Rare,
        Epic,
        Legendary
    }

    [CreateAssetMenu(fileName = "New Item", menuName = "Game/Items/Item Data")]
    public class ItemData : ScriptableObject
    {
        [Header("Basic Info")]
        public string itemName;
        [TextArea(3, 5)]
        public string description;
        public Sprite icon;
        public ItemType itemType;
        public Rarity rarity;

        [Header("Properties")]
        public int maxStack = 99;
        public int buyPrice;
        public int sellPrice;

        [Header("Effects")]
        public int healthRestore;
        public int damageBonus;
        public int defenseBonus;

        public Color GetRarityColor()
        {
            return rarity switch
            {
                Rarity.Common => Color.white,
                Rarity.Uncommon => Color.green,
                Rarity.Rare => Color.blue,
                Rarity.Epic => new Color(0.5f, 0f, 0.5f), // Purple
                Rarity.Legendary => new Color(1f, 0.5f, 0f), // Orange
                _ => Color.white
            };
        }
    }
}
```

### Enemy Configuration
```csharp
using UnityEngine;

namespace MyGame.Data
{
    [CreateAssetMenu(fileName = "New Enemy", menuName = "Game/Enemies/Enemy Config")]
    public class EnemyConfig : ScriptableObject
    {
        [Header("Identity")]
        public string enemyName;
        public GameObject prefab;

        [Header("Stats")]
        public int maxHealth = 50;
        public int damage = 10;
        public float moveSpeed = 3f;
        public float attackRange = 2f;
        public float attackCooldown = 1f;

        [Header("AI")]
        public float detectionRange = 10f;
        public float chaseRange = 15f;
        public bool canPatrol = true;
        public float patrolWaitTime = 2f;

        [Header("Rewards")]
        public int experienceReward = 10;
        public int scoreReward = 100;
        public ItemData[] possibleDrops;
        [Range(0f, 1f)]
        public float dropChance = 0.3f;
    }
}
```

### Game Events (Observer Pattern)
```csharp
using System.Collections.Generic;
using UnityEngine;

namespace MyGame.Data
{
    // Parameterless event
    [CreateAssetMenu(fileName = "New Game Event", menuName = "Game/Events/Game Event")]
    public class GameEvent : ScriptableObject
    {
        private readonly List<GameEventListener> listeners = new();

        public void Raise()
        {
            for (int i = listeners.Count - 1; i >= 0; i--)
            {
                listeners[i].OnEventRaised();
            }
        }

        public void RegisterListener(GameEventListener listener)
        {
            if (!listeners.Contains(listener))
                listeners.Add(listener);
        }

        public void UnregisterListener(GameEventListener listener)
        {
            listeners.Remove(listener);
        }
    }

    // Generic event with data
    public abstract class GameEvent<T> : ScriptableObject
    {
        private readonly List<IGameEventListener<T>> listeners = new();

        public void Raise(T value)
        {
            for (int i = listeners.Count - 1; i >= 0; i--)
            {
                listeners[i].OnEventRaised(value);
            }
        }

        public void RegisterListener(IGameEventListener<T> listener)
        {
            if (!listeners.Contains(listener))
                listeners.Add(listener);
        }

        public void UnregisterListener(IGameEventListener<T> listener)
        {
            listeners.Remove(listener);
        }
    }

    public interface IGameEventListener<T>
    {
        void OnEventRaised(T value);
    }

    [CreateAssetMenu(fileName = "New Int Event", menuName = "Game/Events/Int Event")]
    public class IntGameEvent : GameEvent<int> { }

    [CreateAssetMenu(fileName = "New String Event", menuName = "Game/Events/String Event")]
    public class StringGameEvent : GameEvent<string> { }
}
```

### Event Listener Component
```csharp
using UnityEngine;
using UnityEngine.Events;

namespace MyGame.Data
{
    public class GameEventListener : MonoBehaviour
    {
        [SerializeField] private GameEvent gameEvent;
        [SerializeField] private UnityEvent response;

        private void OnEnable()
        {
            gameEvent?.RegisterListener(this);
        }

        private void OnDisable()
        {
            gameEvent?.UnregisterListener(this);
        }

        public void OnEventRaised()
        {
            response?.Invoke();
        }
    }
}
```

---

## Object Pooling

```csharp
using System.Collections.Generic;
using UnityEngine;

namespace MyGame.Systems
{
    public class ObjectPool<T> where T : Component
    {
        private readonly T prefab;
        private readonly Transform parent;
        private readonly Queue<T> pool = new();
        private readonly List<T> activeObjects = new();

        public int ActiveCount => activeObjects.Count;
        public int PoolCount => pool.Count;

        public ObjectPool(T prefab, Transform parent, int initialSize = 10)
        {
            this.prefab = prefab;
            this.parent = parent;

            for (int i = 0; i < initialSize; i++)
            {
                CreateNew();
            }
        }

        private T CreateNew()
        {
            var obj = Object.Instantiate(prefab, parent);
            obj.gameObject.SetActive(false);
            pool.Enqueue(obj);
            return obj;
        }

        public T Get()
        {
            if (pool.Count == 0)
            {
                CreateNew();
            }

            var obj = pool.Dequeue();
            obj.gameObject.SetActive(true);
            activeObjects.Add(obj);
            return obj;
        }

        public T Get(Vector3 position, Quaternion rotation)
        {
            var obj = Get();
            obj.transform.SetPositionAndRotation(position, rotation);
            return obj;
        }

        public void Return(T obj)
        {
            if (!activeObjects.Contains(obj)) return;

            obj.gameObject.SetActive(false);
            obj.transform.SetParent(parent);
            activeObjects.Remove(obj);
            pool.Enqueue(obj);
        }

        public void ReturnAll()
        {
            foreach (var obj in activeObjects.ToArray())
            {
                Return(obj);
            }
        }
    }

    // MonoBehaviour wrapper for inspector setup
    public class BulletPool : MonoBehaviour
    {
        [SerializeField] private Bullet bulletPrefab;
        [SerializeField] private int initialSize = 20;

        private ObjectPool<Bullet> pool;

        public static BulletPool Instance { get; private set; }

        private void Awake()
        {
            Instance = this;
            pool = new ObjectPool<Bullet>(bulletPrefab, transform, initialSize);
        }

        public Bullet GetBullet(Vector3 position, Quaternion rotation)
        {
            return pool.Get(position, rotation);
        }

        public void ReturnBullet(Bullet bullet)
        {
            pool.Return(bullet);
        }
    }
}
```

---

## Audio Manager

```csharp
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.Audio;

namespace MyGame.Systems
{
    public interface IAudioManager
    {
        void PlaySFX(AudioClip clip, float volume = 1f);
        void PlaySFXAtPosition(AudioClip clip, Vector3 position, float volume = 1f);
        void PlayMusic(AudioClip clip, bool loop = true);
        void StopMusic();
        void SetMasterVolume(float volume);
        void SetMusicVolume(float volume);
        void SetSFXVolume(float volume);
    }

    public class AudioManager : MonoBehaviour, IAudioManager
    {
        [Header("Audio Mixer")]
        [SerializeField] private AudioMixer audioMixer;

        [Header("Audio Sources")]
        [SerializeField] private AudioSource musicSource;
        [SerializeField] private AudioSource sfxSource;

        [Header("Pooling")]
        [SerializeField] private int sfxPoolSize = 10;

        private readonly Queue<AudioSource> sfxPool = new();
        private readonly List<AudioSource> activeSfxSources = new();

        private void Awake()
        {
            // Create SFX pool
            for (int i = 0; i < sfxPoolSize; i++)
            {
                var source = gameObject.AddComponent<AudioSource>();
                source.playOnAwake = false;
                source.outputAudioMixerGroup = sfxSource.outputAudioMixerGroup;
                sfxPool.Enqueue(source);
            }

            LoadVolumeSettings();
        }

        public void PlaySFX(AudioClip clip, float volume = 1f)
        {
            if (clip == null) return;
            sfxSource.PlayOneShot(clip, volume);
        }

        public void PlaySFXAtPosition(AudioClip clip, Vector3 position, float volume = 1f)
        {
            if (clip == null) return;

            var source = GetPooledSource();
            if (source == null)
            {
                AudioSource.PlayClipAtPoint(clip, position, volume);
                return;
            }

            source.transform.position = position;
            source.clip = clip;
            source.volume = volume;
            source.spatialBlend = 1f; // 3D sound
            source.Play();

            StartCoroutine(ReturnToPoolAfterPlay(source, clip.length));
        }

        public void PlayMusic(AudioClip clip, bool loop = true)
        {
            if (clip == null) return;

            musicSource.clip = clip;
            musicSource.loop = loop;
            musicSource.Play();
        }

        public void StopMusic()
        {
            musicSource.Stop();
        }

        public void SetMasterVolume(float volume)
        {
            audioMixer.SetFloat("MasterVolume", LinearToDecibel(volume));
            PlayerPrefs.SetFloat("MasterVolume", volume);
        }

        public void SetMusicVolume(float volume)
        {
            audioMixer.SetFloat("MusicVolume", LinearToDecibel(volume));
            PlayerPrefs.SetFloat("MusicVolume", volume);
        }

        public void SetSFXVolume(float volume)
        {
            audioMixer.SetFloat("SFXVolume", LinearToDecibel(volume));
            PlayerPrefs.SetFloat("SFXVolume", volume);
        }

        private AudioSource GetPooledSource()
        {
            if (sfxPool.Count == 0) return null;

            var source = sfxPool.Dequeue();
            activeSfxSources.Add(source);
            return source;
        }

        private System.Collections.IEnumerator ReturnToPoolAfterPlay(AudioSource source, float delay)
        {
            yield return new WaitForSeconds(delay);

            source.Stop();
            source.clip = null;
            activeSfxSources.Remove(source);
            sfxPool.Enqueue(source);
        }

        private void LoadVolumeSettings()
        {
            SetMasterVolume(PlayerPrefs.GetFloat("MasterVolume", 1f));
            SetMusicVolume(PlayerPrefs.GetFloat("MusicVolume", 1f));
            SetSFXVolume(PlayerPrefs.GetFloat("SFXVolume", 1f));
        }

        private float LinearToDecibel(float linear)
        {
            return linear > 0 ? Mathf.Log10(linear) * 20f : -80f;
        }
    }
}
```

---

## Save System

```csharp
using System;
using System.IO;
using UnityEngine;

namespace MyGame.Systems
{
    public interface ISaveSystem
    {
        void Save<T>(string key, T data);
        T Load<T>(string key, T defaultValue = default);
        bool HasSave(string key);
        void DeleteSave(string key);
        void DeleteAllSaves();
    }

    [Serializable]
    public class SaveData
    {
        public int level;
        public int experience;
        public int gold;
        public float playTime;
        public Vector3Serializable playerPosition;
        public string[] unlockedItems;
        public DateTime lastSaved;
    }

    [Serializable]
    public struct Vector3Serializable
    {
        public float x, y, z;

        public Vector3Serializable(Vector3 v)
        {
            x = v.x;
            y = v.y;
            z = v.z;
        }

        public Vector3 ToVector3() => new(x, y, z);
    }

    public class SaveSystem : MonoBehaviour, ISaveSystem
    {
        private string SavePath => Application.persistentDataPath;

        public void Save<T>(string key, T data)
        {
            try
            {
                var json = JsonUtility.ToJson(data, true);
                var path = GetFilePath(key);
                File.WriteAllText(path, json);
                Debug.Log($"Saved to {path}");
            }
            catch (Exception e)
            {
                Debug.LogError($"Failed to save {key}: {e.Message}");
            }
        }

        public T Load<T>(string key, T defaultValue = default)
        {
            var path = GetFilePath(key);

            if (!File.Exists(path))
            {
                return defaultValue;
            }

            try
            {
                var json = File.ReadAllText(path);
                return JsonUtility.FromJson<T>(json);
            }
            catch (Exception e)
            {
                Debug.LogError($"Failed to load {key}: {e.Message}");
                return defaultValue;
            }
        }

        public bool HasSave(string key)
        {
            return File.Exists(GetFilePath(key));
        }

        public void DeleteSave(string key)
        {
            var path = GetFilePath(key);
            if (File.Exists(path))
            {
                File.Delete(path);
            }
        }

        public void DeleteAllSaves()
        {
            var files = Directory.GetFiles(SavePath, "*.json");
            foreach (var file in files)
            {
                File.Delete(file);
            }
        }

        private string GetFilePath(string key)
        {
            return Path.Combine(SavePath, $"{key}.json");
        }
    }
}
```

---

## UI System

### HUD Controller
```csharp
using UnityEngine;
using UnityEngine.UI;
using TMPro;

namespace MyGame.UI
{
    public class HUDController : MonoBehaviour
    {
        [Header("Health")]
        [SerializeField] private Slider healthBar;
        [SerializeField] private TextMeshProUGUI healthText;

        [Header("Score")]
        [SerializeField] private TextMeshProUGUI scoreText;
        [SerializeField] private TextMeshProUGUI highScoreText;

        [Header("Animation")]
        [SerializeField] private Animator scoreAnimator;
        [SerializeField] private float healthLerpSpeed = 5f;

        private float targetHealthPercent;

        private void Start()
        {
            // Subscribe to events
            if (GameManager.Instance != null)
            {
                GameManager.Instance.OnScoreChanged += UpdateScore;
            }

            var playerHealth = FindObjectOfType<PlayerHealth>();
            if (playerHealth != null)
            {
                playerHealth.OnHealthChanged.AddListener(UpdateHealth);
            }

            // Initialize
            UpdateHighScore();
        }

        private void Update()
        {
            // Smooth health bar animation
            if (healthBar != null)
            {
                healthBar.value = Mathf.Lerp(
                    healthBar.value,
                    targetHealthPercent,
                    healthLerpSpeed * Time.deltaTime
                );
            }
        }

        private void UpdateHealth(int current, int max)
        {
            targetHealthPercent = (float)current / max;

            if (healthText != null)
            {
                healthText.text = $"{current}/{max}";
            }
        }

        private void UpdateScore(int score)
        {
            if (scoreText != null)
            {
                scoreText.text = score.ToString("N0");
            }

            if (scoreAnimator != null)
            {
                scoreAnimator.SetTrigger("ScoreChanged");
            }

            UpdateHighScore();
        }

        private void UpdateHighScore()
        {
            if (highScoreText != null && GameManager.Instance != null)
            {
                highScoreText.text = $"Best: {GameManager.Instance.HighScore:N0}";
            }
        }

        private void OnDestroy()
        {
            if (GameManager.Instance != null)
            {
                GameManager.Instance.OnScoreChanged -= UpdateScore;
            }
        }
    }
}
```

### Pause Menu
```csharp
using UnityEngine;
using UnityEngine.UI;
using UnityEngine.InputSystem;

namespace MyGame.UI
{
    public class PauseMenuController : MonoBehaviour
    {
        [Header("Panels")]
        [SerializeField] private GameObject pausePanel;
        [SerializeField] private GameObject settingsPanel;

        [Header("Settings")]
        [SerializeField] private Slider masterVolumeSlider;
        [SerializeField] private Slider musicVolumeSlider;
        [SerializeField] private Slider sfxVolumeSlider;

        private IAudioManager audioManager;

        private void Start()
        {
            pausePanel.SetActive(false);
            settingsPanel.SetActive(false);

            if (ServiceLocator.TryGet(out IAudioManager audio))
            {
                audioManager = audio;
            }

            LoadSettings();
        }

        public void OnPauseInput(InputAction.CallbackContext context)
        {
            if (context.performed)
            {
                TogglePause();
            }
        }

        public void TogglePause()
        {
            if (GameManager.Instance == null) return;

            if (GameManager.Instance.CurrentState == GameState.Playing)
            {
                Pause();
            }
            else if (GameManager.Instance.CurrentState == GameState.Paused)
            {
                Resume();
            }
        }

        public void Pause()
        {
            pausePanel.SetActive(true);
            settingsPanel.SetActive(false);
            GameManager.Instance.SetState(GameState.Paused);

            Cursor.lockState = CursorLockMode.None;
            Cursor.visible = true;
        }

        public void Resume()
        {
            pausePanel.SetActive(false);
            settingsPanel.SetActive(false);
            GameManager.Instance.SetState(GameState.Playing);

            Cursor.lockState = CursorLockMode.Locked;
            Cursor.visible = false;
        }

        public void OpenSettings()
        {
            pausePanel.SetActive(false);
            settingsPanel.SetActive(true);
        }

        public void CloseSettings()
        {
            settingsPanel.SetActive(false);
            pausePanel.SetActive(true);
        }

        public void ReturnToMainMenu()
        {
            Time.timeScale = 1f;
            GameManager.Instance?.LoadScene("MainMenu");
        }

        public void OnMasterVolumeChanged(float value)
        {
            audioManager?.SetMasterVolume(value);
        }

        public void OnMusicVolumeChanged(float value)
        {
            audioManager?.SetMusicVolume(value);
        }

        public void OnSFXVolumeChanged(float value)
        {
            audioManager?.SetSFXVolume(value);
        }

        private void LoadSettings()
        {
            masterVolumeSlider.value = PlayerPrefs.GetFloat("MasterVolume", 1f);
            musicVolumeSlider.value = PlayerPrefs.GetFloat("MusicVolume", 1f);
            sfxVolumeSlider.value = PlayerPrefs.GetFloat("SFXVolume", 1f);
        }
    }
}
```

---

## Testing

### Edit Mode Tests
```csharp
using NUnit.Framework;
using MyGame.Data;
using UnityEngine;

namespace MyGame.Tests.EditMode
{
    public class ItemDataTests
    {
        [Test]
        public void GetRarityColor_Common_ReturnsWhite()
        {
            var item = ScriptableObject.CreateInstance<ItemData>();
            item.rarity = Rarity.Common;

            var color = item.GetRarityColor();

            Assert.AreEqual(Color.white, color);
        }

        [Test]
        public void GetRarityColor_Legendary_ReturnsOrange()
        {
            var item = ScriptableObject.CreateInstance<ItemData>();
            item.rarity = Rarity.Legendary;

            var color = item.GetRarityColor();

            Assert.AreEqual(new Color(1f, 0.5f, 0f), color);
        }
    }

    public class SaveDataTests
    {
        [Test]
        public void Vector3Serializable_RoundTrip_PreservesValues()
        {
            var original = new Vector3(1.5f, 2.5f, 3.5f);
            var serializable = new Vector3Serializable(original);
            var result = serializable.ToVector3();

            Assert.AreEqual(original, result);
        }

        [Test]
        public void SaveData_JsonSerialization_Works()
        {
            var data = new SaveData
            {
                level = 5,
                experience = 1000,
                gold = 500,
                playTime = 3600f,
                playerPosition = new Vector3Serializable(new Vector3(10, 0, 20))
            };

            var json = JsonUtility.ToJson(data);
            var restored = JsonUtility.FromJson<SaveData>(json);

            Assert.AreEqual(data.level, restored.level);
            Assert.AreEqual(data.experience, restored.experience);
            Assert.AreEqual(data.gold, restored.gold);
            Assert.AreEqual(data.playTime, restored.playTime);
        }
    }
}
```

### Play Mode Tests
```csharp
using System.Collections;
using NUnit.Framework;
using UnityEngine;
using UnityEngine.TestTools;
using MyGame.Player;

namespace MyGame.Tests.PlayMode
{
    public class PlayerHealthTests
    {
        private GameObject playerObject;
        private PlayerHealth playerHealth;

        [SetUp]
        public void SetUp()
        {
            playerObject = new GameObject("Player");
            playerHealth = playerObject.AddComponent<PlayerHealth>();
        }

        [TearDown]
        public void TearDown()
        {
            Object.Destroy(playerObject);
        }

        [UnityTest]
        public IEnumerator TakeDamage_ReducesHealth()
        {
            yield return null; // Wait for Start()

            var initialHealth = playerHealth.CurrentHealth;
            playerHealth.TakeDamage(10);

            Assert.AreEqual(initialHealth - 10, playerHealth.CurrentHealth);
        }

        [UnityTest]
        public IEnumerator TakeDamage_WhenHealthZero_TriggersDeathEvent()
        {
            yield return null;

            bool deathTriggered = false;
            playerHealth.OnDeath.AddListener(() => deathTriggered = true);

            playerHealth.TakeDamage(1000);

            Assert.IsTrue(deathTriggered);
            Assert.IsFalse(playerHealth.IsAlive);
        }

        [UnityTest]
        public IEnumerator Heal_IncreasesHealth()
        {
            yield return null;

            playerHealth.TakeDamage(50);
            var healthAfterDamage = playerHealth.CurrentHealth;

            playerHealth.Heal(25);

            Assert.AreEqual(healthAfterDamage + 25, playerHealth.CurrentHealth);
        }

        [UnityTest]
        public IEnumerator Invincibility_PreventsConsecutiveDamage()
        {
            yield return null;

            playerHealth.TakeDamage(10);
            var healthAfterFirstHit = playerHealth.CurrentHealth;

            playerHealth.TakeDamage(10); // Should be blocked

            Assert.AreEqual(healthAfterFirstHit, playerHealth.CurrentHealth);
        }
    }
}
```

---

## Build & Commands

### Build Commands
```bash
# Build from command line
# Windows
"C:\Program Files\Unity\Hub\Editor\2022.3.0f1\Editor\Unity.exe" \
  -quit -batchmode \
  -projectPath "C:\Projects\MyGame" \
  -buildTarget Win64 \
  -buildPath "Builds/Windows/MyGame.exe"

# macOS
/Applications/Unity/Hub/Editor/2022.3.0f1/Unity.app/Contents/MacOS/Unity \
  -quit -batchmode \
  -projectPath ~/Projects/MyGame \
  -buildTarget StandaloneOSX \
  -buildPath Builds/macOS/MyGame.app

# Run tests
Unity -runTests \
  -projectPath /path/to/project \
  -testResults results.xml \
  -testPlatform EditMode

# Create package
Unity -exportPackage Assets/MyPlugin MyPlugin.unitypackage
```

### Editor Scripts
```csharp
#if UNITY_EDITOR
using UnityEditor;
using UnityEngine;

namespace MyGame.Editor
{
    public static class BuildScript
    {
        [MenuItem("Build/Build Windows")]
        public static void BuildWindows()
        {
            var scenes = new[]
            {
                "Assets/_Project/Scenes/MainMenu.unity",
                "Assets/_Project/Scenes/GameLevel.unity"
            };

            BuildPipeline.BuildPlayer(
                scenes,
                "Builds/Windows/MyGame.exe",
                BuildTarget.StandaloneWindows64,
                BuildOptions.None
            );
        }

        [MenuItem("Build/Build WebGL")]
        public static void BuildWebGL()
        {
            var scenes = new[]
            {
                "Assets/_Project/Scenes/MainMenu.unity",
                "Assets/_Project/Scenes/GameLevel.unity"
            };

            BuildPipeline.BuildPlayer(
                scenes,
                "Builds/WebGL",
                BuildTarget.WebGL,
                BuildOptions.None
            );
        }
    }
}
#endif
```

---

## Best Practices

### Performance Tips
1. **Object Pooling**: Reuse objects instead of Instantiate/Destroy
2. **Avoid Find**: Cache references in Awake/Start
3. **Update Optimization**: Use events instead of polling in Update
4. **String Operations**: Use StringBuilder, avoid concatenation
5. **Physics**: Use layers and avoid complex colliders
6. **Draw Calls**: Batch materials, use sprite atlases

### Code Guidelines
1. **Single Responsibility**: One component = one purpose
2. **Composition over Inheritance**: Favor multiple components
3. **ScriptableObjects**: Use for configuration and data
4. **Events**: Decouple systems with UnityEvents or C# events
5. **Serialization**: Mark fields [SerializeField] for inspector access

### Common Mistakes to Avoid
```csharp
// BAD: Finding objects every frame
void Update()
{
    var player = GameObject.Find("Player"); // Expensive!
}

// GOOD: Cache references
private Transform player;
void Start() => player = GameObject.Find("Player").transform;

// BAD: String comparison
if (gameObject.tag == "Player") { }

// GOOD: Use CompareTag
if (gameObject.CompareTag("Player")) { }

// BAD: Allocating in Update
void Update()
{
    var list = new List<Enemy>(); // GC allocation every frame!
}

// GOOD: Reuse collections
private readonly List<Enemy> enemies = new();
void Update()
{
    enemies.Clear();
    // Reuse...
}
```

---

## Framework Comparison

| Feature | Unity | Unreal Engine | Godot |
|---------|-------|---------------|-------|
| Language | C# | C++/Blueprints | GDScript/C# |
| Learning Curve | Medium | Steep | Easy |
| 2D Support | Excellent | Good | Excellent |
| 3D Support | Excellent | Excellent | Good |
| Mobile | Excellent | Good | Good |
| VR/AR | Excellent | Excellent | Limited |
| Performance | Good | Excellent | Good |
| Asset Store | Massive | Large | Growing |
| Open Source | No | Partial | Yes |
| Free Tier | Yes (< $100K) | Yes (< $1M) | Yes |

---

## References

- [Unity Documentation](https://docs.unity3d.com/)
- [Unity Learn](https://learn.unity.com/)
- [Unity Manual](https://docs.unity3d.com/Manual/)
- [Unity Scripting API](https://docs.unity3d.com/ScriptReference/)
- [Unity Best Practices](https://unity.com/how-to)
- [Game Programming Patterns](https://gameprogrammingpatterns.com/)
- [Unity Community Forums](https://forum.unity.com/)
