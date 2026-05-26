# Phase 16 State Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** provider-specific state adoption の前に、deterministic fake state と idempotency/validation primitives を共有基盤として用意する。

**Architecture:** adapter ごとに重複した in-memory state を作らず、concurrency-safe store、deterministic ID、reset、idempotency key、validation helper、report hook を整える。内部 provider logic は再現せず、公開 API に見える resource lifecycle の土台に限定する。

**Tech Stack:** Go 1.26.3, in-memory stores, deterministic IDs, httptest, report metadata.

---

## Files

- Create: `internal/state/store.go`
- Create: `internal/state/store_test.go`
- Create: `internal/state/idempotency.go`
- Create: `internal/state/idempotency_test.go`
- Create: `docs/state-model.md`
- Modify: `internal/report/report.go`
- Modify: `tasks/status.md`

## Task P16-T01: Deterministic State Store

**Status:** pending

- [ ] Write failing tests for create/retrieve/list/update/delete semantics.
- [ ] Implement concurrency-safe in-memory store with deterministic IDs per adapter and resource type.
- [ ] Add reset support for test isolation.
- [ ] Run `/usr/local/go/bin/go test ./internal/state -v`.

## Task P16-T02: Idempotency And Validation Primitives

**Status:** pending

- [ ] Write failing tests for idempotency replay, conflict detection, and missing required field errors.
- [ ] Implement provider-neutral primitives that adapters can map to provider-shaped errors.
- [ ] Document where provider-specific error shape remains adapter-owned.
- [ ] Run `/usr/local/go/bin/go test ./internal/state -v`.

## Task P16-T03: State Coverage Reporting Hooks

**Status:** pending

- [ ] Add report metadata for stateful resources, idempotency support, and reset behavior.
- [ ] Add text/JSON report tests for state coverage.
- [ ] Document state model limitations.
- [ ] Run `/usr/local/go/bin/go test ./internal/report ./internal/state -v`.

## Phase 16 Exit

- [ ] Deterministic state store exists and is tested.
- [ ] Idempotency and validation primitives are available for adapters.
- [ ] Report can identify stateful workflow coverage.
- [ ] Adapter-wide migration is intentionally left for Phase 17.
