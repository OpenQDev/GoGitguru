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
	FileName   string `json:"fileName"`
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
	fileName := request.FileName
	status := "accepted"
	err = addDependencyToList(apiCfg, r, dependency, fileName)
	if err != nil {
		status = "exists"
	}

	response := HandlerAddDependencyResponse{
		status,
	}

	RespondWithJSON(w, 202, response)
}

func addDependencyToList(apiCfg *ApiConfig, r *http.Request, DependencyName string, DependencyFile string) error {
	dependency := database.InsertDependencyParams{
		DependencyName: DependencyName,
		DependencyFile: DependencyFile,
	}
	_result, err := apiCfg.DB.InsertDependency(r.Context(), dependency)
	println(_result.DependencyName, _result.DependencyFile)
	return err
}
