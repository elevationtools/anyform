
package anyform

import (
	"context"
  "fmt"
	"time"

  commonutil "github.com/elevationtools/anyform/common/util"
  daglib "github.com/elevationtools/anyform/lib/dag/mubased"
)

type StageSpec struct {
  DependsOn []string `json:"depends_on"`
}

type OrchestratorSpec struct {
  ImplDir string `json:"impl_dir"`
  Stages map[string]StageSpec `json:"stages"`
  Config map[string]any `json:"config"`
}

type Stage struct {
  Name string
  Spec *StageSpec
  DependsOn []*Stage
  RequiredBy []*Stage
}

func NewStage(name string, spec *StageSpec) *Stage {
  return &Stage{
    Name: name,
    Spec: spec,
  }
}

type Orchestrator struct {
  Spec OrchestratorSpec
  Stages map[string]*Stage
}

func NewOrchestrator(config *AnyformConfig, util *Util) (*Orchestrator, error) {
  orc := &Orchestrator{
    Stages: map[string]*Stage{},
  }
  err := util.LoadJsonnetFile(config.OrchestratorSpecFile, &orc.Spec)
  if err != nil { return nil, err }

  for stageName, iStageSpec := range orc.Spec.Stages {
		stageSpec := iStageSpec
    orc.Stages[stageName] = NewStage(stageName, &stageSpec)
  }

  for _, stage := range orc.Stages {
    for _, depName := range stage.Spec.DependsOn {
      dep, found := orc.Stages[depName]
      if !found {
        return nil, fmt.Errorf(
          "Invalid spec: stage %v depends on undefined stage %v", stage.Name, depName)
      }
      stage.DependsOn = append(stage.DependsOn, dep)
      dep.RequiredBy = append(dep.RequiredBy, stage)
    }
  }

  return orc, nil
}

func (orc* Orchestrator) Up() error {
  fmt.Println("lib-anyform Up!")
  fmt.Println(commonutil.Must(ToJSONString(orc.Spec)))

  dag := daglib.NewDag()

  for _, iStage := range orc.Stages {
		stage := iStage

    parentNames := []string{}
    for _, parentStage := range stage.DependsOn {
      parentNames = append(parentNames, parentStage.Name)
    }

    dag.AddVertex(stage.Name, parentNames, func (ctx context.Context) error {
      fmt.Printf("%v: bringing up\n", stage.Name)
			time.Sleep(1 * time.Second)
      fmt.Printf("%v: done\n", stage.Name)
      return nil
    })
  }

  return dag.Run(context.TODO())
}

func (orc* Orchestrator) Down() error {
  dag := daglib.NewDag()

  for _, iStage := range orc.Stages {
		stage := iStage

    childNames := []string{}
    for _, childStage := range stage.RequiredBy {
      childNames = append(childNames, childStage.Name)
    }

    dag.AddVertex(stage.Name, childNames, func (ctx context.Context) error {
      fmt.Printf("%v: bringing down\n", stage.Name)
			time.Sleep(1 * time.Second)
      fmt.Printf("%v: done\n", stage.Name)
      return nil
		})
	}

  return dag.Run(context.TODO())
}