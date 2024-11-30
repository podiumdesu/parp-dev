package fileutil

import (
	"encoding/json"
	"io/ioutil"
)

func ReadFile(filename string, v interface{}) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, v); err != nil {
		return nil, err
	}
	return data, nil
}

// WriteFile marshals the provided struct into JSON data and writes it to a file.
func WriteFile(filename string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}
	return nil
}
