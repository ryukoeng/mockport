# Phase 10 Community And Maintenance Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 公開OSSとして継続メンテできる運用面を整え、外部 contributor と adopter の負担を下げる。

**Architecture:** Repository automation、dependency update、compatibility policy、roadmap、triage labels を用意する。機能追加より maintenance loop を作る。

**Tech Stack:** GitHub Actions, Dependabot, Markdown, GitHub labels/issues.

---

## Files

- Create: `.github/dependabot.yml`
- Create: `ROADMAP.md`
- Create: `docs/compatibility-policy.md`
- Create: `docs/maintainer-guide.md`
- Create: `scripts/check-maintenance-policy.sh`
- Modify: `README.md`
- Modify: `tasks/status.md`

## Task P10-T01: Maintenance Policy

**Status:** pending

- [ ] Add compatibility policy for adapter maturity levels and breaking changes.
- [ ] Add maintainer guide for release process, issue triage, and security reports.
- [ ] Add roadmap with near-term adapters and non-goals.
- [ ] Run maintenance policy check.

## Task P10-T02: Dependency And CI Maintenance

**Status:** pending

- [ ] Add Dependabot config for GitHub Actions, Go modules, and npm wrapper.
- [ ] Add scheduled CI smoke for Docker path if feasible.
- [ ] Add stale-free issue triage guidance without auto-closing early users.
- [ ] Run CI static checks locally.

## Task P10-T03: Contribution Quality Bar

**Status:** pending

- [ ] Document adapter acceptance criteria: metadata, scenarios, examples, report coverage, AI-safe behavior.
- [ ] Add checklist for new adapter PRs.
- [ ] Add support policy for experimental vs partial adapters.
- [ ] Run public trust and docs checks.

## Phase 10 Exit

- [ ] Compatibility and maintenance policy are documented.
- [ ] Dependabot is configured.
- [ ] Adapter contribution quality bar is explicit.
- [ ] Public roadmap exists and avoids overpromising.
