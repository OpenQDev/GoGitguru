package util

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
)

func DumpJson(t *testing.T, rr *httptest.ResponseRecorder) {
	bodyBytes, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
}
