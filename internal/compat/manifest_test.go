package compat

import (
	"net/http"
	"strings"
	"testing"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

func TestManifestValidation(t *testing.T) {
	manifest := Manifest{
		Adapter:         "stripe",
		ProviderVersion: "2026-05-26",
		SDKVersions:     []SDKVersion{{Name: "stripe-go", Version: "v83.0.0"}},
		Levels:          []Level{LevelWire, LevelSDK, LevelError},
		Endpoints: []Endpoint{
			{ID: "checkout_sessions_create", Method: http.MethodPost, Path: "/v1/checkout/sessions", Supported: true, Levels: []Level{LevelWire}},
		},
		Scenarios: []Scenario{
			{Name: "payment_success", BuiltIn: true, Supported: true, Levels: []Level{LevelWire}},
		},
	}

	if err := manifest.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestManifestValidationRejectsDuplicateEndpointIDs(t *testing.T) {
	manifest := Manifest{
		Adapter:         "stripe",
		ProviderVersion: "2026-05-26",
		Endpoints: []Endpoint{
			{ID: "duplicate", Method: http.MethodGet, Path: "/v1/a"},
			{ID: "duplicate", Method: http.MethodGet, Path: "/v1/b"},
		},
	}

	err := manifest.Validate()
	if err == nil || !strings.Contains(err.Error(), "duplicate endpoint id") {
		t.Fatalf("Validate() error = %v, want duplicate endpoint id", err)
	}
}

func TestManifestValidationRejectsInvalidLevels(t *testing.T) {
	manifest := Manifest{
		Adapter:         "stripe",
		ProviderVersion: "2026-05-26",
		Levels:          []Level{"complete"},
	}

	err := manifest.Validate()
	if err == nil || !strings.Contains(err.Error(), "invalid compatibility level") {
		t.Fatalf("Validate() error = %v, want invalid compatibility level", err)
	}
}

func TestManifestFromAdapterMetadata(t *testing.T) {
	meta := adapter.Metadata{
		Name:         "openai",
		Maturity:     "experimental",
		Capabilities: []string{"chat_completions"},
		Scenarios:    []adapter.Scenario{{Name: "chat_success", Supported: true}},
		Endpoints:    []adapter.Endpoint{{Method: http.MethodPost, Path: "/openai/v1/chat/completions", SupportedScenarios: []string{"chat_success"}}},
	}

	manifest := FromMetadata(meta)
	if manifest.Adapter != "openai" {
		t.Fatalf("adapter = %q, want openai", manifest.Adapter)
	}
	if manifest.ProviderVersion == "" {
		t.Fatal("provider version is empty")
	}
	if len(manifest.Endpoints) != 1 || manifest.Endpoints[0].ID == "" {
		t.Fatalf("endpoints = %#v, want generated id", manifest.Endpoints)
	}
	if len(manifest.Scenarios) != 1 || !manifest.Scenarios[0].BuiltIn {
		t.Fatalf("scenarios = %#v, want built-in scenario", manifest.Scenarios)
	}
}
