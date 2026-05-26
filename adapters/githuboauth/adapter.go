package githuboauth

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/albert-einshutoin/mockport/internal/adapter"
	"github.com/albert-einshutoin/mockport/internal/state"
)

type Adapter struct{}

func New() Adapter { return Adapter{} }

func (a Adapter) Name() string { return "github-oauth" }

func (a Adapter) Register(mux *http.ServeMux, cfg adapter.Config) error {
	basePath := cfg.BasePath
	if basePath == "" {
		basePath = "/github"
	}
	r := &routes{basePath: strings.TrimRight(basePath, "/"), cfg: cfg, store: state.NewStore()}
	mux.HandleFunc(r.basePath+"/", r.handle)
	return nil
}

func (a Adapter) FakeEnv(cfg adapter.Config) map[string]string {
	basePath := cfg.BasePath
	if basePath == "" {
		basePath = "/github"
	}
	return map[string]string{
		"GITHUB_OAUTH_BASE_URL":      "http://localhost:43101" + basePath,
		"GITHUB_OAUTH_CLIENT_ID":     "mockport_github_client",
		"GITHUB_OAUTH_CLIENT_SECRET": "mockport_github_secret",
	}
}

func (a Adapter) Metadata() adapter.Metadata {
	return adapter.Metadata{
		Name:            "github-oauth",
		Maturity:        "workflow-compatible",
		ProviderVersion: "2022-11-28",
		Levels:          []string{"wire", "client", "workflow", "state", "error"},
		Capabilities:    []string{"oauth_authorize", "oauth_token", "user_profile", "user_emails", "user_orgs"},
		StatefulResources: []string{
			"oauth_code",
			"oauth_token",
			"user_identity",
		},
		Reset: true,
		Scenarios: []adapter.Scenario{
			{Name: "oauth_success", Supported: true},
			{Name: "invalid_code", Supported: true},
			{Name: "expired_token", Supported: true},
			{Name: "scope_missing", Supported: true},
			{Name: "redirect_uri_mismatch", Supported: true},
		},
		Endpoints: []adapter.Endpoint{
			{Method: http.MethodGet, Path: "/github/login/oauth/authorize", SupportedScenarios: []string{"oauth_success", "invalid_code", "redirect_uri_mismatch"}, Notes: "GitHub OAuth-like authorize redirect"},
			{Method: http.MethodPost, Path: "/github/login/oauth/access_token", SupportedScenarios: []string{"oauth_success", "invalid_code", "expired_token", "scope_missing", "redirect_uri_mismatch"}, Notes: "GitHub OAuth-like token exchange"},
			{Method: http.MethodGet, Path: "/github/user", SupportedScenarios: []string{"oauth_success", "expired_token", "scope_missing"}, Notes: "GitHub-like user profile"},
			{Method: http.MethodGet, Path: "/github/user/emails", SupportedScenarios: []string{"oauth_success", "scope_missing"}, Notes: "GitHub-like user emails subset"},
			{Method: http.MethodGet, Path: "/github/user/orgs", SupportedScenarios: []string{"oauth_success", "scope_missing"}, Notes: "GitHub-like user orgs subset"},
		},
	}
}

type routes struct {
	basePath string
	cfg      adapter.Config
	store    *state.Store
}

func (r *routes) handle(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, r.basePath)
	switch {
	case req.Method == http.MethodGet && path == "/login/oauth/authorize":
		redirectURI := req.URL.Query().Get("redirect_uri")
		if redirectURI == "" {
			redirectURI = "http://localhost/callback"
		}
		scope := req.URL.Query().Get("scope")
		if unsupportedScope(scope) {
			redirectWithQuery(w, req, redirectURI, map[string]string{
				"error":             "unsupported_scope",
				"error_description": "Mockport simulated unsupported scope",
				"state":             req.URL.Query().Get("state"),
			})
			return
		}
		code := r.createCode(req)
		redirectWithQuery(w, req, redirectURI, map[string]string{"code": code, "state": req.URL.Query().Get("state")})
	case req.Method == http.MethodPost && path == "/login/oauth/access_token":
		r.writeToken(w, req)
	case req.Method == http.MethodGet && path == "/user":
		r.writeUser(w, req)
	case req.Method == http.MethodGet && path == "/user/emails":
		r.writeEmails(w, req)
	case req.Method == http.MethodGet && path == "/user/orgs":
		r.writeOrgs(w, req)
	default:
		http.NotFound(w, req)
	}
}

func (r *routes) createCode(req *http.Request) string {
	scope := req.URL.Query().Get("scope")
	if scope == "" {
		scope = "read:user"
	}
	resource, _ := r.store.Create("github-oauth", "oauth_code", map[string]any{
		"client_id":    req.URL.Query().Get("client_id"),
		"redirect_uri": req.URL.Query().Get("redirect_uri"),
		"scope":        scope,
		"user":         "mockport-user",
		"expires_at":   "2999-01-01T00:00:00Z",
	})
	return resource.ID
}

func (r *routes) writeToken(w http.ResponseWriter, req *http.Request) {
	switch normalizeScenario(r.cfg.Scenario) {
	case "invalid_code":
		writeOAuthError(w, http.StatusBadRequest, "bad_verification_code", "The code passed is incorrect or expired.")
	case "expired_token":
		writeOAuthError(w, http.StatusUnauthorized, "expired_token", "Mockport simulated expired token.")
	case "scope_missing":
		writeOAuthError(w, http.StatusForbidden, "scope_missing", "Mockport simulated missing scope.")
	case "redirect_uri_mismatch":
		writeOAuthError(w, http.StatusBadRequest, "redirect_uri_mismatch", "The redirect_uri does not match the authorization request.")
	default:
		_ = req.ParseForm()
		code := req.Form.Get("code")
		if code == "" {
			writeJSON(w, http.StatusOK, map[string]interface{}{"access_token": "gho_mockport", "token_type": "bearer", "scope": "read:user"})
			return
		}
		codeResource, ok := r.store.Get("github-oauth", "oauth_code", code)
		if !ok {
			writeOAuthError(w, http.StatusBadRequest, "bad_verification_code", "The code passed is incorrect or expired.")
			return
		}
		if want, _ := codeResource.Data["redirect_uri"].(string); want != "" && req.Form.Get("redirect_uri") != "" && req.Form.Get("redirect_uri") != want {
			writeOAuthError(w, http.StatusBadRequest, "redirect_uri_mismatch", "The redirect_uri does not match the authorization request.")
			return
		}
		token, _ := r.store.Create("github-oauth", "oauth_token", map[string]any{
			"client_id":  codeResource.Data["client_id"],
			"scope":      codeResource.Data["scope"],
			"user":       codeResource.Data["user"],
			"expires_at": "2999-01-01T00:00:00Z",
		})
		writeJSON(w, http.StatusOK, map[string]interface{}{"access_token": token.ID, "token_type": "bearer", "scope": codeResource.Data["scope"]})
	}
}

func (r *routes) writeUser(w http.ResponseWriter, req *http.Request) {
	switch normalizeScenario(r.cfg.Scenario) {
	case "expired_token":
		writeAPIError(w, http.StatusUnauthorized, "Bad credentials")
	case "scope_missing":
		writeAPIError(w, http.StatusForbidden, "Resource not accessible by token")
	default:
		resource, ok := r.tokenResource(req)
		if !ok {
			writeAPIError(w, http.StatusUnauthorized, "Bad credentials")
			return
		}
		body := map[string]interface{}{
			"login": "mockport-user",
			"id":    43101,
			"name":  "Mockport User",
			"email": "mockport@example.test",
			"scope": resource.Data["scope"],
		}
		writeJSON(w, http.StatusOK, body)
	}
}

func (r *routes) writeEmails(w http.ResponseWriter, req *http.Request) {
	resource, ok := r.tokenResource(req)
	if !ok {
		writeAPIError(w, http.StatusUnauthorized, "Bad credentials")
		return
	}
	if !hasScope(resource.Data["scope"], "user:email") {
		writeAPIError(w, http.StatusForbidden, "Resource not accessible by token")
		return
	}
	writeJSON(w, http.StatusOK, []map[string]interface{}{{
		"email":      "mockport@example.test",
		"primary":    true,
		"verified":   true,
		"visibility": "public",
	}})
}

func (r *routes) writeOrgs(w http.ResponseWriter, req *http.Request) {
	resource, ok := r.tokenResource(req)
	if !ok {
		writeAPIError(w, http.StatusUnauthorized, "Bad credentials")
		return
	}
	if !hasScope(resource.Data["scope"], "read:org") {
		writeAPIError(w, http.StatusForbidden, "Resource not accessible by token")
		return
	}
	writeJSON(w, http.StatusOK, []map[string]interface{}{{
		"login":       "mockport-org",
		"id":          431010,
		"description": "Mockport fake organization",
	}})
}

func (r *routes) tokenResource(req *http.Request) (state.Resource, bool) {
	token := bearerToken(req)
	if token == "" {
		return state.Resource{}, false
	}
	return r.store.Get("github-oauth", "oauth_token", token)
}

func bearerToken(r *http.Request) string {
	value := r.Header.Get("Authorization")
	return strings.TrimSpace(strings.TrimPrefix(value, "Bearer "))
}

func unsupportedScope(value string) bool {
	for _, scope := range splitScopes(value) {
		switch scope {
		case "", "read:user", "user:email", "read:org":
		default:
			return true
		}
	}
	return false
}

func hasScope(value any, want string) bool {
	for _, scope := range splitScopes(value) {
		if scope == want {
			return true
		}
	}
	return false
}

func splitScopes(value any) []string {
	scope, _ := value.(string)
	return strings.FieldsFunc(scope, func(r rune) bool {
		return r == ' ' || r == ','
	})
}

func redirectWithQuery(w http.ResponseWriter, req *http.Request, redirectURI string, values map[string]string) {
	parsed, err := url.Parse(redirectURI)
	if err != nil {
		http.Error(w, "invalid redirect_uri", http.StatusBadRequest)
		return
	}
	query := parsed.Query()
	for key, value := range values {
		if value != "" {
			query.Set(key, value)
		}
	}
	parsed.RawQuery = query.Encode()
	http.Redirect(w, req, parsed.String(), http.StatusFound)
}

func normalizeScenario(s string) string {
	if s == "" {
		return "oauth_success"
	}
	return s
}

func writeOAuthError(w http.ResponseWriter, status int, code, description string) {
	writeJSON(w, status, map[string]string{
		"error":             code,
		"error_description": description,
		"error_uri":         "https://docs.github.com/apps/oauth-apps/maintaining-oauth-apps/troubleshooting-oauth-app-access-token-request-errors",
	})
}

func writeAPIError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"message":           message,
		"documentation_url": "https://docs.github.com/rest",
		"status":            http.StatusText(status),
	})
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
