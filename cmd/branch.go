package cmd

import (
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
			
			_, err := gitops.Run("rev-parse", "--verify", name)

			if err == nil {
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