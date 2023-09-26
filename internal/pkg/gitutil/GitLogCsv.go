package gitutil

// git -C . log --reverse --pretty=format:'%H-;-%an-;-%ae-;-%at-;-%ct%n%s' --numstat --since=2020-01-01

import (
	"log"
	"os/exec"
	"path/filepath"
)

func GitLogCsv(prefixPath string, repo string, fromCommitDate string) []byte {
	if fromCommitDate == "" {
		fromCommitDate = "2020-01-01"
	}

	fullRepoPath := filepath.Join(prefixPath, repo)

	cmd := exec.Command("git", "-C", fullRepoPath, "log", "--reverse", "--pretty=format:'%H-;-%an-;-%ae-;-%at-;-%ct%n%s'", "--numstat", "--since="+fromCommitDate)

	out, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	return out

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
