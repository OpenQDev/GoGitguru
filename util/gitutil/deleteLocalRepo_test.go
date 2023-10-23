package gitutil

import (
	"os"
	"path/filepath"
	"testing"
	"util/testhelpers"
)

func TestDeleteLocalRepo(t *testing.T) {
	tests := DeleteLocalRepoTestCases()

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory
			tempDir, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Fatalf("Failed to create temp directory: %s", err)
			}

			// Create a test repo directory inside the temp directory
			testRepo := filepath.Join(tempDir, tt.organization, tt.repo)
			err = os.MkdirAll(testRepo, 0755)
			if err != nil {
				t.Fatalf("Failed to create test repo directory: %s", err)
			}

			// Call the function to test
			err = DeleteLocalRepo(tempDir, tt.organization, tt.repo)
			if err != nil {
				t.Fatalf("DeleteLocalRepoAndTarball failed: %s", err)
			}

			// Check if the test repo directory has been deleted
			if _, err := os.Stat(testRepo); !os.IsNotExist(err) {
				t.Fatalf("Test repo directory was not deleted")
			}
		})
	}
}
