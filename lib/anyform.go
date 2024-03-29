
package anyform

type Anyform struct {
  globe *Globe
}

func NewDefaultAnyform() *Anyform {
  return &Anyform{
    globe: NewDefaultGlobe(),
  }
}

func (af* Anyform) NewOrchestrator() (*Orchestrator, error) {
  return NewOrchestrator(af.globe)
}

