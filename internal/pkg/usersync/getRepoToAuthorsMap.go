package usersync

type RepoToAuthorCommitTuples struct {
	Repos map[string][]AuthorCommitTuple
}

type AuthorCommitTuple struct {
	Author     string
	CommitHash string
}

// Create a map with repoUrl as key and array of authors as value
func GetRepoToAuthorsMap(repoAuthorCommits []UserSync) RepoToAuthorCommitTuples {
	repoToAuthorCommitTuples := RepoToAuthorCommitTuples{}

	for _, repoAuthorCommit := range repoAuthorCommits {
		if repoAuthorCommit.Repo.NotNull {
			authorCommitTuple := AuthorCommitTuple{
				Author:     repoAuthorCommit.Author.Email,
				CommitHash: repoAuthorCommit.CommitHash,
			}

			repoToAuthorCommitTuples.Repos[repoAuthorCommit.Repo.URL] = append(
				repoToAuthorCommitTuples.Repos[repoAuthorCommit.Repo.URL], authorCommitTuple,
			)
		}
	}

	return repoToAuthorCommitTuples
}
