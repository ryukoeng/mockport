package cli

import (
	"strings"
	"testing"
)

func TestRootCommandShowsHelp(t *testing.T) {
	cmd, out := newTestCommand(t, "--help")

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute help: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "Secret-free service emulation") {
		t.Fatalf("help output missing product description: %q", got)
	}
	if !strings.Contains(got, "healthcheck") {
		t.Fatalf("help output missing healthcheck command: %q", got)
	}
}

func TestVersionCommandPrintsVersion(t *testing.T) {
	cmd, out := newTestCommand(t, "version")

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute version: %v", err)
	}

	if got := strings.TrimSpace(out.String()); got != "mockport dev" {
		t.Fatalf("version output = %q, want %q", got, "mockport dev")
	}
}
