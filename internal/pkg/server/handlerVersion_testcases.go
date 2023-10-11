package server

type HandlerVersionTest struct {
	name                 string
	expectedStatusCode   int
	expectedResponseBody HandlerVersionResponse
}

func shouldReturn200AndCorrectVersion() HandlerVersionTest {
	const SHOULD_RETURN_200_AND_CORRECT_VERSION = "should return 200 and version 1.0.0"
	successResponse := HandlerVersionResponse{Version: "1.0.0"}

	return HandlerVersionTest{
		name:                 SHOULD_RETURN_200_AND_CORRECT_VERSION,
		expectedStatusCode:   200,
		expectedResponseBody: successResponse,
	}
}

func HandlerVersionTestCases() []HandlerVersionTest {
	return []HandlerVersionTest{
		shouldReturn200AndCorrectVersion(),
	}
}
