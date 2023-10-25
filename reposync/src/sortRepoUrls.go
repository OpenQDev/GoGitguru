package reposync

import (
	"sort"
	"strings"

	"github.com/OpenQDev/GoGitguru/database"
)

func sortRepoUrls(repoUrlObjects []database.RepoUrl) []string {
	repoUrls := make([]string, len(repoUrlObjects))

	for i, repo := range repoUrlObjects {
		// since sort.Strings uses case-sensitive lexicographic ordering, we must lowercase
		repoUrls[i] = strings.ToLower(repo.Url)
	}

	sort.Strings(repoUrls)
	return repoUrls
}
