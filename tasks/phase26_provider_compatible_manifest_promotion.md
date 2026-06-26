# Phase 26 Provider-compatible Manifest Promotion Implementation Plan

[日本語版](phase26_provider_compatible_manifest_promotion.ja.md)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Versioned compatibility manifests under `compat/manifests/` mirror runtime `Metadata()` claims so maturity changes are reviewable diffs and CI blocks drift or silent promotion.

**Architecture:** `scripts/gen-compat-manifests` writes JSON from `compat.FromMetadata()` for the five published adapters. `scripts/check-compat-manifests.sh` diffs checked-in files against regeneration output. `scripts/validate-compatibility-report.mjs` requires report maturity to match the checked-in manifest maturity.

**Tech Stack:** Go generator tool, shell drift check, Node.js report validator, existing `internal/compat` model.

---

## Files

- Create: `internal/builtins/builtins.go`
- Create: `scripts/gen-compat-manifests/main.go`
- Create: `scripts/check-compat-manifests.sh`
- Create: `compat/manifests/stripe.json`
- Create: `compat/manifests/openai.json`
- Create: `compat/manifests/github-oauth.json`
- Create: `compat/manifests/slack.json`
- Create: `compat/manifests/line.json`
- Modify: `internal/cli/builtin.go`
- Modify: `scripts/check-compatibility-release.sh`
- Modify: `scripts/validate-compatibility-report.mjs`
- Modify: `.github/workflows/ci.yml`
- Modify: `docs/compatibility-model.md`
- Modify: `tasks/status.md`

## Task P26-T01: Manifest Generator And Drift Check

**Status:** done

- [x] Add `internal/builtins` with shared `Adapters()` and `ManifestAdapters()` lists.
- [x] Add `scripts/gen-compat-manifests` with `--out` flag and deterministic JSON output.
- [x] Add `scripts/check-compat-manifests.sh` to diff checked-in manifests against regeneration output.
- [x] Wire manifest drift check into CI and `scripts/check-compatibility-release.sh`.

## Task P26-T02: Initial Manifests

**Status:** done

- [x] Generate and commit `stripe.json`, `openai.json`, `github-oauth.json`, `slack.json`, and `line.json` from current `Metadata()`.

## Task P26-T03: Promotion Gate

**Status:** done

- [x] Extend `scripts/validate-compatibility-report.mjs` to require report maturity to match checked-in manifest maturity.
- [x] Block maturity increases when `promotion_eligible` is false relative to the checked-in manifest baseline.
- [x] Document versioned manifest workflow in `docs/compatibility-model.md`.

## Phase 26 Exit

- [x] Versioned manifests exist for every published compatibility-report adapter (`stripe`, `openai`, `github-oauth`, `slack`, `line`).
- [x] Release checks fail on manifest drift or report/manifest maturity mismatch.
- [x] `provider-compatible` actual promotion remains out of scope (#92); the gate makes future promotion reviewable.
