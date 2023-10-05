package gitutil

// git -C . log --reverse --pretty=format:"%H-;-%an-;-%ae-;-%at-;-%ct%n%s" --numstat --since=2020-01-01
// git -C . rev-parse --is-inside-work-tree

import (
	"fmt"
	"io"
	"main/internal/database"
	"main/internal/pkg/logger"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
)

type GitLog struct {
	CommitHash    string
	AuthorName    string
	AuthorEmail   string
	AuthorDate    int64
	CommitDate    int64
	CommitMessage string
	FilesChanged  int64
	Insertions    int64
	Deletions     int64
}

func StoreGitLogs(prefixPath string, repo string, repoUrl string, fromCommitDate string, db *database.Queries) (int, error) {
	fullRepoPath := filepath.Join(prefixPath, repo)

	defaultCommitStartDate := "2020-01-01"
	if fromCommitDate == "" {
		fromCommitDate = defaultCommitStartDate
	}

	// Parse the fromCommitDate string into a time.Time value
	startDate, err := time.Parse("2006-01-02", fromCommitDate)
	if err != nil {
		return 0, err
	}

	r, err := git.PlainOpen(fullRepoPath)
	if err != nil {
		return 0, fmt.Errorf("not a git repository")
	}

	// Get the HEAD reference
	ref, err := r.Head()
	if err != nil {
		return 0, err
	}

	// Get the commit history from the start date
	// Unfortunately go-git does NOT have a way to order from oldest to newest commit
	log, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &startDate})
	if err != nil {
		return 0, err
	}

	// Iterate through the commit history
	commitCount := 0
	var values []string
	var args []interface{}
	for {
		commit, err := log.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return 0, err
			}
		}

		stats, _ := commit.Stats()

		filesChanged := 0
		insertions := 0
		deletions := 0
		for _, stat := range stats {
			insertions += stat.Addition
			deletions += stat.Deletion
			filesChanged++
		}

		// Create a placeholder for each commit
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", commitCount*10+1, commitCount*10+2, commitCount*10+3, commitCount*10+4, commitCount*10+5, commitCount*10+6, commitCount*10+7, commitCount*10+8, commitCount*10+9, commitCount*10+10))

		// Append the actual values
		args = append(args, commit.Hash.String(), commit.Author.Name, commit.Author.Email, commit.Author.When.Unix(), commit.Committer.When.Unix(), strings.TrimRight(commit.Message, "\n"), int32(insertions), int32(deletions), int32(filesChanged), repoUrl)

		if commitCount%100 == 0 {
			logger.LogGreenDebug("stored %d commits for %s", commitCount, repoUrl)
		}
		commitCount++
	}

	err = BatchInsertCommits(args, values)
	if err != nil {
		return 0, err
	}

	return commitCount, nil
}

func BatchInsertCommits(args []interface{}, values []string) error {
	return nil
}
