# Maintainer Guide

This guide keeps the public OSS workflow explicit while Mockport is in preview.

## Release Process

1. Confirm `tasks/status.md` has the current phase ready for release.
2. Run local verification:

```bash
bash scripts/check-public-trust.sh
bash scripts/check-public-env.sh
bash scripts/check-distribution.sh
bash scripts/check-maintenance-policy.sh
/usr/local/go/bin/go test ./...
/usr/local/go/bin/go vet ./...
```

3. Build and verify release archives:

```bash
tmpdir="$(mktemp -d)"
scripts/build-release-archives.sh 0.1.0-alpha "$tmpdir"
scripts/verify-release-artifacts.sh 0.1.0-alpha "$tmpdir"
```

4. Create an annotated tag and let GitHub Actions publish release archives and GHCR images.
5. Download the GitHub Release artifacts and run `scripts/verify-release-artifacts.sh`.
6. Pull the GHCR image and run Docker smoke from the published image.
7. Record release URL, image tag, digest, known limitations, and verification evidence in `docs/releases/`.

## Issue Triage

- Do not auto-close stale issues while Mockport is in early public preview.
- Ask for adapter name, Mockport version or commit, sanitized config, expected behavior, actual behavior, and test evidence.
- Remove or redact any real provider secret, production URL, customer payload, or unsanitized fixture.
- Treat safety issues involving leaked credentials or unsafe examples as security-sensitive until ruled out.

## Security Reports

Follow `SECURITY.md`. Do not request real secrets or production provider traffic to reproduce a report.

## Dependency And CI Maintenance

- Dependabot covers GitHub Actions, Go modules, and the experimental npm wrapper.
- GitHub Actions should use Node.js 24-compatible action releases where available.
- If a Node.js 20 deprecation warning remains because an upstream action has no Node.js 24-compatible release, document the action, warning date, and upstream tracking reason in the relevant phase notes.
- Scheduled smoke checks should validate the Docker-first path without publishing images.

Test-only SDK dependencies are intentionally pinned later in Phase 15 when the SDK contract harness exists.

## Adapter Contribution Quality Bar

Adapter pull requests must include:

- Adapter metadata: name, base path, maturity, capabilities, endpoints, and scenarios.
- HTTP tests for success, auth error, rate limit, and at least one meaningful failure scenario.
- Example config and user-facing docs when the adapter is exposed publicly.
- Report coverage for supported and unsupported behavior.
- AI-safe fake credentials and local-only URLs.
- No real provider secrets, production URLs, customer payloads, or unsanitized fixtures.

Keep provider-shaped helpers local until `docs/adapter-helper-policy.md` says the abstraction threshold has been met.

Future compatibility levels will be defined in Phase 14. Until then, adapter claims must stay at scenario-compatible or experimental unless stronger evidence exists.
