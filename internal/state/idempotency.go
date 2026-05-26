package state

import (
	"fmt"
	"strings"
	"sync"
)

type IdempotentResponse struct {
	Status int            `json:"status"`
	Body   map[string]any `json:"body,omitempty"`
}

type idempotencyRecord struct {
	fingerprint string
	response    IdempotentResponse
}

type IdempotencyStore struct {
	mu      sync.Mutex
	records map[string]idempotencyRecord
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{records: map[string]idempotencyRecord{}}
}

func (s *IdempotencyStore) Remember(scope, key, fingerprint string, response IdempotentResponse) (bool, IdempotentResponse, error) {
	if strings.TrimSpace(key) == "" {
		return false, cloneResponse(response), nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	recordKey := scope + "\x00" + key
	if record, ok := s.records[recordKey]; ok {
		if record.fingerprint != fingerprint {
			return false, IdempotentResponse{}, &IdempotencyConflictError{Scope: scope, Key: key}
		}
		return true, cloneResponse(record.response), nil
	}
	s.records[recordKey] = idempotencyRecord{fingerprint: fingerprint, response: cloneResponse(response)}
	return false, cloneResponse(response), nil
}

type IdempotencyConflictError struct {
	Scope string
	Key   string
}

func (e *IdempotencyConflictError) Error() string {
	return fmt.Sprintf("idempotency conflict for %s key %s", e.Scope, e.Key)
}

type ValidationError struct {
	MissingFields []string
}

func (e *ValidationError) Error() string {
	return "missing required fields: " + strings.Join(e.MissingFields, ", ")
}

func RequireFields(payload map[string]any, fields ...string) error {
	var missing []string
	for _, field := range fields {
		value, ok := payload[field]
		if !ok || isEmpty(value) {
			missing = append(missing, field)
		}
	}
	if len(missing) > 0 {
		return &ValidationError{MissingFields: missing}
	}
	return nil
}

func isEmpty(value any) bool {
	switch typed := value.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(typed) == ""
	default:
		return false
	}
}

func cloneResponse(response IdempotentResponse) IdempotentResponse {
	response.Body = cloneData(response.Body)
	return response
}
