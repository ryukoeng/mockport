package cli

import (
	"fmt"
	"net/http"

	"github.com/albert-einshutoin/mockport/adapters/stripe"
	"github.com/albert-einshutoin/mockport/internal/adapter"
	"github.com/albert-einshutoin/mockport/internal/config"
	"github.com/albert-einshutoin/mockport/internal/report"
	"github.com/albert-einshutoin/mockport/internal/server"
	"github.com/spf13/cobra"
)

func newRunCommand() *cobra.Command {
	var configPath string
	var check bool
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run Mockport server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadFile(configPath)
			if err != nil {
				return err
			}
			printSafetyWarnings(cmd, cfg)
			if check {
				fmt.Fprintln(cmd.OutOrStdout(), "Config check passed")
				return nil
			}
			reg := adapter.NewRegistry()
			reg.Register(stripe.New())
			handler, err := server.NewConfiguredHandler(cfg, reg, report.NewRecorder())
			if err != nil {
				return err
			}
			addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
			return http.ListenAndServe(addr, handler)
		},
	}
	cmd.Flags().StringVar(&configPath, "config", "mockport.yml", "Path to mockport.yml")
	cmd.Flags().BoolVar(&check, "check", false, "Validate config and print safety warnings without starting the server")
	return cmd
}

func printSafetyWarnings(cmd *cobra.Command, cfg config.Config) {
	if len(cfg.SafetyWarnings) == 0 {
		return
	}
	out := cmd.OutOrStdout()
	fmt.Fprintln(out, "[MOCKPORT SECURITY WARNING]")
	for _, warning := range cfg.SafetyWarnings {
		fmt.Fprintf(out, "- %s: %s (%s)\n", warning.Field, warning.Message, warning.Category)
	}
}
