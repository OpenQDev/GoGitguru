package server

import (
	"fmt"
	"net/http"

	"github.com/OpenQDev/GoGitguru/util/marshaller"
)

type HandlerStatusRequest struct {
	RepoUrls []string `json:"repo_urls"`
}

type HandlerStatusResponse struct {
	Url            string `json:"url"`
	Status         string `json:"status"`
	PendingAuthors int    `json:"pending_authors"`
}

func (apiCfg *ApiConfig) HandlerStatus(w http.ResponseWriter, r *http.Request) {

	response := []HandlerStatusResponse{}

	var body HandlerStatusRequest
	err := marshaller.ReaderToType(r.Body, &body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request: %s", err))
		return
	}

	repoStatuses, err := apiCfg.DB.GetReposStatus(r.Context(), body.RepoUrls)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error in GetReposStatus: %s", err))
		return
	}

	for _, repoStatus := range repoStatuses {
		response = append(response, HandlerStatusResponse{
			Url:            repoStatus.Url,
			Status:         string(repoStatus.Status),
			PendingAuthors: int(repoStatus.PendingAuthors),
		})
	}

	RespondWithJSON(w, 202, response)
}
