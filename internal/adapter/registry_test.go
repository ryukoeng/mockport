package adapter

import (
	"net/http"
	"strings"
	"testing"
)

type fakeAdapter struct{ name string }

func (a fakeAdapter) Name() string                                  { return a.name }
func (a fakeAdapter) Register(mux *http.ServeMux, cfg Config) error { return nil }
func (a fakeAdapter) FakeEnv(cfg Config) map[string]string {
	return map[string]string{"FAKE_URL": "http://localhost"}
}
func (a fakeAdapter) Metadata() Metadata {
	return Metadata{Name: a.name, Maturity: MaturityExperimental}
}

func TestRegistryReturnsRegisteredAdapter(t *testing.T) {
	reg := NewRegistry()
	if err := reg.Register(fakeAdapter{name: "stripe"}); err != nil {
		t.Fatalf("register adapter: %v", err)
	}

	got, ok := reg.Get("stripe")
	if !ok {
		t.Fatal("expected registered adapter")
	}
	if got.Name() != "stripe" {
		t.Fatalf("adapter name = %q, want stripe", got.Name())
	}
}

func TestValidateMaturity(t *testing.T) {
	for _, maturity := range []Maturity{
		MaturityExperimental,
		MaturityPartial,
		MaturitySDKCompatible,
		MaturityWorkflowCompatible,
		MaturityProviderCompatible,
	} {
		if !ValidateMaturity(maturity) {
			t.Fatalf("expected maturity %q to be valid", maturity)
		}
	}
	if ValidateMaturity(Maturity("complete")) {
		t.Fatal("unexpected valid maturity for complete")
	}
}

func TestCompatibilityLevelsAreTypedConstants(t *testing.T) {
	meta := Metadata{
		Name:     "fake",
		Maturity: MaturityWorkflowCompatible,
		Levels:   []Level{LevelWire, LevelSDK, LevelWorkflow, LevelState, LevelError},
	}
	if meta.Maturity != MaturityWorkflowCompatible {
		t.Fatalf("maturity = %q", meta.Maturity)
	}
	if len(meta.Levels) != 5 || meta.Levels[0] != LevelWire {
		t.Fatalf("levels = %#v", meta.Levels)
	}
}

func TestValidateMetadataRejectsInvalidContracts(t *testing.T) {
	tests := []struct {
		name string
		meta Metadata
		want string
	}{
		{
			name: "invalid maturity",
			meta: Metadata{Name: "stripe", Maturity: Maturity("complete"), ProviderVersion: "2025-01-01", Levels: []Level{LevelWire}},
			want: "invalid maturity",
		},
		{
			name: "invalid level",
			meta: Metadata{Name: "stripe", Maturity: MaturityExperimental, ProviderVersion: "2025-01-01", Levels: []Level{"complete"}},
			want: "invalid compatibility level",
		},
		{
			name: "missing provider version",
			meta: Metadata{Name: "stripe", Maturity: MaturityExperimental, Levels: []Level{LevelWire}},
			want: "provider version is required",
		},
		{
			name: "duplicate scenario",
			meta: Metadata{
				Name: "stripe", Maturity: MaturityExperimental, ProviderVersion: "2025-01-01", Levels: []Level{LevelWire},
				Scenarios: []Scenario{{Name: "success"}, {Name: "success"}},
			},
			want: "duplicate scenario",
		},
		{
			name: "duplicate endpoint",
			meta: Metadata{
				Name: "stripe", Maturity: MaturityExperimental, ProviderVersion: "2025-01-01", Levels: []Level{LevelWire},
				Endpoints: []Endpoint{{Method: http.MethodGet, Path: "/v1/a"}, {Method: http.MethodGet, Path: "/v1/a"}},
			},
			want: "duplicate endpoint",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetadata(tt.meta)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("ValidateMetadata() error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestRegistryRejectsDuplicateAndInvalidAdapters(t *testing.T) {
	reg := NewRegistry()
	if err := reg.Register(fakeAdapter{name: "stripe"}); err != nil {
		t.Fatalf("register first: %v", err)
	}
	if err := reg.Register(fakeAdapter{name: "stripe"}); err == nil {
		t.Fatal("duplicate register returned nil error")
	}
	if err := reg.Register(fakeAdapter{name: ""}); err == nil {
		t.Fatal("empty-name register returned nil error")
	}
}
