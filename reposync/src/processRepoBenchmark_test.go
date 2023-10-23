package reposync

import (
	"os"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/gitutil"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/setup"
)

func BenchmarkProcessRepo(b *testing.B) {
	// Setup
	env := setup.ExtractAndVerifyEnvironment("../../../.env")
	dbUrl := env.DbUrl

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
	gitutil.CloneRepo(prefixPath, organization, repo)

	for i := 0; i < b.N; i++ {
		err := ProcessRepo(prefixPath, organization, repo, repoUrl, db)
		if err != nil {
			b.Fatal(err)
		}
	}
}
