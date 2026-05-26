package state

import (
	"fmt"
	"sync"
	"testing"
)

func TestStoreCreatesRetrievesListsUpdatesAndDeletesResources(t *testing.T) {
	store := NewStore()

	created, err := store.Create("stripe", "checkout_session", map[string]any{"amount_total": 1200})
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	if created.ID != "stripe_checkout_session_000001" {
		t.Fatalf("id = %q, want deterministic first id", created.ID)
	}

	got, ok := store.Get("stripe", "checkout_session", created.ID)
	if !ok {
		t.Fatal("created resource not found")
	}
	if got.Data["amount_total"] != 1200 {
		t.Fatalf("amount_total = %#v", got.Data["amount_total"])
	}

	updated, err := store.Update("stripe", "checkout_session", created.ID, map[string]any{"status": "complete"})
	if err != nil {
		t.Fatalf("update resource: %v", err)
	}
	if updated.Data["amount_total"] != 1200 || updated.Data["status"] != "complete" {
		t.Fatalf("updated data = %#v", updated.Data)
	}

	list := store.List("stripe", "checkout_session")
	if len(list) != 1 || list[0].ID != created.ID {
		t.Fatalf("list = %#v", list)
	}

	if !store.Delete("stripe", "checkout_session", created.ID) {
		t.Fatal("delete returned false")
	}
	if _, ok := store.Get("stripe", "checkout_session", created.ID); ok {
		t.Fatal("deleted resource still exists")
	}
}

func TestStoreResetClearsResourcesAndCounters(t *testing.T) {
	store := NewStore()
	first, err := store.Create("openai", "chat_completion", map[string]any{"model": "gpt-4o-mini"})
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	store.Reset("openai", "chat_completion")
	second, err := store.Create("openai", "chat_completion", map[string]any{"model": "gpt-4o-mini"})
	if err != nil {
		t.Fatalf("create second: %v", err)
	}
	if first.ID != second.ID {
		t.Fatalf("id after reset = %q, want %q", second.ID, first.ID)
	}
}

func TestStoreIsConcurrencySafeForDeterministicIDs(t *testing.T) {
	store := NewStore()
	const count = 25

	var wg sync.WaitGroup
	for i := range count {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if _, err := store.Create("slack", "message", map[string]any{"index": i}); err != nil {
				t.Errorf("create %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	seen := map[string]bool{}
	for _, resource := range store.List("slack", "message") {
		seen[resource.ID] = true
	}
	for i := 1; i <= count; i++ {
		id := fmt.Sprintf("slack_message_%06d", i)
		if !seen[id] {
			t.Fatalf("missing deterministic id %s from %#v", id, seen)
		}
	}
}
