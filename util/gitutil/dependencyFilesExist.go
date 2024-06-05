package gitutil

func GitDependencyFiles(repoDir string, dependencyFiles []string) ([]string, error) {
	var dependencyPaths []string

	for _, dependencyFile := range dependencyFiles {
		files, err := LogDependencyFiles(repoDir, dependencyFile)
		if err != nil {
			continue
		}

		dependencyPaths = append(dependencyPaths, files...)
	}

	return dependencyPaths, nil
}
