
package util

import (
	"encoding/json"
)

func ToJSONString(value any) (string, error) {
  res, err := json.MarshalIndent(value, "", "  ")
  if err != nil { return "", err }
  return string(res), nil
}

func FromJSONBytes[T any](jsonBytes []byte, out T) error {
  return json.Unmarshal(jsonBytes, out)
}
