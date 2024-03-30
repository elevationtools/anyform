
package anyform

import (
  "fmt"
  "os/exec"
)

type GomplateCliStageStamper struct {
  globe *Globe
}

func NewGomplateCliStageStamper(globe *Globe) *GomplateCliStageStamper {
  return &GomplateCliStageStamper{globe: globe}
}

func (gc *GomplateCliStageStamper) Stamp(inputDir string, outputDir string) error {
  outBytes, err := exec.Command(gc.globe.Config.Gomplate,
      "--context", "cfg=" + gc.globe.Config.Orchestrator.ConfigJsonFile,
      "--input-dir", inputDir,
      "--output-dir", outputDir).
      CombinedOutput()
  if err != nil {
    return fmt.Errorf("Error(%w) output: %v", err, string(outBytes))
  }
	return nil
}
