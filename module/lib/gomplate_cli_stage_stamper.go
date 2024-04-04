
package anyform

import (
	"context"
	"fmt"
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
		ctx context.Context, stageName string, inputDir string,
		outputDir string, logDir string, envVars []string) error {
  cmd := exec.Command(gc.globe.Config.Gomplate,
      "--context", "cfg=" + gc.globe.Config.Orchestrator.ConfigJsonFile,
      "--input-dir=" + inputDir,
      "--output-dir", outputDir)
  cmd.Env = append(cmd.Environ(), envVars...)

	cmd.Env = append(
		cmd.Environ(), "GOMPLATE_LOG_FORMAT=simple")
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
		// TODO(ux): reformat the hard-to-read gomplate golang text template error.
		logFilePath := filepath.Join(logDir, "stamp-stdout_stderr")
		err2 := WriteFile(logFilePath, outBytes, false)
		if err2 != nil {
			// Ugh, well at least print to stderr that we couldn't write the log file.
			// Then ignore the error and hope for the best.
			fmt.Fprintf(os.Stderr, "%v", err2)
		}
    return Errorf("running stamper: %w: output: %v", err, string(outBytes))
  }
	return nil
}
