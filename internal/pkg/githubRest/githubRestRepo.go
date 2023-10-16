package githubRest

type GithubRestRepo struct {
	GithubRestID    int    `json:"id"`
	GithubGraphqlID string `json:"node_id"`
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	Private         bool   `json:"private"`
	Owner           struct {
		Login      string `json:"login"`
		ID         int    `json:"id"`
		NodeID     string `json:"node_id"`
		AvatarURL  string `json:"avatar_url"`
		GravatarID string `json:"gravatar_id"`
		URL        string `json:"url"`
	} `json:"owner"`
	Description     string `json:"description"`
	Fork            bool   `json:"fork"`
	URL             string `json:"url"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	PushedAt        string `json:"pushed_at"`
	Homepage        string `json:"homepage"`
	Size            int    `json:"size"`
	StargazersCount int    `json:"stargazers_count"`
	WatchersCount   int    `json:"watchers_count"`
	Language        string `json:"language"`
	HasIssues       bool   `json:"has_issues"`
	HasProjects     bool   `json:"has_projects"`
	HasDiscussions  bool   `json:"has_discussions"`
	ForksCount      int    `json:"forks_count"`
	Archived        bool   `json:"archived"`
	Disabled        bool   `json:"disabled"`
	OpenIssuesCount int    `json:"open_issues_count"`
	License         struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	} `json:"license"`
	Forks         int    `json:"forks"`
	OpenIssues    int    `json:"open_issues"`
	Watchers      int    `json:"watchers"`
	Visibility    string `json:"visibility"`
	DefaultBranch string `json:"default_branch"`
}
