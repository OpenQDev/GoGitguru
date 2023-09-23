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
	logBlue("Cloning https://github.com/%s/%s.git...", organization, repo)
	err = util.CloneRepo(prefixPath, organization, repo)

	if err != nil {
		logRed("failed to clone: %s", err)
	}
	logGreen("%s/%s successfully cloned!", organization, repo)

	// Upload the repo to S3
	fmt.Printf("\033[94mUploading %s/%s to S3...\033[0m", organization, repo)
	fmt.Println()
	err = util.UploadTarballToS3(prefixPath, organization, repo)
	if err != nil {
		logRed("failed to upload to S3: %s", err)
	}
	logGreen("%s/%s uploaded to S3!", organization, repo)
}

func logBlue(format string, a ...interface{}) {
	fmt.Printf("\033[94m"+format+"\033[0m", a...)
	fmt.Println()
}

func logGreen(format string, a ...interface{}) {
	fmt.Printf("\033[32m"+format+"\033[0m", a...)
	fmt.Println()
}

func logRed(format string, a ...interface{}) {
	log.Fatalf("\033[91m"+format, a)
}
