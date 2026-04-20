package cmd

import (
	"github.com/ramadanny/gits/internal/gitops"
	"github.com/ramadanny/gits/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "undo",
		Short: "Undo the last commit but keep files staged",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Undoing the last commit...")
			_, err := gitops.Run("reset", "--soft", "HEAD~1")
			if err != nil {
				logger.Error("Failed to undo commit. Make sure there is a previous commit to undo.")
				return
			}
			logger.Info("Successfully undid the last commit.")
			logger.Info("Your files are still safe in the staging area.")
		},
	})
}