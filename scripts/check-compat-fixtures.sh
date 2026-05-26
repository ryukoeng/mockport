#!/usr/bin/env bash
set -euo pipefail

require_file() {
  local path="$1"
  if [[ ! -f "$path" ]]; then
    echo "missing file: $path" >&2
    exit 1
  fi
}

require_text() {
  local path="$1"
  local text="$2"
  if ! grep -Fq "$text" "$path"; then
    echo "missing text in $path: $text" >&2
    exit 1
  fi
}

require_file "compat/fixtures/README.md"
require_file "compat/fixtures/schema.example.json"
require_file "docs/fixture-policy.md"
require_file "docs/scenario-policy.md"

require_text "compat/fixtures/README.md" "provider_version"
require_text "compat/fixtures/README.md" "source.retrieved_at"
require_text "compat/fixtures/README.md" "User-defined scenarios"
require_text "docs/fixture-policy.md" "SDK contract evidence is required"
require_text "docs/scenario-policy.md" "Built-in scenarios"
require_text "docs/scenario-policy.md" "User-defined scenarios"

go test ./internal/security -run 'Fixture|SchemaExample|CompatibilityFixtureFiles' -v
