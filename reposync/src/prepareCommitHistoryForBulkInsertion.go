package reposync

import (
	"io"
	"strings"

	"github.com/OpenQDev/GoGitguru/util/logger"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func PrepareCommitHistoryForBulkInsertion(numberOfCommits int, log object.CommitIter, params GitLogParams) (CommitObject, error) {
	var (
		commitHash    = make([]string, numberOfCommits)
		author        = make([]string, numberOfCommits)
		authorEmail   = make([]string, numberOfCommits)
		authorDate    = make([]int64, numberOfCommits)
		committerDate = make([]int64, numberOfCommits)
		message       = make([]string, numberOfCommits)
		insertions    = make([]int32, numberOfCommits)
		deletions     = make([]int32, numberOfCommits)
		filesChanged  = make([]int32, numberOfCommits)
		repoUrls      = make([]string, numberOfCommits)
	)

	commitCount := 0
	for {
		if commitCount >= numberOfCommits {
			break
		}

		commit, err := log.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return CommitObject{}, err
			}
		}

		// TODO: Git stats is run JIT and is extremely time consuming, as it diffs each patch to the prior.
		// TODO: Ignoring for now, but this is how it's done.

		// stats, err := commit.Stats()
		// if err != nil {
		// 	logger.LogFatalRedAndExit("error computing stats for repository %s: %s", params.repoUrl, err)
		// }
		// for _, stat := range stats {
		// 	totalInsertions += stat.Addition
		// 	totalDeletions += stat.Deletion
		// 	totalFilesChanged++
		// }

		totalFilesChanged := 0
		totalInsertions := 0
		totalDeletions := 0

		commit.Author.Email = strings.Trim(commit.Author.Email, "\"")
		commit.Author.Email = strings.Trim(commit.Author.Email, ".")

		commitHash[commitCount] = commit.Hash.String()
		author[commitCount] = commit.Author.Name
		authorEmail[commitCount] = commit.Author.Email
		authorDate[commitCount] = int64(commit.Author.When.Unix())
		committerDate[commitCount] = int64(commit.Committer.When.Unix())
		message[commitCount] = strings.TrimRight(commit.Message, "\n")
		insertions[commitCount] = int32(totalInsertions)
		deletions[commitCount] = int32(totalDeletions)
		filesChanged[commitCount] = int32(totalFilesChanged)
		repoUrls[commitCount] = params.repoUrl

		if commitCount != 0 && commitCount%100 == 0 {
			logger.LogGreenDebug("process %d commits for %s", commitCount, params.repoUrl)
		}
		commitCount++
	}

	commitObject := CommitObject{
		CommitHash:    commitHash,
		Author:        author,
		AuthorEmail:   authorEmail,
		AuthorDate:    authorDate,
		CommitterDate: committerDate,
		Message:       message,
		Insertions:    insertions,
		Deletions:     deletions,
		FilesChanged:  filesChanged,
		RepoUrls:      repoUrls,
	}

	return commitObject, nil
}
