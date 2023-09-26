package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/internal/database"
	"main/internal/pkg/gitutil"
	"main/internal/pkg/logger"
	"main/internal/pkg/s3util"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

func startSyncing() {
	// Load environment variables
	godotenv.Load(".env")
	dbUrl := os.Getenv("DB_URL")
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	prefixPath := "repos"

	if dbUrl == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	// Initialize the database and S3 uploader
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("can't connect to DB:", err)
	}

	queries := database.New(conn)

	// Create a session using SharedConfigEnable
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	},
	)

	if err != nil {
		logger.LogError("error connecting AWS session:", err)
	}

	uploader := s3manager.NewUploader(sess)

	// Fetch all repository URLs
	repoUrls, err := queries.GetRepoURLs(context.Background())

	if err != nil {
		logger.LogFatalRedAndExit("error getting repo urls: %s ", err)
	}

	for _, repoUrl := range repoUrls {
		// If the item does not exist, clone the repository and upload the tarball to S3
		// Extract the organization and the repository from the github url
		organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrl)

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
			err = s3util.CompressAndUploadToS3(prefixPath, organization, repoUrl.Url, uploader)
			if err != nil {
				logger.LogError("error uploading tarball for %s to s3: %s", repoUrl, err)
			}
		} else {
			gitutil.CloneRepoAndUploadTarballToS3(uploader, prefixPath, organization, repo)
		}
	}
}
