# Phase 22 Provider-compatible Release Track Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 高忠実度互換を release train として運用し、provider/API/SDK version ごとの compatibility を継続的に公開する。

**Architecture:** Compatibility score、SDK contract、fixture coverage、known gaps を release notes と docs に出す。互換性は promise ではなく測定可能な artifact として扱い、provider updates に追従する maintenance loop を作る。

**Tech Stack:** GitHub Actions scheduled jobs, release docs, compatibility manifests, generated reports.

---

## Files

- Create: `.github/workflows/compatibility.yml`
- Create: `docs/compatibility-reports/README.md`
- Create: `scripts/generate-compatibility-report.sh`
- Create: `scripts/check-compatibility-release.sh`
- Modify: `docs/site/support-matrix.md`
- Modify: `CHANGELOG.md`
- Modify: `tasks/status.md`

## Task P22-T01: Compatibility CI

**Status:** pending

- [ ] Add scheduled and manual compatibility workflow.
- [ ] Run SDK contracts and fixture checks in the workflow.
- [ ] Upload compatibility report artifact.
- [ ] Run static workflow checks locally.

## Task P22-T02: Generated Compatibility Reports

**Status:** pending

- [ ] Generate markdown/JSON reports from compatibility manifests.
- [ ] Include adapter score, tested SDK versions, provider API version, known gaps.
- [ ] Add docs index for compatibility reports.
- [ ] Run report generation check.

## Task P22-T03: Provider-compatible Release Criteria

**Status:** pending

- [ ] Define release labels: `experimental`, `sdk-compatible`, `workflow-compatible`, `provider-compatible`.
- [ ] Require SDK contract pass and minimum score before promoting adapter maturity.
- [ ] Add release checklist to changelog/release docs.
- [ ] Run full compatibility release check.

## Phase 22 Exit

- [ ] Compatibility CI runs on demand and schedule.
- [ ] Compatibility reports are generated from manifests.
- [ ] Adapter maturity is promoted only with contract evidence.
- [ ] Release notes show compatibility scores and known gaps.
