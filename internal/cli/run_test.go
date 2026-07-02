package cli

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestRunCommandRejectsMissingConfig(t *testing.T) {
	cmd, out := newTestCommand(t, "run", "--config", "missing.yml")

	err := cmd.Execute()
	if err == nil {
		t.Fatal("run returned nil error for missing config")
	}
	errText := err.Error()
	if !strings.Contains(errText, "load config") {
		t.Fatalf("error = %q, want load config", errText)
	}
	if strings.Contains(out.String(), "Usage:") {
		t.Fatalf("runtime error should not print usage:\n%s", out.String())
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

	cmd, out := newTestCommand(t, "run", "--config", configPath, "--check")

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

	cmd, _ := newTestCommand(t, "run", "--config", configPath, "--check")

	err := cmd.Execute()
	if err == nil {
		t.Fatal("run --check strict returned nil error")
	}
	errText := err.Error()
	if !strings.Contains(errText, "strict mode rejected unsafe config fields") {
		t.Fatalf("error = %q", errText)
	}
	if strings.Contains(errText, "sk_live_secret_should_not_print") {
		t.Fatalf("error leaked secret value: %q", errText)
	}
}

func TestRunCheckHostOverrideRecordsPublicBindWarning(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "safe.yml")
	if err := os.WriteFile(configPath, []byte(`version: "0.1"
server:
  host: 127.0.0.1
  port: 43101
mode: ai-safe
adapters:
  stripe:
    enabled: true
    base_path: /stripe
`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cmd, out := newTestCommand(t, "run", "--config", configPath, "--host", "0.0.0.0", "--check")

	if err := cmd.Execute(); err != nil {
		t.Fatalf("run --check --host: %v", err)
	}
	got := out.String()
	for _, want := range []string{"MOCKPORT SECURITY WARNING", "server.host", "public_bind", "Config check passed"} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q:\n%s", want, got)
		}
	}
}

func TestRunCheckStrictRejectsUnsafeHostOverride(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "safe.yml")
	if err := os.WriteFile(configPath, []byte(`version: "0.1"
server:
  host: 127.0.0.1
  port: 43101
mode: strict
adapters:
  stripe:
    enabled: true
    base_path: /stripe
`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cmd, _ := newTestCommand(t, "run", "--config", configPath, "--host", "0.0.0.0", "--check")

	err := cmd.Execute()
	if err == nil {
		t.Fatal("run --check strict --host returned nil error")
	}
	errText := err.Error()
	if !strings.Contains(errText, "server.host") {
		t.Fatalf("error = %q, want server.host", errText)
	}
}

// S1: --host オーバーライド経路で Validate が2回呼ばれても警告が重複しないことを固定する。
func TestRunCheckHostOverrideDoesNotDuplicateWarnings(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "with-scenarios.yml")
	if err := os.WriteFile(configPath, []byte(`version: "0.1"
server:
  host: 127.0.0.1
  port: 43101
mode: ai-safe
adapters:
  stripe:
    enabled: true
    base_path: /stripe
scenarios:
  payment_success:
    adapter: stripe
`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cmd, out := newTestCommand(t, "run", "--config", configPath, "--host", "0.0.0.0", "--check")

	if err := cmd.Execute(); err != nil {
		t.Fatalf("run --check --host: %v", err)
	}
	got := out.String()
	if n := strings.Count(got, "unsupported_config"); n != 1 {
		t.Fatalf("expected unsupported_config warning exactly once, got %d:\n%s", n, got)
	}
	if n := strings.Count(got, "public_bind"); n != 1 {
		t.Fatalf("expected public_bind warning exactly once, got %d:\n%s", n, got)
	}
}

// S2: safety warnings は stderr に出力され、--check サマリは stdout に出ることを固定する。
func TestRunCheckEmitsSafetyWarningsToStderr(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "with-scenarios.yml")
	if err := os.WriteFile(configPath, []byte(`version: "0.1"
server:
  host: 127.0.0.1
  port: 43101
mode: ai-safe
adapters:
  stripe:
    enabled: true
    base_path: /stripe
scenarios:
  payment_success:
    adapter: stripe
`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cmd := NewRootCommand()
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"run", "--config", configPath, "--check"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("run --check: %v", err)
	}

	errText := stderr.String()
	outText := stdout.String()

	for _, want := range []string{"MOCKPORT SECURITY WARNING", "scenarios", "unsupported_config"} {
		if !strings.Contains(errText, want) {
			t.Fatalf("stderr missing %q:\nstderr=%s", want, errText)
		}
		if strings.Contains(outText, want) {
			t.Fatalf("stdout should not contain safety warning %q:\nstdout=%s", want, outText)
		}
	}
	if !strings.Contains(outText, "Config check passed") {
		t.Fatalf("stdout missing Config check passed:\nstdout=%s", outText)
	}
	if strings.Contains(errText, "Config check passed") {
		t.Fatalf("Config check passed should not go to stderr:\nstderr=%s", errText)
	}
}

func TestServeHTTPShutsDownWhenContextIsCanceled(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- serveHTTP(ctx, listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
	}()

	resp, err := http.Get("http://" + listener.Addr().String())
	if err != nil {
		t.Fatalf("get server: %v", err)
	}
	_ = resp.Body.Close()
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("serveHTTP returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("serveHTTP did not shut down")
	}
}

func TestNewHTTPServerUsesBoundedTimeouts(t *testing.T) {
	server := newHTTPServer(http.NotFoundHandler())

	if server.ReadHeaderTimeout != serverReadHeaderTimeout {
		t.Fatalf("ReadHeaderTimeout = %s, want %s", server.ReadHeaderTimeout, serverReadHeaderTimeout)
	}
	if server.ReadTimeout != serverReadTimeout {
		t.Fatalf("ReadTimeout = %s, want %s", server.ReadTimeout, serverReadTimeout)
	}
	if server.IdleTimeout != serverIdleTimeout {
		t.Fatalf("IdleTimeout = %s, want %s", server.IdleTimeout, serverIdleTimeout)
	}
	if server.MaxHeaderBytes != serverMaxHeaderBytes {
		t.Fatalf("MaxHeaderBytes = %d, want %d", server.MaxHeaderBytes, serverMaxHeaderBytes)
	}
}

func TestFormatListenErrorIncludesAddress(t *testing.T) {
	err := formatListenError("127.0.0.1:43101", fmt.Errorf("listen: %w", syscall.EADDRINUSE))
	errText := err.Error()
	if !strings.Contains(errText, "127.0.0.1:43101") || !strings.Contains(errText, "address already in use") {
		t.Fatalf("error = %q", errText)
	}
}
