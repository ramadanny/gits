package cmd

import (
	"os"

	"github.com/ramadanny/gits/internal/gitops"
	"github.com/ramadanny/gits/internal/logger"
	"github.com/spf13/cobra"
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage remote repository URLs",
}

func init() {
	rootCmd.AddCommand(remoteCmd)

	remoteCmd.AddCommand(&cobra.Command{
		Use:   "add [url]",
		Short: "Add new origin URL",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := os.Stat(".git"); os.IsNotExist(err) {
				logger.Info("No git repository found. Initializing...")
				gitops.Run("init")
			}
			_, err := gitops.Run("remote", "add", "origin", args[0])
			if err != nil {
				logger.Error("Failed to add remote. Make sure origin does not already exist.")
				return
			}
			logger.Info("Remote origin added successfully: " + args[0])
		},
	})

	remoteCmd.AddCommand(&cobra.Command{
		Use:   "set [url]",
		Short: "Change existing origin URL",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_, err := gitops.Run("remote", "set-url", "origin", args[0])
			if err != nil {
				logger.Error("Failed to change remote. Make sure origin exists.")
				return
			}
			logger.Info("Remote origin changed successfully to: " + args[0])
		},
	})

	remoteCmd.AddCommand(&cobra.Command{
		Use:   "check",
		Short: "Check current remote URLs",
		Run: func(cmd *cobra.Command, args []string) {
			out, err := gitops.Run("remote", "-v")
			if err != nil || out == "" {
				logger.Info("No remote URL found in this project.")
				return
			}
			logger.Info("Remote Repositories:\n" + out)
		},
	})
}