package sync

import (
	"main/internal/pkg/testhelpers"
	"reflect"
	"testing"
)

func TestGenerateBatchAuthors(t *testing.T) {
	tests := []struct {
		name           string
		input          map[string][]AuthorCommitTuple
		batchSize      int
		expectedOutput [][]interface{}
	}{
		{
			name: "Single author, single repo",
			input: map[string][]AuthorCommitTuple{
				"https://github.com/test/repo": {
					{"test@example.com", "commit1"},
					{"test2@example.com", "commit2"},
					{"test3@example.com", "commit3"},
				},
			},
			batchSize: 2,
			expectedOutput: [][]interface{}{
				{"https://github.com/test/repo", []AuthorCommitTuple{
					{"test@example.com", "commit1"},
					{"test2@example.com", "commit2"},
				}},
				{"https://github.com/test/repo", []AuthorCommitTuple{
					{"test3@example.com", "commit3"},
				}},
			},
		},
		{
			name: "Single author, two repos",
			input: map[string][]AuthorCommitTuple{
				"https://github.com/test/repo": {
					{"test@example.com", "commit1"},
					{"test2@example.com", "commit2"},
					{"test3@example.com", "commit3"},
				},
				"https://github.com/test/repo2": {
					{"author123@example.com", "commit4"},
					{"author12sdfdsf@example.com", "commit5"},
					{"authosdfsdf@example.com", "commit6"},
				},
			},
			batchSize: 2,
			expectedOutput: [][]interface{}{
				{"https://github.com/test/repo", []AuthorCommitTuple{
					{"test@example.com", "commit1"},
					{"test2@example.com", "commit2"},
				}},
				{"https://github.com/test/repo", []AuthorCommitTuple{
					{"test3@example.com", "commit3"},
				}},
				{"https://github.com/test/repo2", []AuthorCommitTuple{
					{"author123@example.com", "commit4"},
					{"author12sdfdsf@example.com", "commit5"},
				}},
				{"https://github.com/test/repo2", []AuthorCommitTuple{
					{"authosdfsdf@example.com", "commit6"},
				}},
			},
		},
	}

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			result := GenerateBatchAuthors(tt.input, tt.batchSize)
			if !reflect.DeepEqual(result, tt.expectedOutput) {
				t.Errorf("batchAuthors() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
