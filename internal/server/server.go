package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/albert-einshutoin/mockport/internal/adapter"
	"github.com/albert-einshutoin/mockport/internal/config"
	"github.com/albert-einshutoin/mockport/internal/report"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	return mux
}

func NewConfiguredHandler(cfg config.Config, reg *adapter.Registry, rec *report.Recorder) (http.Handler, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	if rec == nil {
		rec = report.NewRecorder()
	}
	rec.SetMode(cfg.Mode)

	var adapterStatuses []report.AdapterStatus
	for warningIdx := range cfg.SafetyWarnings {
		warning := cfg.SafetyWarnings[warningIdx]
		rec.RecordSafetyWarning(warning.Field, warning.Category, warning.Message)
	}

	for name, adapterCfg := range cfg.Adapters {
		if !adapterCfg.Enabled {
			continue
		}
		registered, ok := reg.Get(name)
		if !ok {
			return nil, fmt.Errorf("adapter %s is enabled but not registered", name)
		}
		adapterStatuses = append(adapterStatuses, report.AdapterStatus{Name: name, BasePath: adapterCfg.BasePath, Enabled: true})
		if err := registered.Register(mux, adapter.Config{
			BasePath:             adapterCfg.BasePath,
			Scenario:             adapterCfg.Scenario,
			FakeSecret:           adapterCfg.FakeSecret,
			WebhookTargetURL:     adapterCfg.Webhook.TargetURL,
			WebhookSigningSecret: adapterCfg.Webhook.SigningSecret,
		}); err != nil {
			return nil, err
		}
	}
	rec.SetAdapters(adapterStatuses)
	mux.HandleFunc("/_mockport/report", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rec.Snapshot())
	})

	return recordMiddleware(mux, rec), nil
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func recordMiddleware(next http.Handler, rec *report.Recorder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sr, r)
		if r.URL.Path != "/_mockport/report" {
			rec.RecordRequest(r.Method, r.URL.Path, sr.status)
		}
	})
}
