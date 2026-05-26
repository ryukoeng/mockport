# Phase 10 Public Docs And Discovery Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 公開後の初見ユーザーが、価値、制約、導入方法、adapter coverage を短時間で判断できる状態にする。

**Architecture:** docs/site を単なる markdown 群から公開可能な information architecture に整える。adapter ごとの support matrix と known limitations を明示し、README は短く保つ。文言は provider-compatible 目標と矛盾しないようにする。

**Tech Stack:** Markdown, static docs, link checks, shell checks.

---

## Files

- Create: `docs/site/support-matrix.md`
- Create: `docs/site/examples.md`
- Create: `docs/site/limitations.md`
- Create: `docs/site/comparison.md`
- Create: `scripts/check-doc-links.sh`
- Modify: `docs/site/index.md`
- Modify: `README.md`
- Modify: `tasks/status.md`

## Task P10-T01: Public Docs Information Architecture

**Status:** pending

- [ ] Add docs index sections for install, quickstart, adapters, support matrix, examples, limitations, reports, distribution.
- [ ] Add `support-matrix.md` with endpoint/scenario/maturity table for each adapter.
- [ ] Add `limitations.md` with: "Mockport targets provider-compatible local APIs for selected workflows. It does not reproduce provider internals or undocumented behavior."
- [ ] Run docs path check.

## Task P10-T02: Example-driven Onboarding

**Status:** pending

- [ ] Add `examples.md` linking Stripe, OpenAI, GitHub OAuth, Slack, and multi-adapter examples.
- [ ] Add copy-paste commands for each example.
- [ ] Verify each example config loads with Go tests or shell check.
- [ ] Run smoke script for multi-adapter path.

## Task P10-T03: Public Positioning

**Status:** pending

- [ ] Add `comparison.md` explaining Mockport vs full provider sandbox, WireMock, and hand-written test doubles.
- [ ] Keep claims factual and tied to implemented behavior.
- [ ] Update README intro to point readers to docs instead of becoming too long.
- [ ] Run markdown/link checks.

## Phase 10 Exit

- [ ] Public docs explain install, examples, support matrix, limitations, reports, and distribution.
- [ ] Adapter limitations are visible before adoption.
- [ ] Example commands are tested or covered by smoke checks.
- [ ] README is concise and points to docs.
