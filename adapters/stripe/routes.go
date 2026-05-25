package stripe

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

type routes struct {
	basePath string
	cfg      adapter.Config
}

func (rt *routes) handle(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, rt.basePath)
	switch {
	case r.Method == http.MethodPost && path == "/v1/checkout/sessions":
		rt.writeCheckoutSession(w)
	case r.Method == http.MethodGet && strings.HasPrefix(path, "/v1/checkout/sessions/"):
		rt.writeJSON(w, http.StatusOK, map[string]interface{}{
			"id":             strings.TrimPrefix(path, "/v1/checkout/sessions/"),
			"object":         "checkout.session",
			"payment_status": "paid",
		})
	case r.Method == http.MethodPost && path == "/v1/payment_intents":
		rt.writePaymentIntent(w)
	case r.Method == http.MethodGet && strings.HasPrefix(path, "/v1/payment_intents/"):
		rt.writeJSON(w, http.StatusOK, map[string]interface{}{
			"id":     strings.TrimPrefix(path, "/v1/payment_intents/"),
			"object": "payment_intent",
			"status": "succeeded",
		})
	case r.Method == http.MethodPost && path == "/test/webhook/send":
		rt.sendWebhook(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (rt *routes) writeCheckoutSession(w http.ResponseWriter) {
	switch normalizeScenario(rt.cfg.Scenario) {
	case scenarioPaymentFailed:
		rt.writeStripeError(w, http.StatusPaymentRequired, "card_error", "card_declined", "Mockport simulated card decline")
	case scenarioAuthError:
		rt.writeStripeError(w, http.StatusUnauthorized, "invalid_request_error", "invalid_api_key", "Mockport simulated auth error")
	case scenarioRateLimited:
		rt.writeStripeError(w, http.StatusTooManyRequests, "rate_limit_error", "rate_limited", "Mockport simulated rate limit")
	case scenarioTimeout:
		rt.writeStripeError(w, http.StatusGatewayTimeout, "api_error", "mockport_timeout", "Mockport simulated timeout")
	default:
		rt.writeJSON(w, http.StatusOK, map[string]interface{}{
			"id":             "cs_test_mockport",
			"object":         "checkout.session",
			"payment_status": "paid",
		})
	}
}

func (rt *routes) writePaymentIntent(w http.ResponseWriter) {
	switch normalizeScenario(rt.cfg.Scenario) {
	case scenarioPaymentFailed:
		rt.writeStripeError(w, http.StatusPaymentRequired, "card_error", "card_declined", "Mockport simulated card decline")
	case scenarioAuthError:
		rt.writeStripeError(w, http.StatusUnauthorized, "invalid_request_error", "invalid_api_key", "Mockport simulated auth error")
	case scenarioRateLimited:
		rt.writeStripeError(w, http.StatusTooManyRequests, "rate_limit_error", "rate_limited", "Mockport simulated rate limit")
	case scenarioTimeout:
		rt.writeStripeError(w, http.StatusGatewayTimeout, "api_error", "mockport_timeout", "Mockport simulated timeout")
	default:
		rt.writeJSON(w, http.StatusOK, map[string]interface{}{
			"id":     "pi_mockport_success",
			"object": "payment_intent",
			"status": "succeeded",
		})
	}
}

func (rt *routes) writeStripeError(w http.ResponseWriter, status int, typ, code, message string) {
	rt.writeJSON(w, status, errorBody{Error: stripeError{Type: typ, Code: code, Message: message}})
}

func (rt *routes) writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
