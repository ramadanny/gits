package cmd

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/ramadanny/gits/internal/config"
	"github.com/ramadanny/gits/internal/gitops"
	"github.com/ramadanny/gits/internal/logger"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var (
	force       bool
	setUpstream bool
	pushAll     bool
)

func init() {
	commitCmd := &cobra.Command{
		Use:   "commit [path] [message]",
		Short: "Add and commit with safety checks. Use 'auto' as message for AI-generated commit.",
		Args:  cobra.RangeArgs(1, 2),
		Run:   executeCommit,
	}

	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "Push existing commits to remote.",
		Run:   executePush,
	}
	
	pushCmd.Flags().BoolVarP(&force, "force", "f", false, "Force push")
	pushCmd.Flags().BoolVarP(&setUpstream, "set-upstream", "u", false, "Set upstream tracking")
	pushCmd.Flags().BoolVarP(&pushAll, "all", "a", false, "Push to all registered remotes simultaneously")

	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(pushCmd)
}

func executeCommit(cmd *cobra.Command, args []string) {
	targetPath := args[0]
	message := ""
	if len(args) > 1 {
		message = args[1]
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
		gitops.Run("add", targetPath)
		finalMessage := message

		if finalMessage == "" || strings.ToLower(finalMessage) == "auto" {
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
		}

		gitops.Run("commit", "-m", finalMessage)
		logger.Info("Successfully committed. Use 'gits push' to upload.")
	} else {
		logger.Info("No uncommitted changes detected in the specified path.")
	}
}

func executePush(cmd *cobra.Command, args []string) {
	remotes, _ := gitops.Run("remote")
	if remotes == "" {
		logger.Error("Remote \"origin\" not found.")
		return
	}

	username := config.Get("username")
	token := config.Get("token")
	if username == "" || token == "" {
		logger.Error("Error: Please run \"gits setup\" first.")
		return
	}

	var targetRemotes []string
	if pushAll {
		remotesOut, _ := gitops.Run("remote")
		for _, r := range strings.Split(strings.TrimSpace(remotesOut), "\n") {
			if r != "" {
				targetRemotes = append(targetRemotes, strings.TrimSpace(r))
			}
		}
	} else {
		targetRemotes = append(targetRemotes, "origin")
	}

	for _, remote := range targetRemotes {
		logger.Info(fmt.Sprintf("Pushing to %s...", remote))
		
		pushArgs := []string{}
		if username != "" && token != "" {
			authStr := base64.StdEncoding.EncodeToString([]byte(username + ":" + token))
			pushArgs = append(pushArgs, "-c", "http.extraHeader=Authorization: Basic "+authStr)
		}
		
		pushArgs = append(pushArgs, "push")

		if force {
			pushArgs = append(pushArgs, "-f")
		}
		if setUpstream {
			pushArgs = append(pushArgs, "-u", remote, "HEAD")
		} else {
			pushArgs = append(pushArgs, remote, "HEAD")
		}

		pushCmdExec := exec.Command("git", pushArgs...)
		pushCmdExec.Stdout = os.Stdout
		pushCmdExec.Stderr = os.Stderr
		err := pushCmdExec.Run()

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to push to %s.", remote))
		} else {
			logger.Info(fmt.Sprintf("Successfully pushed to %s.", remote))
		}
	}
}