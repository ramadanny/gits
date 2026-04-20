package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
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
	pushCmd := &cobra.Command{
		Use:   "push [path] [message]",
		Short: "Push with safety checks. Use 'auto' as message for AI-generated commit.",
		Args:  cobra.RangeArgs(1, 2),
		Run:   executePush,
	}
	pushCmd.Flags().BoolVarP(&force, "force", "f", false, "Force push")
	pushCmd.Flags().BoolVarP(&setUpstream, "set-upstream", "u", false, "Set upstream tracking")
	pushCmd.Flags().BoolVarP(&pushAll, "all", "a", false, "Push to all registered remotes simultaneously")
	rootCmd.AddCommand(pushCmd)
}

func executePush(cmd *cobra.Command, args []string) {
	targetPath := args[0]
	message := ""
	if len(args) > 1 {
		message = args[1]
	}

	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		logger.Error("Error: Not a git repository. Please run 'gits remote add <url>' first to initialize.")
		return
	}

	_, err := gitops.Run("status")
	if err != nil && (strings.Contains(err.Error(), "dubious ownership") || strings.Contains(err.Error(), "safe.directory")) {
		logger.Info("Detected dubious ownership. Automatically marking directory as safe.")
		cwd, _ := os.Getwd()
		gitops.Run("config", "--global", "--add", "safe.directory", cwd)
		logger.Info("Directory marked as safe successfully.")
	}

	if _, err := os.Stat(".gitignore"); os.IsNotExist(err) {
		logger.Info("No .gitignore found. Generating default template.")
		defaultIgnore := "node_modules/\n.env\n.env.*\n!.env.example\ndist/\nbuild/\n.DS_Store\ncoverage/\n"
		os.WriteFile(".gitignore", []byte(defaultIgnore), 0644)
		logger.Info("Created default .gitignore.")
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
		isIgnored := true
		gitignoreContent, _ := os.ReadFile(".gitignore")
		ignoreStr := string(gitignoreContent)

		for _, secret := range detectedSecrets {
			if !strings.Contains(ignoreStr, secret) {
				isIgnored = false
				logger.Error(fmt.Sprintf("CRITICAL: Sensitive file \"%s\" detected and NOT ignored.", secret))
			}
		}
		if !isIgnored {
			logger.Error("\nPush aborted to prevent data leak.")
			logger.Info("Please add the above files to .gitignore.")
			return
		}
	}

	remotes, _ := gitops.Run("remote")
	if remotes == "" {
		logger.Error("Remote \"origin\" not found. Add it using \"gits remote add <url>\".")
		return
	}

	username := config.Get("username")
	token := config.Get("token")
	if username == "" || token == "" {
		logger.Error("Error: Please run \"gits setup\" first.")
		return
	}

	statusOut, _ := gitops.Run("status", "--porcelain")
	hasChanges := len(strings.TrimSpace(statusOut)) > 0
	committedJustNow := false

	if hasChanges {
		gitops.Run("add", targetPath)
		finalMessage := message

		if finalMessage == "" || strings.ToLower(finalMessage) == "auto" {
			geminiKey := config.Get("ramadanny-gits-gemini-key")
			if geminiKey == "" {
				logger.Error("Error: Gemini API Key not found. Please run \"gits set gemini <key>\" first.")
				return
			}

			logger.Info("Analyzing changes with Gemini AI...")
			diff, _ := gitops.Run("diff", "--cached")

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
				logger.Error("Gemini returned empty response or blocked by safety settings.")
				return
			}
		}

		gitops.Run("commit", "-m", finalMessage)
		committedJustNow = true
	} else {
		logger.Info("No uncommitted changes detected. Proceeding to push existing commits...")
	}

	remoteUrlOut, _ := gitops.Run("remote", "get-url", "origin")
	remoteUrl := strings.TrimSpace(remoteUrlOut)

	if strings.HasPrefix(remoteUrl, "https://") && username != "" && token != "" {
		authStr := fmt.Sprintf("https://%s:%s@", username, token)
		remoteUrl = strings.Replace(remoteUrl, "https://", authStr, 1)
	} else if remoteUrl == "" {
		remoteUrl = "origin"
	}

	pushArgs := []string{"push"}
	if force {
		pushArgs = append(pushArgs, "-f")
	}
	if setUpstream {
		pushArgs = append(pushArgs, "-u", remoteUrl, "HEAD")
	} else {
		pushArgs = append(pushArgs, remoteUrl, "HEAD")
	}

	safeUrlParts := strings.Split(remoteUrl, "@")
	safeUrlToPrint := remoteUrl
	if len(safeUrlParts) > 1 {
		safeUrlToPrint = safeUrlParts[len(safeUrlParts)-1]
	}
	logger.Info(fmt.Sprintf("Pushing to %s...", safeUrlToPrint))

	out, err := gitops.Run(pushArgs...)
	if err != nil {
		logger.Error("Failed to push.")
		if committedJustNow {
			logger.Info("\nRolling back commit to keep files staged due to push failure...")
			gitops.Run("reset", "--soft", "HEAD~1")
		}
	} else {
		logger.Info("Successfully pushed.")
		if out != "" {
			fmt.Println(out)
		}
	}
}