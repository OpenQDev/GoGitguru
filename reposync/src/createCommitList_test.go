package reposync

import (
	"path/filepath"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestCreateCommitList(t *testing.T) {
	// BEFORE ALL
	prefixPath := "mock"

	// ARRANGE - TESTS
	tests := CreateCommitListTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelpers.CheckTestSkip(t, testhelpers.Targets(
				testhelpers.RUN_ALL_TESTS,
			), tt.name)

			// BEFORE EACH

			// ARRANGE - LOCAL
			repoDir := filepath.Join(prefixPath, tt.organization, tt.repo)
			// ACT
			commitList, err := CreateCommitList(repoDir)

			for index, resultCommit := range tt.expectedResult {
				hash := resultCommit.Hash.String()
				if commitList[index].Hash.String() != hash {
					t.Errorf("Expected %s, got %s", hash, commitList[index].Hash.String())
				}
				if commitList[index].Author.Name != resultCommit.Author.Name {
					t.Errorf("Expected %s, got %s", resultCommit.Author.Name, commitList[0].Author.Name)
				}
				if commitList[index].Author.Email != resultCommit.Author.Email {
					t.Errorf("Expected %s, got %s", resultCommit.Author.Email, commitList[0].Author.Email)
				}
			}

			if err != nil {
				t.Errorf("there was an error processing repo: %s", err)
			}

			// ASSERT

			assert.Nil(t, err)
		})
	}
}
