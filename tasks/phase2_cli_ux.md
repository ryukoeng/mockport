# Phase 2 CLI UX Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 空ディレクトリから `mockport init` と `mockport up/run` で 2 分以内に Stripe adapter を起動できる CLI UX にする。

**Architecture:** 既存の `internal/cli` を拡張し、生成ファイル保護、`--force`、`up`、`report` を追加する。Docker 実行は `docker compose -f docker-compose.mockport.yml up` の薄い wrapper に留め、複雑な Docker SDK は入れない。

**Tech Stack:** Go 1.26.3, Cobra, `os/exec`, `httptest`, temp directories.

---

## Files

- Modify: `internal/cli/init.go`
- Modify: `internal/cli/init_test.go`
- Create: `internal/cli/up.go`
- Create: `internal/cli/up_test.go`
- Create: `internal/cli/report.go`
- Create: `internal/cli/report_test.go`
- Modify: `internal/cli/root.go`
- Create: `scripts/smoke-empty-dir.sh`
- Modify: `README.md`
- Modify: `tasks/status.md`

## Task P2-T01: Protect Generated Files From Accidental Overwrite

**Status:** pending

- [ ] Write failing test: in `internal/cli/init_test.go`, create `mockport.yml` with custom content, run `mockport init --adapter stripe`, assert command returns error and file content is unchanged.
- [ ] Run: `/usr/local/go/bin/go test ./internal/cli -run Init -v`; expected failure because overwrite protection is missing.
- [ ] Implement: make `init` check `mockport.yml`, `.env.mockport`, `docker-compose.mockport.yml` before writing.
- [ ] Run: `/usr/local/go/bin/go test ./internal/cli -run Init -v`; expected PASS.
- [ ] Update `tasks/status.md`: P2-T01 `done`.

## Task P2-T02: Add `mockport init --force`

**Status:** pending

- [ ] Write failing test: existing generated files are overwritten when `mockport init --adapter stripe --force` is used.
- [ ] Run targeted CLI test and confirm failure.
- [ ] Implement `--force` bool flag in `internal/cli/init.go`.
- [ ] Verify generated files contain deterministic Stripe config after force.
- [ ] Run `/usr/local/go/bin/go test ./internal/cli -run Init -v`.

## Task P2-T03: Add `mockport up`

**Status:** pending

- [ ] Write failing test in `internal/cli/up_test.go`: inject command runner, run `mockport up`, assert command is `docker compose -f docker-compose.mockport.yml up`.
- [ ] Verify RED: `up` command is missing.
- [ ] Implement `internal/cli/up.go` with a small runner abstraction for tests.
- [ ] Register `up` in `root.go`.
- [ ] Run `/usr/local/go/bin/go test ./internal/cli -run Up -v`.

## Task P2-T04: Add `mockport report --url`

**Status:** pending

- [ ] Write failing test in `internal/cli/report_test.go`: `httptest.Server` returns report JSON, CLI prints `Mockport Report`, mode, adapter, request, safety counts.
- [ ] Verify RED: `report` command is missing.
- [ ] Implement `internal/cli/report.go` with default URL `http://localhost:43101/_mockport/report`.
- [ ] Register `report` in `root.go`.
- [ ] Run `/usr/local/go/bin/go test ./internal/cli -run Report -v`.

## Task P2-T05: Add User-facing Success Output

**Status:** pending

- [ ] Write failing tests asserting `mockport init` output includes exact next commands: `docker compose -f docker-compose.mockport.yml up` and `curl http://localhost:43101/health`.
- [ ] Implement concise next-step output.
- [ ] Run `/usr/local/go/bin/go test ./internal/cli -run Init -v`.
- [ ] Update README quickstart to match output.

## Task P2-T06: Empty-directory E2E Smoke

**Status:** pending

- [ ] Create `scripts/smoke-empty-dir.sh` that builds `mockport:local`, runs binary `init --adapter stripe`, starts Docker compose, curls `/health`, curls checkout, curls report, then cleans up.
- [ ] Add README section documenting the smoke command.
- [ ] Run: `bash scripts/smoke-empty-dir.sh`.
- [ ] Run full verification: `/usr/local/go/bin/go test ./...`, `/usr/local/go/bin/go vet ./...`, `docker build -t mockport:local -f docker/Dockerfile .`.

## Phase 2 Exit

- [ ] Empty directory flow completes in under 2 minutes.
- [ ] Existing files are protected unless `--force` is used.
- [ ] `mockport up` delegates to Docker Compose.
- [ ] `mockport report --url` renders report output.
- [ ] README commands match implementation.
- [ ] `tasks/status.md` Phase 2 is `done`.
