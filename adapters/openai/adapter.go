package openai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

type Adapter struct{}

func New() Adapter { return Adapter{} }

func (a Adapter) Name() string { return "openai" }

func (a Adapter) Register(mux *http.ServeMux, cfg adapter.Config) error {
	basePath := cfg.BasePath
	if basePath == "" {
		basePath = "/openai"
	}
	r := &routes{basePath: strings.TrimRight(basePath, "/"), cfg: cfg}
	mux.HandleFunc(r.basePath+"/", r.handle)
	return nil
}

func (a Adapter) FakeEnv(cfg adapter.Config) map[string]string {
	basePath := cfg.BasePath
	if basePath == "" {
		basePath = "/openai"
	}
	return map[string]string{
		"OPENAI_BASE_URL": "http://localhost:43101" + basePath + "/v1",
		"OPENAI_API_KEY":  "mockport_openai_key",
	}
}

func (a Adapter) Metadata() adapter.Metadata {
	scenarios := []adapter.Scenario{
		{Name: "chat_success", Supported: true},
		{Name: "stream_success", Supported: true},
		{Name: "rate_limited", Supported: true},
		{Name: "context_length_exceeded", Supported: true},
		{Name: "auth_error", Supported: true},
	}
	return adapter.Metadata{
		Name:         "openai",
		Maturity:     "experimental",
		Capabilities: []string{"models", "chat_completions", "responses"},
		Scenarios:    scenarios,
		Endpoints: []adapter.Endpoint{
			{Method: http.MethodGet, Path: "/openai/v1/models", SupportedScenarios: []string{"chat_success"}, Notes: "OpenAI-like model list"},
			{Method: http.MethodPost, Path: "/openai/v1/chat/completions", SupportedScenarios: []string{"chat_success", "stream_success", "rate_limited", "context_length_exceeded", "auth_error"}, Notes: "OpenAI-like chat completion"},
			{Method: http.MethodPost, Path: "/openai/v1/responses", SupportedScenarios: []string{"chat_success", "rate_limited", "context_length_exceeded", "auth_error"}, Notes: "OpenAI-like responses endpoint"},
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
	case req.Method == http.MethodGet && path == "/v1/models":
		writeJSON(w, http.StatusOK, map[string]interface{}{"object": "list", "data": []map[string]string{{"id": "gpt-mockport", "object": "model"}}})
	case req.Method == http.MethodPost && path == "/v1/chat/completions":
		r.writeCompletion(w, "chat.completion")
	case req.Method == http.MethodPost && path == "/v1/responses":
		r.writeCompletion(w, "response")
	default:
		http.NotFound(w, req)
	}
}

func (r *routes) writeCompletion(w http.ResponseWriter, object string) {
	switch normalizeScenario(r.cfg.Scenario) {
	case "auth_error":
		writeError(w, http.StatusUnauthorized, "invalid_api_key", "Mockport simulated invalid API key")
	case "rate_limited":
		writeError(w, http.StatusTooManyRequests, "rate_limited", "Mockport simulated rate limit")
	case "context_length_exceeded":
		writeError(w, http.StatusBadRequest, "context_length_exceeded", "Mockport simulated context length error")
	case "stream_success":
		if object == "chat.completion" {
			writeChatCompletionStream(w)
			return
		}
		writeJSON(w, http.StatusOK, completionBody(object))
	default:
		writeJSON(w, http.StatusOK, completionBody(object))
	}
}

func completionBody(object string) map[string]interface{} {
	return map[string]interface{}{
		"id":     "mockport_openai_response",
		"object": object,
		"choices": []map[string]interface{}{{
			"index": 0,
			"message": map[string]string{
				"role":    "assistant",
				"content": "Mockport response",
			},
		}},
	}
}

func writeChatCompletionStream(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	chunk := map[string]interface{}{
		"id":     "mockport_openai_response",
		"object": "chat.completion.chunk",
		"choices": []map[string]interface{}{{
			"index": 0,
			"delta": map[string]string{
				"content": "Mockport response",
			},
		}},
	}
	data, _ := json.Marshal(chunk)
	_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
	_, _ = fmt.Fprint(w, "data: [DONE]\n\n")
}

func normalizeScenario(s string) string {
	if s == "" {
		return "chat_success"
	}
	return s
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]interface{}{"error": map[string]string{"type": "mockport_error", "code": code, "message": message}})
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
