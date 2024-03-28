
package anyform

import (
  //"fmt"
)

type Anyform struct {
  Config *AnyformConfig
  Util *Util
}

func NewAnyform() *Anyform {
  config := DefaultConfig()
  return &Anyform{
    Config: config,
    Util: NewUtil(config),
  }
}

func (af* Anyform) NewOrchestrator() (*Orchestrator, error) {
  return NewOrchestrator(af.Config, af.Util)
}
