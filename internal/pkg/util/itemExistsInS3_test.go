package util

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

func TestItemExistsInS3(t *testing.T) {
	godotenv.Load("../../../.env")
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Create a session using SharedConfigEnable
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	},
	)
	if err != nil {
		t.Fatalf("failed to create session, %v", err)
	}

	uploader := s3manager.NewUploader(sess)

	tests := []struct {
		name    string
		bucket  string
		key     string
		want    bool
		wantErr bool
	}{
		{
			name:    "Item exists",
			bucket:  "openqrepos",
			key:     "OpenQDev/OpenQ-Contracts.tar.gz",
			want:    true,
			wantErr: false,
		},
		{
			name:    "Item does not exist",
			bucket:  "openqrepos",
			key:     "non_existent_key",
			want:    false,
			wantErr: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ItemExistsInS3(uploader, tt.bucket, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ItemExistsInS3() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ItemExistsInS3() = %v, want %v", got, tt.want)
			}
		})
	}
}
