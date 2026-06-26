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

for f in "${manifests[@]}"; do
  diff -u "$f" "$OUT_DIR/$(basename "$f")" || {
    echo "manifest drift: $f — run 'go run ./scripts/gen-compat-manifests' and commit" >&2
    exit 1
  }
done
