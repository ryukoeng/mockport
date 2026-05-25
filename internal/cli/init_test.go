package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitGeneratesStripeFiles(t *testing.T) {
	dir := t.TempDir()
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	}()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir temp: %v", err)
	}

	cmd := NewRootCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"init", "--adapter", "stripe"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute init: %v", err)
	}

	for _, name := range []string{"mockport.yml", ".env.mockport", "docker-compose.mockport.yml"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("expected generated file %s: %v", name, err)
		}
	}
	envData, err := os.ReadFile(filepath.Join(dir, ".env.mockport"))
	if err != nil {
		t.Fatalf("read env: %v", err)
	}
	env := string(envData)
	if !strings.Contains(env, "STRIPE_API_URL=http://localhost:43101/stripe") {
		t.Fatalf("env missing Stripe URL: %s", env)
	}
	if !strings.Contains(env, "STRIPE_SECRET_KEY=mockport_stripe_secret") {
		t.Fatalf("env missing fake Stripe key: %s", env)
	}
}
