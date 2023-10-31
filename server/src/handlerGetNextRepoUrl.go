package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

type HandlerGetNextRepoUrlResponse struct {
	RepoUrl string `json:"repo_url"`
}

func (apiCfg *ApiConfig) HandlerGetNextRepoUrl(w http.ResponseWriter, r *http.Request) {
	var response HandlerGetNextRepoUrlResponse
	db, _ := setup.GetSQLDatbase(apiCfg.DBURL)
	repoUrl, err := GetDueURL(db, "10")
	if err != nil {
		logger.LogError("errror: %s", err)
	}

	response = HandlerGetNextRepoUrlResponse{
		RepoUrl: repoUrl,
	}

	RespondWithJSON(w, 200, response)
}

func GetDueURL(db *sql.DB, repoSyncInterval string) (string, error) {
	var url string

	// Start a new transaction
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer a rollback in case anything fails
	defer tx.Rollback()

	// Prepare the SQL statement
	query := `SELECT url FROM repo_urls WHERE status IN ('synced', 'pending', 'failed') OR updated_at IS NULL ORDER BY RANDOM() LIMIT 1 FOR UPDATE`
	row := tx.QueryRowContext(context.Background(), query)

	// Execute the query and scan the result
	err = row.Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to execute query: %w", err)
	}

	// Prepare the update statement
	update := `UPDATE repo_urls SET status = 'queued', updated_at = NOW() WHERE url = $1`
	_, err = tx.ExecContext(context.Background(), update, url)
	if err != nil {
		return "", fmt.Errorf("failed to update url: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return url, nil
}
