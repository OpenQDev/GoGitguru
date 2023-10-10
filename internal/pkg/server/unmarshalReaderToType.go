package server

import (
	"encoding/json"
	"io"
	"main/internal/pkg/logger"
)

// UnmarshalReader reads data from the provided io.Reader and unmarshals it into the given target.
// The target must be a pointer to a variable where the unmarshaled data should be stored.
// Returns an error if reading from the io.Reader fails or if the unmarshaling process encounters an error.
func UnmarshalReader[T any](r io.Reader, target *T) {
	bodyBytes, err := io.ReadAll(r)
	if err != nil {
		logger.LogFatalRedAndExit("failed to marshal response to %T: %s", target, err)
	}

	if err := json.Unmarshal(bodyBytes, target); err != nil {
		logger.LogFatalRedAndExit("failed to marshal response to %T: %s", target, err)
	}
}
