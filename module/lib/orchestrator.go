package anyform

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
  "time"

	util "github.com/elevationtools/anyform/module/common/util"
	daglib "github.com/elevationtools/anyform/module/lib/dag/mubased"
)

// Orchestrator ////////////////////////////////////////////////////////////////

type OrchestratorSpec struct {
  ImplDir string `json:"impl_dir"`
  Stages map[string]StageSpec `json:"stages"`
  InnerCfg map[string]any `json:"config"`
}

type Orchestrator struct {
  globe *Globe
  Spec OrchestratorSpec
  Stages map[string]*Stage
}

func NewOrchestrator(globe *Globe) (*Orchestrator, error) {
  orc := &Orchestrator{
    globe: globe,
    Stages: map[string]*Stage{},
  }

  err := orc.globe.ConfigLoader.Load(
      orc.globe.Config.OrchestratorSpecFile, &orc.Spec)
  if err != nil { return nil, err }

  err = orc.MaybeUpdateConfigJsonFile()
  if err != nil { return nil, err }

  for stageName, iStageSpec := range orc.Spec.Stages {
    stageSpec := iStageSpec
    orc.Stages[stageName] = NewStage(stageName, orc.globe, &orc.Spec, &stageSpec)
  }

	// Build stage DAG.
  for _, stage := range orc.Stages {
    for _, depName := range stage.spec.DependsOn {
      dep, found := orc.Stages[depName]
      if !found {
        return nil, Errorf(
          "invalid DAG: stage %v depends on undefined stage %v", stage.Name, depName)
      }
      stage.DependsOn = append(stage.DependsOn, dep)
      dep.RequiredBy = append(dep.RequiredBy, stage)
    }
  }

  return orc, nil
}

func (orc* Orchestrator) MaybeUpdateConfigJsonFile() error {
  depFiles, err := orc.globe.ConfigLoader.GetTransitiveDeps(
      orc.globe.Config.OrchestratorSpecFile)
  if err != nil { return err }

  configsMaxModTime := time.UnixMilli(0)
  for _, f := range depFiles {
    info, err := os.Stat(f)
    if err != nil {
			return Errorf("Stat()ing config file '%v': %w", f, err)
		}
		modTime := info.ModTime()
		if modTime.After(configsMaxModTime) { configsMaxModTime = modTime }
  }

	jsonFilePath := orc.globe.Config.Orchestrator.ConfigJsonFile
	jsonFileInfo, err := os.Stat(jsonFilePath)
	if err != nil {
		// Ignore because the file might just be missing, and if it's not perhaps
		// assuming it's missing and overwriting it will solve the problem.
	} else if !configsMaxModTime.After(jsonFileInfo.ModTime()) {
		slog.Debug("config json file already up to date")
		return nil
	}

	jsonString, err := util.ToJSONString(orc.Spec.InnerCfg)
	if err != nil { return Errorf("converting InnerCfg to JSON: %w", err) }

	dir := filepath.Dir(jsonFilePath)
	err = os.MkdirAll(dir, 0750)
	if err != nil { return Errorf("mkdir -p %v: %w", dir, err) }
	
	err = os.WriteFile(jsonFilePath, []byte(jsonString), 0660)
	if err != nil { return Errorf("writing %v: %w", jsonFilePath, err) }

	return nil
}

func (orc* Orchestrator) Up(ctx context.Context) error {
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

  return dag.Run(ctx, !orc.globe.Config.Interactive)
}

func (orc* Orchestrator) Down(ctx context.Context) error {
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

  return dag.Run(ctx, !orc.globe.Config.Interactive)
}

func (orc* Orchestrator) Clean(ctx context.Context) error {
  return os.RemoveAll(orc.globe.Config.Orchestrator.GenfilesDir)
}
