package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/ramadanny/gits/internal/logger"
	"github.com/spf13/cobra"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func init() {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update GitS via an interactive installer wizard",
		Long:  "Connects to GitHub to list releases, lets you choose a version and binary, and handles the installation safely.",
		Run: func(cmd *cobra.Command, args []string) {
			net.DefaultResolver = &net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					d := net.Dialer{}
					return d.DialContext(ctx, "udp", "8.8.8.8:53")
				},
			}

			reader := bufio.NewReader(os.Stdin)

			askYesNo := func(prompt string) bool {
				fmt.Print(prompt)
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(strings.ToLower(input))
				return input == "" || input == "y" || input == "yes"
			}

			logger.Info("\x1b[34m[*] Fetching available releases from GitHub...\x1b[0m")
			resp, err := http.Get("https://api.github.com/repos/ramadanny/gits/releases")
			if err != nil {
				logger.Error(fmt.Sprintf("[!] Failed to contact GitHub API: %v", err))
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				logger.Error(fmt.Sprintf("[!] GitHub API returned error status: %s", resp.Status))
				return
			}

			var releases []GitHubRelease
			if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
				logger.Error(fmt.Sprintf("[!] Failed to decode release data: %v", err))
				return
			}

			if len(releases) == 0 {
				logger.Error("[!] No releases found in the repository.")
				return
			}

			fmt.Println("\nAvailable versions:")
			maxReleases := 5
			if len(releases) < maxReleases {
				maxReleases = len(releases)
			}
			for i := 0; i < maxReleases; i++ {
				fmt.Printf("  %d. %s\n", i+1, releases[i].TagName)
			}

			fmt.Printf("\n\x1b[34m[?]\x1b[0m Select version number [1-%d]: ", maxReleases)
			verInput, _ := reader.ReadString('\n')
			verIdx, err := strconv.Atoi(strings.TrimSpace(verInput))
			if err != nil || verIdx < 1 || verIdx > maxReleases {
				logger.Error("[!] Invalid version selection.")
				return
			}
			selectedRelease := releases[verIdx-1]

			ext := ""
			if runtime.GOOS == "windows" {
				ext = ".exe"
			}
			detectedAsset := fmt.Sprintf("gits-%s-%s%s", runtime.GOOS, runtime.GOARCH, ext)

			var downloadURL string
			var selectedAsset string

			assetExists := false
			for _, asset := range selectedRelease.Assets {
				if asset.Name == detectedAsset {
					assetExists = true
					downloadURL = asset.DownloadURL
					selectedAsset = asset.Name
					break
				}
			}

			useAuto := false
			if assetExists {
				useAuto = askYesNo(fmt.Sprintf("\x1b[34m[?]\x1b[0m Use auto-detected binary [%s]? [Y/n]: ", detectedAsset))
			}

			if useAuto {
				downloadURL = selectedRelease.Assets[0].DownloadURL 
				for _, asset := range selectedRelease.Assets {
					if asset.Name == detectedAsset {
						downloadURL = asset.DownloadURL
						selectedAsset = asset.Name
					}
				}
			} else {
				fmt.Println("\nAvailable binaries for " + selectedRelease.TagName + ":")
				for i, asset := range selectedRelease.Assets {
					fmt.Printf("  %d. %s\n", i+1, asset.Name)
				}
				fmt.Printf("\n\x1b[34m[?]\x1b[0m Select binary number [1-%d]: ", len(selectedRelease.Assets))
				assetInput, _ := reader.ReadString('\n')
				assetIdx, err := strconv.Atoi(strings.TrimSpace(assetInput))
				if err != nil || assetIdx < 1 || assetIdx > len(selectedRelease.Assets) {
					logger.Error("[!] Invalid binary selection.")
					return
				}
				downloadURL = selectedRelease.Assets[assetIdx-1].DownloadURL
				selectedAsset = selectedRelease.Assets[assetIdx-1].Name
			}

			logger.Info(fmt.Sprintf("\n\x1b[34m[*] Downloading %s...\x1b[0m", selectedAsset))
			execPath, err := os.Executable()
			if err != nil {
				logger.Error(fmt.Sprintf("[!] Failed to detect current installation path: %v", err))
				return
			}
			execDir := filepath.Dir(execPath)

			tmpFile, err := os.CreateTemp(execDir, "gits-update-tmp-")
			if err != nil {
				logger.Error(fmt.Sprintf("[!] Failed to create temporary file: %v", err))
				return
			}
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			assetResp, err := http.Get(downloadURL)
			if err != nil {
				logger.Error(fmt.Sprintf("[!] Failed to download binary package: %v", err))
				return
			}
			defer assetResp.Body.Close()

			if assetResp.StatusCode != http.StatusOK {
				logger.Error(fmt.Sprintf("[!] Failed to download asset, status: %s", assetResp.Status))
				return
			}

			_, err = io.Copy(tmpFile, assetResp.Body)
			if err != nil {
				logger.Error(fmt.Sprintf("[!] Failed to save binary file: %v", err))
				return
			}

			tmpFile.Close()

			logger.Info("\x1b[34m[*] Installing system update...\x1b[0m")
			if err := replaceBinary(tmpFile.Name(), execPath); err != nil {
				if strings.Contains(err.Error(), "permission denied") {
					if runtime.GOOS != "windows" {
						logger.Error("\n[!] Permission Denied. Please re-run the command with elevated privileges:\nsudo gits update")
					} else {
						logger.Error("\n[!] Permission Denied. Please ensure your terminal is running as Administrator.")
					}
				} else {
					logger.Error(fmt.Sprintf("[!] Failed to replace system binary: %v", err))
				}
				return
			}

			logger.Info(fmt.Sprintf("\x1b[32m[+] Success. GitS has been updated to version %s.\x1b[0m", selectedRelease.TagName))
		},
	}

	rootCmd.AddCommand(updateCmd)
}

func replaceBinary(newPath, currentPath string) error {
	oldPath := currentPath + ".old"
	_ = os.Remove(oldPath)

	err := os.Rename(currentPath, oldPath)
	if err != nil {
		return err
	}

	err = os.Rename(newPath, currentPath)
	if err != nil {
		_ = os.Rename(oldPath, currentPath)
		return err
	}

	if runtime.GOOS != "windows" {
		_ = os.Remove(oldPath)
	}

	return nil
}