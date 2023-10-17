package reposync

import (
	"bytes"
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
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	DateFormat       = "2006-01-02"
	DefaultStartDate = "2020-01-01"
)

type GitLogParams struct {
	prefixPath     string
	repo           string
	repoUrl        string
	fromCommitDate string
	db             *database.Queries
}

func StoreGitLogs(params GitLogParams) (int, error) {
	startDate, err := parseDate(params.fromCommitDate)
	if err != nil {
		return 0, err
	}

	r, err := openGitRepo(params.prefixPath, params.repo)
	if err != nil {
		return 0, err
	}

	log, err := getCommitHistory(r, startDate)
	if err != nil {
		return 0, err
	}

	numberOfCommits, err := getNumberOfCommits(params.prefixPath, params.repo)
	if err != nil {
		return 0, err
	}

	fmt.Printf("%s has %d commits\n", params.repoUrl, numberOfCommits)

	commitCount, err := processCommitHistory(numberOfCommits, log, params)
	if err != nil {
		return 0, err
	}

	return commitCount, nil
}

func processCommitHistory(numberOfCommits int, log object.CommitIter, params GitLogParams) (int, error) {
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
		repoUrls[commitCount] = params.repoUrl

		if commitCount%100 == 0 {
			logger.LogGreenDebug("process %d commits for %s", commitCount, params.repoUrl)
		}
		commitCount++
	}

	err := BatchInsertCommits(
		params.db,
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
		return 0, fmt.Errorf("error storing commits for %s: %s", params.repoUrl, err)
	}
	return commitCount, nil
}

func getNumberOfCommits(prefixPath string, repo string) (int, error) {
	fullRepoPath := filepath.Join(prefixPath, repo)
	cmd := exec.Command("git", "-C", fullRepoPath, "rev-list", "--count", "HEAD")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return 0, err
	}

	return strconv.Atoi(strings.TrimSpace(out.String()))
}

func getCommitHistory(r *git.Repository, startDate time.Time) (object.CommitIter, error) {
	ref, err := r.Head()
	if err != nil {
		return nil, err
	}
	return r.Log(&git.LogOptions{From: ref.Hash(), Since: &startDate})
}

func openGitRepo(prefixPath string, repo string) (*git.Repository, error) {
	fullRepoPath := filepath.Join(prefixPath, repo)
	return git.PlainOpen(fullRepoPath)
}

func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		dateStr = DefaultStartDate
	}
	return time.Parse(DateFormat, dateStr)
}
