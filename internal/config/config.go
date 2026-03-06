package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// ConfigPath returns the path to the config file
func ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".responsewatch", "config.yaml")
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	configPath := ConfigPath()
	if configPath == "" {
		return fmt.Errorf("unable to determine home directory")
	}
	configDir := filepath.Dir(configPath)
	return os.MkdirAll(configDir, 0700)
}

// APIConfig holds API-related configuration
type APIConfig struct {
	BaseURL string `yaml:"base_url"`
	Timeout int    `yaml:"timeout"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Token        string    `yaml:"token"`
	RefreshToken string    `yaml:"refresh_token"`
	ExpiresAt    time.Time `yaml:"expires_at"`
}

// UserConfig holds user profile information
type UserConfig struct {
	Email string `yaml:"email"`
	Name  string `yaml:"name"`
}

// OutputConfig holds output formatting configuration
type OutputConfig struct {
	Format string `yaml:"format"` // table | json
	Color  bool   `yaml:"color"`
}

// Config represents the complete configuration
type Config struct {
	API    APIConfig    `yaml:"api"`
	Auth   AuthConfig   `yaml:"auth"`
	User   UserConfig   `yaml:"user"`
	Output OutputConfig `yaml:"output"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		API: APIConfig{
			BaseURL: "https://response-watch.web.app/api",
			Timeout: 30,
		},
		Auth: AuthConfig{},
		User: UserConfig{},
		Output: OutputConfig{
			Format: "table",
			Color:  true,
		},
	}
}

// Load loads the configuration from the config file
func Load() (*Config, error) {
	configPath := ConfigPath()
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// Save saves the configuration to the config file
func (c *Config) Save() error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	configPath := ConfigPath()
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// IsAuthenticated returns true if the user has a valid token
func (c *Config) IsAuthenticated() bool {
	return c.Auth.Token != "" && c.Auth.ExpiresAt.After(time.Now())
}

// ClearAuth removes authentication data
func (c *Config) ClearAuth() {
	c.Auth.Token = ""
	c.Auth.RefreshToken = ""
	c.Auth.ExpiresAt = time.Time{}
	c.User.Email = ""
	c.User.Name = ""
}
