package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize go dotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Define the repo name
	repo := "OpenQ-Workflows"
	organization := "OpenQDev"

	// At the end of function execution, delete the repo and .tar.gz repo from the local filesystem
	defer os.RemoveAll(repo)
	defer os.RemoveAll(repo + ".tar.gz")

	// Clone the git repo
	_, err = exec.Command("git", "clone", fmt.Sprintf("https://github.com/%s/%s.git", organization, repo)).Output()
	if err != nil {
		log.Fatal("Failed to clone: ", err)
	}

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
		log.Fatal("Failed to start AWS session:", err)
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Create a tarball of the .git directory
	err = exec.Command("tar", "-czf", repo+".tar.gz", repo+"/.git").Run()
	if err != nil {
		log.Fatal("Failed to create tarball:", err)
	}

	// Open the tarball
	tarball, err := os.Open(repo + ".tar.gz")
	if err != nil {
		log.Fatal("Failed to open tarball:", err)
	}

	// Upload the file to S3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("openqrepos"),
		Key:    aws.String(repo + ".tar.gz"),
		Body:   tarball,
	})
	if err != nil {
		log.Fatal("Failed to upload to S3:", err)
	}

	// The file is automatically closed by the uploader
}
