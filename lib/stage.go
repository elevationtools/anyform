package anyform

import (
  "context"
  "fmt"
  "log/slog"
  "os"
  "os/exec"
  "path/filepath"

  commonutil "github.com/elevationtools/anyform/common/util"
)

// Stage ///////////////////////////////////////////////////////////////////////

type StageSpec struct {
  DependsOn []string `json:"depends_on"`
}

type StageStateFile struct {
  LastCommand string `json:"last_command"`
}

type Stage struct {
  Name string
  DependsOn []*Stage
  RequiredBy []*Stage

  globe *Globe
  orchestratorSpec *OrchestratorSpec
  spec *StageSpec

  stateFilePath string
	stageImplDir string
}

func NewStage(name string, globe *Globe,
              orchestratorSpec *OrchestratorSpec, spec *StageSpec) *Stage {
  s := &Stage{
    Name: name,
    DependsOn: []*Stage{},
    RequiredBy: []*Stage{},

    globe: globe,
    orchestratorSpec: orchestratorSpec,
    spec: spec,
  }

  s.stateFilePath = 
      filepath.Join(s.globe.Config.Orchestrator.GenfilesDir, s.Name, "state")
	s.stageImplDir = filepath.Join(s.orchestratorSpec.ImplDir, s.Name)

  return s
}

func (s *Stage) Up(ctx context.Context) error {
	slog.Info("stage up starting", "stage", s.Name)
  autd, err := s.alreadyUpToDate("up")
	if err != nil { return err }
	if autd {
    slog.Info("stage up skipping, already up to date", "stage", s.Name)
    return nil
  }

  err = s.Stamp()
  if err != nil {
    slog.Warn("stage up stamping failed", "stage", s.Name, "error", err)
    return err
  }

  err = s.RunCmd(ctx, "up")
  if err != nil {
    slog.Warn("stage up running failed", "stage", s.Name, "error", err)
    return err
  }

  err = commonutil.ToJSONFile(StageStateFile{LastCommand: "up"},
                              s.stateFilePath)
  if err != nil { return fmt.Errorf("writing %v: %w", s.stateFilePath, err) }
  
  slog.Info("stage up done", "stage", s.Name)
  return nil
}

// Returns true if the command doesn't need to be run because:
// - The state file's "last_command" is the given command.
// - The CONFIG_JSON_FILE is older than the state file.
// - TODO: The implementation of this stage or any parent stages has changed.
// Always returns false on errors, which probably allows ignoring errors.
func (s *Stage) alreadyUpToDate(command string) (bool, error) {
  var stateFile StageStateFile
  err := commonutil.FromJSONFile(s.stateFilePath, &stateFile)
  if err != nil { return false, err }
  if stateFile.LastCommand != command { return false, nil }
  stateFileInfo, err := os.Stat(s.stateFilePath)
  if err != nil { return false, err }
  configFileInfo, err := os.Stat(s.globe.Config.Orchestrator.ConfigJsonFile)
  if err != nil { return false, err }
  return stateFileInfo.ModTime().After(configFileInfo.ModTime()), nil
}

func (s *Stage) stampDir() string {
  return filepath.Join(s.globe.Config.Orchestrator.GenfilesDir, s.Name, "stamp")
}

func (s *Stage) Stamp() error {
  slog.Debug("stage stamping", "stage", s.Name)
  stampDir := s.stampDir()
   err := os.MkdirAll(stampDir, 0750)
  if err != nil { return fmt.Errorf("mkdir -p '%v': %w", stampDir, err) }
  return s.globe.StageStamper.Stamp(
    filepath.Join(s.orchestratorSpec.ImplDir, s.Name), stampDir)
}

func AbsJoin(elem ...string) string {
  res, err := filepath.Abs(filepath.Join(elem...))
  if err != nil { panic(err) }
  return res
}

func (s *Stage) RunCmd(ctx context.Context, ctlArg string) error {
  logStr := "stage './ctl " + ctlArg + "'"
  slog.Debug(logStr, "stage", s.Name)

  cmd := exec.CommandContext(ctx, AbsJoin(s.stampDir(), "/ctl"), ctlArg)
  cmd.Dir = s.stampDir()
  cmd.Env = append(cmd.Environ(),
    "ANYFORM_STAGE_NAME=" + s.Name,
    "ANYFORM_CONFIG_JSON_FILE=" + AbsJoin(s.globe.Config.Orchestrator.ConfigJsonFile),
    "ANYFORM_GENFILES=" + AbsJoin(s.globe.Config.Orchestrator.GenfilesDir),
    "ANYFORM_IMPL_DIR=" + AbsJoin(s.orchestratorSpec.ImplDir),
    "ANYFORM_OUTPUT_DIR=" + AbsJoin(s.globe.Config.Orchestrator.OutputDir),
    "ANYFORM_INTERACTIVE=" + func() string {
      if s.globe.Config.Orchestrator.Interactive { return "true" }
      return "false"
    }(),
  )
  err := s.globe.SubprocessRunner.RunCmd(
      "stage=" + s.Name, cmd, filepath.Join(
      s.globe.Config.Orchestrator.GenfilesDir, s.Name, "logs"))
  if err != nil { return fmt.Errorf("stage %v: %w", s.Name, err) }

  slog.Debug(logStr + " completed", "stage", s.Name)
  return nil
}

func (s *Stage) Down(ctx context.Context) error {
  slog.Info("stage down starting", "stage", s.Name)

  err := s.RunCmd(ctx, "down")
  if err != nil { return err }

  slog.Info("stage down done", "stage", s.Name)
  return nil
}