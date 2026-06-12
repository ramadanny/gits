package cmd

import (
	"fmt"
	"math"
	"os"
	"strings"

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
#  ╚═════╝  ╚═╝    ╚═╝    ╚══════╝
	
> A fast CLI tool for Git Push operations.`

	rootCmd.Long = applyRadialGradient(asciiText)
}

func applyRadialGradient(text string) string {
	lines := strings.Split(text, "\n")
	
	height := float64(len(lines))
	width := 0.0
	for _, line := range lines {
		if float64(len(line)) > width {
			width = float64(len(line))
		}
	}

	centerX := width / 2.0
	centerY := height / 2.0

	yRatio := 2.0

	maxDist := math.Sqrt(math.Pow(centerX, 2) + math.Pow(centerY*yRatio, 2))

	cR, cG, cB := 242.0, 230.0, 238.0
	eR, eG, eB := 151.0, 125.0, 255.0

	var result strings.Builder

	for y, line := range lines {
		for x, char := range line {
			if char == ' ' || char == '\n' {
				result.WriteRune(char)
				continue
			}

			dist := math.Sqrt(math.Pow(float64(x)-centerX, 2) + math.Pow((float64(y)-centerY)*yRatio, 2))
			
			t := dist / maxDist
			if t > 1.0 {
				t = 1.0
			}

			r := int(cR + t*(eR-cR))
			g := int(cG + t*(eG-cG))
			b := int(cB + t*(eB-cB))

			result.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm%c\x1b[0m", r, g, b, char))
		}
		if y < len(lines)-1 {
			result.WriteString("\n")
		}
	}
	return result.String()
}