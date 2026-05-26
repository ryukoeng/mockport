package security

import (
	"encoding/json"
	"fmt"
	"strings"
)

type FixtureFinding struct {
	Path   string
	Field  string
	Reason string
}

func ScanFixtureContent(path string, content string) []FixtureFinding {
	var payload map[string]any
	if err := json.Unmarshal([]byte(content), &payload); err != nil {
		return []FixtureFinding{{Path: path, Field: "$", Reason: "invalid JSON"}}
	}

	var findings []FixtureFinding
	source, ok := payload["source"].(map[string]any)
	if !ok {
		findings = append(findings, FixtureFinding{Path: path, Field: "source", Reason: "missing source metadata"})
	} else {
		for _, field := range []string{"type", "title", "url_or_path", "retrieved_at"} {
			if _, ok := source[field]; !ok {
				findings = append(findings, FixtureFinding{Path: path, Field: "source." + field, Reason: "missing source metadata"})
			}
		}
	}
	for _, field := range []string{"provider", "provider_version", "sdk", "request", "response"} {
		if _, ok := payload[field]; !ok {
			findings = append(findings, FixtureFinding{Path: path, Field: field, Reason: "missing required field"})
		}
	}
	scanFixtureValue(path, "$", payload, &findings)
	return findings
}

func scanFixtureValue(path string, field string, value any, findings *[]FixtureFinding) {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			scanFixtureValue(path, field+"."+key, child, findings)
		}
	case []any:
		for idx, child := range typed {
			scanFixtureValue(path, fmt.Sprintf("%s[%d]", field, idx), child, findings)
		}
	case string:
		for _, token := range strings.Fields(typed) {
			token = strings.Trim(token, ".,;:()[]{}\"'")
			switch {
			case LooksLikeSecret(token):
				*findings = append(*findings, FixtureFinding{Path: path, Field: field, Reason: "real-looking provider secret"})
			case LooksLikeExternalServiceURL(token):
				*findings = append(*findings, FixtureFinding{Path: path, Field: field, Reason: "production provider URL"})
			}
		}
	}
}
