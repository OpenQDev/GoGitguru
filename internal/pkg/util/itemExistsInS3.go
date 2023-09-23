package util

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func ItemExistsInS3(uploader *s3manager.Uploader, bucket, key string) (bool, error) {
	_, err := uploader.S3.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
