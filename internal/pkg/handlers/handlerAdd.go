package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/internal/database"
	"main/internal/pkg/logger"
	"main/internal/pkg/util"
	"net/http"
	"strings"
)

type AddResponse struct {
	Accepted       []string `json:"accepted"`
	AlreadyInQueue []string `json:"already_in_queue"`
}

func (apiCfg *ApiConfig) HandlerAdd(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		RepoUrls []string `json:"repo_urls"`
	}

	// Read off the JSON body to bodyBytes for use in error logging if needed
	bodyBytes, _ := io.ReadAll(r.Body)

	// Reset r.Body to the original content
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Now prepare to decode the r.Body
	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))

	// Make struct repoUrls to decode the body into
	repoUrls := parameters{}

	err := decoder.Decode(&repoUrls)
	if err != nil || len(repoUrls.RepoUrls) == 0 {
		msg := fmt.Sprintf("error parsing JSON for: %s", string(bodyBytes))
		logger.LogError(msg)
		util.RespondWithError(w, 400, msg)
		return
	}

	accepted := []string{}
	alreadyInQueue := []string{}

	for _, repoUrl := range repoUrls.RepoUrls {
		repoIsListed, err := isListed(repoUrl, w, r, apiCfg)

		if err != nil {
			msg := fmt.Sprintf("error checking if repo is listed: %s", err)
			logger.LogError(msg)
			util.RespondWithError(w, 500, msg)
			return
		}

		if repoIsListed {
			alreadyInQueue = append(alreadyInQueue, repoUrl)
		} else {
			err := addToList(apiCfg, r, repoUrl, w)
			if err != nil {
				msg := fmt.Sprintf("error adding %s to repo_urls: %s", repoUrl, err)
				logger.LogError(msg)
				util.RespondWithError(w, 500, msg)
				return
			}
			accepted = append(accepted, repoUrl)
		}
	}

	response := AddResponse{
		Accepted:       accepted,
		AlreadyInQueue: alreadyInQueue,
	}

	util.RespondWithJSON(w, 202, response)
}

func addToList(apiCfg *ApiConfig, r *http.Request, repoUrl string, w http.ResponseWriter) error {
	err := apiCfg.DB.InsertRepoURL(r.Context(), database.InsertRepoURLParams{
		Url: repoUrl,
	})

	return err
}

func isListed(repoUrl string, w http.ResponseWriter, r *http.Request, apiCfg *ApiConfig) (bool, error) {
	_, err := apiCfg.DB.GetRepoURL(r.Context(), repoUrl)

	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}
