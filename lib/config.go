
package anyform

import (
  "os"
  "sync"
)

type OrchestratorConfig struct {
  GenfilesDir string `json:"genfiles_dir"`
  OutputDir string `json:"output_dir"`
  Interactive bool `json:"interactive"`
}

type AnyformConfig struct {
  Jsonnet string
  JsonnetDeps string
  OrchestratorSpecFile string
  Orchestrator OrchestratorConfig
}

var defaultConfig *AnyformConfig
var initDefaultConfigOnce sync.Once

func Getenv(envVar string, defaultValue string) string {
  value := os.Getenv(envVar)
  if value == "" {
    return defaultValue
  } else {
    return value
  }
}

func DefaultConfig() *AnyformConfig {
  initDefaultConfigOnce.Do(func() {
    defaultConfig = &AnyformConfig{
      Orchestrator: OrchestratorConfig{
        GenfilesDir: Getenv("ANYFORM_GENFILES_DIR", "genfiles"),
        OutputDir: Getenv("ANYFORM_OUTPUT_DIR", "output"),
        Interactive: true,
      },
      OrchestratorSpecFile: "anyform.jsonnet",
      Jsonnet: Getenv("JSONNET", "jsonnet"),
      JsonnetDeps: Getenv("JSONNET_DEPS", "jsonnet-deps"),
    }
  })
  
  return defaultConfig
}

