package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/ramadanny/gits/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Git repository and generate default .gitsrc.json",
		Run:   executeInit,
	}
	rootCmd.AddCommand(initCmd)
}

func executeInit(cmd *cobra.Command, args []string) {
	out, err := exec.Command("git", "init").CombinedOutput()
	if err != nil {
		fmt.Print(string(out))
		return
	}
	fmt.Print(string(out))

	if _, err := os.Stat(".gitsrc.json"); os.IsNotExist(err) {
		cfg := Config{}
		cfg.Features.ScanSecrets = true
		cfg.Features.ScanTodos = true
		cfg.Features.AutoMarkSafeDirectory = true
		cfg.Commit.AIModel = "gemini-flash-lite-latest"
		cfg.Commit.PromptLanguage = "english"
		cfg.Security.CustomBannedFiles = []string{"config.production.json"}
		cfg.Security.CustomBannedExts = []string{".pfx"}

		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return
		}

		err = os.WriteFile(".gitsrc.json", data, 0644)
		if err == nil {
			logger.Info("Generated default .gitsrc.json configuration file.")
		}
	}
}