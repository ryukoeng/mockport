# Phase 11 Community And Maintenance Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 公開OSSとして継続メンテできる運用面を整え、外部 contributor と adopter の負担を下げる。

**Architecture:** Repository automation、dependency update、roadmap、triage labels を用意する。Compatibility policy は Phase 12 の compatibility model 後に詳細化するため、この Phase では一般的な maintenance policy に絞る。

**Tech Stack:** GitHub Actions, Dependabot, Markdown, GitHub labels/issues.

---

## Files

- Create: `.github/dependabot.yml`
- Create: `ROADMAP.md`
- Create: `docs/maintainer-guide.md`
- Create: `scripts/check-maintenance-policy.sh`
- Modify: `README.md`
- Modify: `tasks/status.md`

## Task P11-T01: Maintenance Policy

**Status:** pending

- [ ] Add maintainer guide for release process, issue triage, and security reports.
- [ ] Add roadmap with near-term adapters, provider-compatible direction, and non-goals.
- [ ] Document no stale auto-close policy for early external users.
- [ ] Run maintenance policy check.

## Task P11-T02: Dependency And CI Maintenance

**Status:** pending

- [ ] Add Dependabot config for GitHub Actions, Go modules, and npm wrapper.
- [ ] Add scheduled CI smoke for Docker path if feasible.
- [ ] Document that test-only SDK dependencies are pinned in the SDK contract harness phase.
- [ ] Run CI static checks locally.

## Task P11-T03: Contribution Quality Bar

**Status:** pending

- [ ] Document adapter acceptance criteria: metadata, scenarios, examples, report coverage, AI-safe behavior.
- [ ] Add checklist for new adapter PRs.
- [ ] Reference future compatibility levels without defining final thresholds yet.
- [ ] Run public trust and docs checks.

## Phase 11 Exit

- [ ] Maintenance policy is documented.
- [ ] Dependabot is configured.
- [ ] Adapter contribution quality bar is explicit.
- [ ] Public roadmap exists and avoids overpromising.
