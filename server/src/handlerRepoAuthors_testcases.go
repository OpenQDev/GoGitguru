package server

import (
	"log"
	"net/http"
	"os"
	"time"
	"github.com/lib/pq"
	"github.com/OpenQDev/GoGitguru/util/marshaller"

	"github.com/DATA-DOG/go-sqlmock"
)

type HandlerRepoAuthorsTestCase struct {
	name               string
	login              string
	expectedStatus     int
	requestBody        HandlerRepoAuthorsRequest
	expectedReturnBody []AuthorInfo
	authorized         bool
	shouldError        bool
	setupMock          func(mock sqlmock.Sqlmock)
}

func fooAuthors() HandlerRepoAuthorsTestCase {
	const login = "FlacoJones"
	const UNAUTHORIZED = "UNAUTHORIZED"
	return HandlerRepoAuthorsTestCase{
		name:               UNAUTHORIZED,
		login:              login,
		expectedStatus:     http.StatusBadRequest,
		authorized:         false,
		requestBody:        HandlerRepoAuthorsRequest{},
		expectedReturnBody: []AuthorInfo{},
		shouldError:        true,
		setupMock:          func(mock sqlmock.Sqlmock) {},
	}
}

func getAllRepoAuthors() HandlerRepoAuthorsTestCase {
	const login = "FlacoJones"
	const GET_ALL_REPO_AUTHORS = "GET_ALL_REPO_AUTHORS"

	since := "2023-07-14T13:45:32-05:00"
	until := "2024-07-25T13:45:32-05:00"

	requestBody := HandlerRepoAuthorsRequest{
		RepoUrls: []string{"https://github.com/OpenQDev/OpenQ-Frontend"},
		Since:   since,
		Until:   until,
	}

	var authorsResponse []AuthorInfo
	jsonFile, err := os.Open("./mocks/mockRepoAuthorsResponse.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	marshaller.JsonFileToArrayOfType(jsonFile, &authorsResponse)

	expectedReturnBody := authorsResponse

	return HandlerRepoAuthorsTestCase{
		name:               GET_ALL_REPO_AUTHORS,
		login:              login,
		expectedStatus:     http.StatusOK,
		authorized:         true,
		requestBody:        requestBody,
		expectedReturnBody: expectedReturnBody,
		shouldError:        false,
		setupMock: func(mock sqlmock.Sqlmock) {

			sinceTime, _ := time.Parse(time.RFC3339, since)
			untilTime, _ := time.Parse(time.RFC3339, until)
			sinceUnix := sinceTime.Unix()
			untilUnix := untilTime.Unix()

			// Define the mock rows
			mockRows := sqlmock.NewRows([]string{
				"author", "author_email", 
				"rest_id", "email", "internal_id", "github_rest_id", "github_graphql_id", "login", "name", "email_2", "avatar_url", "company", "location", "bio", "blog", "hireable", "twitter_username", "followers", "following", "type",
			})

			// Add rows to the mock rows
			c1 := authorsResponse[0]
			row1 := mockRows.AddRow(
				c1.Author, c1.AuthorEmail, 
				c1.RestID, c1.Email, c1.InternalID, c1.GithubRestID, c1.GithubGraphqlID, c1.Login, c1.Name, c1.Email_2, c1.AvatarUrl, c1.Company, c1.Location, c1.Bio, c1.Blog, c1.Hireable, c1.TwitterUsername, c1.Followers, c1.Following, c1.Type, 
			)

			c2 := authorsResponse[1]
			row2 := mockRows.AddRow(
				c2.Author, c2.AuthorEmail, 
				c2.RestID, c2.Email, c2.InternalID, c2.GithubRestID, c2.GithubGraphqlID, c2.Login, c2.Name, c2.Email_2, c2.AvatarUrl, c2.Company, c2.Location, c2.Bio, c2.Blog, c2.Hireable, c2.TwitterUsername, c2.Followers, c2.Following, c2.Type,
			)

			c3 := authorsResponse[2]
			row3 := mockRows.AddRow(
				c3.Author, c3.AuthorEmail, 
				c3.RestID, c3.Email, c3.InternalID, c3.GithubRestID, c3.GithubGraphqlID, c3.Login, c3.Name, c3.Email_2, c3.AvatarUrl, c3.Company, c3.Location, c3.Bio, c3.Blog, c3.Hireable, c3.TwitterUsername, c3.Followers, c3.Following, c3.Type,
			)

			// Expect the query with the mock rows
			mock.ExpectQuery("-- name: GetRepoAuthorsInfo :many").
				WithArgs(pq.Array(requestBody.RepoUrls), sinceUnix, untilUnix).
				WillReturnRows(row1, row2, row3)
		},
	}
}

func HandlerRepoAuthorsTestCases() []HandlerRepoAuthorsTestCase {
	return []HandlerRepoAuthorsTestCase{
		fooAuthors(),
		getAllRepoAuthors(),
	}
}
