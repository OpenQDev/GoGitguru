package util

import (
	"main/internal/database"
	"testing"
)

func TestExtractOrganizationAndRepositoryFromUrl(t *testing.T) {
	tests := []struct {
		name string
		url  string
		org  string
		repo string
	}{
		{
			name: "Test with valid URL",
			url:  "https://github.com/org/repo",
			org:  "org",
			repo: "repo",
		},
		{
			name: "Test with URL without repo",
			url:  "https://github.com/org/",
			org:  "org",
			repo: "",
		},
		{
			name: "Test with URL without org and repo",
			url:  "https://github.com/",
			org:  "",
			repo: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			org, repo := ExtractOrganizationAndRepositoryFromUrl(database.RepoUrl{Url: tt.url})
			if org != tt.org || repo != tt.repo {
				t.Errorf("got %v %v, want %v %v", org, repo, tt.org, tt.repo)
			}
		})
	}
}
