package gitutil

import (
	"main/internal/database"
	"main/internal/pkg/logger"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func BenchmarkProcessRepo(b *testing.B) {
	// Setup
	db, _, err := sqlmock.New()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	queries := database.New(db)

	repo := "OpenQ-Contracts"
	organization := "OpenQDev"
	repoUrl := "https://github.com/OpenQDev/OpenQ-Contracts"

	tmpDir, err := os.MkdirTemp("", "prefixPath")
	if err != nil {
		logger.LogFatalRedAndExit("can't create temp dir: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	prefixPath := tmpDir

	b.ResetTimer()

	// Clone repo to tmp dir. Will be deleted at end of test
	CloneRepo(prefixPath, organization, repo)

	for i := 0; i < b.N; i++ {
		err := ProcessRepo(prefixPath, repo, repoUrl, queries)
		if err != nil {
			b.Fatal(err)
		}
	}
}
