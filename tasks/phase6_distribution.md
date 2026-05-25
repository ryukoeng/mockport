# Phase 6 Distribution Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Mockport を OSS として配布できる状態に近づける。Docker image、release binaries、Homebrew、npm wrapper、docs site の導線を用意する。

**Architecture:** 配布は GitHub Actions と local scripts で再現可能にする。npm wrapper はこの Phase で scaffold するが、core runtime は引き続き Go binary / Docker image を正とする。

**Tech Stack:** GitHub Actions, Docker Buildx, GoReleaser or shell scripts, Homebrew formula template, npm wrapper scaffold, static docs.

---

## Files

- Create: `.github/workflows/release.yml`
- Create: `.github/workflows/docker.yml`
- Create: `scripts/build-release-archives.sh`
- Create: `scripts/test-release-archives.sh`
- Create: `packaging/homebrew/mockport.rb.template`
- Create: `packaging/npm/package.json`
- Create: `packaging/npm/bin/mockport.js`
- Create: `docs/site/index.md`
- Create: `docs/site/quickstart.md`
- Create: `docs/site/adapters.md`
- Modify: `README.md`

## Task P6-T01: Release Build Workflow

**Status:** pending

- [ ] Write static workflow test script that checks `.github/workflows/release.yml` includes linux/darwin, amd64/arm64, checksums, and artifact upload.
- [ ] Create release workflow.
- [ ] Run workflow static check locally.
- [ ] Keep signing/notarization out of scope unless credentials exist.

## Task P6-T02: GHCR Docker Image Workflow

**Status:** pending

- [ ] Write static workflow test asserting GHCR image name, `latest`, semver tag, and Dockerfile path.
- [ ] Create `.github/workflows/docker.yml`.
- [ ] Ensure workflow runs `go test ./...` before publish.
- [ ] Document required GitHub permissions.

## Task P6-T03: Release Archives And Checksums

**Status:** pending

- [ ] Write failing shell test for archive names: `mockport_<version>_<os>_<arch>.tar.gz`.
- [ ] Implement `scripts/build-release-archives.sh`.
- [ ] Implement checksum generation with `sha256sum` or `shasum -a 256`.
- [ ] Run `scripts/test-release-archives.sh`.

## Task P6-T04: Homebrew Formula Template

**Status:** pending

- [ ] Write template test asserting placeholders for version, URL, sha256, and binary install.
- [ ] Add `packaging/homebrew/mockport.rb.template`.
- [ ] Document manual formula update flow.
- [ ] Do not publish tap until release artifacts exist.

## Task P6-T05: npm Wrapper Scaffold

**Status:** pending

- [ ] Write Node test asserting `mockport --help` wrapper delegates to downloaded binary or Docker fallback.
- [ ] Add `packaging/npm/package.json` and `bin/mockport.js`.
- [ ] Keep npm wrapper marked experimental.
- [ ] Document that npm is a wrapper, not the primary runtime.

## Task P6-T06: Docs Site Scaffold

**Status:** pending

- [ ] Add docs site markdown pages for quickstart, adapters, AI-safe, reports, distribution.
- [ ] Add local docs link check script or simple markdown path check.
- [ ] Update README with docs site source and install options.
- [ ] Run docs check and full Go verification.

## Phase 6 Exit

- [ ] Release workflow is present and statically checked.
- [ ] Docker publish workflow is present and statically checked.
- [ ] Release archive script works locally.
- [ ] Homebrew formula template exists.
- [ ] npm wrapper scaffold exists and is clearly marked later/experimental.
- [ ] Docs site scaffold exists with quickstart, adapters, AI-safe, reports, distribution pages.
