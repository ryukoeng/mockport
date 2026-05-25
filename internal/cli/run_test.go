package cli

import (
	"bytes"
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
