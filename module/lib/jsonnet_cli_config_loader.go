// Not currently in use. Can probably be removed.
package anyform

import (
	"os/exec"

	"github.com/elevationtools/anyform/module/common/util"
)

type JsonnetCliConfigLoader struct {
	globe *Globe
}

func NewJsonnetCliConfigLoader(globe *Globe) *JsonnetCliConfigLoader {
	return &JsonnetCliConfigLoader{globe: globe}
}

func (jc *JsonnetCliConfigLoader) Load(inputFilePath string, out any) error {
  outBytes, err := exec.Command(jc.globe.Config.Jsonnet, inputFilePath).
			CombinedOutput()
  if err != nil {
    return Errorf("Error(%w) output: %v", err, string(outBytes))
  }
  return util.FromJSONBytes(outBytes, out)
}
