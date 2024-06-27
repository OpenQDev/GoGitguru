package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/lib"
	"github.com/OpenQDev/GoGitguru/util/logger"
)

type HandlerAddDependencyPatternRequest struct {
	DependencyPatterns []string `json:"dependency_patterns"`
	Creator            string   `json:"creator"`
}

type HandlerAddDependencyPatternResponse struct {
	Accepted string `json:"accepted"`
}

func (apiCfg *ApiConfig) HandlerAddDependencyPattern(w http.ResponseWriter, r *http.Request) {
	// Read off the JSON body to bodyBytes for use in error logging if needed
	bodyBytes, _ := io.ReadAll(r.Body)

	// Reset r.Body to the original content
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Now prepare to decode the r.Body
	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))

	// Make struct repoUrls to decode the body into
	request := HandlerAddDependencyPatternRequest{}

	err := decoder.Decode(&request)
	if err != nil || len(request.DependencyPatterns) == 0 {
		msg := fmt.Sprintf("error parsing JSON for: %s", string(bodyBytes))
		RespondWithError(w, 400, msg)
		return
	}
	bulkUpsertFilePatternsParams := database.BulkUpsertFilePatternsParams{
		Patterns:  request.DependencyPatterns,
		UpdatedAt: lib.Now().Unix(),
		Creator:   request.Creator,
	}
	err = apiCfg.DB.BulkUpsertFilePatterns(r.Context(), bulkUpsertFilePatternsParams)

	if err != nil {
		msg := fmt.Sprintf("error adding %s to dependency_patterns: %s", request.DependencyPatterns, err)
		logger.LogError(msg)
		RespondWithError(w, 500, msg)
		return
	}

	response := HandlerAddDependencyPatternResponse{
		Accepted: "true",
	}

	RespondWithJSON(w, 202, response)
}
