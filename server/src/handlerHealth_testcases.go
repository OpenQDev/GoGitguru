package server

type HandlerHealthTest struct {
	name               string
	expectedStatus     int
	expectedReturnBody HandlerHealthResponse
}

func shouldReturn200AndEmptyStruct() HandlerHealthTest {
	const SHOULD_RETURN_200_AND_EMPTY_STRUCT = "SHOULD_RETURN_200_AND_EMPTY_STRUCT"
	successReturnBody := HandlerHealthResponse{}

	return HandlerHealthTest{
		name:               SHOULD_RETURN_200_AND_EMPTY_STRUCT,
		expectedStatus:     200,
		expectedReturnBody: successReturnBody,
	}
}

func HandlerHealthTestCases() []HandlerHealthTest {
	return []HandlerHealthTest{
		shouldReturn200AndEmptyStruct(),
	}
}
