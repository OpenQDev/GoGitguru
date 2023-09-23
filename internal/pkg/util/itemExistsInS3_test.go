package util

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// Create a mock S3 client so the unit test can run in isolation, independent of environment
type mockS3Client struct {
	s3iface.S3API
}

// Mock just the Head object on it
func (m *mockS3Client) HeadObject(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	// Mock behavior based on input
	if *input.Key == "non_existent_key" {
		return nil, awserr.New(s3.ErrCodeNoSuchKey, "no such key", nil)
	}
	return &s3.HeadObjectOutput{}, nil
}

func TestItemExistsInS3(t *testing.T) {
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
	}

	myMock := mockS3Client{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			itemExists, err := ItemExistsInS3(&myMock, tt.bucket, tt.key)

			if err != nil && tt.wantErr == false {
				t.Errorf("ItemExistsInS3() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if itemExists != tt.want {
				t.Errorf("ItemExistsInS3() = %v, want %v", itemExists, tt.want)
			}

		})
	}
}
