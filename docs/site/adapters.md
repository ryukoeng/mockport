# Adapters

| Adapter | Base path | Maturity | Workflows |
| --- | --- | --- | --- |
| `stripe` | `/stripe` plus SDK-compatible `/v1` alias | `workflow-compatible` | checkout sessions, payment intents, customers, products, prices, subscriptions, invoices, refunds, fake signed webhooks |
| `openai` | `/openai` | `workflow-compatible` | models, chat completions, responses, embeddings, files, batches |
| `github-oauth` | `/github` | `workflow-compatible` | authorize redirect, token exchange, user profile, user emails, user orgs |
| `slack` | `/slack` | `workflow-compatible` | auth test, conversations list/history, message post/update/delete, Events API URL verification/message callback subset |
| `line` | `/line` | `workflow-compatible` | Messaging API send/content/signed webhook/rich menu/channel token workflows, LINE Login OAuth/profile, LIFF helpers, MINI App service messages, LINE Pay request/confirm, Mini Dapp wallet/payment helpers |

Adapters are scenario-driven today and are moving toward provider-compatible local APIs for selected workflows. Use the [support matrix](support-matrix.md) and report behavior matrix to confirm supported paths.

Detailed adapter specifications:

- [LINE adapter](../adapters/line.md)
