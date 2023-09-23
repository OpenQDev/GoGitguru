package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadTarballToS3(prefixPath string, organization string, repo string) error {
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
	err = exec.Command("tar", "-czf", filepath.Join(prefixPath, repo+".tar.gz"), filepath.Join(prefixPath, repo+"/.git")).Run()
	if err != nil {
		return err
	}

	// Open the tarball
	tarball, err := os.Open(filepath.Join(prefixPath, repo+".tar.gz"))
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
