import { simpleGit } from "simple-git";

const git = simpleGit();

export default function branchCmdsprogram) {
    program
        .command("branch <name>")
        .description("Switch to an existing branch or create a new one")
        .action(async (name) => {
            try {
                const branchSummary = await git.branchLocal();

                if (branchSummary.all.includes(name)) {
                    global.log.info(`Switching to existing branch: ${name}...`);
                    await git.checkout(name);
                } else {
                    global.log.info(`Creating and switching to new branch: ${name}...`);
                    await git.checkoutLocalBranch(name);
                }

                global.log.info(`Successfully switched to branch "${name}".`);
            } catch (error) {
                global.log.error("Error switching branch:");
                global.log.error(error.message);
            }
        });
}
