import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDiffHistoryObject(t *testing.T) {
	dependencyHistoryLogs := ""
	dependencySearched := "axios"
	dependencyTodayOutput := []string{"package.json"}

	result := DiffHistoryObject(dependencyHistoryLogs, dependencySearched, dependencyTodayOutput)

	assert.NotNil(t, result)
	assert.IsType(t, DiffHistoryResult{}, result)
	assert.Equal(t, 1, len(result.CommitsSummary))
	assert.Equal(t, "1234567890abcdef", result.CommitsSummary[0].CommitHash)
	assert.Equal(t, time.Now().Format(time.RFC3339), result.CommitsSummary[0].Date.Format(time.RFC3339))
	assert.Equal(t, 1, result.CommitsSummary[0].AddedLines)
	assert.Equal(t, 0, result.CommitsSummary[0].RemovedLines)
	assert.Equal(t, "commit 1234567890abcdef", result.CommitsSummary[0].FullContent)
	assert.Equal(t, 1, len(result.DatesAdded))
	assert.Equal(t, 0, len(result.DatesRemoved))
}
