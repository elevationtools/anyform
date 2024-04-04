package anyform

import (
  "context"
  "fmt"
	"io/fs"
  "log/slog"
  "os"
  "os/exec"
  "path/filepath"
	 "time"

  "github.com/elevationtools/anyform/module/common/util"
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
	stampDir string
	logDir string
  // Used for both stamping and running of the stage ctl command.
  envVars []string
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
	s.stampDir = filepath.Join(s.globe.Config.Orchestrator.GenfilesDir, s.Name, "stamp")
	s.logDir = filepath.Join(
			s.globe.Config.Orchestrator.GenfilesDir, s.Name, "logs")
  s.envVars = []string{
		// DOCS(STAGE_ENVIRONMENT_VARIABLES)
    "ANYFORM_STAGE_NAME=" + s.Name,
    "ANYFORM_STAGE_STAMP_DIR=" + AbsJoin(s.stampDir),
    "ANYFORM_CONFIG_JSON_FILE=" + AbsJoin(s.globe.Config.Orchestrator.ConfigJsonFile),
    "ANYFORM_GENFILES=" + AbsJoin(s.globe.Config.Orchestrator.GenfilesDir),
    "ANYFORM_IMPL_DIR=" + AbsJoin(s.orchestratorSpec.ImplDir),
    "ANYFORM_OUTPUT_DIR=" + AbsJoin(s.globe.Config.Orchestrator.OutputDir),
  }
  return s
}

// Wrap UpImpl so that all errors can be captured and displayed.  The DAG
// swallows errors so they need to be displayed here.
func (s *Stage) Up(ctx context.Context) error {
	slog.Info("stage up starting", "stage", s.Name)
	err := s.UpImpl(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[stage=%v] failed: %v", s.Name, err)
		return err
	}
  slog.Info("stage up done", "stage", s.Name)
	return nil
}

func (s *Stage) UpImpl(ctx context.Context) error {
	if s.alreadyUpToDate("up") { return nil }

  err := s.Stamp(ctx)
  if err != nil { return err }

  err = s.RunStampedCtl(ctx, "up")
  if err != nil { return err }

  err = util.ToJSONFile(StageStateFile{LastCommand: "up"}, s.stateFilePath)
  if err != nil { return Errorf("writing %v: %w", s.stateFilePath, err) }
  
  return nil
}

// Returns true if the command doesn't need to be run because:
// - The state file's "last_command" is the given command.
// - The CONFIG_JSON_FILE is older than the state file.
// - TODO: The implementation of this stage or any parent stages has changed.
// Always returns false on errors, which allows ignoring errors.
func (s *Stage) alreadyUpToDate(command string) bool {
  autd, reason, err := s.alreadyUpToDateImpl(command)
	if err != nil {
		// Ignore errors here because all stages must be idempotent and maybe the
		// error is recoverable
		fmt.Fprintf(os.Stderr,
			 "[stage=%v] warning: unable to determine if operation is already done," +
			 " assuming it's not: %v\n",
				s.Name, err)
			return false
	}
	
	if autd {
    fmt.Printf("[stage=%v] skipping, already done\n", s.Name)
    return true
  }

	slog.Debug("[stage=%v] needs updating, reason: %v\n", s.Name, reason)
	return false
}

// See alreadyUpToDate
func (s *Stage) alreadyUpToDateImpl(command string) (bool, string, error) {
  var stateData StageStateFile
  err := util.FromJSONFile(s.stateFilePath, &stateData)
  if err != nil {
		if os.IsNotExist(err) { return false, fmt.Sprintf("file doesn't exist: %v", s.stateFilePath), nil }
		return false, "", err
	}
  if stateData.LastCommand != command { return false, "different command", nil }
  stateFileInfo, err := os.Stat(s.stateFilePath)
  if err != nil { return false, "", err }
  configFileInfo, err := os.Stat(s.globe.Config.Orchestrator.ConfigJsonFile)
  if err != nil { return false, "", err }
  if configFileInfo.ModTime().After(stateFileInfo.ModTime()) {
    return false, "config file newer", nil
  }
	maxImplFileModTime, err := s.maxImplFileModTime()
	if err != nil { return false, "", err }
  if maxImplFileModTime.After(stateFileInfo.ModTime()) {
    return false, "impl newer", nil
	}
  return true, "", nil
}

func (s *Stage) maxImplFileModTime() (time.Time, error) {
	maxTime := time.UnixMilli(0)
	err := filepath.Walk(s.stageImplDir,
			func (path string, info fs.FileInfo, err error) error {
		if err != nil { return err }
		modTime := info.ModTime()
		if modTime.After(maxTime) { maxTime = modTime }
		return nil
	})
	if err != nil { return maxTime, err }

	for _, parent := range s.DependsOn {
		parentMax, err := parent.maxImplFileModTime()
		if err != nil { return maxTime, err }
		if parentMax.After(maxTime) { maxTime = parentMax }
	}
	return maxTime, nil
}

func (s *Stage) Stamp(ctx context.Context) error {
  slog.Debug("stage stamping", "stage", s.Name)
  err := MkdirAll(s.stampDir)
  if err != nil {
		return Errorf("Making stage stamp dir '%v': %w", s.stampDir, err)
  }
  return s.globe.StageStamper.Stamp(
			ctx, s.Name, s.stageImplDir, s.stampDir, s.logDir, s.envVars)
}

func AbsJoin(elem ...string) string {
  res, err := filepath.Abs(filepath.Join(elem...))
  if err != nil { panic(err) }
  return res
}

func (s *Stage) RunStampedCtl(ctx context.Context, ctlArg string) error {
  logStr := "stage './ctl " + ctlArg + "'"
  slog.Debug(logStr, "stage", s.Name)

  cmd := exec.CommandContext(ctx, AbsJoin(s.stampDir, "/ctl"), ctlArg)
  cmd.Dir = s.stampDir
  cmd.Env = append(cmd.Environ(), s.envVars...)

  err := s.globe.SubprocessRunner.RunCmd("stage=" + s.Name, cmd, s.logDir)
  if err != nil { return Errorf("stage %v: %w", s.Name, err) }

  slog.Debug(logStr + " completed", "stage", s.Name)
  return nil
}

func (s *Stage) Down(ctx context.Context) error {
  slog.Info("stage down starting", "stage", s.Name)

	if s.alreadyUpToDate("down") { return nil }

  err := s.RunStampedCtl(ctx, "down")
  if err != nil { return err }

  err = util.ToJSONFile(StageStateFile{LastCommand: "down"}, s.stateFilePath)
  if err != nil { return Errorf("writing %v: %w", s.stateFilePath, err) }
 
  slog.Info("stage down done", "stage", s.Name)
  return nil
}
