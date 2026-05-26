package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// MaxBodyBytes is the default per-request body limit for adapter handlers.
const MaxBodyBytes int64 = 1 << 20

// ErrRequestBodyTooLarge marks request bodies rejected by size limits.
var ErrRequestBodyTooLarge = errors.New("request body too large")

// WriteJSON writes a JSON response and returns encode/write failures.
func WriteJSON(w http.ResponseWriter, status int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		return fmt.Errorf("encode json response: %w", err)
	}
	return nil
}

// LimitRequestBody installs an http.MaxBytesReader on the request body.
func LimitRequestBody(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	}
}

// IsRequestBodyTooLarge reports whether err came from http.MaxBytesReader.
func IsRequestBodyTooLarge(err error) bool {
	var maxBytesErr *http.MaxBytesError
	return errors.As(err, &maxBytesErr)
}
