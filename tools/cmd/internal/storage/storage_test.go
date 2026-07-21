package storage

import (
	"os"
	"path/filepath"
	"testing"

	"news4coder/internal/subscription"
)

func TestStorageLoadAndSave(t *testing.T) {
	// Override home dir
	originalHome := os.Getenv("HOME")
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	s, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Load non-existent config should return empty config
	cfg, err := s.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load() returned nil")
	}

	// Save config with subscriptions
	cfg.Subscriptions = append(cfg.Subscriptions, subscription.Subscription{
		Name: "Test",
		URL:  "https://example.com",
	})
	if err := s.Save(cfg); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Verify file exists
	configPath := filepath.Join(tmpHome, ".news4coder", "subscriptions.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("subscriptions.json was not created")
	}

	// Reload
	cfg2, err := s.Load()
	if err != nil {
		t.Fatalf("Load() second time failed: %v", err)
	}
	if len(cfg2.Subscriptions) != 1 {
		t.Errorf("expected 1 subscription, got %d", len(cfg2.Subscriptions))
	}
	if cfg2.Subscriptions[0].Name != "Test" {
		t.Errorf("expected subscription name Test, got %s", cfg2.Subscriptions[0].Name)
	}
}
