package cli

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestUpCommandRunsDockerCompose(t *testing.T) {
	var gotName string
	var gotArgs []string
	mockDocker(t,
		func(path string) bool { return path == "docker-compose.mockport.yml" },
		func(ctx context.Context, name string, args ...string) error {
			gotName = name
			gotArgs = append([]string(nil), args...)
			return nil
		},
	)

	cmd, out := newTestCommand(t, "up")
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute up: %v", err)
	}
	_ = out

	if gotName != "docker" {
		t.Fatalf("command name = %q, want docker", gotName)
	}
	wantArgs := []string{"compose", "-f", "docker-compose.mockport.yml", "up"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("args = %#v, want %#v", gotArgs, wantArgs)
	}
}

func TestUpCommandSupportsDetachAndBuildFlags(t *testing.T) {
	var gotArgs []string
	mockDocker(t,
		func(path string) bool { return path == "docker-compose.mockport.yml" },
		func(ctx context.Context, name string, args ...string) error {
			gotArgs = append([]string(nil), args...)
			return nil
		},
	)

	cmd, _ := newTestCommand(t, "up", "--detach", "--build")
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute up: %v", err)
	}

	wantArgs := []string{"compose", "-f", "docker-compose.mockport.yml", "up", "--detach", "--build"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("args = %#v, want %#v", gotArgs, wantArgs)
	}
}

func TestUpCommandSupportsShortDetachFlag(t *testing.T) {
	var gotArgs []string
	mockDocker(t,
		func(path string) bool { return path == "docker-compose.mockport.yml" },
		func(ctx context.Context, name string, args ...string) error {
			gotArgs = append([]string(nil), args...)
			return nil
		},
	)

	cmd, _ := newTestCommand(t, "up", "-d")
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute up: %v", err)
	}

	wantArgs := []string{"compose", "-f", "docker-compose.mockport.yml", "up", "--detach"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("args = %#v, want %#v", gotArgs, wantArgs)
	}
}

func TestUpCommandSuggestsInitWhenComposeFileMissing(t *testing.T) {
	mockDocker(t,
		func(path string) bool { return false },
		func(ctx context.Context, name string, args ...string) error {
			t.Fatal("runCommand should not be called when compose file is missing")
			return nil
		},
	)

	cmd, _ := newTestCommand(t, "up")
	err := cmd.Execute()
	if err == nil {
		t.Fatal("execute up returned nil, want error")
	}
	if !strings.Contains(err.Error(), "docker-compose.mockport.yml") || !strings.Contains(err.Error(), "mockport init") {
		t.Fatalf("error = %q, want compose file and mockport init guidance", err.Error())
	}
}

func TestUpCommandExplainsMissingDocker(t *testing.T) {
	mockDocker(t,
		func(path string) bool { return path == "docker-compose.mockport.yml" },
		func(ctx context.Context, name string, args ...string) error {
			return errors.New("exec: \"docker\": executable file not found in $PATH")
		},
	)

	cmd, _ := newTestCommand(t, "up")
	err := cmd.Execute()
	if err == nil {
		t.Fatal("execute up returned nil, want error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "docker is required") || !strings.Contains(err.Error(), "docker compose") {
		t.Fatalf("error = %q, want Docker guidance", err.Error())
	}
}
