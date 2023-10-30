package server

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
)

type HandlerStatusRequest struct {
	RepoUrls []string `json:"repo_urls"`
}

type HandlerStatusResponse struct {
	Url            string              `json:"url"`
	Status         database.RepoStatus `json:"status"`
	PendingAuthors int                 `json:"pending_authors"`
}

func (apiCfg *ApiConfig) HandlerStatus(w http.ResponseWriter, r *http.Request) {

	response := []HandlerStatusResponse{}

	var body HandlerStatusRequest
	err := marshaller.ReaderToType(r.Body, &body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request: %s", err))
		return
	}

	if len(body.RepoUrls) == 0 {
		RespondWithError(w, http.StatusBadRequest, "repo_urls cannot be empty")
		return
	}

	repoStatuses, err := apiCfg.DB.GetReposStatus(r.Context(), body.RepoUrls)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error in GetReposStatus: %s", err))
		return
	}

	for _, repoStatus := range repoStatuses {
		if slices.Contains(body.RepoUrls, repoStatus.Url) {
			response = append(response, HandlerStatusResponse{
				Url:            repoStatus.Url,
				Status:         repoStatus.Status,
				PendingAuthors: int(repoStatus.PendingAuthors),
			})
		} else {
			response = append(response, HandlerStatusResponse{
				Url:            repoStatus.Url,
				Status:         repoStatus.Status,
				PendingAuthors: int(repoStatus.PendingAuthors),
			})
		}
	}

	// Find the diff of urls that appear in body.RepoUrls but don't appear in repoStatuses.Url
	repoUrls := make(map[string]bool)
	for _, repoStatus := range repoStatuses {
		repoUrls[repoStatus.Url] = true
	}

	for _, url := range body.RepoUrls {
		if _, ok := repoUrls[url]; !ok {
			response = append(response, HandlerStatusResponse{
				Url:            url,
				Status:         database.RepoStatusNotListed,
				PendingAuthors: 0,
			})
		}
	}

	RespondWithJSON(w, 202, response)
}
