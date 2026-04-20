import { simpleGit } from "simple-git";
import fs from "fs";

const git = simpleGit();

export default function remoteCmds(program) {
    const remoteCmd = program.command("remote").description("Manage remote repository URLs");

    remoteCmd
        .command("add <url>")
        .description("Add new origin URL")
        .action(async (url) => {
            try {
                if (!fs.existsSync(".git")) {
                    global.log.info("No git repository found. Initializing...");
                    await git.init();
                }
                await git.addRemote("origin", url);
                global.log.info(`Remote origin added successfully: ${url}`);
            } catch (error) {
                global.log.error("Failed to add remote. Make sure origin does not already exist.");
            }
        });

    remoteCmd
        .command("set <url>")
        .description("Change existing origin URL")
        .action(async (url) => {
            try {
                await git.remote(["set-url", "origin", url]);
                global.log.info(`Remote origin changed successfully to: ${url}`);
            } catch (error) {
                global.log.error("Failed to change remote. Make sure origin exists.");
            }
        });

    remoteCmd
        .command("check")
        .description("Check current remote URLs")
        .action(async () => {
            try {
                const remotes = await git.getRemotes(true);
                if (remotes.length === 0) {
                    global.log.info("No remote URL found in this project.");
                    return;
                }
                global.log.info("Remote Repositories:");
                remotes.forEach((r) => {
                    global.log.info(`   - ${r.name}: ${r.refs.push}`);
                });
            } catch (error) {
                global.log.error("Error: Not a git repository or Git not installed.");
                global.log.error(error.message);
            }
        });
}