package server

import (
	"net/http"
)

type HandlerGetNextRepoUrlResponse struct {
	RepoUrl string `json:"repo_url"`
}

func (apiCfg *ApiConfig) HandlerGetNextRepoUrl(w http.ResponseWriter, r *http.Request) {
	response := HandlerGetNextRepoUrlResponse{}
	RespondWithJSON(w, 200, response)
}
