package actions

import (
	"bytes"
	"strings"
	"testing"
)

func TestShellCommandAction_Run_Success(t *testing.T) {
	cmd := NewShellCommandAction("echo foo")

	err := cmd.Run()
	if err != nil {
		t.Errorf("command failed: %s", err)
	}
}

func TestShellCommandAction_Run_Writer(t *testing.T) {
	want := "foo"
	cmd := NewShellCommandAction("echo " + want)
	buf := new(bytes.Buffer)
	cmd.outputWriter = buf

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
	cmd := NewShellCommandAction("false")

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
	cmd := NewShellCommandAction(command)

	wanted := "[sh]: " + command
	if wanted != cmd.String() {
		t.Errorf("wanted %s, got %s", wanted, cmd)
	}
}
