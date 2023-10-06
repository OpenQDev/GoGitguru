package server

import (
	"encoding/json"
	"fmt"
	"main/internal/pkg/gitutil"
	"main/internal/pkg/logger"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type DependencyHistoryBody struct {
	RepoUrl            string   `json:"repo_url"`
	FilePaths          []string `json:"files_paths"`
	DependencySearched string   `json:"dependency_searched"`
}

type DependencyHistoryReturn struct {
	CommitsSummary []string `json:"commits_summary"`
	DatesAdded     []string `json:"dates_added"`
	DatesRemoved   []string `json:"dates_removed"`
}

func (apiCfg *ApiConfig) HandlerDependencyHistory(w http.ResponseWriter, r *http.Request) {
	var dependencyHistory DependencyHistoryReturn
	var body DependencyHistoryBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		RespondWithError(w, 202, err.Error())
		return
	}

	// Check if any of the fields are their zero value
	if body.RepoUrl == "" || len(body.FilePaths) == 0 || body.DependencySearched == "" {
		RespondWithError(w, 400, "All fields must be present in the request body.")
		return
	}

	_, repo := gitutil.ExtractOrganizationAndRepositoryFromUrl(body.RepoUrl)

	repoDir := filepath.Join("/Users/alo/OpenQ-Fullstack/GoGitguru/repos", repo)

	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		RespondWithError(w, 404, "Repository directory does not exist.")
		return
	}

	checkGitPathExists(body, repoDir, err)

	filesPathsFormatted := formatDependenciesSearched(body)

	fmt.Println(filesPathsFormatted)

	RespondWithJSON(w, 202, dependencyHistory)
}

func checkGitPathExists(body DependencyHistoryBody, repoDir string, err error) {
	gitGrepExists := strings.Join(body.FilePaths, "|")
	gitGrepExists = strings.ReplaceAll(gitGrepExists, "*", "")

	cmd := gitutil.GitPathExists(repoDir, gitGrepExists)

	_, err = cmd.Output()
	if err != nil {
		logger.LogFatalRedAndExit("error checking for path: %s", err)
	}
}

// This formats the dependency_searched array with wildcards around it
func formatDependenciesSearched(body DependencyHistoryBody) string {
	var filesPathsFormatted string
	for _, path := range body.FilePaths {
		filesPathsFormatted += "'**" + path + "**' "
	}
	filesPathsFormatted = strings.TrimSpace(filesPathsFormatted)
	return filesPathsFormatted
}
