package main

import (
	"fmt"
	"log"
	"main/internal/pkg/util"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize go dotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// The prefix path in which subsequent file operations will occur
	prefixPath := "repos"

	// Define the repo name
	organization := ""
	repo := ""
	if len(os.Args) > 2 {
		organization = os.Args[1]
		repo = os.Args[2]
	}

	// At the end of function execution, delete the repo and .tar.gz repo from the local filesystem
	// The great thing with this defer is it will run regardless of the outcomes of subsequent subprocesses
	defer util.DeleteLocalRepoAndTarball(prefixPath, repo)

	// Clone the git repo
	fmt.Printf("\033[94mCloning https://github.com/%s/%s.git...\033[0m\n", organization, repo)
	err = util.CloneRepo(prefixPath, organization, repo)
	if err != nil {
		log.Fatal("\033[91mFailed to clone: ", err, "\033[0m")
	}
	fmt.Printf("\033[32m%s/%s successfully cloned!\033[0m\n", organization, repo)

	// Upload the repo to S3
	fmt.Printf("\033[94mUploading %s/%s to S3...\033[0m", organization, repo)
	fmt.Println()
	err = util.UploadTarballToS3(prefixPath, organization, repo)
	if err != nil {
		log.Fatalf("\033[91mFailed to upload to S3: %s\033[0m", err)
	}
	fmt.Printf("\033[32m%s/%s uploaded to S3!\033[0m\n", organization, repo)
}
