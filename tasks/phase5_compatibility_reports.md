# Phase 5 Compatibility And Reports Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Mockport の report が「何が対応済みで、何が未対応か」を説明できる信頼性 artifact になる。

**Architecture:** Adapter metadata と request recorder を統合し、scenario coverage、unsupported endpoint、request replay metadata、behavior matrix、maturity levels を JSON と text CLI の両方で出す。

**Tech Stack:** Go 1.26.3, JSON report schema, text renderer tests.

---

## Files

- Modify: `internal/report/report.go`
- Modify: `internal/report/recorder.go`
- Modify: `internal/report/recorder_test.go`
- Modify: `internal/server/server.go`
- Modify: `internal/server/report_test.go`
- Create: `internal/report/render.go`
- Create: `internal/report/render_test.go`
- Modify: `internal/cli/report.go`
- Modify: `internal/cli/report_test.go`
- Create: `docs/reporting.md`

## Task P5-T01: Scenario Coverage Report

**Status:** pending

- [ ] Write failing report tests asserting each adapter reports scenario names and supported status.
- [ ] Extend adapter metadata to expose scenario coverage.
- [ ] Add report schema field `scenario_coverage`.
- [ ] Run `/usr/local/go/bin/go test ./internal/report ./internal/server -run Scenario -v`.

## Task P5-T02: Unsupported Endpoint Report

**Status:** pending

- [ ] Write failing HTTP test: request unknown adapter path, then report includes method/path/status and reason `unsupported_endpoint`.
- [ ] Update middleware to record 404 and 405 responses with classification.
- [ ] Avoid recording `/_mockport/report` itself.
- [ ] Run `/usr/local/go/bin/go test ./internal/server -run Unsupported -v`.

## Task P5-T03: Request Replay Log Metadata

**Status:** pending

- [ ] Write failing recorder tests for stable request id, timestamp, method, path, status, adapter, scenario.
- [ ] Implement monotonic request ids and replay-safe metadata.
- [ ] Ensure no request body or secret header is stored by default.
- [ ] Run `/usr/local/go/bin/go test ./internal/report -run Replay -v`.

## Task P5-T04: Behavior Matrix

**Status:** pending

- [ ] Write failing report tests for `behavior_matrix` with endpoint, method, supported scenarios, notes.
- [ ] Generate matrix from adapter metadata.
- [ ] Include Stripe, OpenAI, GitHub OAuth, Slack entries when adapters exist.
- [ ] Run `/usr/local/go/bin/go test ./internal/report ./internal/server -run Behavior -v`.

## Task P5-T05: Adapter Maturity Levels

**Status:** pending

- [ ] Write failing tests asserting allowed maturity values: `experimental`, `partial`, `common-path`, `contract-tested`, `sandbox-verified`.
- [ ] Add maturity validation to adapter metadata.
- [ ] Include maturity in report and README adapter table.
- [ ] Run `/usr/local/go/bin/go test ./internal/adapter ./internal/report -run Maturity -v`.

## Task P5-T06: JSON And Text Report Modes

**Status:** pending

- [ ] Write failing CLI tests for `mockport report --format json` and `--format text`.
- [ ] Implement `internal/report/render.go`.
- [ ] Text output must include adapters, requests, safety, coverage, unsupported endpoints.
- [ ] Run `/usr/local/go/bin/go test ./internal/cli ./internal/report -run Report -v`.

## Phase 5 Exit

- [ ] Report explains supported and unsupported behavior.
- [ ] Report includes scenario coverage and maturity.
- [ ] Unsupported endpoints are visible.
- [ ] Text and JSON report outputs are tested.
- [ ] Docs explain how to interpret report output.
