package stripe

import (
	"net/http"
	"strings"

	"github.com/albert-einshutoin/mockport/internal/adapter"
	"github.com/albert-einshutoin/mockport/internal/state"
)

type Adapter struct{}

func New() Adapter {
	return Adapter{}
}

func (a Adapter) Name() string {
	return "stripe"
}

func (a Adapter) Register(mux *http.ServeMux, cfg adapter.Config) error {
	basePath := cfg.BasePath
	if basePath == "" {
		basePath = "/stripe"
	}
	routes := &routes{
		basePath:    strings.TrimRight(basePath, "/"),
		cfg:         cfg,
		store:       state.NewStore(),
		idempotency: state.NewIdempotencyStore(),
	}
	mux.HandleFunc(routes.basePath+"/", routes.handle)
	return nil
}

func (a Adapter) FakeEnv(cfg adapter.Config) map[string]string {
	basePath := cfg.BasePath
	if basePath == "" {
		basePath = "/stripe"
	}
	secret := cfg.FakeSecret
	if secret == "" {
		secret = "mockport_stripe_secret"
	}
	signingSecret := cfg.WebhookSigningSecret
	if signingSecret == "" {
		signingSecret = "whsec_mockport"
	}
	return map[string]string{
		"STRIPE_API_URL":        "http://localhost:43101" + basePath,
		"STRIPE_SECRET_KEY":     secret,
		"STRIPE_WEBHOOK_SECRET": signingSecret,
	}
}

func (a Adapter) Metadata() adapter.Metadata {
	scenarios := []adapter.Scenario{
		{Name: "payment_success", Supported: true},
		{Name: "payment_failed", Supported: true},
		{Name: "auth_error", Supported: true},
		{Name: "rate_limited", Supported: true},
		{Name: "timeout", Supported: true},
	}
	scenarioNames := []string{"payment_success", "payment_failed", "auth_error", "rate_limited", "timeout"}
	return adapter.Metadata{
		Name:         "stripe",
		Maturity:     "partial",
		Capabilities: []string{"checkout_sessions", "payment_intents", "webhooks"},
		Scenarios:    scenarios,
		StatefulResources: []string{
			"checkout_session",
			"payment_intent",
		},
		Idempotency: true,
		Reset:       true,
		Endpoints: []adapter.Endpoint{
			{Method: http.MethodPost, Path: "/stripe/v1/checkout/sessions", SupportedScenarios: scenarioNames, Notes: "Stripe-like checkout session creation"},
			{Method: http.MethodGet, Path: "/stripe/v1/checkout/sessions", SupportedScenarios: []string{"payment_success"}, Notes: "Deterministic checkout session list"},
			{Method: http.MethodGet, Path: "/stripe/v1/checkout/sessions/{id}", SupportedScenarios: []string{"payment_success"}, Notes: "Deterministic checkout session lookup"},
			{Method: http.MethodPost, Path: "/stripe/v1/payment_intents", SupportedScenarios: scenarioNames, Notes: "Stripe-like payment intent creation"},
			{Method: http.MethodGet, Path: "/stripe/v1/payment_intents", SupportedScenarios: []string{"payment_success"}, Notes: "Deterministic payment intent list"},
			{Method: http.MethodGet, Path: "/stripe/v1/payment_intents/{id}", SupportedScenarios: []string{"payment_success"}, Notes: "Deterministic payment intent lookup"},
			{Method: http.MethodPost, Path: "/stripe/test/webhook/send", SupportedScenarios: []string{"payment_success", "payment_failed"}, Notes: "Sends fake signed webhook to configured target"},
		},
	}
}
