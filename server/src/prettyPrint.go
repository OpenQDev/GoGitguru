package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func PrettyPrint(i interface{}) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}

func PrintResponseBody(resp *http.Response) *http.Response {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	_ = string(bodyBytes)

	// logger.LogGreenDebug("response body %s", responseBody)

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return resp
}
