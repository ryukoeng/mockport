package githuboauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

func TestOAuthCodeTokenAndUserAreStateful(t *testing.T) {
	mux := newGitHubMux(t, adapter.Config{BasePath: "/github", Scenario: "oauth_success"})

	auth := serveGitHubRequest(mux, http.MethodGet, "/github/login/oauth/authorize?redirect_uri=http://app/callback&state=s1&scope=read:user,user:email", "", nil)
	if auth.Code != http.StatusFound {
		t.Fatalf("authorize status = %d, want %d", auth.Code, http.StatusFound)
	}
	location := auth.Header().Get("Location")
	parsed, err := url.Parse(location)
	if err != nil {
		t.Fatalf("parse redirect location: %v", err)
	}
	code := parsed.Query().Get("code")
	if code != "github_oauth_oauth_code_000001" {
		t.Fatalf("code = %q", code)
	}

	token := serveGitHubRequest(mux, http.MethodPost, "/github/login/oauth/access_token", "code="+code, map[string]string{"Accept": "application/json"})
	if token.Code != http.StatusOK {
		t.Fatalf("token status = %d, want %d, body=%s", token.Code, http.StatusOK, token.Body.String())
	}
	var tokenBody map[string]any
	if err := json.Unmarshal(token.Body.Bytes(), &tokenBody); err != nil {
		t.Fatalf("decode token: %v", err)
	}
	if tokenBody["access_token"] != "github_oauth_oauth_token_000001" || tokenBody["scope"] != "read:user,user:email" {
		t.Fatalf("token body = %#v", tokenBody)
	}

	user := serveGitHubRequest(mux, http.MethodGet, "/github/user", "", map[string]string{"Authorization": "Bearer github_oauth_oauth_token_000001"})
	if user.Code != http.StatusOK {
		t.Fatalf("user status = %d, want %d, body=%s", user.Code, http.StatusOK, user.Body.String())
	}
	var userBody map[string]any
	if err := json.Unmarshal(user.Body.Bytes(), &userBody); err != nil {
		t.Fatalf("decode user: %v", err)
	}
	if userBody["login"] != "mockport-user" || userBody["scope"] != "read:user,user:email" {
		t.Fatalf("user body = %#v", userBody)
	}
}

func TestOAuthRejectsUnknownCode(t *testing.T) {
	mux := newGitHubMux(t, adapter.Config{BasePath: "/github", Scenario: "oauth_success"})
	rec := serveGitHubRequest(mux, http.MethodPost, "/github/login/oauth/access_token", "code=unknown", nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestMetadata(t *testing.T) {
	meta := New().Metadata()
	if meta.Name != "github-oauth" || meta.Maturity != "experimental" {
		t.Fatalf("metadata = %#v", meta)
	}
	if !meta.Reset || len(meta.StatefulResources) != 3 {
		t.Fatalf("state metadata = %#v", meta)
	}
}

func performRequest(t *testing.T, cfg adapter.Config, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	mux := newGitHubMux(t, cfg)
	return serveGitHubRequest(mux, method, path, "", nil)
}

func newGitHubMux(t *testing.T, cfg adapter.Config) *http.ServeMux {
	t.Helper()
	mux := http.NewServeMux()
	if err := New().Register(mux, cfg); err != nil {
		t.Fatalf("register adapter: %v", err)
	}
	return mux
}

func serveGitHubRequest(mux http.Handler, method, path, body string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}
