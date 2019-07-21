package lang

import (
	"fmt"
	"runtime"
	"testing"

	"go.starlark.net/starlark"
)

func TestFnOs(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}

	got, err := FnOs(thread, builtin, []starlark.Value{}, []starlark.Tuple{})
	if err != nil {
		t.Fatalf("got error running FnOs: %+v", err)
	}

	// a starlark primitive string value includes the quotes
	want := fmt.Sprintf("\"%s\"", runtime.GOOS)
	if got.String() != want {
		t.Errorf("wanted %s; got %s", want, got)
	}
}
