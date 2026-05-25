#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GO_BIN="${GO_BIN:-/usr/local/go/bin/go}"
IMAGE_TAG="ghcr.io/albert-einshutoin/mockport:latest"
WORK_DIR="$(mktemp -d)"
CONTAINER_NAME="mockport-smoke"

cleanup() {
  (cd "$WORK_DIR" && docker compose -f docker-compose.mockport.yml down >/dev/null 2>&1 || true)
  docker rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true
  rm -rf "$WORK_DIR"
}
trap cleanup EXIT

cd "$ROOT_DIR"
"$GO_BIN" build -o "$WORK_DIR/mockport" ./cmd/mockport
docker build -t "$IMAGE_TAG" -f docker/Dockerfile .

cd "$WORK_DIR"
"$WORK_DIR/mockport" init --adapter stripe
docker compose -f docker-compose.mockport.yml up -d

for _ in $(seq 1 30); do
  if curl -fsS http://localhost:43101/health >/dev/null; then
    break
  fi
  sleep 1
done

curl -fsS http://localhost:43101/health
printf '\n'
curl -fsS -X POST http://localhost:43101/stripe/v1/checkout/sessions
printf '\n'
"$WORK_DIR/mockport" report --url http://localhost:43101/_mockport/report
