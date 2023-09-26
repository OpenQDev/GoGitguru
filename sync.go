package main

import (
	"fmt"
	"main/internal/database"
	"main/internal/pkg/gitutil"
	"main/internal/pkg/logger"
	"main/internal/pkg/s3util"
	"os"
	"os/exec"
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
		// If the item does not exist, clone the repository and upload the tarball to S3
		// Extract the organization and the repository from the github url
		organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrl.Url)

		fmt.Println(organization, repo)

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
			logger.LogBlue("%s/%s found in S3. downloading tar and pulling latest changes", organization, repo)

			// pull down the tar.gz repo from S3 to local filesystem
			err = s3util.DownloadFromS3(downloader, "openqrepos", item)
			if err != nil {
				logger.LogFatalRedAndExit("error downloading %s from S3: %s", item, err)
			}

			// Untar the .tar.gz located at repos/repo.tar.gz
			err = s3util.DecompressDirectory(item, filepath.Join(prefixPath, repo))
			if err != nil {
				logger.LogError("error untarring %s: %s", item, err)
			}

			// Pull the latest changes and re-upload the tarball to S3
			cmd := exec.Command("git", "-C", filepath.Join(prefixPath, repo), "pull")
			// Print any errors from running tar
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()

			if err != nil {
				logger.LogFatalRedAndExit("error pulling changes for %s: %s", repoUrl, err)
			}

			// Re-compress and upload to S3
			err = s3util.CompressAndUploadToS3(prefixPath, organization, repo, uploader)
			if err != nil {
				logger.LogError("error uploading tarball for %s to s3: %s", repoUrl, err)
			}
		} else {
			gitutil.CloneRepoAndUploadTarballToS3(uploader, prefixPath, organization, repo)
		}
	}
}
