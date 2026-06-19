# Quickstart

[日本語版](quickstart.ja.md)

```bash
mockport init --adapter stripe
docker compose -f docker-compose.mockport.yml up
curl http://localhost:43101/health
mockport healthcheck
```

For multiple adapters:

```bash
mockport init --adapter stripe --adapter openai --adapter github-oauth --adapter slack --adapter line
docker compose -f docker-compose.mockport.yml up
```

## Switching scenarios

Besides fixing a scenario in `mockport.yml`, you can switch per request using the `X-Mockport-Scenario` header — no server restart required.

```bash
# Test the Stripe failure path without restarting the server
curl -X POST http://localhost:43101/stripe/v1/checkout/sessions \
  -H "X-Mockport-Scenario: payment_failed" \
  -H "Authorization: Bearer $STRIPE_KEY" \
  -d "mode=payment&success_url=http://localhost/success&cancel_url=http://localhost/cancel"
```

See the [adapter reference](adapters.md) for the list of supported scenarios per adapter.
