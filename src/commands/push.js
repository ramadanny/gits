import { simpleGit } from 'simple-git';
import Conf from 'conf';
import chalk from 'chalk';

const git = simpleGit();
const config = new Conf({ projectName: 'gits-cli' });

export default function registerPushCommand(program) {
  program
    .command('push <message>')
    .description('Execute git add, commit, and push in one command')
    .option('-f, --force', 'Use this flag to force push')
    .action(async (message, options) => {
      try {
        const username = config.get('username');
        const token = config.get('token');

        console.log(chalk.cyan('Staging files...'));
        await git.add('.');

        console.log(chalk.cyan(`Committing changes: "${message}"...`));
        await git.commit(message);

        let remoteUrl = 'origin';
        
        if (username && token) {
          const remotes = await git.getRemotes(true);
          const origin = remotes.find(r => r.name === 'origin');
          
          if (origin && origin.refs.push.startsWith('https://')) {
              remoteUrl = origin.refs.push.replace('https://', `https://${username}:${token}@`);
              console.log(chalk.yellow('Using saved token for authentication.'));
          }
        } else {
          console.log(chalk.gray('Tip: Set username and token using "gits set" for automatic authentication.'));
        }

        console.log(chalk.cyan('Pushing to repository...'));
        
        const status = await git.status();
        const currentBranch = status.current;

        if (options.force) {
          await git.push(['-f', remoteUrl, currentBranch]);
        } else {
          await git.push([remoteUrl, currentBranch]);
        }

        console.log(chalk.green('Successfully pushed to repository.'));
      } catch (error) {
        console.log(chalk.red('Error during push operation:'));
        console.error(chalk.red(error.message));
      }
    });
}