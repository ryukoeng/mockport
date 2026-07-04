package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/albert-einshutoin/mockport/adapters/stripe"
	"github.com/albert-einshutoin/mockport/internal/adapter"
	"github.com/albert-einshutoin/mockport/internal/config"
	"github.com/albert-einshutoin/mockport/internal/report"
)

// (c) /_mockport/report 経路：scenarios: 使用時に safety_warnings に unsupported_config が含まれることを確認
func TestReportEndpointIncludesScenariosUnsupportedWarning(t *testing.T) {
	cfg := config.Config{
		Mode:   "ai-safe",
		Server: config.ServerConfig{Host: "127.0.0.1", Port: 43101},
		Adapters: map[string]config.AdapterConfig{
			"stripe": {Enabled: true, BasePath: "/stripe", Scenario: "payment_success", FakeSecret: "mockport_stripe_secret"},
		},
		Scenarios: map[string]config.Scenario{
			"payment_success": {Adapter: "stripe"},
		},
	}
	if err := config.Validate(&cfg); err != nil {
		t.Fatalf("validate config: %v", err)
	}

	reg := adapter.NewRegistry()
	if err := reg.Register(stripe.New()); err != nil {
		t.Fatalf("register stripe: %v", err)
	}
	handler, err := NewConfiguredHandler(cfg, reg, report.NewRecorder())
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	reportRec := httptest.NewRecorder()
	handler.ServeHTTP(reportRec, httptest.NewRequest(http.MethodGet, "/_mockport/report", nil))

	if reportRec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", reportRec.Code, http.StatusOK)
	}
	var snapshot report.Snapshot
	if err := json.Unmarshal(reportRec.Body.Bytes(), &snapshot); err != nil {
		t.Fatalf("decode report: %v", err)
	}

	found := false
	for _, w := range snapshot.SafetyWarnings {
		if w.Category == "unsupported_config" && w.Field == "scenarios" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected unsupported_config warning in /_mockport/report, got: %+v", snapshot.SafetyWarnings)
	}
}
