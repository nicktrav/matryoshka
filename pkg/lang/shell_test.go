package lang

import (
	"strings"
	"testing"

	"go.starlark.net/starlark"
)

func TestString(t *testing.T) {
	cmd := ShellCmd{"foo"}

	wanted := "<dep.ShellCmd \"foo\">"
	got := cmd.String()

	if wanted != got {
		t.Errorf("wanted %s, got %s", wanted, got)
	}
}

func TestType(t *testing.T) {
	cmd := ShellCmd{"foo"}

	wanted := "dep.ShellCmd"
	got := cmd.Type()

	if wanted != got {
		t.Errorf("wanted %s, got %s", wanted, got)
	}
}

func TestTruth(t *testing.T) {
	cmd := ShellCmd{"foo"}

	if !cmd.Truth() {
		t.Error("wanted cmd.Truth to return true; got false")
	}
}

func TestHash(t *testing.T) {
	cmd := ShellCmd{"foo"}

	_, err := cmd.Hash()

	if err == nil {
		t.Fatal("wanted Hash to return an error")
	}

	if !strings.HasPrefix(err.Error(), "unhashable type") {
		t.Errorf("wanted error to contain 'unhashable type")
	}
}

func TestShellCmd(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}

	want := "foo"
	input := []starlark.Value{starlark.String(want)}
	value, err := FnShell(thread, builtin, input, []starlark.Tuple{})

	if err != nil {
		t.Errorf("error running FnShell on tuple %v", input)
	}

	cmd, ok := value.(ShellCmd)
	if !ok {
		t.Errorf("returned value %s was not a ShellCmd", value)
	}

	got := cmd.Command
	if want != got {
		t.Errorf("wanted %s, got %s", want, got)
	}
}
