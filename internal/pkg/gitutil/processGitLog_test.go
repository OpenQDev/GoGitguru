package gitutil

import (
	"testing"
)

func TestProcessGitLog(t *testing.T) {
	testString := `141bd5216b3e95c0559de5704e97883214048e94-;-FlacoJones-;-andrew@openq.dev-;-1695429111-;-1695429111
initial commit - adds README.md
1       0       .gitignore
1       0       README.md
9       0       go.mod
40      0       go.sum
81      0       main.go`

	output := ProcessGitLog(testString)

	expected := GitLog{
		CommitHash:    "141bd5216b3e95c0559de5704e97883214048e94",
		AuthorName:    "FlacoJones",
		AuthorEmail:   "andrew@openq.dev",
		AuthorData:    "1695429111",
		CommitDate:    "1695429111",
		CommitMessage: "initial commit - adds README.md",
	}

	if output != expected {
		t.Errorf("Expected %v, but got %v", expected, output)
	}
}
