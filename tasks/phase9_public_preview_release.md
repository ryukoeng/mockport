# Phase 9 Public Preview Release Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Mockport の最初の public preview release を GitHub Release と GHCR image として公開する。

**Architecture:** Phase 6 の scaffold を dry-run から実 publish 証跡へ進める。ただし provider-compatible release とは分け、`v0.1.0-alpha` または `v0.1.0-preview` として current scenario-compatible scope と known gaps を明示する。

**Tech Stack:** GitHub Actions, GHCR, Git tags, shell verification, Docker.

---

## Files

- Create: `CHANGELOG.md`
- Create: `docs/releases/v0.1.0-alpha.md`
- Create: `scripts/verify-release-artifacts.sh`
- Modify: `README.md`
- Modify: `docs/site/distribution.md`
- Modify: `tasks/status.md`

## Task P9-T01: Preview Release Readiness Audit

**Status:** pending

- [ ] Write `scripts/verify-release-artifacts.sh` to verify release archive names, checksums, binary version, and GHCR image availability.
- [ ] Add `CHANGELOG.md` with `v0.1.0-alpha` initial preview release notes.
- [ ] Add `docs/releases/v0.1.0-alpha.md` with scope, supported adapters, known limitations, public env safety, and verification commands.
- [ ] Run full pre-release verification locally.

## Task P9-T02: Tag And GitHub Preview Release

**Status:** pending

- [ ] Create annotated tag `v0.1.0-alpha`.
- [ ] Push tag and confirm release workflow completes on GitHub.
- [ ] Confirm GitHub Release has four archives and `checksums.txt`.
- [ ] Run release archive install test from downloaded artifact.

## Task P9-T03: GHCR Preview Publish Verification

**Status:** pending

- [ ] Confirm GHCR image `ghcr.io/albert-einshutoin/mockport:0.1.0-alpha` exists.
- [ ] Confirm `latest` behavior is intentional for preview releases before enabling it.
- [ ] Run Docker pull and smoke test from GHCR image, not local build.
- [ ] Record image digest in `docs/releases/v0.1.0-alpha.md`.

## Task P9-T04: Preview Install Docs Update

**Status:** pending

- [ ] Update README install instructions with actual preview release URLs.
- [ ] Update docs site distribution page with release artifacts and GHCR commands.
- [ ] Keep Homebrew and npm marked not-yet-published unless actually published.
- [ ] Run README install path audit from a temporary directory.

## Phase 9 Exit

- [ ] Preview tag exists on GitHub.
- [ ] GitHub Release contains archives and checksums.
- [ ] GHCR preview image can be pulled and run.
- [ ] README install commands work from a clean environment.
- [ ] Known limitations and non-provider-compatible status are documented.
