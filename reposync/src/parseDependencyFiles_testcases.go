package reposync

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type ParseFileTestCase struct {
	file     *object.File
	fileName string
}

// const list of files to test
func validParseFileTest() []ParseFileTestCase {
	FILE_LIST := []string{
		"package.json",
		".cabal",
		"requirements.txt",
		"pom.xml",
		"build.gradle",
		"composer.json",
		"go.mod",
		"Pipfile",
		"Gemfile",
		"Cargo.toml",
	}

	// get blob from file
	// get file using go-git
	thisRepository, err := git.PlainOpen("../../")
	if err != nil {
		panic(err)
	}

	// get package.json files in this repo
	// get the files from the commit

	// ... just iterates over the commits, printing it
	head, err := thisRepository.Head()
	if err != nil {
		panic(err)
	}
	fmt.Println(head.Hash())
	commit, err := thisRepository.CommitObject(head.Hash())
	if err != nil {
		panic(err)
	}
	testFiles, err := commit.Files()
	testFiles.ForEach(func(file *object.File) error {
		fmt.Println(file.Name)
		return nil
	})
	newFiles := make([]ParseFileTestCase, 0)
	for _, fileName := range FILE_LIST {
		path := fmt.Sprintf("reposync/src/mock/%s", fileName)
		fmt.Println(path)
		file, err := commit.File(path)
		if err != nil {
			panic(err)
		}
		newFiles = append(newFiles,
			ParseFileTestCase{
				file:     file,
				fileName: file.Name,
			})
	}
	if err != nil {
		panic(err)
	}

	//blob to interface

	return newFiles
}

func ParseFileTestCases() []ParseFileTestCase {
	return validParseFileTest()

}
