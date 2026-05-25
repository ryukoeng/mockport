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
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run Mockport server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadFile(configPath)
			if err != nil {
				return err
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
	return cmd
}
