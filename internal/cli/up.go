package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var runCommand = func(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

var fileExists = func(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func newUpCommand() *cobra.Command {
	var detach bool
	var build bool
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Run generated Docker Compose",
		RunE: func(cmd *cobra.Command, args []string) error {
			const composeFile = "docker-compose.mockport.yml"
			if !fileExists(composeFile) {
				return fmt.Errorf("%s not found; run `mockport init` first", composeFile)
			}
			composeArgs := []string{"compose", "-f", composeFile, "up"}
			if detach {
				composeArgs = append(composeArgs, "--detach")
			}
			if build {
				composeArgs = append(composeArgs, "--build")
			}
			if err := runCommand(cmd.Context(), "docker", composeArgs...); err != nil {
				if strings.Contains(err.Error(), "executable file not found") || strings.Contains(err.Error(), "not found in $PATH") {
					return fmt.Errorf("Docker is required to run `mockport up`; install Docker and ensure `docker compose` is available: %w", err)
				}
				return err
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&detach, "detach", "d", false, "Run Docker Compose in the background")
	cmd.Flags().BoolVar(&build, "build", false, "Build images before starting")
	return cmd
}
