package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func isolateConfigEnvironment(t *testing.T) {
	t.Helper()

	keys := configEnvironmentKeys()
	values := make(map[string]string, len(keys))
	exists := make(map[string]bool, len(keys))
	for _, key := range keys {
		values[key], exists[key] = os.LookupEnv(key)
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("unset %s: %v", key, err)
		}
	}
	t.Cleanup(func() {
		for _, key := range keys {
			if exists[key] {
				if err := os.Setenv(key, values[key]); err != nil {
					t.Errorf("restore %s: %v", key, err)
				}
			} else {
				if err := os.Unsetenv(key); err != nil {
					t.Errorf("restore %s: %v", key, err)
				}
			}
		}
	})
}

func TestLoadConfigDefaultsPublicBaseURLToServicePort(t *testing.T) {
	isolateConfigEnvironment(t)

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
	isolateConfigEnvironment(t)

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

func TestLoadConfigExportsAllValuesToEnvironment(t *testing.T) {
	isolateConfigEnvironment(t)

	configDir := t.TempDir()
	configPath := filepath.Join(configDir, ".env")
	const contents = "SERVICE_NAME=url-shortener\nSERVICE_PORT=8123\nPUBLIC_BASE_URL=https://short.example\nMONGODB_CONNECTION_STRING=mongodb://localhost\nMONGODB_USERNAME=mongo-user\nMONGODB_PASSWORD=mongo-pass\nREDIS_HOST=localhost:6379\nREDIS_USERNAME=redis-user\nREDIS_PASSWORD=redis-pass\nREDIS_DB=3\nREDIS_TLS=true\n"
	if err := os.WriteFile(configPath, []byte(contents), 0600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	if _, err := LoadConfig(configDir); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	want := Config{
		ServiceName:   "url-shortener",
		ServerPort:    "8123",
		PublicBaseURL: "https://short.example",
		MongoURI:      "mongodb://localhost",
		MongoUsername: "mongo-user",
		MongoPassword: "mongo-pass",
		RedisAddr:     "localhost:6379",
		RedisUsername: "redis-user",
		RedisPassword: "redis-pass",
		RedisDB:       3,
		RedisTLS:      true,
	}
	wantValue := reflect.ValueOf(want)
	wantType := wantValue.Type()
	for i := 0; i < wantType.NumField(); i++ {
		key := wantType.Field(i).Tag.Get("mapstructure")
		if key == "" || key == "-" {
			continue
		}
		expected := fmt.Sprint(wantValue.Field(i).Interface())
		if got := os.Getenv(key); got != expected {
			t.Errorf("%s = %q, want %q", key, got, expected)
		}
	}
}
