package reposync

import (
	"testing"
	"time"

	"github.com/OpenQDev/GoGitguru/util/gitutil"

	"github.com/stretchr/testify/assert"
)

func TestPrepareCommitHistoryForBulkInsertion(t *testing.T) {
	// ARRANGE - GLOBAL

	organization := "openqdev"
	repo := "openq-drm-testrepo"
	prefixPath := "mock"

	r, err := gitutil.OpenGitRepo(prefixPath, organization, repo)
	if err != nil {
		t.Fatalf("failed to open git repo: %s", err)
	}

	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	log, err := gitutil.GetCommitHistory(r, startDate)
	if err != nil {
		t.Fatalf("failed to get commit history: %s", err)
	}

	params := GitLogParams{
		repoUrl: "https://github.com/OpenQDev/OpenQ-DRM-TestRepo",
	}

	// ACT
	commitObject, err := PrepareCommitHistoryForBulkInsertion(2, log, params)

	// ASSERT
	assert.NoError(t, err)

	assert.Equal(t, []string{"06a12f9c203112a149707ff73e4298749744c358", "9fae86bc8e89895b961d81bd7e9e4e897501c8bb"}, commitObject.CommitHash)
	assert.Equal(t, []string{"FlacoJones", "FlacoJones"}, commitObject.Author)
	assert.Equal(t, []string{"andrew@openq.dev", "andrew@openq.dev"}, commitObject.AuthorEmail)
	assert.Equal(t, []int64{1696277247, 1696277205}, commitObject.AuthorDate)
	assert.Equal(t, []int64{1696277247, 1696277205}, commitObject.CommitterDate)
	assert.Equal(t, []string{"updates README", "initial commit"}, commitObject.Message)
	assert.Equal(t, []int32{0, 0}, commitObject.Insertions)
	assert.Equal(t, []int32{0, 0}, commitObject.Deletions)
	assert.Equal(t, []int32{0, 0}, commitObject.FilesChanged)
	assert.Equal(t, []string{"https://github.com/OpenQDev/OpenQ-DRM-TestRepo", "https://github.com/OpenQDev/OpenQ-DRM-TestRepo"}, commitObject.RepoUrls)

	// cap (size of array) and len (number of elements in array) must match - or else pq.Array will attempt to insert cap - len # of empty commits
	assert.Equal(t, cap(commitObject.CommitHash), len(commitObject.CommitHash))
	assert.Equal(t, cap(commitObject.Author), len(commitObject.Author))
	assert.Equal(t, cap(commitObject.AuthorEmail), len(commitObject.AuthorEmail))
	assert.Equal(t, cap(commitObject.AuthorDate), len(commitObject.AuthorDate))
	assert.Equal(t, cap(commitObject.CommitterDate), len(commitObject.CommitterDate))
	assert.Equal(t, cap(commitObject.Insertions), len(commitObject.Insertions))
	assert.Equal(t, cap(commitObject.Deletions), len(commitObject.Deletions))
	assert.Equal(t, cap(commitObject.FilesChanged), len(commitObject.FilesChanged))
	assert.Equal(t, cap(commitObject.Message), len(commitObject.Message))
	assert.Equal(t, cap(commitObject.RepoUrls), len(commitObject.RepoUrls))
}
