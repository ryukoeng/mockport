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

func TestStoreZeroValueIsUsable(t *testing.T) {
	var store Store

	created, err := store.Create("stripe", "checkout_session", map[string]any{"amount_total": 1200})
	if err != nil {
		t.Fatalf("create with zero-value store: %v", err)
	}
	if created.ID != "stripe_checkout_session_000001" {
		t.Fatalf("id = %q, want deterministic first id", created.ID)
	}
	if got, ok := store.Get("stripe", "checkout_session", created.ID); !ok || got.ID != created.ID {
		t.Fatalf("get from zero-value store = %#v, %v", got, ok)
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

func TestStoreTakeReturnsAndDeletesResource(t *testing.T) {
	store := NewStore()
	created, err := store.Create("github-oauth", "oauth_code", map[string]any{
		"metadata": map[string]any{"client_id": "mockport_github_client"},
	})
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}

	taken, ok := store.Take("github-oauth", "oauth_code", created.ID)
	if !ok {
		t.Fatal("take returned false")
	}
	taken.Data["metadata"].(map[string]any)["client_id"] = "mutated"

	if _, ok := store.Get("github-oauth", "oauth_code", created.ID); ok {
		t.Fatal("taken resource still exists")
	}
	if _, ok := store.Take("github-oauth", "oauth_code", created.ID); ok {
		t.Fatal("second take returned true")
	}
}

func TestStoreCapsResourcesPerScope(t *testing.T) {
	store := NewStore()
	for i := 0; i < MaxResourcesPerScope+5; i++ {
		if _, err := store.Create("stripe", "checkout_session", map[string]any{"index": i}); err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
	}

	resources := store.List("stripe", "checkout_session")
	if len(resources) != MaxResourcesPerScope {
		t.Fatalf("resource count = %d, want %d", len(resources), MaxResourcesPerScope)
	}
	if resources[0].ID != "stripe_checkout_session_000006" {
		t.Fatalf("first retained id = %q, want stripe_checkout_session_000006", resources[0].ID)
	}
	if _, ok := store.Get("stripe", "checkout_session", "stripe_checkout_session_000001"); ok {
		t.Fatal("oldest resource was retained after cap")
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

func TestStoreDeepClonesNestedData(t *testing.T) {
	store := NewStore()
	created, err := store.Create("stripe", "checkout_session", map[string]any{
		"metadata": map[string]any{"order": "one"},
		"items":    []any{map[string]any{"sku": "sku_1"}},
	})
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	created.Data["metadata"].(map[string]any)["order"] = "mutated"
	created.Data["items"].([]any)[0].(map[string]any)["sku"] = "mutated"

	got, ok := store.Get("stripe", "checkout_session", created.ID)
	if !ok {
		t.Fatal("resource not found")
	}
	if got.Data["metadata"].(map[string]any)["order"] != "one" {
		t.Fatalf("metadata was mutated through clone: %#v", got.Data)
	}
	if got.Data["items"].([]any)[0].(map[string]any)["sku"] != "sku_1" {
		t.Fatalf("items were mutated through clone: %#v", got.Data)
	}

	got.Data["metadata"].(map[string]any)["order"] = "get-mutated"
	again, _ := store.Get("stripe", "checkout_session", created.ID)
	if again.Data["metadata"].(map[string]any)["order"] != "one" {
		t.Fatalf("stored data mutated through get clone: %#v", again.Data)
	}
}

func TestStoreUpdateDeepClonesPatchData(t *testing.T) {
	store := NewStore()
	created, err := store.Create("line", "rich_menu", map[string]any{"name": "menu"})
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	patch := map[string]any{
		"metadata": map[string]any{"order": "one"},
		"items":    []any{map[string]any{"sku": "sku_1"}},
	}

	updated, err := store.Update("line", "rich_menu", created.ID, patch)
	if err != nil {
		t.Fatalf("update resource: %v", err)
	}
	patch["metadata"].(map[string]any)["order"] = "patch-mutated"
	patch["items"].([]any)[0].(map[string]any)["sku"] = "patch-mutated"
	updated.Data["metadata"].(map[string]any)["order"] = "result-mutated"
	updated.Data["items"].([]any)[0].(map[string]any)["sku"] = "result-mutated"

	got, ok := store.Get("line", "rich_menu", created.ID)
	if !ok {
		t.Fatal("resource not found")
	}
	if got.Data["metadata"].(map[string]any)["order"] != "one" {
		t.Fatalf("metadata was mutated through update patch/result: %#v", got.Data)
	}
	if got.Data["items"].([]any)[0].(map[string]any)["sku"] != "sku_1" {
		t.Fatalf("items were mutated through update patch/result: %#v", got.Data)
	}
}
