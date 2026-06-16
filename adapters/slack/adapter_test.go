package slack

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

func TestAuthTest(t *testing.T) {
	rec := performRequest(t, adapter.Config{BasePath: "/slack", Scenario: "message_success"}, http.MethodPost, "/slack/api/auth.test")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode auth.test: %v", err)
	}
	if body["ok"] != true {
		t.Fatalf("ok = %v, want true", body["ok"])
	}
	if body["team_id"] != "T_MOCKPORT" || body["bot_id"] != "B_MOCKPORT" {
		t.Fatalf("auth body = %#v", body)
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
		{"auth", "auth_error", http.StatusOK, false},
		{"rate", "rate_limited", http.StatusTooManyRequests, false},
		{"delivery", "delivery_failed", http.StatusOK, false},
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
	message, ok := posted["message"].(map[string]any)
	if !ok || message["text"] != "hello" || message["team"] != "T_MOCKPORT" {
		t.Fatalf("posted message = %#v", posted["message"])
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

func TestConversationsListUpdateAndDelete(t *testing.T) {
	mux := newSlackMux(t, adapter.Config{BasePath: "/slack", Scenario: "message_success"})

	list := serveSlackRequest(mux, http.MethodPost, "/slack/api/conversations.list", "")
	if list.Code != http.StatusOK || !strings.Contains(list.Body.String(), `"id":"C_MOCKPORT"`) {
		t.Fatalf("list status/body = %d %s", list.Code, list.Body.String())
	}

	post := serveSlackRequest(mux, http.MethodPost, "/slack/api/chat.postMessage", "channel=C_MOCKPORT&text=hello")
	var posted map[string]any
	if err := json.Unmarshal(post.Body.Bytes(), &posted); err != nil {
		t.Fatalf("decode post: %v", err)
	}
	ts, _ := posted["ts"].(string)

	update := serveSlackRequest(mux, http.MethodPost, "/slack/api/chat.update", "channel=C_MOCKPORT&ts="+ts+"&text=edited")
	if update.Code != http.StatusOK || !strings.Contains(update.Body.String(), `"text":"edited"`) {
		t.Fatalf("update status/body = %d %s", update.Code, update.Body.String())
	}

	del := serveSlackRequest(mux, http.MethodPost, "/slack/api/chat.delete", "channel=C_MOCKPORT&ts="+ts)
	if del.Code != http.StatusOK || !strings.Contains(del.Body.String(), `"ok":true`) {
		t.Fatalf("delete status/body = %d %s", del.Code, del.Body.String())
	}

	history := serveSlackRequest(mux, http.MethodPost, "/slack/api/conversations.history", "channel=C_MOCKPORT")
	if history.Code != http.StatusOK || strings.Contains(history.Body.String(), `"text":"edited"`) {
		t.Fatalf("history after delete status/body = %d %s", history.Code, history.Body.String())
	}
}

func TestSlackErrorsAndHeaders(t *testing.T) {
	auth := performRequest(t, adapter.Config{BasePath: "/slack", Scenario: "auth_error"}, http.MethodPost, "/slack/api/auth.test")
	if auth.Code != http.StatusOK {
		t.Fatalf("auth error status = %d, body=%s", auth.Code, auth.Body.String())
	}
	assertSlackError(t, auth, "invalid_auth")

	rate := performRequest(t, adapter.Config{BasePath: "/slack", Scenario: "rate_limited"}, http.MethodPost, "/slack/api/chat.postMessage")
	if rate.Code != http.StatusTooManyRequests || rate.Header().Get("Retry-After") != "1" {
		t.Fatalf("rate status/header/body = %d %q %s", rate.Code, rate.Header().Get("Retry-After"), rate.Body.String())
	}
	assertSlackError(t, rate, "ratelimited")

	mux := newSlackMux(t, adapter.Config{BasePath: "/slack", Scenario: "message_success"})
	unknownChannel := serveSlackRequest(mux, http.MethodPost, "/slack/api/chat.postMessage", "channel=C_UNKNOWN&text=hello")
	if unknownChannel.Code != http.StatusOK {
		t.Fatalf("unknown channel status = %d, body=%s", unknownChannel.Code, unknownChannel.Body.String())
	}
	assertSlackError(t, unknownChannel, "channel_not_found")

	unsupported := serveSlackRequest(mux, http.MethodPost, "/slack/api/views.open", "")
	if unsupported.Code != http.StatusNotFound {
		t.Fatalf("unsupported status = %d, body=%s", unsupported.Code, unsupported.Body.String())
	}
}

func TestEventsURLVerificationAndCallback(t *testing.T) {
	mux := newSlackMux(t, adapter.Config{BasePath: "/slack", Scenario: "message_success", FakeSecret: "mockport_slack_signing_secret"})

	challenge := `{"type":"url_verification","challenge":"challenge-123"}`
	verify := serveSlackSignedRequest(mux, http.MethodPost, "/slack/events", challenge, "mockport_slack_signing_secret")
	if verify.Code != http.StatusOK || !strings.Contains(verify.Body.String(), `"challenge":"challenge-123"`) {
		t.Fatalf("verify status/body = %d %s", verify.Code, verify.Body.String())
	}

	callback := `{"type":"event_callback","event":{"type":"message","channel":"C_MOCKPORT","text":"event hello","user":"U_EVENT"}}`
	event := serveSlackSignedRequest(mux, http.MethodPost, "/slack/events", callback, "mockport_slack_signing_secret")
	if event.Code != http.StatusOK || !strings.Contains(event.Body.String(), `"ok":true`) {
		t.Fatalf("event status/body = %d %s", event.Code, event.Body.String())
	}

	history := serveSlackRequest(mux, http.MethodPost, "/slack/api/conversations.history", "channel=C_MOCKPORT")
	if !strings.Contains(history.Body.String(), `"text":"event hello"`) {
		t.Fatalf("history after event = %s", history.Body.String())
	}

	bad := serveSlackSignedRequest(mux, http.MethodPost, "/slack/events", callback, "wrong_secret")
	if bad.Code != http.StatusUnauthorized {
		t.Fatalf("bad signature status = %d, body=%s", bad.Code, bad.Body.String())
	}
	assertSlackError(t, bad, "invalid_signature")

	old := serveSlackSignedRequestWithTimestamp(mux, http.MethodPost, "/slack/events", callback, "mockport_slack_signing_secret", strconv.FormatInt(time.Now().Add(-10*time.Minute).Unix(), 10))
	if old.Code != http.StatusUnauthorized {
		t.Fatalf("old signature status = %d, body=%s", old.Code, old.Body.String())
	}
	assertSlackError(t, old, "invalid_signature")
}

func TestMetadata(t *testing.T) {
	meta := New().Metadata()
	if meta.Name != "slack" || meta.Maturity != "workflow-compatible" {
		t.Fatalf("metadata = %#v", meta)
	}
	if meta.ProviderVersion != "2025-02-01" || len(meta.Levels) < 5 || len(meta.Endpoints) < 6 {
		t.Fatalf("compat metadata = %#v", meta)
	}
	if !meta.Reset || len(meta.StatefulResources) != 4 {
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

func serveSlackSignedRequest(mux http.Handler, method, path, body, secret string) *httptest.ResponseRecorder {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	return serveSlackSignedRequestWithTimestamp(mux, method, path, body, secret, timestamp)
}

func serveSlackSignedRequestWithTimestamp(mux http.Handler, method, path, body, secret, timestamp string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Slack-Request-Timestamp", timestamp)
	req.Header.Set("X-Slack-Signature", slackSignature(secret, timestamp, body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func slackSignature(secret, timestamp, body string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte("v0:" + timestamp + ":" + body))
	return "v0=" + hex.EncodeToString(mac.Sum(nil))
}

func assertSlackError(t *testing.T, rec *httptest.ResponseRecorder, want string) {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if body["ok"] != false || body["error"] != want {
		t.Fatalf("error body = %#v", body)
	}
}
