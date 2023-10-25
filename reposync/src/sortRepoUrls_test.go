package reposync

import (
	"testing"

	"github.com/OpenQDev/GoGitguru/database"
)

func TestSortRepoUrls(t *testing.T) {
	repoUrlObjects := []database.RepoUrl{
		{Url: "https://github.com/Z"},
		{Url: "https://github.com/A"},
		{Url: "https://github.com/c"},
		{Url: "https://github.com/b"},
	}

	expected := []string{
		"https://github.com/a",
		"https://github.com/b",
		"https://github.com/c",
		"https://github.com/z",
	}

	result := sortRepoUrls(repoUrlObjects)

	for i, url := range result {
		if url != expected[i] {
			t.Errorf("Expected %s, but got %s", expected[i], url)
		}
	}
}
