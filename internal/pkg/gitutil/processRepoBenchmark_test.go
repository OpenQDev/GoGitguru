package gitutil

import (
	"main/internal/pkg/logger"
	"main/internal/pkg/setup"
	"os"
	"testing"
)

func BenchmarkProcessRepo(b *testing.B) {
	// Setup
	_, dbUrl, _, _, _, _, _, _, _ := setup.ExtractAndVerifyEnvironment("../../../.env")

	db, err := setup.GetDatbase(dbUrl)
	if err != nil {
		logger.LogFatalRedAndExit("unable to run BenchmarkProcessRepo. unable to connect to database: %s", err)
	}

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
		err := ProcessRepo(prefixPath, repo, repoUrl, db)
		if err != nil {
			b.Fatal(err)
		}
	}
}
