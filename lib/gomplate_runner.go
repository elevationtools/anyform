
package anyform

import (
  "fmt"
  "os/exec"
)

type GomplateRunner interface {
  Run(inputDir string, outputDir string) error
}

type CliGomplateRunner struct {
  locator *Locator
}

func NewCliGomplateRunner(locator *Locator) *CliGomplateRunner {
  return &CliGomplateRunner{locator: locator}
}

func (cgr *CliGomplateRunner) Run(inputDir string, outputDir string) error {
  outBytes, err := exec.Command(cgr.locator.Config.Gomplate,
      "--context", "cfg=" + cgr.locator.Config.Orchestrator.ConfigJsonFile,
      "--input-dir", inputDir,
      "--output-dir", outputDir).
      CombinedOutput()
  if err != nil {
    return fmt.Errorf("Error(%w) output: %v", err, string(outBytes))
  }
	return nil
}
