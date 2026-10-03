package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigReadsEnvFileAndEnvironmentOverrides(t *testing.T) {
	configDir := t.TempDir()
	configFile := filepath.Join(configDir, ".env")
	if err := os.WriteFile(configFile, []byte("SERVICE_NAME=url-shortener\nSERVICE_PORT=8091\nMONGO_URI=mongodb://file\n"), 0600); err != nil {
		t.Fatalf("write config file: %v", err)
	}
	t.Setenv("SERVICE_PORT", "8092")
	t.Setenv("REDIS_ADDR", "redis-from-env")

	config, err := LoadConfig(configDir)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if config.ServerPort != "8092" {
		t.Errorf("ServerPort = %q, want %q", config.ServerPort, "8092")
	}
	if config.ServiceName != "url-shortener" {
		t.Errorf("ServiceName = %q, want %q", config.ServiceName, "url-shortener")
	}
	if config.MongoURI != "mongodb://file" {
		t.Errorf("MongoURI = %q, want %q", config.MongoURI, "mongodb://file")
	}
	if config.RedisAddr != "redis-from-env" {
		t.Errorf("RedisAddr = %q, want %q", config.RedisAddr, "redis-from-env")
	}
}

func TestLoadConfigDefaultsPortWithoutEnvFile(t *testing.T) {
	t.Setenv("SERVICE_PORT", "")

	config, err := LoadConfig(t.TempDir())
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if config.ServerPort != "8089" {
		t.Errorf("ServerPort = %q, want %q", config.ServerPort, "8089")
	}
}
