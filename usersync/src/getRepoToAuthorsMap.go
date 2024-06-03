package usersync

import "fmt"

type RepoToAuthorCommitTuples struct {
	Repos map[string][]AuthorCommitTuple
}

type AuthorCommitTuple struct {
	Author     string
	CommitHash string
}

// Create a map with repoUrl as key and array of authors as value
func getRepoToAuthorsMap(repoAuthorCommits []UserSync) RepoToAuthorCommitTuples {
	repoToAuthorCommitTuples := RepoToAuthorCommitTuples{Repos: make(map[string][]AuthorCommitTuple)}
	for _, repoAuthorCommit := range repoAuthorCommits {
		fmt.Println("repoAuthorCommits", repoAuthorCommit.RepoUrl)
		if repoAuthorCommit.RepoUrl != "" {
			authorCommitTuple := AuthorCommitTuple{
				Author:     repoAuthorCommit.AuthorEmail,
				CommitHash: repoAuthorCommit.CommitHash,
			}

			repoToAuthorCommitTuples.Repos[repoAuthorCommit.RepoUrl] = append(
				repoToAuthorCommitTuples.Repos[repoAuthorCommit.RepoUrl], authorCommitTuple,
			)
		}
	}

	return repoToAuthorCommitTuples
}
