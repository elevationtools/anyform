
package anyform

import (
  "context"
  "fmt"
  "log/slog"
  "os"
  "os/exec"
  "path/filepath"
)

// Stage ///////////////////////////////////////////////////////////////////////

type StageSpec struct {
  DependsOn []string `json:"depends_on"`
}

type Stage struct {
  Name string
  DependsOn []*Stage
  RequiredBy []*Stage

  locator *Locator
  orchestratorSpec *OrchestratorSpec
  spec *StageSpec
}

func NewStage(name string, locator *Locator,
              orchestratorSpec *OrchestratorSpec, spec *StageSpec) *Stage {
  return &Stage{
    Name: name,
    DependsOn: []*Stage{},
    RequiredBy: []*Stage{},

    locator: locator,
    orchestratorSpec: orchestratorSpec,
    spec: spec,
  }
}

func (s *Stage) Up(ctx context.Context) error {
  err := s.Stamp()
  if err != nil {
    slog.Warn("stage up stamping failed", "stage", s.Name, "error", err)
    return err
  }

  err = s.RunCmd(ctx, "up")
  if err != nil {
    slog.Warn("stage up running failed", "stage", s.Name, "error", err)
    return err
  }
  
  slog.Info("stage up done", "stage", s.Name)
  return nil
}

func (s *Stage) stampDir() string {
  return filepath.Join(s.locator.Config.Orchestrator.GenfilesDir, s.Name)
}

func (s *Stage) Stamp() error {
  slog.Info("stage stamping", "stage", s.Name)
  stampDir := s.stampDir()
  slog.Debug(fmt.Sprintf("mkdir -p %v", stampDir))
  os.MkdirAll(stampDir, 0750)
  return s.locator.GomplateRunner.Run(
    filepath.Join(s.orchestratorSpec.ImplDir, s.Name), stampDir)
}

func AbsJoin(elem ...string) string {
	res, err := filepath.Abs(filepath.Join(elem...))
	if err != nil { panic(err) }
	return res
}

func (s *Stage) RunCmd(ctx context.Context, cmdStr string) error {
  logStr := "stage './ctl " + cmdStr + "'"
  slog.Info(logStr, "stage", s.Name)
  cmd := exec.CommandContext(ctx, AbsJoin(s.stampDir(), "/ctl"), cmdStr)
  cmd.Dir = s.stampDir()
	cmd.Env = append(cmd.Environ(),
		"ANYFORM_STAGE_NAME=" + s.Name,
		"ANYFORM_CONFIG_JSON_FILE=" + AbsJoin(s.locator.Config.Orchestrator.ConfigJsonFile),
		"ANYFORM_GENFILES=" + AbsJoin(s.locator.Config.Orchestrator.GenfilesDir),
		"ANYFORM_IMPL_DIR=" + AbsJoin(s.orchestratorSpec.ImplDir),
		"ANYFORM_OUTPUT_DIR=" + AbsJoin(s.locator.Config.Orchestrator.OutputDir),
		"ANYFORM_INTERACTIVE=" + func() string {
			if s.locator.Config.Orchestrator.Interactive { return "true" }
			return "false"
		}(),
  )
  res, err := cmd.CombinedOutput()
  if err != nil {
    return fmt.Errorf(logStr + ": %w: output: %v", err, string(res))
  }
  slog.Info(logStr + " completed", "stage", s.Name, "result", string(res))
  return nil
}

func (s *Stage) Down(ctx context.Context) error {
  slog.Info("stage down starting", "stage", s.Name)

  err := s.RunCmd(ctx, "down")
  if err != nil { return err }

  slog.Info("stage down done", "stage", s.Name)
  return nil
}
