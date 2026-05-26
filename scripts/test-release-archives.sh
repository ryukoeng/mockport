#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="${VERSION:-0.0.0-test}"
DIST_DIR="$(mktemp -d)"

cleanup() {
  rm -rf "$DIST_DIR"
}
trap cleanup EXIT

cd "$ROOT_DIR"
scripts/build-release-archives.sh "$VERSION" "$DIST_DIR"

for target in linux_amd64 linux_arm64 darwin_amd64 darwin_arm64; do
  archive="$DIST_DIR/mockport_${VERSION}_${target}.tar.gz"
  if [[ ! -f "$archive" ]]; then
    echo "missing archive: $archive" >&2
    exit 1
  fi
done

if [[ ! -f "$DIST_DIR/checksums.txt" ]]; then
  echo "missing checksums.txt" >&2
  exit 1
fi

for target in linux_amd64 linux_arm64 darwin_amd64 darwin_arm64; do
  if ! grep -Fq "mockport_${VERSION}_${target}.tar.gz" "$DIST_DIR/checksums.txt"; then
    echo "missing checksum entry for $target" >&2
    exit 1
  fi
done

scripts/verify-release-artifacts.sh "$VERSION" "$DIST_DIR"
