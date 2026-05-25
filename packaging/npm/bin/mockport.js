#!/usr/bin/env node
"use strict";

const { spawnSync } = require("node:child_process");
const { existsSync } = require("node:fs");
const { join } = require("node:path");

const args = process.argv.slice(2);
const explicitBinary = process.env.MOCKPORT_BIN;
const packageBinary = join(__dirname, "..", "vendor", process.platform, process.arch, "mockport");
const binary = explicitBinary || packageBinary;

if (existsSync(binary)) {
  const result = spawnSync(binary, args, { stdio: "inherit" });
  process.exit(result.status === null ? 1 : result.status);
}

const dockerArgs = [
  "run",
  "--rm",
  "-p",
  "43101:43101",
  "-v",
  `${process.cwd()}/mockport.yml:/etc/mockport/mockport.yml`,
  "ghcr.io/albert-einshutoin/mockport:latest",
  ...args,
];
const result = spawnSync("docker", dockerArgs, { stdio: "inherit" });
process.exit(result.status === null ? 1 : result.status);
