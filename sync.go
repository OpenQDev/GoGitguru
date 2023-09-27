package main

import (
	"context"
	"main/internal/database"
	"main/internal/pkg/gitutil"
	"main/internal/pkg/logger"
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

	err := db.InsertRepoURL(context.Background(), "https://github.com/OpenQDev/OpenQ-Workflows")
	if err != nil {
		logger.LogFatalRedAndExit("Failed to insert repo url: %s", err)
	}

	for _, repoUrl := range repoUrls {
		repoUrl := repoUrl.Url

		organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrl)

		defer gitutil.DeleteLocalRepo(prefixPath, repo)

		gitutil.CloneRepo(prefixPath, organization, repo)

		gitutil.ProcessRepo(prefixPath, repo, repoUrl, db)
	}
}
