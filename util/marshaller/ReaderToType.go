package marshaller

import (
	"encoding/json"
	"io"
)

// ReaderToType reads data from the provided io.Reader and unmarshals it into the given target.
// The target must be a pointer to a variable where the unmarshaled data should be stored.
// Returns an error if reading from the io.Reader fails or if the unmarshaling process encounters an error.
func ReaderToType[T any](r io.Reader, target *T) error {
	bodyBytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bodyBytes, target); err != nil {
		return err
	}

	return nil
}
