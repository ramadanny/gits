package cmd

import (
	"fmt"
	"strings"

	"github.com/ramadanny/gits/internal/gitops"
	"github.com/ramadanny/gits/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	stackCmd := &cobra.Command{
		Use:     "stack",
		Aliases: []string{"log"},
		Short:   "View your unpushed local commit stack",
		Run: func(cmd *cobra.Command, args []string) {
			branch, _ := gitops.Run("branch", "--show-current")
			out, err := gitops.Run("log", "--format=%h|%s", "origin/"+branch+"..HEAD")
			if err != nil || out == "" {
				out, _ = gitops.Run("log", "--format=%h|%s", "--not", "--remotes")
			}

			out = strings.TrimSpace(out)
			if out == "" {
				logger.Info("No unpushed commits in the stack.")
				return
			}

			commits := strings.Split(out, "\n")
			for i, commit := range commits {
				parts := strings.SplitN(commit, "|", 2)
				if len(parts) < 2 {
					continue
				}
				hash, msg := parts[0], parts[1]

				logger.Info(fmt.Sprintf("\x1b[33m%d.\x1b[0m %s", i+1, msg))

				filesOut, _ := gitops.Run("diff-tree", "--no-commit-id", "--name-only", "-r", hash)
				files := strings.Split(strings.TrimSpace(filesOut), "\n")

				var validFiles []string
				for _, f := range files {
					if strings.TrimSpace(f) != "" {
						validFiles = append(validFiles, f)
					}
				}

				for j, file := range validFiles {
					prefix := "├── "
					if j == len(validFiles)-1 {
						prefix = "└── "
					}
					logger.Info(fmt.Sprintf("   %s%s", prefix, file))
				}
			}
		},
	}

	rootCmd.AddCommand(stackCmd)
}