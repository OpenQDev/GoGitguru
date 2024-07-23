package marshaller

import (
	"bytes"
	"encoding/json"
)

func JsonToArrayOfType[T any](jsonBytes []byte, target *T) error {
	decoder := json.NewDecoder(bytes.NewReader(jsonBytes))

	for decoder.More() {
		if err := decoder.Decode(&target); err != nil {
			return err
		}
	}
	return nil
}
