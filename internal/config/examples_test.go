package config

import "testing"

func TestExampleConfigsLoad(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		adapters []string
	}{
		{
			name:     "openai",
			path:     "../../examples/openai-chat/mockport.yml",
			adapters: []string{"openai"},
		},
		{
			name:     "github-oauth",
			path:     "../../examples/github-oauth/mockport.yml",
			adapters: []string{"github-oauth"},
		},
		{
			name:     "slack",
			path:     "../../examples/slack-message/mockport.yml",
			adapters: []string{"slack"},
		},
		{
			name:     "multi-adapter",
			path:     "../../examples/multi-adapter/mockport.yml",
			adapters: []string{"stripe", "openai", "github-oauth", "slack"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadFile(tt.path)
			if err != nil {
				t.Fatalf("LoadFile() error = %v", err)
			}
			for _, adapterName := range tt.adapters {
				adapterCfg, ok := cfg.Adapters[adapterName]
				if !ok {
					t.Fatalf("adapter %q missing from example config", adapterName)
				}
				if !adapterCfg.Enabled {
					t.Fatalf("adapter %q is disabled", adapterName)
				}
			}
		})
	}
}
