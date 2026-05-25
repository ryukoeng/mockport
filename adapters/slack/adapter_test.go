package slack

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

func TestAuthTest(t *testing.T) {
	rec := performRequest(t, adapter.Config{BasePath: "/slack", Scenario: "message_success"}, http.MethodPost, "/slack/api/auth.test")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode auth.test: %v", err)
	}
	if body["ok"] != true {
		t.Fatalf("ok = %v, want true", body["ok"])
	}
}

func TestPostMessageScenarios(t *testing.T) {
	tests := []struct {
		name     string
		scenario string
		status   int
		ok       bool
	}{
		{"success", "message_success", http.StatusOK, true},
		{"auth", "auth_error", http.StatusUnauthorized, false},
		{"rate", "rate_limited", http.StatusTooManyRequests, false},
		{"delivery", "delivery_failed", http.StatusBadGateway, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := performRequest(t, adapter.Config{BasePath: "/slack", Scenario: tt.scenario}, http.MethodPost, "/slack/api/chat.postMessage")
			if rec.Code != tt.status {
				t.Fatalf("status = %d, want %d", rec.Code, tt.status)
			}
		})
	}
}

func TestMetadata(t *testing.T) {
	meta := New().Metadata()
	if meta.Name != "slack" || meta.Maturity != "experimental" {
		t.Fatalf("metadata = %#v", meta)
	}
}

func performRequest(t *testing.T, cfg adapter.Config, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	if err := New().Register(mux, cfg); err != nil {
		t.Fatalf("register adapter: %v", err)
	}
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}
