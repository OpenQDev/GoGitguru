package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
)

type DependencyHistoryRequest struct {
	RepoUrl            string   `json:"repo_url"`
	FilePaths          []string `json:"files_paths"`
	DependencySearched string   `json:"dependency_searched"`
}

type DependencyHistoryResponse struct {
	DatesAdded   []string `json:"dates_added"`
	DatesRemoved []string `json:"dates_removed"`
}

func (apiCfg *ApiConfig) HandlerDependencyHistory(w http.ResponseWriter, r *http.Request) {
	var dependencyHistoryResponse DependencyHistoryResponse
	fmt.Println("starting")
	var body DependencyHistoryRequest
	err := marshaller.ReaderToType(r.Body, &body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request in dependency-history: %s", err))
		return
	}

	repoDependencyParams := database.GetRepoDependenciesParams{
		DependencyName: body.DependencySearched,
		Url:            body.RepoUrl,
		Column3:        body.FilePaths,
	}
	fmt.Println("repoDependencyParams", repoDependencyParams)
	dependencyResult, err := apiCfg.DB.GetRepoDependencies(r.Context(), repoDependencyParams)
	fmt.Println("dependencyResult", dependencyResult)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting dependencies: %s", err))
		return
	}
	datesAdded := []string{}
	datesRemoved := []string{}
	fmt.Println("dependencyResult", dependencyResult)
	for _, dependency := range dependencyResult {
		if dependency.FirstUseDate.Int64 != 0 {
			t := time.Unix(dependency.FirstUseDate.Int64, 0).UTC()
			formattedDate := t.Format(time.RFC3339)
			datesAdded = append(datesAdded, formattedDate)

		}
		if dependency.LastUseDate.Int64 != 0 {
			//get date in this format 2023-01-10T17:55:48Z from timestamp
			t := time.Unix(dependency.LastUseDate.Int64, 0).UTC()
			formattedDate := t.Format(time.RFC3339)
			datesRemoved = append(datesRemoved, formattedDate)
		}
	}

	dependencyHistoryResponse = DependencyHistoryResponse{
		DatesAdded:   datesAdded,
		DatesRemoved: datesRemoved,
	}

	RespondWithJSON(w, http.StatusOK, dependencyHistoryResponse)
}
