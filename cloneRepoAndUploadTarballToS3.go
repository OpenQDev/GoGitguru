package main

import (
	"main/internal/pkg/logger"
	"main/internal/pkg/util"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

func cloneRepoAndUploadTarballToS3(organization string, repo string) {
	// Initialize go dotenv
	err := godotenv.Load()
	if err != nil {
		logger.LogRed("error loading .env file: ", err)
	}

	// The prefix path in which subsequent file operations will occur
	prefixPath := "repos"

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

	logger.LogBlue("initializing AWS session...")

	// Get AWS API key and secret from environment variables
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Create a session using SharedConfigEnable
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	},
	)
	if err != nil {
		logger.LogRed("error initializing AWS session:", err)
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the repo to S3
	logger.LogBlue("uploading %s/%s to S3...", organization, repo)
	err = util.UploadTarballToS3(prefixPath, organization, repo, uploader)
	if err != nil {
		logger.LogRed("failed to upload to S3: %s", err)
	}
	logger.LogGreen("%s/%s uploaded to S3!", organization, repo)
}
