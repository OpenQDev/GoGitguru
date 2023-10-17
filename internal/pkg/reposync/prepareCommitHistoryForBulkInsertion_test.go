package reposync

import (
	"main/internal/pkg/gitutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrepareCommitHistoryForBulkInsertion(t *testing.T) {

	repo := "OpenQ-DRM-TestRepo"
	prefixPath := "mock"

	r, err := gitutil.OpenGitRepo(prefixPath, repo)
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
	assert.Equal(t, []string{"updates README", "initial commit"}, commitObject.Message)
	assert.Equal(t, []string{"https://github.com/OpenQDev/OpenQ-DRM-TestRepo", "https://github.com/OpenQDev/OpenQ-DRM-TestRepo"}, commitObject.RepoUrls)

	assert.Equal(t, 2, cap(commitObject.CommitHash))
	assert.Equal(t, 2, cap(commitObject.Author))
	assert.Equal(t, 2, cap(commitObject.AuthorEmail))
	assert.Equal(t, 2, cap(commitObject.Message))
	assert.Equal(t, 2, cap(commitObject.RepoUrls))

	assert.Equal(t, 2, len(commitObject.CommitHash))
	assert.Equal(t, 2, len(commitObject.Author))
	assert.Equal(t, 2, len(commitObject.AuthorEmail))
	assert.Equal(t, 2, len(commitObject.Message))
	assert.Equal(t, 2, len(commitObject.RepoUrls))
}
