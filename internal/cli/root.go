package cli

import "github.com/spf13/cobra"

func Execute() error {
	return NewRootCommand().Execute()
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mockport",
		Short: "Secret-free service emulation for local and CI integration tests",
	}
	cmd.AddCommand(newVersionCommand())
	cmd.AddCommand(newRunCommand())
	cmd.AddCommand(newInitCommand())
	return cmd
}
