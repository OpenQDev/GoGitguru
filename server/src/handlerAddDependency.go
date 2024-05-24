package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/OpenQDev/GoGitguru/database"
)

type HandlerAddDependencyRequest struct {
	Dependency string `json:"dependency"`
	FileNames   []string `json:"fileNames"`
	Repository string `json:"repository"`
}

type HandlerAddDependencyResponse struct {
	Status string `json:"status"`
}

func (apiCfg *ApiConfig) HandlerAddDependency(w http.ResponseWriter, r *http.Request) {
	// Read off the JSON body to bodyBytes for use in error logging if needed
	bodyBytes, _ := io.ReadAll(r.Body)

	// Reset r.Body to the original content
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Now prepare to decode the r.Body
	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))

	// Make struct repoUrls to decode the body into
	request := HandlerAddDependencyRequest{}

	err := decoder.Decode(&request)
	if err != nil {
		msg := fmt.Sprintf("error parsing JSON for: %s", string(bodyBytes))
		RespondWithError(w, 400, msg)
		return
	}

	dependency := request.Dependency
	fileNames := request.FileNames
	repository := request.Repository
	status := "accepted"
	println("dependency", dependency)
	println("fileNames", fileNames[0])
	dependencyParams := database.BulkInsertDependenciesParams{
		DependencyName: dependency,
		Column2: 	  fileNames,
	}
insertDepsResult, err := apiCfg.DB.BulkInsertDependencies(r.Context(), dependencyParams)
	if err != nil {
		status = "exists"
	}
println("insertDepsResult", insertDepsResult)

println("depsIds",  repository)
	repoDependencyParams := database.InitializeRepoDependenciesParams{
		Url: repository,
		Column2: insertDepsResult,
	}



	err= apiCfg.DB.InitializeRepoDependencies(r.Context(), repoDependencyParams);
	if err != nil {
		println("error initializing repo dependencies", err.Error())
		status = "error"
	}


	response := HandlerAddDependencyResponse{
		status,
	}

	RespondWithJSON(w, 202, response)
}

