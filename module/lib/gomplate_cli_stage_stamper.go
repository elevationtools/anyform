
package anyform

import (
	"context"
	"os"
  "os/exec"
	"path/filepath"
)

type GomplateCliStageStamper struct {
  globe *Globe
}

func NewGomplateCliStageStamper(globe *Globe) *GomplateCliStageStamper {
  return &GomplateCliStageStamper{globe: globe}
}

func (gc *GomplateCliStageStamper) Stamp(
		ctx context.Context, inputDir string, outputDir string) error {
  cmd := exec.Command(gc.globe.Config.Gomplate,
      "--context", "cfg=" + gc.globe.Config.Orchestrator.ConfigJsonFile,
      "--input-dir=" + inputDir,
      "--output-dir", outputDir)

	// If the stage template directory contains a .gomplate.yaml file then
  // set GOMPLATE_CONFIG to point at it.
	configPath := filepath.Join(inputDir, ".gomplate.yaml")
	file, err := os.Open(configPath)
	if err == nil {
		file.Close()
	  cmd.Env = append(cmd.Environ(), "GOMPLATE_CONFIG=" + configPath)
	}

  outBytes, err := cmd.CombinedOutput()
  if err != nil {
    return Errorf("Error(%w) output: %v", err, string(outBytes))
  }
	return nil
}
