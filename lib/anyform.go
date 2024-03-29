
package anyform

import (
  //"fmt"
	"log/slog"
)

type Anyform struct {
  globe *Globe
}

func NewDefaultAnyform() *Anyform {
	slog.SetLogLoggerLevel(slog.LevelDebug)
  return &Anyform{
    globe: NewDefaultGlobe(),
  }
}

func (af* Anyform) NewOrchestrator() (*Orchestrator, error) {
  return NewOrchestrator(af.globe)
}
