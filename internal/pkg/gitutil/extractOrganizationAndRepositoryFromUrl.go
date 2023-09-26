package gitutil

import (
	"main/internal/pkg/logger"
	"net/url"
	"strings"
)

func ExtractOrganizationAndRepositoryFromUrl(repoUrl string) (string, string) {
	parsedUrl, err := url.Parse(repoUrl)
	if err != nil {
		logger.LogError("Error parsing URL: %s", err)
		return "", ""
	}

	pathParts := strings.Split(parsedUrl.Path, "/")
	organization := pathParts[len(pathParts)-2]
	repo := pathParts[len(pathParts)-1]

	return organization, repo
}
