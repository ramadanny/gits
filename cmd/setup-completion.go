package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ramadanny/gits/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	autoCmd := &cobra.Command{
		Use:   "setup-completion",
		Short: "Automatically install shell autocompletion for your terminal",
		Run:   executeSetupCompletion,
	}
	rootCmd.AddCommand(autoCmd)
}

func executeSetupCompletion(cmd *cobra.Command, args []string) {
	prompt := &survey.Select{
		Message: "Select your default terminal shell to install autocompletion:",
		Options: []string{"zsh", "bash", "fish"},
	}

	var shell string
	err := survey.AskOne(prompt, &shell)
	if err != nil {
		logger.Error("Operation aborted.")
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		logger.Error("Could not determine user home directory.")
		return
	}

	var rcFile string
	var injectLine string

	switch shell {
	case "zsh":
		rcFile = filepath.Join(home, ".zshrc")
		injectLine = "\nif command -v gits &> /dev/null; then\n  source <(gits completion zsh)\nfi\n"
	case "bash":
		rcFile = filepath.Join(home, ".bashrc")
		injectLine = "\nif command -v gits &> /dev/null; then\n  source <(gits completion bash)\nfi\n"
	case "fish":
		rcFile = filepath.Join(home, ".config", "fish", "config.fish")
		injectLine = "\nif type -q gits\n  gits completion fish | source\nend\n"
	}

	if shell == "fish" {
		os.MkdirAll(filepath.Dir(rcFile), 0755)
	}

	content, err := os.ReadFile(rcFile)
	if err == nil && strings.Contains(string(content), "gits completion") {
		logger.Info(fmt.Sprintf("Autocompletion is already installed in %s", rcFile))
		return
	}

	f, err := os.OpenFile(rcFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to open %s for writing.", rcFile))
		return
	}
	defer f.Close()

	if _, err := f.WriteString(injectLine); err != nil {
		logger.Error("Failed to write autocompletion script.")
		return
	}

	logger.Info(fmt.Sprintf("Autocompletion successfully installed in %s", rcFile))
	fmt.Printf("\nTo apply changes, please reload your terminal or run:\n")
	fmt.Printf("source %s\n\n", rcFile)
}