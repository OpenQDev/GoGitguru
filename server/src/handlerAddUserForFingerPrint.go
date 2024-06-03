package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/logger"
)

type HandlerAddUserFingerPrintRequest struct {
	Logins       []string `json:"logins"`
	Dependencies []string `json:"dependencies"`
	FileNames    []string `json:"file_names"`
}

type HandlerAddUserFingerPrintResponse struct {
	Accepted []string `json:"accepted"`
}

func (apiCfg *ApiConfig) HandlerAddUserFingerPrint(w http.ResponseWriter, r *http.Request) {
	// Read off the JSON body to bodyBytes for use in error logging if needed
	bodyBytes, _ := io.ReadAll(r.Body)

	// Reset r.Body to the original content
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Now prepare to decode the r.Body
	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))

	// Make struct repoUrls to decode the body into
	request := HandlerAddUserFingerPrintRequest{}

	err := decoder.Decode(&request)
	if err != nil || len(request.Logins) == 0 {
		msg := fmt.Sprintf("error parsing JSON for: %s", string(bodyBytes))
		RespondWithError(w, 400, msg)
		return
	}

	accepted := []string{}

	depsResults := make([][]int32, len(request.Dependencies))
	for index, dependency := range request.Dependencies {
		dependencyParams := database.BulkInsertDependenciesParams{
			DependencyName: dependency,
			Column2:        request.FileNames,
		}
		insertDepsResult, err := apiCfg.DB.BulkInsertDependencies(r.Context(), dependencyParams)
		depsResults[index] = insertDepsResult
		if err != nil {
			RespondWithError(w, 500, fmt.Sprintf("error inserting dependencies: %s", err))
		}
	}
	for _, login := range request.Logins {

		if err != nil {
			msg := fmt.Sprintf("error adding %s to repo_urls: %s", login, err)
			logger.LogError(msg)
			RespondWithError(w, 500, msg)
			return
		}

		accepted = append(accepted, login)
		for index := range request.Dependencies {
			userDependencyParams := database.InitializeUserDependenciesParams{
				Login:   login,
				Column2: depsResults[index],
			}

			err = apiCfg.DB.InitializeUserDependencies(r.Context(), userDependencyParams)
			if err != nil {
				RespondWithError(w, 500, fmt.Sprintf("error initializing repo dependencies: %s", err))
			}

		}
	}

	response := HandlerAddUserFingerPrintResponse{
		Accepted: accepted,
	}

	RespondWithJSON(w, 202, response)
}
