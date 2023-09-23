package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

func UploadTarballToS3(prefixPath string, organization string, repo string, uploader s3manageriface.UploaderAPI) error {
	fmt.Println(prefixPath)
	// Create a tarball of the .git directory
	err := exec.Command("tar", "-czf", filepath.Join(prefixPath, repo+".tar.gz"), filepath.Join(prefixPath, repo+"/.git")).Run()
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
