# Phase 12 SDK Contract Harness Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 公式 SDK または実運用に近い client が Mockport に接続して動くことを adapter ごとに検証する。

**Architecture:** npm/Rust を primary runtime にはしないが、SDK contract harness では provider 公式 SDK を test client として使う。CI では network-free local Mockport に接続し、base URL と fake env だけで主要 workflow が通ることを証明する。

**Tech Stack:** Go server, Node SDK smoke where provider SDK is Node-first, shell scripts, Docker Compose, fixture snapshots.

---

## Files

- Create: `contract/sdk/README.md`
- Create: `contract/sdk/package.json`
- Create: `contract/sdk/stripe-smoke.test.js`
- Create: `contract/sdk/openai-smoke.test.js`
- Create: `contract/sdk/slack-smoke.test.js`
- Create: `contract/sdk/github-oauth-smoke.test.js`
- Create: `scripts/run-sdk-contracts.sh`
- Modify: `.github/workflows/ci.yml`
- Modify: `docs/compatibility-model.md`
- Modify: `tasks/status.md`

## Task P12-T01: SDK Contract Workspace

**Status:** pending

- [ ] Add a dedicated SDK contract package outside runtime code.
- [ ] Add scripts that start Mockport locally and run SDK smoke tests against `localhost:43101`.
- [ ] Keep SDK dependencies test-only and out of the Go runtime.
- [ ] Run `(cd contract/sdk && npm test)` after expected initial failure.

## Task P12-T02: Stripe SDK Contract

**Status:** pending

- [ ] Write failing SDK smoke for checkout session create/retrieve and payment intent create/retrieve.
- [ ] Adjust Stripe adapter request parsing/response shape until official SDK accepts responses.
- [ ] Add Stripe SDK version to compatibility manifest.
- [ ] Run `bash scripts/run-sdk-contracts.sh stripe`.

## Task P12-T03: OpenAI SDK Contract

**Status:** pending

- [ ] Write failing SDK smoke for models list, chat completions, responses, and streaming if SDK supports override cleanly.
- [ ] Adjust response and streaming shape until SDK accepts it.
- [ ] Add OpenAI SDK version to compatibility manifest.
- [ ] Run `bash scripts/run-sdk-contracts.sh openai`.

## Task P12-T04: Slack And GitHub OAuth Client Contracts

**Status:** pending

- [ ] Add Slack SDK or raw compatible client smoke for auth.test and chat.postMessage.
- [ ] Add GitHub OAuth client smoke for authorize redirect, token exchange, and user profile.
- [ ] Add SDK/client versions to compatibility manifest.
- [ ] Run `bash scripts/run-sdk-contracts.sh`.

## Phase 12 Exit

- [ ] Official SDK/client smoke tests pass against local Docker-first Mockport.
- [ ] SDK versions are recorded.
- [ ] Runtime remains Go/Docker-first.
- [ ] Compatibility report reflects SDK contract results.
