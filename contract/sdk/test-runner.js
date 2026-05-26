#!/usr/bin/env node
"use strict";

const { runSmokePlaceholder } = require("./smoke-placeholder.test.js");

const allowedProviders = new Set(["all", "stripe", "openai", "github-oauth", "slack"]);

function parseArgs(argv) {
  const options = {
    baseURL: process.env.MOCKPORT_BASE_URL || "http://127.0.0.1:43101",
    provider: process.env.MOCKPORT_PROVIDER || "all",
    offline: false,
    json: false,
  };

  for (let index = 0; index < argv.length; index += 1) {
    const arg = argv[index];
    if (arg === "--offline") {
      options.offline = true;
    } else if (arg === "--json") {
      options.json = true;
    } else if (arg === "--provider") {
      options.provider = argv[index + 1];
      index += 1;
    } else if (arg.startsWith("--provider=")) {
      options.provider = arg.slice("--provider=".length);
    } else if (arg === "--base-url") {
      options.baseURL = argv[index + 1];
      index += 1;
    } else if (arg.startsWith("--base-url=")) {
      options.baseURL = arg.slice("--base-url=".length);
    } else {
      throw new Error(`unknown argument: ${arg}`);
    }
  }

  if (!allowedProviders.has(options.provider)) {
    throw new Error(`unsupported provider: ${options.provider}`);
  }
  return options;
}

async function main() {
  const options = parseArgs(process.argv.slice(2));
  const result = await runSmokePlaceholder(options);
  if (options.json) {
    process.stdout.write(`${JSON.stringify(result)}\n`);
  } else {
    process.stdout.write(`sdk-contracts provider=${result.provider} status=${result.status} baseURL=${result.baseURL}\n`);
  }
}

main().catch((error) => {
  process.stderr.write(`${error.message}\n`);
  process.exit(1);
});
