package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DependencyHistoryBody struct {
	RepoUrl            string   `json:"repo_url"`
	FilePaths          []string `json:"files_paths"`
	DependencySearched string   `json:"dependency_searched"`
}

func (apiCfg *ApiConfig) HandlerDependencyHistory(w http.ResponseWriter, r *http.Request) {
	var body DependencyHistoryBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		RespondWithError(w, 202, err.Error())
		return
	}

	// Check if any of the fields are their zero value
	if body.RepoUrl == "" || len(body.FilePaths) == 0 || body.DependencySearched == "" {
		RespondWithError(w, 400, "All fields must be present in the request body.")
		return
	}

	fmt.Println(body)

	RespondWithJSON(w, 202, struct{}{})
}
