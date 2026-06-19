package zohooauth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	zohoadapter "github.com/albert-einshutoin/mockport/adapters/zohooauth"
	"github.com/albert-einshutoin/mockport/internal/adapter"
)

func newZohoMuxForScenario(t *testing.T, cfg adapter.Config) *http.ServeMux {
	t.Helper()
	mux := http.NewServeMux()
	if err := zohoadapter.New().Register(mux, cfg); err != nil {
		t.Fatalf("Register: %v", err)
	}
	return mux
}

func TestZohoOAuthHeaderOverridesConfig(t *testing.T) {
	mux := newZohoMuxForScenario(t, adapter.Config{BasePath: "/zoho", Scenario: "oauth_success"})
	form := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {"somecode"},
		"client_id":    {"mockport_zoho_client"},
		"redirect_uri": {"http://localhost/callback"},
	}.Encode()
	req := httptest.NewRequest(http.MethodPost, "/zoho/oauth/v2/token", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Mockport-Scenario", "invalid_code")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	// Zoho は 200 でエラーを返す
	if rec.Code != http.StatusOK {
		t.Errorf("want 200, got %d", rec.Code)
	}
	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	errStr, _ := body["error"].(string)
	if errStr != "invalid_code" {
		t.Errorf("want error=invalid_code, got %q", errStr)
	}
}

func TestZohoOAuthUnknownScenarioReturns400(t *testing.T) {
	mux := newZohoMuxForScenario(t, adapter.Config{BasePath: "/zoho", Scenario: "oauth_success"})
	form := url.Values{
		"grant_type": {"authorization_code"},
		"code":       {"somecode"},
	}.Encode()
	req := httptest.NewRequest(http.MethodPost, "/zoho/oauth/v2/token", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Mockport-Scenario", "totally_fake")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	errStr, _ := body["error"].(string)
	if errStr != "unknown_mockport_scenario" {
		t.Errorf("want error=unknown_mockport_scenario, got %q", errStr)
	}
}
