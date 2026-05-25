package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCommandRejectsMissingConfig(t *testing.T) {
	cmd := NewRootCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"run", "--config", "missing.yml"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("run returned nil error for missing config")
	}
	if !strings.Contains(err.Error(), "load config") {
		t.Fatalf("error = %q, want load config", err.Error())
	}
}

func TestRunCheckPrintsAISafeWarningsWithoutSecretValues(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "unsafe.yml")
	if err := os.WriteFile(configPath, []byte(`version: "0.1"
server:
  host: 127.0.0.1
  port: 43101
mode: ai-safe
adapters:
  stripe:
    enabled: true
    base_path: /stripe
    scenario: payment_success
    fake_secret: sk_live_secret_should_not_print
    api_url: https://api.stripe.com
`), 0o644); err != nil {
		t.Fatalf("write unsafe config: %v", err)
	}

	cmd := NewRootCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"run", "--config", configPath, "--check"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("run --check: %v", err)
	}

	got := out.String()
	for _, want := range []string{
		"MOCKPORT SECURITY WARNING",
		"stripe.fake_secret",
		"real_looking_secret",
		"stripe.api_url",
		"external_url",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q:\n%s", want, got)
		}
	}
	for _, leaked := range []string{"sk_live_secret_should_not_print", "https://api.stripe.com"} {
		if strings.Contains(got, leaked) {
			t.Fatalf("output leaked unsafe value %q:\n%s", leaked, got)
		}
	}
}

func TestRunCheckStrictRejectsUnsafeConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "unsafe.yml")
	if err := os.WriteFile(configPath, []byte(`version: "0.1"
server:
  host: 127.0.0.1
  port: 43101
mode: strict
adapters:
  stripe:
    enabled: true
    base_path: /stripe
    fake_secret: sk_live_secret_should_not_print
`), 0o644); err != nil {
		t.Fatalf("write unsafe config: %v", err)
	}

	cmd := NewRootCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"run", "--config", configPath, "--check"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("run --check strict returned nil error")
	}
	if !strings.Contains(err.Error(), "strict mode rejected unsafe config fields") {
		t.Fatalf("error = %q", err.Error())
	}
	if strings.Contains(err.Error(), "sk_live_secret_should_not_print") {
		t.Fatalf("error leaked secret value: %q", err.Error())
	}
}
