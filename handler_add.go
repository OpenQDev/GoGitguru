package main

import (
	"encoding/json"
	"main/internal/database"
	"main/internal/pkg/logger"
	"net/http"
	"strings"
)

type Response struct {
	Accepted       []string `json:"accepted"`
	AlreadyInQueue []string `json:"already_in_queue"`
}

func (apiCfg *apiConfig) addHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		RepoUrls []string `json:"repo_urls"`
	}

	decoder := json.NewDecoder(r.Body)

	repoUrls := parameters{}

	err := decoder.Decode(&repoUrls)
	if err != nil {
		respondWithError(w, 400, "Error parsin JSON")
	}

	accepted := []string{}
	alreadyInQueue := []string{}

	for _, repoUrl := range repoUrls.RepoUrls {
		if !isListed(repoUrl, w, r, apiCfg) {
			err := apiCfg.DB.InsertRepoURL(r.Context(), database.InsertRepoURLParams{
				Url: repoUrl,
			})
			if err != nil {
				// handle error
				http.Error(w, err.Error(), http.StatusInternalServerError)
				logger.LogFatalRedAndExit("Error inserting repo url", err)
				continue
			}
			accepted = append(accepted, repoUrl)
		} else {
			alreadyInQueue = append(alreadyInQueue, repoUrl)
		}
	}

	response := Response{
		Accepted:       accepted,
		AlreadyInQueue: alreadyInQueue,
	}

	json.NewEncoder(w).Encode(response)
}

func isListed(repoUrl string, w http.ResponseWriter, r *http.Request, apiCfg *apiConfig) bool {
	_, err := apiCfg.DB.GetRepoURL(r.Context(), repoUrl)

	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return false
		} else {
			// handle error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.LogError("Error checking for repo url", err)
		}
	}

	return true
}
