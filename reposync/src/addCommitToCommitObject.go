package reposync

import (
	"strings"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func AddCommitToCommitObject(c *object.Commit, commitObject *database.BulkInsertCommitsParams, commitCount int) {

	totalFilesChanged := 0
	c.Author.Email = strings.Trim(c.Author.Email, "\"")
	c.Author.Email = strings.Trim(c.Author.Email, ".")
	commitObject.Commithashes[commitCount] = c.Hash.String()
	commitObject.Authors[commitCount] = c.Author.Name
	commitObject.Authoremails[commitCount] = c.Author.Email
	commitObject.Authordates[commitCount] = int64(c.Author.When.Unix())
	commitObject.Committerdates[commitCount] = int64(c.Committer.When.Unix())
	commitObject.Messages[commitCount] = strings.TrimRight(c.Message, "\n")
	commitObject.Fileschanged[commitCount] = int32(totalFilesChanged)
	if commitCount != 0 && commitCount%100 == 0 {
		logger.LogGreenDebug("process %d commits for %s", commitCount)
	}

}
