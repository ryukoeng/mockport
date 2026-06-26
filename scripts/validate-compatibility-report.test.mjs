import { test } from "node:test";
import assert from "node:assert/strict";
import { validateReport } from "./validate-compatibility-report.mjs";

// promotion_eligible is the true value computed by Go's CanPromote. These tests
// check that the validator enforces it and also rejects self-declarations that
// contradict it (the provenance guard).
function baseAdapter(overrides = {}) {
  return {
    name: "stripe",
    maturity: "workflow-compatible",
    measured_level: "workflow",
    score: 100,
    promotion_eligible: true,
    sdk_coverage: 100,
    state_coverage: 100,
    error_coverage: 100,
    known_gaps: ["gap"],
    ...overrides,
  };
}

function contractEvidence(overrides = {}) {
  return {
    fixtures: ["compat/fixtures/stripe/checkout_session_create.json"],
    sdk_contracts: ["contract/sdk/stripe"],
    known_gaps: ["docs/compatibility-reports/latest.json#stripe"],
    ...overrides,
  };
}

// Replace the first adapter with the one under test; the rest are passing
// adapters so the required set of five is satisfied.
function reportWith(adapter) {
  const names = ["stripe", "openai", "github-oauth", "slack", "line"];
  const adapters = names.map((name, i) =>
    i === 0 ? { ...adapter, name } : baseAdapter({ name })
  );
  return { adapters };
}

function defaultManifests(maturity = "workflow-compatible") {
  return Object.fromEntries(
    ["stripe", "openai", "github-oauth", "slack", "line"].map((name) => [
      name,
      { adapter: name, maturity },
    ]),
  );
}

function validateWith(report, manifestOverrides = {}) {
  const manifests = { ...defaultManifests(), ...manifestOverrides };
  validateReport(report, { manifests });
}

test("accepts a report where every adapter is promotion-eligible and consistent", () => {
  assert.doesNotThrow(() => validateWith(reportWith(baseAdapter())));
});

test("accepts provider-compatible with complete contract evidence", () => {
  const adapter = baseAdapter({
    maturity: "provider-compatible",
    measured_level: "contract",
    contract_evidence: contractEvidence(),
  });
  assert.doesNotThrow(() =>
    validateWith(reportWith(adapter), { stripe: { adapter: "stripe", maturity: "provider-compatible" } }),
  );
});

test("rejects an adapter that is not promotion-eligible", () => {
  const adapter = baseAdapter({
    maturity: "provider-compatible",
    measured_level: "contract",
    promotion_eligible: false,
  });
  assert.throws(() => validateWith(reportWith(adapter)), /does not meet CanPromote/);
});

test("rejects an adapter missing the promotion_eligible field", () => {
  const adapter = baseAdapter();
  delete adapter.promotion_eligible;
  assert.throws(() => validateWith(reportWith(adapter)), /does not meet CanPromote/);
});

// Provenance guard: even with promotion_eligible=true, contradictory combinations are rejected.
test("rejects promotion_eligible=true that contradicts the maturity score floor", () => {
  const adapter = baseAdapter({ maturity: "provider-compatible", measured_level: "contract", score: 0 });
  assert.throws(() => validateWith(reportWith(adapter)), /score 0 < 80/);
});

test("rejects promotion_eligible=true with an impossible measured_level", () => {
  const adapter = baseAdapter({ maturity: "provider-compatible", measured_level: "wire", score: 100 });
  assert.throws(() => validateWith(reportWith(adapter)), /measured_level is "wire"/);
});

test("rejects provider-compatible without complete contract evidence", () => {
  const adapter = baseAdapter({
    maturity: "provider-compatible",
    measured_level: "contract",
  });
  assert.throws(() => validateWith(reportWith(adapter)), /contract_evidence is incomplete/);

  const partial = baseAdapter({
    maturity: "provider-compatible",
    measured_level: "contract",
    contract_evidence: contractEvidence({ sdk_contracts: [" "] }),
  });
  assert.throws(() => validateWith(reportWith(partial)), /contract_evidence is incomplete/);
});

test("rejects workflow-compatible whose state coverage is not 100", () => {
  const adapter = baseAdapter({ state_coverage: 0 });
  assert.throws(() => validateWith(reportWith(adapter)), /state_coverage is 0/);
});

test("rejects an adapter missing known gaps", () => {
  const adapter = baseAdapter({ known_gaps: [] });
  assert.throws(() => validateWith(reportWith(adapter)), /must publish known gaps/);
});

test("rejects a report missing a required adapter", () => {
  const report = { adapters: [baseAdapter({ name: "stripe" }), baseAdapter({ name: "openai" })] };
  assert.throws(() => validateReport(report, { manifests: defaultManifests() }), /at least five adapters|missing adapter/);
});

test("rejects report maturity that does not match the checked-in manifest", () => {
  const adapter = baseAdapter({
    maturity: "provider-compatible",
    measured_level: "contract",
    score: 100,
    contract_evidence: contractEvidence(),
  });
  assert.throws(
    () => validateWith(reportWith(adapter), { stripe: { adapter: "stripe", maturity: "workflow-compatible" } }),
    /maturity .* does not match manifest maturity/,
  );
});

test("rejects maturity increase when promotion_eligible is false", () => {
  const adapter = baseAdapter({
    maturity: "provider-compatible",
    measured_level: "contract",
    score: 100,
    promotion_eligible: false,
  });
  assert.throws(
    () => validateWith(reportWith(adapter), { stripe: { adapter: "stripe", maturity: "workflow-compatible" } }),
    /does not meet CanPromote|maturity .* exceeds manifest maturity/,
  );
});

test("skips manifest consistency for adapters outside the published manifest set", () => {
  const report = {
    adapters: [
      ...["stripe", "openai", "github-oauth", "slack", "line"].map((name) => baseAdapter({ name })),
      baseAdapter({ name: "zoho-oauth" }),
    ],
  };
  assert.doesNotThrow(() => validateWith(report));
});
