import { simpleGit } from "simple-git";
import Conf from "conf";
import fs from "fs";

const git = simpleGit();
const config = new Conf({ projectName: "ramadanny-gits" });

const BANNED_FILES = [".env", "id_rsa", "id_ed25519", "credentials.json"];
const BANNED_EXTENSIONS = [".key", ".pem"];

export default function pushCmds(program) {
    program
        .command("push <path> <message>")
        .description("Push with safety checks, auto-gitignore, and status summary")
        .option("-f, --force", "Force push")
        .option("-u, --set-upstream", "Set upstream tracking")
        .option("-a, --all", "Push to all registered remotes simultaneously")
        .action(async (targetPath, message, options) => {
            try {
                if (!fs.existsSync(".gitignore")) {
                    global.log.info("No .gitignore found. Generating default template...");
                    const defaultIgnore =
                        "node_modules/\n.env\n.env.*\n!.env.example\ndist/\nbuild/\n.DS_Store\ncoverage/\n";
                    fs.writeFileSync(".gitignore", defaultIgnore);
                    global.log.info("Created default .gitignore.");
                }

                global.log.info("Scanning for sensitive data...");
                const filesInDir = fs.readdirSync(".");
                const detectedSecrets = filesInDir.filter((file) => {
                    const isBannedFile = BANNED_FILES.includes(file);
                    const hasBannedExt = BANNED_EXTENSIONS.some((ext) => file.endsWith(ext));
                    if (file.includes(".env") && file.endsWith(".example")) return false;
                    return isBannedFile || hasBannedExt;
                });

                if (detectedSecrets.length > 0) {
                    let isIgnored = true;
                    const gitignoreContent = fs.readFileSync(".gitignore", "utf8");
                    for (const secret of detectedSecrets) {
                        if (!gitignoreContent.includes(secret)) {
                            isIgnored = false;
                            global.log.error(
                                `CRITICAL: Sensitive file "${secret}" detected and NOT ignored!`
                            );
                        }
                    }
                    if (!isIgnored) {
                        global.log.error("\nPush aborted to prevent data leak.");
                        global.log.info("Please add the files above to .gitignore.");
                        return;
                    }
                }

                if (!fs.existsSync(".git")) {
                    global.log.info("No git repo found. Initializing...");
                    await git.init();
                }

                const username = config.get("username");
                const token = config.get("token");
                if (!username || !token) {
                    global.log.error('Error: Run "gits setup" first.');
                    return;
                }

                if (targetPath !== "." && !fs.existsSync(targetPath)) {
                    global.log.error(`Error: Path "${targetPath}" not found.`);
                    return;
                }

                const statusBeforeAdd = await git.status();
                const modifiedCount = statusBeforeAdd.modified.length;
                const createdCount = statusBeforeAdd.not_added.length;
                const deletedCount = statusBeforeAdd.deleted.length;

                if (modifiedCount === 0 && createdCount === 0 && deletedCount === 0) {
                    global.log.info("No changes detected. Nothing to push.");
                    return;
                }

                global.log.info("\nStatus Summary:");
                global.log.info(`   Modified files: ${modifiedCount}`);
                global.log.info(`   New files:      ${createdCount}`);
                global.log.info(`   Deleted files:  ${deletedCount}\n`);

                await git.add(targetPath);
                await git.commit(message);

                const remotes = await git.getRemotes(true);
                if (remotes.length === 0) {
                    global.log.error(
                        'Remote "origin" not found. Add it using "gits remote add <url>"'
                    );
                    return;
                }

                const status = await git.status();
                const branch = status.current || "main";
                const targetRemotes = options.all
                    ? remotes
                    : [remotes.find((r) => r.name === "origin")].filter(Boolean);

                for (const remote of targetRemotes) {
                    const remoteUrl = remote.refs.push.replace(
                        "https://",
                        `https://${username}:${token}@`
                    );
                    const pushArgs = [remoteUrl, branch];
                    if (options.force) pushArgs.unshift("-f");
                    if (options.setUpstream) pushArgs.unshift("-u");

                    global.log.info(`Pushing to ${remote.name}/${branch}...`);
                    try {
                        await git.push(pushArgs);
                        global.log.info(`Successfully pushed to ${remote.name}.`);
                    } catch (pushErr) {
                        global.log.error(`Failed to push to ${remote.name}: ${pushErr.message}`);
                    }
                }
            } catch (error) {
                global.log.error("\nOperation failed:");
                global.log.error(error.message);
            }
        });
}
