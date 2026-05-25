"use strict";

const assert = require("node:assert");
const { readFileSync } = require("node:fs");
const { join } = require("node:path");

const wrapper = readFileSync(join(__dirname, "..", "bin", "mockport.js"), "utf8");

assert(wrapper.includes("MOCKPORT_BIN"));
assert(wrapper.includes("spawnSync(binary"));
assert(wrapper.includes("docker"));
assert(wrapper.includes("ghcr.io/albert-einshutoin/mockport:latest"));
