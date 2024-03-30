
package anyform

import (
	"fmt"
  "os"
  "path/filepath"

	jsonnet "github.com/google/go-jsonnet"

	"github.com/elevationtools/anyform/common/util"
)

type JsonnetLibConfigLoader struct {
	globe *Globe
  vm *jsonnet.VM
}

func NewJsonnetLibConfigLoader(globe *Globe) *JsonnetLibConfigLoader {
	jl := &JsonnetLibConfigLoader{
    globe: globe,
    vm: jsonnet.MakeVM(),
  }

	jl.vm.Importer(&jsonnet.FileImporter{
		JPaths: filepath.SplitList(os.Getenv("JSONNET_PATH")),
	})

  return jl
}

func (jl *JsonnetLibConfigLoader) Load(inputFilePath string, out any) error {
  jsonString, err := jl.vm.EvaluateFile(inputFilePath)
  if err != nil { return fmt.Errorf("jsonnet EvaluateFile: %w", err) }
  return util.FromJSONString(jsonString, out)
}

