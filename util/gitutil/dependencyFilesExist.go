package gitutil

import (
	"runtime"
	"strings"
)

func GitDependencyFiles(repoDir string, dependencyFiles []string) ([]string, error) {
	var dependencyPaths []string

	for _, dependencyFile := range dependencyFiles {
		if runtime.GOOS != "windows" {
			cmd := LogDependencyFiles(repoDir, dependencyFile)
			out, err := cmd.CombinedOutput()
			if err != nil {
				continue
			}
			outStr := string(out)

			files := strings.Split(strings.TrimSpace(outStr), "\n")
			dependencyPaths = append(dependencyPaths, files...)
		} else {
			files, err := DANGEROUS_WINDOWS_COMPAT_LogDependencyFiles(repoDir, dependencyFile)
			if err != nil {
				continue
			}
			dependencyPaths = append(dependencyPaths, files...)
		}
	}

	return dependencyPaths, nil
}
