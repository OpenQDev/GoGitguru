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

	response := HandlerStatusRequest{}

	var body HandlerGithubUserCommitsRequest
	err := marshaller.ReaderToType(r.Body, &body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request: %s", err))
		return
	}

	fmt.Println("body", body)

	RespondWithJSON(w, 202, response)
}
