package gitutil

type DeleteLocalRepoTest struct {
	name string
	repo string
}

func successfulDelete() DeleteLocalRepoTest {
	return DeleteLocalRepoTest{
		name: "Successful delete",
		repo: "testRepo",
	}
}

func DeleteLocalRepoTestCases() []DeleteLocalRepoTest {
	return []DeleteLocalRepoTest{
		successfulDelete(),
	}
}
