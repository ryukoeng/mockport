# Limitations

[日本語版](limitations.ja.md)

Mockport targets provider-compatible local APIs for selected workflows. It does not reproduce provider internals or undocumented behavior.

## Current Preview Scope

Current mainline adapters are workflow-compatible for selected local and CI integration paths. They are not a substitute for provider sandboxes or production validation.

## What Mockport Does Not Reproduce

- Real payment processing, fraud systems, settlement, or billing networks.
- Real AI inference, model quality, provider tokenization, or private scheduling behavior.
- Real GitHub organization, enterprise, or permission policy.
- Real Slack workspace delivery, enterprise policy, or full directory state.
- Real LINE Login UI, LIFF browser runtime, provider webhook redelivery, quota/rate-bucket enforcement, regional policy, or Dapp Portal behavior.
- Undocumented provider behavior.

## Unimplemented Configuration Blocks

The `scenarios:` block in `mockport.yml` is parsed but **not implemented** — it is silently
ignored at runtime. Mockport will emit a warning at startup (and in `--check` output and
`/_mockport/report`) when this block is present.

For response switching and error-case simulation, use:

- Built-in scenarios via the adapter's `scenario:` field in `mockport.yml`
- The `X-Mockport-Scenario` request header (see issue #80)

See [scenario-policy.md](../scenario-policy.md) for future plans on user-defined scenarios.

## How To Evaluate Support

Use:

- [Support matrix](support-matrix.md)
- `/_mockport/report`
- `mockport report --format json`
- Adapter examples
- Public env safety checks

Unsupported endpoints and approximations should be visible before adoption.
