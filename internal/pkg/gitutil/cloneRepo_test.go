package gitutil

import (
	"os"
	"testing"
)

func TestCloneRepo(t *testing.T) {
	prefixPath, _ := os.MkdirTemp("", "repos")

	tests := CloneRepoTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CloneRepo(prefixPath, tt.organization, tt.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloneRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Defer deletion of the repo after each test
			defer DeleteLocalRepo(prefixPath, tt.repo)
		})
	}
}
