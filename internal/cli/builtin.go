package cli

import (
	"github.com/albert-einshutoin/mockport/adapters/githuboauth"
	"github.com/albert-einshutoin/mockport/adapters/line"
	"github.com/albert-einshutoin/mockport/adapters/openai"
	"github.com/albert-einshutoin/mockport/adapters/slack"
	"github.com/albert-einshutoin/mockport/adapters/stripe"
	"github.com/albert-einshutoin/mockport/internal/adapter"
)

func builtinAdapters() []adapter.Adapter {
	return []adapter.Adapter{stripe.New(), openai.New(), githuboauth.New(), slack.New(), line.New()}
}

func builtinAdapterFor(name string) (adapter.Adapter, bool) {
	for _, adapterImpl := range builtinAdapters() {
		if adapterImpl.Name() == name {
			return adapterImpl, true
		}
	}
	return nil, false
}
