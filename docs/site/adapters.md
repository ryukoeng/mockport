# Adapters

| Adapter | Base path | Maturity | Workflows |
| --- | --- | --- | --- |
| `stripe` | `/stripe` plus SDK-compatible `/v1` alias | `workflow-compatible` | checkout sessions, payment intents, customers, products, prices, subscriptions, invoices, refunds, fake signed webhooks |
| `openai` | `/openai` | `workflow-compatible` | models, chat completions, responses, embeddings, files, batches |
| `github-oauth` | `/github` | `workflow-compatible` | authorize redirect, token exchange, user profile, user emails, user orgs |
| `slack` | `/slack` | `experimental` | auth test, message posting |

Adapters are scenario-driven today and are moving toward provider-compatible local APIs for selected workflows. Use the [support matrix](support-matrix.md) and report behavior matrix to confirm supported paths.
