package actions

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/nicktrav/matryoshka/pkg/lang"
)

// ShellCommandAction is an Action that will run a command using a Bourne shell (i.e. `sh`)
// TODO(nickt): support varargs, i.e. sh(foo, bar, baz), will run "sh -c foo bar baz"
// TODO(nickt): Handle sudo
type ShellCommandAction struct {

	// command is the shell command to run
	command string

	// shell is the type of shell to run
	shell string

	// login determines whether the shell should be a login shell
	login bool

	// outputWriter is a writer to use for capturing stdout and stderr
	outputWriter io.Writer

	// debug determines whether the command will output debug information to
	// the outputWriter
	debug bool
}

// NewShellCommandAction constructs and returns a new ShellCommandAction
// from the given command
func NewShellCommandAction(cmd *lang.ShellCmd) *ShellCommandAction {
	return &ShellCommandAction{
		command:      cmd.Command,
		shell:        cmd.Shell,
		login:        cmd.Login,
		outputWriter: os.Stderr,
	}
}

// Run executes the command as a Shell sub-process. If the command cannot be run,
// the error is returned.
func (s *ShellCommandAction) Run() error {
	var args []string
	if s.login {
		args = append(args, "-i")
	}
	args = append(args, "-c", s.command)

	cmd := exec.Command(s.shell, args...)

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
