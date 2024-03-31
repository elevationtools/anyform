
package anyform

// Currently doesn't work because gomplate's config struct and config file
// parser are "internal".

/*

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	gomplate "github.com/hairyhenderson/gomplate/v4"
	gomplateconfig "github.com/hairyhenderson/gomplate/v4/internal/config"
)

type GomplateLibStageStamper struct {
  globe *Globe
}

func NewGomplateLibStageStamper(globe *Globe) *GomplateLibStageStamper {
  return &GomplateLibStageStamper{globe: globe}
}

func (gc *GomplateLibStageStamper) Stamp(
		ctx context.Context, inputDir string, outputDir string) error {
	var cfg *gomplateconfig.Config

	configFile, err := os.Open(filepath.Join(inputDir, ".gomplate.yaml"))
	if err == nil {
		defer configFile.Close()
		cfg, err = gomplateconfig.Parse(configFile)
		if err != nil { return fmt.Errorf("Parsing .gomplate.yml: %w", err) }
	} else {
		cfg = &gomplateconfig.Config{}
		cfg.ApplyDefaults()
	}

	cfg.InputDir = inputDir
	cfg.OutputDir = outputDir

	gomplate.Run(ctx, cfg)
	return nil
}

*/
