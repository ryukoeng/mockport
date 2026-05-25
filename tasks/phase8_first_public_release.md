# Phase 8 First Public Release Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Mockport の最初の public release を GitHub Release と GHCR image として実際に公開する。

**Architecture:** Phase 6 の scaffold を dry-run から実 publish 証跡へ進める。タグ、release notes、checksums、Docker image digest、install verification を repository に記録する。

**Tech Stack:** GitHub Actions, GHCR, Git tags, shell verification, Docker.

---

## Files

- Create: `CHANGELOG.md`
- Create: `docs/releases/v0.1.0.md`
- Create: `scripts/verify-release-artifacts.sh`
- Modify: `README.md`
- Modify: `docs/site/distribution.md`
- Modify: `tasks/status.md`

## Task P8-T01: Release Readiness Audit

**Status:** pending

- [ ] Write `scripts/verify-release-artifacts.sh` to verify release archive names, checksums, binary version, and GHCR image availability.
- [ ] Add `CHANGELOG.md` with `v0.1.0` initial release notes.
- [ ] Add `docs/releases/v0.1.0.md` with scope, supported adapters, known limitations, verification commands.
- [ ] Run full pre-release verification locally.

## Task P8-T02: Tag And GitHub Release

**Status:** pending

- [ ] Create annotated tag `v0.1.0`.
- [ ] Push tag and confirm release workflow completes on GitHub.
- [ ] Confirm GitHub Release has four archives and `checksums.txt`.
- [ ] Run release archive install test from downloaded artifact.

## Task P8-T03: GHCR Publish Verification

**Status:** pending

- [ ] Confirm GHCR image `ghcr.io/albert-einshutoin/mockport:0.1.0` exists.
- [ ] Confirm `latest` points to the same release generation.
- [ ] Run Docker pull and smoke test from GHCR image, not local build.
- [ ] Record image digest in `docs/releases/v0.1.0.md`.

## Task P8-T04: Release Docs Update

**Status:** pending

- [ ] Update README install instructions with actual release URLs.
- [ ] Update docs site distribution page with release artifacts and GHCR commands.
- [ ] Keep Homebrew and npm marked not-yet-published unless actually published.
- [ ] Run README install path audit from a temporary directory.

## Phase 8 Exit

- [ ] `v0.1.0` tag exists on GitHub.
- [ ] GitHub Release contains archives and checksums.
- [ ] GHCR image can be pulled and run.
- [ ] README install commands work from a clean environment.
- [ ] Known limitations are documented in release notes.
