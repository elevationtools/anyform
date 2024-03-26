
package anyform

import (
	"encoding/json"
  "fmt"
  "os/exec"
)

func Must[T any](result T, err error) T {
  if err != nil { panic(err) }
  return result
}

func Must1(err error) {
  if err != nil { panic(err) }
}

func ToJSONString(value any) string {
  return string(Must(json.MarshalIndent(value, "", "  ")))
}

func FromJSONBytes[T any](jsonBytes []byte, out T) error {
  return json.Unmarshal(jsonBytes, &out)
}

type Util struct {
  Config *AnyformConfig
}

func NewUtil(config *AnyformConfig) *Util {
  return &Util{
    Config: config,
  }
}

func (u *Util) LoadJsonnetFile(path string, out any) error {
  outBytes, err := exec.Command(u.Config.Jsonnet, path).CombinedOutput()
  if err != nil {
    fmt.Printf("Error output: %v", string(outBytes))
    return err
  }
  return FromJSONBytes(outBytes, out)
}