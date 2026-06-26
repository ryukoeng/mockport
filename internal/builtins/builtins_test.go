package builtins

import (
	"slices"
	"testing"
)

func TestAdaptersReturnsUniqueNames(t *testing.T) {
	adapters := Adapters()
	names := make([]string, 0, len(adapters))
	for _, adapterImpl := range adapters {
		name := adapterImpl.Name()
		if slices.Contains(names, name) {
			t.Fatalf("duplicate adapter name: %q", name)
		}
		names = append(names, name)
	}
	if len(names) != 6 {
		t.Fatalf("len(Adapters()) = %d, want 6", len(names))
	}
}

func TestManifestAdaptersAreRequiredSubset(t *testing.T) {
	required := []string{"stripe", "openai", "github-oauth", "slack", "line"}
	manifestAdapters := ManifestAdapters()
	if len(manifestAdapters) != len(required) {
		t.Fatalf("len(ManifestAdapters()) = %d, want %d", len(manifestAdapters), len(required))
	}
	for i, adapterImpl := range manifestAdapters {
		if adapterImpl.Name() != required[i] {
			t.Fatalf("manifest adapter[%d] = %q, want %q", i, adapterImpl.Name(), required[i])
		}
	}
	allNames := make([]string, 0, len(Adapters()))
	for _, adapterImpl := range Adapters() {
		allNames = append(allNames, adapterImpl.Name())
	}
	for _, name := range required {
		if !slices.Contains(allNames, name) {
			t.Fatalf("required manifest adapter %q missing from Adapters()", name)
		}
	}
}
