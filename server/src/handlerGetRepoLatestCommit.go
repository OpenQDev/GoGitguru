package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/OpenQDev/GoGitguru/util/marshaller"
)

type RepoLatestCommit struct {
	RepoName          string    `json:"repo_name"`
	LastPushEventTime time.Time `json:"last_push_event_time"`
}

type RepoLatestCommitRequest struct {
	RepoName string `json:"repo_name"`
}

func (apiCfg *ApiConfig) HandlerGetRepoLatestCommit(w http.ResponseWriter, r *http.Request) {
	var body RepoLatestCommitRequest
	err := marshaller.ReaderToType(r.Body, &body)
	if err != nil {
		http.Error(w, "failed to read body of request", http.StatusBadRequest)
		return
	}

	repoName := strings.ToLower(body.RepoName)
	if repoName == "" {
		http.Error(w, "repo_name is required", http.StatusBadRequest)
		return
	}

	var repoLatestCommit RepoLatestCommit
	resp, err := apiCfg.DB.GetRepoLatestCommit(r.Context(), repoName)
	if err != nil {
		http.Error(w, "repo not found", http.StatusNotFound)
		return
	}

	RespondWithJSON(w, http.StatusOK, repoLatestCommit)
}
