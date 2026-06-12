package cmd

import (
	"fmt"
	"strings"

	"github.com/ramadanny/gits/internal/gitops"
	"github.com/ramadanny/gits/internal/logger"
	"github.com/spf13/cobra"
)

var stackDiffFlag bool

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

				logger.Info(fmt.Sprintf("\n\x1b[7;1;36m  STACK %d: %s  \x1b[0m \x1b[1;37m%s\x1b[0m", i+1, hash, msg))

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

				if stackDiffFlag {
					diffOut, _ := gitops.Run("show", "--color=always", "--pretty=format:", "--patch", hash)
					diffTrimmed := strings.TrimSpace(diffOut)
					if diffTrimmed != "" {
						fmt.Println()
						diffLines := strings.Split(diffTrimmed, "\n")
						for _, line := range diffLines {
							fmt.Printf("      %s\n", line)
						}
						fmt.Println()
					}
				}
			}
		},
	}

	stackCmd.Flags().BoolVarP(&stackDiffFlag, "diff", "d", false, "Show detailed diff for each unpushed commit")

	rootCmd.AddCommand(stackCmd)
}