package adapter

import (
	"net/http"
	"testing"
)

type fakeAdapter struct{ name string }

func (a fakeAdapter) Name() string                                  { return a.name }
func (a fakeAdapter) Register(mux *http.ServeMux, cfg Config) error { return nil }
func (a fakeAdapter) FakeEnv(cfg Config) map[string]string {
	return map[string]string{"FAKE_URL": "http://localhost"}
}
func (a fakeAdapter) Metadata() Metadata {
	return Metadata{Name: a.name, Maturity: "experimental"}
}

func TestRegistryReturnsRegisteredAdapter(t *testing.T) {
	reg := NewRegistry()
	reg.Register(fakeAdapter{name: "stripe"})

	got, ok := reg.Get("stripe")
	if !ok {
		t.Fatal("expected registered adapter")
	}
	if got.Name() != "stripe" {
		t.Fatalf("adapter name = %q, want stripe", got.Name())
	}
}

func TestValidateMaturity(t *testing.T) {
	for _, maturity := range []string{"experimental", "partial", "sdk-compatible", "workflow-compatible", "provider-compatible"} {
		if !ValidateMaturity(maturity) {
			t.Fatalf("expected maturity %q to be valid", maturity)
		}
	}
	if ValidateMaturity("complete") {
		t.Fatal("unexpected valid maturity for complete")
	}
}
