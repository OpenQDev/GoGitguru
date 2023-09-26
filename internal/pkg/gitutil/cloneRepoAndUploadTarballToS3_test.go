package gitutil

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

type mockUploader struct {
	s3manageriface.UploaderAPI
	bucket *string
	key    *string
}

// Implement the methods of the s3manageriface.UploaderAPI interface
func (m *mockUploader) Upload(input *s3manager.UploadInput, foo ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	// Assign the bucket and key on our spyUploader to ensure that it's being called with the appropriate args
	m.bucket = input.Bucket
	m.key = input.Key

	output := s3manager.UploadOutput{}
	return &output, nil
}

func TestCloneRepoAndUploadTarballToS3(t *testing.T) {
	// Setup
	organization := "OpenQDev"
	repo := "OpenQ-Workflows"

	prefixPath, err := os.MkdirTemp("", "repos")
	if err != nil {
		t.Fatal(err)
	}

	defer DeleteLocalRepoAndTarball(prefixPath, repo)

	uploader := &mockUploader{}

	// Call the function
	CloneRepoAndUploadTarballToS3(uploader, prefixPath, organization, repo)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Add more assertions as needed
	if *uploader.bucket != "openqrepos" {
		t.Errorf("Expected uploader.bucket to be 'openqrepos', but it was %s", *uploader.bucket)
	}
	if *uploader.key != fmt.Sprintf("%s/%s.tar.gz", organization, repo) {
		t.Errorf("Expected uploader.key to be %s/%s.tar.gz, but it was %s", organization, repo, *uploader.key)
	}
}
