package server

import (
	"fmt"
	"net/http"
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
	DatesAdded   []time.Time `json:"dates_added"`
	DatesRemoved []time.Time `json:"dates_removed"`
}

func (apiCfg *ApiConfig) HandlerDependencyHistory(w http.ResponseWriter, r *http.Request) {
	var dependencyHistoryResponse DependencyHistoryResponse

	var body DependencyHistoryRequest
	err := marshaller.ReaderToType(r.Body, &body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request: %s", err))
		return
	}

	organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(body.RepoUrl)

	repoDir := fmt.Sprintf("./repos/%s/%s", organization, repo)

	datesAddedCommits, datesRemovedCommits, err := gitutil.GitDependencyHistory(repoDir, body.DependencySearched, body.FilePaths)
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
