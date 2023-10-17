package gitutil

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func GetNumberOfCommits(prefixPath string, repo string) (int, error) {
	fullRepoPath := filepath.Join(prefixPath, repo)
	cmd := exec.Command("git", "-C", fullRepoPath, "rev-list", "--count", "HEAD")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return 0, err
	}

	return strconv.Atoi(strings.TrimSpace(out.String()))
}
