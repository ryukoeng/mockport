# Public Env Safety

Mockport generated env files are intended to be safe to commit when they remain unchanged.

## Safe Values

Safe Mockport env values use:

- `mockport_` fake credential prefixes.
- `whsec_mockport` for fake webhook signing.
- `http://localhost:43101` or other local-only URLs.

Examples:

```env
STRIPE_API_URL=http://localhost:43101/stripe
STRIPE_SECRET_KEY=mockport_stripe_secret
STRIPE_WEBHOOK_SECRET=whsec_mockport
OPENAI_BASE_URL=http://localhost:43101/openai/v1
OPENAI_API_KEY=mockport_openai_key
SLACK_BOT_TOKEN=mockport_slack_token
```

## Unsafe Values

Do not commit:

- Real provider API keys or tokens.
- Production provider API URLs.
- Customer payloads.
- Captured webhook secrets.
- Ambiguous placeholders that may later be replaced with real credentials.

## Check

```bash
bash scripts/check-public-env.sh
```

This check scans public env examples and docs for real-looking provider credentials, production provider URLs, and ambiguous placeholders.
