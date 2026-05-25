# Slack Message Example

This example runs the Slack-like Mockport adapter with fake local credentials.

```bash
docker build -t mockport:local -f docker/Dockerfile .
mockport run --config examples/slack-message/mockport.yml
```

Use these values in the application under test:

```env
SLACK_API_URL=http://localhost:43101/slack/api
SLACK_BOT_TOKEN=mockport_slack_token
```

Smoke test:

```bash
curl -X POST http://localhost:43101/slack/api/auth.test
curl -X POST http://localhost:43101/slack/api/chat.postMessage
curl http://localhost:43101/_mockport/report
```
