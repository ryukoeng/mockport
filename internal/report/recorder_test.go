package report

import (
	"net/http"
	"testing"
)

func TestRecorderStoresRequestsAndSafety(t *testing.T) {
	rec := NewRecorder()
	rec.RecordRequest(http.MethodPost, "/stripe/v1/checkout/sessions", 200)
	rec.RecordSafetyWarning("STRIPE_SECRET_KEY", "real-looking Stripe key")

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
}
