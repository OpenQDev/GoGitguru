package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
)

type HandlerFirstCommitRequest struct {
	RepoUrl string `json:"repo_url"`
	Login   string `json:"login"`
}

type HandlerFirstCommitResponse struct {
	Commit CommitWithAuthorInfo `json:"commit"`
}

func (apiCfg *ApiConfig) HandlerFirstCommit(w http.ResponseWriter, r *http.Request) {
	var body HandlerFirstCommitRequest
	err := marshaller.ReaderToType(r.Body, &body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request: %s", err))
		return
	}

	params := database.GetFirstCommitParams{
		Column1: body.RepoUrl,
		Column2: body.Login,
	}

	commit, err := apiCfg.DB.GetFirstCommit(context.Background(), params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request: %s", err))
		return
	}

	RespondWithJSON(w, 200, commit)
}
