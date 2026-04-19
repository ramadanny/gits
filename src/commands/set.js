import Conf from 'conf';
import chalk from 'chalk';

const config = new Conf({ projectName: 'gits-cli' });

export default function registerSetCommand(program) {
  const setCmd = program.command('set').description('Set global credential configuration');

  setCmd
    .command('username <username>')
    .description('Save your Git platform username')
    .action((username) => {
      config.set('username', username);
      console.log(chalk.green(`Username saved successfully: ${username}`));
    });

  setCmd
    .command('token <token>')
    .description('Save your Personal Access Token (PAT)')
    .action((token) => {
      config.set('token', token);
      console.log(chalk.green('Token saved successfully.'));
    });
}