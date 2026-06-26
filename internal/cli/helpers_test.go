package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/spf13/cobra"
)

func newTestCommand(t *testing.T, args ...string) (*cobra.Command, *bytes.Buffer) {
	t.Helper()
	cmd := NewRootCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)
	return cmd, &out
}

func mockDocker(t *testing.T, fileExistsFn func(string) bool, runFn func(context.Context, string, ...string) error) {
	t.Helper()
	oldRun, oldExists := runCommand, fileExists
	t.Cleanup(func() {
		runCommand, fileExists = oldRun, oldExists
	})
	if fileExistsFn != nil {
		fileExists = fileExistsFn
	}
	if runFn != nil {
		runCommand = runFn
	}
}
