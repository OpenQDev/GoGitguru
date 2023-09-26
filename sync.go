package main

import (
	"fmt"
	"main/internal/database"
	"main/internal/pkg/gitutil"
	"main/internal/pkg/logger"
	"main/internal/pkg/s3util"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

func startSyncing(
	uploader *s3manageriface.UploaderAPI,
	database *database.Queries,
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
		// If the item does not exist, clone the repository and upload the tarball to S3
		// Extract the organization and the repository from the github url
		organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrl.Url)

		// At the end of function execution, delete the repo and .tar.gz repo from the local filesystem
		// The great thing with this defer is it will run regardless of the outcomes of subsequent subprocesses
		defer gitutil.DeleteLocalRepoAndTarball(prefixPath, repo)

		// Check if the item exists in S3
		item := fmt.Sprintf("%s/%s.tar.gz", organization, repo)
		exists, err := s3util.ItemExistsInS3(uploader.S3, "openqrepos", item)

		if err != nil {
			logger.LogError("error checking if item %s for repository %s exists in S3: ", item, repoUrl, err)
		}

		if exists {
			// If the item exists, pull the latest changes and re-upload the tarball to S3
			cmd := exec.Command("git", "-C", filepath.Join(prefixPath, repoUrl.Url), "pull")
			cmd.Run()
			err = s3util.CompressAndUploadToS3(prefixPath, organization, repo, *uploader)
			if err != nil {
				logger.LogError("error uploading tarball for %s to s3: %s", repoUrl, err)
			}
		} else {
			gitutil.CloneRepoAndUploadTarballToS3(*uploader, prefixPath, organization, repo)
		}
	}
}
