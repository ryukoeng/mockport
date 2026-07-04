# Reports

[日本語版](reports.ja.md)

Mockport exposes a run report at:

```bash
curl http://localhost:43101/_mockport/report
```

The report includes safety status, enabled adapters, request metadata, scenario coverage, behavior matrix entries, and unsupported endpoint attempts.

## Request history

Request history keeps metadata for the most recent 1000 requests recorded during a run. When that limit is exceeded, older entries are pruned from the front so the report always returns the newest requests in chronological order. The same bounded history feeds `unsupportedEndpoints` in the report payload.

For CLI output:

```bash
mockport report --format text
mockport report --format json
```
