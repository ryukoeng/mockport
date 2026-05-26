package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkWriteJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		if err := WriteJSON(rec, http.StatusOK, map[string]string{"status": "ok"}); err != nil {
			b.Fatal(err)
		}
	}
}
