# Phase 18 GitHub OAuth Provider Compatibility Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** GitHub OAuth-like adapter を OAuth client/workflow-compatible local API へ引き上げる。

**Architecture:** GitHub OAuth の authorize/token/user profile workflow を fake identity/state 上で再現する。GitHub org/enterprise の完全権限モデルや内部 policy は再現しない。

**Tech Stack:** Go GitHub OAuth adapter, OAuth client contract tests, sanitized fixtures, compatibility manifests.

---

## Files

- Create: `compat/fixtures/github/*`
- Create: `contract/sdk/github-oauth-smoke.test.js`
- Modify: `adapters/githuboauth/*`
- Modify: `docs/site/support-matrix.md`
- Modify: `tasks/status.md`

## Task P18-T01: GitHub OAuth Client Contract Baseline

**Status:** pending

- [ ] Write failing client smoke for authorize redirect with state.
- [ ] Write failing client smoke for token exchange.
- [ ] Write failing client smoke for user profile and user emails subset.
- [ ] Run `bash scripts/run-sdk-contracts.sh github-oauth`.

## Task P18-T02: GitHub OAuth State And Scope Fidelity

**Status:** pending

- [ ] Persist auth codes, access tokens, scopes, expiry, and fake user identities.
- [ ] Add invalid code, expired token, missing scope, redirect URI mismatch fixtures.
- [ ] Add user emails and org membership subset if required by common OAuth workflows.
- [ ] Update support matrix and compatibility score.

## Task P18-T03: GitHub OAuth Error Fidelity

**Status:** pending

- [ ] Match public OAuth error response shapes for token endpoint.
- [ ] Match common API auth error shapes for user endpoints.
- [ ] Report unsupported scopes and endpoints.
- [ ] Run adapter tests, client contracts, and compatibility report.

## Phase 18 Exit

- [ ] GitHub OAuth adapter is at least `workflow-compatible`.
- [ ] OAuth client contract passes for supported workflows.
- [ ] Scope and token state behavior are documented.
- [ ] Known gaps are explicit.
