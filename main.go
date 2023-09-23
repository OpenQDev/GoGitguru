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
	organization := ""
	repo := ""
	if len(os.Args) > 2 {
		organization = os.Args[1]
		repo = os.Args[2]
	}

	// At the end of function execution, delete the repo and .tar.gz repo from the local filesystem
	// The great thing with this defer is it will run regardless of the outcomes of subsequent subprocesses
	defer os.RemoveAll(repo)
	defer os.RemoveAll(repo + ".tar.gz")

	// Clone the git repo
	fmt.Printf("\033[94mCloning https://github.com/%s/%s.git...\033[0m\n", organization, repo)
	err = cloneRepo(repo, organization)
	if err != nil {
		log.Fatal("\033[91mFailed to clone: ", err, "\033[0m")
	}
	fmt.Printf("\033[32m%s/%s successfully cloned!\033[0m\n", organization, repo)

	// Upload the repo to S3
	fmt.Printf("\033[94mUploading %s/%s to S3...\033[0m", organization, repo)
	fmt.Println()
	err = uploadRepoToS3(organization, repo)
	if err != nil {
		log.Fatalf("\033[91mFailed to upload to S3: %s\033[0m", err)
	}
	fmt.Printf("\033[32m%s/%s uploaded to S3!\033[0m\n", organization, repo)
}

func cloneRepo(repo string, organization string) error {
	cmd := exec.Command("git", "clone", fmt.Sprintf("https://github.com/%s/%s.git", organization, repo))

	// This allows you to see the stdout and stderr of the command being run on the host machine
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}

func uploadRepoToS3(organization string, repo string) error {
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
		return err
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Create a tarball of the .git directory
	err = exec.Command("tar", "-czf", repo+".tar.gz", repo+"/.git").Run()
	if err != nil {
		return err
	}

	// Open the tarball
	tarball, err := os.Open(repo + ".tar.gz")
	if err != nil {
		return err
	}

	// Upload the file to S3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("openqrepos"),
		Key:    aws.String(fmt.Sprintf("%s/%s.tar.gz", organization, repo)),
		Body:   tarball,
	})
	return err
}
