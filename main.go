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
		fmt.Println("Error loading .env file")
		log.Fatal(err)
	}

	// Clone the git repo

	defer os.RemoveAll("OpenQ-Workflows")
	defer os.RemoveAll("OpenQ-Workflows.tar.gz")

	_, err = exec.Command("git", "clone", "https://github.com/OpenQDev/OpenQ-Workflows.git").Output()
	if err != nil {
		fmt.Println("Failed to clone")
		log.Fatal(err)
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
		fmt.Println("Failed to start AWS session")
		log.Fatal(err)
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Create a tarball of the .git directory
	err = exec.Command("tar", "-czf", "OpenQ-Workflows.tar.gz", "OpenQ-Workflows/.git").Run()
	if err != nil {
		fmt.Println("Failed to create tarball")
		log.Fatal(err)
	}

	// Open the tarball
	tarball, err := os.Open("OpenQ-Workflows.tar.gz")
	if err != nil {
		fmt.Println("Failed to open tarball")
		log.Fatal(err)
	}

	// Upload the file to S3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("openqrepos"),
		Key:    aws.String("OpenQ-Workflows.tar.gz"),
		Body:   tarball,
	})
	if err != nil {
		fmt.Println("Failed to upload to S3")
		log.Fatal(err)
	}

	// The file is automatically closed by the uploader

}
