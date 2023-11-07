package server

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
)

type DependencyHistoryRequest struct {
	RepoUrl            string   `json:"repo_url"`
	FilePaths          []string `json:"files_paths"`
	DependencySearched string   `json:"dependency_searched"`
}

type DependencyHistoryResponse struct {
	DatesAdded   []int64 `json:"dates_added"`
	DatesRemoved []int64 `json:"dates_removed"`
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

	repoDir := filepath.Join(prefixPath, organization, repo)

	if !gitutil.IsGitRepository(prefixPath, organization, repo) {
		RespondWithError(w, http.StatusNotFound, fmt.Sprintf("directory %s is not a git repository", repoDir))
		return
	}

	// "package.json" -> ["util/package.json", "app/package.json"]
	allFilePaths, err := gitutil.GitDependencyFiles(repoDir, body.FilePaths[0])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to determine if file paths exist dependency-history: %s", err))
		return
	}

	datesAddedCommits, datesRemovedCommits, err := gitutil.GitDependencyHistory(repoDir, body.DependencySearched, allFilePaths)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting dependency history: %s", err))
		return
	}

	dependencyHistoryResponse = DependencyHistoryResponse{
		DatesAdded:   datesAddedCommits,
		DatesRemoved: datesRemovedCommits,
	}

	RespondWithJSON(w, 200, dependencyHistoryResponse)
}
