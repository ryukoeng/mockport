package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadValidConfigAppliesDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mockport.yml")
	err := os.WriteFile(path, []byte("version: \"0.1\"\nserver:\n  port: 43101\nmode: ai-safe\nadapters:\n  stripe:\n    enabled: true\n    base_path: /stripe\n    scenario: payment_success\n    fake_secret: mockport_stripe_secret\n"), 0o644)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Fatalf("host = %q, want default 0.0.0.0", cfg.Server.Host)
	}
	if cfg.Server.Port != 43101 {
		t.Fatalf("port = %d, want 43101", cfg.Server.Port)
	}
}

func TestLoadRejectsInvalidPort(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mockport.yml")
	err := os.WriteFile(path, []byte("version: \"0.1\"\nserver:\n  port: 0\nmode: ai-safe\n"), 0o644)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}

	if _, err := LoadFile(path); err == nil {
		t.Fatal("LoadFile returned nil error for invalid port")
	}
}

func TestAISafeRecordsRealLookingSecretWarning(t *testing.T) {
	cfg := Config{
		Server: ServerConfig{Port: 43101},
		Mode:   "ai-safe",
		Adapters: map[string]AdapterConfig{
			"stripe": {Enabled: true, BasePath: "/stripe", FakeSecret: "sk_live_123"},
		},
	}

	if err := Validate(&cfg); err != nil {
		t.Fatalf("validate ai-safe config: %v", err)
	}
	if len(cfg.SafetyWarnings) != 1 {
		t.Fatalf("warnings = %d, want 1", len(cfg.SafetyWarnings))
	}
}

func TestStrictRejectsRealLookingSecret(t *testing.T) {
	cfg := Config{
		Server: ServerConfig{Port: 43101},
		Mode:   "strict",
		Adapters: map[string]AdapterConfig{
			"stripe": {Enabled: true, BasePath: "/stripe", FakeSecret: "sk_live_123"},
		},
	}

	if err := Validate(&cfg); err == nil {
		t.Fatal("strict config with real-looking secret returned nil error")
	}
}

func TestStrictRejectsExternalServiceURL(t *testing.T) {
	cfg := Config{
		Server: ServerConfig{Port: 43101},
		Mode:   "strict",
		Adapters: map[string]AdapterConfig{
			"stripe": {Enabled: true, BasePath: "/stripe", APIURL: "https://api.stripe.com"},
		},
	}

	if err := Validate(&cfg); err == nil {
		t.Fatal("strict config with external service URL returned nil error")
	}
}
