package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ramadanny/gits/internal/config"
	"github.com/ramadanny/gits/internal/logger"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set global credential configuration",
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.AddCommand(&cobra.Command{
		Use:   "username [username]",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config.Set("username", args[0])
			logger.Info("Username saved successfully: " + args[0])
		},
	})

	setCmd.AddCommand(&cobra.Command{
		Use:   "token [token]",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config.Set("token", args[0])
			logger.Info("Token saved successfully.")
		},
	})

	setCmd.AddCommand(&cobra.Command{
		Use:   "gemini [key]",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config.Set("ramadanny-gits-gemini-key", args[0])
			logger.Info("Gemini API Key saved successfully.")
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "setup",
		Short: "Interactive setup for username, token, and Gemini API",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("\x1b[37mEnter your GitHub Username: \x1b[0m")
			usn, _ := reader.ReadString('\n')
			
			fmt.Print("\x1b[37mEnter your GitHub Personal Access Token: \x1b[0m")
			tok, _ := reader.ReadString('\n')
			
			fmt.Print("\x1b[37mEnter your Gemini API Key (optional): \x1b[0m")
			gemini, _ := reader.ReadString('\n')

			config.Set("username", strings.TrimSpace(usn))
			config.Set("token", strings.TrimSpace(tok))
			if strings.TrimSpace(gemini) != "" {
				config.Set("ramadanny-gits-gemini-key", strings.TrimSpace(gemini))
			}
			logger.Info("\nConfiguration saved successfully!")
		},
	})
}