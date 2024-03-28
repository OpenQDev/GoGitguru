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
type RepoDependencyHistoryObject struct {
	DependencyId     []int32
	DateFirstPresent []int64
	DateLastRemoved  []int64
	RepoUrls         []string
}
