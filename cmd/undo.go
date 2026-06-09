package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ramadanny/gits/internal/gitops"
	"github.com/ramadanny/gits/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	undoCmd := &cobra.Command{
		Use:     "undo [number|all]",
		Aliases: []string{"unstack"},
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			target := "1"
			if len(args) > 0 {
				target = strings.ToLower(args[0])
			}

			branch, _ := gitops.Run("branch", "--show-current")
			out, err := gitops.Run("log", "--format=%h", "origin/"+branch+"..HEAD")
			if err != nil || out == "" {
				out, _ = gitops.Run("log", "--format=%h", "--not", "--remotes")
			}

			var commits []string
			for _, c := range strings.Split(strings.TrimSpace(out), "\n") {
				if c != "" {
					commits = append(commits, c)
				}
			}

			if len(commits) == 0 {
				logger.Info("No unpushed commits to undo.")
				return
			}

			var resetCount int
			if target == "all" {
				resetCount = len(commits)
			} else {
				num, err := strconv.Atoi(target)
				if err != nil || num < 1 || num > len(commits) {
					logger.Error("Invalid stack number.")
					return
				}
				resetCount = num
			}

			_, err = gitops.Run("reset", "--soft", fmt.Sprintf("HEAD~%d", resetCount))
			if err != nil {
				logger.Error("Failed to undo commits.")
				return
			}
			logger.Info(fmt.Sprintf("Successfully undid %d commit(s). Files remain staged.", resetCount))
		},
	}

	rootCmd.AddCommand(undoCmd)
}