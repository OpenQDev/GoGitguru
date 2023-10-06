package usersync

import (
	"database/sql"
	"main/internal/database"
	"reflect"
	"testing"
)

func TestConvertToUserSync(t *testing.T) {
	tests := []struct {
		name           string
		input          []database.GetLatestUncheckedCommitPerAuthorRow
		expectedOutput []UserSync
	}{
		{
			name: "Single author, single repo",
			input: []database.GetLatestUncheckedCommitPerAuthorRow{
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
			},
			expectedOutput: []UserSync{
				{
					CommitHash: "abc123",
					Author: struct {
						Email   string
						NotNull bool
					}{
						Email:   "test@example.com",
						NotNull: true,
					},
					Repo: struct {
						URL     string
						NotNull bool
					}{
						URL:     "https://github.com/test/repo",
						NotNull: true,
					},
				},
			},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertToUserSync(tt.input)
			if !reflect.DeepEqual(result, tt.expectedOutput) {
				t.Errorf("convertToUserSync() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
