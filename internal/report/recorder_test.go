package report

import (
	"net/http"
	"testing"
)

func TestRecorderStoresRequestsAndSafety(t *testing.T) {
	rec := NewRecorder()
	rec.RecordRequest(http.MethodPost, "/stripe/v1/checkout/sessions", 200)
	rec.RecordSafetyWarning("STRIPE_SECRET_KEY", "real_looking_secret", "real-looking Stripe key")

	snapshot := rec.Snapshot()
	if len(snapshot.Requests) != 1 {
		t.Fatalf("request count = %d, want 1", len(snapshot.Requests))
	}
	if snapshot.Requests[0].Path != "/stripe/v1/checkout/sessions" {
		t.Fatalf("path = %q", snapshot.Requests[0].Path)
	}
	if len(snapshot.SafetyWarnings) != 1 {
		t.Fatalf("warning count = %d, want 1", len(snapshot.SafetyWarnings))
	}
	if snapshot.Safety.Safe {
		t.Fatal("safety summary safe = true, want false")
	}
	if snapshot.Safety.RealLookingSecrets != 1 {
		t.Fatalf("real-looking secret count = %d, want 1", snapshot.Safety.RealLookingSecrets)
	}
}

func TestRecorderStoresReplayMetadataAndUnsupportedEndpoints(t *testing.T) {
	rec := NewRecorder()
	rec.RecordRequestWithDetails(http.MethodPost, "/stripe/not-supported", 404, "stripe", "payment_success", "unsupported_endpoint")

	snapshot := rec.Snapshot()
	if len(snapshot.Requests) != 1 {
		t.Fatalf("request count = %d, want 1", len(snapshot.Requests))
	}
	req := snapshot.Requests[0]
	if req.ID != 1 {
		t.Fatalf("request id = %d, want 1", req.ID)
	}
	if req.Timestamp == "" {
		t.Fatal("request timestamp is empty")
	}
	if req.Adapter != "stripe" || req.Scenario != "payment_success" {
		t.Fatalf("request metadata = %#v", req)
	}
	if len(snapshot.UnsupportedEndpoints) != 1 {
		t.Fatalf("unsupported endpoints = %d, want 1", len(snapshot.UnsupportedEndpoints))
	}
}

func TestRecorderStoresCompatibility(t *testing.T) {
	rec := NewRecorder()
	rec.SetCompatibility([]CompatibilityStatus{{
		Adapter:         "stripe",
		Level:           "wire",
		Score:           80,
		ProviderVersion: "2026-05-26",
		SDKVersions:     []string{"stripe-go@v83.0.0"},
	}})

	snapshot := rec.Snapshot()
	if len(snapshot.Compatibility) != 1 {
		t.Fatalf("compatibility count = %d, want 1", len(snapshot.Compatibility))
	}
	if snapshot.Compatibility[0].Adapter != "stripe" || snapshot.Compatibility[0].Score != 80 {
		t.Fatalf("compatibility = %#v", snapshot.Compatibility[0])
	}
}
