package reposync

import (
	"slices"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func AddFirstLastCommitDateByEmail(usersToRepoUrl *UsersToRepoUrl, commit *object.Commit) {
	alreadyHasEmail := slices.Contains(usersToRepoUrl.AuthorEmails, commit.Author.Email)
	if alreadyHasEmail {

		for index, email := range usersToRepoUrl.AuthorEmails {
			if email == commit.Author.Email {

				if commit.Author.When.Unix() < usersToRepoUrl.FirstCommitDates[index] {
					usersToRepoUrl.FirstCommitDates[index] = commit.Author.When.Unix()
				}
				if commit.Author.When.Unix() > usersToRepoUrl.LastCommitDates[index] {
					usersToRepoUrl.LastCommitDates[index] = commit.Author.When.Unix()
				}
				break
			}
		}
	}
	if !alreadyHasEmail {
		usersToRepoUrl.AuthorEmails = append(usersToRepoUrl.AuthorEmails, commit.Author.Email)
		usersToRepoUrl.FirstCommitDates = append(usersToRepoUrl.FirstCommitDates, commit.Author.When.Unix())
		usersToRepoUrl.LastCommitDates = append(usersToRepoUrl.LastCommitDates, commit.Author.When.Unix())
	}

}
