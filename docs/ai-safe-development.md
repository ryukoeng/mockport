# AI-safe Development

Mockport does not try to hide real secrets from AI tools. It makes real secrets unnecessary for local and CI integration testing.

## Default Mode

`ai-safe` is the default mode:

```yaml
mode: ai-safe
```

In `ai-safe` mode, Mockport warns when configuration contains real-looking provider secrets or live provider URLs.

```bash
mockport run --config examples/unsafe-config/mockport.yml --check
```

Expected output includes warning categories and field names, but not full secret values or live URLs:

```txt
[MOCKPORT SECURITY WARNING]
- stripe.fake_secret: real-looking secret detected (real_looking_secret)
- stripe.api_url: external live service URL detected (external_url)
Config check passed
```

## Strict Mode

`strict` mode fails before the server starts if unsafe fields are detected:

```yaml
mode: strict
```

Use strict mode in CI when a real-looking secret or live provider URL should break the job.

## Safe Fake Values

Use fake local values:

```env
STRIPE_API_URL=http://localhost:43101/stripe
STRIPE_SECRET_KEY=mockport_stripe_secret
STRIPE_WEBHOOK_SECRET=whsec_mockport
```

Mockport treats these prefixes as local fake values:

```txt
mockport_
local_
fake_
dummy_
```

## Unsafe Examples

Mockport flags real-looking provider secrets, including common live/test API key prefixes, cloud access key prefixes, GitHub token prefixes, Slack token prefixes, Google API key prefixes, and non-Mockport webhook signing secret prefixes.

Mockport also flags live provider URLs for supported and planned providers.

## Report

`/_mockport/report` includes a safety summary:

```json
{
  "mode": "ai-safe",
  "safety": {
    "mode": "ai-safe",
    "safe": false,
    "real_looking_secrets": 1,
    "external_urls": 1,
    "public_env_safe": false
  }
}
```

No full secret value should appear in CLI output or report output.
