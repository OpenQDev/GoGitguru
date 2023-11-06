package server

import (
	"database/sql"
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

type HandlerFirstCommitResponse struct {
	AuthorDate int `json:"author_date"`
}

func (apiCfg *ApiConfig) HandlerFirstCommit(w http.ResponseWriter, r *http.Request) {
	var body HandlerFirstCommitRequest
	err := marshaller.ReaderToType(r.Body, &body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request: %s", err))
		return
	}

	params := database.GetFirstCommitParams{
		Login:   body.Login,
		RepoUrl: sql.NullString{String: body.RepoUrl, Valid: true},
	}

	commit, err := apiCfg.DB.GetFirstCommit(r.Context(), params)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			RespondWithJSON(w, 200, nil)
		} else {
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request: %s", err))
		}
		return
	}

	response := HandlerFirstCommitResponse{AuthorDate: int(commit.AuthorDate.Int64)}

	RespondWithJSON(w, 200, response)
}
