# Compatibility Model

Mockport compatibility is measured against public provider API behavior, SDK contracts, selected workflows, fake state, and known gaps. It does not mean provider internals, undocumented behavior, model quality, fraud checks, billing networks, delivery guarantees, or production network effects are reproduced.

## Implementation Boundary

Mockport implements selected public API and SDK/client behavior that can be exercised deterministically in local tests. The implementation line is:

- In scope: public endpoint path/method/status/header/response shape, pinned SDK or client contract behavior, deterministic fake state, selected workflow lifecycles, common error envelopes, retry/rate-limit hints, sanitized fixtures, reports, and explicit known gaps.
- Out of scope: provider internals, undocumented behavior, external provider calls, real account policy, production delivery guarantees, fraud or risk engines, billing networks, settlement, real model quality, tokenization parity, hosted tools, enterprise enforcement, UI-level login flows, and provider console review workflows.

No adapter should add broad provider surface area only to make examples pass. New behavior should be added when it is part of a selected workflow, has official-reference grounding, and can be backed by tests, fixtures, and known-gap documentation.

## Levels

| Level | Meaning |
| --- | --- |
| `wire` | Request path, method, status, headers, and response shape are represented for selected endpoints. |
| `sdk` | Selected official SDK calls pass against local Mockport with pinned SDK versions. |
| `workflow` | A selected user workflow works across multiple requests. |
| `state` | Fake deterministic state supports create/retrieve/list/update or equivalent lifecycle paths. |
| `error` | Common provider error shapes, status codes, and retry/rate-limit hints are represented. |
| `contract` | Manifest, fixture, SDK, workflow, state, and known-gap evidence are all present for the selected surface. |

## Score Inputs

Compatibility score is deterministic and offline. It combines:

- Endpoint coverage.
- Built-in scenario coverage.
- SDK evidence.
- Fake state support.
- Error behavior support.

User-defined scenarios do not raise provider compatibility score unless promoted into a built-in scenario with tests, docs, and sanitized fixture evidence.

## SDK Contract Harness

SDK contract tests live under `contract/sdk`. This workspace is test-only and is intentionally separate from the Go runtime and the experimental npm wrapper.

The Phase 15 foundation runs a live placeholder contract against local Mockport health. Provider-specific tracks add real SDK calls later without contacting external provider APIs.

## Maturity Promotion

| Maturity | Minimum evidence |
| --- | --- |
| `experimental` | Adapter exists with explicit metadata and known gaps. |
| `partial` | Common scenario-compatible paths are implemented and reported. |
| `sdk-compatible` | SDK level evidence exists and selected SDK contracts pass. |
| `workflow-compatible` | Workflow, state, and error evidence exists for selected workflows. |
| `provider-compatible` | Contract level evidence exists with manifest, fixtures, SDK contracts, workflow/state/error coverage, score, and known gaps. |

Adapters must not be promoted only because local app-specific behavior works. Unsupported and approximate behavior must stay visible in reports.
