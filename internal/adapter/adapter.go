package adapter

import "net/http"

type Config struct {
	BasePath             string
	Scenario             string
	FakeSecret           string
	WebhookTargetURL     string
	WebhookSigningSecret string
}

type Adapter interface {
	Name() string
	Register(mux *http.ServeMux, cfg Config) error
	FakeEnv(cfg Config) map[string]string
	Metadata() Metadata
}

type Metadata struct {
	Name              string
	Maturity          string
	Capabilities      []string
	Scenarios         []Scenario
	Endpoints         []Endpoint
	StatefulResources []string
	Idempotency       bool
	Reset             bool
}

type Scenario struct {
	Name      string
	Supported bool
}

type Endpoint struct {
	Method             string
	Path               string
	SupportedScenarios []string
	Notes              string
}

func ValidateMaturity(maturity string) bool {
	switch maturity {
	case "experimental", "partial", "sdk-compatible", "workflow-compatible", "provider-compatible":
		return true
	default:
		return false
	}
}
