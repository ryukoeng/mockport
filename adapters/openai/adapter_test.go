package openai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
