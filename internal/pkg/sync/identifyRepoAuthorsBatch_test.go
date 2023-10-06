package sync

import (
	"bytes"
	"errors"
	"main/internal/pkg/gitutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGitUtil struct {
	mock.Mock
}

func (m *MockGitUtil) GithubGetCommitAuthors(query string, ghAccessToken string) (*gitutil.CommitAuthorsResponse, error) {
	args := m.Called(query, ghAccessToken)
	return args.Get(0).(*gitutil.CommitAuthorsResponse), args.Error(1)
}

func TestIdentifyRepoAuthorsBatch(t *testing.T) {
	mockGitUtil := new(MockGitUtil)
	ghAccessToken := "testToken"
	repoUrl := "https://github.com/openqdev/openq-workflows"
	authorList := []string{"author1", "author2"}

	// Define test cases
	tests := []struct {
		name          string
		mockResponse  *gitutil.CommitAuthorsResponse
		mockError     error
		expectedError string
	}{
		{
			name: "Valid response",
			mockResponse: &gitutil.CommitAuthorsResponse{
				Data: map[string]gitutil.Commit{
					"commit1": {
						Author: gitutil.Author{
							User: gitutil.User{
								Login: "author1",
							},
						},
					},
				},
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name:          "Error response",
			mockResponse:  nil,
			mockError:     errors.New("mock error"),
			expectedError: "error occured while fetching from GraphQL API: mock error",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGitUtil.On("GithubGetCommitAuthors", mock.Anything, ghAccessToken).Return(tt.mockResponse, tt.mockError)

			// Redirect logger output to buffer for testing
			var buf bytes.Buffer

			IdentifyRepoAuthorsBatch(repoUrl, authorList, ghAccessToken)

			if tt.expectedError != "" {
				assert.Contains(t, buf.String(), tt.expectedError)
			}
		})
	}
}
