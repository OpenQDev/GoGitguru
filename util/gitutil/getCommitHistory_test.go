package gitutil

import (
	"io"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetCommitHistory(t *testing.T) {
	// ARRANGE - GLOBAL
	r, err := git.PlainOpen("./mock/openqdev/openq-drm-testrepo")
	if err != nil {
		t.Fatalf("Failed to open repo: %v", err)
	}

	// ARRANGE - TESTS
	tests := GetCommitHistoryTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				"INVALID",
			), tt.name)
			// ACT
			commitIter, err := GetCommitHistory(r, tt.startDate)

			// ASSERT
			// Check that the function did not return an error
			assert.NoError(t, err)

			commitCount := 0
			for {
				commit, err := commitIter.Next()
				if err != nil {
					if err == io.EOF {
						break
					} else {
						return CommitObject{}, err
					}
				}

				commitCount++
			}

		})
	}
}
