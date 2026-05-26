package openai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/albert-einshutoin/mockport/internal/adapter"
	"github.com/albert-einshutoin/mockport/internal/state"
)

type Adapter struct{}

func New() Adapter { return Adapter{} }

func (a Adapter) Name() string { return "openai" }

func (a Adapter) Register(mux *http.ServeMux, cfg adapter.Config) error {
	basePath := cfg.BasePath
	if basePath == "" {
		basePath = "/openai"
	}
	r := &routes{basePath: strings.TrimRight(basePath, "/"), cfg: cfg, store: state.NewStore()}
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
		StatefulResources: []string{
			"chat_completion",
			"response",
		},
		Reset: true,
		Endpoints: []adapter.Endpoint{
			{Method: http.MethodGet, Path: "/openai/v1/models", SupportedScenarios: []string{"chat_success"}, Notes: "OpenAI-like model list"},
			{Method: http.MethodPost, Path: "/openai/v1/chat/completions", SupportedScenarios: []string{"chat_success", "stream_success", "rate_limited", "context_length_exceeded", "auth_error"}, Notes: "OpenAI-like chat completion"},
			{Method: http.MethodPost, Path: "/openai/v1/responses", SupportedScenarios: []string{"chat_success", "rate_limited", "context_length_exceeded", "auth_error"}, Notes: "OpenAI-like responses endpoint"},
			{Method: http.MethodGet, Path: "/openai/v1/responses/{id}", SupportedScenarios: []string{"chat_success"}, Notes: "Deterministic response lookup"},
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
	case req.Method == http.MethodGet && path == "/v1/models":
		writeJSON(w, http.StatusOK, map[string]interface{}{"object": "list", "data": []map[string]string{{"id": "gpt-mockport", "object": "model"}}})
	case req.Method == http.MethodPost && path == "/v1/chat/completions":
		r.writeCompletion(w, req, "chat.completion")
	case req.Method == http.MethodPost && path == "/v1/responses":
		r.writeCompletion(w, req, "response")
	case req.Method == http.MethodGet && strings.HasPrefix(path, "/v1/responses/"):
		r.writeResponseLookup(w, strings.TrimPrefix(path, "/v1/responses/"))
	default:
		http.NotFound(w, req)
	}
}

func (r *routes) writeCompletion(w http.ResponseWriter, req *http.Request, object string) {
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
		r.writeStatefulCompletion(w, req, object)
	default:
		r.writeStatefulCompletion(w, req, object)
	}
}

func (r *routes) writeStatefulCompletion(w http.ResponseWriter, req *http.Request, object string) {
	payload, err := decodePayload(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "Request body must be JSON")
		return
	}
	if len(payload) == 0 {
		payload["model"] = "gpt-mockport"
		if object == "chat.completion" {
			payload["messages"] = []any{map[string]any{"role": "user", "content": "Mockport request"}}
		} else {
			payload["input"] = "Mockport request"
		}
	}
	required := []string{"model"}
	resourceType := "response"
	if object == "chat.completion" {
		required = append(required, "messages")
		resourceType = "chat_completion"
	} else {
		required = append(required, "input")
	}
	if err := state.RequireFields(payload, required...); err != nil {
		writeError(w, http.StatusBadRequest, "missing_required_field", err.Error())
		return
	}

	body := completionBody(object)
	body["model"] = payload["model"]
	resource, err := r.store.Create("openai", resourceType, body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "mockport_state_error", err.Error())
		return
	}
	response := resource.Data
	response["id"] = resource.ID
	writeJSON(w, http.StatusOK, response)
}

func (r *routes) writeResponseLookup(w http.ResponseWriter, id string) {
	if resource, ok := r.store.Get("openai", "response", id); ok {
		body := resource.Data
		body["id"] = resource.ID
		writeJSON(w, http.StatusOK, body)
		return
	}
	writeError(w, http.StatusNotFound, "not_found", "Mockport response not found")
}

func completionBody(object string) map[string]interface{} {
	body := map[string]interface{}{
		"object": object,
		"choices": []map[string]interface{}{{
			"index": 0,
			"message": map[string]string{
				"role":    "assistant",
				"content": "Mockport response",
			},
		}},
	}
	if object == "response" {
		body["output_text"] = "Mockport response"
	}
	return body
}

func decodePayload(req *http.Request) (map[string]any, error) {
	if req.Body == nil {
		return map[string]any{}, nil
	}
	var payload map[string]any
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		if err.Error() == "EOF" {
			return map[string]any{}, nil
		}
		return nil, err
	}
	if payload == nil {
		payload = map[string]any{}
	}
	return payload, nil
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
