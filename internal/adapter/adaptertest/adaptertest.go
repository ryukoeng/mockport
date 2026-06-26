// Package adaptertest provides shared helpers for adapter HTTP tests.
package adaptertest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

// NewMux registers an adapter on a fresh mux and fails the test on error.
func NewMux(t *testing.T, a adapter.Adapter, cfg adapter.Config) *http.ServeMux {
	t.Helper()
	mux := http.NewServeMux()
	if err := a.Register(mux, cfg); err != nil {
		t.Fatalf("register adapter: %v", err)
	}
	return mux
}

// Serve performs one request against mux and returns the recorder.
func Serve(mux http.Handler, method, path string, body io.Reader, header http.Header) *httptest.ResponseRecorder {
	return ServeWithRemote(mux, method, path, body, header, "")
}

// ServeWithRemote is Serve with an explicit RemoteAddr (for loopback-only reset tests).
func ServeWithRemote(mux http.Handler, method, path string, body io.Reader, header http.Header, remoteAddr string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	if remoteAddr != "" {
		req.RemoteAddr = remoteAddr
	}
	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

// ConcurrentStatusCodes runs n copies of request concurrently (released by a
// shared start signal) and returns the resulting status codes.
func ConcurrentStatusCodes(n int, request func() int) []int {
	var wg sync.WaitGroup
	start := make(chan struct{})
	codes := make([]int, n)
	for i := range codes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			codes[i] = request()
		}()
	}
	close(start)
	wg.Wait()
	return codes
}

// ConcurrentResults runs n copies of fn concurrently and collects each result.
func ConcurrentResults[T any](n int, fn func() T) []T {
	var wg sync.WaitGroup
	start := make(chan struct{})
	results := make([]T, n)
	for i := range results {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			results[i] = fn()
		}()
	}
	close(start)
	wg.Wait()
	return results
}

// AssertJSONField fails unless the response body JSON has want at the
// dot-separated path (e.g. "error.code").
func AssertJSONField(t *testing.T, rec *httptest.ResponseRecorder, path, want string) {
	t.Helper()
	var decoded any
	if err := json.Unmarshal(rec.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("decode JSON body: %v", err)
	}
	got, ok := jsonPath(decoded, path)
	if !ok {
		t.Fatalf("JSON path %q not found in body: %s", path, rec.Body.String())
	}
	gotStr, ok := got.(string)
	if !ok {
		t.Fatalf("JSON path %q = %#v (%T), want string %q", path, got, got, want)
	}
	if gotStr != want {
		t.Fatalf("JSON path %q = %q, want %q", path, gotStr, want)
	}
}

func jsonPath(v any, path string) (any, bool) {
	if path == "" {
		return v, true
	}
	parts := strings.Split(path, ".")
	cur := v
	for _, part := range parts {
		obj, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		next, ok := obj[part]
		if !ok {
			return nil, false
		}
		cur = next
	}
	return cur, true
}
