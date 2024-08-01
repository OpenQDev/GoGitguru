package reposync

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type ParseFileTestCase struct {
	file         *object.File
	fileName     string
	dependencies []string
}

// const list of files to test
func validParseFileTest() []ParseFileTestCase {

	NO_FILE_LIST := []ParseFileTestCase{
		{
			file:     nil,
			fileName: "package.json",
			dependencies: []string{
				"find-config",
				"hardhat",
				"@nomiclabs/hardhat-ethers",
				"aurora",
				"prettier",
				"prettier-plugin-solidity",
				"near",
			},
		},
		{
			file:     nil,
			fileName: ".cabal",
			dependencies: []string{
				"base",
				"containers",
				"mtl",
				"transformers",
				"near",
			},
		}, {
			file:     nil,
			fileName: "build.gradle",
			dependencies: []string{
				"org.apache.commons:commons-math",
				"junit:junit",
				"com.google.guava:guava",
				"org.slf4j:slf4j-api",
				"org.apache.commons:commons-lang",
				"com.fasterxml.jackson.core:jackson-databind",
				"near",
			},
		},
		{
			file:     nil,
			fileName: "Cargo.toml",
			dependencies: []string{
				"serde",
				"serde_json",
				"reqwest",
				"tokio",
				"near",
			},
		},
		{
			file:     nil,
			fileName: "composer.json",
			dependencies: []string{
				"symfony/property-info",
				"php",
				"ext-json",
				"near",
			},
		},
		{
			file:     nil,
			fileName: "Gemfile",
			dependencies: []string{
				"sqlite3",
				"puma",
				"sass-rails",
				"near",
			},
		},
		{
			file:     nil,
			fileName: "go.mod",
			dependencies: []string{
				"github.com/stretchr/testify",
				"github.com/OpenQDev/GoGitguru",
				"near",
			},
		},
		{
			file:     nil,
			fileName: "Pipfile",
			dependencies: []string{
				"pip-review-req-multi",
				"pip-sync-req-multi",
				"setuptools",
				"near",
			},
		}, {
			file:     nil,
			fileName: "requirements.txt",
			dependencies: []string{
				"example",
				"example2",
				"example3",
			},
		}, {
			file:     nil,
			fileName: "pom.xml",
			dependencies: []string{
				"junit:junit",
				"org.apache.commons:commons-math",
				"org.apache.maven.plugins:maven-compiler-plugin",
				"near",
			},
		},
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
	commit, err := thisRepository.CommitObject(head.Hash())
	if err != nil {
		panic(err)
	}
	testFiles, err := commit.Files()
	testFiles.ForEach(func(file *object.File) error {
		return nil
	})
	newFiles := make([]ParseFileTestCase, 0)
	for _, mockFile := range NO_FILE_LIST {

		path := fmt.Sprintf("reposync/mock/%s", mockFile.fileName)

		file, err := commit.File(path)
		if err != nil {
			panic(err)
		}

		newFiles = append(newFiles,
			ParseFileTestCase{
				file:         file,
				fileName:     mockFile.fileName,
				dependencies: mockFile.dependencies,
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
