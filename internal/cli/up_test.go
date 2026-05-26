package cli

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestUpCommandRunsDockerCompose(t *testing.T) {
	var gotName string
	var gotArgs []string
	oldRunner := runCommand
	oldFileExists := fileExists
	t.Cleanup(func() {
		runCommand = oldRunner
		fileExists = oldFileExists
	})
	fileExists = func(path string) bool { return path == "docker-compose.mockport.yml" }
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

func TestUpCommandSupportsDetachAndBuildFlags(t *testing.T) {
	var gotArgs []string
	oldRunner := runCommand
	oldFileExists := fileExists
	t.Cleanup(func() {
		runCommand = oldRunner
		fileExists = oldFileExists
	})
	fileExists = func(path string) bool { return path == "docker-compose.mockport.yml" }
	runCommand = func(ctx context.Context, name string, args ...string) error {
		gotArgs = append([]string(nil), args...)
		return nil
	}

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"up", "--detach", "--build"})
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
	oldRunner := runCommand
	oldFileExists := fileExists
	t.Cleanup(func() {
		runCommand = oldRunner
		fileExists = oldFileExists
	})
	fileExists = func(path string) bool { return path == "docker-compose.mockport.yml" }
	runCommand = func(ctx context.Context, name string, args ...string) error {
		gotArgs = append([]string(nil), args...)
		return nil
	}

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"up", "-d"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute up: %v", err)
	}

	wantArgs := []string{"compose", "-f", "docker-compose.mockport.yml", "up", "--detach"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("args = %#v, want %#v", gotArgs, wantArgs)
	}
}

func TestUpCommandSuggestsInitWhenComposeFileMissing(t *testing.T) {
	oldRunner := runCommand
	oldFileExists := fileExists
	t.Cleanup(func() {
		runCommand = oldRunner
		fileExists = oldFileExists
	})
	fileExists = func(path string) bool { return false }
	runCommand = func(ctx context.Context, name string, args ...string) error {
		t.Fatal("runCommand should not be called when compose file is missing")
		return nil
	}

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"up"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("execute up returned nil, want error")
	}
	if !strings.Contains(err.Error(), "docker-compose.mockport.yml") || !strings.Contains(err.Error(), "mockport init") {
		t.Fatalf("error = %q, want compose file and mockport init guidance", err.Error())
	}
}

func TestUpCommandExplainsMissingDocker(t *testing.T) {
	oldRunner := runCommand
	oldFileExists := fileExists
	t.Cleanup(func() {
		runCommand = oldRunner
		fileExists = oldFileExists
	})
	fileExists = func(path string) bool { return path == "docker-compose.mockport.yml" }
	runCommand = func(ctx context.Context, name string, args ...string) error {
		return errors.New("exec: \"docker\": executable file not found in $PATH")
	}

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"up"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("execute up returned nil, want error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "docker is required") || !strings.Contains(err.Error(), "docker compose") {
		t.Fatalf("error = %q, want Docker guidance", err.Error())
	}
}
