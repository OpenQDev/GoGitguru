package main

import (
	"main/internal/database"
	"main/internal/pkg/gitutil"
	"time"
)

func startSyncing(
	db *database.Queries,
	prefixPath string,
	concurrency int,
	timeBetweenSyncs time.Duration) {
	// TODO Fetch RepoUrls from DB
	// Fetch all repository URLs
	// repoUrls, err := queries.GetRepoURLs(context.Background())

	// if err != nil {
	// 	logger.LogFatalRedAndExit("error getting repo urls: %s ", err)
	// }

	repoUrls := []database.RepoUrl{
		{
			Url: "https://github.com/OpenQDev/OpenQ-Workflows",
		},
	}

	for _, repoUrl := range repoUrls {
		organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrl.Url)

		gitutil.CloneRepo(prefixPath, organization, repo)
	}
}
