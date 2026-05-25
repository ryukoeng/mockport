package slack

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

type Adapter struct{}

func New() Adapter { return Adapter{} }

func (a Adapter) Name() string { return "slack" }

func (a Adapter) Register(mux *http.ServeMux, cfg adapter.Config) error {
	basePath := cfg.BasePath
	if basePath == "" {
		basePath = "/slack"
	}
	r := &routes{basePath: strings.TrimRight(basePath, "/"), cfg: cfg}
	mux.HandleFunc(r.basePath+"/", r.handle)
	return nil
}

func (a Adapter) FakeEnv(cfg adapter.Config) map[string]string {
	basePath := cfg.BasePath
	if basePath == "" {
		basePath = "/slack"
	}
	return map[string]string{
		"SLACK_API_URL":   "http://localhost:43101" + basePath + "/api",
		"SLACK_BOT_TOKEN": "mockport_slack_token",
	}
}

func (a Adapter) Metadata() adapter.Metadata {
	return adapter.Metadata{
		Name:         "slack",
		Maturity:     "experimental",
		Capabilities: []string{"auth_test", "chat_post_message"},
		Scenarios: []adapter.Scenario{
			{Name: "message_success", Supported: true},
			{Name: "auth_error", Supported: true},
			{Name: "rate_limited", Supported: true},
			{Name: "delivery_failed", Supported: true},
		},
		Endpoints: []adapter.Endpoint{
			{Method: http.MethodPost, Path: "/slack/api/auth.test", SupportedScenarios: []string{"message_success", "auth_error"}, Notes: "Slack-like auth test"},
			{Method: http.MethodPost, Path: "/slack/api/chat.postMessage", SupportedScenarios: []string{"message_success", "auth_error", "rate_limited", "delivery_failed"}, Notes: "Slack-like message post"},
		},
	}
}

type routes struct {
	basePath string
	cfg      adapter.Config
}

func (r *routes) handle(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, r.basePath)
	switch {
	case req.Method == http.MethodPost && path == "/api/auth.test":
		r.writeAuthTest(w)
	case req.Method == http.MethodPost && path == "/api/chat.postMessage":
		r.writePostMessage(w)
	default:
		http.NotFound(w, req)
	}
}

func (r *routes) writeAuthTest(w http.ResponseWriter) {
	if normalizeScenario(r.cfg.Scenario) == "auth_error" {
		writeJSON(w, http.StatusUnauthorized, map[string]interface{}{"ok": false, "error": "invalid_auth"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"ok": true, "user_id": "U_MOCKPORT", "team": "Mockport"})
}

func (r *routes) writePostMessage(w http.ResponseWriter) {
	switch normalizeScenario(r.cfg.Scenario) {
	case "auth_error":
		writeJSON(w, http.StatusUnauthorized, map[string]interface{}{"ok": false, "error": "invalid_auth"})
	case "rate_limited":
		writeJSON(w, http.StatusTooManyRequests, map[string]interface{}{"ok": false, "error": "ratelimited"})
	case "delivery_failed":
		writeJSON(w, http.StatusBadGateway, map[string]interface{}{"ok": false, "error": "message_delivery_failed"})
	default:
		writeJSON(w, http.StatusOK, map[string]interface{}{"ok": true, "channel": "C_MOCKPORT", "ts": "1710000000.000001"})
	}
}

func normalizeScenario(s string) string {
	if s == "" {
		return "message_success"
	}
	return s
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
