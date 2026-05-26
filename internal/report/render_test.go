package report

import (
	"strings"
	"testing"
)

func TestRenderTextIncludesTrustFields(t *testing.T) {
	text := RenderText(Snapshot{
		Mode:                 "ai-safe",
		Safety:               SafetySummary{Mode: "ai-safe", Safe: false, RealLookingSecrets: 1, ExternalURLs: 1, PublicEnvSafe: false},
		Adapters:             []AdapterStatus{{Name: "stripe", BasePath: "/stripe", Enabled: true, Maturity: "partial"}},
		Requests:             []Request{{ID: 1, Method: "POST", Path: "/stripe/v1/not-supported", Status: 404, Reason: "unsupported_endpoint"}},
		ScenarioCoverage:     []ScenarioCoverage{{Adapter: "stripe", Scenarios: []ScenarioSupport{{Name: "payment_success", Supported: true}}}},
		BehaviorMatrix:       []BehaviorMatrixEntry{{Adapter: "stripe", Method: "POST", Path: "/stripe/v1/checkout/sessions", Maturity: "partial"}},
		Compatibility:        []CompatibilityStatus{{Adapter: "stripe", Level: "wire", Score: 80, ProviderVersion: "2026-05-26", UnsupportedEndpoints: []string{"/stripe/v1/not-supported"}}},
		UnsupportedEndpoints: []UnsupportedEndpoint{{Method: "POST", Path: "/stripe/v1/not-supported", Status: 404, Reason: "unsupported_endpoint"}},
		SafetyWarnings:       []SafetyWarning{{Field: "stripe.fake_secret", Category: "real_looking_secret", Message: "real-looking secret detected"}},
	})

	for _, want := range []string{
		"Safety: safe=false real-looking-secrets=1 external-urls=1",
		"Public env safe-to-commit: false",
		"stripe enabled at /stripe maturity=partial",
		"reason=unsupported_endpoint",
		"Scenario coverage:",
		"Behavior matrix:",
		"Compatibility:",
		"stripe level=wire score=80 provider=2026-05-26 unsupported=1",
		"Unsupported endpoints: 1",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("rendered text missing %q:\n%s", want, text)
		}
	}
}
