package line

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

func TestMessagingPushAndProfile(t *testing.T) {
	mux := newLineMux(t, adapter.Config{BasePath: "/line", Scenario: "line_success"})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/line/v2/bot/message/push", strings.NewReader(`{"to":"Umockport","messages":[{"type":"text","text":"hello"}]}`))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("push status = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var push map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &push); err != nil {
		t.Fatalf("decode push: %v", err)
	}
	if _, ok := push["sentMessages"].([]any); !ok {
		t.Fatalf("push response missing sentMessages: %#v", push)
	}

	profileRec := httptest.NewRecorder()
	mux.ServeHTTP(profileRec, httptest.NewRequest(http.MethodGet, "/line/v2/bot/profile/Umockport", nil))
	if profileRec.Code != http.StatusOK {
		t.Fatalf("profile status = %d, want %d", profileRec.Code, http.StatusOK)
	}
	var profile map[string]any
	if err := json.Unmarshal(profileRec.Body.Bytes(), &profile); err != nil {
		t.Fatalf("decode profile: %v", err)
	}
	if profile["userId"] != "Umockport" || profile["displayName"] == "" {
		t.Fatalf("profile = %#v", profile)
	}
}

func TestMessagingAPICoreEndpoints(t *testing.T) {
	mux := newLineMux(t, adapter.Config{BasePath: "/line", Scenario: "line_success"})

	cases := []struct {
		method string
		path   string
		body   string
		status int
		want   string
	}{
		{http.MethodPost, "/line/v2/bot/message/validate/push", `{"messages":[{"type":"text","text":"hello"}]}`, http.StatusOK, `{}`},
		{http.MethodPost, "/line/v2/bot/message/multicast", `{"to":["Umockport"],"messages":[{"type":"text","text":"hello"}]}`, http.StatusOK, `{}`},
		{http.MethodPost, "/line/v2/bot/message/broadcast", `{"messages":[{"type":"text","text":"hello"}]}`, http.StatusOK, `{}`},
		{http.MethodPost, "/line/v2/bot/message/narrowcast", `{"messages":[{"type":"text","text":"hello"}]}`, http.StatusAccepted, `{}`},
		{http.MethodPost, "/line/v2/bot/chat/markAsRead", `{"markAsReadToken":"mockport_read"}`, http.StatusOK, `{}`},
		{http.MethodPost, "/line/v2/bot/chat/loading/start", `{"chatId":"Umockport","loadingSeconds":5}`, http.StatusAccepted, `{}`},
		{http.MethodGet, "/line/v2/bot/info", ``, http.StatusOK, `"basicId":"@mockport"`},
		{http.MethodGet, "/line/v2/bot/message/quota", ``, http.StatusOK, `"type":"limited"`},
		{http.MethodGet, "/line/v2/bot/message/delivery/reply?date=20260528", ``, http.StatusOK, `"status":"ready"`},
		{http.MethodGet, "/line/v2/bot/followers/ids?limit=2", ``, http.StatusOK, `"userIds"`},
	}

	for _, tc := range cases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			mux.ServeHTTP(rec, req)
			if rec.Code != tc.status {
				t.Fatalf("status = %d, want %d: %s", rec.Code, tc.status, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), tc.want) {
				t.Fatalf("body = %s, want to contain %s", rec.Body.String(), tc.want)
			}
			if rec.Header().Get("X-Line-Request-Id") == "" {
				t.Fatalf("missing X-Line-Request-Id")
			}
		})
	}
}

func TestWebhookEndpointSettingsAndContent(t *testing.T) {
	mux := newLineMux(t, adapter.Config{BasePath: "/line", Scenario: "line_success"})

	setRec := httptest.NewRecorder()
	setReq := httptest.NewRequest(http.MethodPut, "/line/v2/bot/channel/webhook/endpoint", strings.NewReader(`{"endpoint":"https://example.com/webhook"}`))
	setReq.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(setRec, setReq)
	if setRec.Code != http.StatusOK {
		t.Fatalf("set webhook status = %d, want %d: %s", setRec.Code, http.StatusOK, setRec.Body.String())
	}

	getRec := httptest.NewRecorder()
	mux.ServeHTTP(getRec, httptest.NewRequest(http.MethodGet, "/line/v2/bot/channel/webhook/endpoint", nil))
	if getRec.Code != http.StatusOK || !strings.Contains(getRec.Body.String(), `"endpoint":"https://example.com/webhook"`) {
		t.Fatalf("get webhook = status %d body %s", getRec.Code, getRec.Body.String())
	}

	contentRec := httptest.NewRecorder()
	mux.ServeHTTP(contentRec, httptest.NewRequest(http.MethodGet, "/line/v2/bot/message/mock-message/content", nil))
	if contentRec.Code != http.StatusOK || contentRec.Header().Get("Content-Type") == "" || contentRec.Body.Len() == 0 {
		t.Fatalf("content response = status %d content-type %q len %d", contentRec.Code, contentRec.Header().Get("Content-Type"), contentRec.Body.Len())
	}

	transcodingRec := httptest.NewRecorder()
	mux.ServeHTTP(transcodingRec, httptest.NewRequest(http.MethodGet, "/line/v2/bot/message/mock-message/content/transcoding", nil))
	if transcodingRec.Code != http.StatusOK || !strings.Contains(transcodingRec.Body.String(), `"status":"succeeded"`) {
		t.Fatalf("transcoding = status %d body %s", transcodingRec.Code, transcodingRec.Body.String())
	}
}

func TestChannelAccessTokenEndpoints(t *testing.T) {
	mux := newLineMux(t, adapter.Config{BasePath: "/line", Scenario: "line_success"})

	form := url.Values{"grant_type": {"client_credentials"}, "client_id": {"mockport_line_channel"}, "client_secret": {"mockport_line_secret"}}
	tokenRec := httptest.NewRecorder()
	tokenReq := httptest.NewRequest(http.MethodPost, "/line/oauth2/v3/token", strings.NewReader(form.Encode()))
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	mux.ServeHTTP(tokenRec, tokenReq)
	if tokenRec.Code != http.StatusOK {
		t.Fatalf("stateless token status = %d, want %d: %s", tokenRec.Code, http.StatusOK, tokenRec.Body.String())
	}
	var token map[string]any
	if err := json.Unmarshal(tokenRec.Body.Bytes(), &token); err != nil {
		t.Fatalf("decode token: %v", err)
	}
	accessToken, _ := token["access_token"].(string)
	if accessToken == "" || token["expires_in"].(float64) != 900 {
		t.Fatalf("token = %#v", token)
	}

	verifyRec := httptest.NewRecorder()
	mux.ServeHTTP(verifyRec, httptest.NewRequest(http.MethodGet, "/line/oauth2/v2.1/verify?access_token="+url.QueryEscape(accessToken), nil))
	if verifyRec.Code != http.StatusOK || !strings.Contains(verifyRec.Body.String(), `"client_id":"mockport_line_channel"`) {
		t.Fatalf("verify = status %d body %s", verifyRec.Code, verifyRec.Body.String())
	}

	revokeRec := httptest.NewRecorder()
	revokeReq := httptest.NewRequest(http.MethodPost, "/line/v2/oauth/revoke", strings.NewReader(url.Values{"access_token": {accessToken}}.Encode()))
	revokeReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	mux.ServeHTTP(revokeRec, revokeReq)
	if revokeRec.Code != http.StatusOK || strings.TrimSpace(revokeRec.Body.String()) != "" {
		t.Fatalf("revoke = status %d body %q", revokeRec.Code, revokeRec.Body.String())
	}
}

func TestRichMenuAndChatResources(t *testing.T) {
	mux := newLineMux(t, adapter.Config{BasePath: "/line", Scenario: "line_success"})

	createRec := httptest.NewRecorder()
	createReq := httptest.NewRequest(http.MethodPost, "/line/v2/bot/richmenu", strings.NewReader(`{"size":{"width":2500,"height":1686},"selected":false,"name":"Nice rich menu","chatBarText":"Tap","areas":[{"bounds":{"x":0,"y":0,"width":2500,"height":1686},"action":{"type":"postback","data":"a=b"}}]}`))
	createReq.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusOK {
		t.Fatalf("create rich menu status = %d, want %d: %s", createRec.Code, http.StatusOK, createRec.Body.String())
	}
	var created map[string]any
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode rich menu: %v", err)
	}
	richMenuID, _ := created["richMenuId"].(string)
	if richMenuID == "" {
		t.Fatalf("richMenuId missing: %#v", created)
	}

	uploadRec := httptest.NewRecorder()
	uploadReq := httptest.NewRequest(http.MethodPost, "/line/v2/bot/richmenu/"+richMenuID+"/content", strings.NewReader("image"))
	uploadReq.Header.Set("Content-Type", "image/png")
	mux.ServeHTTP(uploadRec, uploadReq)
	if uploadRec.Code != http.StatusOK {
		t.Fatalf("upload rich menu status = %d, want %d: %s", uploadRec.Code, http.StatusOK, uploadRec.Body.String())
	}

	defaultRec := httptest.NewRecorder()
	mux.ServeHTTP(defaultRec, httptest.NewRequest(http.MethodPost, "/line/v2/bot/user/all/richmenu/"+richMenuID, nil))
	if defaultRec.Code != http.StatusOK {
		t.Fatalf("set default rich menu status = %d, want %d: %s", defaultRec.Code, http.StatusOK, defaultRec.Body.String())
	}

	aliasRec := httptest.NewRecorder()
	aliasReq := httptest.NewRequest(http.MethodPost, "/line/v2/bot/richmenu/alias", strings.NewReader(`{"richMenuAliasId":"richmenu-alias-a","richMenuId":"`+richMenuID+`"}`))
	aliasReq.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(aliasRec, aliasReq)
	if aliasRec.Code != http.StatusOK {
		t.Fatalf("create rich menu alias status = %d, want %d: %s", aliasRec.Code, http.StatusOK, aliasRec.Body.String())
	}

	groupRec := httptest.NewRecorder()
	mux.ServeHTTP(groupRec, httptest.NewRequest(http.MethodGet, "/line/v2/bot/group/Cmockport/summary", nil))
	if groupRec.Code != http.StatusOK || !strings.Contains(groupRec.Body.String(), `"groupId":"Cmockport"`) {
		t.Fatalf("group summary = status %d body %s", groupRec.Code, groupRec.Body.String())
	}
}

func TestLoginTokenAndProfile(t *testing.T) {
	mux := newLineMux(t, adapter.Config{BasePath: "/line", Scenario: "line_success"})

	authorize := httptest.NewRecorder()
	mux.ServeHTTP(authorize, httptest.NewRequest(http.MethodGet, "/line/oauth2/v2.1/authorize?client_id=mockport_line_channel&redirect_uri=http://localhost/callback&state=abc&scope=profile%20openid", nil))
	if authorize.Code != http.StatusFound {
		t.Fatalf("authorize status = %d, want %d", authorize.Code, http.StatusFound)
	}
	location, err := url.Parse(authorize.Header().Get("Location"))
	if err != nil {
		t.Fatalf("parse redirect: %v", err)
	}
	code := location.Query().Get("code")
	if code == "" || location.Query().Get("state") != "abc" {
		t.Fatalf("redirect location = %s", authorize.Header().Get("Location"))
	}

	form := url.Values{"grant_type": {"authorization_code"}, "code": {code}, "redirect_uri": {"http://localhost/callback"}}
	tokenRec := httptest.NewRecorder()
	tokenReq := httptest.NewRequest(http.MethodPost, "/line/oauth2/v2.1/token", strings.NewReader(form.Encode()))
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	mux.ServeHTTP(tokenRec, tokenReq)
	if tokenRec.Code != http.StatusOK {
		t.Fatalf("token status = %d, want %d: %s", tokenRec.Code, http.StatusOK, tokenRec.Body.String())
	}
	var token map[string]any
	if err := json.Unmarshal(tokenRec.Body.Bytes(), &token); err != nil {
		t.Fatalf("decode token: %v", err)
	}
	accessToken, _ := token["access_token"].(string)
	if accessToken == "" || token["token_type"] != "Bearer" {
		t.Fatalf("token = %#v", token)
	}

	profileRec := httptest.NewRecorder()
	profileReq := httptest.NewRequest(http.MethodGet, "/line/v2/profile", nil)
	profileReq.Header.Set("Authorization", "Bearer "+accessToken)
	mux.ServeHTTP(profileRec, profileReq)
	if profileRec.Code != http.StatusOK {
		t.Fatalf("login profile status = %d, want %d", profileRec.Code, http.StatusOK)
	}
}

func TestLinePayRequestAndConfirm(t *testing.T) {
	mux := newLineMux(t, adapter.Config{BasePath: "/line", Scenario: "line_success"})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/line/v3/payments/request", strings.NewReader(`{"amount":1200,"currency":"JPY","orderId":"order-1","packages":[{"id":"pkg-1","amount":1200,"products":[{"name":"ticket","quantity":1,"price":1200}]}]}`))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("pay request status = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode pay request: %v", err)
	}
	if body["returnCode"] != "0000" {
		t.Fatalf("pay request = %#v", body)
	}
	info, _ := body["info"].(map[string]any)
	transactionID, _ := info["transactionId"].(string)
	if transactionID == "" {
		t.Fatalf("transaction id missing: %#v", body)
	}

	confirm := httptest.NewRecorder()
	confirmReq := httptest.NewRequest(http.MethodPost, "/line/v3/payments/"+transactionID+"/confirm", strings.NewReader(`{"amount":1200,"currency":"JPY"}`))
	confirmReq.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(confirm, confirmReq)
	if confirm.Code != http.StatusOK {
		t.Fatalf("confirm status = %d, want %d: %s", confirm.Code, http.StatusOK, confirm.Body.String())
	}
}

func TestMiniAppServiceMessageAndMiniDappPayment(t *testing.T) {
	mux := newLineMux(t, adapter.Config{BasePath: "/line", Scenario: "line_success"})

	tokenRec := httptest.NewRecorder()
	tokenReq := httptest.NewRequest(http.MethodPost, "/line/message/v3/notifier/token", strings.NewReader(`{"liffAccessToken":"mockport_liff_token"}`))
	tokenReq.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(tokenRec, tokenReq)
	if tokenRec.Code != http.StatusOK {
		t.Fatalf("notifier token status = %d, want %d: %s", tokenRec.Code, http.StatusOK, tokenRec.Body.String())
	}
	var notifier map[string]any
	if err := json.Unmarshal(tokenRec.Body.Bytes(), &notifier); err != nil {
		t.Fatalf("decode notifier token: %v", err)
	}
	notificationToken, _ := notifier["notificationToken"].(string)
	if notificationToken == "" || notifier["remainingCount"].(float64) == 0 {
		t.Fatalf("notifier = %#v", notifier)
	}

	sendRec := httptest.NewRecorder()
	sendReq := httptest.NewRequest(http.MethodPost, "/line/message/v3/notifier/send?target=service", strings.NewReader(`{"templateName":"reservation_complete_en","notificationToken":"`+notificationToken+`","params":{"username":"Brown"}}`))
	sendReq.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(sendRec, sendReq)
	if sendRec.Code != http.StatusOK {
		t.Fatalf("service message status = %d, want %d: %s", sendRec.Code, http.StatusOK, sendRec.Body.String())
	}

	paymentRec := httptest.NewRecorder()
	paymentReq := httptest.NewRequest(http.MethodPost, "/line/mini-dapp/v1/payments", strings.NewReader(`{"itemId":"item-1","amount":"10","currency":"KAIA"}`))
	paymentReq.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(paymentRec, paymentReq)
	if paymentRec.Code != http.StatusOK {
		t.Fatalf("mini dapp payment status = %d, want %d: %s", paymentRec.Code, http.StatusOK, paymentRec.Body.String())
	}
}

func TestLineErrorScenarios(t *testing.T) {
	auth := performLineRequest(t, adapter.Config{BasePath: "/line", Scenario: "auth_error"}, http.MethodPost, "/line/v2/bot/message/push", `{"to":"Umockport","messages":[{"type":"text","text":"hello"}]}`)
	if auth.Code != http.StatusUnauthorized {
		t.Fatalf("auth status = %d, want %d", auth.Code, http.StatusUnauthorized)
	}
	rate := performLineRequest(t, adapter.Config{BasePath: "/line", Scenario: "rate_limited"}, http.MethodPost, "/line/v2/bot/message/push", `{"to":"Umockport","messages":[{"type":"text","text":"hello"}]}`)
	if rate.Code != http.StatusTooManyRequests {
		t.Fatalf("rate status = %d, want %d", rate.Code, http.StatusTooManyRequests)
	}
	pay := performLineRequest(t, adapter.Config{BasePath: "/line", Scenario: "pay_failed"}, http.MethodPost, "/line/v3/payments/request", `{"amount":1200,"currency":"JPY","orderId":"order-1"}`)
	if pay.Code != http.StatusOK || !strings.Contains(pay.Body.String(), `"returnCode":"1169"`) {
		t.Fatalf("pay failure = status %d body %s", pay.Code, pay.Body.String())
	}
}

func performLineRequest(t *testing.T, cfg adapter.Config, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	mux := newLineMux(t, cfg)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rec, req)
	return rec
}

func newLineMux(t *testing.T, cfg adapter.Config) *http.ServeMux {
	t.Helper()
	mux := http.NewServeMux()
	if err := New().Register(mux, cfg); err != nil {
		t.Fatalf("register line adapter: %v", err)
	}
	return mux
}
