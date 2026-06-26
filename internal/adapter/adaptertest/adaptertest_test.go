package adaptertest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConcurrentStatusCodesReturnsNResults(t *testing.T) {
	codes := ConcurrentStatusCodes(5, func() int { return http.StatusTeapot })
	if len(codes) != 5 {
		t.Fatalf("len(codes) = %d, want 5", len(codes))
	}
	for i, code := range codes {
		if code != http.StatusTeapot {
			t.Fatalf("codes[%d] = %d, want %d", i, code, http.StatusTeapot)
		}
	}
}

func TestAssertJSONFieldNestedPath(t *testing.T) {
	rec := httptest.NewRecorder()
	rec.Body.WriteString(`{"error":{"code":"card_declined","type":"card_error"}}`)

	AssertJSONField(t, rec, "error.code", "card_declined")
}

func TestServeSetsHeadersAndBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody string
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body: %v", err)
		}
		gotBody = string(b)
		if r.Header.Get("X-Test") != "yes" {
			t.Errorf("X-Test header = %q, want yes", r.Header.Get("X-Test"))
		}
		w.WriteHeader(http.StatusOK)
	})

	header := http.Header{}
	header.Set("X-Test", "yes")
	rec := Serve(mux, http.MethodPost, "/test", strings.NewReader("payload"), header)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if gotMethod != http.MethodPost || gotPath != "/test" || gotBody != "payload" {
		t.Fatalf("got method=%q path=%q body=%q", gotMethod, gotPath, gotBody)
	}
}
