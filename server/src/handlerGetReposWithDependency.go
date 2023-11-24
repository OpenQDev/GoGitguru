package server

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

type DependencyHistoryByDependencyRequest struct {
	DependencySearched string `json:"dependency_searched"`
}

type DependencyHistoryByDependencyResponse struct {
	RepoUrls []string `json:"repo_urls"`
}

var FilePaths = []string{"package.json",
	".config.",
	".yaml",
	".yml",
	"truffle",
	".toml",
	"network",
	"hardhat",
	"deploy",
	"go.mod",
	"composer.json"}

func (apiCfg *ApiConfig) HandlerDependencyHistoryByDependency(w http.ResponseWriter, r *http.Request) {
	var body DependencyHistoryByDependencyRequest
	err := marshaller.ReaderToType(r.Body, &body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request in dependency-history: %s", err))
		return
	}

	//prefixPath := "./repos"
	//organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(body.RepoUrl)

	//organization = strings.ToLower(organization)

	env := setup.ExtractAndVerifyEnvironment(".env")

	db, _ := setup.GetDatbase(env.DbUrl)
	take := int32(1000000)
	repoUrlObjects, err := db.GetGithubReposByBatch(context.Background(), database.GetGithubReposByBatchParams{Limit: take, Offset: 0})
	repoUrlsWithDependency := make([]string, 0)
	repoIdsWithDependencyToBeIndexed := make([]int32, 0)
	firstCommitDateValuesToBeIndexed := make([]int64, 0)
	dateAddedValuesToBeIndexed := make([]int64, 0)
	dateRemovedValuesToBeIndexed := make([]int64, 0)

	type RepoDependencyData struct {
		github_repo_id    int32
		dependency_name   string
		first_commit_date int64
		date_added        int64
		date_removed      int64
	}

	dependencyParams := database.InsertDependenciesParams{
		DependencyName:  body.DependencySearched,
		DependencyFiles: FilePaths,
	}
	res, err := db.InsertDependencies(context.Background(), dependencyParams)
	if err != nil && err.Error() != "sql: no rows in result set" {
		fmt.Println("failed to insert dependency", err)
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to insert dependency: %s", err))
		return
	}
	fmt.Println("inserted dependency", res)

	repoDependencyDataToBeIndexed := make([]RepoDependencyData, 0)
	indexedReposWithDependency, err := checkIndexedReposWithDependency(db, body.DependencySearched, repoUrlObjects)
	for _, repoUrlObject := range repoUrlObjects {
		// if any of the indexed repos have the same id as repoUrlObject, then skip
		hasIndexedRepoToDependency := getHasRepoToDependency(indexedReposWithDependency, repoUrlObject)
		if hasIndexedRepoToDependency {
			repoUrlsWithDependency = append(repoUrlsWithDependency, repoUrlObject.Url)
			continue
		}

		if err != nil {
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to check indexed repos with dependency: %s", err))
			return
		}

		prefixPath := "./repos"

		organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrlObject.Url)

		organization = strings.ToLower(organization)
		repo = strings.ToLower(repo)

		repoDir := filepath.Join(prefixPath, organization, repo)

		if !gitutil.IsGitRepository(prefixPath, organization, repo) {
			continue
		}

		allFilePaths, err := gitutil.GitDependencyFiles(repoDir, FilePaths)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to determine if file paths exist dependency-history: %s", err))
			return
		}

		datesAddedCommits, datesRemovedCommits, err := gitutil.GitDependencyHistory(repoDir, body.DependencySearched, allFilePaths)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting dependency history: %s", err))
			return
		}
		hasNonIndexedRepoToDependency, lastDateAdded, lastDateRemoved := checkNonIndexedReposWithDependency(db, body.DependencySearched, repoUrlObjects, datesAddedCommits, datesRemovedCommits)

		if hasNonIndexedRepoToDependency {
			repoUrlsWithDependency = append(repoUrlsWithDependency, repoUrlObject.Url)
			repoIdsWithDependencyToBeIndexed = append(repoIdsWithDependencyToBeIndexed, repoUrlObject.InternalID)
			repoDependencyDataToBeIndexed = append(repoDependencyDataToBeIndexed, RepoDependencyData{
				github_repo_id:    repoUrlObject.InternalID,
				dependency_name:   body.DependencySearched,
				first_commit_date: lastDateAdded,
				date_added:        lastDateAdded,
				date_removed:      lastDateRemoved,
			})
			firstCommitDateValuesToBeIndexed = append(firstCommitDateValuesToBeIndexed, lastDateAdded)
			dateAddedValuesToBeIndexed = append(dateAddedValuesToBeIndexed, lastDateAdded)
			dateRemovedValuesToBeIndexed = append(dateRemovedValuesToBeIndexed, lastDateRemoved)

		}

	}

	err = BulkInsertRepoToDependencies(db, repoIdsWithDependencyToBeIndexed, body.DependencySearched, firstCommitDateValuesToBeIndexed, dateAddedValuesToBeIndexed, dateRemovedValuesToBeIndexed)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting dependency history: %s", err))
		return
	}
	repoDependenciesResponse := DependencyHistoryByDependencyResponse{
		RepoUrls: repoUrlsWithDependency,
	}
	RespondWithJSON(w, 200, repoDependenciesResponse)
}
