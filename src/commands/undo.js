import { simpleGit } from 'simple-git';

const git = simpleGit();

export default function registerUndoCommand(program) {
  program
    .command('undo')
    .description('Undo the last commit but keep files staged')
    .action(async () => {
      try {
        global.log.info('Undoing the last commit...');
        await git.reset(['--soft', 'HEAD~1']);
        
        global.log.info('Successfully undid the last commit.');
        global.log.info('Your files are still safe in the staging area.');
      } catch (error) {
        global.log.error('Failed to undo commit. Make sure there is a previous commit to undo.');
      }
    });
}