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

func TestReportEndpointReturnsRequestsAndSafety(t *testing.T) {
	cfg := config.Config{
		Mode:   "ai-safe",
		Server: config.ServerConfig{Host: "0.0.0.0", Port: 43101},
		Adapters: map[string]config.AdapterConfig{
			"stripe": {Enabled: true, BasePath: "/stripe", Scenario: "payment_success", FakeSecret: "sk_live_123"},
		},
	}
	if err := config.Validate(&cfg); err != nil {
		t.Fatalf("validate config: %v", err)
	}
	reg := adapter.NewRegistry()
	reg.Register(stripe.New())
	recorder := report.NewRecorder()
	handler, err := NewConfiguredHandler(cfg, reg, recorder)
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/stripe/v1/checkout/sessions", nil))
	reportRec := httptest.NewRecorder()
	handler.ServeHTTP(reportRec, httptest.NewRequest(http.MethodGet, "/_mockport/report", nil))

	if reportRec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", reportRec.Code, http.StatusOK)
	}
	var snapshot report.Snapshot
	if err := json.Unmarshal(reportRec.Body.Bytes(), &snapshot); err != nil {
		t.Fatalf("decode report: %v", err)
	}
	if snapshot.Mode != "ai-safe" {
		t.Fatalf("mode = %q", snapshot.Mode)
	}
	if len(snapshot.Adapters) != 1 || snapshot.Adapters[0].Name != "stripe" {
		t.Fatalf("adapters = %#v", snapshot.Adapters)
	}
	if len(snapshot.Requests) != 1 || snapshot.Requests[0].Path != "/stripe/v1/checkout/sessions" {
		t.Fatalf("requests = %#v", snapshot.Requests)
	}
	if len(snapshot.SafetyWarnings) != 1 {
		t.Fatalf("safety warnings = %#v", snapshot.SafetyWarnings)
	}
}
