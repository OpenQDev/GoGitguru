package gitutil

import (
	"os/exec"
)

/*
git clone https://github.com/OpenQDev/OpenQ-DRM repos/OpenQ-DRM
*/
func GitCloneCommand(cloneString string, cloneDestination string) *exec.Cmd {
	return exec.Command("git", "clone", cloneString, cloneDestination)
}

/*
git -C . rev-parse --is-inside-work-tree
*/
func CheckIsAGitRepo(fullRepoPath string) *exec.Cmd {
	return exec.Command("git", "-C", fullRepoPath, "rev-parse", "--is-inside-work-tree")
}

/*
git -C . log --reverse --pretty=format":%H-;-%an-;-%ae-;-%at-;-%ct%n%s" --numstat
*/
func GitLogCommand(fullRepoPath string, fromCommitDate string) *exec.Cmd {
	return exec.Command("git", "-C", fullRepoPath, "log", "--reverse", "--pretty=format:%H-;-%an-;-%ae-;-%at-;-%ct%n%s", "--numstat", "--since="+fromCommitDate)
}

/*
git -C . log -n 20 --reverse --pretty=format":%H-;-%an-;-%ae-;-%at-;-%ct%n%s" --numstat
*/
func GitLog20Command(fullRepoPath string, fromCommitDate string) *exec.Cmd {
	return exec.Command("git", "-C", fullRepoPath, "log", "-n", "20", "--reverse", "--pretty=format:%H-;-%an-;-%ae-;-%at-;-%ct%n%s", "--numstat", "--since="+fromCommitDate)
}
