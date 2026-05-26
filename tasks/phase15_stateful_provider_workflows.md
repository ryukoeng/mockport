# Phase 15 Stateful Provider Workflows Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 主要 provider workflow を fake state 上で成立させ、作成した resource を retrieve/list/update できる Docker-first API 環境にする。

**Architecture:** adapter ごとに in-memory deterministic store を導入し、必要なら future phase で file-backed state に拡張する。内部 provider logic は再現せず、公開 API に見える resource lifecycle と error semantics を再現する。

**Tech Stack:** Go 1.26.3, in-memory stores, deterministic IDs, httptest, SDK contract tests.

---

## Files

- Create: `internal/state/store.go`
- Create: `internal/state/store_test.go`
- Modify: `adapters/stripe/*`
- Modify: `adapters/openai/*`
- Modify: `adapters/githuboauth/*`
- Modify: `adapters/slack/*`
- Modify: `internal/report/report.go`
- Modify: `tasks/status.md`

## Task P15-T01: Deterministic State Store

**Status:** pending

- [ ] Write failing tests for create/retrieve/list/update/delete semantics.
- [ ] Implement concurrency-safe in-memory store with deterministic IDs per adapter.
- [ ] Add reset support for test isolation.
- [ ] Run `/usr/local/go/bin/go test ./internal/state -v`.

## Task P15-T02: Stripe Stateful Workflows

**Status:** pending

- [ ] Make checkout session and payment intent create/retrieve/list stateful.
- [ ] Add idempotency key handling for create endpoints.
- [ ] Add validation error shapes close to Stripe for missing required fields.
- [ ] Run adapter tests and Stripe SDK contract.

## Task P15-T03: OpenAI Stateful/Conversation Workflows

**Status:** pending

- [ ] Persist response IDs and chat completion IDs for retrieve-compatible fake workflows where applicable.
- [ ] Add streaming chunks with SDK-compatible shape.
- [ ] Add validation for model, messages/input, and unsupported parameters.
- [ ] Run adapter tests and OpenAI SDK contract.

## Task P15-T04: OAuth And Messaging State

**Status:** pending

- [ ] Persist GitHub OAuth codes, tokens, scopes, and fake user identities.
- [ ] Persist Slack messages by channel and timestamp.
- [ ] Add list/retrieve helper endpoints where provider workflows expect them.
- [ ] Run Slack/GitHub client contracts.

## Phase 15 Exit

- [ ] Major workflows are stateful across create/retrieve/list paths.
- [ ] Idempotency and validation errors are represented for supported workflows.
- [ ] SDK contract tests exercise stateful behavior.
- [ ] Report identifies stateful workflow coverage.
