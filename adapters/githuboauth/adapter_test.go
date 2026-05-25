package githuboauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

func TestAuthorizeRedirect(t *testing.T) {
	rec := performRequest(t, adapter.Config{BasePath: "/github", Scenario: "oauth_success"}, http.MethodGet, "/github/login/oauth/authorize?redirect_uri=http://app/callback&state=s1")
	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusFound)
	}
	if rec.Header().Get("Location") == "" {
		t.Fatal("missing Location header")
	}
}

func TestAccessTokenScenarios(t *testing.T) {
	tests := []struct {
		name     string
		scenario string
		status   int
	}{
		{"success", "oauth_success", http.StatusOK},
		{"invalid", "invalid_code", http.StatusBadRequest},
		{"expired", "expired_token", http.StatusUnauthorized},
		{"scope", "scope_missing", http.StatusForbidden},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := performRequest(t, adapter.Config{BasePath: "/github", Scenario: tt.scenario}, http.MethodPost, "/github/login/oauth/access_token")
			if rec.Code != tt.status {
				t.Fatalf("status = %d, want %d", rec.Code, tt.status)
			}
		})
	}
}

func TestUser(t *testing.T) {
	rec := performRequest(t, adapter.Config{BasePath: "/github", Scenario: "oauth_success"}, http.MethodGet, "/github/user")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode user: %v", err)
	}
	if body["login"] != "mockport-user" {
		t.Fatalf("login = %v", body["login"])
	}
}

func TestMetadata(t *testing.T) {
	meta := New().Metadata()
	if meta.Name != "github-oauth" || meta.Maturity != "experimental" {
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
