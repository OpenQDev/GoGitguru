package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/OpenQDev/GoGitguru/database"
)

type HandlerUserDependencyHistoryRequest struct {
	Dependency string `json:"dependency"`
	FileName   string `json:"fileName"`
	Login      string `json:"login"`
}
type HandlerUserDependencyHistoryResponse struct{}

func (apiConfig *ApiConfig) HandlerUserDependencyHistory(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)

	// Reset r.Body to the original content
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Now prepare to decode the r.Body
	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))

	// Make struct repoUrls to decode the body into
	request := HandlerUserDependencyHistoryRequest{}

	err := decoder.Decode(&request)
	if err != nil {
		msg := fmt.Sprintf("error parsing JSON for: %s", string(bodyBytes))
		RespondWithError(w, 400, msg)
		return
	}

	login := request.Login
	dependency := request.Dependency
	fileName := request.FileName
	dependencyParams := database.GetDependencyParams{
		DependencyName: dependency,
		DependencyFile: fileName,
	}
	dependencyExists, err := apiConfig.DB.GetDependency(context.Background(), dependencyParams)

	userExists, err := apiConfig.DB.CheckGithubUserExists(context.Background(), login)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if userExists {
		user, err := apiConfig.DB.GetGithubUser(context.Background(), login)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, ConvertToReturnUserWithDependency(ConvertDatabaseInsertUserParamsToServerUser(user), dependencyExists))
		return
	}

}

type ReturnUserAndDependency struct {
	DependencyFile  string `json:"dependency_file"`
	DependencyName  string `json:"dependency_name"`
	DependencyID    int32  `json:"dependency_id"`
	UserID          int    `json:"internal_id"`
	GithubRestID    int    `json:"github_rest_id"`
	GithubGraphqlID string `json:"github_graphql_id"`
	Login           string `json:"login"`
}

func ConvertToReturnUserWithDependency(user User, dependency database.Dependency) ReturnUserAndDependency {
	return ReturnUserAndDependency{
		DependencyFile:  dependency.DependencyFile,
		DependencyName:  dependency.DependencyName,
		DependencyID:    dependency.InternalID,
		UserID:          user.InternalID,
		GithubRestID:    user.GithubRestID,
		GithubGraphqlID: user.GithubGraphqlID,
		Login:           user.Login,
	}
}
