package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_GetValue(t *testing.T) {
	config := &Config{
		Version:  "1.0.0",
		Registry: "https://example.com",
		Installed: InstalledItems{
			Languages:  []string{"go", "python"},
			Frameworks: []string{"react"},
			Workflows:  []string{"create-prd"},
		},
	}

	tests := []struct {
		key     string
		want    interface{}
		wantErr bool
	}{
		{"version", "1.0.0", false},
		{"registry", "https://example.com", false},
		{"installed.languages", []string{"go", "python"}, false},
		{"installed.frameworks", []string{"react"}, false},
		{"installed.workflows", []string{"create-prd"}, false},
		{"invalid", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, err := config.GetValue(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValue(%q) error = %v, wantErr %v", tt.key, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				switch v := tt.want.(type) {
				case string:
					if got != v {
						t.Errorf("GetValue(%q) = %v, want %v", tt.key, got, v)
					}
				case []string:
					gotSlice, ok := got.([]string)
					if !ok {
						t.Errorf("GetValue(%q) returned non-slice type", tt.key)
						return
					}
					if len(gotSlice) != len(v) {
						t.Errorf("GetValue(%q) = %v, want %v", tt.key, gotSlice, v)
					}
				}
			}
		})
	}
}

func TestConfig_GetValue_DefaultRegistry(t *testing.T) {
	config := &Config{
		Version:  "1.0.0",
		Registry: "", // Empty registry
	}

	got, err := config.GetValue("registry")
	if err != nil {
		t.Errorf("GetValue(registry) unexpected error: %v", err)
	}
	if got != DefaultRegistry {
		t.Errorf("GetValue(registry) = %v, want %v (default)", got, DefaultRegistry)
	}
}

func TestConfig_SetValue(t *testing.T) {
	_ = &Config{
		Version:  "1.0.0",
		Registry: "",
	}

	tests := []struct {
		key     string
		value   string
		wantErr bool
		check   func(*Config) bool
	}{
		{
			key:     "version",
			value:   "2.0.0",
			wantErr: false,
			check:   func(c *Config) bool { return c.Version == "2.0.0" },
		},
		{
			key:     "registry",
			value:   "https://new.example.com",
			wantErr: false,
			check:   func(c *Config) bool { return c.Registry == "https://new.example.com" },
		},
		{
			key:     "installed.languages",
			value:   "go,python,rust",
			wantErr: false,
			check: func(c *Config) bool {
				return len(c.Installed.Languages) == 3 &&
					c.Installed.Languages[0] == "go" &&
					c.Installed.Languages[1] == "python" &&
					c.Installed.Languages[2] == "rust"
			},
		},
		{
			key:     "installed.frameworks",
			value:   "react,nextjs",
			wantErr: false,
			check: func(c *Config) bool {
				return len(c.Installed.Frameworks) == 2
			},
		},
		{
			key:     "installed.workflows",
			value:   "all",
			wantErr: false,
			check:   func(c *Config) bool { return c.Installed.Workflows[0] == "all" },
		},
		{
			key:     "invalid",
			value:   "test",
			wantErr: true,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			// Create fresh config for each test
			cfg := &Config{Version: "1.0.0"}
			err := cfg.SetValue(tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetValue(%q, %q) error = %v, wantErr %v", tt.key, tt.value, err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil && !tt.check(cfg) {
				t.Errorf("SetValue(%q, %q) did not update config correctly", tt.key, tt.value)
			}
		})
	}
}

func TestConfig_GetAllValues(t *testing.T) {
	config := &Config{
		Version:  "1.0.0",
		Registry: "https://example.com",
		Installed: InstalledItems{
			Languages:  []string{"go"},
			Frameworks: []string{"react"},
			Workflows:  []string{"all"},
		},
	}

	values := config.GetAllValues()

	if values["version"] != "1.0.0" {
		t.Errorf("GetAllValues()[version] = %v, want 1.0.0", values["version"])
	}
	if values["registry"] != "https://example.com" {
		t.Errorf("GetAllValues()[registry] = %v, want https://example.com", values["registry"])
	}

	langs, ok := values["installed.languages"].([]string)
	if !ok || len(langs) != 1 || langs[0] != "go" {
		t.Errorf("GetAllValues()[installed.languages] = %v, want [go]", values["installed.languages"])
	}
}

func TestConfig_GetAllValues_DefaultRegistry(t *testing.T) {
	config := &Config{
		Version:  "1.0.0",
		Registry: "", // Empty
	}

	values := config.GetAllValues()
	if values["registry"] != DefaultRegistry {
		t.Errorf("GetAllValues()[registry] = %v, want %v (default)", values["registry"], DefaultRegistry)
	}
}

func TestValidConfigKeys(t *testing.T) {
	expectedKeys := []string{
		"version",
		"registry",
		"installed.languages",
		"installed.frameworks",
		"installed.workflows",
		"installed.skills",
	}

	if len(ValidConfigKeys) != len(expectedKeys) {
		t.Errorf("ValidConfigKeys has %d items, want %d", len(ValidConfigKeys), len(expectedKeys))
	}

	for _, expected := range expectedKeys {
		found := false
		for _, key := range ValidConfigKeys {
			if key == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ValidConfigKeys missing %q", expected)
		}
	}
}

func TestConfig_HasLanguage(t *testing.T) {
	config := &Config{
		Installed: InstalledItems{
			Languages: []string{"go", "python"},
		},
	}

	if !config.HasLanguage("go") {
		t.Error("HasLanguage(go) = false, want true")
	}
	if !config.HasLanguage("python") {
		t.Error("HasLanguage(python) = false, want true")
	}
	if config.HasLanguage("rust") {
		t.Error("HasLanguage(rust) = true, want false")
	}
}

func TestConfig_HasFramework(t *testing.T) {
	config := &Config{
		Installed: InstalledItems{
			Frameworks: []string{"react", "nextjs"},
		},
	}

	if !config.HasFramework("react") {
		t.Error("HasFramework(react) = false, want true")
	}
	if config.HasFramework("vue") {
		t.Error("HasFramework(vue) = true, want false")
	}
}

func TestConfig_HasWorkflow(t *testing.T) {
	config := &Config{
		Installed: InstalledItems{
			Workflows: []string{"all"},
		},
	}

	// "all" should match any workflow
	if !config.HasWorkflow("create-prd") {
		t.Error("HasWorkflow(create-prd) with 'all' = false, want true")
	}

	config.Installed.Workflows = []string{"create-prd", "code-review"}
	if !config.HasWorkflow("create-prd") {
		t.Error("HasWorkflow(create-prd) = false, want true")
	}
	if config.HasWorkflow("security-audit") {
		t.Error("HasWorkflow(security-audit) = true, want false")
	}
}

func TestLoadConfig_NotExist(t *testing.T) {
	// Change to temp directory
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}
	defer func() { _ = os.Chdir(oldDir) }()

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() should return error when config doesn't exist")
	}
}

func TestLoadConfig_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}
	defer func() { _ = os.Chdir(oldDir) }()

	// Create a valid config file
	configContent := `version: "1.0.0"
registry: "https://github.com/ar4mirez/aicof"
installed:
  languages:
    - go
  frameworks:
    - react
  workflows:
    - all
`
	err := os.WriteFile(filepath.Join(tmpDir, "aicof.yaml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if config.Version != "1.0.0" {
		t.Errorf("config.Version = %q, want %q", config.Version, "1.0.0")
	}
	if len(config.Installed.Languages) != 1 || config.Installed.Languages[0] != "go" {
		t.Errorf("config.Installed.Languages = %v, want [go]", config.Installed.Languages)
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}
	defer func() { _ = os.Chdir(oldDir) }()

	config := &Config{
		Version:  "2.0.0",
		Registry: "https://example.com",
		Installed: InstalledItems{
			Languages:  []string{"rust"},
			Frameworks: []string{"axum"},
			Workflows:  []string{"all"},
		},
	}

	err := config.Save(tmpDir)
	if err != nil {
		t.Fatalf("config.Save() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filepath.Join(tmpDir, "aicof.yaml")); os.IsNotExist(err) {
		t.Error("SaveConfig() did not create aicof.yaml")
	}

	// Load and verify
	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() after save error = %v", err)
	}
	if loaded.Version != "2.0.0" {
		t.Errorf("loaded.Version = %q, want %q", loaded.Version, "2.0.0")
	}
}

func TestDefaultRegistry(t *testing.T) {
	if DefaultRegistry == "" {
		t.Error("DefaultRegistry should not be empty")
	}
}

func TestNewConfig(t *testing.T) {
	config := NewConfig("1.0.0")

	if config.Version != "1.0.0" {
		t.Errorf("NewConfig() Version = %q, want %q", config.Version, "1.0.0")
	}
	if config.Registry != DefaultRegistry {
		t.Errorf("NewConfig() Registry = %q, want %q", config.Registry, DefaultRegistry)
	}
	if len(config.Installed.Languages) != 0 {
		t.Errorf("NewConfig() Languages should be empty, got %v", config.Installed.Languages)
	}
	if len(config.Installed.Frameworks) != 0 {
		t.Errorf("NewConfig() Frameworks should be empty, got %v", config.Installed.Frameworks)
	}
	if len(config.Installed.Workflows) != 1 || config.Installed.Workflows[0] != "all" {
		t.Errorf("NewConfig() Workflows = %v, want [all]", config.Installed.Workflows)
	}
}

func TestConfigExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Should return false for non-existent config
	if ConfigExists(tmpDir) {
		t.Error("ConfigExists() = true for empty dir, want false")
	}

	// Create aicof.yaml
	if err := os.WriteFile(filepath.Join(tmpDir, "aicof.yaml"), []byte("version: 1.0.0"), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
	if !ConfigExists(tmpDir) {
		t.Error("ConfigExists() = false with aicof.yaml, want true")
	}

	// Test with alternative config file
	tmpDir2 := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir2, ".aicof.yaml"), []byte("version: 1.0.0"), 0644); err != nil {
		t.Fatalf("Failed to write alt config file: %v", err)
	}
	if !ConfigExists(tmpDir2) {
		t.Error("ConfigExists() = false with .aicof.yaml, want true")
	}
}

func TestConfig_AddLanguage(t *testing.T) {
	config := &Config{}

	config.AddLanguage("go")
	if len(config.Installed.Languages) != 1 || config.Installed.Languages[0] != "go" {
		t.Errorf("AddLanguage() = %v, want [go]", config.Installed.Languages)
	}

	// AddLanguage should also add the corresponding skill
	if !config.HasSkill("go-guide") {
		t.Error("AddLanguage(go) should also add go-guide skill")
	}

	// Adding duplicate should not add again
	config.AddLanguage("go")
	if len(config.Installed.Languages) != 1 {
		t.Errorf("AddLanguage() should not add duplicate, got %v", config.Installed.Languages)
	}

	config.AddLanguage("python")
	if len(config.Installed.Languages) != 2 {
		t.Errorf("AddLanguage() should have 2 languages, got %v", config.Installed.Languages)
	}
	if !config.HasSkill("python-guide") {
		t.Error("AddLanguage(python) should also add python-guide skill")
	}
}

func TestConfig_AddFramework(t *testing.T) {
	config := &Config{}

	config.AddFramework("react")
	if len(config.Installed.Frameworks) != 1 || config.Installed.Frameworks[0] != "react" {
		t.Errorf("AddFramework() = %v, want [react]", config.Installed.Frameworks)
	}

	// AddFramework should also add the corresponding skill
	if !config.HasSkill("react") {
		t.Error("AddFramework(react) should also add react skill")
	}

	// Adding duplicate should not add again
	config.AddFramework("react")
	if len(config.Installed.Frameworks) != 1 {
		t.Errorf("AddFramework() should not add duplicate, got %v", config.Installed.Frameworks)
	}

	config.AddFramework("django")
	if len(config.Installed.Frameworks) != 2 {
		t.Errorf("AddFramework() should have 2 frameworks, got %v", config.Installed.Frameworks)
	}
	if !config.HasSkill("django") {
		t.Error("AddFramework(django) should also add django skill")
	}
}

func TestConfig_AddWorkflow(t *testing.T) {
	config := &Config{}

	config.AddWorkflow("create-prd")
	if len(config.Installed.Workflows) != 1 || config.Installed.Workflows[0] != "create-prd" {
		t.Errorf("AddWorkflow() = %v, want [create-prd]", config.Installed.Workflows)
	}

	// Adding duplicate should not add again
	config.AddWorkflow("create-prd")
	if len(config.Installed.Workflows) != 1 {
		t.Errorf("AddWorkflow() should not add duplicate, got %v", config.Installed.Workflows)
	}

	// Test with "all" - should not add since "all" matches everything
	config.Installed.Workflows = []string{"all"}
	config.AddWorkflow("code-review")
	if len(config.Installed.Workflows) != 1 {
		t.Errorf("AddWorkflow() with 'all' should not add, got %v", config.Installed.Workflows)
	}
}

func TestConfig_RemoveLanguage(t *testing.T) {
	config := &Config{
		Installed: InstalledItems{
			Languages: []string{"go", "python", "rust"},
			Skills:    []string{"go-guide", "python-guide", "rust-guide"},
		},
	}

	config.RemoveLanguage("python")
	if len(config.Installed.Languages) != 2 {
		t.Errorf("RemoveLanguage() should have 2 languages, got %v", config.Installed.Languages)
	}
	if config.HasLanguage("python") {
		t.Error("RemoveLanguage() should have removed python")
	}
	if config.HasSkill("python-guide") {
		t.Error("RemoveLanguage() should also remove python-guide skill")
	}

	// Removing non-existent should not error
	config.RemoveLanguage("nonexistent")
	if len(config.Installed.Languages) != 2 {
		t.Errorf("RemoveLanguage() should still have 2 languages, got %v", config.Installed.Languages)
	}
}

func TestConfig_RemoveFramework(t *testing.T) {
	config := &Config{
		Installed: InstalledItems{
			Frameworks: []string{"react", "nextjs"},
			Skills:     []string{"react", "nextjs"},
		},
	}

	config.RemoveFramework("react")
	if len(config.Installed.Frameworks) != 1 {
		t.Errorf("RemoveFramework() should have 1 framework, got %v", config.Installed.Frameworks)
	}
	if config.HasFramework("react") {
		t.Error("RemoveFramework() should have removed react")
	}
	if config.HasSkill("react") {
		t.Error("RemoveFramework() should also remove react skill")
	}
}

func TestConfig_RemoveWorkflow(t *testing.T) {
	config := &Config{
		Installed: InstalledItems{
			Workflows: []string{"create-prd", "code-review"},
		},
	}

	config.RemoveWorkflow("create-prd")
	if len(config.Installed.Workflows) != 1 {
		t.Errorf("RemoveWorkflow() should have 1 workflow, got %v", config.Installed.Workflows)
	}
}

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"", []string{}},
		{"go", []string{"go"}},
		{"go,python,rust", []string{"go", "python", "rust"}},
		{"go, python, rust", []string{"go", "python", "rust"}},
		{"  go  ,  python  ", []string{"go", "python"}},
		{",,,", []string{}}, // Empty values
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			config := &Config{}
			if err := config.SetValue("installed.languages", tt.input); err != nil {
				t.Fatalf("SetValue failed: %v", err)
			}
			if len(config.Installed.Languages) != len(tt.want) {
				t.Errorf("splitAndTrim(%q) got %v, want %v", tt.input, config.Installed.Languages, tt.want)
			}
		})
	}
}

func TestConfig_MigrateLanguagesToSkills(t *testing.T) {
	// Simulate a legacy config with languages but no corresponding skills
	config := &Config{
		Installed: InstalledItems{
			Languages: []string{"go", "python"},
			Skills:    []string{"commit-message"},
		},
	}

	migrated := config.MigrateLanguagesToSkills()
	if !migrated {
		t.Error("MigrateLanguagesToSkills() should return true when migration is needed")
	}
	if !config.HasSkill("go-guide") {
		t.Error("MigrateLanguagesToSkills() should add go-guide skill")
	}
	if !config.HasSkill("python-guide") {
		t.Error("MigrateLanguagesToSkills() should add python-guide skill")
	}
	if !config.HasSkill("commit-message") {
		t.Error("MigrateLanguagesToSkills() should preserve existing skills")
	}

	// Running again should return false (no new migration)
	migrated = config.MigrateLanguagesToSkills()
	if migrated {
		t.Error("MigrateLanguagesToSkills() should return false when already migrated")
	}
}

func TestConfig_MigrateFrameworksToSkills(t *testing.T) {
	// Simulate a legacy config with frameworks but no corresponding skills
	config := &Config{
		Installed: InstalledItems{
			Frameworks: []string{"react", "django"},
			Skills:     []string{"commit-message"},
		},
	}

	migrated := config.MigrateFrameworksToSkills()
	if !migrated {
		t.Error("MigrateFrameworksToSkills() should return true when migration is needed")
	}
	if !config.HasSkill("react") {
		t.Error("MigrateFrameworksToSkills() should add react skill")
	}
	if !config.HasSkill("django") {
		t.Error("MigrateFrameworksToSkills() should add django skill")
	}
	if !config.HasSkill("commit-message") {
		t.Error("MigrateFrameworksToSkills() should preserve existing skills")
	}

	// Running again should return false (no new migration)
	migrated = config.MigrateFrameworksToSkills()
	if migrated {
		t.Error("MigrateFrameworksToSkills() should return false when already migrated")
	}
}

func TestConfig_MigrateLanguagesToSkills_Empty(t *testing.T) {
	config := &Config{
		Installed: InstalledItems{
			Languages: []string{},
			Skills:    []string{},
		},
	}

	migrated := config.MigrateLanguagesToSkills()
	if migrated {
		t.Error("MigrateLanguagesToSkills() should return false with no languages")
	}
}

func TestLoadConfigFrom_AltFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create alternative config file
	configContent := `version: "2.0.0"
registry: "https://example.com"
`
	err := os.WriteFile(filepath.Join(tmpDir, ".aicof.yaml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	config, err := LoadConfigFrom(tmpDir)
	if err != nil {
		t.Fatalf("LoadConfigFrom() error = %v", err)
	}
	if config.Version != "2.0.0" {
		t.Errorf("config.Version = %q, want %q", config.Version, "2.0.0")
	}
}
