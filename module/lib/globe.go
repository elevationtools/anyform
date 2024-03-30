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
	"os/exec"
)

type ConfigLoader interface {
	Load(inputFilePath string, out any) error
  GetTransitiveDeps(inputFilePath string) ([]string, error)
}

type StageStamper interface {
  Stamp(inputDir string, outputDir string) error
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

