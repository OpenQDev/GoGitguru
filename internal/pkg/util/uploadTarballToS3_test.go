package util

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

type mockUploader struct {
	s3manageriface.UploaderAPI
	input *s3manager.UploadInput
}

// Implement the methods of the s3manageriface.UploaderAPI interface
func (m *mockUploader) Upload(input *s3manager.UploadInput, foo ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	fmt.Println(input)

	output := s3manager.UploadOutput{}
	return &output, nil
}

func TestUploadTarballToS3(t *testing.T) {
	// Setup
	organization := "testOrganization"
	repo := "testRepo"

	// Pre-populate a tempDir at testPrefixPath/testRepo with a README.md file
	prefixPath, _ := os.MkdirTemp("", "repos")
	gitDir := filepath.Join(prefixPath, repo, ".git")
	err := os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create a README.md file in the .git directory
	readme, err := os.Create(filepath.Join(gitDir, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	readme.Close()

	uploader := &mockUploader{}

	// Call the function
	err = UploadTarballToS3(prefixPath, organization, repo, uploader)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Add more assertions as needed
}
