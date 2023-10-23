package server

type HandlerDependencyHistoryTestCase struct {
	name           string
	authorized     bool
	shouldError    bool
	expectedStatus int
}

func HandlerDependencyHistoryTestCases() []HandlerDependencyHistoryTestCase {
	return []HandlerDependencyHistoryTestCase{}
}
