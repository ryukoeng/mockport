package state

import "testing"

func BenchmarkStoreCreateList(b *testing.B) {
	store := NewStore()
	for i := 0; i < b.N; i++ {
		if _, err := store.Create("stripe", "payment_intent", map[string]any{"amount": i}); err != nil {
			b.Fatal(err)
		}
		_ = store.List("stripe", "payment_intent")
	}
}
