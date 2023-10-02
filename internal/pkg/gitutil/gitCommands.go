package gitutil

import (
	"os/exec"
)

func GitCloneCommand(cloneString string, cloneDestination string) *exec.Cmd {
	return exec.Command("git", "clone", cloneString, cloneDestination)
}

func CheckIsAGitRepo(fullRepoPath string) *exec.Cmd {
	return exec.Command("git", "-C", fullRepoPath, "rev-parse", "--is-inside-work-tree")
}

func GitLogCommand(fullRepoPath string, fromCommitDate string) *exec.Cmd {
	return exec.Command("git", "-C", fullRepoPath, "log", "--reverse", "--pretty=format:%H-;-%an-;-%ae-;-%at-;-%ct%n%s", "--numstat", "--since="+fromCommitDate)
}
