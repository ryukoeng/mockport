# Roadmap

Mockport is a Docker-first local API environment for AI-native development and CI. This roadmap is intentionally scoped to public preview work and provider-compatible direction without promising full provider internals.

## Current Release

- `v0.1.0-alpha`: public preview with Docker/GHCR, GitHub release archives, AI-safe env policy, and scenario-compatible adapters.

## Near Term

- Phase 12: sanitized fixture, provider spec snapshot, and scenario policy.
- Phase 13: public preview contract cleanup for `mockport up`, OpenAI streaming, and adapter helper boundaries.
- Phase 14: compatibility level model, compatibility report shape, and provisional promotion rules.
- Phase 15: SDK contract harness foundation.
- Phase 16: shared state foundation before provider-specific compatibility tracks.

## Public Preview Follow-up

- [#6](https://github.com/albert-einshutoin/mockport/issues/6): Add SSE-compatible streaming response for OpenAI `stream_success`.
- [#8](https://github.com/albert-einshutoin/mockport/issues/8): Improve `mockport up` Docker Compose UX with clearer errors and `--detach` / `--build`.
- [#5](https://github.com/albert-einshutoin/mockport/issues/5): Clarify built-in scenario policy versus user-defined scenarios.
- [#7](https://github.com/albert-einshutoin/mockport/issues/7): Track adapter helper duplication before adding more adapters.

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
