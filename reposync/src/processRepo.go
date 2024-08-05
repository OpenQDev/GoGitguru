package reposync

import (
	"context"
	"time"

	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/logger"
)

func ProcessRepo(prefixPath string, organization string, repo string, repoUrl string, startDate time.Time, db *database.Queries) ([]GithubUser, []string, error) {
	logger.LogGreenDebug("beginning to process %s", repoUrl)

	db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
		Status: database.RepoStatusSyncingRepo,
		Url:    repoUrl,
	})

	commitCount, commitObject, repoUrlsUpdated, err := StoreGitLogsAndDepsHistoryForRepo(GitLogParams{prefixPath, organization, repo, repoUrl, startDate, db})
	if err != nil {
		db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
			Status: database.RepoStatusFailed,
			Url:    repoUrl,
		})
		return []GithubUser{}, []string{}, err
	}

	if commitCount == 0 {
		logger.LogBlue("no new commits in repo %s", repoUrl)
	}

	// Create a map to store unique emails with their associated info
	uniqueEmails := make(map[string]GithubUser)

	// Iterate through commitObject to collect unique emails and their associated info
	for i, email := range commitObject.Authoremails {
		if _, exists := uniqueEmails[email]; !exists {
			uniqueEmails[email] = GithubUser{
				AuthorEmail: email,
				AuthorDate:  time.Unix(commitObject.Authordates[i], 0),
				RepoUrl:     commitObject.Repourl.String,
				CommitHash:  commitObject.Commithashes[i],
			}
		}
	}

	// Create an array of structs with unique emails and their associated info
	var emailList []GithubUser
	for _, user := range uniqueEmails {
		emailList = append(emailList, user)
	}

	uniqueReposUpdated := make(map[string]bool)
	for _, repoUrl := range repoUrlsUpdated {
		if _, exists := uniqueReposUpdated[repoUrl]; !exists {
			uniqueReposUpdated[repoUrl] = true
		}
	}
	var reposUpdated []string
	for repo := range uniqueReposUpdated {
		reposUpdated = append(reposUpdated, repo)
	}

	// Print emailList
	logger.LogBlue("Email List:")
	for _, user := range emailList {
		logger.LogBlue("Email: %s, Author Date: %s, Repo URL: %s, Commit Hash: %s", user.AuthorEmail, user.AuthorDate.Format(time.RFC3339), user.RepoUrl, user.CommitHash)
	}

	db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
		Status: database.RepoStatusSynced,
		Url:    repoUrl,
	})

	logger.LogBlue("Successfully stored %d commits for %s in the database.", commitCount, repoUrl)

	return emailList, reposUpdated, nil
}
