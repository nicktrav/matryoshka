package lang

import (
	"strings"
	"testing"

	"go.starlark.net/starlark"
)

func TestString(t *testing.T) {
	cmd := ShellCmd{"foo"}

	expected := "<dep.ShellCmd \"foo\">"
	actual := cmd.String()

	if expected != actual {
		t.Errorf("wanted %s, got %s", expected, actual)
	}
}

func TestType(t *testing.T) {
	cmd := ShellCmd{"foo"}

	expected := "dep.ShellCmd"
	actual := cmd.Type()

	if expected != actual {
		t.Errorf("wanted %s, got %s", expected, actual)
	}
}

func TestTruth(t *testing.T) {
	cmd := ShellCmd{"foo"}

	if !cmd.Truth() {
		t.Error("Truth did not return true")
	}
}

func TestHash(t *testing.T) {
	cmd := ShellCmd{"foo"}

	_, err := cmd.Hash()

	if err == nil {
		t.Fatal("expected Hash to return an error")
	}

	if !strings.Contains(err.Error(), "unhashable type") {
		t.Errorf("expected error to contain 'unhashable type")
	}
}

func TestShellCmd(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}

	expectedCommand := "foo"
	input := []starlark.Value{starlark.String(expectedCommand)}
	value, err := FnShell(thread, builtin, input, []starlark.Tuple{})

	if err != nil {
		t.Errorf("error running FnShell on tuple %v", input)
	}

	cmd, ok := value.(ShellCmd)
	if !ok {
		t.Errorf("returned value %s was not a ShellCmd", value)
	}

	if expectedCommand != cmd.Command {
		t.Errorf("wanted %s, got %s", expectedCommand, cmd.Command)
	}
}
