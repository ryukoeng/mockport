package compat

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

type Level string

const (
	LevelWire     Level = "wire"
	LevelSDK      Level = "sdk"
	LevelWorkflow Level = "workflow"
	LevelError    Level = "error"
	LevelState    Level = "state"
	LevelContract Level = "contract"
)

type SDKVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Manifest struct {
	Adapter         string                `json:"adapter"`
	ProviderVersion string                `json:"provider_version"`
	Maturity        string                `json:"maturity,omitempty"`
	SDKVersions     []SDKVersion          `json:"sdk_versions,omitempty"`
	Levels          []Level               `json:"levels,omitempty"`
	Endpoints       []Endpoint            `json:"endpoints,omitempty"`
	Scenarios       []Scenario            `json:"scenarios,omitempty"`
	Unsupported     []UnsupportedBehavior `json:"unsupported_behavior,omitempty"`
}

type Endpoint struct {
	ID        string  `json:"id"`
	Method    string  `json:"method"`
	Path      string  `json:"path"`
	Supported bool    `json:"supported"`
	Levels    []Level `json:"levels,omitempty"`
}

type Scenario struct {
	Name      string  `json:"name"`
	BuiltIn   bool    `json:"built_in"`
	Supported bool    `json:"supported"`
	Levels    []Level `json:"levels,omitempty"`
}

type UnsupportedBehavior struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
}

func (m Manifest) Validate() error {
	if strings.TrimSpace(m.Adapter) == "" {
		return fmt.Errorf("adapter is required")
	}
	if strings.TrimSpace(m.ProviderVersion) == "" {
		return fmt.Errorf("provider_version is required")
	}
	for _, level := range m.Levels {
		if !ValidLevel(level) {
			return fmt.Errorf("invalid compatibility level: %s", level)
		}
	}
	ids := map[string]bool{}
	for _, endpoint := range m.Endpoints {
		if strings.TrimSpace(endpoint.ID) == "" {
			return fmt.Errorf("endpoint id is required")
		}
		if ids[endpoint.ID] {
			return fmt.Errorf("duplicate endpoint id: %s", endpoint.ID)
		}
		ids[endpoint.ID] = true
		for _, level := range endpoint.Levels {
			if !ValidLevel(level) {
				return fmt.Errorf("invalid compatibility level: %s", level)
			}
		}
	}
	scenarios := map[string]bool{}
	for _, scenario := range m.Scenarios {
		if strings.TrimSpace(scenario.Name) == "" {
			return fmt.Errorf("scenario name is required")
		}
		if scenarios[scenario.Name] {
			return fmt.Errorf("duplicate scenario name: %s", scenario.Name)
		}
		scenarios[scenario.Name] = true
		for _, level := range scenario.Levels {
			if !ValidLevel(level) {
				return fmt.Errorf("invalid compatibility level: %s", level)
			}
		}
	}
	return nil
}

func ValidLevel(level Level) bool {
	switch level {
	case LevelWire, LevelSDK, LevelWorkflow, LevelError, LevelState, LevelContract:
		return true
	default:
		return false
	}
}

func FromMetadata(meta adapter.Metadata) Manifest {
	manifest := Manifest{
		Adapter:         meta.Name,
		ProviderVersion: "unspecified",
		Maturity:        meta.Maturity,
		Levels:          []Level{LevelWire},
	}
	for _, endpoint := range meta.Endpoints {
		manifest.Endpoints = append(manifest.Endpoints, Endpoint{
			ID:        endpointID(endpoint.Method, endpoint.Path),
			Method:    endpoint.Method,
			Path:      endpoint.Path,
			Supported: len(endpoint.SupportedScenarios) > 0 || endpoint.Method == http.MethodGet,
			Levels:    []Level{LevelWire},
		})
	}
	for _, scenario := range meta.Scenarios {
		manifest.Scenarios = append(manifest.Scenarios, Scenario{
			Name:      scenario.Name,
			BuiltIn:   true,
			Supported: scenario.Supported,
			Levels:    []Level{LevelWire},
		})
	}
	return manifest
}

func endpointID(method, path string) string {
	value := strings.ToLower(method + "_" + strings.Trim(path, "/"))
	value = strings.ReplaceAll(value, "/", "_")
	value = strings.ReplaceAll(value, "-", "_")
	value = regexp.MustCompile(`[^a-z0-9_]+`).ReplaceAllString(value, "_")
	value = strings.Trim(value, "_")
	if value == "" {
		return "endpoint"
	}
	return value
}
