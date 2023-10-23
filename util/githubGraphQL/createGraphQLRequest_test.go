package githubGraphQL

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGraphQLRequest(t *testing.T) {
	url := "https://test-url.com"
	query := "{testQuery}"
	headers := map[string]string{
		"Authorization": "Bearer testToken",
		"Content-Type":  "application/json",
	}

	req, err := createGraphQLRequest(url, query, headers)

	assert.NoError(t, err)
	assert.Equal(t, http.MethodPost, req.Method)
	assert.Equal(t, url, req.URL.String())
	assert.Equal(t, "Bearer testToken", req.Header.Get("Authorization"))
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

	// Additional checks can be added here to validate the request body
	// Check if the payload is what was passed in
	body, _ := io.ReadAll(req.Body)
	assert.Equal(t, `{"query":"{testQuery}"}`, string(body))
}
