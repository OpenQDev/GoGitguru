package githubGraphQL

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type GraphQLPayload struct {
	Query string `json:"query"`
}

func createGraphQLRequest(url string, query string, headers map[string]string) (*http.Request, error) {
	payloadObj := GraphQLPayload{Query: query}
	payloadBytes, err := json.Marshal(payloadObj)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return req, nil
}
