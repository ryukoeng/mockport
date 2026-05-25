package report

import "sync"

type Recorder struct {
	mu             sync.Mutex
	mode           string
	adapters       []AdapterStatus
	requests       []Request
	safetyWarnings []SafetyWarning
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

func (r *Recorder) RecordRequest(method, path string, status int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests = append(r.requests, Request{Method: method, Path: path, Status: status})
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
		Mode:           r.mode,
		Safety:         safetySummary(r.mode, r.safetyWarnings),
		Adapters:       append([]AdapterStatus(nil), r.adapters...),
		Requests:       append([]Request(nil), r.requests...),
		SafetyWarnings: append([]SafetyWarning(nil), r.safetyWarnings...),
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
