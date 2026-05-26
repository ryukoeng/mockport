package githuboauth

import (
	"encoding/json"
	"net/http"
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
		Name:         "github-oauth",
		Maturity:     "experimental",
		Capabilities: []string{"oauth_authorize", "oauth_token", "user_profile"},
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
		},
		Endpoints: []adapter.Endpoint{
			{Method: http.MethodGet, Path: "/github/login/oauth/authorize", SupportedScenarios: []string{"oauth_success", "invalid_code"}, Notes: "GitHub OAuth-like authorize redirect"},
			{Method: http.MethodPost, Path: "/github/login/oauth/access_token", SupportedScenarios: []string{"oauth_success", "invalid_code", "expired_token", "scope_missing"}, Notes: "GitHub OAuth-like token exchange"},
			{Method: http.MethodGet, Path: "/github/user", SupportedScenarios: []string{"oauth_success", "expired_token", "scope_missing"}, Notes: "GitHub-like user profile"},
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
		code := r.createCode(req)
		sep := "?"
		if strings.Contains(redirectURI, "?") {
			sep = "&"
		}
		http.Redirect(w, req, redirectURI+sep+"code="+code+"&state="+req.URL.Query().Get("state"), http.StatusFound)
	case req.Method == http.MethodPost && path == "/login/oauth/access_token":
		r.writeToken(w, req)
	case req.Method == http.MethodGet && path == "/user":
		r.writeUser(w, req)
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
		"scope": scope,
		"user":  "mockport-user",
	})
	return resource.ID
}

func (r *routes) writeToken(w http.ResponseWriter, req *http.Request) {
	switch normalizeScenario(r.cfg.Scenario) {
	case "invalid_code":
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_verification_code"})
	case "expired_token":
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "expired_token"})
	case "scope_missing":
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "scope_missing"})
	default:
		_ = req.ParseForm()
		code := req.Form.Get("code")
		if code == "" {
			writeJSON(w, http.StatusOK, map[string]interface{}{"access_token": "gho_mockport", "token_type": "bearer", "scope": "read:user"})
			return
		}
		codeResource, ok := r.store.Get("github-oauth", "oauth_code", code)
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_verification_code"})
			return
		}
		token, _ := r.store.Create("github-oauth", "oauth_token", map[string]any{
			"scope": codeResource.Data["scope"],
			"user":  codeResource.Data["user"],
		})
		writeJSON(w, http.StatusOK, map[string]interface{}{"access_token": token.ID, "token_type": "bearer", "scope": codeResource.Data["scope"]})
	}
}

func (r *routes) writeUser(w http.ResponseWriter, req *http.Request) {
	switch normalizeScenario(r.cfg.Scenario) {
	case "expired_token":
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Mockport simulated expired token"})
	case "scope_missing":
		writeJSON(w, http.StatusForbidden, map[string]string{"message": "Mockport simulated missing scope"})
	default:
		body := map[string]interface{}{"login": "mockport-user", "id": 43101, "name": "Mockport User"}
		if token := bearerToken(req); token != "" {
			if resource, ok := r.store.Get("github-oauth", "oauth_token", token); ok {
				body["scope"] = resource.Data["scope"]
			}
		}
		writeJSON(w, http.StatusOK, body)
	}
}

func bearerToken(r *http.Request) string {
	value := r.Header.Get("Authorization")
	return strings.TrimSpace(strings.TrimPrefix(value, "Bearer "))
}

func normalizeScenario(s string) string {
	if s == "" {
		return "oauth_success"
	}
	return s
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
