package gitutil

import (
	"reflect"
	"testing"
)

func TestProcessGitLogs(t *testing.T) {
	log1 := `141bd5216b3e95c0559de5704e97883214048e94-;-FlacoJones-;-andrew@openq.dev-;-1695429111-;-1695429111
initial commit - clones repo, TAR and GZIP .git directory, upload to S3, delete from local
1       0       .gitignore
1       0       README.md
9       0       go.mod
40      0       go.sum
81      0       main.go`

	log2 := `141bd5216b3e95c0559de5704e97883214048e95-;-FlacoJones-;-andrew@openq.dev-;-1695429111-;-1695429111
initial commit - clones repo, TAR and GZIP .git directory, upload to S3, delete from local
1       0       .gitignore
1       0       README.md
9       0       go.mod
40      0       go.sum
81      0       main.go`

	testString := log1 + "\n\n" + log2

	output := ProcessGitLogs(testString)

	expected := []GitLog{
		{
			CommitHash:    "141bd5216b3e95c0559de5704e97883214048e94",
			AuthorName:    "FlacoJones",
			AuthorEmail:   "andrew@openq.dev",
			AuthorDate:    1695429111,
			CommitDate:    1695429111,
			CommitMessage: "initial commit - clones repo, TAR and GZIP .git directory, upload to S3, delete from local",
			FilesChanged:  5,
		},
		{
			CommitHash:    "141bd5216b3e95c0559de5704e97883214048e95",
			AuthorName:    "FlacoJones",
			AuthorEmail:   "andrew@openq.dev",
			AuthorDate:    1695429111,
			CommitDate:    1695429111,
			CommitMessage: "initial commit - clones repo, TAR and GZIP .git directory, upload to S3, delete from local",
			FilesChanged:  5,
		},
	}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v, but got %v", expected, output)
	}
}
