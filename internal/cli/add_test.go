package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddAdaptersUpdatesConfig(t *testing.T) {
	dir := chdirTemp(t)
	configPath := filepath.Join(dir, "mockport.yml")
	if err := os.WriteFile(configPath, []byte(generatedConfig([]adapterSpec{mustSpec(t, "stripe")})), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cmd := NewRootCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"add", "openai", "github-oauth", "slack", "line", "--config", configPath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute add: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	configText := string(data)
	for _, want := range []string{"stripe:", "openai:", "github-oauth:", "slack:", "line:"} {
		if !strings.Contains(configText, want) {
			t.Fatalf("config missing %q:\n%s", want, configText)
		}
	}
}

func mustSpec(t *testing.T, name string) adapterSpec {
	t.Helper()
	spec, ok := adapterSpecFor(name)
	if !ok {
		t.Fatalf("missing spec for %s", name)
	}
	return spec
}
