package usersync

type UserSync struct {
	CommitHash string
	Author     struct {
		Email   string
		NotNull bool
	}
	Repo struct {
		URL     string
		NotNull bool
	}
}
