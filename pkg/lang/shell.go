package lang

import (
	"fmt"

	"go.starlark.net/starlark"
)

const (
	shellArg     = starlark.String("shell")
	defaultShell = "bash"

	loginArg = starlark.String("login")
)

// ShellCmd represents the `shell()` builtin function and represents a command
// to run as a shell subprocess.
//
// The structure of `shell` is as follows:
//
//   shell(
//     'echo 42',   // the shell command to run
//     shell='bash' // the shell to run the command in (defaults to 'bash')
//     login=False  // the shell is a login shell (defaults to 'False')
//   )
//
// TODO(nickt): Support varargs
type ShellCmd struct {
	// Command is the shell command to execute.
	Command string

	// Shell is the shell to run the command in.
	Shell string

	// Login indicates whether this command should run in a login shell or not.
	Login bool
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
		return nil, fmt.Errorf("expected %v to be ok type String", value)
	}

	shell := defaultShell
	login := false
	for _, kwarg := range kwargs {
		key := kwarg.Index(0)
		value := kwarg.Index(1)
		switch key {

		case shellArg:
			s, ok := starlark.AsString(value)
			if !ok {
				return nil, fmt.Errorf("shell: argument to shell is not a string")
			}
			shell = s

		case loginArg:
			s, ok := value.(starlark.Bool)
			if !ok {
				return nil, fmt.Errorf("shell: argument to login is not a boolean")
			}
			login = s == starlark.True
		}
	}

	cmd := ShellCmd{
		Command: string(strValue),
		Shell:   shell,
		Login:   login,
	}

	return cmd, nil
}
