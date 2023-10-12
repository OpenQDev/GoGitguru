package sync

import (
	"main/internal/pkg/testhelpers"
	"reflect"
	"testing"
)

func TestGetRepoToAuthorsMap(t *testing.T) {
	tests := []struct {
		name           string
		input          []UserSync
		expectedOutput map[string][]AuthorCommitTuple
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
				{
					CommitHash: "otherCommitHash",
					Author: struct {
						Email   string
						NotNull bool
					}{
						Email:   "otherperson@example.com",
						NotNull: true,
					},
					Repo: struct {
						URL     string
						NotNull bool
					}{
						URL:     "https://github.com/example/other-repo",
						NotNull: true,
					},
				},
			},
			expectedOutput: map[string][]AuthorCommitTuple{
				"https://github.com/example/repo":  []AuthorCommitTuple{AuthorCommitTuple{Author: "test@example.com", CommitHash: "abc123"}},
				"https://github.com/example/repo2": []AuthorCommitTuple{AuthorCommitTuple{Author: "otherperson@example.com", CommitHash: "otherCommitHash"}},
			},
		},
	}

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			result := GetRepoToAuthorsMap(tt.input)
			if !reflect.DeepEqual(result, tt.expectedOutput) {
				t.Errorf("getRepoToAuthorsMap() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
