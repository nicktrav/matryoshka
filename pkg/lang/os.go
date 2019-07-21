package lang

import (
	"runtime"

	"go.starlark.net/starlark"
)

// FnOs returns the operating system of the caller.
//
// The structure of `os` is as follows:
//
//   os()
//
// The functions returns the operating system, as inferred by the Go runtime.
func FnOs(t *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.String(runtime.GOOS), nil
}
