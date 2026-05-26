package stripe

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/albert-einshutoin/mockport/internal/adapter"
	"github.com/albert-einshutoin/mockport/internal/state"
)

type routes struct {
	basePath    string
	cfg         adapter.Config
	store       *state.Store
	idempotency *state.IdempotencyStore
}

func (rt *routes) handle(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, rt.basePath)
	switch {
	case r.Method == http.MethodPost && path == "/v1/checkout/sessions":
		rt.writeCheckoutSession(w, r)
	case r.Method == http.MethodGet && path == "/v1/checkout/sessions":
		rt.writeList(w, "checkout_session")
	case r.Method == http.MethodGet && strings.HasPrefix(path, "/v1/checkout/sessions/"):
		rt.writeResource(w, "checkout_session", strings.TrimPrefix(path, "/v1/checkout/sessions/"), fallbackCheckoutSession)
	case r.Method == http.MethodPost && path == "/v1/payment_intents":
		rt.writePaymentIntent(w, r)
	case r.Method == http.MethodGet && path == "/v1/payment_intents":
		rt.writeList(w, "payment_intent")
	case r.Method == http.MethodGet && strings.HasPrefix(path, "/v1/payment_intents/"):
		rt.writeResource(w, "payment_intent", strings.TrimPrefix(path, "/v1/payment_intents/"), fallbackPaymentIntent)
	case r.Method == http.MethodPost && path == "/test/webhook/send":
		rt.sendWebhook(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (rt *routes) writeCheckoutSession(w http.ResponseWriter, r *http.Request) {
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
		fields := formFields(r)
		body := map[string]interface{}{
			"object":         "checkout.session",
			"payment_status": "paid",
		}
		if clientReferenceID := fields.Get("client_reference_id"); clientReferenceID != "" {
			body["client_reference_id"] = clientReferenceID
		}
		rt.createStatefulResource(w, r, "checkout_session", body)
	}
}

func (rt *routes) writePaymentIntent(w http.ResponseWriter, r *http.Request) {
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
		fields := formFields(r)
		body := map[string]interface{}{
			"object": "payment_intent",
			"status": "succeeded",
		}
		if amount := fields.Get("amount"); amount != "" {
			if parsed, err := strconv.Atoi(amount); err == nil {
				body["amount"] = parsed
			}
		}
		if currency := fields.Get("currency"); currency != "" {
			body["currency"] = currency
		}
		rt.createStatefulResource(w, r, "payment_intent", body)
	}
}

func (rt *routes) createStatefulResource(w http.ResponseWriter, r *http.Request, resourceType string, body map[string]interface{}) {
	scope := "stripe:" + resourceType
	fingerprint := requestFingerprint(r)
	if replayed, replayResponse, err := rt.idempotency.Lookup(scope, r.Header.Get("Idempotency-Key"), fingerprint); err != nil {
		var conflict *state.IdempotencyConflictError
		if errors.As(err, &conflict) {
			rt.writeStripeError(w, http.StatusConflict, "idempotency_error", "idempotency_key_in_use", "Keys for idempotent requests can only be reused with the same parameters")
			return
		}
		rt.writeStripeError(w, http.StatusInternalServerError, "api_error", "mockport_idempotency_error", err.Error())
		return
	} else if replayed {
		rt.writeJSON(w, replayResponse.Status, replayResponse.Body)
		return
	}

	resource, err := rt.store.Create("stripe", resourceType, body)
	if err != nil {
		rt.writeStripeError(w, http.StatusInternalServerError, "api_error", "mockport_state_error", err.Error())
		return
	}
	response := resource.Data
	response["id"] = resource.ID
	idempotentResponse := state.IdempotentResponse{Status: http.StatusOK, Body: response}
	replayed, replayResponse, err := rt.idempotency.Remember(scope, r.Header.Get("Idempotency-Key"), fingerprint, idempotentResponse)
	if err != nil {
		var conflict *state.IdempotencyConflictError
		if errors.As(err, &conflict) {
			rt.writeStripeError(w, http.StatusConflict, "idempotency_error", "idempotency_key_in_use", "Keys for idempotent requests can only be reused with the same parameters")
			return
		}
		rt.writeStripeError(w, http.StatusInternalServerError, "api_error", "mockport_idempotency_error", err.Error())
		return
	}
	if replayed {
		rt.writeJSON(w, replayResponse.Status, replayResponse.Body)
		return
	}
	rt.writeJSON(w, http.StatusOK, response)
}

func (rt *routes) writeResource(w http.ResponseWriter, resourceType, id string, fallback func(string) map[string]interface{}) {
	if resource, ok := rt.store.Get("stripe", resourceType, id); ok {
		body := resource.Data
		body["id"] = resource.ID
		rt.writeJSON(w, http.StatusOK, body)
		return
	}
	rt.writeJSON(w, http.StatusOK, fallback(id))
}

func (rt *routes) writeList(w http.ResponseWriter, resourceType string) {
	var data []map[string]interface{}
	for _, resource := range rt.store.List("stripe", resourceType) {
		body := resource.Data
		body["id"] = resource.ID
		data = append(data, body)
	}
	rt.writeJSON(w, http.StatusOK, map[string]interface{}{
		"object": "list",
		"data":   data,
	})
}

func fallbackCheckoutSession(id string) map[string]interface{} {
	return map[string]interface{}{
		"id":             id,
		"object":         "checkout.session",
		"payment_status": "paid",
	}
}

func fallbackPaymentIntent(id string) map[string]interface{} {
	return map[string]interface{}{
		"id":     id,
		"object": "payment_intent",
		"status": "succeeded",
	}
}

func formFields(r *http.Request) interface{ Get(string) string } {
	if err := r.ParseForm(); err != nil {
		return emptyValues{}
	}
	return r.Form
}

type emptyValues struct{}

func (emptyValues) Get(string) string { return "" }

func requestFingerprint(r *http.Request) string {
	_ = r.ParseForm()
	return r.Method + " " + r.URL.Path + " " + r.Form.Encode()
}

func (rt *routes) writeStripeError(w http.ResponseWriter, status int, typ, code, message string) {
	rt.writeJSON(w, status, errorBody{Error: stripeError{Type: typ, Code: code, Message: message}})
}

func (rt *routes) writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
