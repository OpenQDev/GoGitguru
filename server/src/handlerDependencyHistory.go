package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/marshaller"
)

type DependencyHistoryRequest struct {
	RepoUrl            string   `json:"repo_url"`
	FilePaths          []string `json:"files_paths"`
	DependencySearched string   `json:"dependency_searched"`
}

type DependencyHistoryResponse struct{}

func (apiCfg *ApiConfig) HandlerDependencyHistory(w http.ResponseWriter, r *http.Request) {
	var dependencyHistory gitutil.DiffHistoryResult

	var body DependencyHistoryRequest
	err := marshaller.ReaderToType(r.Body, &body)
	if err != nil {
		RespondWithError(w, 200, err.Error())
		return
	}

	// Check if any of the fields are their zero value
	if body.RepoUrl == "" || len(body.FilePaths) == 0 || body.DependencySearched == "" {
		RespondWithError(w, 400, "All fields must be present in the request body.")
		return
	}

	organization, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(body.RepoUrl)

	repoDir := filepath.Join(apiCfg.PrefixPath, organization, repo)

	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		RespondWithError(w, 404, "Repository directory does not exist.")
		return
	}

	err = checkGitPathExists(body, repoDir)
	if err != nil {
		logger.LogError("not a git repository", err)
	}

	filesPathsFormatted := formatDependenciesSearched(body)

	logger.LogGreenDebug("filesPathsFormatted: %s", filesPathsFormatted)

	dependencyHistoryCmd := gitutil.GitDepFileHistory(repoDir, body.DependencySearched, filesPathsFormatted)

	dependencyHistoryOutput, err := dependencyHistoryCmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	logger.LogGreenDebug("dependencyHistoryCmd: %s", dependencyHistoryOutput)

	if len(dependencyHistoryOutput) == 0 {
		RespondWithJSON(w, 200, dependencyHistory)
	}

	cmd := gitutil.GitDependencySearch(repoDir, body.DependencySearched, filesPathsFormatted)
	dependencySearchOutput, err := cmd.Output()
	if err != nil {
		logger.LogError("error in GitDependencySearch: %s", err)
	}

	logger.LogGreenDebug("GitDependencySearch: %s", dependencyHistoryOutput)

	dependencyHistory = gitutil.DiffHistoryObject(string(dependencyHistoryOutput), body.DependencySearched, string(dependencySearchOutput))

	RespondWithJSON(w, 200, dependencyHistory)
}

func checkGitPathExists(body DependencyHistoryRequest, repoDir string) error {
	return nil
}

// This formats the dependency_searched array with wildcards around it
func formatDependenciesSearched(body DependencyHistoryRequest) string {
	var filesPathsFormatted string
	for _, path := range body.FilePaths {
		filesPathsFormatted += "'**" + path + "**' "
	}
	filesPathsFormatted = strings.TrimSpace(filesPathsFormatted)
	return filesPathsFormatted
}
