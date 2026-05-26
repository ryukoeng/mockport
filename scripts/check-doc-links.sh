#!/usr/bin/env bash
set -euo pipefail

required_pages=(
  docs/site/index.md
  docs/site/quickstart.md
  docs/site/adapters.md
  docs/site/support-matrix.md
  docs/site/examples.md
  docs/site/limitations.md
  docs/site/comparison.md
  docs/site/ai-safe.md
  docs/site/reports.md
  docs/site/distribution.md
)

for page in "${required_pages[@]}"; do
  if [[ ! -f "$page" ]]; then
    echo "missing docs page: $page" >&2
    exit 1
  fi
done

require_link() {
  local file="$1"
  local link="$2"
  if ! grep -Fq "($link)" "$file"; then
    echo "missing link in $file: $link" >&2
    exit 1
  fi
  local target
  target="$(dirname "$file")/$link"
  if [[ ! -f "$target" ]]; then
    echo "broken link in $file: $link" >&2
    exit 1
  fi
}

for link in quickstart.md adapters.md support-matrix.md examples.md limitations.md comparison.md ai-safe.md reports.md distribution.md; do
  require_link docs/site/index.md "$link"
done

for example in stripe-checkout openai-chat github-oauth slack-message multi-adapter; do
  if ! grep -Fq "../examples/${example}/README.md" docs/site/examples.md; then
    echo "missing example link: $example" >&2
    exit 1
  fi
done

if ! grep -Fq "Mockport targets provider-compatible local APIs for selected workflows. It does not reproduce provider internals or undocumented behavior." docs/site/limitations.md; then
  echo "limitations page is missing provider-compatible boundary text" >&2
  exit 1
fi
