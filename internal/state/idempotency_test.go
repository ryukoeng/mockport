package state

import (
	"errors"
	"testing"
)

func TestIdempotencyStoreReplaysMatchingRequest(t *testing.T) {
	store := NewIdempotencyStore()
	response := IdempotentResponse{Status: 200, Body: map[string]any{"id": "cs_test_123"}}

	replayed, got, err := store.Remember("stripe", "key-1", "POST /v1/checkout/sessions amount=1200", response)
	if err != nil {
		t.Fatalf("remember first request: %v", err)
	}
	if replayed {
		t.Fatal("first request replayed = true, want false")
	}

	replayed, got, err = store.Remember("stripe", "key-1", "POST /v1/checkout/sessions amount=1200", IdempotentResponse{Status: 500})
	if err != nil {
		t.Fatalf("remember replay request: %v", err)
	}
	if !replayed {
		t.Fatal("second matching request replayed = false, want true")
	}
	if got.Status != response.Status || got.Body["id"] != "cs_test_123" {
		t.Fatalf("replayed response = %#v", got)
	}
}

func TestIdempotencyStoreRejectsConflictingRequest(t *testing.T) {
	store := NewIdempotencyStore()
	if _, _, err := store.Remember("stripe", "key-1", "amount=1200", IdempotentResponse{Status: 200}); err != nil {
		t.Fatalf("remember first request: %v", err)
	}

	_, _, err := store.Remember("stripe", "key-1", "amount=999", IdempotentResponse{Status: 200})
	var conflict *IdempotencyConflictError
	if !errors.As(err, &conflict) {
		t.Fatalf("err = %v, want IdempotencyConflictError", err)
	}
	if conflict.Scope != "stripe" || conflict.Key != "key-1" {
		t.Fatalf("conflict = %#v", conflict)
	}
}

func TestRequireFieldsReportsMissingRequiredFields(t *testing.T) {
	err := RequireFields(map[string]any{
		"amount":   1200,
		"currency": "",
	}, "amount", "currency", "customer")

	var missing *ValidationError
	if !errors.As(err, &missing) {
		t.Fatalf("err = %v, want ValidationError", err)
	}
	if got, want := missing.MissingFields, []string{"currency", "customer"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("missing fields = %#v, want %#v", got, want)
	}
}
