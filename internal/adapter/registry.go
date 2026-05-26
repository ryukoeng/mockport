package adapter

import (
	"fmt"
	"strings"
)

// Registry stores built-in adapters by unique adapter name.
type Registry struct {
	adapters map[string]Adapter
}

// NewRegistry creates an empty adapter registry.
func NewRegistry() *Registry {
	return &Registry{adapters: map[string]Adapter{}}
}

// Register adds an adapter and rejects invalid or duplicate names.
func (r *Registry) Register(adapter Adapter) error {
	if adapter == nil {
		return fmt.Errorf("adapter is nil")
	}
	name := strings.TrimSpace(adapter.Name())
	if name == "" {
		return fmt.Errorf("adapter name is required")
	}
	if _, exists := r.adapters[name]; exists {
		return fmt.Errorf("duplicate adapter: %s", name)
	}
	r.adapters[name] = adapter
	return nil
}

// Get returns the adapter registered for name.
func (r *Registry) Get(name string) (Adapter, bool) {
	adapter, ok := r.adapters[name]
	return adapter, ok
}
