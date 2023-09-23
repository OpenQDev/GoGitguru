package s3util

import (
	"fmt"
	"main/internal/pkg/logger"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

func CompressAndUploadToS3(prefixPath string, organization string, repo string, uploader s3manageriface.UploaderAPI) error {
	tarPath := filepath.Join(prefixPath, repo+".tar.gz")

	path, err := CompressDirectory(tarPath, filepath.Join(prefixPath, repo+"/.git"))
	if err != nil {
		logger.LogError("error tarring and gzipping", err)
	}
	fmt.Println(path)

	// Open the tarball
	tarball, err := os.Open(tarPath)
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
