package anyform

import (
	"bufio"
	"fmt"
	"io"
	"os"
  "os/exec"
	"path/filepath"
	"sync"
	"time"
)

type DefaultSubprocessRunner struct {
  globe *Globe
	mu sync.Mutex
}

func NewDefaultSubprocessRunner(globe *Globe) *DefaultSubprocessRunner {
  dc := &DefaultSubprocessRunner{}
  dc.globe = globe
  return dc
}

// TODO: consider also creating separate, unadultered stdout and stderr streams.
// TODO: consider appending instead of truncating log files.
func (dc *DefaultSubprocessRunner) RunCmd(
    label string, cmd *exec.Cmd, logDir string) error {
	err := os.MkdirAll(logDir, 0750)
	if err != nil { return Errorf("mkdir -p '%v': %w", logDir, err) }
	logFilePath := filepath.Join(logDir, "stdout_stderr")
	logFile, err := os.Create(logFilePath)
	if err != nil { return Errorf("creating log file '%v': %w", logFilePath, err) }

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil { return Errorf("creating stdout pipe: %w", err) }
	stderrPipe, err := cmd.StderrPipe()
	if err != nil { return Errorf("creating stderr pipe: %w", err) }

	err = cmd.Start()
	if err != nil { return Errorf("starting subprocess: %w", err) }

	// Relay subprocess's stdout and stderr to this process's stdout and the
	// logfile.
  // TODO(bug): STDOUT_STDERR_ORDER_ISSUE This strategy can cause the relative
  // order of stdout and stderr to get mixed up for writes that happen directly
  // back-to-back.  This has been seen in the integration test where a bash
  // script does two back-to-back echos to stdout and stderr. Another approach
  // could be to use a single goroutine and select within it, but this may
  // require significantly more code since the bufio Scanner probably couldn't
  // be used.
	var wg sync.WaitGroup
	wg.Add(2)
	var stdoutErr, stderrErr error
	go dc.relayStream(label, "stdout", stdoutPipe, os.Stdout, logFile, &stdoutErr, &wg)
	go dc.relayStream(label, "stderr", stderrPipe, os.Stderr, logFile, &stderrErr, &wg)
	wg.Wait()

	if stdoutErr != nil { return Errorf("relaying stdout: %w", stdoutErr) }
	if stderrErr != nil { return Errorf("relaying stderr: %w", stderrErr) }
	err = cmd.Wait()
	if err != nil { return Errorf("running subprocess: %w", err) }

  return nil
}

func (dc *DefaultSubprocessRunner) relayStream(
		label string, streamName string, stream io.ReadCloser,
		thisProcessStdStream *os.File, logFile *os.File,
		err *error, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		timestamp := time.Now().Format(time.RFC3339)
		dc.mu.Lock()
		fmt.Fprintf(thisProcessStdStream, "[%v %v] %v\n", timestamp, label, scanner.Text())
		fmt.Fprintf(logFile, "[%v %v] %v\n", timestamp, streamName, scanner.Text())
		dc.mu.Unlock()
	}
	*err = scanner.Err()
}
