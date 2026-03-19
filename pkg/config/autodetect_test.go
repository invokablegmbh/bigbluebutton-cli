package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAutodetectFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	propsPath := filepath.Join(tmpDir, "bbb-web.properties")

	content := `# BigBlueButton properties
bigbluebutton.web.serverURL=https://bbb.example.com
securitySalt=my-super-secret-salt
other.property=value
`
	if err := os.WriteFile(propsPath, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	cfg, err := autodetectFromFile(propsPath)
	if err != nil {
		t.Fatalf("autodetectFromFile error: %v", err)
	}
	if cfg.URL != "https://bbb.example.com" {
		t.Errorf("URL = %s, want https://bbb.example.com", cfg.URL)
	}
	if cfg.Secret != "my-super-secret-salt" {
		t.Errorf("Secret = %s, want my-super-secret-salt", cfg.Secret)
	}
}

func TestAutodetectFromFile_MissingURL(t *testing.T) {
	tmpDir := t.TempDir()
	propsPath := filepath.Join(tmpDir, "bbb-web.properties")

	content := `securitySalt=my-secret`
	if err := os.WriteFile(propsPath, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	_, err := autodetectFromFile(propsPath)
	if err == nil {
		t.Fatal("expected error when URL is missing")
	}
}

func TestAutodetectFromFile_MissingSecret(t *testing.T) {
	tmpDir := t.TempDir()
	propsPath := filepath.Join(tmpDir, "bbb-web.properties")

	content := `bigbluebutton.web.serverURL=https://bbb.example.com`
	if err := os.WriteFile(propsPath, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	_, err := autodetectFromFile(propsPath)
	if err == nil {
		t.Fatal("expected error when secret is missing")
	}
}

func TestAutodetectFromFile_FileNotFound(t *testing.T) {
	_, err := autodetectFromFile("/nonexistent/path/bbb-web.properties")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestAutodetectFromFile_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	propsPath := filepath.Join(tmpDir, "bbb-web.properties")

	if err := os.WriteFile(propsPath, []byte(""), 0644); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	_, err := autodetectFromFile(propsPath)
	if err == nil {
		t.Fatal("expected error for empty file")
	}
}

func TestAutodetectFromFile_CommentsAndWhitespace(t *testing.T) {
	tmpDir := t.TempDir()
	propsPath := filepath.Join(tmpDir, "bbb-web.properties")

	content := `# Comment line

# Another comment
  bigbluebutton.web.serverURL = https://bbb.example.com
  securitySalt = secret123
`
	if err := os.WriteFile(propsPath, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	cfg, err := autodetectFromFile(propsPath)
	if err != nil {
		t.Fatalf("autodetectFromFile error: %v", err)
	}
	if cfg.URL != "https://bbb.example.com" {
		t.Errorf("URL = %s, want https://bbb.example.com", cfg.URL)
	}
	if cfg.Secret != "secret123" {
		t.Errorf("Secret = %s, want secret123", cfg.Secret)
	}
}
