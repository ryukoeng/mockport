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
	// Messaging API のエラー本文は {"message": ...} 形式で専用コードフィールドが無いため、
	// 共通コードは message 先頭の固定プレフィックスとして判別する（厳密にプレフィックス一致）。
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
	if !strings.HasPrefix(msg, "unknown_mockport_scenario:") {
		t.Errorf("want message prefixed with unknown_mockport_scenario:, got %q", msg)
	}
}

// TestLINEOAuthUnknownScenarioUsesErrorCode は OAuth エンドポイント（{"error": ...}
// 形式）で共通コードが機械可読な error フィールドへ入ることを厳密一致で固定する。
func TestLINEOAuthUnknownScenarioUsesErrorCode(t *testing.T) {
	mux := newLineMuxForScenario(t, adapter.Config{BasePath: "/line", Scenario: "line_success"})
	req := httptest.NewRequest(http.MethodPost, "/line/oauth2/v2.1/token", strings.NewReader("grant_type=authorization_code&code=x"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Mockport-Scenario", "not_a_real_scenario")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rec.Code)
	}
	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if code, _ := resp["error"].(string); code != "unknown_mockport_scenario" {
		t.Errorf("want error=unknown_mockport_scenario, got %q", code)
	}
}

// TestLINEPayUnknownScenarioReturnsCommonCode は LINE Pay エンドポイント
// （{"returnCode": ..., "returnMessage": ...} 形式）で returnCode の数値契約を壊さず、
// 共通コードが returnMessage 先頭プレフィックスで判別できることを固定する。
func TestLINEPayUnknownScenarioReturnsCommonCode(t *testing.T) {
	mux := newLineMuxForScenario(t, adapter.Config{BasePath: "/line", Scenario: "line_success"})
	req := httptest.NewRequest(http.MethodPost, "/line/v3/payments/request", strings.NewReader(`{"amount":1000}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Mockport-Scenario", "not_a_real_scenario")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if code, _ := resp["returnCode"].(string); code != "2101" {
		t.Errorf("want returnCode=2101 (numeric contract preserved), got %q", code)
	}
	msg, _ := resp["returnMessage"].(string)
	if !strings.HasPrefix(msg, "unknown_mockport_scenario:") {
		t.Errorf("want returnMessage prefixed with unknown_mockport_scenario:, got %q", msg)
	}
}
