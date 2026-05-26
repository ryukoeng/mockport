# Limitations

Mockport targets provider-compatible local APIs for selected workflows. It does not reproduce provider internals or undocumented behavior.

## Current Preview Scope

Current adapters are scenario-compatible or partial. They are useful for local and CI integration paths, but they are not a substitute for provider sandboxes or production validation.

## What Mockport Does Not Reproduce

- Real payment processing, fraud systems, settlement, or billing networks.
- Real AI inference, model quality, provider tokenization, or private scheduling behavior.
- Real GitHub organization, enterprise, or permission policy.
- Real Slack workspace delivery, enterprise policy, or full directory state.
- Undocumented provider behavior.

## How To Evaluate Support

Use:

- [Support matrix](support-matrix.md)
- `/_mockport/report`
- `mockport report --format json`
- Adapter examples
- Public env safety checks

Unsupported endpoints and approximations should be visible before adoption.
