package server

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
)

type DependencyHistoryRequest struct {
	RepoUrls           []string `json:"repo_urls"`
	FilePaths          []string `json:"files_paths"`
	DependencySearched string   `json:"dependency_searched"`
}

type DependencyHistoryResponseMember struct {
	DatesAdded   int64  `json:"date_added"`
	DatesRemoved int64  `json:"date_removed"`
	RepoUrl      string `json:"repo_url"`
}

func (apiCfg *ApiConfig) HandlerDependencyHistory(w http.ResponseWriter, r *http.Request) {
	var dependencyHistoryResponse []DependencyHistoryResponseMember
	var body DependencyHistoryRequest
	err := marshaller.ReaderToType(r.Body, &body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request in dependency-history: %s", err))
		return
	}

	repoDependencyParams := database.GetRepoDependenciesParams{
		DependencyName: body.DependencySearched,
		Column2:        body.RepoUrls,
		Column3:        body.FilePaths,
	}

	dependencyResult, err := apiCfg.DB.GetRepoDependencies(r.Context(), repoDependencyParams)

	if err != nil {
		fmt.Println(err)
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting dependencies: %s", err))
		return
	}
	for _, dependency := range dependencyResult {
		hasRepoUrl := slices.ContainsFunc(dependencyHistoryResponse, func(dependencyHistoryResponse DependencyHistoryResponseMember) bool {
			return dependencyHistoryResponse.RepoUrl == dependency.Url.String
		})

		if !hasRepoUrl {
			dependencyHistoryResponse = append(dependencyHistoryResponse, DependencyHistoryResponseMember{
				DatesAdded:   dependency.FirstUseDate.Int64,
				DatesRemoved: dependency.LastUseDate.Int64,
				RepoUrl:      dependency.Url.String,
			})
		} else {
			correctSliceIndex := slices.IndexFunc(dependencyHistoryResponse, func(dependencyHistoryResponse DependencyHistoryResponseMember) bool {
				return dependencyHistoryResponse.RepoUrl == dependency.Url.String
			})
			dependencyUsedBefore := dependency.FirstUseDate.Int64 != 0
			dependencyAlreadyCurrentlyUsed := dependencyHistoryResponse[correctSliceIndex].DatesAdded != 0

			if dependencyUsedBefore && (dependency.FirstUseDate.Int64 < dependencyHistoryResponse[correctSliceIndex].DatesAdded || dependencyHistoryResponse[correctSliceIndex].DatesAdded == 0) {
				dependencyHistoryResponse[correctSliceIndex].DatesAdded = dependency.FirstUseDate.Int64
			}
			newDependencyDateMoreRelevant := dependency.LastUseDate.Int64 > dependencyHistoryResponse[correctSliceIndex].DatesRemoved
			if dependencyUsedBefore && !dependencyAlreadyCurrentlyUsed && newDependencyDateMoreRelevant {
				dependencyHistoryResponse[correctSliceIndex].DatesRemoved = dependency.LastUseDate.Int64
			}
		}

	}

	RespondWithJSON(w, http.StatusOK, dependencyHistoryResponse)
}
