package server

type User struct {
	InternalID      int    `json:"internal_id"`
	GithubRestID    int    `json:"github_rest_id"`
	GithubGraphqlID string `json:"github_graphql_id"`
	Login           string `json:"login"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	AvatarURL       string `json:"avatar_url"`
	Company         string `json:"company"`
	Location        string `json:"location"`
	Bio             string `json:"bio"`
	Blog            string `json:"blog"`
	Hireable        bool   `json:"hireable"`
	TwitterUsername string `json:"twitter_username"`
	Followers       int    `json:"followers"`
	Following       int    `json:"following"`
	Type            string `json:"type"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}
