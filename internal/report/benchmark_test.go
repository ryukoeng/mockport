package report

import (
	"net/http"
	"testing"
)

func BenchmarkRecorderSnapshot(b *testing.B) {
	rec := NewRecorder()
	for i := 0; i < 100; i++ {
		rec.RecordRequest(http.MethodGet, "/health", http.StatusOK)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rec.Snapshot()
	}
}
