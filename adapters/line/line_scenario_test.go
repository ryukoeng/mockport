package line_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	lineadapter "github.com/albert-einshutoin/mockport/adapters/line"
	"github.com/albert-einshutoin/mockport/internal/adapter"
)

func newLineMuxForScenario(t *testing.T, cfg adapter.Config) *http.ServeMux {
	t.Helper()
	mux := http.NewServeMux()
	if err := lineadapter.New().Register(mux, cfg); err != nil {
		t.Fatalf("Register: %v", err)
	}
	return mux
}

func TestLINEHeaderOverridesConfig(t *testing.T) {
	// config=line_success でサーバーを起動し、ヘッダで auth_error を指定
	mux := newLineMuxForScenario(t, adapter.Config{BasePath: "/line", Scenario: "line_success"})
	body := strings.NewReader(`{"to":"Umockport","messages":[{"type":"text","text":"hello"}]}`)
	req := httptest.NewRequest(http.MethodPost, "/line/v2/bot/message/push", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Mockport-Scenario", "auth_error")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", rec.Code)
	}
}

func TestLINEUnknownScenarioReturns400(t *testing.T) {
	mux := newLineMuxForScenario(t, adapter.Config{BasePath: "/line", Scenario: "line_success"})
	body := strings.NewReader(`{"to":"Umockport","messages":[{"type":"text","text":"hello"}]}`)
	req := httptest.NewRequest(http.MethodPost, "/line/v2/bot/message/push", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Mockport-Scenario", "not_a_real_scenario")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	msg, _ := resp["message"].(string)
	if !strings.Contains(msg, "unknown_mockport_scenario") {
		t.Errorf("want message containing unknown_mockport_scenario, got %q", msg)
	}
}
