package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ramadanny/gits/internal/config"
	"github.com/ramadanny/gits/internal/gitops"
	"github.com/ramadanny/gits/internal/logger"
	"github.com/spf13/cobra"
)

var (
	force       bool
	setUpstream bool
	pushAll     bool
)

func init() {
	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "Push existing commits to remote.",
		Run:   executePush,
	}

	pushCmd.Flags().BoolVarP(&force, "force", "f", false, "Force push")
	pushCmd.Flags().BoolVarP(&setUpstream, "set-upstream", "u", false, "Set upstream tracking")
	pushCmd.Flags().BoolVarP(&pushAll, "all", "a", false, "Push to all registered remotes simultaneously")

	rootCmd.AddCommand(pushCmd)
}

func executePush(cmd *cobra.Command, args []string) {
	remotes, _ := gitops.Run("remote")
	if remotes == "" {
		logger.Error("Remote \"origin\" not found.")
		return
	}

	username := config.Get("username")
	token := config.Get("token")
	if username == "" || token == "" {
		logger.Error("Error: Please run \"gits setup\" first.")
		return
	}

	var targetRemotes []string
	if pushAll {
		remotesOut, _ := gitops.Run("remote")
		for _, r := range strings.Split(strings.TrimSpace(remotesOut), "\n") {
			if r != "" {
				targetRemotes = append(targetRemotes, strings.TrimSpace(r))
			}
		}
	} else {
		targetRemotes = append(targetRemotes, "origin")
	}

	for _, remote := range targetRemotes {
		logger.Info(fmt.Sprintf("Pushing to %s...", remote))

		pushArgs := []string{}
		if username != "" && token != "" {
			authStr := base64.StdEncoding.EncodeToString([]byte(username + ":" + token))
			pushArgs = append(pushArgs, "-c", "http.extraHeader=Authorization: Basic "+authStr)
		}

		pushArgs = append(pushArgs, "push")

		if force {
			pushArgs = append(pushArgs, "-f")
		}
		if setUpstream {
			pushArgs = append(pushArgs, "-u", remote, "HEAD")
		} else {
			pushArgs = append(pushArgs, remote, "HEAD")
		}

		pushCmdExec := exec.Command("git", pushArgs...)
		pushCmdExec.Stdout = os.Stdout
		pushCmdExec.Stderr = os.Stderr
		err := pushCmdExec.Run()

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to push to %s.", remote))
		} else {
			logger.Info(fmt.Sprintf("Successfully pushed to %s.", remote))
		}
	}
}