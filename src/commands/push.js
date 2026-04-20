// src/commands/push.js
import { simpleGit } from "simple-git";
import Conf from "conf";
import fs from "fs";
import { GoogleGenerativeAI } from "@google/generative-ai";

const git = simpleGit();
const config = new Conf({ projectName: "ramadanny-gits" });

const BANNED_FILES = [".env", "id_rsa", "id_ed25519", "credentials.json"];
const BANNED_EXTENSIONS = [".key", ".pem"];

function formatBytes(bytes, decimals = 2) {
    if (!+bytes) return "0 Bytes";
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
}

export default function pushCmds(program) {
    program
        .command("push <path> [message]")
        .description("Push with safety checks. Use 'auto' as message for AI-generated commit.")
        .option("-f, --force", "Force push")
        .option("-u, --set-upstream", "Set upstream tracking")
        .option("-a, --all", "Push to all registered remotes simultaneously")
        .action(async (targetPath, message, options) => {
            try {
                if (!fs.existsSync(".git")) {
                    global.log.error(
                        "Error: Not a git repository. Please run 'gits remote add <url>' first to initialize."
                    );
                    return;
                }

                try {
                    await git.status();
                } catch (err) {
                    if (
                        err.message.includes("dubious ownership") ||
                        err.message.includes("safe.directory")
                    ) {
                        global.log.info(
                            "Detected dubious ownership. Automatically marking directory as safe."
                        );
                        await git.raw([
                            "config",
                            "--global",
                            "--add",
                            "safe.directory",
                            process.cwd(),
                        ]);
                        global.log.info("Directory marked as safe successfully.");
                    }
                }

                if (!fs.existsSync(".gitignore")) {
                    global.log.info("No .gitignore found. Generating default template.");
                    const defaultIgnore =
                        "node_modules/\n.env\n.env.*\n!.env.example\ndist/\nbuild/\n.DS_Store\ncoverage/\n";
                    fs.writeFileSync(".gitignore", defaultIgnore);
                    global.log.info("Created default .gitignore.");
                }

                global.log.info("Scanning for sensitive data.");
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
                                `CRITICAL: Sensitive file "${secret}" detected and NOT ignored.`
                            );
                        }
                    }
                    if (!isIgnored) {
                        global.log.error("\nPush aborted to prevent data leak.");
                        global.log.info("Please add the above files to .gitignore.");
                        return;
                    }
                }

                const remotes = await git.getRemotes(true);
                if (remotes.length === 0) {
                    global.log.error(
                        'Remote "origin" not found. Add it using "gits remote add <url>".'
                    );
                    return;
                }

                const username = config.get("username");
                const token = config.get("token");
                if (!username || !token) {
                    global.log.error('Error: Please run "gits setup" first.');
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

                let hasChanges = modifiedCount > 0 || createdCount > 0 || deletedCount > 0;
                let committedJustNow = false;

                if (hasChanges) {
                    let totalSize = 0;
                    const filesToProcess = [
                        ...statusBeforeAdd.modified,
                        ...statusBeforeAdd.not_added,
                    ];
                    for (const file of filesToProcess) {
                        try {
                            const stat = fs.statSync(file);
                            totalSize += stat.size;
                        } catch (e) {}
                    }

                    global.log.info("\nStatus Summary:");
                    global.log.info(`   Modified files: ${modifiedCount}`);
                    global.log.info(`   New files:      ${createdCount}`);
                    global.log.info(`   Deleted files:  ${deletedCount}`);
                    global.log.info(`   Total Size:     ${formatBytes(totalSize)}\n`);

                    await git.add(targetPath);

                    let finalMessage = message;

                    if (!finalMessage || finalMessage.toLowerCase() === "auto") {
                        const geminiKey = config.get("ramadanny-gits-gemini-key");
                        if (!geminiKey) {
                            global.log.error(
                                'Error: Gemini API Key not found. Please run "gits set gemini <key>" first.'
                            );
                            return;
                        }

                        global.log.info("Analyzing changes with Gemini AI...");
                        const diff = await git.diff(["--cached"]);

                        if (!diff) {
                            global.log.error(
                                "Error: No diff output available to generate commit message."
                            );
                            return;
                        }

                        try {
                            const genAI = new GoogleGenerativeAI(geminiKey);
                            const model = genAI.getGenerativeModel({
                                model: "gemini-flash-lite-latest",
                            });

                            const prompt = `Analyze the provided git diff and generate a concise, professional commit message.
You MUST adhere strictly to the Conventional Commits specification.
Use one of the following types based on the changes:
- feat: A new feature
- fix: A bug fix
- docs: Documentation only changes
- style: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- refactor: A code change that neither fixes a bug nor adds a feature
- perf: A code change that improves performance
- test: Adding missing tests or correcting existing tests
- chore: Changes to the build process or auxiliary tools and libraries such as documentation generation

Rules:
1. Output ONLY the raw commit message string.
2. Do NOT include markdown formatting, backticks, or quotes.
3. Do NOT add any conversational text or explanations.
4. Keep the first line (summary) under 72 characters.

Diff:
${diff.slice(0, 10000)}`;

                            const result = await model.generateContent(prompt);
                            finalMessage = result.response.text().trim();
                            global.log.info(`message: ${finalMessage}`);
                        } catch (aiError) {
                            global.log.error(
                                `Failed to generate message from Gemini: ${aiError.message}`
                            );
                            return;
                        }
                    }

                    await git.commit(finalMessage);
                    committedJustNow = true;
                } else {
                    global.log.info(
                        "No uncommitted changes detected. Proceeding to push existing commits..."
                    );
                }

                const status = await git.status();
                const branch = status.current || "main";
                const targetRemotes = options.all
                    ? remotes
                    : [remotes.find((r) => r.name === "origin")].filter(Boolean);

                let pushFailed = false;

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
                        pushFailed = true;
                        global.log.error(`Failed to push to ${remote.name}: ${pushErr.message}`);
                    }
                }

                if (pushFailed && committedJustNow) {
                    global.log.info(
                        "\nRolling back commit to keep files staged due to push failure..."
                    );
                    try {
                        await git.reset(["--soft", "HEAD~1"]);
                        global.log.info(
                            "Rollback successful. Fix the conflicts (e.g. git pull) and try again."
                        );
                    } catch (resetErr) {
                        global.log.error("Failed to rollback commit.");
                    }
                }
            } catch (error) {
                global.log.error("\nOperation failed:");
                global.log.error(error.message);
            }
        });
}