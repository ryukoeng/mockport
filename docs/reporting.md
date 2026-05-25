# Reporting

Mockport reports are trust artifacts. They explain what ran, what was safe, what is supported, and what was not supported.

## JSON Report

```bash
mockport report --format json
```

Important fields:

- `safety`: AI-safe summary, including real-looking secret and external URL counts.
- `adapters`: enabled adapters, base paths, capabilities, and maturity.
- `requests`: replay-safe request metadata. Request bodies and secret headers are not stored by default.
- `scenario_coverage`: supported scenarios per adapter.
- `behavior_matrix`: supported endpoints and their scenarios.
- `unsupported_endpoints`: requests that returned unsupported endpoint classifications.

## Text Report

```bash
mockport report --format text
```

Text output is meant for local development and CI logs. JSON output is better for tools.

## Maturity Levels

Allowed adapter maturity values:

```txt
experimental
partial
common-path
contract-tested
sandbox-verified
```

Stripe compatibility is `partial`. OpenAI-compatible, GitHub OAuth-like, and Slack-like adapters start as `experimental`.

All adapters are scenario-driven, not full provider clones. Use the behavior matrix and unsupported endpoint list as the source of truth for what a test run actually exercised.
