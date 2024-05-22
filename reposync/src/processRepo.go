package reposync

import (
	"context"
	"path/filepath"
	"time"

	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/OpenQDev/GoGitguru/util/logger"
)

func ProcessRepo(prefixPath string, organization string, repo string, repoUrl string, startDate time.Time, db *database.Queries) error {
	logger.LogGreenDebug("beginning to process %s", repoUrl)
	depFilePathsFE := []string{"package.json", ".config.", ".yaml", ".yml", "truffle", ".toml", "network", "hardhat", "deploy", "go.mod", "composer.json"}
	dependencies := []string{"dependencies", "devDependencies", "peerDependencies", "optionalDependencies"}
 
	db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
		Status: database.RepoStatusSyncingRepo,
		Url:    repoUrl,
	})



	repoDir := filepath.Join(prefixPath, organization, repo)
	gitutil.GitDependencyHistory(repoDir, dependencies, depFilePathsFE)
	commitCount, err := StoreGitLogsForRepo(GitLogParams{prefixPath, organization, repo, repoUrl, startDate, db})
	if err != nil {
		db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
			Status: database.RepoStatusFailed,
			Url:    repoUrl,
		})
		return err
	}

	if commitCount == 0 {
		logger.LogBlue("no new commits in repo %s", repoUrl)
	}

	db.UpdateStatusAndUpdatedAt(context.Background(), database.UpdateStatusAndUpdatedAtParams{
		Status: database.RepoStatusSynced,
		Url:    repoUrl,
	})

	logger.LogBlue("Successfully stored %d commits for %s in the database.", commitCount, repoUrl)

	return nil
}
