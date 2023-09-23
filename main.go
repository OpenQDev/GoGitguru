package main

import (
	"main/internal/pkg/logger"
	"main/internal/pkg/util"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize go dotenv
	err := godotenv.Load()
	if err != nil {
		logger.LogRed("error loading .env file: ", err)
	}

	// The prefix path in which subsequent file operations will occur
	prefixPath := "repos"

	// Define the repo and organization name
	organization := ""
	repo := ""

	// If args for organization and repo are provided from the command line, use them
	if len(os.Args) > 2 {
		organization = os.Args[1]
		repo = os.Args[2]
	}

	// At the end of function execution, delete the repo and .tar.gz repo from the local filesystem
	// The great thing with this defer is it will run regardless of the outcomes of subsequent subprocesses
	defer util.DeleteLocalRepoAndTarball(prefixPath, repo)

	// Clone the git repo
	logger.LogBlue("cloning https://github.com/%s/%s.git...", organization, repo)
	err = util.CloneRepo(prefixPath, organization, repo)

	if err != nil {
		logger.LogRed("failed to clone: %s", err)
	}
	logger.LogGreen("%s/%s successfully cloned!", organization, repo)

	// Upload the repo to S3
	logger.LogBlue("uploading %s/%s to S3...", organization, repo)
	err = util.UploadTarballToS3(prefixPath, organization, repo)
	if err != nil {
		logger.LogRed("failed to upload to S3: %s", err)
	}
	logger.LogGreen("%s/%s uploaded to S3!", organization, repo)
}
