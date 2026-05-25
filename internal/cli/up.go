package cli

import (
	"context"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var runCommand = func(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func newUpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Run generated Docker Compose",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(cmd.Context(), "docker", "compose", "-f", "docker-compose.mockport.yml", "up")
		},
	}
}
