package compat

import "testing"

func TestCalculateScoreUsesEndpointScenarioSDKStateAndErrorCoverage(t *testing.T) {
	score := CalculateScore(Manifest{
		Adapter:         "stripe",
		ProviderVersion: "2026-05-26",
		Levels:          []Level{LevelWire, LevelSDK, LevelState, LevelError},
		SDKVersions:     []SDKVersion{{Name: "stripe-go", Version: "v83.0.0"}},
		Endpoints: []Endpoint{
			{ID: "supported", Supported: true},
			{ID: "unsupported", Supported: false},
		},
		Scenarios: []Scenario{
			{Name: "built_in_supported", BuiltIn: true, Supported: true},
			{Name: "built_in_unsupported", BuiltIn: true, Supported: false},
			{Name: "local_custom", BuiltIn: false, Supported: true},
		},
	})

	if score.EndpointCoverage != 50 {
		t.Fatalf("endpoint coverage = %d, want 50", score.EndpointCoverage)
	}
	if score.ScenarioCoverage != 50 {
		t.Fatalf("scenario coverage = %d, want 50", score.ScenarioCoverage)
	}
	if score.SDKCoverage != 100 || score.StateCoverage != 100 || score.ErrorCoverage != 100 {
		t.Fatalf("score = %#v, want sdk/state/error coverage", score)
	}
	if score.Total != 80 {
		t.Fatalf("total = %d, want 80", score.Total)
	}
	if score.Level != string(LevelState) {
		t.Fatalf("level = %q, want state", score.Level)
	}
}

func TestCanPromoteRequiresContractEvidenceForProviderCompatible(t *testing.T) {
	manifest := Manifest{
		Adapter:         "stripe",
		ProviderVersion: "2026-05-26",
		Levels:          []Level{LevelWire, LevelSDK, LevelWorkflow, LevelState, LevelError},
		SDKVersions:     []SDKVersion{{Name: "stripe-go", Version: "v83.0.0"}},
		Endpoints:       []Endpoint{{ID: "one", Supported: true}},
		Scenarios:       []Scenario{{Name: "payment_success", BuiltIn: true, Supported: true}},
	}
	score := CalculateScore(manifest)

	if CanPromote(manifest, score, "provider-compatible") {
		t.Fatal("provider-compatible promotion passed without contract level")
	}
	manifest.Levels = append(manifest.Levels, LevelContract)
	score = CalculateScore(manifest)
	if !CanPromote(manifest, score, "provider-compatible") {
		t.Fatalf("provider-compatible promotion failed with contract evidence: %#v", score)
	}
}
