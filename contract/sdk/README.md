# Mockport SDK Contract Harness

This workspace is test-only. It keeps future official SDK contract tests outside the Go runtime and outside the experimental npm wrapper.

## Commands

```bash
npm test
```

Runs the offline harness sanity check.

```bash
MOCKPORT_BASE_URL=http://127.0.0.1:43101 npm run test:live -- --provider all --json
```

Runs the live placeholder contract against a local Mockport server. Provider-specific tracks add real SDK tests here later.

Supported provider selectors:

- `all`
- `stripe`
- `openai`
- `github-oauth`
- `slack`

No test in this workspace may call an external provider API.
