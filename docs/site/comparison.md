# Comparison

[日本語版](comparison.ja.md)

Mockport is a local, Docker-first emulator for selected SaaS workflows. It is not a drop-in replacement for every mock server, provider sandbox, or contract-test tool. This page compares Mockport to the alternatives developers actually evaluate.

## How to choose

- **Stripe only, no state needed, and you want broad OpenAPI-derived endpoint coverage** → [stripe-mock](https://github.com/stripe/stripe-mock) (Stripe's official mock server).
- **Run AWS services locally** → [LocalStack](https://localstack.cloud/) (Mockport does not overlap here — LocalStack targets AWS APIs, not SaaS providers like Stripe).
- **You already have an OpenAPI spec and want a quick mock** → [Stoplight Prism](https://github.com/stoplightio/prism).
- **Keep mocks inside JavaScript tests (browser or Node.js)** → [MSW](https://github.com/mswjs/msw).
- **Multiple SaaS providers (payments + LLM + chat + OAuth) in one process, with state, built-in scenarios, webhook delivery, and secret-safe defaults** → Mockport.
- **Fully custom HTTP stubs you own end-to-end** → WireMock or hand-written doubles (see below).
- **Authoritative provider behavior before production** → provider sandboxes (see below).

## Feature comparison

| | Mockport | stripe-mock | Prism | MSW | WireMock |
| --- | --- | --- | --- | --- | --- |
| Target | Stripe, OpenAI, Slack, GitHub OAuth, LINE, Zoho OAuth | Stripe only | Any API with an OpenAPI (or Postman) spec | Any HTTP target (JavaScript runtimes only) | Any HTTP API (hand-written stubs) |
| Stateful create → retrieve | ✅ | ❌ | ❌ | Depends on your handlers | Depends on your stubs |
| Built-in error/scenario switching | ✅ (`scenario:` field, `X-Mockport-Scenario`) | ❌ | Example/response selection only | Hand-written | Hand-written |
| Webhook / event delivery | ✅ (Stripe, LINE, Slack helpers) | ❌ | ❌ | — | ✅ (when configured) |
| Language-agnostic HTTP client | ✅ | ✅ | ✅ | ❌ (JavaScript only) | ✅ |
| OpenAPI-driven mock generation | ❌ | ✅ (from Stripe's OpenAPI) | ✅ | — | ❌ |
| Secret-safe public env policy | ✅ (`ai-safe` mode) | — | — | — | — |
| Docker-first single process | ✅ | ✅ | ✅ | ❌ (in-process library) | ✅ |

Cells marked `—` were not confirmed from primary sources for that tool in this pass.

### Source notes (primary references)

| Tool | References used for the table |
| --- | --- |
| stripe-mock | [stripe/stripe-mock README](https://github.com/stripe/stripe-mock/blob/master/README.md) — stateless server, OpenAPI-derived responses, no configurable error scenarios, no webhook sender |
| LocalStack | [LocalStack AWS services docs](https://docs.localstack.cloud/aws/services/) — AWS service emulation; SaaS APIs such as Stripe are out of scope |
| Prism | [stoplightio/prism README](https://github.com/stoplightio/prism/blob/master/README.md) — OpenAPI/Postman mock server; dynamic responses from spec examples; no built-in persistent state |
| MSW | [mswjs/msw README](https://github.com/mswjs/msw/blob/main/README.md) — request interception in browser and Node.js; JavaScript/TypeScript runtimes |
| WireMock | [WireMock stubbing docs](https://wiremock.org/docs/stubbing/), [webhooks and callbacks](https://wiremock.org/docs/webhooks-and-callbacks/) — generic stub definitions; webhook delivery is opt-in configuration |

## When Mockport is not the best fit

Be direct about scope:

- **You need every Stripe endpoint** — Mockport supports selected workflows only. See [support matrix](support-matrix.md) and adapter specs; prefer stripe-mock for OpenAPI-wide Stripe coverage.
- **You need provider billing math, fraud logic, or real authorization rules** — Mockport does not reproduce provider internals ([ROADMAP](../../ROADMAP.md) Non-Goals).
- **You need OpenAPI spec → mock with zero adapter code** — use Prism (or generate stubs yourself); Mockport ships curated adapters, not a generic spec compiler.
- **Your tests are JavaScript-only and should stay in-process** — MSW is simpler than running a sidecar HTTP server.
- **You need authoritative provider behavior** — use the real provider sandbox for final validation; Mockport is for fast local and CI integration tests.
- **You need AWS service emulation** — use LocalStack; Mockport does not emulate AWS APIs.

## Mockport vs WireMock

WireMock is a general HTTP mocking tool: you define stubs, matchers, and optional webhook callbacks yourself. Mockport is provider-shaped — adapters encode selected service workflows, fake credentials, safety warnings, reports, and support matrices so teams share the same local provider surface without rewriting stubs per project.

## Mockport vs Hand-written Test Doubles

Hand-written doubles are quick for one codebase but drift quickly. Mockport centralizes adapter behavior, examples, reports, and compatibility evidence so multiple apps and CI jobs can share the same local provider API.

## Mockport vs Provider Sandboxes

Provider sandboxes are authoritative for provider behavior. Mockport is local, Docker-first, secret-free, and deterministic. Use Mockport for fast local and CI integration tests, then validate critical paths against the real provider sandbox before production.

## Positioning summary

Mockport is best when:

- You want Docker-first local APIs for multiple SaaS providers in one process.
- You want fake env values that are safe to commit.
- You want deterministic external-service scenarios in CI with explicit unsupported behavior.
- You want adapter coverage and known gaps reported before adoption.

Mockport is not a full clone of provider internals.
