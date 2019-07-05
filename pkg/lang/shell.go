package lang

import (
	"errors"
	"fmt"

	"go.starlark.net/starlark"
)

// ShellCmd represents the `shell()` builtin function and represents a command
// to run as a shell subprocess.
//
// The structure of `shell` is as follows:
//
//   shell(
//     'echo 42', // the shell command to run
//   )
//
// TODO(nickt): Support varargs
type ShellCmd struct {

	// Command is the shell command to execute.
	Command string
}

// String returns the string representation of the ShellCmd.
func (s ShellCmd) String() string {
	return fmt.Sprintf("<dep.ShellCmd %q>", s.Command)
}

// Type returns a short description about ShellCmd's type.
func (s ShellCmd) Type() string { return "dep.ShellCmd" }

// Freeze does nothing for a ShellCmd.
func (s ShellCmd) Freeze() {}

// Truth always returns true for a ShellCmd.
func (s ShellCmd) Truth() starlark.Bool { return starlark.True }

// Hash is currently not implemented by ShellCmd.
func (s ShellCmd) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: %s", s.Type())
}

// FnShell implements the signature for a builtin function and implements
// the functionality of the `shell` function.
//
// FnShell transforms the arguments into a ShellCmd object after performing
// validation on the arguments.
// TODO(nickt): Validate arguments
func FnShell(t *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	value := args.Index(0)
	strValue, ok := value.(starlark.String)
	if !ok {
		return nil, errors.New(fmt.Sprintf("expected %v to be ok type String", value))
	}

	shell := ShellCmd{
		Command: string(strValue),
	}

	return shell, nil
}
