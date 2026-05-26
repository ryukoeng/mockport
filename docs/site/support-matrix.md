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
| `stripe` | `workflow-compatible` | checkout sessions, payment intents, customers, products, prices, subscriptions, invoices, refunds, fake webhook sender | success, failure, auth error, rate limit, timeout | Stripe SDK `22.1.1` contract passes for supported local workflows; fake state, list/retrieve, validation errors, and idempotency-key replay are available. Known gaps: no fraud, billing network, tax, disputes, Connect, or full Billing lifecycle. |
| `openai` | `experimental` | models, chat completions, responses | chat success, SSE stream success, rate limit, context length, auth error | Stateful chat/response IDs and response retrieval are available; `stream_success` returns Server-Sent Events for chat completions; no real model inference. |
| `github-oauth` | `experimental` | authorize, access token, user profile | oauth success, invalid code, expired token, missing scope | Stateful fake codes, tokens, scopes, and user identity are available; no GitHub org/enterprise policy. |
| `slack` | `experimental` | auth.test, chat.postMessage, conversations.history | message success, auth error, rate limit, delivery failure | Stateful fake messages can be posted and read through channel history; no real workspace delivery. |

## Planned Compatibility Track

Provider-compatible work follows the model in `docs/compatibility-model.md` and is staged after public preview:

1. Compatibility manifests and scores.
2. SDK contract harness foundation.
3. Shared state foundation.
4. Adapter state adoption.
5. Provider-specific compatibility tracks.
