package cmd

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/ramadanny/gits/internal/config"
	"github.com/ramadanny/gits/internal/gitops"
	"github.com/ramadanny/gits/internal/logger"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var (
	commitMessageFlag string
	autoCommitFlag    bool
)

func init() {
	commitCmd := &cobra.Command{
		Use:   "commit [paths...]",
		Short: "Add multiple files and commit. Use -m for message or --auto for AI.",
		Args:  cobra.MinimumNArgs(1),
		Run:   executeCommit,
	}

	commitCmd.Flags().StringVarP(&commitMessageFlag, "message", "m", "", "Manual commit message")
	commitCmd.Flags().BoolVarP(&autoCommitFlag, "auto", "a", false, "Generate commit message using Gemini AI")

	rootCmd.AddCommand(commitCmd)
}

func executeCommit(cmd *cobra.Command, args []string) {
	if !autoCommitFlag && commitMessageFlag == "" {
		logger.Error("Error: Please provide a commit message using -m \"message\" or use --auto for AI.")
		return
	}

	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		logger.Error("Error: Not a git repository.")
		return
	}

	_, err := gitops.Run("status")
	if err != nil && (strings.Contains(err.Error(), "dubious ownership") || strings.Contains(err.Error(), "safe.directory")) {
		logger.Info("Detected dubious ownership.")
		fmt.Print("\x1b[37mDo you want GitS to automatically mark this directory as safe globally? [Y/n]: \x1b[0m")
		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(strings.ToLower(choice))
		if choice == "y" || choice == "" {
			cwd, _ := os.Getwd()
			gitops.Run("config", "--global", "--add", "safe.directory", cwd)
			logger.Info("Directory marked as safe successfully.")
		} else {
			logger.Error("Operation aborted.")
			return
		}
	}

	if _, err := os.Stat(".gitignore"); os.IsNotExist(err) {
		logger.Info("No .gitignore found. Generating default template.")
		defaultIgnore := "node_modules/\n.env\n.env.*\n!.env.example\ndist/\nbuild/\n.DS_Store\ncoverage/\n"
		os.WriteFile(".gitignore", []byte(defaultIgnore), 0644)
	}

	logger.Info("Scanning for sensitive data.")
	files, _ := os.ReadDir(".")
	bannedFiles := []string{".env", "id_rsa", "id_ed25519", "credentials.json"}
	bannedExts := []string{".key", ".pem"}

	var detectedSecrets []string
	for _, file := range files {
		name := file.Name()
		isBanned := false
		for _, b := range bannedFiles {
			if name == b {
				isBanned = true
			}
		}
		for _, ext := range bannedExts {
			if strings.HasSuffix(name, ext) {
				isBanned = true
			}
		}
		if strings.Contains(name, ".env") && strings.HasSuffix(name, ".example") {
			isBanned = false
		}
		if isBanned {
			detectedSecrets = append(detectedSecrets, name)
		}
	}

	if len(detectedSecrets) > 0 {
		isSafe := true
		for _, secret := range detectedSecrets {
			_, errTracked := gitops.Run("ls-files", "--error-unmatch", secret)
			_, errIgnore := gitops.Run("check-ignore", "-q", secret)
			if errTracked == nil || errIgnore != nil {
				isSafe = false
				logger.Error(fmt.Sprintf("CRITICAL: Sensitive file \"%s\" is tracked or not safely ignored.", secret))
			}
		}
		if !isSafe {
			logger.Error("\nCommit aborted to prevent data leak.")
			return
		}
	}

	statusOut, _ := gitops.Run("status", "--porcelain")
	hasChanges := len(strings.TrimSpace(statusOut)) > 0

	if hasChanges {
		for _, targetPath := range args {
			_, errAdd := gitops.Run("add", targetPath)
			if errAdd != nil {
				logger.Error(fmt.Sprintf("[!] Failed to add path: %s", targetPath))
				return
			}
		}

		if !analyzeTodos() {
			logger.Error("\n[!] Commit aborted or no files left to commit.")
			return
		}

		finalMessage := ""

		if autoCommitFlag {
			geminiKey := config.Get("ramadanny-gits-gemini-key")
			if geminiKey == "" {
				logger.Error("Error: Gemini API Key not found.")
				return
			}

			logger.Info("Analyzing changes with Gemini AI...")
			diff, _ := gitops.Run("diff", "--cached")

			if len(diff) > 12000 {
				diff = diff[:12000] + "\n\n...[DIFF TRUNCATED DUE TO LENGTH]"
			}

			ctx := context.Background()
			net.DefaultResolver = &net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					d := net.Dialer{}
					return d.DialContext(ctx, "udp", "8.8.8.8:53")
				},
			}

			client, err := genai.NewClient(ctx, option.WithAPIKey(geminiKey))
			if err != nil {
				logger.Error(fmt.Sprintf("AI Client Error: %v", err))
				return
			}
			defer client.Close()

			model := client.GenerativeModel("gemini-flash-lite-latest")
			prompt := fmt.Sprintf(`Analyze the provided git diff and generate a concise, professional commit message.
You MUST adhere strictly to the Conventional Commits specification.
Use one of the following types based on the changes:
- feat: A new feature
- fix: A bug fix
- docs: Documentation only changes
- style: Changes that do not affect the meaning of the code
- refactor: A code change that neither fixes a bug nor adds a feature
- perf: A code change that improves performance
- test: Adding missing tests or correcting existing tests
- chore: Changes to the build process

Rules:
1. Output ONLY the raw commit message string.
2. Do NOT include markdown formatting, backticks, or quotes.
3. Do NOT add any conversational text or explanations.
4. Keep the first line (summary) under 72 characters.

Diff:
%s`, diff)

			resp, err := model.GenerateContent(ctx, genai.Text(prompt))

			if err != nil {
				logger.Error(fmt.Sprintf("Failed to generate message from Gemini: %v", err))
				return
			}

			if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil && len(resp.Candidates[0].Content.Parts) > 0 {
				finalMessage = fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
				finalMessage = strings.TrimSpace(finalMessage)
				logger.Info(fmt.Sprintf("message: %s", finalMessage))
			} else {
				logger.Error("Gemini returned empty response.")
				return
			}
		} else {
			finalMessage = commitMessageFlag
		}

		out, err := gitops.Run("commit", "-m", finalMessage)
		if err != nil {
			fmt.Println(strings.TrimSpace(string(out)))
		} else {
			logger.Info("Successfully committed. Use 'gits push' to upload.")
		}
	} else {
		logger.Info("No uncommitted changes detected in the specified paths.")
	}
}

func analyzeTodos() bool {
	out, err := exec.Command("git", "diff", "--cached", "--name-only").Output()
	if err != nil {
		return true
	}

	files := strings.Split(strings.TrimSpace(string(out)), "\n")
	var options []string
	fileMap := make(map[string]string)
	seenFiles := make(map[string]bool)

	for _, file := range files {
		if file == "" {
			continue
		}

		f, err := os.Open(file)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(f)
		lineNum := 1
		for scanner.Scan() {
			line := scanner.Text()
			upperLine := strings.ToUpper(line)

			if idx := strings.Index(upperLine, "TODO"); idx != -1 {
				snippet := strings.TrimSpace(line)
				if len(snippet) > 50 {
					snippet = snippet[:47] + "..."
				}

				if !seenFiles[file] {
					label := fmt.Sprintf("%s (line %d:%d) -> %s", file, lineNum, idx+1, snippet)
					options = append(options, label)
					fileMap[label] = file
					seenFiles[file] = true
				}
			}
			lineNum++
		}
		f.Close()
	}

	if len(options) == 0 {
		return true
	}

	fmt.Printf("\n\x1b[33m[!] Found TODOs in staged files:\x1b[0m\n")

	var selected []string

	prompt := &survey.MultiSelect{
		Message: "Select files to UNSTACK:",
		Options: options,
	}

	err = survey.AskOne(prompt, &selected)
	if err != nil {
		return false
	}

	for _, sel := range selected {
		fileToUnstack := fileMap[sel]
		exec.Command("git", "reset", "HEAD", fileToUnstack).Run()
		fmt.Printf("\x1b[34m[*] Unstacked %s\x1b[0m\n", fileToUnstack)
	}

	out, _ = exec.Command("git", "diff", "--cached", "--name-only").Output()
	if strings.TrimSpace(string(out)) == "" {
		return false
	}

	return true
}