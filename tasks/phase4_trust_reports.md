# Phase 4 Trust Reports And Adapter Contracts Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** adapter を増やす前に、Mockport が「何をサポートし、何をサポートしないか」を report で説明できる trust foundation を作る。

**Architecture:** Adapter metadata と request recorder を統合し、scenario coverage、unsupported endpoint、request replay metadata、behavior matrix、maturity levels を JSON と text CLI の両方で出す。Phase 5 の追加 adapter はこの contract に必ず乗せる。

**Tech Stack:** Go 1.26.3, JSON report schema, text renderer tests, adapter metadata contract.

---

## Files

- Modify: `internal/adapter/adapter.go`
- Modify: `internal/adapter/registry.go`
- Modify: `internal/adapter/registry_test.go`
- Modify: `adapters/stripe/adapter.go`
- Modify: `adapters/stripe/adapter_test.go`
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

## Task P4-T01: Adapter Metadata Contract

**Status:** pending

- [ ] Write failing adapter tests asserting each adapter exposes name, maturity, capabilities, scenarios, and supported endpoints.
- [ ] Extend adapter types with metadata without introducing dynamic plugins.
- [ ] Add Stripe metadata: maturity `partial`, scenarios `payment_success`, `payment_failed`, `auth_error`, `rate_limited`, `timeout`.
- [ ] Add endpoint metadata for Stripe checkout session, payment intent, and webhook send.
- [ ] Run `/usr/local/go/bin/go test ./internal/adapter ./adapters/stripe -v`.

## Task P4-T02: Scenario Coverage Report

**Status:** pending

- [ ] Write failing report tests asserting Stripe reports scenario names and supported status.
- [ ] Add report schema field `scenario_coverage`.
- [ ] Generate scenario coverage from adapter metadata.
- [ ] Run `/usr/local/go/bin/go test ./internal/report ./internal/server -run Scenario -v`.

## Task P4-T03: Unsupported Endpoint Report

**Status:** pending

- [ ] Write failing HTTP test: request unknown adapter path, then report includes method/path/status and reason `unsupported_endpoint`.
- [ ] Update middleware to record 404 and 405 responses with classification.
- [ ] Avoid recording `/_mockport/report` itself.
- [ ] Run `/usr/local/go/bin/go test ./internal/server -run Unsupported -v`.

## Task P4-T04: Request Replay Log Metadata

**Status:** pending

- [ ] Write failing recorder tests for stable request id, timestamp, method, path, status, adapter, scenario.
- [ ] Implement monotonic request ids and replay-safe metadata.
- [ ] Ensure no request body or secret header is stored by default.
- [ ] Run `/usr/local/go/bin/go test ./internal/report -run Replay -v`.

## Task P4-T05: Behavior Matrix And Maturity Levels

**Status:** pending

- [ ] Write failing report tests for `behavior_matrix` with endpoint, method, supported scenarios, notes, and maturity.
- [ ] Validate allowed maturity values: `experimental`, `partial`, `common-path`, `contract-tested`, `sandbox-verified`.
- [ ] Include maturity in report and README adapter table.
- [ ] Run `/usr/local/go/bin/go test ./internal/adapter ./internal/report ./internal/server -run 'Behavior|Maturity' -v`.

## Task P4-T06: JSON And Text Report Modes

**Status:** pending

- [ ] Write failing CLI tests for `mockport report --format json` and `--format text`.
- [ ] Implement `internal/report/render.go`.
- [ ] Text output must include adapters, requests, safety, coverage, unsupported endpoints, and maturity.
- [ ] Create `docs/reporting.md` with interpretation examples.
- [ ] Run `/usr/local/go/bin/go test ./internal/cli ./internal/report -run Report -v`.

## Phase 4 Exit

- [ ] Report explains supported and unsupported behavior.
- [ ] Report includes scenario coverage and adapter maturity.
- [ ] Unsupported endpoints are visible.
- [ ] Text and JSON report outputs are tested.
- [ ] Docs explain how to interpret report output.
- [ ] Phase 5 adapters can be added without inventing a new report contract.
