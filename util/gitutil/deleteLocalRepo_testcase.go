package gitutil

type DeleteLocalRepoTest struct {
	name         string
	organization string
	repo         string
}

func successfulDelete() DeleteLocalRepoTest {
	return DeleteLocalRepoTest{
		name:         "Successful delete",
		organization: "testOrganization",
		repo:         "testRepo",
	}
}

func DeleteLocalRepoTestCases() []DeleteLocalRepoTest {
	return []DeleteLocalRepoTest{
		successfulDelete(),
	}
}
