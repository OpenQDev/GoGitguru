package gitutil

import (
	"fmt"
	"strings"
)

func GitDependencyHistory(repoDir string, dependencySearched string, depFilePaths []string) (string, error) {
	var depHistoryString string
	depFilePathsFormatted := ""

	for _, path := range depFilePaths {
		depFilePathsFormatted += fmt.Sprintf("'**%s**' ", path)
	}
	depFilePathsFormatted = strings.TrimSpace(depFilePathsFormatted)
	fmt.Println("depFilePathsFormatted", depFilePathsFormatted)

	cmd := GitDepFileHistory(repoDir, dependencySearched, depFilePathsFormatted)
	fmt.Println("cmd", cmd)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("err", err)
		return depHistoryString, err
	}

	outStr := string(out)

	fmt.Println("outStr", outStr)

	return outStr, nil
}
