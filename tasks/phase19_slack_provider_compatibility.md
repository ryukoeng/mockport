# Phase 19 Slack Provider Compatibility Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Slack-like adapter を Slack client/workflow-compatible local API へ引き上げる。

**Architecture:** Slack messaging workflow、conversation state、message lifecycle、basic events を fake workspace 上で再現する。Slack の enterprise policy、real delivery、workspace directory 全体は再現しない。

**Tech Stack:** Go Slack adapter, Slack SDK/client contract tests, sanitized fixtures, compatibility manifests.

---

## Files

- Create: `compat/fixtures/slack/*`
- Create: `contract/sdk/slack-smoke.test.js`
- Modify: `adapters/slack/*`
- Modify: `docs/site/support-matrix.md`
- Modify: `tasks/status.md`

## Task P19-T01: Slack Client Contract Baseline

**Status:** pending

- [ ] Write failing client smoke for auth.test.
- [ ] Write failing client smoke for chat.postMessage.
- [ ] Write failing client smoke for conversations list/history if client workflow needs it.
- [ ] Run `bash scripts/run-sdk-contracts.sh slack`.

## Task P19-T02: Slack Messaging And Conversation State

**Status:** pending

- [ ] Persist channels, users, messages, timestamps, and bot identity.
- [ ] Add message update/delete and conversation history subset.
- [ ] Add event callback subset with signature verification where useful.
- [ ] Update support matrix and compatibility score.

## Task P19-T03: Slack Error And Rate Limit Fidelity

**Status:** pending

- [ ] Add fixtures for invalid_auth, channel_not_found, not_in_channel, rate_limited, delivery_failed.
- [ ] Implement response shape and retry headers close to Slack public API.
- [ ] Report unsupported methods and scopes.
- [ ] Run adapter tests, client contracts, and compatibility report.

## Phase 19 Exit

- [ ] Slack adapter is at least `workflow-compatible`.
- [ ] Slack client contract passes for supported workflows.
- [ ] Messaging state and errors are documented by fixtures.
- [ ] Known gaps are explicit.
