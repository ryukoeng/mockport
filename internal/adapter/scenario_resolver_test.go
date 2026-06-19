package adapter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

func TestScenarioResolver(t *testing.T) {
	knownScenarios := []adapter.Scenario{
		{Name: "payment_success", Supported: true},
		{Name: "payment_failed", Supported: true},
		{Name: "auth_error", Supported: true},
	}
	meta := adapter.Metadata{
		Name:            "stripe",
		Maturity:        adapter.MaturityWorkflowCompatible,
		ProviderVersion: "2025-10-29.clover",
		Scenarios:       knownScenarios,
	}

	tests := []struct {
		name           string
		configScenario string
		defaultName    string
		headerValue    string
		want           string
		wantErr        bool
	}{
		{
			name:           "ヘッダなし・config値あり→config値を返す",
			configScenario: "payment_failed",
			defaultName:    "payment_success",
			headerValue:    "",
			want:           "payment_failed",
		},
		{
			name:           "ヘッダなし・config空→default値を返す",
			configScenario: "",
			defaultName:    "payment_success",
			headerValue:    "",
			want:           "payment_success",
		},
		{
			name:           "ヘッダあり既知→ヘッダ値を返す",
			configScenario: "payment_success",
			defaultName:    "payment_success",
			headerValue:    "payment_failed",
			want:           "payment_failed",
		},
		{
			name:           "ヘッダあり未知→ErrUnknownScenario",
			configScenario: "payment_success",
			defaultName:    "payment_success",
			headerValue:    "nonexistent_scenario",
			wantErr:        true,
		},
		{
			name:           "ヘッダ前後空白トリム→正常解決",
			configScenario: "payment_success",
			defaultName:    "payment_success",
			headerValue:    "  payment_failed  ",
			want:           "payment_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := adapter.Config{Scenario: tt.configScenario}
			resolver := adapter.NewScenarioResolver(cfg, tt.defaultName, meta)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.headerValue != "" {
				req.Header.Set(adapter.ScenarioHeader, tt.headerValue)
			}

			got, err := resolver.Resolve(req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("want error, got nil (scenario=%q)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
