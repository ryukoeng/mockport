# Contributing

## Setup

Use Go 1.26.3.

```bash
/usr/local/go/bin/go test ./...
/usr/local/go/bin/go vet ./...
/usr/local/go/bin/go build ./cmd/mockport
```

## TDD

Production changes must follow TDD:

1. Write the failing test.
2. Run the narrow test and confirm it fails for the expected reason.
3. Implement the smallest change that passes.
4. Run the narrow test again.
5. Run the full verification for the touched Phase.

## Public Trust Checks

Run these before opening a pull request:

```bash
bash scripts/check-public-trust.sh
bash scripts/check-distribution.sh
/usr/local/go/bin/go test ./...
/usr/local/go/bin/go vet ./...
```

## Adapter Changes

Adapter pull requests must include:

- HTTP tests for success and failure scenarios.
- Metadata/report coverage.
- Example config or docs when user-facing.
- AI-safe behavior for fake credentials and local base URLs.
- Clear unsupported behavior in reports or docs.

## Pull Requests

Include test evidence in the PR body. Do not include real provider secrets, production URLs, customer payloads, or unsanitized fixtures.
