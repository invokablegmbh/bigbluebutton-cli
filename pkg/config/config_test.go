package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_WithFlags(t *testing.T) {
	cfg, err := Load("https://bbb.example.com", "my-secret", "", false)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.URL != "https://bbb.example.com" {
		t.Errorf("URL = %s, want https://bbb.example.com", cfg.URL)
	}
	if cfg.Secret != "my-secret" {
		t.Errorf("Secret = %s, want my-secret", cfg.Secret)
	}
}

func TestLoad_WithEnvVars(t *testing.T) {
	os.Setenv("BBB_URL", "https://env.example.com")
	os.Setenv("BBB_SECRET", "env-secret")
	defer os.Unsetenv("BBB_URL")
	defer os.Unsetenv("BBB_SECRET")

	cfg, err := Load("", "", "", false)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.URL != "https://env.example.com" {
		t.Errorf("URL = %s, want https://env.example.com", cfg.URL)
	}
	if cfg.Secret != "env-secret" {
		t.Errorf("Secret = %s, want env-secret", cfg.Secret)
	}
}

func TestLoad_FlagsPrecedeEnv(t *testing.T) {
	os.Setenv("BBB_URL", "https://env.example.com")
	os.Setenv("BBB_SECRET", "env-secret")
	defer os.Unsetenv("BBB_URL")
	defer os.Unsetenv("BBB_SECRET")

	cfg, err := Load("https://flag.example.com", "flag-secret", "", false)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.URL != "https://flag.example.com" {
		t.Errorf("URL = %s, want https://flag.example.com", cfg.URL)
	}
	if cfg.Secret != "flag-secret" {
		t.Errorf("Secret = %s, want flag-secret", cfg.Secret)
	}
}

func TestLoad_ConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(configPath, []byte("url: https://file.example.com\nsecret: file-secret\n"), 0644)
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	cfg, err := Load("", "", configPath, false)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.URL != "https://file.example.com" {
		t.Errorf("URL = %s, want https://file.example.com", cfg.URL)
	}
	if cfg.Secret != "file-secret" {
		t.Errorf("Secret = %s, want file-secret", cfg.Secret)
	}
}

func TestLoad_MissingConfig_NonInteractive(t *testing.T) {
	_, err := Load("", "", "", false)
	if err == nil {
		t.Fatal("expected error when no config provided in non-interactive mode")
	}
}

func TestLoad_URLNormalization(t *testing.T) {
	cfg, err := Load("https://bbb.example.com/", "secret", "", false)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.URL != "https://bbb.example.com" {
		t.Errorf("URL = %s, want trailing slash removed", cfg.URL)
	}
}

func TestLoad_PartialConfig_FlagAndEnv(t *testing.T) {
	os.Setenv("BBB_SECRET", "env-secret")
	defer os.Unsetenv("BBB_SECRET")

	cfg, err := Load("https://flag.example.com", "", "", false)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.URL != "https://flag.example.com" {
		t.Errorf("URL = %s, want https://flag.example.com", cfg.URL)
	}
	if cfg.Secret != "env-secret" {
		t.Errorf("Secret = %s, want env-secret", cfg.Secret)
	}
}
