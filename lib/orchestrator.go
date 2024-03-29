package anyform

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	util "github.com/elevationtools/anyform/common/util"
	daglib "github.com/elevationtools/anyform/lib/dag/mubased"
)

// Orchestrator ////////////////////////////////////////////////////////////////

type OrchestratorSpec struct {
  ImplDir string `json:"impl_dir"`
  Stages map[string]StageSpec `json:"stages"`
  InnerCfg map[string]any `json:"config"`
}

type Orchestrator struct {
  locator *Locator
  Spec OrchestratorSpec
  Stages map[string]*Stage
}

func NewOrchestrator(locator *Locator) (*Orchestrator, error) {
  orc := &Orchestrator{
    locator: locator,
    Stages: map[string]*Stage{},
  }
  err := orc.locator.JsonnetRunner.Run(
      orc.locator.Config.OrchestratorSpecFile, &orc.Spec)
  if err != nil { return nil, err }

  for stageName, iStageSpec := range orc.Spec.Stages {
    stageSpec := iStageSpec
    orc.Stages[stageName] = NewStage(stageName, orc.locator, &orc.Spec, &stageSpec)
  }

  for _, stage := range orc.Stages {
    for _, depName := range stage.spec.DependsOn {
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

func (orc* Orchestrator) Up(ctx context.Context) error {
	err := orc.WriteConfigJsonFile()
	if err != nil { return err }

  dag := daglib.NewDag()

  for _, iStage := range orc.Stages {
    stage := iStage

    parentNames := []string{}
    for _, parentStage := range stage.DependsOn {
      parentNames = append(parentNames, parentStage.Name)
    }

    dag.AddVertex(stage.Name, parentNames, func (ctx context.Context) error {
      return stage.Up(ctx)
    })
  }

  return dag.Run(ctx)
}

func (orc* Orchestrator) WriteConfigJsonFile() error {
	outFilePath := orc.locator.Config.Orchestrator.ConfigJsonFile

	jsonString, err := util.ToJSONString(orc.Spec.InnerCfg)
	if err != nil { return fmt.Errorf("converting InnerCfg to JSON: %w", err) }

	dir := filepath.Dir(outFilePath)
	err = os.MkdirAll(dir, 0750)
	if err != nil { return fmt.Errorf("mkdir -p %v: %w", dir, err) }

	
	err = os.WriteFile(outFilePath, []byte(jsonString), 0660)
	if err != nil { return fmt.Errorf("writing %v: %w", outFilePath, err) }

	return nil
}

func (orc* Orchestrator) Down(ctx context.Context) error {
	err := orc.WriteConfigJsonFile()
	if err != nil { return err }

  dag := daglib.NewDag()

  for _, iStage := range orc.Stages {
    stage := iStage

    childNames := []string{}
    for _, childStage := range stage.RequiredBy {
      childNames = append(childNames, childStage.Name)
    }

    dag.AddVertex(stage.Name, childNames, func (ctx context.Context) error {
      return stage.Down(ctx)
    })
  }

  return dag.Run(ctx)
}