# Phase 5 Additional Adapters Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Phase 4 の adapter metadata/report contract に沿って、OpenAI-compatible、GitHub OAuth-like、Slack-like adapter を built-in adapter として追加する。

**Architecture:** 既存の built-in registry 方針を維持し、各 adapter は Phase 4 の metadata contract を実装する。各 adapter は common success path、auth error、rate limit、timeout 相当の代表シナリオを持ち、report に対応範囲が表示される。

**Tech Stack:** Go 1.26.3, `net/http`, `httptest`, JSON fixtures, adapter metadata contract.

---

## Files

- Create: `adapters/openai/*`
- Create: `adapters/githuboauth/*`
- Create: `adapters/slack/*`
- Modify: `internal/cli/init.go`
- Modify: `internal/cli/init_test.go`
- Create: `internal/cli/add.go`
- Create: `internal/cli/add_test.go`
- Create: `examples/openai-chat/*`
- Create: `examples/github-oauth/*`
- Create: `examples/slack-message/*`
- Modify: `README.md`
- Modify: `docs/reporting.md`

## Task P5-T01: OpenAI-compatible Adapter

**Status:** pending

- [ ] Write failing HTTP tests for `GET /openai/v1/models`, `POST /openai/v1/chat/completions`, `POST /openai/v1/responses`.
- [ ] Add scenario tests for `chat_success`, `stream_success`, `rate_limited`, `context_length_exceeded`, `auth_error`.
- [ ] Implement `adapters/openai` with deterministic JSON responses.
- [ ] Implement Phase 4 metadata contract for OpenAI.
- [ ] Register OpenAI in `mockport run`.
- [ ] Run `/usr/local/go/bin/go test ./adapters/openai ./internal/server -v`.

## Task P5-T02: GitHub OAuth-like Adapter

**Status:** pending

- [ ] Write failing tests for `GET /github/login/oauth/authorize`, `POST /github/login/oauth/access_token`, `GET /github/user`.
- [ ] Add scenarios: `oauth_success`, `invalid_code`, `expired_token`, `scope_missing`.
- [ ] Implement redirects/token/user JSON with fake deterministic values.
- [ ] Implement Phase 4 metadata contract for GitHub OAuth.
- [ ] Register GitHub OAuth in server.
- [ ] Run `/usr/local/go/bin/go test ./adapters/githuboauth ./internal/server -v`.

## Task P5-T03: Slack-like Adapter

**Status:** pending

- [ ] Write failing tests for `POST /slack/api/auth.test`, `POST /slack/api/chat.postMessage`.
- [ ] Add scenarios: `message_success`, `auth_error`, `rate_limited`, `delivery_failed`.
- [ ] Implement Slack-like JSON bodies with `ok`, `error`, `channel`, `ts`.
- [ ] Implement Phase 4 metadata contract for Slack.
- [ ] Register Slack in server.
- [ ] Run `/usr/local/go/bin/go test ./adapters/slack ./internal/server -v`.

## Task P5-T04: Multi-adapter Init And Add

**Status:** pending

- [ ] Write failing CLI tests: `mockport init --adapter stripe --adapter openai` generates both adapter configs and env vars.
- [ ] Write failing tests for `mockport add openai github-oauth slack`.
- [ ] Implement repeated adapter flags and `add` command for config updates.
- [ ] Preserve existing files unless `--force` is passed.
- [ ] Run `/usr/local/go/bin/go test ./internal/cli -run 'Init|Add' -v`.

## Task P5-T05: Adapter Examples

**Status:** pending

- [ ] Add one example directory per new adapter.
- [ ] Add example config load tests or smoke script entries for each config.
- [ ] Update README supported adapters and adapter maturity table.
- [ ] Run `/usr/local/go/bin/go test ./...`.

## Task P5-T06: Cross-adapter Smoke Coverage

**Status:** pending

- [ ] Create a multi-adapter config enabling Stripe, OpenAI, GitHub OAuth, and Slack.
- [ ] Add smoke script that starts Docker and calls one success endpoint per adapter.
- [ ] Assert report contains all adapters and their metadata.
- [ ] Run Docker smoke and full verification.

## Phase 5 Exit

- [ ] OpenAI, GitHub OAuth, and Slack-like adapters are registered.
- [ ] Each adapter has success, auth error, rate limit or equivalent failure scenario tests.
- [ ] Each adapter implements metadata/report contract.
- [ ] Examples exist and config loads.
- [ ] Multi-adapter init/add works.
- [ ] Cross-adapter Docker smoke passes.
