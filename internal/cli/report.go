package cli

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/albert-einshutoin/mockport/internal/report"
	"github.com/spf13/cobra"
)

func newReportCommand() *cobra.Command {
	var reportURL string
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Print Mockport request and safety report",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := http.Get(reportURL)
			if err != nil {
				return fmt.Errorf("fetch report: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("fetch report: status %d", resp.StatusCode)
			}

			var snapshot report.Snapshot
			if err := json.NewDecoder(resp.Body).Decode(&snapshot); err != nil {
				return fmt.Errorf("decode report: %w", err)
			}
			printReport(cmd, snapshot)
			return nil
		},
	}
	cmd.Flags().StringVar(&reportURL, "url", "http://localhost:43101/_mockport/report", "Mockport report endpoint URL")
	return cmd
}

func printReport(cmd *cobra.Command, snapshot report.Snapshot) {
	out := cmd.OutOrStdout()
	fmt.Fprintln(out, "Mockport Report")
	fmt.Fprintln(out)
	fmt.Fprintf(out, "Mode: %s\n", snapshot.Mode)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Adapters:")
	for _, adapter := range snapshot.Adapters {
		if adapter.Enabled {
			fmt.Fprintf(out, "- %s enabled at %s\n", adapter.Name, adapter.BasePath)
		}
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Requests:")
	for _, request := range snapshot.Requests {
		fmt.Fprintf(out, "- %s %s -> %d\n", request.Method, request.Path, request.Status)
	}
	fmt.Fprintln(out)
	fmt.Fprintf(out, "Safety warnings: %d\n", len(snapshot.SafetyWarnings))
}
