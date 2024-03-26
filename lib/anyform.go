
package anyform

import (
  "fmt"
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

func (af* Anyform) Up() {
  fmt.Println("lib-anyform Up!")
  orc := NewOrchestrator(af.Config, af.Util)
  fmt.Println(ToJSONString(orc.Spec))
}

func (af* Anyform) Down() {
  fmt.Println("lib-anyform Down!")
}

func (af* Anyform) NewOrchestrator() *Orchestrator {
	orc := &Orchestrator{}
  Must1(af.Util.LoadJsonnetFile(af.Config.OrchestratorSpecFile, &orc.Spec))
  fmt.Println(ToJSONString(orc.Spec))
	return orc
}
