package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/internal/database"
	"main/internal/pkg/logger"
	"main/internal/pkg/util"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
		logger.LogFatalRedAndExit("error getting repo urls: ", err)
	}

	for _, repoUrl := range repoUrls {
		// If the item does not exist, clone the repository and upload the tarball to S3
		// Extract the organization and the repository from the github url
		urlParts := strings.Split(repoUrl.Url, "/")
		organization := urlParts[len(urlParts)-2]
		repo := urlParts[len(urlParts)-1]

		// Check if the item exists in S3
		exists, err := util.ItemExistsInS3(uploader.S3, "openqrepos", fmt.Sprintf("%s/%s", organization, repo))

		if err != nil {
			logger.LogError("error checking if item %s exists in S3: ", repoUrl, err)
		}

		if !exists {
			util.CloneRepoAndUploadTarballToS3(organization, repo)
		} else {
			// If the item exists, pull the latest changes and re-upload the tarball to S3
			cmd := exec.Command("git", "-C", filepath.Join(prefixPath, repoUrl.Url), "pull")
			cmd.Run()
			util.UploadTarballToS3(prefixPath, organization, repoUrl.Url, uploader)
		}
	}
}
