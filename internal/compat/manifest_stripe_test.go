package compat

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/albert-einshutoin/mockport/adapters/stripe"
)

func TestStripeCheckedInManifestMatchesFromMetadataFields(t *testing.T) {
	repoRoot, err := FindRepoRoot(".")
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	checkedIn, err := os.ReadFile(filepath.Join(repoRoot, "compat", "manifests", "stripe.json"))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	var got Manifest
	if err := json.Unmarshal(checkedIn, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	want := FromMetadata(stripe.New().Metadata())

	if want.Adapter != got.Adapter {
		t.Fatalf("adapter: want %q got %q", want.Adapter, got.Adapter)
	}
	if want.ProviderVersion != got.ProviderVersion {
		t.Fatalf("provider_version: want %q got %q", want.ProviderVersion, got.ProviderVersion)
	}
	if want.Maturity != got.Maturity {
		t.Fatalf("maturity: want %q got %q", want.Maturity, got.Maturity)
	}
	if !reflect.DeepEqual(want.SDKVersions, got.SDKVersions) {
		t.Fatalf("sdk_versions: want %#v got %#v", want.SDKVersions, got.SDKVersions)
	}
	if !reflect.DeepEqual(want.Levels, got.Levels) {
		t.Fatalf("levels: want %#v got %#v", want.Levels, got.Levels)
	}
	if !reflect.DeepEqual(want.Scenarios, got.Scenarios) {
		t.Fatalf("scenarios: want %#v got %#v", want.Scenarios, got.Scenarios)
	}
	if !reflect.DeepEqual(want.StateEvidence, got.StateEvidence) {
		t.Fatalf("state_evidence: want %#v got %#v", want.StateEvidence, got.StateEvidence)
	}
	if len(want.Endpoints) != len(got.Endpoints) {
		t.Fatalf("endpoint count: want %d got %d", len(want.Endpoints), len(got.Endpoints))
	}
	for i := range want.Endpoints {
		if !reflect.DeepEqual(want.Endpoints[i], got.Endpoints[i]) {
			t.Fatalf("endpoint[%d]: want %#v got %#v", i, want.Endpoints[i], got.Endpoints[i])
		}
	}
}
