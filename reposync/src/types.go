package reposync

type GitLog struct {
	CommitHash    string
	AuthorName    string
	AuthorEmail   string
	AuthorDate    int64
	CommitDate    int64
	CommitMessage string
	FilesChanged  int64
	Insertions    int64
	Deletions     int64
}

type CommitObject struct {
	CommitHash    []string
	Author        []string
	AuthorEmail   []string
	AuthorDate    []int64
	CommitterDate []int64
	Message       []string
	Insertions    []int32
	Deletions     []int32
	FilesChanged  []int32
	RepoUrls      []string
}

type DependencyWithUpdatedTime struct {
	DependencyName string
	DependencyFile string
	UpdatedAt      int64
	InternalID     int32
}

const organization = "openqdev"
const repo = "openq-drm-testrepo"
