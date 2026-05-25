package stripe

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func (rt *routes) sendWebhook(w http.ResponseWriter, r *http.Request) {
	if rt.cfg.WebhookTargetURL == "" {
		rt.writeStripeError(w, http.StatusBadRequest, "invalid_request_error", "missing_webhook_target", "webhook target URL is not configured")
		return
	}
	secret := rt.cfg.WebhookSigningSecret
	if secret == "" {
		secret = "whsec_mockport"
	}
	eventType := "checkout.session.completed"
	if normalizeScenario(rt.cfg.Scenario) == scenarioPaymentFailed {
		eventType = "payment_intent.payment_failed"
	}
	payload, err := json.Marshal(map[string]interface{}{
		"id":   "evt_mockport",
		"type": eventType,
		"data": map[string]interface{}{
			"object": map[string]interface{}{
				"id":     "cs_test_mockport",
				"object": "checkout.session",
			},
		},
	})
	if err != nil {
		rt.writeStripeError(w, http.StatusInternalServerError, "api_error", "webhook_encode_failed", "failed to encode webhook payload")
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, rt.cfg.WebhookTargetURL, bytes.NewReader(payload))
	if err != nil {
		rt.writeStripeError(w, http.StatusBadRequest, "invalid_request_error", "invalid_webhook_target", "webhook target URL is invalid")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Stripe-Signature", signPayload(secret, nowUnix(), payload))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		rt.writeStripeError(w, http.StatusBadGateway, "api_error", "webhook_send_failed", "failed to send webhook")
		return
	}
	defer resp.Body.Close()

	rt.writeJSON(w, http.StatusAccepted, map[string]interface{}{
		"sent":        true,
		"target_url":  rt.cfg.WebhookTargetURL,
		"event_type":  eventType,
		"status_code": resp.StatusCode,
	})
}
