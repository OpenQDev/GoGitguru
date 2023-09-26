package main

import (
	"main/internal/database"
	"main/internal/pkg/gitutil"
	"main/internal/pkg/logger"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func startSyncing(
	downloader *s3manager.Downloader,
	uploader *s3manager.Uploader,
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

		// process git log

		// delete repo
		err := os.RemoveAll(filepath.Join(prefixPath, repo))
		if err != nil {
			logger.LogError("failed to delete repo")
		}
	}
}
