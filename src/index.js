#!/usr/bin/env node

import { Command } from 'commander';
import fs from 'fs';
import registerPushCommand from './commands/push.js';
import registerSetCommand from './commands/set.js';
import registerRemoteCommand from './commands/remote.js';

// Read version from package.json (Notice the path is now '../package.json')
const pkg = JSON.parse(fs.readFileSync(new URL('../package.json', import.meta.url), 'utf8'));

const program = new Command();

program
  .name('gits')
  .description('A fast CLI tool for Git Push operations.')
  .version(pkg.version, '-v, --version', 'Output the current version');

// Register all commands
registerPushCommand(program);
registerSetCommand(program);
registerRemoteCommand(program);

program.parse(process.argv);