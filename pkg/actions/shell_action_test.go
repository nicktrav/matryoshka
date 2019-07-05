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
	expected := "foo"
	cmd := NewShellCommandAction("echo " + expected)
	buf := new(bytes.Buffer)
	cmd.outputWriter = buf

	err := cmd.Run()
	if err != nil {
		t.Errorf("command failed: %s", err)
	}

	stdout := strings.TrimSpace(buf.String())
	if stdout != expected {
		t.Errorf("expected STDOUT to be '%s'; got '%s'", expected, stdout)
	}
}

func TestShellCommandAction_Run_Fail(t *testing.T) {
	cmd := NewShellCommandAction("false")

	err := cmd.Run()
	if err == nil {
		t.Fatal("expected command to fail with non-zero exit code")
	}

	expected := "shell_action: exit status 1"
	if expected != err.Error() {
		t.Errorf("expected error message %s; got %s", expected, err.Error())
	}
}

func TestShellCommandAction_String(t *testing.T) {
	cmd := NewShellCommandAction("foo bar")

	expected := "[sh]: foo bar"
	actual := cmd.String()

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}
