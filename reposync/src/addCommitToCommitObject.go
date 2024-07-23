package reposync

import (
	"strings"
	"unicode/utf8"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func replaceInvalidUTF8(s string) string {
	if !utf8.ValidString(s) {
		s = strings.Map(func(r rune) rune {
			if !utf8.ValidRune(r) {
				return -1 // Replace invalid runes with the replacement character
			}
			return r
		}, s)
	}
	return s
}

func AddCommitToCommitObject(c *object.Commit, commitObject *database.BulkInsertCommitsParams, commitCount int) int {

	c.Author.Email = replaceInvalidUTF8(strings.Trim(c.Author.Email, "\""))
	c.Author.Email = replaceInvalidUTF8(strings.Trim(c.Author.Email, "."))
	commitObject.Commithashes = append(commitObject.Commithashes, replaceInvalidUTF8(c.Hash.String()))
	commitObject.Authors = append(commitObject.Authors, replaceInvalidUTF8(c.Author.Name))
	commitObject.Authoremails = append(commitObject.Authoremails, c.Author.Email)
	commitObject.Authordates = append(commitObject.Authordates, int64(c.Author.When.Unix()))
	commitObject.Committerdates = append(commitObject.Committerdates, c.Committer.When.Unix())
	// Check if the message is UTF-8 encoded

	commitObject.Messages = append(commitObject.Messages, replaceInvalidUTF8(strings.TrimRight(c.Message, "\n")))
	commitObject.Fileschanged = append(commitObject.Fileschanged, 0)
	if commitCount != 0 && commitCount%100 == 0 {
		logger.LogGreenDebug("process %d commits", commitCount)
	}
	commitCount++
	return commitCount

}
