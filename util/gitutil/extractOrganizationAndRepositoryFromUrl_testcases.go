package gitutil

type ExtractOrganizationAndRepositoryFromUrlTest struct {
	name string
	url  string
	org  string
	repo string
}

func validUrlTest() ExtractOrganizationAndRepositoryFromUrlTest {
	const VALID = "VALID"
	return ExtractOrganizationAndRepositoryFromUrlTest{
		name: VALID,
		url:  "https://github.com/org/repo",
		org:  "org",
		repo: "repo",
	}
}

func urlWithoutRepoTest() ExtractOrganizationAndRepositoryFromUrlTest {
	const NO_REPO = "NO_REPO"
	return ExtractOrganizationAndRepositoryFromUrlTest{
		name: NO_REPO,
		url:  "https://github.com/org/",
		org:  "org",
		repo: "",
	}
}

func urlWithoutOrgAndRepoTest() ExtractOrganizationAndRepositoryFromUrlTest {
	const NO_ORG_OR_REPO = "NO_ORG_OR_REPO"
	return ExtractOrganizationAndRepositoryFromUrlTest{
		name: NO_ORG_OR_REPO,
		url:  "https://github.com/",
		org:  "",
		repo: "",
	}
}

func ExtractOrganizationAndRepositoryFromUrlTestCases() []ExtractOrganizationAndRepositoryFromUrlTest {
	return []ExtractOrganizationAndRepositoryFromUrlTest{
		validUrlTest(),
		urlWithoutRepoTest(),
		urlWithoutOrgAndRepoTest(),
	}
}
