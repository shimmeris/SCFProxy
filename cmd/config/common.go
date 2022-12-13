package config

import (
	"encoding/json"
	"os"
)

func save(v interface{}, path string) error {
	data, _ := json.MarshalIndent(v, "", "    ")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}
