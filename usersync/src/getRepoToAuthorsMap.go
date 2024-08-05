package usersync

type RepoToAuthorCommitTuples struct {
	Repos map[string][]AuthorCommitTuple
}

type AuthorCommitTuple struct {
	Author     string
	CommitHash string
}

// Create a map with repoUrl as key and array of authors as value
func getRepoToAuthorsMap(repoAuthorCommit Message) RepoToAuthorCommitTuples {
	repoToAuthorCommitTuples := RepoToAuthorCommitTuples{Repos: make(map[string][]AuthorCommitTuple)}

	if repoAuthorCommit.Repo_URL != "" {
		authorCommitTuple := AuthorCommitTuple{
			Author:     repoAuthorCommit.Author_Email,
			CommitHash: repoAuthorCommit.CommitHash,
		}

		repoToAuthorCommitTuples.Repos[repoAuthorCommit.Repo_URL] = append(
			repoToAuthorCommitTuples.Repos[repoAuthorCommit.Repo_URL], authorCommitTuple,
		)
	}

	return repoToAuthorCommitTuples
}
