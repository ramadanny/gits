package cmd

import (
	"strings"

	"github.com/ramadanny/gits/internal/gitops"
	"github.com/ramadanny/gits/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "branch [name]",
		Short: "Switch to an existing branch or create a new one",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			branches, _ := gitops.Run("branch", "--list")

			if strings.Contains(branches, name) {
				logger.Info("Switching to existing branch: " + name + "...")
				gitops.Run("checkout", name)
			} else {
				logger.Info("Creating and switching to new branch: " + name + "...")
				gitops.Run("checkout", "-b", name)
			}
			logger.Info("Successfully switched to branch \"" + name + "\".")
		},
	})
}