# Roadmap

Mockport is a Docker-first local API environment for AI-native development and CI. This roadmap is intentionally scoped to public preview work and provider-compatible direction without promising full provider internals.

## Current Release

- `v0.1.0-alpha`: public preview with Docker/GHCR, GitHub release archives, AI-safe env policy, and scenario-compatible adapters.

## Near Term

- Phase 11: community and maintenance automation.
- Phase 12: compatibility level model and compatibility report shape.
- Phase 13: sanitized fixture and provider spec snapshot policy.
- Phase 14: SDK contract harness foundation.

## Adapter Direction

Current adapters:

- Stripe-like payments.
- OpenAI-compatible API.
- GitHub OAuth-like API.
- Slack-like messaging API.

Candidate adapters after the compatibility foundation:

- LINE Messaging-like API.
- SendGrid-like email API.

## Compatibility Direction

Mockport aims for provider-compatible local APIs for selected workflows. Compatibility is measured by documented endpoint behavior, SDK contract tests, fake state, error shape, and reportable gaps.

Mockport does not reproduce provider internal logic, undocumented behavior, or production network effects.

## Non-Goals

- Proxying real provider traffic.
- Accepting real provider secrets in public examples.
- Claiming full provider compatibility before SDK and workflow contract evidence exists.
- Publishing npm or Homebrew as primary channels before Docker and Go binary release paths are stable.
