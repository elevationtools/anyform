
package anyform

import (
	"fmt"
	"os/exec"

	"github.com/elevationtools/anyform/common/util"
)

type JsonnetRunner interface {
	Run(inputFilePath string, out any) error
}

type CliJsonnetRunner struct {
	locator *Locator
}

func NewCliJsonnetRunner(locator *Locator) *CliJsonnetRunner {
	return &CliJsonnetRunner{locator: locator}
}

func (cgr *CliJsonnetRunner) Run(inputFilePath string, out any) error {
  outBytes, err := exec.Command(cgr.locator.Config.Jsonnet, inputFilePath).
			CombinedOutput()
  if err != nil {
    return fmt.Errorf("Error(%w) output: %v", err, string(outBytes))
  }
  return util.FromJSONBytes(outBytes, out)
}
