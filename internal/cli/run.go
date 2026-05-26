package cli

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/albert-einshutoin/mockport/adapters/githuboauth"
	"github.com/albert-einshutoin/mockport/adapters/openai"
	"github.com/albert-einshutoin/mockport/adapters/slack"
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
			reg, err := defaultRegistry()
			if err != nil {
				return err
			}
			handler, err := server.NewConfiguredHandler(cfg, reg, report.NewRecorder())
			if err != nil {
				return err
			}
			addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				return formatListenError(addr, err)
			}
			ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()
			return serveHTTP(ctx, listener, handler)
		},
	}
	cmd.Flags().StringVar(&configPath, "config", "mockport.yml", "Path to mockport.yml")
	cmd.Flags().BoolVar(&check, "check", false, "Validate config and print safety warnings without starting the server")
	return cmd
}

func defaultRegistry() (*adapter.Registry, error) {
	reg := adapter.NewRegistry()
	for _, adapterImpl := range []adapter.Adapter{stripe.New(), openai.New(), githuboauth.New(), slack.New()} {
		if err := reg.Register(adapterImpl); err != nil {
			return nil, err
		}
	}
	return reg, nil
}

func serveHTTP(ctx context.Context, listener net.Listener, handler http.Handler) error {
	server := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	errc := make(chan error, 1)
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errc <- err
			return
		}
		errc <- nil
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}
		return <-errc
	}
}

func formatListenError(addr string, err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "address already in use") {
		return fmt.Errorf("listen on %s: address already in use; choose another port or stop the existing process: %w", addr, err)
	}
	return fmt.Errorf("listen on %s: %w", addr, err)
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
