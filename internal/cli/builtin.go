package cli

import (
	"github.com/albert-einshutoin/mockport/internal/adapter"
	"github.com/albert-einshutoin/mockport/internal/builtins"
)

func builtinAdapters() []adapter.Adapter {
	return builtins.Adapters()
}

func builtinAdapterFor(name string) (adapter.Adapter, bool) {
	for _, adapterImpl := range builtinAdapters() {
		if adapterImpl.Name() == name {
			return adapterImpl, true
		}
	}
	return nil, false
}
