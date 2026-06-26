package cli

import (
	"strings"
	"testing"
)

func TestHelpServiceShowsAdapterImplementationAndSpec(t *testing.T) {
	cmd, out := newTestCommand(t, "help", "stripe")

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute service help: %v", err)
	}

	got := out.String()
	for _, want := range []string{
		"Mockport service: stripe",
		"base_path: /stripe",
		"default_scenario: payment_success",
		"maturity: workflow-compatible",
		"capabilities: checkout_sessions",
		"stateful_resources: checkout_session",
		"STRIPE_API_URL=http://localhost:43101/stripe",
		"payment_failed (supported)",
		"POST /stripe/v1/checkout/sessions",
		"mockport add stripe",
		"docs/adapters/stripe.md",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("service help missing %q:\n%s", want, got)
		}
	}
}

func TestHelpServiceSupportsEveryBuiltInService(t *testing.T) {
	for _, service := range supportedServiceNames() {
		t.Run(service, func(t *testing.T) {
			cmd, out := newTestCommand(t, "help", service)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute service help: %v", err)
			}
			if got := out.String(); !strings.Contains(got, "Mockport service: "+service) {
				t.Fatalf("service help missing title for %s:\n%s", service, got)
			}
		})
	}
}

func TestHelpCommandStillShowsCommandHelp(t *testing.T) {
	cmd, out := newTestCommand(t, "help", "add")

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute command help: %v", err)
	}

	got := out.String()
	for _, want := range []string{
		"Add adapter config to mockport.yml",
		"add [adapter...]",
		"--config",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("command help missing %q:\n%s", want, got)
		}
	}
}

func TestHelpServiceRejectsUnsupportedService(t *testing.T) {
	cmd, _ := newTestCommand(t, "help", "unknown")

	err := cmd.Execute()
	if err == nil {
		t.Fatal("help returned nil error for unsupported service")
	}
	errText := err.Error()
	if !strings.Contains(errText, `unsupported service "unknown"`) {
		t.Fatalf("error = %q, want unsupported service", errText)
	}
	if !strings.Contains(errText, "github-oauth, line, openai, slack, stripe, zoho-oauth") {
		t.Fatalf("error missing supported services: %q", errText)
	}
}
