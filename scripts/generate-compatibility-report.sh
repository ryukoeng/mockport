#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="${1:-docs/compatibility-reports}"
PORT="${MOCKPORT_COMPATIBILITY_PORT:-43102}"
BASE_URL="http://127.0.0.1:${PORT}"
WORK_DIR="$(mktemp -d)"

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
mkdir -p "$OUT_DIR"
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
curl -fsS "$BASE_URL/_mockport/report" >"$WORK_DIR/runtime-report.json"

REPORT_DATE="${MOCKPORT_COMPATIBILITY_DATE:-$(date -u +%Y-%m-%d)}"
node "$ROOT_DIR/scripts/render-compatibility-report.mjs" "$WORK_DIR/runtime-report.json" "$OUT_DIR" "$REPORT_DATE"
