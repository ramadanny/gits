import Conf from "conf";
import readline from "readline";

const config = new Conf({ projectName: "ramadanny-gits" });

export default function setCmds(program) {
    const setCmd = program.command("set").description("Set global credential configuration");

    setCmd
        .command("username <username>")
        .description("Save your Git platform username")
        .action((username) => {
            config.set("username", username);
            global.log.info(`Username saved successfully: ${username}`);
        });

    setCmd
        .command("token <token>")
        .description("Save your Personal Access Token (PAT)")
        .action((token) => {
            config.set("token", token);
            global.log.info("Token saved successfully.");
        });

    program
        .command("setup")
        .description("Interactive setup for username and token")
        .action(() => {
            const rl = readline.createInterface({
                input: process.stdin,
                output: process.stdout,
            });

            // Menggunakan ANSI escape code \x1b[37m untuk putih langsung di string pertanyaan
            rl.question("\x1b[37mEnter your GitHub Username: \x1b[0m", (usn) => {
                rl.question("\x1b[37mEnter your GitHub Personal Access Token: \x1b[0m", (tok) => {
                    config.set("username", usn);
                    config.set("token", tok);
                    global.log.info("\nConfiguration saved successfully!");
                    rl.close();
                });
            });
        });
}
