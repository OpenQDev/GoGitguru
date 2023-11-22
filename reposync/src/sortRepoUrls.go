package reposync

import (
	"sort"
	"strings"

	"github.com/OpenQDev/GoGitguru/database"
)

func sortRepoUrls(repoUrlObjects []database.RepoUrl) []string {
	pendingUrls := make([]string, 0)
	otherUrls := make([]string, 0)

	for _, repo := range repoUrlObjects {
		if repo.Status == database.RepoStatusPending {
			pendingUrls = append(pendingUrls, strings.ToLower(repo.Url))
		} else {
			otherUrls = append(otherUrls, strings.ToLower(repo.Url))
		}
	}

	sort.Strings(pendingUrls)
	sort.Strings(otherUrls)

	return append(pendingUrls, otherUrls...)
}
