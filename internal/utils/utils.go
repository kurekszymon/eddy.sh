package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/kurekszymon/eddy.sh/internal/logger"
	"github.com/kurekszymon/eddy.sh/internal/types"
)

func PromptConfirm(prompt string, error_message string, codes ...int) {
	logger.Prompt(prompt)
	logger.Prompt("Type [Y] or [y] to continue")

	var i string
	fmt.Scan(&i)

	if i != "Y" && i != "y" {
		if len(error_message) > 0 {

			logger.Error(error_message)
		}

		if len(codes) > 0 {
			os.Exit(codes[0])
		}
	}
}

func PrintInstallErrors(errors_group ...map[string]error) {
	for _, errors := range errors_group {
		if len(errors) > 0 {
			for toolName, err := range errors {
				message := fmt.Sprintf("Error installing %s: %v\n", toolName, err)
				logger.Error(message)
			}
		}
	}
}

func DetermineVersion(version string, repo types.GHRepo) (string, error) {
	if version == "latest" {
		ver, err := getLatestReleaseFromGithub(repo.Owner, repo.Name)
		if err != nil {
			return "", err
		}
		// Remove leading "v" if present (e.g., "v0.40.3" -> "0.40.3")
		ver = strings.TrimPrefix(ver, "v")
		return ver, nil
	}
	return version, nil
}

func getLatestReleaseFromGithub(owner string, repo string) (string, error) {
	repoURL := fmt.Sprintf("https://github.com/%s/%s", owner, repo)
	latestURL := repoURL + "/releases/latest"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(latestURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest release for %s", repoURL)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusSeeOther && resp.StatusCode != http.StatusTemporaryRedirect {
		return "", fmt.Errorf("failed to fetch latest release: unexpected status code: %d", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("failed to fetch latest release: no redirect location found")
	}

	parts := strings.Split(strings.Trim(location, "/"), "/")
	tag := parts[len(parts)-1]
	return tag, nil
}
