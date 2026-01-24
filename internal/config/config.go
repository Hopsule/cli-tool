package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	APIURL string `mapstructure:"api_url"`
	Token  string `mapstructure:"token"`
	Project string `mapstructure:"project"`
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
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")

	// Environment variables
	viper.SetEnvPrefix("DECISION")
	viper.AutomaticEnv()
	viper.BindEnv("api_url", "DECISION_API_URL")
	viper.BindEnv("token", "DECISION_TOKEN")
	viper.BindEnv("project", "DECISION_PROJECT")

	// Defaults
	viper.SetDefault("api_url", "http://localhost:8080")

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
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
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
	viper.Set("token", c.Token)
	viper.Set("project", c.Project)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return viper.WriteConfigAs(configFile)
}
