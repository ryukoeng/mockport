# Changelog

## v0.1.0-alpha - 2026-05-26

Initial public preview release.

### Included

- Docker-first Mockport runtime with AI-safe configuration checks.
- Stripe-like payment adapter for checkout sessions, payment intents, webhook sending, and common error scenarios.
- Experimental OpenAI-compatible, GitHub OAuth-like, and Slack-like adapters.
- `/_mockport/report` for request history, scenario coverage, behavior matrix, and safety findings.
- GitHub Release archives for Linux and macOS on amd64 and arm64.
- GHCR image published as `ghcr.io/albert-einshutoin/mockport:0.1.0-alpha`.

### Known Limits

- This is scenario-compatible, not full provider-compatible.
- Provider SDK contract coverage starts in later phases.
- Homebrew and npm are not published distribution channels yet.

