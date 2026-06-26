#!/usr/bin/env bash
set -euo pipefail

OUT_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$OUT_DIR"
}
trap cleanup EXIT

go run ./scripts/gen-compat-manifests --out "$OUT_DIR"

shopt -s nullglob
manifests=(compat/manifests/*.json)
if ((${#manifests[@]} == 0)); then
  echo "no checked-in manifests under compat/manifests/" >&2
  exit 1
fi

generated=("$OUT_DIR"/*.json)
if ((${#generated[@]} == 0)); then
  echo "no generated manifests under $OUT_DIR" >&2
  exit 1
fi

for f in "${generated[@]}"; do
  checked_in="compat/manifests/$(basename "$f")"
  if [[ ! -f "$checked_in" ]]; then
    echo "missing checked-in manifest: $checked_in" >&2
    exit 1
  fi
done

for f in "${manifests[@]}"; do
  generated_file="$OUT_DIR/$(basename "$f")"
  if [[ ! -f "$generated_file" ]]; then
    echo "stale checked-in manifest with no generated counterpart: $f" >&2
    exit 1
  fi
  diff -u "$f" "$generated_file" || {
    echo "manifest drift: $f — run 'go run ./scripts/gen-compat-manifests' and commit" >&2
    exit 1
  }
done
