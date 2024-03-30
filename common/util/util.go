
package util

import (
	"os"
	"encoding/json"
)

func ToJSONBytes(value any) ([]byte, error) {
  res, err := json.MarshalIndent(value, "", "  ")
  if err != nil { return nil, err }
  return res, nil
}

func ToJSONString(value any) (string, error) {
	res, err := ToJSONBytes(value)
	return string(res), err
}

func ToJSONFile(value any, filePath string) error {
	data, err := ToJSONBytes(value)
	if err != nil { return err }
	return os.WriteFile(filePath, data, 0660)
}

func FromJSONBytes[T any](jsonBytes []byte, out T) error {
  return json.Unmarshal(jsonBytes, out)
}

func FromJSONString[T any](jsonString string, out T) error {
  return FromJSONBytes([]byte(jsonString), out)
}

func FromJSONFile[T any](filePath string, out T) error {
	data, err := os.ReadFile(filePath) 
	if err != nil { return err }
	return FromJSONBytes(data, out)
}

