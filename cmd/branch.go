package cmd

import (
	"github.com/ramadanny/gits/internal/gitops"
	"github.com/ramadanny/gits/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "branch [name]",
		Short: "List branches, switch to an existing branch, or create a new one",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				out, _ := gitops.Run("branch")
				logger.Info("\n" + out)
				return
			}

			name := args[0]
			_, err := gitops.Run("rev-parse", "--verify", name)

			var checkoutErr error
			if err == nil {
				logger.Info("Switching to existing branch: " + name + "...")
				_, checkoutErr = gitops.Run("checkout", name)
			} else {
				logger.Info("Creating and switching to new branch: " + name + "...")
				_, checkoutErr = gitops.Run("checkout", "-b", name)
			}

			if checkoutErr != nil {
				logger.Error("Failed to switch branch. Make sure you don't have uncommitted changes that conflict.")
			} else {
				logger.Info("Successfully switched to branch \"" + name + "\".")
			}
		},
	})
}