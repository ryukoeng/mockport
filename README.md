# Mockport

Secret-free service emulation for AI-native development.

Mockport runs local Docker-based emulators for external services. The Minimal MVP supports a Stripe-like payment adapter so local development, CI, and AI coding workflows can test payment integration paths without real Stripe secrets.

## Quickstart

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

## Supported Adapters

MVP:

- Stripe-like payment adapter

Later:

- OpenAI-compatible API
- GitHub OAuth
- LINE Messaging or Slack
- SendGrid

## AI-safe By Default

Mockport warns when real-looking credentials or real external service URLs are detected. In `strict` mode, unsafe configuration fails startup.

Mockport is not a full clone of external services. It focuses on local integration testing scenarios: success, failure, auth error, rate limit, timeout, and webhook/callback.
