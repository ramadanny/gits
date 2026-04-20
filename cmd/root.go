package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/ramadanny/gits/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gits",
	Short:   "A fast CLI tool for Git Push operations.",
	Version: "0.0.2",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.Init)
	
	asciiText := `
#      █████████   ███   █████     █████████ 
#     ███▒▒▒▒▒███ ▒▒▒   ▒▒███     ███▒▒▒▒▒███
#    ███     ▒▒▒  ████  ███████  ▒███    ▒▒▒ 
#   ▒███         ▒▒███ ▒▒▒███▒   ▒▒█████████ 
#   ▒███    █████ ▒███   ▒███     ▒▒▒▒▒▒▒▒███
#   ▒▒███  ▒███   ▒███   ▒███ ███ ███    ▒███
#    ▒▒█████████  █████  ▒▒█████ ▒▒█████████ 
#     ▒▒▒▒▒▒▒▒▒  ▒▒▒▒▒    ▒▒▒▒▒   ▒▒▒▒▒▒▒▒▒  `

	rootCmd.Long = applyGradient(asciiText) + "\n\n\x1b[37mA fast CLI tool for Git Push operations.\x1b[0m"
}

func applyGradient(text string) string {
	grad := []int{0, 0, 139, 255, 255, 255}
	
	lines := strings.Split(text, "\n")
	var result []string
	for i, line := range lines {
		t := 0.0
		if len(lines) > 1 {
			t = float64(i) / float64(len(lines)-1)
		}
		r := int(float64(grad[0]) + t*float64(grad[3]-grad[0]))
		g := int(float64(grad[1]) + t*float64(grad[4]-grad[1]))
		b := int(float64(grad[2]) + t*float64(grad[5]-grad[2]))
		result = append(result, fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, line))
	}
	return strings.Join(result, "\n")
}

