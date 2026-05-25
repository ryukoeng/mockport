# Multi-adapter Example

This example runs Stripe-like, OpenAI-compatible, GitHub OAuth-like, and Slack-like adapters in one Mockport process.

```bash
docker build -t mockport:local -f docker/Dockerfile .
docker compose -f examples/multi-adapter/docker-compose.yml up
```

Smoke test:

```bash
curl -X POST http://localhost:43101/stripe/v1/checkout/sessions
curl http://localhost:43101/openai/v1/models
curl http://localhost:43101/github/user
curl -X POST http://localhost:43101/slack/api/auth.test
curl http://localhost:43101/_mockport/report
```
