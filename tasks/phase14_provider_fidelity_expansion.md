# Phase 14 Provider Fidelity Expansion Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 主要 provider の公開 API surface を selected workflows から major endpoint-compatible へ広げる。

**Architecture:** AI 仕様駆動で provider docs/OpenAPI/SDK fixtures から coverage backlog を作り、endpoint ごとに contract fixture、TDD、compatibility report を追加する。Full-compatible は全 endpoint 全 edge case ではなく、公開された主要 API 面の高忠実度互換として扱う。

**Tech Stack:** Go adapters, fixture snapshots, generated coverage manifests, SDK contract harness, compatibility scoring.

---

## Files

- Create: `compat/fixtures/stripe/*`
- Create: `compat/fixtures/openai/*`
- Create: `compat/fixtures/slack/*`
- Create: `compat/fixtures/github/*`
- Create: `scripts/check-compat-fixtures.sh`
- Modify: `docs/site/support-matrix.md`
- Modify: `tasks/status.md`

## Task P14-T01: Fixture And Spec Snapshot Policy

**Status:** pending

- [ ] Define sanitized fixture format for request, response, headers, status, provider version, source note.
- [ ] Add fixture checker that rejects real-looking secrets and external URLs.
- [ ] Document how AI-generated endpoint implementations must cite docs/spec/fixture source.
- [ ] Run `bash scripts/check-compat-fixtures.sh`.

## Task P14-T02: Stripe Major Surface Expansion

**Status:** pending

- [ ] Add customers, prices, products, subscriptions, invoices, refunds, and webhooks coverage backlog.
- [ ] Implement one endpoint group at a time with failing tests first.
- [ ] Add SDK contract coverage for each endpoint group.
- [ ] Update support matrix and compatibility score.

## Task P14-T03: OpenAI Major Surface Expansion

**Status:** pending

- [ ] Add embeddings, files, batches, assistants/responses tool-call subset, model errors, and streaming edge fixtures.
- [ ] Implement endpoint groups with TDD and SDK contracts.
- [ ] Keep model inference fake and deterministic.
- [ ] Update support matrix and compatibility score.

## Task P14-T04: GitHub/Slack Major Surface Expansion

**Status:** pending

- [ ] Add GitHub OAuth scopes, user emails, org membership subset, and token errors.
- [ ] Add Slack conversations, message update/delete, user lookup, and event callback subset.
- [ ] Implement endpoint groups with TDD and client contracts.
- [ ] Update support matrix and compatibility score.

## Phase 14 Exit

- [ ] Major provider surfaces have coverage backlogs.
- [ ] Fixture policy prevents leaking real secrets.
- [ ] Endpoint groups are added with TDD and contract evidence.
- [ ] Support matrix is generated or updated from compatibility manifests.
