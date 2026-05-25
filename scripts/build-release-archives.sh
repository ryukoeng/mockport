#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:?usage: scripts/build-release-archives.sh <version> [dist-dir] [goos] [goarch]}"
DIST_DIR="${2:-dist}"
ONLY_GOOS="${3:-}"
ONLY_GOARCH="${4:-}"

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
if [[ -z "${GO_BIN:-}" ]]; then
  if command -v go >/dev/null 2>&1; then
    GO_BIN="go"
  else
    GO_BIN="/usr/local/go/bin/go"
  fi
fi
CHECKSUM_BIN="${CHECKSUM_BIN:-}"

checksum() {
  local file="$1"
  if [[ -n "$CHECKSUM_BIN" ]]; then
    "$CHECKSUM_BIN" "$file"
  elif command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$file"
  else
    shasum -a 256 "$file"
  fi
}

build_one() {
  local goos="$1"
  local goarch="$2"
  local name="mockport_${VERSION}_${goos}_${goarch}"
  local work_dir="$DIST_DIR/$name"
  local binary="$work_dir/mockport"
  if [[ "$goos" == "windows" ]]; then
    binary="$work_dir/mockport.exe"
  fi

  rm -rf "$work_dir"
  mkdir -p "$work_dir"
  (
    cd "$ROOT_DIR"
    GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 "$GO_BIN" build \
      -ldflags "-X github.com/albert-einshutoin/mockport/internal/cli.version=${VERSION}" \
      -o "$binary" ./cmd/mockport
  )
  cp "$ROOT_DIR/README.md" "$work_dir/README.md"
  tar -C "$DIST_DIR" -czf "$DIST_DIR/${name}.tar.gz" "$name"
  rm -rf "$work_dir"
  checksum "$DIST_DIR/${name}.tar.gz" >>"$DIST_DIR/checksums.txt"
}

mkdir -p "$DIST_DIR"
: >"$DIST_DIR/checksums.txt"

if [[ -n "$ONLY_GOOS" || -n "$ONLY_GOARCH" ]]; then
  if [[ -z "$ONLY_GOOS" || -z "$ONLY_GOARCH" ]]; then
    echo "goos and goarch must be provided together" >&2
    exit 1
  fi
  build_one "$ONLY_GOOS" "$ONLY_GOARCH"
else
  build_one linux amd64
  build_one linux arm64
  build_one darwin amd64
  build_one darwin arm64
fi
