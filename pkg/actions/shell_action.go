package actions

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// ShellCommandAction is an Action that will run a command using a Bourne shell (i.e. `sh`)
// TODO(nickt): support varargs, i.e. sh(foo, bar, baz), will run "sh -c foo bar baz"
// TODO(nickt): Handle sudo
type ShellCommandAction struct {

	// command is the shell command to run
	command string

	// outputWriter is a writer to use for capturing stdout and stderr
	outputWriter io.Writer

	// debug determines whether the command will output debug information to
	// the outputWriter
	debug bool
}

// NewShellCommandAction constructs and returns a new ShellCommandAction
// from the given command
func NewShellCommandAction(command string) *ShellCommandAction {
	return &ShellCommandAction{
		command:      command,
		outputWriter: os.Stderr,
	}
}

// Run executes the command as a Shell sub-process. If the command cannot be run,
// the error is returned.
func (s *ShellCommandAction) Run() error {
	cmd := exec.Command("sh", "-c", s.command)

	if s.debug {
		cmd.Stdout = s.outputWriter
		cmd.Stderr = s.outputWriter
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("shell_action: %s", err)
	}

	return nil
}

func (s *ShellCommandAction) Debug() {
	s.debug = true
}

// String prints a string representation of the current command
func (s ShellCommandAction) String() string {
	return "[sh]: " + s.command
}
