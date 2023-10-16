package gitutil

import (
	"main/internal/database"
	"main/internal/pkg/logger"
	"main/internal/pkg/testhelpers"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestProcessRepo(t *testing.T) {
	// ARRANGE - GLOBAL
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("can't create mock DB: %s", err)
	}

	queries := database.New(db)

	tmpDir, err := os.MkdirTemp("", "prefixPath")
	if err != nil {
		logger.LogFatalRedAndExit("can't create temp dir: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	prefixPath := tmpDir

	// ARRANGE - TESTS
	tests := ProcessRepoTestCases()

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE - LOCAL
			CloneRepo(prefixPath, tt.organization, tt.repo)

			tt.setupMock(mock, tt.gitLogs, tt.repoUrl)

			// ACT
			ProcessRepo(prefixPath, tt.repo, tt.repoUrl, queries)
			if err != nil {
				t.Errorf("there was an error processing repo: %s", err)
			}

			// ASSERT

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			assert.Nil(t, err)
		})
	}
}
