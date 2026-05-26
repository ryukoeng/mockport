# Phase 12 Fixture, Spec, And Scenario Policy Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** compatibility engine と SDK contract 実装の前に、fixture/spec snapshot/scenario の安全性、出典、更新ルールを定義する。

**Architecture:** AI 仕様駆動で provider docs/OpenAPI/SDK fixtures を使う前提として、sanitized fixture format、source metadata、secret leakage check、provider version tracking、built-in scenario と user-defined scenario の境界を用意する。これにより後続 Phase の互換実装を測定可能かつ安全にする。

**Tech Stack:** JSON fixtures, shell/static checks, Go security scanner, docs.

---

## Files

- Create: `compat/fixtures/README.md`
- Create: `compat/fixtures/schema.example.json`
- Create: `scripts/check-compat-fixtures.sh`
- Create: `docs/fixture-policy.md`
- Create: `docs/scenario-policy.md`
- Modify: `internal/security/secrets.go`
- Modify: `internal/security/secrets_test.go`
- Modify: `tasks/status.md`

## Task P12-T01: Fixture Format And Scenario Policy

**Status:** pending

- [ ] Define sanitized fixture format for request, response, headers, status, provider version, SDK version, and source note.
- [ ] Add example fixture with fake local credentials only.
- [ ] Document which fields must be redacted or normalized.
- [ ] Clarify built-in scenario policy versus user-defined scenarios before scoring scenario coverage. See [issue #5](https://github.com/albert-einshutoin/mockport/issues/5).
- [ ] Run fixture format check.

## Task P12-T02: Fixture Safety Check

**Status:** pending

- [ ] Write failing scanner tests for real-looking secrets inside fixture files.
- [ ] Add fixture checker that rejects real-looking secrets, production provider URLs, and missing source metadata.
- [ ] Add CI-ready shell command.
- [ ] Run `/usr/local/go/bin/go test ./internal/security -v` and `bash scripts/check-compat-fixtures.sh`.

## Task P12-T03: Spec Snapshot Policy

**Status:** pending

- [ ] Document how AI-generated endpoint implementations must cite docs/spec/fixture source.
- [ ] Define update policy when provider docs or SDK versions change.
- [ ] Define when a fixture is enough and when SDK contract evidence is required.
- [ ] Run docs/static checks.

## Phase 12 Exit

- [ ] Fixture format is documented.
- [ ] Fixture checker prevents secret leakage.
- [ ] Provider docs/spec/SDK source policy is explicit.
- [ ] Built-in scenarios and user-defined scenarios have an explicit policy.
- [ ] Later fidelity work has a safe evidence base.
