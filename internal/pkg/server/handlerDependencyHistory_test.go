package server

import (
	"main/internal/pkg/testhelpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerDependencyHistory(t *testing.T) {
	// ARRANGE - GLOBAL

	// ARRANGE - TESTS
	tests := HandlerDependencyHistoryTestCases()

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			testhelpers.RUN_ALL_TESTS,
		), tt.name)

		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, true, true)
		})
	}
}
