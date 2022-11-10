package utils

import "encoding/json"

type name struct{}

func Unmarshal[T any](data []byte, result *T) error {
	return json.Unmarshal(data, result)
}
