package actions

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nicktrav/matryoshka/pkg/lang"
)

func TestShellCommandAction_Run_Success(t *testing.T) {
	cmd := newCommand("echo foo")

	err := cmd.Run()
	if err != nil {
		t.Errorf("command failed: %s", err)
	}
}

func TestShellCommandAction_Run_Debug(t *testing.T) {
	want := "foo"
	cmd := newCommand("echo " + want)

	// replace stderr with our own buffer, to capture the output of the command
	buf := new(bytes.Buffer)
	cmd.outputWriter = buf

	// enable debugging
	cmd.Debug()

	err := cmd.Run()
	if err != nil {
		t.Errorf("command failed: %s", err)
	}

	stdout := buf.String()
	if !strings.Contains(stdout, want) {
		t.Errorf("wanted STDOUT to contain '%s'; got '%s'", want, stdout)
	}
}

func TestShellCommandAction_Run_Fail(t *testing.T) {
	cmd := newCommand("false")

	err := cmd.Run()
	if err == nil {
		t.Fatal("wanted command to fail with non-zero exit code")
	}

	want := "shell_action"
	if !strings.HasPrefix(err.Error(), want) {
		t.Errorf("wanted error prefix %s; got %s", want, err.Error())
	}
}

func TestShellCommandAction_String(t *testing.T) {
	command := "foo bar"
	cmd := newCommand(command)

	wanted := "[sh]: " + command
	if wanted != cmd.String() {
		t.Errorf("wanted %s, got %s", wanted, cmd)
	}
}

func TestNewShellCommandAction(t *testing.T) {
	inputs := []*lang.ShellCmd{
		{Command: "foo", Shell: "bash", Login: false},
		{Command: "foo", Shell: "sh", Login: false},
		{Command: "foo", Shell: "bash", Login: true},
	}

	type testResult struct {
		command string
		shell   string
		login   bool
	}
	want := []testResult{
		{"foo", "bash", false},
		{"foo", "sh", false},
		{"foo", "bash", true},
	}

	for i, input := range inputs {
		command := want[i].command
		if command != input.Command {
			t.Errorf("wanted command %s; got %s", command, input.Command)
		}

		shell := want[i].shell
		if shell != input.Shell {
			t.Errorf("wanted shell %s; got %s", shell, input.Shell)
		}

		login := want[i].login
		if login != input.Login {
			t.Errorf("wanted login %v; got %v", login, input.Login)
		}
	}
}

// newCommand returns a pointer to a new ShellCommandAction running the given
// command in a "sh" shell.
func newCommand(command string) *ShellCommandAction {
	return NewShellCommandAction(&lang.ShellCmd{Command: command, Shell: "sh"})
}
