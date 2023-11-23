package server

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

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

	database, _ := setup.GetDatbase(env.DbUrl)
	repoUrlObjects, err := database.GetGithubReposByBatch(context.Background(), 9223372036854775807,
		0,
	)
	repoUrlsWithDependency := make([]string, 0)

	for _, repoUrlObject := range repoUrlObjects {
		prefixPath := "./repos"
		organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(repoUrlObject.Url)

		organization = strings.ToLower(organization)
		repo = strings.ToLower(repo)

		repoDir := filepath.Join(prefixPath, organization, repo)

		if !gitutil.IsGitRepository(prefixPath, organization, repo) {
			RespondWithError(w, http.StatusNotFound, fmt.Sprintf("directory %s is not a git repository", repoDir))
			return
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
		// Convert Unix time to ISO strings
		var lastDateAdded, lastDateRemoved int64 = 0, 0
		var lastDateAddedIso, lastDateRemovedIso string
		datesAddedISO := make([]string, len(datesAddedCommits))
		for i, v := range datesAddedCommits {
			if v >= lastDateAdded {
				lastDateAdded = v
				lastDateAddedIso = time.Unix(v, 0).Format(time.RFC3339)
			}
			datesAddedISO[i] = time.Unix(v, 0).Format(time.RFC3339)
		}

		datesRemovedISO := make([]string, len(datesRemovedCommits))
		for i, v := range datesRemovedCommits {
			if v >= lastDateRemoved {
				lastDateRemoved = v
				lastDateRemovedIso = time.Unix(v, 0).Format(time.RFC3339)
			}

			datesRemovedISO[i] = time.Unix(v, 0).Format(time.RFC3339)
		}
		fmt.Println(lastDateAddedIso, lastDateRemovedIso)
		hasDependency := lastDateAdded > lastDateRemoved || (lastDateAdded != 0 && lastDateRemoved == 0)
		if hasDependency {
			repoUrlsWithDependency = append(repoUrlsWithDependency, repoUrlObject.Url)
		}

	}

	fmt.Println("there")
	repoDependenciesResponse := DependencyHistoryByDependencyResponse{
		RepoUrls: repoUrlsWithDependency,
	}
	RespondWithJSON(w, 200, repoDependenciesResponse)
}
