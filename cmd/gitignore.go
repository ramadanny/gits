package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ramadanny/gits/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gitignore [languages...]",
		Short: "Generate .gitignore rules from public templates",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			query := strings.Join(args, ",")
			url := "https://www.toptal.com/developers/gitignore/api/" + query

			resp, err := http.Get(url)
			if err != nil {
				logger.Error("Failed to communicate with gitignore API.")
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				logger.Error("Invalid language provided or template not found.")
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error("Failed to read response from API.")
				return
			}

			f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				logger.Error("Failed to open or create .gitignore file.")
				return
			}
			defer f.Close()

			if _, err := f.WriteString(fmt.Sprintf("\n%s\n", string(body))); err != nil {
				logger.Error("Failed to write templates to .gitignore file.")
				return
			}

			logger.Info("Templates successfully appended to .gitignore file.")
		},
	})
}