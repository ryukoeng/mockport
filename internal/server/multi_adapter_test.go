package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/albert-einshutoin/mockport/adapters/githuboauth"
	"github.com/albert-einshutoin/mockport/adapters/openai"
	"github.com/albert-einshutoin/mockport/adapters/slack"
	"github.com/albert-einshutoin/mockport/adapters/stripe"
	"github.com/albert-einshutoin/mockport/internal/adapter"
	"github.com/albert-einshutoin/mockport/internal/config"
	"github.com/albert-einshutoin/mockport/internal/report"
)

func TestConfiguredHandlerServesMultipleAdapters(t *testing.T) {
	cfg := config.Config{
		Mode:   "ai-safe",
		Server: config.ServerConfig{Host: "0.0.0.0", Port: 43101},
		Adapters: map[string]config.AdapterConfig{
			"stripe":       {Enabled: true, BasePath: "/stripe", Scenario: "payment_success", FakeSecret: "mockport_stripe_secret"},
			"openai":       {Enabled: true, BasePath: "/openai", Scenario: "chat_success", FakeSecret: "mockport_openai_key"},
			"github-oauth": {Enabled: true, BasePath: "/github", Scenario: "oauth_success", FakeSecret: "mockport_github_secret"},
			"slack":        {Enabled: true, BasePath: "/slack", Scenario: "message_success", FakeSecret: "mockport_slack_token"},
		},
	}
	if err := config.Validate(&cfg); err != nil {
		t.Fatalf("validate config: %v", err)
	}
	reg := adapter.NewRegistry()
	reg.Register(stripe.New())
	reg.Register(openai.New())
	reg.Register(githuboauth.New())
	reg.Register(slack.New())
	handler, err := NewConfiguredHandler(cfg, reg, report.NewRecorder())
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	for _, tc := range []struct {
		method string
		path   string
		status int
	}{
		{http.MethodPost, "/stripe/v1/checkout/sessions", http.StatusOK},
		{http.MethodGet, "/openai/v1/models", http.StatusOK},
		{http.MethodGet, "/github/login/oauth/authorize", http.StatusFound},
		{http.MethodPost, "/slack/api/auth.test", http.StatusOK},
	} {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
		if rec.Code != tc.status {
			t.Fatalf("%s %s status = %d, want %d", tc.method, tc.path, rec.Code, tc.status)
		}
	}

	reportRec := httptest.NewRecorder()
	handler.ServeHTTP(reportRec, httptest.NewRequest(http.MethodGet, "/_mockport/report", nil))
	var snapshot report.Snapshot
	if err := json.Unmarshal(reportRec.Body.Bytes(), &snapshot); err != nil {
		t.Fatalf("decode report: %v", err)
	}
	if len(snapshot.Adapters) != 4 {
		t.Fatalf("adapter count = %d, want 4: %#v", len(snapshot.Adapters), snapshot.Adapters)
	}
	if len(snapshot.ScenarioCoverage) != 4 {
		t.Fatalf("coverage count = %d, want 4", len(snapshot.ScenarioCoverage))
	}
}
