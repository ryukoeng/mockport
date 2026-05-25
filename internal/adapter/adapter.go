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
}
