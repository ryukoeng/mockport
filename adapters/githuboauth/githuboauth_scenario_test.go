package githuboauth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	githuboauthadapter "github.com/albert-einshutoin/mockport/adapters/githuboauth"
	"github.com/albert-einshutoin/mockport/internal/adapter"
)

func newGitHubMuxForScenario(t *testing.T, cfg adapter.Config) *http.ServeMux {
	t.Helper()
	mux := http.NewServeMux()
	if err := githuboauthadapter.New().Register(mux, cfg); err != nil {
		t.Fatalf("Register: %v", err)
	}
	return mux
}

func TestGitHubOAuthHeaderOverridesConfig(t *testing.T) {
	mux := newGitHubMuxForScenario(t, adapter.Config{BasePath: "/github", Scenario: "oauth_success"})
	form := url.Values{"code": {"somecode"}, "client_id": {"mockport_github_client"}}.Encode()
	req := httptest.NewRequest(http.MethodPost, "/github/login/oauth/access_token", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Mockport-Scenario", "invalid_code")
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
	if errStr != "bad_verification_code" {
		t.Errorf("want error=bad_verification_code, got %q", errStr)
	}
}

func TestGitHubOAuthUnknownScenarioReturns400(t *testing.T) {
	mux := newGitHubMuxForScenario(t, adapter.Config{BasePath: "/github", Scenario: "oauth_success"})
	form := url.Values{"code": {"somecode"}, "client_id": {"mockport_github_client"}}.Encode()
	req := httptest.NewRequest(http.MethodPost, "/github/login/oauth/access_token", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Mockport-Scenario", "definitely_fake")
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
