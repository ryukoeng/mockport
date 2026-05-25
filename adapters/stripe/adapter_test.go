package stripe

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

func TestCheckoutSessionSuccess(t *testing.T) {
	rec := performStripeRequest(t, adapter.Config{BasePath: "/stripe", Scenario: "payment_success"}, http.MethodPost, "/stripe/v1/checkout/sessions")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["object"] != "checkout.session" || body["payment_status"] != "paid" {
		t.Fatalf("unexpected body: %#v", body)
	}
}

func TestCheckoutSessionPaymentFailed(t *testing.T) {
	rec := performStripeRequest(t, adapter.Config{BasePath: "/stripe", Scenario: "payment_failed"}, http.MethodPost, "/stripe/v1/checkout/sessions")
	if rec.Code != http.StatusPaymentRequired {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusPaymentRequired)
	}
	assertStripeErrorCode(t, rec, "card_declined")
}

func TestPaymentIntentScenarios(t *testing.T) {
	tests := []struct {
		name     string
		scenario string
		wantCode int
		wantErr  string
	}{
		{"success", "payment_success", http.StatusOK, ""},
		{"failed", "payment_failed", http.StatusPaymentRequired, "card_declined"},
		{"auth", "auth_error", http.StatusUnauthorized, "invalid_api_key"},
		{"rate limited", "rate_limited", http.StatusTooManyRequests, "rate_limited"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := performStripeRequest(t, adapter.Config{BasePath: "/stripe", Scenario: tt.scenario}, http.MethodPost, "/stripe/v1/payment_intents")
			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantCode)
			}
			if tt.wantErr != "" {
				assertStripeErrorCode(t, rec, tt.wantErr)
			}
		})
	}
}

func TestGetPaymentIntent(t *testing.T) {
	rec := performStripeRequest(t, adapter.Config{BasePath: "/stripe", Scenario: "payment_success"}, http.MethodGet, "/stripe/v1/payment_intents/pi_mockport")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestTimeoutScenario(t *testing.T) {
	rec := performStripeRequest(t, adapter.Config{BasePath: "/stripe", Scenario: "timeout"}, http.MethodPost, "/stripe/v1/checkout/sessions")
	if rec.Code != http.StatusGatewayTimeout {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusGatewayTimeout)
	}
	assertStripeErrorCode(t, rec, "mockport_timeout")
}

func TestWebhookSender(t *testing.T) {
	received := make(chan *http.Request, 1)
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Stripe-Signature") == "" {
			t.Error("missing Stripe-Signature header")
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode webhook body: %v", err)
		}
		if body["type"] != "checkout.session.completed" {
			t.Errorf("event type = %v, want checkout.session.completed", body["type"])
		}
		received <- r
		w.WriteHeader(http.StatusNoContent)
	}))
	defer target.Close()

	rec := performStripeRequest(t, adapter.Config{
		BasePath:             "/stripe",
		Scenario:             "payment_success",
		WebhookTargetURL:     target.URL,
		WebhookSigningSecret: "whsec_mockport",
	}, http.MethodPost, "/stripe/test/webhook/send")
	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusAccepted, rec.Body.String())
	}
	select {
	case <-received:
	default:
		t.Fatal("webhook target did not receive request")
	}
}

func performStripeRequest(t *testing.T, cfg adapter.Config, method, path string) *httptest.ResponseRecorder {
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

func assertStripeErrorCode(t *testing.T, rec *httptest.ResponseRecorder, want string) {
	t.Helper()
	var body errorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode error body: %v", err)
	}
	if body.Error.Code != want {
		t.Fatalf("error code = %q, want %q", body.Error.Code, want)
	}
}
