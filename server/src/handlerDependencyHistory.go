package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
)

type DependencyHistoryRequest struct {
	RepoUrl            string   `json:"repo_url"`
	FilePaths          []string `json:"files_paths"`
	DependencySearched string   `json:"dependency_searched"`
}

type DependencyHistoryResponse struct {
	DatesAdded   []string `json:"dates_added"`
	DatesRemoved []string `json:"dates_removed"`
}

func (apiCfg *ApiConfig) HandlerDependencyHistory(w http.ResponseWriter, r *http.Request) {
	var dependencyHistoryResponse DependencyHistoryResponse

	var body DependencyHistoryRequest
	err := marshaller.ReaderToType(r.Body, &body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request in dependency-history: %s", err))
		return
	}

	prefixPath := "./repos"
	organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(body.RepoUrl)

	organization = strings.ToLower(organization)
	repo = strings.ToLower(repo)

	repoDir := filepath.Join(prefixPath, organization, repo)

	fmt.Println("is git repo for ", repoDir, gitutil.IsGitRepository(prefixPath, organization, repo))

	if !gitutil.IsGitRepository(prefixPath, organization, repo) {
		err := gitutil.CloneRepo(prefixPath, organization, repo)
		if err != nil {
			RespondWithError(w, http.StatusNotFound, fmt.Sprintf("directory %s is not a git repository", repoDir))
			return
		}
	}

	// "package.json" -> ["util/package.json", "app/package.json"]
	allFilePaths, err := gitutil.GitDependencyFiles(repoDir, body.FilePaths)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to determine if file paths exist dependency-history: %s", err))
		return
	}
	
	dependencies := []string{"eslint", "ethers", "threejs", body.DependencySearched}

	datesAddedCommits, datesRemovedCommits, err := gitutil.GitDependencyHistory(repoDir, dependencies, allFilePaths)
	currentDatesAdded := datesAddedCommits[body.DependencySearched]
	currentDatesRemoved := datesRemovedCommits[body.DependencySearched]
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting dependency history: %s", err))
		return
	}

	// Convert Unix time to ISO strings
	datesAddedISO := make([]string, len(currentDatesAdded))
	for i, v := range currentDatesAdded {
		datesAddedISO[i] = time.Unix(v, 0).Format(time.RFC3339)
	}

	datesRemovedISO := make([]string, len(currentDatesRemoved))
	for i, v := range currentDatesRemoved {
		datesRemovedISO[i] = time.Unix(v, 0).Format(time.RFC3339)
	}

	dependencyHistoryResponse = DependencyHistoryResponse{
		DatesAdded:   datesAddedISO,
		DatesRemoved: datesRemovedISO,
	}

	RespondWithJSON(w, http.StatusOK, dependencyHistoryResponse)
}
