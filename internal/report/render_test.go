package report

import (
	"strings"
	"testing"
)

func TestRenderTextIncludesTrustFields(t *testing.T) {
	text := RenderText(Snapshot{
		Mode:                 "ai-safe",
		Safety:               SafetySummary{Mode: "ai-safe", Safe: false, RealLookingSecrets: 1, ExternalURLs: 1},
		Adapters:             []AdapterStatus{{Name: "stripe", BasePath: "/stripe", Enabled: true, Maturity: "partial"}},
		Requests:             []Request{{ID: 1, Method: "POST", Path: "/stripe/v1/not-supported", Status: 404, Reason: "unsupported_endpoint"}},
		ScenarioCoverage:     []ScenarioCoverage{{Adapter: "stripe", Scenarios: []ScenarioSupport{{Name: "payment_success", Supported: true}}}},
		BehaviorMatrix:       []BehaviorMatrixEntry{{Adapter: "stripe", Method: "POST", Path: "/stripe/v1/checkout/sessions", Maturity: "partial"}},
		UnsupportedEndpoints: []UnsupportedEndpoint{{Method: "POST", Path: "/stripe/v1/not-supported", Status: 404, Reason: "unsupported_endpoint"}},
		SafetyWarnings:       []SafetyWarning{{Field: "stripe.fake_secret", Category: "real_looking_secret", Message: "real-looking secret detected"}},
	})

	for _, want := range []string{
		"Safety: safe=false real-looking-secrets=1 external-urls=1",
		"stripe enabled at /stripe maturity=partial",
		"reason=unsupported_endpoint",
		"Scenario coverage:",
		"Behavior matrix:",
		"Unsupported endpoints: 1",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("rendered text missing %q:\n%s", want, text)
		}
	}
}
