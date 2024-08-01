package server

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"

	"github.com/OpenQDev/GoGitguru/util/setup"
)

type HandlerSetCurrentRepoUrlsToPendingV2Response struct {
	Accepted bool `json:"accepted"`
}

func (apiCfg *ApiConfig) HandlerSetCurrentRepoUrlsToPendingV2(w http.ResponseWriter, r *http.Request) {
	env := setup.ExtractAndVerifyEnvironment("../../.env")
	// Get the contents of the authorization header
	authHeader := r.Header.Get("Authorization")

	// Match the authorization header with env.api key
	if authHeader != env.GitguruApiKey {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// TODO add 401 for env.GitguruApiKey

	// Read off the JSON body to bodyBytes for use in error logging if needed
	bodyBytes, _ := io.ReadAll(r.Body)

	// Reset r.Body to the original content
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Now prepare to decode the r.Body

	response := HandlerSetCurrentRepoUrlsToPendingV2Response{
		Accepted: true,
	}

	RespondWithJSON(w, 202, response)

	tx, err := apiCfg.Conn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		fmt.Println("failed to begin transaction: %w", err)
	}

	query := `BEGIN; TRUNCATE TABLE repo_urls_v2; INSERT INTO repo_urls_v2 (url, status) SELECT url, 'pending'::repo_status FROM repo_urls; UPDATE users_to_dependencies SET resync_all = TRUE;`
	tx.Exec(query)
	tx.Commit()
}
