package reposync

import "slices"

func GetFirstAndLastUserCommitByEmail(usersToReposObject UsersToRepoUrl, emails []string, firstCommitDate int64, lastCommitDate int64) (int64, int64) {
	for index, userEmail := range usersToReposObject.AuthorEmails {
		if slices.Contains(emails, userEmail) {

			resultFirstCommitDate := firstCommitDate
			resultLastCommitDate := lastCommitDate

			if firstCommitDate > usersToReposObject.FirstCommitDates[index] || firstCommitDate == 0 {
				resultFirstCommitDate = usersToReposObject.FirstCommitDates[index]
			}
			if lastCommitDate < usersToReposObject.LastCommitDates[index] || lastCommitDate == 0 {
				resultLastCommitDate = usersToReposObject.LastCommitDates[index]
			}
			return resultFirstCommitDate, resultLastCommitDate
		}
	}
	return firstCommitDate, lastCommitDate
}
