package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const scenariosYAML = `version: "0.1"
server:
  host: 127.0.0.1
  port: 43101
mode: ai-safe
adapters:
  stripe:
    enabled: true
    base_path: /stripe
    fake_secret: mockport_stripe_secret
scenarios:
  payment_success:
    adapter: stripe
    response:
      status: 200
`

// (a) mockport run 起動時の出力経路：scenarios: 使用時に警告が stdout/stderr に出ることを確認
func TestRunCheckPrintsScenariosUnsupportedWarning(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "with-scenarios.yml")
	if err := os.WriteFile(configPath, []byte(scenariosYAML), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cmd, out := newTestCommand(t, "run", "--config", configPath, "--check")

	if err := cmd.Execute(); err != nil {
		t.Fatalf("run --check: %v", err)
	}

	got := out.String()
	for _, want := range []string{
		"MOCKPORT SECURITY WARNING",
		"scenarios",
		"unsupported_config",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q:\n%s", want, got)
		}
	}
	if !strings.Contains(got, "Config check passed") {
		t.Fatalf("missing Config check passed in output:\n%s", got)
	}
}

// (b) --check 経路：strict モードでも scenarios 警告では起動拒否しないことを確認
func TestRunCheckStrictModeDoesNotRejectScenariosWarning(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "strict-scenarios.yml")
	strictYAML := `version: "0.1"
server:
  host: 127.0.0.1
  port: 43101
mode: strict
adapters:
  stripe:
    enabled: true
    base_path: /stripe
    fake_secret: mockport_stripe_secret
scenarios:
  payment_success:
    adapter: stripe
`
	if err := os.WriteFile(configPath, []byte(strictYAML), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cmd, out := newTestCommand(t, "run", "--config", configPath, "--check")

	if err := cmd.Execute(); err != nil {
		t.Fatalf("strict mode must not reject scenarios warning, got: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "unsupported_config") {
		t.Fatalf("expected unsupported_config in output:\n%s", got)
	}
	if !strings.Contains(got, "Config check passed") {
		t.Fatalf("missing Config check passed in output:\n%s", got)
	}
}
