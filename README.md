<div align="center">
  <h1>🚀 GitS: Automated & Secure Git CLI</h1>
  
<img src="https://cdn.jsdelivr.net/gh/ramadanny/cdn@main/gits-logo.png" alt="Header Image" width="100%"/>

  <p><strong>An advanced, high-performance CLI wrapper for Git, engineered for automated push operations, robust repository security, and AI-driven commit generation.</strong></p>
</div>

---

## 🎯 Overview

**GitS** is a powerful Command Line Interface (CLI) built with Go that abstracts complex Git workflows into streamlined, safe commands. Designed for senior developers and teams who value speed and security, GitS eliminates manual friction by automating pre-push security checks, state configuration, and even the formulation of Conventional Commits using Large Language Models (LLMs).

---

## ✨ Core Functions & Features

- 🛡️ **Pre-Push Security Scanner**: Automatically scans for exposed secrets (e.g., `.env`, `id_rsa`, `*.pem`, `credentials.json`) and blocks the push if these files are not properly ignored in `.gitignore`.
- 🤖 **AI-Powered Commits**: Integrates with Google's Gemini AI (`gemini-flash-lite-latest`) to analyze staged `git diff` outputs and automatically generate highly descriptive, Conventional Commits-compliant messages.
- 🛠️ **Auto-Remediation & Initialization**: Automatically resolves "dubious ownership" git errors by marking directories as safe, and automatically scaffolds standard `.gitignore` files if missing.
- 🗄️ **Centralized Credential Store**: Secure, stateful configuration management stored locally in your home directory (`/.config/ramadanny-gits`) using `Viper`.
- 🔄 **Multi-Remote Push**: Execute pushes to multiple registered remotes concurrently using a single `--all` flag.
- ↩️ **Safety Net (Soft Undo)**: Instantly revert your last commit without losing the actual file changes in your staging area.

---

## 📂 Project Structure

The architecture follows a clean, modular Go project layout utilizing the Cobra CLI framework:

```text
gits/
├── cmd/                  # CLI command definitions (push, branch, remote, set, undo)
├── internal/
│   ├── config/           # Configuration state management (Viper)
│   ├── gitops/           # Core wrapper for executing git binaries
│   └── logger/           # Custom stdout/stderr formatting and UI
├── .github/workflows/    # CI/CD pipelines for automated multi-platform releases
├── install.sh            # Interactive bash installation script
├── main.go               # Application entry point
└── go.mod                # Go module dependencies
```

---

## ⚙️ Prerequisites

- **Git**: Installed and accessible in your system's PATH.
- **Gemini API Key**: Required for the AI auto-commit feature (obtainable from Google AI Studio).

---

## 📦 Installation

To install GitS across any supported OS and architecture (Linux, macOS, Windows, Android/Termux), simply run the following command in your terminal:

```bash
curl -fsSL [https://raw.githubusercontent.com/ramadanny/gits/main/install.sh](https://raw.githubusercontent.com/ramadanny/gits/main/install.sh) | bash
```

---

## 💻 Usage Guide

### 1. Initial Setup
Configure your Git platform credentials and Gemini API key securely. GitS uses a Personal Access Token (PAT) for remote operations instead of passwords.

```bash
# Interactive wizard:
gits setup

# Or set them manually:
gits set username `your-username`
gits set token `your-personal-access-token`
gits set gemini `your-gemini-api-key
````

### 2. Pushing Code
The core command handles adding, scanning, committing, and pushing in one pipeline.

**Manual Commit Message:**
```bash
gits push . "feat: add user authentication"
```

**AI-Generated Commit Message:**
```bash
gits push . auto
```
*(GitS will analyze your code changes and generate a semantic commit message via Gemini).*

**Additional Push Flags:**
- `-f, --force`: Force push to the remote.
- `-u, --set-upstream`: Set upstream tracking for the current branch.
- `-a, --all`: Push to all registered remotes simultaneously.

### 3. Branch Management
Seamlessly switch to an existing branch or create a new one if it doesn't exist.
```bash
gits branch feature/new-payment-gateway
```

### 4. Remote Management
Manage your repository origins without touching standard git syntax.
```bash
gits remote check
gits remote add [https://github.com/user/repo.git](https://github.com/user/repo.git)
gits remote set [https://github.com/user/repo.git](https://github.com/user/repo.git)
```

### 5. Undo Operations
Made a mistake? Revert the last commit while keeping your files safely staged.
```bash
gits undo
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