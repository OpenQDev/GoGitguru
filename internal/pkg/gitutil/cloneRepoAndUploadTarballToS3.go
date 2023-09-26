package gitutil

import (
	"main/internal/pkg/logger"
	"main/internal/pkg/s3util"

	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

func CloneRepoAndUploadTarballToS3(uploader s3manageriface.UploaderAPI, prefixPath string, organization string, repo string) {
	// Clone the git repo
	logger.LogBlue("cloning https://github.com/%s/%s.git...", organization, repo)
	err := CloneRepo(prefixPath, organization, repo)
	if err != nil {
		logger.LogFatalRedAndExit("failed to clone: %s", err)
	}
	logger.LogGreen("%s/%s successfully cloned!", organization, repo)

	logger.LogBlue("initializing AWS session...")

	// TAR, GZIP, and Upload the repo to S3
	logger.LogBlue("uploading %s/%s to S3...", organization, repo)
	err = s3util.CompressAndUploadToS3(prefixPath, organization, repo, uploader)
	if err != nil {
		logger.LogFatalRedAndExit("failed to upload to S3: %s", err)
	}
	logger.LogGreen("%s/%s uploaded to S3!", organization, repo)
}
