package security

import "testing"

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

func TestLooksLikeExternalServiceURL(t *testing.T) {
	if !LooksLikeExternalServiceURL("https://api.stripe.com/v1/checkout/sessions") {
		t.Fatal("expected Stripe live API URL to be dangerous")
	}
	if LooksLikeExternalServiceURL("http://localhost:43101/stripe") {
		t.Fatal("expected localhost URL to be allowed")
	}
}
