package gitutil

// git -C . log --reverse --pretty=format:"%H-;-%an-;-%ae-;-%at-;-%ct%n%s" --numstat --since=2020-01-01
// git -C . rev-parse --is-inside-work-tree

import (
	"main/internal/pkg/logger"
	"os/exec"
	"path/filepath"
)

func FormatGitLogs(prefixPath string, repo string, fromCommitDate string) []GitLog {
	fullRepoPath := filepath.Join(prefixPath, repo)

	cmdCheck := exec.Command("git", "-C", fullRepoPath, "rev-parse", "--is-inside-work-tree")
	err := cmdCheck.Run()
	if err != nil {
		logger.LogFatalRedAndExit("%s/%s is not a git repository: %s", prefixPath, repo, err)
	}

	if fromCommitDate == "" {
		fromCommitDate = "2020-01-01"
	}

	cmd := exec.Command("git", "-C", fullRepoPath, "log", "--reverse", "--pretty=format:%H-;-%an-;-%ae-;-%at-;-%ct%n%s", "--numstat", "--since="+fromCommitDate)

	out, err := cmd.Output()

	if err != nil {
		logger.LogFatalRedAndExit("error running git log in %s: %s", fullRepoPath, err)
	}

	gitLogs := ProcessGitLogs(string(out))

	return gitLogs

	// lines := strings.Split(string(out), "\n")
	// csvFile, err := os.Create("git_log.csv")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer csvFile.Close()

	// writer := csv.NewWriter(csvFile)
	// defer writer.Flush()

	// for _, line := range lines {
	// 	writer.Write(strings.Split(line, ";"))
	// }
}
