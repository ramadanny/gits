package cmd

import (
	"os"

	"github.com/ramadanny/gits/internal/config"
	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "gits",
	Short:   "A fast CLI tool for Git Push operations.",
	Version: Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.Init)

	asciiText := `
#  ██████╗  ██╗ ████████╗ ███████╗
# ██╔════╝  ██║ ╚══██╔══╝ ██╔════╝
# ██║  ███╗ ██║    ██║    ███████╗
# ██║   ██║ ██║    ██║    ╚════██║
# ╚██████╔╝ ██║    ██║    ███████║
#  ╚═════╝  ╚═╝    ╚═╝    ╚══════╝`

	rootCmd.Long = "\x1b[38;2;135;206;235m" + asciiText + "\x1b[0m\n\n\x1b[37mA fast CLI tool for Git Push operations.\x1b[0m"
}