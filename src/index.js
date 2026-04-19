#!/usr/bin/env node

import { Command } from "commander";
import fs from "fs";

global.log = {
    info: (msg) => console.log(`\x1b[37m${msg}\x1b[0m`),
    error: (msg) => console.error(`\x1b[31m${msg}\x1b[0m`),
};

const gradients = [
    { c1: [100, 100, 100], c2: [200, 200, 200] },
    { c1: [255, 50, 150], c2: [255, 200, 220] },
    { c1: [150, 100, 200], c2: [250, 200, 250] },
    { c1: [255, 153, 255], c2: [153, 204, 255] },
];

function applyGradient(text, gradient) {
    const lines = text.split("\n");
    const { c1, c2 } = gradient;

    return lines
        .map((line, index) => {
            const t = lines.length > 1 ? index / (lines.length - 1) : 0;
            const r = Math.round(c1[0] + t * (c2[0] - c1[0]));
            const g = Math.round(c1[1] + t * (c2[1] - c1[1]));
            const b = Math.round(c1[2] + t * (c2[2] - c1[2]));
            return `\x1b[38;2;${r};${g};${b}m${line}\x1b[0m`;
        })
        .join("\n");
}

const asciiText = `
      █████████   ███   █████     █████████ 
     ███▒▒▒▒▒███ ▒▒▒   ▒▒███     ███▒▒▒▒▒███
    ███     ▒▒▒  ████  ███████  ▒███    ▒▒▒ 
   ▒███         ▒▒███ ▒▒▒███▒   ▒▒█████████ 
   ▒███    █████ ▒███   ▒███     ▒▒▒▒▒▒▒▒███
   ▒▒███  ▒███   ▒███   ▒███ ███ ███    ▒███
    ▒▒█████████  █████  ▒▒█████ ▒▒█████████ 
     ▒▒▒▒▒▒▒▒▒  ▒▒▒▒▒    ▒▒▒▒▒   ▒▒▒▒▒▒▒▒▒  `;

const randomGradient = gradients[Math.floor(Math.random() * gradients.length)];
const coloredAscii = applyGradient(asciiText, randomGradient);
const fullDescription = `${coloredAscii}\n\n\x1b[37mA fast CLI tool for Git Push operations.\x1b[0m`;

import pushCmds from "./commands/push.js";
import setCmds from "./commands/set.js";
import remoteCmds from "./commands/remote.js";
import branchCmds from "./commands/branch.js";
import undoCmds from "./commands/undo.js";

const pkg = JSON.parse(fs.readFileSync(new URL("../package.json", import.meta.url), "utf8"));
const program = new Command();

program
    .name("gits")
    .description(fullDescription)
    .version(pkg.version, "-v, --version", "Output the current version");

pushCmds(program);
setCmds(program);
remoteCmds(program);
branchCmds(program);
undoCmds(program);

program.parse(process.argv);
