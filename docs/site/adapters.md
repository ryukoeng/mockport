# Adapters

| Adapter | Base path | Maturity | Workflows |
| --- | --- | --- | --- |
| `stripe` | `/stripe` | `partial` | checkout sessions, payment intents, fake signed webhooks |
| `openai` | `/openai` | `experimental` | models, chat completions, responses |
| `github-oauth` | `/github` | `experimental` | authorize redirect, token exchange, user profile |
| `slack` | `/slack` | `experimental` | auth test, message posting |

Adapters are scenario-driven and intentionally incomplete. Use the report behavior matrix to confirm supported paths.
