package server

import (
	"testing"

	"github.com/albert-einshutoin/mockport/adapters/githuboauth"
	"github.com/albert-einshutoin/mockport/adapters/line"
	"github.com/albert-einshutoin/mockport/adapters/openai"
	"github.com/albert-einshutoin/mockport/adapters/slack"
	"github.com/albert-einshutoin/mockport/adapters/stripe"
	"github.com/albert-einshutoin/mockport/internal/adapter"
)

func TestBuiltInAdapterMetadataConformance(t *testing.T) {
	for _, adapterImpl := range []adapter.Adapter{stripe.New(), openai.New(), githuboauth.New(), slack.New(), line.New()} {
		t.Run(adapterImpl.Name(), func(t *testing.T) {
			meta := adapterImpl.Metadata()
			if err := adapter.ValidateMetadata(meta); err != nil {
				t.Fatalf("metadata invalid: %v", err)
			}
			scenarios := map[string]bool{}
			for _, scenario := range meta.Scenarios {
				scenarios[scenario.Name] = true
			}
			for _, endpoint := range meta.Endpoints {
				if len(endpoint.SupportedScenarios) == 0 {
					t.Fatalf("endpoint %s %s has no supported scenarios", endpoint.Method, endpoint.Path)
				}
				for _, scenario := range endpoint.SupportedScenarios {
					if !scenarios[scenario] {
						t.Fatalf("endpoint %s %s references unknown scenario %q", endpoint.Method, endpoint.Path, scenario)
					}
				}
			}
		})
	}
}
