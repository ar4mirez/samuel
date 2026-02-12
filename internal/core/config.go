package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// ConfigFileName is the default config file name
	ConfigFileName = "aicof.yaml"
	// AltConfigFileName is the hidden config file name
	AltConfigFileName = ".aicof.yaml"
)

// Config represents the project's AICoF configuration
type Config struct {
	Version   string         `yaml:"version"`
	Installed InstalledItems `yaml:"installed"`
	Registry  string         `yaml:"registry,omitempty"`
}

// InstalledItems tracks what components are installed
type InstalledItems struct {
	Languages  []string `yaml:"languages,omitempty"`
	Frameworks []string `yaml:"frameworks,omitempty"`
	Workflows  []string `yaml:"workflows,omitempty"`
	Skills     []string `yaml:"skills,omitempty"`
}

// NewConfig creates a new config with defaults
func NewConfig(version string) *Config {
	return &Config{
		Version: version,
		Installed: InstalledItems{
			Languages:  []string{},
			Frameworks: []string{},
			Workflows:  []string{"all"},
			Skills:     []string{},
		},
		Registry: DefaultRegistry,
	}
}

// LoadConfig loads config from the current directory
func LoadConfig() (*Config, error) {
	return LoadConfigFrom(".")
}

// LoadConfigFrom loads config from a specific directory
func LoadConfigFrom(dir string) (*Config, error) {
	// Try primary config file
	configPath := filepath.Join(dir, ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Try alternative config file
		configPath = filepath.Join(dir, AltConfigFileName)
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return nil, os.ErrNotExist
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Save writes the config to the specified directory
func (c *Config) Save(dir string) error {
	configPath := filepath.Join(dir, ConfigFileName)

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// ConfigExists checks if a config file exists in the directory
func ConfigExists(dir string) bool {
	configPath := filepath.Join(dir, ConfigFileName)
	if _, err := os.Stat(configPath); err == nil {
		return true
	}

	altConfigPath := filepath.Join(dir, AltConfigFileName)
	if _, err := os.Stat(altConfigPath); err == nil {
		return true
	}

	return false
}

// HasLanguage checks if a language is installed
func (c *Config) HasLanguage(name string) bool {
	for _, lang := range c.Installed.Languages {
		if lang == name {
			return true
		}
	}
	return false
}

// HasFramework checks if a framework is installed
func (c *Config) HasFramework(name string) bool {
	for _, fw := range c.Installed.Frameworks {
		if fw == name {
			return true
		}
	}
	return false
}

// HasWorkflow checks if a workflow is installed
func (c *Config) HasWorkflow(name string) bool {
	for _, wf := range c.Installed.Workflows {
		if wf == "all" || wf == name {
			return true
		}
	}
	return false
}

// HasSkill checks if a skill is installed
func (c *Config) HasSkill(name string) bool {
	for _, s := range c.Installed.Skills {
		if s == name {
			return true
		}
	}
	return false
}

// AddLanguage adds a language to the installed list.
// Also registers the corresponding language guide skill.
func (c *Config) AddLanguage(name string) {
	if !c.HasLanguage(name) {
		c.Installed.Languages = append(c.Installed.Languages, name)
	}
	// Also track as a skill
	skillName := LanguageToSkillName(name)
	c.AddSkill(skillName)
}

// AddFramework adds a framework to the installed list
func (c *Config) AddFramework(name string) {
	if !c.HasFramework(name) {
		c.Installed.Frameworks = append(c.Installed.Frameworks, name)
	}
}

// AddWorkflow adds a workflow to the installed list
func (c *Config) AddWorkflow(name string) {
	if !c.HasWorkflow(name) {
		c.Installed.Workflows = append(c.Installed.Workflows, name)
	}
}

// AddSkill adds a skill to the installed list
func (c *Config) AddSkill(name string) {
	if !c.HasSkill(name) {
		c.Installed.Skills = append(c.Installed.Skills, name)
	}
}

// RemoveLanguage removes a language from the installed list.
// Also removes the corresponding language guide skill.
func (c *Config) RemoveLanguage(name string) {
	c.Installed.Languages = removeFromSlice(c.Installed.Languages, name)
	skillName := LanguageToSkillName(name)
	c.RemoveSkill(skillName)
}

// RemoveFramework removes a framework from the installed list
func (c *Config) RemoveFramework(name string) {
	c.Installed.Frameworks = removeFromSlice(c.Installed.Frameworks, name)
}

// RemoveWorkflow removes a workflow from the installed list
func (c *Config) RemoveWorkflow(name string) {
	c.Installed.Workflows = removeFromSlice(c.Installed.Workflows, name)
}

// RemoveSkill removes a skill from the installed list
func (c *Config) RemoveSkill(name string) {
	c.Installed.Skills = removeFromSlice(c.Installed.Skills, name)
}

func removeFromSlice(slice []string, item string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}

// MigrateLanguagesToSkills ensures all installed languages have corresponding skills.
// This handles backward compatibility for configs created before language-as-skills migration.
func (c *Config) MigrateLanguagesToSkills() bool {
	migrated := false
	for _, lang := range c.Installed.Languages {
		skillName := LanguageToSkillName(lang)
		if !c.HasSkill(skillName) {
			c.Installed.Skills = append(c.Installed.Skills, skillName)
			migrated = true
		}
	}
	return migrated
}

// GlobalConfig represents global CLI settings stored in ~/.config/aicof/
type GlobalConfig struct {
	DefaultTemplate   string   `yaml:"default_template,omitempty"`
	DefaultLanguages  []string `yaml:"default_languages,omitempty"`
	DefaultFrameworks []string `yaml:"default_frameworks,omitempty"`
	CachePath         string   `yaml:"cache_path,omitempty"`
}

// GetGlobalConfigPath returns the path to the global config directory
func GetGlobalConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "aicof"), nil
}

// GetCachePath returns the path to the cache directory
func GetCachePath() (string, error) {
	globalPath, err := GetGlobalConfigPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(globalPath, "cache"), nil
}

// EnsureCacheDir creates the cache directory if it doesn't exist
func EnsureCacheDir() (string, error) {
	cachePath, err := GetCachePath()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return "", err
	}

	return cachePath, nil
}

// ValidConfigKeys lists all valid configuration keys
var ValidConfigKeys = []string{
	"version",
	"registry",
	"installed.languages",
	"installed.frameworks",
	"installed.workflows",
	"installed.skills",
}

// GetValue retrieves a configuration value by key
func (c *Config) GetValue(key string) (any, error) {
	switch key {
	case "version":
		return c.Version, nil
	case "registry":
		if c.Registry == "" {
			return DefaultRegistry, nil
		}
		return c.Registry, nil
	case "installed.languages":
		return c.Installed.Languages, nil
	case "installed.frameworks":
		return c.Installed.Frameworks, nil
	case "installed.workflows":
		return c.Installed.Workflows, nil
	case "installed.skills":
		return c.Installed.Skills, nil
	default:
		return nil, fmt.Errorf("unknown config key: %s", key)
	}
}

// SetValue sets a configuration value by key
func (c *Config) SetValue(key, value string) error {
	switch key {
	case "version":
		c.Version = value
	case "registry":
		c.Registry = value
	case "installed.languages":
		c.Installed.Languages = splitAndTrim(value)
	case "installed.frameworks":
		c.Installed.Frameworks = splitAndTrim(value)
	case "installed.workflows":
		c.Installed.Workflows = splitAndTrim(value)
	case "installed.skills":
		c.Installed.Skills = splitAndTrim(value)
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
	return nil
}

// GetAllValues returns all config values as a map
func (c *Config) GetAllValues() map[string]any {
	registry := c.Registry
	if registry == "" {
		registry = DefaultRegistry
	}
	return map[string]any{
		"version":              c.Version,
		"registry":             registry,
		"installed.languages":  c.Installed.Languages,
		"installed.frameworks": c.Installed.Frameworks,
		"installed.workflows":  c.Installed.Workflows,
		"installed.skills":     c.Installed.Skills,
	}
}

// splitAndTrim splits a comma-separated string and trims whitespace
func splitAndTrim(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
