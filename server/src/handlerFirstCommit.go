package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
)

type HandlerFirstCommitRequest struct {
	RepoUrl string `json:"repo_url"`
	Login   string `json:"login"`
}

func (apiCfg *ApiConfig) HandlerFirstCommit(w http.ResponseWriter, r *http.Request) {
	var body HandlerFirstCommitRequest
	err := marshaller.ReaderToType(r.Body, &body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request: %s", err))
		return
	}

	params := database.GetFirstCommitParams{
		Column1: strings.ToLower(body.RepoUrl),
		Column2: strings.ToLower(body.Login),
	}

	fmt.Println("params", params)

	commit, err := apiCfg.DB.GetFirstCommit(context.Background(), params)
	fmt.Println("commit", commit)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			RespondWithJSON(w, 200, nil)
		} else {
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request: %s", err))
		}
		return
	}

	RespondWithJSON(w, 200, commit)
}
