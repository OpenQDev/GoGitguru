package server

import (
	"bytes"
	"encoding/json"
	"io"
	"main/internal/pkg/logger"
)

func TypeToReader[T any](source T) io.Reader {
	jsonData, err := json.Marshal(source)
	if err != nil {
		logger.LogFatalRedAndExit("failed to marshal response to %T: %s", source, err)
	}
	return bytes.NewReader(jsonData)
}
