package gitutil

import (
	"fmt"
	"strings"
)

func GitDependencyHistorySnapshot(repoDir string, dependencySearched string, depFilePaths string) ([]string, error) {
	var depFiles []string

	cmd := GitDepFileToday(repoDir, dependencySearched, depFilePaths)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("err", err)
		return depFiles, err
	}

	outStr := string(out)

	depFiles = strings.Split(strings.TrimSpace(outStr), "\n")

	return depFiles, nil
}
