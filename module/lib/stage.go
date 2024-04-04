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
	ctlPath string
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
	s.ctlPath = filepath.Join(s.stampDir, CtlFileName)
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
	err := s.UpImpl(ctx)
	if err != nil {
		s.stderr("failed: %v", err)
		return err
	}
	return nil
}

func (s *Stage) UpImpl(ctx context.Context) error {
	s.stdout("starting")
	if s.alreadyUpToDate("up") { return nil }

	s.stdout("stamping")
  err := s.Stamp(ctx)
  if err != nil { return err }

	s.stdout("running 'ctl up'")
  err = s.RunStampedCtl(ctx, "up")
  if err != nil { return err }

  err = util.ToJSONFile(StageStateFile{LastCommand: "up"}, s.stateFilePath)
  if err != nil { return Errorf("writing %v: %w", s.stateFilePath, err) }
  
	s.stdout("done")
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
		s.stderr(
				"warning: unable to determine if '%v' is already done," +
			 	" assuming it's not: %v", command, err)
			return false
	}
	
	if autd {
    s.stdout("skipping, already done")
    return true
  }

	slog.Debug(fmt.Sprintf(
			"[stage=%v] needs updating, reason: %v", s.Name, reason))
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

  // TODO(fragile): figure out how to make a relative path in a platform
  // independent way.  filepath.Join(".", "foo") doesn't work.
  cmd := exec.CommandContext(ctx, "./" + CtlFileName, ctlArg)
  cmd.Dir = s.stampDir
  cmd.Env = append(cmd.Environ(), s.envVars...)

  err := s.globe.SubprocessRunner.RunCmd(s.Name, cmd, s.logDir)
  if err != nil { return Errorf("running ctl: %w", err) }

  slog.Debug(logStr + " completed", "stage", s.Name)
  return nil
}

func (s *Stage) Down(ctx context.Context) error {
	err := s.DownImpl(ctx)
	if err != nil {
		s.stderr("failed: %v", err)
		return err
	}
	return nil
}

func (s *Stage) DownImpl(ctx context.Context) error {
	s.stdout("starting")

	if s.alreadyUpToDate("down") { return nil }

	_, err := os.Stat(s.ctlPath)
	if err != nil {
		if !os.IsNotExist(err) { return err }
		// Optimistically try to stamp, maybe we'll get lucky and it will work!
		// It will work if:
		// 1) all dependencies have run successfully.
		// 2) stamping happens to not depend on dependencies.
		// Because of (2) we don't bother checking for (1).
		stampErr := s.Stamp(ctx)
		if stampErr != nil {
			return Errorf("optimistic stamp failed: %v", stampErr)
		}
		s.stdout("optimistic stamp succeeded")
	}

	s.stdout("running 'ctl down'")
  err = s.RunStampedCtl(ctx, "down")
  if err != nil { return err }

  err = util.ToJSONFile(StageStateFile{LastCommand: "down"}, s.stateFilePath)
  if err != nil { return Errorf("writing '%v': %w", s.stateFilePath, err) }
 
	s.stdout("done")
  return nil
}

// Print to stdout.
// Already adds prefix and newline.
func (s *Stage) stdout(format string, args... any) {
	s.fprintf(os.Stdout, format, args...)
}

func (s *Stage) stderr(format string, args... any) {
	s.fprintf(os.Stderr, format, args...)
}

func (s *Stage) fprintf(file *os.File, format string, args... any) {
	args = append([]any{s.Name}, args...)
	fmt.Fprintf(file, "[%v] " + format + "\n", args...)
}
