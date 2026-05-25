package config

import (
	"fmt"
	"strings"

	"github.com/albert-einshutoin/mockport/internal/security"
)

func Validate(cfg *Config) error {
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port %d", cfg.Server.Port)
	}

	switch cfg.Mode {
	case "ai-safe", "local", "strict":
	default:
		return fmt.Errorf("unsupported mode %q", cfg.Mode)
	}

	var warnings []SafetyWarning
	for name, adapter := range cfg.Adapters {
		if adapter.Enabled && adapter.BasePath == "" {
			return fmt.Errorf("adapter %s missing base_path", name)
		}
		checks := map[string]string{
			name + ".fake_secret":            adapter.FakeSecret,
			name + ".api_url":                adapter.APIURL,
			name + ".webhook.target_url":     adapter.Webhook.TargetURL,
			name + ".webhook.signing_secret": adapter.Webhook.SigningSecret,
		}
		for field, value := range checks {
			if value == "" {
				continue
			}
			if security.LooksLikeSecret(value) {
				warnings = append(warnings, SafetyWarning{Field: field, Category: "real_looking_secret", Message: "real-looking secret detected"})
			}
			if security.LooksLikeExternalServiceURL(value) {
				warnings = append(warnings, SafetyWarning{Field: field, Category: "external_url", Message: "external live service URL detected"})
			}
		}
	}

	if cfg.Mode == "strict" && len(warnings) > 0 {
		fields := make([]string, 0, len(warnings))
		for _, warning := range warnings {
			fields = append(fields, warning.Field)
		}
		return fmt.Errorf("strict mode rejected unsafe config fields: %s", strings.Join(fields, ", "))
	}
	cfg.SafetyWarnings = warnings
	return nil
}
