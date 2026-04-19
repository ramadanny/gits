import { simpleGit } from 'simple-git';
import chalk from 'chalk';

const git = simpleGit();

export default function registerRemoteCommand(program) {
  const remoteCmd = program.command('remote').description('Manage remote repository URLs');

  remoteCmd
    .command('add <url>')
    .description('Add new origin URL')
    .action(async (url) => {
      try {
        await git.addRemote('origin', url);
        console.log(chalk.green(`Remote origin added successfully: ${url}`));
      } catch (error) {
        console.log(chalk.red('Failed to add remote. Make sure origin does not already exist.'));
      }
    });

  remoteCmd
    .command('set <url>')
    .description('Change existing origin URL')
    .action(async (url) => {
      try {
        await git.remote(['set-url', 'origin', url]);
        console.log(chalk.green(`Remote origin changed successfully to: ${url}`));
      } catch (error) {
        console.log(chalk.red('Failed to change remote. Make sure origin exists.'));
      }
    });

  remoteCmd
    .command('check')
    .description('Check current remote URLs')
    .action(async () => {
      try {
        const remotes = await git.getRemotes(true);
        if (remotes.length === 0) {
          console.log(chalk.yellow('No remote URL found in this project.'));
          return;
        }
        console.log(chalk.cyan('Remote Repositories:'));
        remotes.forEach(r => {
          console.log(chalk.white(`   - ${r.name}: ${r.refs.push}`));
        });
      } catch (error) {
        console.log(chalk.red('Failed to read remotes.'));
      }
    });
}