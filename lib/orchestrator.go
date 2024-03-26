
package anyform

import (
	//"fmt"
)

type OrchestratorSpec struct {
	ImplDir string `json:"impl_dir"`
	Config map[string]any `json:"config"`
}

type Orchestrator struct {
	Spec OrchestratorSpec
}

func NewOrchestrator(config *AnyformConfig, util *Util) *Orchestrator {
	orc := &Orchestrator{}
  Must1(util.LoadJsonnetFile(config.OrchestratorSpecFile, &orc.Spec))
	return orc
}
