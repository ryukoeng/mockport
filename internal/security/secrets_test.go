package security

import (
	"strings"
	"testing"
)

func TestLooksLikeSecret(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"stripe live", "sk_live_123", true},
		{"stripe test", "sk_test_123", true},
		{"aws access key", "AKIAIOSFODNN7EXAMPLE", true},
		{"github token", "github_pat_abc", true},
		{"mockport fake", "mockport_stripe_secret", false},
		{"mockport webhook fake", "whsec_mockport", false},
		{"local fake", "local_openai_key", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LooksLikeSecret(tt.value)
			if got != tt.want {
				t.Fatalf("LooksLikeSecret(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestRedactSecret(t *testing.T) {
	if got := RedactSecret("mockport_stripe_secret"); got != "mockport_...cret" {
		t.Fatalf("RedactSecret fake = %q", got)
	}
	if got := RedactSecret("sk_live_123456789"); got != "[real-looking secret redacted]" {
		t.Fatalf("RedactSecret real-looking = %q", got)
	}
}

func TestRedactValue(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"short", "short", "[redacted]"},
		{"long fake", "mockport_stripe_secret", "mockport_...cret"},
		{"real secret", "sk_live_123456789", "[real-looking secret redacted]"},
		{"env assignment", "STRIPE_SECRET_KEY=sk_live_123456789", "STRIPE_SECRET_KEY=[real-looking secret redacted]"},
		{"provider url", "https://api.stripe.com/v1/checkout/sessions", "[external service URL redacted]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RedactValue(tt.value); got != tt.want {
				t.Fatalf("RedactValue(%q) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

func TestRedactMessage(t *testing.T) {
	got := RedactMessage("bad sk_live_123 and https://api.stripe.com")
	if strings.Contains(got, "sk_live_123") || strings.Contains(got, "https://api.stripe.com") {
		t.Fatalf("message leaked unsafe value: %q", got)
	}
}

func TestLooksLikeExternalServiceURL(t *testing.T) {
	if !LooksLikeExternalServiceURL("https://api.stripe.com/v1/checkout/sessions") {
		t.Fatal("expected Stripe live API URL to be dangerous")
	}
	if LooksLikeExternalServiceURL("http://localhost:43101/stripe") {
		t.Fatal("expected localhost URL to be allowed")
	}
}

func TestScanPublicEnvRejectsRealProviderSecretsAndURLs(t *testing.T) {
	env := strings.Join([]string{
		"STRIPE_SECRET_KEY=sk_live_123",
		"OPENAI_BASE_URL=https://api.openai.com/v1",
		"SLACK_BOT_TOKEN=xoxb-real-token",
		"GITHUB_OAUTH_CLIENT_SECRET=github_pat_real",
	}, "\n")

	findings := ScanPublicEnv(env)
	if len(findings) != 4 {
		t.Fatalf("finding count = %d, want 4: %#v", len(findings), findings)
	}
	for _, finding := range findings {
		if finding.Line == 0 || finding.Key == "" || finding.Reason == "" {
			t.Fatalf("incomplete finding: %#v", finding)
		}
	}
}

func TestScanPublicEnvAllowsMockportExamples(t *testing.T) {
	env := strings.Join([]string{
		"STRIPE_API_URL=http://localhost:43101/stripe",
		"STRIPE_SECRET_KEY=mockport_stripe_secret",
		"STRIPE_WEBHOOK_SECRET=whsec_mockport",
		"OPENAI_BASE_URL=http://localhost:43101/openai/v1",
		"OPENAI_API_KEY=mockport_openai_key",
		"GITHUB_OAUTH_CLIENT_SECRET=mockport_github_secret",
		"SLACK_BOT_TOKEN=mockport_slack_token",
	}, "\n")

	if findings := ScanPublicEnv(env); len(findings) != 0 {
		t.Fatalf("ScanPublicEnv() findings = %#v, want none", findings)
	}
}
