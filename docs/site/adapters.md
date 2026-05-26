# Adapters

| Adapter | Base path | Maturity | Workflows |
| --- | --- | --- | --- |
| `stripe` | `/stripe` | `partial` | checkout sessions, payment intents, fake signed webhooks |
| `openai` | `/openai` | `experimental` | models, chat completions, responses |
| `github-oauth` | `/github` | `experimental` | authorize redirect, token exchange, user profile |
| `slack` | `/slack` | `experimental` | auth test, message posting |

Adapters are scenario-driven today and are moving toward provider-compatible local APIs for selected workflows. Use the [support matrix](support-matrix.md) and report behavior matrix to confirm supported paths.
