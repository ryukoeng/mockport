# Phase 12 Compatibility Engine Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Mockport の互換性を「感覚」ではなく、adapter ごとの contract、coverage、score として測れる基盤にする。

**Architecture:** 既存の Phase 4 metadata/report contract を拡張し、wire/schema/sdk/workflow/state/error の互換レベルを adapter ごとに表現する。Full-compatible は provider 内部再現ではなく、公開 API 契約、SDK、主要 workflow、fake state、主要 error を local Docker API として高忠実度に再現することとして定義する。

**Tech Stack:** Go 1.26.3, JSON/YAML compatibility manifests, report model, CLI report output, shell/static checks.

---

## Files

- Create: `internal/compat/manifest.go`
- Create: `internal/compat/manifest_test.go`
- Create: `internal/compat/score.go`
- Create: `internal/compat/score_test.go`
- Create: `docs/compatibility-model.md`
- Modify: `internal/adapter/adapter.go`
- Modify: `internal/report/report.go`
- Modify: `internal/report/render.go`
- Modify: `tasks/status.md`

## Task P12-T01: Compatibility Level Model

**Status:** pending

- [ ] Write failing tests for levels: `wire`, `sdk`, `workflow`, `error`, `state`, `contract`.
- [ ] Implement compatibility manifest structs with adapter, provider version, SDK versions, endpoints, scenarios, unsupported behavior.
- [ ] Clarify built-in scenario policy versus user-defined scenarios before scoring scenario coverage. See [issue #5](https://github.com/albert-einshutoin/mockport/issues/5).
- [ ] Add validation for required fields and duplicate endpoint ids.
- [ ] Run `/usr/local/go/bin/go test ./internal/compat -v`.

## Task P12-T02: Compatibility Score

**Status:** pending

- [ ] Write failing score tests for endpoint coverage, scenario coverage, SDK smoke status, state support, and error support.
- [ ] Implement deterministic score calculation without network access.
- [ ] Expose score in report JSON.
- [ ] Run `/usr/local/go/bin/go test ./internal/compat ./internal/report -v`.

## Task P12-T03: Compatibility Report Output

**Status:** pending

- [ ] Add report fields for compatibility level, score, SDK versions tested, provider API version, unsupported endpoints.
- [ ] Add text report rendering.
- [ ] Update docs to define what Mockport means by provider-compatible local API.
- [ ] Run `/usr/local/go/bin/go test ./...`.

## Phase 12 Exit

- [ ] Compatibility is represented as a structured manifest.
- [ ] Report shows compatibility score and levels.
- [ ] Full-compatible is explicitly defined as external API/SDK/workflow compatibility, not provider internals.
- [ ] Unsupported and approximate behavior remains visible.
