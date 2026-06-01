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

require_file ".github/workflows/compatibility.yml"
require_text ".github/workflows/compatibility.yml" "workflow_dispatch"
require_text ".github/workflows/compatibility.yml" "schedule:"
require_text ".github/workflows/compatibility.yml" "bash scripts/run-sdk-contracts.sh stripe"
require_text ".github/workflows/compatibility.yml" "bash scripts/run-sdk-contracts.sh openai"
require_text ".github/workflows/compatibility.yml" "bash scripts/run-sdk-contracts.sh github-oauth"
require_text ".github/workflows/compatibility.yml" "bash scripts/run-sdk-contracts.sh slack"
require_text ".github/workflows/compatibility.yml" "actions/upload-artifact@v7"

require_file "scripts/generate-compatibility-report.sh"
require_file "docs/compatibility-reports/README.md"
require_file "docs/compatibility-reports/latest.md"
require_file "docs/compatibility-reports/latest.json"

require_text "docs/compatibility-reports/README.md" "provider-compatible"
require_text "docs/compatibility-reports/latest.md" "Compatibility Report"
require_text "docs/compatibility-reports/latest.md" "Known Gaps"
require_text "docs/compatibility-reports/latest.json" '"generated_by": "scripts/generate-compatibility-report.sh"'
require_text "docs/site/support-matrix.md" "provider-compatible"
require_text "CHANGELOG.md" "Compatibility release track"
require_text "CHANGELOG.md" "compatibility scores"

node <<'NODE'
const fs = require("fs");
const report = JSON.parse(fs.readFileSync("docs/compatibility-reports/latest.json", "utf8"));
if (!Array.isArray(report.adapters) || report.adapters.length < 5) {
  throw new Error("compatibility report must include at least five adapters");
}
const requiredAdapters = new Set(["stripe", "openai", "github-oauth", "slack", "line"]);
for (const name of requiredAdapters) {
  if (!report.adapters.some((adapter) => adapter.name === name)) {
    throw new Error(`compatibility report missing adapter: ${name}`);
  }
}
for (const adapter of report.adapters) {
  if (!adapter.name || !adapter.maturity || !Number.isInteger(adapter.score)) {
    throw new Error(`invalid adapter report entry: ${JSON.stringify(adapter)}`);
  }
  if (adapter.maturity === "provider-compatible" && adapter.score < 80) {
    throw new Error(`${adapter.name} is provider-compatible with score ${adapter.score}`);
  }
  if (adapter.maturity === "provider-compatible" && adapter.measured_level !== "contract") {
    throw new Error(`${adapter.name} is provider-compatible without contract-level evidence`);
  }
  if (adapter.maturity === "workflow-compatible" && adapter.score < 60) {
    throw new Error(`${adapter.name} is workflow-compatible with score ${adapter.score}`);
  }
  if (adapter.maturity === "sdk-compatible" && adapter.score < 40) {
    throw new Error(`${adapter.name} is sdk-compatible with score ${adapter.score}`);
  }
  if (!Array.isArray(adapter.known_gaps)) {
    throw new Error(`${adapter.name} missing known_gaps array`);
  }
  if (adapter.known_gaps.length === 0) {
    throw new Error(`${adapter.name} must publish known gaps`);
  }
}
NODE
