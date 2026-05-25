package stripe

import (
	"net/http"
	"strings"

	"github.com/albert-einshutoin/mockport/internal/adapter"
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
	routes := &routes{basePath: strings.TrimRight(basePath, "/"), cfg: cfg}
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
