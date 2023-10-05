package sync

import (
	"database/sql"
	"main/internal/database"
	"testing"
)

func TestConvertToUserSync(t *testing.T) {
	// Mock data
	mockData := []database.GetLatestUncheckedCommitPerAuthorRow{
		{
			CommitHash: "abc123",
			AuthorEmail: sql.NullString{
				String: "test@example.com",
				Valid:  true,
			},
			RepoUrl: sql.NullString{
				String: "https://github.com/test/repo",
				Valid:  true,
			},
		},
		// Add more mock data as needed
	}

	// Call the function with mock data
	result := convertToUserSync(mockData)

	// Check the length of the result
	if len(result) != len(mockData) {
		t.Errorf("Expected length %d, got %d", len(mockData), len(result))
	}

	// Check the contents of the result
	for i, userSync := range result {
		if userSync.CommitHash != mockData[i].CommitHash {
			t.Errorf("Expected CommitHash %s, got %s", mockData[i].CommitHash, userSync.CommitHash)
		}
		if userSync.Author.Email != mockData[i].AuthorEmail.String {
			t.Errorf("Expected AuthorEmail %s, got %s", mockData[i].AuthorEmail.String, userSync.Author.Email)
		}
		if userSync.Repo.URL != mockData[i].RepoUrl.String {
			t.Errorf("Expected RepoUrl %s, got %s", mockData[i].RepoUrl.String, userSync.Repo.URL)
		}
	}
}
