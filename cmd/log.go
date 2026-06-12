package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/spf13/cobra"
)

var (
	logDiff   bool
	logLimit  int
	logSince  string
	logAuthor string
	logSearch string
)

func init() {
	logCmd := &cobra.Command{
		Use:   "log",
		Short: "View colorful and beautifully formatted commit history.",
		Run:   executeLog,
	}

	logCmd.Flags().BoolVarP(&logDiff, "diff", "d", false, "Show detailed diff and file stats")
	logCmd.Flags().IntVarP(&logLimit, "limit", "l", 0, "Limit the number of commits to show")
	logCmd.Flags().StringVarP(&logSince, "since", "s", "", "Show commits since specific time (e.g., 1w, 2d, 5h)")
	logCmd.Flags().StringVarP(&logAuthor, "author", "a", "", "Filter commits by author name")
	logCmd.Flags().StringVarP(&logSearch, "search", "g", "", "Search inside commit messages")

	rootCmd.AddCommand(logCmd)
}

func parseTimeFormat(t string) string {
	if t == "" {
		return ""
	}

	re := regexp.MustCompile(`^(\d+)([wdhms])$`)
	matches := re.FindStringSubmatch(t)

	if len(matches) == 3 {
		val := matches[1]
		unit := matches[2]

		switch unit {
		case "w":
			return val + " weeks ago"
		case "d":
			return val + " days ago"
		case "h":
			return val + " hours ago"
		case "m":
			return val + " minutes ago"
		case "s":
			return val + " seconds ago"
		}
	}

	return t
}

func executeLog(cmd *cobra.Command, args []string) {
	gitArgs := []string{"log", "--color=always"}

	if logDiff {
		gitArgs = append(gitArgs, "--stat", "-p", "--pretty=format:%n%C(reverse bold cyan)  COMMIT: %h  %Creset %C(bold yellow) %an %Creset %C(bold green) %ar %Creset%n%C(bold white) %s %Creset%n")
	} else {
		gitArgs = append(gitArgs, "--graph", "--abbrev-commit", "--decorate", "--format=format:%C(bold blue)%h%C(reset) - %C(bold green)(%ar)%C(reset) %C(white)%s%C(reset) %C(dim white)- %an%C(reset)%C(bold yellow)%d%C(reset)")
	}

	if logLimit > 0 {
		gitArgs = append(gitArgs, "-n", fmt.Sprint(logLimit))
	}

	parsedTime := parseTimeFormat(logSince)
	if parsedTime != "" {
		gitArgs = append(gitArgs, "--since="+parsedTime)
	}

	if logAuthor != "" {
		gitArgs = append(gitArgs, "--author="+logAuthor)
	}

	if logSearch != "" {
		gitArgs = append(gitArgs, "--grep="+logSearch, "-i")
	}

	cmdExec := exec.Command("git", gitArgs...)
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	cmdExec.Stdin = os.Stdin
	cmdExec.Run()
}