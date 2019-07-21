package lang

import (
	"strings"
	"testing"

	"go.starlark.net/starlark"
)

func TestString(t *testing.T) {
	cmd := ShellCmd{Command: "foo"}

	wanted := "<dep.ShellCmd \"foo\">"
	got := cmd.String()

	if wanted != got {
		t.Errorf("wanted %s, got %s", wanted, got)
	}
}

func TestType(t *testing.T) {
	cmd := ShellCmd{Command: "foo"}

	wanted := "dep.ShellCmd"
	got := cmd.Type()

	if wanted != got {
		t.Errorf("wanted %s, got %s", wanted, got)
	}
}

func TestTruth(t *testing.T) {
	cmd := ShellCmd{Command: "foo"}

	if !cmd.Truth() {
		t.Error("wanted cmd.Truth to return true; got false")
	}
}

func TestHash(t *testing.T) {
	cmd := ShellCmd{Command: "foo"}

	_, err := cmd.Hash()

	if err == nil {
		t.Fatal("wanted Hash to return an error")
	}

	if !strings.HasPrefix(err.Error(), "unhashable type") {
		t.Errorf("wanted error to contain 'unhashable type")
	}
}

func TestFnShell_defaults(t *testing.T) {
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

	got = cmd.Shell
	if got != defaultShell {
		t.Errorf("wanted default shell %s; got %s", defaultShell, got)
	}

	if cmd.Login {
		t.Errorf("wanted login shell to be false; was true")
	}
}

func TestFnShell_nonDefaultShell(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}

	args := []starlark.Value{starlark.String("foo")}
	nonDefaultShell := "zsh"
	kwargs := []starlark.Tuple{{starlark.String("shell"), starlark.String("zsh")}}

	value, err := FnShell(thread, builtin, args, kwargs)
	if err != nil {
		t.Errorf("error running FnShell with args %+v, kwargs %+v", args, kwargs)
	}

	cmd, ok := value.(ShellCmd)
	if !ok {
		t.Errorf("returned value %s was not a ShellCmd", value)
	}

	got := cmd.Shell
	if got != nonDefaultShell {
		t.Errorf("wanted default shell %s; got %s", nonDefaultShell, got)
	}
}

func TestFnShell_nonDefaultLogin(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}

	args := []starlark.Value{starlark.String("foo")}
	kwargs := []starlark.Tuple{{starlark.String("login"), starlark.Bool(true)}}

	value, err := FnShell(thread, builtin, args, kwargs)
	if err != nil {
		t.Errorf("error running FnShell with args %+v, kwargs %+v", args, kwargs)
	}

	cmd, ok := value.(ShellCmd)
	if !ok {
		t.Errorf("returned value %s was not a ShellCmd", value)
	}

	if !cmd.Login {
		t.Error("wanted login to be true; got false")
	}
}
