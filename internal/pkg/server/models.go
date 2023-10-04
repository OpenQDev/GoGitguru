package server

type User struct {
	InternalID      int    `json:"internal_id"`
	GithubRestID    int    `json:"id"`
	GithubGraphqlID string `json:"node_id"`
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

type Repo struct {
	GithubRestID    int    `json:"id"`
	GithubGraphqlID string `json:"node_id"`
	URL             string `json:"url"`
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	Private         bool   `json:"private"`
	Owner           struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	} `json:"owner"`
	Description     string `json:"description"`
	Homepage        string `json:"homepage"`
	Fork            bool   `json:"fork"`
	ForksCount      int    `json:"forks_count"`
	Archived        bool   `json:"archived"`
	Disabled        bool   `json:"disabled"`
	License         string `json:"license"`
	Language        string `json:"language"`
	StargazersCount int    `json:"stargazers_count"`
	WatchersCount   int    `json:"watchers_count"`
	OpenIssuesCount int    `json:"open_issues_count"`
	HasIssues       bool   `json:"has_issues"`
	HasDiscussions  bool   `json:"has_discussions"`
	HasProjects     bool   `json:"has_projects"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	PushedAt        string `json:"pushed_at"`
	Visibility      string `json:"visibility"`
	Size            int    `json:"size"`
	DefaultBranch   string `json:"default_branch"`
}
