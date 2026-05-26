package slack

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

func TestAuthTest(t *testing.T) {
	rec := performRequest(t, adapter.Config{BasePath: "/slack", Scenario: "message_success"}, http.MethodPost, "/slack/api/auth.test")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode auth.test: %v", err)
	}
	if body["ok"] != true {
		t.Fatalf("ok = %v, want true", body["ok"])
	}
}

func TestPostMessageScenarios(t *testing.T) {
	tests := []struct {
		name     string
		scenario string
		status   int
		ok       bool
	}{
		{"success", "message_success", http.StatusOK, true},
		{"auth", "auth_error", http.StatusUnauthorized, false},
		{"rate", "rate_limited", http.StatusTooManyRequests, false},
		{"delivery", "delivery_failed", http.StatusBadGateway, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := performRequest(t, adapter.Config{BasePath: "/slack", Scenario: tt.scenario}, http.MethodPost, "/slack/api/chat.postMessage")
			if rec.Code != tt.status {
				t.Fatalf("status = %d, want %d", rec.Code, tt.status)
			}
		})
	}
}

func TestPostMessagePersistsConversationHistory(t *testing.T) {
	mux := newSlackMux(t, adapter.Config{BasePath: "/slack", Scenario: "message_success"})

	post := serveSlackRequest(mux, http.MethodPost, "/slack/api/chat.postMessage", "channel=C_TEST&text=hello")
	if post.Code != http.StatusOK {
		t.Fatalf("post status = %d, want %d, body=%s", post.Code, http.StatusOK, post.Body.String())
	}
	var posted map[string]any
	if err := json.Unmarshal(post.Body.Bytes(), &posted); err != nil {
		t.Fatalf("decode post: %v", err)
	}
	if posted["channel"] != "C_TEST" || posted["ts"] != "slack_message_000001" {
		t.Fatalf("posted = %#v", posted)
	}

	history := serveSlackRequest(mux, http.MethodGet, "/slack/api/conversations.history?channel=C_TEST", "")
	if history.Code != http.StatusOK {
		t.Fatalf("history status = %d, want %d, body=%s", history.Code, http.StatusOK, history.Body.String())
	}
	var body struct {
		OK       bool             `json:"ok"`
		Messages []map[string]any `json:"messages"`
	}
	if err := json.Unmarshal(history.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode history: %v", err)
	}
	if !body.OK || len(body.Messages) != 1 || body.Messages[0]["text"] != "hello" {
		t.Fatalf("history = %#v", body)
	}
}

func TestMetadata(t *testing.T) {
	meta := New().Metadata()
	if meta.Name != "slack" || meta.Maturity != "experimental" {
		t.Fatalf("metadata = %#v", meta)
	}
	if !meta.Reset || len(meta.StatefulResources) != 1 {
		t.Fatalf("state metadata = %#v", meta)
	}
}

func performRequest(t *testing.T, cfg adapter.Config, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	mux := newSlackMux(t, cfg)
	return serveSlackRequest(mux, method, path, "")
}

func newSlackMux(t *testing.T, cfg adapter.Config) *http.ServeMux {
	t.Helper()
	mux := http.NewServeMux()
	if err := New().Register(mux, cfg); err != nil {
		t.Fatalf("register adapter: %v", err)
	}
	return mux
}

func serveSlackRequest(mux http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}
