package report

import (
	"strings"
	"sync"
	"time"
)

type Recorder struct {
	mu             sync.Mutex
	mode           string
	adapters       []AdapterStatus
	requests       []Request
	safetyWarnings []SafetyWarning
	coverage       []ScenarioCoverage
	matrix         []BehaviorMatrixEntry
	nextID         int64
}

func NewRecorder() *Recorder {
	return &Recorder{}
}

func (r *Recorder) SetMode(mode string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.mode = mode
}

func (r *Recorder) SetAdapters(adapters []AdapterStatus) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adapters = append([]AdapterStatus(nil), adapters...)
}

func (r *Recorder) SetScenarioCoverage(coverage []ScenarioCoverage) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.coverage = append([]ScenarioCoverage(nil), coverage...)
}

func (r *Recorder) SetBehaviorMatrix(matrix []BehaviorMatrixEntry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.matrix = append([]BehaviorMatrixEntry(nil), matrix...)
}

func (r *Recorder) RecordRequest(method, path string, status int) {
	r.RecordRequestWithDetails(method, path, status, "", "", "")
}

func (r *Recorder) RecordRequestWithDetails(method, path string, status int, adapter, scenario, reason string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	r.requests = append(r.requests, Request{
		ID:        r.nextID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Method:    method,
		Path:      path,
		Status:    status,
		Adapter:   adapter,
		Scenario:  scenario,
		Reason:    reason,
	})
}

func (r *Recorder) RecordSafetyWarning(field, category, message string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.safetyWarnings = append(r.safetyWarnings, SafetyWarning{Field: field, Category: category, Message: message})
}

func (r *Recorder) Snapshot() Snapshot {
	r.mu.Lock()
	defer r.mu.Unlock()
	return Snapshot{
		Mode:                 r.mode,
		Safety:               safetySummary(r.mode, r.safetyWarnings),
		Adapters:             append([]AdapterStatus(nil), r.adapters...),
		Requests:             append([]Request(nil), r.requests...),
		SafetyWarnings:       append([]SafetyWarning(nil), r.safetyWarnings...),
		ScenarioCoverage:     append([]ScenarioCoverage(nil), r.coverage...),
		BehaviorMatrix:       append([]BehaviorMatrixEntry(nil), r.matrix...),
		UnsupportedEndpoints: unsupportedEndpoints(r.requests),
	}
}

func safetySummary(mode string, warnings []SafetyWarning) SafetySummary {
	summary := SafetySummary{Mode: mode, Safe: len(warnings) == 0}
	for _, warning := range warnings {
		switch warning.Category {
		case "real_looking_secret":
			summary.RealLookingSecrets++
		case "external_url":
			summary.ExternalURLs++
		}
	}
	return summary
}

func unsupportedEndpoints(requests []Request) []UnsupportedEndpoint {
	var unsupported []UnsupportedEndpoint
	for _, request := range requests {
		if strings.TrimSpace(request.Reason) == "" {
			continue
		}
		unsupported = append(unsupported, UnsupportedEndpoint{
			Method: request.Method,
			Path:   request.Path,
			Status: request.Status,
			Reason: request.Reason,
		})
	}
	return unsupported
}
