package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/minio/selfupdate"
	"github.com/spf13/cobra"
	"github.com/stef16robbe/stamp/internal/ui"
)

const (
	repoOwner = "stef16robbe"
	repoName  = "stamp"
)

type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update stamp to the latest version",
	Long:  `Check for updates and download the latest version of stamp from GitHub releases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(ui.Muted("Checking for updates..."))

		release, err := getLatestRelease()
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}

		latestVersion := strings.TrimPrefix(release.TagName, "v")
		currentVersion := strings.TrimPrefix(Version, "v")

		if latestVersion == currentVersion {
			fmt.Println(ui.Success(fmt.Sprintf("Already up to date (%s)", Version)))
			return nil
		}

		if Version == "dev" {
			fmt.Println(ui.Warning("Running development build, skipping update"))
			return nil
		}

		assetName := getAssetName()
		var downloadURL string
		for _, asset := range release.Assets {
			if asset.Name == assetName {
				downloadURL = asset.BrowserDownloadURL
				break
			}
		}

		if downloadURL == "" {
			return fmt.Errorf("no binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
		}

		fmt.Printf("Updating %s -> %s\n", ui.Muted(Version), ui.Bold(release.TagName))

		if err := doUpdate(downloadURL); err != nil {
			return fmt.Errorf("failed to apply update: %w", err)
		}

		fmt.Println(ui.Success(fmt.Sprintf("Updated to %s", release.TagName)))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func getLatestRelease() (*githubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no releases found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func getAssetName() string {
	name := fmt.Sprintf("stamp_%s_%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return name
}

func doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	return selfupdate.Apply(resp.Body, selfupdate.Options{})
}
