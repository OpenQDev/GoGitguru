package testhelpers

import (
	"slices"
	"testing"
)

const RUN_ALL_TESTS = "RUN_ALL_TESTS"
const RUN_NO_TESTS = ""

func CheckTestSkip(t *testing.T, targets []string, target string) {
	if targets[0] == RUN_ALL_TESTS {
		return
	}

	if !slices.Contains(targets, target) {
		t.Skipf("skipping %s", target)
	}
}

func Targets(args ...string) []string {
	return args
}
