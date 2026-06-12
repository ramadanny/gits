<div align="center">
  <h1>⚡ GitS: No more complex Git commands, Just pure automation.</h1>
  
<img src="https://cdn.jsdelivr.net/gh/ramadanny/cdn@main/gits-logo.png" alt="Header Image" width="100%"/>

  <p><strong>An advanced, high-performance CLI wrapper for Git, engineered for automated push operations, robust repository security, and AI-driven commit generation.</strong></p>
</div>

---

## 🎯 Overview

**GitS** abstracts complex Git workflows into a streamlined, secure, and developer-friendly experience. Built with Go and the Cobra framework, **GitS** introduces an atomic workflow that separates local staging and commit management from remote synchronization. It features automated secret scanning, prompt-based remediation for environment anomalies, and seamless integration with Google's Gemini AI to formulate semantic, Conventional Commits-compliant messages.

---

## ✨ Core Functions & Features

- 📦 **Atomic Commit Stacking**: Separate your local save points from remote pushes. Work completely offline, layer multiple commits, and push them all at once.
- 🤖 **AI-Powered Commits**: Integrates with Google's Gemini AI (`gemini-flash-lite-latest`) to analyze staged `git diff` outputs. It auto-prioritizes feature logic over style/formatting changes and handles large diffs gracefully via intelligent truncation.
- 🛡️ **Advanced Security Scanner**: Prevents data leaks by scanning for exposed credentials (`.env`, `id_rsa`, `*.pem`, `credentials.json`). It directly queries the Git index (`ls-files` and `check-ignore`), stopping tracking even if files are accidentally staged.
- 🗂️ **Visual Stack Management**: View your unpushed local commit stack as an interactive tree structure showing exactly which files were modified per commit.
- ↩️ **Granular Unstacking (Undo)**: Safely rollback a specific number of local commits or empty the entire stack at once, keeping all your file changes safely staged in the staging area.
- 🛠️ **Smart Environmental Remediation**: Automatically detects "dubious ownership" anomalies and safely prompts the user before applying configuration modifications.
- 🔒 **Secure Authorization Headers**: Eliminates credential exposure by injecting Personal Access Tokens (PAT) dynamically via Git HTTP extraHeaders instead of plain-text URL interpolation.

---

## 📂 Project Structure

The architecture follows a modular, production-ready Go layout:

```text
gits/
├── cmd/                  # CLI command definitions (commit, push, stack, undo, branch, remote, set)
├── internal/
│   ├── config/           # Local configuration state management (Viper)
│   ├── gitops/           # Safe wrapper for executing system git binaries
│   └── logger/           # Custom stdout/stderr UI formatting
├── main.go               # Application entry point
├── go.mod                # Go module dependencies
└── install.sh            # Multi-platform interactive installation script
```

---

## 📦 Installation

To install **GitS** across any supported operating system and architecture, run the following command in your terminal:

```bash
curl -fsSL https://raw.githubusercontent.com/ramadanny/gits/main/install.sh | bash
```

### 🔑 Critical Permissions Configuration

Depending on your environment, you **must** grant executable permissions to the installed binary before running it for the first time to avoid `permission denied` errors.

#### 🤖 Android (Termux)
Run the following command to grant execution rights within the Termux binary path:
```bash
chmod +x $PREFIX/bin/gits
```

#### 🐧 Linux / 🍏 macOS
Run the following command using elevated privileges to update the system binary path:
```bash
sudo chmod +x /usr/local/bin/gits
```

#### 🪟 Windows
The binary (`gits.exe`) is downloaded to your current working directory. Ensure it is moved to a directory included in your system's `PATH` environment variable. No manual `chmod` is required.

---

## 💻 Complete Usage Guide

### 1. Initial Setup
Configure your Git credentials and Gemini API key securely. **GitS** uses a Personal Access Token (PAT) for all remote actions.

```bash
# Interactive wizard:
gits setup

# Manual configuration:
gits set username <your-github-username>
gits set token <your-personal-access-token>
gits set gemini <your-gemini-api-key>
```

### 2. Staging & Committing Code
**GitS** decouples committing from pushing, allowing you to build a local stack of changes.

**AI-Generated Semantic Commit:**
```bash
gits commit . auto
```
*GitS will analyze your staged diff, isolate logical changes from formatting, and automatically commit using the Conventional Commits specification.*

**Specific Path Commit:**
```bash
gits commit src/controllers/ auto
```

**Manual Commit Message:**
```bash
gits commit . "feat: implement user authentication session"
```

### 3. Inspecting the Local Stack
Before pushing, inspect your unpushed local commits and see a tree structure of affected files.

```bash
gits stack
# OR using the alias:
gits log
```

*Example Output:*
```text
1. feat: implement user authentication session
   ├── cmd/push.go
   └── cmd/stack.go
2. style: format codebase with prettier
   └── src/controllers/user.js
```

### 4. Unstacking & Undoing Commits
If you need to fix a commit or re-write history before pushing, use the undo command. Files will remain safely staged.

```bash
# Undo the single latest commit:
gits undo

# Undo a specific number of commits from the top of the stack:
gits undo 2

# Clear the entire unpushed local stack:
gits undo all
```
*Note: `gits unstack` can be used as a direct alias for `gits undo`.*

### 5. Pushing the Stack to Remote
Once your stack is ready, push all accumulated local commits to the remote repository concurrently.

```bash
gits push
```

**Available Flags:**
- `-f, --force`: Force push changes to the remote repository.
- `-u, --set-upstream`: Establish upstream tracking for the current branch.
- `-a, --all`: Concurrently push to all registered remote repositories.

### 6. Branch & Remote Utilities
Manage your workspace state cleanly without touching standard git syntax.

```bash
# Precise exact-match branch switching or creation:
gits branch feature/payment-gateway

# Remote configuration audits:
gits remote check
gits remote add <repository-url>
gits remote set <repository-url>
```

---

## 👏 Credits & Acknowledgements

- **Author:** ramadanny
- **Core Technologies:** Written in [Go](https://go.dev/).
- **Libraries:**
  - [Cobra](https://github.com/spf13/cobra) for powerful CLI command routing.
  - [Viper](https://github.com/spf13/viper) for application configuration management.
  - [Google Generative AI Go SDK](https://github.com/google/generative-ai-go) for LLM capabilities.

---
<div align="center">
  <i>Licensed under the Apache License 2.0. See the LICENSE file for details.</i>
</div>