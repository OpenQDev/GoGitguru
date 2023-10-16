package githubGraphQL

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	User  struct {
		GithubRestID    int    `json:"github_rest_id"`
		GithubGraphqlID string `json:"github_graphql_id"`
		Login           string `json:"login"`
		Name            string `json:"name"`
		Email           string `json:"email"`
		AvatarURL       string `json:"avatar_url"`
		Company         string `json:"company"`
		Location        string `json:"location"`
		Hireable        bool   `json:"hireable"`
		Bio             string `json:"bio"`
		Blog            string `json:"blog"`
		TwitterUsername string `json:"twitter_username"`
		Followers       struct {
			TotalCount int `json:"totalCount"`
		} `json:"followers"`
		Following struct {
			TotalCount int `json:"totalCount"`
		} `json:"following"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"user"`
}

type Commit struct {
	Author Author `json:"author"`
}

type CommitAuthorsResponse struct {
	Data *struct {
		RateLimit struct {
			Limit   int    `json:"limit"`
			Used    int    `json:"used"`
			ResetAt string `json:"resetAt"`
		} `json:"rateLimit"`
		Repository map[string]Commit `json:"repository"`
	} `json:"data"`
	Errors *[]struct {
		Path       []string `json:"path"`
		Extensions struct {
			Code         string `json:"code"`
			TypeName     string `json:"typeName"`
			ArgumentName string `json:"argumentName"`
		} `json:"extensions"`
		Locations []struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"locations"`
		Message string `json:"message"`
	} `json:"errors"`
}

type GraphQLPayload struct {
	Query string `json:"query"`
}
