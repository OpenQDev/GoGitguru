package gitutil

import (
	"testing"
)

func TestProcessGitLog(t *testing.T) {

	testString1 := `141bd5216b3e95c0559de5704e97883214048e94-;-FlacoJones-;-andrew@openq.dev-;-1695429111-;-1695429111
initial commit - adds README.md
1       0       .gitignore
1       0       README.md
9       0       go.mod
40      0       go.sum
81      0       main.go`

	testString2 := `b8cc8b9e5252d470e559a55c841778cb99957cdb-;-mktcode-;-kontakt@markus-kottlaender.de-;-1609951502-;-1609951502
initial commit - adds README.md
1       0       .gitignore
1       0       README.md
9       0       go.mod
40      0       go.sum
81      0       main.go`

	tests := []struct {
		gitLogString   string
		expectedGitLog GitLog
	}{
		{
			gitLogString: testString1,
			expectedGitLog: GitLog{
				CommitHash:    "141bd5216b3e95c0559de5704e97883214048e94",
				AuthorName:    "FlacoJones",
				AuthorEmail:   "andrew@openq.dev",
				AuthorDate:    1695429111,
				CommitDate:    1695429111,
				CommitMessage: "initial commit - adds README.md",
				FilesChanged:  5,
				Insertions:    132,
				Deletions:     0,
			},
		},
		{
			gitLogString: testString2,
			expectedGitLog: GitLog{
				CommitHash:    "b8cc8b9e5252d470e559a55c841778cb99957cdb",
				AuthorName:    "mktcode",
				AuthorEmail:   "kontakt@markus-kottlaender.de",
				AuthorDate:    1609951502,
				CommitDate:    1609951502,
				CommitMessage: "initial commit - adds README.md",
				FilesChanged:  5,
				Insertions:    132,
				Deletions:     0,
			},
		},
	}

	for _, tt := range tests {
		output := ProcessGitLog(tt.gitLogString)

		expected := tt.expectedGitLog

		if output != expected {
			t.Errorf("Expected %v, but got %v", expected, output)
		}
	}
}
