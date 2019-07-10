package lang

import (
	"fmt"

	"go.starlark.net/starlark"
)

const (
	argName        = starlark.String("name")
	argDescription = starlark.String("description")
	argRequires    = starlark.String("requires")
	argMet         = starlark.String("met")
	argMeet        = starlark.String("meet")
)

// Dep represents the `dep()` builtin function and models a dependency in the
// graph.
//
// The structure of `dep` is as follows:
//
//   dep(
//
//     name = 'foo-linux',
//
//     description = 'A meaningful description of the dependency',
//
//     requirements = [
//       // a list of other dep variables, either in the current module,
//       // or contained within another module
//       bar, baz,
//     ],
//
//     met = [
//       // a list of actions that must all be satisfied for the dep to be
//       // satisfied
//       shell("true"),
//     ],
//
//     meet = [
//       // a list of actions that will be run to attempt to satisfy the
//       // current dependency
//       shell("true"),
//     ],
//   )
//
type Dep struct {

	// Name is the name of the dependency. Name should be unique across all
	// modules.
	Name string

	// Description is a meaningful, human readable description of the
	// dependency.
	Description string

	// Requirements is a list of the names of other dependencies this dep
	// depends on.
	Requirements []*Dep

	// MetCommand is a list of ShellCmds that should be run, in order, to
	// determine whether this dependency is satisfied. These commands should be
	// lightweight and ideally do not have side-effects. For example, these
	// commands could check for the presence of a binary or directory.
	MetCommands []*ShellCmd

	// MeetCommands is a list of ShellCmds that should be run, in order, to
	// attempt to satisfy this dependency. These commands typically have
	// side-effects and  will install a particular dependency, clone a repo,
	// or create a directory.
	MeetCommands []*ShellCmd
}

// String returns the string representation of the Dep.
func (d Dep) String() string {
	return fmt.Sprintf("<dep.Dep %q>", d.Name)
}

// Type returns a short description about Dep's type.
func (d Dep) Type() string { return "dep.Dep" }

// Freeze does nothing for a Dep.
func (d Dep) Freeze() {}

// Truth always returns true for a Dep.
func (d Dep) Truth() starlark.Bool { return starlark.True }

// Hash is currently not implemented by Dep.
func (d Dep) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: %s", d.Type())
}

// FnDep implements the signature for a builtin function and implements
// the functionality of the `dep` function.
//
// FnDep transforms the arguments into a Dep object after performing
// validation on the keyword arguments.
func FnDep(t *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	dep := &Dep{}

	for _, tuple := range kwargs {
		key := tuple.Index(0)
		value := tuple.Index(1)

		switch key {

		case argName:
			name, err := asString(value)
			if err != nil {
				return nil, err
			}
			dep.Name = name

		case argDescription:
			description, err := asString(value)
			if err != nil {
				return nil, err
			}
			dep.Description = description

		case argRequires:
			requirements, err := asRequirementsList(value)
			if err != nil {
				return nil, err
			}
			dep.Requirements = requirements

		case argMet:
			cmds, err := asCommandList(value)
			if err != nil {
				return nil, err
			}
			dep.MetCommands = cmds

		case argMeet:
			cmds, err := asCommandList(value)
			if err != nil {
				return nil, err
			}
			dep.MeetCommands = cmds

		default:
			continue
		}
	}

	// TODO(nickt): Validation of arguments

	return dep, nil
}

// asString returns the given starlark.Value as a Go string.
// An error is returned if the value is not a string.
func asString(value starlark.Value) (string, error) {
	str, ok := value.(starlark.String)
	if !ok {
		return "", fmt.Errorf("value %v is not a string", value)
	}
	return string(str), nil
}

// asRequirementsList returns the given list value as a slice of string.
// An error is returned if the value is not a list of strings.
func asRequirementsList(value starlark.Value) ([]*Dep, error) {
	list, ok := value.(*starlark.List)
	if !ok {
		return nil, fmt.Errorf("value %v is not a list", value)
	}

	var deps []*Dep

	iter := list.Iterate()
	defer iter.Done()

	var v starlark.Value
	for iter.Next(&v) {
		dep, ok := v.(*Dep)
		if !ok {
			return nil, fmt.Errorf("value %v is not a dep", v)
		}
		deps = append(deps, dep)
	}

	return deps, nil
}

// asCommandList returns the given list value a slice of ShellCmds.
// An error is returned if the value is not a list of ShellCmds.
func asCommandList(value starlark.Value) ([]*ShellCmd, error) {
	list, ok := value.(*starlark.List)
	if !ok {
		return nil, fmt.Errorf("value %+v is not a list", value)
	}

	var commands []*ShellCmd

	iter := list.Iterate()
	defer iter.Done()

	var v starlark.Value
	for iter.Next(&v) {
		command, ok := v.(ShellCmd)
		if !ok {
			return nil, fmt.Errorf("list item %+v is not a command", v)
		}
		commands = append(commands, &command)
	}

	return commands, nil
}
