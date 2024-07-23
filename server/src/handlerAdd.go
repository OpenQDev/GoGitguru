package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/OpenQDev/GoGitguru/util/logger"
)

type HandlerAddRequest struct {
	RepoUrls []string `json:"repo_urls"`
}

type HandlerAddResponse struct {
	Accepted []string `json:"accepted"`
}

func (apiCfg *ApiConfig) HandlerAdd(w http.ResponseWriter, r *http.Request) {
	// Read off the JSON body to bodyBytes for use in error logging if needed
	bodyBytes, _ := io.ReadAll(r.Body)

	// Reset r.Body to the original content
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Now prepare to decode the r.Body
	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))

	// Make struct repoUrls to decode the body into
	request := HandlerAddRequest{}

	err := decoder.Decode(&request)
	if err != nil || len(request.RepoUrls) == 0 {
		msg := fmt.Sprintf("error parsing JSON for: %s", string(bodyBytes))
		RespondWithError(w, 400, msg)
		return
	}

	accepted := []string{}

	for _, repoUrl := range request.RepoUrls {

		err = addToList(apiCfg, r, repoUrl)
		if err != nil {
			msg := fmt.Sprintf("error adding %s to repo_urls: %s", repoUrl, err)
			logger.LogError(msg)
			RespondWithError(w, 500, msg)
			return
		}

		accepted = append(accepted, repoUrl)

	}

	response := HandlerAddResponse{
		Accepted: accepted,
	}

	RespondWithJSON(w, 202, response)
}

func addToList(apiCfg *ApiConfig, r *http.Request, repoUrl string) error {
	err := apiCfg.DB.UpsertRepoURL(r.Context(), repoUrl)

	return err
}
