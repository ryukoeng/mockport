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
	rt.handlePath(w, r, path)
}

func (rt *routes) handleRoot(w http.ResponseWriter, r *http.Request) {
	rt.handlePath(w, r, r.URL.Path)
}

func (rt *routes) handlePath(w http.ResponseWriter, r *http.Request, path string) {
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
	case r.Method == http.MethodPost && path == "/v1/customers":
		rt.writeGenericResource(w, r, "customer", map[string]interface{}{"object": "customer"}, nil)
	case r.Method == http.MethodGet && path == "/v1/customers":
		rt.writeList(w, "customer")
	case r.Method == http.MethodGet && strings.HasPrefix(path, "/v1/customers/"):
		rt.writeResource(w, "customer", strings.TrimPrefix(path, "/v1/customers/"), nil)
	case r.Method == http.MethodPost && path == "/v1/products":
		rt.writeGenericResource(w, r, "product", map[string]interface{}{"object": "product", "active": true}, []string{"name"})
	case r.Method == http.MethodGet && path == "/v1/products":
		rt.writeList(w, "product")
	case r.Method == http.MethodGet && strings.HasPrefix(path, "/v1/products/"):
		rt.writeResource(w, "product", strings.TrimPrefix(path, "/v1/products/"), nil)
	case r.Method == http.MethodPost && path == "/v1/prices":
		rt.writeGenericResource(w, r, "price", map[string]interface{}{"object": "price", "active": true}, []string{"product", "currency", "unit_amount"})
	case r.Method == http.MethodGet && path == "/v1/prices":
		rt.writeList(w, "price")
	case r.Method == http.MethodGet && strings.HasPrefix(path, "/v1/prices/"):
		rt.writeResource(w, "price", strings.TrimPrefix(path, "/v1/prices/"), nil)
	case r.Method == http.MethodPost && path == "/v1/subscriptions":
		rt.writeGenericResource(w, r, "subscription", map[string]interface{}{"object": "subscription", "status": "active"}, []string{"customer"})
	case r.Method == http.MethodGet && path == "/v1/subscriptions":
		rt.writeList(w, "subscription")
	case r.Method == http.MethodGet && strings.HasPrefix(path, "/v1/subscriptions/"):
		rt.writeResource(w, "subscription", strings.TrimPrefix(path, "/v1/subscriptions/"), nil)
	case r.Method == http.MethodPost && path == "/v1/invoices":
		rt.writeGenericResource(w, r, "invoice", map[string]interface{}{"object": "invoice", "status": "draft"}, []string{"customer"})
	case r.Method == http.MethodGet && path == "/v1/invoices":
		rt.writeList(w, "invoice")
	case r.Method == http.MethodGet && strings.HasPrefix(path, "/v1/invoices/"):
		rt.writeResource(w, "invoice", strings.TrimPrefix(path, "/v1/invoices/"), nil)
	case r.Method == http.MethodPost && path == "/v1/refunds":
		rt.writeGenericResource(w, r, "refund", map[string]interface{}{"object": "refund", "status": "succeeded"}, []string{"payment_intent"})
	case r.Method == http.MethodGet && path == "/v1/refunds":
		rt.writeList(w, "refund")
	case r.Method == http.MethodGet && strings.HasPrefix(path, "/v1/refunds/"):
		rt.writeResource(w, "refund", strings.TrimPrefix(path, "/v1/refunds/"), nil)
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
		if !fields.Empty() {
			if missing := missingFields(fields, "amount", "currency"); len(missing) > 0 {
				rt.writeValidationError(w, missing[0])
				return
			}
		}
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

func (rt *routes) writeGenericResource(w http.ResponseWriter, r *http.Request, resourceType string, body map[string]interface{}, required []string) {
	fields := formFields(r)
	if !fields.Empty() {
		if missing := missingFields(fields, required...); len(missing) > 0 {
			rt.writeValidationError(w, missing[0])
			return
		}
	}
	for _, field := range fields.Keys() {
		body[field] = parseFormValue(fields.Get(field))
	}
	rt.createStatefulResource(w, r, resourceType, body)
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
	if fallback != nil && looksLikeLegacyID(resourceType, id) {
		rt.writeJSON(w, http.StatusOK, fallback(id))
		return
	}
	rt.writeStripeError(w, http.StatusNotFound, "invalid_request_error", "resource_missing", "No such "+resourceType+": "+id)
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

func formFields(r *http.Request) formValues {
	if err := r.ParseForm(); err != nil {
		return emptyValues{}
	}
	return urlValues{values: r.Form}
}

type formValues interface {
	Get(string) string
	Keys() []string
	Empty() bool
}

type urlValues struct {
	values map[string][]string
}

func (v urlValues) Get(name string) string {
	values := v.values[name]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (v urlValues) Keys() []string {
	keys := make([]string, 0, len(v.values))
	for key := range v.values {
		keys = append(keys, key)
	}
	return keys
}

func (v urlValues) Empty() bool { return len(v.values) == 0 }

type emptyValues struct{}

func (emptyValues) Get(string) string { return "" }
func (emptyValues) Keys() []string    { return nil }
func (emptyValues) Empty() bool       { return true }

func requestFingerprint(r *http.Request) string {
	_ = r.ParseForm()
	return r.Method + " " + r.URL.Path + " " + r.Form.Encode()
}

func (rt *routes) writeStripeError(w http.ResponseWriter, status int, typ, code, message string) {
	rt.writeJSON(w, status, errorBody{Error: stripeError{Type: typ, Code: code, Message: message}})
}

func (rt *routes) writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Request-Id", "req_mockport")
	w.Header().Set("Stripe-Version", "2025-10-29.clover")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func (rt *routes) writeValidationError(w http.ResponseWriter, field string) {
	rt.writeJSON(w, http.StatusBadRequest, errorBody{Error: stripeError{
		Type:    "invalid_request_error",
		Code:    "parameter_missing",
		Param:   field,
		Message: "Missing required param: " + field,
	}})
}

func missingFields(fields formValues, required ...string) []string {
	var missing []string
	for _, field := range required {
		if fields.Get(field) == "" {
			missing = append(missing, field)
		}
	}
	return missing
}

func parseFormValue(value string) interface{} {
	if parsed, err := strconv.Atoi(value); err == nil {
		return parsed
	}
	return value
}

func looksLikeLegacyID(resourceType, id string) bool {
	switch resourceType {
	case "checkout_session":
		return strings.HasPrefix(id, "cs_")
	case "payment_intent":
		return strings.HasPrefix(id, "pi_")
	default:
		return false
	}
}
