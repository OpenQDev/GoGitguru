package gitutil

// git -C . log --reverse --pretty=format:"%H-;-%an-;-%ae-;-%at-;-%ct%n%s" --numstat --since=2020-01-01
// git -C . rev-parse --is-inside-work-tree

import (
	"fmt"
	"io"
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

func GetGitLogs(prefixPath string, repo string, fromCommitDate string) ([]GitLog, error) {
	fullRepoPath := filepath.Join(prefixPath, repo)

	defaultCommitStartDate := "2020-01-01"
	if fromCommitDate == "" {
		fromCommitDate = defaultCommitStartDate
	}

	// Parse the fromCommitDate string into a time.Time value
	startDate, err := time.Parse("2006-01-02", fromCommitDate)
	if err != nil {
		return nil, err
	}

	r, err := git.PlainOpen(fullRepoPath)
	if err != nil {
		return nil, fmt.Errorf("not a git repository")
	}

	// Get the HEAD reference
	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	// Get the commit history from the start date
	// Unfortunately go-git does NOT have a way to order from oldest to newest commit
	log, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &startDate})
	if err != nil {
		return nil, err
	}

	// Iterate through the commit history
	gitLogs := make([]GitLog, 0)
	for {
		commit, err := log.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
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

		gitLog := GitLog{
			CommitHash:    commit.Hash.String(),
			AuthorName:    commit.Author.Name,
			AuthorEmail:   commit.Author.Email,
			AuthorDate:    commit.Author.When.Unix(),
			CommitDate:    commit.Committer.When.Unix(),
			CommitMessage: strings.TrimRight(commit.Message, "\n"),
			Insertions:    int64(insertions),
			Deletions:     int64(deletions),
			FilesChanged:  int64(filesChanged),
		}

		gitLogs = append(gitLogs, gitLog)
	}

	return gitLogs, nil
}
