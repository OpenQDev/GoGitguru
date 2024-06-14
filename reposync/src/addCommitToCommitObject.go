package reposync

import (
	"strings"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func AddCommitToCommitObject(c *object.Commit, commitObject *database.BulkInsertCommitsParams, commitCount int) int {

	c.Author.Email = strings.Trim(c.Author.Email, "\"")
	c.Author.Email = strings.Trim(c.Author.Email, ".")
	commitObject.Commithashes = append(commitObject.Commithashes, c.Hash.String())
	commitObject.Authors = append(commitObject.Authors, c.Author.Name)
	commitObject.Authoremails = append(commitObject.Authoremails, c.Author.Email)
	commitObject.Authordates = append(commitObject.Authordates, int64(c.Author.When.Unix()))
	commitObject.Committerdates = append(commitObject.Committerdates, c.Committer.When.Unix())
	commitObject.Messages = append(commitObject.Messages, strings.TrimRight(c.Message, "\n"))
	commitObject.Fileschanged = append(commitObject.Fileschanged, 0)
	if commitCount != 0 && commitCount%100 == 0 {
		logger.LogGreenDebug("process %d commits", commitCount)
	}
	commitCount++
	return commitCount

}
