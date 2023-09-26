package gitutil

import (
	"fmt"
	"net/url"
	"strings"
)

func ExtractOrganizationAndRepositoryFromUrl(repoUrl string) (string, string) {
	parsedUrl, err := url.Parse(repoUrl)
	if err != nil {
		fmt.Println("Error parsing URL: ", err)
		return "", ""
	}

	pathParts := strings.Split(parsedUrl.Path, "/")
	organization := pathParts[len(pathParts)-2]
	repo := pathParts[len(pathParts)-1]

	return organization, repo
}
