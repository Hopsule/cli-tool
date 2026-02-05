package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ProjectConfig represents the .hopsule file in a project directory
type ProjectConfig struct {
	Version int            `yaml:"version"`
	Project ProjectInfo    `yaml:"project"`
}

// ProjectInfo contains project identification
type ProjectInfo struct {
	ID           string           `yaml:"id"`
	Slug         string           `yaml:"slug"`
	Name         string           `yaml:"name,omitempty"`
	Organization OrganizationInfo `yaml:"organization"`
}

// OrganizationInfo contains organization identification
type OrganizationInfo struct {
	ID   string `yaml:"id"`
	Slug string `yaml:"slug"`
	Name string `yaml:"name,omitempty"`
}

const (
	HopsuleFileName    = ".hopsule"
	HopsuleFileVersion = 1
)

// LoadProjectConfig loads the .hopsule file from the current directory
// or searches parent directories up to the filesystem root
func LoadProjectConfig() (*ProjectConfig, string, error) {
	// Start from current directory
	dir, err := os.Getwd()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get current directory: %w", err)
	}

	return LoadProjectConfigFrom(dir)
}

// LoadProjectConfigFrom loads the .hopsule file starting from a specific directory
func LoadProjectConfigFrom(startDir string) (*ProjectConfig, string, error) {
	dir := startDir

	for {
		configPath := filepath.Join(dir, HopsuleFileName)
		
		if _, err := os.Stat(configPath); err == nil {
			// Found the file
			data, err := os.ReadFile(configPath)
			if err != nil {
				return nil, "", fmt.Errorf("failed to read %s: %w", configPath, err)
			}

			var cfg ProjectConfig
			if err := yaml.Unmarshal(data, &cfg); err != nil {
				return nil, "", fmt.Errorf("failed to parse %s: %w", configPath, err)
			}

			return &cfg, configPath, nil
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}

	return nil, "", fmt.Errorf("no %s file found in current directory or any parent", HopsuleFileName)
}

// SaveProjectConfig saves the .hopsule file to the specified directory
func SaveProjectConfig(dir string, cfg *ProjectConfig) error {
	if cfg.Version == 0 {
		cfg.Version = HopsuleFileVersion
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	configPath := filepath.Join(dir, HopsuleFileName)
	
	// Add header comment
	header := "# Hopsule Project Configuration\n# This file connects your local project to Hopsule.\n# Do not edit manually unless you know what you're doing.\n\n"
	finalData := append([]byte(header), data...)

	if err := os.WriteFile(configPath, finalData, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", configPath, err)
	}

	return nil
}

// ProjectConfigExists checks if a .hopsule file exists in the current directory
func ProjectConfigExists() bool {
	dir, err := os.Getwd()
	if err != nil {
		return false
	}

	configPath := filepath.Join(dir, HopsuleFileName)
	_, err = os.Stat(configPath)
	return err == nil
}

// GetProjectConfigPath returns the path where .hopsule would be created
func GetProjectConfigPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	return filepath.Join(dir, HopsuleFileName), nil
}
