package cli

import (
	"bytes"
	"context"
	"reflect"
	"testing"
)

func TestUpCommandRunsDockerCompose(t *testing.T) {
	var gotName string
	var gotArgs []string
	oldRunner := runCommand
	t.Cleanup(func() { runCommand = oldRunner })
	runCommand = func(ctx context.Context, name string, args ...string) error {
		gotName = name
		gotArgs = append([]string(nil), args...)
		return nil
	}

	cmd := NewRootCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"up"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute up: %v", err)
	}

	if gotName != "docker" {
		t.Fatalf("command name = %q, want docker", gotName)
	}
	wantArgs := []string{"compose", "-f", "docker-compose.mockport.yml", "up"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("args = %#v, want %#v", gotArgs, wantArgs)
	}
}
