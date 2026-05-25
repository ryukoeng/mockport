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

Mockport flags provider secrets such as:

```txt
sk_live_
sk_test_
AKIA
ASIA
ghp_
github_pat_
xoxb-
xoxp-
AIza
whsec_
```

Mockport also flags live provider URLs such as:

```txt
https://api.stripe.com
https://api.openai.com
https://api.github.com
https://api.line.me
https://slack.com/api
```

## Report

`/_mockport/report` includes a safety summary:

```json
{
  "mode": "ai-safe",
  "safety": {
    "mode": "ai-safe",
    "safe": false,
    "real_looking_secrets": 1,
    "external_urls": 1
  }
}
```

No full secret value should appear in CLI output or report output.
