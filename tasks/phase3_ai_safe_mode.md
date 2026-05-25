# Phase 3 AI-safe Mode Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** AI-safe mode を Mockport の明確な差別化として、警告、strict fail、redaction、report、docs まで一貫させる。

**Architecture:** `internal/security` が検出・分類・redaction を担当し、`internal/config` が mode に応じて警告または error に変換する。`internal/report` は safe/unsafe の summary を持ち、CLI/server は full secret を表示しない。

**Tech Stack:** Go 1.26.3, table-driven tests, JSON report, Cobra output tests.

---

## Files

- Modify: `internal/security/secrets.go`
- Modify: `internal/security/secrets_test.go`
- Create: `internal/security/urls.go`
- Create: `internal/security/redactor.go`
- Modify: `internal/config/validate.go`
- Modify: `internal/config/config_test.go`
- Modify: `internal/report/report.go`
- Modify: `internal/server/report_test.go`
- Modify: `internal/cli/run.go`
- Modify: `internal/cli/run_test.go`
- Create: `docs/ai-safe-development.md`
- Create: `examples/unsafe-config/mockport.yml`

## Task P3-T01: Safety Status Model

**Status:** pending

- [ ] Write failing report test: unsafe config produces report fields `safety.safe=false`, `real_looking_secrets=1`, `external_urls=1`, `mode=ai-safe`.
- [ ] Implement `report.SafetySummary`.
- [ ] Wire config warnings into recorder summary.
- [ ] Run `/usr/local/go/bin/go test ./internal/report ./internal/server -run Safety -v`.

## Task P3-T02: Startup Warnings In `ai-safe`

**Status:** pending

- [ ] Write failing CLI test: `mockport run --config unsafe.yml --check` prints warning names and categories but not secret values.
- [ ] Add `--check` to `run` so tests avoid long-running server startup.
- [ ] Implement warning output for `ai-safe`.
- [ ] Run `/usr/local/go/bin/go test ./internal/cli -run Run -v`.

## Task P3-T03: Strict Mode Exit Behavior

**Status:** pending

- [ ] Write failing CLI test: strict unsafe config returns non-nil error containing `strict mode rejected unsafe config fields`.
- [ ] Ensure `config.Validate` returns category-rich error without secret values.
- [ ] Run `/usr/local/go/bin/go test ./internal/config ./internal/cli -run Strict -v`.

## Task P3-T04: Redaction Coverage

**Status:** pending

- [ ] Add table tests for short secrets, long secrets, fake secrets, env assignment strings, webhook secrets, provider URLs.
- [ ] Implement `security.RedactValue` and `security.RedactMessage`.
- [ ] Replace ad hoc redaction call sites with redactor functions.
- [ ] Run `/usr/local/go/bin/go test ./internal/security -v`.

## Task P3-T05: External URL Guard

**Status:** pending

- [ ] Add config tests for Stripe/OpenAI/GitHub/LINE/Slack real URLs.
- [ ] Move URL detection to `internal/security/urls.go`.
- [ ] In `strict`, reject all known provider URLs.
- [ ] In `ai-safe`, record warning categories.
- [ ] Run `/usr/local/go/bin/go test ./internal/config ./internal/security -v`.

## Task P3-T06: AI-safe Docs And Unsafe Example

**Status:** pending

- [ ] Create `docs/ai-safe-development.md` with safe and unsafe config examples, strict mode behavior, and redaction examples.
- [ ] Create `examples/unsafe-config/mockport.yml`.
- [ ] Run documented `mockport run --config examples/unsafe-config/mockport.yml --check` and confirm warning output.
- [ ] Run full verification: `/usr/local/go/bin/go test ./...`, `/usr/local/go/bin/go vet ./...`.

## Phase 3 Exit

- [ ] `ai-safe` warns but starts.
- [ ] `strict` fails before startup.
- [ ] Report marks safe/unsafe state.
- [ ] No full secret values appear in CLI output, logs, or report.
- [ ] AI-safe docs match implemented behavior.
