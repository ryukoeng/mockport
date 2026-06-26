package httpx

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

type failingResponseWriter struct {
	header http.Header
	status int
}

func (w *failingResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = http.Header{}
	}
	return w.header
}

func (w *failingResponseWriter) WriteHeader(status int) { w.status = status }

func (w *failingResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestWriteJSONReturnsEncodeError(t *testing.T) {
	writer := &failingResponseWriter{}

	err := WriteJSON(writer, http.StatusAccepted, map[string]string{"ok": "true"})
	if err == nil {
		t.Fatalf("WriteJSON() error = %v, want encode json response", err)
	}
	errText := err.Error()
	if !strings.Contains(errText, "encode json response") {
		t.Fatalf("WriteJSON() error = %v, want encode json response", err)
	}
	if writer.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Content-Type = %q", writer.Header().Get("Content-Type"))
	}
	if writer.status != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", writer.status, http.StatusAccepted)
	}
}
