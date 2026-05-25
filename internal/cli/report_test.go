package cli

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReportCommandPrintsTextSummary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/_mockport/report" {
			t.Fatalf("path = %q, want /_mockport/report", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"mode":"ai-safe","adapters":[{"name":"stripe","base_path":"/stripe","enabled":true}],"requests":[{"method":"POST","path":"/stripe/v1/checkout/sessions","status":200}],"safety_warnings":[]}`))
	}))
	defer server.Close()

	cmd := NewRootCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"report", "--url", server.URL + "/_mockport/report"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute report: %v", err)
	}

	got := out.String()
	for _, want := range []string{
		"Mockport Report",
		"Mode: ai-safe",
		"stripe enabled at /stripe",
		"POST /stripe/v1/checkout/sessions -> 200",
		"Safety warnings: 0",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("report output missing %q:\n%s", want, got)
		}
	}
}

func TestReportCommandPrintsJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"mode":"ai-safe","safety":{"mode":"ai-safe","safe":true},"adapters":[],"requests":[],"safety_warnings":[]}`))
	}))
	defer server.Close()

	cmd := NewRootCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"report", "--url", server.URL, "--format", "json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute report: %v", err)
	}
	if !strings.Contains(out.String(), `"mode": "ai-safe"`) {
		t.Fatalf("json report output missing mode:\n%s", out.String())
	}
}
