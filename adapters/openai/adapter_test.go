package openai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

func TestModels(t *testing.T) {
	rec := performRequest(t, adapter.Config{BasePath: "/openai", Scenario: "chat_success"}, http.MethodGet, "/openai/v1/models")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["object"] != "list" {
		t.Fatalf("object = %v, want list", body["object"])
	}
}

func TestChatCompletionScenarios(t *testing.T) {
	tests := []struct {
		name     string
		scenario string
		status   int
		code     string
	}{
		{"success", "chat_success", http.StatusOK, ""},
		{"rate limited", "rate_limited", http.StatusTooManyRequests, "rate_limited"},
		{"context", "context_length_exceeded", http.StatusBadRequest, "context_length_exceeded"},
		{"auth", "auth_error", http.StatusUnauthorized, "invalid_api_key"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := performRequest(t, adapter.Config{BasePath: "/openai", Scenario: tt.scenario}, http.MethodPost, "/openai/v1/chat/completions")
			if rec.Code != tt.status {
				t.Fatalf("status = %d, want %d", rec.Code, tt.status)
			}
			if tt.code != "" {
				assertErrorCode(t, rec, tt.code)
			}
		})
	}
}

func TestResponses(t *testing.T) {
	rec := performRequest(t, adapter.Config{BasePath: "/openai", Scenario: "chat_success"}, http.MethodPost, "/openai/v1/responses")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestResponsesCreateAndRetrieveState(t *testing.T) {
	mux := newOpenAIMux(t, adapter.Config{BasePath: "/openai", Scenario: "chat_success"})
	create := serveOpenAIRequest(mux, http.MethodPost, "/openai/v1/responses", `{"model":"gpt-mockport","input":"hello"}`)
	if create.Code != http.StatusOK {
		t.Fatalf("create status = %d, want %d, body=%s", create.Code, http.StatusOK, create.Body.String())
	}
	var created map[string]any
	if err := json.Unmarshal(create.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	if created["id"] != "openai_response_000001" {
		t.Fatalf("created id = %#v", created["id"])
	}

	retrieve := serveOpenAIRequest(mux, http.MethodGet, "/openai/v1/responses/openai_response_000001", "")
	if retrieve.Code != http.StatusOK {
		t.Fatalf("retrieve status = %d, want %d, body=%s", retrieve.Code, http.StatusOK, retrieve.Body.String())
	}
	var got map[string]any
	if err := json.Unmarshal(retrieve.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode retrieve: %v", err)
	}
	if got["id"] != created["id"] || got["model"] != "gpt-mockport" {
		t.Fatalf("retrieved = %#v", got)
	}
}

func TestChatCompletionUsesStatefulIDsAndValidatesRequiredFields(t *testing.T) {
	mux := newOpenAIMux(t, adapter.Config{BasePath: "/openai", Scenario: "chat_success"})
	first := serveOpenAIRequest(mux, http.MethodPost, "/openai/v1/chat/completions", `{"model":"gpt-mockport","messages":[{"role":"user","content":"hi"}]}`)
	second := serveOpenAIRequest(mux, http.MethodPost, "/openai/v1/chat/completions", `{"model":"gpt-mockport","messages":[{"role":"user","content":"hi again"}]}`)
	if first.Code != http.StatusOK || second.Code != http.StatusOK {
		t.Fatalf("statuses = %d/%d bodies=%s/%s", first.Code, second.Code, first.Body.String(), second.Body.String())
	}
	if !strings.Contains(first.Body.String(), `"id":"openai_chat_completion_000001"`) || !strings.Contains(second.Body.String(), `"id":"openai_chat_completion_000002"`) {
		t.Fatalf("stateful ids not found: %s / %s", first.Body.String(), second.Body.String())
	}

	missing := serveOpenAIRequest(mux, http.MethodPost, "/openai/v1/chat/completions", `{"model":"gpt-mockport"}`)
	if missing.Code != http.StatusBadRequest {
		t.Fatalf("missing status = %d, want %d, body=%s", missing.Code, http.StatusBadRequest, missing.Body.String())
	}
	assertErrorCode(t, missing, "missing_required_field")
}

func TestChatCompletionStreamSuccessReturnsSSE(t *testing.T) {
	rec := performRequest(t, adapter.Config{BasePath: "/openai", Scenario: "stream_success"}, http.MethodPost, "/openai/v1/chat/completions")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	contentType := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/event-stream") {
		t.Fatalf("Content-Type = %q, want text/event-stream", contentType)
	}
	body := rec.Body.String()
	for _, want := range []string{
		"data: {",
		`"object":"chat.completion.chunk"`,
		`"delta":{"content":"Mockport response"}`,
		"data: [DONE]",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q: %s", want, body)
		}
	}
}

func TestChatCompletionSuccessStillReturnsJSON(t *testing.T) {
	rec := performRequest(t, adapter.Config{BasePath: "/openai", Scenario: "chat_success"}, http.MethodPost, "/openai/v1/chat/completions")
	if contentType := rec.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json", contentType)
	}
	if strings.Contains(rec.Body.String(), "data: [DONE]") {
		t.Fatalf("non-stream response used SSE body: %s", rec.Body.String())
	}
}

func TestMetadata(t *testing.T) {
	meta := New().Metadata()
	if meta.Name != "openai" {
		t.Fatalf("name = %q", meta.Name)
	}
	if meta.Maturity != "experimental" {
		t.Fatalf("maturity = %q", meta.Maturity)
	}
	if len(meta.Scenarios) < 5 || len(meta.Endpoints) < 3 {
		t.Fatalf("metadata too small: %#v", meta)
	}
	if !meta.Reset || len(meta.StatefulResources) != 2 {
		t.Fatalf("state metadata = %#v", meta)
	}
}

func performRequest(t *testing.T, cfg adapter.Config, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	mux := newOpenAIMux(t, cfg)
	return serveOpenAIRequest(mux, method, path, "")
}

func newOpenAIMux(t *testing.T, cfg adapter.Config) *http.ServeMux {
	t.Helper()
	mux := http.NewServeMux()
	if err := New().Register(mux, cfg); err != nil {
		t.Fatalf("register adapter: %v", err)
	}
	return mux
}

func serveOpenAIRequest(mux http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func assertErrorCode(t *testing.T, rec *httptest.ResponseRecorder, want string) {
	t.Helper()
	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if body.Error.Code != want {
		t.Fatalf("error code = %q, want %q", body.Error.Code, want)
	}
}
