# Support Matrix

Mockport support is explicit and scenario-driven. Use reports to confirm what a test run exercised.

## Maturity

| Maturity | Meaning |
| --- | --- |
| `experimental` | Early adapter coverage for selected workflows. Expect gaps. |
| `partial` | Common workflows are implemented with documented unsupported behavior. |
| `sdk-compatible` | Selected official SDK calls pass against local Mockport. |
| `workflow-compatible` | Selected workflows include fake state, errors, and replayable behavior. |
| `provider-compatible` | Selected provider workflows are backed by manifests, SDK contracts, fixtures, and known-gap reports. |

## Adapter Coverage

| Adapter | Maturity | Endpoints | Scenarios | Notes |
| --- | --- | --- | --- | --- |
| `stripe` | `partial` | checkout sessions, payment intents, fake webhook sender | success, failure, auth error, rate limit, timeout | Scenario-compatible; not yet SDK/workflow-compatible. |
| `openai` | `experimental` | models, chat completions, responses | chat success, SSE stream success, rate limit, context length, auth error | Deterministic fake responses; `stream_success` returns Server-Sent Events for chat completions; no real model inference. |
| `github-oauth` | `experimental` | authorize, access token, user profile | oauth success, invalid code, expired token, missing scope | Fake OAuth identity; no GitHub org/enterprise policy. |
| `slack` | `experimental` | auth.test, chat.postMessage | message success, auth error, rate limit, delivery failure | Fake messaging; no real workspace delivery. |

## Planned Compatibility Track

Provider-compatible work is staged after public preview:

1. Compatibility manifests and scores.
2. SDK contract harness foundation.
3. Shared state foundation.
4. Adapter state adoption.
5. Provider-specific compatibility tracks.
