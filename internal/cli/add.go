package cli

import (
	"fmt"
	"os"

	"github.com/albert-einshutoin/mockport/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newAddCommand() *cobra.Command {
	var configPath string
	cmd := &cobra.Command{
		Use:   "add [adapter...]",
		Short: "Add adapter config to mockport.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("at least one adapter is required")
			}
			specs, err := specsFor(args)
			if err != nil {
				return err
			}
			data, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("read %s: %w", configPath, err)
			}
			var cfg config.Config
			if err := yaml.Unmarshal(data, &cfg); err != nil {
				return fmt.Errorf("parse %s: %w", configPath, err)
			}
			config.ApplyDefaults(&cfg)
			if cfg.Adapters == nil {
				cfg.Adapters = map[string]config.AdapterConfig{}
			}
			for _, spec := range specs {
				cfg.Adapters[spec.Name] = config.AdapterConfig{
					Enabled:    true,
					BasePath:   spec.BasePath,
					Scenario:   spec.Scenario,
					FakeSecret: spec.FakeSecret,
					Webhook:    spec.Webhook,
				}
			}
			if err := config.Validate(&cfg); err != nil {
				return err
			}
			out, err := yaml.Marshal(cfg)
			if err != nil {
				return fmt.Errorf("encode config: %w", err)
			}
			if err := os.WriteFile(configPath, out, 0o644); err != nil {
				return fmt.Errorf("write %s: %w", configPath, err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Updated %s\n", configPath)
			return nil
		},
	}
	cmd.Flags().StringVar(&configPath, "config", "mockport.yml", "Path to mockport.yml")
	return cmd
}
