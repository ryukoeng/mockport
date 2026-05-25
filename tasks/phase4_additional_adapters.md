# Phase 4 Additional Adapters Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Stripe に加えて OpenAI-compatible、GitHub OAuth-like、Slack-like adapter を built-in adapter として追加する。

**Architecture:** 既存の built-in registry 方針を維持し、adapter interface は必要最小限だけ拡張する。各 adapter は common success path、auth error、rate limit、timeout 相当の代表シナリオを持つ。

**Tech Stack:** Go 1.26.3, `net/http`, `httptest`, JSON fixtures.

---

## Files

- Modify: `internal/adapter/adapter.go`
- Modify: `internal/adapter/registry.go`
- Create: `adapters/openai/*`
- Create: `adapters/githuboauth/*`
- Create: `adapters/slack/*`
- Modify: `internal/cli/init.go`
- Modify: `internal/cli/init_test.go`
- Create: `examples/openai-chat/*`
- Create: `examples/github-oauth/*`
- Create: `examples/slack-message/*`
- Modify: `README.md`

## Task P4-T01: Adapter Scenario Metadata

**Status:** pending

- [ ] Write failing adapter tests asserting each adapter exposes name, maturity, capabilities, and scenarios.
- [ ] Extend adapter types with optional metadata without breaking existing Stripe behavior.
- [ ] Add Stripe metadata: maturity `partial`, scenarios `payment_success`, `payment_failed`, `auth_error`, `rate_limited`, `timeout`.
- [ ] Run `/usr/local/go/bin/go test ./internal/adapter ./adapters/stripe -v`.

## Task P4-T02: OpenAI-compatible Adapter

**Status:** pending

- [ ] Write failing HTTP tests for `GET /openai/v1/models`, `POST /openai/v1/chat/completions`, `POST /openai/v1/responses`.
- [ ] Add scenario tests for `chat_success`, `rate_limited`, `context_length_exceeded`, `auth_error`.
- [ ] Implement `adapters/openai` with deterministic JSON responses.
- [ ] Register OpenAI in `mockport run`.
- [ ] Run `/usr/local/go/bin/go test ./adapters/openai ./internal/server -v`.

## Task P4-T03: GitHub OAuth-like Adapter

**Status:** pending

- [ ] Write failing tests for `GET /github/login/oauth/authorize`, `POST /github/login/oauth/access_token`, `GET /github/user`.
- [ ] Add scenarios: `oauth_success`, `invalid_code`, `expired_token`, `scope_missing`.
- [ ] Implement redirects/token/user JSON with fake deterministic values.
- [ ] Register GitHub OAuth in server.
- [ ] Run `/usr/local/go/bin/go test ./adapters/githuboauth ./internal/server -v`.

## Task P4-T04: Slack-like Adapter

**Status:** pending

- [ ] Write failing tests for `POST /slack/api/auth.test`, `POST /slack/api/chat.postMessage`.
- [ ] Add scenarios: `message_success`, `auth_error`, `rate_limited`, `delivery_failed`.
- [ ] Implement Slack-like JSON bodies with `ok`, `error`, `channel`, `ts`.
- [ ] Register Slack in server.
- [ ] Run `/usr/local/go/bin/go test ./adapters/slack ./internal/server -v`.

## Task P4-T05: Multi-adapter Init And Add

**Status:** pending

- [ ] Write failing CLI tests: `mockport init --adapter stripe --adapter openai` generates both adapter configs and env vars.
- [ ] Write failing tests for `mockport add openai github-oauth slack`.
- [ ] Implement repeated adapter flags and `add` command for config updates.
- [ ] Preserve existing files unless `--force` is passed.
- [ ] Run `/usr/local/go/bin/go test ./internal/cli -run 'Init|Add' -v`.

## Task P4-T06: Adapter Examples

**Status:** pending

- [ ] Add one example directory per new adapter.
- [ ] Add example config load tests or smoke script entries for each config.
- [ ] Update README supported adapters.
- [ ] Run `/usr/local/go/bin/go test ./...` and Docker smoke for enabled multi-adapter config.

## Phase 4 Exit

- [ ] OpenAI, GitHub OAuth, and Slack-like adapters are registered.
- [ ] Each adapter has success, auth error, rate limit or equivalent failure scenario tests.
- [ ] Examples exist and config loads.
- [ ] Multi-adapter init/add works.
