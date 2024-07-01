package reposync

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestAddFirstLastCommitDateByEmail(t *testing.T) {
	usersToRepoUrl := UsersToRepoUrl{
		AuthorEmails:     []string{"christopher.stevers1@gmail.com"},
		FirstCommitDates: []int64{1231231231203},
		LastCommitDates:  []int64{1231231231203},
	}

	commits := []object.Commit{

		{Author: object.Signature{
			Name:  "author",
			Email: "christopher.stevers1@gmail.com",
			When:  time.Unix(1231231231203, 0),
		},
		},

		{Author: object.Signature{
			Name:  "author",
			Email: "christopher.stevers1@gmail.com",
			When:  time.Unix(5, 0),
		}},

		{Author: object.Signature{
			Name:  "author",
			Email: "test.test@test.com",
			When:  time.Unix(23425, 0),
		}},
	}

	for _, commit := range commits {
		AddFirstLastCommitDateByEmail(&usersToRepoUrl, &commit)
	}

	expectedOutput := UsersToRepoUrl{
		AuthorEmails:     []string{"christopher.stevers1@gmail.com", "test.test@test.com"},
		FirstCommitDates: []int64{5, 23425},
		LastCommitDates:  []int64{1231231231203, 23425},
	}

	if !reflect.DeepEqual(usersToRepoUrl, expectedOutput) {
		t.Errorf("Expected %v, got %v", expectedOutput, usersToRepoUrl)
	}
}
