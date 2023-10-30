package gitutil

import (
	"fmt"
	"strings"
)

func GitDependencyFiles(repoDir string, dependencyFile string) ([]string, error) {
	var dependencyPaths []string

	cmd := LogDependencyFiles(repoDir, dependencyFile)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("err", err)
		return dependencyPaths, err
	}

	outStr := string(out)

	// splits output on newlines and trims trailing spaces
	files := strings.Split(strings.TrimSpace(outStr), "\n")

	dependencyPaths = append(dependencyPaths, files...)

	return dependencyPaths, nil
}
