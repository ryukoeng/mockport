package adapter

import (
	"encoding/json"
	"testing"
)

func TestContractEvidenceJSONKeys(t *testing.T) {
	evidence := ContractEvidence{
		Fixtures:     []string{"compat/fixtures/stripe/checkout_session_create.json"},
		SDKContracts: []string{"contract/sdk/stripe"},
		KnownGaps:    []string{"docs/compatibility-reports/latest.json#stripe"},
	}
	raw, err := json.Marshal(evidence)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"fixtures", "sdk_contracts", "known_gaps"} {
		if _, ok := decoded[key]; !ok {
			t.Fatalf("missing json key %q in %s", key, raw)
		}
	}
}

func TestContractEvidenceCloneDeepCopy(t *testing.T) {
	original := ContractEvidence{
		Fixtures:     []string{"a"},
		SDKContracts: []string{"b"},
		KnownGaps:    []string{"c"},
	}
	clone := original.Clone()
	original.Fixtures[0] = "mutated"
	original.SDKContracts[0] = "mutated"
	original.KnownGaps[0] = "mutated"
	if clone.Fixtures[0] != "a" || clone.SDKContracts[0] != "b" || clone.KnownGaps[0] != "c" {
		t.Fatalf("clone shares slices: %#v", clone)
	}
}
