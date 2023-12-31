package marshaller

import (
	"encoding/json"
	"os"
)

// JsonFileToType reads a JSON file and decodes its contents into the given target type.
func JsonFileToType[T any](jsonFile *os.File, target *T) error {
	defer jsonFile.Seek(0, 0)
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(target); err != nil {
		return err
	}

	return nil
}

// JsonFileToArrayOfType reads a JSON file and decodes its contents into the given target slice.
func JsonFileToArrayOfType[T any](jsonFile *os.File, target *[]T) error {
	defer jsonFile.Seek(0, 0)
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(target); err != nil {
		return err
	}

	return nil
}
