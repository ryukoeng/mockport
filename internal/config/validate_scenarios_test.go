package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Step 1: scenarios: ブロック使用時の警告テスト

func TestValidateScenariosBlockEmitsWarning(t *testing.T) {
	cfg := Config{
		Server: ServerConfig{Port: 43101},
		Mode:   "ai-safe",
		Scenarios: map[string]Scenario{
			"payment_success": {Adapter: "stripe"},
		},
	}
	if err := Validate(&cfg); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	found := false
	for _, w := range cfg.SafetyWarnings {
		if w.Category == "unsupported_config" && w.Field == "scenarios" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected unsupported_config warning for scenarios, got: %+v", cfg.SafetyWarnings)
	}
}

func TestValidateNoScenariosBlockNoWarning(t *testing.T) {
	cfg := Config{
		Server:    ServerConfig{Port: 43101},
		Mode:      "ai-safe",
		Scenarios: nil,
	}
	if err := Validate(&cfg); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	for _, w := range cfg.SafetyWarnings {
		if w.Category == "unsupported_config" {
			t.Fatalf("unexpected unsupported_config warning: %+v", w)
		}
	}
}

func TestStrictModeDoesNotRejectUnsupportedConfigWarning(t *testing.T) {
	// strict モードでも unsupported_config 警告では起動拒否しないことを確認
	cfg := Config{
		Server: ServerConfig{Port: 43101},
		Mode:   "strict",
		Scenarios: map[string]Scenario{
			"payment_success": {Adapter: "stripe"},
		},
	}
	if err := Validate(&cfg); err != nil {
		t.Fatalf("strict mode must not reject unsupported_config warning, got: %v", err)
	}
	// 警告は記録される
	found := false
	for _, w := range cfg.SafetyWarnings {
		if w.Category == "unsupported_config" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected unsupported_config warning in strict mode, got: %+v", cfg.SafetyWarnings)
	}
}

func TestStrictModeStillRejectsRealSecretEvenWithScenarios(t *testing.T) {
	// strict モードが既存のシークレット警告で正しく拒否することを確認
	cfg := Config{
		Server: ServerConfig{Port: 43101},
		Mode:   "strict",
		Adapters: map[string]AdapterConfig{
			"stripe": {Enabled: true, BasePath: "/stripe", FakeSecret: "sk_live_123"},
		},
		Scenarios: map[string]Scenario{
			"payment_success": {Adapter: "stripe"},
		},
	}
	err := Validate(&cfg)
	if err == nil {
		t.Fatal("strict mode must reject real-looking secret even when scenarios block is present")
	}
	errText := err.Error()
	if !strings.Contains(errText, "stripe.fake_secret") {
		t.Fatalf("error = %q, want stripe.fake_secret", errText)
	}
}

// Step 2: 未知キー検出テスト

func TestLoadFileUnknownKeyEmitsWarning(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mockport.yml")
	yaml := "version: \"0.1\"\nserver:\n  port: 43101\nmode: ai-safe\n# typo: 'adapteres' instead of 'adapters'\nadapteres:\n  stripe:\n    enabled: true\n"
	if err := os.WriteFile(path, []byte(yaml), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile must succeed even with unknown key, got: %v", err)
	}
	found := false
	for _, w := range cfg.SafetyWarnings {
		if w.Category == "unknown_config_key" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected unknown_config_key warning, got: %+v", cfg.SafetyWarnings)
	}
}

// S1: Validate を同じ cfg に2回呼んでも警告が増殖しないこと（冪等性）を固定する。
func TestValidateIsIdempotent(t *testing.T) {
	cfg := Config{
		Server: ServerConfig{Host: "0.0.0.0", Port: 43101},
		Mode:   "ai-safe",
		Scenarios: map[string]Scenario{
			"payment_success": {Adapter: "stripe"},
		},
	}
	if err := Validate(&cfg); err != nil {
		t.Fatalf("first Validate returned error: %v", err)
	}
	first := append([]SafetyWarning(nil), cfg.SafetyWarnings...)

	if err := Validate(&cfg); err != nil {
		t.Fatalf("second Validate returned error: %v", err)
	}
	if len(cfg.SafetyWarnings) != len(first) {
		t.Fatalf("Validate is not idempotent: first=%d second=%d warnings\nfirst=%+v\nsecond=%+v",
			len(first), len(cfg.SafetyWarnings), first, cfg.SafetyWarnings)
	}
}

// S1: loader が事前に積んだ unknown_config_key 警告は Validate を呼んでも保持され、
// かつ2回呼んでも重複しないことを固定する。
func TestValidatePreservesLoaderWarningsIdempotently(t *testing.T) {
	cfg := Config{
		Server: ServerConfig{Host: "127.0.0.1", Port: 43101},
		Mode:   "ai-safe",
		SafetyWarnings: []SafetyWarning{
			{Field: "config", Category: "unknown_config_key", Message: "config contains unrecognized keys"},
		},
		Scenarios: map[string]Scenario{
			"payment_success": {Adapter: "stripe"},
		},
	}
	if err := Validate(&cfg); err != nil {
		t.Fatalf("first Validate returned error: %v", err)
	}
	if err := Validate(&cfg); err != nil {
		t.Fatalf("second Validate returned error: %v", err)
	}

	countUnknown := 0
	countScenarios := 0
	for _, w := range cfg.SafetyWarnings {
		switch w.Category {
		case "unknown_config_key":
			countUnknown++
		case "unsupported_config":
			countScenarios++
		}
	}
	if countUnknown != 1 {
		t.Fatalf("expected exactly 1 unknown_config_key warning preserved, got %d: %+v", countUnknown, cfg.SafetyWarnings)
	}
	if countScenarios != 1 {
		t.Fatalf("expected exactly 1 unsupported_config warning, got %d: %+v", countScenarios, cfg.SafetyWarnings)
	}
}

func TestLoadFileKnownKeysNoUnknownWarning(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mockport.yml")
	yaml := "version: \"0.1\"\nserver:\n  port: 43101\nmode: ai-safe\n"
	if err := os.WriteFile(path, []byte(yaml), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	for _, w := range cfg.SafetyWarnings {
		if w.Category == "unknown_config_key" {
			t.Fatalf("unexpected unknown_config_key warning: %+v", w)
		}
	}
}
