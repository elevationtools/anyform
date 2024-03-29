
package anyform

import (
  //"fmt"
	"log/slog"
)

type Anyform struct {
  locator *Locator
}

func NewDefaultAnyform() *Anyform {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	loc := NewDefaultLocator()
  return &Anyform{
    locator: loc,
  }
}

func (af* Anyform) NewOrchestrator() (*Orchestrator, error) {
  return NewOrchestrator(af.locator)
}
