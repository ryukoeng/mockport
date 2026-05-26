# Mockport

Secret-free service emulation for AI-native development.

Mockport runs local Docker-based emulators for external services. The Minimal MVP supports a Stripe-like payment adapter so local development, CI, and AI coding workflows can test payment integration paths without real Stripe secrets.

For public preview scope, support matrix, examples, and limitations, see [Mockport Docs](docs/site/index.md).

## Quickstart

No local install required:

```bash
docker build -t mockport:local -f docker/Dockerfile .
docker run --rm -p 43101:43101 \
  -v $(pwd)/configs/mockport.example.yml:/etc/mockport/mockport.yml \
  mockport:local
```

Then verify the local API:

```bash
curl http://localhost:43101/health
curl -X POST http://localhost:43101/stripe/v1/checkout/sessions
curl http://localhost:43101/_mockport/report
```

CLI workflow from a built binary:

```bash
mockport init --adapter stripe
docker compose -f docker-compose.mockport.yml up
```

`mockport init` protects existing generated files by default. Use `--force` only when you intentionally want to replace `mockport.yml`, `.env.mockport`, and `docker-compose.mockport.yml`.

Application `.env`:

```env
STRIPE_API_URL=http://localhost:43101/stripe
STRIPE_SECRET_KEY=mockport_stripe_secret
STRIPE_WEBHOOK_SECRET=whsec_mockport
```

This Mockport env is safe to commit when the generated fake values remain unchanged. See [Public Env Safety](docs/public-env-safety.md).

Test a Stripe-like checkout session:

```bash
curl -X POST http://localhost:43101/stripe/v1/checkout/sessions
```

View the request and safety report:

```bash
mockport report
```

## Source Smoke Test

From a checkout of this repository:

```bash
bash scripts/smoke-empty-dir.sh
```

The smoke test builds the local Docker image, creates a temporary empty directory, runs `mockport init --adapter stripe`, starts Docker Compose, checks `/health`, posts a Stripe-like checkout request, and prints `mockport report`.

## Docker

```bash
docker run --rm -p 43101:43101 \
  -v $(pwd)/configs/mockport.example.yml:/etc/mockport/mockport.yml \
  mockport:local
```

## Install And Distribution

Mockport is Docker-first. Release binaries and packaging scaffolds are prepared for OSS distribution.

| Channel | Status | Notes |
| --- | --- | --- |
| Docker / GHCR | Planned | `ghcr.io/albert-einshutoin/mockport` workflow publishes semver tags and `latest` |
| GitHub release archives | Planned | `mockport_<version>_<os>_<arch>.tar.gz` with `checksums.txt` |
| Homebrew | Template | Formula template is under `packaging/homebrew/` |
| npm | Experimental wrapper | The npm wrapper is experimental; Go binary and Docker remain primary |

Docs site source lives under `docs/site/`.

## Services

Supported:

| Service | Adapter | Base path | Maturity | Supported workflows |
| --- | --- | --- | --- | --- |
| Stripe-like payments | `stripe` | `/stripe` | `partial` | checkout sessions, payment intents, fake signed webhooks, success/failure/auth/rate-limit/timeout scenarios |
| OpenAI-compatible API | `openai` | `/openai` | `experimental` | models, chat completions, responses, auth error, rate limit, context length error |
| GitHub OAuth-like API | `github-oauth` | `/github` | `experimental` | authorize redirect, access token exchange, user profile, invalid code, expired token, missing scope |
| Slack-like messaging API | `slack` | `/slack` | `experimental` | auth test, message posting, auth error, rate limit, delivery failure |

Planned:

| Service | Planned adapter | Target workflows | Status |
| --- | --- | --- | --- |
| LINE Messaging-like API | `line` | message push/reply, webhook signature, delivery failure, rate limit | Later candidate |
| SendGrid-like email API | `sendgrid` | email send success/failure, auth error, rate limit, webhook events | Later candidate |

## AI-safe By Default

Mockport warns when real-looking credentials or real external service URLs are detected. In `strict` mode, unsafe configuration fails startup.

Check a config without starting the server:

```bash
mockport run --config examples/unsafe-config/mockport.yml --check
```

See [AI-safe Development](docs/ai-safe-development.md) for warning categories, strict mode, and redaction behavior.

Mockport is not a full clone of external services. It focuses on local integration testing scenarios: success, failure, auth error, rate limit, timeout, and webhook/callback.
