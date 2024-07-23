package gitutil

import (
	"fmt"
	"runtime"
	"strings"
)

func GitDependencyFiles(repoDir string, dependencyFiles []string) ([]string, error) {
	var dependencyPaths []string

	for _, dependencyFile := range dependencyFiles {
		if runtime.GOOS == "windows" {
			files, err := DANGEROUS_WINDOWS_COMPAT_LogDependencyFiles(repoDir, dependencyFile)
			if err != nil {
				fmt.Println("error", err)
				continue
			}
			dependencyPaths = append(dependencyPaths, files...)

		} else {
			cmd := LogDependencyFiles(repoDir, dependencyFile)
			out, err := cmd.CombinedOutput()
			if err != nil {
				continue
			}
			outStr := string(out)

			files := strings.Split(strings.TrimSpace(outStr), "\n")
			dependencyPaths = append(dependencyPaths, files...)
		}
	}

	return dependencyPaths, nil
}
