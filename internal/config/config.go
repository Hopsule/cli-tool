package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	APIURL       string `mapstructure:"api_url"`
	WebURL       string `mapstructure:"web_url"`
	Token        string `mapstructure:"token"`
	Project      string `mapstructure:"project"`
	Organization string `mapstructure:"organization"`
	User         *User  `mapstructure:"user"`
}

// User represents the authenticated user info stored in config
type User struct {
	ID        string `mapstructure:"id"`
	Email     string `mapstructure:"email"`
	Name      string `mapstructure:"name"`
	AvatarURL string `mapstructure:"avatar_url"`
}

var cfg *Config

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".decision-cli")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// Only the per-user config dir: reading ./config.yaml from the working
	// directory would let any checked-out repo override api_url and receive
	// the user's token.
	viper.AddConfigPath(configDir)
	// config.yaml holds the API token: never world/group readable.
	viper.SetConfigPermissions(0o600)

	// Environment variables
	viper.SetEnvPrefix("DECISION")
	viper.AutomaticEnv()
	viper.BindEnv("api_url", "DECISION_API_URL")
	viper.BindEnv("token", "DECISION_TOKEN")
	viper.BindEnv("project", "DECISION_PROJECT")

	// Defaults
	viper.SetDefault("api_url", "https://api.hopsule.com")
	viper.SetDefault("web_url", "https://app.hopsule.com")

	// Read config file (optional - file may not exist)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is OK - we'll use defaults/env vars
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Ensure config directory exists for future writes
	if err := ensurePrivateDir(configDir); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetConfig returns the loaded config or loads it if not yet loaded
func GetConfig() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}
	return LoadConfig()
}

// SaveConfig saves the current config to file
func SaveConfig(c *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".decision-cli")
	configFile := filepath.Join(configDir, "config.yaml")

	viper.Set("api_url", c.APIURL)
	viper.Set("web_url", c.WebURL)
	viper.Set("token", c.Token)
	viper.Set("project", c.Project)
	viper.Set("organization", c.Organization)
	
	if c.User != nil {
		viper.Set("user.id", c.User.ID)
		viper.Set("user.email", c.User.Email)
		viper.Set("user.name", c.User.Name)
		viper.Set("user.avatar_url", c.User.AvatarURL)
	} else {
		viper.Set("user", nil)
	}

	if err := ensurePrivateDir(configDir); err != nil {
		return err
	}

	viper.SetConfigPermissions(0o600)
	if err := viper.WriteConfigAs(configFile); err != nil {
		return err
	}
	// WriteConfigAs keeps the mode of an existing file (older CLI versions
	// created it 0644), so tighten it explicitly.
	return os.Chmod(configFile, 0o600)
}

// ensurePrivateDir creates dir (0700) and tightens it if it already exists
// with broader permissions from an older CLI version.
func ensurePrivateDir(dir string) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	if err := os.Chmod(dir, 0o700); err != nil {
		return fmt.Errorf("failed to secure config directory: %w", err)
	}
	return nil
}

// ClearAuth clears authentication data from config
func ClearAuth(c *Config) error {
	c.Token = ""
	c.User = nil
	return SaveConfig(c)
}

// IsAuthenticated returns true if user is logged in
func (c *Config) IsAuthenticated() bool {
	return c.Token != "" && c.User != nil
}

// GetWebURL returns the web URL with fallback to default
func (c *Config) GetWebURL() string {
	if c.WebURL != "" {
		return c.WebURL
	}
	return "https://app.hopsule.com"
}

// GetAPIURL returns the API URL with fallback to default
func (c *Config) GetAPIURL() string {
	if c.APIURL != "" {
		return c.APIURL
	}
	return "https://api.hopsule.com"
}
