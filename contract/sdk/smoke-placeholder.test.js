"use strict";

async function runSmokePlaceholder(options) {
  if (options.offline) {
    return {
      provider: options.provider,
      baseURL: options.baseURL,
      status: "offline-ok",
    };
  }

  const response = await fetch(new URL("/health", options.baseURL));
  if (!response.ok) {
    throw new Error(`health check failed: ${response.status}`);
  }
  const body = await response.json();
  if (body.status !== "ok") {
    throw new Error(`unexpected health body: ${JSON.stringify(body)}`);
  }
  return {
    provider: options.provider,
    baseURL: options.baseURL,
    status: "ok",
  };
}

module.exports = { runSmokePlaceholder };
