package util

import (
	"bytes"
	"encoding/json"
	"io"
)

func TypeToReader[T any](source T) (io.Reader, error) {
	jsonData, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(jsonData), nil
}
