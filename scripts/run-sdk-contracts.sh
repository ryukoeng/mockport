#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROVIDER="${1:-all}"
PORT="${MOCKPORT_CONTRACT_PORT:-43101}"
BASE_URL="http://127.0.0.1:${PORT}"
WORK_DIR="$(mktemp -d)"

case "$PROVIDER" in
  all|stripe|openai|github-oauth|slack) ;;
  *)
    echo "unsupported provider: $PROVIDER" >&2
    exit 1
    ;;
esac

if [[ -z "${GO_BIN:-}" ]]; then
  if command -v go >/dev/null 2>&1; then
    GO_BIN="go"
  else
    GO_BIN="/usr/local/go/bin/go"
  fi
fi

cleanup() {
  if [[ -n "${SERVER_PID:-}" ]]; then
    kill "$SERVER_PID" >/dev/null 2>&1 || true
    wait "$SERVER_PID" >/dev/null 2>&1 || true
  fi
  rm -rf "$WORK_DIR"
}
trap cleanup EXIT

cd "$ROOT_DIR"
"$GO_BIN" build -o "$WORK_DIR/mockport" ./cmd/mockport
cp examples/multi-adapter/mockport.yml "$WORK_DIR/mockport.yml"
python3 - "$WORK_DIR/mockport.yml" "$PORT" <<'PY'
import pathlib
import sys

path = pathlib.Path(sys.argv[1])
port = sys.argv[2]
text = path.read_text()
text = text.replace("  port: 43101", f"  port: {port}")
path.write_text(text)
PY

"$WORK_DIR/mockport" run --config "$WORK_DIR/mockport.yml" >"$WORK_DIR/mockport.log" 2>&1 &
SERVER_PID="$!"

for _ in $(seq 1 30); do
  if curl -fsS "$BASE_URL/health" >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

curl -fsS "$BASE_URL/health" >/dev/null

cd "$ROOT_DIR/contract/sdk"
npm ci
npm run test:live -- --provider "$PROVIDER" --base-url "$BASE_URL" --json
