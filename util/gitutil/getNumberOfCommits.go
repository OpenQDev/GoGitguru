package gitutil

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GetNumberOfCommits(prefixPath string, organization string, repo string, startDate time.Time) (int, error) {
	fullRepoPath := filepath.Join(prefixPath, organization, repo)

	cmd := exec.Command("git", "-C", fullRepoPath, "rev-list", "--count", "--since", startDate.String(), "HEAD")

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
