
package anyform

import (
  "os"
	"path/filepath"
)

type OrchestratorConfig struct {
  GenfilesDir string `json:"genfiles_dir"`
  OutputDir string `json:"output_dir"`
  ConfigJsonFile string `json:"config_json_file"`
}

type AnyformConfig struct {
  Orchestrator OrchestratorConfig
  OrchestratorSpecFile string
  Jsonnet string
  JsonnetDeps string
  Gomplate string
	Interactive bool
}

func Getenv(envVar string, defaultValue string) string {
  value := os.Getenv(envVar)
  if value == "" {
    return defaultValue
  } else {
    return value
  }
}

func NewDefaultAnyformConfig() *AnyformConfig {
  ac := &AnyformConfig{}

	ac.Orchestrator.GenfilesDir = Getenv("ANYFORM_GENFILES_DIR", "genfiles")
	ac.Orchestrator.OutputDir = Getenv("ANYFORM_OUTPUT_DIR", "output")
	ac.Orchestrator.ConfigJsonFile = Getenv("ANYFORM_CONFIG_JSON_FILE",
			filepath.Join(ac.Orchestrator.GenfilesDir, "config.json"))

	ac.OrchestratorSpecFile = "anyform.jsonnet"
	ac.Jsonnet = Getenv("JSONNET", "jsonnet")
	ac.JsonnetDeps = Getenv("JSONNET_DEPS", "jsonnet-deps")
	ac.Gomplate = Getenv("GOMPLATE", "gomplate")
	// Only exactly "false" will make it non-interactive.
	ac.Interactive = Getenv("INTERACTIVE", "true") != "false"

	return ac
}

const CtlFileName = "ctl"
