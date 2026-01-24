package config

import (
	"testing"
)

func TestConfigStruct(t *testing.T) {
	cfg := &Config{
		APIURL:       "http://localhost:8080",
		Project:      "test-project",
		Organization: "test-org",
		Token:        "test-token",
	}

	if cfg.APIURL != "http://localhost:8080" {
		t.Errorf("Expected APIURL to be http://localhost:8080, got %s", cfg.APIURL)
	}

	if cfg.Project != "test-project" {
		t.Errorf("Expected Project to be test-project, got %s", cfg.Project)
	}

	if cfg.Organization != "test-org" {
		t.Errorf("Expected Organization to be test-org, got %s", cfg.Organization)
	}

	if cfg.Token != "test-token" {
		t.Errorf("Expected Token to be test-token, got %s", cfg.Token)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := &Config{}

	if cfg.APIURL != "" {
		t.Errorf("Expected empty APIURL for default config, got %s", cfg.APIURL)
	}

	if cfg.Project != "" {
		t.Errorf("Expected empty Project for default config, got %s", cfg.Project)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				APIURL:       "http://localhost:8080",
				Project:      "test",
				Organization: "org",
				Token:        "token",
			},
			wantErr: false,
		},
		{
			name: "missing api url",
			cfg: &Config{
				Project:      "test",
				Organization: "org",
				Token:        "token",
			},
			wantErr: true,
		},
		{
			name: "missing project",
			cfg: &Config{
				APIURL:       "http://localhost:8080",
				Organization: "org",
				Token:        "token",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Test that LoadConfig doesn't panic and returns a config
	cfg, err := LoadConfig()
	if err != nil {
		// It's OK if config file doesn't exist
		t.Logf("LoadConfig returned error (expected if no config file): %v", err)
	}
	
	if cfg != nil {
		// If we got a config, it should have the default API URL
		if cfg.APIURL == "" {
			t.Errorf("Expected default APIURL, got empty string")
		}
	}
}

func TestGetConfig(t *testing.T) {
	// Reset global config
	cfg = nil
	
	// First call should load config
	c1, err := GetConfig()
	if err != nil && c1 == nil {
		t.Logf("GetConfig returned error (expected if no config file): %v", err)
		return
	}
	
	// Second call should return cached config
	c2, err := GetConfig()
	if err != nil {
		t.Errorf("GetConfig should not error on second call: %v", err)
	}
	
	if c1 != c2 {
		t.Errorf("GetConfig should return cached config")
	}
}
