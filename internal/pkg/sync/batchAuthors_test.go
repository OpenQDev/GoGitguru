package sync

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBatchAuthors(t *testing.T) {
	tests := []struct {
		name           string
		input          map[string][]string
		batchSize      int
		expectedOutput [][]interface{}
	}{
		{
			name: "Single author, single repo",
			input: map[string][]string{
				"https://github.com/test/repo": {"test@example.com", "test2@example.com", "test3@example.com"},
			},
			batchSize: 2,
			expectedOutput: [][]interface{}{
				{"https://github.com/test/repo", []string{"test@example.com", "test2@example.com"}},
				{"https://github.com/test/repo", []string{"test3@example.com"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BatchAuthors(tt.input, tt.batchSize)
			fmt.Println(result)
			if !reflect.DeepEqual(result, tt.expectedOutput) {
				t.Errorf("batchAuthors() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
