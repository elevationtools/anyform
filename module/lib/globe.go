// Globe is a combination of the service locator and dependency injection
// patterns. Compared to pure dependency injection:
// - Advantage: Faster writing of code because it allows injecting dependencies
//   without having to weave multiple individual dependencies through the
//   initialization path.
// - Disadvantage: It's harder to determine if a given piece of code uses a
//   given dependency or not.
// Compared to direct importing:
// - Advantage: Actually allows dependency injection.
// - Advantage: Easily replace the implementation with another.
// - Disadvantage: same as above.
package anyform

import (
	"context"
	"os/exec"
)

// TODO: add context.Context to all (most?) interfaces here.

type ConfigLoader interface {
	// Load inputFilePath into out (which should obviously be a pointer type).
	Load(inputFilePath string, out any) error
	// Returns inputFilePath and all of its transitive dependencies.  The max
	// modification time of this list is used to determine if anything has changed
	// since the last run.
  GetTransitiveDeps(inputFilePath string) ([]string, error)
}

// inputDir and outputDir can be relative paths.
// Implementations must set PWD=inputDir during template evaluation, which means
// other paths (like outputDir) may need to be converted to absolute paths.
type StageStamper interface {
  Stamp(ctx context.Context, stageName string, inputDir string,
			 outputDir string, logDir string, envVars []string) error
}

type SubprocessRunner interface {
	RunCmd(name string, cmd *exec.Cmd, logDir string) error
}

type Globe struct {
	Config* AnyformConfig
	StageStamper StageStamper
	ConfigLoader ConfigLoader
	SubprocessRunner SubprocessRunner 
}

func NewDefaultGlobe() *Globe {
	globe := &Globe{}
	globe.Config = NewDefaultAnyformConfig()
	globe.StageStamper = NewGomplateCliStageStamper(globe)
	globe.ConfigLoader = NewJsonnetLibConfigLoader(globe)
	globe.SubprocessRunner = NewDefaultSubprocessRunner(globe)
	return globe
}

