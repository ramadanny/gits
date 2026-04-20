# <div align='center'>GitS: Automated & Secure Git CLI</div>

<div align="center">

<img src="https://cdn.jsdelivr.net/gh/ramadanny/cdn@main/gits-logo.png" alt="Header Image" width="100%"/>

</div>

An advanced, intuitive Command Line Interface (CLI) wrapper for Git, designed to automate push operations, manage remotes, and enforce security policies (like preventing accidental credential leaks).

Built with Node.js, `gits` abstracts complex Git workflows into simple, safe commands while ensuring your repository remains clean and secure.

---

## 🎯 Project Overview

**GitS** aims to reduce friction in day-to-day Git operations. While traditional Git requires users to manually track status, handle credential management via credential helpers, and carefully avoid pushing sensitive files, `gits` provides an all-in-one execution pipeline.

When you run a push via `gits`, it automatically:

- Validates the presence of a `.gitignore` (generating one if missing).
- Scans for exposed secrets (`.env`, SSH keys) and prevents pushes if they are tracked.
- Checks Git status to prevent empty commits.
- Uses a globally configured Personal Access Token (PAT) to securely authenticate remote operations.
- **AI-Powered Commit Messages:** Analyzes your code changes using Google Gemini AI to generate meaningful commit messages automatically.

---

## ✨ Key Features

- **Automated Security Scanning:** Pre-push checks to detect and block banned files (`.env`, `id_rsa`, `.pem`, etc.) from being leaked to remotes.
- **AI-Generated Messages:** Automatically generates professional, Conventional Commits-compliant messages (feat, fix, docs, style, refactor, perf, test, chore) based on staged diffs.
- **Auto-Initialization:** Automatically initializes Git and standard `.gitignore` files if they are missing.
- **Centralized Configuration:** Built-in credential management utilizing global system configuration stores for Git platform tokens and Gemini API keys.
- **Multi-Remote Push:** Push to all registered remotes simultaneously with a single flag.
- **Commit Safety Net:** Soft undo functionality to instantly revert your last commit without losing staged changes.
- **Beautiful Console Output:** Colorful UI, informative status summaries, and dynamic ASCII art branding.

---

## 🛠 Tech Stack

- **Runtime:** [Node.js](https://nodejs.org/) (Requires v18+ / ES Modules)
- **CLI Framework:** [Commander.js](https://github.com/tj/commander.js/)
- **Git Interoperability:** [simple-git](https://github.com/steveukx/git-js)
- **AI Integration:** [@google/generative-ai](https://www.npmjs.com/package/@google/generative-ai)
- **Configuration Storage:** [conf](https://github.com/sindresorhus/conf)
- **Styling:** Custom ANSI escape codes

---

## 📂 Project Structure

```text
gits/
├── src/
│   ├── index.js               # Application entry point, CLI router, & UI logic
│   └── commands/              # Modular command definitions
│       ├── branch.js          # Branch creation and switching
│       ├── push.js            # Core push pipeline, security scanner, & AI logic
│       ├── remote.js          # Remote origin management (add/set/check)
│       ├── set.js             # Credential & Gemini API configuration
│       └── undo.js            # Commit reversion logic
├── .gitignore                 # Project ignore rules
├── LICENSE                    # Apache 2.0 License definitions
├── package.json               # Node.js dependencies and script registry
└── README.md                  # Project documentation
```

### Architectural Flow

1. **Routing:** `src/index.js` acts as the dispatcher, rendering the banner and delegating arguments to the appropriate command module via `commander`.
2. **Context:** Commands leverage `simple-git` for synchronous-like promise wrappers over the local Git binary.
3. **State:** User state (tokens/usernames/api keys) is securely persisted in the user's home directory under `ramadanny-gits` using the `conf` package.

---

## ⚙️ Prerequisites

Before installing `gits`, ensure your environment meets the following requirements:

- **Node.js**: `v18.0.0` or higher.
- **Git**: Installed and available in your system's PATH.
- **Gemini API Key**: Obtainable from [Google AI Studio](https://aistudio.google.com/).

---

## 📦 Installation

Currently, the tool can be installed globally from the source directory.

1. **Clone the repository:**

    ```bash
    git clone https://github.com/ramadanny/gits.git
    cd gits
    ```

2. **Install dependencies:**

    ```bash
    npm install
    ```

3. **Verify Installation:**
    ```bash
    gits --version
    ```

---

## 🔧 Configuration

Before pushing code, you must configure your Git platform credentials. `gits` uses a Personal Access Token (PAT) rather than prompting for a password.

### Interactive Setup

The easiest way to get started is the interactive wizard:

```bash
gits setup
```

_You will be prompted to enter your GitHub/GitLab username, your Personal Access Token, and optionally your Gemini API Key._

### Manual Configuration

You can also set these values individually:

```bash
gits set username <your-username>
gits set token <your-personal-access-token>
gits set gemini <your-gemini-api-key>
```

---

## 💻 Usage & Commands

### 1. Push Code (`gits push`)

The crown jewel of the CLI. Adds, commits, checks security, and pushes.

**Manual Commit:**

```bash
gits push index.js "Initial Commit"
```

**Auto Commit (AI-Generated):**

```bash
gits push index.js auto
```

_Analyzes changes and automatically generates a commit message following the Conventional Commits standard._

**Options:**

- `-f, --force`: Force push to the remote.
- `-u, --set-upstream`: Set upstream tracking for the branch.
- `-a, --all`: Push to all registered remotes at once.

### 2. Manage Remotes (`gits remote`)

View and configure your remote connections.

```bash
# View current remotes
gits remote check

# Add a new origin
gits remote add https://github.com/user/repo.git

# Modify the existing origin
gits remote set https://github.com/user/repo.git
```

### 3. Manage Branches (`gits branch`)

Seamlessly switch branches or create a new one if it doesn't exist.

```bash
gits branch feature/new-login
```

### 4. Undo Last Commit (`gits undo`)

Made a mistake in your commit message? Missed a file? This command performs a `git reset --soft HEAD~1`, stripping the last commit but keeping all your file changes staged.

```bash
gits undo
```

---

## 🛡 Security Features

`gits` prevents catastrophic data leaks by intercepting push operations.

**Banned Files:**
The push command will abort if it detects untracked or unignored files matching:

- `.env` (Except `.env.example`)
- `id_rsa` / `id_ed25519`
- `credentials.json`
- `*.key` / `*.pem`

If you trigger this fail-safe, `gits` will explicitly list the vulnerable files and instruct you to add them to your `.gitignore` before allowing the push to continue.

---

## 🚑 Troubleshooting

**Error: "Run 'gits setup' first."**

- **Cause:** You are trying to push without configuring your token.
- **Fix:** Run `gits setup` and provide your credentials.

**Error: "CRITICAL: Sensitive file detected and NOT ignored!"**

- **Cause:** You have a file like `.env` in your directory, and it is missing from `.gitignore`.
- **Fix:** Add the exact filename to `.gitignore` or delete the file.

**Error: "Failed to push... Authentication failed"**

- **Cause:** Your Personal Access Token has expired or lacks `repo` scopes.
- **Fix:** Generate a new PAT on your Git provider and run `gits set token <new-token>`.

**Error: "Gemini API Key not found"**

- **Cause:** Attempting to use auto-commit without setting the API key.
- **Fix:** Run `gits set gemini <key>` or `gits setup`.

**Error: "Failed to generate message from Gemini"**

- **Cause:** Invalid API key or model availability (Tool defaults to `gemini-pro`).
- **Fix:** Ensure your key is valid and has access to the Gemini API.

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!

1. Fork the project.
2. Create your feature branch: `gits branch feature/AmazingFeature`
3. Commit your changes (Feel free to use `gits push . auto`).
4. Push to the branch.
5. Open a Pull Request.

---

## 📄 License

This project is licensed under the **Apache License 2.0**.

You may use, reproduce, and distribute the work, provided you include prominent notices of modifications and retain attribution. See the [LICENSE](./LICENSE) file for complete details.

---

## 👏 Credits

- **Author:** ramadanny
- Powered by [simple-git](https://www.npmjs.com/package/simple-git), [commander](https://www.npmjs.com/package/commander), and [Google Generative AI](https://www.npmjs.com/package/@google/generative-ai).
