package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigDefaultsPublicBaseURLToServicePort(t *testing.T) {
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, ".env")
	if err := os.WriteFile(configPath, []byte("SERVICE_PORT=8123\n"), 0600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	cfg, err := LoadConfig(configDir)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cfg.PublicBaseURL != "http://localhost:8123" {
		t.Errorf("PublicBaseURL = %q, want %q", cfg.PublicBaseURL, "http://localhost:8123")
	}
}

func TestLoadConfigPreservesConfiguredPublicBaseURL(t *testing.T) {
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, ".env")
	if err := os.WriteFile(configPath, []byte("SERVICE_PORT=8123\nPUBLIC_BASE_URL=https://short.example\n"), 0600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	cfg, err := LoadConfig(configDir)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cfg.PublicBaseURL != "https://short.example" {
		t.Errorf("PublicBaseURL = %q, want %q", cfg.PublicBaseURL, "https://short.example")
	}
}
