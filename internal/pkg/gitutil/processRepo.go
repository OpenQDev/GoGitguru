package gitutil

import "main/internal/database"

func ProcessRepo(prefixPath string, repo string, repoUrl string, db *database.Queries) {
	StoreCommits(GetFormattedGitLogs(prefixPath, repo, ""), repoUrl, db)
}
