package adapter

type Registry struct {
	adapters map[string]Adapter
}

func NewRegistry() *Registry {
	return &Registry{adapters: map[string]Adapter{}}
}

func (r *Registry) Register(adapter Adapter) {
	r.adapters[adapter.Name()] = adapter
}

func (r *Registry) Get(name string) (Adapter, bool) {
	adapter, ok := r.adapters[name]
	return adapter, ok
}
