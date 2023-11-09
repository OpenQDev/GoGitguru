package gitutil

import (
	"strings"

	"github.com/OpenQDev/GoGitguru/util/logger"
)

func GitDependencyFiles(repoDir string, dependencyFiles []string) ([]string, error) {
	var dependencyPaths []string

	for _, dependencyFile := range dependencyFiles {
		cmd := LogDependencyFiles(repoDir, dependencyFile)

		out, err := cmd.CombinedOutput()
		if err != nil {
			logger.LogBlue("dependency file %s not found in %s", dependencyFile, repoDir)
			continue
		}

		outStr := string(out)

		// splits output on newlines and trims trailing spaces
		files := strings.Split(strings.TrimSpace(outStr), "\n")

		dependencyPaths = append(dependencyPaths, files...)
	}

	return dependencyPaths, nil
}
