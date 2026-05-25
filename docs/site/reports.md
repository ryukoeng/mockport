# Reports

Mockport exposes a run report at:

```bash
curl http://localhost:43101/_mockport/report
```

The report includes safety status, enabled adapters, request metadata, scenario coverage, behavior matrix entries, and unsupported endpoint attempts.

For CLI output:

```bash
mockport report --format text
mockport report --format json
```
