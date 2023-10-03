package gitutil

// git -C . log --reverse --pretty=format:"%H-;-%an-;-%ae-;-%at-;-%ct%n%s" --numstat --since=2020-01-01
// git -C . rev-parse --is-inside-work-tree

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"main/internal/database"
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

func GetGitLogs(prefixPath string, repo string, repoUrl string, fromCommitDate string, db *database.Queries) (int, error) {
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

		params := database.InsertCommitParams{
			CommitHash:    commit.Hash.String(),
			Author:        sql.NullString{String: commit.Author.Name, Valid: commit.Author.Name != ""},
			AuthorEmail:   sql.NullString{String: commit.Author.Email, Valid: commit.Author.Email != ""},
			AuthorDate:    sql.NullInt64{Int64: commit.Author.When.Unix(), Valid: commit.Author.When.Unix() != 0},
			CommitterDate: sql.NullInt64{Int64: commit.Committer.When.Unix(), Valid: commit.Committer.When.Unix() != 0},
			Message:       sql.NullString{String: strings.TrimRight(commit.Message, "\n"), Valid: strings.TrimRight(commit.Message, "\n") != ""},
			Insertions:    sql.NullInt32{Int32: int32(insertions), Valid: true},
			Deletions:     sql.NullInt32{Int32: int32(deletions), Valid: true},
			FilesChanged:  sql.NullInt32{Int32: int32(filesChanged), Valid: true},
			RepoUrl:       sql.NullString{String: repoUrl, Valid: repoUrl != ""},
		}

		_, err = db.InsertCommit(context.Background(), params)
		if err != nil {
			return 0, fmt.Errorf("error loading commit number %d for %s: %s", commitCount, repoUrl, err)
		}

		commitCount++
	}

	return commitCount, nil
}
