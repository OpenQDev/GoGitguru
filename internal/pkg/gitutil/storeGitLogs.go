package gitutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"main/internal/database"
	"main/internal/pkg/logger"
	"os/exec"
	"path/filepath"
	"strconv"
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

	// Execute the command to get the number of commits
	cmd := exec.Command("git", "-C", fullRepoPath, "rev-list", "--count", "HEAD")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return 0, err
	}

	numberOfCommits, err := strconv.Atoi(strings.TrimSpace(out.String()))
	if err != nil {
		return 0, err
	}

	fmt.Printf("%s has %d commits\n", repoUrl, numberOfCommits)

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

		totalFilesChanged := 0
		totalInsertions := 0
		totalDeletions := 0
		for _, stat := range stats {
			totalInsertions += stat.Addition
			totalDeletions += stat.Deletion
			totalFilesChanged++
		}

		commitHash[commitCount] = commit.Hash.String()
		author[commitCount] = commit.Author.Name
		authorEmail[commitCount] = commit.Author.Email
		authorDate[commitCount] = int64(commit.Author.When.Unix())
		committerDate[commitCount] = int64(commit.Committer.When.Unix())
		message[commitCount] = strings.TrimRight(commit.Message, "\n")
		insertions[commitCount] = int32(totalInsertions)
		deletions[commitCount] = int32(totalDeletions)
		filesChanged[commitCount] = int32(totalFilesChanged)
		repoUrls[commitCount] = repoUrl

		if commitCount%100 == 0 {
			logger.LogGreenDebug("process %d commits for %s", commitCount, repoUrl)
		}
		commitCount++
	}

	err = BatchInsertCommits(
		db,
		commitHash,
		author,
		authorEmail,
		authorDate,
		committerDate,
		message,
		insertions,
		deletions,
		filesChanged,
		repoUrls,
	)
	if err != nil {
		return 0, fmt.Errorf("error storing commits for %s: %s", repoUrl, err)
	}

	return commitCount, nil
}

func BatchInsertCommits(
	db *database.Queries,
	commitHash []string,
	author []string,
	authorEmail []string,
	authorDate []int64,
	committerDate []int64,
	message []string,
	insertions []int32,
	deletions []int32,
	filesChanged []int32,
	repoUrls []string,
) error {
	params := database.BulkInsertCommitsParams{
		Column1:  commitHash,
		Column2:  author,
		Column3:  authorEmail,
		Column4:  authorDate,
		Column5:  committerDate,
		Column6:  message,
		Column7:  insertions,
		Column8:  deletions,
		Column9:  filesChanged,
		Column10: repoUrls,
	}
	return db.BulkInsertCommits(context.Background(), params)
}
