package usersync

import (
	"reflect"
	"testing"
)

func TestGetRepoToAuthorsMap(t *testing.T) {
	tests := []struct {
		name           string
		input          []UserSync
		expectedOutput map[string][]string
	}{
		{
			name: "Single author, single repo",
			input: []UserSync{
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
						URL:     "https://github.com/example/repo",
						NotNull: true,
					},
				},
			},
			expectedOutput: map[string][]string{
				"https://github.com/example/repo": {"test@example.com"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRepoToAuthorsMap(tt.input)
			if !reflect.DeepEqual(result, tt.expectedOutput) {
				t.Errorf("getRepoToAuthorsMap() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}