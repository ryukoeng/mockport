package compat

import (
	"net/http"
	"testing"
)

func BenchmarkCalculateScore(b *testing.B) {
	manifest := Manifest{
		Adapter:         "stripe",
		ProviderVersion: "2025-10-29.clover",
		Levels:          []Level{LevelWire, LevelSDK, LevelWorkflow, LevelState, LevelError},
		SDKVersions:     []SDKVersion{{Name: "stripe", Version: "22.1.1"}},
		Endpoints:       []Endpoint{{ID: "create", Method: http.MethodPost, Path: "/v1/payment_intents", Supported: true}},
		Scenarios:       []Scenario{{Name: "payment_success", BuiltIn: true, Supported: true}},
	}
	for i := 0; i < b.N; i++ {
		_ = CalculateScore(manifest)
	}
}
