#!/usr/bin/env bash
set -euo pipefail

files=()
while IFS= read -r file; do
  files+=("$file")
done < <(find examples -name '.env.mockport.example' -type f | sort)

for file in README.md docs/public-env-safety.md docs/ai-safe-development.md docs/reporting.md docs/site/*.md; do
  [[ -f "$file" ]] && files+=("$file")
done

for file in "${files[@]}"; do
  if grep -En 'sk_(live|test)_|AKIA|ASIA|ghp_|github_pat_|xox[bp]-|AIza' "$file"; then
    echo "real-looking provider secret found in $file" >&2
    exit 1
  fi
  if grep -En 'whsec_' "$file" | grep -Fv 'whsec_mockport'; then
    echo "real-looking webhook secret found in $file" >&2
    exit 1
  fi
  if grep -En 'https://api\.stripe\.com|https://api\.openai\.com|https://api\.github\.com|https://api\.line\.me|https://slack\.com/api' "$file"; then
    echo "production provider URL found in $file" >&2
    exit 1
  fi
  if grep -Ein 'change[-_ ]?me|replace[-_ ]?me|changeme' "$file"; then
    echo "ambiguous placeholder found in $file" >&2
    exit 1
  fi
done
